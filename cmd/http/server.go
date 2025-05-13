package main

import (
	"net/http"
	"os"

	"github.com/Toflex/directory_v2/cmd/http/router"
	"github.com/Toflex/directory_v2/cmd/http/runtime"
	"github.com/Toflex/directory_v2/database/database"
	"github.com/Toflex/directory_v2/pkg/configuration"
	"github.com/samber/do/v2"
	log "github.com/sirupsen/logrus"
)

type Config struct {
	DbName string `env:"DB_NAME"`
	DbUser string `env:"DB_USER"`
	DbPass string `env:"DB_PASS"`
	DbHost string `env:"DB_HOST"`
	DbPort uint   `env:"DB_PORT"`
}

func initLog() *log.Logger {
	l := log.New()
	l.SetFormatter(&log.JSONFormatter{})
	l.SetOutput(os.Stdout)
	l.SetLevel(log.InfoLevel)

	return l
}

func main() {
	l := initLog()
	configuration.LoadBaseConfiguration()

	// initialize Runtime Dependency
	runtime.InitializeDI()

	// close database
	db := do.MustInvoke[*database.Client](runtime.Injector)
	defer db.Close()

	defer runtime.Injector.Shutdown()

	// initialize routes
	router.InitializeRoutes()

	port := configuration.Instance.Port
	l.Printf("connect to http://localhost:%s/ for GraphQL playground", port)
	l.Fatal(http.ListenAndServe(":"+port, nil))
}
