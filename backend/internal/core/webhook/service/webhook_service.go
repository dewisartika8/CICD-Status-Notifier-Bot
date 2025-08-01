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
	// Validate payload
	if payload.WorkflowRun == nil {
		return fmt.Errorf("invalid workflow run payload: workflow_run is nil")
	}

	// Extract workflow information
	workflowInfo := s.extractWorkflowInfo(payload)

	// Create build event
	buildEvent, err := s.createWorkflowBuildEvent(ctx, webhookEvent, payload, workflowInfo)
	if err != nil {
		return fmt.Errorf(errFailedToCreateBuildEvent, err)
	}

	// Create and send notifications
	if err := s.createAndSendWorkflowNotifications(ctx, buildEvent, webhookEvent.ProjectID(), payload, workflowInfo); err != nil {
		return fmt.Errorf(errFailedToCreateNotification, err)
	}

	return nil
}

// workflowInfo holds extracted workflow information
type workflowInfo struct {
	Branch      string
	CommitSHA   string
	BuildURL    string
	BuildStatus buildDomain.BuildStatus
	EventType   buildDomain.EventType
}

// extractWorkflowInfo extracts workflow information from payload
func (s *webhookService) extractWorkflowInfo(payload dto.GitHubActionsPayload) workflowInfo {
	// Extract branch with fallback
	branch := "main"
	if payload.WorkflowRun.HeadBranch != "" {
		branch = payload.WorkflowRun.HeadBranch
	}

	// Extract commit SHA and build URL
	commitSHA := payload.WorkflowRun.HeadSha
	buildURL := payload.WorkflowRun.HTMLURL

	// Use repository commit URL if available
	repoURL := s.safeRepositoryURL(payload)
	if repoURL != "" && commitSHA != "" {
		buildURL = repoURL + commitURLPath + commitSHA
	}

	// Determine build status
	buildStatus := s.determineBuildStatus(payload.WorkflowRun.Conclusion)

	// Determine event type
	eventType := s.determineEventType(payload.Action)

	return workflowInfo{
		Branch:      branch,
		CommitSHA:   commitSHA,
		BuildURL:    buildURL,
		BuildStatus: buildStatus,
		EventType:   eventType,
	}
}

// determineBuildStatus determines build status from workflow conclusion
func (s *webhookService) determineBuildStatus(conclusion string) buildDomain.BuildStatus {
	switch conclusion {
	case "success":
		return buildDomain.BuildStatusSuccess
	case "failure":
		return buildDomain.BuildStatusFailed
	case "cancelled":
		return buildDomain.BuildStatusCancelled
	default:
		return buildDomain.BuildStatusInProgress
	}
}

// determineEventType determines event type from workflow action
func (s *webhookService) determineEventType(action string) buildDomain.EventType {
	switch action {
	case "completed":
		return buildDomain.EventTypeBuildCompleted
	case "requested":
		return buildDomain.EventTypeBuildStarted
	default:
		return buildDomain.EventTypeBuildCompleted
	}
}

// createWorkflowBuildEvent creates a build event for workflow runs
func (s *webhookService) createWorkflowBuildEvent(ctx context.Context, webhookEvent *domain.WebhookEvent, payload dto.GitHubActionsPayload, info workflowInfo) (*buildDomain.BuildEvent, error) {
	buildEventReq := buildDto.CreateBuildEventRequest{
		ProjectID:     webhookEvent.ProjectID(),
		EventType:     info.EventType,
		Status:        info.BuildStatus,
		Branch:        info.Branch,
		CommitSHA:     info.CommitSHA,
		CommitMessage: "", // Not available in workflow run payload
		AuthorName:    "",
		AuthorEmail:   "",
		BuildURL:      info.BuildURL,
	}

	return s.BuildService.CreateBuildEvent(ctx, buildEventReq)
}

// createAndSendWorkflowNotifications creates and sends notifications for workflow events
func (s *webhookService) createAndSendWorkflowNotifications(ctx context.Context, buildEvent *buildDomain.BuildEvent, projectID value_objects.ID, payload dto.GitHubActionsPayload, info workflowInfo) error {
	if buildEvent == nil || s.NotificationLogService == nil {
		return nil
	}

	// Build notification message
	statusText := s.buildStatusText(info.BuildStatus)
	message := fmt.Sprintf("🔔 %s %s for %s on branch %s",
		payload.WorkflowRun.Name, statusText, s.safeRepositoryName(payload), info.Branch)

	// Create notifications
	notifications, err := s.NotificationLogService.CreateNotificationForBuildEvent(
		ctx,
		buildEvent.ID(),
		projectID,
		message,
	)
	if err != nil {
		return err
	}

	// Immediately process the created notifications (same as original behavior)
	if len(notifications) > 0 {
		for _, notification := range notifications {
			if err := s.NotificationLogService.SendNotification(ctx, notification.ID()); err != nil {
				// Log the error but don't fail the entire webhook processing
				// The notification will remain pending and can be retried later
				continue
			}
		}
	}

	return nil
}

