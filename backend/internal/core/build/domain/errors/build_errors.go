package errors

import (
	shared_errors "github.com/dewisartika8/cicd-status-notifier-bot/internal/core/shared/domain/errors"
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
	ErrBuildEventNotFound = shared_errors.NewDomainError(
		ErrCodeBuildEventNotFound,
		"build event not found",
	)

	ErrInvalidBuildEvent = shared_errors.NewDomainError(
		ErrCodeInvalidBuildEvent,
		"build event is invalid",
	)

	ErrInvalidEventType = shared_errors.NewDomainError(
		ErrCodeInvalidEventType,
		"event type is invalid",
	)

	ErrInvalidBuildStatus = shared_errors.NewDomainError(
		ErrCodeInvalidBuildStatus,
		"build status is invalid",
	)

	ErrInvalidProjectID = shared_errors.NewDomainError(
		ErrCodeInvalidProjectID,
		"project ID is invalid",
	)

	ErrInvalidBranch = shared_errors.NewDomainError(
		ErrCodeInvalidBranch,
		"branch name is invalid",
	)

	ErrBuildProcessingFailed = shared_errors.NewDomainError(
		ErrCodeBuildProcessingFailed,
		"build processing failed",
	)
)

// NewBuildEventNotFoundError creates a build event not found error with context
func NewBuildEventNotFoundError(identifier string) shared_errors.DomainError {
	return shared_errors.NewDomainError(
		ErrCodeBuildEventNotFound,
		"build event not found: "+identifier,
	)
}

// NewBuildProcessingFailedError creates a build processing failed error with context
func NewBuildProcessingFailedError(reason string) shared_errors.DomainError {
	return shared_errors.NewDomainError(
		ErrCodeBuildProcessingFailed,
		"build processing failed: "+reason,
	)
}
