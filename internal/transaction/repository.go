package transaction

import (
	"context"

	"github.com/Toflex/directory_v2/ent"
)

type IRepository interface {
	GetTransactionByID(ctx context.Context, id string) (*ent.Transaction, error)
}

type repository struct {
	db *ent.Client
}

func Newrepository(db *ent.Client) IRepository {
	return &repository{
		db: db,
	}
}

