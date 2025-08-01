package service

import (
	"context"
	"encoding/json"
	"fmt"

	buildDomain "github.com/dewisartika8/cicd-status-notifier-bot/internal/core/build/domain"
	buildDto "github.com/dewisartika8/cicd-status-notifier-bot/internal/core/build/dto"
	buildPort "github.com/dewisartika8/cicd-status-notifier-bot/internal/core/build/port"
	notificationPort "github.com/dewisartika8/cicd-status-notifier-bot/internal/core/notification/port"
	projectPort "github.com/dewisartika8/cicd-status-notifier-bot/internal/core/project/port"
	"github.com/dewisartika8/cicd-status-notifier-bot/internal/core/shared/domain/value_objects"
	"github.com/dewisartika8/cicd-status-notifier-bot/internal/core/webhook/domain"
	"github.com/dewisartika8/cicd-status-notifier-bot/internal/core/webhook/dto"
	"github.com/dewisartika8/cicd-status-notifier-bot/internal/core/webhook/port"
	"github.com/dewisartika8/cicd-status-notifier-bot/pkg/crypto"
)

// Constants for error messages and URL patterns
const (
	commitURLPath                 = "/commit/"
	errFailedToCreateBuildEvent   = "failed to create build event: %w"
	errFailedToCreateNotification = "failed to create notification: %w"
)

// Dep defines the dependencies for WebhookService
type Dep struct {
	WebhookEventRepo       port.WebhookEventRepository
	ProjectService         projectPort.ProjectService
	BuildService           buildPort.BuildEventService
	NotificationLogService notificationPort.NotificationLogService
	SignatureVerifier      crypto.SignatureVerifier
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
	// Check if WorkflowRun is nil (invalid payload)
	if payload.WorkflowRun == nil {
		return fmt.Errorf("invalid workflow run payload: workflow_run is nil")
	}

	// Extract information from workflow run payload
	branch := "main" // Default branch
	if payload.WorkflowRun.HeadBranch != "" {
		branch = payload.WorkflowRun.HeadBranch
	}

	// Get commit information
	var commitSHA, buildURL string

	if payload.WorkflowRun.HeadSha != "" {
		commitSHA = payload.WorkflowRun.HeadSha
	}

	if payload.WorkflowRun.HTMLURL != "" {
		buildURL = payload.WorkflowRun.HTMLURL
	}

	// Use repository info if available
	repoURL := s.safeRepositoryURL(payload)
	if repoURL != "" && commitSHA != "" {
		buildURL = repoURL + commitURLPath + commitSHA
	}

	// Determine build status based on workflow conclusion
	var buildStatus buildDomain.BuildStatus
	switch payload.WorkflowRun.Conclusion {
	case "success":
		buildStatus = buildDomain.BuildStatusSuccess
	case "failure":
		buildStatus = buildDomain.BuildStatusFailed
	case "cancelled":
		buildStatus = buildDomain.BuildStatusCancelled
	default:
		buildStatus = buildDomain.BuildStatusInProgress
	}

	// Determine event type based on workflow action
	var eventType buildDomain.EventType
	switch payload.Action {
	case "completed":
		eventType = buildDomain.EventTypeBuildCompleted
	case "requested":
		eventType = buildDomain.EventTypeBuildStarted
	default:
		eventType = buildDomain.EventTypeBuildCompleted
	}

	// Create build event request
	buildEventReq := buildDto.CreateBuildEventRequest{
		ProjectID:     webhookEvent.ProjectID(),
		EventType:     eventType,
		Status:        buildStatus,
		Branch:        branch,
		CommitSHA:     commitSHA,
		CommitMessage: "", // Not available in workflow run payload
		AuthorName:    "",
		AuthorEmail:   "",
		BuildURL:      buildURL,
	}

	// Create build event
	buildEvent, err := s.BuildService.CreateBuildEvent(ctx, buildEventReq)
	if err != nil {
		return fmt.Errorf(errFailedToCreateBuildEvent, err)
	}

	// Create notification if build event was created successfully
	if buildEvent != nil {
		var statusText string
		switch buildStatus {
		case buildDomain.BuildStatusSuccess:
			statusText = "✅ succeeded"
		case buildDomain.BuildStatusFailed:
			statusText = "❌ failed"
		case buildDomain.BuildStatusCancelled:
			statusText = "⏹️ cancelled"
		default:
			statusText = "🔄 is running"
		}

		message := fmt.Sprintf("🔔 %s %s for %s on branch %s",
			payload.WorkflowRun.Name, statusText, s.safeRepositoryName(payload), branch)

		_, err = s.NotificationLogService.CreateNotificationForBuildEvent(
			ctx,
			buildEvent.ID(),
			webhookEvent.ProjectID(),
			message,
		)
		if err != nil {
			return fmt.Errorf(errFailedToCreateNotification, err)
		}
	}

	return nil
}

