package wallet

import (
	"context"

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
