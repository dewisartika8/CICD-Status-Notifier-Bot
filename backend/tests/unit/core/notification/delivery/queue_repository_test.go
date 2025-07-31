package delivery

import (
	"context"
	"testing"
	"time"

	"github.com/dewisartika8/cicd-status-notifier-bot/internal/adapter/repository/memory"
	"github.com/dewisartika8/cicd-status-notifier-bot/internal/core/notification/domain"
	"github.com/dewisartika8/cicd-status-notifier-bot/internal/core/notification/port"
	"github.com/dewisartika8/cicd-status-notifier-bot/internal/core/shared/domain/value_objects"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type QueueRepositoryTestSuite struct {
	suite.Suite
	repo port.DeliveryQueueRepository
	ctx  context.Context
}

func (suite *QueueRepositoryTestSuite) SetupTest() {
	suite.repo = memory.NewInMemoryDeliveryQueueRepository()
	suite.ctx = context.Background()
}

func (suite *QueueRepositoryTestSuite) TestCreateAndGetByID() {
	// Arrange
	notification := domain.NewQueuedNotification(
		value_objects.NewID(),
		domain.NotificationChannelTelegram,
		"123456789",
		"Test message",
		"Test subject",
		1,
		3,
	)

	// Act
	err := suite.repo.Create(suite.ctx, notification)
	assert.NoError(suite.T(), err)

	retrieved, err := suite.repo.GetByID(suite.ctx, notification.ID)

	// Assert
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), notification.ID, retrieved.ID)
	assert.Equal(suite.T(), notification.Channel, retrieved.Channel)
	assert.Equal(suite.T(), notification.Recipient, retrieved.Recipient)
	assert.Equal(suite.T(), notification.Message, retrieved.Message)
}

func (suite *QueueRepositoryTestSuite) TestGetByIDNotFound() {
	// Act
	_, err := suite.repo.GetByID(suite.ctx, value_objects.NewID())

	// Assert
	assert.Error(suite.T(), err)
	assert.Equal(suite.T(), memory.ErrNotificationNotFound, err)
}

func (suite *QueueRepositoryTestSuite) TestGetPendingNotifications() {
	// Arrange
	notification1 := domain.NewQueuedNotification(
		value_objects.NewID(),
		domain.NotificationChannelTelegram,
		"123456789",
		"Test message 1",
		"Test subject",
		1,
		3,
	)
	notification2 := domain.NewQueuedNotification(
		value_objects.NewID(),
		domain.NotificationChannelEmail,
		"test@example.com",
		"Test message 2",
		"Test subject",
		2,
		3,
	)
	notification3 := domain.NewQueuedNotification(
		value_objects.NewID(),
		domain.NotificationChannelSlack,
		"channel1",
		"Test message 3",
		"Test subject",
		1,
		3,
	)
	notification3.Status = domain.DeliveryStatusDelivered // Not pending

	suite.repo.Create(suite.ctx, notification1)
	suite.repo.Create(suite.ctx, notification2)
	suite.repo.Create(suite.ctx, notification3)

	// Act
	pending, err := suite.repo.GetPendingNotifications(suite.ctx, 10)

	// Assert
	assert.NoError(suite.T(), err)
	assert.Len(suite.T(), pending, 2) // Only pending notifications
}

func (suite *QueueRepositoryTestSuite) TestGetPendingByPriority() {
	// Arrange
	lowPriority := domain.NewQueuedNotification(
		value_objects.NewID(),
		domain.NotificationChannelTelegram,
		"123456789",
		"Low priority",
		"Test subject",
		1, // Lower priority
		3,
	)
	highPriority := domain.NewQueuedNotification(
		value_objects.NewID(),
		domain.NotificationChannelEmail,
		"test@example.com",
		"High priority",
		"Test subject",
		5, // Higher priority
		3,
	)

	suite.repo.Create(suite.ctx, lowPriority)
	suite.repo.Create(suite.ctx, highPriority)

	// Act
	pending, err := suite.repo.GetPendingByPriority(suite.ctx, 10)

	// Assert
	assert.NoError(suite.T(), err)
	assert.Len(suite.T(), pending, 2)
	// High priority should come first
	assert.Equal(suite.T(), highPriority.ID, pending[0].ID)
	assert.Equal(suite.T(), lowPriority.ID, pending[1].ID)
}

func (suite *QueueRepositoryTestSuite) TestUpdateStatus() {
	// Arrange
	notification := domain.NewQueuedNotification(
		value_objects.NewID(),
		domain.NotificationChannelTelegram,
		"123456789",
		"Test message",
		"Test subject",
		1,
		3,
	)
	suite.repo.Create(suite.ctx, notification)

	// Act
	err := suite.repo.UpdateStatus(suite.ctx, notification.ID, domain.DeliveryStatusDelivered, "")

	// Assert
	assert.NoError(suite.T(), err)

	updated, err := suite.repo.GetByID(suite.ctx, notification.ID)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), domain.DeliveryStatusDelivered, updated.Status)
}

func (suite *QueueRepositoryTestSuite) TestUpdate() {
	// Arrange
	notification := domain.NewQueuedNotification(
		value_objects.NewID(),
		domain.NotificationChannelTelegram,
		"123456789",
		"Test message",
		"Test subject",
		1,
		3,
	)
	suite.repo.Create(suite.ctx, notification)

	// Modify notification
	notification.Status = domain.DeliveryStatusProcessing
	notification.Message = "Updated message"

	// Act
	err := suite.repo.Update(suite.ctx, notification)

	// Assert
	assert.NoError(suite.T(), err)

	updated, err := suite.repo.GetByID(suite.ctx, notification.ID)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), domain.DeliveryStatusProcessing, updated.Status)
	assert.Equal(suite.T(), "Updated message", updated.Message)
}

