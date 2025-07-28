package dto

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

// EventType represents the type of build event
type EventType string

const (
	EventTypeBuildStarted      EventType = "build_started"
	EventTypeBuildSuccess      EventType = "build_success"
	EventTypeBuildFailed       EventType = "build_failed"
	EventTypeTestStarted       EventType = "test_started"
	EventTypeTestPassed        EventType = "test_passed"
	EventTypeTestFailed        EventType = "test_failed"
	EventTypeDeploymentStarted EventType = "deployment_started"
	EventTypeDeploymentSuccess EventType = "deployment_success"
	EventTypeDeploymentFailed  EventType = "deployment_failed"
)

// BuildStatus represents the status of a build
type BuildStatus string

const (
	BuildStatusSuccess BuildStatus = "success"
	BuildStatusFailed  BuildStatus = "failed"
	BuildStatusPending BuildStatus = "pending"
	BuildStatusRunning BuildStatus = "running"
)

// BuildEvent represents a CI/CD build event entity
type BuildEvent struct {
	ID              uuid.UUID       `json:"id"`
	ProjectID       uuid.UUID       `json:"project_id" validate:"required"`
	EventType       EventType       `json:"event_type" validate:"required"`
	Status          BuildStatus     `json:"status" validate:"required"`
	Branch          string          `json:"branch" validate:"required"`
	CommitSHA       string          `json:"commit_sha,omitempty"`
	CommitMessage   string          `json:"commit_message,omitempty"`
	AuthorName      string          `json:"author_name,omitempty"`
	AuthorEmail     string          `json:"author_email,omitempty"`
	BuildURL        string          `json:"build_url,omitempty"`
	DurationSeconds *int            `json:"duration_seconds,omitempty"`
	WebhookPayload  json.RawMessage `json:"webhook_payload,omitempty"`
	CreatedAt       time.Time       `json:"created_at"`
}

// NewBuildEvent creates a new build event entity
func NewBuildEvent(projectID uuid.UUID, eventType EventType, status BuildStatus, branch string) *BuildEvent {
	return &BuildEvent{
		ID:        uuid.New(),
		ProjectID: projectID,
		EventType: eventType,
		Status:    status,
		Branch:    branch,
		CreatedAt: time.Now(),
	}
}

// Validate performs basic validation on build event entity
func (be *BuildEvent) Validate() error {
	if be.ProjectID == uuid.Nil {
		return ErrInvalidProjectID
	}
	if be.EventType == "" {
		return ErrInvalidEventType
	}
	if be.Status == "" {
		return ErrInvalidBuildStatus
	}
	if be.Branch == "" {
		return ErrInvalidBranch
	}
	return nil
}

// SetCommitInfo sets commit-related information
func (be *BuildEvent) SetCommitInfo(sha, message, authorName, authorEmail string) {
	be.CommitSHA = sha
	be.CommitMessage = message
	be.AuthorName = authorName
	be.AuthorEmail = authorEmail
}

// SetBuildInfo sets build-related information
func (be *BuildEvent) SetBuildInfo(buildURL string, durationSeconds *int) {
	be.BuildURL = buildURL
	be.DurationSeconds = durationSeconds
}

// SetWebhookPayload sets the raw webhook payload
func (be *BuildEvent) SetWebhookPayload(payload json.RawMessage) {
	be.WebhookPayload = payload
}

// IsFailureEvent returns true if the event represents a failure
func (be *BuildEvent) IsFailureEvent() bool {
	return be.Status == BuildStatusFailed ||
		be.EventType == EventTypeBuildFailed ||
		be.EventType == EventTypeTestFailed ||
		be.EventType == EventTypeDeploymentFailed
}

// IsSuccessEvent returns true if the event represents a success
func (be *BuildEvent) IsSuccessEvent() bool {
	return be.Status == BuildStatusSuccess ||
		be.EventType == EventTypeBuildSuccess ||
		be.EventType == EventTypeTestPassed ||
		be.EventType == EventTypeDeploymentSuccess
}
