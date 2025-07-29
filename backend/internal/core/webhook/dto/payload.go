package dto

import (
	"time"

	"github.com/dewisartika8/cicd-status-notifier-bot/internal/core/shared/domain/value_objects"
	"github.com/dewisartika8/cicd-status-notifier-bot/internal/core/webhook/domain"
)

// GitHubActionsPayload represents GitHub webhook payload structure
type GitHubActionsPayload struct {
	Action     string `json:"action"`
	Workflow   string `json:"workflow"`
	Repository struct {
		Name     string `json:"name"`
		FullName string `json:"full_name"`
		HTMLURL  string `json:"html_url"`
	} `json:"repository"`
	Sender struct {
		Login string `json:"login"`
	} `json:"sender"`
	WorkflowRun *WorkflowRun `json:"workflow_run,omitempty"`
	// Add more fields as needed for different event types
}

// WorkflowRun represents GitHub workflow run information
type WorkflowRun struct {
	ID         int64     `json:"id"`
	Name       string    `json:"name"`
	Status     string    `json:"status"`
	Conclusion string    `json:"conclusion"`
	HTMLURL    string    `json:"html_url"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
	RunNumber  int       `json:"run_number"`
	Event      string    `json:"event"`
	HeadBranch string    `json:"head_branch"`
	HeadSha    string    `json:"head_sha"`
}

// ProcessWebhookRequest represents a request to process a webhook
type ProcessWebhookRequest struct {
	ProjectID  value_objects.ID        `json:"project_id" validate:"required"`
	EventType  domain.WebhookEventType `json:"event_type" validate:"required"`
	Signature  string                  `json:"signature" validate:"required"`
	DeliveryID string                  `json:"delivery_id"`
	Body       []byte                  `json:"body" validate:"required"`
	Payload    GitHubActionsPayload    `json:"payload"`
}

// WebhookEventResponse represents webhook event response
type WebhookEventResponse struct {
	ID          string                  `json:"id"`
	ProjectID   string                  `json:"project_id"`
	EventType   domain.WebhookEventType `json:"event_type"`
	DeliveryID  string                  `json:"delivery_id"`
	ProcessedAt *time.Time              `json:"processed_at"`
	CreatedAt   time.Time               `json:"created_at"`
}

// ToWebhookEventResponse converts domain entity to response DTO
func ToWebhookEventResponse(event *domain.WebhookEvent) *WebhookEventResponse {
	return &WebhookEventResponse{
		ID:          event.ID().String(),
		ProjectID:   event.ProjectID().String(),
		EventType:   event.EventType(),
		DeliveryID:  event.DeliveryID(),
		ProcessedAt: event.ProcessedAt(),
		CreatedAt:   event.CreatedAt().ToTime(),
	}
}
