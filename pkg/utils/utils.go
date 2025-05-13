package utils

import (
	"errors"
	"github.com/google/uuid"
	"regexp"
	"strings"
)

func GenerateUUID() string {
	uniqueID, _ := uuid.NewUUID()
	return uniqueID.String()
}

func ValidateEmail(email string) error {
	emailRegex := `^[a-zA-Z0-9.!#$%&'*+/=?^_` + "`" + `{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$`
	// Compile the regex
	re := regexp.MustCompile(emailRegex)

	if re.MatchString(email) {
		return nil
	}

	return errors.New("invalid email address")
}

func IsValidPhoneNumber(phoneNumber string) error {
	e164Regex := `^\+[1-9]\d{1,14}$`
	re := regexp.MustCompile(e164Regex)
	phoneNumber = strings.ReplaceAll(phoneNumber, " ", "")

	if re.Find([]byte(phoneNumber)) != nil {
		return nil
	}

	return errors.New("invalid phone number")
}
