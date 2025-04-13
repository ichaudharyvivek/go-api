package errors

import (
	"fmt"
)

// ApiError represents a custom error with additional fields for context.
type ApiError struct {
	Code    string
	Message string
	Err     error
}

// New creates a new AppError with the specified code, message, and underlying error.
func New(code, message string, err error) *ApiError {
	return &ApiError{
		Code:    code,
		Message: message,
		Err:     err,
	}
}

// ApiError implements the 'error' interface.
func (e *ApiError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("%s: %s: %v", e.Code, e.Message, e.Err)
	}

	return fmt.Sprintf("%s: %s", e.Code, e.Message)
}

// Unwrap returns the underlying error.
// NOTE: Unwrap allows use of errors.Unwrap(), errors.Is(), and errors.As()
func (e *ApiError) Unwrap() error {
	return e.Err
}
