package errors

import (
	"errors"
	"testing"
)

func TestIsDBErrorMessage(t *testing.T) {
	tests := []struct {
		name     string
		msg      string
		expected bool
	}{
		{
			name:     "postgres column error",
			msg:      "pq: column fee_rules.min_fee does not exist",
			expected: true,
		},
		{
			name:     "postgres relation error",
			msg:      "pq: relation \"accounts\" does not exist",
			expected: true,
		},
		{
			name:     "duplicate key constraint error",
			msg:      "duplicate key value violates unique constraint",
			expected: true,
		},
		{
			name:     "foreign key violation",
			msg:      "violates foreign key constraint \"fee_rules_services_fee_rules\"",
			expected: true,
		},
		{
			name:     "ent constraint error",
			msg:      "ent: constraint failed: unique index",
			expected: true,
		},
		{
			name:     "sql driver error",
			msg:      "driver: bad connection",
			expected: true,
		},
		{
			name:     "business validation message",
			msg:      "insufficient fund",
			expected: false,
		},
		{
			name:     "input validation message",
			msg:      "amount is outside allowable limit",
			expected: false,
		},
		{
			name:     "not found message",
			msg:      "service not found",
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := IsDBErrorMessage(tt.msg)
			if got != tt.expected {
				t.Errorf("IsDBErrorMessage(%q) = %v; want %v", tt.msg, got, tt.expected)
			}
		})
	}
}

func TestSanitizeMessage(t *testing.T) {
	dbMsg := "pq: column fee_rules.min_fee does not exist"
	if got := SanitizeMessage(dbMsg); got != "request failed" {
		t.Errorf("SanitizeMessage(%q) = %q; want \"request failed\"", dbMsg, got)
	}

	userMsg := "insufficient funds in treasury wallet"
	if got := SanitizeMessage(userMsg); got != userMsg {
		t.Errorf("SanitizeMessage(%q) = %q; want %q", userMsg, got, userMsg)
	}
}

func TestIsDBError(t *testing.T) {
	if IsDBError(nil) {
		t.Error("IsDBError(nil) should be false")
	}

	err := errors.New("pq: syntax error at or near \"WHERE\"")
	if !IsDBError(err) {
		t.Error("IsDBError should be true for pq error")
	}
}
