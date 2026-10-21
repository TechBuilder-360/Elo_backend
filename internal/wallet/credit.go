package wallet

import (
	"context"

	"github.com/Toflex/directory_v2/pkg/errors"
	"github.com/Toflex/directory_v2/pkg/log"
	"github.com/Toflex/directory_v2/pkg/types"
)

func (s *service) CreditWallet(ctx context.Context, payload *CreditWalletPayload) error {
	logger := log.LoggerInContext(ctx)

	if payload.AffectedBalance == "" {
		payload.AffectedBalance = AvailableBalance
	}

	w, err := s.ValidateWalletForCredit(ctx, &ValidateWalletPayload{
		DBTransaction: payload.DBTransaction,
		OwnerID:       payload.OwnerID,
		WalletType:    payload.WalletType,
		Currency:      types.CurrencyCode(payload.Transaction.Currency),
	})
	if err != nil {
		return err
	}

	// create ledger transaction
	ledger, err := s.creditLedgerEntry(ctx, &LedgerEntryPayload{
		Amount:                payload.CreditAmount,
		TransactionID:         payload.Transaction.ID,
		ReversalTransactionID: payload.ReversalTransactionID,
		Wallet:                w,
		AffectedBalance:       payload.AffectedBalance,
		DBTransaction:         payload.DBTransaction,
	})
	if err != nil {
		logger.WithError(err).Error("failed to create ledger record")
		return errors.New(errors.ErrFailed, "error when crediting wallet")
	}

	// modify wallet balance
	repo := s.repo
	if payload.DBTransaction != nil {
		repo = s.repo.WithTransaction(payload.DBTransaction)
	}

	err = repo.UpdateWalletBalance(ctx, w, ledger)
	if err != nil {
		logger.WithError(err).Error("failed to update wallet balance")
		return errors.New(errors.ErrFailed, "failed to credit wallet")
	}

	return nil
}
