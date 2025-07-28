package repository

import (
	"context"
	"errors"

	"github.com/dewisartika8/cicd-status-notifier-bot/internal/domain/entities"
	"github.com/dewisartika8/cicd-status-notifier-bot/internal/domain/ports"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// projectRepository implements the ProjectRepository interface
type projectRepository struct {
	db *gorm.DB
}

// NewProjectRepository creates a new project repository
func NewProjectRepository(db *gorm.DB) ports.ProjectRepository {
	return &projectRepository{db: db}
}

// Create creates a new project
func (r *projectRepository) Create(ctx context.Context, project *entities.Project) error {
	var model ProjectModel
	model.FromEntity(project)

	if err := r.db.WithContext(ctx).Create(&model).Error; err != nil {
		if isUniqueConstraintError(err) {
			return ports.ErrProjectAlreadyExists
		}
		return err
	}

	// Update entity with generated values
	*project = *model.ToEntity()
	return nil
}

// GetByID retrieves a project by its ID
func (r *projectRepository) GetByID(ctx context.Context, id uuid.UUID) (*entities.Project, error) {
	var model ProjectModel
	if err := r.db.WithContext(ctx).First(&model, "id = ?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ports.ErrProjectNotFound
		}
		return nil, err
	}

	return model.ToEntity(), nil
}

// GetByName retrieves a project by its name
func (r *projectRepository) GetByName(ctx context.Context, name string) (*entities.Project, error) {
	var model ProjectModel
	if err := r.db.WithContext(ctx).First(&model, "name = ?", name).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ports.ErrProjectNotFound
		}
		return nil, err
	}

	return model.ToEntity(), nil
}

// List retrieves all projects
func (r *projectRepository) List(ctx context.Context) ([]*entities.Project, error) {
	var models []ProjectModel
	if err := r.db.WithContext(ctx).Find(&models).Error; err != nil {
		return nil, err
	}

	projects := make([]*entities.Project, len(models))
	for i, model := range models {
		projects[i] = model.ToEntity()
	}

	return projects, nil
}

// Update updates an existing project
func (r *projectRepository) Update(ctx context.Context, project *entities.Project) error {
	var model ProjectModel
	model.FromEntity(project)

	result := r.db.WithContext(ctx).Save(&model)
	if result.Error != nil {
		if isUniqueConstraintError(result.Error) {
			return ports.ErrProjectAlreadyExists
		}
		return result.Error
	}

	if result.RowsAffected == 0 {
		return ports.ErrProjectNotFound
	}

	// Update entity with saved values
	*project = *model.ToEntity()
	return nil
}

// Delete deletes a project by its ID
func (r *projectRepository) Delete(ctx context.Context, id uuid.UUID) error {
	result := r.db.WithContext(ctx).Delete(&ProjectModel{}, "id = ?", id)
	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return ports.ErrProjectNotFound
	}

	return nil
}

// GetWithBuildEvents retrieves a project with its recent build events
func (r *projectRepository) GetWithBuildEvents(ctx context.Context, id uuid.UUID, limit int) (*entities.Project, []*entities.BuildEvent, error) {
	var projectModel ProjectModel
	if err := r.db.WithContext(ctx).First(&projectModel, "id = ?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil, ports.ErrProjectNotFound
		}
		return nil, nil, err
	}

	var buildEventModels []BuildEventModel
	if err := r.db.WithContext(ctx).
		Where("project_id = ?", id).
		Order("created_at DESC").
		Limit(limit).
		Find(&buildEventModels).Error; err != nil {
		return nil, nil, err
	}

	project := projectModel.ToEntity()
	buildEvents := make([]*entities.BuildEvent, len(buildEventModels))
	for i, model := range buildEventModels {
		buildEvents[i] = model.ToEntity()
	}

	return project, buildEvents, nil
}

// isUniqueConstraintError checks if the error is a unique constraint violation
func isUniqueConstraintError(err error) bool {
	if err == nil {
		return false
	}

	errStr := err.Error()
	return contains(errStr, "duplicate key value violates unique constraint") ||
		contains(errStr, "UNIQUE constraint failed") ||
		contains(errStr, "violates unique constraint")
}

// contains checks if a string contains a substring (case-insensitive)
func contains(s, substr string) bool {
	return len(s) >= len(substr) &&
		(s == substr || len(substr) == 0 ||
			findInString(s, substr) >= 0)
}

func findInString(s, substr string) int {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return i
		}
	}
	return -1
}
