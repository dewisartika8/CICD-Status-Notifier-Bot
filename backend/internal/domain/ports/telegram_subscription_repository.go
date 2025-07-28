package ports

import (
	"context"

	"github.com/dewisartika8/cicd-status-notifier-bot/internal/domain/entities"
	"github.com/google/uuid"
)

// TelegramSubscriptionRepository defines the contract for telegram subscription data access
type TelegramSubscriptionRepository interface {
	// Create creates a new telegram subscription
	Create(ctx context.Context, subscription *entities.TelegramSubscription) error

	// GetByID retrieves a subscription by its ID
	GetByID(ctx context.Context, id uuid.UUID) (*entities.TelegramSubscription, error)

	// GetByProjectAndChat retrieves a subscription by project ID and chat ID
	GetByProjectAndChat(ctx context.Context, projectID uuid.UUID, chatID int64) (*entities.TelegramSubscription, error)

	// GetActiveByProjectID retrieves all active subscriptions for a project
	GetActiveByProjectID(ctx context.Context, projectID uuid.UUID) ([]*entities.TelegramSubscription, error)

	// GetByChatID retrieves all subscriptions for a chat
	GetByChatID(ctx context.Context, chatID int64) ([]*entities.TelegramSubscription, error)

	// Update updates an existing subscription
	Update(ctx context.Context, subscription *entities.TelegramSubscription) error

	// Delete deletes a subscription by its ID
	Delete(ctx context.Context, id uuid.UUID) error

	// DeleteByProjectAndChat deletes a subscription by project ID and chat ID
	DeleteByProjectAndChat(ctx context.Context, projectID uuid.UUID, chatID int64) error
}
