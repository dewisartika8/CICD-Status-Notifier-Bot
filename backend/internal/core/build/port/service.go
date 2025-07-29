package port

import (
	"context"

	"github.com/dewisartika8/cicd-status-notifier-bot/internal/core/build/domain"
	"github.com/dewisartika8/cicd-status-notifier-bot/internal/core/build/dto"
	"github.com/dewisartika8/cicd-status-notifier-bot/internal/core/shared/domain/value_objects"
)

// BuildEventService defines the contract for build event business logic
type BuildEventService interface {
	// CreateBuildEvent creates a new build event
	CreateBuildEvent(ctx context.Context, req dto.CreateBuildEventRequest) (*domain.BuildEvent, error)

	// ProcessWebhookEvent processes a webhook event and creates build events
	ProcessWebhookEvent(ctx context.Context, req dto.ProcessWebhookRequest) ([]*domain.BuildEvent, error)

	// GetBuildEvent retrieves a build event by its ID
	GetBuildEvent(ctx context.Context, id value_objects.ID) (*domain.BuildEvent, error)

	// GetBuildEventsByProject retrieves build events for a specific project
	GetBuildEventsByProject(ctx context.Context, projectID value_objects.ID, filters dto.ListBuildEventFilters) ([]*domain.BuildEvent, error)

	// UpdateBuildEventStatus updates the status of a build event
	UpdateBuildEventStatus(ctx context.Context, id value_objects.ID, status domain.BuildStatus, duration *int) error

	// GetLatestBuildEvent gets the latest build event for a project
	GetLatestBuildEvent(ctx context.Context, projectID value_objects.ID) (*domain.BuildEvent, error)

	// GetBuildMetrics retrieves build metrics for a project
	GetBuildMetrics(ctx context.Context, projectID value_objects.ID) (*domain.BuildMetrics, error)

	// ListBuildEvents retrieves build events with filtering and pagination
	ListBuildEvents(ctx context.Context, filters dto.ListBuildEventFilters) ([]*domain.BuildEvent, error)
}
