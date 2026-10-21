package wallet

import (
	"context"

	"entgo.io/ent/dialect/sql"
	"github.com/Toflex/directory_v2/ent"
	"github.com/Toflex/directory_v2/ent/currency"
	"github.com/Toflex/directory_v2/ent/ledger"
	"github.com/Toflex/directory_v2/ent/ledgerowner"
	"github.com/Toflex/directory_v2/ent/transaction"
	"github.com/Toflex/directory_v2/ent/vault"
	"github.com/Toflex/directory_v2/ent/wallet"
	"github.com/Toflex/directory_v2/pkg/types"
)

func (r *repository) GetWallets(
	ctx context.Context,
	ownerID string,
	walletType WalletType,
	filter *WalletFilter,
) ([]WalletResponse, error) {

	query := r.db.Vault.
		Query().
		Where(
			vault.HasOwnerWith(
				ledgerowner.IDEQ(ownerID),
			),
			vault.TypeEQ(vault.Type(walletType)),
		)

	query = query.WithWallets(func(q *ent.WalletQuery) {

		// Apply wallet filters
		if filter != nil && filter.IsFiat != nil {
			q.Where(
				wallet.HasCurrencyWith(
					currency.IsFiatEQ(*filter.IsFiat),
				),
			)
		}

		// Load only the currency fields we need
		q.WithCurrency(func(cq *ent.CurrencyQuery) {
			cq.Select(
				currency.FieldID,
				currency.FieldCode,
				currency.FieldMultiplier,
				currency.FieldIsFiat,
			)
		})
	})

	v, err := query.First(ctx)
	if err != nil {
		return nil, err
	}

	resp := make([]WalletResponse, 0, len(v.Edges.Wallets))

	for _, w := range v.Edges.Wallets {
		currency := w.Edges.Currency

		resp = append(resp, WalletResponse{
			ID:               w.ID,
			AvailableBalance: w.AvailableBalance,
			LedgerBalance:    w.LedgerBalance,
			HoldingBalance:   w.HoldingBalance,
			Active:           w.Active,
			Currency:         types.CurrencyCode(currency.Code),
			Multiplier:       currency.Multiplier,
			IsFiat:           currency.IsFiat,
		})
	}

	return resp, nil
}

func (r *repository) GetWallet(ctx context.Context, walletType, ownerID string) (*WalletResponse, error) {
	v, err := r.db.Vault.Query().
		Where(vault.HasOwnerWith(ledgerowner.IDEQ(ownerID)),
			vault.TypeEQ(vault.Type(walletType))).
		WithWallets(func(q *ent.WalletQuery) {
			q.WithCurrency(func(cq *ent.CurrencyQuery) {
				cq.Select(
					currency.FieldID,
					currency.FieldCode,
					currency.FieldMultiplier,
					currency.FieldIsFiat,
				)
			})
		}).
		First(ctx)
	if err != nil {
		return nil, err
	}

	w := v.Edges.Wallets[0]

	return &WalletResponse{
		ID:               w.ID,
		AvailableBalance: w.AvailableBalance,
		LedgerBalance:    w.LedgerBalance,
		HoldingBalance:   w.HoldingBalance,
		Active:           w.Active,
		Currency:         types.CurrencyCode(w.Edges.Currency.Code),
		Multiplier:       w.Edges.Currency.Multiplier,
		IsFiat:           w.Edges.Currency.IsFiat,
	}, nil
}

func (r *repository) Create(ctx context.Context, payload *createWallet) (*WalletResponse, error) {
	walletType := wallet.TypeCRYPTO
	if payload.Currency.IsFiat {
		walletType = wallet.TypeFIAT
	}

	w, err := r.db.Wallet.Create().
		SetCurrencyID(payload.Currency.ID).
		SetType(walletType).
		SetVaultID(payload.VaultID).
		Save(ctx)
	if err != nil {
		return nil, err
	}

	return &WalletResponse{
		ID:               w.ID,
		AvailableBalance: w.AvailableBalance,
		LedgerBalance:    w.LedgerBalance,
		HoldingBalance:   w.HoldingBalance,
		Active:           w.Active,
		Currency:         types.CurrencyCode(payload.Currency.Code),
		Multiplier:       payload.Currency.Multiplier,
		IsFiat:           payload.Currency.IsFiat,
	}, nil
}

// GetWalletWithCurrency implements [IRepository].
func (r *repository) GetWalletWithCurrency(ctx context.Context, ownerID string, walletType string, currencyCode types.CurrencyCode) (*WalletResponse, error) {
	w, err := r.db.Wallet.Query().
		Where(
			wallet.HasCurrencyWith(currency.CodeEQ(string(currencyCode.Capitalize()))),
			wallet.HasVaultWith(vault.HasOwnerWith(ledgerowner.IDEQ(ownerID)),
				vault.TypeEQ(vault.Type(walletType))),
		).
		WithCurrency(func(q *ent.CurrencyQuery) {
			q.Select(
				currency.FieldID,
				currency.FieldCode,
				currency.FieldMultiplier,
				currency.FieldIsFiat,
			)
		}).First(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, nil
		}
		return nil, err
	}

	return &WalletResponse{
		ID:               w.ID,
		AvailableBalance: w.AvailableBalance,
		LedgerBalance:    w.LedgerBalance,
		HoldingBalance:   w.HoldingBalance,
		Active:           w.Active,
		Currency:         types.CurrencyCode(w.Edges.Currency.Code),
		Multiplier:       w.Edges.Currency.Multiplier,
		IsFiat:           w.Edges.Currency.IsFiat,
	}, nil
}

