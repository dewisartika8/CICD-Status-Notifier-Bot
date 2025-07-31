package telegram

import (
	"strings"

	"github.com/dewisartika8/cicd-status-notifier-bot/internal/core/project/port"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// CommandRouter handles routing of different telegram commands
type CommandRouter struct {
	statusHandler *StatusCommandHandler
	formatter     *ResponseFormatter
}

// NewCommandRouter creates a new command router with all handlers
func NewCommandRouter(projectService port.ProjectService) *CommandRouter {
	return &CommandRouter{
		statusHandler: NewStatusCommandHandler(projectService),
		formatter:     NewResponseFormatter(),
	}
}

// RouteCommand routes incoming telegram commands to appropriate handlers
func (cr *CommandRouter) RouteCommand(msg *tgbotapi.Message) (string, error) {
	if !msg.IsCommand() {
		return cr.formatter.FormatInfo("Invalid input",
			"Please send a valid command. Type /help for available commands."), nil
	}

	command := strings.ToLower(msg.Command())
	args := strings.Fields(msg.CommandArguments())

	switch command {
	case "start":
		return cr.handleStartCommand(msg)
	case "help":
		return cr.handleHelpCommand(msg)
	case "status":
		return cr.handleStatusCommand(msg, args)
	case "projects":
		return cr.statusHandler.HandleProjectsList()
	default:
		return cr.handleUnknownCommand(msg, command)
	}
}

// handleStartCommand handles /start command
func (cr *CommandRouter) handleStartCommand(msg *tgbotapi.Message) (string, error) {
	userName := "there"
	if msg.From.FirstName != "" {
		userName = msg.From.FirstName
	}

	response := "🎉 **Welcome to CICD Status Notifier Bot!**\n\n" +
		"Hello " + userName + "! 👋\n\n" +
		"I'm here to help you monitor your CI/CD pipeline status and get real-time notifications about your builds, deployments, and more.\n\n" +
		"**Quick Start:**\n" +
		"• Type /help to see all available commands\n" +
		"• Use /projects to see monitored projects\n" +
		"• Check /status to see current pipeline status\n\n" +
		"Let's get started! 🚀"

	return response, nil
}

// handleHelpCommand handles /help command
func (cr *CommandRouter) handleHelpCommand(msg *tgbotapi.Message) (string, error) {
	response := "📚 **CICD Status Notifier Bot - Help**\n\n" +
		"**Available Commands:**\n\n" +
		"🏁 */start* - Welcome message and quick introduction\n" +
		"📖 */help* - Show this help message\n\n" +
		"📊 *Project Status Commands:*\n" +
		"• */status* - Get current status of all projects\n" +
		"• */status <project>* - Get detailed status for specific project\n" +
		"• */projects* - List all monitored projects\n\n" +
		"**Usage Examples:**\n" +
		"• `/status` - Get status for all projects\n" +
		"• `/status my-app` - Get status for \"my-app\" project\n" +
		"• `/projects` - List all projects\n\n" +
		"**Need more help?** Contact your system administrator."

	return response, nil
}

// handleStatusCommand handles /status command with optional project name
func (cr *CommandRouter) handleStatusCommand(msg *tgbotapi.Message, args []string) (string, error) {
	if len(args) == 0 {
		// No project specified, show all projects status
		return cr.statusHandler.HandleStatusAllProjects()
	} else if len(args) == 1 {
		// Specific project requested
		projectName := strings.TrimSpace(args[0])
		return cr.statusHandler.HandleStatusSpecificProject(projectName)
	} else {
		// Too many arguments
		return cr.formatter.FormatInvalidCommand("`/status` or `/status <project_name>`"), nil
	}
}

// handleUnknownCommand handles unknown commands
func (cr *CommandRouter) handleUnknownCommand(msg *tgbotapi.Message, command string) (string, error) {
	return cr.formatter.FormatError("Unknown command",
		"Command `/"+command+"` is not recognized.\n\n"+
			"Type /help for available commands."), nil
}
