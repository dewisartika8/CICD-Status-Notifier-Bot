/*
File: Modul API Public Route
@author -riff-
Date : 28-02-2025
*/
package router

import (
	"github.com/gofiber/fiber/v2"

	h "github.com/dewisartika8/cicd-status-notifier-bot/internal/adapter/handler/health"
)

type Dep struct {
	App           *fiber.App
	HealthHandler *h.HealthHandler
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
	root := r.App.Group("")
	r.HealthHandler.RegisterRoutes(root)
}
