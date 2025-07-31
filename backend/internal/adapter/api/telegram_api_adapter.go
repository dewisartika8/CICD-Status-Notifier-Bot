package api

import (
	"log"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

	"github.com/dewisartika8/cicd-status-notifier-bot/internal/config"
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