// processPushEvent processes push events
func (s *webhookService) processPushEvent(ctx context.Context, webhookEvent *domain.WebhookEvent, payload dto.GitHubActionsPayload) error {
	branch := s.extractBranchFromRef(payload.Ref)
	commitInfo := s.extractCommitInfo(payload)

	buildEventReq := buildDto.CreateBuildEventRequest{
		ProjectID:     webhookEvent.ProjectID(),
		EventType:     buildDomain.EventTypePush,
		Status:        buildDomain.BuildStatusSuccess,
		Branch:        branch,
		CommitSHA:     commitInfo.SHA,
		CommitMessage: commitInfo.Message,
		AuthorName:    commitInfo.AuthorName,
		AuthorEmail:   commitInfo.AuthorEmail,
		BuildURL:      commitInfo.BuildURL,
		WebhookPayload: func() []byte {
			data, _ := json.Marshal(payload)
			return data
		}(),
	}

	buildEvent, err := s.BuildService.CreateBuildEvent(ctx, buildEventReq)
	if err != nil {
		return fmt.Errorf(errFailedToCreateBuildEvent, err)
	}

	// Create notification if build event was created successfully
	if buildEvent != nil {
		message := s.buildNotificationMessage(payload, branch, commitInfo)
		_, err = s.NotificationLogService.CreateNotificationForBuildEvent(
			ctx,
			buildEvent.ID(),
			webhookEvent.ProjectID(),
			message,
		)
		if err != nil {
			return fmt.Errorf(errFailedToCreateNotification, err)
		}
	}

	return nil
}

// extractBranchFromRef extracts branch name from git ref
func (s *webhookService) extractBranchFromRef(ref string) string {
	if ref == "" {
		return "main"
	}
	if len(ref) > 11 && ref[:11] == "refs/heads/" {
		return ref[11:]
	}
	return "main"
}

// commitInfo holds commit information
type commitInfo struct {
	SHA         string
	Message     string
	AuthorName  string
	AuthorEmail string
	BuildURL    string
}

// extractCommitInfo extracts commit information from payload with enhanced nil safety
func (s *webhookService) extractCommitInfo(payload dto.GitHubActionsPayload) (result commitInfo) {
	// Use named return and defer to ensure we always return something valid
	defer func() {
		if r := recover(); r != nil {
			// If panic occurs, return the fallback commit info
			result = s.createFallbackCommitInfo(payload)
		}
	}()

	// Try HeadCommit first
	if payload.HeadCommit != nil {
		result = commitInfo{
			SHA:         s.safeString(payload.HeadCommit.ID),
			Message:     s.safeString(payload.HeadCommit.Message),
			AuthorName:  s.safeString(payload.HeadCommit.Author.Name),
			AuthorEmail: s.safeString(payload.HeadCommit.Author.Email),
			BuildURL:    s.safeString(payload.HeadCommit.URL),
		}
		return result
	}

	// Try Commits array as fallback
	if len(payload.Commits) > 0 {
		lastCommit := payload.Commits[len(payload.Commits)-1]
		buildURL := ""
		commitID := s.safeString(lastCommit.ID)
		repoURL := s.safeRepositoryURL(payload)
		if repoURL != "" && commitID != "" {
			buildURL = repoURL + commitURLPath + commitID
		}

		result = commitInfo{
			SHA:         commitID,
			Message:     s.safeString(lastCommit.Message),
			AuthorName:  s.safeString(lastCommit.Author.Name),
			AuthorEmail: s.safeString(lastCommit.Author.Email),
			BuildURL:    buildURL,
		}
		return result
	}

	// Final fallback
	return s.createFallbackCommitInfo(payload)
}

// safeString returns empty string if input is empty or contains only whitespace
func (s *webhookService) safeString(str string) string {
	if str == "" {
		return ""
	}
	// You could add additional sanitization here if needed
	return str
}

// safeRepositoryName safely extracts repository name from payload
func (s *webhookService) safeRepositoryName(payload dto.GitHubActionsPayload) string {
	if payload.Repository.FullName != "" {
		return payload.Repository.FullName
	}
	if payload.Repository.Name != "" {
		return payload.Repository.Name
	}
	return "Unknown Repository"
}

