package service

import (
	"context"
	"encoding/json"

	buildPort "github.com/dewisartika8/cicd-status-notifier-bot/internal/core/build/port"
	projectPort "github.com/dewisartika8/cicd-status-notifier-bot/internal/core/project/port"
	"github.com/dewisartika8/cicd-status-notifier-bot/internal/core/shared/domain/value_objects"
	"github.com/dewisartika8/cicd-status-notifier-bot/internal/core/webhook/domain"
	"github.com/dewisartika8/cicd-status-notifier-bot/internal/core/webhook/dto"
	"github.com/dewisartika8/cicd-status-notifier-bot/internal/core/webhook/port"
	"github.com/dewisartika8/cicd-status-notifier-bot/pkg/crypto"
)

// Dep defines the dependencies for WebhookService
type Dep struct {
	WebhookEventRepo  port.WebhookEventRepository
	ProjectService    projectPort.ProjectService
	BuildService      buildPort.BuildEventService
	SignatureVerifier crypto.SignatureVerifier
}

// webhookService handles webhook business logic
type webhookService struct {
	Dep
}

// NewWebhookService creates a new webhook service
func NewWebhookService(d Dep) port.WebhookService {
	return &webhookService{
		Dep: d,
	}
}

// ProcessWebhook processes an incoming webhook request
func (s *webhookService) ProcessWebhook(ctx context.Context, req dto.ProcessWebhookRequest) (*domain.WebhookEvent, error) {
	// 1. Verify project exists
	project, err := s.ProjectService.GetProject(ctx, req.ProjectID)
	if err != nil {
		return nil, domain.NewWebhookProjectNotFoundError(req.ProjectID.String())
	}

	// 2. Verify webhook signature
	if !s.SignatureVerifier.VerifySignature(project.WebhookSecret(), req.Signature, req.Body) {
		return nil, domain.ErrWebhookInvalidSignature
	}

	// 3. Check if this webhook has already been processed (idempotency)
	if req.DeliveryID != "" {
		exists, err := s.WebhookEventRepo.ExistsByDeliveryID(ctx, req.DeliveryID)
		if err != nil {
			return nil, domain.NewWebhookProcessingFailedError("failed to check duplicate delivery")
		}
		if exists {
			// Return existing webhook event
			return s.WebhookEventRepo.GetByDeliveryID(ctx, req.DeliveryID)
		}
	}

	// 4. Convert payload to JSON string
	payloadBytes, err := json.Marshal(req.Payload)
	if err != nil {
		return nil, domain.NewWebhookInvalidPayloadError("failed to marshal payload")
	}

	// 5. Create webhook event domain entity
	webhookEvent, err := domain.NewWebhookEvent(
		req.ProjectID,
		req.EventType,
		string(payloadBytes),
		req.Signature,
		req.DeliveryID,
	)
	if err != nil {
		return nil, err
	}

	// 6. Store webhook event
	if err := s.WebhookEventRepo.Create(ctx, webhookEvent); err != nil {
		return nil, domain.NewWebhookProcessingFailedError("failed to store webhook event")
	}

	// 7. Process the webhook based on event type
	if err := s.processWebhookEvent(ctx, webhookEvent, req.Payload); err != nil {
		// Log error but don't fail the webhook processing
		// The webhook event is already stored, so we can retry processing later
		return webhookEvent, nil
	}

	// 8. Mark as processed
	webhookEvent.MarkAsProcessed()
	if err := s.WebhookEventRepo.Update(ctx, webhookEvent); err != nil {
		// Log error but don't fail - the main processing is done
		return webhookEvent, nil
	}

	return webhookEvent, nil
}

// VerifyWebhookSignature verifies the webhook signature
func (s *webhookService) VerifyWebhookSignature(secret, signature string, body []byte) bool {
	return s.SignatureVerifier.VerifySignature(secret, signature, body)
}

// GetWebhookEvent retrieves a webhook event by its ID
func (s *webhookService) GetWebhookEvent(ctx context.Context, id value_objects.ID) (*domain.WebhookEvent, error) {
	return s.WebhookEventRepo.GetByID(ctx, id)
}

// GetWebhookEventsByProject retrieves webhook events for a specific project
func (s *webhookService) GetWebhookEventsByProject(ctx context.Context, projectID value_objects.ID, limit, offset int) ([]*domain.WebhookEvent, error) {
	return s.WebhookEventRepo.GetByProjectID(ctx, projectID, limit, offset)
}

// ReprocessFailedWebhooks reprocesses failed webhook events
func (s *webhookService) ReprocessFailedWebhooks(ctx context.Context, limit int) error {
	unprocessedEvents, err := s.WebhookEventRepo.GetUnprocessedEvents(ctx, limit)
	if err != nil {
		return domain.NewWebhookProcessingFailedError("failed to get unprocessed events")
	}

	for _, event := range unprocessedEvents {
		// Parse the payload
		var payload dto.GitHubActionsPayload
		if err := json.Unmarshal([]byte(event.Payload()), &payload); err != nil {
			continue // Skip invalid payloads
		}

		// Process the event
		if err := s.processWebhookEvent(ctx, event, payload); err != nil {
			continue // Skip failed processing
		}

		// Mark as processed
		event.MarkAsProcessed()
		if err := s.WebhookEventRepo.Update(ctx, event); err != nil {
			continue // Log error but continue
		}
	}

	return nil
}

// processWebhookEvent processes the webhook event based on its type
func (s *webhookService) processWebhookEvent(ctx context.Context, webhookEvent *domain.WebhookEvent, payload dto.GitHubActionsPayload) error {
	switch webhookEvent.EventType() {
	case domain.WorkflowRunEvent:
		return s.processWorkflowRunEvent(ctx, webhookEvent, payload)
	case domain.PushEvent:
		return s.processPushEvent(ctx, webhookEvent, payload)
	case domain.PullRequestEvent:
		return s.processPullRequestEvent(ctx, webhookEvent, payload)
	default:
		return domain.NewWebhookInvalidEventError(string(webhookEvent.EventType()))
	}
}

// processWorkflowRunEvent processes workflow_run events
func (s *webhookService) processWorkflowRunEvent(ctx context.Context, webhookEvent *domain.WebhookEvent, payload dto.GitHubActionsPayload) error {
	// Implementation for workflow run processing
	// This would typically create build events, send notifications, etc.
	// For now, we just log the event type
	return nil
}

// processPushEvent processes push events
func (s *webhookService) processPushEvent(ctx context.Context, webhookEvent *domain.WebhookEvent, payload dto.GitHubActionsPayload) error {
	// Implementation for push event processing
	return nil
}

// processPullRequestEvent processes pull request events
func (s *webhookService) processPullRequestEvent(ctx context.Context, webhookEvent *domain.WebhookEvent, payload dto.GitHubActionsPayload) error {
	// Implementation for pull request event processing
	return nil
}
