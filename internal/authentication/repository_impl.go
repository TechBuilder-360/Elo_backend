package authentication

import (
	"context"

	"github.com/Toflex/directory_v2/ent"
	"github.com/Toflex/directory_v2/ent/user"
	"github.com/Toflex/directory_v2/pkg/utils"
)

func (r *repository) Create(ctx context.Context, payload Onboarding) (*string, error) {
	u, err := r.db.DBClient.User.
		Create().
		SetFirstName(payload.FirstName).
		SetLastName(payload.LastName).
		SetAvatar(utils.AddressToString(payload.Avatar)).
		SetEmailAddress(payload.EmailAddress).
		SetEmailVerifiedAt(payload.EmailVerifiedAt).
		SetEmailVerified(payload.EmailVerified).
		SetDisplayName(utils.AddressToString(payload.DisplayName)).
		SetPhoneNumber(utils.AddressToString(payload.PhoneNumber)).
		Save(ctx)

	if err != nil {
		return nil, err
	}

	return &u.ID, nil
}

// GetUserByEmail implements IRepository.
func (r *repository) GetUserByEmail(ctx context.Context, email string) (*User, error) {
	u, err := r.db.DBClient.User.
		Query().
		Where(user.EmailAddressEQ(email)).
		Only(ctx)

	if err != nil {
		if ent.IsNotFound(err) {
			return nil, nil
		}
		return nil, err
	}

	return &User{*u}, nil
}

// GetUserById implements IRepository.
func (r *repository) GetUserByID(ctx context.Context, id string) (*User, error) {
	u, err := r.db.DBClient.User.
		Query().
		Where(user.IDEQ(id)).
		Only(ctx)

	if err != nil {
		return nil, err
	}

	return &User{*u}, nil
}
