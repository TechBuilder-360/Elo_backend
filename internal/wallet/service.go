package wallet

import (
	"context"

	"github.com/Toflex/directory_v2/ent"
	"github.com/Toflex/directory_v2/graph/model"
	"github.com/Toflex/directory_v2/internal/currency"
	"github.com/samber/do/v2"
)

type IService interface {
	GetWallets(ctx context.Context, b *ent.Business, walletType string) ([]*model.Wallet, error)
	GetWallet(ctx context.Context, b *ent.Business, walletType, currencyCode string) (*model.Wallet, error)
	AddWallet(ctx context.Context, b *ent.Business, walletType, currencyCode string) (*model.Wallet, error)
}

type service struct {
	db           *ent.Client
	repo         IRepository
	currencyRepo currency.IRepository
}

func Newservice(i do.Injector) IService {
	db := do.MustInvoke[*ent.Client](i)
	return &service{
		db:           db,
		repo:         Newrepository(db),
		currencyRepo: currency.Newrepository(db),
	}
}
