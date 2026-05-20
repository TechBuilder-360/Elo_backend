package repository

import (
	"context"

	"github.com/Toflex/directory_v2/ent"
)

type IUnitOfWorkRepository interface {
	Begin(ctx context.Context, db *ent.Client) (*ent.Tx, error)
	Commit(db *ent.Tx) error
	Rollback(db *ent.Tx) error
}

type UnitOfWorkRepository struct {
}

func NewUnitOfWorkRepository(db *ent.Client) IUnitOfWorkRepository {
	return &UnitOfWorkRepository{}
}

func (*UnitOfWorkRepository) Begin(ctx context.Context, db *ent.Client) (*ent.Tx, error) {
	tx, err := db.Tx(ctx)
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
