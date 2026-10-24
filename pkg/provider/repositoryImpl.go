package provider

import (
	"context"
	"time"

	"github.com/Toflex/directory_v2/ent"
	fee "github.com/Toflex/directory_v2/ent/feerule"
	"github.com/Toflex/directory_v2/ent/provider"
	s "github.com/Toflex/directory_v2/ent/service"
)

func (r *repository) WithTransaction(db *ent.Client) IRepository {
	return &repository{
		db: db,
	}
}

func (r *repository) FindByIdentifier(ctx context.Context, identifier string) (*ent.Service, error) {
	now := time.Now()
	serv, err := r.db.Service.
		Query().
		Where(s.IdentifierEqualFold(identifier)).
		WithFeeRules(func(q *ent.FeeRuleQuery) {
			q.WithTiers().
				Where(
					fee.IsActive(true),
					fee.EffectiveFromLTE(now),
					fee.Or(
						fee.EffectiveToIsNil(),
						fee.EffectiveToGTE(now),
					),
				).
				Order(ent.Asc(fee.FieldMinValue))
		}).
		First(ctx)

	if err != nil {
		return nil, err
	}

	return serv, nil
}

func (r *repository) GetProvider(ctx context.Context, identifier string) (*ActiveProvider, error) {
	activeProvider, err := r.db.Provider.
		Query().
		Where(provider.Or(
			provider.NameEqualFold(identifier),
			provider.SlugEqualFold(identifier),
		)).
		First(ctx)

	if err != nil {
		return nil, err
	}

	return &ActiveProvider{
		Name:   activeProvider.Name,
		Slug:   activeProvider.Slug,
		Active: activeProvider.Active,
	}, nil
}
