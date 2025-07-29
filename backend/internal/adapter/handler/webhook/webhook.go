package webhook

import (
	"github.com/dewisartika8/cicd-status-notifier-bot/internal/core/webhook/port"
	"github.com/sirupsen/logrus"
)

// WebhookHandler handles webhook-related HTTP requests
type WebhookHandler struct {
	webhookService port.WebhookService
	logger         *logrus.Logger
}

// NewWebhookHandler creates a new webhook handler
func NewWebhookHandler(webhookService port.WebhookService, logger *logrus.Logger) *WebhookHandler {
	return &WebhookHandler{
		webhookService: webhookService,
		logger:         logger,
	}
}
