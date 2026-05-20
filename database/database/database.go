package database

import (
	"fmt"

	"github.com/Toflex/directory_v2/ent"
	"github.com/Toflex/directory_v2/pkg/configuration"
	"github.com/Toflex/directory_v2/pkg/log"
	"github.com/samber/do/v2"

	_ "github.com/lib/pq"
)

type config struct {
	DbName string `env:"DB_NAME"`
	DbUser string `env:"DB_USER"`
	DbPass string `env:"DB_PASS"`
	DbHost string `env:"DB_HOST"`
	DbPort uint   `env:"DB_PORT"`
}

// type Client struct {
// 	DBClient *ent.Client
// }

var dbInstance *ent.Client

func initializeDB() {
	conf := &config{}
	conf = configuration.Load(conf).(*config)

	uri := fmt.Sprintf("host=%s port=%d user=%s dbname=%s password=%s sslmode=disable", conf.DbHost, conf.DbPort, conf.DbUser, conf.DbName, conf.DbPass)
	client, err := ent.Open("postgres", uri)
	if err != nil {
		log.Panic("failed opening connection to postgres: %v", err)
	}

	dbInstance = client
}

func NewClient(i do.Injector) (*ent.Client, error) {
	if dbInstance == nil {
		initializeDB()
	}

	return dbInstance, nil
}

func DBInstance() *ent.Client {
	return dbInstance
}

// func (c *Client) MigrateDBSchema() {
// 	// Run the auto migration tool.
// 	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
// 	defer cancel()

// 	if err := c.DBClient.Schema.Create(ctx); err != nil {
// 		log.Panic("failed creating schema resources: %v", err)
// 	}
// }

// func (c *Client) Close() {
// 	err := c.DBClient.Close()
// 	if err != nil {
// 		log.Errorf("Failed to close DB client: %v", err)
// 	}
// }
