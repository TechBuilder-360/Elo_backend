package wallet

import (
	"context"

	"github.com/Toflex/directory_v2/ent"
	"github.com/Toflex/directory_v2/ent/currency"
	"github.com/Toflex/directory_v2/ent/ledgerowner"
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
