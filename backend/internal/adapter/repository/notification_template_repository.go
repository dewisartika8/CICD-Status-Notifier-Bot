package repository

import (
	"context"
	"fmt"

	"github.com/dewisartika8/cicd-status-notifier-bot/internal/core/notification/domain"
	"github.com/dewisartika8/cicd-status-notifier-bot/internal/core/notification/port"
	"github.com/dewisartika8/cicd-status-notifier-bot/internal/core/shared/domain/value_objects"
	"gorm.io/gorm"
)

// NotificationTemplateRepository implements the notification template repository interface
type NotificationTemplateRepository struct {
	db *gorm.DB
}

// NewNotificationTemplateRepository creates a new notification template repository
func NewNotificationTemplateRepository(db *gorm.DB) port.NotificationTemplateRepository {
	return &NotificationTemplateRepository{
		db: db,
	}
}

// Create creates a new notification template
func (r *NotificationTemplateRepository) Create(ctx context.Context, template *domain.NotificationTemplate) error {
	model := &domain.NotificationTemplateModel{}
	model.FromEntity(template)

	if err := r.db.WithContext(ctx).Create(model).Error; err != nil {
		return fmt.Errorf("failed to create notification template: %w", err)
	}

	return nil
}

// GetByID retrieves a notification template by its ID
func (r *NotificationTemplateRepository) GetByID(ctx context.Context, id value_objects.ID) (*domain.NotificationTemplate, error) {
	var model domain.NotificationTemplateModel

	err := r.db.WithContext(ctx).Where(queryByID, id.String()).First(&model).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, domain.ErrNotificationTemplateNotFound
		}
		return nil, fmt.Errorf("failed to get notification template: %w", err)
	}

	return model.ToEntity(), nil
}

// GetByTypeAndChannel retrieves a notification template by type and channel
func (r *NotificationTemplateRepository) GetByTypeAndChannel(ctx context.Context, templateType domain.NotificationTemplateType, channel domain.NotificationChannel) (*domain.NotificationTemplate, error) {
	var model domain.NotificationTemplateModel

	err := r.db.WithContext(ctx).
		Where(queryByTemplateType, string(templateType)).
		Where(queryByChannel, string(channel)).
		Where(queryByIsActive, true).
		First(&model).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, domain.ErrNotificationTemplateNotFound
		}
		return nil, fmt.Errorf("failed to get notification template by type and channel: %w", err)
	}

	return model.ToEntity(), nil
}

// GetByType retrieves all notification templates for a specific type
func (r *NotificationTemplateRepository) GetByType(ctx context.Context, templateType domain.NotificationTemplateType) ([]*domain.NotificationTemplate, error) {
	var models []domain.NotificationTemplateModel

	err := r.db.WithContext(ctx).Where(queryByTemplateType, string(templateType)).Find(&models).Error
	if err != nil {
		return nil, fmt.Errorf("failed to get notification templates by type: %w", err)
	}

	templates := make([]*domain.NotificationTemplate, len(models))
	for i, model := range models {
		templates[i] = model.ToEntity()
	}

	return templates, nil
}

// GetByChannel retrieves all notification templates for a specific channel
func (r *NotificationTemplateRepository) GetByChannel(ctx context.Context, channel domain.NotificationChannel) ([]*domain.NotificationTemplate, error) {
	var models []domain.NotificationTemplateModel

	err := r.db.WithContext(ctx).Where(queryByChannel, string(channel)).Find(&models).Error
	if err != nil {
		return nil, fmt.Errorf("failed to get notification templates by channel: %w", err)
	}

	templates := make([]*domain.NotificationTemplate, len(models))
	for i, model := range models {
		templates[i] = model.ToEntity()
	}

	return templates, nil
}

// GetActiveTemplates retrieves all active notification templates
func (r *NotificationTemplateRepository) GetActiveTemplates(ctx context.Context) ([]*domain.NotificationTemplate, error) {
	var models []domain.NotificationTemplateModel

	err := r.db.WithContext(ctx).Where(queryByIsActive, true).Find(&models).Error
	if err != nil {
		return nil, fmt.Errorf("failed to get active notification templates: %w", err)
	}

	templates := make([]*domain.NotificationTemplate, len(models))
	for i, model := range models {
		templates[i] = model.ToEntity()
	}

	return templates, nil
}

// Update updates an existing notification template
func (r *NotificationTemplateRepository) Update(ctx context.Context, template *domain.NotificationTemplate) error {
	model := &domain.NotificationTemplateModel{}
	model.FromEntity(template)

	result := r.db.WithContext(ctx).Where(queryByID, template.ID().String()).Updates(model)
	if result.Error != nil {
		return fmt.Errorf("failed to update notification template: %w", result.Error)
	}

	if result.RowsAffected == 0 {
		return domain.ErrNotificationTemplateNotFound
	}

	return nil
}

// Delete deletes a notification template by its ID
func (r *NotificationTemplateRepository) Delete(ctx context.Context, id value_objects.ID) error {
	result := r.db.WithContext(ctx).Where(queryByID, id.String()).Delete(&domain.NotificationTemplateModel{})
	if result.Error != nil {
		return fmt.Errorf("failed to delete notification template: %w", result.Error)
	}

	if result.RowsAffected == 0 {
		return domain.ErrNotificationTemplateNotFound
	}

	return nil
}

// Count returns the total number of notification templates
func (r *NotificationTemplateRepository) Count(ctx context.Context, templateType *domain.NotificationTemplateType, channel *domain.NotificationChannel, isActive *bool) (int64, error) {
	query := r.db.WithContext(ctx).Model(&domain.NotificationTemplateModel{})

	if templateType != nil {
		query = query.Where(queryByTemplateType, string(*templateType))
	}

	if channel != nil {
		query = query.Where(queryByChannel, string(*channel))
	}

	if isActive != nil {
		query = query.Where(queryByIsActive, *isActive)
	}

	var count int64
	err := query.Count(&count).Error
	if err != nil {
		return 0, fmt.Errorf("failed to count notification templates: %w", err)
	}

	return count, nil
}
