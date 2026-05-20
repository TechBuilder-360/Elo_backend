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

func (r *repository) GetProvider(ctx context.Context, name string) (*ActiveProvider, error) {
	activeProvider := new(ActiveProvider)
	err := r.db.Provider.
		Query().
		Where(provider.NameEqualFold(name)).
		Select(provider.FieldName, provider.FieldSlug, provider.FieldActive).
		Scan(ctx, &activeProvider)

	if err != nil {
		return nil, err
	}

	return activeProvider, nil
}
