package service

import (
	"context"
	"fmt"

	"github.com/dewisartika8/cicd-status-notifier-bot/internal/core/bot/domain"
	"github.com/dewisartika8/cicd-status-notifier-bot/internal/core/bot/dto"
	"github.com/dewisartika8/cicd-status-notifier-bot/internal/core/bot/port"
)

// BotServiceImpl implements the BotService interface
type BotServiceImpl struct {
	telegramAPI         port.TelegramAPI
	commandValidator    port.CommandValidator
	commandRouter       port.CommandRouter
	projectService      port.ProjectService
	subscriptionService port.SubscriptionService
}

// NewBotService creates a new bot service instance
func NewBotService(
	telegramAPI port.TelegramAPI,
	commandValidator port.CommandValidator,
	commandRouter port.CommandRouter,
	projectService port.ProjectService,
	subscriptionService port.SubscriptionService,
) port.BotService {
	service := &BotServiceImpl{
		telegramAPI:         telegramAPI,
		commandValidator:    commandValidator,
		commandRouter:       commandRouter,
		projectService:      projectService,
		subscriptionService: subscriptionService,
	}

	// Register command handlers
	service.registerHandlers()

	return service
}

// registerHandlers registers all command handlers with the router
func (bs *BotServiceImpl) registerHandlers() {
	bs.commandRouter.RegisterHandler("start", &startCommandHandler{botService: bs})
	bs.commandRouter.RegisterHandler("help", &helpCommandHandler{botService: bs})
	bs.commandRouter.RegisterHandler("status", &statusCommandHandler{botService: bs})
	bs.commandRouter.RegisterHandler("subscribe", &subscribeCommandHandler{botService: bs})
	bs.commandRouter.RegisterHandler("unsubscribe", &unsubscribeCommandHandler{botService: bs})
}

// HandleCommand handles incoming commands
func (bs *BotServiceImpl) HandleCommand(ctx context.Context, commandCtx *domain.CommandContext) error {
	// Validate command
	if err := bs.commandValidator.ValidateCommand(commandCtx); err != nil {
		errorMsg := fmt.Sprintf("❌ Error: %s", err.Error())
		return bs.SendMessage(ctx, commandCtx.ChatID, errorMsg)
	}

	// Route to appropriate handler
	return bs.commandRouter.RouteCommand(commandCtx)
}

// HandleStartCommand handles /start command
func (bs *BotServiceImpl) HandleStartCommand(ctx context.Context, req *dto.StartCommandRequest) (*dto.StartCommandResponse, error) {
	welcomeText := fmt.Sprintf(`🎉 *Welcome to CICD Status Notifier Bot!*

Hello %s! 👋

I'm here to help you monitor your CI/CD pipeline status and get real-time notifications about your builds, deployments, and more.

*Quick Start:*
• Type /help to see all available commands
• Use /subscribe to get notifications for your projects
• Check /status to see current pipeline status

Let's get started! 🚀`, req.UserFirstName)

	if err := bs.SendFormattedMessage(ctx, req.ChatID, welcomeText, "Markdown"); err != nil {
		return nil, fmt.Errorf("failed to send welcome message: %w", err)
	}

	return &dto.StartCommandResponse{
		WelcomeMessage: welcomeText,
		UserFirstName:  req.UserFirstName,
		ChatID:         req.ChatID,
	}, nil
}

// HandleHelpCommand handles /help command
func (bs *BotServiceImpl) HandleHelpCommand(ctx context.Context, req *port.HelpCommandRequest) (*dto.HelpCommandResponse, error) {
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

	if err := bs.SendFormattedMessage(ctx, req.ChatID, helpText, "Markdown"); err != nil {
		return nil, fmt.Errorf("failed to send help message: %w", err)
	}

	commands := []dto.CommandHelp{
		{Command: "/start", Description: "Welcome message and quick introduction", Usage: "/start", Category: "Basic"},
		{Command: "/help", Description: "Show this help message", Usage: "/help", Category: "Basic"},
		{Command: "/status", Description: "Get current pipeline status", Usage: "/status [project]", Category: "Pipeline"},
		{Command: "/subscribe", Description: "Subscribe to project notifications", Usage: "/subscribe <project>", Category: "Notification"},
		{Command: "/unsubscribe", Description: "Unsubscribe from project notifications", Usage: "/unsubscribe <project>", Category: "Notification"},
	}

	examples := []dto.UsageExample{
		{Command: "/status my-app", Description: "Get status for 'my-app' project"},
		{Command: "/subscribe my-app", Description: "Subscribe to 'my-app' notifications"},
		{Command: "/unsubscribe my-app", Description: "Unsubscribe from 'my-app'"},
	}

	return &dto.HelpCommandResponse{
		HelpText:      helpText,
		Commands:      commands,
		UsageExamples: examples,
	}, nil
}

// HandleStatusCommand handles /status command
func (bs *BotServiceImpl) HandleStatusCommand(ctx context.Context, req *dto.StatusCommandRequest) (*dto.StatusCommandResponse, error) {
	projectName := "all"
	if req.ProjectName != "" {
		projectName = req.ProjectName
	}

	// For now, return a mock response. This will be connected to real project service later
	response := fmt.Sprintf("📊 *Pipeline Status for: %s*\n\n✅ Status: All systems operational", projectName)

	if err := bs.SendFormattedMessage(ctx, req.ChatID, response, "Markdown"); err != nil {
		return nil, fmt.Errorf("failed to send status message: %w", err)
	}

	return &dto.StatusCommandResponse{
		ProjectName: projectName,
		Status:      "operational",
		Message:     response,
	}, nil
}

