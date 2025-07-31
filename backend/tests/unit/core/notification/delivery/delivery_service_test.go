package delivery

import (
	"context"
	"testing"
	"time"

	"github.com/dewisartika8/cicd-status-notifier-bot/internal/adapter/repository/memory"
	"github.com/dewisartika8/cicd-status-notifier-bot/internal/core/notification/domain"
	"github.com/dewisartika8/cicd-status-notifier-bot/internal/core/notification/dto"
	"github.com/dewisartika8/cicd-status-notifier-bot/internal/core/notification/port"
	"github.com/dewisartika8/cicd-status-notifier-bot/internal/core/notification/service/delivery"
	"github.com/dewisartika8/cicd-status-notifier-bot/internal/core/shared/domain/value_objects"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

// Constants for test data
const (
	TestMessage   = "Build completed successfully"
	TestSubject   = "Build Status"
	TestRecipient = "123456789"
	TestEmail     = "test@example.com"
	TestFailedMsg = "Build failed"
	TestDBError   = "database connection failed"
)

// Mock implementations for retry service
type MockRetryService struct {
	mock.Mock
}

func (m *MockRetryService) CalculateRetryDelay(ctx context.Context, channel domain.NotificationChannel, attemptNumber int) (time.Duration, error) {
	args := m.Called(ctx, channel, attemptNumber)
	return args.Get(0).(time.Duration), args.Error(1)
}

func (m *MockRetryService) ShouldRetryNotification(ctx context.Context, channel domain.NotificationChannel, attemptCount int, lastError error) (bool, error) {
	args := m.Called(ctx, channel, attemptCount, lastError)
	return args.Bool(0), args.Error(1)
}

// Add other required methods with empty implementations for now
func (m *MockRetryService) CreateRetryConfiguration(ctx context.Context, req dto.CreateRetryConfigurationRequest) (*domain.RetryConfiguration, error) {
	return nil, nil
}
func (m *MockRetryService) GetRetryConfiguration(ctx context.Context, id value_objects.ID) (*domain.RetryConfiguration, error) {
	return nil, nil
}
func (m *MockRetryService) GetRetryConfigurationByChannel(ctx context.Context, channel domain.NotificationChannel) (*domain.RetryConfiguration, error) {
	return nil, nil
}
func (m *MockRetryService) UpdateRetryConfiguration(ctx context.Context, id value_objects.ID, req dto.UpdateRetryConfigurationRequest) (*domain.RetryConfiguration, error) {
	return nil, nil
}
func (m *MockRetryService) ActivateRetryConfiguration(ctx context.Context, id value_objects.ID) error {
	return nil
}
func (m *MockRetryService) DeactivateRetryConfiguration(ctx context.Context, id value_objects.ID) error {
	return nil
}
func (m *MockRetryService) DeleteRetryConfiguration(ctx context.Context, id value_objects.ID) error {
	return nil
}
func (m *MockRetryService) ListActiveRetryConfigurations(ctx context.Context) ([]*domain.RetryConfiguration, error) {
	return nil, nil
}
func (m *MockRetryService) InitializeDefaultRetryConfigurations(ctx context.Context) error {
	return nil
}
func (m *MockRetryService) ProcessRetryableNotification(ctx context.Context, req dto.ProcessRetryableNotificationRequest) (*dto.ProcessRetryableNotificationResponse, error) {
	return nil, nil
}

type MockDeliveryChannel struct {
	mock.Mock
}

func (m *MockDeliveryChannel) Send(ctx context.Context, recipient, subject, message string) (string, error) {
	args := m.Called(ctx, recipient, subject, message)
	return args.String(0), args.Error(1)
}

func (m *MockDeliveryChannel) GetChannelType() domain.NotificationChannel {
	args := m.Called()
	return args.Get(0).(domain.NotificationChannel)
}

func (m *MockDeliveryChannel) IsAvailable(ctx context.Context) bool {
	args := m.Called(ctx)
	return args.Bool(0)
}

func (m *MockDeliveryChannel) GetMaxRetries() int {
	args := m.Called()
	return args.Int(0)
}

func (m *MockDeliveryChannel) GetRateLimitInfo() (int, time.Duration) {
	args := m.Called()
	return args.Int(0), args.Get(1).(time.Duration)
}

// DeliveryServiceTestSuite defines the test suite
type DeliveryServiceTestSuite struct {
	suite.Suite
	queueRepo    port.DeliveryQueueRepository
	rateLimiter  domain.RateLimiter
	retryService *MockRetryService
	service      port.NotificationDeliveryService
	ctx          context.Context
}

func (suite *DeliveryServiceTestSuite) SetupTest() {
	suite.queueRepo = memory.NewInMemoryDeliveryQueueRepository()
	suite.rateLimiter = memory.NewInMemoryRateLimiter()
	suite.retryService = new(MockRetryService)
	suite.ctx = context.Background()

	// Create the actual service
	suite.service = delivery.NewNotificationDeliveryService(
		suite.queueRepo,
		suite.rateLimiter,
		suite.retryService,
	)
}

func (suite *DeliveryServiceTestSuite) TestQueueNotificationSuccess() {
	// Arrange
	notification := domain.NewQueuedNotification(
		value_objects.NewID(),
		domain.NotificationChannelTelegram,
		TestRecipient,
		TestMessage,
		TestSubject,
		1,
		3,
	)

	// Act
	err := suite.service.QueueNotification(suite.ctx, notification)

	// Assert
	assert.NoError(suite.T(), err)

	// Verify notification was saved
	saved, err := suite.queueRepo.GetByID(suite.ctx, notification.ID)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), notification.ID, saved.ID)
	assert.Equal(suite.T(), domain.DeliveryStatusPending, saved.Status)
}

