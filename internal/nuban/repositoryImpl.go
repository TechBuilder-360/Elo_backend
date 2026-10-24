package nuban

import (
	"context"

	"github.com/Toflex/directory_v2/ent"
	"github.com/Toflex/directory_v2/ent/nubanbank"
	"github.com/Toflex/directory_v2/ent/nubanbeneficiary"
	"github.com/Toflex/directory_v2/ent/nubanblacklist"
)

func (r *repository) GetNubanBlacklist(ctx context.Context, accountNumber, bankID string) (*beneficiaryBlacklist, error) {
	blacklist, err := r.db.NubanBlacklist.Query().
		Where(nubanblacklist.AccountNumberEQ(accountNumber), nubanblacklist.NubanBankIDEQ(bankID)).First(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, nil
		}
		return nil, err
	}

	response := beneficiaryBlacklist{blacklist}

	return &response, nil
}

func (r *repository) GetBankBySlug(ctx context.Context, slug string) (*ent.NubanBank, error) {
	bank, err := r.db.NubanBank.Query().
		Where(nubanbank.SlugEqualFold(slug)).WithProviderBanks().First(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, nil
		}
		return nil, err
	}

	return bank, nil
}

func (r *repository) GetBankByID(ctx context.Context, id string) (*ent.NubanBank, error) {
	bank, err := r.db.NubanBank.Query().
		Where(nubanbank.IDEQ(id)).WithProviderBanks().First(ctx)
	if err != nil {
		return nil, err
	}

	return bank, nil
}

func (r *repository) GetBeneficiaryByAccountNumber(ctx context.Context, ownerID, accountNumber, bankID string) (*beneficiary, error) {
	ben, err := r.db.NubanBeneficiary.
		Query().
		Where(nubanbeneficiary.AccountNumber(accountNumber),
			nubanbeneficiary.NubanBankID(bankID),
			nubanbeneficiary.OwnerIDEQ(ownerID),
		).First(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, nil
		}

		return nil, err
	}

	return &beneficiary{
		ID:            ben.ID,
		AccountName:   ben.AccountName,
		AccountNumber: ben.AccountNumber,
		BankName:      ben.BankName,
		BankID:        ben.NubanBankID,
		Disabled:      ben.Disabled,
	}, nil
}

func (r *repository) CreateBeneficiary(ctx context.Context, payload *createBeneficiary) (*beneficiary, error) {
	ben, err := r.db.NubanBeneficiary.Create().
		SetAccountName(payload.AccountName).
		SetAccountNumber(payload.AccountNumber).
		SetBankName(payload.BankName).
		SetNubanBankID(payload.BankID).
		SetOwnerID(payload.OwnerID).
		Save(ctx)

	if err != nil {
		return nil, err
	}

	return &beneficiary{
		ID:            ben.ID,
		AccountName:   ben.AccountName,
		AccountNumber: ben.AccountNumber,
		BankName:      ben.BankName,
		BankID:        ben.NubanBankID,
		Disabled:      ben.Disabled,
	}, nil
}

func (r *repository) WithTransaction(tx *ent.Tx) IRepository {
	return &repository{
		db: tx.Client(),
	}
}
