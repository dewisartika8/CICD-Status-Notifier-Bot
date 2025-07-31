package delivery

import (
	"context"
	"fmt"
	"sync"

	"github.com/dewisartika8/cicd-status-notifier-bot/internal/core/notification/domain"
	"github.com/dewisartika8/cicd-status-notifier-bot/internal/core/notification/port"
)

// notificationDeliveryService implements the NotificationDeliveryService interface
type notificationDeliveryService struct {
	queueRepo     port.DeliveryQueueRepository
	rateLimiter   domain.RateLimiter
	retryService  port.RetryService
	channels      map[domain.NotificationChannel]port.DeliveryChannel
	channelsMutex sync.RWMutex
}

// NewNotificationDeliveryService creates a new notification delivery service
func NewNotificationDeliveryService(
	queueRepo port.DeliveryQueueRepository,
	rateLimiter domain.RateLimiter,
	retryService port.RetryService,
) port.NotificationDeliveryService {
	return &notificationDeliveryService{
		queueRepo:    queueRepo,
		rateLimiter:  rateLimiter,
		retryService: retryService,
		channels:     make(map[domain.NotificationChannel]port.DeliveryChannel),
	}
}

// QueueNotification adds a notification to the delivery queue
func (s *notificationDeliveryService) QueueNotification(ctx context.Context, notification *domain.QueuedNotification) error {
	if notification == nil {
		return fmt.Errorf("notification cannot be nil")
	}

	// Validate notification data
	if notification.Channel == "" {
		return fmt.Errorf("notification channel cannot be empty")
	}

	if notification.Recipient == "" {
		return fmt.Errorf("notification recipient cannot be empty")
	}

	if notification.Message == "" {
		return fmt.Errorf("notification message cannot be empty")
	}

	// Save to repository
	return s.queueRepo.Create(ctx, notification)
}

// ProcessQueue processes pending notifications in the queue
func (s *notificationDeliveryService) ProcessQueue(ctx context.Context, batchSize int) error {
	// Get pending notifications by priority
	notifications, err := s.queueRepo.GetPendingByPriority(ctx, batchSize)
	if err != nil {
		return fmt.Errorf("failed to get pending notifications: %w", err)
	}

	if len(notifications) == 0 {
		return nil // No notifications to process
	}

	// Process each notification
	for _, notification := range notifications {
		if err := s.processNotification(ctx, notification); err != nil {
			// Log error but continue processing other notifications
			continue
		}
	}

	return nil
}

// ProcessRetryQueue processes failed notifications for retry
func (s *notificationDeliveryService) ProcessRetryQueue(ctx context.Context, batchSize int) error {
	// Get failed notifications that can be retried
	failedNotifications, err := s.queueRepo.GetFailedNotifications(ctx, batchSize)
	if err != nil {
		return fmt.Errorf("failed to get failed notifications: %w", err)
	}

	if len(failedNotifications) == 0 {
		return nil // No notifications to retry
	}

	// Process each failed notification
	for _, notification := range failedNotifications {
		// Check if notification should be retried
		shouldRetry, err := s.retryService.ShouldRetryNotification(ctx, notification.Channel, notification.AttemptCount, fmt.Errorf(notification.LastError))
		if err != nil {
			continue
		}

		if !shouldRetry {
			// Mark as expired or cancelled
			notification.Status = domain.DeliveryStatusExpired
			s.queueRepo.Update(ctx, notification)
			continue
		}

		// Calculate retry delay
		delay, err := s.retryService.CalculateRetryDelay(ctx, notification.Channel, notification.AttemptCount)
		if err != nil {
			continue
		}

		// Schedule for retry
		notification.ScheduleRetry(delay)
		s.queueRepo.Update(ctx, notification)
	}

	return nil
}

// GetQueueStats returns queue statistics
func (s *notificationDeliveryService) GetQueueStats(ctx context.Context) (map[string]interface{}, error) {
	stats := make(map[string]interface{})

	// Get pending count
	pendingCount, err := s.queueRepo.GetPendingCount(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get pending count: %w", err)
	}
	stats["pending_count"] = pendingCount

	// Get queue stats by status
	queueStats, err := s.queueRepo.GetQueueStats(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get queue stats: %w", err)
	}

	for status, count := range queueStats {
		stats[status] = count
	}

	return stats, nil
}

