package business

import (
	"context"

	"github.com/Toflex/directory_v2/ent"
	"github.com/Toflex/directory_v2/ent/business"
	"github.com/Toflex/directory_v2/ent/manager"
	biz "github.com/Toflex/directory_v2/pkg/business"
	"github.com/Toflex/directory_v2/pkg/errors"
	"github.com/Toflex/directory_v2/pkg/log"
	"github.com/Toflex/directory_v2/pkg/util"
)

func (s *service) GetBusiness(ctx context.Context, user *ent.User, id string, logger log.Entry) (*BusinessResult, error) {
	business, err := s.db.Business.Query().Where(business.IDEQ(id), business.HasManagesWith(manager.UserID(user.ID))).WithLocations().First(ctx)
	if err != nil {
		logger.WithError(err).WithField("business_id", id).Error("failed to fetch business")
		return nil, errors.New(errors.ErrNotFound, "not found")
	}

	address := BusinessAddress{}
	if len(business.Edges.Locations) > 0 {
		address = BusinessAddress{
			City:    business.Edges.Locations[0].City,
			Street:  business.Edges.Locations[0].Address,
			State:   business.Edges.Locations[0].State,
			Country: business.Edges.Locations[0].Country,
			ZipCode: business.Edges.Locations[0].ZipCode,
		}
	}

	return &BusinessResult{
		ID:                     business.ID,
		Name:                   business.Name,
		Logo:                   &business.Logo,
		Email:                  &business.Email,
		About:                  &business.About,
		OnSite:                 business.OnSite,
		Number:                 util.AddressToString(business.RegistrationNumber),
		Industry:               business.Category,
		CountryOfIncorporation: util.AddressToString(business.CountryOfIncorporation),
		DateOfIncorporation:    util.AddressToString(business.DateOfIncorporation),
		Address:                address,
	}, nil
}

func (s *service) Businesses(ctx context.Context, user *ent.User, logger log.Entry) ([]*MyBusinessResult, error) {
	managers, err := s.db.Manager.Query().
		Where(manager.UserID(user.ID)).
		WithBusiness().
		WithRole().
		All(ctx)
	if err != nil {
		logger.WithError(err).Error("failed to fetch businesses")
		return nil, errors.New(errors.ErrNotFound, "not found")
	}

	results := make([]*MyBusinessResult, 0, len(managers))
	for _, managerEntry := range managers {
		if managerEntry.Edges.Business == nil {
			continue
		}

		roleName := ""
		if managerEntry.Edges.Role != nil {
			roleName = managerEntry.Edges.Role.Name
		}

		business := managerEntry.Edges.Business
		logo := &business.Logo
		results = append(results, &MyBusinessResult{
			ID:   business.ID,
			Name: business.Name,
			Logo: logo,
			Role: &roleName,
		})
	}

	return results, nil
}

func (s *service) GetCategory(ctx context.Context) []string {
	return biz.GetCategory()
}
