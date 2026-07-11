package business

import (
	"context"
	"fmt"
	"strconv"

	"github.com/Toflex/directory_v2/internal/manager"
	unitofwork "github.com/Toflex/directory_v2/internal/unit_of_work"
	rbac "github.com/Toflex/directory_v2/pkg/RBAC"
	"github.com/Toflex/directory_v2/pkg/errors"
	"github.com/Toflex/directory_v2/pkg/log"
	"github.com/Toflex/directory_v2/pkg/util"
)

func (s *service) CreateBusiness(ctx context.Context, payload CreateBusinessRequest, logger log.Entry) error {
	if !payload.User.Verified {
		return errors.New(errors.ErrFailed, "user not verified, please present a valid means of identification")
	}

	// Add random string to business name until they gets verified
	payload.Name = fmt.Sprintf("%s_%s", payload.Name, strconv.Itoa(util.RandomInt()))

	uw := unitofwork.NewUnitOfWorkRepository(s.db)
	tx, err := uw.Begin(ctx)
	if err != nil {
		return errors.New(errors.ErrFailed, "failed to begin transaction")
	}

	// rollback uncommited transaction in case of error
	defer tx.Rollback()

	businessTx := s.repo.WithTransaction(tx)

	body := createbusiness{
		Name:         payload.Name,
		Category:     payload.Industry,
		About:        payload.About,
		Email:        payload.Email,
		RegisteredBy: payload.User.ID,
	}

	business, err := businessTx.Create(ctx, body)
	if err != nil {
		logger.WithError(err).
			WithField("business", body.Name).
			WithField("registered by", body.RegisteredBy).
			Error("failed to create business on db")
		return errors.New(errors.ErrFailed, "request failed")
	}

	location := createBusinssLocation{
		Name:         payload.Name,
		Street:       payload.Address.Street,
		City:         payload.Address.City,
		State:        payload.Address.State,
		Country:      payload.Address.Country,
		ZipCode:      payload.Address.ZipCode,
		IsHeadOffice: true,
		Business:     business,
	}
	err = businessTx.CreateBusinssLocation(ctx, location)
	if err != nil {
		logger.WithError(err).
			WithField("business", body.Name).
			WithField("registered by", body.RegisteredBy).
			Error("failed to create business location on db")
		return errors.New(errors.ErrFailed, "request failed")
	}

	// Member role is set to authorized representative if the user is the one registering the business, else it is set to member
	role, err := s.rbacRepository.GetRole(ctx, rbac.AdminRole.ToString())
	if err != nil {
		logger.WithError(err).
			WithField("business", body.Name).
			WithField("registered by", body.RegisteredBy).
			Error("failed to get role")
		return errors.New(errors.ErrFailed, "request failed")
	}

	err = s.managerRepository.WithTransaction(tx).AddManager(ctx, &manager.Manager{
		UserID:     payload.User.ID,
		BusinessID: business.ID,
		IsOwner:    false,
		RoleID:     role.ID,
	})
	if err != nil {
		logger.WithError(err).
			WithField("business", body.Name).
			WithField("registered by", body.RegisteredBy).
			Error("failed to add user to business members on db")
		return errors.New(errors.ErrFailed, "request failed")
	}

	err = tx.Commit()
	if err != nil {
		logger.WithError(err).
			WithField("business", body.Name).
			WithField("registered by", body.RegisteredBy).
			Error("failed to commit transaction")
		return errors.New(errors.ErrFailed, "request failed")
	}

	return nil
}
