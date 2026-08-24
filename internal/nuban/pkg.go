package nuban

import (
	"context"
	"time"
)

type AccountType string

const (
	StaticNUBANAccountType  AccountType = "static"
	DynamicNUBANAccountType AccountType = "dynamic"
)

type NUBANProvider interface {
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
