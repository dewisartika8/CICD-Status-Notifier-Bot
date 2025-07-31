package domain

import (
	"time"

	"github.com/dewisartika8/cicd-status-notifier-bot/internal/core/shared/domain/value_objects"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// TelegramSubscriptionModel represents the GORM model for telegram subscriptions
type TelegramSubscriptionModel struct {
	ID        uuid.UUID `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()"`
	ProjectID uuid.UUID `gorm:"type:uuid;not null;index:idx_telegram_subscriptions_project"`
	ChatID    int64     `gorm:"type:bigint;not null;index:idx_telegram_subscriptions_chat"`
	IsActive  bool      `gorm:"type:boolean;not null;default:true;index:idx_telegram_subscriptions_active"`
	CreatedAt time.Time `gorm:"type:timestamp with time zone;not null;default:now()"`
	UpdatedAt time.Time `gorm:"type:timestamp with time zone;not null;default:now()"`
}

// TableName returns the table name for the TelegramSubscriptionModel
func (TelegramSubscriptionModel) TableName() string {
	return "telegram_subscriptions"
}

// BeforeCreate hook to set timestamps
func (tsm *TelegramSubscriptionModel) BeforeCreate(tx *gorm.DB) error {
	now := time.Now()
	if tsm.CreatedAt.IsZero() {
		tsm.CreatedAt = now
	}
	if tsm.UpdatedAt.IsZero() {
		tsm.UpdatedAt = now
	}
	return nil
}

// BeforeUpdate hook to update timestamp
func (tsm *TelegramSubscriptionModel) BeforeUpdate(tx *gorm.DB) error {
	tsm.UpdatedAt = time.Now()
	return nil
}

// ToEntity converts GORM model to domain entity
func (tsm *TelegramSubscriptionModel) ToEntity() *TelegramSubscription {
	id, _ := value_objects.NewIDFromString(tsm.ID.String())
	projectID, _ := value_objects.NewIDFromString(tsm.ProjectID.String())

	params := RestoreTelegramSubscriptionParams{
		ID:        id,
		ProjectID: projectID,
		ChatID:    tsm.ChatID,
		IsActive:  tsm.IsActive,
		CreatedAt: value_objects.NewTimestampFromTime(tsm.CreatedAt),
		UpdatedAt: value_objects.NewTimestampFromTime(tsm.UpdatedAt),
	}

	return RestoreTelegramSubscription(params)
}

// FromEntity converts domain entity to GORM model
func (tsm *TelegramSubscriptionModel) FromEntity(entity *TelegramSubscription) {
	// Parse UUID from string
	if id, err := uuid.Parse(entity.ID().String()); err == nil {
		tsm.ID = id
	}
	if projectID, err := uuid.Parse(entity.ProjectID().String()); err == nil {
		tsm.ProjectID = projectID
	}

	tsm.ChatID = entity.ChatID()
	tsm.IsActive = entity.IsActive()
	tsm.CreatedAt = entity.CreatedAt().ToTime()
	tsm.UpdatedAt = entity.UpdatedAt().ToTime()
}
