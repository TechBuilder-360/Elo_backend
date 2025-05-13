package runtime

import (
	"github.com/Toflex/directory_v2/database/database"
	r "github.com/Toflex/directory_v2/database/redis"
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
}
