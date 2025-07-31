package domain

import (
	"errors"
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
	ErrCodeMaxRetriesExceeded           = "MAX_RETRIES_EXCEEDED"
	ErrCodeNotificationSendFailed       = "NOTIFICATION_SEND_FAILED"
	ErrCodeTelegramSubscriptionNotFound = "TELEGRAM_SUBSCRIPTION_NOT_FOUND"
	ErrCodeInvalidTelegramChatID        = "INVALID_TELEGRAM_CHAT_ID"
	ErrCodeSubscriptionAlreadyExists    = "SUBSCRIPTION_ALREADY_EXISTS"
	ErrCodeMaxRetryAttemptsExceeded     = "MAX_RETRY_ATTEMPTS_EXCEEDED"
	ErrCodeInvalidMessage               = "INVALID_MESSAGE"
	ErrCodeInvalidProjectID             = "INVALID_PROJECT_ID"
	// Template-specific error codes
	ErrCodeTemplateNotFound        = "TEMPLATE_NOT_FOUND"
	ErrCodeInvalidTemplateType     = "INVALID_TEMPLATE_TYPE"
	ErrCodeInvalidTemplateSubject  = "INVALID_TEMPLATE_SUBJECT"
	ErrCodeInvalidTemplateBody     = "INVALID_TEMPLATE_BODY"
	ErrCodeTemplateRenderError     = "TEMPLATE_RENDER_ERROR"
	ErrCodeTemplateAlreadyActive   = "TEMPLATE_ALREADY_ACTIVE"
	ErrCodeTemplateAlreadyInactive = "TEMPLATE_ALREADY_INACTIVE"
	ErrCodeTemplateInactive        = "TEMPLATE_INACTIVE"
	// Retry configuration error codes
	ErrCodeInvalidRetryConfiguration         = "INVALID_RETRY_CONFIGURATION"
	ErrCodeRetryConfigurationNotFound        = "RETRY_CONFIGURATION_NOT_FOUND"
	ErrCodeRetryConfigurationAlreadyActive   = "RETRY_CONFIGURATION_ALREADY_ACTIVE"
	ErrCodeRetryConfigurationAlreadyInactive = "RETRY_CONFIGURATION_ALREADY_INACTIVE"
)

// Repository layer error variables - for repository implementations
var (
	ErrNotificationTemplateNotFound = errors.New("notification template not found")
	ErrRetryConfigurationNotFound   = errors.New("retry configuration not found")
)

// Generic CRUD error message constants - reusable across all services
const (
	ErrMsgCreate     = "failed to create %s: %w"
	ErrMsgGet        = "failed to get %s: %w"
	ErrMsgUpdate     = "failed to update %s: %w"
	ErrMsgDelete     = "failed to delete %s: %w"
	ErrMsgValidate   = "failed to validate %s: %w"
	ErrMsgProcess    = "failed to process %s: %w"
	ErrMsgSend       = "failed to send %s: %w"
	ErrMsgRender     = "failed to render %s: %w"
	ErrMsgPersist    = "failed to persist %s: %w"
	ErrMsgActivate   = "failed to activate %s: %w"
	ErrMsgDeactivate = "failed to deactivate %s: %w"
)

// Specific service error message constants that cannot be generalized
const (
	ErrMsgCreateNotificationLog    = "failed to create notification log: %w"
	ErrMsgPersistNotificationLog   = "failed to persist notification log: %w"
	ErrMsgGetNotificationLog       = "failed to get notification log: %w"
	ErrMsgUpdateNotificationLog    = "failed to update notification log: %w"
	ErrMsgGetTelegramSubscriptions = "failed to get telegram subscriptions: %w"
	ErrMsgGetPendingNotifications  = "failed to get pending notifications: %w"
	ErrMsgGetFailedNotifications   = "failed to get failed notifications: %w"
	ErrMsgGetNotificationStats     = "failed to get notification stats: %w"
	ErrMsgSendTelegramNotification = "failed to send telegram notification: %w"
	ErrMsgSendEmailNotification    = "failed to send email notification: %w"
	ErrMsgSendSlackNotification    = "failed to send slack notification: %w"
	ErrMsgSendWebhookNotification  = "failed to send webhook notification: %w"
	ErrMsgMarkNotificationAsSent   = "failed to mark notification as sent: %w"
	ErrMsgMarkNotificationAsFailed = "failed to mark notification as failed: %w"
	ErrMsgMarkNotificationAsRetry  = "failed to mark notification as retrying: %w"
)

