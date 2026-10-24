package wallet

import (
	"context"
	"testing"

	"github.com/Toflex/directory_v2/ent"
	"github.com/Toflex/directory_v2/ent/ledgerowner"
	entvault "github.com/Toflex/directory_v2/ent/vault"
	"github.com/Toflex/directory_v2/graph/model"
	"github.com/Toflex/directory_v2/internal/currency"
	"github.com/Toflex/directory_v2/internal/vault"
	"github.com/Toflex/directory_v2/pkg/types"
)

type mockRepo struct {
	getOwnerWalletOverviewFn   func(ctx context.Context, ownerID string, walletType WalletType) (*OwnerVaultOverview, error)
	getWalletWithAtomicLockFn func(ctx context.Context, tx *ent.Tx, ownerID string, walletType string, currencyCode types.CurrencyCode) (*WalletResponse, error)
	getWalletWithCurrencyFn   func(ctx context.Context, ownerID, walletType string, currencyCode types.CurrencyCode) (*WalletResponse, error)
	createVaultFn             func(ctx context.Context, ownerID string, vaultType WalletType) (*ent.Vault, error)
	createFn                  func(ctx context.Context, payload *createWallet) (*WalletResponse, error)
	addLedgerRecordFn         func(ctx context.Context, payload *LedgerEntry) error
	updateWalletBalanceFn     func(ctx context.Context, wallet *WalletResponse, payload *LedgerEntryResult) error
	getLedgerbyTransactionIDFn func(ctx context.Context, transactionID string) (*ent.Transaction, error)
}

func (m *mockRepo) GetWallets(ctx context.Context, ownerID string, walletType WalletType, filter *WalletFilter) ([]WalletResponse, error) {
	return nil, nil
}
func (m *mockRepo) GetWallet(ctx context.Context, walletType, ownerID string) (*WalletResponse, error) {
	return nil, nil
}
func (m *mockRepo) GetWalletByOwnerID(ctx context.Context, walletID, ownerID string) (*WalletResponse, error) {
	return nil, nil
}
func (m *mockRepo) GetWalletWithCurrency(ctx context.Context, ownerID, walletType string, currencyCode types.CurrencyCode) (*WalletResponse, error) {
	if m.getWalletWithCurrencyFn != nil {
		return m.getWalletWithCurrencyFn(ctx, ownerID, walletType, currencyCode)
	}
	return nil, nil
}
func (m *mockRepo) GetOwnerWalletOverview(ctx context.Context, ownerID string, walletType WalletType) (*OwnerVaultOverview, error) {
	if m.getOwnerWalletOverviewFn != nil {
		return m.getOwnerWalletOverviewFn(ctx, ownerID, walletType)
	}
	return nil, nil
}
func (m *mockRepo) Create(ctx context.Context, payload *createWallet) (*WalletResponse, error) {
	if m.createFn != nil {
		return m.createFn(ctx, payload)
	}
	return nil, nil
}
func (m *mockRepo) CreateVault(ctx context.Context, ownerID string, vaultType WalletType) (*ent.Vault, error) {
	if m.createVaultFn != nil {
		return m.createVaultFn(ctx, ownerID, vaultType)
	}
	return nil, nil
}
func (m *mockRepo) GetWalletWithAtomicLock(ctx context.Context, tx *ent.Tx, ownerID string, walletType string, currencyCode types.CurrencyCode) (*WalletResponse, error) {
	if m.getWalletWithAtomicLockFn != nil {
		return m.getWalletWithAtomicLockFn(ctx, tx, ownerID, walletType, currencyCode)
	}
	return nil, nil
}
func (m *mockRepo) WithTransaction(tx *ent.Tx) IRepository {
	return m
}
func (m *mockRepo) AddLedgerRecord(ctx context.Context, payload *LedgerEntry) error {
	if m.addLedgerRecordFn != nil {
		return m.addLedgerRecordFn(ctx, payload)
	}
	return nil
}
func (m *mockRepo) UpdateWalletBalance(ctx context.Context, wallet *WalletResponse, payload *LedgerEntryResult) error {
	if m.updateWalletBalanceFn != nil {
		return m.updateWalletBalanceFn(ctx, wallet, payload)
	}
	return nil
}
func (m *mockRepo) GetLedgerbyTransactionID(ctx context.Context, transactionID string) (*ent.Transaction, error) {
	if m.getLedgerbyTransactionIDFn != nil {
		return m.getLedgerbyTransactionIDFn(ctx, transactionID)
	}
	return nil, nil
}

