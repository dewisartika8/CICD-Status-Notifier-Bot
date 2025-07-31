package telegram

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/gofiber/fiber/v2"
)

type BotService struct {
	bot *tgbotapi.BotAPI
}

func NewBotService() *BotService {
	botToken := os.Getenv("TELEGRAM_BOT_TOKEN")
	if botToken == "" {
		log.Fatal("TELEGRAM_BOT_TOKEN not set")
	}

	bot, err := tgbotapi.NewBotAPI(botToken)
	if err != nil {
		log.Fatal(err)
	}

	bot.Debug = true
	log.Printf("Authorized on account %s", bot.Self.UserName)

	return &BotService{bot: bot}
}

func (bs *BotService) StartTelegramBot() {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := bs.bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message != nil {
			bs.handleCommand(update.Message)
		}
	}
}

func (bs *BotService) handleCommand(msg *tgbotapi.Message) {
	if !msg.IsCommand() {
		response := "Please send a valid command. Type /help for available commands."
		bs.sendMessage(msg.Chat.ID, response)
		return
	}

	// Parse command context
	ctx := &CommandContext{
		Command:  strings.ToLower(msg.Command()),
		Args:     strings.Fields(msg.CommandArguments()),
		UserID:   msg.From.ID,
		ChatID:   msg.Chat.ID,
		Username: msg.From.UserName,
	}

	// Validate command
	validator := NewCommandValidator()
	if err := validator.ValidateCommand(ctx); err != nil {
		bs.sendMessage(msg.Chat.ID, fmt.Sprintf("❌ Error: %s", err.Error()))
		return
	}

	// Route to appropriate handler
	switch ctx.Command {
	case "start":
		bs.handleStartCommand(msg)
	case "help":
		bs.handleHelpCommand(msg)
	case "status":
		bs.handleStatusCommand(msg, ctx.Args)
	case "subscribe":
		bs.handleSubscribeCommand(msg, ctx.Args)
	case "unsubscribe":
		bs.handleUnsubscribeCommand(msg, ctx.Args)
	default:
		bs.handleUnknownCommand(msg, ctx.Command)
	}
}

func (bs *BotService) handleStartCommand(msg *tgbotapi.Message) {
	welcomeText := fmt.Sprintf(`🎉 *Welcome to CICD Status Notifier Bot!*

Hello %s! 👋

I'm here to help you monitor your CI/CD pipeline status and get real-time notifications about your builds, deployments, and more.

*Quick Start:*
• Type /help to see all available commands
• Use /subscribe to get notifications for your projects
• Check /status to see current pipeline status

Let's get started! 🚀`, msg.From.FirstName)

	bs.sendMessage(msg.Chat.ID, welcomeText)
}

func (bs *BotService) handleHelpCommand(msg *tgbotapi.Message) {
	helpText := `📚 *CICD Status Notifier Bot - Help*

*Available Commands:*

🏁 */start* - Welcome message and quick introduction
📖 */help* - Show this help message

📊 *Pipeline Commands:*
• */status* [project] - Get current pipeline status
• */status all* - Get status for all projects

🔔 *Notification Commands:*
• */subscribe* <project> - Subscribe to project notifications
• */unsubscribe* <project> - Unsubscribe from project notifications
• */list* - List your subscribed projects

*Usage Examples:*
• ` + "`/status my-app`" + ` - Get status for "my-app" project
• ` + "`/subscribe my-app`" + ` - Subscribe to "my-app" notifications
• ` + "`/unsubscribe my-app`" + ` - Unsubscribe from "my-app"

*Need more help?* Contact your system administrator.`

	bs.sendMessage(msg.Chat.ID, helpText)
}

func (bs *BotService) sendMessage(chatID int64, text string) {
	msg := tgbotapi.NewMessage(chatID, text)
	msg.ParseMode = "Markdown"
	bs.bot.Send(msg)
}

type WebhookHandler struct {
	botService *BotService
}

