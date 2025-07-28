package ports

import (
	"context"

	"github.com/dewisartika8/cicd-status-notifier-bot/internal/core/project/domain/entities"
	"github.com/dewisartika8/cicd-status-notifier-bot/internal/core/shared/domain/value_objects"
)

// ProjectRepository defines the contract for project persistence
type ProjectRepository interface {
	// Create creates a new project
	Create(ctx context.Context, project *entities.Project) error

	// GetByID retrieves a project by its ID
	GetByID(ctx context.Context, id value_objects.ID) (*entities.Project, error)

	// GetByName retrieves a project by its name
	GetByName(ctx context.Context, name string) (*entities.Project, error)

	// List retrieves all projects with optional filtering
	List(ctx context.Context, filters ListProjectFilters) ([]*entities.Project, error)

	// Update updates an existing project
	Update(ctx context.Context, project *entities.Project) error

	// Delete soft deletes a project by its ID
	Delete(ctx context.Context, id value_objects.ID) error

	// ExistsByName checks if a project with the given name exists
	ExistsByName(ctx context.Context, name string) (bool, error)

	// Count returns the total number of projects
	Count(ctx context.Context, filters ListProjectFilters) (int64, error)
}

// ListProjectFilters defines filters for listing projects
type ListProjectFilters struct {
	Status   *entities.ProjectStatus
	Name     *string
	Limit    int
	Offset   int
	OrderBy  string
	OrderDir string // "asc" or "desc"
}
