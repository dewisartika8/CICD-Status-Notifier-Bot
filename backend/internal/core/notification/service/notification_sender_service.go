package service

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/dewisartika8/cicd-status-notifier-bot/internal/core/notification/domain"
	"github.com/dewisartika8/cicd-status-notifier-bot/internal/core/notification/port"
	"github.com/sirupsen/logrus"
)

// NotificationSenderAdapter implements the NotificationSender interface
type NotificationSenderDep struct {
	TelegramBotToken string
	EmailConfig      EmailConfig
	SlackConfig      SlackConfig
	Logger           *logrus.Logger
}

// EmailConfig holds email configuration
type EmailConfig struct {
	SMTPHost     string
	SMTPPort     int
	SMTPUsername string
	SMTPPassword string
	FromEmail    string
}

// SlackConfig holds Slack configuration
type SlackConfig struct {
	WebhookURL string
	BotToken   string
}

type notificationSenderService struct {
	NotificationSenderDep
}

// NewNotificationSenderAdapter creates a new notification sender adapter
func NewNotificationSenderAdapter(d NotificationSenderDep) port.NotificationSender {
	return &notificationSenderService{
		NotificationSenderDep: d,
	}
}

// SendTelegramNotification sends a notification through Telegram
func (n *notificationSenderService) SendTelegramNotification(ctx context.Context, chatID int64, message string) (messageID string, err error) {
	if n.TelegramBotToken == "" {
		return "", domain.NewNotificationSendFailedError(fmt.Errorf("telegram bot token not configured"))
	}

	if chatID == 0 {
		return "", domain.ErrInvalidTelegramChatID
	}

	if message == "" {
		return "", domain.ErrInvalidMessage
	}

	// Implement Telegram API call here
	// For now, we'll simulate the call
	n.Logger.Info(ctx, "Sending Telegram notification", map[string]interface{}{
		"chatID":  chatID,
		"message": message,
	})

	// TODO: Implement actual Telegram API call
	// Example implementation would be:
	// url := fmt.Sprintf("https://api.telegram.org/bot%s/sendMessage", n.telegramBotToken)
	// payload := map[string]interface{}{
	//     "chat_id": chatID,
	//     "text":    message,
	// }
	// ... make HTTP request and parse response for message ID

	// Return mock message ID for now
	return strconv.FormatInt(chatID, 10) + "_msg_1", nil
}

// SendEmailNotification sends a notification through email
func (n *notificationSenderService) SendEmailNotification(ctx context.Context, email, subject, message string) error {
	if email == "" {
		return domain.ErrInvalidRecipient
	}

	if message == "" {
		return domain.ErrInvalidMessage
	}

	// Validate email format
	if !strings.Contains(email, "@") {
		return domain.NewInvalidRecipientError(email)
	}

	n.Logger.Info(ctx, "Sending email notification", map[string]interface{}{
		"email":   email,
		"subject": subject,
		"message": message,
	})

	// TODO: Implement actual email sending
	// For now, just log the action
	return nil
}

// SendSlackNotification sends a notification through Slack
func (n *notificationSenderService) SendSlackNotification(ctx context.Context, channel, message string) (messageID string, err error) {
	if channel == "" {
		return "", domain.ErrInvalidRecipient
	}

	if message == "" {
		return "", domain.ErrInvalidMessage
	}

	n.Logger.Info(ctx, "Sending Slack notification", map[string]interface{}{
		"channel": channel,
		"message": message,
	})

	// TODO: Implement actual Slack API call
	// For now, return mock message ID
	return fmt.Sprintf("slack_msg_%s_1", channel), nil
}

// SendWebhookNotification sends a notification through webhook
func (n *notificationSenderService) SendWebhookNotification(ctx context.Context, webhookURL, message string) error {
	if webhookURL == "" {
		return domain.ErrInvalidRecipient
	}

	if message == "" {
		return domain.ErrInvalidMessage
	}

	n.Logger.Info(ctx, "Sending webhook notification", map[string]interface{}{
		"webhookURL": webhookURL,
		"message":    message,
	})

	// TODO: Implement actual webhook call
	// Example:
	// req, err := http.NewRequestWithContext(ctx, "POST", webhookURL, strings.NewReader(message))
	// if err != nil {
	//     return domain.NewNotificationSendFailedError(err)
	// }
	//
	// client := &http.Client{}
	// resp, err := client.Do(req)
	// if err != nil {
	//     return domain.NewNotificationSendFailedError(err)
	// }
	// defer resp.Body.Close()
	//
	// if resp.StatusCode != http.StatusOK {
	//     return domain.NewNotificationSendFailedError(fmt.Errorf("webhook returned status %d", resp.StatusCode))
	// }

	// For now, just simulate success
	_ = http.StatusOK // Prevent unused import error
	return nil
}
