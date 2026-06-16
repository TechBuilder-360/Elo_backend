package business

import (
	"context"

	"github.com/Toflex/directory_v2/ent"
)

type IRepository interface {
	WithTransaction(db *ent.Client) IRepository
	GetBusinessByName(ctx context.Context, name string) (*businessResult, error)
	GetBusinessByID(ctx context.Context, id string) (*businessResult, error)
	Create(ctx context.Context, payload interface{}) error
}

type repository struct {
	db *ent.Client
}

func Newrepository(db *ent.Client) IRepository {
	return &repository{
		db: db,
	}
}
