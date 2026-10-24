package transaction

import (
	"context"

	"github.com/Toflex/directory_v2/ent"
	"github.com/Toflex/directory_v2/ent/transaction"
)

func (r *repository) GetTransactionByID(ctx context.Context, id string) (*ent.Transaction, error) {
	t, err := r.db.Transaction.
		Query().
		Where(
			transaction.Or(
				transaction.IDEQ(id),
				transaction.ReferenceEQ(id),
			),
		).
		WithFees().
		WithNubanTransfer(func(nq *ent.NubanTransferQuery) {
			nq.WithNubanBank().WithBeneficiaries()
		}).
		WithLedgerEntry(func(lq *ent.LedgerQuery) {
			lq.WithWallet(func(wq *ent.WalletQuery) {
				wq.WithVault(func(vq *ent.VaultQuery) {
					vq.WithOwner()
				})
			})
		}).
		First(ctx)

	if ent.IsNotFound(err) {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	return t, nil
}
