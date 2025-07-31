package log

import (
	"context"
	"fmt"

	"github.com/dewisartika8/cicd-status-notifier-bot/internal/core/notification/domain"
	"github.com/dewisartika8/cicd-status-notifier-bot/internal/core/notification/port"
	"github.com/dewisartika8/cicd-status-notifier-bot/internal/core/shared/domain/value_objects"
	"github.com/sirupsen/logrus"
)

type Dep struct {
	NotificationRepo         port.NotificationLogRepository
	TelegramSubscriptionRepo port.TelegramSubscriptionRepository
	NotificationSender       port.NotificationSender
	Logger                   *logrus.Logger
}

// notificationLogService implements notification log business logic
type notificationLogService struct {
	Dep
}

// NewNotificationLogService creates a new notification log service
func NewNotificationLogService(d Dep) port.NotificationLogService {
	return &notificationLogService{
		Dep: d,
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
		"channel":        channel,
		"recipient":      recipient,
	}).Info("Creating notification log")

	// Create new notification log entity
	log, err := domain.NewNotificationLog(buildEventID, projectID, channel, recipient, message, 3) // Default maxRetries = 3
	if err != nil {
		s.Logger.WithError(err).Error("Failed to create notification log entity")
		return nil, fmt.Errorf(domain.ErrMsgCreateNotificationLog, err)
	}

	// Persist the notification log
	if err := s.NotificationRepo.Create(ctx, log); err != nil {
		s.Logger.WithError(err).Error("Failed to persist notification log")
		return nil, fmt.Errorf(domain.ErrMsgPersistNotificationLog, err)
	}

	s.Logger.WithField("log_id", log.ID().String()).Info("Notification log created successfully")
	return log, nil
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

	// Get telegram subscriptions for this project
	subscriptions, err := s.TelegramSubscriptionRepo.GetByProjectID(ctx, projectID)
	if err != nil {
		s.Logger.WithError(err).Error("Failed to get telegram subscriptions")
		return nil, fmt.Errorf(domain.ErrMsgGetTelegramSubscriptions, err)
	}

	var notifications []*domain.NotificationLog

	// Create notification for each subscription
	for _, subscription := range subscriptions {
		// Create notification log for telegram
		log, err := s.CreateNotificationLog(
			ctx,
			buildEventID,
			projectID,
			domain.NotificationChannelTelegram,
			fmt.Sprintf("%d", subscription.ChatID()),
			message,
		)
		if err != nil {
			s.Logger.WithError(err).WithField("chat_id", subscription.ChatID()).Error("Failed to create notification log")
			continue
		}

		notifications = append(notifications, log)
	}

	s.Logger.WithFields(logrus.Fields{
		"build_event_id":      buildEventID.String(),
		"notifications_count": len(notifications),
	}).Info("Created notifications for build event")

	return notifications, nil
}

// SendNotification sends a notification and updates the log
func (s *notificationLogService) SendNotification(ctx context.Context, notificationLogID value_objects.ID) error {
	s.Logger.WithField("log_id", notificationLogID.String()).Info("Sending notification")

	// Get the notification log
	log, err := s.getNotificationLog(ctx, notificationLogID)
	if err != nil {
		return err
	}

	// Send notification through appropriate channel
	messageID, err := s.sendNotificationByChannel(ctx, log)
	if err != nil {
		return s.handleSendFailure(ctx, log, err)
	}

	// Mark notification as sent and update
	return s.markNotificationAsSent(ctx, log, messageID)
}

// getNotificationLog retrieves and validates notification log
func (s *notificationLogService) getNotificationLog(ctx context.Context, notificationLogID value_objects.ID) (*domain.NotificationLog, error) {
	log, err := s.NotificationRepo.GetByID(ctx, notificationLogID)
	if err != nil {
		s.Logger.WithError(err).Error(domain.LogMsgGetNotificationLog)
		return nil, fmt.Errorf(domain.ErrMsgGetNotificationLog, err)
	}
	return log, nil
}

