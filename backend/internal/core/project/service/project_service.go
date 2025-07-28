package service

import (
	"context"

	"github.com/dewisartika8/cicd-status-notifier-bot/internal/domain/entities"
	"github.com/dewisartika8/cicd-status-notifier-bot/internal/domain/ports"
	"github.com/google/uuid"
)

// ProjectService handles project business logic
type ProjectService struct {
	repo ports.ProjectRepository
}

// NewProjectService creates a new project service
func NewProjectService(repo ports.ProjectRepository) *ProjectService {
	return &ProjectService{
		repo: repo,
	}
}

// CreateProject creates a new project
func (s *ProjectService) CreateProject(ctx context.Context, project *entities.Project) error {
	// Add business logic validation here
	if err := project.Validate(); err != nil {
		return err
	}

	return s.repo.Create(ctx, project)
}

// GetProject retrieves a project by ID
func (s *ProjectService) GetProject(ctx context.Context, id uuid.UUID) (*entities.Project, error) {
	return s.repo.GetByID(ctx, id)
}

// GetProjectByName retrieves a project by name
func (s *ProjectService) GetProjectByName(ctx context.Context, name string) (*entities.Project, error) {
	return s.repo.GetByName(ctx, name)
}

// ListProjects retrieves all projects
func (s *ProjectService) ListProjects(ctx context.Context) ([]*entities.Project, error) {
	return s.repo.List(ctx)
}

// UpdateProject updates an existing project
func (s *ProjectService) UpdateProject(ctx context.Context, project *entities.Project) error {
	// Add business logic validation here
	if err := project.Validate(); err != nil {
		return err
	}

	return s.repo.Update(ctx, project)
}

// DeleteProject deletes a project
func (s *ProjectService) DeleteProject(ctx context.Context, id uuid.UUID) error {
	return s.repo.Delete(ctx, id)
}
