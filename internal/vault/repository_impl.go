package vault

import (
	"context"

	"github.com/Toflex/directory_v2/ent"
	"github.com/Toflex/directory_v2/ent/business"
	"github.com/Toflex/directory_v2/ent/ledgerowner"
	"github.com/Toflex/directory_v2/ent/vault"
)

// GetByID implements IRepository.
func (r *repository) GetByID(ctx context.Context, id string) (*ent.Vault, error) {
	v, err := r.db.
		Vault.
		Query().
		Where(vault.IDEQ(id)).
		First(ctx)

	if err != nil {
		if ent.IsNotFound(err) {
			return nil, nil
		}

		return nil, err
	}

	return v, nil
}

// GetBusinessVaultByType implements IRepository.
func (r *repository) GetBusinessVaultByType(ctx context.Context, b *ent.Business, vaultType string) (*ent.Vault, error) {
	v, err := r.db.
		Vault.
		Query().
		Where(
			vault.HasOwnerWith(ledgerowner.HasBusinessWith(business.IDEQ(b.ID))),
			vault.TypeEQ(vault.Type(vaultType)),
		).
		First(ctx)

	if err != nil {
		if ent.IsNotFound(err) {
			return nil, nil
		}

		return nil, err
	}

	return v, nil
}

func (r *repository) Create(ctx context.Context, ownerID string) (*ent.Vault, error) {
	return r.CreateWithType(ctx, ownerID, string(vault.DefaultType))
}

func (r *repository) CreateWithType(ctx context.Context, ownerID, vaultType string) (*ent.Vault, error) {
	vt := vault.Type(vaultType)
	if vaultType == "" {
		vt = vault.DefaultType
	}

	v, err := r.db.Vault.Create().SetType(vt).SetOwnerID(ownerID).Save(ctx)
	if err != nil || v == nil {
		return nil, err
	}

	return v, nil
}

func (r *repository) GetByOwner(ctx context.Context, ownerID, vaultType string) (*ent.Vault, error) {
	// create default Treasury vault
	v, err := r.db.Vault.Query().
		Where(vault.HasOwnerWith(
			ledgerowner.IDEQ(ownerID)),
			vault.TypeEQ(vault.Type(vaultType))).
		First(ctx)
	if err != nil || v == nil {
		if ent.IsNotFound(err) {
			return nil, nil
		}
		return nil, err
	}

	return v, nil
}

func (r *repository) WithTransaction(tx *ent.Tx) IRepository {
	return &repository{
		db: tx.Client(),
	}
}