// HandleSubscribeCommand handles /subscribe command
func (bs *BotServiceImpl) HandleSubscribeCommand(ctx context.Context, req *dto.SubscribeCommandRequest) (*dto.SubscribeCommandResponse, error) {
	// For now, return a mock response. This will be connected to real subscription service later
	response := fmt.Sprintf("🔔 Successfully subscribed to notifications for project: *%s*", req.ProjectName)

	if err := bs.SendFormattedMessage(ctx, req.ChatID, response, "Markdown"); err != nil {
		return nil, fmt.Errorf("failed to send subscription message: %w", err)
	}

	return &dto.SubscribeCommandResponse{
		ProjectName: req.ProjectName,
		Success:     true,
		Message:     response,
	}, nil
}

// HandleUnsubscribeCommand handles /unsubscribe command
func (bs *BotServiceImpl) HandleUnsubscribeCommand(ctx context.Context, req *dto.UnsubscribeCommandRequest) (*dto.UnsubscribeCommandResponse, error) {
	// For now, return a mock response. This will be connected to real subscription service later
	response := fmt.Sprintf("🔕 Successfully unsubscribed from notifications for project: *%s*", req.ProjectName)

	if err := bs.SendFormattedMessage(ctx, req.ChatID, response, "Markdown"); err != nil {
		return nil, fmt.Errorf("failed to send unsubscription message: %w", err)
	}

	return &dto.UnsubscribeCommandResponse{
		ProjectName: req.ProjectName,
		Success:     true,
		Message:     response,
	}, nil
}

// SendMessage sends a plain text message
func (bs *BotServiceImpl) SendMessage(ctx context.Context, chatID int64, message string) error {
	return bs.telegramAPI.SendMessage(chatID, message)
}

// SendFormattedMessage sends a formatted message
func (bs *BotServiceImpl) SendFormattedMessage(ctx context.Context, chatID int64, message string, parseMode string) error {
	if parseMode == "Markdown" {
		return bs.telegramAPI.SendMessageWithMarkdown(chatID, message)
	}
	return bs.telegramAPI.SendMessage(chatID, message)
}

// SetWebhook sets the webhook URL
func (bs *BotServiceImpl) SetWebhook(ctx context.Context, webhookURL string) error {
	return bs.telegramAPI.SetWebhook(webhookURL)
}

// DeleteWebhook deletes the webhook
func (bs *BotServiceImpl) DeleteWebhook(ctx context.Context) error {
	return bs.telegramAPI.DeleteWebhook()
}

// Command handlers implementation
type startCommandHandler struct {
	botService *BotServiceImpl
}

func (h *startCommandHandler) Handle(ctx *domain.CommandContext) error {
	req := &dto.StartCommandRequest{
		ChatID:        ctx.ChatID,
		UserFirstName: ctx.Username, // Using username as firstname for now
	}
	_, err := h.botService.HandleStartCommand(context.Background(), req)
	return err
}

type helpCommandHandler struct {
	botService *BotServiceImpl
}

func (h *helpCommandHandler) Handle(ctx *domain.CommandContext) error {
	req := &port.HelpCommandRequest{
		ChatID: ctx.ChatID,
		UserID: ctx.UserID,
	}
	_, err := h.botService.HandleHelpCommand(context.Background(), req)
	return err
}

type statusCommandHandler struct {
	botService *BotServiceImpl
}

func (h *statusCommandHandler) Handle(ctx *domain.CommandContext) error {
	projectName := ""
	if len(ctx.Args) > 0 {
		projectName = ctx.Args[0]
	}

	req := &dto.StatusCommandRequest{
		ProjectName: projectName,
		ChatID:      ctx.ChatID,
		UserID:      ctx.UserID,
	}
	_, err := h.botService.HandleStatusCommand(context.Background(), req)
	return err
}

type subscribeCommandHandler struct {
	botService *BotServiceImpl
}

func (h *subscribeCommandHandler) Handle(ctx *domain.CommandContext) error {
	if len(ctx.Args) == 0 {
		return fmt.Errorf("project name is required")
	}

	req := &dto.SubscribeCommandRequest{
		ProjectName: ctx.Args[0],
		ChatID:      ctx.ChatID,
		UserID:      ctx.UserID,
		Username:    ctx.Username,
	}
	_, err := h.botService.HandleSubscribeCommand(context.Background(), req)
	return err
}

type unsubscribeCommandHandler struct {
	botService *BotServiceImpl
}

func (h *unsubscribeCommandHandler) Handle(ctx *domain.CommandContext) error {
	if len(ctx.Args) == 0 {
		return fmt.Errorf("project name is required")
	}

	req := &dto.UnsubscribeCommandRequest{
		ProjectName: ctx.Args[0],
		ChatID:      ctx.ChatID,
		UserID:      ctx.UserID,
	}
	_, err := h.botService.HandleUnsubscribeCommand(context.Background(), req)
	return err
}
