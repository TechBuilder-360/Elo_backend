package mock

import (
	"context"
	"math/rand"

	"github.com/Toflex/directory_v2/internal/nuban"
	"github.com/Toflex/directory_v2/pkg/constant"
	"github.com/Toflex/directory_v2/pkg/types"
	"github.com/Toflex/directory_v2/pkg/util"
)

// ResolveAccount implements [nuban.NUBANTransferProvider].
func (m *mock) ResolveAccount(ctx context.Context, payload nuban.ResolveAccountPayload) (*nuban.ResolveAccountResult, error) {
	return &nuban.ResolveAccountResult{
		AccountName:   "Test Account",
		AccountNumber: payload.AccountNumber,
		BankCode:      payload.BankCode,
	}, nil
}

func randomStatus() types.TransactionStatus {
	status := []types.TransactionStatus{constant.TransactionFailedStatus, constant.TransactionPendingStatus, constant.TransactionSuccessStatus}
	return status[rand.Intn(len(status))]
}

// Transfer implements [nuban.NUBANTransferProvider].
func (m *mock) Transfer(ctx context.Context, payload *nuban.NubanTransferPayload) (*nuban.NubanTransferResult, error) {

	finalStatus := randomStatus()
	if finalStatus == constant.TransactionFailedStatus {
		return &nuban.NubanTransferResult{
			Status:    finalStatus,
			Reference: payload.Reference,
		}, nil
	}

	return &nuban.NubanTransferResult{
		Status:    finalStatus,
		PartnerID: util.GenerateCUID(),
		Reference: payload.Reference,
	}, nil
}

// GetTransactionStatus implements [nuban.NUBANTransferProvider].
func (m *mock) GetTransactionStatus(ctx context.Context, id string) (*nuban.NubanTransferResult, error) {
	return &nuban.NubanTransferResult{
		Status:    randomStatus(),
		PartnerID: util.GenerateCUID(),
		Reference: id,
	}, nil
}

var _ nuban.NUBANTransferProvider = (*mock)(nil)
