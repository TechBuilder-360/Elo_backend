package crypto

type CryptoAddressResponse struct {
	Address string  `json:"address"`
	Coin    Coin    `json:"coin"`
	Network Network `json:"network"`
	ID      string  `json:"id"`
}