// buildStatusText returns the status text with emoji for notifications
func (s *webhookService) buildStatusText(status buildDomain.BuildStatus) string {
	switch status {
	case buildDomain.BuildStatusSuccess:
		return "✅ succeeded"
	case buildDomain.BuildStatusFailed:
		return "❌ failed"
	case buildDomain.BuildStatusCancelled:
		return "⏹️ cancelled"
	default:
		return "🔄 is running"
	}
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
	if buildEvent != nil && s.NotificationLogService != nil {
		message := s.buildNotificationMessage(payload, branch, commitInfo)
		notifications, err := s.NotificationLogService.CreateNotificationForBuildEvent(
			ctx,
			buildEvent.ID(),
			webhookEvent.ProjectID(),
			message,
		)
		if err != nil {
			return fmt.Errorf(errFailedToCreateNotification, err)
		}

		// Immediately process the created notifications
		if len(notifications) > 0 {
			for _, notification := range notifications {
				if err := s.NotificationLogService.SendNotification(ctx, notification.ID()); err != nil {
					// Log the error but don't fail the entire webhook processing
					// The notification will remain pending and can be retried later
					continue
				}
			}
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
	status := s.determinePRStatus(payload.Action, pr.State)

	// Create build event request
	buildEventReq := s.createPRBuildEventRequest(webhookEvent, payload, pr, status)

	// Create build event
	buildEvent, err := s.BuildService.CreateBuildEvent(ctx, buildEventReq)
	if err != nil {
		return fmt.Errorf(errFailedToCreateBuildEvent, err)
	}

	// Create notification if build event was created successfully
	if buildEvent != nil && s.NotificationLogService != nil {
		message := s.createPRNotificationMessage(payload, pr)
		notifications, err := s.NotificationLogService.CreateNotificationForBuildEvent(
			ctx,
			buildEvent.ID(),
			webhookEvent.ProjectID(),
			message,
		)
		if err != nil {
			return fmt.Errorf(errFailedToCreateNotification, err)
		}

		// Immediately process the created notifications
		if len(notifications) > 0 {
			for _, notification := range notifications {
				if err := s.NotificationLogService.SendNotification(ctx, notification.ID()); err != nil {
					// Log the error but don't fail the entire webhook processing
					// The notification will remain pending and can be retried later
					continue
				}
			}
		}
	}

	return nil
}

// determinePRStatus determines the build status based on PR action and state
func (s *webhookService) determinePRStatus(action, state string) buildDomain.BuildStatus {
	status := buildDomain.BuildStatusPending
	if action == "closed" {
		if state == "merged" {
			status = buildDomain.BuildStatusSuccess
		} else {
			status = buildDomain.BuildStatusCancelled
		}
	}
	return status
}

// createPRBuildEventRequest creates build event request for pull request
func (s *webhookService) createPRBuildEventRequest(webhookEvent *domain.WebhookEvent, payload dto.GitHubActionsPayload, pr *dto.PullRequest, status buildDomain.BuildStatus) buildDto.CreateBuildEventRequest {
	branch := pr.Head.Ref
	commitSHA := pr.Head.SHA
	commitMessage := fmt.Sprintf("Pull Request: %s", pr.Title)
	authorName := pr.User.Name
	authorEmail := pr.User.Email
	buildURL := pr.HTMLURL

	return buildDto.CreateBuildEventRequest{
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
}

// createPRNotificationMessage creates notification message for pull request
func (s *webhookService) createPRNotificationMessage(payload dto.GitHubActionsPayload, pr *dto.PullRequest) string {
	actionText := s.getPRActionText(payload.Action, pr.State)
	return fmt.Sprintf("📋 *Pull Request %s*\n*Project:* %s\n*Title:* %s\n*Branch:* %s → %s\n*Author:* %s",
		actionText, s.safeRepositoryName(payload), pr.Title, pr.Head.Ref, pr.Base.Ref, pr.User.Name)
}

// getPRActionText returns the action text for pull request notifications
func (s *webhookService) getPRActionText(action, state string) string {
	switch action {
	case "closed":
		if state == "merged" {
			return "merged"
		}
		return "closed"
	case "reopened":
		return "reopened"
	case "synchronize":
		return "updated"
	default:
		return "opened"
	}
}
