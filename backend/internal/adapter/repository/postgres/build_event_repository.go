package postgres

import (
	"context"
	"errors"
	"time"

	"github.com/dewisartika8/cicd-status-notifier-bot/internal/core/build/domain"
	"github.com/dewisartika8/cicd-status-notifier-bot/internal/core/build/dto"
	"github.com/dewisartika8/cicd-status-notifier-bot/internal/core/build/port"
	"github.com/dewisartika8/cicd-status-notifier-bot/internal/core/shared/domain/value_objects"
	"github.com/dewisartika8/cicd-status-notifier-bot/pkg/exception"
	"gorm.io/gorm"
)

// buildEventRepositoryImpl implements the BuildEventRepository interface
type buildEventRepositoryImpl struct {
	db *gorm.DB
}

// NewBuildEventRepository creates a new build event repository
func NewBuildEventRepository(db *gorm.DB) port.BuildEventRepository {
	if db == nil {
		panic("database connection cannot be nil")
	}
	return &buildEventRepositoryImpl{db: db}
}

// Create creates a new build event
func (r *buildEventRepositoryImpl) Create(ctx context.Context, buildEvent *domain.BuildEvent) error {
	if buildEvent == nil {
		return errors.New("build event cannot be nil")
	}

	var model domain.BuildEventModel
	model.FromEntity(buildEvent)

	if err := r.db.WithContext(ctx).Create(&model).Error; err != nil {
		return err
	}

	return nil
}

// GetByID retrieves a build event by its ID
func (r *buildEventRepositoryImpl) GetByID(ctx context.Context, id value_objects.ID) (*domain.BuildEvent, error) {
	var model domain.BuildEventModel
	if err := r.db.WithContext(ctx).First(&model, queryByID, id.Value()).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, exception.ErrBuildEventNotFound
		}
		return nil, err
	}

	return model.ToEntity(), nil
}

// GetByProjectID retrieves build events for a specific project
func (r *buildEventRepositoryImpl) GetByProjectID(ctx context.Context, projectID value_objects.ID, filters dto.ListBuildEventFilters) ([]*domain.BuildEvent, error) {
	query := r.db.WithContext(ctx).Where(queryByProjectID, projectID.Value())

	query = r.applyFilters(query, filters)
	query = r.applyPagination(query, filters)

	var models []domain.BuildEventModel
	if err := query.Order(orderByCreatedAtDesc).Find(&models).Error; err != nil {
		return nil, err
	}

	return r.modelsToEntities(models), nil
}

// List retrieves build events with optional filtering
func (r *buildEventRepositoryImpl) List(ctx context.Context, filters dto.ListBuildEventFilters) ([]*domain.BuildEvent, error) {
	query := r.db.WithContext(ctx).Model(&domain.BuildEventModel{})

	if filters.ProjectID != nil {
		query = query.Where(queryByProjectID, filters.ProjectID.Value())
	}

	query = r.applyFilters(query, filters)
	query = r.applyPagination(query, filters)

	var models []domain.BuildEventModel
	if err := query.Order(orderByCreatedAtDesc).Find(&models).Error; err != nil {
		return nil, err
	}

	return r.modelsToEntities(models), nil
}

// Update updates an existing build event
func (r *buildEventRepositoryImpl) Update(ctx context.Context, buildEvent *domain.BuildEvent) error {
	if buildEvent == nil {
		return errors.New("build event cannot be nil")
	}

	var model domain.BuildEventModel
	model.FromEntity(buildEvent)

	result := r.db.WithContext(ctx).
		Where(queryByID, buildEvent.ID().Value()).
		Updates(&model)

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return exception.ErrBuildEventNotFound
	}

	return nil
}

// Delete deletes a build event by its ID
func (r *buildEventRepositoryImpl) Delete(ctx context.Context, id value_objects.ID) error {
	result := r.db.WithContext(ctx).
		Delete(&domain.BuildEventModel{}, queryByID, id.Value())

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return exception.ErrBuildEventNotFound
	}

	return nil
}

// GetLatestByProjectID gets the latest build event for a project
func (r *buildEventRepositoryImpl) GetLatestByProjectID(ctx context.Context, projectID value_objects.ID) (*domain.BuildEvent, error) {
	var model domain.BuildEventModel
	if err := r.db.WithContext(ctx).
		Where(queryByProjectID, projectID.Value()).
		Order(orderByCreatedAtDesc).
		First(&model).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, exception.ErrBuildEventNotFound
		}
		return nil, err
	}

	return model.ToEntity(), nil
}

// Count returns the total number of build events
func (r *buildEventRepositoryImpl) Count(ctx context.Context, filters dto.ListBuildEventFilters) (int64, error) {
	query := r.db.WithContext(ctx).Model(&domain.BuildEventModel{})

	if filters.ProjectID != nil {
		query = query.Where(queryByProjectID, filters.ProjectID.Value())
	}

	query = r.applyFilters(query, filters)

	var count int64
	if err := query.Count(&count).Error; err != nil {
		return 0, err
	}

	return count, nil
}

