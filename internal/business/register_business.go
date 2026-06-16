package business

import (
	"context"
	"errors"

	"github.com/Toflex/directory_v2/pkg/log"
	"github.com/Toflex/directory_v2/pkg/util"
)

func (s *service) CreateBusiness(ctx context.Context, payload CreateBusinessRequest, logger log.Entry) error {
	err := payload.Validate()
	if err != nil {
		return err
	}

	user, err := s.userRepository.GetByID(ctx, payload.User.ID)
	if err != nil {
		logger.WithField("ID", payload.User.ID).WithError(err).Error("user request failed")
		return errors.New("request failed")
	}

	if user == nil {
		return errors.New("user not found!")
	}

	if !user.Verified {
		return errors.New("user not verified, please present a valid means of identification")
	}

	// TODO: Prevent user from adding a new business if there exist a business registered but not verified (KYB not completed)
	if err := util.ValidateEmail(payload.Email); err != nil {
		logger.WithError(err).Error("business email is invalid")
		return err
	}

	business, err := s.repo.GetBusinessByName(ctx, payload.Name)
	if err != nil {
		logger.WithError(err).Error("failed to fetch business by name")
		return errors.New("request failed")
	}

	if business != nil {
		return errors.New("business name already registered")
	}

	// body := createBusiness{
	// 	Name:     util.ToTitleCase(payload.Name),
	// 	Category: payload.Category,
	// 	Email:    strings.ToLower(payload.Email),
	// }

	// err = s.repo.Create(ctx, body)
	// if err != nil {
	// 	logger.WithError(err).Error("failed to create business on db")
	// 	return errors.New("request failed")
	// }

	return nil
}
