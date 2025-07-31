package telegram

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/gofiber/fiber/v2"

	"github.com/dewisartika8/cicd-status-notifier-bot/internal/config"
	"github.com/dewisartika8/cicd-status-notifier-bot/internal/core/bot/domain"
	"github.com/dewisartika8/cicd-status-notifier-bot/internal/core/bot/port"
)

// TelegramAPIAdapter implements TelegramAPI interface
type TelegramAPIAdapter struct {
	bot *tgbotapi.BotAPI
}

// NewTelegramAPIAdapter creates a new Telegram API adapter
func NewTelegramAPIAdapter(cfg *config.AppConfig) port.TelegramAPI {
	botToken := cfg.Telegram.BotToken
	if botToken == "" {
		log.Fatal("Telegram bot token not configured. Please set TELEGRAM_BOT_TOKEN environment variable or configure it in config file")
	}

	bot, err := tgbotapi.NewBotAPI(botToken)
	if err != nil {
		log.Fatal(err)
	}

	bot.Debug = true
	log.Printf("Authorized on account %s", bot.Self.UserName)

	return &TelegramAPIAdapter{bot: bot}
}

// SendMessage sends a plain text message
func (t *TelegramAPIAdapter) SendMessage(chatID int64, text string) error {
	msg := tgbotapi.NewMessage(chatID, text)
	_, err := t.bot.Send(msg)
	return err
}

// SendMessageWithMarkdown sends a message with markdown formatting
func (t *TelegramAPIAdapter) SendMessageWithMarkdown(chatID int64, text string) error {
	msg := tgbotapi.NewMessage(chatID, text)
	msg.ParseMode = "Markdown"
	_, err := t.bot.Send(msg)
	return err
}

// SetWebhook sets the webhook URL
func (t *TelegramAPIAdapter) SetWebhook(webhookURL string) error {
	webhook, err := tgbotapi.NewWebhook(webhookURL + "/api/v1/telegram/webhook")
	if err != nil {
		return err
	}
	_, err = t.bot.Request(webhook)
	return err
}

// DeleteWebhook deletes the webhook
func (t *TelegramAPIAdapter) DeleteWebhook() error {
	webhook := tgbotapi.DeleteWebhookConfig{}
	_, err := t.bot.Request(webhook)
	return err
}

// WebhookHandler handles Telegram webhooks
type WebhookHandler struct {
	botService       port.BotService
	commandValidator port.CommandValidator
}

// NewWebhookHandler creates a new webhook handler
func NewWebhookHandler(botService port.BotService, commandValidator port.CommandValidator) *WebhookHandler {
	return &WebhookHandler{
		botService:       botService,
		commandValidator: commandValidator,
	}
}

// HandleTelegramWebhook handles incoming Telegram webhook
func (wh *WebhookHandler) HandleTelegramWebhook(c *fiber.Ctx) error {
	var update tgbotapi.Update

	if err := json.Unmarshal(c.Body(), &update); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid JSON payload",
		})
	}

	// Handle bot updates
	if update.Message != nil {
		if err := wh.handleCommand(update.Message); err != nil {
			log.Printf("Error handling command: %v", err)
		}
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status": "ok",
	})
}

// handleCommand processes incoming commands
func (wh *WebhookHandler) handleCommand(msg *tgbotapi.Message) error {
	if !msg.IsCommand() {
		response := "Please send a valid command. Type /help for available commands."
		return wh.botService.SendMessage(nil, msg.Chat.ID, response)
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
	return wh.botService.HandleCommand(nil, ctx)
}

// SetWebhook sets the webhook URL
func (wh *WebhookHandler) SetWebhook(webhookURL string) error {
	webhook, err := tgbotapi.NewWebhook(webhookURL + "/api/v1/telegram/webhook")
	if err != nil {
		return err
	}
	// We need access to the bot API for this operation
	// This would need to be injected or passed from the adapter
	_ = webhook // Suppress unused variable warning for now
	return fmt.Errorf("SetWebhook not implemented in webhook handler - use TelegramAPIAdapter instead")
}

// LegacyBotService provides backward compatibility with existing code
type LegacyBotService struct {
	bot        *tgbotapi.BotAPI
	botService port.BotService
}

// NewLegacyBotService creates a legacy bot service for backward compatibility
func NewLegacyBotService(botService port.BotService, cfg *config.AppConfig) *LegacyBotService {
	botToken := cfg.Telegram.BotToken
	if botToken == "" {
		log.Fatal("Telegram bot token not configured. Please set TELEGRAM_BOT_TOKEN environment variable or configure it in config file")
	}

	bot, err := tgbotapi.NewBotAPI(botToken)
	if err != nil {
		log.Fatal(err)
	}

	bot.Debug = true
	log.Printf("Authorized on account %s", bot.Self.UserName)

	return &LegacyBotService{
		bot:        bot,
		botService: botService,
	}
}

// StartTelegramBot starts the Telegram bot with polling
func (lbs *LegacyBotService) StartTelegramBot() {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := lbs.bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message != nil {
			lbs.handleCommand(update.Message)
		}
	}
}

// handleCommand processes incoming commands for legacy service
func (lbs *LegacyBotService) handleCommand(msg *tgbotapi.Message) {
	if !msg.IsCommand() {
		response := "Please send a valid command. Type /help for available commands."
		lbs.sendMessage(msg.Chat.ID, response)
		return
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
	if err := lbs.botService.HandleCommand(nil, ctx); err != nil {
		lbs.sendMessage(msg.Chat.ID, fmt.Sprintf("❌ Error: %s", err.Error()))
	}
}

// sendMessage sends a message using the legacy bot
func (lbs *LegacyBotService) sendMessage(chatID int64, text string) {
	msg := tgbotapi.NewMessage(chatID, text)
	msg.ParseMode = "Markdown"
	lbs.bot.Send(msg)
}