func (suite *QueueRepositoryTestSuite) TestGetFailedNotifications() {
	// Arrange
	failedNotification := domain.NewQueuedNotification(
		value_objects.NewID(),
		domain.NotificationChannelTelegram,
		"123456789",
		"Failed message",
		"Test subject",
		1,
		3,
	)
	failedNotification.MarkAsFailed("Test error")

	pendingNotification := domain.NewQueuedNotification(
		value_objects.NewID(),
		domain.NotificationChannelEmail,
		"test@example.com",
		"Pending message",
		"Test subject",
		1,
		3,
	)

	suite.repo.Create(suite.ctx, failedNotification)
	suite.repo.Create(suite.ctx, pendingNotification)

	// Act
	failed, err := suite.repo.GetFailedNotifications(suite.ctx, 10)

	// Assert
	assert.NoError(suite.T(), err)
	assert.Len(suite.T(), failed, 1)
	assert.Equal(suite.T(), failedNotification.ID, failed[0].ID)
}

func (suite *QueueRepositoryTestSuite) TestGetPendingCount() {
	// Arrange
	notification1 := domain.NewQueuedNotification(
		value_objects.NewID(),
		domain.NotificationChannelTelegram,
		"123456789",
		"Test message 1",
		"Test subject",
		1,
		3,
	)
	notification2 := domain.NewQueuedNotification(
		value_objects.NewID(),
		domain.NotificationChannelEmail,
		"test@example.com",
		"Test message 2",
		"Test subject",
		2,
		3,
	)
	notification2.Status = domain.DeliveryStatusDelivered

	suite.repo.Create(suite.ctx, notification1)
	suite.repo.Create(suite.ctx, notification2)

	// Act
	count, err := suite.repo.GetPendingCount(suite.ctx)

	// Assert
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), int64(1), count) // Only one pending
}

func (suite *QueueRepositoryTestSuite) TestGetQueueStats() {
	// Arrange
	pendingNotification := domain.NewQueuedNotification(
		value_objects.NewID(),
		domain.NotificationChannelTelegram,
		"123456789",
		"Pending message",
		"Test subject",
		1,
		3,
	)
	deliveredNotification := domain.NewQueuedNotification(
		value_objects.NewID(),
		domain.NotificationChannelEmail,
		"test@example.com",
		"Delivered message",
		"Test subject",
		2,
		3,
	)
	deliveredNotification.Status = domain.DeliveryStatusDelivered

	suite.repo.Create(suite.ctx, pendingNotification)
	suite.repo.Create(suite.ctx, deliveredNotification)

	// Act
	stats, err := suite.repo.GetQueueStats(suite.ctx)

	// Assert
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), int64(1), stats[string(domain.DeliveryStatusPending)])
	assert.Equal(suite.T(), int64(1), stats[string(domain.DeliveryStatusDelivered)])
}

func (suite *QueueRepositoryTestSuite) TestDeleteProcessedNotifications() {
	// Arrange
	oldNotification := domain.NewQueuedNotification(
		value_objects.NewID(),
		domain.NotificationChannelTelegram,
		"123456789",
		"Old message",
		"Test subject",
		1,
		3,
	)
	oldNotification.Status = domain.DeliveryStatusDelivered
	oldNotification.UpdatedAt = time.Now().Add(-2 * time.Hour) // 2 hours ago

	recentNotification := domain.NewQueuedNotification(
		value_objects.NewID(),
		domain.NotificationChannelEmail,
		"test@example.com",
		"Recent message",
		"Test subject",
		2,
		3,
	)
	recentNotification.Status = domain.DeliveryStatusDelivered
	recentNotification.UpdatedAt = time.Now().Add(-30 * time.Minute) // 30 minutes ago

	suite.repo.Create(suite.ctx, oldNotification)
	suite.repo.Create(suite.ctx, recentNotification)

	// Act - Delete notifications older than 1 hour
	err := suite.repo.DeleteProcessedNotifications(suite.ctx, time.Hour)

	// Assert
	assert.NoError(suite.T(), err)

	// Old notification should be deleted
	_, err = suite.repo.GetByID(suite.ctx, oldNotification.ID)
	assert.Error(suite.T(), err)

	// Recent notification should still exist
	_, err = suite.repo.GetByID(suite.ctx, recentNotification.ID)
	assert.NoError(suite.T(), err)
}

func (suite *QueueRepositoryTestSuite) TestDelete() {
	// Arrange
	notification := domain.NewQueuedNotification(
		value_objects.NewID(),
		domain.NotificationChannelTelegram,
		"123456789",
		"Test message",
		"Test subject",
		1,
		3,
	)
	suite.repo.Create(suite.ctx, notification)

	// Act
	err := suite.repo.Delete(suite.ctx, notification.ID)

	// Assert
	assert.NoError(suite.T(), err)

	// Should not be found
	_, err = suite.repo.GetByID(suite.ctx, notification.ID)
	assert.Error(suite.T(), err)
	assert.Equal(suite.T(), memory.ErrNotificationNotFound, err)
}

func TestQueueRepositoryTestSuite(t *testing.T) {
	suite.Run(t, new(QueueRepositoryTestSuite))
}
