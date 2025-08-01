package service_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/dewisartika8/cicd-status-notifier-bot/internal/core/bot/domain"
	"github.com/dewisartika8/cicd-status-notifier-bot/internal/core/bot/dto"
	"github.com/dewisartika8/cicd-status-notifier-bot/internal/core/bot/port"
	"github.com/dewisartika8/cicd-status-notifier-bot/internal/core/bot/service"
)

// Constants for mock types to avoid duplication
const (
	commandContextType = "*domain.CommandContext"
	stringType         = "string"
)

// Mock implementations for testing
type MockTelegramAPI struct {
	mock.Mock
}

func (m *MockTelegramAPI) SendMessage(chatID int64, text string) error {
	args := m.Called(chatID, text)
	return args.Error(0)
}

func (m *MockTelegramAPI) SendMessageWithMarkdown(chatID int64, text string) error {
	// This method specifically handles markdown formatting
	// In real implementation, it would parse and format markdown
	args := m.Called(chatID, text)
	if args.Error(0) != nil {
		return args.Error(0)
	}
	return nil
}

func (m *MockTelegramAPI) SetWebhook(webhookURL string) error {
	args := m.Called(webhookURL)
	return args.Error(0)
}

func (m *MockTelegramAPI) DeleteWebhook() error {
	args := m.Called()
	return args.Error(0)
}

type MockCommandValidator struct {
	mock.Mock
}

func (m *MockCommandValidator) ValidateCommand(ctx *domain.CommandContext) error {
	args := m.Called(ctx)
	return args.Error(0)
}

func (m *MockCommandValidator) AddAllowedUser(userID int64) {
	m.Called(userID)
}

type MockCommandRouter struct {
	mock.Mock
}

func (m *MockCommandRouter) RegisterHandler(command string, handler domain.CommandHandler) {
	m.Called(command, handler)
}

func (m *MockCommandRouter) RouteCommand(ctx *domain.CommandContext) error {
	args := m.Called(ctx)
	return args.Error(0)
}

type MockProjectService struct {
	mock.Mock
}

func (m *MockProjectService) GetProject(ctx context.Context, name string) (*port.Project, error) {
	args := m.Called(ctx, name)
	return args.Get(0).(*port.Project), args.Error(1)
}

func (m *MockProjectService) GetAllProjects(ctx context.Context) ([]*port.Project, error) {
	args := m.Called(ctx)
	return args.Get(0).([]*port.Project), args.Error(1)
}

func (m *MockProjectService) GetProjectStatus(ctx context.Context, projectName string) (*dto.StatusCommandResponse, error) {
	args := m.Called(ctx, projectName)
	return args.Get(0).(*dto.StatusCommandResponse), args.Error(1)
}

type MockSubscriptionService struct {
	mock.Mock
}

func (m *MockSubscriptionService) Subscribe(ctx context.Context, req *dto.SubscribeCommandRequest) error {
	args := m.Called(ctx, req)
	return args.Error(0)
}

func (m *MockSubscriptionService) Unsubscribe(ctx context.Context, req *dto.UnsubscribeCommandRequest) error {
	args := m.Called(ctx, req)
	return args.Error(0)
}

func (m *MockSubscriptionService) IsSubscribed(ctx context.Context, chatID int64, projectName string) (bool, error) {
	args := m.Called(ctx, chatID, projectName)
	return args.Bool(0), args.Error(1)
}

func (m *MockSubscriptionService) GetSubscriptions(ctx context.Context, chatID int64) ([]*port.Subscription, error) {
	args := m.Called(ctx, chatID)
	return args.Get(0).([]*port.Subscription), args.Error(1)
}

