package entities_test

import (
	"testing"
	"time"

	"github.com/dewisartika8/cicd-status-notifier-bot/internal/core/project/domain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestProject_NewProject(t *testing.T) {
	name := "test-project"
	repoURL := "https://github.com/user/repo"
	secret := "webhook-secret"

	project, err := domain.NewProject(name, repoURL, secret, nil)
	assert.NoError(t, err)
	assert.Equal(t, name, project.Name())
	assert.Equal(t, repoURL, project.RepositoryURL())
	assert.Equal(t, secret, project.WebhookSecret())
	assert.Equal(t, domain.ProjectStatusActive, project.Status())
	assert.NotZero(t, project.CreatedAt())
	assert.NotZero(t, project.UpdatedAt())
}

func TestProject_Validate(t *testing.T) {
	tests := []struct {
		name          string
		nameVal       string
		repoURLVal    string
		expectedError bool
	}{
		{"valid project", "test", "https://github.com/user/repo", false},
		{"empty name", "", "https://github.com/user/repo", true},
		{"empty repository URL", "test", "", true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			project, err := domain.NewProject(tt.nameVal, tt.repoURLVal, "secret", nil)
			if tt.expectedError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, project)
			}
		})
	}
}

func TestProject_Update(t *testing.T) {
	project, err := domain.NewProject("original", "https://github.com/user/repo", "secret", nil)
	require.NoError(t, err)
	originalCreatedAt := project.CreatedAt()
	originalUpdatedAt := project.UpdatedAt()

	time.Sleep(time.Millisecond * 10)
	chatID := int64(123456789)
	err = project.UpdateName("updated")
	require.NoError(t, err)
	err = project.UpdateRepositoryURL("https://github.com/user/updated")
	require.NoError(t, err)
	err = project.UpdateTelegramChatID(&chatID)
	require.NoError(t, err)
	assert.Equal(t, "updated", project.Name())
	assert.Equal(t, "https://github.com/user/updated", project.RepositoryURL())
	assert.Equal(t, &chatID, project.TelegramChatID())
	assert.Equal(t, domain.ProjectStatusActive, project.Status())
	assert.Equal(t, originalCreatedAt, project.CreatedAt())
	assert.True(t, project.UpdatedAt().After(originalUpdatedAt) || !project.UpdatedAt().Before(originalUpdatedAt))
}
