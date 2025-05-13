package model

type Business struct {
	ID                 string
	Category           string
	CountryID          string
	Name               string
	LogoURL            *string
	PhoneNumber        *string
	SupportPhoneNumber *string
	EmailAddress       string
	Website            *string
}
