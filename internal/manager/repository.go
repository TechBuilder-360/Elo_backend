package manager

import (
	"context"

	"github.com/Toflex/directory_v2/ent"
)

type IRepository interface {
	AddManager(ctx context.Context, manager *Manager) error
	WithTransaction(tx *ent.Tx) IRepository
}

type repository struct {
	db *ent.Client
}

func NewRepository(db *ent.Client) IRepository {
	return &repository{db: db}
}
