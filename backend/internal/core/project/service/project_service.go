package service

import (
	"context"

	"github.com/dewisartika8/cicd-status-notifier-bot/internal/core/project/domain"
	"github.com/dewisartika8/cicd-status-notifier-bot/internal/core/project/dto"
	"github.com/dewisartika8/cicd-status-notifier-bot/internal/core/project/port"
	"github.com/dewisartika8/cicd-status-notifier-bot/internal/core/shared/domain/value_objects"
)

type Dep struct {
	ProjectRepo port.ProjectRepository
}

// ProjectService implements the project business logic
type projectService struct {
	Dep
}

// NewProjectService creates a new project service
func NewProjectService(d Dep) port.ProjectService {
	return &projectService{
		Dep: d,
	}
}

// CreateProject creates a new project
func (s *projectService) CreateProject(ctx context.Context, req dto.CreateProjectRequest) (*domain.Project, error) {
	// Check if project with same name already exists
	if exists, err := s.ProjectRepo.ExistsByName(ctx, req.Name); err != nil {
		return nil, err
	} else if exists {
		return nil, domain.NewProjectAlreadyExistsError(req.Name)
	}

	// Check if project with same repository URL already exists
	if exists, err := s.ProjectRepo.ExistsByRepositoryURL(ctx, req.RepositoryURL); err != nil {
		return nil, err
	} else if exists {
		return nil, domain.ErrProjectAlreadyExists
	}

	// Create domain entity
	project, err := domain.NewProject(req.Name, req.RepositoryURL, req.WebhookSecret, req.TelegramChatID)
	if err != nil {
		return nil, err
	}

	// Persist to repository
	if err := s.ProjectRepo.Create(ctx, project); err != nil {
		return nil, err
	}

	return project, nil
}

// GetProject retrieves a project by its ID
func (s *projectService) GetProject(ctx context.Context, id value_objects.ID) (*domain.Project, error) {
	project, err := s.ProjectRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if project == nil {
		return nil, domain.NewProjectNotFoundError(id.String())
	}

	return project, nil
}

// GetProjectByName retrieves a project by its name
func (s *projectService) GetProjectByName(ctx context.Context, name string) (*domain.Project, error) {
	project, err := s.ProjectRepo.GetByName(ctx, name)
	if err != nil {
		return nil, err
	}

	if project == nil {
		return nil, domain.NewProjectNotFoundError(name)
	}

	return project, nil
}

// GetProjectByRepositoryURL retrieves a project by its repository URL
func (s *projectService) GetProjectByRepositoryURL(ctx context.Context, repositoryURL string) (*domain.Project, error) {
	project, err := s.ProjectRepo.GetByRepositoryURL(ctx, repositoryURL)
	if err != nil {
		return nil, err
	}

	if project == nil {
		return nil, domain.NewProjectNotFoundError(repositoryURL)
	}

	return project, nil
}

// UpdateProject updates an existing project with the provided request data.
// It validates that unique constraints are not violated before applying updates.
func (s *projectService) UpdateProject(ctx context.Context, id value_objects.ID, req dto.UpdateProjectRequest) (*domain.Project, error) {
	// Retrieve the existing project
	project, err := s.GetProject(ctx, id)
	if err != nil {
		return nil, err
	}

	// Apply updates with validation
	if err := s.updateProjectName(ctx, project, req.Name, id); err != nil {
		return nil, err
	}

	if err := s.updateProjectRepositoryURL(ctx, project, req.RepositoryURL, id); err != nil {
		return nil, err
	}

	if err := s.updateProjectWebhookSecret(project, req.WebhookSecret); err != nil {
		return nil, err
	}

	if err := s.updateProjectTelegramChatID(project, req.TelegramChatID); err != nil {
		return nil, err
	}

	// Persist all changes to the repository
	if err := s.ProjectRepo.Update(ctx, project); err != nil {
		return nil, err
	}

	return project, nil
}

// updateProjectName updates the project name if provided and validates uniqueness
func (s *projectService) updateProjectName(ctx context.Context, project *domain.Project, newName *string, currentID value_objects.ID) error {
	if newName == nil {
		return nil
	}

	// Validate that the new name doesn't conflict with existing projects
	if err := s.validateUniqueProjectName(ctx, *newName, currentID); err != nil {
		return err
	}

	return project.UpdateName(*newName)
}

// updateProjectRepositoryURL updates the repository URL if provided and validates uniqueness
func (s *projectService) updateProjectRepositoryURL(ctx context.Context, project *domain.Project, newURL *string, currentID value_objects.ID) error {
	if newURL == nil {
		return nil
	}

	// Validate that the new repository URL doesn't conflict with existing projects
	if err := s.validateUniqueRepositoryURL(ctx, *newURL, currentID); err != nil {
		return err
	}

	return project.UpdateRepositoryURL(*newURL)
}

