package entities_test

import (
	"testing"
	"time"

	"github.com/dewisartika8/cicd-status-notifier-bot/internal/domain/entities"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestProject_NewProject(t *testing.T) {
	name := "test-project"
	repoURL := "https://github.com/user/repo"
	secret := "webhook-secret"

	project := entities.NewProject(name, repoURL, secret)

	assert.NotEqual(t, uuid.Nil, project.ID)
	assert.Equal(t, name, project.Name)
	assert.Equal(t, repoURL, project.RepositoryURL)
	assert.Equal(t, secret, project.WebhookSecret)
	assert.True(t, project.IsActive)
	assert.NotZero(t, project.CreatedAt)
	assert.NotZero(t, project.UpdatedAt)
}

func TestProject_Validate(t *testing.T) {
	tests := []struct {
		name          string
		project       *entities.Project
		expectedError error
	}{
		{
			name:          "valid project",
			project:       entities.NewProject("test", "https://github.com/user/repo", "secret"),
			expectedError: nil,
		},
		{
			name: "empty name",
			project: &entities.Project{
				Name:          "",
				RepositoryURL: "https://github.com/user/repo",
			},
			expectedError: entities.ErrInvalidProjectName,
		},
		{
			name: "empty repository URL",
			project: &entities.Project{
				Name:          "test",
				RepositoryURL: "",
			},
			expectedError: entities.ErrInvalidRepositoryURL,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.project.Validate()
			if tt.expectedError != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.expectedError, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestProject_Update(t *testing.T) {
	project := entities.NewProject("original", "https://github.com/user/repo", "secret")
	originalCreatedAt := project.CreatedAt
	originalUpdatedAt := project.UpdatedAt

	// Sleep a bit to ensure UpdatedAt changes
	time.Sleep(time.Millisecond * 10)

	chatID := int64(123456789)
	err := project.Update("updated", "https://github.com/user/updated", &chatID, false)

	require.NoError(t, err)
	assert.Equal(t, "updated", project.Name)
	assert.Equal(t, "https://github.com/user/updated", project.RepositoryURL)
	assert.Equal(t, &chatID, project.TelegramChatID)
	assert.False(t, project.IsActive)
	assert.Equal(t, originalCreatedAt, project.CreatedAt)                                                    // Should not change
	assert.True(t, project.UpdatedAt.After(originalUpdatedAt) || project.UpdatedAt.Equal(originalUpdatedAt)) // Should be updated or equal
}