// RegisterDeliveryChannel registers a delivery channel
func (s *notificationDeliveryService) RegisterDeliveryChannel(channel port.DeliveryChannel) error {
	if channel == nil {
		return fmt.Errorf("delivery channel cannot be nil")
	}

	channelType := channel.GetChannelType()
	if channelType == "" {
		return fmt.Errorf("delivery channel type cannot be empty")
	}

	s.channelsMutex.Lock()
	defer s.channelsMutex.Unlock()

	s.channels[channelType] = channel
	return nil
}

// UnregisterDeliveryChannel unregisters a delivery channel
func (s *notificationDeliveryService) UnregisterDeliveryChannel(channelType domain.NotificationChannel) error {
	s.channelsMutex.Lock()
	defer s.channelsMutex.Unlock()

	delete(s.channels, channelType)
	return nil
}

// SendNotification sends a notification immediately (bypassing queue)
func (s *notificationDeliveryService) SendNotification(ctx context.Context, channel domain.NotificationChannel, recipient, subject, message string) (string, error) {
	// Check rate limit first
	allowed, err := s.CheckRateLimit(ctx, channel, recipient)
	if err != nil {
		return "", fmt.Errorf("failed to check rate limit: %w", err)
	}

	if !allowed {
		return "", fmt.Errorf("rate limit exceeded for channel %s and recipient %s", channel, recipient)
	}

	// Get delivery channel
	s.channelsMutex.RLock()
	deliveryChannel, exists := s.channels[channel]
	s.channelsMutex.RUnlock()

	if !exists {
		return "", fmt.Errorf("delivery channel %s not registered", channel)
	}

	// Check if channel is available
	if !deliveryChannel.IsAvailable(ctx) {
		return "", fmt.Errorf("delivery channel %s is not available", channel)
	}

	// Send notification
	messageID, err := deliveryChannel.Send(ctx, recipient, subject, message)
	if err != nil {
		return "", fmt.Errorf("failed to send notification via %s: %w", channel, err)
	}

	return messageID, nil
}

// CheckRateLimit checks if a notification can be sent based on rate limiting
func (s *notificationDeliveryService) CheckRateLimit(ctx context.Context, channel domain.NotificationChannel, recipient string) (bool, error) {
	return s.rateLimiter.Allow(ctx, recipient, channel)
}

// processNotification processes a single notification
func (s *notificationDeliveryService) processNotification(ctx context.Context, notification *domain.QueuedNotification) error {
	// Mark as processing
	if err := s.queueRepo.UpdateStatus(ctx, notification.ID, domain.DeliveryStatusProcessing, ""); err != nil {
		return fmt.Errorf("failed to mark notification as processing: %w", err)
	}

	// Check rate limit
	allowed, err := s.CheckRateLimit(ctx, notification.Channel, notification.Recipient)
	if err != nil {
		s.queueRepo.UpdateStatus(ctx, notification.ID, domain.DeliveryStatusFailed, fmt.Sprintf("rate limit check failed: %v", err))
		return err
	}

	if !allowed {
		// Schedule for retry later
		delay, _ := s.retryService.CalculateRetryDelay(ctx, notification.Channel, notification.AttemptCount)
		notification.ScheduleRetry(delay)
		s.queueRepo.Update(ctx, notification)
		return fmt.Errorf("rate limit exceeded")
	}

	// Send notification
	messageID, err := s.SendNotification(ctx, notification.Channel, notification.Recipient, notification.Subject, notification.Message)
	if err != nil {
		// Mark as failed
		s.queueRepo.UpdateStatus(ctx, notification.ID, domain.DeliveryStatusFailed, err.Error())
		return err
	}

	// Mark as delivered
	if err := s.queueRepo.UpdateStatus(ctx, notification.ID, domain.DeliveryStatusDelivered, ""); err != nil {
		return fmt.Errorf("failed to mark notification as delivered: %w", err)
	}

	// Store message ID if needed (could be added to domain model later)
	_ = messageID

	return nil
}
