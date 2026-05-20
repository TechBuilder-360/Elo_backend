package business

import (
	"context"

	"github.com/Toflex/directory_v2/ent"
	"github.com/Toflex/directory_v2/internal/user"
	"github.com/Toflex/directory_v2/pkg/log"
	"github.com/samber/do/v2"
)

type IService interface {
	CreateBusiness(ctx context.Context, payload CreateBusinessRequest, logger log.Entry) error
}

type service struct {
	db             *ent.Client
	repo           IRepository
	userRepository user.IRepository
}

func NewService(i do.Injector) IService {
	db := do.MustInvoke[*ent.Client](i)
	return &service{
		db:             db,
		repo:           Newrepository(db),
		userRepository: user.NewRepository(db),
	}
}
