package runtime

import (
	"github.com/Toflex/directory_v2/database/database"
	r "github.com/Toflex/directory_v2/database/redis"
	"github.com/Toflex/directory_v2/pkg/queue"
	s "github.com/Toflex/directory_v2/pkg/queue/server"
	"github.com/samber/do/v2"
)

var injector = do.New()

func GetInjector() do.Injector {
	return injector
}

func InitializeDI() {
	// provides Database connection
	// database initialization
	do.Provide(injector, database.NewClient)

	// provides redis connection
	// redis initialization
	do.Provide(injector, r.NewClient)

	// provides async queue client
	do.Provide(injector, queue.NewClient)
	// provides Asynq server for background workers
	do.Provide(injector, s.NewServer)

	// provides Server HTTP Engine
	// GO-GIN server
	do.Provide(injector, server)

	// Initialize Temporal Worker
	initializeTemporal(injector)

	// Initialize Services
	initializeService(injector)
}
