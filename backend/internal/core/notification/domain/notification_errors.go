package domain

import (
	"fmt"

	"github.com/dewisartika8/cicd-status-notifier-bot/pkg/exception"
)

// Notification-specific error codes
const (
	ErrCodeNotificationLogNotFound      = "NOTIFICATION_LOG_NOT_FOUND"
	ErrCodeInvalidNotificationLog       = "INVALID_NOTIFICATION_LOG"
	ErrCodeInvalidNotificationChannel   = "INVALID_NOTIFICATION_CHANNEL"
	ErrCodeInvalidNotificationStatus    = "INVALID_NOTIFICATION_STATUS"
	ErrCodeInvalidRecipient             = "INVALID_RECIPIENT"
	ErrCodeNotificationSendFailed       = "NOTIFICATION_SEND_FAILED"
	ErrCodeTelegramSubscriptionNotFound = "TELEGRAM_SUBSCRIPTION_NOT_FOUND"
	ErrCodeInvalidTelegramChatID        = "INVALID_TELEGRAM_CHAT_ID"
	ErrCodeSubscriptionAlreadyExists    = "SUBSCRIPTION_ALREADY_EXISTS"
	ErrCodeMaxRetryAttemptsExceeded     = "MAX_RETRY_ATTEMPTS_EXCEEDED"
	ErrCodeInvalidMessage               = "INVALID_MESSAGE"
	ErrCodeInvalidProjectID             = "INVALID_PROJECT_ID"
)

// Notification-specific domain errors
var (
	ErrNotificationLogNotFound = exception.NewDomainError(
		ErrCodeNotificationLogNotFound,
		"notification log not found",
	)

	ErrInvalidNotificationLog = exception.NewDomainError(
		ErrCodeInvalidNotificationLog,
		"notification log is invalid",
	)

	ErrInvalidNotificationChannel = exception.NewDomainError(
		ErrCodeInvalidNotificationChannel,
		"notification channel is invalid",
	)

	ErrInvalidNotificationStatus = exception.NewDomainError(
		ErrCodeInvalidNotificationStatus,
		"notification status is invalid",
	)

	ErrInvalidRecipient = exception.NewDomainError(
		ErrCodeInvalidRecipient,
		"notification recipient is invalid",
	)

	ErrNotificationSendFailed = exception.NewDomainError(
		ErrCodeNotificationSendFailed,
		"failed to send notification",
	)

	ErrTelegramSubscriptionNotFound = exception.NewDomainError(
		ErrCodeTelegramSubscriptionNotFound,
		"telegram subscription not found",
	)

	ErrInvalidTelegramChatID = exception.NewDomainError(
		ErrCodeInvalidTelegramChatID,
		"telegram chat ID is invalid",
	)

	ErrSubscriptionAlreadyExists = exception.NewDomainError(
		ErrCodeSubscriptionAlreadyExists,
		"subscription already exists for this project and chat",
	)

	ErrMaxRetryAttemptsExceeded = exception.NewDomainError(
		ErrCodeMaxRetryAttemptsExceeded,
		"maximum retry attempts exceeded",
	)

	ErrInvalidMessage = exception.NewDomainError(
		ErrCodeInvalidMessage,
		"notification message is invalid",
	)

	ErrInvalidProjectID = exception.NewDomainError(
		ErrCodeInvalidProjectID,
		"project ID is required and cannot be nil",
	)

	ErrInvalidChatID = exception.NewDomainError(
		ErrCodeInvalidTelegramChatID,
		"chat ID is required and cannot be zero",
	)
)

// Helper functions to create domain errors with context
func NewNotificationSendFailedError(cause error) error {
	return exception.NewDomainErrorWithCause(
		ErrCodeNotificationSendFailed,
		"failed to send notification",
		cause,
	)
}

func NewInvalidRecipientError(recipient string) error {
	return exception.NewDomainError(
		ErrCodeInvalidRecipient,
		fmt.Sprintf("invalid recipient: %s", recipient),
	)
}
