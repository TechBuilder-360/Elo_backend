package currency

import (
	"context"

	"github.com/Toflex/directory_v2/ent"
	"github.com/Toflex/directory_v2/ent/currency"
	"github.com/Toflex/directory_v2/pkg/types"
)

// GetCurrencyByCode implements [IRepository].
func (r *repository) GetCurrencyByCode(ctx context.Context, currencyCode types.CurrencyCode) (*ent.Currency, error) {
	return r.db.Currency.Query().
		Where(currency.CodeEQ(currencyCode.Capitalize().ToString())).
		First(ctx)
}

// GetCurrencyByCode implements [IRepository].
func (r *repository) Currencies(ctx context.Context, filter *CurrencyFilter) ([]*ent.Currency, error) {
	query := r.db.Currency.Query().Where(currency.ActiveEQ(true))

	if filter != nil {
		query = query.Where(currency.IsFiatEQ(filter.IsFiat))
	}

	return query.All(ctx)
}
