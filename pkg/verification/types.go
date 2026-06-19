package verification

import "github.com/Toflex/directory_v2/pkg/types"

type VerificationResult struct {
	Provider       types.Provider  `json:"provider"`
	ReferenceID    string          `json:"reference_id"`
	Metadata       interface{}     `json:"meta_map"`
	BVN            *BVN            `json:"bvn"`
	NationalID     *NationalID     `json:"nin"`
	Passport       *Passport       `json:"passport"`
	DriversLicense *DriversLicense `json:"drivers_license"`
	VoterID        *VoterID        `json:"voter_id"`
}

type VerificationRequest struct {
	Entity types.EntityType `json:"entity"`
	ID     string           `json:"id"`
}

type VerificationResponse struct {
	URL    string `json:"url"`
	Status string `json:"status"`
}

type userDetails struct {
	FirstName                    string `json:"first_name"`
	LastName                     string `json:"last_name"`
	MiddleName                   string `json:"middle_name"`
	Number                       string `json:"number"`
	PhoneNumber                  string `json:"phone_number"`
	Nationality                  string `json:"nationality"`
	NationalIdentificationNumber string `json:"national_identification_number"`
}

type VoterID struct {
	FirstName   string `json:"first_name"`
	LastName    string `json:"last_name"`
	Number      string `json:"number"`
	Expiration  string `json:"expiration"`
	IssueDate   string `json:"issue_date"`
	Nationality string `json:"nationality"`
	DateOfBirth string `json:"date_of_birth"`
	FrontImage  string `json:"front_image"`
	BackImage   string `json:"back_image"`
}

type NationalID struct {
	FirstName   string `json:"first_name"`
	LastName    string `json:"last_name"`
	Number      string `json:"number"`
	Nationality string `json:"nationality"`
	DateOfBirth string `json:"date_of_birth"`
	FrontImage  string `json:"front_image"`
	BackImage   string `json:"back_image"`
}

type BVN struct {
	FirstName                    string `json:"first_name"`
	LastName                     string `json:"last_name"`
	MiddleName                   string `json:"middle_name"`
	Number                       string `json:"number"`
	PhoneNumber                  string `json:"phone_number"`
	Nationality                  string `json:"nationality"`
	NationalIdentificationNumber string `json:"national_identification_number"`
	DateOfBirth                  string `json:"date_of_birth"`
}

type DriversLicense struct {
	FirstName   string `json:"first_name"`
	LastName    string `json:"last_name"`
	Number      string `json:"number"`
	Expiration  string `json:"expiration"`
	IssueDate   string `json:"issue_date"`
	Nationality string `json:"nationality"`
	DateOfBirth string `json:"date_of_birth"`
	FrontImage  string `json:"front_image"`
	BackImage   string `json:"back_image"`
}

type Passport struct {
	FirstName   string `json:"first_name"`
	LastName    string `json:"last_name"`
	Number      string `json:"number"`
	Expiration  string `json:"expiration"`
	IssueDate   string `json:"issue_date"`
	Nationality string `json:"nationality"`
	DateOfBirth string `json:"date_of_birth"`
	FrontImage  string `json:"front_image"`
	BackImage   string `json:"back_image"`
}
