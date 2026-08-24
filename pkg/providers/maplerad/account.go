package maplerad

import (
	"context"
	"fmt"

	httpclient "github.com/Toflex/directory_v2/pkg/apm"
	"github.com/Toflex/directory_v2/pkg/constant"
	"github.com/Toflex/directory_v2/pkg/errors"
	"github.com/Toflex/directory_v2/pkg/log"
)

type accountRequest struct {
	Currency    string `json:"currency"`
	AccountName string `json:"account_name"`
}

type accountResult struct {
	Status  bool        `json:"status"`
	Message string      `json:"message"`
	Data    accountData `json:"data"`
}

type accountData struct {
	ID            string `json:"id"`
	BankName      string `json:"bank_name"`
	AccountNumber string `json:"account_number"`
	AccountName   string `json:"account_name"`
	Currency      string `json:"currency"`
	CreatedAt     string `json:"created_at"`
}

func (m *maplerad) requestStaticAccount(ctx context.Context, payload *accountRequest) (*accountResult, error) {
	logger := log.LoggerInContext(ctx)
	logger.WithField("partner", constant.Maplerad)

	result := new(accountResult)

	// create a resty client
	client := httpclient.HTTPClient()

	resp, err := client.R().
		SetHeader("Content-Type", "application/json").
		SetAuthToken(m.config.SecretKey).
		SetBody(payload).
		SetResult(result).
		SetError(&errorResponse{}).
		Post(fmt.Sprintf("%s/v1/merchant/collections", m.config.BaseURL))

	if err != nil {
		logger.WithError(err).Error("failed to generate account")
		return nil, err
	}

	if resp.IsError() {
		logger.WithField("response", resp.String()).Error("failed to generate account")
		return nil, errors.New(errors.ErrFailed, "request failed")
	}

	return result, nil
}
