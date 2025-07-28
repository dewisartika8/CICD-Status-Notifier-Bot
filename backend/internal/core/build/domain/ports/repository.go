package ports

import (
	"context"
	"time"

	"github.com/dewisartika8/cicd-status-notifier-bot/internal/core/build/domain/entities"
	"github.com/dewisartika8/cicd-status-notifier-bot/internal/core/shared/domain/value_objects"
)

// BuildEventRepository defines the contract for build event persistence
type BuildEventRepository interface {
	// Create creates a new build event
	Create(ctx context.Context, buildEvent *entities.BuildEvent) error

	// GetByID retrieves a build event by its ID
	GetByID(ctx context.Context, id value_objects.ID) (*entities.BuildEvent, error)

	// GetByProjectID retrieves build events for a specific project
	GetByProjectID(ctx context.Context, projectID value_objects.ID, filters ListBuildEventFilters) ([]*entities.BuildEvent, error)

	// List retrieves build events with optional filtering
	List(ctx context.Context, filters ListBuildEventFilters) ([]*entities.BuildEvent, error)

	// Update updates an existing build event
	Update(ctx context.Context, buildEvent *entities.BuildEvent) error

	// Delete deletes a build event by its ID
	Delete(ctx context.Context, id value_objects.ID) error

	// GetLatestByProjectID gets the latest build event for a project
	GetLatestByProjectID(ctx context.Context, projectID value_objects.ID) (*entities.BuildEvent, error)

	// Count returns the total number of build events
	Count(ctx context.Context, filters ListBuildEventFilters) (int64, error)

	// GetBuildMetrics retrieves build metrics for a project
	GetBuildMetrics(ctx context.Context, projectID value_objects.ID) (*entities.BuildMetrics, error)
}

// ListBuildEventFilters defines filters for listing build events
type ListBuildEventFilters struct {
	ProjectID *value_objects.ID
	EventType *entities.EventType
	Status    *entities.BuildStatus
	Branch    *string
	DateFrom  *time.Time
	DateTo    *time.Time
	Limit     int
	Offset    int
	OrderBy   string
	OrderDir  string // "asc" or "desc"
}
