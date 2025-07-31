package postgres

import (
	"context"
	"errors"

	"github.com/dewisartika8/cicd-status-notifier-bot/internal/core/project/domain"
	"github.com/dewisartika8/cicd-status-notifier-bot/internal/core/project/dto"
	"github.com/dewisartika8/cicd-status-notifier-bot/internal/core/project/port"
	"github.com/dewisartika8/cicd-status-notifier-bot/internal/core/shared/domain/value_objects"
	"github.com/dewisartika8/cicd-status-notifier-bot/pkg/exception"
	"gorm.io/gorm"
)

// projectRepository implements the ProjectRepository interface
type projectRepository struct {
	db *gorm.DB
}

// NewProjectRepository creates a new project repository
func NewProjectRepository(db *gorm.DB) port.ProjectRepository {
	return &projectRepository{db: db}
}

// Create creates a new project
func (r *projectRepository) Create(ctx context.Context, project *domain.Project) error {
	var projectModel domain.ProjectModel
	projectModel.FromEntity(project)

	if err := r.db.WithContext(ctx).Create(&projectModel).Error; err != nil {
		if isUniqueConstraintError(err) {
			return exception.ErrProjectAlreadyExists
		}
		return err
	}

	return nil
}

// GetByID retrieves a project by its ID
func (r *projectRepository) GetByID(ctx context.Context, id value_objects.ID) (*domain.Project, error) {
	var projectModel domain.ProjectModel
	if err := r.db.WithContext(ctx).Where(queryByID, id.String()).First(&projectModel).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, domain.ErrProjectNotFound
		}
		return nil, err
	}
	return projectModel.ToEntity()
}

// GetByName retrieves a project by its name
func (r *projectRepository) GetByName(ctx context.Context, name string) (*domain.Project, error) {
	var projectModel domain.ProjectModel
	if err := r.db.WithContext(ctx).Where(queryByName, name).First(&projectModel).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, domain.ErrProjectNotFound
		}
		return nil, err
	}
	return projectModel.ToEntity()
}

// GetByRepositoryURL retrieves a project by its repository URL
func (r *projectRepository) GetByRepositoryURL(ctx context.Context, repositoryURL string) (*domain.Project, error) {
	var projectModel domain.ProjectModel
	if err := r.db.WithContext(ctx).Where(queryByRepositoryURL, repositoryURL).First(&projectModel).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, domain.ErrProjectNotFound
		}
		return nil, err
	}
	return projectModel.ToEntity()
}

// List retrieves projects with optional filtering
func (r *projectRepository) List(ctx context.Context, filters dto.ListProjectFilters) ([]*domain.Project, error) {
	var projectModels []domain.ProjectModel

	// Build query with filters
	query := r.buildListQuery(ctx, filters)

	// Execute query
	if err := query.Find(&projectModels).Error; err != nil {
		return nil, err
	}

	// Convert models to entities
	return r.convertModelsToEntities(projectModels)
}

// buildListQuery constructs the query with all filters, sorting and pagination applied
func (r *projectRepository) buildListQuery(ctx context.Context, filters dto.ListProjectFilters) *gorm.DB {
	query := r.db.WithContext(ctx).Model(&domain.ProjectModel{})

	// Apply filters
	query = r.applyFilters(query, filters)

	// Apply sorting
	query = r.applySorting(query, filters)

	// Apply pagination
	query = r.applyPagination(query, filters)

	return query
}

// applyFilters applies all filter conditions to the query
func (r *projectRepository) applyFilters(query *gorm.DB, filters dto.ListProjectFilters) *gorm.DB {
	if filters.Status != nil {
		query = query.Where(queryByStatus, string(*filters.Status))
	}

	if filters.Name != nil {
		query = query.Where("name ILIKE ?", "%"+*filters.Name+"%")
	}

	if filters.RepositoryURL != nil {
		query = query.Where("repository_url ILIKE ?", "%"+*filters.RepositoryURL+"%")
	}

	if filters.HasTelegramChat != nil {
		if *filters.HasTelegramChat {
			query = query.Where("telegram_chat_id IS NOT NULL")
		} else {
			query = query.Where("telegram_chat_id IS NULL")
		}
	}

	return query
}

