package verification

import (
	"github.com/Toflex/directory_v2/ent"
)

type IRepository interface {
}

type repository struct {
	db *ent.Client
}

func NewRepository(db *ent.Client) IRepository {
	return &repository{
		db: db,
	}
}
