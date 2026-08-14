package account

import (
	"context"
	"time"

	"github.com/Toflex/directory_v2/pkg/types"
)

type StaticAccountRequest struct {
	AccountName string
	Reference   string
	Currency    types.CurrencyCode
}

type StaticAccountResult struct {
	AccountName      string
	AccountNumber    string
	BankName         string
	Reference        string
	PartnerReference string
}

type DynamicAccountResult struct {
	AccountName      string
	AccountNumber    string
	BankName         string
	Reference        string
	PartnerReference string
	ExpirationTime   time.Time
}

type DynamicAccountRequest struct {
	AccountName string
	Reference   string
	Currency    types.CurrencyCode
}

type StaticAccount interface {
	GenerateAccount(context.Context, StaticAccountRequest) (*StaticAccountResult, error)
}

type DynamicAccount interface {
	GenerateAccount(context.Context, DynamicAccountRequest) (*DynamicAccountResult, error)
}
