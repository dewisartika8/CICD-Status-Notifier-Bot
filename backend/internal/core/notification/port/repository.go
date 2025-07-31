package port

import (
	"context"
	"time"

	"github.com/dewisartika8/cicd-status-notifier-bot/internal/core/notification/domain"
	"github.com/dewisartika8/cicd-status-notifier-bot/internal/core/shared/domain/value_objects"
)

// NotificationTemplateRepository defines the contract for notification template data access
type NotificationTemplateRepository interface {
	// Create creates a new notification template
	Create(ctx context.Context, template *domain.NotificationTemplate) error

	// GetByID retrieves a notification template by its ID
	GetByID(ctx context.Context, id value_objects.ID) (*domain.NotificationTemplate, error)

	// GetByTypeAndChannel retrieves a notification template by type and channel
	GetByTypeAndChannel(ctx context.Context, templateType domain.NotificationTemplateType, channel domain.NotificationChannel) (*domain.NotificationTemplate, error)

	// GetByType retrieves all notification templates for a specific type
	GetByType(ctx context.Context, templateType domain.NotificationTemplateType) ([]*domain.NotificationTemplate, error)

	// GetByChannel retrieves all notification templates for a specific channel
	GetByChannel(ctx context.Context, channel domain.NotificationChannel) ([]*domain.NotificationTemplate, error)

	// GetActiveTemplates retrieves all active notification templates
	GetActiveTemplates(ctx context.Context) ([]*domain.NotificationTemplate, error)

	// Update updates an existing notification template
	Update(ctx context.Context, template *domain.NotificationTemplate) error

	// Delete deletes a notification template by its ID
	Delete(ctx context.Context, id value_objects.ID) error

	// Count returns the total number of notification templates
	Count(ctx context.Context, templateType *domain.NotificationTemplateType, channel *domain.NotificationChannel, isActive *bool) (int64, error)
}

// NotificationLogRepository defines the contract for notification log data access
type NotificationLogRepository interface {
	// Create creates a new notification log
	Create(ctx context.Context, log *domain.NotificationLog) error

	// GetByID retrieves a notification log by its ID
	GetByID(ctx context.Context, id value_objects.ID) (*domain.NotificationLog, error)

	// GetByBuildEventID retrieves all notification logs for a build event
	GetByBuildEventID(ctx context.Context, buildEventID value_objects.ID) ([]*domain.NotificationLog, error)

	// GetByProjectID retrieves notification logs for a specific project
	GetByProjectID(ctx context.Context, projectID value_objects.ID, limit, offset int) ([]*domain.NotificationLog, error)

	// GetByRecipient retrieves notification logs for a specific recipient
	GetByRecipient(ctx context.Context, recipient string, limit, offset int) ([]*domain.NotificationLog, error)

	// Update updates an existing notification log
	Update(ctx context.Context, log *domain.NotificationLog) error

	// Delete deletes a notification log by its ID
	Delete(ctx context.Context, id value_objects.ID) error

	// GetFailedNotifications retrieves failed notifications for retry
	GetFailedNotifications(ctx context.Context, limit int) ([]*domain.NotificationLog, error)

	// GetPendingNotifications retrieves pending notifications
	GetPendingNotifications(ctx context.Context, limit int) ([]*domain.NotificationLog, error)

	// Count returns the total number of notification logs matching the criteria
	Count(ctx context.Context, projectID *value_objects.ID, status *domain.NotificationStatus) (int64, error)

	// GetNotificationStats retrieves notification statistics for a project
	GetNotificationStats(ctx context.Context, projectID value_objects.ID) (*domain.NotificationStats, error)
}

// TelegramSubscriptionRepository defines the contract for telegram subscription data access
type TelegramSubscriptionRepository interface {
	// Create creates a new telegram subscription
	Create(ctx context.Context, subscription *domain.TelegramSubscription) error

	// GetByID retrieves a telegram subscription by its ID
	GetByID(ctx context.Context, id value_objects.ID) (*domain.TelegramSubscription, error)

	// GetByProjectID retrieves telegram subscriptions for a specific project
	GetByProjectID(ctx context.Context, projectID value_objects.ID) ([]*domain.TelegramSubscription, error)

	// GetByChatID retrieves a telegram subscription by chat ID
	GetByChatID(ctx context.Context, chatID int64) (*domain.TelegramSubscription, error)

	// GetByProjectAndChatID retrieves a specific subscription by project and chat ID
	GetByProjectAndChatID(ctx context.Context, projectID value_objects.ID, chatID int64) (*domain.TelegramSubscription, error)

	// Update updates an existing telegram subscription
	Update(ctx context.Context, subscription *domain.TelegramSubscription) error

	// Delete deletes a telegram subscription by its ID
	Delete(ctx context.Context, id value_objects.ID) error

	// GetActiveSubscriptions retrieves all active telegram subscriptions
	GetActiveSubscriptions(ctx context.Context) ([]*domain.TelegramSubscription, error)

	// GetActiveSubscriptionsByProject retrieves active subscriptions for a project
	GetActiveSubscriptionsByProject(ctx context.Context, projectID value_objects.ID) ([]*domain.TelegramSubscription, error)

	// ExistsByProjectAndChatID checks if a subscription exists for project and chat
	ExistsByProjectAndChatID(ctx context.Context, projectID value_objects.ID, chatID int64) (bool, error)

	// Count returns the total number of telegram subscriptions
	Count(ctx context.Context, projectID *value_objects.ID, isActive *bool) (int64, error)
}

