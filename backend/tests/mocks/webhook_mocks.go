package mocks

import (
	"context"

	"github.com/stretchr/testify/mock"

	"github.com/dewisartika8/cicd-status-notifier-bot/internal/core/shared/domain/value_objects"
	"github.com/dewisartika8/cicd-status-notifier-bot/internal/core/webhook/domain"
)

// MockWebhookEventRepository provides a mock implementation of webhook event repository
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

// MockSignatureVerifier provides a mock implementation of signature verifier
type MockSignatureVerifier struct {
	mock.Mock
}

func (m *MockSignatureVerifier) VerifySignature(secret, signature string, body []byte) bool {
	args := m.Called(secret, signature, body)
	return args.Bool(0)
}
