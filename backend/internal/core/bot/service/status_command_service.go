package service

import (
	"context"
	"fmt"
	"strings"

	projectDomain "github.com/dewisartika8/cicd-status-notifier-bot/internal/core/project/domain"
	projectPort "github.com/dewisartika8/cicd-status-notifier-bot/internal/core/project/port"
)

// StatusCommandService handles status-related commands in Clean Architecture style
type StatusCommandService struct {
	projectService projectPort.ProjectService
}

// NewStatusCommandService creates a new status command service
func NewStatusCommandService(projectService projectPort.ProjectService) *StatusCommandService {
	return &StatusCommandService{
		projectService: projectService,
	}
}

// HandleStatusAllProjects handles the /status command for all projects
func (s *StatusCommandService) HandleStatusAllProjects() (string, error) {
	ctx := context.Background()
	projects, err := s.projectService.GetActiveProjects(ctx)
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
		case projectDomain.ProjectStatusActive:
			statusIcon = "✅"
			statusText = "Active"
			successCount++
		case projectDomain.ProjectStatusInactive:
			statusIcon = "🟡"
			statusText = "Inactive"
		case projectDomain.ProjectStatusArchived:
			statusIcon = "📦"
			statusText = "Archived"
		default:
			statusIcon = "❓"
			statusText = "Unknown"
		}

		// Get chat status
		chatStatus := "No"
		if project.TelegramChatID() != nil {
			chatStatus = "Yes"
		}

		response.WriteString(fmt.Sprintf("• **%s** %s _%s_\n", project.Name(), statusIcon, statusText))
		response.WriteString(fmt.Sprintf("  📦 Repository: %s\n", project.RepositoryURL()))
		response.WriteString(fmt.Sprintf("  📞 Notifications: %s\n\n", chatStatus))
	}

	// Summary
	totalProjects := len(projects)
	response.WriteString(fmt.Sprintf("📈 **Summary:**\n"))
	response.WriteString(fmt.Sprintf("   • Total Projects: %d\n", totalProjects))
	response.WriteString(fmt.Sprintf("   • Active: %d\n", successCount))
	response.WriteString(fmt.Sprintf("   • Issues: %d\n", failedCount))

	// Add recent activity placeholder
	response.WriteString("📈 **Recent Activity**: No build data available yet")

	return response.String(), nil
}

// HandleStatusSpecificProject handles the /status command for a specific project
func (s *StatusCommandService) HandleStatusSpecificProject(projectName string) (string, error) {
	if strings.TrimSpace(projectName) == "" {
		return "❌ **Invalid command**\n\n" +
			"Please specify a project name.\n\n" +
			"*Usage:* `/status <project-name>`\n" +
			"*Example:* `/status my-awesome-app`", nil
	}

	ctx := context.Background()
	project, err := s.projectService.GetProjectByName(ctx, projectName)
	if err != nil {
		return fmt.Sprintf("❌ **Project not found**\n\n"+
			"The project `%s` was not found in the system.\n\n"+
			"Use `/projects` to see available projects.", projectName), nil
	}

	var statusIcon string
	var statusText string

	switch project.Status() {
	case projectDomain.ProjectStatusActive:
		statusIcon = "✅"
		statusText = "Active"
	case projectDomain.ProjectStatusInactive:
		statusIcon = "🟡"
		statusText = "Inactive"
	case projectDomain.ProjectStatusArchived:
		statusIcon = "📦"
		statusText = "Archived"
	default:
		statusIcon = "❓"
		statusText = "Unknown"
	}

	// Get chat status
	chatStatus := "No"
	if project.TelegramChatID() != nil {
		chatStatus = "Yes"
	}

	var response strings.Builder
	response.WriteString(fmt.Sprintf("📊 **Project Status: %s**\n\n", project.Name()))
	response.WriteString(fmt.Sprintf("**Status:** %s %s\n", statusIcon, statusText))
	response.WriteString(fmt.Sprintf("**Repository:** %s\n", project.RepositoryURL()))
	response.WriteString(fmt.Sprintf("**Notifications:** %s\n", chatStatus))
	response.WriteString(fmt.Sprintf("**Created:** %s\n\n", project.CreatedAt().ToTime().Format("2006-01-02 15:04:05")))

	// Add build status placeholder (will be implemented later with build service)
	response.WriteString("🔧 **Latest Build Status**: No build data available yet\n\n")

	// Add quick actions
	response.WriteString("🚀 **Quick Actions:**\n")
	if project.TelegramChatID() == nil {
		response.WriteString("• Use `/subscribe` to get notifications\n")
	} else {
		response.WriteString("• Use `/unsubscribe` to stop notifications\n")
	}
	response.WriteString("• Use `/projects` to see all projects")

	// Add recent activity placeholder
	response.WriteString("📈 **Recent Activity**: No build data available yet")

	return response.String(), nil
}

