package runtime

import (
	"github.com/Toflex/directory_v2/internal/authentication"
	"github.com/Toflex/directory_v2/internal/business"
	"github.com/samber/do/v2"
)

func initializeService(i do.Injector) {
	do.Provide(i, func(i do.Injector) (authentication.IService, error) {
		return authentication.NewService(i), nil
	})
	do.Provide(i, func(i do.Injector) (business.IService, error) {
		return business.NewService(i), nil
	})
}
