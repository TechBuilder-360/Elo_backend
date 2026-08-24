package crypto

import (
	"context"

	"github.com/Toflex/directory_v2/ent"
	"github.com/Toflex/directory_v2/graph/model"
	"github.com/Toflex/directory_v2/internal/wallet"
	"github.com/samber/do/v2"
)

type IService interface {
	GenerateCryptoAddress(ctx context.Context, b *ent.Business, body *StablecoinRequest) (*CryptoAddressResponse, error)
	GetCryptoAddress(ctx context.Context, ownerID string, filter *model.StablecoinFilter) ([]CryptoAddressResponse, error)
}

type service struct {
	db            *ent.Client
	repo          IRepository
	walletService wallet.IService
}

func NewService(i do.Injector) IService {
	db := do.MustInvoke[*ent.Client](i)
	return &service{
		db:            db,
		repo:          NewRepository(db),
		walletService: do.MustInvoke[wallet.IService](i),
	}
}
