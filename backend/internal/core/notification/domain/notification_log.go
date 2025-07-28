package domain

import (
	"time"

	"github.com/dewisartika8/cicd-status-notifier-bot/internal/domain/entities"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// NotificationLogModel represents the GORM model for notification logs
type NotificationLogModel struct {
	ID           uuid.UUID  `gorm:"type:uuid;primary_key;default:uuid_generate_v4()"`
	BuildEventID uuid.UUID  `gorm:"type:uuid;not null;index"`
	ChatID       int64      `gorm:"type:bigint;not null;index"`
	MessageID    *int       `gorm:"type:integer"`
	Status       string     `gorm:"type:varchar(20);not null;index"`
	ErrorMessage string     `gorm:"type:text"`
	SentAt       *time.Time `gorm:"type:timestamp with time zone"`
	CreatedAt    time.Time  `gorm:"type:timestamp with time zone;default:now();index:idx_notification_logs_created_at,sort:desc"`

	// Relationships
	BuildEvent BuildEventModel `gorm:"foreignKey:BuildEventID"`
}

// TableName returns the table name for the NotificationLogModel
func (NotificationLogModel) TableName() string {
	return "notification_logs"
}

// BeforeCreate hook to generate UUID if not set
func (nl *NotificationLogModel) BeforeCreate(tx *gorm.DB) error {
	if nl.ID == uuid.Nil {
		nl.ID = uuid.New()
	}
	return nil
}

// ToEntity converts GORM model to domain entity
func (nl *NotificationLogModel) ToEntity() *entities.NotificationLog {
	return &entities.NotificationLog{
		ID:           nl.ID,
		BuildEventID: nl.BuildEventID,
		ChatID:       nl.ChatID,
		MessageID:    nl.MessageID,
		Status:       entities.NotificationStatus(nl.Status),
		ErrorMessage: nl.ErrorMessage,
		SentAt:       nl.SentAt,
		CreatedAt:    nl.CreatedAt,
	}
}

// FromEntity converts domain entity to GORM model
func (nl *NotificationLogModel) FromEntity(entity *entities.NotificationLog) {
	nl.ID = entity.ID
	nl.BuildEventID = entity.BuildEventID
	nl.ChatID = entity.ChatID
	nl.MessageID = entity.MessageID
	nl.Status = string(entity.Status)
	nl.ErrorMessage = entity.ErrorMessage
	nl.SentAt = entity.SentAt
	nl.CreatedAt = entity.CreatedAt
}
