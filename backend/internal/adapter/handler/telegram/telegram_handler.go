package telegram

import (
	"github.com/gofiber/fiber/v2"

	"github.com/dewisartika8/cicd-status-notifier-bot/internal/adapter/api"
	"github.com/dewisartika8/cicd-status-notifier-bot/internal/adapter/handler/webhook"
	"github.com/dewisartika8/cicd-status-notifier-bot/internal/config"
	"github.com/dewisartika8/cicd-status-notifier-bot/internal/core/bot/domain"
	"github.com/dewisartika8/cicd-status-notifier-bot/internal/core/bot/port"
	"github.com/dewisartika8/cicd-status-notifier-bot/internal/core/bot/service"
)

type TelegramHandler struct {
	botService     port.BotService
	webhookHandler *webhook.TelegramWebhookHandler
	telegramAPI    port.TelegramAPI
}

func NewTelegramHandler(cfg *config.AppConfig) *TelegramHandler {
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
		nil, // subscriptionService - to be implemented
	)

	// Create webhook handler
	webhookHandler := webhook.NewTelegramWebhookHandler(botService, commandValidator)

	return &TelegramHandler{
		botService:     botService,
		webhookHandler: webhookHandler,
		telegramAPI:    telegramAPI,
	}
}

func (h *TelegramHandler) RegisterRoutes(api fiber.Router) {
	telegram := api.Group("/telegram")

	// Webhook endpoint
	telegram.Post("/webhook", h.webhookHandler.HandleTelegramWebhook)

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
