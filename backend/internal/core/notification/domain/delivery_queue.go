package domain

import (
	"context"
	"time"

	"github.com/dewisartika8/cicd-status-notifier-bot/internal/core/shared/domain/value_objects"
)

// DeliveryStatus represents the status of a notification delivery
type DeliveryStatus value_objects.Status

const (
	DeliveryStatusPending    DeliveryStatus = "pending"
	DeliveryStatusProcessing DeliveryStatus = "processing"
	DeliveryStatusDelivered  DeliveryStatus = "delivered"
	DeliveryStatusFailed     DeliveryStatus = "failed"
	DeliveryStatusRetrying   DeliveryStatus = "retrying"
	DeliveryStatusCancelled  DeliveryStatus = "cancelled"
	DeliveryStatusExpired    DeliveryStatus = "expired"
)

// IsValid checks if the delivery status is valid
func (s DeliveryStatus) IsValid() bool {
	switch s {
	case DeliveryStatusPending, DeliveryStatusProcessing, DeliveryStatusDelivered,
		DeliveryStatusFailed, DeliveryStatusRetrying, DeliveryStatusCancelled,
		DeliveryStatusExpired:
		return true
	default:
		return false
	}
}

// String returns the string representation of delivery status
func (s DeliveryStatus) String() string {
	return string(s)
}

// QueuedNotification represents a notification in the delivery queue
type QueuedNotification struct {
	ID             value_objects.ID    `json:"id"`
	NotificationID value_objects.ID    `json:"notification_id"`
	Channel        NotificationChannel `json:"channel"`
	Recipient      string              `json:"recipient"`
	Message        string              `json:"message"`
	Subject        string              `json:"subject"`
	Priority       int                 `json:"priority"`
	ScheduledAt    time.Time           `json:"scheduled_at"`
	AttemptCount   int                 `json:"attempt_count"`
	MaxAttempts    int                 `json:"max_attempts"`
	Status         DeliveryStatus      `json:"status"`
	LastError      string              `json:"last_error"`
	CreatedAt      time.Time           `json:"created_at"`
	UpdatedAt      time.Time           `json:"updated_at"`
}

// NewQueuedNotification creates a new queued notification
func NewQueuedNotification(
	notificationID value_objects.ID,
	channel NotificationChannel,
	recipient, message, subject string,
	priority int,
	maxAttempts int,
) *QueuedNotification {
	now := time.Now()
	return &QueuedNotification{
		ID:             value_objects.NewID(),
		NotificationID: notificationID,
		Channel:        channel,
		Recipient:      recipient,
		Message:        message,
		Subject:        subject,
		Priority:       priority,
		ScheduledAt:    now,
		AttemptCount:   0,
		MaxAttempts:    maxAttempts,
		Status:         DeliveryStatusPending,
		CreatedAt:      now,
		UpdatedAt:      now,
	}
}

// IsRetryable checks if the notification can be retried
func (qn *QueuedNotification) IsRetryable() bool {
	return qn.AttemptCount < qn.MaxAttempts &&
		(qn.Status == DeliveryStatusFailed || qn.Status == DeliveryStatusRetrying)
}

// ShouldBeProcessed checks if the notification should be processed now
func (qn *QueuedNotification) ShouldBeProcessed() bool {
	return qn.Status == DeliveryStatusPending && !time.Now().Before(qn.ScheduledAt)
}

// MarkAsProcessing marks the notification as being processed
func (qn *QueuedNotification) MarkAsProcessing() {
	qn.Status = DeliveryStatusProcessing
	qn.UpdatedAt = time.Now()
}

// MarkAsDelivered marks the notification as successfully delivered
func (qn *QueuedNotification) MarkAsDelivered() {
	qn.Status = DeliveryStatusDelivered
	qn.UpdatedAt = time.Now()
}

// MarkAsFailed marks the notification as failed and increments attempt count
func (qn *QueuedNotification) MarkAsFailed(errorMessage string) {
	qn.Status = DeliveryStatusFailed
	qn.LastError = errorMessage
	qn.AttemptCount++
	qn.UpdatedAt = time.Now()
}

// ScheduleRetry schedules the notification for retry with delay
func (qn *QueuedNotification) ScheduleRetry(delay time.Duration) {
	qn.Status = DeliveryStatusRetrying
	qn.ScheduledAt = time.Now().Add(delay)
	qn.UpdatedAt = time.Now()
}

// NotificationQueue defines the interface for notification queue operations
type NotificationQueue interface {
	// Enqueue adds a notification to the queue
	Enqueue(ctx context.Context, notification *QueuedNotification) error

	// Dequeue retrieves the next notification from the queue
	Dequeue(ctx context.Context) (*QueuedNotification, error)

	// DequeueByPriority retrieves notifications by priority order
	DequeueByPriority(ctx context.Context, limit int) ([]*QueuedNotification, error)

	// UpdateStatus updates the status of a queued notification
	UpdateStatus(ctx context.Context, id value_objects.ID, status DeliveryStatus, errorMessage string) error

	// GetPendingCount returns the number of pending notifications
	GetPendingCount(ctx context.Context) (int64, error)

	// GetQueuedNotification retrieves a queued notification by ID
	GetQueuedNotification(ctx context.Context, id value_objects.ID) (*QueuedNotification, error)

	// DeleteProcessedNotifications removes successfully delivered notifications older than specified duration
	DeleteProcessedNotifications(ctx context.Context, olderThan time.Duration) error

	// GetFailedNotifications retrieves failed notifications for retry
	GetFailedNotifications(ctx context.Context, limit int) ([]*QueuedNotification, error)
}
