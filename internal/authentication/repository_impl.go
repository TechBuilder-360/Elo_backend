package authentication

import (
	"context"
	"strings"
	"time"

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
		SetNillableDisplayName(payload.DisplayName).
		SetNillablePhoneNumber(payload.PhoneNumber).
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

func (r *repository) Update(ctx context.Context, user *User, opt *UpdateUser) error {
	urs := r.db.User.UpdateOneID(user.ID)

	if opt.EmailVerified {
		urs.SetEmailVerified(true).
			SetEmailVerifiedAt(time.Now())
	}

	if opt.Avartar {
		urs.SetAvatar(*user.Avatar)
	}

	if opt.Disable {
		urs.SetDisabled(user.Disabled).
			SetDisableReason(*user.DisableReason)
	}

	if opt.Password {
		urs.SetPassword(user.Password)
	}

	if opt.PersonalInfo {
		urs.SetFirstName(user.FirstName).
			SetLastName(user.LastName).
			SetMiddleName(user.MiddleName).
			SetPhoneNumber(*user.PhoneNumber)
	}

	if opt.Verified {
		urs.SetVerified(user.Verified)
	}

	if opt.Avartar {
		urs.SetAvatar(*user.Avatar)
	}

	_, err := urs.Save(ctx)

	return err
}
