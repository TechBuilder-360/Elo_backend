package wallet

import (
	"context"

	"github.com/Toflex/directory_v2/ent"
	"github.com/Toflex/directory_v2/graph/model"
	"github.com/Toflex/directory_v2/internal/currency"
	"github.com/Toflex/directory_v2/internal/vault"
	"github.com/Toflex/directory_v2/pkg/types"
	"github.com/Toflex/directory_v2/pkg/util"
	"github.com/samber/do/v2"
)

type IService interface {
	GetWallets(ctx context.Context, ownerID, walletType string) ([]*model.Wallet, error)
	GetWallet(ctx context.Context, ownerID, walletType, currencyCode string) (*model.Wallet, error)
	AddWallet(ctx context.Context, ownerID, walletType, currencyCode string) (*model.Wallet, error)
}

type service struct {
	db              *ent.Client
	repo            IRepository
	vaultService    vault.IService
	currencyService currency.IService
}

func Newservice(i do.Injector) IService {
	db := do.MustInvoke[*ent.Client](i)
	currency := do.MustInvoke[currency.IService](i)
	return &service{
		db:              db,
		repo:            Newrepository(db),
		vaultService:    vault.NewService(db),
		currencyService: currency,
	}
}

func (wallet *WalletResponse) ToWallet() model.Wallet {
	return model.Wallet{
		Currency: wallet.Currency.ToString(),
		AvailableBalance: util.ToMajorUnit(types.ToMajor{
			Amount:    wallet.AvailableBalance,
			Precision: uint(wallet.Multiplier),
		}),
		LedgerBalance: util.ToMajorUnit(types.ToMajor{
			Amount:    wallet.LedgerBalance,
			Precision: uint(wallet.Multiplier),
		}),
		HoldingBalance: util.ToMajorUnit(types.ToMajor{
			Amount:    wallet.HoldingBalance,
			Precision: uint(wallet.Multiplier),
		}),
		Active: wallet.Active,
		ID:     wallet.ID,
		IsFiat: wallet.IsFiat,
	}
}
