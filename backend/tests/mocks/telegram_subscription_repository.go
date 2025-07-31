package mocks

import (
	"context"

	"github.com/dewisartika8/cicd-status-notifier-bot/internal/core/notification/domain"
	"github.com/dewisartika8/cicd-status-notifier-bot/internal/core/shared/domain/value_objects"
	"github.com/stretchr/testify/mock"
)

// TelegramSubscriptionRepository is a mock of port.TelegramSubscriptionRepository interface
type TelegramSubscriptionRepository struct {
	mock.Mock
}

// NewTelegramSubscriptionRepository creates a new mock instance
func NewTelegramSubscriptionRepository(t mock.TestingT) *TelegramSubscriptionRepository {
	mock := &TelegramSubscriptionRepository{}
	mock.Test(t)
	return mock
}

// Create provides a mock function with given fields: ctx, subscription
func (m *TelegramSubscriptionRepository) Create(ctx context.Context, subscription *domain.TelegramSubscription) error {
	ret := m.Called(ctx, subscription)
	return ret.Error(0)
}

// GetByID provides a mock function with given fields: ctx, id
func (m *TelegramSubscriptionRepository) GetByID(ctx context.Context, id value_objects.ID) (*domain.TelegramSubscription, error) {
	ret := m.Called(ctx, id)

	var r0 *domain.TelegramSubscription
	if rf, ok := ret.Get(0).(func(context.Context, value_objects.ID) *domain.TelegramSubscription); ok {
		r0 = rf(ctx, id)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*domain.TelegramSubscription)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, value_objects.ID) error); ok {
		r1 = rf(ctx, id)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetByProjectID provides a mock function with given fields: ctx, projectID
func (m *TelegramSubscriptionRepository) GetByProjectID(ctx context.Context, projectID value_objects.ID) ([]*domain.TelegramSubscription, error) {
	ret := m.Called(ctx, projectID)

	var r0 []*domain.TelegramSubscription
	if rf, ok := ret.Get(0).(func(context.Context, value_objects.ID) []*domain.TelegramSubscription); ok {
		r0 = rf(ctx, projectID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*domain.TelegramSubscription)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, value_objects.ID) error); ok {
		r1 = rf(ctx, projectID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetByChatID provides a mock function with given fields: ctx, chatID
func (m *TelegramSubscriptionRepository) GetByChatID(ctx context.Context, chatID int64) (*domain.TelegramSubscription, error) {
	ret := m.Called(ctx, chatID)

	var r0 *domain.TelegramSubscription
	if rf, ok := ret.Get(0).(func(context.Context, int64) *domain.TelegramSubscription); ok {
		r0 = rf(ctx, chatID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*domain.TelegramSubscription)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, int64) error); ok {
		r1 = rf(ctx, chatID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetByProjectAndChatID provides a mock function with given fields: ctx, projectID, chatID
func (m *TelegramSubscriptionRepository) GetByProjectAndChatID(ctx context.Context, projectID value_objects.ID, chatID int64) (*domain.TelegramSubscription, error) {
	ret := m.Called(ctx, projectID, chatID)

	var r0 *domain.TelegramSubscription
	if rf, ok := ret.Get(0).(func(context.Context, value_objects.ID, int64) *domain.TelegramSubscription); ok {
		r0 = rf(ctx, projectID, chatID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*domain.TelegramSubscription)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, value_objects.ID, int64) error); ok {
		r1 = rf(ctx, projectID, chatID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Update provides a mock function with given fields: ctx, subscription
func (m *TelegramSubscriptionRepository) Update(ctx context.Context, subscription *domain.TelegramSubscription) error {
	ret := m.Called(ctx, subscription)
	return ret.Error(0)
}

// Delete provides a mock function with given fields: ctx, id
func (m *TelegramSubscriptionRepository) Delete(ctx context.Context, id value_objects.ID) error {
	ret := m.Called(ctx, id)
	return ret.Error(0)
}

// GetActiveSubscriptions provides a mock function with given fields: ctx
func (m *TelegramSubscriptionRepository) GetActiveSubscriptions(ctx context.Context) ([]*domain.TelegramSubscription, error) {
	ret := m.Called(ctx)

	var r0 []*domain.TelegramSubscription
	if rf, ok := ret.Get(0).(func(context.Context) []*domain.TelegramSubscription); ok {
		r0 = rf(ctx)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*domain.TelegramSubscription)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context) error); ok {
		r1 = rf(ctx)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetActiveSubscriptionsByProject provides a mock function with given fields: ctx, projectID
func (m *TelegramSubscriptionRepository) GetActiveSubscriptionsByProject(ctx context.Context, projectID value_objects.ID) ([]*domain.TelegramSubscription, error) {
	ret := m.Called(ctx, projectID)

	var r0 []*domain.TelegramSubscription
	if rf, ok := ret.Get(0).(func(context.Context, value_objects.ID) []*domain.TelegramSubscription); ok {
		r0 = rf(ctx, projectID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*domain.TelegramSubscription)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, value_objects.ID) error); ok {
		r1 = rf(ctx, projectID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// ExistsByProjectAndChatID provides a mock function with given fields: ctx, projectID, chatID
func (m *TelegramSubscriptionRepository) ExistsByProjectAndChatID(ctx context.Context, projectID value_objects.ID, chatID int64) (bool, error) {
	ret := m.Called(ctx, projectID, chatID)

	var r0 bool
	if rf, ok := ret.Get(0).(func(context.Context, value_objects.ID, int64) bool); ok {
		r0 = rf(ctx, projectID, chatID)
	} else {
		r0 = ret.Get(0).(bool)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, value_objects.ID, int64) error); ok {
		r1 = rf(ctx, projectID, chatID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Count provides a mock function with given fields: ctx, projectID, isActive
func (m *TelegramSubscriptionRepository) Count(ctx context.Context, projectID *value_objects.ID, isActive *bool) (int64, error) {
	ret := m.Called(ctx, projectID, isActive)

	var r0 int64
	if rf, ok := ret.Get(0).(func(context.Context, *value_objects.ID, *bool) int64); ok {
		r0 = rf(ctx, projectID, isActive)
	} else {
		r0 = ret.Get(0).(int64)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, *value_objects.ID, *bool) error); ok {
		r1 = rf(ctx, projectID, isActive)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}
