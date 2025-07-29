package events

import (
	"github.com/dewisartika8/cicd-status-notifier-bot/internal/core/shared/domain/value_objects"
)

// DomainEvent represents a domain event
type DomainEvent interface {
	EventID() value_objects.ID
	EventType() string
	AggregateID() value_objects.ID
	OccurredOn() value_objects.Timestamp
	Version() int
}

// BaseDomainEvent provides common implementation for domain events
type BaseDomainEvent struct {
	eventID     value_objects.ID
	eventType   string
	aggregateID value_objects.ID
	occurredOn  value_objects.Timestamp
	version     int
}

// NewBaseDomainEvent creates a new base domain event
func NewBaseDomainEvent(eventType string, aggregateID value_objects.ID, version int) BaseDomainEvent {
	return BaseDomainEvent{
		eventID:     value_objects.NewID(),
		eventType:   eventType,
		aggregateID: aggregateID,
		occurredOn:  value_objects.NewTimestamp(),
		version:     version,
	}
}

// EventID returns the event ID
func (e BaseDomainEvent) EventID() value_objects.ID {
	return e.eventID
}

// EventType returns the event type
func (e BaseDomainEvent) EventType() string {
	return e.eventType
}

// AggregateID returns the aggregate ID
func (e BaseDomainEvent) AggregateID() value_objects.ID {
	return e.aggregateID
}

// OccurredOn returns when the event occurred
func (e BaseDomainEvent) OccurredOn() value_objects.Timestamp {
	return e.occurredOn
}

// Version returns the event version
func (e BaseDomainEvent) Version() int {
	return e.version
}
