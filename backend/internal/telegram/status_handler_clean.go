package telegram

import (
	"context"
	"fmt"
	"strings"

	"github.com/dewisartika8/cicd-status-notifier-bot/internal/core/project/domain"
	"github.com/dewisartika8/cicd-status-notifier-bot/internal/core/project/port"
)

type StatusCommandHandler struct {
	projectService port.ProjectService
}

func NewStatusCommandHandler(projectService port.ProjectService) *StatusCommandHandler {
	return &StatusCommandHandler{
		projectService: projectService,
	}
}

// Task 2.3.1: Implement /status command for all projects
func (h *StatusCommandHandler) HandleStatusAllProjects() (string, error) {
	ctx := context.Background()
	projects, err := h.projectService.GetActiveProjects(ctx)
	if err != nil {
		return "❌ **Error fetching project status**\n\n" +
			"Unable to retrieve project information at the moment. Please try again later.", nil
	}

	if len(projects) == 0 {
		return "📊 **Overall Project Status**\n\n" +
			"ℹ️ No projects are currently being monitored.\n\n" +
			"Use `/projects add <name>` to start monitoring a project.", nil
	}

	var response strings.Builder
	response.WriteString("📊 **Overall Project Status**\n\n")

	successCount := 0
	failedCount := 0

	for _, project := range projects {
		var statusIcon string
		var statusText string

		switch project.Status() {
		case domain.ProjectStatusActive:
			statusIcon = "✅"
			statusText = "Active"
			successCount++
		case domain.ProjectStatusInactive:
			statusIcon = "❌"
			statusText = "Inactive"
			failedCount++
		case domain.ProjectStatusArchived:
			statusIcon = "📦"
			statusText = "Archived"
		default:
			statusIcon = "⚪"
			statusText = "Unknown"
		}

		// Format timestamp
		lastUpdated := project.UpdatedAt().Value().Format("2006-01-02 15:04 UTC")

		response.WriteString(fmt.Sprintf("%s **%s**: %s\n", statusIcon, project.Name(), statusText))
		response.WriteString(fmt.Sprintf("   🕐 Last updated: %s\n\n", lastUpdated))
	}

	// Add summary
	response.WriteString(fmt.Sprintf("📈 **Summary**: %d projects (%d active, %d inactive)",
		len(projects), successCount, failedCount))

	return response.String(), nil
}

// Task 2.3.2: Implement /status <project> for specific project
func (h *StatusCommandHandler) HandleStatusSpecificProject(projectName string) (string, error) {
	ctx := context.Background()

	if strings.TrimSpace(projectName) == "" {
		return "❌ **Invalid command**\n\n" +
			"Please specify a project name.\n" +
			"Usage: `/status <project_name>`", nil
	}

	project, err := h.projectService.GetProjectByName(ctx, projectName)
	if err != nil {
		return fmt.Sprintf("❌ **Project not found**\n\n"+
			"Project `%s` is not being monitored.\n\n"+
			"Use `/projects` to see all monitored projects.", projectName), nil
	}

	var statusIcon string
	var statusText string

	switch project.Status() {
	case domain.ProjectStatusActive:
		statusIcon = "✅"
		statusText = "Active"
	case domain.ProjectStatusInactive:
		statusIcon = "❌"
		statusText = "Inactive"
	case domain.ProjectStatusArchived:
		statusIcon = "📦"
		statusText = "Archived"
	default:
		statusIcon = "⚪"
		statusText = "Unknown"
	}

	// Format timestamps
	createdAt := project.CreatedAt().Value().Format("2006-01-02 15:04 UTC")
	lastUpdated := project.UpdatedAt().Value().Format("2006-01-02 15:04 UTC")

	var response strings.Builder
	response.WriteString(fmt.Sprintf("📊 **Project Status: %s**\n\n", project.Name()))
	response.WriteString(fmt.Sprintf("%s **Status**: %s\n", statusIcon, statusText))
	response.WriteString(fmt.Sprintf("🌐 **Repository**: %s\n", project.RepositoryURL()))

	if project.TelegramChatID() != nil {
		response.WriteString(fmt.Sprintf("💬 **Telegram Chat**: %d\n", *project.TelegramChatID()))
	} else {
		response.WriteString("💬 **Telegram Chat**: Not configured\n")
	}

	response.WriteString(fmt.Sprintf("📅 **Created**: %s\n", createdAt))
	response.WriteString(fmt.Sprintf("🕐 **Last Updated**: %s\n\n", lastUpdated))

	// Add recent activity placeholder (will be implemented later with build history)
	response.WriteString("📈 **Recent Activity**: No build data available yet")

	return response.String(), nil
}

// Task 2.3.3: Implement /projects command
func (h *StatusCommandHandler) HandleProjectsList() (string, error) {
	ctx := context.Background()
	projects, err := h.projectService.GetActiveProjects(ctx)
	if err != nil {
		return "❌ **Error fetching projects**\n\n" +
			"Unable to retrieve project list at the moment. Please try again later.", nil
	}

	if len(projects) == 0 {
		return "📋 **Monitored Projects**\n\n" +
			"ℹ️ No projects are currently being monitored.\n\n" +
			"Contact your administrator to add projects.", nil
	}

	var response strings.Builder
	response.WriteString("📋 **Monitored Projects**\n\n")

	// Group projects by status
	activeProjects := []*domain.Project{}
	inactiveProjects := []*domain.Project{}
	archivedProjects := []*domain.Project{}

	for _, project := range projects {
		switch project.Status() {
		case domain.ProjectStatusActive:
			activeProjects = append(activeProjects, project)
		case domain.ProjectStatusInactive:
			inactiveProjects = append(inactiveProjects, project)
		case domain.ProjectStatusArchived:
			archivedProjects = append(archivedProjects, project)
		}
	}

	// Display active projects
	if len(activeProjects) > 0 {
		response.WriteString("✅ **Active Projects:**\n")
		for _, project := range activeProjects {
			chatStatus := "No"
			if project.TelegramChatID() != nil {
				chatStatus = "Yes"
			}
			response.WriteString(fmt.Sprintf("   • `%s` - Notifications: %s\n", project.Name(), chatStatus))
		}
		response.WriteString("\n")
	}

	// Display inactive projects
	if len(inactiveProjects) > 0 {
		response.WriteString("❌ **Inactive Projects:**\n")
		for _, project := range inactiveProjects {
			response.WriteString(fmt.Sprintf("   • `%s`\n", project.Name()))
		}
		response.WriteString("\n")
	}

	// Display archived projects
	if len(archivedProjects) > 0 {
		response.WriteString("📦 **Archived Projects:**\n")
		for _, project := range archivedProjects {
			response.WriteString(fmt.Sprintf("   • `%s`\n", project.Name()))
		}
		response.WriteString("\n")
	}

	// Add summary and usage instructions
	response.WriteString(fmt.Sprintf("📊 **Total**: %d projects (%d active, %d inactive, %d archived)\n\n",
		len(projects), len(activeProjects), len(inactiveProjects), len(archivedProjects)))

	response.WriteString("💡 **Quick Commands:**\n")
	response.WriteString("   • `/status <project>` - Get detailed project status\n")
	response.WriteString("   • `/subscribe <project>` - Subscribe to notifications\n")
	response.WriteString("   • `/unsubscribe <project>` - Unsubscribe from notifications")

	return response.String(), nil
}
