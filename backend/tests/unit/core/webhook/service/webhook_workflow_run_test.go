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
	projectDto "github.com/dewisartika8/cicd-status-notifier-bot/internal/core/project/dto"
	"github.com/dewisartika8/cicd-status-notifier-bot/internal/core/shared/domain/value_objects"
	"github.com/dewisartika8/cicd-status-notifier-bot/internal/core/webhook/domain"
	"github.com/dewisartika8/cicd-status-notifier-bot/internal/core/webhook/dto"
	"github.com/dewisartika8/cicd-status-notifier-bot/internal/core/webhook/service"
	"github.com/dewisartika8/cicd-status-notifier-bot/tests/mocks"
)

// Test constants
const (
	workflowTestSignature        = "test-signature"
	workflowTestDeliveryID       = "test-delivery-id"
	workflowTestWebhookSecret    = "test-webhook-secret"
	workflowTestProjectName      = "Test Project"
	workflowTestRepoURL          = "https://github.com/test/repo"
	workflowWebhookEventType     = "*domain.WebhookEvent"
	workflowBuildEventReqType    = "buildDto.CreateBuildEventRequest"
	workflowWebhookEventReturned = "Webhook event should be returned"
)

// Mock implementations for project service TDD testing
type MockProjectServiceTDD struct {
	mock.Mock
}

func (m *MockProjectServiceTDD) CreateProject(ctx context.Context, req projectDto.CreateProjectRequest) (*projectDomain.Project, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*projectDomain.Project), args.Error(1)
}

func (m *MockProjectServiceTDD) GetProject(ctx context.Context, id value_objects.ID) (*projectDomain.Project, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*projectDomain.Project), args.Error(1)
}

func (m *MockProjectServiceTDD) GetProjectByName(ctx context.Context, name string) (*projectDomain.Project, error) {
	args := m.Called(ctx, name)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*projectDomain.Project), args.Error(1)
}

func (m *MockProjectServiceTDD) GetProjectByRepositoryURL(ctx context.Context, repositoryURL string) (*projectDomain.Project, error) {
	args := m.Called(ctx, repositoryURL)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*projectDomain.Project), args.Error(1)
}

func (m *MockProjectServiceTDD) UpdateProject(ctx context.Context, id value_objects.ID, req projectDto.UpdateProjectRequest) (*projectDomain.Project, error) {
	args := m.Called(ctx, id, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*projectDomain.Project), args.Error(1)
}

func (m *MockProjectServiceTDD) UpdateProjectStatus(ctx context.Context, id value_objects.ID, status projectDomain.ProjectStatus) (*projectDomain.Project, error) {
	args := m.Called(ctx, id, status)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*projectDomain.Project), args.Error(1)
}

func (m *MockProjectServiceTDD) DeleteProject(ctx context.Context, id value_objects.ID) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockProjectServiceTDD) ListProjects(ctx context.Context, filters projectDto.ListProjectFilters) ([]*projectDomain.Project, error) {
	args := m.Called(ctx, filters)
	return args.Get(0).([]*projectDomain.Project), args.Error(1)
}

func (m *MockProjectServiceTDD) GetActiveProjects(ctx context.Context) ([]*projectDomain.Project, error) {
	args := m.Called(ctx)
	return args.Get(0).([]*projectDomain.Project), args.Error(1)
}

func (m *MockProjectServiceTDD) GetProjectsWithTelegramChat(ctx context.Context) ([]*projectDomain.Project, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*projectDomain.Project), args.Error(1)
}

