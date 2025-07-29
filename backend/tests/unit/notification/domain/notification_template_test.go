package domain_test

import (
	"testing"

	"github.com/dewisartika8/cicd-status-notifier-bot/internal/core/notification/domain"
	"github.com/dewisartika8/cicd-status-notifier-bot/internal/core/shared/domain/value_objects"
	"github.com/stretchr/testify/assert"
)

func TestNotificationTemplate_NewNotificationTemplate(t *testing.T) {
	tests := []struct {
		name         string
		templateType domain.NotificationTemplateType
		channel      domain.NotificationChannel
		subject      string
		bodyTemplate string
		wantErr      bool
		expectedErr  string
	}{
		{
			name:         "valid_telegram_build_success_template",
			templateType: domain.TemplateTypeBuildSuccess,
			channel:      domain.NotificationChannelTelegram,
			subject:      "",
			bodyTemplate: "🎉 Build Success: {{.ProjectName}}",
			wantErr:      false,
		},
		{
			name:         "valid_email_build_failure_template",
			templateType: domain.TemplateTypeBuildFailure,
			channel:      domain.NotificationChannelEmail,
			subject:      "[BUILD FAILED] {{.ProjectName}}",
			bodyTemplate: "Build failed for project {{.ProjectName}}",
			wantErr:      false,
		},
		{
			name:         "invalid_template_type",
			templateType: domain.NotificationTemplateType("invalid"),
			channel:      domain.NotificationChannelTelegram,
			subject:      "",
			bodyTemplate: "Test template",
			wantErr:      true,
			expectedErr:  "invalid template type",
		},
		{
			name:         "invalid_channel",
			templateType: domain.TemplateTypeBuildSuccess,
			channel:      domain.NotificationChannel("invalid"),
			subject:      "",
			bodyTemplate: "Test template",
			wantErr:      true,
			expectedErr:  "notification channel is invalid",
		},
		{
			name:         "email_missing_subject",
			templateType: domain.TemplateTypeBuildSuccess,
			channel:      domain.NotificationChannelEmail,
			subject:      "",
			bodyTemplate: "Test template",
			wantErr:      true,
			expectedErr:  "subject is required for email",
		},
		{
			name:         "empty_body_template",
			templateType: domain.TemplateTypeBuildSuccess,
			channel:      domain.NotificationChannelTelegram,
			subject:      "",
			bodyTemplate: "",
			wantErr:      true,
			expectedErr:  "template body is invalid or empty",
		},
		{
			name:         "invalid_template_syntax",
			templateType: domain.TemplateTypeBuildSuccess,
			channel:      domain.NotificationChannelTelegram,
			subject:      "",
			bodyTemplate: "{{.InvalidSyntax",
			wantErr:      true,
			expectedErr:  "template compilation failed",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			template, err := domain.NewNotificationTemplate(
				tt.templateType,
				tt.channel,
				tt.subject,
				tt.bodyTemplate,
			)

			if tt.wantErr {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedErr)
				assert.Nil(t, template)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, template)
				assert.Equal(t, tt.templateType, template.TemplateType())
				assert.Equal(t, tt.channel, template.Channel())
				assert.Equal(t, tt.subject, template.Subject())
				assert.Equal(t, tt.bodyTemplate, template.BodyTemplate())
				assert.True(t, template.IsActive())
				assert.False(t, template.ID().IsNil())
			}
		})
	}
}

