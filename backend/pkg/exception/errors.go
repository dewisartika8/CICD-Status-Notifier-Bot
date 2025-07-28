package exception

import "errors"

// Define custom error messages
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
