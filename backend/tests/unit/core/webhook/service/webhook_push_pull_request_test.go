package service_test

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	buildDomain "github.com/dewisartika8/cicd-status-notifier-bot/internal/core/build/domain"
	buildDto "github.com/dewisartika8/cicd-status-notifier-bot/internal/core/build/dto"
	notificationDomain "github.com/dewisartika8/cicd-status-notifier-bot/internal/core/notification/domain"
	projectDomain "github.com/dewisartika8/cicd-status-notifier-bot/internal/core/project/domain"
	"github.com/dewisartika8/cicd-status-notifier-bot/internal/core/shared/domain/value_objects"
	"github.com/dewisartika8/cicd-status-notifier-bot/internal/core/webhook/domain"
	"github.com/dewisartika8/cicd-status-notifier-bot/internal/core/webhook/dto"
	"github.com/dewisartika8/cicd-status-notifier-bot/internal/core/webhook/service"
	"github.com/dewisartika8/cicd-status-notifier-bot/tests/mocks"
)

// Test constants for push/PR tests
const (
	pushTestUserEmail      = "testuser@example.com"
	pushTestUserName       = "Test User"
	pushTestFeatureMessage = "Add new feature"
	pushTestCommitURL      = "https://github.com/test/repo/commit/abc123def456789"
	pushTestPullRequestURL = "https://github.com/test/repo/pull/123"
	pushTestFeatureBranch  = "feature-branch"
	pushTestRepoURL        = "https://github.com/test/repo"
	pushTestWebhookSecret  = "test-webhook-secret"
	pushTestSignature      = "test-signature"
	pushTestDeliveryID     = "test-delivery-id"
	pushTestProjectName    = "Test Project"
	pushWebhookEventType   = "*domain.WebhookEvent"
)

