package port

import (
	"context"

	"github.com/dewisartika8/cicd-status-notifier-bot/internal/core/notification/domain"
	"github.com/dewisartika8/cicd-status-notifier-bot/internal/core/shared/domain/value_objects"
)

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
