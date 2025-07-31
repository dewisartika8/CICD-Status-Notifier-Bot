package service_test

import (
	"context"
	"testing"

	"github.com/dewisartika8/cicd-status-notifier-bot/internal/core/notification/domain"
	"github.com/dewisartika8/cicd-status-notifier-bot/internal/core/notification/service"
	"github.com/dewisartika8/cicd-status-notifier-bot/internal/core/shared/domain/value_objects"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

const (
	testProjectName = "test-project"
)

// MockNotificationTemplateRepository is a mock implementation of NotificationTemplateRepository
type MockNotificationTemplateRepository struct {
	mock.Mock
}

func (m *MockNotificationTemplateRepository) Create(ctx context.Context, template *domain.NotificationTemplate) error {
	args := m.Called(ctx, template)
	return args.Error(0)
}

func (m *MockNotificationTemplateRepository) GetByID(ctx context.Context, id value_objects.ID) (*domain.NotificationTemplate, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(*domain.NotificationTemplate), args.Error(1)
}

func (m *MockNotificationTemplateRepository) GetByTypeAndChannel(ctx context.Context, templateType domain.NotificationTemplateType, channel domain.NotificationChannel) (*domain.NotificationTemplate, error) {
	args := m.Called(ctx, templateType, channel)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.NotificationTemplate), args.Error(1)
}

func (m *MockNotificationTemplateRepository) GetByType(ctx context.Context, templateType domain.NotificationTemplateType) ([]*domain.NotificationTemplate, error) {
	args := m.Called(ctx, templateType)
	return args.Get(0).([]*domain.NotificationTemplate), args.Error(1)
}

func (m *MockNotificationTemplateRepository) GetByChannel(ctx context.Context, channel domain.NotificationChannel) ([]*domain.NotificationTemplate, error) {
	args := m.Called(ctx, channel)
	return args.Get(0).([]*domain.NotificationTemplate), args.Error(1)
}

func (m *MockNotificationTemplateRepository) GetActiveTemplates(ctx context.Context) ([]*domain.NotificationTemplate, error) {
	args := m.Called(ctx)
	return args.Get(0).([]*domain.NotificationTemplate), args.Error(1)
}

func (m *MockNotificationTemplateRepository) Update(ctx context.Context, template *domain.NotificationTemplate) error {
	// This is the Update method, which is different from Create.
	args := m.Called(ctx, template)
	// You could add additional logic here if needed for testing Update specifically.
	return args.Error(0)
}

func (m *MockNotificationTemplateRepository) Delete(ctx context.Context, id value_objects.ID) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockNotificationTemplateRepository) Count(ctx context.Context, templateType *domain.NotificationTemplateType, channel *domain.NotificationChannel, isActive *bool) (int64, error) {
	args := m.Called(ctx, templateType, channel, isActive)
	return args.Get(0).(int64), args.Error(1)
}

func TestNotificationFormatterServiceFormatNotification(t *testing.T) {
	ctx := context.Background()
	mockRepo := new(MockNotificationTemplateRepository)
	logger := logrus.New()

	formatterService := service.NewNotificationFormatterService(service.NotificationFormatterDep{
		TemplateRepo: mockRepo,
		Logger:       logger,
	})

	tests := []struct {
		name            string
		templateType    domain.NotificationTemplateType
		channel         domain.NotificationChannel
		params          domain.TemplateParams
		mockTemplate    *domain.NotificationTemplate
		mockError       error
		expectedSubject string
		expectedBody    string
		wantErr         bool
		expectedErrMsg  string
	}{
		{
			name:         "successful_telegram_build_success_formatting",
			templateType: domain.TemplateTypeBuildSuccess,
			channel:      domain.NotificationChannelTelegram,
			params: domain.TemplateParams{
				ProjectName: testProjectName,
				BuildBranch: "main",
				BuildStatus: "success",
			},
			mockTemplate: func() *domain.NotificationTemplate {
				template, _ := domain.NewNotificationTemplate(
					domain.TemplateTypeBuildSuccess,
					domain.NotificationChannelTelegram,
					"",
					"🎉 Build Success: {{.ProjectName}} on {{.BuildBranch}}",
				)
				return template
			}(),
			mockError:       nil,
			expectedSubject: "",
			expectedBody:    "🎉 Build Success: test-project on main",
			wantErr:         false,
		},
		{
			name:         "successful_email_build_failure_formatting",
			templateType: domain.TemplateTypeBuildFailure,
			channel:      domain.NotificationChannelEmail,
			params: domain.TemplateParams{
				ProjectName:  testProjectName,
				BuildBranch:  "develop",
				ErrorMessage: "compilation error",
			},
			mockTemplate: func() *domain.NotificationTemplate {
				template, _ := domain.NewNotificationTemplate(
					domain.TemplateTypeBuildFailure,
					domain.NotificationChannelEmail,
					"[BUILD FAILED] {{.ProjectName}} - {{.BuildBranch}}",
					"Build failed for {{.ProjectName}}. Error: {{.ErrorMessage}}",
				)
				return template
			}(),
			mockError:       nil,
			expectedSubject: "[BUILD FAILED] test-project - develop",
			expectedBody:    "Build failed for test-project. Error: compilation error",
			wantErr:         false,
		},
		{
			name:           "template_not_found",
			templateType:   domain.TemplateTypeBuildSuccess,
			channel:        domain.NotificationChannelTelegram,
			params:         domain.TemplateParams{},
			mockTemplate:   nil,
			mockError:      domain.ErrTemplateNotFound,
			wantErr:        true,
			expectedErrMsg: "template not found",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup mock expectations
			mockRepo.On("GetByTypeAndChannel", ctx, tt.templateType, tt.channel).Return(tt.mockTemplate, tt.mockError)

			subject, body, err := formatterService.FormatNotification(ctx, tt.templateType, tt.channel, tt.params)

			if tt.wantErr {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedErrMsg)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedSubject, subject)
				assert.Equal(t, tt.expectedBody, body)
			}

			// Reset mock for next test
			mockRepo.ExpectedCalls = nil
			mockRepo.Calls = nil
		})
	}
}

