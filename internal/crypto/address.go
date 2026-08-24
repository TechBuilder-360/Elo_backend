package crypto

import (
	"context"
	"fmt"

	"github.com/Toflex/directory_v2/ent"
	"github.com/Toflex/directory_v2/graph/model"
	"github.com/Toflex/directory_v2/pkg/errors"
	"github.com/Toflex/directory_v2/pkg/log"
	"github.com/Toflex/directory_v2/pkg/provider"
	"github.com/Toflex/directory_v2/pkg/util"
)

// GenerateCryptoAddress implements [IService].
func (s *service) GenerateCryptoAddress(ctx context.Context, b *ent.Business, body *StablecoinRequest) (*CryptoAddressResponse, error) {
	logger := log.LoggerInContext(ctx)

	// verify that wallet belongs to business
	w, err := s.walletService.GetWalletOwnerWithID(ctx, b.Edges.Owner.ID, body.WalletID)
	if err != nil || w == nil {
		logger.WithError(err).Error("wallet could not be found for business")
		return nil, errors.New(errors.ErrFailed, "wallet not found")
	}

	if !w.Active {
		return nil, errors.New(errors.ErrFailed, "wallet is not active")
	}

	// Get Network
	network, err := s.repo.GetNetwork(ctx, body.Network)
	if err != nil || network == nil {
		logger.WithField("network", body.Network).WithError(err).Error("failed to fetch stablecoin network")
		return nil, errors.New(errors.ErrFailed, "network not found")
	}

	if !network.Active {
		return nil, errors.New(errors.ErrFailed, "network is inactive")
	}

	// Get Coin
	coin, err := s.repo.GetCoin(ctx, w.Currency)
	if err != nil || coin == nil {
		logger.WithField("coin", w.Currency).WithError(err).Error("failed to fetch stablecoin coin")
		return nil, errors.New(errors.ErrFailed, "coin not found")
	}

	if !coin.Active {
		return nil, errors.New(errors.ErrFailed, "coin is inactive")
	}

	sw, _ := s.repo.GetWalletByCoinAndNetwork(ctx, body.WalletID, coin.ID, network.ID)
	if sw != nil {
		return &CryptoAddressResponse{
			Address: sw.Address,
			Coin:    Coin(sw.Coin),
			Network: Network(sw.Network),
			ID:      sw.ID,
		}, nil
	}

	// Is network supported
	supportedNetwork, err := s.repo.GetSupportedNetwork(ctx, network.ID, coin.ID)
	if err != nil || supportedNetwork == nil {
		logger.WithField("network", body.Network).WithError(err).Error("failed to fetch stablecoin supported network")
		return nil, errors.New(errors.ErrFailed, "network coin not supported")
	}

	if !supportedNetwork.Active {
		return nil, errors.New(errors.ErrFailed, "network coin is inactive")
	}

	if !supportedNetwork.CanReceive {
		return nil, errors.New(errors.ErrFailed, "network --> coin not available")
	}

	identifier := fmt.Sprintf("stablecoin-address-%s", network.Slug)
	serv, err := s.getServiceProvider(ctx, identifier)
	if err != nil {
		logger.WithField("identifier", identifier).WithError(err).Error("failed to fetch stablecoin address provider")
		return nil, errors.New(errors.ErrFailed, "request failed")
	}

	impl, ok := provider.GetImpl(serv.ActiveProvider.Slug)
	if !ok {
		logger.WithField("provider", serv.ActiveProvider.Slug).WithError(err).Error("failed to fetch stablecoin address implementation")
		return nil, errors.New(errors.ErrFailed, "request failed")
	}

	c, ok := provider.ConformsTo[CryptoProvider](impl)
	if !ok {
		logger.WithField("provider", serv.ActiveProvider.Slug).WithError(err).Error("failed to fetch stablecoin address implementation")
		return nil, errors.New(errors.ErrFailed, "request failed")
	}

	// request for address from provider
	request := CryptoRequest{
		Coin:      Coin(coin.Code),
		Network:   Network(network.Slug),
		Reference: util.GenerateCUID(),
	}
	result, err := c.GenerateCryptoAddress(ctx, request)
	if err != nil {
		logger.WithField("payload", request).WithError(err).Error("failed to generate stablecoin address")
		return nil, errors.New(errors.ErrFailed, "request failed")
	}

	wallet, err := s.repo.CreateWallet(ctx, StablecoinWallet{
		Coin:              coin.Code,
		Network:           network.Slug,
		CoinID:            coin.ID,
		NetworkID:         network.ID,
		Provider:          serv.ActiveProvider.Slug,
		ProviderReference: result.PartnerReference,
		Address:           result.Address,
		WalletID:          w.ID,
	})
	if err != nil {
		logger.WithError(err).Error("failed to store stablecoin address")
		return nil, errors.New(errors.ErrFailed, "request failed")
	}

	return &CryptoAddressResponse{
		Address: wallet.Address,
		Coin:    Coin(wallet.Coin),
		Network: Network(wallet.Network),
		ID:      wallet.ID,
	}, nil

}

func (s *service) getServiceProvider(ctx context.Context, identifier string) (provider.ServiceLocator, error) {
	return provider.NewService().GetServiceByIdentifier(ctx, identifier)
}

func (s *service) GetCryptoAddress(ctx context.Context, ownerID string, filter *model.StablecoinFilter) ([]CryptoAddressResponse, error) {
	logger := log.LoggerInContext(ctx)

	result := make([]CryptoAddressResponse, 0)

	sws, err := s.repo.GetAddressByWallet(ctx, ownerID, filter)
	if err != nil {
		logger.WithError(err).Error("caould not fetch stablecoin wallets")
		return result, nil
	}

	for _, wallet := range sws {
		result = append(result, CryptoAddressResponse{
			Address: wallet.Address,
			Coin:    Coin(wallet.Coin),
			Network: Network(wallet.Network),
			ID:      wallet.ID,
		})
	}

	return result, nil
}
