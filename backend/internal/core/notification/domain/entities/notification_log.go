package entities

import (
	"github.com/dewisartika8/cicd-status-notifier-bot/internal/core/shared/domain/errors"
	"github.com/dewisartika8/cicd-status-notifier-bot/internal/core/shared/domain/value_objects"
)

// NotificationStatus represents the status of a notification
type NotificationStatus value_objects.Status

const (
	NotificationStatusPending  NotificationStatus = "pending"
	NotificationStatusSent     NotificationStatus = "sent"
	NotificationStatusFailed   NotificationStatus = "failed"
	NotificationStatusRetrying NotificationStatus = "retrying"
)

// NotificationChannel represents the notification channel type
type NotificationChannel value_objects.Status

const (
	NotificationChannelTelegram NotificationChannel = "telegram"
	NotificationChannelEmail    NotificationChannel = "email"
	NotificationChannelSlack    NotificationChannel = "slack"
	NotificationChannelWebhook  NotificationChannel = "webhook"
)

// NotificationLog represents a notification log domain entity
type NotificationLog struct {
	id           value_objects.ID
	buildEventID value_objects.ID
	projectID    value_objects.ID
	channel      NotificationChannel
	recipient    string
	message      string
	status       NotificationStatus
	errorMessage string
	retryCount   int
	sentAt       *value_objects.Timestamp
	createdAt    value_objects.Timestamp
}

// NewNotificationLog creates a new notification log entity
func NewNotificationLog(
	buildEventID, projectID value_objects.ID,
	channel NotificationChannel,
	recipient, message string,
) (*NotificationLog, error) {
	notificationLog := &NotificationLog{
		id:           value_objects.NewID(),
		buildEventID: buildEventID,
		projectID:    projectID,
		channel:      channel,
		recipient:    recipient,
		message:      message,
		status:       NotificationStatusPending,
		retryCount:   0,
		createdAt:    value_objects.NewTimestamp(),
	}

	if err := notificationLog.validate(); err != nil {
		return nil, err
	}

	return notificationLog, nil
}

// RestoreNotificationLog restores a notification log from persistence
func RestoreNotificationLog(params RestoreNotificationLogParams) *NotificationLog {
	return &NotificationLog{
		id:           params.ID,
		buildEventID: params.BuildEventID,
		projectID:    params.ProjectID,
		channel:      params.Channel,
		recipient:    params.Recipient,
		message:      params.Message,
		status:       params.Status,
		errorMessage: params.ErrorMessage,
		retryCount:   params.RetryCount,
		sentAt:       params.SentAt,
		createdAt:    params.CreatedAt,
	}
}

// RestoreNotificationLogParams holds parameters for restoring a notification log
type RestoreNotificationLogParams struct {
	ID           value_objects.ID
	BuildEventID value_objects.ID
	ProjectID    value_objects.ID
	Channel      NotificationChannel
	Recipient    string
	Message      string
	Status       NotificationStatus
	ErrorMessage string
	RetryCount   int
	SentAt       *value_objects.Timestamp
	CreatedAt    value_objects.Timestamp
}

// ID returns the notification log ID
func (nl *NotificationLog) ID() value_objects.ID {
	return nl.id
}

// BuildEventID returns the build event ID
func (nl *NotificationLog) BuildEventID() value_objects.ID {
	return nl.buildEventID
}

// ProjectID returns the project ID
func (nl *NotificationLog) ProjectID() value_objects.ID {
	return nl.projectID
}

// Channel returns the notification channel
func (nl *NotificationLog) Channel() NotificationChannel {
	return nl.channel
}

// Recipient returns the notification recipient
func (nl *NotificationLog) Recipient() string {
	return nl.recipient
}

// Message returns the notification message
func (nl *NotificationLog) Message() string {
	return nl.message
}

// Status returns the notification status
func (nl *NotificationLog) Status() NotificationStatus {
	return nl.status
}

// ErrorMessage returns the error message if any
func (nl *NotificationLog) ErrorMessage() string {
	return nl.errorMessage
}

// RetryCount returns the retry count
func (nl *NotificationLog) RetryCount() int {
	return nl.retryCount
}

// SentAt returns when the notification was sent
func (nl *NotificationLog) SentAt() *value_objects.Timestamp {
	return nl.sentAt
}

// CreatedAt returns when the notification was created
func (nl *NotificationLog) CreatedAt() value_objects.Timestamp {
	return nl.createdAt
}

// MarkAsSent marks the notification as sent
func (nl *NotificationLog) MarkAsSent() {
	nl.status = NotificationStatusSent
	now := value_objects.NewTimestamp()
	nl.sentAt = &now
	nl.errorMessage = ""
}

// MarkAsFailed marks the notification as failed
func (nl *NotificationLog) MarkAsFailed(errorMessage string) {
	nl.status = NotificationStatusFailed
	nl.errorMessage = errorMessage
}

// MarkAsRetrying marks the notification as retrying
func (nl *NotificationLog) MarkAsRetrying() {
	nl.status = NotificationStatusRetrying
	nl.retryCount++
}

// CanRetry checks if the notification can be retried
func (nl *NotificationLog) CanRetry() bool {
	return nl.status == NotificationStatusFailed && nl.retryCount < 3
}

// IsCompleted checks if the notification is completed (sent or failed with no retries)
func (nl *NotificationLog) IsCompleted() bool {
	return nl.status == NotificationStatusSent ||
		(nl.status == NotificationStatusFailed && !nl.CanRetry())
}

// validate performs domain validation
func (nl *NotificationLog) validate() error {
	if nl.buildEventID.IsNil() {
		return errors.NewDomainError("INVALID_BUILD_EVENT_ID", "build event ID cannot be empty")
	}

	if nl.projectID.IsNil() {
		return errors.NewDomainError("INVALID_PROJECT_ID", "project ID cannot be empty")
	}

	if nl.channel == "" {
		return errors.NewDomainError("INVALID_NOTIFICATION_CHANNEL", "notification channel cannot be empty")
	}

	if nl.recipient == "" {
		return errors.NewDomainError("INVALID_RECIPIENT", "recipient cannot be empty")
	}

	if nl.message == "" {
		return errors.NewDomainError("INVALID_MESSAGE", "message cannot be empty")
	}

	return nil
}
