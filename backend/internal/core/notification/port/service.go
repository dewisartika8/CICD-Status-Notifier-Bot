package port

import (
	"context"
	"time"

	"github.com/dewisartika8/cicd-status-notifier-bot/internal/core/notification/domain"
	"github.com/dewisartika8/cicd-status-notifier-bot/internal/core/notification/dto"
	"github.com/dewisartika8/cicd-status-notifier-bot/internal/core/shared/domain/value_objects"
)

// NotificationTemplateService defines the contract for notification template business logic
type NotificationTemplateService interface {
	// CreateNotificationTemplate creates a new notification template
	CreateNotificationTemplate(
		ctx context.Context,
		templateType domain.NotificationTemplateType,
		channel domain.NotificationChannel,
		subject, bodyTemplate string,
	) (*domain.NotificationTemplate, error)

	// GetNotificationTemplate retrieves a notification template by its ID
	GetNotificationTemplate(ctx context.Context, id value_objects.ID) (*domain.NotificationTemplate, error)

	// GetTemplateByTypeAndChannel retrieves a template by type and channel
	GetTemplateByTypeAndChannel(
		ctx context.Context,
		templateType domain.NotificationTemplateType,
		channel domain.NotificationChannel,
	) (*domain.NotificationTemplate, error)

	// UpdateNotificationTemplate updates an existing notification template
	UpdateNotificationTemplate(
		ctx context.Context,
		id value_objects.ID,
		subject, bodyTemplate string,
	) (*domain.NotificationTemplate, error)

	// ActivateTemplate activates a notification template
	ActivateTemplate(ctx context.Context, id value_objects.ID) error

	// DeactivateTemplate deactivates a notification template
	DeactivateTemplate(ctx context.Context, id value_objects.ID) error

	// DeleteNotificationTemplate deletes a notification template
	DeleteNotificationTemplate(ctx context.Context, id value_objects.ID) error

	// GetActiveTemplates retrieves all active notification templates
	GetActiveTemplates(ctx context.Context) ([]*domain.NotificationTemplate, error)

	// InitializeDefaultTemplates creates default templates for all channels and types
	InitializeDefaultTemplates(ctx context.Context) error
}

// NotificationFormatterService defines the contract for notification formatting
type NotificationFormatterService interface {
	// FormatNotification formats a notification using templates
	FormatNotification(
		ctx context.Context,
		templateType domain.NotificationTemplateType,
		channel domain.NotificationChannel,
		params domain.TemplateParams,
	) (subject, body string, err error)

	// FormatNotificationWithTemplate formats a notification using a specific template
	FormatNotificationWithTemplate(
		ctx context.Context,
		template *domain.NotificationTemplate,
		params domain.TemplateParams,
	) (subject, body string, err error)

	// ValidateTemplate validates a template by trying to compile and render it
	ValidateTemplate(
		templateType domain.NotificationTemplateType,
		channel domain.NotificationChannel,
		subject, bodyTemplate string,
		testParams domain.TemplateParams,
	) error

	// GetAvailableTemplateVariables returns available template variables for a template type
	GetAvailableTemplateVariables(templateType domain.NotificationTemplateType) []string

	// FormatEmoji adds emoji formatting based on build status and channel
	FormatEmoji(status string, channel domain.NotificationChannel) string
}

