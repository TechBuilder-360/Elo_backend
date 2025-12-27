package types

type Partner string
type JWTKey string

func (p Partner) ToString() string {
	return string(p)
}
