package dto

import (
	"github.com/dewisartika8/cicd-status-notifier-bot/internal/core/project/domain"
	"github.com/dewisartika8/cicd-status-notifier-bot/internal/core/shared/domain/value_objects"
)

// CreateProjectRequest represents the request to create a new project
type CreateProjectRequest struct {
	Name           string `json:"name" validate:"required,min=1,max=100"`
	RepositoryURL  string `json:"repository_url" validate:"required,url"`
	WebhookSecret  string `json:"webhook_secret" validate:"required,min=10"`
	TelegramChatID *int64 `json:"telegram_chat_id,omitempty"`
}

// UpdateProjectRequest represents the request to update a project
type UpdateProjectRequest struct {
	Name           *string `json:"name,omitempty" validate:"omitempty,min=1,max=100"`
	RepositoryURL  *string `json:"repository_url,omitempty" validate:"omitempty,url"`
	WebhookSecret  *string `json:"webhook_secret,omitempty" validate:"omitempty,min=10"`
	TelegramChatID *int64  `json:"telegram_chat_id,omitempty"`
}

// ProjectResponse represents the response containing project information
type ProjectResponse struct {
	ID             string                  `json:"id"`
	Name           string                  `json:"name"`
	RepositoryURL  string                  `json:"repository_url"`
	Status         string                  `json:"status"`
	TelegramChatID *int64                  `json:"telegram_chat_id,omitempty"`
	CreatedAt      value_objects.Timestamp `json:"created_at"`
	UpdatedAt      value_objects.Timestamp `json:"updated_at"`
}

// ListProjectFilters represents filters for listing projects
type ListProjectFilters struct {
	Status          *domain.ProjectStatus `json:"status,omitempty"`
	Name            *string               `json:"name,omitempty"`
	RepositoryURL   *string               `json:"repository_url,omitempty"`
	HasTelegramChat *bool                 `json:"has_telegram_chat,omitempty"`
	Limit           *int                  `json:"limit,omitempty" validate:"omitempty,min=1,max=100"`
	Offset          *int                  `json:"offset,omitempty" validate:"omitempty,min=0"`
	SortBy          *string               `json:"sort_by,omitempty" validate:"omitempty,oneof=name repository_url status created_at updated_at"`
	SortOrder       *string               `json:"sort_order,omitempty" validate:"omitempty,oneof=asc desc"`
}

// ListProjectResponse represents the response for listing projects
type ListProjectResponse struct {
	Projects []*ProjectResponse `json:"projects"`
	Total    int64              `json:"total"`
	Limit    int                `json:"limit"`
	Offset   int                `json:"offset"`
}

// ProjectStatusUpdateRequest represents the request to update project status
type ProjectStatusUpdateRequest struct {
	Status domain.ProjectStatus `json:"status" validate:"required,oneof=active inactive archived"`
}

// ToProjectResponse converts a domain project to a response DTO
func ToProjectResponse(project *domain.Project) *ProjectResponse {
	return &ProjectResponse{
		ID:             project.ID().String(),
		Name:           project.Name(),
		RepositoryURL:  project.RepositoryURL(),
		Status:         string(project.Status()),
		TelegramChatID: project.TelegramChatID(),
		CreatedAt:      project.CreatedAt(),
		UpdatedAt:      project.UpdatedAt(),
	}
}

// ToProjectResponseList converts a slice of domain projects to response DTOs
func ToProjectResponseList(projects []*domain.Project) []*ProjectResponse {
	responses := make([]*ProjectResponse, len(projects))
	for i, project := range projects {
		responses[i] = ToProjectResponse(project)
	}
	return responses
}
