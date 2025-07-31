package domain

import (
	"fmt"
	"strings"
	"text/template"

	"github.com/dewisartika8/cicd-status-notifier-bot/internal/core/shared/domain/value_objects"
)

// NotificationTemplateType represents the type of notification template
type NotificationTemplateType value_objects.Status

const (
	TemplateTypeBuildSuccess NotificationTemplateType = "build_success"
	TemplateTypeBuildFailure NotificationTemplateType = "build_failure"
	TemplateTypeBuildStarted NotificationTemplateType = "build_started"
	TemplateTypeDeployment   NotificationTemplateType = "deployment"
)

// IsValid checks if the template type is valid
func (t NotificationTemplateType) IsValid() bool {
	switch t {
	case TemplateTypeBuildSuccess, TemplateTypeBuildFailure, TemplateTypeBuildStarted, TemplateTypeDeployment:
		return true
	default:
		return false
	}
}

// String returns the string representation of the template type
func (t NotificationTemplateType) String() string {
	return string(t)
}

// NotificationTemplate represents a notification template domain entity
type NotificationTemplate struct {
	id               value_objects.ID
	templateType     NotificationTemplateType
	channel          NotificationChannel
	subject          string
	bodyTemplate     string
	compiledTemplate *template.Template
	isActive         bool
	createdAt        value_objects.Timestamp
	updatedAt        value_objects.Timestamp
}

// TemplateParams holds parameters for template substitution
type TemplateParams struct {
	ProjectName   string
	BuildStatus   string
	BuildBranch   string
	BuildCommit   string
	BuildDuration string
	BuildURL      string
	ErrorMessage  string
	Timestamp     string
	Environment   string
}

// NewNotificationTemplate creates a new notification template entity
func NewNotificationTemplate(
	templateType NotificationTemplateType,
	channel NotificationChannel,
	subject, bodyTemplate string,
) (*NotificationTemplate, error) {
	// Validate template type
	if !templateType.IsValid() {
		return nil, NewInvalidTemplateTypeError(string(templateType))
	}

	// Validate channel
	if !channel.IsValid() {
		return nil, ErrInvalidNotificationChannel
	}

	// Validate subject for channels that require it
	if needsSubject(channel) && strings.TrimSpace(subject) == "" {
		return nil, NewInvalidTemplateSubjectError("subject is required for " + string(channel))
	}

	// Validate body template
	if strings.TrimSpace(bodyTemplate) == "" {
		return nil, ErrInvalidTemplateBody
	}

	// Compile template to validate syntax
	compiledTemplate, err := template.New(string(templateType)).Parse(bodyTemplate)
	if err != nil {
		return nil, NewInvalidTemplateBodyError(fmt.Sprintf("template compilation failed: %v", err))
	}

	now := value_objects.NewTimestamp()
	id := value_objects.NewID()

	template := &NotificationTemplate{
		id:               id,
		templateType:     templateType,
		channel:          channel,
		subject:          subject,
		bodyTemplate:     bodyTemplate,
		compiledTemplate: compiledTemplate,
		isActive:         true,
		createdAt:        now,
		updatedAt:        now,
	}

	return template, nil
}

// RestoreNotificationTemplate restores a notification template from persistence
func RestoreNotificationTemplate(params RestoreNotificationTemplateParams) (*NotificationTemplate, error) {
	// Compile template
	compiledTemplate, err := template.New(params.TemplateType.String()).Parse(params.BodyTemplate)
	if err != nil {
		return nil, NewInvalidTemplateBodyError(fmt.Sprintf("template compilation failed: %v", err))
	}

	return &NotificationTemplate{
		id:               params.ID,
		templateType:     params.TemplateType,
		channel:          params.Channel,
		subject:          params.Subject,
		bodyTemplate:     params.BodyTemplate,
		compiledTemplate: compiledTemplate,
		isActive:         params.IsActive,
		createdAt:        params.CreatedAt,
		updatedAt:        params.UpdatedAt,
	}, nil
}

// RestoreNotificationTemplateParams holds parameters for restoring a notification template
type RestoreNotificationTemplateParams struct {
	ID           value_objects.ID
	TemplateType NotificationTemplateType
	Channel      NotificationChannel
	Subject      string
	BodyTemplate string
	IsActive     bool
	CreatedAt    value_objects.Timestamp
	UpdatedAt    value_objects.Timestamp
}

