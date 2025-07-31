package service

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
	testProjectNameFormatter = "test-project"
)

// MockNotificationTemplateRepositoryFormatter is a mock implementation of NotificationTemplateRepository
type MockNotificationTemplateRepositoryFormatter struct {
	mock.Mock
}

func (m *MockNotificationTemplateRepositoryFormatter) Create(ctx context.Context, template *domain.NotificationTemplate) error {
	args := m.Called(ctx, template)
	return args.Error(0)
}

func (m *MockNotificationTemplateRepositoryFormatter) GetByID(ctx context.Context, id value_objects.ID) (*domain.NotificationTemplate, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.NotificationTemplate), args.Error(1)
}

func (m *MockNotificationTemplateRepositoryFormatter) GetByTypeAndChannel(ctx context.Context, templateType domain.NotificationTemplateType, channel domain.NotificationChannel) (*domain.NotificationTemplate, error) {
	args := m.Called(ctx, templateType, channel)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.NotificationTemplate), args.Error(1)
}

func (m *MockNotificationTemplateRepositoryFormatter) GetByType(ctx context.Context, templateType domain.NotificationTemplateType) ([]*domain.NotificationTemplate, error) {
	args := m.Called(ctx, templateType)
	return args.Get(0).([]*domain.NotificationTemplate), args.Error(1)
}

func (m *MockNotificationTemplateRepositoryFormatter) GetByChannel(ctx context.Context, channel domain.NotificationChannel) ([]*domain.NotificationTemplate, error) {
	args := m.Called(ctx, channel)
	return args.Get(0).([]*domain.NotificationTemplate), args.Error(1)
}

func (m *MockNotificationTemplateRepositoryFormatter) GetActiveTemplates(ctx context.Context) ([]*domain.NotificationTemplate, error) {
	args := m.Called(ctx)
	return args.Get(0).([]*domain.NotificationTemplate), args.Error(1)
}

func (m *MockNotificationTemplateRepositoryFormatter) Update(ctx context.Context, template *domain.NotificationTemplate) error {
	// Simulate update by marking the template as updated in mock, or just call the mock with a different method name for clarity
	args := m.Called(ctx, template)
	// Optionally, you could add a log or state change here to differentiate from Create
	return args.Error(0)
}

func (m *MockNotificationTemplateRepositoryFormatter) Delete(ctx context.Context, id value_objects.ID) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockNotificationTemplateRepositoryFormatter) Count(ctx context.Context, templateType *domain.NotificationTemplateType, channel *domain.NotificationChannel, isActive *bool) (int64, error) {
	args := m.Called(ctx, templateType, channel, isActive)
	return args.Get(0).(int64), args.Error(1)
}

func TestFormatNotification(t *testing.T) {
	ctx := context.Background()
	mockRepo := new(MockNotificationTemplateRepositoryFormatter)
	logger := logrus.New()

	formatterService := service.NewNotificationFormatterService(service.NotificationFormatterDep{
		TemplateRepo: mockRepo,
		Logger:       logger,
	})

	// Test successful formatting
	template, err := domain.NewNotificationTemplate(
		domain.TemplateTypeBuildSuccess,
		domain.NotificationChannelTelegram,
		"",
		"🎉 Build Success: {{.ProjectName}} on {{.BuildBranch}}",
	)
	assert.NoError(t, err)

	mockRepo.On("GetByTypeAndChannel", ctx, domain.TemplateTypeBuildSuccess, domain.NotificationChannelTelegram).Return(template, nil)

	params := domain.TemplateParams{
		ProjectName: testProjectNameFormatter,
		BuildBranch: "main",
		BuildStatus: "success",
	}

	subject, body, err := formatterService.FormatNotification(ctx, domain.TemplateTypeBuildSuccess, domain.NotificationChannelTelegram, params)

	assert.NoError(t, err)
	assert.Equal(t, "", subject)
	assert.Contains(t, body, testProjectNameFormatter)
	assert.Contains(t, body, "main")

	mockRepo.AssertExpectations(t)
}

