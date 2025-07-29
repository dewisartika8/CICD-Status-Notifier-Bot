package domain

import (
	"time"

	"github.com/dewisartika8/cicd-status-notifier-bot/internal/core/shared/domain/value_objects"
)

// NotificationTemplateModel represents the database model for notification templates
type NotificationTemplateModel struct {
	ID           string    `gorm:"column:id;primaryKey" json:"id"`
	TemplateType string    `gorm:"column:template_type;not null" json:"template_type"`
	Channel      string    `gorm:"column:channel;not null" json:"channel"`
	Subject      string    `gorm:"column:subject;not null" json:"subject"`
	BodyTemplate string    `gorm:"column:body_template;not null" json:"body_template"`
	IsActive     bool      `gorm:"column:is_active;default:true" json:"is_active"`
	CreatedAt    time.Time `gorm:"column:created_at;autoCreateTime" json:"created_at"`
	UpdatedAt    time.Time `gorm:"column:updated_at;autoUpdateTime" json:"updated_at"`
}

// TableName returns the table name for notification templates
func (NotificationTemplateModel) TableName() string {
	return "notification_templates"
}

// ToEntity converts the model to domain entity
func (m *NotificationTemplateModel) ToEntity() *NotificationTemplate {
	id, _ := value_objects.NewIDFromString(m.ID)
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
	m.ID = template.ID().String()
	m.TemplateType = string(template.TemplateType())
	m.Channel = string(template.Channel())
	m.Subject = template.Subject()
	m.BodyTemplate = template.BodyTemplate()
	m.IsActive = template.IsActive()
	m.CreatedAt = template.CreatedAt().ToTime()
	m.UpdatedAt = template.UpdatedAt().ToTime()
}
