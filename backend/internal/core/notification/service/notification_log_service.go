package service

import (
	"context"
	"fmt"

	"github.com/dewisartika8/cicd-status-notifier-bot/internal/core/notification/domain"
	"github.com/dewisartika8/cicd-status-notifier-bot/internal/core/notification/port"
	"github.com/dewisartika8/cicd-status-notifier-bot/internal/core/shared/domain/value_objects"
	"github.com/sirupsen/logrus"
)

// Error message constants for better maintainability
const (
	errMsgCreateNotificationLog    = "failed to create notification log: %w"
	errMsgPersistNotificationLog   = "failed to persist notification log: %w"
	errMsgGetNotificationLog       = "failed to get notification log: %w"
	errMsgUpdateNotificationLog    = "failed to update notification log: %w"
	errMsgGetTelegramSubscriptions = "failed to get telegram subscriptions: %w"
	errMsgGetPendingNotifications  = "failed to get pending notifications: %w"
	errMsgGetFailedNotifications   = "failed to get failed notifications: %w"
	errMsgGetNotificationStats     = "failed to get notification stats: %w"
	errMsgSendTelegramNotification = "failed to send telegram notification: %w"
	errMsgSendEmailNotification    = "failed to send email notification: %w"
	errMsgSendSlackNotification    = "failed to send slack notification: %w"
	errMsgSendWebhookNotification  = "failed to send webhook notification: %w"
	errMsgMarkNotificationAsSent   = "failed to mark notification as sent: %w"
	errMsgMarkNotificationAsFailed = "failed to mark notification as failed: %w"
	errMsgMarkNotificationAsRetry  = "failed to mark notification as retrying: %w"
)

// Log message constants to avoid duplication
const (
	logMsgGetNotificationLog    = "Failed to get notification log"
	logMsgUpdateNotificationLog = "Failed to update notification log"
	logMsgSendNotification      = "Failed to send notification"
	logMsgMarkNotificationFail  = "Failed to mark notification as failed"
	logMsgMarkNotificationSent  = "Failed to mark notification as sent"
)

type NotificationLogDep struct {
	NotificationRepo         port.NotificationLogRepository
	TelegramSubscriptionRepo port.TelegramSubscriptionRepository
	NotificationSender       port.NotificationSender
	Logger                   *logrus.Logger
}

// notificationLogService implements notification log business logic
type notificationLogService struct {
	NotificationLogDep
}

// NewNotificationLogService creates a new notification log service
func NewNotificationLogService(d NotificationLogDep) port.NotificationLogService {
	return &notificationLogService{
		NotificationLogDep: d,
	}
}

// CreateNotificationLog creates a new notification log
func (s *notificationLogService) CreateNotificationLog(
	ctx context.Context,
	buildEventID, projectID value_objects.ID,
	channel domain.NotificationChannel,
	recipient, message string,
) (*domain.NotificationLog, error) {
	s.Logger.WithFields(logrus.Fields{
		"build_event_id": buildEventID.String(),
		"project_id":     projectID.String(),
		"channel":        string(channel),
		"recipient":      recipient,
	}).Info("Creating notification log")

	// Create new notification log entity
	notificationLog, err := domain.NewNotificationLog(buildEventID, projectID, channel, recipient, message)
	if err != nil {
		s.Logger.WithError(err).Error("Failed to create notification log entity")
		return nil, fmt.Errorf(errMsgCreateNotificationLog, err)
	}

	// Persist the notification log
	if err := s.NotificationRepo.Create(ctx, notificationLog); err != nil {
		s.Logger.WithError(err).Error("Failed to persist notification log")
		return nil, fmt.Errorf(errMsgPersistNotificationLog, err)
	}

	s.Logger.WithField("notification_log_id", notificationLog.ID().String()).Info("Notification log created successfully")
	return notificationLog, nil
}

