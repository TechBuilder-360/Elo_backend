package wallet

import (
	"github.com/Toflex/directory_v2/ent"
	"github.com/Toflex/directory_v2/pkg/types"
)

const (
	TreasuryWalletType WalletType = "TREASURY"
)

type WalletResponse struct {
	ID               string
	Type             string
	AvailableBalance int64
	LedgerBalance    int64
	HoldingBalance   int64
	Owner            string
	Identifier       string
	Active           bool
	Currency         types.CurrencyCode
	Multiplier       int64
}

type createWallet struct {
	Type       WalletType
	IsBusiness bool
	Identifier string
	Currency   *ent.Currency
}
