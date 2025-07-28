package repository

import (
	"context"
	"errors"

	"github.com/dewisartika8/cicd-status-notifier-bot/internal/domain/entities"
	"github.com/dewisartika8/cicd-status-notifier-bot/internal/domain/ports"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// buildEventRepository implements the BuildEventRepository interface
type buildEventRepository struct {
	db *gorm.DB
}

// NewBuildEventRepository creates a new build event repository
func NewBuildEventRepository(db *gorm.DB) ports.BuildEventRepository {
	return &buildEventRepository{db: db}
}

// Create creates a new build event
func (r *buildEventRepository) Create(ctx context.Context, buildEvent *entities.BuildEvent) error {
	var model BuildEventModel
	model.FromEntity(buildEvent)

	if err := r.db.WithContext(ctx).Create(&model).Error; err != nil {
		return err
	}

	// Update entity with generated values
	*buildEvent = *model.ToEntity()
	return nil
}

// GetByID retrieves a build event by its ID
func (r *buildEventRepository) GetByID(ctx context.Context, id uuid.UUID) (*entities.BuildEvent, error) {
	var model BuildEventModel
	if err := r.db.WithContext(ctx).First(&model, "id = ?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ports.ErrBuildEventNotFound
		}
		return nil, err
	}

	return model.ToEntity(), nil
}

// GetByProjectID retrieves build events for a specific project with pagination
func (r *buildEventRepository) GetByProjectID(ctx context.Context, projectID uuid.UUID, limit, offset int) ([]*entities.BuildEvent, error) {
	var models []BuildEventModel
	if err := r.db.WithContext(ctx).
		Where("project_id = ?", projectID).
		Order("created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&models).Error; err != nil {
		return nil, err
	}

	buildEvents := make([]*entities.BuildEvent, len(models))
	for i, model := range models {
		buildEvents[i] = model.ToEntity()
	}

	return buildEvents, nil
}

// GetLatestByProjectID retrieves the latest build event for a project
func (r *buildEventRepository) GetLatestByProjectID(ctx context.Context, projectID uuid.UUID) (*entities.BuildEvent, error) {
	var model BuildEventModel
	if err := r.db.WithContext(ctx).
		Where("project_id = ?", projectID).
		Order("created_at DESC").
		First(&model).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ports.ErrBuildEventNotFound
		}
		return nil, err
	}

	return model.ToEntity(), nil
}

// CountByStatus counts build events by status for a project
func (r *buildEventRepository) CountByStatus(ctx context.Context, projectID uuid.UUID, status entities.BuildStatus) (int64, error) {
	var count int64
	if err := r.db.WithContext(ctx).
		Model(&BuildEventModel{}).
		Where("project_id = ? AND status = ?", projectID, string(status)).
		Count(&count).Error; err != nil {
		return 0, err
	}

	return count, nil
}

// GetMetrics retrieves build metrics for a project
func (r *buildEventRepository) GetMetrics(ctx context.Context, projectID uuid.UUID) (*ports.ProjectMetrics, error) {
	// Get total builds count
	var totalBuilds int64
	if err := r.db.WithContext(ctx).
		Model(&BuildEventModel{}).
		Where("project_id = ?", projectID).
		Count(&totalBuilds).Error; err != nil {
		return nil, err
	}

	if totalBuilds == 0 {
		return &ports.ProjectMetrics{
			ProjectID:      projectID,
			TotalBuilds:    0,
			SuccessRate:    0,
			FailureRate:    0,
			AvgDuration:    0,
			BuildsByStatus: make(map[entities.BuildStatus]int64),
			BuildsByType:   make(map[entities.EventType]int64),
		}, nil
	}

	// Get builds by status
	var statusResults []struct {
		Status string
		Count  int64
	}
	if err := r.db.WithContext(ctx).
		Model(&BuildEventModel{}).
		Select("status, COUNT(*) as count").
		Where("project_id = ?", projectID).
		Group("status").
		Scan(&statusResults).Error; err != nil {
		return nil, err
	}

	buildsByStatus := make(map[entities.BuildStatus]int64)
	for _, result := range statusResults {
		buildsByStatus[entities.BuildStatus(result.Status)] = result.Count
	}

	// Get builds by type
	var typeResults []struct {
		EventType string
		Count     int64
	}
	if err := r.db.WithContext(ctx).
		Model(&BuildEventModel{}).
		Select("event_type, COUNT(*) as count").
		Where("project_id = ?", projectID).
		Group("event_type").
		Scan(&typeResults).Error; err != nil {
		return nil, err
	}

	buildsByType := make(map[entities.EventType]int64)
	for _, result := range typeResults {
		buildsByType[entities.EventType(result.EventType)] = result.Count
	}

	// Calculate success and failure rates
	successCount := buildsByStatus[entities.BuildStatusSuccess]
	failureCount := buildsByStatus[entities.BuildStatusFailed]

	successRate := float64(successCount) / float64(totalBuilds) * 100
	failureRate := float64(failureCount) / float64(totalBuilds) * 100

	// Calculate average duration
	var avgDuration float64
	if err := r.db.WithContext(ctx).
		Model(&BuildEventModel{}).
		Select("AVG(duration_seconds) as avg_duration").
		Where("project_id = ? AND duration_seconds IS NOT NULL", projectID).
		Scan(&avgDuration).Error; err != nil {
		return nil, err
	}

	return &ports.ProjectMetrics{
		ProjectID:      projectID,
		TotalBuilds:    totalBuilds,
		SuccessRate:    successRate,
		FailureRate:    failureRate,
		AvgDuration:    avgDuration,
		BuildsByStatus: buildsByStatus,
		BuildsByType:   buildsByType,
	}, nil
}