// SendNotification sends a notification and updates the log
func (s *notificationLogService) SendNotification(ctx context.Context, notificationLogID value_objects.ID) error {
	s.Logger.WithField("notification_log_id", notificationLogID.String()).Info("Sending notification")

	// Get the notification log
	notificationLog, err := s.NotificationRepo.GetByID(ctx, notificationLogID)
	if err != nil {
		s.Logger.WithError(err).Error(logMsgGetNotificationLog)
		return fmt.Errorf(errMsgGetNotificationLog, err)
	}

	// Send notification based on channel
	var messageID *string
	var sendErr error

	switch notificationLog.Channel() {
	case domain.NotificationChannelTelegram:
		messageID, sendErr = s.sendTelegramNotification(ctx, notificationLog)
	case domain.NotificationChannelEmail:
		sendErr = s.sendEmailNotification(ctx, notificationLog)
	case domain.NotificationChannelSlack:
		messageID, sendErr = s.sendSlackNotification(ctx, notificationLog)
	case domain.NotificationChannelWebhook:
		sendErr = s.sendWebhookNotification(ctx, notificationLog)
	default:
		sendErr = fmt.Errorf("unsupported notification channel: %s", notificationLog.Channel())
	}

	// Update notification status based on send result
	if sendErr != nil {
		s.Logger.WithError(sendErr).Error(logMsgSendNotification)
		if err := notificationLog.MarkAsFailed(sendErr.Error()); err != nil {
			s.Logger.WithError(err).Error(logMsgMarkNotificationFail)
		}
	} else {
		s.Logger.WithField("message_id", messageID).Info("Notification sent successfully")
		if err := notificationLog.MarkAsSent(messageID); err != nil {
			s.Logger.WithError(err).Error(logMsgMarkNotificationSent)
		}
	}

	// Update the notification log in repository
	if err := s.NotificationRepo.Update(ctx, notificationLog); err != nil {
		s.Logger.WithError(err).Error(logMsgUpdateNotificationLog)
		return fmt.Errorf(errMsgUpdateNotificationLog, err)
	}

	return sendErr
}

// GetNotificationLog retrieves a notification log by its ID
func (s *notificationLogService) GetNotificationLog(ctx context.Context, id value_objects.ID) (*domain.NotificationLog, error) {
	s.Logger.WithField("id", id.String()).Info("Getting notification log")

	notificationLog, err := s.NotificationRepo.GetByID(ctx, id)
	if err != nil {
		s.Logger.WithError(err).Error(logMsgGetNotificationLog)
		return nil, fmt.Errorf(errMsgGetNotificationLog, err)
	}

	return notificationLog, nil
}

// GetNotificationLogsByBuildEvent retrieves notification logs for a build event
func (s *notificationLogService) GetNotificationLogsByBuildEvent(ctx context.Context, buildEventID value_objects.ID) ([]*domain.NotificationLog, error) {
	s.Logger.WithField("build_event_id", buildEventID.String()).Info("Getting notification logs by build event")

	notificationLogs, err := s.NotificationRepo.GetByBuildEventID(ctx, buildEventID)
	if err != nil {
		s.Logger.WithError(err).Error("Failed to get notification logs by build event")
		return nil, fmt.Errorf("failed to get notification logs by build event: %w", err)
	}

	return notificationLogs, nil
}

// GetNotificationLogsByProject retrieves notification logs for a project
func (s *notificationLogService) GetNotificationLogsByProject(ctx context.Context, projectID value_objects.ID, limit, offset int) ([]*domain.NotificationLog, error) {
	s.Logger.WithFields(logrus.Fields{
		"project_id": projectID.String(),
		"limit":      limit,
		"offset":     offset,
	}).Info("Getting notification logs by project")

	notificationLogs, err := s.NotificationRepo.GetByProjectID(ctx, projectID, limit, offset)
	if err != nil {
		s.Logger.WithError(err).Error("Failed to get notification logs by project")
		return nil, fmt.Errorf("failed to get notification logs by project: %w", err)
	}

	return notificationLogs, nil
}