func TestProcessPushEventTDD(t *testing.T) {
	// Test data setup
	projectID := value_objects.NewID()

	pushPayload := dto.GitHubActionsPayload{
		Action:  "push",
		Ref:     "refs/heads/main",
		Before:  "0000000000000000000000000000000000000000",
		After:   "abc123def456789",
		Created: true,
		Deleted: false,
		Forced:  false,
		Compare: "https://github.com/test/repo/compare/000000...abc123",
		Repository: struct {
			Name     string `json:"name"`
			FullName string `json:"full_name"`
			HTMLURL  string `json:"html_url"`
		}{
			Name:     "test-repo",
			FullName: "test/repo",
			HTMLURL:  pushTestRepoURL,
		},
		Sender: struct {
			Login string `json:"login"`
			Email string `json:"email,omitempty"`
		}{
			Login: "testuser",
			Email: pushTestUserEmail,
		},
		HeadCommit: &dto.Commit{
			ID:        "abc123def456789",
			Message:   pushTestFeatureMessage,
			Timestamp: time.Now().Format(time.RFC3339),
			URL:       pushTestCommitURL,
			Author: dto.User{
				Login: "testuser",
				Email: pushTestUserEmail,
				Name:  pushTestUserName,
			},
			Committer: dto.User{
				Login: "testuser",
				Email: pushTestUserEmail,
				Name:  pushTestUserName,
			},
			Added:    []string{"file1.go"},
			Modified: []string{"file2.go"},
			Removed:  []string{},
		},
		Commits: []dto.Commit{
			{
				ID:        "abc123def456789",
				Message:   pushTestFeatureMessage,
				Timestamp: time.Now().Format(time.RFC3339),
				Author: dto.User{
					Login: "testuser",
					Email: pushTestUserEmail,
					Name:  pushTestUserName,
				},
				Committer: dto.User{
					Login: "testuser",
					Email: pushTestUserEmail,
					Name:  pushTestUserName,
				},
			},
		},
		Pusher: &dto.User{
			Login: "testuser",
			Email: pushTestUserEmail,
			Name:  pushTestUserName,
		},
	}

	t.Run("TDD_Red_ProcessPushEvent_should_create_build_event_and_notification", func(t *testing.T) {
		// Arrange - Setup mocks
		mockWebhookRepo := &mocks.MockWebhookEventRepository{}
		mockProjectService := &MockProjectServiceTDD{}
		mockBuildService := &MockBuildEventServiceTDD{}
		mockNotificationService := &MockNotificationLogServiceTDD{}
		mockSignatureVerifier := &mocks.MockSignatureVerifier{}

		// Setup webhook service
		webhookService := service.NewWebhookService(service.Dep{
			WebhookEventRepo:       mockWebhookRepo,
			ProjectService:         mockProjectService,
			BuildService:           mockBuildService,
			NotificationLogService: mockNotificationService,
			SignatureVerifier:      mockSignatureVerifier,
		})

		// Create expected project
		expectedProject, err := projectDomain.NewProject(pushTestProjectName, pushTestRepoURL, pushTestWebhookSecret, nil)
		require.NoError(t, err)

		// Create expected build event
		expectedBuildEvent, err := buildDomain.NewBuildEvent(buildDomain.BuildEventParams{
			ProjectID:     projectID,
			EventType:     buildDomain.EventTypePush,
			Status:        buildDomain.BuildStatusSuccess,
			Branch:        "main",
			CommitSHA:     "abc123def456789",
			CommitMessage: pushTestFeatureMessage,
			AuthorName:    pushTestUserName,
			AuthorEmail:   pushTestUserEmail,
			BuildURL:      pushTestCommitURL,
		})
		require.NoError(t, err)

		// Create expected notification logs
		expectedNotifications := []*notificationDomain.NotificationLog{}
		notificationLog, err := notificationDomain.NewNotificationLog(
			expectedBuildEvent.ID(),
			projectID,
			notificationDomain.NotificationChannelTelegram,
			"123456789",
			"Test push notification message",
			3,
		)
		require.NoError(t, err)
		expectedNotifications = append(expectedNotifications, notificationLog)

		// Setup mock expectations for project service
		mockProjectService.On("GetProject", mock.Anything, projectID).
			Return(expectedProject, nil).Once()

		// Setup mock expectations for signature verification
		mockSignatureVerifier.On("VerifySignature", pushTestWebhookSecret, pushTestSignature, mock.Anything).
			Return(true).Once()

		// Setup mock expectations for duplicate check
		mockWebhookRepo.On("ExistsByDeliveryID", mock.Anything, pushTestDeliveryID).
			Return(false, nil).Once()

		// Setup mock expectations for webhook creation
		mockWebhookRepo.On("Create", mock.Anything, mock.AnythingOfType(pushWebhookEventType)).
			Return(nil).Once()

		// Setup mock expectations for webhook update
		mockWebhookRepo.On("Update", mock.Anything, mock.AnythingOfType(pushWebhookEventType)).
			Return(nil).Once()

		// Setup mock expectations for build service
		mockBuildService.On("CreateBuildEvent", mock.Anything, mock.MatchedBy(func(req buildDto.CreateBuildEventRequest) bool {
			return req.ProjectID == projectID &&
				req.EventType == buildDomain.EventTypePush &&
				req.Status == buildDomain.BuildStatusSuccess &&
				req.Branch == "main" &&
				req.CommitSHA == "abc123def456789" &&
				req.CommitMessage == pushTestFeatureMessage &&
				req.AuthorName == pushTestUserName &&
				req.AuthorEmail == pushTestUserEmail &&
				req.BuildURL == pushTestCommitURL
		})).Return(expectedBuildEvent, nil).Once()

		// Setup mock expectations for notification service
		mockNotificationService.On("CreateNotificationForBuildEvent",
			mock.Anything,
			expectedBuildEvent.ID(),
			projectID,
			mock.AnythingOfType("string")).
			Return(expectedNotifications, nil).Once()

		// Create webhook request
		webhookRequest := dto.ProcessWebhookRequest{
			ProjectID:  projectID,
			EventType:  domain.PushEvent,
			Signature:  pushTestSignature,
			DeliveryID: pushTestDeliveryID,
			Body:       []byte(`{"action": "push"}`),
			Payload:    pushPayload,
		}

		// Act - Process the webhook
		result, err := webhookService.ProcessWebhook(context.Background(), webhookRequest)

		// Assert - Verify the result (This will fail initially - RED phase)
		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, projectID, result.ProjectID())
		assert.Equal(t, domain.PushEvent, result.EventType())
		assert.True(t, result.IsProcessed())

		// Verify all mocks were called as expected
		mockProjectService.AssertExpectations(t)
		mockSignatureVerifier.AssertExpectations(t)
		mockWebhookRepo.AssertExpectations(t)
		mockBuildService.AssertExpectations(t)
		mockNotificationService.AssertExpectations(t)
	})
}

