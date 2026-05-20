package user

import (
	"context"

	"github.com/Toflex/directory_v2/ent"
)

type IRepository interface {
	GetByID(context.Context, string) (*User, error)
	GetByEmail(ctx context.Context, email string) (*User, error)
}

type repository struct {
	db *ent.Client
}

func NewRepository(db *ent.Client) IRepository {
	return &repository{
		db: db,
	}
}
