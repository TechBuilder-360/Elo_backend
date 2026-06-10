package verification

import (
	"context"
	"encoding/json"

	"github.com/Toflex/directory_v2/ent"
	"github.com/Toflex/directory_v2/ent/requestverification"
	"github.com/Toflex/directory_v2/ent/verification"
	unitofwork "github.com/Toflex/directory_v2/internal/unit_of_work"
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

func (s *Service) ProcessVerification(ctx context.Context, payload *VerificationResult) error {
	verificationRequest, err := s.db.RequestVerification.Query().
		Where(requestverification.ReferenceID(payload.ReferenceID)).
		WithUser().
		WithBusiness().
		First(ctx)
	if err != nil {
		return err
	}

	if verificationRequest.Status != requestverification.StatusPENDING {
		return nil
	}

	metadata, err := ToMap(payload.Metadata)
	if err != nil {
		return err
	}

	var (
		data             map[string]interface{}
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

	switch {
	case payload.BVN != nil:
		verificationType = verification.VerificationTypeBVN
		data, err = ToMap(payload.BVN)
	case payload.NationalID != nil:
		verificationType = verification.VerificationTypeNIN
		data, err = ToMap(payload.NationalID)
	case payload.Passport != nil:
		verificationType = verification.VerificationTypePASSPORT
		data, err = ToMap(payload.Passport)
	case payload.DriversLicense != nil:
		verificationType = verification.VerificationTypeDRIVERS_LICENSE
		data, err = ToMap(payload.DriversLicense)
	case payload.VoterID != nil:
		verificationType = verification.VerificationTypeVOTER_ID
		data, err = ToMap(payload.VoterID)
	default:
		return nil
	}

	if err != nil {
		return err
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

	_, err = tx.Verification.Create().
		SetProvider(string(payload.Provider)).
		SetVerificationType(verificationType).
		SetStatus(verification.StatusVERIFIED).
		SetReferenceID(payload.ReferenceID).
		SetMetadata(metadata).
		SetData(data).
		AddUser(user).
		AddBusiness(business).
		Save(ctx)
	if err != nil {
		return err
	}

	_, err = tx.RequestVerification.UpdateOne(verificationRequest).
		SetStatus(requestverification.StatusVERIFIED).
		Save(ctx)

	//Verify details
	switch entityType {
	case constant.UserEntityType:
		{

		}
	case constant.BusinessEntityType:
		{
		}
	}

	// end db transaction
	return uw.Commit(tx)
}

func (s *Service) verifyUserInformation() {

}
