package ports

import (
	"context"
	"encoding/json"

	"github.com/dewisartika8/cicd-status-notifier-bot/internal/core/build/domain/entities"
	"github.com/dewisartika8/cicd-status-notifier-bot/internal/core/shared/domain/value_objects"
)

// BuildEventService defines the contract for build event business logic
type BuildEventService interface {
	// CreateBuildEvent creates a new build event
	CreateBuildEvent(ctx context.Context, req CreateBuildEventRequest) (*entities.BuildEvent, error)

	// ProcessWebhookEvent processes a webhook event and creates build events
	ProcessWebhookEvent(ctx context.Context, req ProcessWebhookRequest) ([]*entities.BuildEvent, error)

	// GetBuildEvent retrieves a build event by its ID
	GetBuildEvent(ctx context.Context, id value_objects.ID) (*entities.BuildEvent, error)

	// GetBuildEventsByProject retrieves build events for a specific project
	GetBuildEventsByProject(ctx context.Context, projectID value_objects.ID, filters ListBuildEventFilters) ([]*entities.BuildEvent, error)

	// UpdateBuildEventStatus updates the status of a build event
	UpdateBuildEventStatus(ctx context.Context, id value_objects.ID, status entities.BuildStatus, duration *int) error

	// GetLatestBuildEvent gets the latest build event for a project
	GetLatestBuildEvent(ctx context.Context, projectID value_objects.ID) (*entities.BuildEvent, error)

	// GetBuildMetrics retrieves build metrics for a project
	GetBuildMetrics(ctx context.Context, projectID value_objects.ID) (*entities.BuildMetrics, error)

	// ListBuildEvents retrieves build events with filtering and pagination
	ListBuildEvents(ctx context.Context, filters ListBuildEventFilters) ([]*entities.BuildEvent, error)
}

// CreateBuildEventRequest represents a request to create a build event
type CreateBuildEventRequest struct {
	ProjectID       value_objects.ID
	EventType       entities.EventType
	Status          entities.BuildStatus
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
