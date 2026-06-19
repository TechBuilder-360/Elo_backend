package user

import (
	"github.com/Toflex/directory_v2/database/redis"
	"github.com/Toflex/directory_v2/ent"
	"github.com/samber/do/v2"
)

type IService interface{}

type Service struct {
	repo  IRepository
	cache *redis.Client
}

func NewService(i do.Injector) IService {
	db := do.MustInvoke[*ent.Client](i)
	return &Service{
		repo:  NewRepository(db),
		cache: do.MustInvoke[*redis.Client](i),
	}
}