func TestNotificationFormatterServiceFormatNotificationWithTemplate(t *testing.T) {
	ctx := context.Background()
	mockRepo := new(MockNotificationTemplateRepository)
	logger := logrus.New()

	formatterService := service.NewNotificationFormatterService(service.NotificationFormatterDep{
		TemplateRepo: mockRepo,
		Logger:       logger,
	})

	template, err := domain.NewNotificationTemplate(
		domain.TemplateTypeBuildSuccess,
		domain.NotificationChannelTelegram,
		"",
		"🎉 {{.ProjectName}} build {{.BuildStatus}} on {{.BuildBranch}}",
	)
	assert.NoError(t, err)

	params := domain.TemplateParams{
		ProjectName: testProjectName,
		BuildBranch: "main",
		BuildStatus: "success",
	}

	subject, body, err := formatterService.FormatNotificationWithTemplate(ctx, template, params)

	assert.NoError(t, err)
	assert.Equal(t, "", subject)
	assert.Equal(t, "🎉 test-project build success on main", body)
}

func TestNotificationFormatterServiceValidateTemplate(t *testing.T) {
	mockRepo := new(MockNotificationTemplateRepository)
	logger := logrus.New()

	formatterService := service.NewNotificationFormatterService(service.NotificationFormatterDep{
		TemplateRepo: mockRepo,
		Logger:       logger,
	})

	testParams := domain.TemplateParams{
		ProjectName: "test",
		BuildBranch: "main",
	}

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
			name:         "valid_telegram_template",
			templateType: domain.TemplateTypeBuildSuccess,
			channel:      domain.NotificationChannelTelegram,
			subject:      "",
			bodyTemplate: "🎉 Build Success: {{.ProjectName}}",
			wantErr:      false,
		},
		{
			name:         "valid_email_template",
			templateType: domain.TemplateTypeBuildFailure,
			channel:      domain.NotificationChannelEmail,
			subject:      "[FAILED] {{.ProjectName}}",
			bodyTemplate: "Build failed for {{.ProjectName}}",
			wantErr:      false,
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
		{
			name:         "email_missing_subject",
			templateType: domain.TemplateTypeBuildSuccess,
			channel:      domain.NotificationChannelEmail,
			subject:      "",
			bodyTemplate: "Valid body template",
			wantErr:      true,
			expectedErr:  "subject is required for email",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := formatterService.ValidateTemplate(tt.templateType, tt.channel, tt.subject, tt.bodyTemplate, testParams)

			if tt.wantErr {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedErr)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestNotificationFormatterServiceGetAvailableTemplateVariables(t *testing.T) {
	mockRepo := new(MockNotificationTemplateRepository)
	logger := logrus.New()

	formatterService := service.NewNotificationFormatterService(service.NotificationFormatterDep{
		TemplateRepo: mockRepo,
		Logger:       logger,
	})

	tests := []struct {
		name         string
		templateType domain.NotificationTemplateType
		expected     []string
	}{
		{
			name:         "build_success_variables",
			templateType: domain.TemplateTypeBuildSuccess,
			expected:     []string{"ProjectName", "BuildStatus", "BuildBranch", "BuildCommit", "BuildDuration", "BuildURL", "Timestamp"},
		},
		{
			name:         "build_failure_variables",
			templateType: domain.TemplateTypeBuildFailure,
			expected:     []string{"ProjectName", "BuildStatus", "BuildBranch", "BuildCommit", "BuildDuration", "BuildURL", "ErrorMessage", "Timestamp"},
		},
		{
			name:         "deployment_variables",
			templateType: domain.TemplateTypeDeployment,
			expected:     []string{"ProjectName", "BuildStatus", "BuildBranch", "BuildCommit", "BuildDuration", "BuildURL", "Environment", "Timestamp"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			variables := formatterService.GetAvailableTemplateVariables(tt.templateType)

			assert.ElementsMatch(t, tt.expected, variables)
		})
	}
}

func TestNotificationFormatterServiceFormatEmoji(t *testing.T) {
	mockRepo := new(MockNotificationTemplateRepository)
	logger := logrus.New()

	formatterService := service.NewNotificationFormatterService(service.NotificationFormatterDep{
		TemplateRepo: mockRepo,
		Logger:       logger,
	})

	tests := []struct {
		name     string
		status   string
		channel  domain.NotificationChannel
		expected string
	}{
		{"success_telegram", "success", domain.NotificationChannelTelegram, "✅"},
		{"failure_telegram", "failure", domain.NotificationChannelTelegram, "❌"},
		{"running_telegram", "running", domain.NotificationChannelTelegram, "⏳"},
		{"success_slack", "success", domain.NotificationChannelSlack, ":white_check_mark:"},
		{"failure_slack", "failure", domain.NotificationChannelSlack, ":x:"},
		{"running_slack", "running", domain.NotificationChannelSlack, ":hourglass_flowing_sand:"},
		{"success_email", "success", domain.NotificationChannelEmail, "[SUCCESS]"},
		{"failure_email", "failure", domain.NotificationChannelEmail, "[FAILED]"},
		{"unknown_status", "unknown", domain.NotificationChannelTelegram, ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := formatterService.FormatEmoji(tt.status, tt.channel)
			assert.Equal(t, tt.expected, result)
		})
	}
}
