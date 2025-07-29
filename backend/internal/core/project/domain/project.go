package domain

import (
	"net/url"
	"strings"

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

// validate validates the project entity
func (p *Project) validate() error {
	if p.name == "" {
		return ErrInvalidProjectName
	}

	if p.repositoryURL == "" {
		return ErrInvalidRepositoryURL
	}

	// Validate repository URL format
	if _, err := url.Parse(p.repositoryURL); err != nil {
		return ErrInvalidRepositoryURL
	}

	// Validate webhook secret is not empty
	if strings.TrimSpace(p.webhookSecret) == "" {
		return ErrInvalidWebhookSecret
	}

	// Validate telegram chat ID if provided
	if p.telegramChatID != nil && *p.telegramChatID == 0 {
		return ErrInvalidTelegramChat
	}

	return nil
}

// Getters
func (p *Project) ID() value_objects.ID {
	return p.id
}

func (p *Project) Name() string {
	return p.name
}

func (p *Project) RepositoryURL() string {
	return p.repositoryURL
}

func (p *Project) WebhookSecret() string {
	return p.webhookSecret
}

func (p *Project) TelegramChatID() *int64 {
	return p.telegramChatID
}

func (p *Project) Status() ProjectStatus {
	return p.status
}

func (p *Project) CreatedAt() value_objects.Timestamp {
	return p.createdAt
}

func (p *Project) UpdatedAt() value_objects.Timestamp {
	return p.updatedAt
}

// Business logic methods

// UpdateName updates the project name
func (p *Project) UpdateName(name string) error {
	trimmedName := strings.TrimSpace(name)
	if trimmedName == "" {
		return ErrInvalidProjectName
	}

	p.name = trimmedName
	p.updatedAt = value_objects.NewTimestamp()
	return nil
}

// UpdateRepositoryURL updates the repository URL
func (p *Project) UpdateRepositoryURL(repositoryURL string) error {
	trimmedURL := strings.TrimSpace(repositoryURL)
	if trimmedURL == "" {
		return ErrInvalidRepositoryURL
	}

	// Validate URL format
	if _, err := url.Parse(trimmedURL); err != nil {
		return ErrInvalidRepositoryURL
	}

	p.repositoryURL = trimmedURL
	p.updatedAt = value_objects.NewTimestamp()
	return nil
}

// UpdateWebhookSecret updates the webhook secret
func (p *Project) UpdateWebhookSecret(webhookSecret string) error {
	if strings.TrimSpace(webhookSecret) == "" {
		return ErrInvalidWebhookSecret
	}

	p.webhookSecret = webhookSecret
	p.updatedAt = value_objects.NewTimestamp()
	return nil
}

// UpdateTelegramChatID updates the telegram chat ID
func (p *Project) UpdateTelegramChatID(chatID *int64) error {
	if chatID != nil && *chatID == 0 {
		return ErrInvalidTelegramChat
	}

	p.telegramChatID = chatID
	p.updatedAt = value_objects.NewTimestamp()
	return nil
}

// SetStatus sets the project status
func (p *Project) SetStatus(status ProjectStatus) {
	p.status = status
	p.updatedAt = value_objects.NewTimestamp()
}

// Activate sets the project status to active
func (p *Project) Activate() {
	p.SetStatus(ProjectStatusActive)
}

// Deactivate sets the project status to inactive
func (p *Project) Deactivate() {
	p.SetStatus(ProjectStatusInactive)
}

// Archive sets the project status to archived
func (p *Project) Archive() {
	p.SetStatus(ProjectStatusArchived)
}

// IsActive checks if the project is active
func (p *Project) IsActive() bool {
	return p.status == ProjectStatusActive
}

// IsInactive checks if the project is inactive
func (p *Project) IsInactive() bool {
	return p.status == ProjectStatusInactive
}

// IsArchived checks if the project is archived
func (p *Project) IsArchived() bool {
	return p.status == ProjectStatusArchived
}

// CanReceiveNotifications checks if the project can receive notifications
func (p *Project) CanReceiveNotifications() bool {
	return p.IsActive() && p.telegramChatID != nil
}

// ValidateWebhookSecret validates if the provided secret matches the project's webhook secret
func (p *Project) ValidateWebhookSecret(secret string) bool {
	return p.webhookSecret == secret
}

// ProjectDBData holds the data needed to reconstruct a Project entity from database
type ProjectDBData struct {
	ID             value_objects.ID
	Name           string
	RepositoryURL  string
	WebhookSecret  string
	TelegramChatID *int64
	Status         ProjectStatus
	CreatedAt      value_objects.Timestamp
	UpdatedAt      value_objects.Timestamp
}

// NewProjectFromDB creates a project entity from database data for repository reconstruction.
// This function is used to reconstruct a Project entity from persisted data without validation,
// as the data is assumed to be already validated when it was first stored.
//
// Parameters:
//   - data: ProjectDBData containing all the necessary fields to reconstruct the entity
//
// Returns a pointer to the reconstructed Project entity.
func NewProjectFromDB(data ProjectDBData) *Project {
	return &Project{
		id:             data.ID,
		name:           data.Name,
		repositoryURL:  data.RepositoryURL,
		webhookSecret:  data.WebhookSecret,
		telegramChatID: data.TelegramChatID,
		status:         data.Status,
		createdAt:      data.CreatedAt,
		updatedAt:      data.UpdatedAt,
	}
}
