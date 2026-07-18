package currency

import (
	"context"

	"github.com/Toflex/directory_v2/graph/model"
	"github.com/Toflex/directory_v2/pkg/errors"
	"github.com/Toflex/directory_v2/pkg/log"
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
