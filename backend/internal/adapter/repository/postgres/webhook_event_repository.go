package postgres

import (
	"context"
	"errors"

	"github.com/dewisartika8/cicd-status-notifier-bot/internal/core/shared/domain/value_objects"
	"github.com/dewisartika8/cicd-status-notifier-bot/internal/core/webhook/domain"
	"github.com/dewisartika8/cicd-status-notifier-bot/internal/core/webhook/port"
	"github.com/dewisartika8/cicd-status-notifier-bot/pkg/exception"
	"gorm.io/gorm"
)

const (
	errWebhookEventNotFound = "webhook event not found"
	errCodeWebhookNotFound  = "WEBHOOK_EVENT_NOT_FOUND"
	errCodeDBCreateFailed   = "DB_CREATE_FAILED"
	errCodeDBReadFailed     = "DB_READ_FAILED"
	errCodeDBUpdateFailed   = "DB_UPDATE_FAILED"
	errCodeDBDeleteFailed   = "DB_DELETE_FAILED"
)

// webhookEventRepository implements WebhookEventRepository interface
type webhookEventRepository struct {
	db *gorm.DB
}

// NewWebhookEventRepository creates a new webhook event repository
func NewWebhookEventRepository(db *gorm.DB) port.WebhookEventRepository {
	return &webhookEventRepository{
		db: db,
	}
}

// Create stores a new webhook event
func (r *webhookEventRepository) Create(ctx context.Context, event *domain.WebhookEvent) error {
	model := &domain.WebhookEventModel{}
	model.FromEntity(event)

	if err := r.db.WithContext(ctx).Create(model).Error; err != nil {
		return exception.NewDomainErrorWithCause(errCodeDBCreateFailed, "failed to create webhook event", err)
	}

	return nil
}

// GetByID retrieves a webhook event by its ID
func (r *webhookEventRepository) GetByID(ctx context.Context, id value_objects.ID) (*domain.WebhookEvent, error) {
	var model domain.WebhookEventModel
	err := r.db.WithContext(ctx).Where(queryByID, id.String()).First(&model).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, exception.NewDomainError(errCodeWebhookNotFound, errWebhookEventNotFound)
		}
		return nil, exception.NewDomainErrorWithCause(errCodeDBReadFailed, "failed to get webhook event", err)
	}

	return model.ToEntity()
}

// GetByDeliveryID retrieves a webhook event by its delivery ID
func (r *webhookEventRepository) GetByDeliveryID(ctx context.Context, deliveryID string) (*domain.WebhookEvent, error) {
	var model domain.WebhookEventModel
	err := r.db.WithContext(ctx).Where(queryByDeliveryID, deliveryID).First(&model).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, exception.NewDomainError(errCodeWebhookNotFound, errWebhookEventNotFound)
		}
		return nil, exception.NewDomainErrorWithCause(errCodeDBReadFailed, "failed to get webhook event", err)
	}

	return model.ToEntity()
}

// GetByProjectID retrieves webhook events for a specific project
func (r *webhookEventRepository) GetByProjectID(ctx context.Context, projectID value_objects.ID, limit, offset int) ([]*domain.WebhookEvent, error) {
	var models []domain.WebhookEventModel
	err := r.db.WithContext(ctx).
		Where(queryByProjectID, projectID.String()).
		Order("created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&models).Error

	if err != nil {
		return nil, exception.NewDomainErrorWithCause(errCodeDBReadFailed, "failed to get webhook events", err)
	}

	events := make([]*domain.WebhookEvent, len(models))
	for i, model := range models {
		event, err := model.ToEntity()
		if err != nil {
			return nil, err
		}
		events[i] = event
	}

	return events, nil
}

// Update updates an existing webhook event
func (r *webhookEventRepository) Update(ctx context.Context, event *domain.WebhookEvent) error {
	model := &domain.WebhookEventModel{}
	model.FromEntity(event)

	result := r.db.WithContext(ctx).Where(queryByID, event.ID().String()).Updates(model)
	if result.Error != nil {
		return exception.NewDomainErrorWithCause(errCodeDBUpdateFailed, "failed to update webhook event", result.Error)
	}

	if result.RowsAffected == 0 {
		return exception.NewDomainError(errCodeWebhookNotFound, errWebhookEventNotFound)
	}

	return nil
}

// Delete removes a webhook event
func (r *webhookEventRepository) Delete(ctx context.Context, id value_objects.ID) error {
	result := r.db.WithContext(ctx).Where(queryByID, id.String()).Delete(&domain.WebhookEventModel{})
	if result.Error != nil {
		return exception.NewDomainErrorWithCause(errCodeDBDeleteFailed, "failed to delete webhook event", result.Error)
	}

	if result.RowsAffected == 0 {
		return exception.NewDomainError(errCodeWebhookNotFound, errWebhookEventNotFound)
	}

	return nil
}

// ExistsByDeliveryID checks if a webhook event with the given delivery ID exists
func (r *webhookEventRepository) ExistsByDeliveryID(ctx context.Context, deliveryID string) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).
		Model(&domain.WebhookEventModel{}).
		Where(queryByDeliveryID, deliveryID).
		Count(&count).Error

	if err != nil {
		return false, exception.NewDomainErrorWithCause(errCodeDBReadFailed, "failed to check webhook event existence", err)
	}

	return count > 0, nil
}

// GetUnprocessedEvents retrieves unprocessed webhook events
func (r *webhookEventRepository) GetUnprocessedEvents(ctx context.Context, limit int) ([]*domain.WebhookEvent, error) {
	var models []domain.WebhookEventModel
	err := r.db.WithContext(ctx).
		Where("processed_at IS NULL").
		Order("created_at ASC").
		Limit(limit).
		Find(&models).Error

	if err != nil {
		return nil, exception.NewDomainErrorWithCause(errCodeDBReadFailed, "failed to get unprocessed webhook events", err)
	}

	events := make([]*domain.WebhookEvent, len(models))
	for i, model := range models {
		event, err := model.ToEntity()
		if err != nil {
			return nil, err
		}
		events[i] = event
	}

	return events, nil
}
