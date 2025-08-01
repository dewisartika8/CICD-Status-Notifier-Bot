package services_test

import (
	"context"
	"testing"
	"time"

	"github.com/dewisartika8/cicd-status-notifier-bot/internal/adapter/repository/memory"
	"github.com/dewisartika8/cicd-status-notifier-bot/internal/core/notification/domain"
	"github.com/dewisartika8/cicd-status-notifier-bot/internal/core/notification/dto"
	"github.com/dewisartika8/cicd-status-notifier-bot/internal/core/notification/service/delivery"
	"github.com/dewisartika8/cicd-status-notifier-bot/internal/core/shared/domain/value_objects"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

const (
	IntTestMessage      = "Integration test message"
	IntTestSubject      = "Integration test"
	IntTestRecipient    = "int-test-user"
	TestMessage         = "Test Message"
	TestSubject         = "Test Subject"
	TestUser            = "test-user"
	TestEmail           = "test@example.com"
	TestTelegramMsgID   = "tg-msg-123"
	TestEmailMsgID      = "email-msg-123"
	TestMsgID           = "msg-123"
	TestGenericSubject  = "Test"
	RateLimitTestEmail  = "rate-limit-test@example.com"
	LowPriorityUser     = "low-priority-user"
	HighPriorityUser    = "high-priority-user"
	LowPriorityMessage  = "Low priority message"
	HighPriorityMessage = "High priority message"
	NotRegisteredError  = "not registered"

	// Test configuration constants
	MaxRetries          = 3
	MaxRateLimit        = 30
	EmailRateLimit      = 10
	ProcessQueueLimit   = 10
	LowPriority         = 1
	HighPriority        = 5
	DefaultAttemptCount = 1
)

// Mock delivery channel for integration testing
type MockIntegrationDeliveryChannel struct {
	mock.Mock
	channelType domain.NotificationChannel
}

func NewMockIntegrationDeliveryChannel(channelType domain.NotificationChannel) *MockIntegrationDeliveryChannel {
	return &MockIntegrationDeliveryChannel{
		channelType: channelType,
	}
}

func (m *MockIntegrationDeliveryChannel) Send(ctx context.Context, recipient, subject, message string) (string, error) {
	args := m.Called(ctx, recipient, subject, message)
	return args.String(0), args.Error(1)
}

func (m *MockIntegrationDeliveryChannel) GetChannelType() domain.NotificationChannel {
	return m.channelType
}

func (m *MockIntegrationDeliveryChannel) IsAvailable(ctx context.Context) bool {
	return true
}

func (m *MockIntegrationDeliveryChannel) GetMaxRetries() int {
	return MaxRetries
}

func (m *MockIntegrationDeliveryChannel) GetRateLimitInfo() (int, time.Duration) {
	return MaxRateLimit, time.Minute
}

// Mock retry service for integration testing
type MockIntegrationRetryService struct{}

func (m *MockIntegrationRetryService) CalculateRetryDelay(ctx context.Context, channel domain.NotificationChannel, attemptNumber int) (time.Duration, error) {
	return time.Second * time.Duration(attemptNumber), nil
}

func (m *MockIntegrationRetryService) ShouldRetryNotification(ctx context.Context, channel domain.NotificationChannel, attemptCount int, lastError error) (bool, error) {
	return attemptCount < MaxRetries, nil
}

// Implement all other required methods with empty implementations
func (m *MockIntegrationRetryService) CreateRetryConfiguration(ctx context.Context, req dto.CreateRetryConfigurationRequest) (*domain.RetryConfiguration, error) {
	return nil, nil
}
func (m *MockIntegrationRetryService) GetRetryConfiguration(ctx context.Context, id value_objects.ID) (*domain.RetryConfiguration, error) {
	return nil, nil
}
func (m *MockIntegrationRetryService) GetRetryConfigurationByChannel(ctx context.Context, channel domain.NotificationChannel) (*domain.RetryConfiguration, error) {
	return nil, nil
}
func (m *MockIntegrationRetryService) UpdateRetryConfiguration(ctx context.Context, id value_objects.ID, req dto.UpdateRetryConfigurationRequest) (*domain.RetryConfiguration, error) {
	return nil, nil
}
func (m *MockIntegrationRetryService) ActivateRetryConfiguration(ctx context.Context, id value_objects.ID) error {
	return nil
}
func (m *MockIntegrationRetryService) DeactivateRetryConfiguration(ctx context.Context, id value_objects.ID) error {
	return nil
}
func (m *MockIntegrationRetryService) DeleteRetryConfiguration(ctx context.Context, id value_objects.ID) error {
	return nil
}
func (m *MockIntegrationRetryService) ListActiveRetryConfigurations(ctx context.Context) ([]*domain.RetryConfiguration, error) {
	return nil, nil
}
func (m *MockIntegrationRetryService) InitializeDefaultRetryConfigurations(ctx context.Context) error {
	return nil
}
func (m *MockIntegrationRetryService) ProcessRetryableNotification(ctx context.Context, req dto.ProcessRetryableNotificationRequest) (*dto.ProcessRetryableNotificationResponse, error) {
	return nil, nil
}

// Integration test for the complete delivery flow
func TestNotificationDeliveryIntegration(t *testing.T) {
	// Arrange - Setup complete system
	queueRepo := memory.NewInMemoryDeliveryQueueRepository()
	rateLimiter := memory.NewInMemoryRateLimiter()
	retryService := &MockIntegrationRetryService{}
	deliveryService := delivery.NewNotificationDeliveryService(queueRepo, rateLimiter, retryService)

	// Setup mock delivery channel
	mockChannel := NewMockIntegrationDeliveryChannel(domain.NotificationChannelTelegram)
	mockChannel.On("Send", mock.Anything, IntTestRecipient, IntTestSubject, IntTestMessage).Return(TestMsgID, nil)

	// Register delivery channel
	err := deliveryService.RegisterDeliveryChannel(mockChannel)
	assert.NoError(t, err)

	ctx := context.Background()

	// Act - Queue a notification
	notification := domain.NewQueuedNotification(
		value_objects.NewID(),
		domain.NotificationChannelTelegram,
		IntTestRecipient,
		IntTestMessage,
		IntTestSubject,
		1,
		3,
	)

	err = deliveryService.QueueNotification(ctx, notification)
	assert.NoError(t, err)

	// Process the queue
	err = deliveryService.ProcessQueue(ctx, ProcessQueueLimit)
	assert.NoError(t, err)

	// Assert - Verify the notification was processed
	stats, err := deliveryService.GetQueueStats(ctx)
	assert.NoError(t, err)
	assert.Contains(t, stats, "pending_count")

	// Verify delivery channel was called
	mockChannel.AssertExpectations(t)
}

// Integration test for rate limiting
func TestRateLimitingIntegration(t *testing.T) {
	// Arrange
	queueRepo := memory.NewInMemoryDeliveryQueueRepository()
	rateLimiter := memory.NewInMemoryRateLimiter()
	retryService := &MockIntegrationRetryService{}
	deliveryService := delivery.NewNotificationDeliveryService(queueRepo, rateLimiter, retryService)

	ctx := context.Background()
	channel := domain.NotificationChannelEmail // Email has lower limit (10)
	recipient := RateLimitTestEmail

	// Act - Make requests up to the limit
	for i := 0; i < EmailRateLimit; i++ {
		allowed, err := deliveryService.CheckRateLimit(ctx, channel, recipient)
		assert.NoError(t, err)
		assert.True(t, allowed, "Request %d should be allowed", i+1)
	}

	// The 11th request should be denied
	allowed, err := deliveryService.CheckRateLimit(ctx, channel, recipient)
	assert.NoError(t, err)
	assert.False(t, allowed, "Request %d should be denied due to rate limit", EmailRateLimit+1)
}

// Integration test for queue processing with priority
func TestQueuePriorityProcessingIntegration(t *testing.T) {
	// Arrange
	queueRepo := memory.NewInMemoryDeliveryQueueRepository()
	rateLimiter := memory.NewInMemoryRateLimiter()
	retryService := &MockIntegrationRetryService{}
	deliveryService := delivery.NewNotificationDeliveryService(queueRepo, rateLimiter, retryService)

	ctx := context.Background()

	// Create notifications with different priorities
	lowPriority := domain.NewQueuedNotification(
		value_objects.NewID(),
		domain.NotificationChannelTelegram,
		LowPriorityUser,
		LowPriorityMessage,
		TestGenericSubject,
		LowPriority,
		MaxRetries,
	)

	highPriority := domain.NewQueuedNotification(
		value_objects.NewID(),
		domain.NotificationChannelTelegram,
		HighPriorityUser,
		HighPriorityMessage,
		TestGenericSubject,
		HighPriority,
		MaxRetries,
	)

	// Queue notifications (low priority first)
	err := deliveryService.QueueNotification(ctx, lowPriority)
	assert.NoError(t, err)

	err = deliveryService.QueueNotification(ctx, highPriority)
	assert.NoError(t, err)

	// Act - Get pending notifications by priority
	pendingNotifications, err := queueRepo.GetPendingByPriority(ctx, ProcessQueueLimit)
	assert.NoError(t, err)

	// Assert - High priority should come first
	assert.Len(t, pendingNotifications, 2)
	assert.Equal(t, highPriority.ID, pendingNotifications[0].ID, "High priority notification should be first")
	assert.Equal(t, lowPriority.ID, pendingNotifications[1].ID, "Low priority notification should be second")
}

// Integration test for delivery channel management
func TestDeliveryChannelManagementIntegration(t *testing.T) {
	// Arrange
	queueRepo := memory.NewInMemoryDeliveryQueueRepository()
	rateLimiter := memory.NewInMemoryRateLimiter()
	retryService := &MockIntegrationRetryService{}
	deliveryService := delivery.NewNotificationDeliveryService(queueRepo, rateLimiter, retryService)

	ctx := context.Background()

	// Setup mock channels
	telegramChannel := NewMockIntegrationDeliveryChannel(domain.NotificationChannelTelegram)
	emailChannel := NewMockIntegrationDeliveryChannel(domain.NotificationChannelEmail)

	// Act - Register channels
	err := deliveryService.RegisterDeliveryChannel(telegramChannel)
	assert.NoError(t, err)

	err = deliveryService.RegisterDeliveryChannel(emailChannel)
	assert.NoError(t, err)

	// Test direct sending through registered channels
	telegramChannel.On("Send", mock.Anything, TestUser, TestSubject, TestMessage).Return(TestTelegramMsgID, nil)
	emailChannel.On("Send", mock.Anything, TestEmail, TestSubject, TestMessage).Return(TestEmailMsgID, nil)

	// Send through Telegram channel
	messageID, err := deliveryService.SendNotification(ctx, domain.NotificationChannelTelegram, TestUser, TestSubject, TestMessage)
	assert.NoError(t, err)
	assert.Equal(t, TestTelegramMsgID, messageID)

	// Send through Email channel
	messageID, err = deliveryService.SendNotification(ctx, domain.NotificationChannelEmail, TestEmail, TestSubject, TestMessage)
	assert.NoError(t, err)
	assert.Equal(t, TestEmailMsgID, messageID)

	// Test unregistering a channel
	err = deliveryService.UnregisterDeliveryChannel(domain.NotificationChannelTelegram)
	assert.NoError(t, err)

	// Sending through unregistered channel should fail
	_, err = deliveryService.SendNotification(ctx, domain.NotificationChannelTelegram, TestUser, TestSubject, TestMessage)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), NotRegisteredError)

	// Verify mock expectations
	telegramChannel.AssertExpectations(t)
	emailChannel.AssertExpectations(t)
}
