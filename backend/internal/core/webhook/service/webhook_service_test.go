package service

import (
	"context"
	"testing"

	"github.com/dewisartika8/cicd-status-notifier-bot/internal/core/shared/domain/value_objects"
	"github.com/dewisartika8/cicd-status-notifier-bot/internal/core/webhook/domain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// Mock implementations for testing
type MockWebhookEventRepository struct {
	mock.Mock
}

func (m *MockWebhookEventRepository) Create(ctx context.Context, event *domain.WebhookEvent) error {
	args := m.Called(ctx, event)
	return args.Error(0)
}

func (m *MockWebhookEventRepository) GetByID(ctx context.Context, id value_objects.ID) (*domain.WebhookEvent, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(*domain.WebhookEvent), args.Error(1)
}

func (m *MockWebhookEventRepository) GetByDeliveryID(ctx context.Context, deliveryID string) (*domain.WebhookEvent, error) {
	args := m.Called(ctx, deliveryID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.WebhookEvent), args.Error(1)
}

func (m *MockWebhookEventRepository) GetByProjectID(ctx context.Context, projectID value_objects.ID, limit, offset int) ([]*domain.WebhookEvent, error) {
	args := m.Called(ctx, projectID, limit, offset)
	return args.Get(0).([]*domain.WebhookEvent), args.Error(1)
}

func (m *MockWebhookEventRepository) Update(ctx context.Context, event *domain.WebhookEvent) error {
	updateArgs := m.Called(ctx, event)
	return updateArgs.Error(0)
}

func (m *MockWebhookEventRepository) Delete(ctx context.Context, id value_objects.ID) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockWebhookEventRepository) ExistsByDeliveryID(ctx context.Context, deliveryID string) (bool, error) {
	args := m.Called(ctx, deliveryID)
	return args.Bool(0), args.Error(1)
}

func (m *MockWebhookEventRepository) GetUnprocessedEvents(ctx context.Context, limit int) ([]*domain.WebhookEvent, error) {
	args := m.Called(ctx, limit)
	return args.Get(0).([]*domain.WebhookEvent), args.Error(1)
}

type MockSignatureVerifier struct {
	mock.Mock
}

func (m *MockSignatureVerifier) VerifySignature(secret, signature string, body []byte) bool {
	args := m.Called(secret, signature, body)
	return args.Bool(0)
}

func TestWebhookServiceVerifySignature(t *testing.T) {
	mockVerifier := &MockSignatureVerifier{}
	service := &webhookService{
		Dep: Dep{
			SignatureVerifier: mockVerifier,
		},
	}

	secret := "test-secret"
	signature := "test-signature"
	body := []byte("test-body")

	mockVerifier.On("VerifySignature", secret, signature, body).Return(true)

	result := service.VerifyWebhookSignature(secret, signature, body)

	assert.True(t, result)
	mockVerifier.AssertExpectations(t)
}

func TestWebhookServiceGetWebhookEvent(t *testing.T) {
	mockRepo := &MockWebhookEventRepository{}
	service := &webhookService{
		Dep: Dep{
			WebhookEventRepo: mockRepo,
		},
	}

	eventID := value_objects.NewID()
	expectedEvent, _ := domain.NewWebhookEvent(
		value_objects.NewID(),
		domain.WorkflowRunEvent,
		`{"action": "completed"}`,
		"sha256=test",
		"delivery-123",
	)

	mockRepo.On("GetByID", mock.Anything, eventID).Return(expectedEvent, nil)

	result, err := service.GetWebhookEvent(context.Background(), eventID)

	assert.NoError(t, err)
	assert.Equal(t, expectedEvent, result)
	mockRepo.AssertExpectations(t)
}

func TestWebhookServiceGetWebhookEventsByProject(t *testing.T) {
	mockRepo := &MockWebhookEventRepository{}
	service := &webhookService{
		Dep: Dep{
			WebhookEventRepo: mockRepo,
		},
	}

	projectID := value_objects.NewID()
	limit := 10
	offset := 0
	expectedEvents := []*domain.WebhookEvent{}

	mockRepo.On("GetByProjectID", mock.Anything, projectID, limit, offset).Return(expectedEvents, nil)

	result, err := service.GetWebhookEventsByProject(context.Background(), projectID, limit, offset)

	assert.NoError(t, err)
	assert.Equal(t, expectedEvents, result)
	mockRepo.AssertExpectations(t)
}

func TestNewWebhookService(t *testing.T) {
	mockRepo := &MockWebhookEventRepository{}
	mockVerifier := &MockSignatureVerifier{}

	dep := Dep{
		WebhookEventRepo:  mockRepo,
		SignatureVerifier: mockVerifier,
	}

	service := NewWebhookService(dep)

	assert.NotNil(t, service)
}
