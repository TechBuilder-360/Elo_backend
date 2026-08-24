package crypto

import (
	"context"

	"github.com/Toflex/directory_v2/ent"
	"github.com/Toflex/directory_v2/graph/model"
)

type IRepository interface {
	GetSupportedNetwork(ctx context.Context, networkID, coinID string) (*SupportedNetwork, error)
	GetNetworksByCoinID(ctx context.Context, coinID string) ([]*SupportedNetwork, error)
	GetCoin(ctx context.Context, coin string) (*CryptoCoin, error)
	GetNetwork(ctx context.Context, network string) (*CryptoNetwork, error)
	CreateWallet(ctx context.Context, payload StablecoinWallet) (*StablecoinWalletResult, error)
	GetWalletByCoinAndNetwork(ctx context.Context, walletID, coinID, networkID string) (*StablecoinWalletResult, error)
	GetAddressByWallet(ctx context.Context, ownerID string, filter *model.StablecoinFilter) ([]StablecoinWalletResult, error)
}

type repository struct {
	db *ent.Client
}

func NewRepository(db *ent.Client) IRepository {
	return &repository{
		db: db,
	}
}
