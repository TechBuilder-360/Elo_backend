package atlas

import (
	"context"
	"fmt"
	"log"
	"os"

	"ariga.io/atlas-go-sdk/atlasexec"
	"github.com/Toflex/directory_v2/pkg/configuration"
)

type config struct {
	DbName    string `env:"DB_NAME"`
	DbUser    string `env:"DB_USER"`
	DbPass    string `env:"DB_PASS"`
	DbHost    string `env:"DB_HOST"`
	DbPort    uint   `env:"DB_PORT"`
	DbSSLMode string `env:"DB_SSL_MODE"`
	DBSchema  string `env:"DB_SCHEMA"`
}

func AtlasMigration() {
	conf := &config{}
	conf = configuration.Load(conf).(*config)

	// Define the execution context, supplying a migration directory
	// and potentially an `atlas.hcl` configuration file using `atlasexec.WithHCL`.
	workdir, err := atlasexec.NewWorkingDir(
		atlasexec.WithMigrations(
			os.DirFS("ent/migrate/migrations"),
		),
	)
	if err != nil {
		log.Fatalf("failed to load working directory: %v", err)
	}
	// atlasexec works on a temporary directory, so we need to close it
	defer workdir.Close()

	// Initialize the client.
	client, err := atlasexec.NewClient(workdir.Path(), "atlas")
	if err != nil {
		log.Fatalf("failed to initialize client: %v", err)
	}

	// Run `atlas migrate apply` on DB
	uri := fmt.Sprintf("postgresql://%s:%s@%s:%d/%s?sslmode=%s",
		conf.DbUser, conf.DbPass, conf.DbHost, conf.DbPort, conf.DbName, conf.DbSSLMode)
	migrationsDir := workdir.Path("migrations")
	dirURL := fmt.Sprintf("file://%s", migrationsDir)
	log.Printf("Atlas working dir=%q, migrations dir=%q", workdir.Path(), migrationsDir)

	status, statusErr := client.MigrateStatus(context.Background(), &atlasexec.MigrateStatusParams{
		URL:             uri,
		DirURL:          dirURL,
		RevisionsSchema: conf.DBSchema,
	})
	if statusErr != nil {
		log.Printf("Atlas migration status check failed: %v", statusErr)
	} else {
		log.Printf("Atlas migration status: %s, current=%q, next=%q, pending=%d", status.Status, status.Current, status.Next, len(status.Pending))
	}

	res, err := client.MigrateApply(context.Background(), &atlasexec.MigrateApplyParams{
		URL:             uri,
		DirURL:          dirURL,
		AllowDirty:      true,
		RevisionsSchema: conf.DBSchema,
	})
	if err != nil {
		log.Fatalf("failed to apply migrations: %v", err)
	}

	if len(res.Applied) == 0 {
		log.Printf("No atlas migrations applied; database is already up to date")
	} else {
		log.Printf("Applied %d migrations\n", len(res.Applied))
	}
}
