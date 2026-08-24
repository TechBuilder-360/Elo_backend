package mock

import (
	"context"

	"crypto/ecdsa"
	"crypto/rand"

	pkg "github.com/Toflex/directory_v2/internal/crypto"
	"github.com/Toflex/directory_v2/pkg/util"

	"github.com/ethereum/go-ethereum/crypto"
)

func FakeEthereumAddress() (string, error) {
	privateKey, err := ecdsa.GenerateKey(crypto.S256(), rand.Reader)
	if err != nil {
		return "", err
	}

	return crypto.PubkeyToAddress(privateKey.PublicKey).Hex(), nil
}

func (m *mock) GenerateCryptoAddress(ctx context.Context, payload pkg.CryptoRequest) (*pkg.CryptoResponse, error) {
	address, err := FakeEthereumAddress()
	if err != nil {
		return nil, err
	}

	return &pkg.CryptoResponse{
		Address:          address,
		Coin:             payload.Coin,
		Network:          payload.Network,
		PartnerReference: util.GenerateRandomString(15), // Generate a random partner reference
	}, nil
}
