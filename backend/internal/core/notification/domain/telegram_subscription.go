package domain

import (
	"fmt"
	"strconv"

	"github.com/dewisartika8/cicd-status-notifier-bot/internal/core/shared/domain/value_objects"
)

// TelegramSubscription represents a Telegram subscription domain entity
type TelegramSubscription struct {
	id         value_objects.ID
	projectID  value_objects.ID
	chatID     int64
	userID     *int64
	username   string
	eventTypes []string
	isActive   bool
	createdAt  value_objects.Timestamp
	updatedAt  value_objects.Timestamp
}

// NewTelegramSubscription creates a new telegram subscription entity
func NewTelegramSubscription(projectID value_objects.ID, chatID int64) (*TelegramSubscription, error) {
	subscription := &TelegramSubscription{
		id:         value_objects.NewID(),
		projectID:  projectID,
		chatID:     chatID,
		userID:     nil,
		username:   "",
		eventTypes: []string{},
		isActive:   true,
		createdAt:  value_objects.NewTimestamp(),
		updatedAt:  value_objects.NewTimestamp(),
	}

	if err := subscription.validate(); err != nil {
		return nil, err
	}

	return subscription, nil
}

// RestoreTelegramSubscription restores a telegram subscription from persistence
func RestoreTelegramSubscription(params RestoreTelegramSubscriptionParams) *TelegramSubscription {
	return &TelegramSubscription{
		id:         params.ID,
		projectID:  params.ProjectID,
		chatID:     params.ChatID,
		userID:     params.UserID,
		username:   params.Username,
		eventTypes: params.EventTypes,
		isActive:   params.IsActive,
		createdAt:  params.CreatedAt,
		updatedAt:  params.UpdatedAt,
	}
}

// RestoreTelegramSubscriptionParams holds parameters for restoring a telegram subscription
type RestoreTelegramSubscriptionParams struct {
	ID         value_objects.ID
	ProjectID  value_objects.ID
	ChatID     int64
	UserID     *int64
	Username   string
	EventTypes []string
	IsActive   bool
	CreatedAt  value_objects.Timestamp
	UpdatedAt  value_objects.Timestamp
}

// ID returns the subscription ID
func (ts *TelegramSubscription) ID() value_objects.ID {
	return ts.id
}

// ProjectID returns the project ID
func (ts *TelegramSubscription) ProjectID() value_objects.ID {
	return ts.projectID
}

// ChatID returns the Telegram chat ID
func (ts *TelegramSubscription) ChatID() int64 {
	return ts.chatID
}

// UserID returns the Telegram user ID
func (ts *TelegramSubscription) UserID() *int64 {
	return ts.userID
}

// Username returns the Telegram username
func (ts *TelegramSubscription) Username() string {
	return ts.username
}

// EventTypes returns the subscribed event types
func (ts *TelegramSubscription) EventTypes() []string {
	return ts.eventTypes
}

// IsActive returns whether the subscription is active
func (ts *TelegramSubscription) IsActive() bool {
	return ts.isActive
}

// CreatedAt returns the creation timestamp
func (ts *TelegramSubscription) CreatedAt() value_objects.Timestamp {
	return ts.createdAt
}

// UpdatedAt returns the last update timestamp
func (ts *TelegramSubscription) UpdatedAt() value_objects.Timestamp {
	return ts.updatedAt
}

// Activate activates the subscription
func (ts *TelegramSubscription) Activate() {
	if !ts.isActive {
		ts.isActive = true
		ts.updatedAt = value_objects.NewTimestamp()
	}
}

// Deactivate deactivates the subscription
func (ts *TelegramSubscription) Deactivate() {
	if ts.isActive {
		ts.isActive = false
		ts.updatedAt = value_objects.NewTimestamp()
	}
}

// UpdateChatID updates the chat ID (useful for chat migrations)
func (ts *TelegramSubscription) UpdateChatID(newChatID int64) error {
	if newChatID == 0 {
		return ErrInvalidTelegramChatID
	}

	if ts.chatID != newChatID {
		ts.chatID = newChatID
		ts.updatedAt = value_objects.NewTimestamp()
	}

	return nil
}

// GetChatIDString returns the chat ID as a string for notification purposes
func (ts *TelegramSubscription) GetChatIDString() string {
	return strconv.FormatInt(ts.chatID, 10)
}

// validate validates the telegram subscription entity
func (ts *TelegramSubscription) validate() error {
	if ts.projectID.IsNil() {
		return ErrInvalidProjectID
	}

	// Validate chat ID format - should be non-zero integer
	// Telegram chat IDs can be negative (for groups) or positive (for users)
	if ts.chatID == 0 {
		return ErrInvalidChatID
	}

	return nil
}

// String returns a string representation of the subscription
func (ts *TelegramSubscription) String() string {
	status := "inactive"
	if ts.isActive {
		status = "active"
	}

	return fmt.Sprintf("TelegramSubscription{ID: %s, ProjectID: %s, ChatID: %d, Status: %s}",
		ts.id.String(), ts.projectID.String(), ts.chatID, status)
}
