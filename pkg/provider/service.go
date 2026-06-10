package provider

import (
	"context"

	"github.com/Toflex/directory_v2/database/database"
	"github.com/Toflex/directory_v2/ent"
	"github.com/Toflex/directory_v2/ent/schema"
	"github.com/Toflex/directory_v2/pkg/log"
)

type IService interface {
	GetServiceByIdentifier(ctx context.Context, identifier string) (ServiceLocator, error)
}

type service struct {
	db   *ent.Client
	repo IRepository
}

func NewService() IService {
	db := database.DBInstance()
	return &service{
		db:   db,
		repo: Newrepository(db),
	}
}

func (s *service) GetServiceByIdentifier(ctx context.Context, identifier string) (ServiceLocator, error) {
	logger := log.LoggerInContext(ctx)
	serviceLocator := ServiceLocator{}

	// get service by identifier
	serv, err := s.repo.FindByIdentifier(ctx, identifier)
	if err != nil {
		logger.WithField("identifier", identifier).WithError(err).Error("failed to find service by identifier")
		return serviceLocator, err
	}

	// Get active provider
	provider, err := s.repo.GetProvider(ctx, serv.Provider)
	if err != nil {
		logger.WithField("Provider", serv.Provider).WithError(err).Error("failed to get provider")
		return serviceLocator, err
	}

	serviceLocator = ServiceLocator{
		Name:       serv.Name,
		Identifier: serv.Identifier,
		ActiveProvider: ActiveProvider{
			Name:   provider.Name,
			Slug:   provider.Slug,
			Active: provider.Active,
		},
		Fee: Fee{
			Type:  FeeType(serv.Fee.Type),
			Value: serv.Fee.Value,
			Min:   serv.Fee.Min,
			Max:   serv.Fee.Max,
			Tiers: func(tier []schema.Tier) []Tier {
				tiers := make([]Tier, 0)
				if len(tier) == 0 {
					return tiers
				}

				for _, t := range tier {
					tiers = append(tiers, Tier{
						From:  t.From,
						To:    t.To,
						Type:  FeeType(t.Type),
						Value: t.Value,
					})
				}

				return tiers
			}(serv.Fee.Tiers),
		},
		MinValue: int64(serv.Min),
		MaxValue: int64(serv.Max),
	}

	return serviceLocator, nil
}
