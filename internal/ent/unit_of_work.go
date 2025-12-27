package repository

import (
	"context"

	"github.com/Toflex/directory_v2/database/database"
)

type IUnitOfWorkRepository interface {
	Begin(ctx context.Context, db *database.Client) (*database.Client, error)
	Commit(db *database.Client) error
	Rollback(db *database.Client) error
}

type UnitOfWorkRepository struct {
	db *database.Client
}

func NewUnitOfWorkRepository(db *database.Client) IUnitOfWorkRepository {
	return &UnitOfWorkRepository{
		db: db,
	}
}

func (*UnitOfWorkRepository) Begin(ctx context.Context, db *database.Client) (*database.Client, error) {
	tx, err := db.DBClient.Tx(ctx)
	if err != nil {
		return nil, err
	}

	return &database.Client{DBClient: tx.Client(), DBTransaction: tx}, nil
}

func (*UnitOfWorkRepository) Commit(db *database.Client) error {
	return db.DBTransaction.Commit()
}

func (*UnitOfWorkRepository) Rollback(db *database.Client) error {
	return db.DBTransaction.Rollback()
}
