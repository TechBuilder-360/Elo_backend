package verification

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/Toflex/directory_v2/ent"
	"github.com/Toflex/directory_v2/ent/requestverification"
	"github.com/Toflex/directory_v2/ent/verification"
	"github.com/Toflex/directory_v2/internal/email"
	unitofwork "github.com/Toflex/directory_v2/internal/unit_of_work"
	"github.com/Toflex/directory_v2/pkg/configuration"
	"github.com/Toflex/directory_v2/pkg/constant"
	"github.com/Toflex/directory_v2/pkg/types"
)

func ToMap(v interface{}) (map[string]interface{}, error) {
	if v == nil {
		return nil, nil
	}

	if m, ok := v.(map[string]interface{}); ok {
		return m, nil
	}

	jsonData, err := json.Marshal(v)
	if err != nil {
		return nil, err
	}

	var m map[string]interface{}
	if err := json.Unmarshal(jsonData, &m); err != nil {
		return nil, err
	}

	return m, nil
}

func (s *Service) processSuccessVerification(ctx context.Context, payload *VerificationResult, verificationRequest *ent.RequestVerification) error {
	metadata, err := ToMap(payload.Metadata)
	if err != nil {
		return err
	}

	var (
		data             = make(map[string]interface{})
		verificationType verification.VerificationType
		entityType       types.EntityType
		user             *ent.User
		business         *ent.Business
	)

	if len(verificationRequest.Edges.User) > 0 {
		entityType = constant.UserEntityType
		user = verificationRequest.Edges.User[0]
	}

	if len(verificationRequest.Edges.Business) > 0 {
		entityType = constant.BusinessEntityType
		business = verificationRequest.Edges.Business[0]
	}

	var found bool

	if payload.BVN.Number != "" {
		verificationType = verification.VerificationTypeBVN
		data, err = ToMap(payload.BVN)
		if err != nil {
			return err
		}
		found = true
	} else if payload.NationalID.Number != "" {
		verificationType = verification.VerificationTypeNIN
		data, err = ToMap(payload.NationalID)
		if err != nil {
			return err
		}
		found = true
	} else if payload.Passport.Number != "" {
		verificationType = verification.VerificationTypePASSPORT
		data, err = ToMap(payload.Passport)
		if err != nil {
			return err
		}
		found = true
	} else if payload.DriversLicense.Number != "" {
		verificationType = verification.VerificationTypeDRIVERS_LICENSE
		data, err = ToMap(payload.DriversLicense)
		if err != nil {
			return err
		}
		found = true
	} else if payload.VoterID.Number != "" {
		verificationType = verification.VerificationTypeVOTER_ID
		data, err = ToMap(payload.VoterID)
		if err != nil {
			return err
		}
		found = true
	}

	if !found {
		// nothing to save
		return errors.New("identity documents not found")
	}

	if verificationType != verification.VerificationTypeBVN {
		img, ok := data["front_image"].(string)
		if ok && img != "" {
			fileName := fmt.Sprintf("front_%s", string(verificationType))
			url, err := s.cloudinary.UploadUserFile(ctx, img, fileName, user.ID, "verification")
			if err != nil {
				return err
			}

			data["front_image"] = url
		}

		img, ok = data["back_image"].(string)
		if ok && img != "" {
			fileName := fmt.Sprintf("back_%s", string(verificationType))
			url, err := s.cloudinary.UploadUserFile(ctx, img, fileName, user.ID, "verification")
			if err != nil {
				return err
			}

			data["back_image"] = url
		}
	}

	// begin db transaction
	uw := unitofwork.NewUnitOfWorkRepository(s.db)
	tx, err := uw.Begin(ctx)
	if err != nil {
		return err
	}

	defer func() {
		if err != nil {
			_ = uw.Rollback(tx)
		}
	}()

	vc := tx.Verification.Create().
		SetProvider(string(payload.Provider)).
		SetProviderReference(payload.ReferenceID).
		SetVerificationType(verificationType).
		SetStatus(verification.StatusVERIFIED).
		SetReferenceID(payload.ReferenceID)

	if metadata != nil {
		vc = vc.SetMetadata(metadata)
	}
	if len(data) > 0 {
		vc = vc.SetData(data)
	}

	if user != nil {
		vc = vc.AddUser(user)
	}
	if business != nil {
		vc = vc.AddBusiness(business)
	}

	_, err = vc.Save(ctx)
	if err != nil {
		return err
	}

	// update verification request
	_, err = tx.RequestVerification.UpdateOne(verificationRequest).
		SetStatus(requestverification.StatusVERIFIED).
		SetMessage(payload.Message).
		Save(ctx)

	if err != nil {
		return err
	}

	// update user to verified
	if user != nil {
		_, err = tx.User.UpdateOne(user).
			SetVerified(true).
			Save(ctx)
		if err != nil {
			return err
		}
	}

	//Verify details
	switch entityType {
	case constant.UserEntityType:
		{

		}
	case constant.BusinessEntityType:
		{
		}
	}

	// send success mail notice
	if user != nil {
		err = email.NewUserVerificationTask(email.VerificationMailPayload{
			ToMail:             user.EmailAddress,
			FullName:           fmt.Sprintf("%s %s", user.FirstName, user.LastName),
			VerificationStatus: string(payload.Status),
			StatusMessage:      payload.Message,
			DashboardURL:       configuration.Instance.DashboardUrl,
		})
		if err != nil {
			return err
		}
	}

	// end db transaction
	return uw.Commit(tx)
}

func (s *Service) ProcessVerification(ctx context.Context, payload *VerificationResult) error {
	verificationRequest, err := s.db.RequestVerification.Query().
		Where(requestverification.ReferenceIDEQ(payload.ReferenceID)).
		WithUser().
		WithBusiness().
		First(ctx)
	if err != nil {
		return err
	}

	if verificationRequest.Status != requestverification.StatusPENDING {
		return nil
	}

	switch payload.Status {
	case constant.Success:
		return s.processSuccessVerification(ctx, payload, verificationRequest)
	case constant.Failed:
		return s.processFailedVerification(ctx, payload, verificationRequest)
	default:
		return nil
	}
}

func (s *Service) processFailedVerification(ctx context.Context, payload *VerificationResult, verificationRequest *ent.RequestVerification) error {
	// begin db transaction
	uw := unitofwork.NewUnitOfWorkRepository(s.db)
	tx, err := uw.Begin(ctx)
	if err != nil {
		return err
	}

	defer func() {
		if err != nil {
			_ = uw.Rollback(tx)
		}
	}()

	_, err = tx.RequestVerification.UpdateOne(verificationRequest).
		SetStatus(requestverification.StatusFAILED).
		SetMessage(payload.Message).
		Save(ctx)

	if err != nil {
		return err
	}

	var user *ent.User
	if len(verificationRequest.Edges.User) > 0 {
		user = verificationRequest.Edges.User[0]
	}

	if user != nil {
		err = email.NewUserVerificationTask(email.VerificationMailPayload{
			ToMail:             user.EmailAddress,
			FullName:           fmt.Sprintf("%s %s", user.FirstName, user.LastName),
			VerificationStatus: string(payload.Status),
			StatusMessage:      payload.Message,
			DashboardURL:       configuration.Instance.DashboardUrl,
		})
		if err != nil {
			return err
		}
	}
	if err != nil {
		return err
	}

	return uw.Commit(tx)
}

func (s *Service) verifyUserInformation() {

}
