package formatter

import (
	"bytes"
	"context"
	"fmt"
	"strings"
	"text/template"

	"github.com/dewisartika8/cicd-status-notifier-bot/internal/core/notification/domain"
	"github.com/dewisartika8/cicd-status-notifier-bot/internal/core/notification/port"
	"github.com/sirupsen/logrus"
)

// Resource type constants
const (
	resourceTemplate     = "notification template"
	resourceNotification = "notification"
)

type Dep struct {
	TemplateRepo port.NotificationTemplateRepository
	Logger       *logrus.Logger
}

// notificationFormatterService implements the NotificationFormatterService interface
type notificationFormatterService struct {
	Dep
}

// NewNotificationFormatterService creates a new notification formatter service
func NewNotificationFormatterService(d Dep) port.NotificationFormatterService {
	return &notificationFormatterService{
		Dep: d,
	}
}

// FormatNotification formats a notification using templates
func (s *notificationFormatterService) FormatNotification(
	ctx context.Context,
	templateType domain.NotificationTemplateType,
	channel domain.NotificationChannel,
	params domain.TemplateParams,
) (subject, body string, err error) {
	s.Logger.WithFields(logrus.Fields{
		"template_type": templateType,
		"channel":       channel,
		"project_name":  params.ProjectName,
	}).Info(domain.LogMsgFormatNotification)

	// Get the template from repository
	tmpl, err := s.TemplateRepo.GetByTypeAndChannel(ctx, templateType, channel)
	if err != nil {
		s.Logger.WithError(err).Error("Failed to get template from repository")
		return "", "", fmt.Errorf(domain.ErrMsgGet, resourceTemplate, err)
	}

	// Format using the template
	return s.FormatNotificationWithTemplate(ctx, tmpl, params)
}

// FormatNotificationWithTemplate formats a notification using a specific template
func (s *notificationFormatterService) FormatNotificationWithTemplate(
	ctx context.Context,
	tmpl *domain.NotificationTemplate,
	params domain.TemplateParams,
) (subject, body string, err error) {
	s.Logger.WithFields(logrus.Fields{
		"template_id":   tmpl.ID().String(),
		"template_type": tmpl.TemplateType(),
		"channel":       tmpl.Channel(),
		"project_name":  params.ProjectName,
	}).Info(domain.LogMsgExecuteTemplate)

	// Render subject
	subject, err = s.renderTemplate("subject", tmpl.Subject(), params)
	if err != nil {
		s.Logger.WithError(err).Error("Failed to render subject template")
		return "", "", fmt.Errorf(domain.ErrMsgProcess, resourceNotification, err)
	}

	// Render body
	body, err = s.renderTemplate("body", tmpl.BodyTemplate(), params)
	if err != nil {
		s.Logger.WithError(err).Error("Failed to render body template")
		return "", "", fmt.Errorf(domain.ErrMsgProcess, resourceNotification, err)
	}

	// Add emoji formatting based on channel
	body = s.addEmojiFormatting(body, tmpl.Channel(), params)

	s.Logger.WithFields(logrus.Fields{
		"template_id": tmpl.ID().String(),
		"subject":     subject,
	}).Info("Notification formatted successfully")

	return subject, body, nil
}

// ValidateTemplate validates a template by trying to compile and render it
func (s *notificationFormatterService) ValidateTemplate(
	templateType domain.NotificationTemplateType,
	channel domain.NotificationChannel,
	subject, bodyTemplate string,
	testParams domain.TemplateParams,
) error {
	s.Logger.WithFields(logrus.Fields{
		"template_type": templateType,
		"channel":       channel,
	}).Info(domain.LogMsgValidateTemplate)

	// Check if email requires subject
	if channel == domain.NotificationChannelEmail && strings.TrimSpace(subject) == "" {
		s.Logger.Error("Email channel requires subject")
		return fmt.Errorf("subject is required for email notifications")
	}

	// Validate subject template
	if err := s.validateTemplateString("subject", subject, testParams); err != nil {
		s.Logger.WithError(err).Error("Subject template validation failed")
		return fmt.Errorf("template compilation failed: %w", err)
	}

	// Validate body template
	if err := s.validateTemplateString("body", bodyTemplate, testParams); err != nil {
		s.Logger.WithError(err).Error("Body template validation failed")
		return fmt.Errorf("template compilation failed: %w", err)
	}

	s.Logger.Info("Template validation successful")
	return nil
}

