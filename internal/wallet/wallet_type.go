package wallet

import "strings"

func GetWalletType(wt string) WalletType {
	switch strings.ToUpper(wt) {
	case string(TreasuryWalletType):
		return TreasuryWalletType
	default:
		return TreasuryWalletType
	}
}

func ValidateWalletType(wt string) bool {
	switch strings.ToUpper(wt) {
	case string(TreasuryWalletType):
		return true
	default:
		return false
	}
}
