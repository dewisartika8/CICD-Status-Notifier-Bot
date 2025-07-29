package domain

import (
	"time"

	"github.com/dewisartika8/cicd-status-notifier-bot/internal/core/shared/domain/value_objects"
)

// ProjectModel represents the database model for project
type ProjectModel struct {
	ID             string    `gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	Name           string    `gorm:"uniqueIndex;not null;size:100"`
	RepositoryURL  string    `gorm:"uniqueIndex;not null;size:500"`
	WebhookSecret  string    `gorm:"not null;size:255"`
	TelegramChatID *int64    `gorm:"index"`
	Status         string    `gorm:"not null;default:'active';size:20"`
	CreatedAt      time.Time `gorm:"not null"`
	UpdatedAt      time.Time `gorm:"not null"`
}

// TableName specifies the table name for GORM
func (ProjectModel) TableName() string {
	return "projects"
}

// ToEntity converts database model to domain entity
func (m *ProjectModel) ToEntity() (*Project, error) {
	id, err := value_objects.NewIDFromString(m.ID)
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
	m.ID = project.ID().String()
	m.Name = project.Name()
	m.RepositoryURL = project.RepositoryURL()
	m.WebhookSecret = project.WebhookSecret()
	m.TelegramChatID = project.TelegramChatID()
	m.Status = string(project.Status())
	m.CreatedAt = project.CreatedAt().ToTime()
	m.UpdatedAt = project.UpdatedAt().ToTime()
}
