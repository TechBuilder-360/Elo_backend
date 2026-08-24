package nuban

import (
	"context"

	"github.com/Toflex/directory_v2/ent"
	"github.com/samber/do/v2"
)

type IService interface {
	GenerateBusinessStaticNuban(ctx context.Context, ownerID string) (*StaticNubanResponse, error)
	GenerateBusinessDynamicNuban(ctx context.Context, ownerID string) (*DynamicNubanResponse, error)
}

type service struct {
	db   *ent.Client
	repo IRepository
}

func NewService(i do.Injector) IService {
	db := do.MustInvoke[*ent.Client](i)
	return &service{
		db:   db,
		repo: NewRepository(db),
	}
}
