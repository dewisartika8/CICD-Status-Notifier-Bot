package entities_test

import (
	"encoding/json"
	"testing"

	"github.com/dewisartika8/cicd-status-notifier-bot/internal/domain/entities"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestBuildEvent_NewBuildEvent(t *testing.T) {
	projectID := uuid.New()
	eventType := entities.EventTypeBuildSuccess
	status := entities.BuildStatusSuccess
	branch := "main"

	buildEvent := entities.NewBuildEvent(projectID, eventType, status, branch)

	assert.NotEqual(t, uuid.Nil, buildEvent.ID)
	assert.Equal(t, projectID, buildEvent.ProjectID)
	assert.Equal(t, eventType, buildEvent.EventType)
	assert.Equal(t, status, buildEvent.Status)
	assert.Equal(t, branch, buildEvent.Branch)
	assert.NotZero(t, buildEvent.CreatedAt)
}

func TestBuildEvent_Validate(t *testing.T) {
	tests := []struct {
		name          string
		buildEvent    *entities.BuildEvent
		expectedError error
	}{
		{
			name:          "valid build event",
			buildEvent:    entities.NewBuildEvent(uuid.New(), entities.EventTypeBuildSuccess, entities.BuildStatusSuccess, "main"),
			expectedError: nil,
		},
		{
			name: "invalid project ID",
			buildEvent: &entities.BuildEvent{
				ProjectID: uuid.Nil,
				EventType: entities.EventTypeBuildSuccess,
				Status:    entities.BuildStatusSuccess,
				Branch:    "main",
			},
			expectedError: entities.ErrInvalidProjectID,
		},
		{
			name: "invalid event type",
			buildEvent: &entities.BuildEvent{
				ProjectID: uuid.New(),
				EventType: "",
				Status:    entities.BuildStatusSuccess,
				Branch:    "main",
			},
			expectedError: entities.ErrInvalidEventType,
		},
		{
			name: "invalid status",
			buildEvent: &entities.BuildEvent{
				ProjectID: uuid.New(),
				EventType: entities.EventTypeBuildSuccess,
				Status:    "",
				Branch:    "main",
			},
			expectedError: entities.ErrInvalidBuildStatus,
		},
		{
			name: "invalid branch",
			buildEvent: &entities.BuildEvent{
				ProjectID: uuid.New(),
				EventType: entities.EventTypeBuildSuccess,
				Status:    entities.BuildStatusSuccess,
				Branch:    "",
			},
			expectedError: entities.ErrInvalidBranch,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.buildEvent.Validate()
			if tt.expectedError != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.expectedError, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestBuildEvent_SetCommitInfo(t *testing.T) {
	buildEvent := entities.NewBuildEvent(uuid.New(), entities.EventTypeBuildSuccess, entities.BuildStatusSuccess, "main")

	sha := "abc123"
	message := "Initial commit"
	authorName := "John Doe"
	authorEmail := "john@example.com"

	buildEvent.SetCommitInfo(sha, message, authorName, authorEmail)

	assert.Equal(t, sha, buildEvent.CommitSHA)
	assert.Equal(t, message, buildEvent.CommitMessage)
	assert.Equal(t, authorName, buildEvent.AuthorName)
	assert.Equal(t, authorEmail, buildEvent.AuthorEmail)
}

func TestBuildEvent_SetBuildInfo(t *testing.T) {
	buildEvent := entities.NewBuildEvent(uuid.New(), entities.EventTypeBuildSuccess, entities.BuildStatusSuccess, "main")

	buildURL := "https://github.com/user/repo/actions/runs/123"
	duration := 300

	buildEvent.SetBuildInfo(buildURL, &duration)

	assert.Equal(t, buildURL, buildEvent.BuildURL)
	assert.Equal(t, &duration, buildEvent.DurationSeconds)
}

func TestBuildEvent_SetWebhookPayload(t *testing.T) {
	buildEvent := entities.NewBuildEvent(uuid.New(), entities.EventTypeBuildSuccess, entities.BuildStatusSuccess, "main")

	payload := map[string]interface{}{
		"action": "completed",
		"workflow_run": map[string]interface{}{
			"id":     123,
			"status": "completed",
		},
	}

	payloadBytes, err := json.Marshal(payload)
	require.NoError(t, err)

	buildEvent.SetWebhookPayload(json.RawMessage(payloadBytes))

	assert.NotNil(t, buildEvent.WebhookPayload)

	// Verify we can unmarshal it back
	var unmarshaledPayload map[string]interface{}
	err = json.Unmarshal(buildEvent.WebhookPayload, &unmarshaledPayload)
	require.NoError(t, err)
	assert.Equal(t, "completed", unmarshaledPayload["action"])
}

func TestBuildEvent_IsFailureEvent(t *testing.T) {
	tests := []struct {
		name       string
		buildEvent *entities.BuildEvent
		expected   bool
	}{
		{
			name:       "build failed status",
			buildEvent: entities.NewBuildEvent(uuid.New(), entities.EventTypeBuildSuccess, entities.BuildStatusFailed, "main"),
			expected:   true,
		},
		{
			name:       "build failed event type",
			buildEvent: entities.NewBuildEvent(uuid.New(), entities.EventTypeBuildFailed, entities.BuildStatusSuccess, "main"),
			expected:   true,
		},
		{
			name:       "test failed event type",
			buildEvent: entities.NewBuildEvent(uuid.New(), entities.EventTypeTestFailed, entities.BuildStatusSuccess, "main"),
			expected:   true,
		},
		{
			name:       "deployment failed event type",
			buildEvent: entities.NewBuildEvent(uuid.New(), entities.EventTypeDeploymentFailed, entities.BuildStatusSuccess, "main"),
			expected:   true,
		},
		{
			name:       "successful build",
			buildEvent: entities.NewBuildEvent(uuid.New(), entities.EventTypeBuildSuccess, entities.BuildStatusSuccess, "main"),
			expected:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, tt.buildEvent.IsFailureEvent())
		})
	}
}

func TestBuildEvent_IsSuccessEvent(t *testing.T) {
	tests := []struct {
		name       string
		buildEvent *entities.BuildEvent
		expected   bool
	}{
		{
			name:       "build success status",
			buildEvent: entities.NewBuildEvent(uuid.New(), entities.EventTypeBuildFailed, entities.BuildStatusSuccess, "main"),
			expected:   true,
		},
		{
			name:       "build success event type",
			buildEvent: entities.NewBuildEvent(uuid.New(), entities.EventTypeBuildSuccess, entities.BuildStatusFailed, "main"),
			expected:   true,
		},
		{
			name:       "test passed event type",
			buildEvent: entities.NewBuildEvent(uuid.New(), entities.EventTypeTestPassed, entities.BuildStatusFailed, "main"),
			expected:   true,
		},
		{
			name:       "deployment success event type",
			buildEvent: entities.NewBuildEvent(uuid.New(), entities.EventTypeDeploymentSuccess, entities.BuildStatusFailed, "main"),
			expected:   true,
		},
		{
			name:       "failed build",
			buildEvent: entities.NewBuildEvent(uuid.New(), entities.EventTypeBuildFailed, entities.BuildStatusFailed, "main"),
			expected:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, tt.buildEvent.IsSuccessEvent())
		})
	}
}
