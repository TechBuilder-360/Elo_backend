package authentication

import (
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/Toflex/directory_v2/ent"
	"github.com/Toflex/directory_v2/pkg/errors"
	"github.com/Toflex/directory_v2/pkg/log"
	"github.com/Toflex/directory_v2/pkg/util"
)

type Registration struct {
	DisplayName  *string `json:"display_name,omitempty"`
	EmailAddress string  `json:"email_address"`
	FirstName    string  `json:"first_name"`
	LastName     string  `json:"last_name"`
	Password     string  `json:"password"`
	PhoneNumber  *string `json:"phone_number,omitempty"`
	Avatar       *string `json:"avatar,omitempty"`
}

type Onboarding struct {
	DisplayName     *string
	EmailAddress    string
	FirstName       string
	LastName        string
	Password        string
	PhoneNumber     *string
	Avatar          *string
	EmailVerifiedAt time.Time
	EmailVerified   bool
}

func (r *Registration) Validate() error {
	re := regexp.MustCompile(`[^a-zA-Z0-9-]`)

	// validate display name field
	displayName := util.AddressToString(r.DisplayName)
	if len(displayName) == 0 {
		displayName = fmt.Sprintf("%s_%s", r.FirstName, util.RandomString(4))
	}

	if len(displayName) > 10 {
		displayName = displayName[:10]
	}

	displayName = strings.ToLower(displayName)

	// validate email field
	err := util.ValidateEmail(r.EmailAddress)
	if err != nil {
		return errors.New(errors.ErrInvalidInput, err.Error())
	}

	r.EmailAddress = strings.ToLower(r.EmailAddress)

	// validate first name
	firstname := strings.ReplaceAll(strings.ToLower(r.FirstName), " ", "")
	if re.MatchString(firstname) {
		return errors.New(errors.ErrInvalidInput, "First name contains special character")
	}

	// validate last name
	lastname := strings.ReplaceAll(strings.ToLower(r.LastName), " ", "")
	if re.MatchString(lastname) {
		errors.New(errors.ErrInvalidInput, "Last name contains special character")
	}

	// validate phone number
	if r.PhoneNumber == nil {
		err = util.IsValidPhoneNumber(util.AddressToString(r.PhoneNumber))
		if err != nil {
			return errors.New(errors.ErrInvalidInput, err.Error())
		}
	}

	// validate avatar
	if r.Avatar == nil {
		avatar := strings.ReplaceAll(strings.ToLower(*r.Avatar), " ", "")
		r.Avatar = &avatar
	}

	// password validation
	isValid, err := util.PasswordStrength(r.Password)
	if err != nil {
		log.WithError(err).WithField("email", r.EmailAddress).Error("error: occurred while validating password strength.")
		return errors.New(errors.ErrFailed, "invalid password")
	}

	if !isValid {
		log.WithField("email", r.EmailAddress).Error("error: password does not meet strength requirements.")
		return errors.New(errors.ErrFailed, "password does not meet strength requirements")
	}

	r.DisplayName = &displayName
	r.FirstName = util.ToTitleCase(firstname)
	r.LastName = util.ToTitleCase(lastname)
	r.Password = util.EncryptPassword(r.Password)

	return nil
}

type User struct {
	ent.User
}

type Login struct {
	Otp        string
	Identifier string
}

type CacheOTP struct {
	UserID       string
	EncryptedOTP string
	Attempts     uint
}

type JWTToken struct {
	UserID string
}

type OTPRequest struct {
	Email    string
	Password string
}

func (r *OTPRequest) Validate() error {
	// validate email field
	email := r.Email
	err := util.ValidateEmail(email)
	if err != nil {
		return errors.New(errors.ErrInvalidInput, err.Error())
	}

	r.Email = strings.ToLower(email)

	// password validation
	r.Password = util.EncryptPassword(r.Password)

	return nil
}
