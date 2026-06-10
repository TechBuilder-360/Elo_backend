package seed

import (
	"context"

	"entgo.io/ent/dialect/sql"
	"github.com/Toflex/directory_v2/ent"
	"github.com/Toflex/directory_v2/ent/schema"
	"github.com/Toflex/directory_v2/ent/service"
)

func seedServices(ctx context.Context, db *ent.Client) error {
	return db.Service.CreateBulk(
		db.Service.Create().
			SetName("Send Email").
			SetIdentifier("email-provider").
			SetMin(0).
			SetMax(0).
			SetActive(true).
			SetProvider("brevo").SetFee(&schema.Fee{}),
		db.Service.Create().
			SetName("Identity Verification").
			SetIdentifier("identity-verification").
			SetMin(0).
			SetMax(0).
			SetActive(true).
			SetProvider("dojah").SetFee(&schema.Fee{}),
	).OnConflict(
		sql.ConflictColumns(service.FieldIdentifier),
		sql.DoNothing(),
	).
		Exec(ctx)
}