type mockVaultService struct {
	createVaultWithTypeFn func(ctx context.Context, ownerID, vaultType string) (*vault.Vault, error)
	getOwnerVaultFn       func(ctx context.Context, ownerID, vaultType string) (*vault.Vault, error)
}

func (m *mockVaultService) CreateVault(ctx context.Context, ownerID string) (*vault.Vault, error) {
	return nil, nil
}
func (m *mockVaultService) CreateVaultWithType(ctx context.Context, ownerID, vaultType string) (*vault.Vault, error) {
	if m.createVaultWithTypeFn != nil {
		return m.createVaultWithTypeFn(ctx, ownerID, vaultType)
	}
	return nil, nil
}
func (m *mockVaultService) GetVault(ctx context.Context, ownerID string) (*vault.Vault, error) {
	return nil, nil
}
func (m *mockVaultService) GetOwnerVault(ctx context.Context, ownerID, vaultType string) (*vault.Vault, error) {
	if m.getOwnerVaultFn != nil {
		return m.getOwnerVaultFn(ctx, ownerID, vaultType)
	}
	return nil, nil
}

type mockCurrencyService struct {
	getCurrencyByCodeFn func(ctx context.Context, code types.CurrencyCode) (*currency.Currency, error)
}

func (m *mockCurrencyService) GetCurrencies(ctx context.Context, filter *model.CurrencyFilter) ([]*model.Currency, error) {
	return nil, nil
}
func (m *mockCurrencyService) GetCurrencyByCode(ctx context.Context, code types.CurrencyCode) (*currency.Currency, error) {
	if m.getCurrencyByCodeFn != nil {
		return m.getCurrencyByCodeFn(ctx, code)
	}
	return nil, nil
}

func TestDebitWallet_Conditions(t *testing.T) {
	t.Run("fails when owner is not active", func(t *testing.T) {
		repo := &mockRepo{
			getOwnerWalletOverviewFn: func(ctx context.Context, ownerID string, walletType WalletType) (*OwnerVaultOverview, error) {
				return &OwnerVaultOverview{
					LedgerOwnerStatus: ledgerowner.StatusSUSPENDED.String(),
					VaultStatus:       entvault.StatusACTIVE.String(),
				}, nil
			},
		}

		svc := &service{repo: repo}
		_, err := svc.ValidateWalletForDebit(context.Background(), &ValidateWalletPayload{
			OwnerID:    "owner-1",
			WalletType: "TREASURY",
			Currency:   currency.NGN,
		})

		if err == nil {
			t.Fatal("expected error when owner is not active, got nil")
		}
		if err.Error() != "ledger owner is not active" {
			t.Fatalf("expected 'ledger owner is not active', got '%v'", err.Error())
		}
	})

	t.Run("fails when vault is not active", func(t *testing.T) {
		repo := &mockRepo{
			getOwnerWalletOverviewFn: func(ctx context.Context, ownerID string, walletType WalletType) (*OwnerVaultOverview, error) {
				return &OwnerVaultOverview{
					LedgerOwnerStatus: ledgerowner.StatusACTIVE.String(),
					VaultStatus:       entvault.StatusSUSPENDED.String(),
				}, nil
			},
		}

		svc := &service{repo: repo}
		_, err := svc.ValidateWalletForDebit(context.Background(), &ValidateWalletPayload{
			OwnerID:    "owner-1",
			WalletType: "TREASURY",
			Currency:   currency.NGN,
		})

		if err == nil {
			t.Fatal("expected error when vault is not active, got nil")
		}
		if err.Error() != "vault is not active" {
			t.Fatalf("expected 'vault is not active', got '%v'", err.Error())
		}
	})

	t.Run("fails when wallet is not active", func(t *testing.T) {
		repo := &mockRepo{
			getOwnerWalletOverviewFn: func(ctx context.Context, ownerID string, walletType WalletType) (*OwnerVaultOverview, error) {
				return &OwnerVaultOverview{
					LedgerOwnerStatus: ledgerowner.StatusACTIVE.String(),
					VaultStatus:       entvault.StatusACTIVE.String(),
				}, nil
			},
			getWalletWithAtomicLockFn: func(ctx context.Context, tx *ent.Tx, ownerID string, walletType string, currencyCode types.CurrencyCode) (*WalletResponse, error) {
				return &WalletResponse{
					ID:     "wallet-1",
					Active: false,
				}, nil
			},
		}

		svc := &service{repo: repo}
		_, err := svc.ValidateWalletForDebit(context.Background(), &ValidateWalletPayload{
			OwnerID:    "owner-1",
			WalletType: "TREASURY",
			Currency:   currency.NGN,
		})

		if err == nil {
			t.Fatal("expected error when wallet is not active, got nil")
		}
		if err.Error() != "wallet is not active" {
			t.Fatalf("expected 'wallet is not active', got '%v'", err.Error())
		}
	})

	t.Run("succeeds when all conditions are active", func(t *testing.T) {
		repo := &mockRepo{
			getOwnerWalletOverviewFn: func(ctx context.Context, ownerID string, walletType WalletType) (*OwnerVaultOverview, error) {
				return &OwnerVaultOverview{
					LedgerOwnerStatus: ledgerowner.StatusACTIVE.String(),
					VaultStatus:       entvault.StatusACTIVE.String(),
				}, nil
			},
			getWalletWithAtomicLockFn: func(ctx context.Context, tx *ent.Tx, ownerID string, walletType string, currencyCode types.CurrencyCode) (*WalletResponse, error) {
				return &WalletResponse{
					ID:               "wallet-1",
					Active:           true,
					AvailableBalance: 1000,
				}, nil
			},
		}

		svc := &service{repo: repo}
		w, err := svc.ValidateWalletForDebit(context.Background(), &ValidateWalletPayload{
			OwnerID:    "owner-1",
			WalletType: "TREASURY",
			Currency:   currency.NGN,
		})

		if err != nil {
			t.Fatalf("expected success, got error: %v", err)
		}
		if w.ID != "wallet-1" {
			t.Fatalf("expected wallet-1, got %v", w.ID)
		}
	})

	t.Run("debit fails on insufficient balance", func(t *testing.T) {
		repo := &mockRepo{
			getOwnerWalletOverviewFn: func(ctx context.Context, ownerID string, walletType WalletType) (*OwnerVaultOverview, error) {
				return &OwnerVaultOverview{
					LedgerOwnerStatus: ledgerowner.StatusACTIVE.String(),
					VaultStatus:       entvault.StatusACTIVE.String(),
				}, nil
			},
			getWalletWithAtomicLockFn: func(ctx context.Context, tx *ent.Tx, ownerID string, walletType string, currencyCode types.CurrencyCode) (*WalletResponse, error) {
				return &WalletResponse{
					ID:               "wallet-1",
					Active:           true,
					AvailableBalance: 50,
				}, nil
			},
		}

		svc := &service{repo: repo}
		err := svc.DebitWallet(context.Background(), &DebitWalletPayload{
			OwnerID:         "owner-1",
			WalletType:      "TREASURY",
			DebitAmount:     100,
			AffectedBalance: AvailableBalance,
			Transaction:     &ent.Transaction{Currency: "NGN"},
		})

		if err == nil {
			t.Fatal("expected insufficient fund error, got nil")
		}
		if err.Error() != "insufficient fund" {
			t.Fatalf("expected 'insufficient fund', got '%v'", err.Error())
		}
	})
}