// GetBuildMetrics retrieves build metrics for a project
func (r *buildEventRepositoryImpl) GetBuildMetrics(ctx context.Context, projectID value_objects.ID) (*domain.BuildMetrics, error) {
	totalBuilds, err := r.getBuildCount(ctx, projectID, "")
	if err != nil {
		return nil, err
	}

	successfulBuilds, err := r.getBuildCount(ctx, projectID, string(domain.BuildStatusSuccess))
	if err != nil {
		return nil, err
	}

	failedBuilds, err := r.getFailedBuildCount(ctx, projectID)
	if err != nil {
		return nil, err
	}

	avgDuration, err := r.getAverageDuration(ctx, projectID)
	if err != nil {
		return nil, err
	}

	lastBuildStatus, lastBuildTime, err := r.getLatestBuildInfo(ctx, projectID)
	if err != nil {
		return nil, err
	}

	return domain.RestoreBuildMetrics(
		projectID,
		totalBuilds,
		successfulBuilds,
		failedBuilds,
		time.Duration(avgDuration)*time.Second,
		lastBuildStatus,
		lastBuildTime,
	), nil
}

// Helper methods for better code organization and reusability

func (r *buildEventRepositoryImpl) applyFilters(query *gorm.DB, filters dto.ListBuildEventFilters) *gorm.DB {
	if filters.EventType != nil {
		query = query.Where(queryByEventType, string(*filters.EventType))
	}
	if filters.Status != nil {
		query = query.Where(queryByStatus, string(*filters.Status))
	}
	if filters.Branch != nil {
		query = query.Where(queryByBranch, *filters.Branch)
	}
	if filters.DateFrom != nil {
		query = query.Where(queryCreatedAtGTE, *filters.DateFrom)
	}
	if filters.DateTo != nil {
		query = query.Where(queryCreatedAtLTE, *filters.DateTo)
	}
	return query
}

func (r *buildEventRepositoryImpl) applyPagination(query *gorm.DB, filters dto.ListBuildEventFilters) *gorm.DB {
	if filters.Limit > 0 {
		query = query.Limit(filters.Limit)
	}
	if filters.Offset > 0 {
		query = query.Offset(filters.Offset)
	}
	return query
}

func (r *buildEventRepositoryImpl) modelsToEntities(models []domain.BuildEventModel) []*domain.BuildEvent {
	buildEvents := make([]*domain.BuildEvent, len(models))
	for i, model := range models {
		buildEvents[i] = model.ToEntity()
	}
	return buildEvents
}

func (r *buildEventRepositoryImpl) getBuildCount(ctx context.Context, projectID value_objects.ID, status string) (int64, error) {
	query := r.db.WithContext(ctx).
		Model(&domain.BuildEventModel{}).
		Where(queryByProjectID, projectID.Value())

	if status != "" {
		query = query.Where(queryByStatus, status)
	}

	var count int64
	if err := query.Count(&count).Error; err != nil {
		return 0, err
	}

	return count, nil
}

func (r *buildEventRepositoryImpl) getFailedBuildCount(ctx context.Context, projectID value_objects.ID) (int64, error) {
	var count int64
	if err := r.db.WithContext(ctx).
		Model(&domain.BuildEventModel{}).
		Where("project_id = ? AND (status = ? OR status = ?)",
			projectID.Value(), string(domain.BuildStatusFailed), string(domain.BuildStatusCancelled)).
		Count(&count).Error; err != nil {
		return 0, err
	}

	return count, nil
}

func (r *buildEventRepositoryImpl) getAverageDuration(ctx context.Context, projectID value_objects.ID) (float64, error) {
	var avgDuration float64
	if err := r.db.WithContext(ctx).
		Model(&domain.BuildEventModel{}).
		Where("project_id = ? AND duration_seconds IS NOT NULL", projectID.Value()).
		Select("AVG(duration_seconds)").
		Scan(&avgDuration).Error; err != nil {
		return 0, err
	}

	return avgDuration, nil
}

func (r *buildEventRepositoryImpl) getLatestBuildInfo(ctx context.Context, projectID value_objects.ID) (domain.BuildStatus, value_objects.Timestamp, error) {
	var latestModel domain.BuildEventModel
	if err := r.db.WithContext(ctx).
		Where(queryByProjectID, projectID.Value()).
		Order(orderByCreatedAtDesc).
		First(&latestModel).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// Return zero values if no records found
			return domain.BuildStatus(""), value_objects.Timestamp{}, nil
		}
		return domain.BuildStatus(""), value_objects.Timestamp{}, err
	}

	lastBuildStatus := domain.BuildStatus(latestModel.Status)
	lastBuildTime := value_objects.NewTimestampFromTime(latestModel.CreatedAt)

	return lastBuildStatus, lastBuildTime, nil
}
