package transaction

import (
	"context"
	"time"

	"github.com/Toflex/directory_v2/ent"
	"github.com/Toflex/directory_v2/ent/currency"
	"github.com/Toflex/directory_v2/graph/model"
	"github.com/Toflex/directory_v2/pkg/errors"
	"github.com/Toflex/directory_v2/pkg/log"
	"github.com/Toflex/directory_v2/pkg/types"
	"github.com/Toflex/directory_v2/pkg/util"
)

func (s *service) GetTransactionByID(ctx context.Context, id string, ownerID string) (*model.TransactionDetail, error) {
	logger := log.LoggerInContext(ctx)

	t, err := s.repo.GetTransactionByID(ctx, id)
	if err != nil {
		logger.WithError(err).WithField("transaction_id", id).Error("failed to query transaction")
		return nil, errors.New(errors.ErrFailed, "request failed")
	}

	if t == nil {
		return nil, errors.New(errors.ErrNotFound, "transaction not found")
	}

	if ownerID != "" && !s.belongsToOwner(t, ownerID) {
		logger.WithFields(map[string]interface{}{
			"transaction_id": id,
			"owner_id":       ownerID,
		}).Warning("transaction does not belong to the owner")
		return nil, errors.New(errors.ErrNotFound, "transaction not found")
	}

	return s.toTransactionDetail(ctx, t), nil
}

func (s *service) belongsToOwner(t *ent.Transaction, ownerID string) bool {
	if le := t.Edges.LedgerEntry; le != nil {
		if w := le.Edges.Wallet; w != nil {
			if v := w.Edges.Vault; v != nil {
				if o := v.Edges.Owner; o != nil && o.ID == ownerID {
					return true
				}
			}
		}
	}

	if nt := t.Edges.NubanTransfer; nt != nil {
		for _, b := range nt.Edges.Beneficiaries {
			if b.OwnerID == ownerID {
				return true
			}
		}
	}

	return false
}

func (s *service) getCurrencyMultiplier(ctx context.Context, currencyCode string) int64 {
	if s.db == nil {
		return 100
	}
	cur, err := s.db.Currency.Query().Where(currency.CodeEQ(currencyCode)).First(ctx)
	if err != nil || cur == nil || cur.Multiplier <= 0 {
		return 100
	}
	return cur.Multiplier
}

func (s *service) toTransactionDetail(ctx context.Context, t *ent.Transaction) *model.TransactionDetail {
	curMultiplier := s.getCurrencyMultiplier(ctx, t.Currency)

	amount := util.ToMajorUnit(types.ToMajor{
		Amount:    t.Amount,
		Precision: uint(curMultiplier),
	})

	netAmount := util.ToMajorUnit(types.ToMajor{
		Amount:    t.NetAmount,
		Precision: uint(curMultiplier),
	})

	var summary *string
	if t.Summary != "" {
		summary = &t.Summary
	}

	fees := make([]*model.TransactionFeeDetail, 0, len(t.Edges.Fees))
	for _, f := range t.Edges.Fees {
		feeMultiplier := s.getCurrencyMultiplier(ctx, f.Currency)
		feeAmount := util.ToMajorUnit(types.ToMajor{
			Amount:    f.Amount,
			Precision: uint(feeMultiplier),
		})

		fees = append(fees, &model.TransactionFeeDetail{
			ID:       f.ID,
			Amount:   feeAmount,
			Rate:     f.Rate,
			Type:     string(f.Type),
			Currency: f.Currency,
		})
	}

	var transferDetail model.TransferDetail
	if nt := t.Edges.NubanTransfer; nt != nil {
		transferAmount := util.ToMajorUnit(types.ToMajor{
			Amount:    nt.Amount,
			Precision: uint(curMultiplier),
		})

		var bankName, bankCode *string
		if nt.Edges.NubanBank != nil {
			bankName = &nt.Edges.NubanBank.Name
			bankCode = &nt.Edges.NubanBank.CbnCode
		}

		transferDetail = &model.NubanTransferDetail{
			ID:                     nt.ID,
			Amount:                 transferAmount,
			RecipientAccountName:   nt.RecipientAccountName,
			RecipientAccountNumber: nt.RecipientAccountNumber,
			BankName:               bankName,
			BankCode:               bankCode,
			SessionID:              nt.SessionID,
			ProviderReference:      nt.ProviderReference,
		}
	}

	return &model.TransactionDetail{
		ID:        t.ID,
		Reference: t.Reference,
		Type:      t.Type.String(),
		Channel:   t.Channel.String(),
		Status:    t.Status.String(),
		Currency:  t.Currency,
		Summary:   summary,
		Amount:    amount,
		NetAmount: &netAmount,
		Entry:     t.Entry.String(),
		CreatedAt: t.CreatedAt.Format(time.RFC3339),
		Fees:      fees,
		Transfer:  transferDetail,
	}
}
