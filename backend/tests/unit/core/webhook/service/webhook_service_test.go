package service_test

import (
	"context"
	"testing"

	"github.com/dewisartika8/cicd-status-notifier-bot/internal/core/shared/domain/value_objects"
	"github.com/dewisartika8/cicd-status-notifier-bot/internal/core/webhook/domain"
	"github.com/dewisartika8/cicd-status-notifier-bot/internal/core/webhook/service"
	"github.com/dewisartika8/cicd-status-notifier-bot/tests/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestWebhookServiceVerifySignature(t *testing.T) {
	mockVerifier := &mocks.MockSignatureVerifier{}
	webhookService := service.NewWebhookService(service.Dep{
		SignatureVerifier: mockVerifier,
	})

	secret := "test-secret"
	signature := "test-signature"
	body := []byte("test-body")

	mockVerifier.On("VerifySignature", secret, signature, body).Return(true)

	result := webhookService.VerifyWebhookSignature(secret, signature, body)

	assert.True(t, result)
	mockVerifier.AssertExpectations(t)
}

func TestWebhookServiceGetWebhookEvent(t *testing.T) {
	mockRepo := &mocks.MockWebhookEventRepository{}
	webhookService := service.NewWebhookService(service.Dep{
		WebhookEventRepo: mockRepo,
	})

	eventID := value_objects.NewID()
	expectedEvent, _ := domain.NewWebhookEvent(
		value_objects.NewID(),
		domain.WorkflowRunEvent,
		`{"action": "completed"}`,
		"sha256=test",
		"delivery-123",
	)

	mockRepo.On("GetByID", mock.Anything, eventID).Return(expectedEvent, nil)

	result, err := webhookService.GetWebhookEvent(context.Background(), eventID)

	assert.NoError(t, err)
	assert.Equal(t, expectedEvent, result)
	mockRepo.AssertExpectations(t)
}

func TestWebhookServiceGetWebhookEventsByProject(t *testing.T) {
	mockRepo := &mocks.MockWebhookEventRepository{}
	webhookService := service.NewWebhookService(service.Dep{
		WebhookEventRepo: mockRepo,
	})

	projectID := value_objects.NewID()
	limit := 10
	offset := 0
	expectedEvents := []*domain.WebhookEvent{}

	mockRepo.On("GetByProjectID", mock.Anything, projectID, limit, offset).Return(expectedEvents, nil)

	result, err := webhookService.GetWebhookEventsByProject(context.Background(), projectID, limit, offset)

	assert.NoError(t, err)
	assert.Equal(t, expectedEvents, result)
	mockRepo.AssertExpectations(t)
}

func TestNewWebhookService(t *testing.T) {
	mockRepo := &mocks.MockWebhookEventRepository{}
	mockVerifier := &mocks.MockSignatureVerifier{}

	dep := service.Dep{
		WebhookEventRepo:  mockRepo,
		SignatureVerifier: mockVerifier,
	}

	webhookService := service.NewWebhookService(dep)

	assert.NotNil(t, webhookService)
}
