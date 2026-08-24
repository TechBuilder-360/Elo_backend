package seed

import (
	"context"
	"strings"

	"entgo.io/ent/dialect/sql"
	"github.com/Toflex/directory_v2/ent"
	"github.com/Toflex/directory_v2/ent/currency"
	"github.com/Toflex/directory_v2/pkg/util"
)

func seedStablecoins(ctx context.Context, db *ent.Client) error {
	networks := map[string]string{
		"eth":      "Ethereum",
		"bsc":      "Binance Smart Chain",
		"polygon":  "Polygon",
		"sol":      "Solana",
		"optimism": "Optimism",
	}
	builders := make([]*ent.StablecoinNetworkCreate, 0)

	for k, v := range networks {
		builders = append(builders, db.StablecoinNetwork.Create().SetName(util.ToTitleCase(v)).SetSlug(strings.ToLower(k)).SetLogoURL("").SetActive(true))
	}

	err := db.StablecoinNetwork.CreateBulk(
		builders...,
	).OnConflict(
		sql.ConflictColumns("slug"),
		sql.DoNothing(),
	).
		Exec(ctx)
	if err != nil {
		return err
	}

	supportedNetworks := map[string][]string{
		"USDC":  {"eth", "bsc", "polygon", "sol", "optimism"},
		"USDT":  {"eth", "bsc", "polygon", "sol"},
		"PYUSD": {"eth", "sol"},
	}

	sn, err := db.StablecoinNetwork.Query().All(ctx)
	if err != nil {
		return err
	}

	networkBuilders := make([]*ent.StablecoinSupportedNetworkCreate, 0)

	for k, v := range supportedNetworks {
		coin, err := db.Currency.Query().Where(currency.CodeEQ(k)).Only(ctx)
		if err != nil {
			return err
		}

		for _, n := range sn {
			if util.Contains(v, n.Slug) {
				networkBuilders = append(networkBuilders, db.StablecoinSupportedNetwork.Create().SetCurrency(coin).SetStablecoinNetwork(n).SetActive(true))
			}
		}
	}

	return db.StablecoinSupportedNetwork.CreateBulk(
		networkBuilders...,
	).OnConflict(
		sql.ConflictColumns("network_id", "coin_id"),
		sql.DoNothing(),
	).
		Exec(ctx)
}
