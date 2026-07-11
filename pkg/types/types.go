package types

type Provider string
type JWTKey string
type EntityType string
type VerificationStatus string
type DocumentType string

func (p Provider) ToString() string {
	return string(p)
}
