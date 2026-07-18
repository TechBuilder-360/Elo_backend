package currency

import (
	"context"

	"github.com/Toflex/directory_v2/ent"
	"github.com/Toflex/directory_v2/pkg/types"
)

type IRepository interface {
	GetCurrencyByCode(ctx context.Context, currencyCode types.CurrencyCode) (*ent.Currency, error)
}

type repository struct {
	db *ent.Client
}

func Newrepository(db *ent.Client) IRepository {
	return &repository{
		db: db,
	}
}
