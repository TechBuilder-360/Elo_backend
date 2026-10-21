package seed

import (
	"context"

	"entgo.io/ent/dialect/sql"
	"github.com/Toflex/directory_v2/ent"
)

func SeedNubanBank(ctx context.Context, db *ent.Client) error {
	return db.NubanBank.CreateBulk(
		db.NubanBank.
			Create().
			SetCbnCode("000").
			SetName("ELO Bank").
			SetSlug("elo").
			SetShortName("ELO").
			SetActive(true),
	).OnConflict(
		sql.ConflictColumns("cbn_code"),
		sql.DoNothing(),
	).
		Exec(ctx)
}