// sendNotificationByChannel sends notification using the appropriate channel
func (s *notificationLogService) sendNotificationByChannel(ctx context.Context, log *domain.NotificationLog) (string, error) {
	switch log.Channel() {
	case domain.NotificationChannelTelegram:
		return s.sendTelegramNotification(ctx, log)
	case domain.NotificationChannelEmail:
		return s.sendEmailNotification(ctx, log)
	case domain.NotificationChannelSlack:
		return s.sendSlackNotification(ctx, log)
	case domain.NotificationChannelWebhook:
		return s.sendWebhookNotification(ctx, log)
	default:
		return "", fmt.Errorf("unsupported notification channel: %s", log.Channel())
	}
}

// sendTelegramNotification handles Telegram-specific notification sending
func (s *notificationLogService) sendTelegramNotification(ctx context.Context, log *domain.NotificationLog) (string, error) {
	chatID, err := s.parseTelegramChatID(log.Recipient())
	if err != nil {
		return "", err
	}

	messageID, err := s.NotificationSender.SendTelegramNotification(ctx, chatID, log.Message())
	if err != nil {
		s.Logger.WithError(err).Error("Failed to send telegram notification")
		return "", fmt.Errorf(domain.ErrMsgSendTelegramNotification, err)
	}

	return messageID, nil
}

// sendEmailNotification handles Email-specific notification sending
func (s *notificationLogService) sendEmailNotification(ctx context.Context, log *domain.NotificationLog) (string, error) {
	err := s.NotificationSender.SendEmailNotification(ctx, log.Recipient(), log.Message(), log.Message())
	if err != nil {
		s.Logger.WithError(err).Error("Failed to send email notification")
		return "", fmt.Errorf(domain.ErrMsgSendEmailNotification, err)
	}

	return "email-sent", nil // Email doesn't return messageID
}

// sendSlackNotification handles Slack-specific notification sending
func (s *notificationLogService) sendSlackNotification(ctx context.Context, log *domain.NotificationLog) (string, error) {
	messageID, err := s.NotificationSender.SendSlackNotification(ctx, log.Recipient(), log.Message())
	if err != nil {
		s.Logger.WithError(err).Error("Failed to send slack notification")
		return "", fmt.Errorf(domain.ErrMsgSendSlackNotification, err)
	}

	return messageID, nil
}

// sendWebhookNotification handles Webhook-specific notification sending
func (s *notificationLogService) sendWebhookNotification(ctx context.Context, log *domain.NotificationLog) (string, error) {
	err := s.NotificationSender.SendWebhookNotification(ctx, log.Recipient(), log.Message())
	if err != nil {
		s.Logger.WithError(err).Error("Failed to send webhook notification")
		return "", fmt.Errorf(domain.ErrMsgSendWebhookNotification, err)
	}

	return "webhook-sent", nil // Webhook doesn't return messageID
}

// parseTelegramChatID parses and validates Telegram chat ID from recipient string
func (s *notificationLogService) parseTelegramChatID(recipient string) (int64, error) {
	var chatID int64
	if _, err := fmt.Sscanf(recipient, "%d", &chatID); err != nil {
		s.Logger.WithError(err).Error("Failed to parse telegram chat ID")
		return 0, fmt.Errorf("failed to parse telegram chat ID: %w", err)
	}
	return chatID, nil
}

// handleSendFailure handles notification send failure by marking as failed
func (s *notificationLogService) handleSendFailure(ctx context.Context, log *domain.NotificationLog, sendErr error) error {
	log.MarkAsFailed(sendErr.Error())
	if updateErr := s.NotificationRepo.Update(ctx, log); updateErr != nil {
		s.Logger.WithError(updateErr).Error(domain.LogMsgMarkNotificationFail)
		// Return the original send error, but log the update error
	}
	return sendErr
}

