package entities

import (
	"time"

	"github.com/google/uuid"
)

// TelegramSubscription represents a telegram subscription entity
type TelegramSubscription struct {
	ID         uuid.UUID   `json:"id"`
	ProjectID  uuid.UUID   `json:"project_id" validate:"required"`
	ChatID     int64       `json:"chat_id" validate:"required"`
	UserID     *int64      `json:"user_id,omitempty"`
	Username   string      `json:"username,omitempty"`
	EventTypes []EventType `json:"event_types"`
	IsActive   bool        `json:"is_active"`
	CreatedAt  time.Time   `json:"created_at"`
}

// NewTelegramSubscription creates a new telegram subscription entity
func NewTelegramSubscription(projectID uuid.UUID, chatID int64, userID *int64, username string) *TelegramSubscription {
	return &TelegramSubscription{
		ID:        uuid.New(),
		ProjectID: projectID,
		ChatID:    chatID,
		UserID:    userID,
		Username:  username,
		EventTypes: []EventType{
			EventTypeBuildSuccess,
			EventTypeBuildFailed,
			EventTypeDeploymentSuccess,
			EventTypeDeploymentFailed,
		}, // Default event types
		IsActive:  true,
		CreatedAt: time.Now(),
	}
}

// Validate performs basic validation on telegram subscription entity
func (ts *TelegramSubscription) Validate() error {
	if ts.ProjectID == uuid.Nil {
		return ErrInvalidProjectID
	}
	if ts.ChatID == 0 {
		return ErrInvalidChatID
	}
	return nil
}

// UpdateEventTypes updates the event types for the subscription
func (ts *TelegramSubscription) UpdateEventTypes(eventTypes []EventType) {
	ts.EventTypes = eventTypes
}

// Subscribe enables the subscription
func (ts *TelegramSubscription) Subscribe() {
	ts.IsActive = true
}

// Unsubscribe disables the subscription
func (ts *TelegramSubscription) Unsubscribe() {
	ts.IsActive = false
}

// IsSubscribedTo checks if the subscription is active for a specific event type
func (ts *TelegramSubscription) IsSubscribedTo(eventType EventType) bool {
	if !ts.IsActive {
		return false
	}

	for _, et := range ts.EventTypes {
		if et == eventType {
			return true
		}
	}
	return false
}
