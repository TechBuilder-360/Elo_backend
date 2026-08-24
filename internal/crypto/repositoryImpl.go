package crypto

import (
	"context"

	"github.com/Toflex/directory_v2/ent"
	"github.com/Toflex/directory_v2/ent/currency"
	"github.com/Toflex/directory_v2/ent/ledgerowner"
	"github.com/Toflex/directory_v2/ent/stablecoinnetwork"
	"github.com/Toflex/directory_v2/ent/stablecoinsupportednetwork"
	"github.com/Toflex/directory_v2/ent/stablecoinwallet"
	"github.com/Toflex/directory_v2/ent/vault"
	"github.com/Toflex/directory_v2/ent/wallet"
	"github.com/Toflex/directory_v2/graph/model"
)

func (r *repository) GetSupportedNetwork(ctx context.Context, networkID, coinID string) (*SupportedNetwork, error) {
	network, err := r.db.StablecoinSupportedNetwork.Query().Where(stablecoinsupportednetwork.CoinIDEQ(coinID), stablecoinsupportednetwork.NetworkIDEQ(networkID)).Only(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, nil
		}

		return nil, err
	}

	return &SupportedNetwork{
		NetworkID:  network.NetworkID,
		CoinID:     network.CoinID,
		Active:     network.Active,
		CanSend:    network.CanSend,
		CanReceive: network.CanReceive,
	}, nil
}

func (r *repository) GetNetworksByCoinID(ctx context.Context, coinID string) ([]*SupportedNetwork, error) {
	networks, err := r.db.StablecoinSupportedNetwork.Query().Where(stablecoinsupportednetwork.CoinIDEQ(coinID)).All(ctx)
	if err != nil {
		return nil, err
	}

	var supportedNetworks []*SupportedNetwork
	for _, network := range networks {
		supportedNetworks = append(supportedNetworks, &SupportedNetwork{
			NetworkID:  network.NetworkID,
			CoinID:     network.CoinID,
			Active:     network.Active,
			CanSend:    network.CanSend,
			CanReceive: network.CanReceive,
		})
	}

	return supportedNetworks, nil
}

func (r *repository) GetCoin(ctx context.Context, coin string) (*CryptoCoin, error) {
	cur, err := r.db.Currency.Query().Where(currency.CodeEqualFold(coin), currency.IsFiatEQ(false)).Only(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, nil
		}

		return nil, err
	}

	return &CryptoCoin{
		ID:         cur.ID,
		Name:       cur.Name,
		Code:       cur.Code,
		Multiplier: (cur.Multiplier),
		Active:     cur.Active,
	}, nil
}

func (r *repository) GetNetwork(ctx context.Context, network string) (*CryptoNetwork, error) {
	chain, err := r.db.StablecoinNetwork.Query().Where(stablecoinnetwork.Or(
		stablecoinnetwork.NameEqualFold(network), stablecoinnetwork.SlugContainsFold(network),
	)).Only(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, nil
		}

		return nil, err
	}

	return &CryptoNetwork{
		ID:      chain.ID,
		Name:    chain.Name,
		Active:  chain.Active,
		Slug:    chain.Slug,
		LogoURL: chain.LogoURL,
	}, nil
}

func (r *repository) CreateWallet(ctx context.Context, payload StablecoinWallet) (*StablecoinWalletResult, error) {
	wallet, err := r.db.StablecoinWallet.Create().
		SetAddress(payload.Address).
		SetCoin(payload.Coin).
		SetCoinID(payload.CoinID).
		SetNetwork(payload.Network).
		SetNetworkID(payload.NetworkID).
		SetProviderReference(payload.ProviderReference).
		SetProvider(payload.Provider).
		SetWalletID(payload.WalletID).
		Save(ctx)
	if err != nil {
		return nil, err
	}

	return &StablecoinWalletResult{
		Address: wallet.Address,
		Coin:    wallet.Coin,
		Network: wallet.Network,
		ID:      wallet.ID,
	}, nil
}

func (r *repository) GetWalletByCoinAndNetwork(ctx context.Context, walletID, coinID, networkID string) (*StablecoinWalletResult, error) {
	sw, err := r.db.StablecoinWallet.Query().Where(
		stablecoinwallet.HasWalletWith(wallet.IDEQ(walletID)),
		stablecoinwallet.CoinIDEQ(coinID),
		stablecoinwallet.NetworkIDEQ(networkID),
		stablecoinwallet.Disabled(false),
	).Only(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, nil
		}

		return nil, err
	}

	return &StablecoinWalletResult{
		Address: sw.Address,
		Coin:    sw.Coin,
		Network: sw.Network,
		ID:      sw.ID,
	}, nil
}

func (r *repository) GetAddressByWallet(ctx context.Context, ownerID string, filter *model.StablecoinFilter) ([]StablecoinWalletResult, error) {
	query := r.db.StablecoinWallet.Query().Where(
		stablecoinwallet.HasWalletWith(wallet.HasVaultWith(vault.HasOwnerWith(ledgerowner.IDEQ(ownerID)))),
		stablecoinwallet.Disabled(false),
	)

	if filter != nil {
		if filter.WalletID != nil {
			query.Where(stablecoinwallet.HasWalletWith(wallet.IDEQ(*filter.WalletID)))
		}

		if filter.Coin != nil {
			query.Where(stablecoinwallet.HasCurrencyWith(currency.CodeEQ(*filter.Coin)))
		}

		if filter.Network != nil {
			query.Where(stablecoinwallet.HasStablecoinNetworksWith(stablecoinnetwork.SlugEQ(*filter.Network)))
		}
	}

	addresses, err := query.All(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, nil
		}

		return nil, err
	}

	result := make([]StablecoinWalletResult, 0)

	for _, address := range addresses {
		result = append(result, StablecoinWalletResult{
			Address: address.Address,
			Coin:    address.Coin,
			Network: address.Network,
			ID:      address.ID,
		})
	}

	return result, nil
}
