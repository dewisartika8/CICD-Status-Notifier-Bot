/*
File: Modul API Public Route
@author -riff-
Date : 28-02-2025
*/
package router

import (
	"github.com/gofiber/fiber/v2"

	h "github.com/dewisartika8/cicd-status-notifier-bot/internal/adapter/handler/health"
	w "github.com/dewisartika8/cicd-status-notifier-bot/internal/adapter/handler/webhook"
)

type Dep struct {
	App            *fiber.App
	HealthHandler  *h.HealthHandler
	WebhookHandler *w.WebhookHandler
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

	// Webhook routes
	webhooks := api.Group("/webhooks")
	r.WebhookHandler.RegisterRoutes(webhooks)
}