func (suite *DeliveryServiceTestSuite) TestQueueNotificationValidation() {
	// Test nil notification
	err := suite.service.QueueNotification(suite.ctx, nil)
	assert.Error(suite.T(), err)
	assert.Contains(suite.T(), err.Error(), "notification cannot be nil")

	// Test empty channel
	notification := domain.NewQueuedNotification(
		value_objects.NewID(),
		"", // empty channel
		TestRecipient,
		TestMessage,
		TestSubject,
		1,
		3,
	)
	err = suite.service.QueueNotification(suite.ctx, notification)
	assert.Error(suite.T(), err)
	assert.Contains(suite.T(), err.Error(), "channel cannot be empty")

	// Test empty recipient
	notification.Channel = domain.NotificationChannelTelegram
	notification.Recipient = "" // empty recipient
	err = suite.service.QueueNotification(suite.ctx, notification)
	assert.Error(suite.T(), err)
	assert.Contains(suite.T(), err.Error(), "recipient cannot be empty")

	// Test empty message
	notification.Recipient = TestRecipient
	notification.Message = "" // empty message
	err = suite.service.QueueNotification(suite.ctx, notification)
	assert.Error(suite.T(), err)
	assert.Contains(suite.T(), err.Error(), "message cannot be empty")
}

func (suite *DeliveryServiceTestSuite) TestCheckRateLimitAllowed() {
	// Arrange
	channel := domain.NotificationChannelTelegram
	recipient := TestRecipient

	// Act
	allowed, err := suite.service.CheckRateLimit(suite.ctx, channel, recipient)

	// Assert
	assert.NoError(suite.T(), err)
	assert.True(suite.T(), allowed)
}

