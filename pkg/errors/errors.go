package errors

import (
	"fmt"
)

// Predefined codes (example)
const (
	ErrNotFound     = "NOT_FOUND"
	ErrInvalidInput = "INVALID_INPUT"
	ErrInternal     = "INTERNAL_ERROR"
	ErrFailed       = "REQUEST_FAILED"
)

// CustomError is a structured error type with a code and message
type CustomError struct {
	Code    string
	Message string
	Err     error
}

// Error implements the error interface
func (e *CustomError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("[%s] %s: %v", e.Code, e.Message, e.Err)
	}
	return fmt.Sprintf("[%s] %s", e.Code, e.Message)
}

// Unwrap provides compatibility for errors.Unwrap
func (e *CustomError) Unwrap() error {
	return e
}

// New creates a new CustomError with a code and message
func New(code, message string) error {
	return &CustomError{
		Code:    code,
		Message: message,
	}
}

// Wrap wraps an existing error with code and message
func Wrap(err error, code, message string) error {
	return &CustomError{
		Code:    code,
		Message: message,
		Err:     err,
	}
}

// Unwrap provides compatibility for errors.Unwrap
// func (e *CustomError) Unwrap() error {
// 	return e
// }
