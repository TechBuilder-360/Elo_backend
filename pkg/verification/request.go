package verification

import (
	"context"

	"entgo.io/ent/dialect/sql"
	"github.com/Toflex/directory_v2/ent"
	"github.com/Toflex/directory_v2/ent/business"
	"github.com/Toflex/directory_v2/ent/requestverification"
	"github.com/Toflex/directory_v2/ent/user"
	"github.com/Toflex/directory_v2/pkg/constant"
	"github.com/Toflex/directory_v2/pkg/errors"
	"github.com/Toflex/directory_v2/pkg/log"
	"github.com/Toflex/directory_v2/pkg/provider"
	"github.com/Toflex/directory_v2/pkg/util"
)

func (s *Service) RequestVerificationLink(ctx context.Context, payload *VerificationRequest) (_ *VerificationResponse, err error) {
	logger := log.LoggerInContext(ctx)
	var (
		verification *ent.RequestVerification
		usr          *ent.User
		biz          *ent.Business
	)

	switch payload.Entity {
	case constant.UserEntityType:
		{
			usr, err = s.db.User.Query().Where(user.IDEQ(payload.ID)).First(ctx)
			if err != nil {
				return nil, errors.New(errors.ErrInvalidInput, "user not found")
			}
		}
	case constant.BusinessEntityType:
		{
			biz, err = s.db.Business.Query().Where(business.IDEQ(payload.ID)).First(ctx)
			if err != nil {
				return nil, errors.New(errors.ErrInvalidInput, "business not found")
			}
		}
	}

	switch payload.Entity {
	case constant.UserEntityType:
		{
			verification, err = s.db.RequestVerification.Query().Where(requestverification.HasUserWith(
				user.IDEQ(payload.ID)), requestverification.StatusNotIn(requestverification.StatusEXPIRED, requestverification.StatusREJECTED, requestverification.StatusFAILED)).
				Order(requestverification.ByCreatedAt(sql.OrderDesc())).First(ctx)
		}
	case constant.BusinessEntityType:
		{
			verification, err = s.db.RequestVerification.Query().Where(requestverification.HasBusinessWith(
				business.IDEQ(payload.ID)), requestverification.StatusNotIn(requestverification.StatusEXPIRED, requestverification.StatusFAILED)).
				Order(requestverification.ByCreatedAt(sql.OrderDesc())).First(ctx)
		}
	default:
		return nil, errors.New(errors.ErrInvalidInput, "invalid entity")
	}

	if verification != nil {
		switch verification.Status {
		case requestverification.StatusIN_PROGRESS, requestverification.StatusPENDING, requestverification.StatusVERIFIED:
			{
				return &VerificationResponse{
					Link:   verification.Link,
					Status: string(verification.Status),
				}, nil
			}
		}
	}

	serv, err := s.serviceLocator.GetServiceByIdentifier(ctx, constant.IdentityVerificationServiceIdentifier)
	if err != nil {
		logger.WithError(err).WithField("identifier", constant.IdentityVerificationServiceIdentifier).Error("service provider not found!")
		return nil, errors.New(errors.ErrFailed, "request failed")
	}

	impl, ok := provider.GetImpl(serv.ActiveProvider.Slug)
	if !ok {
		logger.WithField("provider", serv.ActiveProvider.Slug).Error("provider not implemented")
		return nil, errors.New(errors.ErrFailed, "request failed")
	}

	verifier, ok := provider.ConformsTo[Verifier](impl)
	if !ok {
		logger.WithField("provider", serv.ActiveProvider.Slug).Error("provider not implemented service type")
		return nil, errors.New(errors.ErrFailed, "request failed")
	}

	var (
		result      *VerifyResult
		referenceId = util.GenerateUUID()
	)

	// store validation request
	switch payload.Entity {
	case constant.UserEntityType:
		{
			result, err = verifier.VerifyUser(ctx, referenceId)
			if err != nil {
				logger.WithError(err).Error("failed to generate user verification link")
				return nil, errors.New(errors.ErrFailed, "request failed")
			}

			_, err = s.db.RequestVerification.Create().
				SetLink(result.Link).
				SetProviderLink(result.ProviderLink).
				SetReferenceID(result.ReferenceID).
				AddUser(usr).
				Save(ctx)
		}
	case constant.BusinessEntityType:
		{
			result, err = verifier.VerifyBusiness(ctx, referenceId)
			if err != nil {
				logger.WithError(err).Error("failed to generate business verification link")
				return nil, errors.New(errors.ErrFailed, "request failed")
			}

			_, err = s.db.RequestVerification.Create().
				SetLink(result.Link).
				SetProviderLink(result.ProviderLink).
				SetReferenceID(result.ReferenceID).
				AddBusiness(biz).
				Save(ctx)
		}
	}

	if err != nil {
		logger.WithError(err).Error("failed to store validation request")
		return nil, errors.New(errors.ErrFailed, "request failed")
	}

	return &VerificationResponse{
		Status: string(requestverification.StatusIN_PROGRESS),
		Link:   result.Link,
	}, nil
}
