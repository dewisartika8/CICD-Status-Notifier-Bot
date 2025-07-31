package webhook

import (
	"encoding/json"
	"log"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/gofiber/fiber/v2"

	"github.com/dewisartika8/cicd-status-notifier-bot/internal/core/bot/domain"
	"github.com/dewisartika8/cicd-status-notifier-bot/internal/core/bot/port"
)

// TelegramWebhookHandler handles Telegram webhooks
type TelegramWebhookHandler struct {
	botService       port.BotService
	commandValidator port.CommandValidator
}

// NewTelegramWebhookHandler creates a new webhook handler
func NewTelegramWebhookHandler(botService port.BotService, commandValidator port.CommandValidator) *TelegramWebhookHandler {
	return &TelegramWebhookHandler{
		botService:       botService,
		commandValidator: commandValidator,
	}
}

// HandleTelegramWebhook handles incoming Telegram webhook
func (h *TelegramWebhookHandler) HandleTelegramWebhook(c *fiber.Ctx) error {
	var update tgbotapi.Update

	if err := json.Unmarshal(c.Body(), &update); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid JSON payload",
		})
	}

	// Handle bot updates
	if update.Message != nil {
		if err := h.handleCommand(update.Message); err != nil {
			log.Printf("Error handling command: %v", err)
		}
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status": "ok",
	})
}

// handleCommand processes incoming commands
func (h *TelegramWebhookHandler) handleCommand(msg *tgbotapi.Message) error {
	if !msg.IsCommand() {
		response := "Please send a valid command. Type /help for available commands."
		return h.botService.SendMessage(nil, msg.Chat.ID, response)
	}

	// Parse command context
	ctx := &domain.CommandContext{
		Command:  strings.ToLower(msg.Command()),
		Args:     strings.Fields(msg.CommandArguments()),
		UserID:   msg.From.ID,
		ChatID:   msg.Chat.ID,
		Username: msg.From.UserName,
	}

	// Handle command through bot service
	return h.botService.HandleCommand(nil, ctx)
}
