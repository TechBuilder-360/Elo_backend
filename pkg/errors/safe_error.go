package errors

import (
	"runtime/debug"
	"strings"
)

// SafeError is a structured error type with a code and message
type SafeError struct {
	Code    ErrorCode
	Message string
	Err     error
	Stack   []byte
}

// Error implements the error interface
func (e *SafeError) Error() string {
	return e.Message
}

// Unwrap provides compatibility for errors.Unwrap
func (e *SafeError) Unwrap() error {
	return e.Err
}

// New creates a new SafeError with a code and message
func New(code ErrorCode, message string) error {
	return &SafeError{
		Code:    code,
		Message: message,
	}
}

// Wrap wraps an existing error without error code
func Wrap(err error, message string) error {
	return &SafeError{
		Code:    ErrFailed,
		Message: message,
		Err:     err,
		Stack:   debug.Stack(),
	}
}

// Wrap wraps an existing error with error code and message
func WrapCode(err error, code ErrorCode, message string) error {
	return &SafeError{
		Code:    code,
		Message: message,
		Err:     err,
		Stack:   debug.Stack(),
	}
}

// IsDBError checks if an error or its chain represents a database or SQL driver error.
func IsDBError(err error) bool {
	if err == nil {
		return false
	}
	return IsDBErrorMessage(err.Error())
}

// IsDBErrorMessage reports whether an error message string contains internal database patterns.
func IsDBErrorMessage(msg string) bool {
	lower := strings.ToLower(msg)
	patterns := []string{
		"pq:",
		"sql:",
		"ent:",
		"relation \"",
		"column \"",
		"table \"",
		"duplicate key",
		"violates foreign key",
		"violates unique",
		"violates not-null",
		"violates check",
		"syntax error at or near",
		"driver: bad connection",
		"connection refused",
		"connection reset",
		"deadlock detected",
	}
	for _, p := range patterns {
		if strings.Contains(lower, p) {
			return true
		}
	}
	return false
}

// SanitizeMessage ensures database or raw internal driver error messages are not leaked to API responses.
// If the message contains database error patterns, it replaces it with "request failed".
func SanitizeMessage(msg string) string {
	if IsDBErrorMessage(msg) {
		return "request failed"
	}
	return msg
}