// Getters
func (nt *NotificationTemplate) ID() value_objects.ID {
	return nt.id
}

func (nt *NotificationTemplate) TemplateType() NotificationTemplateType {
	return nt.templateType
}

func (nt *NotificationTemplate) Channel() NotificationChannel {
	return nt.channel
}

func (nt *NotificationTemplate) Subject() string {
	return nt.subject
}

func (nt *NotificationTemplate) BodyTemplate() string {
	return nt.bodyTemplate
}

func (nt *NotificationTemplate) IsActive() bool {
	return nt.isActive
}

func (nt *NotificationTemplate) CreatedAt() value_objects.Timestamp {
	return nt.createdAt
}

func (nt *NotificationTemplate) UpdatedAt() value_objects.Timestamp {
	return nt.updatedAt
}

// Business methods

// UpdateTemplate updates the template content
func (nt *NotificationTemplate) UpdateTemplate(subject, bodyTemplate string) error {
	// Validate subject for channels that require it
	if needsSubject(nt.channel) && strings.TrimSpace(subject) == "" {
		return NewInvalidTemplateSubjectError("subject is required for " + string(nt.channel))
	}

	// Validate body template
	if strings.TrimSpace(bodyTemplate) == "" {
		return ErrInvalidTemplateBody
	}

	// Compile template to validate syntax
	compiledTemplate, err := template.New(nt.templateType.String()).Parse(bodyTemplate)
	if err != nil {
		return NewInvalidTemplateBodyError(fmt.Sprintf("template compilation failed: %v", err))
	}

	// Update fields
	nt.subject = subject
	nt.bodyTemplate = bodyTemplate
	nt.compiledTemplate = compiledTemplate
	nt.updatedAt = value_objects.NewTimestamp()

	return nil
}

// Activate activates the template
func (nt *NotificationTemplate) Activate() error {
	if nt.isActive {
		return ErrTemplateAlreadyActive
	}

	nt.isActive = true
	nt.updatedAt = value_objects.NewTimestamp()
	return nil
}

// Deactivate deactivates the template
func (nt *NotificationTemplate) Deactivate() error {
	if !nt.isActive {
		return ErrTemplateAlreadyInactive
	}

	nt.isActive = false
	nt.updatedAt = value_objects.NewTimestamp()
	return nil
}

// RenderTemplate renders the template with provided parameters
func (nt *NotificationTemplate) RenderTemplate(params TemplateParams) (subject, body string, err error) {
	if !nt.isActive {
		return "", "", ErrTemplateInactive
	}

	// Render body template
	var bodyBuffer strings.Builder
	if err := nt.compiledTemplate.Execute(&bodyBuffer, params); err != nil {
		return "", "", NewTemplateRenderError(fmt.Sprintf("failed to render body template: %v", err))
	}

	// Render subject template if it contains template variables
	renderedSubject := nt.subject
	if strings.Contains(nt.subject, "{{") {
		subjectTemplate, err := template.New("subject").Parse(nt.subject)
		if err != nil {
			return "", "", NewTemplateRenderError(fmt.Sprintf("failed to parse subject template: %v", err))
		}

		var subjectBuffer strings.Builder
		if err := subjectTemplate.Execute(&subjectBuffer, params); err != nil {
			return "", "", NewTemplateRenderError(fmt.Sprintf("failed to render subject template: %v", err))
		}
		renderedSubject = subjectBuffer.String()
	}

	return renderedSubject, bodyBuffer.String(), nil
}

// Helper functions

// needsSubject returns true if the channel requires a subject
func needsSubject(channel NotificationChannel) bool {
	switch channel {
	case NotificationChannelEmail:
		return true
	default:
		return false
	}
}

