package domain_test

import (
	"testing"
	"time"

	"github.com/dewisartika8/cicd-status-notifier-bot/internal/core/notification/domain"
	"github.com/dewisartika8/cicd-status-notifier-bot/internal/core/shared/domain/value_objects"
	"github.com/stretchr/testify/assert"
)

func TestNewNotificationMetrics(t *testing.T) {
	metrics := domain.NewNotificationMetrics()

	assert.NotNil(t, metrics)
	assert.Equal(t, 0, metrics.DeliveryAttempts())
	assert.Equal(t, 0, metrics.TotalRetries())
	assert.Equal(t, time.Duration(0), metrics.AverageDeliveryTime())
	assert.Nil(t, metrics.LastAttemptAt())
	assert.Nil(t, metrics.FirstAttemptAt())
	assert.Nil(t, metrics.DeliveredAt())
	assert.Nil(t, metrics.FailedAt())
}

func TestNotificationMetrics_RecordAttempt(t *testing.T) {
	metrics := domain.NewNotificationMetrics()

	metrics.RecordAttempt()

	assert.Equal(t, 1, metrics.DeliveryAttempts())
	assert.NotNil(t, metrics.LastAttemptAt())
	assert.NotNil(t, metrics.FirstAttemptAt())
	assert.Equal(t, metrics.FirstAttemptAt(), metrics.LastAttemptAt())
}

func TestNotificationMetrics_RecordRetry(t *testing.T) {
	metrics := domain.NewNotificationMetrics()

	metrics.RecordRetry()

	assert.Equal(t, 1, metrics.TotalRetries())
	assert.Equal(t, 1, metrics.DeliveryAttempts())
	assert.NotNil(t, metrics.LastAttemptAt())
}

func TestNotificationMetrics_RecordDelivery(t *testing.T) {
	metrics := domain.NewNotificationMetrics()

	// Record first attempt
	metrics.RecordAttempt()

	// Small delay
	time.Sleep(1 * time.Millisecond)

	// Record delivery
	metrics.RecordDelivery()

	assert.NotNil(t, metrics.DeliveredAt())
	assert.Greater(t, metrics.AverageDeliveryTime(), time.Duration(0))
}

func TestNotificationMetrics_RecordFailure(t *testing.T) {
	metrics := domain.NewNotificationMetrics()

	metrics.RecordFailure()

	assert.NotNil(t, metrics.FailedAt())
}

func TestNewNotificationLog(t *testing.T) {
	buildEventID := value_objects.NewID()
	projectID := value_objects.NewID()
	channel := domain.NotificationChannelTelegram
	recipient := "123456789"
	message := "Test notification"
	maxRetries := 3

	log, err := domain.NewNotificationLog(
		buildEventID,
		projectID,
		channel,
		recipient,
		message,
		maxRetries,
	)

	assert.NoError(t, err)
	assert.NotNil(t, log)
	assert.Equal(t, buildEventID, log.BuildEventID())
	assert.Equal(t, projectID, log.ProjectID())
	assert.Equal(t, channel, log.Channel())
	assert.Equal(t, recipient, log.Recipient())
	assert.Equal(t, message, log.Message())
	assert.Equal(t, maxRetries, log.MaxRetries())
	assert.Equal(t, domain.NotificationStatusPending, log.Status())
	assert.Equal(t, 0, log.RetryCount())
	assert.NotNil(t, log.Metrics())
	assert.NotNil(t, log.Metadata())
	assert.False(t, log.ID().IsNil())
}

func TestNewNotificationLog_ValidationErrors(t *testing.T) {
	tests := []struct {
		name         string
		buildEventID value_objects.ID
		projectID    value_objects.ID
		channel      domain.NotificationChannel
		recipient    string
		message      string
		maxRetries   int
		wantErr      bool
	}{
		{
			name:         "empty_build_event_id",
			buildEventID: value_objects.ID{},
			projectID:    value_objects.NewID(),
			channel:      domain.NotificationChannelTelegram,
			recipient:    "123456789",
			message:      "Test message",
			maxRetries:   3,
			wantErr:      true,
		},
		{
			name:         "empty_project_id",
			buildEventID: value_objects.NewID(),
			projectID:    value_objects.ID{},
			channel:      domain.NotificationChannelTelegram,
			recipient:    "123456789",
			message:      "Test message",
			maxRetries:   3,
			wantErr:      true,
		},
		{
			name:         "invalid_channel",
			buildEventID: value_objects.NewID(),
			projectID:    value_objects.NewID(),
			channel:      domain.NotificationChannel("invalid"),
			recipient:    "123456789",
			message:      "Test message",
			maxRetries:   3,
			wantErr:      true,
		},
		{
			name:         "empty_recipient",
			buildEventID: value_objects.NewID(),
			projectID:    value_objects.NewID(),
			channel:      domain.NotificationChannelTelegram,
			recipient:    "",
			message:      "Test message",
			maxRetries:   3,
			wantErr:      true,
		},
		{
			name:         "empty_message",
			buildEventID: value_objects.NewID(),
			projectID:    value_objects.NewID(),
			channel:      domain.NotificationChannelTelegram,
			recipient:    "123456789",
			message:      "",
			maxRetries:   3,
			wantErr:      true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			log, err := domain.NewNotificationLog(
				tt.buildEventID,
				tt.projectID,
				tt.channel,
				tt.recipient,
				tt.message,
				tt.maxRetries,
			)

			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, log)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, log)
			}
		})
	}
}

