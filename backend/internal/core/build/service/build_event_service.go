package service

import (
	"context"

	"github.com/dewisartika8/cicd-status-notifier-bot/internal/core/build/port"
	"github.com/dewisartika8/cicd-status-notifier-bot/internal/domain/entities"
	"github.com/google/uuid"
)

type Dep struct {
	BuildEventRepo port.BuildEventRepository
}

// BuildEventService handles build event business logic
type buildEventService struct {
	Dep
}

// NewBuildEventService creates a new build event service
func NewBuildEventService(dep Dep) port.BuildEventService {
	return &buildEventService{
		Dep: dep,
	}
}

// CreateBuildEvent creates a new build event
func (s *buildEventService) CreateBuildEvent(ctx context.Context, buildEvent *entities.BuildEvent) error {
	// Add business logic validation here
	if err := buildEvent.Validate(); err != nil {
		return err
	}

	return s.BuildEventRepo.Create(ctx, buildEvent)
}

// GetBuildEvent retrieves a build event by ID
func (s *buildEventService) GetBuildEvent(ctx context.Context, id uuid.UUID) (*entities.BuildEvent, error) {
	return s.BuildEventRepo.GetByID(ctx, id)
}

// GetBuildEventsByProject retrieves build events for a project
func (s *buildEventService) GetBuildEventsByProject(ctx context.Context, projectID uuid.UUID, limit, offset int) ([]*entities.BuildEvent, error) {
	return s.BuildEventRepo.GetByProjectID(ctx, projectID, limit, offset)
}

// GetLatestBuildEvent retrieves the latest build event for a project
func (s *buildEventService) GetLatestBuildEvent(ctx context.Context, projectID uuid.UUID) (*entities.BuildEvent, error) {
	return s.BuildEventRepo.GetLatestByProjectID(ctx, projectID)
}

// GetProjectMetrics retrieves metrics for a project
func (s *buildEventService) GetProjectMetrics(ctx context.Context, projectID uuid.UUID) (*ports.ProjectMetrics, error) {
	return s.BuildEventRepo.GetMetrics(ctx, projectID)
}
