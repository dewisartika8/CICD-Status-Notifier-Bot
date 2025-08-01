package webhook_test

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/dewisartika8/cicd-status-notifier-bot/internal/adapter/handler/webhook"
	buildDomain "github.com/dewisartika8/cicd-status-notifier-bot/internal/core/build/domain"
	buildDto "github.com/dewisartika8/cicd-status-notifier-bot/internal/core/build/dto"
	notificationDomain "github.com/dewisartika8/cicd-status-notifier-bot/internal/core/notification/domain"
	projectDomain "github.com/dewisartika8/cicd-status-notifier-bot/internal/core/project/domain"
	projectDto "github.com/dewisartika8/cicd-status-notifier-bot/internal/core/project/dto"
	"github.com/dewisartika8/cicd-status-notifier-bot/internal/core/shared/domain/value_objects"
	webhookDomain "github.com/dewisartika8/cicd-status-notifier-bot/internal/core/webhook/domain"
	"github.com/dewisartika8/cicd-status-notifier-bot/internal/core/webhook/dto"
	"github.com/dewisartika8/cicd-status-notifier-bot/internal/core/webhook/service"
	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

// Test constants
const (
	testProjectName           = "Test Repository"
	testRepoURL               = "https://github.com/testowner/testrepo"
	testWebhookSecret         = "test-webhook-secret-123"
	testSignature             = "sha256=test-signature-hash"
	testDeliveryID            = "test-delivery-12345"
	testCommitMessage         = "Add new feature"
	testAuthorName            = "John Doe"
	testAuthorEmail           = "john.doe@example.com"
	testBuildURL              = "https://github.com/testowner/testrepo/commit/abc123def456789"
	testRef                   = "refs/heads/main"
	testCompareURL            = "https://github.com/testowner/testrepo/compare/000000...abc123"
	testFullName              = "testowner/testrepo"
	testPushToMainMessage     = "Push to main"
	testAPIV1                 = "/api/v1"
	testWebhooksPath          = "/webhooks"
	testWebhookEventDomain    = "*domain.WebhookEvent"
	testWebhookGitHubEndpoint = "/api/v1/webhooks/github/"
	testContentType           = "Content-Type"
	testApplicationJSON       = "application/json"
	testHubSignature256       = "X-Hub-Signature-256"
	testGitHubEvent           = "X-GitHub-Event"
	testGitHubDelivery        = "X-GitHub-Delivery"
	testPushEvent             = "push"
)

// Mock implementations for E2E testing
type MockWebhookEventRepository struct {
	mock.Mock
}

func (m *MockWebhookEventRepository) Create(ctx context.Context, event *webhookDomain.WebhookEvent) error {
	// Create webhook event
	args := m.Called(ctx, event)
	return args.Error(0)
}

func (m *MockWebhookEventRepository) Update(ctx context.Context, event *webhookDomain.WebhookEvent) error {
	// Update webhook event
	args := m.Called(ctx, event)
	return args.Error(0)
}

func (m *MockWebhookEventRepository) Delete(ctx context.Context, id value_objects.ID) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockWebhookEventRepository) GetByID(ctx context.Context, id value_objects.ID) (*webhookDomain.WebhookEvent, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*webhookDomain.WebhookEvent), args.Error(1)
}

func (m *MockWebhookEventRepository) GetByDeliveryID(ctx context.Context, deliveryID string) (*webhookDomain.WebhookEvent, error) {
	args := m.Called(ctx, deliveryID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*webhookDomain.WebhookEvent), args.Error(1)
}

func (m *MockWebhookEventRepository) ExistsByDeliveryID(ctx context.Context, deliveryID string) (bool, error) {
	args := m.Called(ctx, deliveryID)
	return args.Bool(0), args.Error(1)
}

func (m *MockWebhookEventRepository) GetByProjectID(ctx context.Context, projectID value_objects.ID, limit, offset int) ([]*webhookDomain.WebhookEvent, error) {
	args := m.Called(ctx, projectID, limit, offset)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*webhookDomain.WebhookEvent), args.Error(1)
}

func (m *MockWebhookEventRepository) GetUnprocessedEvents(ctx context.Context, limit int) ([]*webhookDomain.WebhookEvent, error) {
	args := m.Called(ctx, limit)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*webhookDomain.WebhookEvent), args.Error(1)
}

type MockProjectService struct {
	mock.Mock
}

func (m *MockProjectService) GetProject(ctx context.Context, id value_objects.ID) (*projectDomain.Project, error) {
	// Get project by ID
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*projectDomain.Project), args.Error(1)
}

func (m *MockProjectService) CreateProject(ctx context.Context, req projectDto.CreateProjectRequest) (*projectDomain.Project, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*projectDomain.Project), args.Error(1)
}

func (m *MockProjectService) UpdateProject(ctx context.Context, id value_objects.ID, req projectDto.UpdateProjectRequest) (*projectDomain.Project, error) {
	args := m.Called(ctx, id, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*projectDomain.Project), args.Error(1)
}