func TestProcessPullRequestEventTDD(t *testing.T) {
	// Test data setup
	projectID := value_objects.NewID()

	pullRequestPayload := dto.GitHubActionsPayload{
		Action: "opened",
		Number: 123,
		Repository: struct {
			Name     string `json:"name"`
			FullName string `json:"full_name"`
			HTMLURL  string `json:"html_url"`
		}{
			Name:     "test-repo",
			FullName: "test/repo",
			HTMLURL:  pushTestRepoURL,
		},
		Sender: struct {
			Login string `json:"login"`
			Email string `json:"email,omitempty"`
		}{
			Login: "testuser",
			Email: pushTestUserEmail,
		},
		PullRequest: &dto.PullRequest{
			ID:       12345,
			Number:   123,
			State:    "open",
			Title:    pushTestFeatureMessage,
			Body:     "This PR adds a new feature to the application",
			HTMLURL:  pushTestPullRequestURL,
			DiffURL:  "https://github.com/test/repo/pull/123.diff",
			PatchURL: "https://github.com/test/repo/pull/123.patch",
			User: dto.User{
				Login: "testuser",
				Email: pushTestUserEmail,
				Name:  pushTestUserName,
			},
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
			Merged:    false,
			Head: dto.Branch{
				Label: "testuser:feature-branch",
				Ref:   pushTestFeatureBranch,
				SHA:   "abc123def456789",
				User: dto.User{
					Login: "testuser",
					Email: pushTestUserEmail,
					Name:  pushTestUserName,
				},
			},
			Base: dto.Branch{
				Label: "test:main",
				Ref:   "main",
				SHA:   "def456abc123456",
				User: dto.User{
					Login: "test",
				},
			},
		},
	}

	t.Run("TDD_Red_ProcessPullRequestEvent_should_create_build_event_and_notification", func(t *testing.T) {
		// Arrange - Setup mocks
		mockWebhookRepo := &mocks.MockWebhookEventRepository{}
		mockProjectService := &MockProjectServiceTDD{}
		mockBuildService := &MockBuildEventServiceTDD{}
		mockNotificationService := &MockNotificationLogServiceTDD{}
		mockSignatureVerifier := &mocks.MockSignatureVerifier{}

		// Setup webhook service
		webhookService := service.NewWebhookService(service.Dep{
			WebhookEventRepo:       mockWebhookRepo,
			ProjectService:         mockProjectService,
			BuildService:           mockBuildService,
			NotificationLogService: mockNotificationService,
			SignatureVerifier:      mockSignatureVerifier,
		})

		// Create expected project
		expectedProject, err := projectDomain.NewProject(pushTestProjectName, pushTestRepoURL, pushTestWebhookSecret, nil)
		require.NoError(t, err)

		// Create expected build event
		expectedBuildEvent, err := buildDomain.NewBuildEvent(buildDomain.BuildEventParams{
			ProjectID:     projectID,
			EventType:     buildDomain.EventTypePullRequest,
			Status:        buildDomain.BuildStatusPending,
			Branch:        pushTestFeatureBranch,
			CommitSHA:     "abc123def456789",
			CommitMessage: "Pull Request: " + pushTestFeatureMessage,
			AuthorName:    pushTestUserName,
			AuthorEmail:   pushTestUserEmail,
			BuildURL:      pushTestPullRequestURL,
		})
		require.NoError(t, err)

		// Create expected notification logs
		expectedNotifications := []*notificationDomain.NotificationLog{}
		notificationLog, err := notificationDomain.NewNotificationLog(
			expectedBuildEvent.ID(),
			projectID,
			notificationDomain.NotificationChannelTelegram,
			"123456789",
			"Test pull request notification message",
			3,
		)
		require.NoError(t, err)
		expectedNotifications = append(expectedNotifications, notificationLog)

		// Setup mock expectations for project service
		mockProjectService.On("GetProject", mock.Anything, projectID).
			Return(expectedProject, nil).Once()

		// Setup mock expectations for signature verification
		mockSignatureVerifier.On("VerifySignature", pushTestWebhookSecret, pushTestSignature, mock.Anything).
			Return(true).Once()

		// Setup mock expectations for duplicate check
		mockWebhookRepo.On("ExistsByDeliveryID", mock.Anything, pushTestDeliveryID).
			Return(false, nil).Once()

		// Setup mock expectations for webhook creation
		mockWebhookRepo.On("Create", mock.Anything, mock.AnythingOfType(pushWebhookEventType)).
			Return(nil).Once()

		// Setup mock expectations for webhook update
		mockWebhookRepo.On("Update", mock.Anything, mock.AnythingOfType(pushWebhookEventType)).
			Return(nil).Once()

		// Setup mock expectations for build service
		mockBuildService.On("CreateBuildEvent", mock.Anything, mock.MatchedBy(func(req buildDto.CreateBuildEventRequest) bool {
			return req.ProjectID == projectID &&
				req.EventType == buildDomain.EventTypePullRequest &&
				req.Status == buildDomain.BuildStatusPending &&
				req.Branch == pushTestFeatureBranch &&
				req.CommitSHA == "abc123def456789" &&
				req.CommitMessage == "Pull Request: "+pushTestFeatureMessage &&
				req.AuthorName == pushTestUserName &&
				req.AuthorEmail == pushTestUserEmail &&
				req.BuildURL == pushTestPullRequestURL
		})).Return(expectedBuildEvent, nil).Once()

		// Setup mock expectations for notification service
		mockNotificationService.On("CreateNotificationForBuildEvent",
			mock.Anything,
			expectedBuildEvent.ID(),
			projectID,
			mock.AnythingOfType("string")).
			Return(expectedNotifications, nil).Once()

		// Create webhook request
		webhookRequest := dto.ProcessWebhookRequest{
			ProjectID:  projectID,
			EventType:  domain.PullRequestEvent,
			Signature:  pushTestSignature,
			DeliveryID: pushTestDeliveryID,
			Body:       []byte(`{"action": "opened"}`),
			Payload:    pullRequestPayload,
		}

		// Act - Process the webhook
		result, err := webhookService.ProcessWebhook(context.Background(), webhookRequest)

		// Assert - Verify the result (This will fail initially - RED phase)
		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, projectID, result.ProjectID())
		assert.Equal(t, domain.PullRequestEvent, result.EventType())
		assert.True(t, result.IsProcessed())

		// Verify all mocks were called as expected
		mockProjectService.AssertExpectations(t)
		mockSignatureVerifier.AssertExpectations(t)
		mockWebhookRepo.AssertExpectations(t)
		mockBuildService.AssertExpectations(t)
		mockNotificationService.AssertExpectations(t)
	})
}
