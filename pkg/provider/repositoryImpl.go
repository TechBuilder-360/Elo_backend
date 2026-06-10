package provider

import (
	"context"

	"github.com/Toflex/directory_v2/ent"
	"github.com/Toflex/directory_v2/ent/provider"
	s "github.com/Toflex/directory_v2/ent/service"
)

func (r *repository) WithTransaction(db *ent.Client) IRepository {
	return &repository{
		db: db,
	}
}

func (r *repository) FindByIdentifier(ctx context.Context, identifier string) (*ent.Service, error) {
	serv, err := r.db.Service.
		Query().
		Where(s.IdentifierEqualFold(identifier)).
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
