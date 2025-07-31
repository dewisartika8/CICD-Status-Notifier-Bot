package memory

import (
	"context"
	"errors"
	"sort"
	"sync"
	"time"

	"github.com/dewisartika8/cicd-status-notifier-bot/internal/core/notification/domain"
	"github.com/dewisartika8/cicd-status-notifier-bot/internal/core/notification/port"
	"github.com/dewisartika8/cicd-status-notifier-bot/internal/core/shared/domain/value_objects"
)

// inMemoryDeliveryQueueRepository implements DeliveryQueueRepository using in-memory storage
type inMemoryDeliveryQueueRepository struct {
	notifications map[value_objects.ID]*domain.QueuedNotification
	mutex         sync.RWMutex
}

// NewInMemoryDeliveryQueueRepository creates a new in-memory delivery queue repository
func NewInMemoryDeliveryQueueRepository() port.DeliveryQueueRepository {
	return &inMemoryDeliveryQueueRepository{
		notifications: make(map[value_objects.ID]*domain.QueuedNotification),
	}
}

func (r *inMemoryDeliveryQueueRepository) Create(ctx context.Context, notification *domain.QueuedNotification) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	r.notifications[notification.ID] = notification
	return nil
}

func (r *inMemoryDeliveryQueueRepository) GetByID(ctx context.Context, id value_objects.ID) (*domain.QueuedNotification, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	notification, exists := r.notifications[id]
	if !exists {
		return nil, ErrNotificationNotFound
	}

	return notification, nil
}

func (r *inMemoryDeliveryQueueRepository) GetPendingNotifications(ctx context.Context, limit int) ([]*domain.QueuedNotification, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	var pending []*domain.QueuedNotification
	count := 0

	for _, notification := range r.notifications {
		if notification.Status == domain.DeliveryStatusPending && count < limit {
			pending = append(pending, notification)
			count++
		}
	}

	return pending, nil
}

func (r *inMemoryDeliveryQueueRepository) GetPendingByPriority(ctx context.Context, limit int) ([]*domain.QueuedNotification, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	var pending []*domain.QueuedNotification

	for _, notification := range r.notifications {
		if notification.Status == domain.DeliveryStatusPending && notification.ShouldBeProcessed() {
			pending = append(pending, notification)
		}
	}

	// Sort by priority (higher priority first) and creation time
	sort.Slice(pending, func(i, j int) bool {
		if pending[i].Priority == pending[j].Priority {
			return pending[i].CreatedAt.Before(pending[j].CreatedAt)
		}
		return pending[i].Priority > pending[j].Priority
	})

	// Limit results
	if len(pending) > limit {
		pending = pending[:limit]
	}

	return pending, nil
}

func (r *inMemoryDeliveryQueueRepository) GetFailedNotifications(ctx context.Context, limit int) ([]*domain.QueuedNotification, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	var failed []*domain.QueuedNotification
	count := 0

	for _, notification := range r.notifications {
		if (notification.Status == domain.DeliveryStatusFailed || notification.Status == domain.DeliveryStatusRetrying) &&
			notification.IsRetryable() && count < limit {
			failed = append(failed, notification)
			count++
		}
	}

	return failed, nil
}

func (r *inMemoryDeliveryQueueRepository) Update(ctx context.Context, notification *domain.QueuedNotification) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	if _, exists := r.notifications[notification.ID]; !exists {
		return ErrNotificationNotFound
	}

	notification.UpdatedAt = time.Now()
	r.notifications[notification.ID] = notification
	return nil
}

func (r *inMemoryDeliveryQueueRepository) UpdateStatus(ctx context.Context, id value_objects.ID, status domain.DeliveryStatus, errorMessage string) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	notification, exists := r.notifications[id]
	if !exists {
		return ErrNotificationNotFound
	}

	notification.Status = status
	notification.LastError = errorMessage
	notification.UpdatedAt = time.Now()

	if status == domain.DeliveryStatusFailed {
		notification.AttemptCount++
	}

	return nil
}

func (r *inMemoryDeliveryQueueRepository) Delete(ctx context.Context, id value_objects.ID) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	if _, exists := r.notifications[id]; !exists {
		return ErrNotificationNotFound
	}

	delete(r.notifications, id)
	return nil
}

func (r *inMemoryDeliveryQueueRepository) DeleteProcessedNotifications(ctx context.Context, olderThan time.Duration) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	cutoff := time.Now().Add(-olderThan)

	for id, notification := range r.notifications {
		if notification.Status == domain.DeliveryStatusDelivered && notification.UpdatedAt.Before(cutoff) {
			delete(r.notifications, id)
		}
	}

	return nil
}

func (r *inMemoryDeliveryQueueRepository) GetPendingCount(ctx context.Context) (int64, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	count := int64(0)
	for _, notification := range r.notifications {
		if notification.Status == domain.DeliveryStatusPending {
			count++
		}
	}

	return count, nil
}

func (r *inMemoryDeliveryQueueRepository) GetQueueStats(ctx context.Context) (map[string]int64, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	stats := make(map[string]int64)

	for _, notification := range r.notifications {
		stats[string(notification.Status)]++
	}

	return stats, nil
}

var ErrNotificationNotFound = errors.New("notification not found")
