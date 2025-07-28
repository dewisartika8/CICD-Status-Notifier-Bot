package models_test

import (
	"testing"

	"github.com/dewisartika8/cicd-status-notifier-bot/internal/domain/entities"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestProjectEntity_BusinessLogic(t *testing.T) {
	// Test project creation
	project := entities.NewProject("test-project", "https://github.com/user/repo", "secret")

	assert.NotEqual(t, uuid.Nil, project.ID)
	assert.Equal(t, "test-project", project.Name)
	assert.True(t, project.IsActive)

	// Test validation
	err := project.Validate()
	assert.NoError(t, err)

	// Test update
	chatID := int64(123456789)
	err = project.Update("updated-name", "https://github.com/updated/repo", &chatID, false)
	assert.NoError(t, err)
	assert.Equal(t, "updated-name", project.Name)
	assert.False(t, project.IsActive)
}

func TestBuildEventEntity_BusinessLogic(t *testing.T) {
	projectID := uuid.New()

	// Test build event creation
	buildEvent := entities.NewBuildEvent(projectID, entities.EventTypeBuildSuccess, entities.BuildStatusSuccess, "main")

	assert.NotEqual(t, uuid.Nil, buildEvent.ID)
	assert.Equal(t, projectID, buildEvent.ProjectID)
	assert.Equal(t, entities.EventTypeBuildSuccess, buildEvent.EventType)

	// Test validation
	err := buildEvent.Validate()
	assert.NoError(t, err)

	// Test setting commit info
	buildEvent.SetCommitInfo("abc123", "Initial commit", "John Doe", "john@example.com")
	assert.Equal(t, "abc123", buildEvent.CommitSHA)
	assert.Equal(t, "John Doe", buildEvent.AuthorName)

	// Test business logic methods
	assert.True(t, buildEvent.IsSuccessEvent())
	assert.False(t, buildEvent.IsFailureEvent())
}

func TestTelegramSubscriptionEntity_BusinessLogic(t *testing.T) {
	projectID := uuid.New()
	chatID := int64(123456789)

	// Test subscription creation
	subscription := entities.NewTelegramSubscription(projectID, chatID, nil, "testuser")

	assert.NotEqual(t, uuid.Nil, subscription.ID)
	assert.Equal(t, projectID, subscription.ProjectID)
	assert.Equal(t, chatID, subscription.ChatID)
	assert.True(t, subscription.IsActive)

	// Test validation
	err := subscription.Validate()
	assert.NoError(t, err)

	// Test subscription logic
	assert.True(t, subscription.IsSubscribedTo(entities.EventTypeBuildSuccess))
	assert.True(t, subscription.IsSubscribedTo(entities.EventTypeBuildFailed))

	// Test unsubscribe
	subscription.Unsubscribe()
	assert.False(t, subscription.IsActive)
	assert.False(t, subscription.IsSubscribedTo(entities.EventTypeBuildSuccess))
}

func TestNotificationLogEntity_BusinessLogic(t *testing.T) {
	buildEventID := uuid.New()
	chatID := int64(123456789)

	// Test notification log creation
	log := entities.NewNotificationLog(buildEventID, chatID)

	assert.NotEqual(t, uuid.Nil, log.ID)
	assert.Equal(t, buildEventID, log.BuildEventID)
	assert.Equal(t, chatID, log.ChatID)
	assert.Equal(t, entities.NotificationStatusPending, log.Status)

	// Test validation
	err := log.Validate()
	assert.NoError(t, err)

	// Test status changes
	assert.True(t, log.IsPending())
	assert.False(t, log.IsSent())
	assert.False(t, log.IsFailed())

	// Test mark as sent
	messageID := 42
	log.MarkAsSent(messageID)
	assert.True(t, log.IsSent())
	assert.Equal(t, &messageID, log.MessageID)
	assert.NotNil(t, log.SentAt)

	// Test mark as failed
	log.MarkAsFailed("Network error")
	assert.True(t, log.IsFailed())
	assert.Equal(t, "Network error", log.ErrorMessage)
	assert.Nil(t, log.MessageID)
}
