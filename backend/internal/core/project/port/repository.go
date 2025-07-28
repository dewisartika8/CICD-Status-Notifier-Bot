package port

import (
	"context"

	"github.com/dewisartika8/cicd-status-notifier-bot/internal/core/project/domain"
	"github.com/dewisartika8/cicd-status-notifier-bot/internal/core/project/dto"
	"github.com/dewisartika8/cicd-status-notifier-bot/internal/core/shared/domain/value_objects"
)

// ProjectRepository defines the contract for project persistence
type ProjectRepository interface {
	// Create creates a new project
	Create(ctx context.Context, project *domain.Project) error

	// GetByID retrieves a project by its ID
	GetByID(ctx context.Context, id value_objects.ID) (*domain.Project, error)

	// GetByName retrieves a project by its name
	GetByName(ctx context.Context, name string) (*domain.Project, error)

	// GetByRepositoryURL retrieves a project by its repository URL
	GetByRepositoryURL(ctx context.Context, repositoryURL string) (*domain.Project, error)

	// List retrieves projects with optional filtering
	List(ctx context.Context, filters dto.ListProjectFilters) ([]*domain.Project, error)

	// Update updates an existing project
	Update(ctx context.Context, project *domain.Project) error

	// Delete deletes a project by its ID
	Delete(ctx context.Context, id value_objects.ID) error

	// ExistsByName checks if a project with the given name exists
	ExistsByName(ctx context.Context, name string) (bool, error)

	// ExistsByRepositoryURL checks if a project with the given repository URL exists
	ExistsByRepositoryURL(ctx context.Context, repositoryURL string) (bool, error)

	// Count returns the total number of projects
	Count(ctx context.Context, filters dto.ListProjectFilters) (int64, error)

	// GetActiveProjects retrieves all active projects
	GetActiveProjects(ctx context.Context) ([]*domain.Project, error)

	// GetProjectsWithTelegramChat retrieves projects that have telegram chat configured
	GetProjectsWithTelegramChat(ctx context.Context) ([]*domain.Project, error)
}
