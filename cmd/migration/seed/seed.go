package seed

import (
	"context"
	"time"

	"github.com/Toflex/directory_v2/ent"
	"github.com/Toflex/directory_v2/pkg/log"
)

func Seeder(db *ent.Client) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute*5)
	defer cancel()

	errs := make([]error, 0)

	errs = append(errs, seedServices(ctx, db))
	errs = append(errs, seedProviders(ctx, db))
	errs = append(errs, seedRoles(ctx, db))
	errs = append(errs, seedPermissions(ctx, db))
	errs = append(errs, seedRolePermissions(ctx, db))
	errs = append(errs, seedKYBDocuments(ctx, db))

	for _, e := range errs {
		if e != nil {
			log.Error("%s", e.Error())
		}
	}
}
