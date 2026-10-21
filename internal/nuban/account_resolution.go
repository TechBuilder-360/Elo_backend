package nuban

import (
	"context"

	"github.com/Toflex/directory_v2/pkg/errors"
	"github.com/Toflex/directory_v2/pkg/log"
	"github.com/Toflex/directory_v2/pkg/provider"
)

const (
	NubanBankTransferIdentifier string = "ng-bank-transfer"
)

func (s *service) ResolveAccountName(ctx context.Context, payload ResolveAccountNamePayload) (*ResolveAccountNameResult, error) {
	logger := log.LoggerInContext(ctx)

	p, err := s.getServiceProvider(ctx, NubanBankTransferIdentifier)
	if err != nil {
		logger.WithError(err).WithField("identifier", NubanBankTransferIdentifier).Error("failed to get service identifier")
		return nil, errors.New(errors.ErrFailed, "request failed")
	}

	impl, ok := provider.GetProviderImpl(p.ActiveProvider.Slug)
	if !ok {
		logger.WithField("provider", p.ActiveProvider.Slug).Error("provider impl not found")
		return nil, errors.New(errors.ErrFailed, "request failed")
	}

	account, ok := provider.ConformsTo[NUBANTransferProvider](impl)
	if !ok {
		logger.WithField("provider", p.ActiveProvider.Slug).Error("provider does not impl NUBANTransferProvider")
		return nil, errors.New(errors.ErrFailed, "request failed")
	}

	result, err := account.ResolveAccount(ctx, ResolveAccountPayload{
		AccountNumber: payload.AccountNumber,
		BankCode:      payload.BankCode,
	})

	if err != nil {
		logger.WithError(err).Error("account name resolution failed")
		return nil, errors.New(errors.ErrFailed, "account name could not be resolved")
	}

	return &ResolveAccountNameResult{
		AccountNumber: result.AccountNumber,
		AccountName:   result.AccountName,
		BankCode:      result.BankCode,
	}, nil
}

func (s *service) getServiceProvider(ctx context.Context, identifier string) (provider.ServiceLocator, error) {
	return provider.NewService().GetServiceByIdentifier(ctx, identifier)
}