func TestCreditWallet_Conditions(t *testing.T) {
	t.Run("credits normal wallet when all conditions are active", func(t *testing.T) {
		var creditedWalletID string
		var ledgerCreated bool

		repo := &mockRepo{
			getOwnerWalletOverviewFn: func(ctx context.Context, ownerID string, walletType WalletType) (*OwnerVaultOverview, error) {
				return &OwnerVaultOverview{
					LedgerOwnerStatus: ledgerowner.StatusACTIVE.String(),
					VaultStatus:       entvault.StatusACTIVE.String(),
				}, nil
			},
			getWalletWithAtomicLockFn: func(ctx context.Context, tx *ent.Tx, ownerID string, walletType string, currencyCode types.CurrencyCode) (*WalletResponse, error) {
				return &WalletResponse{
					ID:               "treasury-wallet-1",
					Active:           true,
					AvailableBalance: 100,
				}, nil
			},
			addLedgerRecordFn: func(ctx context.Context, payload *LedgerEntry) error {
				ledgerCreated = true
				creditedWalletID = payload.Wallet.ID
				return nil
			},
			updateWalletBalanceFn: func(ctx context.Context, wallet *WalletResponse, payload *LedgerEntryResult) error {
				return nil
			},
		}

		svc := &service{repo: repo}
		err := svc.CreditWallet(context.Background(), &CreditWalletPayload{
			OwnerID:         "owner-1",
			WalletType:      "TREASURY",
			CreditAmount:    500,
			Transaction:     &ent.Transaction{ID: "tx-1", Currency: "NGN"},
			AffectedBalance: AvailableBalance,
		})

		if err != nil {
			t.Fatalf("expected credit to succeed, got %v", err)
		}
		if !ledgerCreated {
			t.Fatal("expected ledger record to be created")
		}
		if creditedWalletID != "treasury-wallet-1" {
			t.Fatalf("expected credit to treasury-wallet-1, got %v", creditedWalletID)
		}
	})

	t.Run("credits FROZEN vault when owner is not active", func(t *testing.T) {
		var creditedWalletID string
		var ledgerCreated bool

		repo := &mockRepo{
			getOwnerWalletOverviewFn: func(ctx context.Context, ownerID string, walletType WalletType) (*OwnerVaultOverview, error) {
				if walletType == FrozenWalletType {
					return &OwnerVaultOverview{
						VaultID: "frozen-vault-1",
					}, nil
				}
				return &OwnerVaultOverview{
					LedgerOwnerStatus: ledgerowner.StatusSUSPENDED.String(),
					VaultStatus:       entvault.StatusACTIVE.String(),
				}, nil
			},
			getWalletWithCurrencyFn: func(ctx context.Context, ownerID, walletType string, currencyCode types.CurrencyCode) (*WalletResponse, error) {
				if walletType == string(FrozenWalletType) {
					return &WalletResponse{
						ID:               "frozen-wallet-1",
						Active:           true,
						AvailableBalance: 0,
					}, nil
				}
				return nil, nil
			},
			addLedgerRecordFn: func(ctx context.Context, payload *LedgerEntry) error {
				ledgerCreated = true
				creditedWalletID = payload.Wallet.ID
				return nil
			},
			updateWalletBalanceFn: func(ctx context.Context, wallet *WalletResponse, payload *LedgerEntryResult) error {
				return nil
			},
		}

		svc := &service{repo: repo}
		err := svc.CreditWallet(context.Background(), &CreditWalletPayload{
			OwnerID:         "owner-1",
			WalletType:      "TREASURY",
			CreditAmount:    500,
			Transaction:     &ent.Transaction{ID: "tx-2", Currency: "NGN"},
			AffectedBalance: AvailableBalance,
		})

		if err != nil {
			t.Fatalf("expected credit to succeed, got %v", err)
		}
		if !ledgerCreated {
			t.Fatal("expected ledger record to be created")
		}
		if creditedWalletID != "frozen-wallet-1" {
			t.Fatalf("expected credit to FROZEN wallet, got %v", creditedWalletID)
		}
	})

	t.Run("credits FROZEN vault when vault is not active", func(t *testing.T) {
		var creditedWalletID string

		repo := &mockRepo{
			getOwnerWalletOverviewFn: func(ctx context.Context, ownerID string, walletType WalletType) (*OwnerVaultOverview, error) {
				if walletType == FrozenWalletType {
					return &OwnerVaultOverview{
						VaultID: "frozen-vault-1",
					}, nil
				}
				return &OwnerVaultOverview{
					LedgerOwnerStatus: ledgerowner.StatusACTIVE.String(),
					VaultStatus:       entvault.StatusSUSPENDED.String(),
				}, nil
			},
			getWalletWithCurrencyFn: func(ctx context.Context, ownerID, walletType string, currencyCode types.CurrencyCode) (*WalletResponse, error) {
				if walletType == string(FrozenWalletType) {
					return &WalletResponse{
						ID:               "frozen-wallet-1",
						Active:           true,
						AvailableBalance: 0,
					}, nil
				}
				return nil, nil
			},
			addLedgerRecordFn: func(ctx context.Context, payload *LedgerEntry) error {
				creditedWalletID = payload.Wallet.ID
				return nil
			},
			updateWalletBalanceFn: func(ctx context.Context, wallet *WalletResponse, payload *LedgerEntryResult) error {
				return nil
			},
		}

		svc := &service{repo: repo}
		err := svc.CreditWallet(context.Background(), &CreditWalletPayload{
			OwnerID:         "owner-1",
			WalletType:      "TREASURY",
			CreditAmount:    500,
			Transaction:     &ent.Transaction{ID: "tx-3", Currency: "NGN"},
			AffectedBalance: AvailableBalance,
		})

		if err != nil {
			t.Fatalf("expected credit to succeed, got %v", err)
		}
		if creditedWalletID != "frozen-wallet-1" {
			t.Fatalf("expected credit to FROZEN wallet, got %v", creditedWalletID)
		}
	})

	t.Run("credits FROZEN vault when target wallet is not active", func(t *testing.T) {
		var creditedWalletID string

		repo := &mockRepo{
			getOwnerWalletOverviewFn: func(ctx context.Context, ownerID string, walletType WalletType) (*OwnerVaultOverview, error) {
				if walletType == FrozenWalletType {
					return &OwnerVaultOverview{
						VaultID: "frozen-vault-1",
					}, nil
				}
				return &OwnerVaultOverview{
					LedgerOwnerStatus: ledgerowner.StatusACTIVE.String(),
					VaultStatus:       entvault.StatusACTIVE.String(),
				}, nil
			},
			getWalletWithAtomicLockFn: func(ctx context.Context, tx *ent.Tx, ownerID string, walletType string, currencyCode types.CurrencyCode) (*WalletResponse, error) {
				if walletType == string(TreasuryWalletType) {
					return &WalletResponse{
						ID:     "treasury-wallet-1",
						Active: false, // inactive!
					}, nil
				}
				return nil, nil
			},
			getWalletWithCurrencyFn: func(ctx context.Context, ownerID, walletType string, currencyCode types.CurrencyCode) (*WalletResponse, error) {
				if walletType == string(FrozenWalletType) {
					return &WalletResponse{
						ID:               "frozen-wallet-1",
						Active:           true,
						AvailableBalance: 0,
					}, nil
				}
				return nil, nil
			},
			addLedgerRecordFn: func(ctx context.Context, payload *LedgerEntry) error {
				creditedWalletID = payload.Wallet.ID
				return nil
			},
			updateWalletBalanceFn: func(ctx context.Context, wallet *WalletResponse, payload *LedgerEntryResult) error {
				return nil
			},
		}

		svc := &service{repo: repo}
		err := svc.CreditWallet(context.Background(), &CreditWalletPayload{
			OwnerID:         "owner-1",
			WalletType:      "TREASURY",
			CreditAmount:    500,
			Transaction:     &ent.Transaction{ID: "tx-4", Currency: "NGN"},
			AffectedBalance: AvailableBalance,
		})

		if err != nil {
			t.Fatalf("expected credit to succeed, got %v", err)
		}
		if creditedWalletID != "frozen-wallet-1" {
			t.Fatalf("expected credit to FROZEN wallet, got %v", creditedWalletID)
		}
	})

	t.Run("creates FROZEN vault and wallet if they do not exist when routing credit", func(t *testing.T) {
		var vaultCreated bool
		var walletCreated bool
		var creditedWalletID string

		repo := &mockRepo{
			getOwnerWalletOverviewFn: func(ctx context.Context, ownerID string, walletType WalletType) (*OwnerVaultOverview, error) {
				if walletType == FrozenWalletType {
					return &OwnerVaultOverview{}, nil // no vault ID yet
				}
				// Owner is suspended
				return &OwnerVaultOverview{
					LedgerOwnerStatus: ledgerowner.StatusSUSPENDED.String(),
					VaultStatus:       entvault.StatusACTIVE.String(),
				}, nil
			},
			createVaultFn: func(ctx context.Context, ownerID string, vaultType WalletType) (*ent.Vault, error) {
				vaultCreated = true
				return &ent.Vault{ID: "new-frozen-vault"}, nil
			},
			createFn: func(ctx context.Context, payload *createWallet) (*WalletResponse, error) {
				walletCreated = true
				return &WalletResponse{
					ID:     "new-frozen-wallet",
					Active: true,
				}, nil
			},
			addLedgerRecordFn: func(ctx context.Context, payload *LedgerEntry) error {
				creditedWalletID = payload.Wallet.ID
				return nil
			},
			updateWalletBalanceFn: func(ctx context.Context, wallet *WalletResponse, payload *LedgerEntryResult) error {
				return nil
			},
		}

		currSvc := &mockCurrencyService{
			getCurrencyByCodeFn: func(ctx context.Context, code types.CurrencyCode) (*currency.Currency, error) {
				return &currency.Currency{ID: "curr-ngn"}, nil
			},
		}

		svc := &service{
			repo:            repo,
			currencyService: currSvc,
		}

		err := svc.CreditWallet(context.Background(), &CreditWalletPayload{
			OwnerID:         "owner-1",
			WalletType:      "TREASURY",
			CreditAmount:    500,
			Transaction:     &ent.Transaction{ID: "tx-5", Currency: "NGN"},
			AffectedBalance: AvailableBalance,
		})

		if err != nil {
			t.Fatalf("expected credit to succeed, got %v", err)
		}
		if !vaultCreated {
			t.Fatal("expected FROZEN vault to be created")
		}
		if !walletCreated {
			t.Fatal("expected wallet in FROZEN vault to be created")
		}
		if creditedWalletID != "new-frozen-wallet" {
			t.Fatalf("expected credit to new-frozen-wallet, got %v", creditedWalletID)
		}
	})
}