// RetryConfigurationRepository defines the interface for retry configuration persistence
type RetryConfigurationRepository interface {
	// Create saves a new retry configuration
	Create(ctx context.Context, config *domain.RetryConfiguration) error

	// GetByID retrieves a retry configuration by ID
	GetByID(ctx context.Context, id value_objects.ID) (*domain.RetryConfiguration, error)

	// GetActiveConfigurations retrieves all active retry configurations
	GetActiveConfigurations(ctx context.Context) ([]*domain.RetryConfiguration, error)

	// GetByChannel retrieves retry configuration for a specific channel
	GetByChannel(ctx context.Context, channel domain.NotificationChannel) (*domain.RetryConfiguration, error)

	// Update saves changes to an existing retry configuration
	Update(ctx context.Context, config *domain.RetryConfiguration) error

	// Delete removes a retry configuration
	Delete(ctx context.Context, id value_objects.ID) error

	// BulkCreate saves multiple retry configurations
	BulkCreate(ctx context.Context, configs []*domain.RetryConfiguration) error
}

// DeliveryQueueRepository defines the interface for delivery queue persistence
type DeliveryQueueRepository interface {
	// Create saves a new queued notification
	Create(ctx context.Context, notification *domain.QueuedNotification) error

	// GetByID retrieves a queued notification by ID
	GetByID(ctx context.Context, id value_objects.ID) (*domain.QueuedNotification, error)

	// GetPendingNotifications retrieves pending notifications ready for processing
	GetPendingNotifications(ctx context.Context, limit int) ([]*domain.QueuedNotification, error)

	// GetPendingByPriority retrieves pending notifications ordered by priority
	GetPendingByPriority(ctx context.Context, limit int) ([]*domain.QueuedNotification, error)

	// GetFailedNotifications retrieves failed notifications that can be retried
	GetFailedNotifications(ctx context.Context, limit int) ([]*domain.QueuedNotification, error)

	// Update saves changes to an existing queued notification
	Update(ctx context.Context, notification *domain.QueuedNotification) error

	// UpdateStatus updates only the status and error message of a notification
	UpdateStatus(ctx context.Context, id value_objects.ID, status domain.DeliveryStatus, errorMessage string) error

	// Delete removes a queued notification
	Delete(ctx context.Context, id value_objects.ID) error

	// DeleteProcessedNotifications removes successfully delivered notifications older than specified duration
	DeleteProcessedNotifications(ctx context.Context, olderThan time.Duration) error

	// GetPendingCount returns the count of pending notifications
	GetPendingCount(ctx context.Context) (int64, error)

	// GetQueueStats returns queue statistics by status
	GetQueueStats(ctx context.Context) (map[string]int64, error)
}

// RateLimiterRepository defines the interface for rate limiter persistence
type RateLimiterRepository interface {
	// GetEntry retrieves a rate limit entry
	GetEntry(ctx context.Context, key string, channel domain.NotificationChannel) (*domain.RateLimitEntry, error)

	// SetEntry saves or updates a rate limit entry
	SetEntry(ctx context.Context, entry *domain.RateLimitEntry) error

	// DeleteEntry removes a rate limit entry
	DeleteEntry(ctx context.Context, key string, channel domain.NotificationChannel) error

	// GetRule retrieves a rate limiting rule for a channel
	GetRule(ctx context.Context, channel domain.NotificationChannel) (*domain.RateLimitRule, error)

	// SetRule saves or updates a rate limiting rule
	SetRule(ctx context.Context, rule *domain.RateLimitRule) error

	// DeleteRule removes a rate limiting rule
	DeleteRule(ctx context.Context, channel domain.NotificationChannel) error

	// GetAllRules retrieves all rate limiting rules
	GetAllRules(ctx context.Context) ([]*domain.RateLimitRule, error)

	// CleanupExpiredEntries removes expired rate limit entries
	CleanupExpiredEntries(ctx context.Context) error
}
