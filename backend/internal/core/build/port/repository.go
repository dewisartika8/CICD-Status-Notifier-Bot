package port

import (
	"context"

	"github.com/dewisartika8/cicd-status-notifier-bot/internal/domain/entities"
	"github.com/google/uuid"
)

// BuildEventRepository defines the contract for build event data access
type BuildEventRepository interface {
	// Create creates a new build event
	Create(ctx context.Context, buildEvent *entities.BuildEvent) error

	// GetByID retrieves a build event by its ID
	GetByID(ctx context.Context, id uuid.UUID) (*entities.BuildEvent, error)

	// GetByProjectID retrieves build events for a specific project with pagination
	GetByProjectID(ctx context.Context, projectID uuid.UUID, limit, offset int) ([]*entities.BuildEvent, error)

	// GetLatestByProjectID retrieves the latest build event for a project
	GetLatestByProjectID(ctx context.Context, projectID uuid.UUID) (*entities.BuildEvent, error)

	// CountByStatus counts build events by status for a project
	CountByStatus(ctx context.Context, projectID uuid.UUID, status entities.BuildStatus) (int64, error)

	// GetMetrics retrieves build metrics for a project
	GetMetrics(ctx context.Context, projectID uuid.UUID) (*ProjectMetrics, error)
}
