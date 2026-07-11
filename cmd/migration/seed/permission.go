package seed

import (
	"context"

	"entgo.io/ent/dialect/sql"
	"github.com/Toflex/directory_v2/ent"
	"github.com/Toflex/directory_v2/ent/role"
	rbac "github.com/Toflex/directory_v2/pkg/RBAC"
)

func seedRoles(ctx context.Context, db *ent.Client) error {
	return db.Role.CreateBulk(
		db.Role.Create().SetName("Super Admin").SetDescription("Super Admin has full access"),
		db.Role.Create().SetName("Admin").SetDescription("Admin has limited access"),
		db.Role.Create().SetName("Manager").SetDescription("Manager has intermediate access"),
	).OnConflict(
		sql.ConflictColumns(role.FieldName),
		sql.DoNothing(),
	).Exec(ctx)
}

func seedRolePermissions(ctx context.Context, db *ent.Client) error {
	// Implementation for seeding role permissions
	roles, err := db.Role.Query().All(ctx)
	if err != nil {
		return err
	}

	permissions, err := db.Permission.Query().All(ctx)
	if err != nil {
		return err
	}

	var rolesMap = make(map[string]string)
	for _, r := range roles {
		rolesMap[r.Name] = r.ID
	}

	var permissionsMap = make(map[string]string)
	for _, p := range permissions {
		permissionsMap[p.Name] = p.ID
	}

	var rolePermissions = []*ent.RolePermissionCreate{}

	for roleName, permissions := range rbac.RolePermissions {
		roleID, ok := rolesMap[string(roleName)]
		if !ok {
			continue
		}

		for _, permission := range permissions {
			permissionID, ok := permissionsMap[string(permission)]
			if !ok {
				continue
			}

			rolePermissions = append(rolePermissions, db.RolePermission.Create().
				SetRoleID(roleID).
				SetPermissionID(permissionID))
		}
	}

	return db.RolePermission.CreateBulk(rolePermissions...).OnConflict(
		sql.ConflictColumns("role_id", "permission_id"),
		sql.DoNothing(),
	).Exec(ctx)
}

func seedPermissions(ctx context.Context, db *ent.Client) error {
	var permissions = []*ent.PermissionCreate{}

	for permission, description := range rbac.AllPermissions {
		permissions = append(permissions, db.Permission.Create().SetName(string(permission)).SetDescription(description))
	}

	return db.Permission.CreateBulk(permissions...).OnConflict(
		sql.ConflictColumns("name"),
		sql.DoNothing(),
	).Exec(ctx)
}
