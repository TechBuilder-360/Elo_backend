package wallet

import (
	"context"

	"github.com/Toflex/directory_v2/ent"
	"github.com/Toflex/directory_v2/pkg/types"
)

type IRepository interface {
	GetWallets(ctx context.Context, ownerID string, walletType WalletType) ([]WalletResponse, error)
	GetWallet(ctx context.Context, walletID, ownerID string) (*WalletResponse, error)
	GetWalletWithCurrency(ctx context.Context, ownerID, walletType string, currencyCode types.CurrencyCode) (*WalletResponse, error)
	Create(ctx context.Context, payload *createWallet) (*WalletResponse, error)
}

type repository struct {
	db *ent.Client
}

func Newrepository(db *ent.Client) IRepository {
	return &repository{
		db: db,
	}
}
