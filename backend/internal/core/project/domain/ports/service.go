package ports

import (
	"context"

	"github.com/dewisartika8/cicd-status-notifier-bot/internal/core/project/domain/entities"
	"github.com/dewisartika8/cicd-status-notifier-bot/internal/core/shared/domain/value_objects"
)

// ProjectService defines the contract for project business logic
type ProjectService interface {
	// CreateProject creates a new project
	CreateProject(ctx context.Context, req CreateProjectRequest) (*entities.Project, error)

	// GetProject retrieves a project by its ID
	GetProject(ctx context.Context, id value_objects.ID) (*entities.Project, error)

	// GetProjectByName retrieves a project by its name
	GetProjectByName(ctx context.Context, name string) (*entities.Project, error)

	// ListProjects retrieves all projects with optional filtering
	ListProjects(ctx context.Context, filters ListProjectFilters) ([]*entities.Project, error)

	// UpdateProject updates an existing project
	UpdateProject(ctx context.Context, id value_objects.ID, req UpdateProjectRequest) (*entities.Project, error)

	// ActivateProject activates a project
	ActivateProject(ctx context.Context, id value_objects.ID) error

	// DeactivateProject deactivates a project
	DeactivateProject(ctx context.Context, id value_objects.ID) error

	// ArchiveProject archives a project
	ArchiveProject(ctx context.Context, id value_objects.ID) error

	// DeleteProject deletes a project
	DeleteProject(ctx context.Context, id value_objects.ID) error

	// ProjectExists checks if a project exists by name
	ProjectExists(ctx context.Context, name string) (bool, error)
}

// CreateProjectRequest represents a request to create a project
type CreateProjectRequest struct {
	Name           string
	RepositoryURL  string
	WebhookSecret  string
	TelegramChatID *int64
}

// UpdateProjectRequest represents a request to update a project
type UpdateProjectRequest struct {
	Name           *string
	RepositoryURL  *string
	WebhookSecret  *string
	TelegramChatID *int64
}
