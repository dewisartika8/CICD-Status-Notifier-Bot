package repository

import (
	"context"
	"fmt"

	"github.com/dewisartika8/cicd-status-notifier-bot/internal/core/notification/domain"
	"github.com/dewisartika8/cicd-status-notifier-bot/internal/core/notification/port"
	"github.com/dewisartika8/cicd-status-notifier-bot/internal/core/shared/domain/value_objects"
	"gorm.io/gorm"
)

// NotificationLogRepository implements the notification log repository interface
type NotificationLogRepository struct {
	db *gorm.DB
}

// NewNotificationLogRepository creates a new notification log repository
func NewNotificationLogRepository(db *gorm.DB) port.NotificationLogRepository {
	return &NotificationLogRepository{
		db: db,
	}
}

// Create creates a new notification log
func (r *NotificationLogRepository) Create(ctx context.Context, log *domain.NotificationLog) error {
	model := &domain.NotificationLogModel{}
	model.FromEntity(log)

	if err := r.db.WithContext(ctx).Create(model).Error; err != nil {
		return fmt.Errorf("failed to create notification log: %w", err)
	}

	return nil
}

// GetByID retrieves a notification log by its ID
func (r *NotificationLogRepository) GetByID(ctx context.Context, id value_objects.ID) (*domain.NotificationLog, error) {
	var model domain.NotificationLogModel

	err := r.db.WithContext(ctx).Where(queryByID, id.String()).First(&model).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, domain.ErrNotificationLogNotFound
		}
		return nil, fmt.Errorf("failed to get notification log: %w", err)
	}

	return model.ToEntity(), nil
}

// GetByBuildEventID retrieves all notification logs for a build event
func (r *NotificationLogRepository) GetByBuildEventID(ctx context.Context, buildEventID value_objects.ID) ([]*domain.NotificationLog, error) {
	var models []domain.NotificationLogModel

	err := r.db.WithContext(ctx).Where(queryByBuildEventID, buildEventID.String()).Find(&models).Error
	if err != nil {
		return nil, fmt.Errorf("failed to get notification logs by build event ID: %w", err)
	}

	logs := make([]*domain.NotificationLog, len(models))
	for i, model := range models {
		logs[i] = model.ToEntity()
	}

	return logs, nil
}

// GetByProjectID retrieves notification logs for a specific project
func (r *NotificationLogRepository) GetByProjectID(ctx context.Context, projectID value_objects.ID, limit, offset int) ([]*domain.NotificationLog, error) {
	var models []domain.NotificationLogModel

	query := r.db.WithContext(ctx).Where(queryByProjectID, projectID.String())
	if limit > 0 {
		query = query.Limit(limit)
	}
	if offset > 0 {
		query = query.Offset(offset)
	}

	err := query.Find(&models).Error
	if err != nil {
		return nil, fmt.Errorf("failed to get notification logs by project ID: %w", err)
	}

	logs := make([]*domain.NotificationLog, len(models))
	for i, model := range models {
		logs[i] = model.ToEntity()
	}

	return logs, nil
}

// GetByRecipient retrieves notification logs for a specific recipient
func (r *NotificationLogRepository) GetByRecipient(ctx context.Context, recipient string, limit, offset int) ([]*domain.NotificationLog, error) {
	var models []domain.NotificationLogModel

	query := r.db.WithContext(ctx).Where(queryByRecipient, recipient)
	if limit > 0 {
		query = query.Limit(limit)
	}
	if offset > 0 {
		query = query.Offset(offset)
	}

	err := query.Find(&models).Error
	if err != nil {
		return nil, fmt.Errorf("failed to get notification logs by recipient: %w", err)
	}

	logs := make([]*domain.NotificationLog, len(models))
	for i, model := range models {
		logs[i] = model.ToEntity()
	}

	return logs, nil
}

// Update updates an existing notification log
func (r *NotificationLogRepository) Update(ctx context.Context, log *domain.NotificationLog) error {
	model := &domain.NotificationLogModel{}
	model.FromEntity(log)

	result := r.db.WithContext(ctx).Model(model).Where(queryByID, model.ID).Updates(model)
	if result.Error != nil {
		return fmt.Errorf("failed to update notification log: %w", result.Error)
	}

	if result.RowsAffected == 0 {
		return domain.ErrNotificationLogNotFound
	}

	return nil
}

