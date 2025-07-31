package domain_test

import (
	"testing"
	"time"

	"github.com/dewisartika8/cicd-status-notifier-bot/internal/core/notification/domain"
	"github.com/dewisartika8/cicd-status-notifier-bot/internal/core/shared/domain/value_objects"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewTelegramSubscription(t *testing.T) {
	projectID := value_objects.NewID()
	chatID := int64(123456789)

	t.Run("valid subscription creation", func(t *testing.T) {
		subscription, err := domain.NewTelegramSubscription(projectID, chatID)

		require.NoError(t, err)
		assert.NotNil(t, subscription)
		assert.Equal(t, projectID, subscription.ProjectID())
		assert.Equal(t, chatID, subscription.ChatID())
		assert.True(t, subscription.IsActive())
		assert.Empty(t, subscription.EventTypes())
		assert.Empty(t, subscription.Username())
		assert.Nil(t, subscription.UserID())
		assert.False(t, subscription.ID().IsNil())
		assert.False(t, subscription.CreatedAt().IsZero())
		assert.False(t, subscription.UpdatedAt().IsZero())
	})

	t.Run("invalid chat ID - zero", func(t *testing.T) {
		subscription, err := domain.NewTelegramSubscription(projectID, 0)

		assert.Error(t, err)
		assert.Nil(t, subscription)
		assert.Contains(t, err.Error(), "chat ID is required and cannot be zero")
	})

	t.Run("invalid project ID - nil", func(t *testing.T) {
		nilProjectID := value_objects.ID{}
		subscription, err := domain.NewTelegramSubscription(nilProjectID, chatID)

		assert.Error(t, err)
		assert.Nil(t, subscription)
	})
}

func TestTelegramSubscriptionBusinessMethods(t *testing.T) {
	projectID := value_objects.NewID()
	chatID := int64(123456789)
	subscription, err := domain.NewTelegramSubscription(projectID, chatID)
	require.NoError(t, err)

	t.Run("activate subscription", func(t *testing.T) {
		// Deactivate first
		subscription.Deactivate()
		assert.False(t, subscription.IsActive())

		// Then activate
		oldUpdatedAt := subscription.UpdatedAt()
		time.Sleep(1 * time.Millisecond) // Ensure different timestamp
		subscription.Activate()

		assert.True(t, subscription.IsActive())
		assert.True(t, subscription.UpdatedAt().After(oldUpdatedAt))
	})

	t.Run("deactivate subscription", func(t *testing.T) {
		// Ensure it's active first
		subscription.Activate()
		assert.True(t, subscription.IsActive())

		// Then deactivate
		oldUpdatedAt := subscription.UpdatedAt()
		time.Sleep(1 * time.Millisecond) // Ensure different timestamp
		subscription.Deactivate()

		assert.False(t, subscription.IsActive())
		assert.True(t, subscription.UpdatedAt().After(oldUpdatedAt))
	})

	t.Run("update chat ID - valid", func(t *testing.T) {
		newChatID := int64(987654321)
		oldUpdatedAt := subscription.UpdatedAt()
		time.Sleep(1 * time.Millisecond) // Ensure different timestamp

		err := subscription.UpdateChatID(newChatID)

		require.NoError(t, err)
		assert.Equal(t, newChatID, subscription.ChatID())
		assert.True(t, subscription.UpdatedAt().After(oldUpdatedAt))
	})

	t.Run("update chat ID - invalid zero", func(t *testing.T) {
		err := subscription.UpdateChatID(0)

		assert.Error(t, err)
		assert.Equal(t, domain.ErrInvalidTelegramChatID, err)
	})

	t.Run("get chat ID as string", func(t *testing.T) {
		chatIDStr := subscription.GetChatIDString()
		assert.Equal(t, "987654321", chatIDStr) // From previous test
	})
}

func TestTelegramSubscriptionRestoreFromPersistence(t *testing.T) {
	id := value_objects.NewID()
	projectID := value_objects.NewID()
	chatID := int64(123456789)
	userID := int64(987654321)
	username := "testuser"
	eventTypes := []string{"build_success", "build_failed"}
	isActive := true
	createdAt := value_objects.NewTimestamp()
	updatedAt := value_objects.NewTimestamp()

	params := domain.RestoreTelegramSubscriptionParams{
		ID:         id,
		ProjectID:  projectID,
		ChatID:     chatID,
		UserID:     &userID,
		Username:   username,
		EventTypes: eventTypes,
		IsActive:   isActive,
		CreatedAt:  createdAt,
		UpdatedAt:  updatedAt,
	}

	subscription := domain.RestoreTelegramSubscription(params)

	assert.Equal(t, id, subscription.ID())
	assert.Equal(t, projectID, subscription.ProjectID())
	assert.Equal(t, chatID, subscription.ChatID())
	assert.Equal(t, &userID, subscription.UserID())
	assert.Equal(t, username, subscription.Username())
	assert.Equal(t, eventTypes, subscription.EventTypes())
	assert.Equal(t, isActive, subscription.IsActive())
	assert.Equal(t, createdAt, subscription.CreatedAt())
	assert.Equal(t, updatedAt, subscription.UpdatedAt())
}

func TestTelegramSubscriptionStringRepresentation(t *testing.T) {
	projectID := value_objects.NewID()
	chatID := int64(123456789)
	subscription, err := domain.NewTelegramSubscription(projectID, chatID)
	require.NoError(t, err)

	str := subscription.String()
	assert.Contains(t, str, "TelegramSubscription")
	assert.Contains(t, str, subscription.ID().String())
	assert.Contains(t, str, projectID.String())
	assert.Contains(t, str, "123456789")
	assert.Contains(t, str, "active")

	// Test inactive subscription
	subscription.Deactivate()
	str = subscription.String()
	assert.Contains(t, str, "inactive")
}

func TestTelegramSubscriptionValidationEdgeCases(t *testing.T) {
	projectID := value_objects.NewID()

	t.Run("negative chat ID should be valid for groups", func(t *testing.T) {
		chatID := int64(-123456789) // Negative for groups
		subscription, err := domain.NewTelegramSubscription(projectID, chatID)

		require.NoError(t, err)
		assert.NotNil(t, subscription)
		assert.Equal(t, chatID, subscription.ChatID())
	})

	t.Run("large positive chat ID should be valid", func(t *testing.T) {
		chatID := int64(9223372036854775807) // Max int64
		subscription, err := domain.NewTelegramSubscription(projectID, chatID)

		require.NoError(t, err)
		assert.NotNil(t, subscription)
		assert.Equal(t, chatID, subscription.ChatID())
	})
}
