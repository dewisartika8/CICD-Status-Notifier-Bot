package domain

import (
	"strconv"
	"time"

	"github.com/dewisartika8/cicd-status-notifier-bot/internal/core/shared/domain/value_objects"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// NotificationLogModel represents the GORM model for notification logs
type NotificationLogModel struct {
	ID           uuid.UUID  `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()"`
	BuildEventID uuid.UUID  `gorm:"type:uuid;not null;index:idx_notification_logs_build_event_id"`
	ChatID       int64      `gorm:"type:bigint;not null;index:idx_notification_logs_chat_id"`
	MessageID    *int       `gorm:"type:integer"`
	Status       string     `gorm:"type:varchar(20);not null;index:idx_notification_logs_status"`
	ErrorMessage string     `gorm:"type:text"`
	SentAt       *time.Time `gorm:"type:timestamp with time zone"`
	CreatedAt    time.Time  `gorm:"type:timestamp with time zone;not null;default:now();index:idx_notification_logs_created_at"`

	// Additional columns from migration 002
	RetryCount int        `gorm:"type:integer;not null;default:0;column:retry_count;index:idx_notification_logs_retry_count"`
	Channel    string     `gorm:"type:varchar(50);column:channel;index:idx_notification_logs_channel"`
	TemplateID *uuid.UUID `gorm:"type:uuid;column:template_id;index:idx_notification_logs_template"`
}

// TableName returns the table name for the NotificationLogModel
func (NotificationLogModel) TableName() string {
	return "notification_logs"
}

// BeforeCreate hook to set timestamps
func (nlm *NotificationLogModel) BeforeCreate(tx *gorm.DB) error {
	now := time.Now()
	if nlm.CreatedAt.IsZero() {
		nlm.CreatedAt = now
	}
	return nil
}

// BeforeUpdate hook to update timestamp
func (nlm *NotificationLogModel) BeforeUpdate(tx *gorm.DB) error {
	// No UpdatedAt field in current schema
	return nil
}

// ToEntity converts GORM model to domain entity
func (nlm *NotificationLogModel) ToEntity() *NotificationLog {
	id, _ := value_objects.NewIDFromString(nlm.ID.String())
	buildEventID, _ := value_objects.NewIDFromString(nlm.BuildEventID.String())
	// Note: ProjectID is not in the current database schema but exists in domain
	// For now, we'll generate a default project ID
	projectID := value_objects.NewID()

	params := RestoreNotificationLogParams{
		ID:           id,
		BuildEventID: buildEventID,
		ProjectID:    projectID,
		Channel:      NotificationChannel(nlm.Channel),
		Recipient:    strconv.FormatInt(nlm.ChatID, 10), // Convert ChatID to string as recipient
		Message:      "Notification",                    // Default message since not stored in DB
		Status:       NotificationStatus(nlm.Status),
		ErrorMessage: nlm.ErrorMessage,
		RetryCount:   nlm.RetryCount,
		MessageID:    convertIntToStringPointer(nlm.MessageID),
		CreatedAt:    value_objects.NewTimestampFromTime(nlm.CreatedAt),
		UpdatedAt:    value_objects.NewTimestampFromTime(nlm.CreatedAt), // Use CreatedAt since no UpdatedAt in DB
	}

	if nlm.SentAt != nil {
		sentAt := value_objects.NewTimestampFromTime(*nlm.SentAt)
		params.SentAt = &sentAt
	}

	if nlm.TemplateID != nil {
		templateID, _ := value_objects.NewIDFromString(nlm.TemplateID.String())
		params.TemplateID = &templateID
	}

	return RestoreNotificationLog(params)
}

// convertIntToStringPointer converts *int to *string
func convertIntToStringPointer(i *int) *string {
	if i == nil {
		return nil
	}
	str := strconv.Itoa(*i)
	return &str
}

// FromEntity converts domain entity to GORM model
func (nlm *NotificationLogModel) FromEntity(entity *NotificationLog) {
	// Parse UUID from string
	if id, err := uuid.Parse(entity.ID().String()); err == nil {
		nlm.ID = id
	}
	if buildEventID, err := uuid.Parse(entity.BuildEventID().String()); err == nil {
		nlm.BuildEventID = buildEventID
	}

	// Convert recipient string back to ChatID (assuming it's a number)
	if chatID, err := strconv.ParseInt(entity.Recipient(), 10, 64); err == nil {
		nlm.ChatID = chatID
	}

	nlm.Channel = string(entity.Channel())
	nlm.Status = string(entity.Status())
	nlm.ErrorMessage = entity.ErrorMessage()
	nlm.RetryCount = entity.RetryCount()
	nlm.MessageID = convertStringToIntPointer(entity.MessageID())
	nlm.CreatedAt = entity.CreatedAt().ToTime()

	if entity.SentAt() != nil {
		sentAtTime := entity.SentAt().ToTime()
		nlm.SentAt = &sentAtTime
	}

	if entity.TemplateID() != nil {
		templateID, _ := uuid.Parse(entity.TemplateID().String())
		nlm.TemplateID = &templateID
	}
}

// convertStringToIntPointer converts *string to *int
func convertStringToIntPointer(s *string) *int {
	if s == nil {
		return nil
	}
	if i, err := strconv.Atoi(*s); err == nil {
		return &i
	}
	return nil
}