func (m *MockProjectServiceTDD) ActivateProject(ctx context.Context, id value_objects.ID) (*projectDomain.Project, error) {
	args := m.Called(ctx, id)
	if project := args.Get(0); project != nil {
		return project.(*projectDomain.Project), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockProjectServiceTDD) DeactivateProject(ctx context.Context, id value_objects.ID) (*projectDomain.Project, error) {
	args := m.Called(ctx, id)
	project := args.Get(0)
	if project != nil {
		return project.(*projectDomain.Project), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockProjectServiceTDD) ArchiveProject(ctx context.Context, id value_objects.ID) (*projectDomain.Project, error) {
	args := m.Called(ctx, id)
	result := args.Get(0)
	if result != nil {
		return result.(*projectDomain.Project), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockProjectServiceTDD) ValidateWebhookSecret(ctx context.Context, id value_objects.ID, secret string) (bool, error) {
	args := m.Called(ctx, id, secret)
	return args.Bool(0), args.Error(1)
}

func (m *MockProjectServiceTDD) CountProjects(ctx context.Context, filters projectDto.ListProjectFilters) (int64, error) {
	args := m.Called(ctx, filters)
	return args.Get(0).(int64), args.Error(1)
}

type MockBuildEventServiceTDD struct {
	mock.Mock
}

func (m *MockBuildEventServiceTDD) CreateBuildEvent(ctx context.Context, req buildDto.CreateBuildEventRequest) (*buildDomain.BuildEvent, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*buildDomain.BuildEvent), args.Error(1)
}

func (m *MockBuildEventServiceTDD) ProcessWebhookEvent(ctx context.Context, req buildDto.ProcessWebhookRequest) ([]*buildDomain.BuildEvent, error) {
	args := m.Called(ctx, req)
	return args.Get(0).([]*buildDomain.BuildEvent), args.Error(1)
}

func (m *MockBuildEventServiceTDD) GetBuildEvent(ctx context.Context, id value_objects.ID) (*buildDomain.BuildEvent, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*buildDomain.BuildEvent), args.Error(1)
}

func (m *MockBuildEventServiceTDD) GetBuildEventsByProject(ctx context.Context, projectID value_objects.ID, filters buildDto.ListBuildEventFilters) ([]*buildDomain.BuildEvent, error) {
	args := m.Called(ctx, projectID, filters)
	return args.Get(0).([]*buildDomain.BuildEvent), args.Error(1)
}

func (m *MockBuildEventServiceTDD) UpdateBuildEventStatus(ctx context.Context, id value_objects.ID, status buildDomain.BuildStatus, duration *int) error {
	args := m.Called(ctx, id, status, duration)
	return args.Error(0)
}

func (m *MockBuildEventServiceTDD) GetLatestBuildEvent(ctx context.Context, projectID value_objects.ID) (*buildDomain.BuildEvent, error) {
	args := m.Called(ctx, projectID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*buildDomain.BuildEvent), args.Error(1)
}

func (m *MockBuildEventServiceTDD) GetBuildMetrics(ctx context.Context, projectID value_objects.ID) (*buildDomain.BuildMetrics, error) {
	args := m.Called(ctx, projectID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*buildDomain.BuildMetrics), args.Error(1)
}

func (m *MockBuildEventServiceTDD) ListBuildEvents(ctx context.Context, filters buildDto.ListBuildEventFilters) ([]*buildDomain.BuildEvent, error) {
	args := m.Called(ctx, filters)
	return args.Get(0).([]*buildDomain.BuildEvent), args.Error(1)
}

// Mock implementations for notification service TDD testing
type MockNotificationLogServiceTDD struct {
	mock.Mock
}

func (m *MockNotificationLogServiceTDD) CreateNotificationForBuildEvent(
	ctx context.Context,
	buildEventID, projectID value_objects.ID,
	message string,
) ([]*notificationDomain.NotificationLog, error) {
	args := m.Called(ctx, buildEventID, projectID, message)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*notificationDomain.NotificationLog), args.Error(1)
}

func (m *MockNotificationLogServiceTDD) CreateNotificationLog(
	ctx context.Context,
	buildEventID, projectID value_objects.ID,
	channel notificationDomain.NotificationChannel,
	recipient, message string,
) (*notificationDomain.NotificationLog, error) {
	args := m.Called(ctx, buildEventID, projectID, channel, recipient, message)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*notificationDomain.NotificationLog), args.Error(1)
}

func (m *MockNotificationLogServiceTDD) SendNotification(ctx context.Context, notificationLogID value_objects.ID) error {
	args := m.Called(ctx, notificationLogID)
	return args.Error(0)
}

func (m *MockNotificationLogServiceTDD) GetNotificationLog(ctx context.Context, id value_objects.ID) (*notificationDomain.NotificationLog, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*notificationDomain.NotificationLog), args.Error(1)
}

func (m *MockNotificationLogServiceTDD) GetNotificationLogsByBuildEvent(ctx context.Context, buildEventID value_objects.ID) ([]*notificationDomain.NotificationLog, error) {
	args := m.Called(ctx, buildEventID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*notificationDomain.NotificationLog), args.Error(1)
}

func (m *MockNotificationLogServiceTDD) GetNotificationLogsByProject(ctx context.Context, projectID value_objects.ID, limit, offset int) ([]*notificationDomain.NotificationLog, error) {
	args := m.Called(ctx, projectID, limit, offset)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*notificationDomain.NotificationLog), args.Error(1)
}

