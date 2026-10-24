package vault

import (
	"context"

	"github.com/Toflex/directory_v2/pkg/errors"
	"github.com/Toflex/directory_v2/pkg/log"
)

func (s *Service) CreateVault(ctx context.Context, ownerID string) (*Vault, error) {
	return s.CreateVaultWithType(ctx, ownerID, string(Treasury))
}

func (s *Service) CreateVaultWithType(ctx context.Context, ownerID, vaultType string) (*Vault, error) {
	logger := log.LoggerInContext(ctx)

	v, err := s.repo.CreateWithType(ctx, ownerID, vaultType)
	if err != nil {
		logger.WithError(err).Error("failed to create vault")
		return nil, errors.New(errors.ErrFailed, "something went wrong")
	}

	return &Vault{
		ID:     v.ID,
		Type:   string(v.Type),
		Status: string(v.Status),
	}, nil
}

func (s *Service) GetVault(ctx context.Context, id string) (*Vault, error) {
	logger := log.LoggerInContext(ctx)

	v, err := s.repo.GetByID(ctx, id)
	if err != nil {
		logger.WithError(err).Error("failed to fetch vault")
		return nil, errors.New(errors.ErrFailed, "something went wrong")
	}

	return &Vault{
		ID:     v.ID,
		Type:   string(v.Type),
		Status: string(v.Status),
	}, nil
}

func (s *Service) GetOwnerVault(ctx context.Context, ownerID, vaultType string) (*Vault, error) {
	logger := log.LoggerInContext(ctx)

	v, err := s.repo.GetByOwner(ctx, ownerID, vaultType)
	if err != nil {
		logger.WithError(err).Error("failed to fetch vault")
		return nil, errors.New(errors.ErrFailed, "something went wrong")
	}

	if v == nil {
		return s.CreateVaultWithType(ctx, ownerID, vaultType)
	}

	return &Vault{
		ID:     v.ID,
		Type:   string(v.Type),
		Status: string(v.Status),
	}, nil
}
