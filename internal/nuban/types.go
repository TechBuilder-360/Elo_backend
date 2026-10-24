package nuban

import (
	"github.com/Toflex/directory_v2/ent"
	"github.com/Toflex/directory_v2/pkg/types"
)

type ResolveAccountNamePayload struct {
	AccountNumber string
	BankCode      string
}

type ResolveAccountNameResult struct {
	AccountNumber string
	AccountName   string
	BankCode      string
}

type ProcessTransferPayload struct {
	TransactionID string
	Narration     *string
	Amount        uint64
	Provider      string
	Currency      string
	Beneficiary   *nubanBeneficiaryResult
	Bank          *bankResponse
}

type TransferResult struct {
	ID     string
	Status string
}

type beneficiary struct {
	AccountName   string
	AccountNumber string
	BankName      string
	BankID        string
	Nickname      string
	ID            string
	Disabled      bool
}

type createBeneficiary struct {
	AccountName   string
	AccountNumber string
	BankName      string
	BankID        string
	Nickname      string
	OwnerID       string
}

type beneficiaryBlacklist struct {
	*ent.NubanBlacklist
}

type transactionUpdatePayload struct {
	ID               string
	Status           types.TransactionStatus
	SessionID        *string
	PartnerReference *string
}