// applySorting applies sorting to the query based on provided sort criteria
func (r *projectRepository) applySorting(query *gorm.DB, filters dto.ListProjectFilters) *gorm.DB {
	const (
		sortFieldName      = "name"
		sortFieldCreatedAt = "created_at"
		sortOrderAsc       = "asc"
		sortOrderDesc      = "desc"
	)

	if filters.SortBy == nil {
		return query.Order(orderByCreatedAtDesc)
	}

	sortOrder := sortOrderDesc
	if filters.SortOrder != nil {
		sortOrder = *filters.SortOrder
	}

	switch *filters.SortBy {
	case sortFieldName:
		if sortOrder == sortOrderDesc {
			return query.Order("name DESC")
		}
		return query.Order(orderByNameAsc)
	case sortFieldCreatedAt:
		if sortOrder == sortOrderAsc {
			return query.Order("created_at ASC")
		}
		return query.Order(orderByCreatedAtDesc)
	default:
		return query.Order(orderByCreatedAtDesc)
	}
}

// applyPagination applies limit and offset to the query
func (r *projectRepository) applyPagination(query *gorm.DB, filters dto.ListProjectFilters) *gorm.DB {
	if filters.Limit != nil {
		query = query.Limit(*filters.Limit)
	}

	if filters.Offset != nil {
		query = query.Offset(*filters.Offset)
	}

	return query
}

// convertModelsToEntities converts database models to domain entities
func (r *projectRepository) convertModelsToEntities(projectModels []domain.ProjectModel) ([]*domain.Project, error) {
	if len(projectModels) == 0 {
		return make([]*domain.Project, 0), nil
	}

	projects := make([]*domain.Project, len(projectModels))
	for i, projectModel := range projectModels {
		entity, err := projectModel.ToEntity()
		if err != nil {
			return nil, err
		}
		projects[i] = entity
	}

	return projects, nil
}

// Update updates an existing project
func (r *projectRepository) Update(ctx context.Context, project *domain.Project) error {
	var projectModel domain.ProjectModel
	projectModel.FromEntity(project)

	result := r.db.WithContext(ctx).Save(&projectModel)
	if result.Error != nil {
		if isUniqueConstraintError(result.Error) {
			return domain.ErrProjectAlreadyExists
		}
		return result.Error
	}

	if result.RowsAffected == 0 {
		return domain.ErrProjectNotFound
	}

	return nil
}

// Delete deletes a project by its ID
func (r *projectRepository) Delete(ctx context.Context, id value_objects.ID) error {
	result := r.db.WithContext(ctx).Delete(&domain.ProjectModel{}, queryByID, id.String())
	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return domain.ErrProjectNotFound
	}

	return nil
}

// ExistsByName checks if a project with the given name exists
func (r *projectRepository) ExistsByName(ctx context.Context, name string) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&domain.ProjectModel{}).Where(queryByName, name).Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

// ExistsByRepositoryURL checks if a project with the given repository URL exists
func (r *projectRepository) ExistsByRepositoryURL(ctx context.Context, repositoryURL string) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&domain.ProjectModel{}).Where(queryByRepositoryURL, repositoryURL).Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

// Count returns the total number of projects
func (r *projectRepository) Count(ctx context.Context, filters dto.ListProjectFilters) (int64, error) {
	var count int64

	query := r.db.WithContext(ctx).Model(&domain.ProjectModel{})

	// Apply filters
	if filters.Status != nil {
		query = query.Where(queryByStatus, string(*filters.Status))
	}

	if filters.Name != nil {
		query = query.Where("name ILIKE ?", "%"+*filters.Name+"%")
	}

	if filters.RepositoryURL != nil {
		query = query.Where("repository_url ILIKE ?", "%"+*filters.RepositoryURL+"%")
	}

	if filters.HasTelegramChat != nil {
		if *filters.HasTelegramChat {
			query = query.Where("telegram_chat_id IS NOT NULL")
		} else {
			query = query.Where("telegram_chat_id IS NULL")
		}
	}

	err := query.Count(&count).Error
	return count, err
}

// GetActiveProjects retrieves all active projects
func (r *projectRepository) GetActiveProjects(ctx context.Context) ([]*domain.Project, error) {
	filters := dto.ListProjectFilters{
		Status: &[]domain.ProjectStatus{domain.ProjectStatusActive}[0],
	}
	return r.List(ctx, filters)
}

// GetProjectsWithTelegramChat retrieves projects that have telegram chat configured
func (r *projectRepository) GetProjectsWithTelegramChat(ctx context.Context) ([]*domain.Project, error) {
	filters := dto.ListProjectFilters{
		HasTelegramChat: &[]bool{true}[0],
	}
	return r.List(ctx, filters)
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