func TestNotificationLog_MarkAsSent(t *testing.T) {
	log := createTestNotificationLog(t)
	messageID := "msg123"

	err := log.MarkAsSent(&messageID)

	assert.NoError(t, err)
	assert.Equal(t, domain.NotificationStatusSent, log.Status())
	assert.Equal(t, &messageID, log.MessageID())
	assert.NotNil(t, log.SentAt())
	assert.Equal(t, "", log.ErrorMessage())
	assert.Equal(t, 1, log.Metrics().DeliveryAttempts())
}

func TestNotificationLog_MarkAsDelivered(t *testing.T) {
	log := createTestNotificationLog(t)

	err := log.MarkAsDelivered()

	assert.NoError(t, err)
	assert.Equal(t, domain.NotificationStatusDelivered, log.Status())
	assert.Equal(t, "", log.ErrorMessage())
	assert.NotNil(t, log.Metrics().DeliveredAt())
}

func TestNotificationLog_MarkAsFailed(t *testing.T) {
	log := createTestNotificationLog(t)
	errorMsg := "Network timeout"

	err := log.MarkAsFailed(errorMsg)

	assert.NoError(t, err)
	assert.Equal(t, domain.NotificationStatusFailed, log.Status())
	assert.Equal(t, errorMsg, log.ErrorMessage())
	assert.NotNil(t, log.Metrics().FailedAt())
}

func TestNotificationLog_MarkAsRetrying(t *testing.T) {
	log := createTestNotificationLog(t)

	err := log.MarkAsRetrying()

	assert.NoError(t, err)
	assert.Equal(t, domain.NotificationStatusRetrying, log.Status())
	assert.Equal(t, 1, log.RetryCount())
	assert.Equal(t, 1, log.Metrics().TotalRetries())
}

func TestNotificationLog_MarkAsRetrying_MaxRetriesExceeded(t *testing.T) {
	log := createTestNotificationLog(t)

	// Exhaust retries
	for i := 0; i < 3; i++ {
		_ = log.MarkAsRetrying()
	}

	// Try one more retry
	err := log.MarkAsRetrying()

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "maximum retry attempts")
}

func TestNotificationLog_MarkAsExpired(t *testing.T) {
	log := createTestNotificationLog(t)

	err := log.MarkAsExpired()

	assert.NoError(t, err)
	assert.Equal(t, domain.NotificationStatusExpired, log.Status())
}

func TestNotificationLog_MarkAsCancelled(t *testing.T) {
	log := createTestNotificationLog(t)

	err := log.MarkAsCancelled()

	assert.NoError(t, err)
	assert.Equal(t, domain.NotificationStatusCancelled, log.Status())
}

func TestNotificationLog_CanRetry(t *testing.T) {
	log := createTestNotificationLog(t)

	// Initially can't retry (status is pending)
	assert.False(t, log.CanRetry())

	// Mark as failed, now can retry
	_ = log.MarkAsFailed("Test error")
	assert.True(t, log.CanRetry())

	// Exhaust retries
	for i := 0; i < 3; i++ {
		_ = log.MarkAsRetrying()
		_ = log.MarkAsFailed("Test error")
	}

	// Now can't retry
	assert.False(t, log.CanRetry())
}

func TestNotificationLog_ScheduleRetry(t *testing.T) {
	log := createTestNotificationLog(t)
	retryAt := value_objects.NewTimestamp()

	err := log.ScheduleRetry(retryAt)

	assert.NoError(t, err)
	assert.Equal(t, &retryAt, log.NextRetryAt())
}

func TestNotificationLog_ClearRetrySchedule(t *testing.T) {
	log := createTestNotificationLog(t)
	retryAt := value_objects.NewTimestamp()
	_ = log.ScheduleRetry(retryAt)

	log.ClearRetrySchedule()

	assert.Nil(t, log.NextRetryAt())
}

