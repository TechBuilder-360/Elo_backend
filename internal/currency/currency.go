package currency

import (
	"context"

	"github.com/Toflex/directory_v2/graph/model"
	"github.com/Toflex/directory_v2/pkg/errors"
	"github.com/Toflex/directory_v2/pkg/log"
	"github.com/Toflex/directory_v2/pkg/types"
)

// GetCurrencies implements [IService].
func (s *service) GetCurrencies(ctx context.Context) ([]*model.Currency, error) {
	logger := log.LoggerInContext(ctx)
	currencies, err := s.repo.Currencyies(ctx)
	if err != nil {
		logger.WithError(err).Error("failed to fetch currencies")
		return nil, errors.New(errors.ErrFailed, "request failed")
	}

	result := make([]*model.Currency, 0)

	for _, currency := range currencies {
		result = append(result, &model.Currency{
			Name:   currency.Name,
			Code:   currency.Code,
			Symbol: currency.Symbol,
			ID:     currency.ID,
			IsFiat: currency.IsFiat,
		})
	}

	return result, nil
}

// GetCurrencyByCode implements [IService].
func (s *service) GetCurrencyByCode(ctx context.Context, code types.CurrencyCode) (*Currency, error) {
	logger := log.LoggerInContext(ctx)
	currency, err := s.repo.GetCurrencyByCode(ctx, code)
	if err != nil || currency == nil {
		logger.WithError(err).WithField("currency_code", code).Error("failed to fetch currency by code")
		return nil, errors.New(errors.ErrFailed, "request failed")
	}

	if !currency.Active {
		logger.WithField("currency_active", currency.Active).Error("currency is not active")
		return nil, errors.New(errors.ErrFailed, "currency is not available")
	}

	return &Currency{
		Name:       currency.Name,
		Code:       currency.Code,
		Symbol:     currency.Symbol,
		ID:         currency.ID,
		IsFiat:     currency.IsFiat,
		Multiplier: currency.Multiplier,
	}, nil
}
