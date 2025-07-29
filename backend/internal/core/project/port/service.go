package port

import (
	"context"

	"github.com/dewisartika8/cicd-status-notifier-bot/internal/core/project/domain"
	"github.com/dewisartika8/cicd-status-notifier-bot/internal/core/project/dto"
	"github.com/dewisartika8/cicd-status-notifier-bot/internal/core/shared/domain/value_objects"
)

// ProjectService defines the contract for project business logic
type ProjectService interface {
	// CreateProject creates a new project
	CreateProject(ctx context.Context, req dto.CreateProjectRequest) (*domain.Project, error)

	// GetProject retrieves a project by its ID
	GetProject(ctx context.Context, id value_objects.ID) (*domain.Project, error)

	// GetProjectByName retrieves a project by its name
	GetProjectByName(ctx context.Context, name string) (*domain.Project, error)

	// GetProjectByRepositoryURL retrieves a project by its repository URL
	GetProjectByRepositoryURL(ctx context.Context, repositoryURL string) (*domain.Project, error)

	// UpdateProject updates an existing project
	UpdateProject(ctx context.Context, id value_objects.ID, req dto.UpdateProjectRequest) (*domain.Project, error)

	// UpdateProjectStatus updates the project status
	UpdateProjectStatus(ctx context.Context, id value_objects.ID, status domain.ProjectStatus) (*domain.Project, error)

	// DeleteProject deletes a project
	DeleteProject(ctx context.Context, id value_objects.ID) error

	// ListProjects retrieves projects with filtering and pagination
	ListProjects(ctx context.Context, filters dto.ListProjectFilters) ([]*domain.Project, error)

	// GetActiveProjects retrieves all active projects
	GetActiveProjects(ctx context.Context) ([]*domain.Project, error)

	// GetProjectsWithTelegramChat retrieves projects that have telegram chat configured
	GetProjectsWithTelegramChat(ctx context.Context) ([]*domain.Project, error)

	// ActivateProject activates a project
	ActivateProject(ctx context.Context, id value_objects.ID) (*domain.Project, error)

	// DeactivateProject deactivates a project
	DeactivateProject(ctx context.Context, id value_objects.ID) (*domain.Project, error)

	// ArchiveProject archives a project
	ArchiveProject(ctx context.Context, id value_objects.ID) (*domain.Project, error)

	// ValidateWebhookSecret validates webhook secret for a project
	ValidateWebhookSecret(ctx context.Context, id value_objects.ID, secret string) (bool, error)

	// CountProjects returns the total number of projects with filters
	CountProjects(ctx context.Context, filters dto.ListProjectFilters) (int64, error)
}
