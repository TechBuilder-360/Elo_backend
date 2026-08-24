package main

import (
	"github.com/Toflex/directory_v2/cmd/http/runtime"
	"github.com/Toflex/directory_v2/cmd/migration/atlas"
	"github.com/Toflex/directory_v2/cmd/migration/seed"
	"github.com/Toflex/directory_v2/database/database"
	"github.com/Toflex/directory_v2/ent"
	"github.com/Toflex/directory_v2/pkg/configuration"
	"github.com/Toflex/directory_v2/pkg/log"
	"github.com/samber/do/v2"
)

var injector = do.New()

func initializeDI() {
	// provides Database connection
	// database initialization
	do.Provide(injector, database.NewClient)
}

func main() {
	configuration.LoadBaseConfiguration()

	// initialize Runtime Dependency
	initializeDI()

	// close database
	db := do.MustInvoke[*ent.Client](injector)
	defer db.Close()

	// register providers
	runtime.Register(injector)

	// run ATLAS migration
	atlas.AtlasMigration()

	seed.Seeder(db)

	log.Info("Migration completed")
}
