package authentication

import (
	"context"
	"fmt"

	"github.com/Toflex/directory_v2/internal/email"
	"github.com/Toflex/directory_v2/pkg/errors"
	"github.com/Toflex/directory_v2/pkg/log"
)

func (s *Service) RegisterUser(ctx context.Context, body Registration, log log.Entry) (*string, error) {
	// Check if email address exist
	existingUser, err := s.repo.GetUserByEmail(ctx, body.EmailAddress)
	if err != nil {
		log.Error(err.Error())
		return nil, errors.New(errors.ErrFailed, err.Error())
	}

	if existingUser != nil {
		log.WithField("email", body.EmailAddress).Info("Email address already registered.")
		return nil, errors.New(errors.ErrFailed, "account already exist")
	}

	onboarding := Onboarding{
		DisplayName:  body.DisplayName,
		EmailAddress: body.EmailAddress,
		FirstName:    body.FirstName,
		LastName:     body.LastName,
		PhoneNumber:  body.PhoneNumber,
		Avatar:       body.Avatar,
		Password:     body.Password,
	}

	uid, err := s.repo.Create(ctx, onboarding)
	if err != nil {
		log.WithError(err).WithField("email", body.EmailAddress).Error("error: occurred when saving new user.")
		return nil, errors.New(errors.ErrFailed, "registration was not successful")
	}

	//send welcome email
	if err := email.NewEmailWelcomeTask(email.WelcomeEmailPayload{
		ToMail:   body.EmailAddress,
		ToName:   fmt.Sprintf("%s %s", body.FirstName, body.LastName),
		FullName: fmt.Sprintf("%s %s", body.FirstName, body.LastName),
	}); err != nil {
		log.WithError(err).Error("failed to enqueue welcome email task")
	}

	return uid, nil
}
