package provider

import (
	"context"

	"github.com/Toflex/directory_v2/ent"
)

type IRepository interface {
	FindByIdentifier(ctx context.Context, identifier string) (*ent.Service, error)
	GetProvider(ctx context.Context, name string) (*ActiveProvider, error)
	WithTransaction(db *ent.Client) IRepository
}

type repository struct {
	db *ent.Client
}

func Newrepository(db *ent.Client) IRepository {
	return &repository{
		db: db,
	}
}
