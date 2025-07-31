package mocks

import (
	"context"

	"github.com/dewisartika8/cicd-status-notifier-bot/internal/core/notification/domain"
	"github.com/dewisartika8/cicd-status-notifier-bot/internal/core/shared/domain/value_objects"
	"github.com/stretchr/testify/mock"
)

// NotificationLogRepository is a mock of port.NotificationLogRepository interface
type NotificationLogRepository struct {
	mock.Mock
}

// NewNotificationLogRepository creates a new mock instance
func NewNotificationLogRepository(t mock.TestingT) *NotificationLogRepository {
	mock := &NotificationLogRepository{}
	mock.Test(t)
	return mock
}

// Create provides a mock function with given fields: ctx, log
func (m *NotificationLogRepository) Create(ctx context.Context, log *domain.NotificationLog) error {
	ret := m.Called(ctx, log)
	return ret.Error(0)
}

// GetByID provides a mock function with given fields: ctx, id
func (m *NotificationLogRepository) GetByID(ctx context.Context, id value_objects.ID) (*domain.NotificationLog, error) {
	ret := m.Called(ctx, id)

	var r0 *domain.NotificationLog
	if rf, ok := ret.Get(0).(func(context.Context, value_objects.ID) *domain.NotificationLog); ok {
		r0 = rf(ctx, id)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*domain.NotificationLog)
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

// GetByBuildEventID provides a mock function with given fields: ctx, buildEventID
func (m *NotificationLogRepository) GetByBuildEventID(ctx context.Context, buildEventID value_objects.ID) ([]*domain.NotificationLog, error) {
	ret := m.Called(ctx, buildEventID)

	var r0 []*domain.NotificationLog
	if rf, ok := ret.Get(0).(func(context.Context, value_objects.ID) []*domain.NotificationLog); ok {
		r0 = rf(ctx, buildEventID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*domain.NotificationLog)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, value_objects.ID) error); ok {
		r1 = rf(ctx, buildEventID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetByProjectID provides a mock function with given fields: ctx, projectID, limit, offset
func (m *NotificationLogRepository) GetByProjectID(ctx context.Context, projectID value_objects.ID, limit, offset int) ([]*domain.NotificationLog, error) {
	ret := m.Called(ctx, projectID, limit, offset)

	var r0 []*domain.NotificationLog
	if rf, ok := ret.Get(0).(func(context.Context, value_objects.ID, int, int) []*domain.NotificationLog); ok {
		r0 = rf(ctx, projectID, limit, offset)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*domain.NotificationLog)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, value_objects.ID, int, int) error); ok {
		r1 = rf(ctx, projectID, limit, offset)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetByRecipient provides a mock function with given fields: ctx, recipient, limit, offset
func (m *NotificationLogRepository) GetByRecipient(ctx context.Context, recipient string, limit, offset int) ([]*domain.NotificationLog, error) {
	ret := m.Called(ctx, recipient, limit, offset)

	var r0 []*domain.NotificationLog
	if rf, ok := ret.Get(0).(func(context.Context, string, int, int) []*domain.NotificationLog); ok {
		r0 = rf(ctx, recipient, limit, offset)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*domain.NotificationLog)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, string, int, int) error); ok {
		r1 = rf(ctx, recipient, limit, offset)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Update provides a mock function with given fields: ctx, log
func (m *NotificationLogRepository) Update(ctx context.Context, log *domain.NotificationLog) error {
	ret := m.Called(ctx, log)
	return ret.Error(0)
}

// Delete provides a mock function with given fields: ctx, id
func (m *NotificationLogRepository) Delete(ctx context.Context, id value_objects.ID) error {
	ret := m.Called(ctx, id)
	return ret.Error(0)
}

// GetFailedNotifications provides a mock function with given fields: ctx, limit
func (m *NotificationLogRepository) GetFailedNotifications(ctx context.Context, limit int) ([]*domain.NotificationLog, error) {
	ret := m.Called(ctx, limit)

	var r0 []*domain.NotificationLog
	if rf, ok := ret.Get(0).(func(context.Context, int) []*domain.NotificationLog); ok {
		r0 = rf(ctx, limit)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*domain.NotificationLog)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, int) error); ok {
		r1 = rf(ctx, limit)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetPendingNotifications provides a mock function with given fields: ctx, limit
func (m *NotificationLogRepository) GetPendingNotifications(ctx context.Context, limit int) ([]*domain.NotificationLog, error) {
	ret := m.Called(ctx, limit)

	var r0 []*domain.NotificationLog
	if rf, ok := ret.Get(0).(func(context.Context, int) []*domain.NotificationLog); ok {
		r0 = rf(ctx, limit)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*domain.NotificationLog)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, int) error); ok {
		r1 = rf(ctx, limit)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Count provides a mock function with given fields: ctx, projectID, status
func (m *NotificationLogRepository) Count(ctx context.Context, projectID *value_objects.ID, status *domain.NotificationStatus) (int64, error) {
	ret := m.Called(ctx, projectID, status)

	var r0 int64
	if rf, ok := ret.Get(0).(func(context.Context, *value_objects.ID, *domain.NotificationStatus) int64); ok {
		r0 = rf(ctx, projectID, status)
	} else {
		r0 = ret.Get(0).(int64)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, *value_objects.ID, *domain.NotificationStatus) error); ok {
		r1 = rf(ctx, projectID, status)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetNotificationStats provides a mock function with given fields: ctx, projectID
func (m *NotificationLogRepository) GetNotificationStats(ctx context.Context, projectID value_objects.ID) (*domain.NotificationStats, error) {
	ret := m.Called(ctx, projectID)

	var r0 *domain.NotificationStats
	if rf, ok := ret.Get(0).(func(context.Context, value_objects.ID) *domain.NotificationStats); ok {
		r0 = rf(ctx, projectID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*domain.NotificationStats)
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
