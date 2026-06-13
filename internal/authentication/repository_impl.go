package authentication

import (
	"context"
	"strings"

	"github.com/Toflex/directory_v2/ent"
	"github.com/Toflex/directory_v2/ent/user"
	"github.com/Toflex/directory_v2/pkg/util"
)

func (r *repository) Create(ctx context.Context, payload Onboarding) (*string, error) {
	u, err := r.db.User.
		Create().
		SetFirstName(payload.FirstName).
		SetLastName(payload.LastName).
		SetAvatar(util.AddressToString(payload.Avatar)).
		SetEmailAddress(payload.EmailAddress).
		SetEmailVerifiedAt(payload.EmailVerifiedAt).
		SetEmailVerified(payload.EmailVerified).
		SetDisplayName(util.AddressToString(payload.DisplayName)).
		SetPhoneNumber(util.AddressToString(payload.PhoneNumber)).
		SetPassword(payload.Password).
		Save(ctx)

	if err != nil {
		return nil, err
	}

	return &u.ID, nil
}

// GetUserByEmail implements IRepository.
func (r *repository) GetUserByEmail(ctx context.Context, email string) (*User, error) {
	u, err := r.db.User.
		Query().
		Where(user.EmailAddressEQ(strings.ToLower(email))).
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
	u, err := r.db.User.
		Query().
		Where(user.IDEQ(id)).
		Only(ctx)

	if err != nil {
		return nil, err
	}

	return &User{*u}, nil
}
