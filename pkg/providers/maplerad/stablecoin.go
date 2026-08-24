package maplerad

import (
	"context"
	"fmt"

	httpclient "github.com/Toflex/directory_v2/pkg/apm"
	"github.com/Toflex/directory_v2/pkg/constant"
	"github.com/Toflex/directory_v2/pkg/errors"
	"github.com/Toflex/directory_v2/pkg/log"
)

type stablecoinRequest struct {
	Coin  string `json:"coin"`
	Chain string `json:"chain"`
}

type stablecoinResponse struct {
	Status  bool           `json:"status"`
	Message string         `json:"message"`
	Data    stablecoinData `json:"data"`
}

type stablecoinData struct {
	ID              string `json:"id"`
	Address         string `json:"address"`
	Chain           string `json:"chain"`
	Coin            string `json:"coin"`
	Offramp         bool   `json:"offramp"`
	Active          bool   `json:"active"`
	FundDestination string `json:"fund_destination"`
	CreatedAt       string `json:"created_at"`
}

func (m *maplerad) requestAddress(ctx context.Context, payload *stablecoinRequest) (*stablecoinResponse, error) {
	logger := log.LoggerInContext(ctx)
	logger.WithField("partner", constant.Maplerad)

	result := new(stablecoinResponse)

	// create a resty client
	client := httpclient.HTTPClient()

	resp, err := client.R().
		SetHeader("Content-Type", "application/json").
		SetAuthToken(m.config.SecretKey).
		SetBody(payload).
		SetResult(result).
		SetError(&errorResponse{}).
		Post(fmt.Sprintf("%s/v1/merchant/crypto", m.config.BaseURL))

	if err != nil {
		logger.WithError(err).Error("failed to generate stablecoin address")
		return nil, err
	}

	if resp.IsError() {
		logger.WithField("response", resp.String()).Error("failed to generate stablecoin address")
		return nil, errors.New(errors.ErrFailed, "request failed")
	}

	return result, nil
}
