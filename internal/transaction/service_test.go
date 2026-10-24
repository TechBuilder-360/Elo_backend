package transaction

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/Toflex/directory_v2/ent"
	enttransaction "github.com/Toflex/directory_v2/ent/transaction"
	"github.com/Toflex/directory_v2/ent/transactionfee"
	"github.com/Toflex/directory_v2/graph/model"
	"github.com/stretchr/testify/assert"
)

type mockRepo struct {
	getTransactionByIDFn func(ctx context.Context, id string) (*ent.Transaction, error)
}

func (m *mockRepo) GetTransactionByID(ctx context.Context, id string) (*ent.Transaction, error) {
	if m.getTransactionByIDFn != nil {
		return m.getTransactionByIDFn(ctx, id)
	}
	return nil, nil
}

func TestGetTransactionByID_Success(t *testing.T) {
	rate := 1.5
	externalRef := "ext-ref-123"
	now := time.Now()

	mockT := &ent.Transaction{
		ID:                "tx-123",
		Reference:         "REF-001",
		ExternalReference: &externalRef,
		Type:              enttransaction.TypeTRANSFER,
		Channel:           enttransaction.ChannelBANK_TRANSFER,
		Status:            enttransaction.StatusSUCCESS,
		Currency:          "NGN",
		Summary:           "Transfer for supplies",
		Provider:          "mock-bank",
		Amount:            10000, // 100.00 in major units with 100 multiplier
		NetAmount:         9500,  // 95.00 in major units
		Entry:             enttransaction.EntryDEBIT,
		CreatedAt:         now,
		Edges: ent.TransactionEdges{
			Fees: []*ent.TransactionFee{
				{
					ID:       "fee-1",
					Amount:   500, // 5.00
					Rate:     &rate,
					Type:     transactionfee.TypeSERVICE_FEE,
					Currency: "NGN",
				},
			},
			NubanTransfer: &ent.NubanTransfer{
				ID:                     "nt-1",
				Amount:                 10000,
				RecipientAccountName:   "Jane Doe",
				RecipientAccountNumber: "0123456789",
				Edges: ent.NubanTransferEdges{
					NubanBank: &ent.NubanBank{
						Name:    "First Bank",
						CbnCode: "011",
					},
				},
			},
			LedgerEntry: &ent.Ledger{
				Edges: ent.LedgerEdges{
					Wallet: &ent.Wallet{
						Edges: ent.WalletEdges{
							Vault: &ent.Vault{
								Edges: ent.VaultEdges{
									Owner: &ent.LedgerOwner{
										ID: "owner-123",
									},
								},
							},
						},
					},
				},
			},
		},
	}

	repo := &mockRepo{
		getTransactionByIDFn: func(ctx context.Context, id string) (*ent.Transaction, error) {
			assert.Equal(t, "tx-123", id)
			return mockT, nil
		},
	}

	svc := &service{
		repo: repo,
	}

	res, err := svc.GetTransactionByID(context.Background(), "tx-123", "owner-123")
	assert.NoError(t, err)
	assert.NotNil(t, res)

	assert.Equal(t, "tx-123", res.ID)
	assert.Equal(t, "REF-001", res.Reference)
	assert.Equal(t, 100.0, res.Amount)
	assert.Equal(t, 95.0, *res.NetAmount)
	assert.Equal(t, "TRANSFER", res.Type)
	assert.Equal(t, "BANK_TRANSFER", res.Channel)
	assert.Equal(t, "SUCCESS", res.Status)
	assert.Equal(t, "DEBIT", res.Entry)
	assert.Equal(t, "Transfer for supplies", *res.Summary)

	// Fees
	assert.Len(t, res.Fees, 1)
	assert.Equal(t, "fee-1", res.Fees[0].ID)
	assert.Equal(t, 5.0, res.Fees[0].Amount)
	assert.Equal(t, &rate, res.Fees[0].Rate)
	assert.Equal(t, "SERVICE_FEE", res.Fees[0].Type)

	// Nuban transfer
	assert.NotNil(t, res.Transfer)
	transfer, ok := res.Transfer.(*model.NubanTransferDetail)
	assert.True(t, ok)
	assert.Equal(t, "nt-1", transfer.ID)
	assert.Equal(t, 100.0, transfer.Amount)
	assert.Equal(t, "Jane Doe", transfer.RecipientAccountName)
	assert.Equal(t, "0123456789", transfer.RecipientAccountNumber)
	assert.Equal(t, "First Bank", *transfer.BankName)
	assert.Equal(t, "011", *transfer.BankCode)
}

