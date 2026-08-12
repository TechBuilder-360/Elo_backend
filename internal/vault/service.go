package vault

import (
	"context"

	"github.com/Toflex/directory_v2/database/redis"
	"github.com/Toflex/directory_v2/ent"
)

type IService interface {
	CreateVault(ctx context.Context, ownerID string) (*Vault, error)
	GetVault(ctx context.Context, ownerID string) (*Vault, error)
	GetOwnerVault(ctx context.Context, ownerID, vaultType string) (*Vault, error)
}

type Service struct {
	repo  IRepository
	cache *redis.Client
}

func NewService(db *ent.Client) IService {
	return &Service{
		repo: newRepository(db),
	}
}
