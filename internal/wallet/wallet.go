package wallet

import (
	"context"

	"github.com/Toflex/directory_v2/ent"
	"github.com/Toflex/directory_v2/ent/ledgerowner"
	"github.com/Toflex/directory_v2/ent/vault"
	"github.com/Toflex/directory_v2/graph/model"
	"github.com/Toflex/directory_v2/pkg/errors"
	"github.com/Toflex/directory_v2/pkg/log"
	"github.com/Toflex/directory_v2/pkg/types"
)

// GetWallets implements [IService].
func (s *service) GetWallets(ctx context.Context, ownerID, walletType string, filter *model.WalletFilter) ([]*model.Wallet, error) {
	logger := log.LoggerInContext(ctx)
	result, err := s.repo.GetWallets(ctx, ownerID, GetWalletType(walletType), &WalletFilter{
		IsFiat: &filter.IsFiat,
	})
	if err != nil {
		logger.WithError(err).Error("failed to fetch wallets")
		return nil, errors.New(errors.ErrFailed, "request failed")
	}

	wallets := make([]*model.Wallet, 0)

	for _, wallet := range result {
		result := wallet.ToWallet()
		wallets = append(wallets, &result)
	}

	return wallets, nil
}

// GetWallet implements [IService].
func (s *service) GetWallet(ctx context.Context, ownerID, walletType, currencyCode string) (*model.Wallet, error) {
	logger := log.LoggerInContext(ctx)
	wallet, err := s.repo.GetWalletWithCurrency(ctx, ownerID, walletType, types.CurrencyCode(currencyCode))
	if err != nil {
		logger.WithError(err).Error("failed to fetch wallet")
		return nil, errors.New(errors.ErrFailed, "request failed")
	}

	if wallet == nil {
		return nil, errors.New(errors.ErrFailed, "wallet not found")
	}

	result := wallet.ToWallet()

	return &result, nil
}

func (s *service) AddWallet(ctx context.Context, ownerID, walletType, currencyCode string) (*model.Wallet, error) {
	logger := log.LoggerInContext(ctx)

	if !ValidateWalletType(walletType) {
		logger.WithField("wallet_type", walletType).Error("unable to validate wallet type")
		return nil, errors.New(errors.ErrInvalidInput, "invalid wallet type")
	}

	currency, err := s.currencyService.GetCurrencyByCode(ctx, types.CurrencyCode(currencyCode))
	if err != nil {
		logger.WithError(err).Error("unable to fetch currencies")
		return nil, errors.New(errors.ErrFailed, "request failed")
	}

	// check if wallet already exists for the owner with the given currency code
	wallet, err := s.repo.GetWalletWithCurrency(ctx, ownerID, walletType, types.CurrencyCode(currencyCode))
	if err != nil {
		logger.WithError(err).Error("wallet could not be retrived")
		return nil, errors.New(errors.ErrFailed, string(errors.ErrFailed))
	}

	var result model.Wallet

	if wallet != nil {
		result = wallet.ToWallet()
		return &result, nil
	}

	vault, err := s.vaultService.GetOwnerVault(ctx, ownerID, walletType)
	if err != nil {
		logger.WithError(err).Error("unable to fetch vault")
		return nil, errors.New(errors.ErrFailed, "request failed")
	}

	w := &createWallet{
		VaultID:  vault.ID,
		Currency: currency,
	}

	wallet, err = s.repo.Create(ctx, w)
	if err != nil {
		logger.WithError(err).Error("unable to create wallet")
		return nil, errors.New(errors.ErrFailed, "request failed")
	}

	result = wallet.ToWallet()
	return &result, nil
}

func (s *service) GetWalletOwnerWithID(ctx context.Context, ownerID, walletID string) (*model.Wallet, error) {
	wallet, err := s.repo.GetWalletByOwnerID(ctx, walletID, ownerID)
	if err != nil {
		log.LoggerInContext(ctx).WithError(err).Error("Wallet not found")
		return nil, errors.New(errors.ErrNotFound, "wallet not found")
	}

	if wallet == nil {
		return nil, nil
	}

	result := wallet.ToWallet()
	return &result, nil
}

func (s *service) ValidateWalletForDebit(ctx context.Context, payload *ValidateWalletPayload) (*WalletResponse, error) {
	logger := log.LoggerInContext(ctx)

	repo := s.repo
	if payload.DBTransaction != nil {
		repo = s.repo.WithTransaction(payload.DBTransaction)
	}

	walletType := GetWalletType(payload.WalletType)

	overview, err := repo.GetOwnerWalletOverview(ctx, payload.OwnerID, walletType)
	if err != nil {
		logger.WithError(err).Error("failed to get wallet overview for pre-debit action")
		return nil, errors.New(errors.ErrFailed, "request failed")
	}

	if overview == nil {
		logger.WithField("payload", payload).Error("ledger owner not found")
		return nil, errors.New(errors.ErrNotFound, "ledger owner not found")
	}

	// 2. Ledger owner is active
	if overview.LedgerOwnerStatus != ledgerowner.StatusACTIVE.String() {
		logger.WithField("owner_status", overview.LedgerOwnerStatus).Error("wallet owner is not active")
		return nil, errors.New(errors.ErrFailed, "ledger owner is not active")
	}

	// 3. Vault is active
	if overview.VaultStatus != vault.StatusACTIVE.String() {
		logger.WithField("vault_status", overview.VaultStatus).Error("vault is not active")
		return nil, errors.New(errors.ErrFailed, "vault is not active")
	}

	// 1. Wallet is active
	w, err := repo.GetWalletWithAtomicLock(ctx, payload.DBTransaction, payload.OwnerID, string(walletType), payload.Currency)
	if err != nil {
		return nil, err
	}

	if w == nil {
		return nil, errors.New(errors.ErrNotFound, "wallet not found")
	}

	if !w.Active {
		logger.Error("wallet is not active")
		return nil, errors.New(errors.ErrFailed, "wallet is not active")
	}

	return w, nil
}

