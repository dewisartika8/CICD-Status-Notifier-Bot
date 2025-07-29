package domain

import (
	"strings"

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

// IsValid checks if the notification status is valid
func (s NotificationStatus) IsValid() bool {
	switch s {
	case NotificationStatusPending, NotificationStatusSent, NotificationStatusFailed, NotificationStatusRetrying:
		return true
	default:
		return false
	}
}

// NotificationChannel represents the notification channel type
type NotificationChannel value_objects.Status

const (
	NotificationChannelTelegram NotificationChannel = "telegram"
	NotificationChannelEmail    NotificationChannel = "email"
	NotificationChannelSlack    NotificationChannel = "slack"
	NotificationChannelWebhook  NotificationChannel = "webhook"
)

// IsValid checks if the notification channel is valid
func (c NotificationChannel) IsValid() bool {
	switch c {
	case NotificationChannelTelegram, NotificationChannelEmail, NotificationChannelSlack, NotificationChannelWebhook:
		return true
	default:
		return false
	}
}

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
	messageID    *string // For storing external message ID (e.g., Telegram message ID)
	sentAt       *value_objects.Timestamp
	createdAt    value_objects.Timestamp
	updatedAt    value_objects.Timestamp
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
		recipient:    strings.TrimSpace(recipient),
		message:      strings.TrimSpace(message),
		status:       NotificationStatusPending,
		retryCount:   0,
		createdAt:    value_objects.NewTimestamp(),
		updatedAt:    value_objects.NewTimestamp(),
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
		messageID:    params.MessageID,
		sentAt:       params.SentAt,
		createdAt:    params.CreatedAt,
		updatedAt:    params.UpdatedAt,
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
	MessageID    *string
	SentAt       *value_objects.Timestamp
	CreatedAt    value_objects.Timestamp
	UpdatedAt    value_objects.Timestamp
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

// ErrorMessage returns the error message
func (nl *NotificationLog) ErrorMessage() string {
	return nl.errorMessage
}

// RetryCount returns the retry count
func (nl *NotificationLog) RetryCount() int {
	return nl.retryCount
}

// MessageID returns the external message ID
func (nl *NotificationLog) MessageID() *string {
	return nl.messageID
}

// SentAt returns the sent timestamp
func (nl *NotificationLog) SentAt() *value_objects.Timestamp {
	return nl.sentAt
}

// CreatedAt returns the creation timestamp
func (nl *NotificationLog) CreatedAt() value_objects.Timestamp {
	return nl.createdAt
}

// UpdatedAt returns the last update timestamp
func (nl *NotificationLog) UpdatedAt() value_objects.Timestamp {
	return nl.updatedAt
}

// MarkAsSent marks the notification as successfully sent
func (nl *NotificationLog) MarkAsSent(messageID *string) error {
	if nl.status == NotificationStatusSent {
		return nil // Already sent, no-op
	}

	nl.status = NotificationStatusSent
	nl.messageID = messageID
	ts := value_objects.NewTimestamp()
	nl.sentAt = &ts
	nl.updatedAt = value_objects.NewTimestamp()
	nl.errorMessage = "" // Clear any previous error

	return nil
}

// MarkAsFailed marks the notification as failed with an error message
func (nl *NotificationLog) MarkAsFailed(errorMessage string) error {
	if strings.TrimSpace(errorMessage) == "" {
		return ErrInvalidMessage
	}

	nl.status = NotificationStatusFailed
	nl.errorMessage = strings.TrimSpace(errorMessage)
	nl.updatedAt = value_objects.NewTimestamp()

	return nil
}

// MarkAsRetrying marks the notification for retry
func (nl *NotificationLog) MarkAsRetrying() error {
	const maxRetries = 3

	if nl.retryCount >= maxRetries {
		return ErrMaxRetryAttemptsExceeded
	}

	nl.status = NotificationStatusRetrying
	nl.retryCount++
	nl.updatedAt = value_objects.NewTimestamp()

	return nil
}

// CanRetry checks if the notification can be retried
func (nl *NotificationLog) CanRetry() bool {
	const maxRetries = 3
	return nl.status == NotificationStatusFailed && nl.retryCount < maxRetries
}

// UpdateMessage updates the notification message
func (nl *NotificationLog) UpdateMessage(message string) error {
	message = strings.TrimSpace(message)
	if message == "" {
		return ErrInvalidMessage
	}

	nl.message = message
	nl.updatedAt = value_objects.NewTimestamp()

	return nil
}

// validate validates the notification log entity
func (nl *NotificationLog) validate() error {
	if nl.buildEventID.IsNil() {
		return ErrInvalidNotificationLog
	}

	if nl.projectID.IsNil() {
		return ErrInvalidNotificationLog
	}

	if !nl.channel.IsValid() {
		return ErrInvalidNotificationChannel
	}

	if strings.TrimSpace(nl.recipient) == "" {
		return ErrInvalidRecipient
	}

	if strings.TrimSpace(nl.message) == "" {
		return ErrInvalidMessage
	}

	if !nl.status.IsValid() {
		return ErrInvalidNotificationStatus
	}

	// Channel-specific validation
	if err := nl.validateChannelSpecific(); err != nil {
		return err
	}

	return nil
}

// validateChannelSpecific performs channel-specific validation
func (nl *NotificationLog) validateChannelSpecific() error {
	switch nl.channel {
	case NotificationChannelTelegram:
		// Telegram recipient should be a chat ID (numeric string) or username
		if nl.recipient == "" {
			return NewInvalidRecipientError("telegram chat ID cannot be empty")
		}
	case NotificationChannelEmail:
		// Basic email validation (could be enhanced)
		if !strings.Contains(nl.recipient, "@") {
			return NewInvalidRecipientError("invalid email format")
		}
	case NotificationChannelSlack:
		// Slack channel should start with # or be a user ID
		if nl.recipient == "" {
			return NewInvalidRecipientError("slack channel cannot be empty")
		}
	case NotificationChannelWebhook:
		// Webhook URL validation (basic)
		if !strings.HasPrefix(nl.recipient, "http") {
			return NewInvalidRecipientError("webhook URL must start with http")
		}
	}

	return nil
}
