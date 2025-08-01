package domain

import (
	"encoding/json"
	"time"

	"github.com/dewisartika8/cicd-status-notifier-bot/internal/core/shared/domain/value_objects"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// BuildEventModel represents the GORM model for build events
type BuildEventModel struct {
	ID              uuid.UUID       `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()" json:"id"`
	ProjectID       uuid.UUID       `gorm:"type:uuid;not null;index:idx_build_events_project_id" json:"project_id"`
	EventType       string          `gorm:"type:varchar(50);not null;index:idx_build_events_event_type" json:"event_type"`
	Status          string          `gorm:"type:varchar(20);not null;index:idx_build_events_status" json:"status"`
	Branch          string          `gorm:"type:varchar(255);not null;index:idx_build_events_branch" json:"branch"`
	CommitSHA       string          `gorm:"type:varchar(40)" json:"commit_sha"`
	CommitMessage   string          `gorm:"type:text" json:"commit_message"`
	AuthorName      string          `gorm:"type:varchar(255)" json:"author_name"`
	AuthorEmail     string          `gorm:"type:varchar(255)" json:"author_email"`
	BuildURL        string          `gorm:"type:varchar(500)" json:"build_url"`
	DurationSeconds *int            `gorm:"type:integer" json:"duration_seconds"`
	WebhookPayload  json.RawMessage `gorm:"type:jsonb" json:"webhook_payload"`
	CreatedAt       time.Time       `gorm:"type:timestamp with time zone;default:now();index:idx_build_events_created_at" json:"created_at"`
	UpdatedAt       time.Time       `gorm:"type:timestamp with time zone;default:now()" json:"updated_at"`
	DeletedAt       gorm.DeletedAt  `gorm:"index:idx_build_events_deleted_at" json:"deleted_at,omitempty"`

	// Relationships
	Project ProjectModel `gorm:"foreignKey:ProjectID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
}

// TableName returns the table name for the BuildEventModel
func (BuildEventModel) TableName() string {
	return "build_events"
}

// ToEntity converts the GORM model to domain entity
func (m *BuildEventModel) ToEntity() *BuildEvent {
	id := value_objects.NewIDFromUUID(m.ID)
	projectID := value_objects.NewIDFromUUID(m.ProjectID)
	createdAt := value_objects.NewTimestampFromTime(m.CreatedAt)

	return RestoreBuildEvent(RestoreBuildEventParams{
		ID:              id,
		ProjectID:       projectID,
		EventType:       EventType(m.EventType),
		Status:          BuildStatus(m.Status),
		Branch:          m.Branch,
		CommitSHA:       m.CommitSHA,
		CommitMessage:   m.CommitMessage,
		AuthorName:      m.AuthorName,
		AuthorEmail:     m.AuthorEmail,
		BuildURL:        m.BuildURL,
		DurationSeconds: m.DurationSeconds,
		WebhookPayload:  m.WebhookPayload,
		CreatedAt:       createdAt,
	})
}

// FromEntity converts domain entity to GORM model
func (m *BuildEventModel) FromEntity(entity *BuildEvent) {
	m.ID = entity.ID().Value()
	m.ProjectID = entity.ProjectID().Value()
	m.EventType = string(entity.EventType())
	m.Status = string(entity.Status())
	m.Branch = entity.Branch()
	m.CommitSHA = entity.CommitSHA()
	m.CommitMessage = entity.CommitMessage()
	m.AuthorName = entity.AuthorName()
	m.AuthorEmail = entity.AuthorEmail()
	m.BuildURL = entity.BuildURL()
	m.DurationSeconds = entity.DurationSeconds()
	m.WebhookPayload = entity.WebhookPayload()
	m.CreatedAt = entity.CreatedAt().ToTime()
	m.UpdatedAt = time.Now()
}

// ProjectModel represents a simplified project model for relationships
type ProjectModel struct {
	ID   uuid.UUID `gorm:"type:uuid;primary_key;default:uuid_generate_v4()"`
	Name string    `gorm:"type:varchar(255);not null"`
}

// TableName returns the table name for the ProjectModel
func (ProjectModel) TableName() string {
	return "projects"
}