func TestGetTransactionByID_NotFound(t *testing.T) {
	repo := &mockRepo{
		getTransactionByIDFn: func(ctx context.Context, id string) (*ent.Transaction, error) {
			return nil, nil
		},
	}

	svc := &service{
		repo: repo,
	}

	res, err := svc.GetTransactionByID(context.Background(), "non-existent", "owner-123")
	assert.Error(t, err)
	assert.Nil(t, res)
	assert.Contains(t, err.Error(), "transaction not found")
}

func TestGetTransactionByID_DBError(t *testing.T) {
	repo := &mockRepo{
		getTransactionByIDFn: func(ctx context.Context, id string) (*ent.Transaction, error) {
			return nil, fmt.Errorf("pq: connection refused")
		},
	}

	svc := &service{
		repo: repo,
	}

	res, err := svc.GetTransactionByID(context.Background(), "tx-123", "owner-123")
	assert.Error(t, err)
	assert.Nil(t, res)
	// Must not leak internal database error
	assert.NotContains(t, err.Error(), "pq:")
	assert.Equal(t, "request failed", err.Error())
}

func TestGetTransactionByID_OwnerAuthorization(t *testing.T) {
	mockT := &ent.Transaction{
		ID:        "tx-123",
		Reference: "REF-001",
		Edges: ent.TransactionEdges{
			LedgerEntry: &ent.Ledger{
				Edges: ent.LedgerEdges{
					Wallet: &ent.Wallet{
						Edges: ent.WalletEdges{
							Vault: &ent.Vault{
								Edges: ent.VaultEdges{
									Owner: &ent.LedgerOwner{
										ID: "owner-ownerA",
									},
								},
							},
						},
					},
				},
			},
		},
	}

	repo := &mockRepo{
		getTransactionByIDFn: func(ctx context.Context, id string) (*ent.Transaction, error) {
			return mockT, nil
		},
	}

	svc := &service{
		repo: repo,
	}

	// Case 1: Wrong owner should return transaction not found
	res, err := svc.GetTransactionByID(context.Background(), "tx-123", "owner-ownerB")
	assert.Error(t, err)
	assert.Nil(t, res)
	assert.Contains(t, err.Error(), "transaction not found")

	// Case 2: Matching owner succeeds
	res, err = svc.GetTransactionByID(context.Background(), "tx-123", "owner-ownerA")
	assert.NoError(t, err)
	assert.NotNil(t, res)
	assert.Equal(t, "tx-123", res.ID)

	// Case 3: Empty ownerID bypasses owner check (for internal/admin callers)
	res, err = svc.GetTransactionByID(context.Background(), "tx-123", "")
	assert.NoError(t, err)
	assert.NotNil(t, res)
	assert.Equal(t, "tx-123", res.ID)
}

func TestGetTransactionByID_BeneficiaryOwner(t *testing.T) {
	mockT := &ent.Transaction{
		ID:        "tx-456",
		Reference: "REF-002",
		Edges: ent.TransactionEdges{
			NubanTransfer: &ent.NubanTransfer{
				ID: "nt-2",
				Edges: ent.NubanTransferEdges{
					Beneficiaries: []*ent.NubanBeneficiary{
						{
							ID:      "ben-1",
							OwnerID: "owner-beneficiary",
						},
					},
				},
			},
		},
	}

	repo := &mockRepo{
		getTransactionByIDFn: func(ctx context.Context, id string) (*ent.Transaction, error) {
			return mockT, nil
		},
	}

	svc := &service{
		repo: repo,
	}

	// Matching beneficiary owner
	res, err := svc.GetTransactionByID(context.Background(), "tx-456", "owner-beneficiary")
	assert.NoError(t, err)
	assert.NotNil(t, res)
	assert.Equal(t, "tx-456", res.ID)

	// Mismatched owner
	res, err = svc.GetTransactionByID(context.Background(), "tx-456", "different-owner")
	assert.Error(t, err)
	assert.Nil(t, res)
	assert.Contains(t, err.Error(), "transaction not found")
}
