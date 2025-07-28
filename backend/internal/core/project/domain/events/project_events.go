package events

import (
	"github.com/dewisartika8/cicd-status-notifier-bot/internal/core/project/domain/entities"
	"github.com/dewisartika8/cicd-status-notifier-bot/internal/core/shared/domain/events"
	"github.com/dewisartika8/cicd-status-notifier-bot/internal/core/shared/domain/value_objects"
)

// Project domain event types
const (
	ProjectCreatedEventType     = "project.created"
	ProjectUpdatedEventType     = "project.updated"
	ProjectActivatedEventType   = "project.activated"
	ProjectDeactivatedEventType = "project.deactivated"
	ProjectArchivedEventType    = "project.archived"
)

// ProjectCreatedEvent represents a project creation event
type ProjectCreatedEvent struct {
	events.BaseDomainEvent
	Project *entities.Project
}

// NewProjectCreatedEvent creates a new project created event
func NewProjectCreatedEvent(project *entities.Project) *ProjectCreatedEvent {
	return &ProjectCreatedEvent{
		BaseDomainEvent: events.NewBaseDomainEvent(
			ProjectCreatedEventType,
			project.ID(),
			1,
		),
		Project: project,
	}
}

// ProjectUpdatedEvent represents a project update event
type ProjectUpdatedEvent struct {
	events.BaseDomainEvent
	Project *entities.Project
}

// NewProjectUpdatedEvent creates a new project updated event
func NewProjectUpdatedEvent(project *entities.Project) *ProjectUpdatedEvent {
	return &ProjectUpdatedEvent{
		BaseDomainEvent: events.NewBaseDomainEvent(
			ProjectUpdatedEventType,
			project.ID(),
			1,
		),
		Project: project,
	}
}

// ProjectActivatedEvent represents a project activation event
type ProjectActivatedEvent struct {
	events.BaseDomainEvent
	ProjectID value_objects.ID
}

// NewProjectActivatedEvent creates a new project activated event
func NewProjectActivatedEvent(projectID value_objects.ID) *ProjectActivatedEvent {
	return &ProjectActivatedEvent{
		BaseDomainEvent: events.NewBaseDomainEvent(
			ProjectActivatedEventType,
			projectID,
			1,
		),
		ProjectID: projectID,
	}
}

// ProjectDeactivatedEvent represents a project deactivation event
type ProjectDeactivatedEvent struct {
	events.BaseDomainEvent
	ProjectID value_objects.ID
}

// NewProjectDeactivatedEvent creates a new project deactivated event
func NewProjectDeactivatedEvent(projectID value_objects.ID) *ProjectDeactivatedEvent {
	return &ProjectDeactivatedEvent{
		BaseDomainEvent: events.NewBaseDomainEvent(
			ProjectDeactivatedEventType,
			projectID,
			1,
		),
		ProjectID: projectID,
	}
}

// ProjectArchivedEvent represents a project archive event
type ProjectArchivedEvent struct {
	events.BaseDomainEvent
	ProjectID value_objects.ID
}

// NewProjectArchivedEvent creates a new project archived event
func NewProjectArchivedEvent(projectID value_objects.ID) *ProjectArchivedEvent {
	return &ProjectArchivedEvent{
		BaseDomainEvent: events.NewBaseDomainEvent(
			ProjectArchivedEventType,
			projectID,
			1,
		),
		ProjectID: projectID,
	}
}
