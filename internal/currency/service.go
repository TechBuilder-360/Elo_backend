package currency

import (
	"context"

	"github.com/Toflex/directory_v2/ent"
	"github.com/Toflex/directory_v2/graph/model"
	"github.com/samber/do/v2"
)

type IService interface {
	GetCurrencies(ctx context.Context) ([]*model.Currency, error)
}

type service struct {
	db   *ent.Client
	repo IRepository
}

func Newservice(i do.Injector) IService {
	db := do.MustInvoke[*ent.Client](i)
	return &service{
		db:   db,
		repo: Newrepository(db),
	}
}