// updateProjectWebhookSecret updates the webhook secret if provided
func (s *projectService) updateProjectWebhookSecret(project *domain.Project, newSecret *string) error {
	if newSecret == nil {
		return nil
	}

	return project.UpdateWebhookSecret(*newSecret)
}

// updateProjectTelegramChatID updates the telegram chat ID if provided
func (s *projectService) updateProjectTelegramChatID(project *domain.Project, newChatID *int64) error {
	if newChatID == nil {
		return nil
	}

	return project.UpdateTelegramChatID(newChatID)
}

// validateUniqueProjectName ensures the project name is unique across all projects except the current one
func (s *projectService) validateUniqueProjectName(ctx context.Context, name string, excludeID value_objects.ID) error {
	existingProject, err := s.ProjectRepo.GetByName(ctx, name)
	if err != nil {
		// If error is not "not found", return it
		return err
	}

	// If project exists and it's not the current project being updated
	if existingProject != nil && existingProject.ID() != excludeID {
		return domain.NewProjectAlreadyExistsError(name)
	}

	return nil
}

// validateUniqueRepositoryURL ensures the repository URL is unique across all projects except the current one
func (s *projectService) validateUniqueRepositoryURL(ctx context.Context, repositoryURL string, excludeID value_objects.ID) error {
	existingProject, err := s.ProjectRepo.GetByRepositoryURL(ctx, repositoryURL)
	if err != nil {
		// If error is not "not found", return it
		return err
	}

	// If project exists and it's not the current project being updated
	if existingProject != nil && existingProject.ID() != excludeID {
		return domain.NewProjectAlreadyExistsError(repositoryURL)
	}

	return nil
}

// UpdateProjectStatus updates the project status
func (s *projectService) UpdateProjectStatus(ctx context.Context, id value_objects.ID, status domain.ProjectStatus) (*domain.Project, error) {
	project, err := s.GetProject(ctx, id)
	if err != nil {
		return nil, err
	}

	project.SetStatus(status)

	if err := s.ProjectRepo.Update(ctx, project); err != nil {
		return nil, err
	}

	return project, nil
}

// DeleteProject deletes a project
func (s *projectService) DeleteProject(ctx context.Context, id value_objects.ID) error {
	// Check if project exists
	if _, err := s.GetProject(ctx, id); err != nil {
		return err
	}

	return s.ProjectRepo.Delete(ctx, id)
}

// ListProjects retrieves projects with filtering and pagination
func (s *projectService) ListProjects(ctx context.Context, filters dto.ListProjectFilters) ([]*domain.Project, error) {
	return s.ProjectRepo.List(ctx, filters)
}

// GetActiveProjects retrieves all active projects
func (s *projectService) GetActiveProjects(ctx context.Context) ([]*domain.Project, error) {
	return s.ProjectRepo.GetActiveProjects(ctx)
}

// GetProjectsWithTelegramChat retrieves projects that have telegram chat configured
func (s *projectService) GetProjectsWithTelegramChat(ctx context.Context) ([]*domain.Project, error) {
	return s.ProjectRepo.GetProjectsWithTelegramChat(ctx)
}

// ActivateProject activates a project
func (s *projectService) ActivateProject(ctx context.Context, id value_objects.ID) (*domain.Project, error) {
	return s.UpdateProjectStatus(ctx, id, domain.ProjectStatusActive)
}

// DeactivateProject deactivates a project
func (s *projectService) DeactivateProject(ctx context.Context, id value_objects.ID) (*domain.Project, error) {
	return s.UpdateProjectStatus(ctx, id, domain.ProjectStatusInactive)
}

// ArchiveProject archives a project
func (s *projectService) ArchiveProject(ctx context.Context, id value_objects.ID) (*domain.Project, error) {
	return s.UpdateProjectStatus(ctx, id, domain.ProjectStatusArchived)
}

// ValidateWebhookSecret validates webhook secret for a project
func (s *projectService) ValidateWebhookSecret(ctx context.Context, id value_objects.ID, secret string) (bool, error) {
	project, err := s.GetProject(ctx, id)
	if err != nil {
		return false, err
	}

	return project.ValidateWebhookSecret(secret), nil
}

// CountProjects returns the total number of projects with filters
func (s *projectService) CountProjects(ctx context.Context, filters dto.ListProjectFilters) (int64, error) {
	return s.ProjectRepo.Count(ctx, filters)
}
