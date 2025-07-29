package models_test

import (
	"encoding/json"
	"testing"

	builddomain "github.com/dewisartika8/cicd-status-notifier-bot/internal/core/build/domain"
	projectdomain "github.com/dewisartika8/cicd-status-notifier-bot/internal/core/project/domain"
	"github.com/dewisartika8/cicd-status-notifier-bot/internal/core/shared/domain/value_objects"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

const (
	testCommitMessage = "Initial commit"
	testAuthorName    = "John Doe"
	testAuthorEmail   = "john@example.com"
	testBuildURL      = "https://github.com/user/repo/actions/runs/123"
)

func TestProjectModelToEntity(t *testing.T) {
	chatID := int64(123456789)
	now := projectdomain.NewProjectFromDB(projectdomain.ProjectDBData{}).CreatedAt().ToTime()
	model := &projectdomain.ProjectModel{
		ID:             uuid.New(),
		Name:           "test-project",
		RepositoryURL:  "https://github.com/user/repo",
		WebhookSecret:  "secret123",
		TelegramChatID: &chatID,
		Status:         string(projectdomain.ProjectStatusActive),
		CreatedAt:      now,
		UpdatedAt:      now,
	}
	entity, err := model.ToEntity()
	assert.NoError(t, err)
	assert.Equal(t, model.ID.String(), entity.ID().String())
	assert.Equal(t, model.Name, entity.Name())
	assert.Equal(t, model.RepositoryURL, entity.RepositoryURL())
	assert.Equal(t, model.WebhookSecret, entity.WebhookSecret())
	assert.Equal(t, model.TelegramChatID, entity.TelegramChatID())
	assert.Equal(t, projectdomain.ProjectStatusActive, entity.Status())
}

func TestProjectModelFromEntity(t *testing.T) {
	chatID := int64(123456789)
	entity, err := projectdomain.NewProject("test-project", "https://github.com/user/repo", "secret123", &chatID)
	assert.NoError(t, err)
	var model projectdomain.ProjectModel
	model.FromEntity(entity)
	assert.Equal(t, entity.ID().String(), model.ID.String())
	assert.Equal(t, entity.Name(), model.Name)
	assert.Equal(t, entity.RepositoryURL(), model.RepositoryURL)
	assert.Equal(t, entity.WebhookSecret(), model.WebhookSecret)
	assert.Equal(t, entity.TelegramChatID(), model.TelegramChatID)
	assert.Equal(t, string(entity.Status()), model.Status)
}

func TestBuildEventModelToEntity(t *testing.T) {
	id := uuid.New()
	projectID := uuid.New()
	duration := 300
	be, _ := builddomain.NewBuildEvent(builddomain.BuildEventParams{
		ProjectID:      value_objects.NewIDFromUUID(projectID),
		EventType:      builddomain.EventTypeBuildCompleted,
		Status:         builddomain.BuildStatusSuccess,
		Branch:         "main",
		CommitSHA:      "abc123",
		CommitMessage:  testCommitMessage,
		AuthorName:     testAuthorName,
		AuthorEmail:    testAuthorEmail,
		BuildURL:       testBuildURL,
		WebhookPayload: json.RawMessage(`{"action":"completed"}`),
	})
	model := &builddomain.BuildEventModel{
		ID:              id,
		ProjectID:       projectID,
		EventType:       string(builddomain.EventTypeBuildCompleted),
		Status:          string(builddomain.BuildStatusSuccess),
		Branch:          "main",
		CommitSHA:       "abc123",
		CommitMessage:   testCommitMessage,
		AuthorName:      testAuthorName,
		AuthorEmail:     testAuthorEmail,
		BuildURL:        testBuildURL,
		DurationSeconds: &duration,
		WebhookPayload:  json.RawMessage(`{"action":"completed"}`),
		CreatedAt:       be.CreatedAt().ToTime(),
		UpdatedAt:       be.CreatedAt().ToTime(),
	}
	entity := model.ToEntity()
	assert.Equal(t, id, entity.ID().Value())
	assert.Equal(t, projectID, entity.ProjectID().Value())
	assert.Equal(t, builddomain.EventTypeBuildCompleted, entity.EventType())
	assert.Equal(t, builddomain.BuildStatusSuccess, entity.Status())
	assert.Equal(t, model.Branch, entity.Branch())
	assert.Equal(t, model.CommitSHA, entity.CommitSHA())
	assert.Equal(t, model.CommitMessage, entity.CommitMessage())
	assert.Equal(t, model.AuthorName, entity.AuthorName())
	assert.Equal(t, model.AuthorEmail, entity.AuthorEmail())
	assert.Equal(t, model.BuildURL, entity.BuildURL())
	assert.Equal(t, model.DurationSeconds, entity.DurationSeconds())
	assert.NotNil(t, entity.WebhookPayload())
}

func TestBuildEventModelFromEntity(t *testing.T) {
	projectID := value_objects.NewID()
	params := builddomain.BuildEventParams{
		ProjectID:      projectID,
		EventType:      builddomain.EventTypeBuildCompleted,
		Status:         builddomain.BuildStatusSuccess,
		Branch:         "main",
		CommitSHA:      "abc123",
		CommitMessage:  testCommitMessage,
		AuthorName:     testAuthorName,
		AuthorEmail:    testAuthorEmail,
		BuildURL:       testBuildURL,
		WebhookPayload: json.RawMessage(`{"action":"completed"}`),
	}
	entity, _ := builddomain.NewBuildEvent(params)
	entity.SetDuration(300)
	var model builddomain.BuildEventModel
	model.FromEntity(entity)
	assert.Equal(t, entity.ID().Value(), model.ID)
	assert.Equal(t, entity.ProjectID().Value(), model.ProjectID)
	assert.Equal(t, string(entity.EventType()), model.EventType)
	assert.Equal(t, string(entity.Status()), model.Status)
	assert.Equal(t, entity.Branch(), model.Branch)
	assert.Equal(t, entity.CommitSHA(), model.CommitSHA)
	assert.Equal(t, entity.CommitMessage(), model.CommitMessage)
	assert.Equal(t, entity.AuthorName(), model.AuthorName)
	assert.Equal(t, entity.AuthorEmail(), model.AuthorEmail)
	assert.Equal(t, entity.BuildURL(), model.BuildURL)
	assert.Equal(t, entity.DurationSeconds(), model.DurationSeconds)
	assert.NotNil(t, model.WebhookPayload)
}
