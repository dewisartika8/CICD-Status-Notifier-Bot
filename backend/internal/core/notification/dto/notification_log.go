package dto

import (
	"time"

	"github.com/google/uuid"
)

// NotificationStatus represents the status of a notification
type NotificationStatus string

const (
	NotificationStatusSent    NotificationStatus = "sent"
	NotificationStatusFailed  NotificationStatus = "failed"
	NotificationStatusPending NotificationStatus = "pending"
)

// NotificationLog represents a notification log entity
type NotificationLog struct {
	ID           uuid.UUID          `json:"id"`
	BuildEventID uuid.UUID          `json:"build_event_id" validate:"required"`
	ChatID       int64              `json:"chat_id" validate:"required"`
	MessageID    *int               `json:"message_id,omitempty"`
	Status       NotificationStatus `json:"status" validate:"required"`
	ErrorMessage string             `json:"error_message,omitempty"`
	SentAt       *time.Time         `json:"sent_at,omitempty"`
	CreatedAt    time.Time          `json:"created_at"`
}

// NewNotificationLog creates a new notification log entity
func NewNotificationLog(buildEventID uuid.UUID, chatID int64) *NotificationLog {
	return &NotificationLog{
		ID:           uuid.New(),
		BuildEventID: buildEventID,
		ChatID:       chatID,
		Status:       NotificationStatusPending,
		CreatedAt:    time.Now(),
	}
}

// Validate performs basic validation on notification log entity
func (nl *NotificationLog) Validate() error {
	if nl.BuildEventID == uuid.Nil {
		return ErrInvalidBuildEventID
	}
	if nl.ChatID == 0 {
		return ErrInvalidChatID
	}
	if nl.Status == "" {
		return ErrInvalidNotificationStatus
	}
	return nil
}

// MarkAsSent marks the notification as successfully sent
func (nl *NotificationLog) MarkAsSent(messageID int) {
	nl.Status = NotificationStatusSent
	nl.MessageID = &messageID
	now := time.Now()
	nl.SentAt = &now
	nl.ErrorMessage = ""
}

// MarkAsFailed marks the notification as failed with an error message
func (nl *NotificationLog) MarkAsFailed(errorMessage string) {
	nl.Status = NotificationStatusFailed
	nl.ErrorMessage = errorMessage
	nl.MessageID = nil
	nl.SentAt = nil
}

// IsSent returns true if the notification was successfully sent
func (nl *NotificationLog) IsSent() bool {
	return nl.Status == NotificationStatusSent
}

// IsFailed returns true if the notification failed to send
func (nl *NotificationLog) IsFailed() bool {
	return nl.Status == NotificationStatusFailed
}

// IsPending returns true if the notification is pending
func (nl *NotificationLog) IsPending() bool {
	return nl.Status == NotificationStatusPending
}
