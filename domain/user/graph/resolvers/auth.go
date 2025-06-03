package authentication

import (
	"context"

	"github.com/Toflex/directory_v2/domain/user/graph/model"
)

// Resolver defines the mutation resolver
type Resolver struct{}

func (r *Resolver) Login(ctx context.Context, input model.Login) (*model.LoginResult, error) {
	return nil, nil
}
