package template

import (
	"context"
	"fmt"

	"github.com/dewisartika8/cicd-status-notifier-bot/internal/core/notification/domain"
	"github.com/dewisartika8/cicd-status-notifier-bot/internal/core/notification/port"
	"github.com/dewisartika8/cicd-status-notifier-bot/internal/core/shared/domain/value_objects"
	"github.com/sirupsen/logrus"
)

// Resource type constants
const (
	resourceTemplate = "notification template"
)

type Dep struct {
	TemplateRepo port.NotificationTemplateRepository
	Logger       *logrus.Logger
}

// notificationTemplateService implements the NotificationTemplateService interface
type notificationTemplateService struct {
	Dep
}

// NewNotificationTemplateService creates a new notification template service
func NewNotificationTemplateService(d Dep) port.NotificationTemplateService {
	return &notificationTemplateService{
		Dep: d,
	}
}

// CreateNotificationTemplate creates a new notification template
func (s *notificationTemplateService) CreateNotificationTemplate(
	ctx context.Context,
	templateType domain.NotificationTemplateType,
	channel domain.NotificationChannel,
	subject, bodyTemplate string,
) (*domain.NotificationTemplate, error) {
	s.Logger.WithFields(logrus.Fields{
		"template_type": templateType,
		"channel":       channel,
		"subject":       subject,
	}).Info("Creating notification template")

	// Create new notification template entity
	template, err := domain.NewNotificationTemplate(templateType, channel, subject, bodyTemplate)
	if err != nil {
		s.Logger.WithError(err).Error("Failed to create notification template entity")
		return nil, fmt.Errorf(domain.ErrMsgCreate, resourceTemplate, err)
	}

	// Persist the template
	if err := s.TemplateRepo.Create(ctx, template); err != nil {
		s.Logger.WithError(err).Error("Failed to persist notification template")
		return nil, fmt.Errorf(domain.ErrMsgCreate, resourceTemplate, err)
	}

	s.Logger.WithField("template_id", template.ID().String()).Info(domain.TemplateCreated)
	return template, nil
}

// GetNotificationTemplate retrieves a notification template by its ID
func (s *notificationTemplateService) GetNotificationTemplate(ctx context.Context, id value_objects.ID) (*domain.NotificationTemplate, error) {
	s.Logger.WithField("id", id.String()).Info("Getting notification template")

	template, err := s.TemplateRepo.GetByID(ctx, id)
	if err != nil {
		s.Logger.WithError(err).Error(domain.LogMsgGetTemplate)
		return nil, fmt.Errorf(domain.ErrMsgGet, resourceTemplate, err)
	}

	return template, nil
}

// GetTemplateByTypeAndChannel retrieves a template by type and channel
func (s *notificationTemplateService) GetTemplateByTypeAndChannel(
	ctx context.Context,
	templateType domain.NotificationTemplateType,
	channel domain.NotificationChannel,
) (*domain.NotificationTemplate, error) {
	s.Logger.WithFields(logrus.Fields{
		"template_type": templateType,
		"channel":       channel,
	}).Info("Getting template by type and channel")

	template, err := s.TemplateRepo.GetByTypeAndChannel(ctx, templateType, channel)
	if err != nil {
		s.Logger.WithError(err).Error(domain.LogMsgGetTemplate)
		return nil, fmt.Errorf(domain.ErrMsgGet, resourceTemplate, err)
	}

	return template, nil
}

// UpdateNotificationTemplate updates an existing notification template
func (s *notificationTemplateService) UpdateNotificationTemplate(
	ctx context.Context,
	id value_objects.ID,
	subject, bodyTemplate string,
) (*domain.NotificationTemplate, error) {
	s.Logger.WithField("id", id.String()).Info("Updating notification template")

	// Get the template
	template, err := s.TemplateRepo.GetByID(ctx, id)
	if err != nil {
		s.Logger.WithError(err).Error(domain.LogMsgGetTemplate)
		return nil, fmt.Errorf(domain.ErrMsgGet, resourceTemplate, err)
	}

	// Update template content
	if err := template.UpdateTemplate(subject, bodyTemplate); err != nil {
		s.Logger.WithError(err).Error("Failed to update template content")
		return nil, fmt.Errorf(domain.ErrMsgUpdate, resourceTemplate, err)
	}

	// Update the template in repository
	if err := s.TemplateRepo.Update(ctx, template); err != nil {
		s.Logger.WithError(err).Error(domain.LogMsgUpdateTemplate)
		return nil, fmt.Errorf(domain.ErrMsgUpdate, resourceTemplate, err)
	}

	s.Logger.Info(domain.TemplateUpdated)
	return template, nil
}

