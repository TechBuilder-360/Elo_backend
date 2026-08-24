package seed

import (
	"context"
	"strings"

	"entgo.io/ent/dialect/sql"
	"github.com/Toflex/directory_v2/ent"
	"github.com/Toflex/directory_v2/ent/provider"
	"github.com/Toflex/directory_v2/pkg/configuration"
	providerImpl "github.com/Toflex/directory_v2/pkg/provider"
	"github.com/Toflex/directory_v2/pkg/util"
)

func seedProviders(ctx context.Context, db *ent.Client) error {
	serviceProviders := providerImpl.GetProviders()
	builders := make([]*ent.ProviderCreate, 0)

	for _, p := range serviceProviders {
		if configuration.IsProduction() && p.Slug() == "mock" {
			continue
		}

		builders = append(builders, db.Provider.Create().SetName(util.ToTitleCase(p.DisplayName())).SetSlug(strings.ToLower(p.Slug())).SetActive(true))
	}

	return db.Provider.CreateBulk(
		builders...,
	).OnConflict(
		sql.ConflictColumns(provider.FieldSlug),
		sql.DoNothing(),
	).
		Exec(ctx)
}
