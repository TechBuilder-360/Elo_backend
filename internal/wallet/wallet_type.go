package wallet

import "strings"

func GetWalletType(wt string) WalletType {
	switch strings.ToUpper(wt) {
	case string(TreasuryWalletType):
		return TreasuryWalletType
	case string(FrozenWalletType):
		return FrozenWalletType
	case string(HoldingWalletType):
		return HoldingWalletType
	default:
		return TreasuryWalletType
	}
}

func ValidateWalletType(wt string) bool {
	switch strings.ToUpper(wt) {
	case string(TreasuryWalletType), string(FrozenWalletType), string(HoldingWalletType):
		return true
	default:
		return false
	}
}
