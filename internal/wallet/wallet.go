package wallet

import (
	"context"

	"github.com/Toflex/directory_v2/ent"
	"github.com/Toflex/directory_v2/graph/model"
	"github.com/Toflex/directory_v2/pkg/errors"
	"github.com/Toflex/directory_v2/pkg/log"
	"github.com/Toflex/directory_v2/pkg/types"
	"github.com/Toflex/directory_v2/pkg/util"
)

// GetWallets implements [IService].
func (s *service) GetWallets(ctx context.Context, b *ent.Business, walletType string) ([]*model.Wallet, error) {
	logger := log.LoggerInContext(ctx)
	result, err := s.repo.GetWallets(ctx, b.ID, getWalletType(walletType))
	if err != nil {
		logger.WithError(err).Error("failed to fetch wallets")
		return nil, errors.New(errors.ErrFailed, "request failed")
	}

	wallets := make([]*model.Wallet, 0)

	for _, wallet := range result {
		wallets = append(wallets, &model.Wallet{
			Type:     wallet.Type,
			Currency: wallet.Currency.ToString(),
			AvailableBalance: util.ToMajorUnit(types.ToMajor{
				Amount:    wallet.AvailableBalance,
				Precision: uint(wallet.Multiplier),
			}),
			LedgerBalance: util.ToMajorUnit(types.ToMajor{
				Amount:    wallet.LedgerBalance,
				Precision: uint(wallet.Multiplier),
			}),
			HoldingBalance: util.ToMajorUnit(types.ToMajor{
				Amount:    wallet.HoldingBalance,
				Precision: uint(wallet.Multiplier),
			}),
			Active: wallet.Active,
			ID:     wallet.ID,
		})
	}

	return wallets, nil
}

// GetWallet implements [IService].
func (s *service) GetWallet(ctx context.Context, b *ent.Business, walletType, currencyCode string) (*model.Wallet, error) {
	logger := log.LoggerInContext(ctx)
	wallet, err := s.repo.GetWallet(ctx, b.ID, getWalletType(walletType), types.CurrencyCode(currencyCode))
	if err != nil {
		logger.WithError(err).Error("failed to fetch wallet")
		return nil, errors.New(errors.ErrFailed, "request failed")
	}

	if wallet == nil {
		return nil, errors.New(errors.ErrFailed, "wallet not found")
	}

	return &model.Wallet{
		Type:     wallet.Type,
		Currency: wallet.Currency.ToString(),
		AvailableBalance: util.ToMajorUnit(types.ToMajor{
			Amount:    wallet.AvailableBalance,
			Precision: uint(wallet.Multiplier),
		}),
		LedgerBalance: util.ToMajorUnit(types.ToMajor{
			Amount:    wallet.LedgerBalance,
			Precision: uint(wallet.Multiplier),
		}),
		HoldingBalance: util.ToMajorUnit(types.ToMajor{
			Amount:    wallet.HoldingBalance,
			Precision: uint(wallet.Multiplier),
		}),
		Active: wallet.Active,
		ID:     wallet.ID,
	}, nil
}

func (s *service) AddWallet(ctx context.Context, b *ent.Business, walletType, currencyCode string) (*model.Wallet, error) {
	logger := log.LoggerInContext(ctx)

	if !validateWalletType(walletType) {
		logger.WithField("wallet_type", walletType).Error("unable to validate wallet type")
		return nil, errors.New(errors.ErrInvalidInput, "invalid wallet type")
	}

	currency, err := s.currencyRepo.GetCurrencyByCode(ctx, types.CurrencyCode(currencyCode))
	if err != nil {
		logger.WithError(err).Error("unable to fetch currencies")
		return nil, errors.New(errors.ErrFailed, "request failed")
	}

	if currency == nil || !currency.Active {
		logger.WithField("currency_active", currency.Active).Error("currency is not active")
		return nil, errors.New(errors.ErrFailed, "currency is not available")
	}

	wallet, err := s.repo.GetWallet(ctx, b.ID, getWalletType(walletType), types.CurrencyCode(currencyCode))
	if err != nil {
		logger.WithError(err).Error("wallet could not be retrived")
		return nil, errors.New(errors.ErrFailed, string(errors.ErrFailed))
	}

	if wallet != nil {
		return &model.Wallet{
			Type:     wallet.Type,
			Currency: wallet.Currency.ToString(),
			AvailableBalance: util.ToMajorUnit(types.ToMajor{
				Amount:    wallet.AvailableBalance,
				Precision: uint(wallet.Multiplier),
			}),
			LedgerBalance: util.ToMajorUnit(types.ToMajor{
				Amount:    wallet.LedgerBalance,
				Precision: uint(wallet.Multiplier),
			}),
			HoldingBalance: util.ToMajorUnit(types.ToMajor{
				Amount:    wallet.HoldingBalance,
				Precision: uint(wallet.Multiplier),
			}),
			Active: wallet.Active,
			ID:     wallet.ID,
		}, nil
	}

	w := &createWallet{
		Type:       getWalletType(walletType),
		Currency:   currency,
		IsBusiness: true,
		Identifier: b.ID,
	}

	wallet, err = s.repo.Create(ctx, w)
	if err != nil {
		logger.WithError(err).Error("unable to create wallet")
		return nil, errors.New(errors.ErrFailed, "request failed")
	}

	return &model.Wallet{
		Type:     wallet.Type,
		Currency: wallet.Currency.ToString(),
		AvailableBalance: util.ToMajorUnit(types.ToMajor{
			Amount:    wallet.AvailableBalance,
			Precision: uint(wallet.Multiplier),
		}),
		LedgerBalance: util.ToMajorUnit(types.ToMajor{
			Amount:    wallet.LedgerBalance,
			Precision: uint(wallet.Multiplier),
		}),
		HoldingBalance: util.ToMajorUnit(types.ToMajor{
			Amount:    wallet.HoldingBalance,
			Precision: uint(wallet.Multiplier),
		}),
		Active: wallet.Active,
		ID:     wallet.ID,
	}, nil
}
