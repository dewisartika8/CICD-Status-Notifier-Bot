package port

import (
	"context"

	"github.com/dewisartika8/cicd-status-notifier-bot/internal/domain/entities"
	"github.com/google/uuid"
)

// NotificationLogRepository defines the contract for notification log data access
type NotificationLogRepository interface {
	// Create creates a new notification log
	Create(ctx context.Context, log *entities.NotificationLog) error

	// GetByID retrieves a notification log by its ID
	GetByID(ctx context.Context, id uuid.UUID) (*entities.NotificationLog, error)

	// GetByBuildEventID retrieves all notification logs for a build event
	GetByBuildEventID(ctx context.Context, buildEventID uuid.UUID) ([]*entities.NotificationLog, error)

	// GetByChatID retrieves notification logs for a specific chat
	GetByChatID(ctx context.Context, chatID int64, limit, offset int) ([]*entities.NotificationLog, error)

	// Update updates an existing notification log
	Update(ctx context.Context, log *entities.NotificationLog) error

	// GetFailedNotifications retrieves failed notifications for retry
	GetFailedNotifications(ctx context.Context, limit int) ([]*entities.NotificationLog, error)
}
