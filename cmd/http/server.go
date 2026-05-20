package main

import (
	"fmt"
	"os"

	"github.com/Toflex/directory_v2/cmd/http/router"
	"github.com/Toflex/directory_v2/cmd/http/runtime"
	r "github.com/Toflex/directory_v2/database/redis"
	"github.com/Toflex/directory_v2/ent"
	"github.com/Toflex/directory_v2/pkg/apm/sentry"
	"github.com/Toflex/directory_v2/pkg/configuration"
	"github.com/gin-gonic/gin"
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
	db := do.MustInvoke[*ent.Client](runtime.Injector)
	defer db.Close()

	// close redis database
	rdb := do.MustInvoke[*r.Client](runtime.Injector)
	defer rdb.Close()

	defer runtime.Injector.Shutdown()

	// http engine
	engine := do.MustInvoke[*gin.Engine](runtime.Injector)

	// initialize routes
	router.InitializeRoutes(engine)

	// initialize Sentry
	sentry.InitializeSentry(runtime.Injector, l)

	port := configuration.Instance.Port
	addr := fmt.Sprintf("%s:%s", configuration.Instance.BASEURL, port)
	l.Infof("connect to %s", addr)

	engine.Run(addr)
}