func (s *service) ValidateWalletForCredit(ctx context.Context, payload *ValidateWalletPayload) (*WalletResponse, error) {
	logger := log.LoggerInContext(ctx)

	repo := s.repo
	if payload.DBTransaction != nil {
		repo = s.repo.WithTransaction(payload.DBTransaction)
	}

	walletType := GetWalletType(payload.WalletType)

	overview, err := repo.GetOwnerWalletOverview(ctx, payload.OwnerID, walletType)
	if err != nil {
		logger.WithError(err).Error("failed to get wallet overview for pre-credit action")
		return nil, errors.New(errors.ErrFailed, "request failed")
	}

	if overview == nil {
		logger.WithField("payload", payload).Error("ledger owner not found")
		return nil, errors.New(errors.ErrNotFound, "ledger owner not found")
	}

	// Check conditions:
	// 1. Wallet is active
	// 2. Ledger owner is active
	// 3. Vault is active
	ownerActive := overview.LedgerOwnerStatus == ledgerowner.StatusACTIVE.String()
	vaultActive := overview.VaultStatus == vault.StatusACTIVE.String()

	walletActive := false
	var targetWallet *WalletResponse

	if ownerActive && vaultActive {
		w, err := repo.GetWalletWithAtomicLock(ctx, payload.DBTransaction, payload.OwnerID, string(walletType), payload.Currency)
		if err == nil && w != nil && w.Active {
			walletActive = true
			targetWallet = w
		}
	}

	// If all conditions are met, return target wallet
	if walletActive && ownerActive && vaultActive {
		return targetWallet, nil
	}

	// If any of the conditions is not met, credit the user's Vault with type 'FROZEN'
	logger.WithFields(map[string]interface{}{
		"wallet_active": walletActive,
		"owner_active":  ownerActive,
		"vault_active":  vaultActive,
		"owner_id":      payload.OwnerID,
	}).Warning("credit condition not met; routing credit to user's FROZEN vault")

	frozenWallet, err := s.getOrCreateFrozenWallet(ctx, payload.DBTransaction, payload.OwnerID, payload.Currency)
	if err != nil {
		logger.WithError(err).Error("failed to get or create frozen wallet for credit")
		return nil, err
	}

	return frozenWallet, nil
}

func (s *service) getOrCreateFrozenWallet(ctx context.Context, tx *ent.Tx, ownerID string, currencyCode types.CurrencyCode) (*WalletResponse, error) {
	logger := log.LoggerInContext(ctx)

	repo := s.repo
	if tx != nil {
		repo = s.repo.WithTransaction(tx)
	}

	// 1. Check if wallet already exists in FROZEN vault
	var w *WalletResponse
	var err error
	if tx != nil {
		w, err = repo.GetWalletWithAtomicLock(ctx, tx, ownerID, string(FrozenWalletType), currencyCode)
	} else {
		w, err = repo.GetWalletWithCurrency(ctx, ownerID, string(FrozenWalletType), currencyCode)
	}
	if err != nil {
		logger.WithError(err).Error("failed to query frozen wallet")
		return nil, errors.New(errors.ErrFailed, "request failed")
	}
	if w != nil {
		return w, nil
	}

	// 2. Ensure FROZEN vault exists
	var vaultID string
	overview, err := repo.GetOwnerWalletOverview(ctx, ownerID, FrozenWalletType)
	if err != nil {
		logger.WithError(err).Error("failed to get frozen vault overview")
		return nil, errors.New(errors.ErrFailed, "request failed")
	}

	if overview != nil && overview.VaultID != "" {
		vaultID = overview.VaultID
	} else {
		createdVault, err := repo.CreateVault(ctx, ownerID, FrozenWalletType)
		if err != nil {
			logger.WithError(err).Error("failed to create frozen vault")
			return nil, errors.New(errors.ErrFailed, "request failed")
		}
		vaultID = createdVault.ID
	}

	// 3. Get currency details
	curr, err := s.currencyService.GetCurrencyByCode(ctx, currencyCode)
	if err != nil {
		logger.WithError(err).Error("unable to fetch currency for frozen wallet")
		return nil, errors.New(errors.ErrFailed, "request failed")
	}

	// 4. Create wallet in FROZEN vault
	newWallet, err := repo.Create(ctx, &createWallet{
		VaultID:  vaultID,
		Currency: curr,
	})
	if err != nil {
		logger.WithError(err).Error("unable to create wallet in frozen vault")
		return nil, errors.New(errors.ErrFailed, "request failed")
	}

	if tx != nil {
		lockedW, err := repo.GetWalletWithAtomicLock(ctx, tx, ownerID, string(FrozenWalletType), currencyCode)
		if err == nil && lockedW != nil {
			return lockedW, nil
		}
	}

	return newWallet, nil
}
