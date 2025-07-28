package entities

import (
	"github.com/dewisartika8/cicd-status-notifier-bot/internal/core/build/domain/entities"
	"github.com/dewisartika8/cicd-status-notifier-bot/internal/core/shared/domain/errors"
	"github.com/dewisartika8/cicd-status-notifier-bot/internal/core/shared/domain/value_objects"
)

// TelegramSubscription represents a telegram subscription domain entity
type TelegramSubscription struct {
	id         value_objects.ID
	projectID  value_objects.ID
	chatID     int64
	userID     *int64
	username   string
	eventTypes []entities.EventType
	isActive   bool
	createdAt  value_objects.Timestamp
}

// NewTelegramSubscription creates a new telegram subscription entity
func NewTelegramSubscription(
	projectID value_objects.ID,
	chatID int64,
	userID *int64,
	username string,
	eventTypes []entities.EventType,
) (*TelegramSubscription, error) {
	if len(eventTypes) == 0 {
		// Default event types
		eventTypes = []entities.EventType{
			entities.EventTypeBuildCompleted,
			entities.EventTypeDeploymentCompleted,
		}
	}

	subscription := &TelegramSubscription{
		id:         value_objects.NewID(),
		projectID:  projectID,
		chatID:     chatID,
		userID:     userID,
		username:   username,
		eventTypes: eventTypes,
		isActive:   true,
		createdAt:  value_objects.NewTimestamp(),
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
	}
}

// RestoreTelegramSubscriptionParams holds parameters for restoring a telegram subscription
type RestoreTelegramSubscriptionParams struct {
	ID         value_objects.ID
	ProjectID  value_objects.ID
	ChatID     int64
	UserID     *int64
	Username   string
	EventTypes []entities.EventType
	IsActive   bool
	CreatedAt  value_objects.Timestamp
}

// ID returns the subscription ID
func (ts *TelegramSubscription) ID() value_objects.ID {
	return ts.id
}

// ProjectID returns the project ID
func (ts *TelegramSubscription) ProjectID() value_objects.ID {
	return ts.projectID
}

// ChatID returns the telegram chat ID
func (ts *TelegramSubscription) ChatID() int64 {
	return ts.chatID
}

// UserID returns the telegram user ID
func (ts *TelegramSubscription) UserID() *int64 {
	return ts.userID
}

// Username returns the telegram username
func (ts *TelegramSubscription) Username() string {
	return ts.username
}

// EventTypes returns the subscribed event types
func (ts *TelegramSubscription) EventTypes() []entities.EventType {
	return ts.eventTypes
}

// IsActive returns if the subscription is active
func (ts *TelegramSubscription) IsActive() bool {
	return ts.isActive
}

// CreatedAt returns when the subscription was created
func (ts *TelegramSubscription) CreatedAt() value_objects.Timestamp {
	return ts.createdAt
}

// UpdateEventTypes updates the subscribed event types
func (ts *TelegramSubscription) UpdateEventTypes(eventTypes []entities.EventType) error {
	if len(eventTypes) == 0 {
		return errors.NewDomainError("INVALID_EVENT_TYPES", "at least one event type must be specified")
	}
	ts.eventTypes = eventTypes
	return nil
}

// Activate activates the subscription
func (ts *TelegramSubscription) Activate() {
	ts.isActive = true
}

// Deactivate deactivates the subscription
func (ts *TelegramSubscription) Deactivate() {
	ts.isActive = false
}

// IsSubscribedTo checks if the subscription is subscribed to a specific event type
func (ts *TelegramSubscription) IsSubscribedTo(eventType entities.EventType) bool {
	for _, et := range ts.eventTypes {
		if et == eventType {
			return true
		}
	}
	return false
}

// validate performs domain validation
func (ts *TelegramSubscription) validate() error {
	if ts.projectID.IsNil() {
		return errors.NewDomainError("INVALID_PROJECT_ID", "project ID cannot be empty")
	}

	if ts.chatID == 0 {
		return errors.NewDomainError("INVALID_CHAT_ID", "chat ID cannot be zero")
	}

	if len(ts.eventTypes) == 0 {
		return errors.NewDomainError("INVALID_EVENT_TYPES", "at least one event type must be specified")
	}

	return nil
}