// Delete deletes a notification log by its ID
func (r *NotificationLogRepository) Delete(ctx context.Context, id value_objects.ID) error {
	result := r.db.WithContext(ctx).Where(queryByID, id.String()).Delete(&domain.NotificationLogModel{})
	if result.Error != nil {
		return fmt.Errorf("failed to delete notification log: %w", result.Error)
	}

	if result.RowsAffected == 0 {
		return domain.ErrNotificationLogNotFound
	}

	return nil
}

// GetFailedNotifications retrieves failed notifications for retry
func (r *NotificationLogRepository) GetFailedNotifications(ctx context.Context, limit int) ([]*domain.NotificationLog, error) {
	var models []domain.NotificationLogModel

	query := r.db.WithContext(ctx).Where(queryByStatus, "FAILED")
	if limit > 0 {
		query = query.Limit(limit)
	}

	err := query.Find(&models).Error
	if err != nil {
		return nil, fmt.Errorf("failed to get failed notifications: %w", err)
	}

	logs := make([]*domain.NotificationLog, len(models))
	for i, model := range models {
		logs[i] = model.ToEntity()
	}

	return logs, nil
}

// GetPendingNotifications retrieves pending notifications
func (r *NotificationLogRepository) GetPendingNotifications(ctx context.Context, limit int) ([]*domain.NotificationLog, error) {
	var models []domain.NotificationLogModel

	query := r.db.WithContext(ctx).Where(queryByStatus, "PENDING")
	if limit > 0 {
		query = query.Limit(limit)
	}

	err := query.Find(&models).Error
	if err != nil {
		return nil, fmt.Errorf("failed to get pending notifications: %w", err)
	}

	logs := make([]*domain.NotificationLog, len(models))
	for i, model := range models {
		logs[i] = model.ToEntity()
	}

	return logs, nil
}

// Count returns the total number of notification logs matching the criteria
func (r *NotificationLogRepository) Count(ctx context.Context, projectID *value_objects.ID, status *domain.NotificationStatus) (int64, error) {
	var count int64

	query := r.db.WithContext(ctx).Model(&domain.NotificationLogModel{})

	if projectID != nil {
		query = query.Where(queryByProjectID, projectID.String())
	}

	if status != nil {
		query = query.Where(queryByStatus, string(*status))
	}

	err := query.Count(&count).Error
	if err != nil {
		return 0, fmt.Errorf("failed to count notification logs: %w", err)
	}

	return count, nil
}

// GetNotificationStats retrieves notification statistics for a project
func (r *NotificationLogRepository) GetNotificationStats(ctx context.Context, projectID value_objects.ID) (*domain.NotificationStats, error) {
	stats := domain.NewNotificationStats(projectID)

	// Get status counts
	var statusResults []struct {
		Status string
		Count  int64
	}

	err := r.db.WithContext(ctx).Model(&domain.NotificationLogModel{}).
		Select("status, COUNT(*) as count").
		Where(queryByProjectID, projectID.String()).
		Group("status").
		Scan(&statusResults).Error

	if err != nil {
		return nil, fmt.Errorf("failed to get notification status stats: %w", err)
	}

	for _, result := range statusResults {
		stats.UpdateStatusCount(domain.NotificationStatus(result.Status), result.Count)
	}

	// Get channel counts
	var channelResults []struct {
		Channel string
		Count   int64
	}

	err = r.db.WithContext(ctx).Model(&domain.NotificationLogModel{}).
		Select("channel, COUNT(*) as count").
		Where(queryByProjectID, projectID.String()).
		Group("channel").
		Scan(&channelResults).Error

	if err != nil {
		return nil, fmt.Errorf("failed to get notification channel stats: %w", err)
	}

	for _, result := range channelResults {
		if result.Channel != "" {
			stats.UpdateChannelCount(domain.NotificationChannel(result.Channel), result.Count)
		}
	}

	return stats, nil
}
