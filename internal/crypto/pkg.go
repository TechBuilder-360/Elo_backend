package crypto

import "context"

type Coin string
type Network string

const (
	USDT  Coin = "USDT"
	PYUSD Coin = "PYUSD"
	USDC  Coin = "USDC"

	EthereumNetwork Network = "eth"
	SolanaNetwork   Network = "sol"
	TronNetwork     Network = "trx"
	BSCNetwork      Network = "bsc"
	PolygonNetwork  Network = "polygon"
)

type CryptoProvider interface {
	GenerateCryptoAddress(ctx context.Context, payload CryptoRequest) (*CryptoResponse, error)
}

type CryptoRequest struct {
	Coin      Coin    `json:"coin"`
	Network   Network `json:"network"`
	Reference string  `json:"reference"`
}

type CryptoResponse struct {
	Address          string  `json:"address"`
	Coin             Coin    `json:"coin"`
	Network          Network `json:"network"`
	PartnerReference string  `json:"partner_reference"`
}
