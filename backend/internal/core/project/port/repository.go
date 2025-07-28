package port

import (
	"context"

	"github.com/dewisartika8/cicd-status-notifier-bot/internal/domain/entities"
	"github.com/google/uuid"
)

// ProjectRepository defines the contract for project data access
type ProjectRepository interface {
	// Create creates a new project
	Create(ctx context.Context, project *entities.Project) error

	// GetByID retrieves a project by its ID
	GetByID(ctx context.Context, id uuid.UUID) (*entities.Project, error)

	// GetByName retrieves a project by its name
	GetByName(ctx context.Context, name string) (*entities.Project, error)

	// List retrieves all projects
	List(ctx context.Context) ([]*entities.Project, error)

	// Update updates an existing project
	Update(ctx context.Context, project *entities.Project) error

	// Delete deletes a project by its ID
	Delete(ctx context.Context, id uuid.UUID) error

	// GetWithBuildEvents retrieves a project with its recent build events
	GetWithBuildEvents(ctx context.Context, id uuid.UUID, limit int) (*entities.Project, []*entities.BuildEvent, error)
}
