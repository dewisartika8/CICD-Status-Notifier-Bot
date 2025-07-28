package errors

import (
	shared_errors "github.com/dewisartika8/cicd-status-notifier-bot/internal/core/shared/domain/errors"
)

// Project-specific error codes
const (
	ErrCodeProjectNotFound      = "PROJECT_NOT_FOUND"
	ErrCodeProjectAlreadyExists = "PROJECT_ALREADY_EXISTS"
	ErrCodeInvalidProjectName   = "INVALID_PROJECT_NAME"
	ErrCodeInvalidRepositoryURL = "INVALID_REPOSITORY_URL"
	ErrCodeProjectNotActive     = "PROJECT_NOT_ACTIVE"
)

// Project-specific domain errors
var (
	ErrProjectNotFound = shared_errors.NewDomainError(
		ErrCodeProjectNotFound,
		"project not found",
	)

	ErrProjectAlreadyExists = shared_errors.NewDomainError(
		ErrCodeProjectAlreadyExists,
		"project with this name already exists",
	)

	ErrInvalidProjectName = shared_errors.NewDomainError(
		ErrCodeInvalidProjectName,
		"project name is invalid",
	)

	ErrInvalidRepositoryURL = shared_errors.NewDomainError(
		ErrCodeInvalidRepositoryURL,
		"repository URL is invalid",
	)

	ErrProjectNotActive = shared_errors.NewDomainError(
		ErrCodeProjectNotActive,
		"project is not active",
	)
)

// NewProjectNotFoundError creates a project not found error with context
func NewProjectNotFoundError(identifier string) shared_errors.DomainError {
	return shared_errors.NewDomainError(
		ErrCodeProjectNotFound,
		"project not found: "+identifier,
	)
}

// NewProjectAlreadyExistsError creates a project already exists error with context
func NewProjectAlreadyExistsError(name string) shared_errors.DomainError {
	return shared_errors.NewDomainError(
		ErrCodeProjectAlreadyExists,
		"project already exists: "+name,
	)
}