// Specialized error message constants for operations with specific parameters
const (
	ErrMsgCalculateRetryDelay      = "failed to calculate retry delay: %w"
	ErrMsgProcessRetryable         = "failed to process retryable notification: %w"
	ErrMsgInvalidRetryAttempt      = "invalid retry attempt: %w"
	ErrMsgExceededMaxRetries       = "exceeded maximum retry attempts: %w"
	ErrMsgInitializeTelegramBot    = "failed to initialize telegram bot: %w"
	ErrMsgInitializeEmailClient    = "failed to initialize email client: %w"
	ErrMsgInitializeSlackClient    = "failed to initialize slack client: %w"
	ErrMsgInitializeWebhookClient  = "failed to initialize webhook client: %w"
	ErrMsgValidateNotificationData = "failed to validate notification data: %w"
	ErrMsgParseTemplate            = "failed to parse template: %w"
	ErrMsgExecuteTemplate          = "failed to execute template: %w"
)

// Log message constants for service layer logging
const (
	LogMsgGetNotificationLog         = "Failed to get notification log"
	LogMsgUpdateNotificationLog      = "Failed to update notification log"
	LogMsgSendNotification           = "Failed to send notification"
	LogMsgMarkNotificationFail       = "Failed to mark notification as failed"
	LogMsgMarkNotificationSent       = "Failed to mark notification as sent"
	LogMsgMarkNotificationAsRetrying = "Failed to mark notification as retrying"
)

// Retry service log message constants
const (
	LogMsgGetRetryConfig      = "Failed to get retry configuration"
	LogMsgUpdateRetryConfig   = "Failed to update retry configuration"
	LogMsgDeleteRetryConfig   = "Failed to delete retry configuration"
	LogMsgCalculatingDelay    = "Calculating retry delay"
	LogMsgProcessingRetryable = "Processing retryable notification"
	LogMsgRetryDecision       = "Making retry decision"
)

// Sender service log message constants
const (
	LogMsgSendNotificationFailure = "Failed to send notification"
	LogMsgInitializeService       = "Failed to initialize notification service"
	LogMsgValidateData            = "Failed to validate notification data"
)

// Sender service log message constants
const (
	LogMsgSendingTelegram = "Sending telegram notification"
	LogMsgSendingEmail    = "Sending email notification"
	LogMsgSendingSlack    = "Sending slack notification"
	LogMsgSendingWebhook  = "Sending webhook notification"
)

// Template service log message constants
const (
	LogMsgGetTemplate        = "Failed to get notification template"
	LogMsgUpdateTemplate     = "Failed to update notification template"
	LogMsgActivateTemplate   = "Failed to activate notification template"
	LogMsgDeactivateTemplate = "Failed to deactivate notification template"
	LogMsgDeleteTemplate     = "Failed to delete notification template"
	LogMsgRenderTemplate     = "Failed to render notification template"
	LogMsgValidateTemplate   = "Failed to validate notification template"
)

// Subscription service log message constants
const (
	LogMsgGetSubscription      = "Failed to get telegram subscription"
	LogMsgUpdateSubscription   = "Failed to update telegram subscription"
	LogMsgDeleteSubscription   = "Failed to delete telegram subscription"
	LogMsgValidateSubscription = "Failed to validate telegram subscription"
)

// Formatter service log message constants
const (
	LogMsgFormatNotification = "Failed to format notification"
	LogMsgParseTemplate      = "Failed to parse template"
	LogMsgExecuteTemplate    = "Failed to execute template"
)