func NewWebhookHandler(botService *BotService) *WebhookHandler {
	return &WebhookHandler{
		botService: botService,
	}
}

func (wh *WebhookHandler) HandleTelegramWebhook(c *fiber.Ctx) error {
	var update tgbotapi.Update

	if err := json.Unmarshal(c.Body(), &update); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid JSON payload",
		})
	}

	// Handle bot updates
	if update.Message != nil {
		wh.botService.handleCommand(update.Message)
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status": "ok",
	})
}

func (wh *WebhookHandler) SetWebhook(webhookURL string) error {
	webhook := tgbotapi.NewWebhook(webhookURL + "/api/v1/telegram/webhook")
	_, err := wh.botService.bot.Request(webhook)
	return err
}

func (bs *BotService) handleStatusCommand(msg *tgbotapi.Message, args []string) {
	project := "all"
	if len(args) > 0 {
		project = args[0]
	}

	response := fmt.Sprintf("📊 *Pipeline Status for: %s*\n\n✅ Status: All systems operational", project)
	bs.sendMessage(msg.Chat.ID, response)
}

func (bs *BotService) handleSubscribeCommand(msg *tgbotapi.Message, args []string) {
	project := args[0]
	response := fmt.Sprintf("🔔 Successfully subscribed to notifications for project: *%s*", project)
	bs.sendMessage(msg.Chat.ID, response)
}

func (bs *BotService) handleUnsubscribeCommand(msg *tgbotapi.Message, args []string) {
	project := args[0]
	response := fmt.Sprintf("🔕 Successfully unsubscribed from notifications for project: *%s*", project)
	bs.sendMessage(msg.Chat.ID, response)
}

func (bs *BotService) handleUnknownCommand(msg *tgbotapi.Message, command string) {
	response := fmt.Sprintf("❓ Unknown command: *%s*\n\nType /help for available commands.", command)
	bs.sendMessage(msg.Chat.ID, response)
}

type CommandContext struct {
	Command  string
	Args     []string
	UserID   int64
	ChatID   int64
	Username string
}

type CommandValidator struct {
	allowedUsers map[int64]bool
}

func NewCommandValidator() *CommandValidator {
	return &CommandValidator{
		allowedUsers: make(map[int64]bool),
	}
}

func (cv *CommandValidator) ValidateCommand(ctx *CommandContext) error {
	// Validate command exists
	validCommands := []string{"start", "help", "status", "subscribe", "unsubscribe", "list"}
	if !cv.isValidCommand(ctx.Command, validCommands) {
		return errors.New("invalid command")
	}

	// Validate user permissions
	if !cv.hasPermission(ctx.UserID, ctx.Command) {
		return errors.New("insufficient permissions")
	}

	// Validate arguments based on command
	if err := cv.validateArguments(ctx.Command, ctx.Args); err != nil {
		return err
	}

	return nil
}

func (cv *CommandValidator) isValidCommand(command string, validCommands []string) bool {
	for _, valid := range validCommands {
		if command == valid {
			return true
		}
	}
	return false
}

func (cv *CommandValidator) hasPermission(userID int64, command string) bool {
	// For now, allow all users for basic commands
	basicCommands := []string{"start", "help", "status"}
	for _, basic := range basicCommands {
		if command == basic {
			return true
		}
	}

	// For advanced commands, check user permissions
	return cv.allowedUsers[userID]
}

func (cv *CommandValidator) validateArguments(command string, args []string) error {
	switch command {
	case "subscribe", "unsubscribe":
		if len(args) == 0 {
			return errors.New("project name is required")
		}
		if len(args[0]) < 2 {
			return errors.New("project name too short")
		}
	case "status":
		if len(args) > 1 {
			return errors.New("too many arguments for status command")
		}
	}
	return nil
}

func (cv *CommandValidator) AddAllowedUser(userID int64) {
	cv.allowedUsers[userID] = true
}
