package util

import (
	"crypto/rand"
	"encoding/json"
	"errors"
	"fmt"
	"math/big"
	"net/url"
	"regexp"
	"strconv"
	"strings"

	"github.com/google/uuid"
	"github.com/lucsky/cuid"
	"golang.org/x/crypto/bcrypt"
)

func GenerateCUID() string {
	return cuid.New()
}

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

func GenerateReference(prefix string) string {
	return fmt.Sprintf("%s-%s", prefix, GenerateCUID())
}

func AddressToString(txt *string) string {
	if txt == nil {
		return ""
	}

	return *txt
}

func EncryptPassword(text string) string {
	bytes, _ := bcrypt.GenerateFromPassword([]byte(text), bcrypt.DefaultCost)
	return string(bytes)
}

func DecryptPassword(text, hash string) bool {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(text)) == nil
}

func GenerateOTP() string {
	defaultOTP := "748290"

	n, err := rand.Int(rand.Reader, big.NewInt(900000)) // 0–899999
	if err != nil {
		return defaultOTP
	}

	otp := int(n.Int64()) + 100000
	return strconv.Itoa(otp)
}

func ToJSON(v interface{}) string {
	json, err := json.Marshal(v)
	if err != nil {
		return ""
	}
	return string(json)
}

func ValidateURL(link string) error {
	u, err := url.ParseRequestURI(link)
	if err != nil {
		return errors.New("url is invalid!")
	}

	if u.Scheme == "" || u.Host == "" {
		return errors.New("url is invalid!")
	}

	return nil
}

func ToTitleCase(txt string) string {
	return strings.Title(strings.ToLower(txt))
}

func RandomString(n int) string {
	const letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	result := make([]byte, n)
	for i := range result {
		num, err := rand.Int(rand.Reader, big.NewInt(int64(len(letters))))
		if err != nil {
			return ""
		}
		result[i] = letters[num.Int64()]
	}
	return string(result)
}

func PasswordStrength(password string) (bool, error) {
	var (
		hasMinLen  = false
		hasUpper   = false
		hasLower   = false
		hasNumber  = false
		hasSpecial = false
	)

	if len(password) >= 8 {
		hasMinLen = true
	}

	for _, char := range password {
		switch {
		case 'A' <= char && char <= 'Z':
			hasUpper = true
		case 'a' <= char && char <= 'z':
			hasLower = true
		case '0' <= char && char <= '9':
			hasNumber = true
		case strings.ContainsRune("!@#$%^&*()-_=+[]{}|;:'\",.<>?/", char):
			hasSpecial = true
		}
	}

	if !hasMinLen {
		return false, errors.New("password must be at least 8 characters long")
	}
	if !hasUpper {
		return false, errors.New("password must contain at least one uppercase letter")
	}
	if !hasLower {
		return false, errors.New("password must contain at least one lowercase letter")
	}
	if !hasNumber {
		return false, errors.New("password must contain at least one number")
	}
	if !hasSpecial {
		return false, errors.New("password must contain at least one special character")
	}

	return true, nil
}