// Service operation success messages
const (
	RetryConfigCreated             = "retry configuration created successfully"
	RetryConfigUpdated             = "retry configuration updated successfully"
	RetryConfigActivated           = "retry configuration activated successfully"
	RetryConfigDeactivated         = "retry configuration deactivated successfully"
	RetryConfigDeleted             = "retry configuration deleted successfully"
	DefaultRetryConfigsInitialized = "default retry configurations initialized successfully"
	RetryDelayCalculated           = "retry delay calculated successfully"
	RetryDecisionMade              = "retry decision made successfully"
	RetryableNotificationProcessed = "retryable notification processed successfully"
	NotificationSent               = "notification sent successfully"
	TelegramNotificationSent       = "telegram notification sent successfully"
	EmailNotificationSent          = "email notification sent successfully"
	SlackNotificationSent          = "slack notification sent successfully"
	WebhookNotificationSent        = "webhook notification sent successfully"
	NotificationFormatted          = "notification formatted successfully"
	TemplateCreated                = "notification template created successfully"
	TemplateUpdated                = "notification template updated successfully"
	TemplateActivated              = "notification template activated successfully"
	TemplateDeactivated            = "notification template deactivated successfully"
	TemplateDeleted                = "notification template deleted successfully"
	SubscriptionCreated            = "telegram subscription created successfully"
	SubscriptionUpdated            = "telegram subscription updated successfully"
	SubscriptionActivated          = "telegram subscription activated successfully"
	SubscriptionDeactivated        = "telegram subscription deactivated successfully"
	SubscriptionDeleted            = "telegram subscription deleted successfully"
)

// Error type constants for retryable errors
const (
	ErrTypeTimeout          = "timeout"
	ErrTypeNetworkError     = "network error"
	ErrTypeTelegramAPI      = "telegram api error"
	ErrTypeRateLimit        = "rate limit"
	ErrTypeTemporaryFailure = "temporary failure"
	ErrTypeSMTPError        = "smtp error"
	ErrTypeSlackAPI         = "slack api error"
	ErrTypeHTTP5xx          = "http 5xx error"
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

	// Template-specific domain errors
	ErrTemplateNotFound = exception.NewDomainError(
		ErrCodeTemplateNotFound,
		"notification template not found",
	)

	ErrInvalidTemplateBody = exception.NewDomainError(
		ErrCodeInvalidTemplateBody,
		"template body is invalid or empty",
	)

	ErrTemplateAlreadyActive = exception.NewDomainError(
		ErrCodeTemplateAlreadyActive,
		"template is already active",
	)

	ErrTemplateAlreadyInactive = exception.NewDomainError(
		ErrCodeTemplateAlreadyInactive,
		"template is already inactive",
	)

	ErrTemplateInactive = exception.NewDomainError(
		ErrCodeTemplateInactive,
		"template is inactive and cannot be rendered",
	)

	// Retry configuration domain errors
	ErrRetryConfigurationAlreadyActive = exception.NewDomainError(
		ErrCodeRetryConfigurationAlreadyActive,
		"retry configuration is already active",
	)

	ErrRetryConfigurationAlreadyInactive = exception.NewDomainError(
		ErrCodeRetryConfigurationAlreadyInactive,
		"retry configuration is already inactive",
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

func NewInvalidTemplateTypeError(templateType string) error {
	return exception.NewDomainError(
		ErrCodeInvalidTemplateType,
		fmt.Sprintf("invalid template type: %s", templateType),
	)
}

func NewInvalidTemplateSubjectError(message string) error {
	return exception.NewDomainError(
		ErrCodeInvalidTemplateSubject,
		message,
	)
}

func NewInvalidTemplateBodyError(message string) error {
	return exception.NewDomainError(
		ErrCodeInvalidTemplateBody,
		message,
	)
}

func NewTemplateRenderError(message string) error {
	return exception.NewDomainError(
		ErrCodeTemplateRenderError,
		message,
	)
}

func NewInvalidRetryConfigurationError(message string) error {
	return exception.NewDomainError(
		ErrCodeInvalidRetryConfiguration,
		message,
	)
}

func NewRetryConfigurationNotFoundError(configID string) error {
	return exception.NewDomainError(
		ErrCodeRetryConfigurationNotFound,
		fmt.Sprintf("retry configuration with ID %s not found", configID),
	)
}

func NewInvalidRecipientError(recipient string) error {
	return exception.NewDomainError(
		ErrCodeInvalidRecipient,
		fmt.Sprintf("invalid recipient: %s", recipient),
	)
}

func NewMaxRetriesExceededError(maxRetries int) error {
	return exception.NewDomainError(
		ErrCodeMaxRetriesExceeded,
		fmt.Sprintf("maximum retry attempts (%d) exceeded", maxRetries),
	)
}