// GetAvailableTemplateVariables returns available template variables for a template type
func (s *notificationFormatterService) GetAvailableTemplateVariables(templateType domain.NotificationTemplateType) []string {
	s.Logger.WithField("template_type", templateType).Info("Getting available template variables")

	baseVars := []string{
		"ProjectName",
		"Timestamp",
		"BuildURL",
	}

	switch templateType {
	case domain.TemplateTypeBuildSuccess:
		return append(baseVars, []string{
			"BuildStatus",
			"BuildBranch",
			"BuildCommit",
			"BuildDuration",
		}...)

	case domain.TemplateTypeBuildFailure:
		return append(baseVars, []string{
			"BuildStatus",
			"BuildBranch",
			"BuildCommit",
			"BuildDuration",
			"ErrorMessage",
		}...)

	case domain.TemplateTypeDeployment:
		return append(baseVars, []string{
			"BuildStatus",
			"BuildBranch",
			"BuildCommit",
			"BuildDuration",
			"Environment",
		}...)

	case domain.TemplateTypeBuildStarted:
		return append(baseVars, []string{
			"BuildStatus",
			"BuildBranch",
			"BuildCommit",
		}...)

	default:
		return baseVars
	}
}

// FormatEmoji adds emoji formatting based on build status and channel
func (s *notificationFormatterService) FormatEmoji(status string, channel domain.NotificationChannel) string {
	switch strings.ToLower(status) {
	case "success", "passed", "completed":
		switch channel {
		case domain.NotificationChannelTelegram:
			return "✅"
		case domain.NotificationChannelSlack:
			return ":white_check_mark:"
		case domain.NotificationChannelEmail:
			return "[SUCCESS]"
		default:
			return "✅"
		}
	case "failure", "failed", "error":
		switch channel {
		case domain.NotificationChannelTelegram:
			return "❌"
		case domain.NotificationChannelSlack:
			return ":x:"
		case domain.NotificationChannelEmail:
			return "[FAILED]"
		default:
			return "❌"
		}
	case "warning", "unstable":
		switch channel {
		case domain.NotificationChannelTelegram:
			return "⚠️"
		case domain.NotificationChannelSlack:
			return ":warning:"
		case domain.NotificationChannelEmail:
			return "[WARNING]"
		default:
			return "⚠️"
		}
	case "running", "in_progress":
		switch channel {
		case domain.NotificationChannelTelegram:
			return "⏳"
		case domain.NotificationChannelSlack:
			return ":hourglass_flowing_sand:"
		case domain.NotificationChannelEmail:
			return "[RUNNING]"
		default:
			return "⏳"
		}
	case "pending", "queued":
		switch channel {
		case domain.NotificationChannelTelegram:
			return "⏳"
		case domain.NotificationChannelSlack:
			return ":hourglass_flowing_sand:"
		case domain.NotificationChannelEmail:
			return "[PENDING]"
		default:
			return "⏳"
		}
	default:
		return ""
	}
}

// renderTemplate renders a template string with the given parameters
func (s *notificationFormatterService) renderTemplate(name, templateStr string, params domain.TemplateParams) (string, error) {
	// Create a new template
	tmpl, err := template.New(name).Parse(templateStr)
	if err != nil {
		return "", fmt.Errorf(domain.ErrMsgParseTemplate, err)
	}

	// Execute the template
	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, params); err != nil {
		return "", fmt.Errorf(domain.ErrMsgExecuteTemplate, err)
	}

	return buf.String(), nil
}

// validateTemplateString validates a template string by trying to parse and execute it
func (s *notificationFormatterService) validateTemplateString(name, templateStr string, testParams domain.TemplateParams) error {
	// Try to parse the template
	tmpl, err := template.New(name).Parse(templateStr)
	if err != nil {
		return fmt.Errorf("parse error: %w", err)
	}

	// Try to execute with test parameters
	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, testParams); err != nil {
		return fmt.Errorf("execution error: %w", err)
	}

	return nil
}

// addEmojiFormatting adds emoji formatting to the body based on channel and parameters
func (s *notificationFormatterService) addEmojiFormatting(body string, channel domain.NotificationChannel, params domain.TemplateParams) string {
	// This is a simple implementation. In a real scenario, you might want to add more sophisticated formatting
	// based on the build status, channel preferences, etc.

	// Only add emoji if the template doesn't already contain any emoji characters
	// Check for common emoji patterns to avoid double formatting
	commonEmojis := []string{"🎉", "✅", "❌", "⚠️", "🔄", "⏳", "🚀", "💥"}
	hasEmoji := false
	for _, emoji := range commonEmojis {
		if strings.Contains(body, emoji) {
			hasEmoji = true
			break
		}
	}

	// Add status emoji if not already present and no existing emoji detected
	if !hasEmoji && params.BuildStatus != "" {
		emoji := s.FormatEmoji(params.BuildStatus, channel)
		if emoji != "" && !strings.Contains(body, emoji) {
			body = emoji + " " + body
		}
	}

	return body
}