// ActivateTemplate activates a notification template
func (s *notificationTemplateService) ActivateTemplate(ctx context.Context, id value_objects.ID) error {
	s.Logger.WithField("id", id.String()).Info("Activating notification template")

	template, err := s.TemplateRepo.GetByID(ctx, id)
	if err != nil {
		s.Logger.WithError(err).Error(domain.LogMsgGetTemplate)
		return fmt.Errorf(domain.ErrMsgGet, resourceTemplate, err)
	}

	template.Activate()

	if err := s.TemplateRepo.Update(ctx, template); err != nil {
		s.Logger.WithError(err).Error(domain.LogMsgActivateTemplate)
		return fmt.Errorf(domain.ErrMsgActivate, resourceTemplate, err)
	}

	s.Logger.Info(domain.TemplateActivated)
	return nil
}

// DeactivateTemplate deactivates a notification template
func (s *notificationTemplateService) DeactivateTemplate(ctx context.Context, id value_objects.ID) error {
	s.Logger.WithField("id", id.String()).Info("Deactivating notification template")

	template, err := s.TemplateRepo.GetByID(ctx, id)
	if err != nil {
		s.Logger.WithError(err).Error(domain.LogMsgGetTemplate)
		return fmt.Errorf(domain.ErrMsgGet, resourceTemplate, err)
	}

	template.Deactivate()

	if err := s.TemplateRepo.Update(ctx, template); err != nil {
		s.Logger.WithError(err).Error(domain.LogMsgDeactivateTemplate)
		return fmt.Errorf(domain.ErrMsgDeactivate, resourceTemplate, err)
	}

	s.Logger.Info(domain.TemplateDeactivated)
	return nil
}

// DeleteNotificationTemplate deletes a notification template
func (s *notificationTemplateService) DeleteNotificationTemplate(ctx context.Context, id value_objects.ID) error {
	s.Logger.WithField("id", id.String()).Info("Deleting notification template")

	if err := s.TemplateRepo.Delete(ctx, id); err != nil {
		s.Logger.WithError(err).Error(domain.LogMsgDeleteTemplate)
		return fmt.Errorf(domain.ErrMsgDelete, resourceTemplate, err)
	}

	s.Logger.Info(domain.TemplateDeleted)
	return nil
}

// GetActiveTemplates retrieves all active notification templates
func (s *notificationTemplateService) GetActiveTemplates(ctx context.Context) ([]*domain.NotificationTemplate, error) {
	s.Logger.Info("Getting all active notification templates")

	templates, err := s.TemplateRepo.GetActiveTemplates(ctx)
	if err != nil {
		s.Logger.WithError(err).Error("Failed to get active notification templates")
		return nil, fmt.Errorf("failed to get active %s: %w", resourceTemplate, err)
	}

	s.Logger.WithField("count", len(templates)).Info("Retrieved active notification templates")
	return templates, nil
}

