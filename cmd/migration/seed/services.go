package seed

import (
	"context"

	"entgo.io/ent/dialect/sql"
	"github.com/Toflex/directory_v2/ent"
	"github.com/Toflex/directory_v2/ent/service"
)

func seedServices(ctx context.Context, db *ent.Client) error {
	return db.Service.CreateBulk(
		db.Service.Create().
			SetName("Send Email").
			SetIdentifier("email-provider").
			SetMinValue(0).
			SetMaxValue(0).
			SetActive(true).
			SetProvider("brevo"),
		db.Service.Create().
			SetName("Identity Verification").
			SetIdentifier("identity-verification").
			SetMinValue(0).
			SetMaxValue(0).
			SetActive(true).
			SetProvider("dojah"),
		db.Service.Create().
			SetName("Generate Stablecoin Address Solana").
			SetIdentifier("stablecoin-address-sol").
			SetMinValue(0).
			SetMaxValue(0).
			SetActive(true).
			SetProvider("maplerad"),
		db.Service.Create().
			SetName("Generate Stablecoin Address Tron").
			SetIdentifier("stablecoin-address-tron").
			SetMinValue(0).
			SetMaxValue(0).
			SetActive(true).
			SetProvider("maplerad"),
		db.Service.Create().
			SetName("Generate Stablecoin Address Polygon").
			SetIdentifier("stablecoin-address-polygon").
			SetMinValue(0).
			SetMaxValue(0).
			SetActive(true).
			SetProvider("maplerad"),
		db.Service.Create().
			SetName("Generate Stablecoin Address Optimism").
			SetIdentifier("stablecoin-address-optimism").
			SetMinValue(0).
			SetMaxValue(0).
			SetActive(true).
			SetProvider("maplerad"),
		db.Service.Create().
			SetName("Generate Stablecoin Address Ethereum").
			SetIdentifier("stablecoin-address-eth").
			SetMinValue(0).
			SetMaxValue(0).
			SetActive(true).
			SetProvider("maplerad"),
		db.Service.Create().
			SetName("Nigeria Bank Transfer").
			SetIdentifier("ng-bank-transfer").
			SetMinValue(0).
			SetMaxValue(0).
			SetActive(true).
			SetProvider("maplerad"),
	).OnConflict(
		sql.ConflictColumns(service.FieldIdentifier),
		sql.DoNothing(),
	).
		Exec(ctx)
}
