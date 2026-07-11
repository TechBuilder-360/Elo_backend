package util

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"math/big"
	r "math/rand/v2"
	"net/http"
	"net/url"
	"regexp"
	"strconv"
	"strings"

	"github.com/99designs/gqlgen/graphql"
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

func RandomInt() int {
	return r.IntN(193857-10+1) + 667
}

func ValidateEmail(email string) error {
	email = strings.Trim(email, " ")
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

func Contains(arr []string, s string) bool {
	for _, v := range arr {
		if strings.EqualFold(v, s) {
			return true
		}
	}

	return false
}

func URLToBase64(url string) (string, error) {
	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("failed to download: %s", resp.Status)
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return base64.RawStdEncoding.EncodeToString(data), nil
}

func UploadToBase64(upload graphql.Upload) (string, error) {
	data, err := io.ReadAll(upload.File)
	if err != nil {
		return "", err
	}

	return base64.StdEncoding.EncodeToString(data), nil
}

func ValidateDataURI(value string, extension []string) error {
	var dataURLRegex = regexp.MustCompile(
		`^data:([a-zA-Z0-9!#$&^_.+-]+/[a-zA-Z0-9!#$&^_.+-]+);base64,[A-Za-z0-9+/]+={0,2}$`,
	)

	if !dataURLRegex.MatchString(value) {
		return errors.New("invalid data URI")
	}

	return VerifyDocumentType(value, extension)
}

func VerifyDocumentType(data string, extension []string) error {
	ext, err := ExtensionFromDataURI(data)
	if err != nil {
		return err
	}

	if Contains(extension, ext) {
		return nil
	}

	return errors.New("unsupprted file type")
}

func ExtensionFromDataURI(dataURI string) (string, error) {
	if !strings.HasPrefix(dataURI, "data:") {
		return "", fmt.Errorf("invalid data URI")
	}

	header, _, found := strings.Cut(dataURI, ",")
	if !found {
		return "", fmt.Errorf("invalid data URI")
	}

	mimeType := strings.TrimPrefix(header, "data:")
	mimeType = strings.TrimSuffix(mimeType, ";base64")

	idx := strings.Index(mimeType, "/")

	return mimeType[idx+1:], nil
}
