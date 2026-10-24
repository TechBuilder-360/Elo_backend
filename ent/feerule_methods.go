package ent

// AppliesTo determines if the fee rule should apply based on the transaction amount.
// Both bounds are inclusive:
// - If MinValue != nil, amount must be >= *MinValue.
// - If MaxValue != nil, amount must be <= *MaxValue.
func (fr *FeeRule) AppliesTo(amount int64) bool {
	if fr == nil {
		return false
	}
	if fr.MinValue != nil && amount < *fr.MinValue {
		return false
	}
	if fr.MaxValue != nil && amount > *fr.MaxValue {
		return false
	}
	return true
}

// ClampFee applies optional min_fee (floor) and max_fee (cap) guardrails to the computed fee.
func (fr *FeeRule) ClampFee(fee int64) int64 {
	if fr == nil {
		return fee
	}
	if fr.MinFee != nil && fee < *fr.MinFee {
		fee = *fr.MinFee
	}
	if fr.MaxFee != nil && fee > *fr.MaxFee {
		fee = *fr.MaxFee
	}
	return fee
}
