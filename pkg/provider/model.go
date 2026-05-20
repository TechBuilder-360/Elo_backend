package provider

type FeeType string

const (
	TierFeeType       FeeType = "TIER"
	PercentageFeeType FeeType = "PERCENTAGE"
	FlatFeeType       FeeType = "FLAT"
)

type Tier struct {
	From  uint
	To    uint
	Type  FeeType
	Value float64
}

type ActiveProvider struct {
	Name   string
	Slug   string
	Active bool
}

type Fee struct {
	Type  FeeType
	Value float64
	Min   int
	Max   int
	Tiers []Tier
}

type ServiceLocator struct {
	Name           string
	Identifier     string
	ActiveProvider ActiveProvider
	Fee            Fee
	MinValue       int64
	MaxValue       int64
}
