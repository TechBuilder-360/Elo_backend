package wallet

type BalanceType string

const (
	AvailableBalance BalanceType = "AVAILABLE_BALANCE"
	LedgerBalance    BalanceType = "LEDGER_BALANCE"
	HoldingBalance   BalanceType = "HOLDING_BALANCE"
)