func TestBotServiceHandleStartCommand(t *testing.T) {
	tests := []struct {
		name          string
		request       *dto.StartCommandRequest
		mockSetup     func(*MockTelegramAPI)
		expectedError bool
	}{
		{
			name: "should handle start command successfully",
			request: &dto.StartCommandRequest{
				ChatID:        12345,
				UserFirstName: "John",
			},
			mockSetup: func(api *MockTelegramAPI) {
				api.On("SendMessageWithMarkdown", int64(12345), mock.AnythingOfType(stringType)).Return(nil)
			},
			expectedError: false,
		},
		{
			name: "should handle telegram API error",
			request: &dto.StartCommandRequest{
				ChatID:        12345,
				UserFirstName: "John",
			},
			mockSetup: func(api *MockTelegramAPI) {
				api.On("SendMessageWithMarkdown", int64(12345), mock.AnythingOfType(stringType)).Return(assert.AnError)
			},
			expectedError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockAPI := &MockTelegramAPI{}
			mockValidator := &MockCommandValidator{}
			mockRouter := &MockCommandRouter{}
			mockProjectService := &MockProjectService{}
			mockSubscriptionService := &MockSubscriptionService{}

			// Setup router expectations for constructor
			mockRouter.On("RegisterHandler", mock.AnythingOfType("string"), mock.Anything).Return().Times(5)

			tt.mockSetup(mockAPI)

			botService := service.NewBotService(
				mockAPI,
				mockValidator,
				mockRouter,
				mockProjectService,
				mockSubscriptionService,
			)

			response, err := botService.HandleStartCommand(context.Background(), tt.request)

			if tt.expectedError {
				assert.Error(t, err)
				assert.Nil(t, response)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, response)
				assert.Equal(t, tt.request.ChatID, response.ChatID)
				assert.Equal(t, tt.request.UserFirstName, response.UserFirstName)
				assert.Contains(t, response.WelcomeMessage, "Welcome to CICD Status Notifier Bot")
			}

			mockAPI.AssertExpectations(t)
		})
	}
}

func TestBotServiceHandleHelpCommand(t *testing.T) {
	tests := []struct {
		name          string
		request       *port.HelpCommandRequest
		mockSetup     func(*MockTelegramAPI)
		expectedError bool
	}{
		{
			name: "should handle help command successfully",
			request: &port.HelpCommandRequest{
				ChatID: 12345,
				UserID: 67890,
			},
			mockSetup: func(api *MockTelegramAPI) {
				api.On("SendMessageWithMarkdown", int64(12345), mock.AnythingOfType("string")).Return(nil)
			},
			expectedError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockAPI := &MockTelegramAPI{}
			mockValidator := &MockCommandValidator{}
			mockRouter := &MockCommandRouter{}
			mockProjectService := &MockProjectService{}
			mockSubscriptionService := &MockSubscriptionService{}

			// Setup router expectations for constructor
			mockRouter.On("RegisterHandler", mock.AnythingOfType("string"), mock.Anything).Return().Times(5)

			tt.mockSetup(mockAPI)

			botService := service.NewBotService(
				mockAPI,
				mockValidator,
				mockRouter,
				mockProjectService,
				mockSubscriptionService,
			)

			response, err := botService.HandleHelpCommand(context.Background(), tt.request)

			if tt.expectedError {
				assert.Error(t, err)
				assert.Nil(t, response)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, response)
				assert.Contains(t, response.HelpText, "CICD Status Notifier Bot - Help")
				assert.NotEmpty(t, response.Commands)
				assert.NotEmpty(t, response.UsageExamples)
			}

			mockAPI.AssertExpectations(t)
		})
	}
}

