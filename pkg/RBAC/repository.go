package rbac

import (
	"context"

	"github.com/Toflex/directory_v2/ent"
)

type IRepository interface {
	GetRole(ctx context.Context, name string) (*ent.Role, error)
	GetUserRole(ctx context.Context, userID string) (*ent.Role, error)
	GetRolePermissions(ctx context.Context, roleID string) ([]PermissionDetail, error)
}

type repository struct {
	db *ent.Client
}

func NewRepository(db *ent.Client) IRepository {
	return &repository{db: db}
}
