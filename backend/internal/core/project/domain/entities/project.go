package entities

import (
	"strings"

	"github.com/dewisartika8/cicd-status-notifier-bot/internal/core/shared/domain/errors"
	"github.com/dewisartika8/cicd-status-notifier-bot/internal/core/shared/domain/value_objects"
)

// ProjectStatus represents the status of a project
type ProjectStatus value_objects.Status

const (
	ProjectStatusActive   ProjectStatus = "active"
	ProjectStatusInactive ProjectStatus = "inactive"
	ProjectStatusArchived ProjectStatus = "archived"
)

// Project represents a CI/CD project domain entity
type Project struct {
	id             value_objects.ID
	name           string
	repositoryURL  string
	webhookSecret  string
	telegramChatID *int64
	status         ProjectStatus
	createdAt      value_objects.Timestamp
	updatedAt      value_objects.Timestamp
}

// NewProject creates a new project entity
func NewProject(name, repositoryURL, webhookSecret string, telegramChatID *int64) (*Project, error) {
	project := &Project{
		id:             value_objects.NewID(),
		name:           strings.TrimSpace(name),
		repositoryURL:  strings.TrimSpace(repositoryURL),
		webhookSecret:  webhookSecret,
		telegramChatID: telegramChatID,
		status:         ProjectStatusActive,
		createdAt:      value_objects.NewTimestamp(),
		updatedAt:      value_objects.NewTimestamp(),
	}

	if err := project.validate(); err != nil {
		return nil, err
	}

	return project, nil
}

// RestoreProject restores a project from persistence
func RestoreProject(
	id value_objects.ID,
	name, repositoryURL, webhookSecret string,
	telegramChatID *int64,
	status ProjectStatus,
	createdAt, updatedAt value_objects.Timestamp,
) *Project {
	return &Project{
		id:             id,
		name:           name,
		repositoryURL:  repositoryURL,
		webhookSecret:  webhookSecret,
		telegramChatID: telegramChatID,
		status:         status,
		createdAt:      createdAt,
		updatedAt:      updatedAt,
	}
}

// ID returns the project ID
func (p *Project) ID() value_objects.ID {
	return p.id
}

// Name returns the project name
func (p *Project) Name() string {
	return p.name
}

// RepositoryURL returns the repository URL
func (p *Project) RepositoryURL() string {
	return p.repositoryURL
}

// WebhookSecret returns the webhook secret
func (p *Project) WebhookSecret() string {
	return p.webhookSecret
}

// TelegramChatID returns the telegram chat ID
func (p *Project) TelegramChatID() *int64 {
	return p.telegramChatID
}

// Status returns the project status
func (p *Project) Status() ProjectStatus {
	return p.status
}

// CreatedAt returns when the project was created
func (p *Project) CreatedAt() value_objects.Timestamp {
	return p.createdAt
}

// UpdatedAt returns when the project was last updated
func (p *Project) UpdatedAt() value_objects.Timestamp {
	return p.updatedAt
}

// UpdateName updates the project name
func (p *Project) UpdateName(name string) error {
	name = strings.TrimSpace(name)
	if name == "" {
		return errors.NewDomainError("INVALID_PROJECT_NAME", "project name cannot be empty")
	}
	p.name = name
	p.updatedAt = value_objects.NewTimestamp()
	return nil
}

// UpdateRepositoryURL updates the repository URL
func (p *Project) UpdateRepositoryURL(url string) error {
	url = strings.TrimSpace(url)
	if url == "" {
		return errors.NewDomainError("INVALID_REPOSITORY_URL", "repository URL cannot be empty")
	}
	p.repositoryURL = url
	p.updatedAt = value_objects.NewTimestamp()
	return nil
}

// UpdateWebhookSecret updates the webhook secret
func (p *Project) UpdateWebhookSecret(secret string) {
	p.webhookSecret = secret
	p.updatedAt = value_objects.NewTimestamp()
}

// UpdateTelegramChatID updates the telegram chat ID
func (p *Project) UpdateTelegramChatID(chatID *int64) {
	p.telegramChatID = chatID
	p.updatedAt = value_objects.NewTimestamp()
}

// Activate activates the project
func (p *Project) Activate() {
	p.status = ProjectStatusActive
	p.updatedAt = value_objects.NewTimestamp()
}

// Deactivate deactivates the project
func (p *Project) Deactivate() {
	p.status = ProjectStatusInactive
	p.updatedAt = value_objects.NewTimestamp()
}

// Archive archives the project
func (p *Project) Archive() {
	p.status = ProjectStatusArchived
	p.updatedAt = value_objects.NewTimestamp()
}

// IsActive checks if the project is active
func (p *Project) IsActive() bool {
	return p.status == ProjectStatusActive
}

// validate performs domain validation
func (p *Project) validate() error {
	if p.name == "" {
		return errors.NewDomainError("INVALID_PROJECT_NAME", "project name cannot be empty")
	}

	if p.repositoryURL == "" {
		return errors.NewDomainError("INVALID_REPOSITORY_URL", "repository URL cannot be empty")
	}

	if !isValidRepositoryURL(p.repositoryURL) {
		return errors.NewDomainError("INVALID_REPOSITORY_URL", "repository URL format is invalid")
	}

	return nil
}

// isValidRepositoryURL performs basic URL validation for repository
func isValidRepositoryURL(url string) bool {
	// Basic validation - should contain github.com, gitlab.com, or bitbucket.org
	url = strings.ToLower(url)
	return strings.Contains(url, "github.com") ||
		strings.Contains(url, "gitlab.com") ||
		strings.Contains(url, "bitbucket.org") ||
		strings.HasPrefix(url, "https://") ||
		strings.HasPrefix(url, "http://")
}
