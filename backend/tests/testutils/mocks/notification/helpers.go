package notification

import (
	"context"

	"github.com/dewisartika8/cicd-status-notifier-bot/internal/core/notification/port"
	"github.com/dewisartika8/cicd-status-notifier-bot/internal/core/notification/service/subscription"
	"github.com/dewisartika8/cicd-status-notifier-bot/internal/core/shared/domain/value_objects"
	"github.com/sirupsen/logrus"
)

// TestHelpers provides utility functions for testing
type TestHelpers struct{}

// NewTestHelpers creates a new instance of test helpers
func NewTestHelpers() *TestHelpers {
	return &TestHelpers{}
}

// CreateTelegramSubscriptionService creates a service instance with mocked dependencies for testing
func (h *TestHelpers) CreateTelegramSubscriptionService() (port.TelegramSubscriptionService, *MockTelegramRepository) {
	mockRepo := new(MockTelegramRepository)
	logger := logrus.New()
	logger.SetLevel(logrus.FatalLevel) // Suppress logs in tests

	service := subscription.NewTelegramSubscriptionService(subscription.Dep{
		TelegramRepo: mockRepo,
		Logger:       logger,
	})

	return service, mockRepo
}

// ValidationService interface for accessing validation methods
type ValidationService interface {
	ValidateUserPermissions(ctx context.Context, userID int64, projectID value_objects.ID, chatID int64) error
	ValidateProjectExistence(ctx context.Context, projectID value_objects.ID) error
	ValidateDuplicateSubscription(ctx context.Context, projectID value_objects.ID, chatID int64) error
	ValidateChatID(ctx context.Context, chatID int64) error
}

// CreateValidationService creates a service with validation interface access for testing
func (h *TestHelpers) CreateValidationService() (ValidationService, *MockTelegramRepository) {
	mockRepo := new(MockTelegramRepository)
	logger := logrus.New()
	logger.SetLevel(logrus.FatalLevel)

	service := subscription.NewTelegramSubscriptionService(subscription.Dep{
		TelegramRepo: mockRepo,
		Logger:       logger,
	})

	// Cast to validation service interface
	validationService, ok := service.(ValidationService)
	if !ok {
		panic("Service does not implement ValidationService interface")
	}

	return validationService, mockRepo
}

// GenerateTestIDs creates test value objects IDs for consistent testing
func (h *TestHelpers) GenerateTestIDs(count int) []value_objects.ID {
	ids := make([]value_objects.ID, count)
	for i := 0; i < count; i++ {
		ids[i] = value_objects.NewID()
	}
	return ids
}
