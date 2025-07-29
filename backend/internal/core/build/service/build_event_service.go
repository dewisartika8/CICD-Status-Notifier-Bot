package service

import (
	"context"

	"github.com/dewisartika8/cicd-status-notifier-bot/internal/core/build/domain"
	"github.com/dewisartika8/cicd-status-notifier-bot/internal/core/build/dto"
	"github.com/dewisartika8/cicd-status-notifier-bot/internal/core/build/port"
	"github.com/dewisartika8/cicd-status-notifier-bot/internal/core/shared/domain/value_objects"
	"github.com/dewisartika8/cicd-status-notifier-bot/pkg/exception"
)

// Dep for BuildEventService
type Dep struct {
	BuildEventRepo port.BuildEventRepository
}

// buildEventService handles build event business logic
type buildEventService struct {
	Dep
}

// NewBuildEventService creates a new build event service
func NewBuildEventService(d Dep) port.BuildEventService {
	return &buildEventService{
		Dep: d,
	}
}

// CreateBuildEvent creates a new build event
func (s *buildEventService) CreateBuildEvent(ctx context.Context, req dto.CreateBuildEventRequest) (*domain.BuildEvent, error) {
	// Create domain entity from request using the refactored constructor
	buildEvent, err := domain.NewBuildEvent(domain.BuildEventParams{
		ProjectID:      req.ProjectID,
		EventType:      req.EventType,
		Status:         req.Status,
		Branch:         req.Branch,
		CommitSHA:      req.CommitSHA,
		CommitMessage:  req.CommitMessage,
		AuthorName:     req.AuthorName,
		AuthorEmail:    req.AuthorEmail,
		BuildURL:       req.BuildURL,
		WebhookPayload: req.WebhookPayload,
	})
	if err != nil {
		return nil, err
	}

	// Set duration if provided
	if req.DurationSeconds != nil {
		buildEvent.SetDuration(*req.DurationSeconds)
	}

	// Validate business rules
	if err := s.validateBuildEvent(buildEvent); err != nil {
		return nil, err
	}

	// Save to repository
	if err := s.BuildEventRepo.Create(ctx, buildEvent); err != nil {
		return nil, err
	}

	return buildEvent, nil
}

// ProcessWebhookEvent processes a webhook event and creates build events
func (s *buildEventService) ProcessWebhookEvent(ctx context.Context, req dto.ProcessWebhookRequest) ([]*domain.BuildEvent, error) {
	// Implement webhook processing logic based on webhook type
	// Parse the webhook payload and create appropriate build events
	// This is a placeholder implementation
	return nil, exception.ErrBuildEventNotFound
}

// GetBuildEvent retrieves a build event by its ID
func (s *buildEventService) GetBuildEvent(ctx context.Context, id value_objects.ID) (*domain.BuildEvent, error) {
	return s.BuildEventRepo.GetByID(ctx, id)
}

// GetBuildEventsByProject retrieves build events for a specific project
func (s *buildEventService) GetBuildEventsByProject(ctx context.Context, projectID value_objects.ID, filters dto.ListBuildEventFilters) ([]*domain.BuildEvent, error) {
	// Set project ID filter
	filters.ProjectID = &projectID
	return s.BuildEventRepo.GetByProjectID(ctx, projectID, filters)
}

// UpdateBuildEventStatus updates the status of a build event
func (s *buildEventService) UpdateBuildEventStatus(ctx context.Context, id value_objects.ID, status domain.BuildStatus, duration *int) error {
	// Get existing build event
	buildEvent, err := s.BuildEventRepo.GetByID(ctx, id)
	if err != nil {
		return err
	}

	// Update status
	buildEvent.UpdateStatus(status)

	// Set duration if provided
	if duration != nil {
		buildEvent.SetDuration(*duration)
	}

	// Save updated build event
	return s.BuildEventRepo.Update(ctx, buildEvent)
}

// GetLatestBuildEvent gets the latest build event for a project
func (s *buildEventService) GetLatestBuildEvent(ctx context.Context, projectID value_objects.ID) (*domain.BuildEvent, error) {
	return s.BuildEventRepo.GetLatestByProjectID(ctx, projectID)
}

// GetBuildMetrics retrieves build metrics for a project
func (s *buildEventService) GetBuildMetrics(ctx context.Context, projectID value_objects.ID) (*domain.BuildMetrics, error) {
	return s.BuildEventRepo.GetBuildMetrics(ctx, projectID)
}

// ListBuildEvents retrieves build events with filtering and pagination
func (s *buildEventService) ListBuildEvents(ctx context.Context, filters dto.ListBuildEventFilters) ([]*domain.BuildEvent, error) {
	return s.BuildEventRepo.List(ctx, filters)
}

// validateBuildEvent validates business rules for build event
func (s *buildEventService) validateBuildEvent(buildEvent *domain.BuildEvent) error {
	// Add custom business validation logic here
	if buildEvent.Branch() == "" {
		return exception.ErrInvalidBranch
	}

	if buildEvent.ProjectID().IsNil() {
		return exception.ErrInvalidProjectID
	}

	return nil
}
