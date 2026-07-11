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
	"github.com/pariz/gountries"
)

var g *gountries.Query

func init() {
	g = gountries.New()
}

func (b *CreateBusinessRequest) Validate() error {
	if b.Name == "" {
		return errors.New("business name is required")
	}

	b.Name = util.ToTitleCase(b.Name)

	// validate business email
	if err := util.ValidateEmail(b.Email); err != nil {
		return err
	}

	if b.Role.AuthorizedRepresentative {
		if b.Role.AuthorizedRepresentativeEmail == nil {
			return errors.New("authorized representative email is required")
		}

		email := util.AddressToString(b.Role.AuthorizedRepresentativeEmail)
		err := util.ValidateEmail(email)
		if err != nil {
			return err
		}

		email = strings.ToLower(email)
		b.Role.AuthorizedRepresentativeEmail = &email
	}

	if err := b.Address.Validate(); err != nil {
		return err
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

func (b *RegistrationDetail) Validate() error {
	if b == nil {
		return nil
	}

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

	return nil
}

func (ba *BusinessAddress) Validate() error {

	if ba == nil {
		return nil
	}

	if strings.TrimSpace(ba.Street) == "" {
		return errors.New("street address is required")
	}

	if strings.TrimSpace(ba.City) == "" {
		return errors.New("city is required")
	}

	if strings.TrimSpace(ba.State) == "" {
		return errors.New("state is required")
	}

	if strings.TrimSpace(ba.Country) == "" {
		return errors.New("country is required")
	}

	country, err := g.FindCountryByAlpha(ba.Country)
	if err != nil {
		return errors.New("country not valid")
	}

	ba.Country = country.Alpha2

	if strings.TrimSpace(ba.ZipCode) == "" {
		return errors.New("zip code is required")
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

func (b *BusinessDetailRequest) Validate() error {
	var err error
	if b.Name != nil {
		name := (util.ToTitleCase(*b.Name))
		b.Name = &name
	}

	if b.Website != nil {
		err = util.ValidateURL(*b.Website)
		if err != nil {
			return fmt.Errorf("Website %s", err.Error())
		}
	}

	err = b.RegistrationDetail.Validate()
	if err != nil {
		return err
	}

	return nil
}
