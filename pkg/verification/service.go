package verification

import (
	"github.com/Toflex/directory_v2/database/database"
	"github.com/Toflex/directory_v2/ent"
	"github.com/Toflex/directory_v2/pkg/provider"
)

type Service struct {
	db             *ent.Client
	serviceLocator provider.IService
}

func NewService() *Service {
	return &Service{
		db:             database.DBInstance(),
		serviceLocator: provider.NewService(),
	}
}