func (suite *DeliveryServiceTestSuite) TestGetQueueStats() {
	// Arrange - add some notifications
	notification1 := domain.NewQueuedNotification(
		value_objects.NewID(),
		domain.NotificationChannelTelegram,
		TestRecipient,
		TestMessage,
		TestSubject,
		1,
		3,
	)
	notification2 := domain.NewQueuedNotification(
		value_objects.NewID(),
		domain.NotificationChannelEmail,
		TestEmail,
		TestFailedMsg,
		TestSubject,
		2,
		3,
	)

	suite.service.QueueNotification(suite.ctx, notification1)
	suite.service.QueueNotification(suite.ctx, notification2)

	// Act
	stats, err := suite.service.GetQueueStats(suite.ctx)

	// Assert
	assert.NoError(suite.T(), err)
	assert.Contains(suite.T(), stats, "pending_count")
	assert.Equal(suite.T(), int64(2), stats["pending_count"])
}

func (suite *DeliveryServiceTestSuite) TestRegisterDeliveryChannel() {
	// Arrange
	mockChannel := new(MockDeliveryChannel)
	mockChannel.On("GetChannelType").Return(domain.NotificationChannelTelegram)

	// Act
	err := suite.service.RegisterDeliveryChannel(mockChannel)

	// Assert
	assert.NoError(suite.T(), err)
	mockChannel.AssertExpectations(suite.T())
}

func (suite *DeliveryServiceTestSuite) TestRegisterDeliveryChannelValidation() {
	// Test nil channel
	err := suite.service.RegisterDeliveryChannel(nil)
	assert.Error(suite.T(), err)
	assert.Contains(suite.T(), err.Error(), "delivery channel cannot be nil")

	// Test empty channel type
	mockChannel := new(MockDeliveryChannel)
	mockChannel.On("GetChannelType").Return(domain.NotificationChannel(""))

	err = suite.service.RegisterDeliveryChannel(mockChannel)
	assert.Error(suite.T(), err)
	assert.Contains(suite.T(), err.Error(), "channel type cannot be empty")
}

func TestDeliveryServiceTestSuite(t *testing.T) {
	suite.Run(t, new(DeliveryServiceTestSuite))
}

// Test domain model methods directly
func TestQueuedNotificationIsRetryable(t *testing.T) {
	notification := domain.NewQueuedNotification(
		value_objects.NewID(),
		domain.NotificationChannelTelegram,
		TestRecipient,
		"Test message",
		"Test subject",
		1,
		3,
	)

	// Should not be retryable when status is pending
	assert.False(t, notification.IsRetryable())

	// Should be retryable when status is failed and attempts < max
	notification.MarkAsFailed("Test error")
	assert.True(t, notification.IsRetryable())

	// Should not be retryable when max attempts reached
	notification.AttemptCount = 3
	assert.False(t, notification.IsRetryable())
}

func TestQueuedNotificationShouldBeProcessed(t *testing.T) {
	notification := domain.NewQueuedNotification(
		value_objects.NewID(),
		domain.NotificationChannelTelegram,
		TestRecipient,
		"Test message",
		"Test subject",
		1,
		3,
	)

	// Should be processed immediately
	assert.True(t, notification.ShouldBeProcessed())

	// Should not be processed when scheduled for future
	notification.ScheduleRetry(time.Hour)
	assert.False(t, notification.ShouldBeProcessed())

	// Should not be processed when not pending
	notification.Status = domain.DeliveryStatusPending // Reset to pending first
	notification.ScheduledAt = time.Now()              // Reset scheduled time
	notification.MarkAsProcessing()
	assert.False(t, notification.ShouldBeProcessed())
}

func TestDefaultRateLimitRules(t *testing.T) {
	rules := domain.DefaultRateLimitRules()

	assert.Contains(t, rules, domain.NotificationChannelTelegram)
	assert.Contains(t, rules, domain.NotificationChannelEmail)
	assert.Contains(t, rules, domain.NotificationChannelSlack)
	assert.Contains(t, rules, domain.NotificationChannelWebhook)

	telegramRule := rules[domain.NotificationChannelTelegram]
	assert.Equal(t, 30, telegramRule.MaxRequests)
	assert.Equal(t, time.Minute, telegramRule.WindowSize)
	assert.Equal(t, 5, telegramRule.BurstLimit)
}
