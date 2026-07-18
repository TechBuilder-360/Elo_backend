package wallet

import "strings"

func getWalletType(wt string) WalletType {
	switch strings.ToUpper(wt) {
	case string(TreasuryWalletType):
		return TreasuryWalletType
	default:
		return TreasuryWalletType
	}
}

func validateWalletType(wt string) bool {
	switch strings.ToUpper(wt) {
	case string(TreasuryWalletType):
		return true
	default:
		return false
	}
}
