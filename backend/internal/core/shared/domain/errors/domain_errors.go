package errors

import "fmt"

// DomainError represents a domain-specific error
type DomainError struct {
	Code    string
	Message string
	Cause   error
}

// Error implements the error interface
func (e DomainError) Error() string {
	if e.Cause != nil {
		return fmt.Sprintf("%s: %s (caused by: %v)", e.Code, e.Message, e.Cause)
	}
	return fmt.Sprintf("%s: %s", e.Code, e.Message)
}

// Unwrap returns the underlying cause
func (e DomainError) Unwrap() error {
	return e.Cause
}

// NewDomainError creates a new domain error
func NewDomainError(code, message string) DomainError {
	return DomainError{
		Code:    code,
		Message: message,
	}
}

// NewDomainErrorWithCause creates a new domain error with cause
func NewDomainErrorWithCause(code, message string, cause error) DomainError {
	return DomainError{
		Code:    code,
		Message: message,
		Cause:   cause,
	}
}

// Common domain error codes
const (
	ErrCodeValidation    = "VALIDATION_ERROR"
	ErrCodeNotFound      = "NOT_FOUND"
	ErrCodeAlreadyExists = "ALREADY_EXISTS"
	ErrCodeUnauthorized  = "UNAUTHORIZED"
	ErrCodeForbidden     = "FORBIDDEN"
	ErrCodeInternalError = "INTERNAL_ERROR"
)

// Common domain errors
var (
	ErrValidationFailed = NewDomainError(ErrCodeValidation, "validation failed")
	ErrNotFound         = NewDomainError(ErrCodeNotFound, "resource not found")
	ErrAlreadyExists    = NewDomainError(ErrCodeAlreadyExists, "resource already exists")
	ErrUnauthorized     = NewDomainError(ErrCodeUnauthorized, "unauthorized access")
	ErrForbidden        = NewDomainError(ErrCodeForbidden, "forbidden access")
	ErrInternalError    = NewDomainError(ErrCodeInternalError, "internal error occurred")
)
