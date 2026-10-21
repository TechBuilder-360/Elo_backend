package wallet

import (
	"context"

	"github.com/Toflex/directory_v2/ent"
	"github.com/Toflex/directory_v2/pkg/types"
)

type IRepository interface {
	GetWallets(ctx context.Context, ownerID string, walletType WalletType, filter *WalletFilter) ([]WalletResponse, error)
	GetWallet(ctx context.Context, walletType, ownerID string) (*WalletResponse, error)
	GetWalletByOwnerID(ctx context.Context, walletID, ownerID string) (*WalletResponse, error)
	GetWalletWithCurrency(ctx context.Context, ownerID, walletType string, currencyCode types.CurrencyCode) (*WalletResponse, error)
	GetOwnerWalletOverview(ctx context.Context, ownerID string, walletType WalletType) (*OwnerVaultOverview, error)
	Create(ctx context.Context, payload *createWallet) (*WalletResponse, error)
	CreateVault(ctx context.Context, ownerID string, vaultType WalletType) (*ent.Vault, error)
	GetWalletWithAtomicLock(ctx context.Context, tx *ent.Tx, walletType, ownerID string, currencyCode types.CurrencyCode) (*WalletResponse, error)
	WithTransaction(tx *ent.Tx) IRepository
	AddLedgerRecord(ctx context.Context, payload *LedgerEntry) error
	UpdateWalletBalance(ctx context.Context, wallet *WalletResponse, payload *LedgerEntryResult) error
	GetLedgerbyTransactionID(ctx context.Context, transactionID string) (*ent.Transaction, error)
}

type repository struct {
	db *ent.Client
}

func Newrepository(db *ent.Client) IRepository {
	return &repository{
		db: db,
	}
}
