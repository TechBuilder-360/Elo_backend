package wallet

import (
	"context"
	"errors"

	"github.com/Toflex/directory_v2/ent"
)

type LedgerEntryPayload struct {
	Amount                int64
	TransactionID         string
	ReversalTransactionID *string
	Wallet                *WalletResponse
	AffectedBalance       BalanceType
	DBTransaction         *ent.Tx
}

type LedgerEntryResult struct {
	AvailableBalance int64
	HoldingBalance   int64
	LedgerBalance    int64
}

func (s *service) debitLedgerEntry(ctx context.Context, payload *LedgerEntryPayload) (*LedgerEntryResult, error) {
	var (
		result = &LedgerEntryResult{
			AvailableBalance: payload.Wallet.AvailableBalance,
			LedgerBalance:    payload.Wallet.LedgerBalance,
			HoldingBalance:   payload.Wallet.HoldingBalance,
		}
		currenctBalance int64
		previousBalance int64
	)

	// compute new balance
	switch payload.AffectedBalance {
	case AvailableBalance:
		previousBalance = payload.Wallet.AvailableBalance
		currenctBalance = previousBalance - payload.Amount
		result.AvailableBalance = currenctBalance
	case HoldingBalance:
		previousBalance = payload.Wallet.HoldingBalance
		currenctBalance = previousBalance - payload.Amount
		result.HoldingBalance = currenctBalance
	case LedgerBalance:
		previousBalance = payload.Wallet.LedgerBalance
		currenctBalance = previousBalance - payload.Amount
		result.LedgerBalance = currenctBalance
	default:
		return nil, errors.New("invalid affected balance type")
	}

	debitPayload := &LedgerEntry{
		DBTransaction:         payload.DBTransaction,
		TransactionID:         payload.TransactionID,
		ReversalTransactionID: payload.ReversalTransactionID,
		IsReversal:            payload.ReversalTransactionID != nil,
		CurrentBalance:        currenctBalance,
		PreviousBalance:       previousBalance,
		AffectedBalance:       payload.AffectedBalance,
		Debit:                 payload.Amount,
		Wallet:                payload.Wallet,
	}

	err := s.repo.WithTransaction(payload.DBTransaction).AddLedgerRecord(ctx, debitPayload)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (s *service) creditLedgerEntry(ctx context.Context, payload *LedgerEntryPayload) (*LedgerEntryResult, error) {
	var (
		result = &LedgerEntryResult{
			AvailableBalance: payload.Wallet.AvailableBalance,
			LedgerBalance:    payload.Wallet.LedgerBalance,
			HoldingBalance:   payload.Wallet.HoldingBalance,
		}
		currenctBalance int64
		previousBalance int64
	)

	// compute new balance
	switch payload.AffectedBalance {
	case AvailableBalance:
		previousBalance = payload.Wallet.AvailableBalance
		currenctBalance = previousBalance + payload.Amount
		result.AvailableBalance = currenctBalance
	case HoldingBalance:
		previousBalance = payload.Wallet.HoldingBalance
		currenctBalance = previousBalance + payload.Amount
		result.HoldingBalance = currenctBalance
	case LedgerBalance:
		previousBalance = payload.Wallet.LedgerBalance
		currenctBalance = previousBalance + payload.Amount
		result.LedgerBalance = currenctBalance
	default:
		return nil, errors.New("invalid affected balance type")
	}

	creditPayload := &LedgerEntry{
		DBTransaction:         payload.DBTransaction,
		TransactionID:         payload.TransactionID,
		ReversalTransactionID: payload.ReversalTransactionID,
		IsReversal:            payload.ReversalTransactionID != nil,
		CurrentBalance:        currenctBalance,
		PreviousBalance:       previousBalance,
		AffectedBalance:       payload.AffectedBalance,
		Credit:                payload.Amount,
		Wallet:                payload.Wallet,
	}

	err := s.repo.WithTransaction(payload.DBTransaction).AddLedgerRecord(ctx, creditPayload)
	if err != nil {
		return nil, err
	}

	return result, nil
}