// NotificationLogService defines the contract for notification log business logic
type NotificationLogService interface {
	// CreateNotificationLog creates a new notification log
	CreateNotificationLog(
		ctx context.Context,
		buildEventID, projectID value_objects.ID,
		channel domain.NotificationChannel,
		recipient, message string,
	) (*domain.NotificationLog, error)

	// SendNotification sends a notification and updates the log
	SendNotification(ctx context.Context, notificationLogID value_objects.ID) error

	// GetNotificationLog retrieves a notification log by its ID
	GetNotificationLog(ctx context.Context, id value_objects.ID) (*domain.NotificationLog, error)

	// GetNotificationLogsByBuildEvent retrieves notification logs for a build event
	GetNotificationLogsByBuildEvent(ctx context.Context, buildEventID value_objects.ID) ([]*domain.NotificationLog, error)

	// GetNotificationLogsByProject retrieves notification logs for a project
	GetNotificationLogsByProject(ctx context.Context, projectID value_objects.ID, limit, offset int) ([]*domain.NotificationLog, error)

	// UpdateNotificationStatus updates the status of a notification log
	UpdateNotificationStatus(
		ctx context.Context,
		id value_objects.ID,
		status domain.NotificationStatus,
		errorMessage string,
		messageID *string,
	) error

	// RetryFailedNotification retries a failed notification
	RetryFailedNotification(ctx context.Context, notificationLogID value_objects.ID) error

	// ProcessPendingNotifications processes all pending notifications
	ProcessPendingNotifications(ctx context.Context, limit int) error

	// ProcessFailedNotifications processes failed notifications for retry
	ProcessFailedNotifications(ctx context.Context, limit int) error

	// GetNotificationStats retrieves notification statistics for a project
	GetNotificationStats(ctx context.Context, projectID value_objects.ID) (map[domain.NotificationStatus]int64, error)

	// CreateNotificationForBuildEvent creates notifications for all subscribed channels for a build event
	CreateNotificationForBuildEvent(
		ctx context.Context,
		buildEventID, projectID value_objects.ID,
		message string,
	) ([]*domain.NotificationLog, error)
}

// TelegramSubscriptionService defines the contract for telegram subscription business logic
type TelegramSubscriptionService interface {
	// CreateTelegramSubscription creates a new telegram subscription
	CreateTelegramSubscription(
		ctx context.Context,
		projectID value_objects.ID,
		chatID int64,
	) (*domain.TelegramSubscription, error)

	// GetTelegramSubscription retrieves a telegram subscription by its ID
	GetTelegramSubscription(ctx context.Context, id value_objects.ID) (*domain.TelegramSubscription, error)

	// GetTelegramSubscriptionsByProject retrieves telegram subscriptions for a project
	GetTelegramSubscriptionsByProject(ctx context.Context, projectID value_objects.ID) ([]*domain.TelegramSubscription, error)

	// GetTelegramSubscriptionByChatID retrieves a telegram subscription by chat ID
	GetTelegramSubscriptionByChatID(ctx context.Context, chatID int64) (*domain.TelegramSubscription, error)

	// UpdateTelegramSubscription updates a telegram subscription
	UpdateTelegramSubscription(
		ctx context.Context,
		id value_objects.ID,
		chatID *int64,
		isActive *bool,
	) (*domain.TelegramSubscription, error)

	// DeleteTelegramSubscription deletes a telegram subscription
	DeleteTelegramSubscription(ctx context.Context, id value_objects.ID) error

	// ActivateTelegramSubscription activates a telegram subscription
	ActivateTelegramSubscription(ctx context.Context, id value_objects.ID) error

	// DeactivateTelegramSubscription deactivates a telegram subscription
	DeactivateTelegramSubscription(ctx context.Context, id value_objects.ID) error

	// GetActiveSubscriptionsForProject retrieves active telegram subscriptions for a project
	GetActiveSubscriptionsForProject(ctx context.Context, projectID value_objects.ID) ([]*domain.TelegramSubscription, error)

	// GetAllActiveSubscriptions retrieves all active telegram subscriptions
	GetAllActiveSubscriptions(ctx context.Context) ([]*domain.TelegramSubscription, error)

	// CheckSubscriptionExists checks if a subscription exists for a project and chat ID
	CheckSubscriptionExists(ctx context.Context, projectID value_objects.ID, chatID int64) (bool, error)
}

// NotificationSender defines the contract for sending notifications through different channels
type NotificationSender interface {
	// SendTelegramNotification sends a notification through Telegram
	SendTelegramNotification(ctx context.Context, chatID int64, message string) (messageID string, err error)

	// SendEmailNotification sends a notification through email
	SendEmailNotification(ctx context.Context, email, subject, message string) error

	// SendSlackNotification sends a notification through Slack
	SendSlackNotification(ctx context.Context, channel, message string) (messageID string, err error)

	// SendWebhookNotification sends a notification through webhook
	SendWebhookNotification(ctx context.Context, webhookURL, message string) error
}

