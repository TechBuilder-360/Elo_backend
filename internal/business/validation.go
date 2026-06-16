package business

import (
	"errors"
	"fmt"
	"net/http"
	"path/filepath"
	"strings"
	"time"

	"github.com/99designs/gqlgen/graphql"
	"github.com/Toflex/directory_v2/pkg/util"
)

func (b *CreateBusinessRequest) Validate() error {
	if b.Name == "" {
		return errors.New("business name is required")
	}

	b.Name = util.ToTitleCase(b.Name)

	// validate business email
	if err := util.ValidateEmail(b.Email); err != nil {
		return err
	}

	if b.IsRegistered {
		return b.RegistrationDetail.Validate()
	}

	for _, document := range b.OtherDocument {
		return document.Validate()
	}

	return nil
}

func (d *Document) Validate() error {
	err := validateUpload(d.File)
	if err != nil {
		return fmt.Errorf("%s: %s", d.Description, err.Error())
	}

	return nil
}

func (b *BusinessRegistrationDetail) Validate() error {
	if b.Number == "" {
		return errors.New("business registration number is required")
	}

	b.Number = strings.ToUpper(b.Number)

	if b.CountryOfIncorporation == "" {
		return errors.New("business country of incorporation is required")
	}

	_, err := time.Parse("02-01-2006", b.DateOfIncorporation)
	if err != nil {
		return (err)
	}

	err = validateUpload(b.CertificateOfIncorporation)
	if err != nil {
		return fmt.Errorf("Certificate of Incorporation: %s", err.Error())
	}

	err = validateUpload(b.ArticlesOfAssociation)
	if err != nil {
		return fmt.Errorf("Articles of Association: %s", err.Error())
	}

	err = validateUpload(b.StatusCertificate)
	if err != nil {
		return fmt.Errorf("Status Certificate: %s", err.Error())
	}

	return nil
}

func validateUpload(upload graphql.Upload) error {
	if upload.File == nil {
		return errors.New("file is required")
	}

	// Max size: 2MB
	const maxSize = 2 << 20

	if upload.Size > maxSize {
		return errors.New("file too large")
	}

	// Extension validation
	allowedExt := map[string]bool{
		".jpg":  true,
		".jpeg": true,
		".pdf":  true,
	}

	ext := strings.ToLower(
		filepath.Ext(upload.Filename),
	)

	if !allowedExt[ext] {
		return errors.New("invalid file extension")
	}

	// MIME validation
	buffer := make([]byte, 512)

	_, err := upload.File.Read(buffer)

	if err != nil {
		return err
	}

	// reset file reader
	upload.File.Seek(0, 0)

	contentType := http.DetectContentType(buffer)

	allowedMime := map[string]bool{
		"image/jpeg":      true,
		"image/jpg":       true,
		"application/pdf": true,
	}

	if !allowedMime[contentType] {
		return errors.New("invalid file content")
	}

	return nil
}
