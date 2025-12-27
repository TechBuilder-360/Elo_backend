package authentication

import (
	"context"

	"github.com/Toflex/directory_v2/database/database"
)

type IRepository interface {
	Create(ctx context.Context, payload Onboarding) (*string, error)
	GetUserByID(ctx context.Context, id string) (*User, error)
	GetUserByEmail(ctx context.Context, email string) (*User, error)
}

type repository struct {
	db *database.Client
}

func Newrepository(db *database.Client) IRepository {
	return &repository{
		db: db,
	}
}
