package port

import (
	"context"

	"github.com/dewisartika8/cicd-status-notifier-bot/internal/core/shared/domain/value_objects"
	"github.com/dewisartika8/cicd-status-notifier-bot/internal/core/webhook/domain"
)

// WebhookEventRepository defines the contract for webhook event data access
type WebhookEventRepository interface {
	// Create stores a new webhook event
	Create(ctx context.Context, event *domain.WebhookEvent) error

	// GetByID retrieves a webhook event by its ID
	GetByID(ctx context.Context, id value_objects.ID) (*domain.WebhookEvent, error)

	// GetByDeliveryID retrieves a webhook event by its delivery ID
	GetByDeliveryID(ctx context.Context, deliveryID string) (*domain.WebhookEvent, error)

	// GetByProjectID retrieves webhook events for a specific project
	GetByProjectID(ctx context.Context, projectID value_objects.ID, limit, offset int) ([]*domain.WebhookEvent, error)

	// Update updates an existing webhook event
	Update(ctx context.Context, event *domain.WebhookEvent) error

	// Delete removes a webhook event
	Delete(ctx context.Context, id value_objects.ID) error

	// ExistsByDeliveryID checks if a webhook event with the given delivery ID exists
	ExistsByDeliveryID(ctx context.Context, deliveryID string) (bool, error)

	// GetUnprocessedEvents retrieves unprocessed webhook events
	GetUnprocessedEvents(ctx context.Context, limit int) ([]*domain.WebhookEvent, error)
}
