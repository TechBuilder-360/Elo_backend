package nuban

import (
	"context"
	"strings"

	"github.com/Toflex/directory_v2/ent"
	"github.com/Toflex/directory_v2/pkg/errors"
	"github.com/Toflex/directory_v2/pkg/log"
	"github.com/Toflex/directory_v2/pkg/provider"
	"github.com/Toflex/directory_v2/pkg/util"
)

type bankPayload struct {
	Slug     string
	Provider string
}

func (a *activityWrapper) GetBankCode(ctx context.Context, payload *bankPayload) (*bankResponse, error) {
	bank, err := a.svc.GetBankBySlug(ctx, payload.Slug)
	if err != nil {
		return nil, err
	}

	if bank == nil {
		return nil, errors.New(errors.ErrFailed, "bank not found")
	}

	var providerCode = bank.CbnCode
	for _, p := range bank.Edges.ProviderBanks {
		if strings.EqualFold(payload.Provider, p.Provider) {
			providerCode = p.ProviderCode
		}
	}

	return &bankResponse{
		ID:           bank.ID,
		ProviderCode: providerCode,
	}, nil
}

func (a *activityWrapper) GetTransferService(ctx context.Context) (provider.ServiceLocator, error) {
	logger := log.LoggerInContext(ctx)
	serv, err := a.svc.getServiceProvider(ctx, NubanBankTransferIdentifier)
	if err != nil {
		logger.WithError(err).WithField("identifier", NubanBankTransferIdentifier).Error("failed to get service identifier")
		if ent.IsNotFound(err) {
			return provider.ServiceLocator{}, errors.New(errors.ErrFailed, "service not found")
		}
		return provider.ServiceLocator{}, errors.New(errors.ErrFailed, "failed to get transfer service")
	}

	// // Account Name resolution
	_, ok := provider.GetImpl[NUBANTransferProvider](serv.ActiveProvider.Slug)
	if !ok {
		logger.WithField("provider", serv.ActiveProvider.Slug).Error("failed to get Implementation")
		return provider.ServiceLocator{}, errors.New(errors.ErrFailed, "provider not supported")
	}

	return serv, nil
}

func (a *activityWrapper) AccountResolutionActivity(ctx context.Context, payload ResolveAccountNamePayload) (*ResolveAccountNameResult, error) {
	return a.svc.ResolveAccountName(ctx, payload)
}

func (a *activityWrapper) BeneficiaryChecklist(ctx context.Context, payload *nubanBeneficiaryPayload) (*nubanBeneficiaryResult, error) {
	return a.svc.beneficiaryCheck(ctx, payload)
}

func (a *activityWrapper) createTransferTransaction(ctx context.Context, payload *transactionPayload) (*NubanTransferResponse, error) {
	return a.svc.createTransferTransaction(ctx, payload)
}

func (a *activityWrapper) transferToProvider(ctx context.Context, param *ProcessTransferPayload) (*NubanTransferResult, error) {
	p, ok := provider.GetImpl[NUBANTransferProvider](param.Provider)
	if !ok {
		return nil, errors.New(errors.ErrFailed, "provider does not support nuban transfer")
	}

	return p.Transfer(ctx, &NubanTransferPayload{
		Reference:     param.TransactionID,
		Narration:     util.AddressToString(param.Narration),
		Amount:        param.Amount,
		AccountName:   param.Beneficiary.AccountName,
		AccountNumber: param.Beneficiary.AccountNumber,
		BankCode:      param.Bank.ProviderCode,
	})
}

func (a *activityWrapper) transferStatus(ctx context.Context, reference, servProvider string) (*NubanTransferResult, error) {
	p, ok := provider.GetImpl[NUBANTransferProvider](servProvider)
	if !ok {
		return nil, errors.New(errors.ErrFailed, "provider does not support nuban transfer")
	}

	return p.GetTransactionStatus(ctx, reference)
}

func (a *activityWrapper) updateTransaction(ctx context.Context, param *transactionUpdatePayload) error {
	return a.svc.updateTransferTransactionStatus(ctx, param)
}

func (a *activityWrapper) reverseTransfer(ctx context.Context, transactionID string) error {
	return a.svc.reverseTransfer(ctx, transactionID)
}
