package app

import (
	"github.com/dewisartika8/cicd-status-notifier-bot/internal/server/router"
)

func (s *service) createRoutes() {
	router.NewRoutes(router.Dep{
		App:            s.HTTPServer,
		HealthHandler:  s.HealthHandler,
		WebhookHandler: s.WebhookHandler,
	}).RegisterRoutes()
}