func TestNotificationTemplate_RenderTemplate(t *testing.T) {
	tests := []struct {
		name            string
		templateType    domain.NotificationTemplateType
		channel         domain.NotificationChannel
		subject         string
		bodyTemplate    string
		params          domain.TemplateParams
		expectedSubject string
		expectedBody    string
		wantErr         bool
	}{
		{
			name:         "render_telegram_build_success",
			templateType: domain.TemplateTypeBuildSuccess,
			channel:      domain.NotificationChannelTelegram,
			subject:      "",
			bodyTemplate: "🎉 Build Success: {{.ProjectName}} on {{.BuildBranch}}",
			params: domain.TemplateParams{
				ProjectName: "test-project",
				BuildBranch: "main",
			},
			expectedSubject: "",
			expectedBody:    "🎉 Build Success: test-project on main",
			wantErr:         false,
		},
		{
			name:         "render_email_with_dynamic_subject",
			templateType: domain.TemplateTypeBuildFailure,
			channel:      domain.NotificationChannelEmail,
			subject:      "[BUILD FAILED] {{.ProjectName}} - {{.BuildBranch}}",
			bodyTemplate: "Build failed for {{.ProjectName}}. Error: {{.ErrorMessage}}",
			params: domain.TemplateParams{
				ProjectName:  "test-project",
				BuildBranch:  "develop",
				ErrorMessage: "compilation error",
			},
			expectedSubject: "[BUILD FAILED] test-project - develop",
			expectedBody:    "Build failed for test-project. Error: compilation error",
			wantErr:         false,
		},
		{
			name:         "render_with_all_params",
			templateType: domain.TemplateTypeBuildSuccess,
			channel:      domain.NotificationChannelSlack,
			subject:      "",
			bodyTemplate: "Project: {{.ProjectName}}, Status: {{.BuildStatus}}, Duration: {{.BuildDuration}}, URL: {{.BuildURL}}",
			params: domain.TemplateParams{
				ProjectName:   "full-test",
				BuildStatus:   "success",
				BuildDuration: "2m 30s",
				BuildURL:      "https://ci.example.com/build/123",
			},
			expectedSubject: "",
			expectedBody:    "Project: full-test, Status: success, Duration: 2m 30s, URL: https://ci.example.com/build/123",
			wantErr:         false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			template, err := domain.NewNotificationTemplate(
				tt.templateType,
				tt.channel,
				tt.subject,
				tt.bodyTemplate,
			)
			assert.NoError(t, err)

			subject, body, err := template.RenderTemplate(tt.params)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedSubject, subject)
				assert.Equal(t, tt.expectedBody, body)
			}
		})
	}
}

func TestNotificationTemplate_UpdateTemplate(t *testing.T) {
	template, err := domain.NewNotificationTemplate(
		domain.TemplateTypeBuildSuccess,
		domain.NotificationChannelTelegram,
		"",
		"Original template: {{.ProjectName}}",
	)
	assert.NoError(t, err)

	tests := []struct {
		name        string
		newSubject  string
		newBody     string
		wantErr     bool
		expectedErr string
	}{
		{
			name:       "valid_update",
			newSubject: "",
			newBody:    "Updated template: {{.ProjectName}} - {{.BuildStatus}}",
			wantErr:    false,
		},
		{
			name:        "empty_body",
			newSubject:  "",
			newBody:     "",
			wantErr:     true,
			expectedErr: "template body is invalid or empty",
		},
		{
			name:        "invalid_syntax",
			newSubject:  "",
			newBody:     "{{.InvalidSyntax",
			wantErr:     true,
			expectedErr: "template compilation failed",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := template.UpdateTemplate(tt.newSubject, tt.newBody)

			if tt.wantErr {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedErr)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.newBody, template.BodyTemplate())
			}
		})
	}
}

func TestNotificationTemplate_ActivateDeactivate(t *testing.T) {
	template, err := domain.NewNotificationTemplate(
		domain.TemplateTypeBuildSuccess,
		domain.NotificationChannelTelegram,
		"",
		"Test template: {{.ProjectName}}",
	)
	assert.NoError(t, err)
	assert.True(t, template.IsActive())

	// Test deactivate
	err = template.Deactivate()
	assert.NoError(t, err)
	assert.False(t, template.IsActive())

	// Test deactivate already inactive
	err = template.Deactivate()
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "template is already inactive")

	// Test activate
	err = template.Activate()
	assert.NoError(t, err)
	assert.True(t, template.IsActive())

	// Test activate already active
	err = template.Activate()
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "template is already active")
}

