package service_test

import (
	"context"
	"testing"

	"github.com/dewisartika8/cicd-status-notifier-bot/internal/core/notification/domain"
	"github.com/dewisartika8/cicd-status-notifier-bot/internal/core/notification/service/log"
	"github.com/dewisartika8/cicd-status-notifier-bot/internal/core/shared/domain/value_objects"
	"github.com/dewisartika8/cicd-status-notifier-bot/tests/mocks"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

const (
	notificationLogType = "*domain.NotificationLog"
	buildCompletedMsg   = "Build completed"
)

func TestBasicNotificationIntegration(t *testing.T) {
	t.Run("create notification for build event with subscription", func(t *testing.T) {
		mockLogRepo := mocks.NewNotificationLogRepository(t)
		mockSubRepo := mocks.NewTelegramSubscriptionRepository(t)
		logger := logrus.New()
		logger.SetLevel(logrus.FatalLevel)

		service := log.NewNotificationLogService(log.Dep{
			NotificationRepo:         mockLogRepo,
			TelegramSubscriptionRepo: mockSubRepo,
			Logger:                   logger,
		})

		buildEventID := value_objects.NewID()
		projectID := value_objects.NewID()
		message := "Build completed successfully"

		// Mock active subscriptions
		activeSubscriptions := []*domain.TelegramSubscription{
			domain.RestoreTelegramSubscription(domain.RestoreTelegramSubscriptionParams{
				ID:        value_objects.NewID(),
				ProjectID: projectID,
				ChatID:    int64(123456789),
				IsActive:  true,
				CreatedAt: value_objects.NewTimestamp(),
				UpdatedAt: value_objects.NewTimestamp(),
			}),
		}

		mockSubRepo.On("GetActiveSubscriptionsByProject", mock.Anything, projectID).
			Return(activeSubscriptions, nil)
		mockLogRepo.On("Create", mock.Anything, mock.AnythingOfType(notificationLogType)).
			Return(nil)

		result, err := service.CreateNotificationForBuildEvent(context.Background(), buildEventID, projectID, message)

		require.NoError(t, err)
		assert.Len(t, result, 1)
		assert.Equal(t, buildEventID, result[0].BuildEventID())
		assert.Equal(t, projectID, result[0].ProjectID())
		assert.Equal(t, message, result[0].Message())
		assert.Equal(t, domain.NotificationChannelTelegram, result[0].Channel())
		assert.Equal(t, "123456789", result[0].Recipient())

		mockSubRepo.AssertExpectations(t)
		mockLogRepo.AssertExpectations(t)
	})

	t.Run("create notifications for multiple subscriptions", func(t *testing.T) {
		mockLogRepo := mocks.NewNotificationLogRepository(t)
		mockSubRepo := mocks.NewTelegramSubscriptionRepository(t)
		logger := logrus.New()
		logger.SetLevel(logrus.FatalLevel)

		service := log.NewNotificationLogService(log.Dep{
			NotificationRepo:         mockLogRepo,
			TelegramSubscriptionRepo: mockSubRepo,
			Logger:                   logger,
		})

		buildEventID := value_objects.NewID()
		projectID := value_objects.NewID()
		message := "Build failed"

		// Mock multiple active subscriptions
		activeSubscriptions := []*domain.TelegramSubscription{
			domain.RestoreTelegramSubscription(domain.RestoreTelegramSubscriptionParams{
				ID:        value_objects.NewID(),
				ProjectID: projectID,
				ChatID:    int64(123456789),
				IsActive:  true,
				CreatedAt: value_objects.NewTimestamp(),
				UpdatedAt: value_objects.NewTimestamp(),
			}),
			domain.RestoreTelegramSubscription(domain.RestoreTelegramSubscriptionParams{
				ID:        value_objects.NewID(),
				ProjectID: projectID,
				ChatID:    int64(987654321),
				IsActive:  true,
				CreatedAt: value_objects.NewTimestamp(),
				UpdatedAt: value_objects.NewTimestamp(),
			}),
		}

		mockSubRepo.On("GetActiveSubscriptionsByProject", mock.Anything, projectID).
			Return(activeSubscriptions, nil)
		mockLogRepo.On("Create", mock.Anything, mock.AnythingOfType(notificationLogType)).
			Return(nil).Times(2)

		result, err := service.CreateNotificationForBuildEvent(context.Background(), buildEventID, projectID, message)

		require.NoError(t, err)
		assert.Len(t, result, 2)

		// Check that each subscription received a notification
		recipients := make([]string, len(result))
		for i, notification := range result {
			recipients[i] = notification.Recipient()
			assert.Equal(t, buildEventID, notification.BuildEventID())
			assert.Equal(t, projectID, notification.ProjectID())
			assert.Equal(t, message, notification.Message())
			assert.Equal(t, domain.NotificationChannelTelegram, notification.Channel())
		}

		assert.Contains(t, recipients, "123456789")
		assert.Contains(t, recipients, "987654321")

		mockSubRepo.AssertExpectations(t)
		mockLogRepo.AssertExpectations(t)
	})

	t.Run("no active subscriptions for project", func(t *testing.T) {
		mockLogRepo := mocks.NewNotificationLogRepository(t)
		mockSubRepo := mocks.NewTelegramSubscriptionRepository(t)
		logger := logrus.New()
		logger.SetLevel(logrus.FatalLevel)

		service := log.NewNotificationLogService(log.Dep{
			NotificationRepo:         mockLogRepo,
			TelegramSubscriptionRepo: mockSubRepo,
			Logger:                   logger,
		})

		buildEventID := value_objects.NewID()
		projectID := value_objects.NewID()
		message := buildCompletedMsg

		// Mock no active subscriptions
		mockSubRepo.On("GetActiveSubscriptionsByProject", mock.Anything, projectID).
			Return([]*domain.TelegramSubscription{}, nil)

		result, err := service.CreateNotificationForBuildEvent(context.Background(), buildEventID, projectID, message)

		require.NoError(t, err)
		assert.Len(t, result, 0)

		mockSubRepo.AssertExpectations(t)
		mockLogRepo.AssertExpectations(t)
	})

	t.Run("repository error getting subscriptions", func(t *testing.T) {
		mockLogRepo := mocks.NewNotificationLogRepository(t)
		mockSubRepo := mocks.NewTelegramSubscriptionRepository(t)
		logger := logrus.New()
		logger.SetLevel(logrus.FatalLevel)

		service := log.NewNotificationLogService(log.Dep{
			NotificationRepo:         mockLogRepo,
			TelegramSubscriptionRepo: mockSubRepo,
			Logger:                   logger,
		})

		buildEventID := value_objects.NewID()
		projectID := value_objects.NewID()
		message := buildCompletedMsg

		mockSubRepo.On("GetActiveSubscriptionsByProject", mock.Anything, projectID).
			Return(nil, assert.AnError)

		result, err := service.CreateNotificationForBuildEvent(context.Background(), buildEventID, projectID, message)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "failed to get active subscriptions")

		mockSubRepo.AssertExpectations(t)
	})

	t.Run("repository error creating notification log", func(t *testing.T) {
		mockLogRepo := mocks.NewNotificationLogRepository(t)
		mockSubRepo := mocks.NewTelegramSubscriptionRepository(t)
		logger := logrus.New()
		logger.SetLevel(logrus.FatalLevel)

		service := log.NewNotificationLogService(log.Dep{
			NotificationRepo:         mockLogRepo,
			TelegramSubscriptionRepo: mockSubRepo,
			Logger:                   logger,
		})

		buildEventID := value_objects.NewID()
		projectID := value_objects.NewID()
		message := buildCompletedMsg

		activeSubscriptions := []*domain.TelegramSubscription{
			domain.RestoreTelegramSubscription(domain.RestoreTelegramSubscriptionParams{
				ID:        value_objects.NewID(),
				ProjectID: projectID,
				ChatID:    int64(123456789),
				IsActive:  true,
				CreatedAt: value_objects.NewTimestamp(),
				UpdatedAt: value_objects.NewTimestamp(),
			}),
		}

		mockSubRepo.On("GetActiveSubscriptionsByProject", mock.Anything, projectID).
			Return(activeSubscriptions, nil)
		mockLogRepo.On("Create", mock.Anything, mock.AnythingOfType(notificationLogType)).
			Return(assert.AnError)

		result, err := service.CreateNotificationForBuildEvent(context.Background(), buildEventID, projectID, message)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "failed to create notification log")

		mockSubRepo.AssertExpectations(t)
		mockLogRepo.AssertExpectations(t)
	})
}