// InitializeDefaultTemplates creates default templates for all channels and types
func (s *notificationTemplateService) InitializeDefaultTemplates(ctx context.Context) error {
	s.Logger.Info("Initializing default notification templates")

	// Define default templates
	defaultTemplates := []struct {
		templateType domain.NotificationTemplateType
		channel      domain.NotificationChannel
		subject      string
		bodyTemplate string
	}{
		// Success templates
		{
			templateType: domain.TemplateTypeBuildSuccess,
			channel:      domain.NotificationChannelTelegram,
			subject:      "✅ Build Success",
			bodyTemplate: "✅ **Build Successful**\n\n📋 **Project:** {{.ProjectName}}\n🔀 **Branch:** {{.Branch}}\n📦 **Build:** #{{.BuildNumber}}\n⏰ **Duration:** {{.Duration}}\n🔗 **View Details:** {{.BuildURL}}",
		},
		{
			templateType: domain.TemplateTypeBuildSuccess,
			channel:      domain.NotificationChannelEmail,
			subject:      "✅ Build Success - {{.ProjectName}}",
			bodyTemplate: "Good news! Your build has completed successfully.\n\nProject: {{.ProjectName}}\nBranch: {{.Branch}}\nBuild Number: {{.BuildNumber}}\nDuration: {{.Duration}}\nBuild URL: {{.BuildURL}}\n\nBest regards,\nCI/CD Notification Bot",
		},
		{
			templateType: domain.TemplateTypeBuildSuccess,
			channel:      domain.NotificationChannelSlack,
			subject:      "Build Success",
			bodyTemplate: ":white_check_mark: *Build Successful*\n\n*Project:* {{.ProjectName}}\n*Branch:* {{.Branch}}\n*Build:* #{{.BuildNumber}}\n*Duration:* {{.Duration}}\n<{{.BuildURL}}|View Details>",
		},
		{
			templateType: domain.TemplateTypeBuildSuccess,
			channel:      domain.NotificationChannelWebhook,
			subject:      "Build Success",
			bodyTemplate: `{"status": "success", "project": "{{.ProjectName}}", "branch": "{{.Branch}}", "build_number": "{{.BuildNumber}}", "duration": "{{.Duration}}", "build_url": "{{.BuildURL}}"}`,
		},
		// Failure templates
		{
			templateType: domain.TemplateTypeBuildFailure,
			channel:      domain.NotificationChannelTelegram,
			subject:      "❌ Build Failed",
			bodyTemplate: "❌ **Build Failed**\n\n📋 **Project:** {{.ProjectName}}\n🔀 **Branch:** {{.Branch}}\n📦 **Build:** #{{.BuildNumber}}\n⏰ **Duration:** {{.Duration}}\n❗ **Error:** {{.ErrorMessage}}\n🔗 **View Details:** {{.BuildURL}}",
		},
		{
			templateType: domain.TemplateTypeBuildFailure,
			channel:      domain.NotificationChannelEmail,
			subject:      "❌ Build Failed - {{.ProjectName}}",
			bodyTemplate: "Unfortunately, your build has failed.\n\nProject: {{.ProjectName}}\nBranch: {{.Branch}}\nBuild Number: {{.BuildNumber}}\nDuration: {{.Duration}}\nError: {{.ErrorMessage}}\nBuild URL: {{.BuildURL}}\n\nPlease check the build logs for more details.\n\nBest regards,\nCI/CD Notification Bot",
		},
		{
			templateType: domain.TemplateTypeBuildFailure,
			channel:      domain.NotificationChannelSlack,
			subject:      "Build Failed",
			bodyTemplate: ":x: *Build Failed*\n\n*Project:* {{.ProjectName}}\n*Branch:* {{.Branch}}\n*Build:* #{{.BuildNumber}}\n*Duration:* {{.Duration}}\n*Error:* {{.ErrorMessage}}\n<{{.BuildURL}}|View Details>",
		},
		{
			templateType: domain.TemplateTypeBuildFailure,
			channel:      domain.NotificationChannelWebhook,
			subject:      "Build Failed",
			bodyTemplate: `{"status": "failure", "project": "{{.ProjectName}}", "branch": "{{.Branch}}", "build_number": "{{.BuildNumber}}", "duration": "{{.Duration}}", "error": "{{.ErrorMessage}}", "build_url": "{{.BuildURL}}"}`,
		},
		// Started templates
		{
			templateType: domain.TemplateTypeBuildStarted,
			channel:      domain.NotificationChannelTelegram,
			subject:      "🔄 Build Started",
			bodyTemplate: "🔄 **Build Started**\n\n📋 **Project:** {{.ProjectName}}\n🔀 **Branch:** {{.Branch}}\n📦 **Build:** #{{.BuildNumber}}\n🔗 **View Details:** {{.BuildURL}}",
		},
		{
			templateType: domain.TemplateTypeBuildStarted,
			channel:      domain.NotificationChannelEmail,
			subject:      "🔄 Build Started - {{.ProjectName}}",
			bodyTemplate: "Your build has started.\n\nProject: {{.ProjectName}}\nBranch: {{.Branch}}\nBuild Number: {{.BuildNumber}}\nBuild URL: {{.BuildURL}}\n\nBest regards,\nCI/CD Notification Bot",
		},
		{
			templateType: domain.TemplateTypeBuildStarted,
			channel:      domain.NotificationChannelSlack,
			subject:      "Build Started",
			bodyTemplate: ":arrows_counterclockwise: *Build Started*\n\n*Project:* {{.ProjectName}}\n*Branch:* {{.Branch}}\n*Build:* #{{.BuildNumber}}\n<{{.BuildURL}}|View Details>",
		},
		{
			templateType: domain.TemplateTypeBuildStarted,
			channel:      domain.NotificationChannelWebhook,
			subject:      "Build Started",
			bodyTemplate: `{"status": "started", "project": "{{.ProjectName}}", "branch": "{{.Branch}}", "build_number": "{{.BuildNumber}}", "build_url": "{{.BuildURL}}"}`,
		},
	}

	// Create each default template
	for _, tmpl := range defaultTemplates {
		// Check if template already exists
		existing, err := s.TemplateRepo.GetByTypeAndChannel(ctx, tmpl.templateType, tmpl.channel)
		if err == nil && existing != nil {
			s.Logger.WithFields(logrus.Fields{
				"template_type": tmpl.templateType,
				"channel":       tmpl.channel,
			}).Info("Default template already exists, skipping")
			continue
		}

		// Create the template
		template, err := domain.NewNotificationTemplate(
			tmpl.templateType,
			tmpl.channel,
			tmpl.subject,
			tmpl.bodyTemplate,
		)
		if err != nil {
			s.Logger.WithError(err).WithFields(logrus.Fields{
				"template_type": tmpl.templateType,
				"channel":       tmpl.channel,
			}).Error("Failed to create default template entity")
			continue
		}

		// Persist the template
		if err := s.TemplateRepo.Create(ctx, template); err != nil {
			s.Logger.WithError(err).WithFields(logrus.Fields{
				"template_type": tmpl.templateType,
				"channel":       tmpl.channel,
			}).Error("Failed to persist default template")
			continue
		}

		s.Logger.WithFields(logrus.Fields{
			"template_id":   template.ID().String(),
			"template_type": tmpl.templateType,
			"channel":       tmpl.channel,
		}).Info("Default template created successfully")
	}

	s.Logger.Info("Default notification templates initialization completed")
	return nil
}
