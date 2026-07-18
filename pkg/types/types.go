package types

import "strings"

type Provider string
type JWTKey string
type EntityType string
type VerificationStatus string
type DocumentType string
type CurrencyCode string
type ToMinor struct {
	Amount    float64
	Precision uint
}

type ToMajor struct {
	Amount    int64
	Precision uint
}

func (p Provider) ToString() string {
	return string(p)
}

func (c CurrencyCode) ToString() string {
	return string(c)
}

func (c CurrencyCode) Capitalize() CurrencyCode {
	c = CurrencyCode(strings.ToUpper(string(c)))
	return c
}