// UpdateNotificationStatus updates the status of a notification log
func (s *notificationLogService) UpdateNotificationStatus(
	ctx context.Context,
	id value_objects.ID,
	status domain.NotificationStatus,
	errorMessage string,
	messageID *string,
) error {
	s.Logger.WithFields(logrus.Fields{
		"id":            id.String(),
		"status":        string(status),
		"error_message": errorMessage,
	}).Info("Updating notification status")

	// Get the notification log
	notificationLog, err := s.NotificationRepo.GetByID(ctx, id)
	if err != nil {
		s.Logger.WithError(err).Error(logMsgGetNotificationLog)
		return fmt.Errorf(errMsgGetNotificationLog, err)
	}

	// Update status based on the provided status
	switch status {
	case domain.NotificationStatusSent:
		if err := notificationLog.MarkAsSent(messageID); err != nil {
			return fmt.Errorf(errMsgMarkNotificationAsSent, err)
		}
	case domain.NotificationStatusFailed:
		if err := notificationLog.MarkAsFailed(errorMessage); err != nil {
			return fmt.Errorf(errMsgMarkNotificationAsFailed, err)
		}
	case domain.NotificationStatusRetrying:
		if err := notificationLog.MarkAsRetrying(); err != nil {
			return fmt.Errorf(errMsgMarkNotificationAsRetry, err)
		}
	default:
		return fmt.Errorf("unsupported status update: %s", status)
	}

	// Update the notification log in repository
	if err := s.NotificationRepo.Update(ctx, notificationLog); err != nil {
		s.Logger.WithError(err).Error(logMsgUpdateNotificationLog)
		return fmt.Errorf(errMsgUpdateNotificationLog, err)
	}

	s.Logger.Info("Notification status updated successfully")
	return nil
}

// RetryFailedNotification retries a failed notification
func (s *notificationLogService) RetryFailedNotification(ctx context.Context, notificationLogID value_objects.ID) error {
	s.Logger.WithField("notification_log_id", notificationLogID.String()).Info("Retrying failed notification")

	// Get the notification log
	notificationLog, err := s.NotificationRepo.GetByID(ctx, notificationLogID)
	if err != nil {
		s.Logger.WithError(err).Error(logMsgGetNotificationLog)
		return fmt.Errorf(errMsgGetNotificationLog, err)
	}

	// Check if notification can be retried
	if !notificationLog.CanRetry() {
		return fmt.Errorf("notification cannot be retried: max attempts exceeded or invalid status")
	}

	// Mark as retrying
	if err := notificationLog.MarkAsRetrying(); err != nil {
		return fmt.Errorf(errMsgMarkNotificationAsRetry, err)
	}

	// Update the notification log
	if err := s.NotificationRepo.Update(ctx, notificationLog); err != nil {
		s.Logger.WithError(err).Error(logMsgUpdateNotificationLog)
		return fmt.Errorf(errMsgUpdateNotificationLog, err)
	}

	// Attempt to send the notification again
	return s.SendNotification(ctx, notificationLogID)
}

// ProcessPendingNotifications processes all pending notifications
func (s *notificationLogService) ProcessPendingNotifications(ctx context.Context, limit int) error {
	s.Logger.WithField("limit", limit).Info("Processing pending notifications")

	pendingNotifications, err := s.NotificationRepo.GetPendingNotifications(ctx, limit)
	if err != nil {
		s.Logger.WithError(err).Error("Failed to get pending notifications")
		return fmt.Errorf(errMsgGetPendingNotifications, err)
	}

	var processedCount int
	var errorCount int

	for _, notification := range pendingNotifications {
		if err := s.SendNotification(ctx, notification.ID()); err != nil {
			s.Logger.WithFields(logrus.Fields{
				"notification_id": notification.ID().String(),
				"error":           err,
			}).Error("Failed to send pending notification")
			errorCount++
		} else {
			processedCount++
		}
	}

	s.Logger.WithFields(logrus.Fields{
		"total":     len(pendingNotifications),
		"processed": processedCount,
		"errors":    errorCount,
	}).Info("Finished processing pending notifications")

	return nil
}

// ProcessFailedNotifications processes failed notifications for retry
func (s *notificationLogService) ProcessFailedNotifications(ctx context.Context, limit int) error {
	s.Logger.WithField("limit", limit).Info("Processing failed notifications for retry")

	failedNotifications, err := s.NotificationRepo.GetFailedNotifications(ctx, limit)
	if err != nil {
		s.Logger.WithError(err).Error("Failed to get failed notifications")
		return fmt.Errorf(errMsgGetFailedNotifications, err)
	}

	var retriedCount int
	var errorCount int

	for _, notification := range failedNotifications {
		if notification.CanRetry() {
			if err := s.RetryFailedNotification(ctx, notification.ID()); err != nil {
				s.Logger.WithFields(logrus.Fields{
					"notification_id": notification.ID().String(),
					"error":           err,
				}).Error("Failed to retry notification")
				errorCount++
			} else {
				retriedCount++
			}
		}
	}

	s.Logger.WithFields(logrus.Fields{
		"total":   len(failedNotifications),
		"retried": retriedCount,
		"errors":  errorCount,
	}).Info("Finished processing failed notifications")

	return nil
}

