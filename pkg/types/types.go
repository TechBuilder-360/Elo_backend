package types

type Provider string
type JWTKey string
type EntityType string

func (p Provider) ToString() string {
	return string(p)
}
