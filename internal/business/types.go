package business

import "github.com/Toflex/directory_v2/ent"

type createbusiness struct {
	Name     string
	Category string
	About    string
	// CountryOfIncorporation *string
	// RegistrationNumber     *string
	Email        string
	RegisteredBy string
	OnSite       bool
}

type createBusinssLocation struct {
	Name         string
	Street       string
	City         string
	State        string
	Country      string
	ZipCode      string
	Latitude     *float64
	Longitude    *float64
	IsHeadOffice bool
	Business     *ent.Business
}

type businessUpload struct {
	Business           *ent.Business
	RegistrationDetail *RegistrationDetail `json:"registration_detail,omitempty"`
	OtherDocument      []Document          `json:"other_document,omitempty"`
}

type AddBusinessDocument struct {
	Title       string
	Description string
	URL         string
	Type        string
	DocumentID  string
}
