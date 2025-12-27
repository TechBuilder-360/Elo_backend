package model

import "time"

type User struct {
	DisplayName     *string `json:"display_name,omitempty"`
	EmailAddress    string  `json:"email_address"`
	FirstName       string  `json:"first_name"`
	LastName        string  `json:"last_name"`
	PhoneNumber     *string `json:"phone_number,omitempty"`
	Avatar          *string `json:"avatar,omitempty"`
	EmailVerified   bool
	EmailVerifiedAt time.Time
	Active          bool
}
