package authentication

import (
	"context"
	"github.com/Toflex/directory_v2/database/redis"
	"github.com/Toflex/directory_v2/ent"
	"github.com/Toflex/directory_v2/graph/model"
	"github.com/Toflex/directory_v2/pkg/log"
	"github.com/samber/do/v2"
)

type IService interface {
	RegisterUser(ctx context.Context, body Registration, log log.Entry) (*string, error)
	RequestOTP(ctx context.Context, payload OTPRequest, log log.Entry) (*string, error)
	Login(ctx context.Context, payload Login, log log.Entry) (*model.LoginResponse, error)
	VerifyJWT(ctx context.Context, token string) (string, bool)
	Logout(ctx context.Context, user *ent.User)
}

type Service struct {
	repo  IRepository
	cache *redis.Client
}

func NewService(i do.Injector) IService {
	db := do.MustInvoke[*ent.Client](i)
	return &Service{
		repo:  Newrepository(db),
		cache: do.MustInvoke[*redis.Client](i),
	}
}
