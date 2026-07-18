package wallet

import (
	"context"
	"fmt"

	"github.com/Toflex/directory_v2/ent"
	"github.com/Toflex/directory_v2/ent/currency"
	"github.com/Toflex/directory_v2/ent/wallet"
	"github.com/Toflex/directory_v2/pkg/types"
)

func getIdentifier(id string) string {
	return fmt.Sprintf("identifier:%s", id)
}

func (r *repository) GetWallets(ctx context.Context, identifier string, walletType WalletType) ([]WalletResponse, error) {
	wallets, err := r.db.Wallet.Query().
		Where(wallet.IdentifierEQ(getIdentifier(identifier)),
			wallet.TypeEQ(wallet.Type(walletType))).
		WithCurrency(func(q *ent.CurrencyQuery) {
			q.Select(
				currency.FieldCode,
				currency.FieldMultipler,
			)
		}).
		All(ctx)
	if err != nil {
		return nil, err
	}

	resp := make([]WalletResponse, 0, len(wallets))

	for _, w := range wallets {
		wr := WalletResponse{
			ID:               w.ID,
			Type:             string(w.Type),
			AvailableBalance: w.AvailableBalance,
			Owner:            w.Owner.String(),
			Identifier:       w.Identifier,
			Active:           w.Active,
			Currency:         types.CurrencyCode(w.Edges.Currency.Code),
			Multiplier:       w.Edges.Currency.Multipler,
		}
		resp = append(resp, wr)
	}

	return resp, nil
}

func (r *repository) GetWallet(ctx context.Context, identifier string, walletType WalletType, currencyCode types.CurrencyCode) (*WalletResponse, error) {
	w, err := r.db.Wallet.Query().
		Where(wallet.IdentifierEQ(getIdentifier(identifier)),
			wallet.TypeEQ(wallet.Type(walletType)),
			wallet.HasCurrencyWith(currency.CodeEQ(currencyCode.Capitalize().ToString()))).
		WithCurrency(func(q *ent.CurrencyQuery) {
			q.Select(
				currency.FieldCode,
				currency.FieldMultipler,
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
		Type:             string(w.Type),
		AvailableBalance: w.AvailableBalance,
		Owner:            w.Owner.String(),
		Identifier:       w.Identifier,
		Active:           w.Active,
		Currency:         types.CurrencyCode(w.Edges.Currency.Code),
		Multiplier:       w.Edges.Currency.Multipler,
	}, nil
}

func (r *repository) Create(ctx context.Context, payload *createWallet) (*WalletResponse, error) {
	owner := wallet.OwnerUSER
	if payload.IsBusiness {
		owner = wallet.OwnerBUSINESS
	}

	w, err := r.db.Wallet.Create().
		SetCurrency(payload.Currency).
		SetIdentifier(getIdentifier(payload.Identifier)).
		SetOwner(owner).
		SetType(wallet.Type(payload.Type)).
		Save(ctx)
	if err != nil {
		return nil, err
	}

	return &WalletResponse{
		ID:               w.ID,
		Type:             string(w.Type),
		AvailableBalance: w.AvailableBalance,
		Owner:            w.Owner.String(),
		Identifier:       w.Identifier,
		Active:           w.Active,
		Currency:         types.CurrencyCode(payload.Currency.Code),
		Multiplier:       payload.Currency.Multipler,
	}, nil
}