// Constants for message templates
const (
	errorFetchingProjectsMsg = "❌ **Error fetching projects**\n\n" +
		"Unable to retrieve project list at the moment. Please try again later."

	noProjectsMsg = "📋 **Monitored Projects**\n\n" +
		"ℹ️ No projects are currently being monitored.\n\n" +
		"Contact your administrator to add projects."

	projectsHeaderMsg = "📋 **Monitored Projects**\n\n"

	quickCommandsMsg = "🚀 **Quick Commands:**\n" +
		"• `/status` - Overall status\n" +
		"• `/status <project>` - Specific project status\n" +
		"• `/subscribe <project>` - Get notifications\n" +
		"• `/help` - Show all commands"
)

// ProjectGroup represents projects grouped by status
type ProjectGroup struct {
	Active   []*projectDomain.Project
	Inactive []*projectDomain.Project
	Archived []*projectDomain.Project
}

// HandleProjectsList handles the /projects command
func (s *StatusCommandService) HandleProjectsList() (string, error) {
	projects, err := s.fetchProjects()
	if err != nil {
		return errorFetchingProjectsMsg, nil
	}

	if len(projects) == 0 {
		return noProjectsMsg, nil
	}

	return s.buildProjectsListResponse(projects), nil
}

// fetchProjects retrieves all active projects
func (s *StatusCommandService) fetchProjects() ([]*projectDomain.Project, error) {
	ctx := context.Background()
	return s.projectService.GetActiveProjects(ctx)
}

// buildProjectsListResponse constructs the complete response message
func (s *StatusCommandService) buildProjectsListResponse(projects []*projectDomain.Project) string {
	var response strings.Builder
	response.WriteString(projectsHeaderMsg)

	projectGroups := s.groupProjectsByStatus(projects)

	s.appendProjectsByStatus(&response, projectGroups)
	s.appendSummary(&response, len(projects))
	s.appendQuickCommands(&response)

	return response.String()
}

// groupProjectsByStatus groups projects by their status
func (s *StatusCommandService) groupProjectsByStatus(projects []*projectDomain.Project) *ProjectGroup {
	groups := &ProjectGroup{
		Active:   make([]*projectDomain.Project, 0),
		Inactive: make([]*projectDomain.Project, 0),
		Archived: make([]*projectDomain.Project, 0),
	}

	for _, project := range projects {
		switch project.Status() {
		case projectDomain.ProjectStatusActive:
			groups.Active = append(groups.Active, project)
		case projectDomain.ProjectStatusInactive:
			groups.Inactive = append(groups.Inactive, project)
		case projectDomain.ProjectStatusArchived:
			groups.Archived = append(groups.Archived, project)
		}
	}

	return groups
}

// appendProjectsByStatus adds projects to response grouped by status
func (s *StatusCommandService) appendProjectsByStatus(response *strings.Builder, groups *ProjectGroup) {
	s.appendActiveProjects(response, groups.Active)
	s.appendInactiveProjects(response, groups.Inactive)
	s.appendArchivedProjects(response, groups.Archived)
}

// appendActiveProjects adds active projects to the response
func (s *StatusCommandService) appendActiveProjects(response *strings.Builder, projects []*projectDomain.Project) {
	if len(projects) == 0 {
		return
	}

	response.WriteString("✅ **Active Projects:**\n")
	for _, project := range projects {
		notificationStatus := s.getNotificationStatus(project)
		response.WriteString(fmt.Sprintf("   • `%s` - Notifications: %s\n", project.Name(), notificationStatus))
	}
	response.WriteString("\n")
}

// appendInactiveProjects adds inactive projects to the response
func (s *StatusCommandService) appendInactiveProjects(response *strings.Builder, projects []*projectDomain.Project) {
	if len(projects) == 0 {
		return
	}

	response.WriteString("🟡 **Inactive Projects:**\n")
	for _, project := range projects {
		response.WriteString(fmt.Sprintf("   • `%s`\n", project.Name()))
	}
	response.WriteString("\n")
}

// appendArchivedProjects adds archived projects to the response
func (s *StatusCommandService) appendArchivedProjects(response *strings.Builder, projects []*projectDomain.Project) {
	if len(projects) == 0 {
		return
	}

	response.WriteString("📦 **Archived Projects:**\n")
	for _, project := range projects {
		response.WriteString(fmt.Sprintf("   • `%s`\n", project.Name()))
	}
	response.WriteString("\n")
}

// appendSummary adds project count summary to the response
func (s *StatusCommandService) appendSummary(response *strings.Builder, totalCount int) {
	response.WriteString(fmt.Sprintf("📊 **Total:** %d projects\n\n", totalCount))
}

// appendQuickCommands adds quick commands section to the response
func (s *StatusCommandService) appendQuickCommands(response *strings.Builder) {
	response.WriteString(quickCommandsMsg)
}

// getNotificationStatus returns the notification status for a project
func (s *StatusCommandService) getNotificationStatus(project *projectDomain.Project) string {
	if project.TelegramChatID() != nil {
		return "Yes"
	}
	return "No"
}