// markNotificationAsSent marks notification as sent and updates the repository
func (s *notificationLogService) markNotificationAsSent(ctx context.Context, log *domain.NotificationLog, messageID string) error {
	if err := log.MarkAsSent(&messageID); err != nil {
		s.Logger.WithError(err).Error(domain.LogMsgMarkNotificationSent)
		return fmt.Errorf(domain.ErrMsgMarkNotificationAsSent, err)
	}

	if err := s.NotificationRepo.Update(ctx, log); err != nil {
		s.Logger.WithError(err).Error(domain.LogMsgMarkNotificationSent)
		return fmt.Errorf(domain.ErrMsgMarkNotificationAsSent, err)
	}

	s.Logger.WithFields(logrus.Fields{
		"log_id":     log.ID().String(),
		"message_id": messageID,
		"channel":    log.Channel(),
	}).Info("Notification sent successfully")

	return nil
}

// GetNotificationLog retrieves a notification log by its ID
func (s *notificationLogService) GetNotificationLog(ctx context.Context, id value_objects.ID) (*domain.NotificationLog, error) {
	s.Logger.WithField("id", id.String()).Info("Getting notification log")

	log, err := s.NotificationRepo.GetByID(ctx, id)
	if err != nil {
		s.Logger.WithError(err).Error(domain.LogMsgGetNotificationLog)
		return nil, fmt.Errorf(domain.ErrMsgGetNotificationLog, err)
	}

	return log, nil
}

// GetNotificationLogsByBuildEvent retrieves notification logs for a build event
func (s *notificationLogService) GetNotificationLogsByBuildEvent(ctx context.Context, buildEventID value_objects.ID) ([]*domain.NotificationLog, error) {
	s.Logger.WithField("build_event_id", buildEventID.String()).Info("Getting notification logs by build event")

	logs, err := s.NotificationRepo.GetByBuildEventID(ctx, buildEventID)
	if err != nil {
		s.Logger.WithError(err).Error("Failed to get notification logs by build event")
		return nil, fmt.Errorf("failed to get notification logs by build event: %w", err)
	}

	return logs, nil
}

// GetNotificationLogsByProject retrieves notification logs for a project
func (s *notificationLogService) GetNotificationLogsByProject(ctx context.Context, projectID value_objects.ID, limit, offset int) ([]*domain.NotificationLog, error) {
	s.Logger.WithFields(logrus.Fields{
		"project_id": projectID.String(),
		"limit":      limit,
		"offset":     offset,
	}).Info("Getting notification logs by project")

	logs, err := s.NotificationRepo.GetByProjectID(ctx, projectID, limit, offset)
	if err != nil {
		s.Logger.WithError(err).Error("Failed to get notification logs by project")
		return nil, fmt.Errorf("failed to get notification logs by project: %w", err)
	}

	return logs, nil
}

// GetPendingNotifications retrieves pending notifications
func (s *notificationLogService) GetPendingNotifications(ctx context.Context, limit int) ([]*domain.NotificationLog, error) {
	s.Logger.WithField("limit", limit).Info("Getting pending notifications")

	logs, err := s.NotificationRepo.GetPendingNotifications(ctx, limit)
	if err != nil {
		s.Logger.WithError(err).Error("Failed to get pending notifications")
		return nil, fmt.Errorf(domain.ErrMsgGetPendingNotifications, err)
	}

	return logs, nil
}

// GetFailedNotifications retrieves failed notifications
func (s *notificationLogService) GetFailedNotifications(ctx context.Context, limit int) ([]*domain.NotificationLog, error) {
	s.Logger.WithField("limit", limit).Info("Getting failed notifications")

	logs, err := s.NotificationRepo.GetFailedNotifications(ctx, limit)
	if err != nil {
		s.Logger.WithError(err).Error("Failed to get failed notifications")
		return nil, fmt.Errorf(domain.ErrMsgGetFailedNotifications, err)
	}

	return logs, nil
}

// GetNotificationStats retrieves notification statistics
func (s *notificationLogService) GetNotificationStats(ctx context.Context, projectID value_objects.ID) (map[domain.NotificationStatus]int64, error) {
	s.Logger.WithField("project_id", projectID.String()).Info("Getting notification stats")

	stats, err := s.NotificationRepo.GetNotificationStats(ctx, projectID)
	if err != nil {
		s.Logger.WithError(err).Error("Failed to get notification stats")
		return nil, fmt.Errorf(domain.ErrMsgGetNotificationStats, err)
	}

	// Convert NotificationStats to map
	result := make(map[domain.NotificationStatus]int64)
	if stats != nil && stats.StatusCounts != nil {
		for status, count := range stats.StatusCounts {
			result[status] = count
		}
	}

	return result, nil
}

