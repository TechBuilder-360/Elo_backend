package rbac

import (
	"context"

	"github.com/Toflex/directory_v2/ent"
	"github.com/Toflex/directory_v2/ent/manager"
	"github.com/Toflex/directory_v2/ent/permission"
	"github.com/Toflex/directory_v2/ent/role"
	"github.com/Toflex/directory_v2/ent/rolepermission"
)

// GetRolePermissions implements [IRepository].
func (r *repository) GetRolePermissions(ctx context.Context, roleID string) ([]PermissionDetail, error) {
	var result []PermissionDetail
	err := r.db.RolePermission.Query().Where(rolepermission.RoleID(roleID)).
		QueryPermission().
		Select(permission.FieldID, permission.FieldName, permission.FieldDescription).
		Scan(ctx, &result)
	if err != nil {
		return nil, err
	}

	return result, nil
}

// GetUserRole implements [IRepository].
func (r *repository) GetUserRole(ctx context.Context, userID string) (*ent.Role, error) {
	role, err := r.db.Manager.Query().
		Where(manager.UserIDEQ(userID)).QueryRole().
		Only(ctx)
	if err != nil {
		return nil, err
	}
	return role, nil
}

// GetRole implements [IRepository].
func (r *repository) GetRole(ctx context.Context, name string) (*ent.Role, error) {
	role, err := r.db.Role.Query().Where(role.NameEQ(name)).Only(ctx)
	if err != nil {
		return nil, err
	}
	return role, nil
}
