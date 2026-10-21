package wallet

import (
	"context"

	"github.com/Toflex/directory_v2/ent/ledger"
)

func (r *repository) AddLedgerRecord(ctx context.Context, payload *LedgerEntry) error {
	_, err := payload.DBTransaction.Ledger.Create().
		SetTransactionID(payload.TransactionID).
		SetNillableReversalTransactionID(payload.ReversalTransactionID).
		SetCredit(payload.Credit).
		SetDebit(payload.Debit).
		SetAffectedBalance(ledger.AffectedBalance(payload.AffectedBalance)).
		SetIsReversal(payload.IsReversal).
		SetCurrentBalance(payload.CurrentBalance).
		SetPreviousBalance(payload.PreviousBalance).
		SetWalletID(payload.Wallet.ID).
		Save(ctx)

	return err
}
