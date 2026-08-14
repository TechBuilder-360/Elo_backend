package vault

import (
	"context"

	"github.com/Toflex/directory_v2/ent"
)

type IRepository interface {
	Create(ctx context.Context, ownerID string) (*ent.Vault, error)
	GetByID(ctx context.Context, id string) (*ent.Vault, error)
	GetByOwner(ctx context.Context, ownerID, vaultType string) (*ent.Vault, error)
	GetBusinessVaultByType(ctx context.Context, b *ent.Business, vaultType string) (*ent.Vault, error)
	WithTransaction(tx *ent.Tx) IRepository
}

type repository struct {
	db *ent.Client
}

func newRepository(db *ent.Client) IRepository {
	return &repository{
		db: db,
	}
}
