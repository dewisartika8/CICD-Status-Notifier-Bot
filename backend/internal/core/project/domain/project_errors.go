package domain

import (
	"github.com/dewisartika8/cicd-status-notifier-bot/pkg/exception"
)

// Project-specific error codes
const (
	ErrCodeProjectNotFound      = "PROJECT_NOT_FOUND"
	ErrCodeProjectAlreadyExists = "PROJECT_ALREADY_EXISTS"
	ErrCodeInvalidProjectName   = "INVALID_PROJECT_NAME"
	ErrCodeInvalidRepositoryURL = "INVALID_REPOSITORY_URL"
	ErrCodeProjectNotActive     = "PROJECT_NOT_ACTIVE"
	ErrCodeInvalidWebhookSecret = "INVALID_WEBHOOK_SECRET"
	ErrCodeInvalidTelegramChat  = "INVALID_TELEGRAM_CHAT"
)

// Project-specific domain errors
var (
	ErrProjectNotFound = exception.NewDomainError(
		ErrCodeProjectNotFound,
		"project not found",
	)

	ErrProjectAlreadyExists = exception.NewDomainError(
		ErrCodeProjectAlreadyExists,
		"project with this name already exists",
	)

	ErrInvalidProjectName = exception.NewDomainError(
		ErrCodeInvalidProjectName,
		"project name is invalid",
	)

	ErrInvalidRepositoryURL = exception.NewDomainError(
		ErrCodeInvalidRepositoryURL,
		"repository URL is invalid",
	)

	ErrProjectNotActive = exception.NewDomainError(
		ErrCodeProjectNotActive,
		"project is not active",
	)

	ErrInvalidWebhookSecret = exception.NewDomainError(
		ErrCodeInvalidWebhookSecret,
		"webhook secret is invalid",
	)

	ErrInvalidTelegramChat = exception.NewDomainError(
		ErrCodeInvalidTelegramChat,
		"telegram chat ID is invalid",
	)
)

// NewProjectNotFoundError creates a project not found error with context
func NewProjectNotFoundError(identifier string) exception.DomainError {
	return exception.NewDomainError(
		ErrCodeProjectNotFound,
		"project not found: "+identifier,
	)
}

// NewProjectAlreadyExistsError creates a project already exists error with context
func NewProjectAlreadyExistsError(name string) exception.DomainError {
	return exception.NewDomainError(
		ErrCodeProjectAlreadyExists,
		"project already exists: "+name,
	)
}
