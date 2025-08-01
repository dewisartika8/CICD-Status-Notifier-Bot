package domain

import (
	"time"

	"github.com/dewisartika8/cicd-status-notifier-bot/internal/core/shared/domain/value_objects"
	"github.com/google/uuid"
)

// NotificationTemplateModel represents the database model for notification templates
type NotificationTemplateModel struct {
	ID           uuid.UUID `gorm:"column:id;primaryKey;type:uuid;default:uuid_generate_v4()"`
	TemplateType string    `gorm:"column:template_type;not null;index:idx_notification_templates_type;uniqueIndex:unique_template_type_channel" json:"template_type"`
	Channel      string    `gorm:"column:channel;not null;index:idx_notification_templates_channel;uniqueIndex:unique_template_type_channel" json:"channel"`
	Subject      string    `gorm:"column:subject;not null" json:"subject"`
	BodyTemplate string    `gorm:"column:body_template;not null" json:"body_template"`
	IsActive     bool      `gorm:"column:is_active;default:true;index:idx_notification_templates_active" json:"is_active"`
	CreatedAt    time.Time `gorm:"column:created_at;type:timestamp with time zone;default:current_timestamp" json:"created_at"`
	UpdatedAt    time.Time `gorm:"column:updated_at;type:timestamp with time zone;default:current_timestamp" json:"updated_at"`
}

// TableName returns the table name for notification templates
func (NotificationTemplateModel) TableName() string {
	return "notification_templates"
}

// ToEntity converts the model to domain entity
func (m *NotificationTemplateModel) ToEntity() *NotificationTemplate {
	id, _ := value_objects.NewIDFromString(m.ID.String())
	templateType := NotificationTemplateType(m.TemplateType)
	channel := NotificationChannel(m.Channel)
	createdAt := value_objects.NewTimestampFromTime(m.CreatedAt)
	updatedAt := value_objects.NewTimestampFromTime(m.UpdatedAt)

	template, _ := RestoreNotificationTemplate(RestoreNotificationTemplateParams{
		ID:           id,
		TemplateType: templateType,
		Channel:      channel,
		Subject:      m.Subject,
		BodyTemplate: m.BodyTemplate,
		IsActive:     m.IsActive,
		CreatedAt:    createdAt,
		UpdatedAt:    updatedAt,
	})

	return template
}

// FromEntity converts domain entity to model
func (m *NotificationTemplateModel) FromEntity(template *NotificationTemplate) {
	// Parse UUID from string
	if id, err := uuid.Parse(template.ID().String()); err == nil {
		m.ID = id
	}

	m.TemplateType = string(template.TemplateType())
	m.Channel = string(template.Channel())
	m.Subject = template.Subject()
	m.BodyTemplate = template.BodyTemplate()
	m.IsActive = template.IsActive()
	m.CreatedAt = template.CreatedAt().ToTime()
	m.UpdatedAt = template.UpdatedAt().ToTime()
}
