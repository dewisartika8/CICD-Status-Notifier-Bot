package port

import (
	"context"

	"github.com/dewisartika8/cicd-status-notifier-bot/internal/core/build/domain"
	"github.com/dewisartika8/cicd-status-notifier-bot/internal/core/build/dto"
	"github.com/dewisartika8/cicd-status-notifier-bot/internal/core/shared/domain/value_objects"
)

// BuildEventRepository defines the contract for build event persistence
type BuildEventRepository interface {
	// Create creates a new build event
	Create(ctx context.Context, buildEvent *domain.BuildEvent) error

	// GetByID retrieves a build event by its ID
	GetByID(ctx context.Context, id value_objects.ID) (*domain.BuildEvent, error)

	// GetByProjectID retrieves build events for a specific project
	GetByProjectID(ctx context.Context, projectID value_objects.ID, filters dto.ListBuildEventFilters) ([]*domain.BuildEvent, error)

	// List retrieves build events with optional filtering
	List(ctx context.Context, filters dto.ListBuildEventFilters) ([]*domain.BuildEvent, error)

	// Update updates an existing build event
	Update(ctx context.Context, buildEvent *domain.BuildEvent) error

	// Delete deletes a build event by its ID
	Delete(ctx context.Context, id value_objects.ID) error

	// GetLatestByProjectID gets the latest build event for a project
	GetLatestByProjectID(ctx context.Context, projectID value_objects.ID) (*domain.BuildEvent, error)

	// Count returns the total number of build events
	Count(ctx context.Context, filters dto.ListBuildEventFilters) (int64, error)

	// GetBuildMetrics retrieves build metrics for a project
	GetBuildMetrics(ctx context.Context, projectID value_objects.ID) (*domain.BuildMetrics, error)
}
