package account

import (
	"github.com/Toflex/directory_v2/ent"
	"github.com/samber/do/v2"
)

type IService interface {
}

type service struct {
	db   *ent.Client
	repo IRepository
}

func Newservice(i do.Injector) IService {
	db := do.MustInvoke[*ent.Client](i)
	return &service{
		db:   db,
		repo: Newrepository(db),
	}
}
