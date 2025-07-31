package domain

import (
	"encoding/json"

	"github.com/dewisartika8/cicd-status-notifier-bot/internal/core/shared/domain/value_objects"
	"github.com/dewisartika8/cicd-status-notifier-bot/pkg/exception"
)

// BuildStatus represents the status of a build
type BuildStatus value_objects.Status

const (
	BuildStatusPending    BuildStatus = "pending"
	BuildStatusInProgress BuildStatus = "in_progress"
	BuildStatusSuccess    BuildStatus = "success"
	BuildStatusFailed     BuildStatus = "failed"
	BuildStatusCancelled  BuildStatus = "cancelled"
	BuildStatusSkipped    BuildStatus = "skipped"
)

// EventType represents the type of build event
type EventType value_objects.Status

const (
	EventTypePush                EventType = "push"
	EventTypePullRequest         EventType = "pull_request"
	EventTypeRelease             EventType = "release"
	EventTypeBuildStarted        EventType = "build_started"
	EventTypeBuildCompleted      EventType = "build_completed"
	EventTypeDeploymentStarted   EventType = "deployment_started"
	EventTypeDeploymentCompleted EventType = "deployment_completed"
)

// BuildEvent represents a CI/CD build event domain entity
type BuildEvent struct {
	id              value_objects.ID
	projectID       value_objects.ID
	eventType       EventType
	status          BuildStatus
	branch          string
	commitSHA       string
	commitMessage   string
	authorName      string
	authorEmail     string
	buildURL        string
	durationSeconds *int
	webhookPayload  json.RawMessage
	createdAt       value_objects.Timestamp
}

// BuildEventParams contains parameters for creating a build event
type BuildEventParams struct {
	ProjectID      value_objects.ID
	EventType      EventType
	Status         BuildStatus
	Branch         string
	CommitSHA      string
	CommitMessage  string
	AuthorName     string
	AuthorEmail    string
	BuildURL       string
	WebhookPayload json.RawMessage
}

// NewBuildEvent creates a new build event entity with validation
func NewBuildEvent(params BuildEventParams) (*BuildEvent, error) {
	// Validate required parameters
	if err := validateBuildEventParams(params); err != nil {
		return nil, err
	}

	buildEvent := &BuildEvent{
		id:             value_objects.NewID(),
		projectID:      params.ProjectID,
		eventType:      params.EventType,
		status:         params.Status,
		branch:         params.Branch,
		commitSHA:      params.CommitSHA,
		commitMessage:  params.CommitMessage,
		authorName:     params.AuthorName,
		authorEmail:    params.AuthorEmail,
		buildURL:       params.BuildURL,
		webhookPayload: params.WebhookPayload,
		createdAt:      value_objects.NewTimestamp(),
	}

	if err := buildEvent.validate(); err != nil {
		return nil, err
	}

	return buildEvent, nil
}

// RestoreBuildEventParams contains parameters for restoring a build event from persistence
type RestoreBuildEventParams struct {
	ID              value_objects.ID
	ProjectID       value_objects.ID
	EventType       EventType
	Status          BuildStatus
	Branch          string
	CommitSHA       string
	CommitMessage   string
	AuthorName      string
	AuthorEmail     string
	BuildURL        string
	DurationSeconds *int
	WebhookPayload  json.RawMessage
	CreatedAt       value_objects.Timestamp
}

// RestoreBuildEvent restores a build event from persistence data
func RestoreBuildEvent(params RestoreBuildEventParams) *BuildEvent {
	return &BuildEvent{
		id:              params.ID,
		projectID:       params.ProjectID,
		eventType:       params.EventType,
		status:          params.Status,
		branch:          params.Branch,
		commitSHA:       params.CommitSHA,
		commitMessage:   params.CommitMessage,
		authorName:      params.AuthorName,
		authorEmail:     params.AuthorEmail,
		buildURL:        params.BuildURL,
		durationSeconds: params.DurationSeconds,
		webhookPayload:  params.WebhookPayload,
		createdAt:       params.CreatedAt,
	}
}

