package port

import (
	"context"

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
