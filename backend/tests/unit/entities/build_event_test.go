package entities_test

import (
	"encoding/json"
	"testing"

	"github.com/dewisartika8/cicd-status-notifier-bot/internal/core/build/domain"
	"github.com/dewisartika8/cicd-status-notifier-bot/internal/core/shared/domain/value_objects"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestBuildEventNewBuildEvent(t *testing.T) {
	projectID := value_objects.NewID()
	eventType := domain.EventTypeBuildCompleted
	status := domain.BuildStatusSuccess
	branch := "main"

	params := domain.BuildEventParams{
		ProjectID: projectID,
		EventType: eventType,
		Status:    status,
		Branch:    branch,
	}
	buildEvent, err := domain.NewBuildEvent(params)
	assert.NoError(t, err)
	assert.Equal(t, projectID, buildEvent.ProjectID())
	assert.Equal(t, eventType, buildEvent.EventType())
	assert.Equal(t, status, buildEvent.Status())
	assert.Equal(t, branch, buildEvent.Branch())
	assert.NotZero(t, buildEvent.CreatedAt())
}

func TestBuildEventValidate(t *testing.T) {
	tests := []struct {
		name        string
		params      domain.BuildEventParams
		expectError bool
	}{
		{
			name: "valid build event",
			params: domain.BuildEventParams{
				ProjectID: value_objects.NewID(),
				EventType: domain.EventTypeBuildCompleted,
				Status:    domain.BuildStatusSuccess,
				Branch:    "main",
			},
			expectError: false,
		},
		{
			name: "invalid project ID",
			params: domain.BuildEventParams{
				ProjectID: value_objects.ID{},
				EventType: domain.EventTypeBuildCompleted,
				Status:    domain.BuildStatusSuccess,
				Branch:    "main",
			},
			expectError: true,
		},
		{
			name: "invalid event type",
			params: domain.BuildEventParams{
				ProjectID: value_objects.NewID(),
				EventType: "",
				Status:    domain.BuildStatusSuccess,
				Branch:    "main",
			},
			expectError: true,
		},
		{
			name: "invalid status",
			params: domain.BuildEventParams{
				ProjectID: value_objects.NewID(),
				EventType: domain.EventTypeBuildCompleted,
				Status:    "",
				Branch:    "main",
			},
			expectError: true,
		},
		{
			name: "invalid branch",
			params: domain.BuildEventParams{
				ProjectID: value_objects.NewID(),
				EventType: domain.EventTypeBuildCompleted,
				Status:    domain.BuildStatusSuccess,
				Branch:    "",
			},
			expectError: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			buildEvent, err := domain.NewBuildEvent(tt.params)
			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, buildEvent)
			}
		})
	}
}

func TestBuildEventSetDuration(t *testing.T) {
	buildURL := "https://github.com/user/repo/actions/runs/123"
	duration := 300
	params := domain.BuildEventParams{
		ProjectID: value_objects.NewID(),
		EventType: domain.EventTypeBuildCompleted,
		Status:    domain.BuildStatusSuccess,
		Branch:    "main",
		BuildURL:  buildURL,
	}
	buildEvent, _ := domain.NewBuildEvent(params)
	buildEvent.SetDuration(duration)

	assert.Equal(t, buildURL, buildEvent.BuildURL())
	assert.Equal(t, &duration, buildEvent.DurationSeconds())
}

func TestBuildEventSetWebhookPayload(t *testing.T) {
	payload := map[string]interface{}{
		"action": "completed",
		"workflow_run": map[string]interface{}{
			"id":     123,
			"status": "completed",
		},
	}

	payloadBytes, err := json.Marshal(payload)
	require.NoError(t, err)

	params := domain.BuildEventParams{
		ProjectID:      value_objects.NewID(),
		EventType:      domain.EventTypeBuildCompleted,
		Status:         domain.BuildStatusSuccess,
		Branch:         "main",
		WebhookPayload: json.RawMessage(payloadBytes),
	}
	buildEvent, _ := domain.NewBuildEvent(params)

	assert.NotNil(t, buildEvent.WebhookPayload())

	// Verify we can unmarshal it back
	var unmarshaledPayload map[string]interface{}
	err = json.Unmarshal(buildEvent.WebhookPayload(), &unmarshaledPayload)
	require.NoError(t, err)
	assert.Equal(t, "completed", unmarshaledPayload["action"])
}

func TestBuildEventIsFailureEvent(t *testing.T) {
	tests := []struct {
		name       string
		buildEvent *domain.BuildEvent
		expected   bool
	}{
		{
			name: "build failed status",
			buildEvent: func() *domain.BuildEvent {
				params := domain.BuildEventParams{
					ProjectID: value_objects.NewID(),
					EventType: domain.EventTypeBuildCompleted,
					Status:    domain.BuildStatusFailed,
					Branch:    "main",
				}
				be, _ := domain.NewBuildEvent(params)
				return be
			}(),
			expected: true,
		},
		{
			name: "successful build",
			buildEvent: func() *domain.BuildEvent {
				params := domain.BuildEventParams{
					ProjectID: value_objects.NewID(),
					EventType: domain.EventTypeBuildCompleted,
					Status:    domain.BuildStatusSuccess,
					Branch:    "main",
				}
				be, _ := domain.NewBuildEvent(params)
				return be
			}(),
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, tt.buildEvent.IsFailed())
		})
	}
}

func TestBuildEventIsSuccessEvent(t *testing.T) {
	tests := []struct {
		name       string
		buildEvent *domain.BuildEvent
		expected   bool
	}{
		{
			name: "build success status",
			buildEvent: func() *domain.BuildEvent {
				params := domain.BuildEventParams{
					ProjectID: value_objects.NewID(),
					EventType: domain.EventTypeBuildCompleted,
					Status:    domain.BuildStatusSuccess,
					Branch:    "main",
				}
				be, _ := domain.NewBuildEvent(params)
				return be
			}(),
			expected: true,
		},
		{
			name: "failed build",
			buildEvent: func() *domain.BuildEvent {
				params := domain.BuildEventParams{
					ProjectID: value_objects.NewID(),
					EventType: domain.EventTypeBuildCompleted,
					Status:    domain.BuildStatusFailed,
					Branch:    "main",
				}
				be, _ := domain.NewBuildEvent(params)
				return be
			}(),
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, tt.buildEvent.IsSuccessful())
		})
	}
}
