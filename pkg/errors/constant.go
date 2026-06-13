package errors

// Predefined codes (example)

type ErrorCode string

const (
	ErrNotFound        ErrorCode = "NOT FOUND"
	ErrInvalidInput    ErrorCode = "INVALID INPUT"
	ErrInternal        ErrorCode = "INTERNAL ERROR"
	ErrFailed          ErrorCode = "REQUEST FAILED"
	ErrTooManyRequests ErrorCode = "TOO MANY REQUESTS"
	ErrInputRequired   ErrorCode = "INPUT REQUIRED"
)
