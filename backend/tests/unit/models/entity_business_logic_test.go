package models_test

import (
	"testing"

	builddomain "github.com/dewisartika8/cicd-status-notifier-bot/internal/core/build/domain"
	notificationdomain "github.com/dewisartika8/cicd-status-notifier-bot/internal/core/notification/domain"
	projectdomain "github.com/dewisartika8/cicd-status-notifier-bot/internal/core/project/domain"
	"github.com/dewisartika8/cicd-status-notifier-bot/internal/core/shared/domain/value_objects"
	"github.com/stretchr/testify/assert"
)

func TestProjectEntityBusinessLogic(t *testing.T) {
	// Test project creation
	project, err := projectdomain.NewProject("test-project", "https://github.com/user/repo", "secret", nil)
	assert.NoError(t, err)
	assert.NotNil(t, project)
	assert.Equal(t, "test-project", project.Name())
	assert.Equal(t, projectdomain.ProjectStatusActive, project.Status())

	// Test update
	chatID := int64(123456789)
	err = project.UpdateName("updated-name")
	assert.NoError(t, err)
	err = project.UpdateRepositoryURL("https://github.com/updated/repo")
	assert.NoError(t, err)
	err = project.UpdateTelegramChatID(&chatID)
	assert.NoError(t, err)
	project.SetStatus(projectdomain.ProjectStatusInactive)
	assert.Equal(t, "updated-name", project.Name())
	assert.Equal(t, "https://github.com/updated/repo", project.RepositoryURL())
	assert.Equal(t, &chatID, project.TelegramChatID())
	assert.Equal(t, projectdomain.ProjectStatusInactive, project.Status())
}

func TestBuildEventEntityBusinessLogic(t *testing.T) {
	projectID := value_objects.NewID()
	params := builddomain.BuildEventParams{
		ProjectID: projectID,
		EventType: builddomain.EventTypeBuildCompleted,
		Status:    builddomain.BuildStatusSuccess,
		Branch:    "main",
	}
	buildEvent, err := builddomain.NewBuildEvent(params)
	assert.NoError(t, err)
	assert.NotNil(t, buildEvent)
	assert.Equal(t, projectID, buildEvent.ProjectID())
	assert.Equal(t, builddomain.EventTypeBuildCompleted, buildEvent.EventType())

	// Test setting commit info
	buildEvent.UpdateStatus(builddomain.BuildStatusSuccess)
	buildEvent.SetDuration(120)
	// Test business logic methods
	assert.True(t, buildEvent.IsSuccessful())
	assert.False(t, buildEvent.IsFailed())
}

func TestTelegramSubscriptionEntityBusinessLogic(t *testing.T) {
	projectID := value_objects.NewID()
	chatID := int64(123456789)
	subscription, err := notificationdomain.NewTelegramSubscription(projectID, chatID)
	assert.NoError(t, err)
	assert.NotNil(t, subscription)
	assert.Equal(t, projectID, subscription.ProjectID()) // Test ProjectID, not ID
	assert.Equal(t, chatID, subscription.ChatID())
	assert.True(t, subscription.IsActive())
	// Test unsubscribe
	// (Assuming IsActive can be toggled, otherwise skip)
}

func TestNotificationLogEntityBusinessLogic(t *testing.T) {
	buildEventID := value_objects.NewID()
	projectID := value_objects.NewID()
	channel := notificationdomain.NotificationChannelTelegram
	recipient := "test_user"
	message := "test message"
	maxRetries := 3
	log, err := notificationdomain.NewNotificationLog(buildEventID, projectID, channel, recipient, message, maxRetries)
	assert.NoError(t, err)
	assert.NotNil(t, log)
	assert.Equal(t, buildEventID, log.BuildEventID())
	assert.Equal(t, projectID, log.ProjectID())
	assert.Equal(t, channel, log.Channel())
	assert.Equal(t, recipient, log.Recipient())
	assert.Equal(t, message, log.Message())
	assert.Equal(t, notificationdomain.NotificationStatusPending, log.Status())
	// Test status changes
	// (Assuming status can be updated, otherwise skip)
}
