package business

import (
	"github.com/99designs/gqlgen/graphql"
	"github.com/Toflex/directory_v2/ent"
)

type CreateBusinessRequest struct {
	User               *ent.User                   `json:"user"`
	Name               string                      `json:"name"`
	About              string                      `json:"about"`
	Email              string                      `json:"email"`
	OnSite             bool                        `json:"on_site,omitempty"`
	Industry           string                      `json:"industry"`
	IsRegistered       bool                        `json:"is_registered,omitempty"`
	Address            *BusinessAddress            `json:"address"`
	RegistrationDetail *BusinessRegistrationDetail `json:"registration_detail,omitempty"`
	OtherDocument      []*Document                 `json:"other_document,omitempty"`
}

type BusinessAddress struct {
	Number  string `json:"number,omitempty"`
	Street  string `json:"street"`
	State   string `json:"state"`
	Country string `json:"country"`
	ZipCode string `json:"zip_code"`
}

type BusinessRegistrationDetail struct {
	Number                     string         `json:"number"`
	CountryOfIncorporation     string         `json:"country_of_incorporation"`
	DateOfIncorporation        string         `json:"date_of_incorporation"`
	CertificateOfIncorporation graphql.Upload `json:"certificate_of_incorporation"`
	ArticlesOfAssociation      graphql.Upload `json:"articles_of_association"`
	StatusCertificate          graphql.Upload `json:"status_certificate"`
}


type Document struct {
	Description string         `json:"description"`
	File        graphql.Upload `json:"file"`
}

type businessResult struct {
	ent.Business
}
