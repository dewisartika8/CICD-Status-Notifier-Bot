package router

import (
	"github.com/gofiber/fiber/v2"

	h "github.com/dewisartika8/cicd-status-notifier-bot/internal/adapter/handler/health"
	p "github.com/dewisartika8/cicd-status-notifier-bot/internal/adapter/handler/project"
	t "github.com/dewisartika8/cicd-status-notifier-bot/internal/adapter/handler/telegram"
	w "github.com/dewisartika8/cicd-status-notifier-bot/internal/adapter/handler/webhook"
)

type Dep struct {
	App             *fiber.App
	HealthHandler   *h.HealthHandler
	ProjectHandler  *p.Handler
	WebhookHandler  *w.WebhookHandler
	TelegramHandler *t.TelegramHandler
}

type router struct {
	Dep
}

func NewRoutes(d Dep) *router {
	return &router{
		Dep: d,
	}
}

func (r *router) RegisterRoutes() {
	// Health check routes (root level)
	root := r.App.Group("")
	r.HealthHandler.RegisterRoutes(root)

	// API v1 routes
	api := r.App.Group("/api/v1")

	// Project routes
	r.ProjectHandler.RegisterRoutes(api)

	// Telegram bot routes
	r.TelegramHandler.RegisterRoutes(api)

	// Webhook routes
	webhooks := api.Group("/webhooks")
	r.WebhookHandler.RegisterRoutes(webhooks)
}
