package nuban

import (
	"context"

	"github.com/Toflex/directory_v2/ent"
)

type IRepository interface {
	GetBankBySlug(ctx context.Context, slug string) (*ent.NubanBank, error)
	GetBankByID(ctx context.Context, id string) (*ent.NubanBank, error)
	GetBeneficiaryByAccountNumber(ctx context.Context, ownerID, accountNumber, bankID string) (*beneficiary, error)
	CreateBeneficiary(ctx context.Context, payload *createBeneficiary) (*beneficiary, error)
	GetNubanBlacklist(ctx context.Context, accountNumber, bankID string) (*beneficiaryBlacklist, error)
	WithTransaction(tx *ent.Tx) IRepository
}

type repository struct {
	db *ent.Client
}

func NewRepository(db *ent.Client) IRepository {
	return &repository{
		db: db,
	}
}
