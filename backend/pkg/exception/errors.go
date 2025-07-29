package exception

import (
	"errors"
	"fmt"
)

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

	// Generic validation error codes
	ErrCodeInvalidInput  = "INVALID_INPUT"
	ErrCodeInvalidFormat = "INVALID_FORMAT"
	ErrCodeRequiredField = "REQUIRED_FIELD"
	ErrCodeInvalidLength = "INVALID_LENGTH"
	ErrCodeInvalidRange  = "INVALID_RANGE"
)

// Generic domain errors (can be reused across modules)
var (
	ErrValidationFailed = NewDomainError(ErrCodeValidation, "validation failed")
	ErrNotFound         = NewDomainError(ErrCodeNotFound, "resource not found")
	ErrAlreadyExists    = NewDomainError(ErrCodeAlreadyExists, "resource already exists")
	ErrUnauthorized     = NewDomainError(ErrCodeUnauthorized, "unauthorized access")
	ErrForbidden        = NewDomainError(ErrCodeForbidden, "forbidden access")
	ErrInternalError    = NewDomainError(ErrCodeInternalError, "internal error")
	ErrInvalidInput     = NewDomainError(ErrCodeInvalidInput, "invalid input")
	ErrInvalidFormat    = NewDomainError(ErrCodeInvalidFormat, "invalid format")
	ErrRequiredField    = NewDomainError(ErrCodeRequiredField, "required field missing")
)

// Legacy simple errors - will be migrated to DomainError gradually
var (
	ErrProjectNotFound                   = errors.New("project not found")
	ErrProjectAlreadyExists              = errors.New("project with this name already exists")
	ErrBuildEventNotFound                = errors.New("build event not found")
	ErrTelegramSubscriptionNotFound      = errors.New("telegram subscription not found")
	ErrTelegramSubscriptionAlreadyExists = errors.New("telegram subscription already exists for this project and chat")
	ErrNotificationLogNotFound           = errors.New("notification log not found")

	// Project errors
	ErrInvalidProjectName   = errors.New("project name is required and cannot be empty")
	ErrInvalidRepositoryURL = errors.New("repository URL is required and must be valid")
	ErrInvalidProjectID     = errors.New("project ID is required and cannot be nil")

	// Build event errors
	ErrInvalidEventType    = errors.New("event type is required and cannot be empty")
	ErrInvalidBuildStatus  = errors.New("build status is required and cannot be empty")
	ErrInvalidBranch       = errors.New("branch is required and cannot be empty")
	ErrInvalidBuildEventID = errors.New("build event ID is required and cannot be nil")

	// Telegram subscription errors
	ErrInvalidChatID = errors.New("chat ID is required and cannot be zero")

	// Notification log errors
	ErrInvalidNotificationStatus = errors.New("notification status is required and cannot be empty")
)
