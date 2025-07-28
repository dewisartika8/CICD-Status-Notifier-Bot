package models_test

import (
	"testing"

	"github.com/dewisartika8/cicd-status-notifier-bot/internal/adapters/database"
	"github.com/dewisartika8/cicd-status-notifier-bot/internal/domain/entities"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestProjectModel_ToEntity(t *testing.T) {
	id := uuid.New()
	chatID := int64(123456789)

	model := &database.ProjectModel{
		ID:             id,
		Name:           "test-project",
		RepositoryURL:  "https://github.com/user/repo",
		WebhookSecret:  "secret123",
		TelegramChatID: &chatID,
		IsActive:       true,
	}

	entity := model.ToEntity()

	assert.Equal(t, model.ID, entity.ID)
	assert.Equal(t, model.Name, entity.Name)
	assert.Equal(t, model.RepositoryURL, entity.RepositoryURL)
	assert.Equal(t, model.WebhookSecret, entity.WebhookSecret)
	assert.Equal(t, model.TelegramChatID, entity.TelegramChatID)
	assert.Equal(t, model.IsActive, entity.IsActive)
	assert.Equal(t, model.CreatedAt, entity.CreatedAt)
	assert.Equal(t, model.UpdatedAt, entity.UpdatedAt)
}

func TestProjectModel_FromEntity(t *testing.T) {
	entity := entities.NewProject("test-project", "https://github.com/user/repo", "secret123")
	chatID := int64(123456789)
	entity.TelegramChatID = &chatID

	var model database.ProjectModel
	model.FromEntity(entity)

	assert.Equal(t, entity.ID, model.ID)
	assert.Equal(t, entity.Name, model.Name)
	assert.Equal(t, entity.RepositoryURL, model.RepositoryURL)
	assert.Equal(t, entity.WebhookSecret, model.WebhookSecret)
	assert.Equal(t, entity.TelegramChatID, model.TelegramChatID)
	assert.Equal(t, entity.IsActive, model.IsActive)
	assert.Equal(t, entity.CreatedAt, model.CreatedAt)
	assert.Equal(t, entity.UpdatedAt, model.UpdatedAt)
}

func TestBuildEventModel_ToEntity(t *testing.T) {
	id := uuid.New()
	projectID := uuid.New()
	duration := 300

	model := &database.BuildEventModel{
		ID:              id,
		ProjectID:       projectID,
		EventType:       string(entities.EventTypeBuildSuccess),
		Status:          string(entities.BuildStatusSuccess),
		Branch:          "main",
		CommitSHA:       "abc123",
		CommitMessage:   "Initial commit",
		AuthorName:      "John Doe",
		AuthorEmail:     "john@example.com",
		BuildURL:        "https://github.com/user/repo/actions/runs/123",
		DurationSeconds: &duration,
		WebhookPayload:  database.JSONB(`{"action":"completed"}`),
	}

	entity := model.ToEntity()

	assert.Equal(t, model.ID, entity.ID)
	assert.Equal(t, model.ProjectID, entity.ProjectID)
	assert.Equal(t, entities.EventTypeBuildSuccess, entity.EventType)
	assert.Equal(t, entities.BuildStatusSuccess, entity.Status)
	assert.Equal(t, model.Branch, entity.Branch)
	assert.Equal(t, model.CommitSHA, entity.CommitSHA)
	assert.Equal(t, model.CommitMessage, entity.CommitMessage)
	assert.Equal(t, model.AuthorName, entity.AuthorName)
	assert.Equal(t, model.AuthorEmail, entity.AuthorEmail)
	assert.Equal(t, model.BuildURL, entity.BuildURL)
	assert.Equal(t, model.DurationSeconds, entity.DurationSeconds)
	assert.NotNil(t, entity.WebhookPayload)
}

func TestBuildEventModel_FromEntity(t *testing.T) {
	entity := entities.NewBuildEvent(uuid.New(), entities.EventTypeBuildSuccess, entities.BuildStatusSuccess, "main")
	entity.SetCommitInfo("abc123", "Initial commit", "John Doe", "john@example.com")
	entity.SetBuildInfo("https://github.com/user/repo/actions/runs/123", &[]int{300}[0])
	entity.SetWebhookPayload([]byte(`{"action":"completed"}`))

	var model database.BuildEventModel
	model.FromEntity(entity)

	assert.Equal(t, entity.ID, model.ID)
	assert.Equal(t, entity.ProjectID, model.ProjectID)
	assert.Equal(t, string(entity.EventType), model.EventType)
	assert.Equal(t, string(entity.Status), model.Status)
	assert.Equal(t, entity.Branch, model.Branch)
	assert.Equal(t, entity.CommitSHA, model.CommitSHA)
	assert.Equal(t, entity.CommitMessage, model.CommitMessage)
	assert.Equal(t, entity.AuthorName, model.AuthorName)
	assert.Equal(t, entity.AuthorEmail, model.AuthorEmail)
	assert.Equal(t, entity.BuildURL, model.BuildURL)
	assert.Equal(t, entity.DurationSeconds, model.DurationSeconds)
	assert.Equal(t, entity.CreatedAt, model.CreatedAt)
}
