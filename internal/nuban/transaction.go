package nuban

import (
	"context"
	"fmt"

	"github.com/Toflex/directory_v2/ent"
	"github.com/Toflex/directory_v2/ent/nubantransfer"
	"github.com/Toflex/directory_v2/ent/transaction"
	"github.com/Toflex/directory_v2/ent/transactionfee"
	"github.com/Toflex/directory_v2/internal/currency"
	unitofwork "github.com/Toflex/directory_v2/internal/unit_of_work"
	"github.com/Toflex/directory_v2/internal/wallet"
	"github.com/Toflex/directory_v2/pkg/errors"
	"github.com/Toflex/directory_v2/pkg/log"
	"github.com/Toflex/directory_v2/pkg/provider"
	"github.com/Toflex/directory_v2/pkg/util"
)

func (s *service) createTransferTransaction(ctx context.Context, payload *transactionPayload) (*NubanTransferResponse, error) {
	logger := log.LoggerInContext(ctx)
	tx, err := unitofwork.NewUnitOfWorkRepository(s.db).Begin(ctx)
	if err != nil {
		logger.WithError(err).Error("failed to create db transaction")
		return nil, errors.New(errors.ErrFailed, "request failed")
	}

	defer tx.Rollback()

	if payload.Reference == "" {
		payload.Reference = fmt.Sprintf("NUBAN-%s", util.GenerateCUID())
	}

	activeProvider := payload.Service.ActiveProvider.Slug

	// verify minimum and maximum amount
	if payload.Amount == 0 || !util.IsWithinRange(int64(payload.Amount), payload.Service.MinValue, payload.Service.MaxValue) {
		logger.WithField("amount", payload.Amount).Error("transfer amount is outside allowable limit")
		return nil, errors.New(errors.ErrInvalidInput, "transfer amount is outside allowable limit")
	}

	fees, err := provider.CalculateFee(ctx, payload.Service, int64(payload.Amount))
	if err != nil {
		logger.WithError(err).Error("failed to compute fee")
		return nil, err
	}

	totalFee := provider.SumServiceFees(fees)
	transactionFees := make([]*ent.TransactionFeeCreate, 0, len(fees))

	t, err := tx.Transaction.Create().
		SetAmount(int64(payload.Amount)).
		SetChannel(transaction.ChannelBANK_TRANSFER).
		SetStatus(transaction.StatusPENDING).
		SetCurrency(payload.Currency.ToString()).
		SetNillableSummary(payload.Narration).
		SetType(transaction.TypeTRANSFER).
		SetEntry(transaction.EntryDEBIT).
		SetReference(payload.Reference).
		SetProvider(activeProvider).
		SetNetAmount(int64(payload.Amount) + totalFee).
		Save(ctx)
	if err != nil {
		return nil, err
	}

	for _, fee := range fees {
		transactionFees = append(transactionFees,
			tx.TransactionFee.Create().
				SetAmount(fee.Fee).
				SetNillableRate(fee.Rate).
				SetType(transactionfee.Type(fee.FeeType)).
				SetCurrency(currency.NGN.ToString()).
				SetBaseAmount(int64(payload.Amount)).
				SetServiceIdentifier(payload.Service.Identifier).
				SetTransaction(t),
		)
	}

	if len(transactionFees) > 0 {
		_, err = tx.TransactionFee.CreateBulk(transactionFees...).Save(ctx)
		if err != nil {
			return nil, err
		}
	}

	// Nuban Transfer
	_, err = tx.NubanTransfer.Create().
		SetAmount(int64(payload.Amount)).
		SetNubanBankID(payload.Bank.ID).
		SetProvider(activeProvider).
		SetRecipientAccountName(payload.AccountName).
		SetRecipientAccountNumber(payload.AccountNumber).
		SetTransactionID(t.ID).
		AddBeneficiaryIDs(payload.BeneficiaryID).
		Save(ctx)
	if err != nil {
		return nil, err
	}

	// debit wallet
	err = s.walletService.DebitWallet(ctx, &wallet.DebitWalletPayload{
		OwnerID:               payload.OwnerId,
		WalletType:            string(wallet.TreasuryWalletType),
		DebitAmount:           t.NetAmount,
		Transaction:           t,
		ReversalTransactionID: nil,
		IsReversal:            false,
		DBTransaction:         tx,
		AffectedBalance:       wallet.AvailableBalance,
	})
	if err != nil {
		logger.WithError(err).Error("failed to debit wallet")
		return nil, err
	}

	// commit transaction
	err = tx.Commit()
	if err != nil {
		logger.WithError(err).Error("failed to commit transaction")
		return nil, errors.New(errors.ErrFailed, "request failed")
	}

	return &NubanTransferResponse{
		ID:     t.ID,
		Status: t.Status.String(),
	}, nil
}

func (s *service) updateTransferTransactionStatus(ctx context.Context, param *transactionUpdatePayload) error {
	logger := log.LoggerInContext(ctx)

	// Begin a DB Transaction
	tx, err := unitofwork.NewUnitOfWorkRepository(s.db).Begin(ctx)
	if err != nil {
		logger.WithError(err).Error("failed to create db transaction")
		return errors.New(errors.ErrFailed, "request failed")
	}

	defer tx.Rollback()

	_, err = tx.Transaction.
		Update().
		Where(transaction.IDEQ(param.ID)).
		SetNillableExternalReference(param.PartnerReference).
		SetStatus(transaction.Status(param.Status)).
		Save(ctx)

	if err != nil {
		logger.WithError(err).WithField("payload", param).Error("error when updating nuban transaction")
		return err
	}

	_, err = tx.NubanTransfer.
		Update().
		Where(nubantransfer.HasTransactionWith(transaction.IDEQ(param.ID))).
		SetNillableSessionID(param.SessionID).
		SetNillableProviderReference(param.PartnerReference).
		Save(ctx)
	if err != nil {
		logger.WithError(err).WithField("payload", param).Error("error when updating nuban transfer")
		return err
	}

	// commit transaction
	err = tx.Commit()
	if err != nil {
		logger.WithError(err).Error("failed to commit transaction")
		return errors.New(errors.ErrFailed, "request failed")
	}

	return nil
}

func (s *service) reverseTransfer(ctx context.Context, transactionID string) error {
	logger := log.LoggerInContext(ctx)

	tx, err := unitofwork.NewUnitOfWorkRepository(s.db).Begin(ctx)
	if err != nil {
		logger.WithError(err).Error("failed to create db transaction for reversal")
		return errors.New(errors.ErrFailed, "request failed")
	}

	defer tx.Rollback()

	err = s.walletService.ReverseFund(ctx, tx, transactionID)
	if err != nil {
		logger.WithError(err).WithField("transaction_id", transactionID).Error("failed to reverse transfer funds")
		return err
	}

	err = tx.Commit()
	if err != nil {
		logger.WithError(err).Error("failed to commit reversal transaction")
		return errors.New(errors.ErrFailed, "request failed")
	}

	return nil
}