// ProcessPendingNotifications processes all pending notifications
func (s *notificationLogService) ProcessPendingNotifications(ctx context.Context, limit int) error {
	s.Logger.WithField("limit", limit).Info("Processing pending notifications")

	// Get pending notifications
	pendingLogs, err := s.GetPendingNotifications(ctx, limit)
	if err != nil {
		s.Logger.WithError(err).Error("Failed to get pending notifications")
		return err
	}

	// Process each pending notification
	for _, log := range pendingLogs {
		if err := s.SendNotification(ctx, log.ID()); err != nil {
			s.Logger.WithError(err).WithField("log_id", log.ID().String()).Error("Failed to send pending notification")
			// Continue with other notifications even if one fails
		}
	}

	s.Logger.WithFields(logrus.Fields{
		"processed_count": len(pendingLogs),
	}).Info("Completed processing pending notifications")

	return nil
}

// ProcessFailedNotifications processes failed notifications for retry
func (s *notificationLogService) ProcessFailedNotifications(ctx context.Context, limit int) error {
	s.Logger.WithField("limit", limit).Info("Processing failed notifications")

	// Get failed notifications
	failedLogs, err := s.GetFailedNotifications(ctx, limit)
	if err != nil {
		s.Logger.WithError(err).Error("Failed to get failed notifications")
		return err
	}

	// Process each failed notification
	for _, log := range failedLogs {
		// Try to send the notification again
		if err := s.SendNotification(ctx, log.ID()); err != nil {
			s.Logger.WithError(err).WithField("log_id", log.ID().String()).Error("Failed to retry notification")
			// Mark as retrying if retry count hasn't exceeded
			if retryErr := s.MarkNotificationAsRetrying(ctx, log.ID(), log.RetryCount(), err.Error()); retryErr != nil {
				s.Logger.WithError(retryErr).WithField("log_id", log.ID().String()).Error(domain.LogMsgMarkNotificationAsRetrying)
			}
		}
	}

	s.Logger.WithFields(logrus.Fields{
		"processed_count": len(failedLogs),
	}).Info("Completed processing failed notifications")

	return nil
}

// RetryFailedNotification retries a failed notification
func (s *notificationLogService) RetryFailedNotification(ctx context.Context, notificationLogID value_objects.ID) error {
	s.Logger.WithField("log_id", notificationLogID.String()).Info("Retrying failed notification")

	// Get the notification log
	log, err := s.NotificationRepo.GetByID(ctx, notificationLogID)
	if err != nil {
		s.Logger.WithError(err).Error(domain.LogMsgGetNotificationLog)
		return fmt.Errorf(domain.ErrMsgGetNotificationLog, err)
	}

	// Check if notification is in failed status
	if log.Status() != domain.NotificationStatusFailed {
		err := fmt.Errorf("notification is not in failed status, current status: %s", log.Status())
		s.Logger.WithError(err).Error("Cannot retry notification")
		return err
	}

	// Try to send the notification again
	if err := s.SendNotification(ctx, notificationLogID); err != nil {
		s.Logger.WithError(err).Error("Failed to retry notification")
		// Mark as retrying if retry count hasn't exceeded
		if retryErr := s.MarkNotificationAsRetrying(ctx, notificationLogID, log.RetryCount(), err.Error()); retryErr != nil {
			s.Logger.WithError(retryErr).Error(domain.LogMsgMarkNotificationAsRetrying)
		}
		return err
	}

	s.Logger.WithField("log_id", notificationLogID.String()).Info("Notification retried successfully")
	return nil
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
		"log_id": id.String(),
		"status": status,
	}).Info("Updating notification status")

	// Get the notification log
	log, err := s.NotificationRepo.GetByID(ctx, id)
	if err != nil {
		s.Logger.WithError(err).Error(domain.LogMsgGetNotificationLog)
		return fmt.Errorf(domain.ErrMsgGetNotificationLog, err)
	}

	// Update status based on the provided status
	switch status {
	case domain.NotificationStatusSent:
		if err := log.MarkAsSent(messageID); err != nil {
			return err
		}
	case domain.NotificationStatusDelivered:
		if err := log.MarkAsDelivered(); err != nil {
			return err
		}
	case domain.NotificationStatusFailed:
		if err := log.MarkAsFailed(errorMessage); err != nil {
			return err
		}
	case domain.NotificationStatusRetrying:
		if err := log.MarkAsRetrying(); err != nil {
			return err
		}
	case domain.NotificationStatusCancelled:
		if err := log.MarkAsCancelled(); err != nil {
			return err
		}
	case domain.NotificationStatusExpired:
		if err := log.MarkAsExpired(); err != nil {
			return err
		}
	default:
		return fmt.Errorf("unsupported notification status: %s", status)
	}

	// Save the updated log
	if err := s.NotificationRepo.Update(ctx, log); err != nil {
		s.Logger.WithError(err).Error(domain.LogMsgUpdateNotificationLog)
		return fmt.Errorf(domain.ErrMsgUpdateNotificationLog, err)
	}

	s.Logger.WithFields(logrus.Fields{
		"log_id": id.String(),
		"status": status,
	}).Info("Notification status updated successfully")

	return nil
}

