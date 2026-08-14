package business

import (
	"context"
	"time"

	"github.com/Toflex/directory_v2/ent"
	"github.com/Toflex/directory_v2/graph/model"
	"github.com/Toflex/directory_v2/internal/wallet"
	"github.com/Toflex/directory_v2/pkg/errors"
	"github.com/Toflex/directory_v2/pkg/log"
	"github.com/Toflex/directory_v2/pkg/saferoutine"
)

func (s *service) GetBusinessWallets(ctx context.Context, b *ent.Business, walletType string) ([]*model.Wallet, error) {
	logger := log.LoggerInContext(ctx)

	if !wallet.ValidateWalletType(walletType) {
		logger.WithField("wallet_type", walletType).Error("unable to validate wallet type")
		return nil, errors.New(errors.ErrInvalidInput, "invalid wallet type")
	}

	owner, err := b.QueryOwner().First(ctx)
	if err != nil {
		logger.WithError(err).Error("failed to fetch business owner in context")
		return nil, errors.New(errors.ErrFailed, "something went wrong")
	}

	// if business owner is yet to be created
	// Create Owner and Vault in a background job
	if owner == nil {
		saferoutine.Run(func() {
			ctx, cancel := context.WithTimeout(context.Background(), time.Second*20)
			defer cancel()

			s.createBusinessOwner(ctx, b, walletType)
		})

		return []*model.Wallet{}, nil
	}

	ownerID := owner.ID

	return s.walletService.GetWallets(ctx, ownerID, walletType)
}

func (s *service) createBusinessOwner(ctx context.Context, b *ent.Business, walletType string) {
	logger := log.LoggerInContext(ctx)
	owner, err := b.QueryOwner().First(ctx)
	if err != nil {
		logger.WithError(err).Error("failed to fetch business owner in context")
		return
	}

	if owner == nil {
		owner, err := s.repo.CreateOwner(ctx, b)
		if err != nil || owner == nil {
			logger.WithError(err).Error("failed to create business ledger owner")
		}
	}

	vault, err := s.vaultService.GetOwnerVault(ctx, owner.ID, string(wallet.GetWalletType(walletType)))
	if err != nil || vault == nil {
		logger.WithError(err).Error("failed to fetch business vault")
	}

	if vault == nil {
		_, err := s.vaultService.CreateVault(ctx, owner.ID)
		if err != nil {
			logger.WithError(err).Error("failed to create vault")
		}
	}
}

func (s *service) GetBusinessWallet(ctx context.Context, b *ent.Business, walletType string, currencyCode string) (*model.Wallet, error) {
	logger := log.LoggerInContext(ctx)

	if !wallet.ValidateWalletType(walletType) {
		logger.WithField("wallet_type", walletType).Error("unable to validate wallet type")
		return nil, errors.New(errors.ErrInvalidInput, "invalid wallet type")
	}

	owner, err := b.QueryOwner().First(ctx)
	if err != nil {
		logger.WithError(err).Error("failed to fetch business owner in context")
		return nil, errors.New(errors.ErrFailed, "something went wrong")
	}

	// if business owner is yet to be created
	// Create Owner and Vault in a background job
	if owner == nil {
		saferoutine.Run(func() {
			ctx, cancel := context.WithTimeout(context.Background(), time.Second*20)
			defer cancel()

			s.createBusinessOwner(ctx, b, walletType)
		})

		return nil, nil
	}

	ownerID := owner.ID

	return s.walletService.GetWallet(ctx, ownerID, walletType, currencyCode)
}
