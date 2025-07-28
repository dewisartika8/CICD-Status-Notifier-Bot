package database

import (
	"database/sql/driver"
	"time"

	"github.com/dewisartika8/cicd-status-notifier-bot/internal/domain/entities"
	"github.com/google/uuid"
	"github.com/lib/pq"
	"gorm.io/gorm"
)

// StringArray represents a PostgreSQL text array
type StringArray []string

// Value implements driver.Valuer interface for StringArray
func (sa StringArray) Value() (driver.Value, error) {
	if sa == nil {
		return nil, nil
	}
	return pq.Array(sa).Value()
}

// Scan implements sql.Scanner interface for StringArray
func (sa *StringArray) Scan(value interface{}) error {
	var arr pq.StringArray
	if err := arr.Scan(value); err != nil {
		return err
	}
	*sa = StringArray(arr)
	return nil
}

// TelegramSubscriptionModel represents the GORM model for telegram subscriptions
type TelegramSubscriptionModel struct {
	ID         uuid.UUID   `gorm:"type:uuid;primary_key;default:uuid_generate_v4()"`
	ProjectID  uuid.UUID   `gorm:"type:uuid;not null;index"`
	ChatID     int64       `gorm:"type:bigint;not null;index"`
	UserID     *int64      `gorm:"type:bigint"`
	Username   string      `gorm:"type:varchar(255)"`
	EventTypes StringArray `gorm:"type:text[]"`
	IsActive   bool        `gorm:"type:boolean;default:true;index"`
	CreatedAt  time.Time   `gorm:"type:timestamp with time zone;default:now()"`

	// Relationships
	Project ProjectModel `gorm:"foreignKey:ProjectID"`
}

// TableName returns the table name for the TelegramSubscriptionModel
func (TelegramSubscriptionModel) TableName() string {
	return "telegram_subscriptions"
}

// BeforeCreate hook to generate UUID and set default event types
func (ts *TelegramSubscriptionModel) BeforeCreate(tx *gorm.DB) error {
	if ts.ID == uuid.Nil {
		ts.ID = uuid.New()
	}

	// Set default event types if empty
	if len(ts.EventTypes) == 0 {
		ts.EventTypes = StringArray{
			string(entities.EventTypeBuildSuccess),
			string(entities.EventTypeBuildFailed),
			string(entities.EventTypeDeploymentSuccess),
			string(entities.EventTypeDeploymentFailed),
		}
	}

	return nil
}

// ToEntity converts GORM model to domain entity
func (ts *TelegramSubscriptionModel) ToEntity() *entities.TelegramSubscription {
	eventTypes := make([]entities.EventType, len(ts.EventTypes))
	for i, et := range ts.EventTypes {
		eventTypes[i] = entities.EventType(et)
	}

	return &entities.TelegramSubscription{
		ID:         ts.ID,
		ProjectID:  ts.ProjectID,
		ChatID:     ts.ChatID,
		UserID:     ts.UserID,
		Username:   ts.Username,
		EventTypes: eventTypes,
		IsActive:   ts.IsActive,
		CreatedAt:  ts.CreatedAt,
	}
}

// FromEntity converts domain entity to GORM model
func (ts *TelegramSubscriptionModel) FromEntity(entity *entities.TelegramSubscription) {
	ts.ID = entity.ID
	ts.ProjectID = entity.ProjectID
	ts.ChatID = entity.ChatID
	ts.UserID = entity.UserID
	ts.Username = entity.Username
	ts.IsActive = entity.IsActive
	ts.CreatedAt = entity.CreatedAt

	// Convert event types
	eventTypes := make(StringArray, len(entity.EventTypes))
	for i, et := range entity.EventTypes {
		eventTypes[i] = string(et)
	}
	ts.EventTypes = eventTypes
}
