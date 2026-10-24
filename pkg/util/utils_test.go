package util

import "testing"

func TestIsWithinRange(t *testing.T) {
	tests := []struct {
		name     string
		value    int64
		min      int64
		max      int64
		expected bool
	}{
		{
			name:     "within min and max",
			value:    500,
			min:      100,
			max:      1000,
			expected: true,
		},
		{
			name:     "exactly equal to min",
			value:    100,
			min:      100,
			max:      1000,
			expected: true,
		},
		{
			name:     "exactly equal to max",
			value:    1000,
			min:      100,
			max:      1000,
			expected: true,
		},
		{
			name:     "below min",
			value:    99,
			min:      100,
			max:      1000,
			expected: false,
		},
		{
			name:     "above max",
			value:    1001,
			min:      100,
			max:      1000,
			expected: false,
		},
		{
			name:     "no bounds (min=0, max=0)",
			value:    5000,
			min:      0,
			max:      0,
			expected: true,
		},
		{
			name:     "only min bound enforced - valid",
			value:    500,
			min:      100,
			max:      0,
			expected: true,
		},
		{
			name:     "only min bound enforced - below min",
			value:    50,
			min:      100,
			max:      0,
			expected: false,
		},
		{
			name:     "only max bound enforced - valid",
			value:    500,
			min:      0,
			max:      1000,
			expected: true,
		},
		{
			name:     "only max bound enforced - above max",
			value:    1500,
			min:      0,
			max:      1000,
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := IsWithinRange(tt.value, tt.min, tt.max)
			if got != tt.expected {
				t.Errorf("IsWithinRange(%d, %d, %d) = %v; want %v", tt.value, tt.min, tt.max, got, tt.expected)
			}
		})
	}
}
