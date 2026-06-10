package unitofwork

import (
	"context"

	"github.com/Toflex/directory_v2/ent"
)

type IUnitOfWorkRepository interface {
	Begin(ctx context.Context) (*ent.Tx, error)
	Commit(db *ent.Tx) error
	Rollback(db *ent.Tx) error
}

type UnitOfWorkRepository struct {
	db *ent.Client
}

func NewUnitOfWorkRepository(db *ent.Client) IUnitOfWorkRepository {
	return &UnitOfWorkRepository{
		db: db,
	}
}

func (u *UnitOfWorkRepository) Begin(ctx context.Context) (*ent.Tx, error) {
	tx, err := u.db.Tx(ctx)
	if err != nil {
		return nil, err
	}

	return tx, nil
}

func (*UnitOfWorkRepository) Commit(tx *ent.Tx) error {
	return tx.Commit()
}

func (*UnitOfWorkRepository) Rollback(tx *ent.Tx) error {
	return tx.Rollback()
}
