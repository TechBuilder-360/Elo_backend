package errors

import "runtime/debug"

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
