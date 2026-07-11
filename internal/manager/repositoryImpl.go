package manager

import (
	"context"

	"github.com/Toflex/directory_v2/ent"
)

// AddManager implements [IRepository].
func (r *repository) AddManager(ctx context.Context, manager *Manager) error {
	_, err := r.db.Manager.Create().
		SetUserID(manager.UserID).
		SetBusinessID(manager.BusinessID).
		SetIsOwner(manager.IsOwner).
		SetRoleID(manager.RoleID).
		Save(ctx)
	if err != nil {
		return err
	}

	return nil
}

func (r *repository) WithTransaction(tx *ent.Tx) IRepository {
	return &repository{db: tx.Client()}
}
