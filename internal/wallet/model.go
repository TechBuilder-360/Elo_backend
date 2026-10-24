package wallet

import (
	"github.com/Toflex/directory_v2/internal/currency"
	"github.com/Toflex/directory_v2/pkg/types"
)

const (
	TreasuryWalletType WalletType = "TREASURY"
	FrozenWalletType   WalletType = "FROZEN"
	HoldingWalletType  WalletType = "HOLDING"
)

type WalletResponse struct {
	ID               string
	AvailableBalance int64
	LedgerBalance    int64
	HoldingBalance   int64
	Active           bool
	Currency         types.CurrencyCode
	Multiplier       int64
	IsFiat           bool
}

type createWallet struct {
	VaultID  string
	Currency *currency.Currency
}

type OwnerVaultOverview struct {
	LedgerOwnerStatus string
	VaultStatus       string
	VaultID           string
	Wallet            *WalletResponse
}
