package sender

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/dewisartika8/cicd-status-notifier-bot/internal/core/notification/domain"
	"github.com/dewisartika8/cicd-status-notifier-bot/internal/core/notification/port"
	"github.com/sirupsen/logrus"
)

// Resource type constants
const (
	resourceTelegramMsg = "telegram message"
	resourceEmailMsg    = "email message"
	resourceSlackMsg    = "slack message"
	resourceWebhookMsg  = "webhook message"
)

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

type Dep struct {
	TelegramBotToken string
	EmailConfig      EmailConfig
	SlackConfig      SlackConfig
	Logger           *logrus.Logger
}

type notificationSenderService struct {
	Dep
}

// NewNotificationSenderService creates a new notification sender service
func NewNotificationSenderService(d Dep) port.NotificationSender {
	return &notificationSenderService{
		Dep: d,
	}
}

// SendTelegramNotification sends a notification through Telegram
func (s *notificationSenderService) SendTelegramNotification(ctx context.Context, chatID int64, message string) (messageID string, err error) {
	s.Logger.WithFields(logrus.Fields{
		"chat_id": chatID,
		"message": message,
	}).Info(domain.LogMsgSendingTelegram)

	if s.TelegramBotToken == "" {
		err = fmt.Errorf("telegram bot token is not configured")
		s.Logger.WithError(err).Error("Telegram bot token missing")
		return "", fmt.Errorf(domain.ErrMsgSend, resourceTelegramMsg, err)
	}

	// Prepare the request URL
	url := fmt.Sprintf("https://api.telegram.org/bot%s/sendMessage", s.TelegramBotToken)

	// Prepare the request body
	body := fmt.Sprintf(`{
		"chat_id": %d,
		"text": "%s",
		"parse_mode": "HTML"
	}`, chatID, escapeJSON(message))

	// Create HTTP request
	req, err := http.NewRequestWithContext(ctx, "POST", url, strings.NewReader(body))
	if err != nil {
		s.Logger.WithError(err).Error("Failed to create telegram request")
		return "", fmt.Errorf(domain.ErrMsgSend, resourceTelegramMsg, err)
	}

	req.Header.Set("Content-Type", "application/json")

	// Send the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		s.Logger.WithError(err).Error("Failed to send telegram request")
		return "", fmt.Errorf(domain.ErrMsgSend, resourceTelegramMsg, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		err = fmt.Errorf("telegram API returned status: %d", resp.StatusCode)
		s.Logger.WithError(err).Error("Telegram API error")
		return "", fmt.Errorf(domain.ErrMsgSend, resourceTelegramMsg, err)
	}

	// For simplicity, return a placeholder message ID
	// In a real implementation, you would parse the response to get the actual message ID
	messageID = fmt.Sprintf("tg_%d_%s", chatID, resp.Header.Get("Date"))

	s.Logger.WithFields(logrus.Fields{
		"chat_id":    chatID,
		"message_id": messageID,
	}).Info(domain.TelegramNotificationSent)

	return messageID, nil
}

// SendEmailNotification sends a notification through Email
func (s *notificationSenderService) SendEmailNotification(ctx context.Context, to, subject, body string) error {
	s.Logger.WithFields(logrus.Fields{
		"to":      to,
		"subject": subject,
	}).Info(domain.LogMsgSendingEmail)

	if s.EmailConfig.SMTPHost == "" {
		err := fmt.Errorf("email SMTP configuration is not set")
		s.Logger.WithError(err).Error("Email configuration missing")
		return fmt.Errorf(domain.ErrMsgSend, resourceEmailMsg, err)
	}

	// TODO: Implement actual SMTP email sending
	// For now, just simulate success
	messageID := fmt.Sprintf("email_%s_%d", to, len(body))

	s.Logger.WithFields(logrus.Fields{
		"to":         to,
		"message_id": messageID,
	}).Info(domain.EmailNotificationSent)

	return nil
}

// SendSlackNotification sends a notification through Slack
func (s *notificationSenderService) SendSlackNotification(ctx context.Context, channel, message string) (messageID string, err error) {
	s.Logger.WithFields(logrus.Fields{
		"channel": channel,
		"message": message,
	}).Info(domain.LogMsgSendingSlack)

	if s.SlackConfig.WebhookURL == "" && s.SlackConfig.BotToken == "" {
		err = fmt.Errorf("slack configuration is not set")
		s.Logger.WithError(err).Error("Slack configuration missing")
		return "", fmt.Errorf(domain.ErrMsgSend, resourceSlackMsg, err)
	}

	// TODO: Implement actual Slack API call
	// For now, just simulate success
	messageID = fmt.Sprintf("slack_%s_%d", channel, len(message))

	s.Logger.WithFields(logrus.Fields{
		"channel":    channel,
		"message_id": messageID,
	}).Info(domain.SlackNotificationSent)

	return messageID, nil
}

// SendWebhookNotification sends a notification through Webhook
func (s *notificationSenderService) SendWebhookNotification(ctx context.Context, webhookURL, payload string) error {
	s.Logger.WithFields(logrus.Fields{
		"webhook_url": webhookURL,
		"payload":     payload,
	}).Info(domain.LogMsgSendingWebhook)

	if webhookURL == "" {
		err := fmt.Errorf("webhook URL is empty")
		s.Logger.WithError(err).Error("Webhook URL missing")
		return fmt.Errorf(domain.ErrMsgSend, resourceWebhookMsg, err)
	}

	// Create HTTP request
	req, err := http.NewRequestWithContext(ctx, "POST", webhookURL, strings.NewReader(payload))
	if err != nil {
		s.Logger.WithError(err).Error("Failed to create webhook request")
		return fmt.Errorf(domain.ErrMsgSend, resourceWebhookMsg, err)
	}

	req.Header.Set("Content-Type", "application/json")

	// Send the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		s.Logger.WithError(err).Error("Failed to send webhook request")
		return fmt.Errorf(domain.ErrMsgSend, resourceWebhookMsg, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		err = fmt.Errorf("webhook returned status: %d", resp.StatusCode)
		s.Logger.WithError(err).Error("Webhook error")
		return fmt.Errorf(domain.ErrMsgSend, resourceWebhookMsg, err)
	}

	// Generate a message ID based on response
	messageID := fmt.Sprintf("webhook_%s_%d", webhookURL, resp.StatusCode)

	s.Logger.WithFields(logrus.Fields{
		"webhook_url": webhookURL,
		"message_id":  messageID,
		"status_code": resp.StatusCode,
	}).Info(domain.WebhookNotificationSent)

	return nil
}

// escapeJSON escapes special characters in JSON strings
func escapeJSON(s string) string {
	s = strings.ReplaceAll(s, "\\", "\\\\")
	s = strings.ReplaceAll(s, "\"", "\\\"")
	s = strings.ReplaceAll(s, "\n", "\\n")
	s = strings.ReplaceAll(s, "\r", "\\r")
	s = strings.ReplaceAll(s, "\t", "\\t")
	return s
}
