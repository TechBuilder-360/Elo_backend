package seed

import (
	"context"

	"entgo.io/ent/dialect/sql"
	"github.com/Toflex/directory_v2/ent"
	"github.com/Toflex/directory_v2/ent/currency"
)

func SeedCurrencies(ctx context.Context, db *ent.Client) error {
	return db.Currency.CreateBulk(
		db.Currency.Create().
			SetCode("NGN").
			SetName("Naira").
			SetSymbol("₦").
			SetMultiplier(100).
			SetIsFiat(true),
		db.Currency.Create().
			SetCode("USD").
			SetName("US Dollar").
			SetSymbol("$").
			SetMultiplier(100).
			SetIsFiat(true),
		db.Currency.Create().
			SetCode("USDC").
			SetName("USD Coin").
			SetSymbol("USDC").
			SetMultiplier(1000000).
			SetIsFiat(false),
		db.Currency.Create().
			SetCode("USDT").
			SetName("Tether USD").
			SetSymbol("USDT").
			SetMultiplier(1000000).
			SetIsFiat(false),
		db.Currency.Create().
			SetCode("PYUSD").
			SetName("PayPal USD").
			SetSymbol("PYUSD").
			SetMultiplier(1000000).
			SetIsFiat(false),
	).OnConflict(
		sql.ConflictColumns(currency.FieldCode),
		sql.DoNothing(),
	).Exec(ctx)
}
