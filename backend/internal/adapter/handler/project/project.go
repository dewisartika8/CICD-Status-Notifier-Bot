package project

import (
	"github.com/dewisartika8/cicd-status-notifier-bot/internal/core/project/port"
	"github.com/sirupsen/logrus"
)

// ProjectHandler represents the HTTP handler for project endpoints
type ProjectHandlerDep struct {
	ProjectService port.ProjectService
	Logger         *logrus.Logger
}

// Handler struct for organizing handler dependencies
type Handler struct {
	ProjectHandlerDep
}

// NewProjectHandler creates a new project handler instance
func NewProjectHandler(d ProjectHandlerDep) *Handler {
	return &Handler{
		ProjectHandlerDep: d,
	}
}
