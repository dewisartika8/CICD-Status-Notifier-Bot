package domain

import (
	"time"

	"github.com/dewisartika8/cicd-status-notifier-bot/internal/core/shared/domain/value_objects"
	"gorm.io/gorm"
)

// NotificationLogModel represents the GORM model for notification logs
type NotificationLogModel struct {
	ID           string     `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	BuildEventID string     `gorm:"type:uuid;not null;index:idx_notification_logs_build_event"`
	ProjectID    string     `gorm:"type:uuid;not null;index:idx_notification_logs_project"`
	Channel      string     `gorm:"type:varchar(20);not null;index:idx_notification_logs_channel"`
	Recipient    string     `gorm:"type:varchar(255);not null"`
	Message      string     `gorm:"type:text;not null"`
	Status       string     `gorm:"type:varchar(20);not null;index:idx_notification_logs_status"`
	ErrorMessage string     `gorm:"type:text"`
	RetryCount   int        `gorm:"type:int;not null;default:0"`
	MessageID    *string    `gorm:"type:varchar(255)"` // External message ID (e.g., Telegram message ID)
	SentAt       *time.Time `gorm:"type:timestamptz"`
	CreatedAt    time.Time  `gorm:"type:timestamptz;not null;default:now();index:idx_notification_logs_created_at"`
	UpdatedAt    time.Time  `gorm:"type:timestamptz;not null;default:now()"`
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
	if nlm.UpdatedAt.IsZero() {
		nlm.UpdatedAt = now
	}
	return nil
}

// BeforeUpdate hook to update timestamp
func (nlm *NotificationLogModel) BeforeUpdate(tx *gorm.DB) error {
	nlm.UpdatedAt = time.Now()
	return nil
}

// ToEntity converts GORM model to domain entity
func (nlm *NotificationLogModel) ToEntity() *NotificationLog {
	id, _ := value_objects.NewIDFromString(nlm.ID)
	buildEventID, _ := value_objects.NewIDFromString(nlm.BuildEventID)
	projectID, _ := value_objects.NewIDFromString(nlm.ProjectID)

	params := RestoreNotificationLogParams{
		ID:           id,
		BuildEventID: buildEventID,
		ProjectID:    projectID,
		Channel:      NotificationChannel(nlm.Channel),
		Recipient:    nlm.Recipient,
		Message:      nlm.Message,
		Status:       NotificationStatus(nlm.Status),
		ErrorMessage: nlm.ErrorMessage,
		RetryCount:   nlm.RetryCount,
		MessageID:    nlm.MessageID,
		CreatedAt:    value_objects.NewTimestampFromTime(nlm.CreatedAt),
		UpdatedAt:    value_objects.NewTimestampFromTime(nlm.UpdatedAt),
	}

	if nlm.SentAt != nil {
		sentAt := value_objects.NewTimestampFromTime(*nlm.SentAt)
		params.SentAt = &sentAt
	}

	return RestoreNotificationLog(params)
}

// FromEntity converts domain entity to GORM model
func (nlm *NotificationLogModel) FromEntity(entity *NotificationLog) {
	nlm.ID = entity.ID().String()
	nlm.BuildEventID = entity.BuildEventID().String()
	nlm.ProjectID = entity.ProjectID().String()
	nlm.Channel = string(entity.Channel())
	nlm.Recipient = entity.Recipient()
	nlm.Message = entity.Message()
	nlm.Status = string(entity.Status())
	nlm.ErrorMessage = entity.ErrorMessage()
	nlm.RetryCount = entity.RetryCount()
	nlm.MessageID = entity.MessageID()
	nlm.CreatedAt = entity.CreatedAt().ToTime()
	nlm.UpdatedAt = entity.UpdatedAt().ToTime()

	if entity.SentAt() != nil {
		sentAtTime := entity.SentAt().ToTime()
		nlm.SentAt = &sentAtTime
	}
}

// TelegramSubscriptionModel represents the GORM model for telegram subscriptions
type TelegramSubscriptionModel struct {
	ID        string    `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	ProjectID string    `gorm:"type:uuid;not null;index:idx_telegram_subscriptions_project"`
	ChatID    int64     `gorm:"type:bigint;not null;index:idx_telegram_subscriptions_chat"`
	IsActive  bool      `gorm:"type:boolean;not null;default:true;index:idx_telegram_subscriptions_active"`
	CreatedAt time.Time `gorm:"type:timestamptz;not null;default:now()"`
	UpdatedAt time.Time `gorm:"type:timestamptz;not null;default:now()"`
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
	id, _ := value_objects.NewIDFromString(tsm.ID)
	projectID, _ := value_objects.NewIDFromString(tsm.ProjectID)

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
	tsm.ID = entity.ID().String()
	tsm.ProjectID = entity.ProjectID().String()
	tsm.ChatID = entity.ChatID()
	tsm.IsActive = entity.IsActive()
	tsm.CreatedAt = entity.CreatedAt().ToTime()
	tsm.UpdatedAt = entity.UpdatedAt().ToTime()
}
