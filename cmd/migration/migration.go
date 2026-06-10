package main

import (
	"os"

	"github.com/Toflex/directory_v2/cmd/migration/atlas"
	"github.com/Toflex/directory_v2/cmd/migration/seed"
	"github.com/Toflex/directory_v2/database/database"
	"github.com/Toflex/directory_v2/ent"
	"github.com/samber/do/v2"
)

var injector = do.New()

func initializeDI() {
	// provides Database connection
	// database initialization
	do.Provide(injector, database.NewClient)
}

func main() {
	os.Setenv("ENVIRONMENT", "PRODUCTION")
	// initialize Runtime Dependency
	initializeDI()

	// close database
	db := do.MustInvoke[*ent.Client](injector)
	defer db.Close()

	// register providers
	// runtime.Register()

	// run ATLAS migration
	atlas.AtlasMigration()

	seed.Seeder(db)
}
