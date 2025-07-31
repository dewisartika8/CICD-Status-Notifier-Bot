package app

import (
	"context"
	"log"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

	"github.com/dewisartika8/cicd-status-notifier-bot/internal/config"
	"github.com/dewisartika8/cicd-status-notifier-bot/internal/core/bot/domain"
	"github.com/dewisartika8/cicd-status-notifier-bot/internal/core/bot/port"
)

// TelegramBotManager manages the Telegram bot polling
type TelegramBotManager struct {
	bot        *tgbotapi.BotAPI
	botService port.BotService
	logger     interface{}
}

// NewTelegramBotManager creates a new Telegram bot manager
func NewTelegramBotManager(botService port.BotService, cfg *config.AppConfig) *TelegramBotManager {
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

	return &TelegramBotManager{
		bot:        bot,
		botService: botService,
	}
}

// StartTelegramBot starts the Telegram bot with polling
func (tbm *TelegramBotManager) StartTelegramBot(ctx context.Context) {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := tbm.bot.GetUpdatesChan(u)

	log.Println("Telegram bot started with polling...")

	for {
		select {
		case <-ctx.Done():
			log.Println("Telegram bot polling stopped")
			return
		case update := <-updates:
			if update.Message != nil {
				go tbm.handleCommand(update.Message)
			}
		}
	}
}

// handleCommand processes incoming commands
func (tbm *TelegramBotManager) handleCommand(msg *tgbotapi.Message) {
	if !msg.IsCommand() {
		response := "Please send a valid command. Type /help for available commands."
		tbm.sendMessage(msg.Chat.ID, response)
		return
	}

	// Parse command context
	commandCtx := &domain.CommandContext{
		Command:  strings.ToLower(msg.Command()),
		Args:     strings.Fields(msg.CommandArguments()),
		UserID:   msg.From.ID,
		ChatID:   msg.Chat.ID,
		Username: msg.From.UserName,
	}

	// Handle command through bot service
	if err := tbm.botService.HandleCommand(context.Background(), commandCtx); err != nil {
		log.Printf("Error handling command: %v", err)
		tbm.sendMessage(msg.Chat.ID, "❌ Error processing command. Please try again.")
	}
}

// sendMessage sends a message using the bot
func (tbm *TelegramBotManager) sendMessage(chatID int64, text string) {
	msg := tgbotapi.NewMessage(chatID, text)
	msg.ParseMode = "Markdown"
	if _, err := tbm.bot.Send(msg); err != nil {
		log.Printf("Error sending message: %v", err)
	}
}

// StopBot stops the bot gracefully
func (tbm *TelegramBotManager) StopBot() {
	tbm.bot.StopReceivingUpdates()
}