// RetryService defines the interface for retry logic operations
type RetryService interface {
	// CreateRetryConfiguration creates a new retry configuration
	CreateRetryConfiguration(ctx context.Context, req dto.CreateRetryConfigurationRequest) (*domain.RetryConfiguration, error)

	// GetRetryConfiguration retrieves a retry configuration by ID
	GetRetryConfiguration(ctx context.Context, id value_objects.ID) (*domain.RetryConfiguration, error)

	// GetRetryConfigurationByChannel retrieves retry configuration for a channel
	GetRetryConfigurationByChannel(ctx context.Context, channel domain.NotificationChannel) (*domain.RetryConfiguration, error)

	// UpdateRetryConfiguration updates an existing retry configuration
	UpdateRetryConfiguration(ctx context.Context, id value_objects.ID, req dto.UpdateRetryConfigurationRequest) (*domain.RetryConfiguration, error)

	// ActivateRetryConfiguration activates a retry configuration
	ActivateRetryConfiguration(ctx context.Context, id value_objects.ID) error

	// DeactivateRetryConfiguration deactivates a retry configuration
	DeactivateRetryConfiguration(ctx context.Context, id value_objects.ID) error

	// DeleteRetryConfiguration deletes a retry configuration
	DeleteRetryConfiguration(ctx context.Context, id value_objects.ID) error

	// ListActiveRetryConfigurations lists all active retry configurations
	ListActiveRetryConfigurations(ctx context.Context) ([]*domain.RetryConfiguration, error)

	// InitializeDefaultRetryConfigurations sets up default retry configurations
	InitializeDefaultRetryConfigurations(ctx context.Context) error

	// CalculateRetryDelay calculates the delay for a retry attempt
	CalculateRetryDelay(ctx context.Context, channel domain.NotificationChannel, attemptNumber int) (time.Duration, error)

	// ShouldRetryNotification determines if a notification should be retried
	ShouldRetryNotification(ctx context.Context, channel domain.NotificationChannel, attemptCount int, lastError error) (bool, error)

	// ProcessRetryableNotification processes a notification that can be retried
	ProcessRetryableNotification(ctx context.Context, req dto.ProcessRetryableNotificationRequest) (*dto.ProcessRetryableNotificationResponse, error)
}

// DeliveryChannel defines the interface for notification delivery channels (abstraction for Dewi's work)
type DeliveryChannel interface {
	// Send sends a notification through the specific channel
	Send(ctx context.Context, recipient, subject, message string) (messageID string, err error)

	// GetChannelType returns the type of the delivery channel
	GetChannelType() domain.NotificationChannel

	// IsAvailable checks if the delivery channel is available/healthy
	IsAvailable(ctx context.Context) bool

	// GetMaxRetries returns the maximum number of retries for this channel
	GetMaxRetries() int

	// GetRateLimitInfo returns rate limiting information for this channel
	GetRateLimitInfo() (maxRequests int, windowSize time.Duration)
}

// NotificationDeliveryService defines the interface for notification delivery operations
type NotificationDeliveryService interface {
	// QueueNotification adds a notification to the delivery queue
	QueueNotification(ctx context.Context, notification *domain.QueuedNotification) error

	// ProcessQueue processes pending notifications in the queue
	ProcessQueue(ctx context.Context, batchSize int) error

	// ProcessRetryQueue processes failed notifications for retry
	ProcessRetryQueue(ctx context.Context, batchSize int) error

	// GetQueueStats returns queue statistics
	GetQueueStats(ctx context.Context) (map[string]interface{}, error)

	// RegisterDeliveryChannel registers a delivery channel
	RegisterDeliveryChannel(channel DeliveryChannel) error

	// UnregisterDeliveryChannel unregisters a delivery channel
	UnregisterDeliveryChannel(channelType domain.NotificationChannel) error

	// SendNotification sends a notification immediately (bypassing queue)
	SendNotification(ctx context.Context, channel domain.NotificationChannel, recipient, subject, message string) (messageID string, err error)

	// CheckRateLimit checks if a notification can be sent based on rate limiting
	CheckRateLimit(ctx context.Context, channel domain.NotificationChannel, recipient string) (bool, error)
}
