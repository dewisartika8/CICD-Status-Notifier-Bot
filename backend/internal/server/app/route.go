package app

import (
	"github.com/dewisartika8/cicd-status-notifier-bot/internal/server/router"
)

func (s *service) createRoutes() {
	router.NewRoutes(router.Dep{
		App:              s.HTTPServer,
		HealthHandler:    s.HealthHandler,
		ProjectHandler:   s.ProjectHandler,
		WebhookHandler:   s.WebhookHandler,
		TelegramHandler:  s.TelegramHandler,
		DashboardHandler: s.DashboardHandler,
	}).RegisterRoutes()
}
