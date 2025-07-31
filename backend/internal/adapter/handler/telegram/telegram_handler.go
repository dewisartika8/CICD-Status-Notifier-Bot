package telegram

import (
	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"

	"github.com/dewisartika8/cicd-status-notifier-bot/internal/adapter/api"
	"github.com/dewisartika8/cicd-status-notifier-bot/internal/adapter/handler/webhook"
	"github.com/dewisartika8/cicd-status-notifier-bot/internal/config"
	"github.com/dewisartika8/cicd-status-notifier-bot/internal/core/bot/domain"
	"github.com/dewisartika8/cicd-status-notifier-bot/internal/core/bot/port"
	"github.com/dewisartika8/cicd-status-notifier-bot/internal/core/bot/service"
	notificationPort "github.com/dewisartika8/cicd-status-notifier-bot/internal/core/notification/port"
)

type TelegramHandler struct {
	botService          port.BotService
	subscriptionHandler *TelegramSubscriptionHandler
	webhookHandler      *webhook.TelegramWebhookHandler
	telegramAPI         port.TelegramAPI
}

func NewTelegramHandler(cfg *config.AppConfig, subscriptionService notificationPort.TelegramSubscriptionService, logger *logrus.Logger) *TelegramHandler {
	// Initialize clean architecture components
	telegramAPI := api.NewTelegramAPIAdapter(cfg)
	commandValidator := domain.NewCommandValidator()
	commandRouter := domain.NewCommandRouter()

	// Create bot service
	botService := service.NewBotService(
		telegramAPI,
		commandValidator,
		commandRouter,
		nil, // projectService - to be implemented
		nil, // subscriptionService - separate implementation
	)

	// Create handlers
	subscriptionHandler := NewTelegramSubscriptionHandler(subscriptionService, logger)
	webhookHandler := webhook.NewTelegramWebhookHandler(botService, commandValidator)

	return &TelegramHandler{
		botService:          botService,
		subscriptionHandler: subscriptionHandler,
		webhookHandler:      webhookHandler,
		telegramAPI:         telegramAPI,
	}
}

func (h *TelegramHandler) RegisterRoutes(api fiber.Router) {
	telegram := api.Group("/telegram")

	// Webhook endpoint
	telegram.Post("/webhook", h.webhookHandler.HandleTelegramWebhook)

	// Subscription management endpoints
	subscriptions := telegram.Group("/subscriptions")
	subscriptions.Post("/", h.subscriptionHandler.CreateSubscription)
	subscriptions.Get("/:id", h.subscriptionHandler.GetSubscriptionByID)
	subscriptions.Put("/:id", h.subscriptionHandler.UpdateSubscription)
	subscriptions.Delete("/:id", h.subscriptionHandler.DeleteSubscription)
	subscriptions.Get("/active", h.subscriptionHandler.GetActiveSubscriptions)
	subscriptions.Get("/stats", h.subscriptionHandler.GetSubscriptionStats)

	// Project-specific subscription endpoints
	projects := telegram.Group("/projects")
	projects.Get("/:projectId/subscriptions", h.subscriptionHandler.GetSubscriptionsByProject)

	// Admin endpoints (for webhook management)
	telegram.Post("/webhook/set", h.setWebhook)
	telegram.Delete("/webhook", h.deleteWebhook)
}

func (h *TelegramHandler) setWebhook(c *fiber.Ctx) error {
	var req struct {
		WebhookURL string `json:"webhook_url"`
	}

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	if req.WebhookURL == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "webhook_url is required",
		})
	}

	err := h.telegramAPI.SetWebhook(req.WebhookURL)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   "Failed to set webhook",
			"details": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"message":     "Webhook set successfully",
		"webhook_url": req.WebhookURL + "/api/v1/telegram/webhook",
	})
}

func (h *TelegramHandler) deleteWebhook(c *fiber.Ctx) error {
	err := h.telegramAPI.DeleteWebhook()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   "Failed to delete webhook",
			"details": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"message": "Webhook deleted successfully",
	})
}

func (h *TelegramHandler) GetBotService() port.BotService {
	return h.botService
}