// GetDefaultTemplates returns default notification templates for all channels and types
func GetDefaultTemplates() map[NotificationTemplateType]map[NotificationChannel]DefaultTemplate {
	return map[NotificationTemplateType]map[NotificationChannel]DefaultTemplate{
		TemplateTypeBuildSuccess: {
			NotificationChannelTelegram: {
				Subject: "",
				Body: `🎉 *Build Success*

*Project:* {{.ProjectName}}
*Branch:* {{.BuildBranch}}
*Commit:* {{.BuildCommit}}
*Duration:* {{.BuildDuration}}
*Time:* {{.Timestamp}}

✅ Build completed successfully!

[View Build]({{.BuildURL}})`,
			},
			NotificationChannelEmail: {
				Subject: "[BUILD SUCCESS] {{.ProjectName}} - {{.BuildBranch}}",
				Body: `Build completed successfully!

Project: {{.ProjectName}}
Branch: {{.BuildBranch}}
Commit: {{.BuildCommit}}
Duration: {{.BuildDuration}}
Time: {{.Timestamp}}

View Build: {{.BuildURL}}`,
			},
			NotificationChannelSlack: {
				Subject: "",
				Body: `🎉 *Build Success*

*Project:* {{.ProjectName}}
*Branch:* {{.BuildBranch}}
*Commit:* {{.BuildCommit}}
*Duration:* {{.BuildDuration}}
*Time:* {{.Timestamp}}

✅ Build completed successfully!

<{{.BuildURL}}|View Build>`,
			},
		},
		TemplateTypeBuildFailure: {
			NotificationChannelTelegram: {
				Subject: "",
				Body: `🚨 *Build Failed*

*Project:* {{.ProjectName}}
*Branch:* {{.BuildBranch}}
*Commit:* {{.BuildCommit}}
*Duration:* {{.BuildDuration}}
*Time:* {{.Timestamp}}

❌ Build failed!

*Error:* {{.ErrorMessage}}

[View Build]({{.BuildURL}})`,
			},
			NotificationChannelEmail: {
				Subject: "[BUILD FAILED] {{.ProjectName}} - {{.BuildBranch}}",
				Body: `Build failed!

Project: {{.ProjectName}}
Branch: {{.BuildBranch}}
Commit: {{.BuildCommit}}
Duration: {{.BuildDuration}}
Time: {{.Timestamp}}

Error: {{.ErrorMessage}}

View Build: {{.BuildURL}}`,
			},
			NotificationChannelSlack: {
				Subject: "",
				Body: `🚨 *Build Failed*

*Project:* {{.ProjectName}}
*Branch:* {{.BuildBranch}}
*Commit:* {{.BuildCommit}}
*Duration:* {{.BuildDuration}}
*Time:* {{.Timestamp}}

❌ Build failed!

*Error:* {{.ErrorMessage}}

<{{.BuildURL}}|View Build>`,
			},
		},
		TemplateTypeBuildStarted: {
			NotificationChannelTelegram: {
				Subject: "",
				Body: `🔄 *Build Started*

*Project:* {{.ProjectName}}
*Branch:* {{.BuildBranch}}
*Commit:* {{.BuildCommit}}
*Time:* {{.Timestamp}}

⏳ Build is now running...

[View Build]({{.BuildURL}})`,
			},
			NotificationChannelSlack: {
				Subject: "",
				Body: `🔄 *Build Started*

*Project:* {{.ProjectName}}
*Branch:* {{.BuildBranch}}
*Commit:* {{.BuildCommit}}
*Time:* {{.Timestamp}}

⏳ Build is now running...

<{{.BuildURL}}|View Build>`,
			},
		},
		TemplateTypeDeployment: {
			NotificationChannelTelegram: {
				Subject: "",
				Body: `🚀 *Deployment*

*Project:* {{.ProjectName}}
*Environment:* {{.Environment}}
*Branch:* {{.BuildBranch}}
*Commit:* {{.BuildCommit}}
*Time:* {{.Timestamp}}

🎯 Successfully deployed to {{.Environment}}!

[View Build]({{.BuildURL}})`,
			},
			NotificationChannelEmail: {
				Subject: "[DEPLOYMENT] {{.ProjectName}} deployed to {{.Environment}}",
				Body: `Deployment completed successfully!

Project: {{.ProjectName}}
Environment: {{.Environment}}
Branch: {{.BuildBranch}}
Commit: {{.BuildCommit}}
Time: {{.Timestamp}}

View Build: {{.BuildURL}}`,
			},
			NotificationChannelSlack: {
				Subject: "",
				Body: `🚀 *Deployment*

*Project:* {{.ProjectName}}
*Environment:* {{.Environment}}
*Branch:* {{.BuildBranch}}
*Commit:* {{.BuildCommit}}
*Time:* {{.Timestamp}}

🎯 Successfully deployed to {{.Environment}}!

<{{.BuildURL}}|View Build>`,
			},
		},
	}
}

// DefaultTemplate holds default template content
type DefaultTemplate struct {
	Subject string
	Body    string
}