func TestBotServiceHandleStatusCommand(t *testing.T) {
	tests := []struct {
		name          string
		request       *dto.StatusCommandRequest
		mockSetup     func(*MockTelegramAPI)
		expectedError bool
	}{
		{
			name: "should handle status command with project name",
			request: &dto.StatusCommandRequest{
				ProjectName: "my-project",
				ChatID:      12345,
				UserID:      67890,
			},
			mockSetup: func(api *MockTelegramAPI) {
				api.On("SendMessageWithMarkdown", int64(12345), mock.AnythingOfType("string")).Return(nil)
			},
			expectedError: false,
		},
		{
			name: "should handle status command without project name",
			request: &dto.StatusCommandRequest{
				ProjectName: "",
				ChatID:      12345,
				UserID:      67890,
			},
			mockSetup: func(api *MockTelegramAPI) {
				api.On("SendMessageWithMarkdown", int64(12345), mock.AnythingOfType("string")).Return(nil)
			},
			expectedError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockAPI := &MockTelegramAPI{}
			mockValidator := &MockCommandValidator{}
			mockRouter := &MockCommandRouter{}
			mockProjectService := &MockProjectService{}
			mockSubscriptionService := &MockSubscriptionService{}

			// Setup router expectations for constructor
			mockRouter.On("RegisterHandler", mock.AnythingOfType("string"), mock.Anything).Return().Times(5)

			tt.mockSetup(mockAPI)

			botService := service.NewBotService(
				mockAPI,
				mockValidator,
				mockRouter,
				mockProjectService,
				mockSubscriptionService,
			)

			response, err := botService.HandleStatusCommand(context.Background(), tt.request)

			if tt.expectedError {
				assert.Error(t, err)
				assert.Nil(t, response)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, response)

				expectedProject := tt.request.ProjectName
				if expectedProject == "" {
					expectedProject = "all"
				}
				assert.Equal(t, expectedProject, response.ProjectName)
				assert.Equal(t, "operational", response.Status)
				assert.Contains(t, response.Message, "Pipeline Status")
			}

			mockAPI.AssertExpectations(t)
		})
	}
}

func TestBotServiceHandleCommand(t *testing.T) {
	tests := []struct {
		name          string
		commandCtx    *domain.CommandContext
		validatorErr  error
		routerErr     error
		mockSetup     func(*MockTelegramAPI, *MockCommandValidator, *MockCommandRouter)
		expectedError bool
	}{
		{
			name: "should handle valid command successfully",
			commandCtx: &domain.CommandContext{
				Command: "help",
				Args:    []string{},
				UserID:  12345,
				ChatID:  67890,
			},
			validatorErr: nil,
			routerErr:    nil,
			mockSetup: func(api *MockTelegramAPI, validator *MockCommandValidator, router *MockCommandRouter) {
				validator.On("ValidateCommand", mock.AnythingOfType(commandContextType)).Return(nil)
				router.On("RouteCommand", mock.AnythingOfType(commandContextType)).Return(nil)
			},
			expectedError: false,
		},
		{
			name: "should handle validation error",
			commandCtx: &domain.CommandContext{
				Command: "invalid",
				Args:    []string{},
				UserID:  12345,
				ChatID:  67890,
			},
			validatorErr: assert.AnError,
			routerErr:    nil,
			mockSetup: func(api *MockTelegramAPI, validator *MockCommandValidator, router *MockCommandRouter) {
				validator.On("ValidateCommand", mock.AnythingOfType(commandContextType)).Return(assert.AnError)
				api.On("SendMessage", int64(67890), mock.AnythingOfType(stringType)).Return(nil)
			},
			expectedError: false, // Error is handled by sending message
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockAPI := &MockTelegramAPI{}
			mockValidator := &MockCommandValidator{}
			mockRouter := &MockCommandRouter{}
			mockProjectService := &MockProjectService{}
			mockSubscriptionService := &MockSubscriptionService{}

			// Setup router expectations for constructor
			mockRouter.On("RegisterHandler", mock.AnythingOfType("string"), mock.Anything).Return().Times(5)

			tt.mockSetup(mockAPI, mockValidator, mockRouter)

			botService := service.NewBotService(
				mockAPI,
				mockValidator,
				mockRouter,
				mockProjectService,
				mockSubscriptionService,
			)

			err := botService.HandleCommand(context.Background(), tt.commandCtx)

			if tt.expectedError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			mockAPI.AssertExpectations(t)
			mockValidator.AssertExpectations(t)
			mockRouter.AssertExpectations(t)
		})
	}
}
