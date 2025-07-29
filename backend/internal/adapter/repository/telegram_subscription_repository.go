package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/dewisartika8/cicd-status-notifier-bot/internal/core/notification/domain"
	"github.com/dewisartika8/cicd-status-notifier-bot/internal/core/notification/port"
	"github.com/dewisartika8/cicd-status-notifier-bot/internal/core/shared/domain/value_objects"
	"gorm.io/gorm"
)

// TelegramSubscriptionModel represents the GORM model for telegram subscriptions
type TelegramSubscriptionModel struct {
	ID        string `gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	ProjectID string `gorm:"not null;type:uuid"`
	ChatID    int64  `gorm:"not null"`
	IsActive  bool   `gorm:"not null;default:true"`
	CreatedAt int64  `gorm:"autoCreateTime"`
	UpdatedAt int64  `gorm:"autoUpdateTime"`
}

// TableName returns the table name for the TelegramSubscriptionModel
func (TelegramSubscriptionModel) TableName() string {
	return "telegram_subscriptions"
}

// ToEntity converts the model to domain entity
func (tsm *TelegramSubscriptionModel) ToEntity() (*domain.TelegramSubscription, error) {
	id, err := value_objects.NewIDFromString(tsm.ID)
	if err != nil {
		return nil, fmt.Errorf("invalid ID: %w", err)
	}

	projectID, err := value_objects.NewIDFromString(tsm.ProjectID)
	if err != nil {
		return nil, fmt.Errorf("invalid project ID: %w", err)
	}

	return domain.RestoreTelegramSubscription(domain.RestoreTelegramSubscriptionParams{
		ID:        id,
		ProjectID: projectID,
		ChatID:    tsm.ChatID,
		IsActive:  tsm.IsActive,
		CreatedAt: value_objects.NewTimestampFromTime(time.Unix(tsm.CreatedAt, 0)),
		UpdatedAt: value_objects.NewTimestampFromTime(time.Unix(tsm.UpdatedAt, 0)),
	}), nil
}

// FromEntity converts domain entity to model
func (tsm *TelegramSubscriptionModel) FromEntity(entity *domain.TelegramSubscription) {
	tsm.ID = entity.ID().String()
	tsm.ProjectID = entity.ProjectID().String()
	tsm.ChatID = entity.ChatID()
	tsm.IsActive = entity.IsActive()
	tsm.CreatedAt = entity.CreatedAt().Unix()
	tsm.UpdatedAt = entity.UpdatedAt().Unix()
}

const (
	queryTelegramByID        = "id = ?"
	queryTelegramByProjectID = "project_id = ?"
	queryTelegramByChatID    = "chat_id = ?"
	queryTelegramByActive    = "is_active = ?"

	errConvertToEntity = "failed to convert model to entity: %w"
)

// TelegramSubscriptionRepository implements the telegram subscription repository interface
type TelegramSubscriptionRepository struct {
	db *gorm.DB
}

// NewTelegramSubscriptionRepository creates a new telegram subscription repository
func NewTelegramSubscriptionRepository(db *gorm.DB) port.TelegramSubscriptionRepository {
	return &TelegramSubscriptionRepository{
		db: db,
	}
}

// Create creates a new telegram subscription
func (r *TelegramSubscriptionRepository) Create(ctx context.Context, subscription *domain.TelegramSubscription) error {
	model := &TelegramSubscriptionModel{}
	model.FromEntity(subscription)

	if err := r.db.WithContext(ctx).Create(model).Error; err != nil {
		return fmt.Errorf("failed to create telegram subscription: %w", err)
	}

	return nil
}

// GetByID retrieves a telegram subscription by its ID
func (r *TelegramSubscriptionRepository) GetByID(ctx context.Context, id value_objects.ID) (*domain.TelegramSubscription, error) {
	var model TelegramSubscriptionModel

	err := r.db.WithContext(ctx).Where(queryTelegramByID, id.String()).First(&model).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, domain.ErrTelegramSubscriptionNotFound
		}
		return nil, fmt.Errorf("failed to get telegram subscription: %w", err)
	}

	entity, err := model.ToEntity()
	if err != nil {
		return nil, fmt.Errorf(errConvertToEntity, err)
	}

	return entity, nil
}

// GetByProjectID retrieves telegram subscriptions for a specific project
func (r *TelegramSubscriptionRepository) GetByProjectID(ctx context.Context, projectID value_objects.ID) ([]*domain.TelegramSubscription, error) {
	var models []TelegramSubscriptionModel

	err := r.db.WithContext(ctx).Where(queryTelegramByProjectID, projectID.String()).Find(&models).Error
	if err != nil {
		return nil, fmt.Errorf("failed to get telegram subscriptions by project ID: %w", err)
	}

	subscriptions := make([]*domain.TelegramSubscription, len(models))
	for i, model := range models {
		entity, err := model.ToEntity()
		if err != nil {
			return nil, fmt.Errorf(errConvertToEntity, err)
		}
		subscriptions[i] = entity
	}

	return subscriptions, nil
}

// GetByChatID retrieves a telegram subscription by chat ID
func (r *TelegramSubscriptionRepository) GetByChatID(ctx context.Context, chatID int64) (*domain.TelegramSubscription, error) {
	var model TelegramSubscriptionModel

	err := r.db.WithContext(ctx).Where(queryTelegramByChatID, chatID).First(&model).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, domain.ErrTelegramSubscriptionNotFound
		}
		return nil, fmt.Errorf("failed to get telegram subscription by chat ID: %w", err)
	}

	entity, err := model.ToEntity()
	if err != nil {
		return nil, fmt.Errorf(errConvertToEntity, err)
	}

	return entity, nil
}

// GetByProjectAndChatID retrieves a specific subscription by project and chat ID
func (r *TelegramSubscriptionRepository) GetByProjectAndChatID(ctx context.Context, projectID value_objects.ID, chatID int64) (*domain.TelegramSubscription, error) {
	var model TelegramSubscriptionModel

	err := r.db.WithContext(ctx).
		Where(queryTelegramByProjectID+" AND "+queryTelegramByChatID, projectID.String(), chatID).
		First(&model).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, domain.ErrTelegramSubscriptionNotFound
		}
		return nil, fmt.Errorf("failed to get telegram subscription by project and chat ID: %w", err)
	}

	entity, err := model.ToEntity()
	if err != nil {
		return nil, fmt.Errorf(errConvertToEntity, err)
	}

	return entity, nil
}

// Update updates an existing telegram subscription
func (r *TelegramSubscriptionRepository) Update(ctx context.Context, subscription *domain.TelegramSubscription) error {
	model := &TelegramSubscriptionModel{}
	model.FromEntity(subscription)

	result := r.db.WithContext(ctx).Model(model).Where(queryTelegramByID, model.ID).Updates(model)
	if result.Error != nil {
		return fmt.Errorf("failed to update telegram subscription: %w", result.Error)
	}

	if result.RowsAffected == 0 {
		return domain.ErrTelegramSubscriptionNotFound
	}

	return nil
}

// Delete deletes a telegram subscription by its ID
func (r *TelegramSubscriptionRepository) Delete(ctx context.Context, id value_objects.ID) error {
	result := r.db.WithContext(ctx).Where(queryTelegramByID, id.String()).Delete(&TelegramSubscriptionModel{})
	if result.Error != nil {
		return fmt.Errorf("failed to delete telegram subscription: %w", result.Error)
	}

	if result.RowsAffected == 0 {
		return domain.ErrTelegramSubscriptionNotFound
	}

	return nil
}

// GetActiveSubscriptions retrieves all active telegram subscriptions
func (r *TelegramSubscriptionRepository) GetActiveSubscriptions(ctx context.Context) ([]*domain.TelegramSubscription, error) {
	var models []TelegramSubscriptionModel

	err := r.db.WithContext(ctx).Where(queryTelegramByActive, true).Find(&models).Error
	if err != nil {
		return nil, fmt.Errorf("failed to get active telegram subscriptions: %w", err)
	}

	subscriptions := make([]*domain.TelegramSubscription, len(models))
	for i, model := range models {
		entity, err := model.ToEntity()
		if err != nil {
			return nil, fmt.Errorf(errConvertToEntity, err)
		}
		subscriptions[i] = entity
	}

	return subscriptions, nil
}

// GetActiveSubscriptionsByProject retrieves active subscriptions for a project
func (r *TelegramSubscriptionRepository) GetActiveSubscriptionsByProject(ctx context.Context, projectID value_objects.ID) ([]*domain.TelegramSubscription, error) {
	var models []TelegramSubscriptionModel

	err := r.db.WithContext(ctx).
		Where(queryTelegramByProjectID+" AND "+queryTelegramByActive, projectID.String(), true).
		Find(&models).Error
	if err != nil {
		return nil, fmt.Errorf("failed to get active telegram subscriptions by project: %w", err)
	}

	subscriptions := make([]*domain.TelegramSubscription, len(models))
	for i, model := range models {
		entity, err := model.ToEntity()
		if err != nil {
			return nil, fmt.Errorf(errConvertToEntity, err)
		}
		subscriptions[i] = entity
	}

	return subscriptions, nil
}

// ExistsByProjectAndChatID checks if a subscription exists for project and chat
func (r *TelegramSubscriptionRepository) ExistsByProjectAndChatID(ctx context.Context, projectID value_objects.ID, chatID int64) (bool, error) {
	var count int64

	err := r.db.WithContext(ctx).Model(&TelegramSubscriptionModel{}).
		Where(queryTelegramByProjectID+" AND "+queryTelegramByChatID, projectID.String(), chatID).
		Count(&count).Error
	if err != nil {
		return false, fmt.Errorf("failed to check telegram subscription existence: %w", err)
	}

	return count > 0, nil
}

// Count returns the total number of telegram subscriptions
func (r *TelegramSubscriptionRepository) Count(ctx context.Context, projectID *value_objects.ID, isActive *bool) (int64, error) {
	var count int64

	query := r.db.WithContext(ctx).Model(&TelegramSubscriptionModel{})

	if projectID != nil {
		query = query.Where(queryTelegramByProjectID, projectID.String())
	}

	if isActive != nil {
		query = query.Where(queryTelegramByActive, *isActive)
	}

	err := query.Count(&count).Error
	if err != nil {
		return 0, fmt.Errorf("failed to count telegram subscriptions: %w", err)
	}

	return count, nil
}
