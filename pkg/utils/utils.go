package utils

import (
	"crypto/rand"
	"errors"
	"fmt"
	"math/big"
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

func Encrypt(text string) string {
	bytes, _ := bcrypt.GenerateFromPassword([]byte(text), bcrypt.DefaultCost)
	return string(bytes)
}

func Decrypt(text, hash string) bool {
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