// MarkNotificationAsRetrying marks a notification as retrying
func (s *notificationLogService) MarkNotificationAsRetrying(ctx context.Context, id value_objects.ID, retryCount int, errorMessage string) error {
	s.Logger.WithFields(logrus.Fields{
		"id":          id.String(),
		"retry_count": retryCount,
	}).Info("Marking notification as retrying")

	log, err := s.NotificationRepo.GetByID(ctx, id)
	if err != nil {
		s.Logger.WithError(err).Error(domain.LogMsgGetNotificationLog)
		return fmt.Errorf(domain.ErrMsgGetNotificationLog, err)
	}

	// Set error message before marking as retrying
	if errorMessage != "" {
		log.MarkAsFailed(errorMessage)
	}

	if err := log.MarkAsRetrying(); err != nil {
		s.Logger.WithError(err).Error("Failed to mark notification as retrying")
		return fmt.Errorf(domain.ErrMsgMarkNotificationAsRetry, err)
	}

	if err := s.NotificationRepo.Update(ctx, log); err != nil {
		s.Logger.WithError(err).Error(domain.LogMsgMarkNotificationAsRetrying)
		return fmt.Errorf(domain.ErrMsgMarkNotificationAsRetry, err)
	}

	s.Logger.Info("Notification marked as retrying successfully")
	return nil
}

// ProcessNotificationsForBuildEvent processes notifications for a build event
func (s *notificationLogService) ProcessNotificationsForBuildEvent(ctx context.Context, buildEventID value_objects.ID) error {
	s.Logger.WithField("build_event_id", buildEventID.String()).Info("Processing notifications for build event")

	// Get pending notifications for this build event
	logs, err := s.NotificationRepo.GetByBuildEventID(ctx, buildEventID)
	if err != nil {
		s.Logger.WithError(err).Error("Failed to get notification logs by build event")
		return fmt.Errorf("failed to get notification logs by build event: %w", err)
	}

	// Filter only pending notifications
	pendingLogs := make([]*domain.NotificationLog, 0)
	for _, log := range logs {
		if log.Status() == domain.NotificationStatusPending {
			pendingLogs = append(pendingLogs, log)
		}
	}

	// Send each pending notification
	for _, log := range pendingLogs {
		if err := s.SendNotification(ctx, log.ID()); err != nil {
			s.Logger.WithError(err).WithField("log_id", log.ID().String()).Error("Failed to send notification")
			// Continue with other notifications even if one fails
		}
	}

	s.Logger.WithFields(logrus.Fields{
		"build_event_id":   buildEventID.String(),
		"processed_count":  len(pendingLogs),
		"total_logs_count": len(logs),
	}).Info("Completed processing notifications for build event")

	return nil
}
