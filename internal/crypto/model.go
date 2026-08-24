package crypto

type SupportedNetwork struct {
	NetworkID  string `json:"network_id"`
	CoinID     string `json:"coin_id"`
	Active     bool   `json:"active"`
	CanSend    bool   `json:"can_send"`
	CanReceive bool   `json:"can_receive"`
}

type CryptoCoin struct {
	ID         string `json:"id"`
	Name       string `json:"name"`
	Code       string `json:"code"`
	Multiplier int64  `json:"multiplier"`
	Active     bool   `json:"active"`
}

type CryptoNetwork struct {
	ID      string
	Name    string
	Slug    string
	Active  bool
	LogoURL *string
}

type StablecoinRequest struct {
	Network string
	// Coin     string
	WalletID string
}

type StablecoinWallet struct {
	Network           string
	Coin              string
	NetworkID         string
	CoinID            string
	Address           string
	ProviderReference string
	Provider          string
	WalletID          string
}

type StablecoinWalletResult struct {
	ID      string
	Network string
	Coin    string
	Address string
}
