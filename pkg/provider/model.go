package provider

type FeeComputationType string

const (
	TieredFeeType     FeeComputationType = "TIERED"
	PercentageFeeType FeeComputationType = "PERCENTAGE"
	FlatFeeType       FeeComputationType = "FLAT"
)

type Tier struct {
	From            int64
	To              *int64
	ComputationType FeeComputationType
	FlatAmount      int64
	PercentageRate  float64
}

type ActiveProvider struct {
	Name   string
	Slug   string
	Active bool
}

type Fee struct {
	ComputationType FeeComputationType
	FeeType         string
	FlatAmount      int64
	PercentageRate  float64
	Min             *int64
	Max             *int64
	MinFee          *int64
	MaxFee          *int64
	Tiered          []Tier
}

// AppliesTo reports whether the fee rule applies to the given transaction amount.
func (f Fee) AppliesTo(amount int64) bool {
	if f.Min != nil && amount < *f.Min {
		return false
	}
	if f.Max != nil && amount > *f.Max {
		return false
	}
	return true
}

// ClampFee applies optional min_fee and max_fee guardrails to the computed fee.
func (f Fee) ClampFee(fee int64) int64 {
	if f.MinFee != nil && fee < *f.MinFee {
		fee = *f.MinFee
	}
	if f.MaxFee != nil && fee > *f.MaxFee {
		fee = *f.MaxFee
	}
	return fee
}

type ServiceLocator struct {
	Name           string
	Identifier     string
	ActiveProvider ActiveProvider
	Fees           []Fee
	MinValue       int64
	MaxValue       int64
}

type ServiceFee struct {
	FeeType string
	Fee     int64
	Rate    *float64
}
