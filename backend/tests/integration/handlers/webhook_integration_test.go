package handlers_test

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/dewisartika8/cicd-status-notifier-bot/internal/adapter/handler/webhook"
	"github.com/dewisartika8/cicd-status-notifier-bot/internal/core/shared/domain/value_objects"
	"github.com/dewisartika8/cicd-status-notifier-bot/internal/core/webhook/domain"
	"github.com/dewisartika8/cicd-status-notifier-bot/internal/core/webhook/dto"
	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// Mock webhook service for testing
type MockWebhookService struct {
	mock.Mock
}

func (m *MockWebhookService) ProcessWebhook(ctx context.Context, req dto.ProcessWebhookRequest) (*domain.WebhookEvent, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.WebhookEvent), args.Error(1)
}

func (m *MockWebhookService) VerifyWebhookSignature(secret, signature string, body []byte) bool {
	args := m.Called(secret, signature, body)
	return args.Bool(0)
}

func (m *MockWebhookService) GetWebhookEvent(ctx context.Context, id value_objects.ID) (*domain.WebhookEvent, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.WebhookEvent), args.Error(1)
}

func (m *MockWebhookService) GetWebhookEventsByProject(ctx context.Context, projectID value_objects.ID, limit, offset int) ([]*domain.WebhookEvent, error) {
	args := m.Called(ctx, projectID, limit, offset)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*domain.WebhookEvent), args.Error(1)
}

func (m *MockWebhookService) ReprocessFailedWebhooks(ctx context.Context, limit int) error {
	args := m.Called(ctx, limit)
	return args.Error(0)
}

func TestWebhookEndpointIntegration(t *testing.T) {
	// Create a new Fiber app
	app := fiber.New()

	// Create mock service and logger
	mockService := &MockWebhookService{}
	logger := logrus.New()

	// Create webhook handler
	webhookHandler := webhook.NewWebhookHandler(mockService, logger)

	// Setup routes following the health handler pattern
	api := app.Group("/api/v1")
	webhooks := api.Group("/webhooks")
	webhooks.Post("/github/:projectId", webhookHandler.ProcessGitHubWebhook)
	webhooks.Get("/events/:projectId", webhookHandler.GetWebhookEvents)
	webhooks.Get("/events/:projectId/:eventId", webhookHandler.GetWebhookEvent)

	// Test cases
	tests := []struct {
		name           string
		projectID      string
		headers        map[string]string
		payload        dto.GitHubActionsPayload
		expectedStatus int
		expectedBody   string
	}{
		{
			name:      "Missing signature header",
			projectID: "550e8400-e29b-41d4-a716-446655440000",
			headers: map[string]string{
				"X-GitHub-Event": "workflow_run",
			},
			payload: dto.GitHubActionsPayload{
				Action:   "completed",
				Workflow: "CI",
			},
			expectedStatus: http.StatusUnauthorized,
			expectedBody:   "missing X-Hub-Signature-256 header",
		},
		{
			name:      "Missing event type header",
			projectID: "550e8400-e29b-41d4-a716-446655440000",
			headers: map[string]string{
				"X-Hub-Signature-256": "sha256=test",
			},
			payload: dto.GitHubActionsPayload{
				Action:   "completed",
				Workflow: "CI",
			},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   "missing X-GitHub-Event header",
		},
		{
			name:      "Invalid project ID",
			projectID: "invalid-uuid",
			headers: map[string]string{
				"X-Hub-Signature-256": "sha256=test",
				"X-GitHub-Event":      "workflow_run",
				"X-GitHub-Delivery":   "delivery-123",
			},
			payload: dto.GitHubActionsPayload{
				Action:   "completed",
				Workflow: "CI",
			},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   "invalid project_id format",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create request body
			body, _ := json.Marshal(tt.payload)

			// Create HTTP request
			req := httptest.NewRequest("POST", "/api/v1/webhooks/github/"+tt.projectID, bytes.NewReader(body))
			req.Header.Set("Content-Type", "application/json")

			// Set headers
			for key, value := range tt.headers {
				req.Header.Set(key, value)
			}

			// Perform request
			resp, err := app.Test(req)
			assert.NoError(t, err)

			// Check status code
			assert.Equal(t, tt.expectedStatus, resp.StatusCode)

			// Check response body if expected
			if tt.expectedBody != "" {
				var response map[string]interface{}
				err := json.NewDecoder(resp.Body).Decode(&response)
				assert.NoError(t, err)
				assert.Contains(t, response["error"], tt.expectedBody)
			}
		})
	}
}

func TestWebhookEndpointHealthCheck(t *testing.T) {
	// Simple test to verify the application structure is correct
	logger := logrus.New()
	mockService := &MockWebhookService{}

	// Should not panic when creating webhook handler
	assert.NotPanics(t, func() {
		webhook.NewWebhookHandler(mockService, logger)
	})
}
