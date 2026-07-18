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
