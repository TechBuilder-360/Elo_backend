package wallet

import (
	"github.com/Toflex/directory_v2/ent"
	"github.com/Toflex/directory_v2/pkg/types"
)

type WalletType string

type WalletFilter struct {
	IsFiat *bool
}

type DebitWalletPayload struct {
	OwnerID               string
	WalletType            string
	DebitAmount           int64
	Transaction           *ent.Transaction
	ReversalTransactionID *string
	IsReversal            bool
	DBTransaction         *ent.Tx
	AffectedBalance       BalanceType
}

type CreditWalletPayload struct {
	OwnerID               string
	WalletType            string
	CreditAmount          int64
	Transaction           *ent.Transaction
	ReversalTransactionID *string
	IsReversal            bool
	DBTransaction         *ent.Tx
	AffectedBalance       BalanceType
}

type ValidateWalletPayload struct {
	DBTransaction *ent.Tx
	OwnerID       string
	WalletType    string
	Currency      types.CurrencyCode
}

type LedgerEntry struct {
	TransactionID         string
	ReversalTransactionID *string
	IsReversal            bool
	CurrentBalance        int64
	PreviousBalance       int64
	AffectedBalance       BalanceType
	Debit                 int64
	Credit                int64
	Wallet                *WalletResponse
	DBTransaction         *ent.Tx
}