func (m *MockProjectService) UpdateProjectStatus(ctx context.Context, id value_objects.ID, status projectDomain.ProjectStatus) (*projectDomain.Project, error) {
	args := m.Called(ctx, id, status)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*projectDomain.Project), args.Error(1)
}

func (m *MockProjectService) DeleteProject(ctx context.Context, id value_objects.ID) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockProjectService) ListProjects(ctx context.Context, filters projectDto.ListProjectFilters) ([]*projectDomain.Project, error) {
	args := m.Called(ctx, filters)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*projectDomain.Project), args.Error(1)
}

func (m *MockProjectService) GetActiveProjects(ctx context.Context) ([]*projectDomain.Project, error) {
	// Get active projects
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*projectDomain.Project), args.Error(1)
}

func (m *MockProjectService) GetProjectsWithTelegramChat(ctx context.Context) ([]*projectDomain.Project, error) {
	// Get projects with telegram chat configuration
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*projectDomain.Project), args.Error(1)
}

func (m *MockProjectService) ActivateProject(ctx context.Context, id value_objects.ID) (*projectDomain.Project, error) {
	// Activate project by ID
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*projectDomain.Project), args.Error(1)
}

func (m *MockProjectService) DeactivateProject(ctx context.Context, id value_objects.ID) (*projectDomain.Project, error) {
	// Deactivate project by ID
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*projectDomain.Project), args.Error(1)
}

func (m *MockProjectService) ArchiveProject(ctx context.Context, id value_objects.ID) (*projectDomain.Project, error) {
	// Archive project by ID
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*projectDomain.Project), args.Error(1)
}

func (m *MockProjectService) ValidateWebhookSecret(ctx context.Context, id value_objects.ID, secret string) (bool, error) {
	args := m.Called(ctx, id, secret)
	return args.Bool(0), args.Error(1)
}

func (m *MockProjectService) CountProjects(ctx context.Context, filters projectDto.ListProjectFilters) (int64, error) {
	args := m.Called(ctx, filters)
	return args.Get(0).(int64), args.Error(1)
}

func (m *MockProjectService) GetProjectByName(ctx context.Context, name string) (*projectDomain.Project, error) {
	args := m.Called(ctx, name)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*projectDomain.Project), args.Error(1)
}

func (m *MockProjectService) GetProjectByRepositoryURL(ctx context.Context, repositoryURL string) (*projectDomain.Project, error) {
	args := m.Called(ctx, repositoryURL)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*projectDomain.Project), args.Error(1)
}

type MockBuildEventService struct {
	mock.Mock
}

func (m *MockBuildEventService) CreateBuildEvent(ctx context.Context, req buildDto.CreateBuildEventRequest) (*buildDomain.BuildEvent, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*buildDomain.BuildEvent), args.Error(1)
}

func (m *MockBuildEventService) GetBuildEvent(ctx context.Context, id value_objects.ID) (*buildDomain.BuildEvent, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*buildDomain.BuildEvent), args.Error(1)
}

func (m *MockBuildEventService) GetBuildEventsByProject(ctx context.Context, projectID value_objects.ID, filters buildDto.ListBuildEventFilters) ([]*buildDomain.BuildEvent, error) {
	args := m.Called(ctx, projectID, filters)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*buildDomain.BuildEvent), args.Error(1)
}

func (m *MockBuildEventService) GetBuildMetrics(ctx context.Context, projectID value_objects.ID) (*buildDomain.BuildMetrics, error) {
	args := m.Called(ctx, projectID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*buildDomain.BuildMetrics), args.Error(1)
}

func (m *MockBuildEventService) GetLatestBuildEvent(ctx context.Context, projectID value_objects.ID) (*buildDomain.BuildEvent, error) {
	args := m.Called(ctx, projectID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*buildDomain.BuildEvent), args.Error(1)
}

func (m *MockBuildEventService) ListBuildEvents(ctx context.Context, filters buildDto.ListBuildEventFilters) ([]*buildDomain.BuildEvent, error) {
	args := m.Called(ctx, filters)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*buildDomain.BuildEvent), args.Error(1)
}

func (m *MockBuildEventService) ProcessWebhookEvent(ctx context.Context, req buildDto.ProcessWebhookRequest) ([]*buildDomain.BuildEvent, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*buildDomain.BuildEvent), args.Error(1)
}

func (m *MockBuildEventService) UpdateBuildEventStatus(ctx context.Context, id value_objects.ID, status buildDomain.BuildStatus, exitCode *int) error {
	args := m.Called(ctx, id, status, exitCode)
	return args.Error(0)
}

type MockNotificationLogService struct {
	mock.Mock
}

func (m *MockNotificationLogService) CreateNotificationForBuildEvent(ctx context.Context, buildEventID, projectID value_objects.ID, message string) ([]*notificationDomain.NotificationLog, error) {
	args := m.Called(ctx, buildEventID, projectID, message)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*notificationDomain.NotificationLog), args.Error(1)
}

func (m *MockNotificationLogService) CreateNotificationLog(ctx context.Context, buildEventID, projectID value_objects.ID, channel notificationDomain.NotificationChannel, recipient, message string) (*notificationDomain.NotificationLog, error) {
	args := m.Called(ctx, buildEventID, projectID, channel, recipient, message)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*notificationDomain.NotificationLog), args.Error(1)
}

