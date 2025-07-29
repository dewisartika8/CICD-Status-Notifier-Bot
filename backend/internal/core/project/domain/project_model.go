package domain

import (
	"time"

	"github.com/dewisartika8/cicd-status-notifier-bot/internal/core/shared/domain/value_objects"
	"github.com/google/uuid"
)

// ProjectModel represents the database model for project
type ProjectModel struct {
	ID             uuid.UUID `gorm:"primaryKey;type:uuid;default:uuid_generate_v4()"`
	Name           string    `gorm:"uniqueIndex;not null;size:255"`
	RepositoryURL  string    `gorm:"not null;size:500"`
	WebhookSecret  string    `gorm:"size:255"`
	TelegramChatID *int64    `gorm:"index"`
	Status         string    `gorm:"not null;default:'active';size:20"`
	CreatedAt      time.Time `gorm:"type:timestamp with time zone;default:now()"`
	UpdatedAt      time.Time `gorm:"type:timestamp with time zone;default:now()"`
}

// TableName specifies the table name for GORM
func (ProjectModel) TableName() string {
	return "projects"
}

// ToEntity converts database model to domain entity
func (m *ProjectModel) ToEntity() (*Project, error) {
	id, err := value_objects.NewIDFromString(m.ID.String())
	if err != nil {
		return nil, err
	}

	createdAt := value_objects.NewTimestampFromTime(m.CreatedAt)
	updatedAt := value_objects.NewTimestampFromTime(m.UpdatedAt)

	// Create entity from database data
	dbData := ProjectDBData{
		ID:             id,
		Name:           m.Name,
		RepositoryURL:  m.RepositoryURL,
		WebhookSecret:  m.WebhookSecret,
		TelegramChatID: m.TelegramChatID,
		Status:         ProjectStatus(m.Status),
		CreatedAt:      createdAt,
		UpdatedAt:      updatedAt,
	}

	return NewProjectFromDB(dbData), nil
}

// FromEntity converts domain entity to database model
func (m *ProjectModel) FromEntity(project *Project) {
	// Parse UUID from string
	if id, err := uuid.Parse(project.ID().String()); err == nil {
		m.ID = id
	}

	m.Name = project.Name()
	m.RepositoryURL = project.RepositoryURL()
	m.WebhookSecret = project.WebhookSecret()
	m.TelegramChatID = project.TelegramChatID()
	m.Status = string(project.Status())
	m.CreatedAt = project.CreatedAt().ToTime()
	m.UpdatedAt = project.UpdatedAt().ToTime()
}
