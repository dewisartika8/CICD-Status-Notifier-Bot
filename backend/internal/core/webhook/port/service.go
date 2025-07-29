package port

import (
	"context"

	"github.com/dewisartika8/cicd-status-notifier-bot/internal/core/shared/domain/value_objects"
	"github.com/dewisartika8/cicd-status-notifier-bot/internal/core/webhook/domain"
	"github.com/dewisartika8/cicd-status-notifier-bot/internal/core/webhook/dto"
)

// WebhookService defines the contract for webhook business logic
type WebhookService interface {
	// ProcessWebhook processes an incoming webhook request
	ProcessWebhook(ctx context.Context, req dto.ProcessWebhookRequest) (*domain.WebhookEvent, error)

	// VerifyWebhookSignature verifies the webhook signature
	VerifyWebhookSignature(secret, signature string, body []byte) bool

	// GetWebhookEvent retrieves a webhook event by its ID
	GetWebhookEvent(ctx context.Context, id value_objects.ID) (*domain.WebhookEvent, error)

	// GetWebhookEventsByProject retrieves webhook events for a specific project
	GetWebhookEventsByProject(ctx context.Context, projectID value_objects.ID, limit, offset int) ([]*domain.WebhookEvent, error)

	// ReprocessFailedWebhooks reprocesses failed webhook events
	ReprocessFailedWebhooks(ctx context.Context, limit int) error
}