func (r *repository) GetWalletByOwnerID(ctx context.Context, walletID, ownerID string) (*WalletResponse, error) {
	w, err := r.db.Wallet.Query().Where(
		wallet.IDEQ(walletID), wallet.HasVaultWith(vault.HasOwnerWith(ledgerowner.IDEQ(ownerID))),
	).WithCurrency(func(q *ent.CurrencyQuery) {
		q.Select(
			currency.FieldID,
			currency.FieldCode,
			currency.FieldMultiplier,
			currency.FieldIsFiat,
		)
	}).First(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, nil
		}

		return nil, err
	}

	return &WalletResponse{
		ID:               w.ID,
		AvailableBalance: w.AvailableBalance,
		LedgerBalance:    w.LedgerBalance,
		HoldingBalance:   w.HoldingBalance,
		Active:           w.Active,
		Currency:         types.CurrencyCode(w.Edges.Currency.Code),
		Multiplier:       w.Edges.Currency.Multiplier,
		IsFiat:           w.Edges.Currency.IsFiat,
	}, nil
}

// GetOwnerWalletOverview returns the ledger owner status, vault status and wallet details for the given owner and wallet type.
func (r *repository) GetOwnerWalletOverview(ctx context.Context, ownerID string, walletType WalletType) (*OwnerVaultOverview, error) {
	owner, err := r.db.LedgerOwner.Query().
		Where(ledgerowner.IDEQ(ownerID)).
		WithVaults(func(vq *ent.VaultQuery) {
			vq.Where(vault.TypeEQ(vault.Type(walletType)))
		}).
		First(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, nil
		}
		return nil, err
	}

	overview := &OwnerVaultOverview{
		LedgerOwnerStatus: string(owner.Status),
	}

	if len(owner.Edges.Vaults) > 0 {
		v := owner.Edges.Vaults[0]
		overview.VaultStatus = string(v.Status)
		overview.VaultID = v.ID
	}

	return overview, nil
}

func (r *repository) CreateVault(ctx context.Context, ownerID string, vaultType WalletType) (*ent.Vault, error) {
	vt := vault.Type(vaultType)
	if vaultType == "" {
		vt = vault.DefaultType
	}

	return r.db.Vault.Create().
		SetType(vt).
		SetOwnerID(ownerID).
		Save(ctx)
}

func (r *repository) GetWalletWithAtomicLock(ctx context.Context, tx *ent.Tx, ownerID string, walletType string, currencyCode types.CurrencyCode) (*WalletResponse, error) {
	client := r.db
	if tx != nil {
		client = tx.Client()
	}

	q := client.Wallet.Query().
		Where(
			wallet.HasCurrencyWith(currency.CodeEQ(string(currencyCode.Capitalize()))),
			wallet.HasVaultWith(vault.HasOwnerWith(ledgerowner.IDEQ(ownerID)),
				vault.TypeEQ(vault.Type(walletType))),
		).
		WithCurrency(func(q *ent.CurrencyQuery) {
			q.Select(
				currency.FieldID,
				currency.FieldCode,
				currency.FieldMultiplier,
				currency.FieldIsFiat,
			)
		})

	if tx != nil {
		q = q.ForShare(
			sql.WithLockTables(wallet.Table),
			sql.WithLockAction(sql.NoWait),
		)
	}

	w, err := q.First(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, nil
		}
		return nil, err
	}

	return &WalletResponse{
		ID:               w.ID,
		AvailableBalance: w.AvailableBalance,
		LedgerBalance:    w.LedgerBalance,
		HoldingBalance:   w.HoldingBalance,
		Active:           w.Active,
		Currency:         types.CurrencyCode(w.Edges.Currency.Code),
		Multiplier:       w.Edges.Currency.Multiplier,
		IsFiat:           w.Edges.Currency.IsFiat,
	}, nil
}

func (r *repository) WithTransaction(tx *ent.Tx) IRepository {
	return &repository{
		db: tx.Client(),
	}
}

func (r *repository) UpdateWalletBalance(ctx context.Context, w *WalletResponse, payload *LedgerEntryResult) error {
	_, err := r.db.Wallet.UpdateOneID(w.ID).
		SetAvailableBalance(payload.AvailableBalance).
		SetHoldingBalance(payload.HoldingBalance).
		SetLedgerBalance(payload.LedgerBalance).
		Save(ctx)

	return err
}

func (r *repository) GetLedgerbyTransactionID(ctx context.Context, transactionID string) (*ent.Transaction, error) {
	t, err := r.db.Transaction.
		Query().
		Where(
			transaction.HasLedgerEntryWith(ledger.TransactionIDEQ(transactionID))).
		WithLedgerEntry(func(lq *ent.LedgerQuery) {
			lq.WithWallet(func(wq *ent.WalletQuery) {
				wq.WithVault(func(vq *ent.VaultQuery) {
					vq.WithOwner()
				})
			})
		}).
		Only(ctx)

	if ent.IsNotFound(err) {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	return t, nil
}
