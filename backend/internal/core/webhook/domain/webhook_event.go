package domain

import (
	"time"

	"github.com/dewisartika8/cicd-status-notifier-bot/internal/core/shared/domain/value_objects"
)

// WebhookEventType represents the type of webhook event
type WebhookEventType string

const (
	// Supported webhook event types
	WorkflowRunEvent WebhookEventType = "workflow_run"
	PushEvent        WebhookEventType = "push"
	PullRequestEvent WebhookEventType = "pull_request"
)

// WebhookEvent represents a webhook event received from GitHub
type WebhookEvent struct {
	id          value_objects.ID
	projectID   value_objects.ID
	eventType   WebhookEventType
	payload     string // JSON payload as string
	signature   string
	deliveryID  string // GitHub delivery ID
	processedAt *time.Time
	createdAt   value_objects.Timestamp
}

// WebhookEventData holds the data needed to create or reconstruct a WebhookEvent
type WebhookEventData struct {
	ID          value_objects.ID
	ProjectID   value_objects.ID
	EventType   WebhookEventType
	Payload     string
	Signature   string
	DeliveryID  string
	ProcessedAt *time.Time
	CreatedAt   value_objects.Timestamp
}

// NewWebhookEvent creates a new webhook event
func NewWebhookEvent(projectID value_objects.ID, eventType WebhookEventType, payload, signature, deliveryID string) (*WebhookEvent, error) {
	id := value_objects.NewID()
	now := value_objects.NewTimestamp()

	// Validate event type
	if !isValidEventType(eventType) {
		return nil, NewWebhookInvalidEventError(string(eventType))
	}

	// Validate required fields
	if payload == "" {
		return nil, NewWebhookInvalidPayloadError("payload cannot be empty")
	}

	if signature == "" {
		return nil, NewWebhookInvalidPayloadError("signature cannot be empty")
	}

	return &WebhookEvent{
		id:         id,
		projectID:  projectID,
		eventType:  eventType,
		payload:    payload,
		signature:  signature,
		deliveryID: deliveryID,
		createdAt:  now,
	}, nil
}

// NewWebhookEventFromData creates a webhook event from existing data (e.g., from database)
func NewWebhookEventFromData(data WebhookEventData) *WebhookEvent {
	return &WebhookEvent{
		id:          data.ID,
		projectID:   data.ProjectID,
		eventType:   data.EventType,
		payload:     data.Payload,
		signature:   data.Signature,
		deliveryID:  data.DeliveryID,
		processedAt: data.ProcessedAt,
		createdAt:   data.CreatedAt,
	}
}

// MarkAsProcessed marks the webhook event as processed
func (w *WebhookEvent) MarkAsProcessed() {
	now := time.Now()
	w.processedAt = &now
}

// IsProcessed returns true if the webhook event has been processed
func (w *WebhookEvent) IsProcessed() bool {
	return w.processedAt != nil
}

// Getters
func (w *WebhookEvent) ID() value_objects.ID               { return w.id }
func (w *WebhookEvent) ProjectID() value_objects.ID        { return w.projectID }
func (w *WebhookEvent) EventType() WebhookEventType        { return w.eventType }
func (w *WebhookEvent) Payload() string                    { return w.payload }
func (w *WebhookEvent) Signature() string                  { return w.signature }
func (w *WebhookEvent) DeliveryID() string                 { return w.deliveryID }
func (w *WebhookEvent) ProcessedAt() *time.Time            { return w.processedAt }
func (w *WebhookEvent) CreatedAt() value_objects.Timestamp { return w.createdAt }

// isValidEventType checks if the event type is supported
func isValidEventType(eventType WebhookEventType) bool {
	switch eventType {
	case WorkflowRunEvent, PushEvent, PullRequestEvent:
		return true
	default:
		return false
	}
}
