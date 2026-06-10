package verification

import (
	"context"

	"github.com/Toflex/directory_v2/ent"
	"github.com/Toflex/directory_v2/internal/user"
	"github.com/samber/do/v2"
)

type IService interface {
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
		repo:           NewRepository(db),
		userRepository: user.NewRepository(db),
	}
}

func (s *service) HandleUserVerification(ctx context.Context) error {
	return nil
}
