package notification

import (
	"context"

	"github.com/dewisartika8/cicd-status-notifier-bot/internal/core/notification/domain"
	"github.com/dewisartika8/cicd-status-notifier-bot/internal/core/shared/domain/value_objects"
	"github.com/stretchr/testify/mock"
)

// MockTelegramRepository is a consolidated mock for testing Telegram subscription repository
type MockTelegramRepository struct {
	mock.Mock
}

// Create mocks the Create method
func (m *MockTelegramRepository) Create(ctx context.Context, subscription *domain.TelegramSubscription) error {
	args := m.Called(ctx, subscription)
	return args.Error(0)
}

// GetByID mocks the GetByID method
func (m *MockTelegramRepository) GetByID(ctx context.Context, id value_objects.ID) (*domain.TelegramSubscription, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.TelegramSubscription), args.Error(1)
}

// GetByProjectID mocks the GetByProjectID method
func (m *MockTelegramRepository) GetByProjectID(ctx context.Context, projectID value_objects.ID) ([]*domain.TelegramSubscription, error) {
	args := m.Called(ctx, projectID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*domain.TelegramSubscription), args.Error(1)
}

// GetByChatID mocks the GetByChatID method
func (m *MockTelegramRepository) GetByChatID(ctx context.Context, chatID int64) (*domain.TelegramSubscription, error) {
	args := m.Called(ctx, chatID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.TelegramSubscription), args.Error(1)
}

// GetByProjectAndChatID mocks the GetByProjectAndChatID method
func (m *MockTelegramRepository) GetByProjectAndChatID(ctx context.Context, projectID value_objects.ID, chatID int64) (*domain.TelegramSubscription, error) {
	args := m.Called(ctx, projectID, chatID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.TelegramSubscription), args.Error(1)
}

// Update mocks the Update method
func (m *MockTelegramRepository) Update(ctx context.Context, subscription *domain.TelegramSubscription) error {
	args := m.Called(ctx, subscription)
	return args.Error(0)
}

// Delete mocks the Delete method
func (m *MockTelegramRepository) Delete(ctx context.Context, id value_objects.ID) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

// ExistsByProjectAndChatID mocks the ExistsByProjectAndChatID method
func (m *MockTelegramRepository) ExistsByProjectAndChatID(ctx context.Context, projectID value_objects.ID, chatID int64) (bool, error) {
	args := m.Called(ctx, projectID, chatID)
	return args.Bool(0), args.Error(1)
}

// GetActiveSubscriptions mocks the GetActiveSubscriptions method
func (m *MockTelegramRepository) GetActiveSubscriptions(ctx context.Context) ([]*domain.TelegramSubscription, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*domain.TelegramSubscription), args.Error(1)
}

// GetActiveSubscriptionsByProject mocks the GetActiveSubscriptionsByProject method
func (m *MockTelegramRepository) GetActiveSubscriptionsByProject(ctx context.Context, projectID value_objects.ID) ([]*domain.TelegramSubscription, error) {
	args := m.Called(ctx, projectID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*domain.TelegramSubscription), args.Error(1)
}

// Count mocks the Count method
func (m *MockTelegramRepository) Count(ctx context.Context, projectID *value_objects.ID, isActive *bool) (int64, error) {
	args := m.Called(ctx, projectID, isActive)
	return args.Get(0).(int64), args.Error(1)
}
