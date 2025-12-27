package user

import (
	"context"

	"github.com/Toflex/directory_v2/database/database"
)

type IUserRepository interface {
	GetByID(context.Context, string) (*User, error)
	GetByEmail(ctx context.Context, email string) (*User, error)
}

type UserRepository struct {
	db *database.Client
}

func NewUserRepository(db *database.Client) IUserRepository {
	return &UserRepository{
		db: db,
	}
}

func (ur *UserRepository) WithTx(tx *database.Client) IUserRepository {
	return &UserRepository{
		db: tx,
	}
}
