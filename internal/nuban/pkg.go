package nuban

import (
	"context"
	"time"

	"github.com/Toflex/directory_v2/pkg/types"
)

type AccountType string

const (
	StaticNUBANAccountType  AccountType = "static"
	DynamicNUBANAccountType AccountType = "dynamic"
)

type NUBANDepositProvider interface {
	GenerateStaticNUBANAccountNumber(ctx context.Context, payload *StaticNubanRequest) (*StaticNubanResponse, error)
	GenerateDynamicNUBANAccountNumber(ctx context.Context, payload *DynamicNubanRequest) (*DynamicNubanResponse, error)
}

type StaticNubanRequest struct {
	AccountName string   `json:"account_name"`
	Customer    Customer `json:"customer"`
	Reference   string   `json:"reference"`
}

type DynamicNubanRequest struct {
	AccountName string `json:"account_name"`
	Amount      int64  `json:"amount"`
	Reference   string `json:"reference"`
}

type DynamicNubanResponse struct {
	AccountNumber    string      `json:"account_number"`
	AccountName      string      `json:"account_name"`
	BankCode         string      `json:"bank_code"`
	BankName         string      `json:"bank_name"`
	PartnerReference string      `json:"partner_reference"`
	Type             AccountType `json:"type"`
	Amount           int64       `json:"amount"`
	ExpiryDate       time.Time   `json:"expiry_date"`
}

type Customer struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email"`
	Phone     string `json:"phone"`
	BVN       string `json:"bvn"`
}

type StaticNubanResponse struct {
	AccountNumber    string      `json:"account_number"`
	AccountName      string      `json:"account_name"`
	BankCode         string      `json:"bank_code"`
	BankName         string      `json:"bank_name"`
	PartnerReference string      `json:"partner_reference"`
	Type             AccountType `json:"type"`
}

type ResolveAccountPayload struct {
	AccountNumber string
	BankCode      string
}

type ResolveAccountResult struct {
	AccountNumber string
	AccountName   string
	BankCode      string
}

type NubanTransferPayload struct {
	Reference     string
	Narration     string
	Amount        uint64
	AccountName   string
	AccountNumber string
	BankCode      string
}

type NubanTransferResult struct {
	PartnerID string
	SessionID *string
	Reference string
	Status    types.TransactionStatus
}

type NUBANTransferProvider interface {
	ResolveAccount(ctx context.Context, payload ResolveAccountPayload) (*ResolveAccountResult, error)
	Transfer(ctx context.Context, payload *NubanTransferPayload) (*NubanTransferResult, error)
	GetTransactionStatus(ctx context.Context, id string) (*NubanTransferResult, error)
}
