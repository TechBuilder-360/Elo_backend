package ent

import "testing"

func int64Ptr(v int64) *int64 {
	return &v
}

func TestFeeRule_AppliesTo(t *testing.T) {
	tests := []struct {
		name     string
		rule     *FeeRule
		amount   int64
		expected bool
	}{
		{
			name:     "nil rule returns false",
			rule:     nil,
			amount:   1000,
			expected: false,
		},
		{
			name: "unbounded rule applies to any amount",
			rule: &FeeRule{
				MinValue: nil,
				MaxValue: nil,
			},
			amount:   5000,
			expected: true,
		},
		{
			name: "amount within min and max",
			rule: &FeeRule{
				MinValue: int64Ptr(1000),
				MaxValue: int64Ptr(10000),
			},
			amount:   5000,
			expected: true,
		},
		{
			name: "amount exactly at min",
			rule: &FeeRule{
				MinValue: int64Ptr(1000),
				MaxValue: int64Ptr(10000),
			},
			amount:   1000,
			expected: true,
		},
		{
			name: "amount exactly at max",
			rule: &FeeRule{
				MinValue: int64Ptr(1000),
				MaxValue: int64Ptr(10000),
			},
			amount:   1000,
			expected: true,
		},
		{
			name: "amount below min",
			rule: &FeeRule{
				MinValue: int64Ptr(1000),
				MaxValue: int64Ptr(10000),
			},
			amount:   999,
			expected: false,
		},
		{
			name: "amount above max",
			rule: &FeeRule{
				MinValue: int64Ptr(1000),
				MaxValue: int64Ptr(10000),
			},
			amount:   10001,
			expected: false,
		},
		{
			name: "only min bound set - valid",
			rule: &FeeRule{
				MinValue: int64Ptr(5000),
				MaxValue: nil,
			},
			amount:   5000,
			expected: true,
		},
		{
			name: "only min bound set - below min",
			rule: &FeeRule{
				MinValue: int64Ptr(5000),
				MaxValue: nil,
			},
			amount:   4999,
			expected: false,
		},
		{
			name: "only max bound set - valid",
			rule: &FeeRule{
				MinValue: nil,
				MaxValue: int64Ptr(5000),
			},
			amount:   5000,
			expected: true,
		},
		{
			name: "only max bound set - above max",
			rule: &FeeRule{
				MinValue: nil,
				MaxValue: int64Ptr(5000),
			},
			amount:   5001,
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.rule.AppliesTo(tt.amount)
			if got != tt.expected {
				t.Errorf("AppliesTo(%d) = %v; want %v", tt.amount, got, tt.expected)
			}
		})
	}
}

func TestFeeRule_ClampFee(t *testing.T) {
	rule := &FeeRule{
		MinFee: int64Ptr(50),
		MaxFee: int64Ptr(2000),
	}

	if got := rule.ClampFee(20); got != 50 {
		t.Errorf("ClampFee(20) = %d; want 50 (min_fee)", got)
	}

	if got := rule.ClampFee(500); got != 500 {
		t.Errorf("ClampFee(500) = %d; want 500", got)
	}

	if got := rule.ClampFee(3000); got != 2000 {
		t.Errorf("ClampFee(3000) = %d; want 2000 (max_fee)", got)
	}
}
