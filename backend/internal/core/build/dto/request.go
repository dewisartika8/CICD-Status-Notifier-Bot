package dto

import (
	"encoding/json"

	"github.com/dewisartika8/cicd-status-notifier-bot/internal/core/build/domain"
	"github.com/dewisartika8/cicd-status-notifier-bot/internal/core/shared/domain/value_objects"
)

// CreateBuildEventRequest represents a request to create a build event
type CreateBuildEventRequest struct {
	ProjectID       value_objects.ID
	EventType       domain.EventType
	Status          domain.BuildStatus
	Branch          string
	CommitSHA       string
	CommitMessage   string
	AuthorName      string
	AuthorEmail     string
	BuildURL        string
	DurationSeconds *int
	WebhookPayload  json.RawMessage
}

// ProcessWebhookRequest represents a request to process a webhook
type ProcessWebhookRequest struct {
	ProjectID      value_objects.ID
	WebhookPayload json.RawMessage
	Headers        map[string]string
	EventType      string
}
