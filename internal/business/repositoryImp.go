package business

import (
	"context"

	"github.com/Toflex/directory_v2/ent"
	"github.com/Toflex/directory_v2/ent/business"
)

func (r *repository) WithTransaction(db *ent.Client) IRepository {
	return &repository{
		db: db,
	}
}

func (r *repository) GetBusinessByName(ctx context.Context, name string) (*businessResult, error) {
	b, err := r.db.Business.Query().
		Where(business.NameEqualFold(name)).
		Only(ctx)

	if err != nil {
		if !ent.IsNotFound(err) {
			return nil, err
		}

		return nil, err
	}

	return &businessResult{*b}, nil
}

func (r *repository) GetBusinessByID(ctx context.Context, id string) (*businessResult, error) {
	b, err := r.db.Business.Query().
		Where(business.IDEQ(id)).
		Only(ctx)

	if err != nil {
		if !ent.IsNotFound(err) {
			return nil, err
		}

		return nil, err
	}

	return &businessResult{*b}, nil
}

func (r *repository) Create(ctx context.Context, payload createBusiness) error {
	_, err := r.db.Business.
		Create().
		SetName(payload.Name).
		SetCategory(payload.Category).
		SetEmail(payload.Email).
		// SetWebsite(payload.Website).
		Save(ctx)

	return err
}
