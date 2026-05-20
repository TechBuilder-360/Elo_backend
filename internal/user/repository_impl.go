package user

import (
	"context"

	"github.com/Toflex/directory_v2/ent"
	"github.com/Toflex/directory_v2/ent/user"
)

// GetUserByID implements IUserRepository.
func (ur *repository) GetByID(ctx context.Context, id string) (*User, error) {
	result, err := ur.db.
		User.
		Query().
		Where(user.IDEQ(id)).
		First(ctx)

	if err != nil {
		if ent.IsNotFound(err) {
			return nil, nil
		}

		return nil, err
	}

	// Convert ent.User to internal User model
	u := &User{
		User: *result,
	}

	return u, nil
}

// GetUserByEmail implements IUserRepository.
func (ur *repository) GetByEmail(ctx context.Context, email string) (*User, error) {
	var users []*ent.User

	err := ur.db.
		User.
		Query().
		Where(user.EmailAddressEQ(email)).
		Select(user.FieldID, user.FieldFirstName, user.FieldLastName, user.FieldEmailAddress).
		Scan(ctx, &users)

	if err != nil {
		return nil, err
	}

	if len(users) == 0 {
		return nil, nil
	}

	// Convert ent.User to internal User model
	u := &User{
		User: *users[0],
	}

	return u, nil
}
