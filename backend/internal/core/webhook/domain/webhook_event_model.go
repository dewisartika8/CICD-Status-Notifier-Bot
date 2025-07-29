package domain

import (
	"time"

	"github.com/dewisartika8/cicd-status-notifier-bot/internal/core/shared/domain/value_objects"
	"github.com/google/uuid"
)

// WebhookEventModel represents the database model for webhook events
type WebhookEventModel struct {
	ID          uuid.UUID  `gorm:"primaryKey;type:uuid;default:uuid_generate_v4()" json:"id"`
	ProjectID   uuid.UUID  `gorm:"type:uuid;not null;index" json:"project_id"`
	EventType   string     `gorm:"type:varchar(50);not null;index" json:"event_type"`
	Payload     string     `gorm:"type:jsonb;not null" json:"payload"`
	Signature   string     `gorm:"type:varchar(255);not null" json:"signature"`
	DeliveryID  string     `gorm:"type:varchar(255);index" json:"delivery_id"`
	ProcessedAt *time.Time `gorm:"type:timestamp" json:"processed_at"`
	CreatedAt   time.Time  `gorm:"type:timestamp;default:CURRENT_TIMESTAMP" json:"created_at"`
}

// TableName specifies the table name for GORM
func (WebhookEventModel) TableName() string {
	return "webhook_events"
}

// ToEntity converts database model to domain entity
func (m *WebhookEventModel) ToEntity() (*WebhookEvent, error) {
	id, err := value_objects.NewIDFromString(m.ID.String())
	if err != nil {
		return nil, err
	}

	projectID, err := value_objects.NewIDFromString(m.ProjectID.String())
	if err != nil {
		return nil, err
	}

	createdAt := value_objects.NewTimestampFromTime(m.CreatedAt)

	// Create entity from database data
	data := WebhookEventData{
		ID:          id,
		ProjectID:   projectID,
		EventType:   WebhookEventType(m.EventType),
		Payload:     m.Payload,
		Signature:   m.Signature,
		DeliveryID:  m.DeliveryID,
		ProcessedAt: m.ProcessedAt,
		CreatedAt:   createdAt,
	}

	return NewWebhookEventFromData(data), nil
}

// FromEntity converts domain entity to database model
func (m *WebhookEventModel) FromEntity(event *WebhookEvent) {
	// Parse UUID from string
	if id, err := uuid.Parse(event.ID().String()); err == nil {
		m.ID = id
	}
	if projectID, err := uuid.Parse(event.ProjectID().String()); err == nil {
		m.ProjectID = projectID
	}

	m.EventType = string(event.EventType())
	m.Payload = event.Payload()
	m.Signature = event.Signature()
	m.DeliveryID = event.DeliveryID()
	m.ProcessedAt = event.ProcessedAt()
	m.CreatedAt = event.CreatedAt().ToTime()
}
