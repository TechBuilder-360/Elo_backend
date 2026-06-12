package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Toflex/directory_v2/cmd/http/router"
	"github.com/Toflex/directory_v2/cmd/http/runtime"
	r "github.com/Toflex/directory_v2/database/redis"
	"github.com/Toflex/directory_v2/ent"
	"github.com/Toflex/directory_v2/pkg/configuration"
	"github.com/Toflex/directory_v2/pkg/queue"
	qserver "github.com/Toflex/directory_v2/pkg/queue/server"
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

	// register providers
	runtime.Register()

	// close database
	db := do.MustInvoke[*ent.Client](runtime.Injector)
	defer db.Close()

	// close redis database
	rdb := do.MustInvoke[*r.Client](runtime.Injector)
	defer rdb.Close()

	queueClient := do.MustInvoke[*queue.Client](runtime.Injector)
	defer queueClient.Close()

	queueServer := do.MustInvoke[*qserver.Server](runtime.Injector)
	defer queueServer.Shutdown()

	go func() {
		if err := queueServer.Run(); err != nil {
			l.Errorf("asynq worker stopped: %v", err)
		}
	}()

	defer runtime.Injector.Shutdown()

	// http engine
	engine := do.MustInvoke[*gin.Engine](runtime.Injector)
	if configuration.IsProduction() {
		gin.SetMode(gin.ReleaseMode)
	}

	// initialize routes
	router.InitializeRoutes(engine)

	// initialize Sentry
	// sentry.InitializeSentry(runtime.Injector, l)

	port := configuration.Instance.Port
	addr := fmt.Sprintf("%s:%s", configuration.Instance.BASEURL, port)
	l.Infof("connect to %s", addr)

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", port),
		Handler: engine,
	}

	// Start server in goroutine
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	l.Info("Server started")

	// Wait for interrupt signal
	quit := make(chan os.Signal, 1)

	signal.Notify(
		quit,
		syscall.SIGINT,
		syscall.SIGTERM,
	)

	<-quit

	l.Info("Shutting down server...")

	ctx, cancel := context.WithTimeout(
		context.Background(),
		5*time.Second,
	)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown:", err)
	}

	l.Info("Server exiting")
}