func (m *MockNotificationLogService) GetNotificationLog(ctx context.Context, id value_objects.ID) (*notificationDomain.NotificationLog, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*notificationDomain.NotificationLog), args.Error(1)
}

func (m *MockNotificationLogService) GetNotificationLogsByBuildEvent(ctx context.Context, buildEventID value_objects.ID) ([]*notificationDomain.NotificationLog, error) {
	args := m.Called(ctx, buildEventID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*notificationDomain.NotificationLog), args.Error(1)
}

func (m *MockNotificationLogService) GetNotificationLogsByProject(ctx context.Context, projectID value_objects.ID, limit, offset int) ([]*notificationDomain.NotificationLog, error) {
	args := m.Called(ctx, projectID, limit, offset)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*notificationDomain.NotificationLog), args.Error(1)
}

func (m *MockNotificationLogService) GetNotificationStats(ctx context.Context, projectID value_objects.ID) (map[notificationDomain.NotificationStatus]int64, error) {
	args := m.Called(ctx, projectID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(map[notificationDomain.NotificationStatus]int64), args.Error(1)
}

func (m *MockNotificationLogService) ProcessFailedNotifications(ctx context.Context, limit int) error {
	args := m.Called(ctx, limit)
	return args.Error(0)
}

func (m *MockNotificationLogService) ProcessPendingNotifications(ctx context.Context, limit int) error {
	// Process pending notifications
	args := m.Called(ctx, limit)
	return args.Error(0)
}

func (m *MockNotificationLogService) RetryFailedNotification(ctx context.Context, notificationLogID value_objects.ID) error {
	args := m.Called(ctx, notificationLogID)
	return args.Error(0)
}

func (m *MockNotificationLogService) SendNotification(ctx context.Context, notificationLogID value_objects.ID) error {
	// Send notification
	args := m.Called(ctx, notificationLogID)
	return args.Error(0)
}

func (m *MockNotificationLogService) UpdateNotificationStatus(ctx context.Context, id value_objects.ID, status notificationDomain.NotificationStatus, errorMessage string, messageID *string) error {
	args := m.Called(ctx, id, status, errorMessage, messageID)
	return args.Error(0)
}

type MockSignatureVerifier struct {
	mock.Mock
}

func (m *MockSignatureVerifier) VerifySignature(secret, signature string, body []byte) bool {
	args := m.Called(secret, signature, body)
	return args.Bool(0)
}

// Test suite for E2E webhook push scenarios
func TestWebhookPushE2E(t *testing.T) {
	t.Run("GitHub_Push_Event_With_Complete_Payload_Should_Success", func(t *testing.T) {
		// Setup mocks
		mockWebhookRepo := &MockWebhookEventRepository{}
		mockProjectService := &MockProjectService{}
		mockBuildService := &MockBuildEventService{}
		mockNotificationService := &MockNotificationLogService{}
		mockSignatureVerifier := &MockSignatureVerifier{}

		// Create webhook service
		webhookService := service.NewWebhookService(service.Dep{
			WebhookEventRepo:       mockWebhookRepo,
			ProjectService:         mockProjectService,
			BuildService:           mockBuildService,
			NotificationLogService: mockNotificationService,
			SignatureVerifier:      mockSignatureVerifier,
		})

		// Create Fiber app and handler
		app := fiber.New()
		logger := logrus.New()
		webhookHandler := webhook.NewWebhookHandler(webhookService, logger)

		// Register routes
		api := app.Group(testAPIV1)
		webhooks := api.Group(testWebhooksPath)
		webhookHandler.RegisterRoutes(webhooks)

		// Test data
		projectID := value_objects.NewID()

		// Create expected project
		expectedProject, err := projectDomain.NewProject(testProjectName, testRepoURL, testWebhookSecret, nil)
		require.NoError(t, err)

		// Create expected build event
		expectedBuildEvent, err := buildDomain.NewBuildEvent(buildDomain.BuildEventParams{
			ProjectID:     projectID,
			EventType:     buildDomain.EventTypePush,
			Status:        buildDomain.BuildStatusSuccess,
			Branch:        "main",
			CommitSHA:     "abc123def456789",
			CommitMessage: testCommitMessage,
			AuthorName:    testAuthorName,
			AuthorEmail:   testAuthorEmail,
			BuildURL:      testBuildURL,
		})
		require.NoError(t, err)

		// Complete GitHub push payload (realistic scenario)
		pushPayload := dto.GitHubActionsPayload{
			Ref:     testRef,
			Before:  "0000000000000000000000000000000000000000",
			After:   "abc123def456789",
			Created: false,
			Deleted: false,
			Forced:  false,
			Compare: testCompareURL,
			Repository: struct {
				Name     string `json:"name"`
				FullName string `json:"full_name"`
				HTMLURL  string `json:"html_url"`
			}{
				Name:     "testrepo",
				FullName: testFullName,
				HTMLURL:  testRepoURL,
			},
			Sender: struct {
				Login string `json:"login"`
				Email string `json:"email,omitempty"`
			}{
				Login: "johndoe",
				Email: "john.doe@example.com",
			},
			HeadCommit: &dto.Commit{
				ID:        "abc123def456789",
				Message:   testCommitMessage,
				Timestamp: time.Now().Format(time.RFC3339),
				URL:       testBuildURL,
				Author: dto.User{
					Login: "johndoe",
					Email: testAuthorEmail,
					Name:  testAuthorName,
				},
				Committer: dto.User{
					Login: "johndoe",
					Email: testAuthorEmail,
					Name:  testAuthorName,
				},
			},
			Commits: []dto.Commit{
				{
					ID:        "abc123def456789",
					Message:   testCommitMessage,
					Timestamp: time.Now().Format(time.RFC3339),
					Author: dto.User{
						Login: "johndoe",
						Email: testAuthorEmail,
						Name:  testAuthorName,
					},
					Committer: dto.User{
						Login: "johndoe",
						Email: testAuthorEmail,
						Name:  testAuthorName,
					},
				},
			},
			Pusher: &dto.User{
				Login: "johndoe",
				Email: testAuthorEmail,
				Name:  testAuthorName,
			},
		}

		// Setup mock expectations
		mockProjectService.On("GetProject", mock.Anything, projectID).Return(expectedProject, nil).Once()
		mockSignatureVerifier.On("VerifySignature", testWebhookSecret, testSignature, mock.Anything).Return(true).Once()
		mockWebhookRepo.On("ExistsByDeliveryID", mock.Anything, testDeliveryID).Return(false, nil).Once()
		mockWebhookRepo.On("Create", mock.Anything, mock.AnythingOfType(testWebhookEventDomain)).Return(nil).Once()
		mockWebhookRepo.On("Update", mock.Anything, mock.AnythingOfType(testWebhookEventDomain)).Return(nil).Once()

		mockBuildService.On("CreateBuildEvent", mock.Anything, mock.MatchedBy(func(req buildDto.CreateBuildEventRequest) bool {
			return req.ProjectID == projectID &&
				req.EventType == buildDomain.EventTypePush &&
				req.Status == buildDomain.BuildStatusSuccess &&
				req.Branch == "main" &&
				req.CommitSHA == "abc123def456789" &&
				req.CommitMessage == testCommitMessage
		})).Return(expectedBuildEvent, nil).Once()

		mockNotificationService.On("CreateNotificationForBuildEvent", mock.Anything, expectedBuildEvent.ID(), projectID, mock.AnythingOfType("string")).Return([]*notificationDomain.NotificationLog{}, nil).Once()

		// Prepare request
		payloadBytes, err := json.Marshal(pushPayload)
		require.NoError(t, err)

		req := httptest.NewRequest("POST", testWebhookGitHubEndpoint+projectID.String(), bytes.NewReader(payloadBytes))
		req.Header.Set(testContentType, testApplicationJSON)
		req.Header.Set(testHubSignature256, testSignature)
		req.Header.Set(testGitHubEvent, "push")
		req.Header.Set(testGitHubDelivery, testDeliveryID)

		// Execute request
		resp, err := app.Test(req, -1)
		require.NoError(t, err)
		defer resp.Body.Close()

		// Verify response
		assert.Equal(t, http.StatusAccepted, resp.StatusCode)

		// Verify all mocks were called
		mockProjectService.AssertExpectations(t)
		mockSignatureVerifier.AssertExpectations(t)
		mockWebhookRepo.AssertExpectations(t)
		mockBuildService.AssertExpectations(t)
		mockNotificationService.AssertExpectations(t)
	})

	t.Run("GitHub_Push_Event_With_Nil_HeadCommit_Should_Handle_Gracefully", func(t *testing.T) {
		// This test reproduces the nil pointer error scenario
		// Setup mocks
		mockWebhookRepo := &MockWebhookEventRepository{}
		mockProjectService := &MockProjectService{}
		mockBuildService := &MockBuildEventService{}
		mockNotificationService := &MockNotificationLogService{}
		mockSignatureVerifier := &MockSignatureVerifier{}

		webhookService := service.NewWebhookService(service.Dep{
			WebhookEventRepo:       mockWebhookRepo,
			ProjectService:         mockProjectService,
			BuildService:           mockBuildService,
			NotificationLogService: mockNotificationService,
			SignatureVerifier:      mockSignatureVerifier,
		})

		app := fiber.New()
		logger := logrus.New()
		webhookHandler := webhook.NewWebhookHandler(webhookService, logger)

		api := app.Group(testAPIV1)
		webhooks := api.Group(testWebhooksPath)
		webhookHandler.RegisterRoutes(webhooks)

		projectID := value_objects.NewID()

		expectedProject, err := projectDomain.NewProject(testProjectName, testRepoURL, testWebhookSecret, nil)
		require.NoError(t, err)

		expectedBuildEvent, err := buildDomain.NewBuildEvent(buildDomain.BuildEventParams{
			ProjectID:     projectID,
			EventType:     buildDomain.EventTypePush,
			Status:        buildDomain.BuildStatusSuccess,
			Branch:        "main",
			CommitSHA:     "abc123def456789",
			CommitMessage: testCommitMessage,
			AuthorName:    testAuthorName,
			AuthorEmail:   testAuthorEmail,
			BuildURL:      testBuildURL,
		})
		require.NoError(t, err)

		// Push payload with nil HeadCommit but with Commits array (common GitHub scenario)
		pushPayload := dto.GitHubActionsPayload{
			Ref:     testRef,
			Before:  "0000000000000000000000000000000000000000",
			After:   "abc123def456789",
			Created: false,
			Deleted: false,
			Forced:  false,
			Compare: testCompareURL,
			Repository: struct {
				Name     string `json:"name"`
				FullName string `json:"full_name"`
				HTMLURL  string `json:"html_url"`
			}{
				Name:     "testrepo",
				FullName: testFullName,
				HTMLURL:  testRepoURL,
			},
			Sender: struct {
				Login string `json:"login"`
				Email string `json:"email,omitempty"`
			}{
				Login: "johndoe",
				Email: testAuthorEmail,
			},
			HeadCommit: nil, // This is nil - potential source of nil pointer error
			Commits: []dto.Commit{
				{
					ID:        "abc123def456789",
					Message:   testCommitMessage,
					Timestamp: time.Now().Format(time.RFC3339),
					Author: dto.User{
						Login: "johndoe",
						Email: testAuthorEmail,
						Name:  testAuthorName,
					},
					Committer: dto.User{
						Login: "johndoe",
						Email: testAuthorEmail,
						Name:  testAuthorName,
					},
				},
			},
			Pusher: &dto.User{
				Login: "johndoe",
				Email: testAuthorEmail,
				Name:  testAuthorName,
			},
		}

		// Setup mock expectations
		mockProjectService.On("GetProject", mock.Anything, projectID).Return(expectedProject, nil).Once()
		mockSignatureVerifier.On("VerifySignature", testWebhookSecret, testSignature, mock.Anything).Return(true).Once()
		mockWebhookRepo.On("ExistsByDeliveryID", mock.Anything, testDeliveryID).Return(false, nil).Once()
		mockWebhookRepo.On("Create", mock.Anything, mock.AnythingOfType(testWebhookEventDomain)).Return(nil).Once()
		mockWebhookRepo.On("Update", mock.Anything, mock.AnythingOfType(testWebhookEventDomain)).Return(nil).Once()

		mockBuildService.On("CreateBuildEvent", mock.Anything, mock.MatchedBy(func(req buildDto.CreateBuildEventRequest) bool {
			return req.ProjectID == projectID &&
				req.EventType == buildDomain.EventTypePush &&
				req.Status == buildDomain.BuildStatusSuccess &&
				req.Branch == "main" &&
				req.CommitSHA == "abc123def456789" &&
				req.CommitMessage == testCommitMessage
		})).Return(expectedBuildEvent, nil).Once()

		mockNotificationService.On("CreateNotificationForBuildEvent", mock.Anything, expectedBuildEvent.ID(), projectID, mock.AnythingOfType("string")).Return([]*notificationDomain.NotificationLog{}, nil).Once()

		// Prepare request
		payloadBytes, err := json.Marshal(pushPayload)
		require.NoError(t, err)

		req := httptest.NewRequest("POST", testWebhookGitHubEndpoint+projectID.String(), bytes.NewReader(payloadBytes))
		req.Header.Set(testContentType, testApplicationJSON)
		req.Header.Set(testHubSignature256, testSignature)
		req.Header.Set(testGitHubEvent, "push")
		req.Header.Set(testGitHubDelivery, testDeliveryID)

		// Execute request
		resp, err := app.Test(req, -1)
		require.NoError(t, err)
		defer resp.Body.Close()

		// This should not panic and should handle the nil HeadCommit gracefully
		assert.Equal(t, http.StatusAccepted, resp.StatusCode)

		mockProjectService.AssertExpectations(t)
		mockSignatureVerifier.AssertExpectations(t)
		mockWebhookRepo.AssertExpectations(t)
		mockBuildService.AssertExpectations(t)
		mockNotificationService.AssertExpectations(t)
	})

	t.Run("GitHub_Push_Event_With_Nil_HeadCommit_And_Empty_Commits_Should_Handle_Gracefully", func(t *testing.T) {
		// This test covers the most extreme case where both HeadCommit and Commits are nil/empty
		mockWebhookRepo := &MockWebhookEventRepository{}
		mockProjectService := &MockProjectService{}
		mockBuildService := &MockBuildEventService{}
		mockNotificationService := &MockNotificationLogService{}
		mockSignatureVerifier := &MockSignatureVerifier{}

		webhookService := service.NewWebhookService(service.Dep{
			WebhookEventRepo:       mockWebhookRepo,
			ProjectService:         mockProjectService,
			BuildService:           mockBuildService,
			NotificationLogService: mockNotificationService,
			SignatureVerifier:      mockSignatureVerifier,
		})

		app := fiber.New()
		logger := logrus.New()
		webhookHandler := webhook.NewWebhookHandler(webhookService, logger)

		api := app.Group(testAPIV1)
		webhooks := api.Group(testWebhooksPath)
		webhookHandler.RegisterRoutes(webhooks)

		projectID := value_objects.NewID()

		expectedProject, err := projectDomain.NewProject(testProjectName, testRepoURL, testWebhookSecret, nil)
		require.NoError(t, err)

		expectedBuildEvent, err := buildDomain.NewBuildEvent(buildDomain.BuildEventParams{
			ProjectID:     projectID,
			EventType:     buildDomain.EventTypePush,
			Status:        buildDomain.BuildStatusSuccess,
			Branch:        "main",
			CommitSHA:     "abc123def456789",
			CommitMessage: testPushToMainMessage,
			AuthorName:    "johndoe",
			AuthorEmail:   testAuthorEmail,
			BuildURL:      testBuildURL,
		})
		require.NoError(t, err)

		// Push payload with minimal data - both HeadCommit and Commits are nil/empty
		pushPayload := dto.GitHubActionsPayload{
			Ref:     testRef,
			Before:  "0000000000000000000000000000000000000000",
			After:   "abc123def456789",
			Created: false,
			Deleted: false,
			Forced:  false,
			Compare: testCompareURL,
			Repository: struct {
				Name     string `json:"name"`
				FullName string `json:"full_name"`
				HTMLURL  string `json:"html_url"`
			}{
				Name:     "testrepo",
				FullName: testFullName,
				HTMLURL:  testRepoURL,
			},
			Sender: struct {
				Login string `json:"login"`
				Email string `json:"email,omitempty"`
			}{
				Login: "johndoe",
				Email: testAuthorEmail,
			},
			HeadCommit: nil,            // nil
			Commits:    []dto.Commit{}, // empty
			Pusher: &dto.User{
				Login: "johndoe",
				Email: testAuthorEmail,
				Name:  testAuthorName,
			},
		}

		// Setup mock expectations
		mockProjectService.On("GetProject", mock.Anything, projectID).Return(expectedProject, nil).Once()
		mockSignatureVerifier.On("VerifySignature", testWebhookSecret, testSignature, mock.Anything).Return(true).Once()
		mockWebhookRepo.On("ExistsByDeliveryID", mock.Anything, testDeliveryID).Return(false, nil).Once()
		mockWebhookRepo.On("Create", mock.Anything, mock.AnythingOfType(testWebhookEventDomain)).Return(nil).Once()
		mockWebhookRepo.On("Update", mock.Anything, mock.AnythingOfType(testWebhookEventDomain)).Return(nil).Once()

		mockBuildService.On("CreateBuildEvent", mock.Anything, mock.MatchedBy(func(req buildDto.CreateBuildEventRequest) bool {
			return req.ProjectID == projectID &&
				req.EventType == buildDomain.EventTypePush &&
				req.Status == buildDomain.BuildStatusSuccess &&
				req.Branch == "main" &&
				req.CommitSHA == "abc123def456789" &&
				req.CommitMessage == testPushToMainMessage &&
				req.AuthorName == "johndoe" &&
				req.AuthorEmail == testAuthorEmail
		})).Return(expectedBuildEvent, nil).Once()

		mockNotificationService.On("CreateNotificationForBuildEvent", mock.Anything, expectedBuildEvent.ID(), projectID, mock.AnythingOfType("string")).Return([]*notificationDomain.NotificationLog{}, nil).Once()

		// Prepare request
		payloadBytes, err := json.Marshal(pushPayload)
		require.NoError(t, err)

		req := httptest.NewRequest("POST", testWebhookGitHubEndpoint+projectID.String(), bytes.NewReader(payloadBytes))
		req.Header.Set(testContentType, testApplicationJSON)
		req.Header.Set(testHubSignature256, testSignature)
		req.Header.Set(testGitHubEvent, "push")
		req.Header.Set(testGitHubDelivery, testDeliveryID)

		// Execute request
		resp, err := app.Test(req, -1)
		require.NoError(t, err)
		defer resp.Body.Close()

		// This should not panic and should handle the minimal payload gracefully
		assert.Equal(t, http.StatusAccepted, resp.StatusCode)

		mockProjectService.AssertExpectations(t)
		mockSignatureVerifier.AssertExpectations(t)
		mockWebhookRepo.AssertExpectations(t)
		mockBuildService.AssertExpectations(t)
		mockNotificationService.AssertExpectations(t)
	})

	t.Run("GitHub_Push_Event_With_Nil_Pusher_Should_Handle_Gracefully", func(t *testing.T) {
		// This test covers the case where Pusher is nil
		mockWebhookRepo := &MockWebhookEventRepository{}
		mockProjectService := &MockProjectService{}
		mockBuildService := &MockBuildEventService{}
		mockNotificationService := &MockNotificationLogService{}
		mockSignatureVerifier := &MockSignatureVerifier{}

		webhookService := service.NewWebhookService(service.Dep{
			WebhookEventRepo:       mockWebhookRepo,
			ProjectService:         mockProjectService,
			BuildService:           mockBuildService,
			NotificationLogService: mockNotificationService,
			SignatureVerifier:      mockSignatureVerifier,
		})

		app := fiber.New()
		logger := logrus.New()
		webhookHandler := webhook.NewWebhookHandler(webhookService, logger)

		api := app.Group(testAPIV1)
		webhooks := api.Group(testWebhooksPath)
		webhookHandler.RegisterRoutes(webhooks)

		projectID := value_objects.NewID()

		expectedProject, err := projectDomain.NewProject(testProjectName, testRepoURL, testWebhookSecret, nil)
		require.NoError(t, err)

		expectedBuildEvent, err := buildDomain.NewBuildEvent(buildDomain.BuildEventParams{
			ProjectID:     projectID,
			EventType:     buildDomain.EventTypePush,
			Status:        buildDomain.BuildStatusSuccess,
			Branch:        "main",
			CommitSHA:     "abc123def456789",
			CommitMessage: testPushToMainMessage,
			AuthorName:    "",
			AuthorEmail:   "",
			BuildURL:      testBuildURL,
		})
		require.NoError(t, err)

		// Push payload with nil Pusher
		pushPayload := dto.GitHubActionsPayload{
			Ref:     testRef,
			Before:  "0000000000000000000000000000000000000000",
			After:   "abc123def456789",
			Created: false,
			Deleted: false,
			Forced:  false,
			Compare: testCompareURL,
			Repository: struct {
				Name     string `json:"name"`
				FullName string `json:"full_name"`
				HTMLURL  string `json:"html_url"`
			}{
				Name:     "testrepo",
				FullName: testFullName,
				HTMLURL:  testRepoURL,
			},
			Sender: struct {
				Login string `json:"login"`
				Email string `json:"email,omitempty"`
			}{
				Login: "johndoe",
				Email: testAuthorEmail,
			},
			HeadCommit: nil,
			Commits:    []dto.Commit{},
			Pusher:     nil, // nil pusher
		}

		// Setup mock expectations
		mockProjectService.On("GetProject", mock.Anything, projectID).Return(expectedProject, nil).Once()
		mockSignatureVerifier.On("VerifySignature", testWebhookSecret, testSignature, mock.Anything).Return(true).Once()
		mockWebhookRepo.On("ExistsByDeliveryID", mock.Anything, testDeliveryID).Return(false, nil).Once()
		mockWebhookRepo.On("Create", mock.Anything, mock.AnythingOfType(testWebhookEventDomain)).Return(nil).Once()
		mockWebhookRepo.On("Update", mock.Anything, mock.AnythingOfType(testWebhookEventDomain)).Return(nil).Once()

		mockBuildService.On("CreateBuildEvent", mock.Anything, mock.MatchedBy(func(req buildDto.CreateBuildEventRequest) bool {
			return req.ProjectID == projectID &&
				req.EventType == buildDomain.EventTypePush &&
				req.Status == buildDomain.BuildStatusSuccess &&
				req.Branch == "main" &&
				req.CommitSHA == "abc123def456789" &&
				req.CommitMessage == testPushToMainMessage &&
				req.AuthorName == "" &&
				req.AuthorEmail == ""
		})).Return(expectedBuildEvent, nil).Once()

		mockNotificationService.On("CreateNotificationForBuildEvent", mock.Anything, expectedBuildEvent.ID(), projectID, mock.AnythingOfType("string")).Return([]*notificationDomain.NotificationLog{}, nil).Once()

		// Prepare request
		payloadBytes, err := json.Marshal(pushPayload)
		require.NoError(t, err)

		req := httptest.NewRequest("POST", testWebhookGitHubEndpoint+projectID.String(), bytes.NewReader(payloadBytes))
		req.Header.Set(testContentType, testApplicationJSON)
		req.Header.Set(testHubSignature256, testSignature)
		req.Header.Set(testGitHubEvent, "push")
		req.Header.Set(testGitHubDelivery, testDeliveryID)

		// Execute request
		resp, err := app.Test(req, -1)
		require.NoError(t, err)
		defer resp.Body.Close()

		// This should not panic and should handle the nil Pusher gracefully
		assert.Equal(t, http.StatusAccepted, resp.StatusCode)

		mockProjectService.AssertExpectations(t)
		mockSignatureVerifier.AssertExpectations(t)
		mockWebhookRepo.AssertExpectations(t)
		mockBuildService.AssertExpectations(t)
		mockNotificationService.AssertExpectations(t)
	})
}

// Benchmark test to measure performance
func BenchmarkWebhookPushProcessing(b *testing.B) {
	// Setup mocks
	mockWebhookRepo := &MockWebhookEventRepository{}
	mockProjectService := &MockProjectService{}
	mockBuildService := &MockBuildEventService{}
	mockNotificationService := &MockNotificationLogService{}
	mockSignatureVerifier := &MockSignatureVerifier{}

	webhookService := service.NewWebhookService(service.Dep{
		WebhookEventRepo:       mockWebhookRepo,
		ProjectService:         mockProjectService,
		BuildService:           mockBuildService,
		NotificationLogService: mockNotificationService,
		SignatureVerifier:      mockSignatureVerifier,
	})

	projectID := value_objects.NewID()
	expectedProject, _ := projectDomain.NewProject(testProjectName, testRepoURL, testWebhookSecret, nil)
	expectedBuildEvent, _ := buildDomain.NewBuildEvent(buildDomain.BuildEventParams{
		ProjectID:     projectID,
		EventType:     buildDomain.EventTypePush,
		Status:        buildDomain.BuildStatusSuccess,
		Branch:        "main",
		CommitSHA:     "abc123def456789",
		CommitMessage: "Test commit",
		AuthorName:    testAuthorName,
		AuthorEmail:   testAuthorEmail,
		BuildURL:      testBuildURL,
	})

	pushPayload := dto.GitHubActionsPayload{
		Ref:    testRef,
		Before: "0000000000000000000000000000000000000000",
		After:  "abc123def456789",
		Repository: struct {
			Name     string `json:"name"`
			FullName string `json:"full_name"`
			HTMLURL  string `json:"html_url"`
		}{
			Name:     "testrepo",
			FullName: testFullName,
			HTMLURL:  testRepoURL,
		},
		HeadCommit: &dto.Commit{
			ID:      "abc123def456789",
			Message: "Test commit",
			Author: dto.User{
				Name:  testAuthorName,
				Email: testAuthorEmail,
			},
		},
	}

	// Setup mock expectations that can be called multiple times
	mockProjectService.On("GetProject", mock.Anything, projectID).Return(expectedProject, nil)
	mockSignatureVerifier.On("VerifySignature", mock.Anything, mock.Anything, mock.Anything).Return(true)
	mockWebhookRepo.On("ExistsByDeliveryID", mock.Anything, mock.Anything).Return(false, nil)
	mockWebhookRepo.On("Create", mock.Anything, mock.Anything).Return(nil)
	mockWebhookRepo.On("Update", mock.Anything, mock.Anything).Return(nil)
	mockBuildService.On("CreateBuildEvent", mock.Anything, mock.Anything).Return(expectedBuildEvent, nil)
	mockNotificationService.On("CreateNotificationForBuildEvent", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return([]*notificationDomain.NotificationLog{}, nil)

	webhookRequest := dto.ProcessWebhookRequest{
		ProjectID:  projectID,
		EventType:  webhookDomain.PushEvent,
		Signature:  testSignature,
		DeliveryID: testDeliveryID,
		Body:       []byte(`{"ref": "` + testRef + `"}`),
		Payload:    pushPayload,
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := webhookService.ProcessWebhook(context.Background(), webhookRequest)
		if err != nil {
			b.Fatalf("Unexpected error: %v", err)
		}
	}
}

// Helper functions to reduce code duplication

func setupTestApp() (*fiber.App, *webhook.WebhookHandler) {
	app := fiber.New()
	logger := logrus.New()

	// Setup mocks
	mockWebhookRepo := &MockWebhookEventRepository{}
	mockProjectService := &MockProjectService{}
	mockBuildService := &MockBuildEventService{}
	mockNotificationService := &MockNotificationLogService{}
	mockSignatureVerifier := &MockSignatureVerifier{}

	webhookService := service.NewWebhookService(service.Dep{
		WebhookEventRepo:       mockWebhookRepo,
		ProjectService:         mockProjectService,
		BuildService:           mockBuildService,
		NotificationLogService: mockNotificationService,
		SignatureVerifier:      mockSignatureVerifier,
	})

	webhookHandler := webhook.NewWebhookHandler(webhookService, logger)

	// Register routes
	api := app.Group(testAPIV1)
	webhooks := api.Group(testWebhooksPath)
	webhookHandler.RegisterRoutes(webhooks)

	return app, webhookHandler
}

func createTestRequest(projectID value_objects.ID, payload interface{}) *http.Request {
	payloadBytes, _ := json.Marshal(payload)
	req := httptest.NewRequest("POST", testWebhookGitHubEndpoint+projectID.String(), bytes.NewReader(payloadBytes))
	req.Header.Set(testContentType, testApplicationJSON)
	req.Header.Set(testHubSignature256, testSignature)
	req.Header.Set(testGitHubEvent, testPushEvent)
	req.Header.Set(testGitHubDelivery, testDeliveryID)
	return req
}

func createBasicPushPayload() dto.GitHubActionsPayload {
	return dto.GitHubActionsPayload{
		Ref:     testRef,
		Before:  "0000000000000000000000000000000000000000",
		After:   "abc123def456789",
		Created: false,
		Deleted: false,
		Forced:  false,
		Compare: testCompareURL,
		Repository: struct {
			Name     string `json:"name"`
			FullName string `json:"full_name"`
			HTMLURL  string `json:"html_url"`
		}{
			Name:     "testrepo",
			FullName: testFullName,
			HTMLURL:  testRepoURL,
		},
		Sender: struct {
			Login string `json:"login"`
			Email string `json:"email,omitempty"`
		}{
			Login: "johndoe",
			Email: testAuthorEmail,
		},
	}
}
