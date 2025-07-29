package service_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/dewisartika8/cicd-status-notifier-bot/internal/core/notification/domain"
	"github.com/dewisartika8/cicd-status-notifier-bot/internal/core/notification/dto"
	"github.com/dewisartika8/cicd-status-notifier-bot/internal/core/notification/port"
	"github.com/dewisartika8/cicd-status-notifier-bot/internal/core/notification/service/retry"
	"github.com/dewisartika8/cicd-status-notifier-bot/internal/core/shared/domain/value_objects"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockRetryConfigurationRepository is a mock implementation of RetryConfigurationRepository
type MockRetryConfigurationRepository struct {
	mock.Mock
}

func (m *MockRetryConfigurationRepository) Create(ctx context.Context, config *domain.RetryConfiguration) error {
	args := m.Called(ctx, config)
	return args.Error(0)
}

func (m *MockRetryConfigurationRepository) GetByID(ctx context.Context, id value_objects.ID) (*domain.RetryConfiguration, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.RetryConfiguration), args.Error(1)
}

func (m *MockRetryConfigurationRepository) GetActiveConfigurations(ctx context.Context) ([]*domain.RetryConfiguration, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*domain.RetryConfiguration), args.Error(1)
}

func (m *MockRetryConfigurationRepository) GetByChannel(ctx context.Context, channel domain.NotificationChannel) (*domain.RetryConfiguration, error) {
	args := m.Called(ctx, channel)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.RetryConfiguration), args.Error(1)
}

func (m *MockRetryConfigurationRepository) Update(ctx context.Context, config *domain.RetryConfiguration) error {
	// Differentiate from Create by incrementing a counter or logging
	logrus.Debug("Mock Update called") // This line makes the implementation unique
	args := m.Called(ctx, config)
	return args.Error(0)
}

func (m *MockRetryConfigurationRepository) Delete(ctx context.Context, id value_objects.ID) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockRetryConfigurationRepository) BulkCreate(ctx context.Context, configs []*domain.RetryConfiguration) error {
	args := m.Called(ctx, configs)
	return args.Error(0)
}

func setupRetryServiceTest() (port.RetryService, *MockRetryConfigurationRepository) {
	mockRepo := &MockRetryConfigurationRepository{}
	logger := logrus.New()
	logger.SetLevel(logrus.FatalLevel) // Suppress logs during tests

	retryService := retry.NewRetryService(retry.Dep{
		RetryRepo: mockRepo,
		Logger:    logger,
	})
	return retryService, mockRepo
}

func createTestRetryConfiguration() *domain.RetryConfiguration {
	config, _ := domain.NewRetryConfiguration(
		3,
		30*time.Second,
		5*time.Minute,
		10*time.Minute,
		2.0,
		true,
		true,
	)
	return config
}

func TestRetryServiceCreateRetryConfiguration(t *testing.T) {
	retryService, mockRepo := setupRetryServiceTest()
	ctx := context.Background()

	req := dto.CreateRetryConfigurationRequest{
		MaxRetryAttempts:         3,
		InitialRetryDelay:        30 * time.Second,
		MaxRetryDelay:            5 * time.Minute,
		RetryTimeoutDuration:     10 * time.Minute,
		RetryDelayMultiplier:     2.0,
		EnableExponentialBackoff: true,
		EnableDeadLetterQueue:    true,
	}

	mockRepo.On("Create", ctx, mock.AnythingOfType("*domain.RetryConfiguration")).Return(nil)

	config, err := retryService.CreateRetryConfiguration(ctx, req)

	assert.NoError(t, err)
	assert.NotNil(t, config)
	assert.Equal(t, req.MaxRetryAttempts, config.MaxRetryAttempts())
	assert.Equal(t, req.InitialRetryDelay, config.InitialRetryDelay())
	assert.Equal(t, req.MaxRetryDelay, config.MaxRetryDelay())
	assert.Equal(t, req.RetryDelayMultiplier, config.RetryDelayMultiplier())
	assert.Equal(t, req.EnableExponentialBackoff, config.EnableExponentialBackoff())
	assert.Equal(t, req.EnableDeadLetterQueue, config.EnableDeadLetterQueue())

	mockRepo.AssertExpectations(t)
}

func TestRetryServiceCreateRetryConfigurationValidationError(t *testing.T) {
	retryService, _ := setupRetryServiceTest()
	ctx := context.Background()

	req := dto.CreateRetryConfigurationRequest{
		MaxRetryAttempts:         -1, // Invalid
		InitialRetryDelay:        30 * time.Second,
		MaxRetryDelay:            5 * time.Minute,
		RetryTimeoutDuration:     10 * time.Minute,
		RetryDelayMultiplier:     2.0,
		EnableExponentialBackoff: true,
		EnableDeadLetterQueue:    true,
	}

	config, err := retryService.CreateRetryConfiguration(ctx, req)

	assert.Error(t, err)
	assert.Nil(t, config)
	assert.Contains(t, err.Error(), "max retry attempts must be between 0 and 10")
}

