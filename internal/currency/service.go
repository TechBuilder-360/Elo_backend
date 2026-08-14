package currency

import (
	"context"

	"github.com/Toflex/directory_v2/ent"
	"github.com/Toflex/directory_v2/graph/model"
	"github.com/Toflex/directory_v2/pkg/types"
	"github.com/samber/do/v2"
)

type IService interface {
	GetCurrencies(ctx context.Context, filter *model.CurrencyFilter) ([]*model.Currency, error)
	GetCurrencyByCode(ctx context.Context, code types.CurrencyCode) (*Currency, error)
}

type service struct {
	db   *ent.Client
	repo IRepository
}

func Newservice(i do.Injector) IService {
	db := do.MustInvoke[*ent.Client](i)
	return &service{
		db:   db,
		repo: newrepository(db),
	}
}