func TestFormatNotificationWithTemplate(t *testing.T) {
	ctx := context.Background()
	mockRepo := new(MockNotificationTemplateRepositoryFormatter)
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
		ProjectName: testProjectNameFormatter,
		BuildBranch: "main",
		BuildStatus: "success",
	}

	subject, body, err := formatterService.FormatNotificationWithTemplate(ctx, template, params)

	assert.NoError(t, err)
	assert.Equal(t, "", subject)
	assert.Contains(t, body, testProjectNameFormatter)
	assert.Contains(t, body, "success")
	assert.Contains(t, body, "main")
}

func TestValidateTemplate(t *testing.T) {
	mockRepo := new(MockNotificationTemplateRepositoryFormatter)
	logger := logrus.New()

	formatterService := service.NewNotificationFormatterService(service.NotificationFormatterDep{
		TemplateRepo: mockRepo,
		Logger:       logger,
	})

	testParams := domain.TemplateParams{
		ProjectName: "test",
		BuildBranch: "main",
	}

	// Test valid template
	err := formatterService.ValidateTemplate(
		domain.TemplateTypeBuildSuccess,
		domain.NotificationChannelTelegram,
		"",
		"🎉 Build Success: {{.ProjectName}}",
		testParams,
	)
	assert.NoError(t, err)

	// Test invalid template syntax
	err = formatterService.ValidateTemplate(
		domain.TemplateTypeBuildSuccess,
		domain.NotificationChannelTelegram,
		"",
		"{{.InvalidSyntax",
		testParams,
	)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "template compilation failed")
}

func TestGetAvailableTemplateVariables(t *testing.T) {
	mockRepo := new(MockNotificationTemplateRepositoryFormatter)
	logger := logrus.New()

	formatterService := service.NewNotificationFormatterService(service.NotificationFormatterDep{
		TemplateRepo: mockRepo,
		Logger:       logger,
	})

	variables := formatterService.GetAvailableTemplateVariables(domain.TemplateTypeBuildSuccess)
	expected := []string{"ProjectName", "BuildStatus", "BuildBranch", "BuildCommit", "BuildDuration", "BuildURL", "Timestamp"}

	assert.ElementsMatch(t, expected, variables)

	// Test build failure includes ErrorMessage
	variables = formatterService.GetAvailableTemplateVariables(domain.TemplateTypeBuildFailure)
	assert.Contains(t, variables, "ErrorMessage")
}

func TestFormatEmoji(t *testing.T) {
	mockRepo := new(MockNotificationTemplateRepositoryFormatter)
	logger := logrus.New()

	formatterService := service.NewNotificationFormatterService(service.NotificationFormatterDep{
		TemplateRepo: mockRepo,
		Logger:       logger,
	})

	// Test Telegram emojis
	assert.Equal(t, "✅", formatterService.FormatEmoji("success", domain.NotificationChannelTelegram))
	assert.Equal(t, "❌", formatterService.FormatEmoji("failure", domain.NotificationChannelTelegram))
	assert.Equal(t, "⏳", formatterService.FormatEmoji("running", domain.NotificationChannelTelegram))

	// Test Slack emojis
	assert.Equal(t, ":white_check_mark:", formatterService.FormatEmoji("success", domain.NotificationChannelSlack))
	assert.Equal(t, ":x:", formatterService.FormatEmoji("failure", domain.NotificationChannelSlack))

	// Test email prefixes
	assert.Equal(t, "[SUCCESS]", formatterService.FormatEmoji("success", domain.NotificationChannelEmail))
	assert.Equal(t, "[FAILED]", formatterService.FormatEmoji("failure", domain.NotificationChannelEmail))

	// Test unknown status
	assert.Equal(t, "", formatterService.FormatEmoji("unknown", domain.NotificationChannelTelegram))
}
