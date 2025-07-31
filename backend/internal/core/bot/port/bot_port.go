package port

import (
	"context"

	"github.com/dewisartika8/cicd-status-notifier-bot/internal/core/bot/domain"
	"github.com/dewisartika8/cicd-status-notifier-bot/internal/core/bot/dto"
)

// BotService interface defines the contract for bot operations
type BotService interface {
	// Command handling
	HandleCommand(ctx context.Context, commandCtx *domain.CommandContext) error

	// Basic commands
	HandleStartCommand(ctx context.Context, req *dto.StartCommandRequest) (*dto.StartCommandResponse, error)
	HandleHelpCommand(ctx context.Context, req *HelpCommandRequest) (*dto.HelpCommandResponse, error)

	// Status commands
	HandleStatusCommand(ctx context.Context, req *dto.StatusCommandRequest) (*dto.StatusCommandResponse, error)

	// Subscription commands
	HandleSubscribeCommand(ctx context.Context, req *dto.SubscribeCommandRequest) (*dto.SubscribeCommandResponse, error)
	HandleUnsubscribeCommand(ctx context.Context, req *dto.UnsubscribeCommandRequest) (*dto.UnsubscribeCommandResponse, error)

	// Messaging
	SendMessage(ctx context.Context, chatID int64, message string) error
	SendFormattedMessage(ctx context.Context, chatID int64, message string, parseMode string) error

	// Webhook handling
	SetWebhook(ctx context.Context, webhookURL string) error
	DeleteWebhook(ctx context.Context) error
}

// TelegramAPI interface defines the contract for Telegram API operations
type TelegramAPI interface {
	SendMessage(chatID int64, text string) error
	SendMessageWithMarkdown(chatID int64, text string) error
	SetWebhook(webhookURL string) error
	DeleteWebhook() error
}

// CommandValidator interface defines the contract for command validation
type CommandValidator interface {
	ValidateCommand(ctx *domain.CommandContext) error
	AddAllowedUser(userID int64)
}

// CommandRouter interface defines the contract for command routing
type CommandRouter interface {
	RegisterHandler(command string, handler domain.CommandHandler)
	RouteCommand(ctx *domain.CommandContext) error
}

// WebhookHandler interface defines the contract for webhook handling
type WebhookHandler interface {
	HandleTelegramWebhook(ctx context.Context, payload []byte) error
	SetWebhook(ctx context.Context, webhookURL string) error
}

// ProjectService interface for project-related operations
type ProjectService interface {
	GetProject(ctx context.Context, name string) (*Project, error)
	GetAllProjects(ctx context.Context) ([]*Project, error)
	GetProjectStatus(ctx context.Context, projectName string) (*dto.StatusCommandResponse, error)
}

// SubscriptionService interface for subscription-related operations
type SubscriptionService interface {
	Subscribe(ctx context.Context, req *dto.SubscribeCommandRequest) error
	Unsubscribe(ctx context.Context, req *dto.UnsubscribeCommandRequest) error
	IsSubscribed(ctx context.Context, chatID int64, projectName string) (bool, error)
	GetSubscriptions(ctx context.Context, chatID int64) ([]*Subscription, error)
}

// Additional DTOs for interfaces
type HelpCommandRequest struct {
	ChatID int64 `json:"chat_id"`
	UserID int64 `json:"user_id"`
}

type Project struct {
	ID            string `json:"id"`
	Name          string `json:"name"`
	RepositoryURL string `json:"repository_url"`
	IsActive      bool   `json:"is_active"`
}

type Subscription struct {
	ID          string   `json:"id"`
	ProjectName string   `json:"project_name"`
	ChatID      int64    `json:"chat_id"`
	UserID      int64    `json:"user_id"`
	EventTypes  []string `json:"event_types"`
	IsActive    bool     `json:"is_active"`
}