// GetNotificationStats retrieves notification statistics for a project
func (s *notificationLogService) GetNotificationStats(ctx context.Context, projectID value_objects.ID) (map[domain.NotificationStatus]int64, error) {
	s.Logger.WithField("project_id", projectID.String()).Info("Getting notification stats")

	stats, err := s.NotificationRepo.GetNotificationStats(ctx, projectID)
	if err != nil {
		s.Logger.WithError(err).Error("Failed to get notification stats")
		return nil, fmt.Errorf(errMsgGetNotificationStats, err)
	}

	return stats, nil
}

// CreateNotificationForBuildEvent creates notifications for all subscribed channels for a build event
func (s *notificationLogService) CreateNotificationForBuildEvent(
	ctx context.Context,
	buildEventID, projectID value_objects.ID,
	message string,
) ([]*domain.NotificationLog, error) {
	s.Logger.WithFields(logrus.Fields{
		"build_event_id": buildEventID.String(),
		"project_id":     projectID.String(),
	}).Info("Creating notifications for build event")

	var notifications []*domain.NotificationLog

	// Get active Telegram subscriptions for the project
	telegramSubscriptions, err := s.TelegramSubscriptionRepo.GetActiveSubscriptionsByProject(ctx, projectID)
	if err != nil {
		s.Logger.WithError(err).Error("Failed to get telegram subscriptions")
		return nil, fmt.Errorf(errMsgGetTelegramSubscriptions, err)
	}

	// Create notification logs for each Telegram subscription
	for _, subscription := range telegramSubscriptions {
		notificationLog, err := s.CreateNotificationLog(
			ctx,
			buildEventID,
			projectID,
			domain.NotificationChannelTelegram,
			subscription.GetChatIDString(),
			message,
		)
		if err != nil {
			s.Logger.WithFields(logrus.Fields{
				"chat_id": subscription.ChatID(),
				"error":   err,
			}).Error("Failed to create telegram notification log")
			continue
		}
		notifications = append(notifications, notificationLog)
	}

	s.Logger.WithField("total_notifications", len(notifications)).Info("Created notifications for build event")

	return notifications, nil
}

// Helper methods for sending notifications through different channels

func (s *notificationLogService) sendTelegramNotification(ctx context.Context, log *domain.NotificationLog) (*string, error) {
	// Parse chat ID from recipient
	var chatID int64
	if _, err := fmt.Sscanf(log.Recipient(), "%d", &chatID); err != nil {
		return nil, fmt.Errorf("invalid telegram chat ID: %s", log.Recipient())
	}

	messageID, err := s.NotificationSender.SendTelegramNotification(ctx, chatID, log.Message())
	if err != nil {
		return nil, fmt.Errorf(errMsgSendTelegramNotification, err)
	}

	return &messageID, nil
}

func (s *notificationLogService) sendEmailNotification(ctx context.Context, log *domain.NotificationLog) error {
	// For email, we'd need to parse subject from message or have it as a separate field
	subject := "CI/CD Build Notification"
	if err := s.NotificationSender.SendEmailNotification(ctx, log.Recipient(), subject, log.Message()); err != nil {
		return fmt.Errorf(errMsgSendEmailNotification, err)
	}
	return nil
}

func (s *notificationLogService) sendSlackNotification(ctx context.Context, log *domain.NotificationLog) (*string, error) {
	messageID, err := s.NotificationSender.SendSlackNotification(ctx, log.Recipient(), log.Message())
	if err != nil {
		return nil, fmt.Errorf(errMsgSendSlackNotification, err)
	}
	return &messageID, nil
}

func (s *notificationLogService) sendWebhookNotification(ctx context.Context, log *domain.NotificationLog) error {
	if err := s.NotificationSender.SendWebhookNotification(ctx, log.Recipient(), log.Message()); err != nil {
		return fmt.Errorf(errMsgSendWebhookNotification, err)
	}
	return nil
}