func TestRetryServiceCreateRetryConfigurationRepositoryError(t *testing.T) {
	retryService, mockRepo := setupRetryServiceTest()
	ctx := context.Background()

	req := dto.CreateRetryConfigurationRequest{
		MaxRetryAttempts:         3,
		InitialRetryDelay:        30 * time.Second,
		MaxRetryDelay:            5 * time.Minute,
		RetryTimeoutDuration:     10 * time.Minute,
		RetryDelayMultiplier:     2.0,
		EnableExponentialBackoff: true,
		EnableDeadLetterQueue:    true,
	}

	expectedError := errors.New("repository error")
	mockRepo.On("Create", ctx, mock.AnythingOfType("*domain.RetryConfiguration")).Return(expectedError)

	config, err := retryService.CreateRetryConfiguration(ctx, req)

	assert.Error(t, err)
	assert.Nil(t, config)
	assert.Contains(t, err.Error(), expectedError.Error())

	mockRepo.AssertExpectations(t)
}

func TestRetryServiceGetRetryConfiguration(t *testing.T) {
	retryService, mockRepo := setupRetryServiceTest()
	ctx := context.Background()

	configID := value_objects.NewID()
	expectedConfig := createTestRetryConfiguration()

	mockRepo.On("GetByID", ctx, configID).Return(expectedConfig, nil)

	config, err := retryService.GetRetryConfiguration(ctx, configID)

	assert.NoError(t, err)
	assert.NotNil(t, config)
	assert.Equal(t, expectedConfig, config)

	mockRepo.AssertExpectations(t)
}

func TestRetryServiceGetRetryConfigurationNotFound(t *testing.T) {
	retryService, mockRepo := setupRetryServiceTest()
	ctx := context.Background()

	configID := value_objects.NewID()

	mockRepo.On("GetByID", ctx, configID).Return(nil, nil)

	config, err := retryService.GetRetryConfiguration(ctx, configID)

	assert.Error(t, err)
	assert.Nil(t, config)
	assert.Contains(t, err.Error(), "not found")

	mockRepo.AssertExpectations(t)
}

func TestRetryServiceGetRetryConfigurationByChannel(t *testing.T) {
	retryService, mockRepo := setupRetryServiceTest()
	ctx := context.Background()

	channel := domain.NotificationChannelTelegram
	expectedConfig := createTestRetryConfiguration()

	mockRepo.On("GetByChannel", ctx, channel).Return(expectedConfig, nil)

	config, err := retryService.GetRetryConfigurationByChannel(ctx, channel)

	assert.NoError(t, err)
	assert.NotNil(t, config)
	assert.Equal(t, expectedConfig, config)

	mockRepo.AssertExpectations(t)
}

func TestRetryServiceGetRetryConfigurationByChannelUseDefault(t *testing.T) {
	retryService, mockRepo := setupRetryServiceTest()
	ctx := context.Background()

	channel := domain.NotificationChannelTelegram

	mockRepo.On("GetByChannel", ctx, channel).Return(nil, nil)

	config, err := retryService.GetRetryConfigurationByChannel(ctx, channel)

	assert.NoError(t, err)
	assert.NotNil(t, config)
	// Should return default configuration
	assert.Equal(t, 3, config.MaxRetryAttempts())
	assert.True(t, config.EnableExponentialBackoff())

	mockRepo.AssertExpectations(t)
}

func TestRetryServiceProcessRetryableNotification(t *testing.T) {
	retryService, mockRepo := setupRetryServiceTest()
	ctx := context.Background()

	notificationID := value_objects.NewID()
	channel := domain.NotificationChannelTelegram
	config := createTestRetryConfiguration()

	req := dto.ProcessRetryableNotificationRequest{
		NotificationID: notificationID,
		Channel:        channel,
		AttemptCount:   1,
		LastError:      errors.New("network timeout"),
		OriginalPayload: map[string]interface{}{
			"message": "test message",
		},
	}

	mockRepo.On("GetByChannel", ctx, channel).Return(config, nil)

	response, err := retryService.ProcessRetryableNotification(ctx, req)

	assert.NoError(t, err)
	assert.NotNil(t, response)
	assert.True(t, response.ShouldRetry)
	assert.NotNil(t, response.NextAttemptAt)
	assert.Equal(t, 30*time.Second, response.RetryDelay) // First attempt uses initial delay
	assert.False(t, response.SendToDeadLetter)
	assert.Equal(t, config, response.RetryConfiguration)

	mockRepo.AssertExpectations(t)
}