// safeRepositoryURL safely extracts repository HTML URL from payload
func (s *webhookService) safeRepositoryURL(payload dto.GitHubActionsPayload) string {
	if payload.Repository.HTMLURL != "" {
		return payload.Repository.HTMLURL
	}
	return ""
}

// createFallbackCommitInfo creates fallback commit info when HeadCommit and Commits are not available
func (s *webhookService) createFallbackCommitInfo(payload dto.GitHubActionsPayload) commitInfo {
	sha := payload.After
	if sha == "" {
		sha = "unknown"
	}

	message := "Push to " + s.extractBranchFromRef(payload.Ref)

	var authorName, authorEmail string
	if payload.Pusher != nil {
		authorName = payload.Pusher.Name
		authorEmail = payload.Pusher.Email
	}

	buildURL := ""
	// Safely access Repository fields
	repoURL := s.safeRepositoryURL(payload)
	if repoURL != "" && sha != "" && sha != "unknown" {
		buildURL = repoURL + commitURLPath + sha
	}

	return commitInfo{
		SHA:         sha,
		Message:     message,
		AuthorName:  authorName,
		AuthorEmail: authorEmail,
		BuildURL:    buildURL,
	}
}

// buildNotificationMessage builds notification message with safe string handling
func (s *webhookService) buildNotificationMessage(payload dto.GitHubActionsPayload, branch string, commit commitInfo) string {
	// Safely extract repository name with fallback
	repoName := s.safeRepositoryName(payload)

	// Safely extract author name with fallback
	authorName := "Unknown Author"
	if commit.AuthorName != "" {
		authorName = commit.AuthorName
	}

	// Safely extract commit message with fallback
	commitMessage := "No message"
	if commit.Message != "" {
		commitMessage = commit.Message
	}

	return fmt.Sprintf("📤 *Push Event*\n*Project:* %s\n*Branch:* %s\n*Commit:* %s\n*Author:* %s",
		repoName, branch, commitMessage, authorName)
}

// processPullRequestEvent processes pull request events
func (s *webhookService) processPullRequestEvent(ctx context.Context, webhookEvent *domain.WebhookEvent, payload dto.GitHubActionsPayload) error {
	if payload.PullRequest == nil {
		return fmt.Errorf("pull request data is missing")
	}

	pr := payload.PullRequest

	// Determine status based on PR action
	status := buildDomain.BuildStatusPending
	if payload.Action == "closed" {
		if pr.State == "merged" {
			status = buildDomain.BuildStatusSuccess
		} else {
			status = buildDomain.BuildStatusCancelled
		}
	}

	// Use head branch and commit information
	branch := pr.Head.Ref
	commitSHA := pr.Head.SHA
	commitMessage := fmt.Sprintf("Pull Request: %s", pr.Title)
	authorName := pr.User.Name
	authorEmail := pr.User.Email
	buildURL := pr.HTMLURL

	// Create build event request
	buildEventReq := buildDto.CreateBuildEventRequest{
		ProjectID:     webhookEvent.ProjectID(),
		EventType:     buildDomain.EventTypePullRequest,
		Status:        status,
		Branch:        branch,
		CommitSHA:     commitSHA,
		CommitMessage: commitMessage,
		AuthorName:    authorName,
		AuthorEmail:   authorEmail,
		BuildURL:      buildURL,
		WebhookPayload: func() []byte {
			data, _ := json.Marshal(payload)
			return data
		}(),
	}

	// Create build event
	buildEvent, err := s.BuildService.CreateBuildEvent(ctx, buildEventReq)
	if err != nil {
		return fmt.Errorf(errFailedToCreateBuildEvent, err)
	}

	// Create notification if build event was created successfully
	if buildEvent != nil {
		// Create notification message
		actionText := "opened"
		switch payload.Action {
		case "closed":
			if pr.State == "merged" {
				actionText = "merged"
			} else {
				actionText = "closed"
			}
		case "reopened":
			actionText = "reopened"
		case "synchronize":
			actionText = "updated"
		}

		message := fmt.Sprintf("📋 *Pull Request %s*\n*Project:* %s\n*Title:* %s\n*Branch:* %s → %s\n*Author:* %s",
			actionText, s.safeRepositoryName(payload), pr.Title, pr.Head.Ref, pr.Base.Ref, authorName)

		// Send notifications
		_, err = s.NotificationLogService.CreateNotificationForBuildEvent(
			ctx,
			buildEvent.ID(),
			webhookEvent.ProjectID(),
			message,
		)
		if err != nil {
			return fmt.Errorf(errFailedToCreateNotification, err)
		}
	}

	return nil
}
