package runtime

import (
	"github.com/Toflex/directory_v2/database/database"
	r "github.com/Toflex/directory_v2/database/redis"
	"github.com/Toflex/directory_v2/pkg/queue"
	s "github.com/Toflex/directory_v2/pkg/queue/server"
	"github.com/samber/do/v2"
)

var Injector = do.New()

func InitializeDI() {
	// provides Database connection
	// database initialization
	do.Provide(Injector, database.NewClient)

	// provides redis connection
	// redis initialization
	do.Provide(Injector, r.NewClient)

	// provides async queue client
	do.Provide(Injector, queue.NewClient)
	// provides Asynq server for background workers
	do.Provide(Injector, s.NewServer)

	// provides Server HTTP Engine
	// GO-GIN server
	do.Provide(Injector, server)

	// Initialize Services
	initializeService(Injector)
}