func TestNotificationTemplate_RenderInactiveTemplate(t *testing.T) {
	template, err := domain.NewNotificationTemplate(
		domain.TemplateTypeBuildSuccess,
		domain.NotificationChannelTelegram,
		"",
		"Test template: {{.ProjectName}}",
	)
	assert.NoError(t, err)

	// Deactivate template
	err = template.Deactivate()
	assert.NoError(t, err)

	// Try to render inactive template
	params := domain.TemplateParams{ProjectName: "test"}
	_, _, err = template.RenderTemplate(params)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "template is inactive and cannot be rendered")
}

func TestNotificationTemplateType_IsValid(t *testing.T) {
	tests := []struct {
		name         string
		templateType domain.NotificationTemplateType
		expected     bool
	}{
		{"valid_build_success", domain.TemplateTypeBuildSuccess, true},
		{"valid_build_failure", domain.TemplateTypeBuildFailure, true},
		{"valid_build_started", domain.TemplateTypeBuildStarted, true},
		{"valid_deployment", domain.TemplateTypeDeployment, true},
		{"invalid_type", domain.NotificationTemplateType("invalid"), false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.templateType.IsValid()
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestGetDefaultTemplates(t *testing.T) {
	defaultTemplates := domain.GetDefaultTemplates()

	// Test that all required template types exist
	expectedTypes := []domain.NotificationTemplateType{
		domain.TemplateTypeBuildSuccess,
		domain.TemplateTypeBuildFailure,
		domain.TemplateTypeBuildStarted,
		domain.TemplateTypeDeployment,
	}

	for _, templateType := range expectedTypes {
		assert.Contains(t, defaultTemplates, templateType)

		// Test that telegram templates exist for all types
		assert.Contains(t, defaultTemplates[templateType], domain.NotificationChannelTelegram)
		telegramTemplate := defaultTemplates[templateType][domain.NotificationChannelTelegram]
		assert.NotEmpty(t, telegramTemplate.Body)

		// Test that email templates exist for success/failure
		if templateType == domain.TemplateTypeBuildSuccess || templateType == domain.TemplateTypeBuildFailure {
			assert.Contains(t, defaultTemplates[templateType], domain.NotificationChannelEmail)
			emailTemplate := defaultTemplates[templateType][domain.NotificationChannelEmail]
			assert.NotEmpty(t, emailTemplate.Subject)
			assert.NotEmpty(t, emailTemplate.Body)
		}
	}
}

func TestRestoreNotificationTemplate(t *testing.T) {
	id := value_objects.NewID()
	now := value_objects.NewTimestamp()

	params := domain.RestoreNotificationTemplateParams{
		ID:           id,
		TemplateType: domain.TemplateTypeBuildSuccess,
		Channel:      domain.NotificationChannelTelegram,
		Subject:      "",
		BodyTemplate: "Restored template: {{.ProjectName}}",
		IsActive:     true,
		CreatedAt:    now,
		UpdatedAt:    now,
	}

	template, err := domain.RestoreNotificationTemplate(params)
	assert.NoError(t, err)
	assert.NotNil(t, template)
	assert.Equal(t, id, template.ID())
	assert.Equal(t, domain.TemplateTypeBuildSuccess, template.TemplateType())
	assert.Equal(t, domain.NotificationChannelTelegram, template.Channel())
	assert.True(t, template.IsActive())
}

func TestRestoreNotificationTemplate_InvalidTemplate(t *testing.T) {
	id := value_objects.NewID()
	now := value_objects.NewTimestamp()

	params := domain.RestoreNotificationTemplateParams{
		ID:           id,
		TemplateType: domain.TemplateTypeBuildSuccess,
		Channel:      domain.NotificationChannelTelegram,
		Subject:      "",
		BodyTemplate: "{{.InvalidSyntax", // Invalid template syntax
		IsActive:     true,
		CreatedAt:    now,
		UpdatedAt:    now,
	}

	template, err := domain.RestoreNotificationTemplate(params)
	assert.Error(t, err)
	assert.Nil(t, template)
	assert.Contains(t, err.Error(), "template compilation failed")
}
