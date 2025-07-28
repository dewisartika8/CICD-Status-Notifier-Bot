package port

import (
	"context"

	"github.com/google/uuid"

	"github.com/dewisartika8/cicd-status-notifier-bot/internal/core/build/dto"
	project "github.com/dewisartika8/cicd-status-notifier-bot/internal/core/project/dto"
)

type BuildEventService interface {
	CreateBuildEvent(ctx context.Context, buildEvent *dto.BuildEvent)
	GetBuildEvent(ctx context.Context, id uuid.UUID) (*dto.BuildEvent, error)
	GetBuildEventsByProject(ctx context.Context, projectID uuid.UUID, limit, offset int) ([]*dto.BuildEvent, error)
	GetLatestBuildEvent(ctx context.Context, projectID uuid.UUID) (*dto.BuildEvent, error)
	GetProjectMetrics(ctx context.Context, projectID uuid.UUID) (*project.ProjectMetrics, error)
}