func TestNotificationLog_SetExpiration(t *testing.T) {
	log := createTestNotificationLog(t)
	expiresAt := value_objects.NewTimestamp()

	log.SetExpiration(expiresAt)

	assert.Equal(t, &expiresAt, log.ExpiresAt())
}

func TestNotificationLog_IsExpired(t *testing.T) {
	log := createTestNotificationLog(t)

	// Not expired initially
	assert.False(t, log.IsExpired())

	// Set expiration in the past
	pastTime := value_objects.NewTimestampFromTime(time.Now().Add(-1 * time.Hour))
	log.SetExpiration(pastTime)

	assert.True(t, log.IsExpired())
}

func TestNotificationLog_SetTemplateID(t *testing.T) {
	log := createTestNotificationLog(t)
	templateID := value_objects.NewID()

	log.SetTemplateID(templateID)

	assert.Equal(t, &templateID, log.TemplateID())
}

func TestNotificationLog_UpdateMetadata(t *testing.T) {
	log := createTestNotificationLog(t)

	log.UpdateMetadata("key1", "value1")
	log.UpdateMetadata("key2", 42)

	metadata := log.Metadata()
	assert.Equal(t, "value1", metadata["key1"])
	assert.Equal(t, 42, metadata["key2"])

	value, exists := log.GetMetadataValue("key1")
	assert.True(t, exists)
	assert.Equal(t, "value1", value)
}

func TestNotificationLog_RemoveMetadata(t *testing.T) {
	log := createTestNotificationLog(t)
	log.UpdateMetadata("key1", "value1")

	log.RemoveMetadata("key1")

	_, exists := log.GetMetadataValue("key1")
	assert.False(t, exists)
}

func TestRestoreNotificationLog(t *testing.T) {
	id := value_objects.NewID()
	buildEventID := value_objects.NewID()
	projectID := value_objects.NewID()
	now := value_objects.NewTimestamp()
	metrics := domain.NewNotificationMetrics()

	params := domain.RestoreNotificationLogParams{
		ID:           id,
		BuildEventID: buildEventID,
		ProjectID:    projectID,
		Channel:      domain.NotificationChannelTelegram,
		Recipient:    "123456789",
		Message:      "Test message",
		Status:       domain.NotificationStatusSent,
		ErrorMessage: "",
		RetryCount:   1,
		MaxRetries:   3,
		Metrics:      metrics,
		CreatedAt:    now,
		UpdatedAt:    now,
	}

	log := domain.RestoreNotificationLog(params)

	assert.NotNil(t, log)
	assert.Equal(t, id, log.ID())
	assert.Equal(t, buildEventID, log.BuildEventID())
	assert.Equal(t, projectID, log.ProjectID())
	assert.Equal(t, domain.NotificationChannelTelegram, log.Channel())
	assert.Equal(t, domain.NotificationStatusSent, log.Status())
	assert.Equal(t, 1, log.RetryCount())
	assert.Equal(t, 3, log.MaxRetries())
}

func TestRestoreNotificationMetrics(t *testing.T) {
	now := value_objects.NewTimestamp()

	params := domain.RestoreNotificationMetricsParams{
		DeliveryAttempts:    2,
		TotalRetries:        1,
		AverageDeliveryTime: 5 * time.Second,
		LastAttemptAt:       &now,
		FirstAttemptAt:      &now,
		DeliveredAt:         &now,
	}

	metrics := domain.RestoreNotificationMetrics(params)

	assert.NotNil(t, metrics)
	assert.Equal(t, 2, metrics.DeliveryAttempts())
	assert.Equal(t, 1, metrics.TotalRetries())
	assert.Equal(t, 5*time.Second, metrics.AverageDeliveryTime())
	assert.Equal(t, &now, metrics.LastAttemptAt())
	assert.Equal(t, &now, metrics.FirstAttemptAt())
	assert.Equal(t, &now, metrics.DeliveredAt())
}

// Helper function to create a test notification log
func createTestNotificationLog(t *testing.T) *domain.NotificationLog {
	buildEventID := value_objects.NewID()
	projectID := value_objects.NewID()
	channel := domain.NotificationChannelTelegram
	recipient := "123456789"
	message := "Test notification"
	maxRetries := 3

	log, err := domain.NewNotificationLog(
		buildEventID,
		projectID,
		channel,
		recipient,
		message,
		maxRetries,
	)

	assert.NoError(t, err)
	assert.NotNil(t, log)

	return log
}
