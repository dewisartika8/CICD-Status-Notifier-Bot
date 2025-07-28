package dto

import (
	"time"

	"github.com/google/uuid"
)

// Project represents a CI/CD project entity
type Project struct {
	ID             uuid.UUID `json:"id"`
	Name           string    `json:"name" validate:"required,min=1,max=255"`
	RepositoryURL  string    `json:"repository_url" validate:"required,url"`
	WebhookSecret  string    `json:"-"` // Hidden from JSON for security
	TelegramChatID *int64    `json:"telegram_chat_id,omitempty"`
	IsActive       bool      `json:"is_active"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

// ProjectMetrics represents aggregated metrics for a project
type ProjectMetrics struct {
	ProjectID      uuid.UUID                      `json:"project_id"`
	TotalBuilds    int64                          `json:"total_builds"`
	SuccessRate    float64                        `json:"success_rate"`
	FailureRate    float64                        `json:"failure_rate"`
	AvgDuration    float64                        `json:"avg_duration_seconds"`
	BuildsByStatus map[entities.BuildStatus]int64 `json:"builds_by_status"`
	BuildsByType   map[entities.EventType]int64   `json:"builds_by_type"`
}

// NewProject creates a new project entity
func NewProject(name, repositoryURL, webhookSecret string) *Project {
	return &Project{
		ID:            uuid.New(),
		Name:          name,
		RepositoryURL: repositoryURL,
		WebhookSecret: webhookSecret,
		IsActive:      true,
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}
}

// Validate performs basic validation on project entity
func (p *Project) Validate() error {
	if p.Name == "" {
		return ErrInvalidProjectName
	}
	if p.RepositoryURL == "" {
		return ErrInvalidRepositoryURL
	}
	return nil
}

// Update updates project fields while maintaining entity invariants
func (p *Project) Update(name, repositoryURL string, telegramChatID *int64, isActive bool) error {
	if name != "" {
		p.Name = name
	}
	if repositoryURL != "" {
		p.RepositoryURL = repositoryURL
	}
	p.TelegramChatID = telegramChatID
	p.IsActive = isActive
	p.UpdatedAt = time.Now()

	return p.Validate()
}
