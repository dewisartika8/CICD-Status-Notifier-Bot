package domain

import (
	"github.com/dewisartika8/cicd-status-notifier-bot/pkg/exception"
)

// Build-specific error codes
const (
	ErrCodeBuildEventNotFound    = "BUILD_EVENT_NOT_FOUND"
	ErrCodeInvalidBuildEvent     = "INVALID_BUILD_EVENT"
	ErrCodeInvalidEventType      = "INVALID_EVENT_TYPE"
	ErrCodeInvalidBuildStatus    = "INVALID_BUILD_STATUS"
	ErrCodeInvalidProjectID      = "INVALID_PROJECT_ID"
	ErrCodeInvalidBranch         = "INVALID_BRANCH"
	ErrCodeBuildProcessingFailed = "BUILD_PROCESSING_FAILED"
)

// Build-specific domain errors
var (
	ErrBuildEventNotFound = exception.NewDomainError(
		ErrCodeBuildEventNotFound,
		"build event not found",
	)

	ErrInvalidBuildEvent = exception.NewDomainError(
		ErrCodeInvalidBuildEvent,
		"build event is invalid",
	)

	ErrInvalidEventType = exception.NewDomainError(
		ErrCodeInvalidEventType,
		"event type is invalid",
	)

	ErrInvalidBuildStatus = exception.NewDomainError(
		ErrCodeInvalidBuildStatus,
		"build status is invalid",
	)

	ErrInvalidProjectID = exception.NewDomainError(
		ErrCodeInvalidProjectID,
		"project ID is invalid",
	)

	ErrInvalidBranch = exception.NewDomainError(
		ErrCodeInvalidBranch,
		"branch name is invalid",
	)

	ErrBuildProcessingFailed = exception.NewDomainError(
		ErrCodeBuildProcessingFailed,
		"build processing failed",
	)
)

// NewBuildEventNotFoundError creates a build event not found error with context
func NewBuildEventNotFoundError(identifier string) exception.DomainError {
	return exception.NewDomainError(
		ErrCodeBuildEventNotFound,
		"build event not found: "+identifier,
	)
}

// NewBuildProcessingFailedError creates a build processing failed error with context
func NewBuildProcessingFailedError(reason string) exception.DomainError {
	return exception.NewDomainError(
		ErrCodeBuildProcessingFailed,
		"build processing failed: "+reason,
	)
}
