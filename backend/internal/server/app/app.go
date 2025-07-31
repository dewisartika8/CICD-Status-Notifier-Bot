package app

import (
	h "github.com/dewisartika8/cicd-status-notifier-bot/internal/adapter/handler/health"
	p "github.com/dewisartika8/cicd-status-notifier-bot/internal/adapter/handler/project"
	t "github.com/dewisartika8/cicd-status-notifier-bot/internal/adapter/handler/telegram"
	w "github.com/dewisartika8/cicd-status-notifier-bot/internal/adapter/handler/webhook"
	"github.com/dewisartika8/cicd-status-notifier-bot/internal/config"
	"github.com/dewisartika8/cicd-status-notifier-bot/internal/server/middleware"
	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
)

// Error messages
const (
	ErrorInternalServerError = "Internal Server Error"
)

type Dep struct {
	AppConfig       *config.AppConfig
	HealthHandler   *h.HealthHandler
	ProjectHandler  *p.Handler
	WebhookHandler  *w.WebhookHandler
	TelegramHandler *t.TelegramHandler
	Logger          *logrus.Logger
}

type service struct {
	Dep
	HTTPServer         *fiber.App
	TelegramBotManager *TelegramBotManager
}

func Init(d Dep) *service {
	// Create Telegram bot manager
	var telegramBotManager *TelegramBotManager
	if d.TelegramHandler != nil {
		telegramBotManager = NewTelegramBotManager(d.TelegramHandler.GetBotService(), d.AppConfig)
	}

	// Create Fiber app
	return &service{
		Dep: d,
		HTTPServer: fiber.New(fiber.Config{
			AppName: "CI/CD Status Notifier Bot",
			ErrorHandler: func(c *fiber.Ctx, err error) error {
				d.Logger.Errorf("Request error: %v", err)
				return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
					"error": ErrorInternalServerError,
				})
			},
		}),
		TelegramBotManager: telegramBotManager,
	}
}

func (s *service) Run() {
	// Middlewares.
	middleware.FiberMiddleware(s.HTTPServer, s.Logger) // Register Fiber's middleware for app.

	s.createRoutes()

	// Start the server
	s.startServerWithGracefulShutdown()
}
