package repository

import (
	"context"
	"fmt"

	"github.com/dewisartika8/cicd-status-notifier-bot/internal/core/notification/domain"
	"github.com/dewisartika8/cicd-status-notifier-bot/internal/core/notification/port"
	"github.com/dewisartika8/cicd-status-notifier-bot/internal/core/shared/domain/value_objects"
	"gorm.io/gorm"
)

// RetryConfigurationRepository implements the retry configuration repository interface
type RetryConfigurationRepository struct {
	db *gorm.DB
}

// NewRetryConfigurationRepository creates a new retry configuration repository
func NewRetryConfigurationRepository(db *gorm.DB) port.RetryConfigurationRepository {
	return &RetryConfigurationRepository{
		db: db,
	}
}

// Create creates a new retry configuration
func (r *RetryConfigurationRepository) Create(ctx context.Context, config *domain.RetryConfiguration) error {
	model := &domain.RetryConfigurationModel{}
	model.FromEntity(config)

	if err := r.db.WithContext(ctx).Create(model).Error; err != nil {
		return fmt.Errorf("failed to create retry configuration: %w", err)
	}

	return nil
}

// GetByID retrieves a retry configuration by its ID
func (r *RetryConfigurationRepository) GetByID(ctx context.Context, id value_objects.ID) (*domain.RetryConfiguration, error) {
	var model domain.RetryConfigurationModel

	err := r.db.WithContext(ctx).Where("id = ?", id.String()).First(&model).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, domain.ErrRetryConfigurationNotFound
		}
		return nil, fmt.Errorf("failed to get retry configuration: %w", err)
	}

	return model.ToEntity(), nil
}

// GetByChannel retrieves retry configuration for a specific channel
func (r *RetryConfigurationRepository) GetByChannel(ctx context.Context, channel domain.NotificationChannel) (*domain.RetryConfiguration, error) {
	var model domain.RetryConfigurationModel

	err := r.db.WithContext(ctx).Where("channel = ? AND is_active = ?", string(channel), true).First(&model).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, domain.ErrRetryConfigurationNotFound
		}
		return nil, fmt.Errorf("failed to get retry configuration by channel: %w", err)
	}

	return model.ToEntity(), nil
}

// GetByProjectAndChannel retrieves retry configuration by project and channel
func (r *RetryConfigurationRepository) GetByProjectAndChannel(ctx context.Context, projectID value_objects.ID, channel domain.NotificationChannel) (*domain.RetryConfiguration, error) {
	var model domain.RetryConfigurationModel

	err := r.db.WithContext(ctx).Where("project_id = ? AND channel = ? AND is_active = ?", projectID.String(), string(channel), true).First(&model).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, domain.ErrRetryConfigurationNotFound
		}
		return nil, fmt.Errorf("failed to get retry configuration by project and channel: %w", err)
	}

	return model.ToEntity(), nil
}

// GetByProjectID retrieves all retry configurations for a project
func (r *RetryConfigurationRepository) GetByProjectID(ctx context.Context, projectID value_objects.ID) ([]*domain.RetryConfiguration, error) {
	var models []domain.RetryConfigurationModel

	err := r.db.WithContext(ctx).Where("project_id = ?", projectID.String()).Find(&models).Error
	if err != nil {
		return nil, fmt.Errorf("failed to get retry configurations by project: %w", err)
	}

	configs := make([]*domain.RetryConfiguration, len(models))
	for i, model := range models {
		configs[i] = model.ToEntity()
	}

	return configs, nil
}

// GetActiveConfigurations retrieves all active retry configurations
func (r *RetryConfigurationRepository) GetActiveConfigurations(ctx context.Context) ([]*domain.RetryConfiguration, error) {
	var models []domain.RetryConfigurationModel

	err := r.db.WithContext(ctx).Where("is_active = ?", true).Find(&models).Error
	if err != nil {
		return nil, fmt.Errorf("failed to get active retry configurations: %w", err)
	}

	configs := make([]*domain.RetryConfiguration, len(models))
	for i, model := range models {
		configs[i] = model.ToEntity()
	}

	return configs, nil
}

// Update updates an existing retry configuration
func (r *RetryConfigurationRepository) Update(ctx context.Context, config *domain.RetryConfiguration) error {
	model := &domain.RetryConfigurationModel{}
	model.FromEntity(config)

	result := r.db.WithContext(ctx).Where("id = ?", config.ID().String()).Updates(model)
	if result.Error != nil {
		return fmt.Errorf("failed to update retry configuration: %w", result.Error)
	}

	if result.RowsAffected == 0 {
		return domain.ErrRetryConfigurationNotFound
	}

	return nil
}

// Delete deletes a retry configuration by its ID
func (r *RetryConfigurationRepository) Delete(ctx context.Context, id value_objects.ID) error {
	result := r.db.WithContext(ctx).Where("id = ?", id.String()).Delete(&domain.RetryConfigurationModel{})
	if result.Error != nil {
		return fmt.Errorf("failed to delete retry configuration: %w", result.Error)
	}

	if result.RowsAffected == 0 {
		return domain.ErrRetryConfigurationNotFound
	}

	return nil
}

// Count returns the total number of retry configurations
func (r *RetryConfigurationRepository) Count(ctx context.Context, projectID *value_objects.ID, channel *domain.NotificationChannel, isActive *bool) (int64, error) {
	query := r.db.WithContext(ctx).Model(&domain.RetryConfigurationModel{})

	if projectID != nil {
		query = query.Where("project_id = ?", projectID.String())
	}

	if channel != nil {
		query = query.Where("channel = ?", string(*channel))
	}

	if isActive != nil {
		query = query.Where("is_active = ?", *isActive)
	}

	var count int64
	err := query.Count(&count).Error
	if err != nil {
		return 0, fmt.Errorf("failed to count retry configurations: %w", err)
	}

	return count, nil
}

// BulkCreate creates multiple retry configurations
func (r *RetryConfigurationRepository) BulkCreate(ctx context.Context, configs []*domain.RetryConfiguration) error {
	if len(configs) == 0 {
		return nil
	}

	models := make([]*domain.RetryConfigurationModel, 0, len(configs))
	for _, config := range configs {
		model := &domain.RetryConfigurationModel{}
		model.FromEntity(config)
		models = append(models, model)
	}

	if err := r.db.WithContext(ctx).Create(models).Error; err != nil {
		return fmt.Errorf("failed to bulk create retry configurations: %w", err)
	}

	return nil
}
