package authentication

import (
	"regexp"
	"strings"
	"time"

	"github.com/Toflex/directory_v2/ent"
	"github.com/Toflex/directory_v2/pkg/errors"
	"github.com/Toflex/directory_v2/pkg/utils"
)

type Registration struct {
	DisplayName  *string `json:"display_name,omitempty"`
	EmailAddress string  `json:"email_address"`
	FirstName    string  `json:"first_name"`
	LastName     string  `json:"last_name"`
	PhoneNumber  *string `json:"phone_number,omitempty"`
	Avatar       *string `json:"avatar,omitempty"`
}

type Onboarding struct {
	DisplayName     *string
	EmailAddress    string
	FirstName       string
	LastName        string
	PhoneNumber     *string
	Avatar          *string
	EmailVerifiedAt time.Time
	EmailVerified   bool
}

func (r *Registration) Validate() error {
	re := regexp.MustCompile(`[^a-zA-Z0-9-]`)

	// validate display name field
	disaplayName := utils.AddressToString(r.DisplayName)
	if len(disaplayName) == 0 {
		disaplayName = r.FirstName
	}

	if len(disaplayName) > 10 {
		disaplayName = disaplayName[:10]
	}

	// validate email field
	email := r.EmailAddress
	err := utils.ValidateEmail(email)
	if err != nil {
		return errors.New(errors.ErrInvalidInput, err.Error())
	}

	email = strings.ToLower(email)

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
		err = utils.IsValidPhoneNumber(utils.AddressToString(r.PhoneNumber))
		if err != nil {
			return errors.New(errors.ErrInvalidInput, err.Error())
		}
	}

	// validate avatar
	if r.Avatar == nil {
		avatar := strings.ReplaceAll(strings.ToLower(*r.Avatar), " ", "")
		r.Avatar = &avatar
	}

	r.DisplayName = &disaplayName
	r.EmailAddress = email
	r.FirstName = firstname
	r.LastName = lastname

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
