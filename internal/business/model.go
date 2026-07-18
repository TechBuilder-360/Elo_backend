package business

import (
	"github.com/99designs/gqlgen/graphql"
	"github.com/Toflex/directory_v2/ent"
)

type CreateBusinessRequest struct {
	User     *ent.User `json:"user"`
	Name     string    `json:"name"`
	About    string    `json:"about"`
	Email    string    `json:"email"`
	OnSite   bool      `json:"on_site,omitempty"`
	Industry string    `json:"industry"`
	// IsRegistered       bool                        `json:"is_registered,omitempty"`
	Address BusinessAddress `json:"address"`
	// RegistrationDetail *BusinessRegistrationDetail `json:"registration_detail,omitempty"`
	// OtherDocument      []Document                  `json:"other_document,omitempty"`
	Role Role `json:"role"`
}

type Role struct {
	AuthorizedRepresentative      bool    `json:"authorized_representative"`
	AuthorizedRepresentativeEmail *string `json:"authorized_representative_email,omitempty"`
}

type BusinessAddress struct {
	Number  string `json:"number,omitempty"`
	City    string `json:"city"`
	Street  string `json:"street"`
	State   string `json:"state"`
	Country string `json:"country"`
	ZipCode string `json:"zip_code"`
}

type RegistrationDetail struct {
	Number                 string `json:"number"`
	CountryOfIncorporation string `json:"country_of_incorporation"`
	DateOfIncorporation    string `json:"date_of_incorporation"`
}

type Document struct {
	Description string         `json:"description"`
	File        graphql.Upload `json:"file"`
}

type businessResult struct {
	ent.Business
}

type UploadDocumentRequest struct {
	User        *ent.User     `json:"user"`
	Business    *ent.Business `json:"business"`
	DocumentID  string        `json:"document_id"`
	Description string        `json:"description"`
	File        string        `json:"file"`
}

type BusinessDetailRequest struct {
	RegistrationDetail *RegistrationDetail `json:"registration_detail,omitempty"`
	Name               *string             `json:"name,omitempty"`
	About              *string             `json:"about,omitempty"`
	Industry           *string             `json:"industry,omitempty"`
	Website            *string             `json:"website,omitempty"`
}

type MyBusinessResult struct {
	ID   string  `json:"id"`
	Name string  `json:"name"`
	Logo *string `json:"logo,omitempty"`
	Role *string `json:"role,omitempty"`
}

type BusinessResult struct {
	ID                      string   `json:"id"`
	Name                    string   `json:"name"`
	OnSite                  bool     `json:"on_site"`
	Logo                    *string  `json:"logo,omitempty"`
	Email                   *string  `json:"email,omitempty"`
	About                   *string  `json:"about,omitempty"`
	Services                []string `json:"services"`
	Industry                string   `json:"industry"`
	Number                  string   `json:"number"`
	CountryOfIncorporation  string   `json:"country_of_incorporation"`
	DateOfIncorporation     string   `json:"date_of_incorporation"`
	TaxIdentificationNumber string   `json:"tax_identification_number"`
	Address                 BusinessAddress
	// Socials  []*Social `json:"socials"`
}

type DocumentResult struct {
	ID       string
	Name     string
	Required bool
}

type KYBDocument struct {
	ID          string `json:"id"`
	DocumentID  string `json:"document_id"`
	Description string `json:"description"`
	File        string `json:"file"`
}

type UpdateBusiness struct {
	Name                   *string
	About                  *string
	Industry               *string
	Website                *string
	Number                 *string
	CountryOfIncorporation *string
	DateOfIncorporation    *string
}

type CategoryRequest struct {
	Name string
}