func (m *MockNotificationLogServiceTDD) UpdateNotificationStatus(
	ctx context.Context,
	id value_objects.ID,
	status notificationDomain.NotificationStatus,
	errorMessage string,
	messageID *string,
) error {
	args := m.Called(ctx, id, status, errorMessage, messageID)
	return args.Error(0)
}

func (m *MockNotificationLogServiceTDD) RetryFailedNotification(ctx context.Context, notificationLogID value_objects.ID) error {
	callArgs := m.Called(ctx, notificationLogID)
	return callArgs.Error(0)
}

func (m *MockNotificationLogServiceTDD) ProcessPendingNotifications(ctx context.Context, limit int) error {
	args := m.Called(ctx, limit)
	return args.Error(0)
}

func (m *MockNotificationLogServiceTDD) ProcessFailedNotifications(ctx context.Context, limit int) error {
	result := m.Called(ctx, limit)
	return result.Error(0)
}

func (m *MockNotificationLogServiceTDD) GetNotificationStats(ctx context.Context, projectID value_objects.ID) (map[notificationDomain.NotificationStatus]int64, error) {
	args := m.Called(ctx, projectID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(map[notificationDomain.NotificationStatus]int64), args.Error(1)
}

func TestWorkflowRunProcessingTDD(t *testing.T) {
	// Test data setup
	projectID := value_objects.NewID()
	const testBuildURL = "https://github.com/test/repo/actions/runs/123456789"
	const testCommitMessage = "Test commit"
	const testBuildEventReqType = "buildDto.CreateBuildEventRequest"

	workflowRunPayload := dto.GitHubActionsPayload{
		Action: "completed",
		WorkflowRun: &dto.WorkflowRun{
			ID:         123456789,
			Name:       "CI/CD Pipeline",
			Status:     "completed",
			Conclusion: "success",
			HTMLURL:    testBuildURL,
			CreatedAt:  time.Now().Add(-10 * time.Minute),
			UpdatedAt:  time.Now(),
			RunNumber:  42,
			Event:      "push",
			HeadBranch: "main",
			HeadSha:    "abc123def456",
		},
		Repository: struct {
			Name     string `json:"name"`
			FullName string `json:"full_name"`
			HTMLURL  string `json:"html_url"`
		}{
			Name:     "test-repo",
			FullName: "test/repo",
			HTMLURL:  workflowTestRepoURL,
		},
		Sender: struct {
			Login string `json:"login"`
			Email string `json:"email,omitempty"`
		}{
			Login: "testuser",
			Email: "",
		},
	}

	t.Run("TDD_successful_workflow_run_processing", func(t *testing.T) {
		// Arrange - Setup mocks
		mockWebhookRepo := &mocks.MockWebhookEventRepository{}
		mockProjectService := &MockProjectServiceTDD{}
		mockBuildService := &MockBuildEventServiceTDD{}
		mockNotificationService := &MockNotificationLogServiceTDD{}
		mockSignatureVerifier := &mocks.MockSignatureVerifier{}

		// Setup webhook service with all dependencies
		webhookService := service.NewWebhookService(service.Dep{
			WebhookEventRepo:       mockWebhookRepo,
			ProjectService:         mockProjectService,
			BuildService:           mockBuildService,
			NotificationLogService: mockNotificationService,
			SignatureVerifier:      mockSignatureVerifier,
		})

		// Create expected project
		expectedProject, err := projectDomain.NewProject(workflowTestProjectName, workflowTestRepoURL, workflowTestWebhookSecret, nil)
		require.NoError(t, err)

		// Create expected build event
		expectedBuildEvent, err := buildDomain.NewBuildEvent(buildDomain.BuildEventParams{
			ProjectID:     projectID,
			EventType:     buildDomain.EventTypeBuildCompleted,
			Status:        buildDomain.BuildStatusSuccess,
			Branch:        "main",
			CommitSHA:     "abc123def456",
			CommitMessage: testCommitMessage,
			AuthorName:    "testuser",
			AuthorEmail:   "",
			BuildURL:      testBuildURL,
		})
		require.NoError(t, err)

		// Create expected notification
		expectedNotification, err := notificationDomain.NewNotificationLog(
			expectedBuildEvent.ID(),
			projectID,
			notificationDomain.NotificationChannelTelegram,
			"123456789",
			"Test notification message",
			3,
		)
		require.NoError(t, err)

		// Setup mock expectations for project service
		mockProjectService.On("GetProject", mock.Anything, projectID).
			Return(expectedProject, nil).Once()

		// Setup mock expectations for signature verification
		mockSignatureVerifier.On("VerifySignature", workflowTestWebhookSecret, workflowTestSignature, mock.Anything).
			Return(true).Once()

		// Setup mock expectations for duplicate check
		mockWebhookRepo.On("ExistsByDeliveryID", mock.Anything, workflowTestDeliveryID).
			Return(false, nil).Once()

		// Setup mock expectations for webhook creation
		mockWebhookRepo.On("Create", mock.Anything, mock.AnythingOfType(workflowWebhookEventType)).
			Return(nil).Once()

		// Setup mock expectations for webhook update
		mockWebhookRepo.On("Update", mock.Anything, mock.AnythingOfType(workflowWebhookEventType)).
			Return(nil).Once()

		// Setup mock expectations for build service
		mockBuildService.On("CreateBuildEvent", mock.Anything, mock.MatchedBy(func(req buildDto.CreateBuildEventRequest) bool {
			return req.EventType == buildDomain.EventTypeBuildCompleted &&
				req.Status == buildDomain.BuildStatusSuccess &&
				req.Branch == "main" &&
				req.CommitSHA == "abc123def456"
		})).Return(expectedBuildEvent, nil).Once()

		// Setup mock expectations for notification service
		mockNotificationService.On("CreateNotificationForBuildEvent",
			mock.Anything,
			expectedBuildEvent.ID(),
			projectID,
			mock.AnythingOfType("string")).
			Return([]*notificationDomain.NotificationLog{expectedNotification}, nil).Once()

		// Setup mock expectations for sending notification immediately
		mockNotificationService.On("SendNotification",
			mock.Anything,
			expectedNotification.ID()).
			Return(nil).Once()

		// Create ProcessWebhookRequest
		processReq := dto.ProcessWebhookRequest{
			ProjectID:  projectID,
			EventType:  domain.WorkflowRunEvent,
			Payload:    workflowRunPayload,
			Signature:  workflowTestSignature,
			DeliveryID: workflowTestDeliveryID,
			Body:       []byte("test-body"),
		}

		// Act - Process webhook via public method
		webhookEvent, err := webhookService.ProcessWebhook(context.Background(), processReq)

		// Assert - Verify results
		assert.NoError(t, err, "Processing workflow run should not return error")
		assert.NotNil(t, webhookEvent, workflowWebhookEventReturned)
		assert.Equal(t, domain.WorkflowRunEvent, webhookEvent.EventType())

		// Verify all mocks were called as expected
		mockWebhookRepo.AssertExpectations(t)
		mockProjectService.AssertExpectations(t)
		mockSignatureVerifier.AssertExpectations(t)
		mockBuildService.AssertExpectations(t)
		mockNotificationService.AssertExpectations(t)
	})

	t.Run("TDD_workflow_run_with_invalid_payload", func(t *testing.T) {
		// Arrange
		mockWebhookRepo := &mocks.MockWebhookEventRepository{}
		mockProjectService := &MockProjectServiceTDD{}
		mockBuildService := &MockBuildEventServiceTDD{}
		mockNotificationService := &MockNotificationLogServiceTDD{}
		mockSignatureVerifier := &mocks.MockSignatureVerifier{}

		webhookService := service.NewWebhookService(service.Dep{
			WebhookEventRepo:       mockWebhookRepo,
			ProjectService:         mockProjectService,
			BuildService:           mockBuildService,
			NotificationLogService: mockNotificationService,
			SignatureVerifier:      mockSignatureVerifier,
		})

		// Create expected project
		expectedProject, err := projectDomain.NewProject(workflowTestProjectName, workflowTestRepoURL, workflowTestWebhookSecret, nil)
		require.NoError(t, err)

		// Setup mock expectations for project service
		mockProjectService.On("GetProject", mock.Anything, projectID).
			Return(expectedProject, nil).Once()

		// Setup mock expectations for signature verification
		mockSignatureVerifier.On("VerifySignature", workflowTestWebhookSecret, workflowTestSignature, mock.Anything).
			Return(true).Once()

		// Setup mock expectations for duplicate check
		mockWebhookRepo.On("ExistsByDeliveryID", mock.Anything, workflowTestDeliveryID).
			Return(false, nil).Once()

		// Setup mock expectations for webhook creation
		mockWebhookRepo.On("Create", mock.Anything, mock.AnythingOfType(workflowWebhookEventType)).
			Return(nil).Once()

		// Create payload without WorkflowRun data
		invalidPayload := dto.GitHubActionsPayload{
			Action:      "completed",
			WorkflowRun: nil, // Missing workflow run data
		}

		// Create ProcessWebhookRequest
		processReq := dto.ProcessWebhookRequest{
			ProjectID:  projectID,
			EventType:  domain.WorkflowRunEvent,
			Payload:    invalidPayload,
			Signature:  workflowTestSignature,
			DeliveryID: workflowTestDeliveryID,
		}

		// Act
		webhookEvent, err := webhookService.ProcessWebhook(context.Background(), processReq)

		// Assert
		// The webhook service stores the event even if processing fails, so it returns success
		assert.NoError(t, err, "Webhook should be stored successfully even with invalid payload")
		assert.NotNil(t, webhookEvent, workflowWebhookEventReturned)
		assert.Equal(t, domain.WorkflowRunEvent, webhookEvent.EventType())

		// Verify the webhook was stored but not marked as processed
		// (since processing failed due to invalid payload)
		assert.Nil(t, webhookEvent.ProcessedAt(), "Webhook should not be marked as processed due to invalid payload")

		mockWebhookRepo.AssertExpectations(t)
		mockProjectService.AssertExpectations(t)
		mockSignatureVerifier.AssertExpectations(t)
	})

	t.Run("TDD_workflow_run_with_build_service_failure", func(t *testing.T) {
		// Arrange
		mockWebhookRepo := &mocks.MockWebhookEventRepository{}
		mockProjectService := &MockProjectServiceTDD{}
		mockBuildService := &MockBuildEventServiceTDD{}
		mockNotificationService := &MockNotificationLogServiceTDD{}
		mockSignatureVerifier := &mocks.MockSignatureVerifier{}

		webhookService := service.NewWebhookService(service.Dep{
			WebhookEventRepo:       mockWebhookRepo,
			ProjectService:         mockProjectService,
			BuildService:           mockBuildService,
			NotificationLogService: mockNotificationService,
			SignatureVerifier:      mockSignatureVerifier,
		})

		// Create expected project
		expectedProject, err := projectDomain.NewProject(workflowTestProjectName, workflowTestRepoURL, workflowTestWebhookSecret, nil)
		require.NoError(t, err)

		// Setup mock expectations for project service
		mockProjectService.On("GetProject", mock.Anything, projectID).
			Return(expectedProject, nil).Once()

		// Setup mock expectations for signature verification
		mockSignatureVerifier.On("VerifySignature", workflowTestWebhookSecret, workflowTestSignature, mock.Anything).
			Return(true).Once()

		// Setup mock expectations for duplicate check
		mockWebhookRepo.On("ExistsByDeliveryID", mock.Anything, workflowTestDeliveryID).
			Return(false, nil).Once()

		// Setup mock expectations for webhook creation
		mockWebhookRepo.On("Create", mock.Anything, mock.AnythingOfType(workflowWebhookEventType)).
			Return(nil).Once()

		// Setup mock to return error for build service
		mockBuildService.On("CreateBuildEvent", mock.Anything, mock.MatchedBy(func(req buildDto.CreateBuildEventRequest) bool {
			return req.EventType == buildDomain.EventTypeBuildCompleted
		})).Return(nil, assert.AnError).Once()

		// Create ProcessWebhookRequest
		processReq := dto.ProcessWebhookRequest{
			ProjectID:  projectID,
			EventType:  domain.WorkflowRunEvent,
			Payload:    workflowRunPayload,
			Signature:  workflowTestSignature,
			DeliveryID: workflowTestDeliveryID,
		}

		// Act
		webhookEvent, err := webhookService.ProcessWebhook(context.Background(), processReq)

		// Assert
		// The webhook service stores the event even if processing fails, so it returns success
		assert.NoError(t, err, "Webhook should be stored successfully even when build service fails")
		assert.NotNil(t, webhookEvent, workflowWebhookEventReturned)
		assert.Equal(t, domain.WorkflowRunEvent, webhookEvent.EventType())

		// Verify the webhook was stored but not marked as processed
		// (since processing failed due to build service error)
		assert.Nil(t, webhookEvent.ProcessedAt(), "Webhook should not be marked as processed due to build service failure")

		mockWebhookRepo.AssertExpectations(t)
		mockProjectService.AssertExpectations(t)
		mockSignatureVerifier.AssertExpectations(t)
		mockBuildService.AssertExpectations(t)
	})

	t.Run("TDD_workflow_run_notification_message_format", func(t *testing.T) {
		// Arrange
		mockWebhookRepo := &mocks.MockWebhookEventRepository{}
		mockProjectService := &MockProjectServiceTDD{}
		mockBuildService := &MockBuildEventServiceTDD{}
		mockNotificationService := &MockNotificationLogServiceTDD{}
		mockSignatureVerifier := &mocks.MockSignatureVerifier{}

		webhookService := service.NewWebhookService(service.Dep{
			WebhookEventRepo:       mockWebhookRepo,
			ProjectService:         mockProjectService,
			BuildService:           mockBuildService,
			NotificationLogService: mockNotificationService,
			SignatureVerifier:      mockSignatureVerifier,
		})

		// Create expected project
		expectedProject, err := projectDomain.NewProject(workflowTestProjectName, workflowTestRepoURL, workflowTestWebhookSecret, nil)
		require.NoError(t, err)

		// Setup mock expectations for project service
		mockProjectService.On("GetProject", mock.Anything, projectID).
			Return(expectedProject, nil).Once()

		// Setup mock expectations for signature verification
		mockSignatureVerifier.On("VerifySignature", workflowTestWebhookSecret, workflowTestSignature, mock.Anything).
			Return(true).Once()

		// Create expected build event
		expectedBuildEvent, err := buildDomain.NewBuildEvent(buildDomain.BuildEventParams{
			ProjectID:     projectID,
			EventType:     buildDomain.EventTypeBuildCompleted,
			Status:        buildDomain.BuildStatusSuccess,
			Branch:        "main",
			CommitSHA:     "abc123def456",
			CommitMessage: testCommitMessage,
			AuthorName:    "testuser",
			AuthorEmail:   "",
			BuildURL:      testBuildURL,
		})
		require.NoError(t, err)

		// Setup mock expectations for duplicate check
		mockWebhookRepo.On("ExistsByDeliveryID", mock.Anything, workflowTestDeliveryID).
			Return(false, nil).Once()

		// Setup mock expectations for webhook creation
		mockWebhookRepo.On("Create", mock.Anything, mock.AnythingOfType(workflowWebhookEventType)).
			Return(nil).Once()

		// Setup mock expectations for webhook update
		mockWebhookRepo.On("Update", mock.Anything, mock.AnythingOfType(workflowWebhookEventType)).
			Return(nil).Once()

		// Setup mock for build event
		mockBuildService.On("CreateBuildEvent", mock.Anything, mock.MatchedBy(func(req buildDto.CreateBuildEventRequest) bool {
			return req.EventType == buildDomain.EventTypeBuildCompleted &&
				req.Status == buildDomain.BuildStatusSuccess
		})).Return(expectedBuildEvent, nil).Once()

		// Capture the notification message
		var capturedMessage string
		mockNotificationService.On("CreateNotificationForBuildEvent",
			mock.Anything,
			expectedBuildEvent.ID(),
			projectID,
			mock.MatchedBy(func(message string) bool {
				capturedMessage = message
				return true
			})).
			Return([]*notificationDomain.NotificationLog{}, nil).Once()

		// Create ProcessWebhookRequest
		processReq := dto.ProcessWebhookRequest{
			ProjectID:  projectID,
			EventType:  domain.WorkflowRunEvent,
			Payload:    workflowRunPayload,
			Signature:  workflowTestSignature,
			DeliveryID: workflowTestDeliveryID,
		}

		// Act
		webhookEvent, err := webhookService.ProcessWebhook(context.Background(), processReq)

		// Assert
		assert.NoError(t, err, "Processing should not return error")
		assert.NotNil(t, webhookEvent, workflowWebhookEventReturned)

		// Verify notification message format
		assert.Contains(t, capturedMessage, "CI/CD Pipeline", "Message should contain workflow name")
		assert.Contains(t, capturedMessage, "✅", "Message should contain success emoji for successful workflow")
		assert.Contains(t, capturedMessage, "main", "Message should contain branch")
		assert.Contains(t, capturedMessage, "test/repo", "Message should contain repository name")

		mockWebhookRepo.AssertExpectations(t)
		mockProjectService.AssertExpectations(t)
		mockSignatureVerifier.AssertExpectations(t)
		mockBuildService.AssertExpectations(t)
		mockNotificationService.AssertExpectations(t)
	})
}