// validateBuildEventParams validates the build event creation parameters
func validateBuildEventParams(params BuildEventParams) error {
	if params.ProjectID.IsNil() {
		return exception.NewDomainError("INVALID_PROJECT_ID", "project ID cannot be nil")
	}

	if params.Branch == "" {
		return exception.NewDomainError("INVALID_BRANCH", "branch cannot be empty")
	}

	if !isValidEventType(params.EventType) {
		return exception.NewDomainError("INVALID_EVENT_TYPE", "invalid event type")
	}

	if !isValidBuildStatus(params.Status) {
		return exception.NewDomainError("INVALID_BUILD_STATUS", "invalid build status")
	}

	return nil
}

// isValidEventType checks if the event type is valid
func isValidEventType(eventType EventType) bool {
	switch eventType {
	case EventTypePush, EventTypePullRequest, EventTypeRelease,
		EventTypeBuildStarted, EventTypeBuildCompleted,
		EventTypeDeploymentStarted, EventTypeDeploymentCompleted:
		return true
	default:
		return false
	}
}

// isValidBuildStatus checks if the build status is valid
func isValidBuildStatus(status BuildStatus) bool {
	switch status {
	case BuildStatusPending, BuildStatusInProgress, BuildStatusSuccess,
		BuildStatusFailed, BuildStatusCancelled, BuildStatusSkipped:
		return true
	default:
		return false
	}
}

// ID returns the build event ID
func (be *BuildEvent) ID() value_objects.ID {
	return be.id
}

// ProjectID returns the project ID
func (be *BuildEvent) ProjectID() value_objects.ID {
	return be.projectID
}

// EventType returns the event type
func (be *BuildEvent) EventType() EventType {
	return be.eventType
}

// Status returns the build status
func (be *BuildEvent) Status() BuildStatus {
	return be.status
}

// Branch returns the git branch
func (be *BuildEvent) Branch() string {
	return be.branch
}

// CommitSHA returns the commit SHA
func (be *BuildEvent) CommitSHA() string {
	return be.commitSHA
}

// CommitMessage returns the commit message
func (be *BuildEvent) CommitMessage() string {
	return be.commitMessage
}

// AuthorName returns the commit author name
func (be *BuildEvent) AuthorName() string {
	return be.authorName
}

// AuthorEmail returns the commit author email
func (be *BuildEvent) AuthorEmail() string {
	return be.authorEmail
}

// BuildURL returns the build URL
func (be *BuildEvent) BuildURL() string {
	return be.buildURL
}

// DurationSeconds returns the build duration in seconds
func (be *BuildEvent) DurationSeconds() *int {
	return be.durationSeconds
}

// WebhookPayload returns the original webhook payload
func (be *BuildEvent) WebhookPayload() json.RawMessage {
	return be.webhookPayload
}

// CreatedAt returns when the build event was created
func (be *BuildEvent) CreatedAt() value_objects.Timestamp {
	return be.createdAt
}

// UpdateStatus updates the build status
func (be *BuildEvent) UpdateStatus(status BuildStatus) {
	be.status = status
}

// SetDuration sets the build duration
func (be *BuildEvent) SetDuration(seconds int) {
	be.durationSeconds = &seconds
}

// IsCompleted checks if the build is completed (success, failed, or cancelled)
func (be *BuildEvent) IsCompleted() bool {
	return be.status == BuildStatusSuccess ||
		be.status == BuildStatusFailed ||
		be.status == BuildStatusCancelled
}

// IsSuccessful checks if the build was successful
func (be *BuildEvent) IsSuccessful() bool {
	return be.status == BuildStatusSuccess
}

// IsFailed checks if the build failed
func (be *BuildEvent) IsFailed() bool {
	return be.status == BuildStatusFailed
}

// validate performs domain validation
func (be *BuildEvent) validate() error {
	if be.projectID.IsNil() {
		return exception.NewDomainError("INVALID_PROJECT_ID", "project ID cannot be empty")
	}

	if be.eventType == "" {
		return exception.NewDomainError("INVALID_EVENT_TYPE", "event type cannot be empty")
	}

	if be.status == "" {
		return exception.NewDomainError("INVALID_BUILD_STATUS", "build status cannot be empty")
	}

	if be.branch == "" {
		return exception.NewDomainError("INVALID_BRANCH", "branch cannot be empty")
	}

	return nil
}
