package wallet

import (
	"context"

	"github.com/Toflex/directory_v2/ent"
	"github.com/Toflex/directory_v2/ent/transaction"
	"github.com/Toflex/directory_v2/pkg/constant"
	"github.com/Toflex/directory_v2/pkg/errors"
	"github.com/Toflex/directory_v2/pkg/log"
	"github.com/Toflex/directory_v2/pkg/types"
	"github.com/Toflex/directory_v2/pkg/util"
)

func (s *service) DebitWallet(ctx context.Context, payload *DebitWalletPayload) error {
	logger := log.LoggerInContext(ctx)

	if payload.AffectedBalance == "" {
		payload.AffectedBalance = AvailableBalance
	}

	w, err := s.ValidateWalletForDebit(ctx, &ValidateWalletPayload{
		OwnerID:       payload.OwnerID,
		DBTransaction: payload.DBTransaction,
		WalletType:    payload.WalletType,
		Currency:      types.CurrencyCode(payload.Transaction.Currency),
	})
	if err != nil {
		return err
	}

	// check that there is sufficient fund
	if !sufficientBalance(w, payload.AffectedBalance, payload.DebitAmount) {
		logger.WithFields(map[string]interface{}{
			"available_balance": w.AvailableBalance,
			"holding_balance":   w.HoldingBalance,
			"ledger_balance":    w.LedgerBalance,
			"debit_amount":      payload.DebitAmount,
		}).Error("insufficient wallet balance")
		return errors.New(errors.ErrFailed, "insufficient fund")
	}

	// create ledger transaction
	ledger, err := s.debitLedgerEntry(ctx, &LedgerEntryPayload{
		Amount:                payload.DebitAmount,
		TransactionID:         payload.Transaction.ID,
		ReversalTransactionID: payload.ReversalTransactionID,
		Wallet:                w,
		AffectedBalance:       payload.AffectedBalance,
		DBTransaction:         payload.DBTransaction,
	})
	if err != nil {
		logger.WithError(err).Error("failed to create ledger record")
		return errors.New(errors.ErrFailed, "error when debitting wallet")
	}

	// modify wallet balance
	repo := s.repo
	if payload.DBTransaction != nil {
		repo = s.repo.WithTransaction(payload.DBTransaction)
	}

	err = repo.UpdateWalletBalance(ctx, w, ledger)
	if err != nil {
		logger.WithError(err).Error("failed to update wallet balance")
		return errors.New(errors.ErrFailed, "failed to debit wallet")
	}

	return nil
}

func sufficientBalance(wallet *WalletResponse, balanceType BalanceType, amount int64) bool {
	switch balanceType {
	case AvailableBalance:
		return wallet.AvailableBalance >= amount
	case HoldingBalance:
		return wallet.HoldingBalance >= amount
	case LedgerBalance:
		return wallet.LedgerBalance >= amount
	default:
		return false
	}
}

// ReverseFund implements [IService].
func (s *service) ReverseFund(ctx context.Context, tx *ent.Tx, transactionID string) error {
	logger := log.LoggerInContext(ctx)

	repo := s.repo
	if tx != nil {
		repo = s.repo.WithTransaction(tx)
	}

	t, err := repo.GetLedgerbyTransactionID(ctx, transactionID)
	if err != nil {
		logger.WithError(err).Error("failed to fetch transaction ledger")
		return errors.New(errors.ErrFailed, "request failed")
	}

	if t == nil {
		return errors.New(errors.ErrFailed, "transaction ledger record not found")
	}

	if t.Status != transaction.StatusFAILED {
		logger.WithField("transaction_status", t.Status).Error("transaction cannot be reversed")
		return errors.New(errors.ErrFailed, "transaction cannot be reversed")
	}

	if t.Entry != transaction.EntryDEBIT {
		logger.WithField("transaction_status", t.Status).Error("transaction cannot be reversed")
		return errors.New(errors.ErrFailed, "transaction cannot be reversed")
	}

	ledgerEntry := t.Edges.LedgerEntry
	if ledgerEntry == nil {
		return errors.New(errors.ErrFailed, "ledger entry not found")
	}

	wallet := ledgerEntry.Edges.Wallet
	if wallet == nil {
		return errors.New(errors.ErrFailed, "ledger wallet not found")
	}

	ownerID := ""
	walletType := string(TreasuryWalletType)
	if wallet.Edges.Vault != nil {
		walletType = string(wallet.Edges.Vault.Type)
		if wallet.Edges.Vault.Edges.Owner != nil {
			ownerID = wallet.Edges.Vault.Edges.Owner.ID
		}
	}

	return s.reverseTransaction(ctx, reverseTransactionPayload{
		DBTransaction: tx,
		Transaction:   t,
		WalletID:      wallet.ID,
		OwnerID:       ownerID,
		WalletType:    walletType,
	})
}

type reverseTransactionPayload struct {
	DBTransaction *ent.Tx
	Transaction   *ent.Transaction
	WalletID      string
	OwnerID       string
	WalletType    string
}

func (s *service) reverseTransaction(ctx context.Context, param reverseTransactionPayload) error {
	reference := util.GenerateReference("REVERSAL")

	newTx, err := param.DBTransaction.Transaction.Create().
		SetAmount(param.Transaction.NetAmount).
		SetNetAmount(param.Transaction.NetAmount).
		SetChannel(transaction.ChannelINTERNAL).
		SetCurrency(param.Transaction.Currency).
		SetEntry(transaction.EntryCREDIT).
		SetType(transaction.TypeREVERSAL).
		SetStatus(transaction.StatusSUCCESS).
		SetReference(reference).
		SetProvider(string(constant.Elo)).
		Save(ctx)
	if err != nil {
		return err
	}

	return s.CreditWallet(ctx, &CreditWalletPayload{
		OwnerID:               param.OwnerID,
		WalletType:            param.WalletType,
		CreditAmount:          newTx.Amount,
		Transaction:           newTx,
		ReversalTransactionID: &param.Transaction.ID,
		IsReversal:            true,
		DBTransaction:         param.DBTransaction,
		AffectedBalance:       AvailableBalance,
	})
}
