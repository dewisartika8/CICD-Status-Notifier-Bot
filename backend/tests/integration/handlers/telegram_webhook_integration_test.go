package handlers_test

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/dewisartika8/cicd-status-notifier-bot/internal/adapter/handler/webhook"
	"github.com/dewisartika8/cicd-status-notifier-bot/internal/core/bot/domain"
	"github.com/dewisartika8/cicd-status-notifier-bot/internal/core/bot/dto"
	"github.com/dewisartika8/cicd-status-notifier-bot/internal/core/bot/port"
)

const (
	// Test constants
	telegramWebhookEndpoint = "/api/v1/telegram/webhook"
	invalidJSONPayload      = "invalid json"
)

// Mock BotService for integration testing
type MockBotService struct {
	mock.Mock
}

func (m *MockBotService) HandleCommand(ctx context.Context, commandCtx *domain.CommandContext) error {
	args := m.Called(ctx, commandCtx)
	return args.Error(0)
}

func (m *MockBotService) SendMessage(ctx context.Context, chatID int64, message string) error {
	args := m.Called(ctx, chatID, message)
	return args.Error(0)
}

func (m *MockBotService) SendFormattedMessage(ctx context.Context, chatID int64, message string, parseMode string) error {
	args := m.Called(ctx, chatID, message, parseMode)
	return args.Error(0)
}

func (m *MockBotService) SetWebhook(ctx context.Context, webhookURL string) error {
	args := m.Called(ctx, webhookURL)
	return args.Error(0)
}

func (m *MockBotService) DeleteWebhook(ctx context.Context) error {
	args := m.Called(ctx)
	return args.Error(0)
}

// Methods from the interface that aren't used in these tests but need to be implemented
func (m *MockBotService) HandleStartCommand(ctx context.Context, req *dto.StartCommandRequest) (*dto.StartCommandResponse, error) {
	args := m.Called(ctx, req)
	return args.Get(0).(*dto.StartCommandResponse), args.Error(1)
}

func (m *MockBotService) HandleHelpCommand(ctx context.Context, req *port.HelpCommandRequest) (*dto.HelpCommandResponse, error) {
	args := m.Called(ctx, req)
	return args.Get(0).(*dto.HelpCommandResponse), args.Error(1)
}

func (m *MockBotService) HandleStatusCommand(ctx context.Context, req *dto.StatusCommandRequest) (*dto.StatusCommandResponse, error) {
	args := m.Called(ctx, req)
	return args.Get(0).(*dto.StatusCommandResponse), args.Error(1)
}

func (m *MockBotService) HandleSubscribeCommand(ctx context.Context, req *dto.SubscribeCommandRequest) (*dto.SubscribeCommandResponse, error) {
	args := m.Called(ctx, req)
	return args.Get(0).(*dto.SubscribeCommandResponse), args.Error(1)
}

func (m *MockBotService) HandleUnsubscribeCommand(ctx context.Context, req *dto.UnsubscribeCommandRequest) (*dto.UnsubscribeCommandResponse, error) {
	args := m.Called(ctx, req)
	return args.Get(0).(*dto.UnsubscribeCommandResponse), args.Error(1)
}

// Mock CommandValidator for integration testing
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

func TestTelegramWebhookHandlerHandleTelegramWebhook(t *testing.T) {
	tests := []struct {
		name           string
		requestBody    interface{}
		mockSetup      func(*MockBotService, *MockCommandValidator)
		expectedStatus int
		expectedBody   map[string]interface{}
	}{
		{
			name: "should handle valid webhook with command message",
			requestBody: tgbotapi.Update{
				Message: &tgbotapi.Message{
					MessageID: 123,
					From: &tgbotapi.User{
						ID:       12345,
						UserName: "testuser",
					},
					Chat: &tgbotapi.Chat{
						ID: 67890,
					},
					Text: "/help",
					Entities: []tgbotapi.MessageEntity{
						{
							Type:   "bot_command",
							Offset: 0,
							Length: 5,
						},
					},
				},
			},
			mockSetup: func(botService *MockBotService, validator *MockCommandValidator) {
				botService.On("HandleCommand", mock.Anything, mock.MatchedBy(func(ctx *domain.CommandContext) bool {
					return ctx.Command == "help" && ctx.UserID == 12345 && ctx.ChatID == 67890
				})).Return(nil)
			},
			expectedStatus: 200,
			expectedBody:   map[string]interface{}{"status": "ok"},
		},
		{
			name: "should handle valid webhook with non-command message",
			requestBody: tgbotapi.Update{
				Message: &tgbotapi.Message{
					MessageID: 123,
					From: &tgbotapi.User{
						ID:       12345,
						UserName: "testuser",
					},
					Chat: &tgbotapi.Chat{
						ID: 67890,
					},
					Text: "hello bot",
				},
			},
			mockSetup: func(botService *MockBotService, validator *MockCommandValidator) {
				botService.On("SendMessage", mock.Anything, int64(67890), "Please send a valid command. Type /help for available commands.").Return(nil)
			},
			expectedStatus: 200,
			expectedBody:   map[string]interface{}{"status": "ok"},
		},
		{
			name:        "should handle invalid JSON payload",
			requestBody: invalidJSONPayload,
			mockSetup: func(botService *MockBotService, validator *MockCommandValidator) {
				// No mocks needed for invalid JSON test case
			},
			expectedStatus: 400,
			expectedBody:   map[string]interface{}{"error": "Invalid JSON payload"},
		},
		{
			name: "should handle webhook without message",
			requestBody: tgbotapi.Update{
				UpdateID: 123,
			},
			mockSetup: func(botService *MockBotService, validator *MockCommandValidator) {
				// No mocks needed for webhook without message test case
			},
			expectedStatus: 200,
			expectedBody:   map[string]interface{}{"status": "ok"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup mocks
			mockBotService := &MockBotService{}
			mockValidator := &MockCommandValidator{}

			tt.mockSetup(mockBotService, mockValidator)

			// Create webhook handler
			handler := webhook.NewTelegramWebhookHandler(mockBotService, mockValidator)

			// Setup Fiber app
			app := fiber.New()
			app.Post(telegramWebhookEndpoint, handler.HandleTelegramWebhook)

			// Prepare request body
			var bodyReader io.Reader
			if tt.requestBody == invalidJSONPayload {
				bodyReader = bytes.NewBufferString(invalidJSONPayload)
			} else {
				bodyBytes, _ := json.Marshal(tt.requestBody)
				bodyReader = bytes.NewReader(bodyBytes)
			}

			// Create test request
			req := httptest.NewRequest("POST", telegramWebhookEndpoint, bodyReader)
			req.Header.Set("Content-Type", "application/json")

			// Perform request
			resp, err := app.Test(req, -1)
			assert.NoError(t, err)

			// Check status code
			assert.Equal(t, tt.expectedStatus, resp.StatusCode)

			// Check response body
			var responseBody map[string]interface{}
			err = json.NewDecoder(resp.Body).Decode(&responseBody)
			assert.NoError(t, err)
			assert.Equal(t, tt.expectedBody, responseBody)

			// Assert mock expectations
			mockBotService.AssertExpectations(t)
			mockValidator.AssertExpectations(t)
		})
	}
}

func setupMockBotServiceForCommand(t *testing.T, expectedCommand string, expectedArgs []string, expectedHandlerCall bool) *MockBotService {
	mockBotService := &MockBotService{}
	if expectedHandlerCall {
		mockBotService.On("HandleCommand", mock.Anything, mock.MatchedBy(func(ctx *domain.CommandContext) bool {
			if ctx.Command != expectedCommand {
				return false
			}
			if len(ctx.Args) != len(expectedArgs) {
				return false
			}
			for i, arg := range expectedArgs {
				if ctx.Args[i] != arg {
					return false
				}
			}
			return true
		})).Return(nil)
	}
	return mockBotService
}

func createTelegramUpdate(messageText string, messageEntities []tgbotapi.MessageEntity) tgbotapi.Update {
	return tgbotapi.Update{
		Message: &tgbotapi.Message{
			MessageID: 123,
			From: &tgbotapi.User{
				ID:       12345,
				UserName: "testuser",
			},
			Chat: &tgbotapi.Chat{
				ID: 67890,
			},
			Text:     messageText,
			Entities: messageEntities,
		},
	}
}

func performTelegramWebhookTest(t *testing.T, mockBotService *MockBotService, mockValidator *MockCommandValidator, update tgbotapi.Update) *http.Response {
	handler := webhook.NewTelegramWebhookHandler(mockBotService, mockValidator)
	app := fiber.New()
	app.Post(telegramWebhookEndpoint, handler.HandleTelegramWebhook)
	bodyBytes, _ := json.Marshal(update)
	req := httptest.NewRequest("POST", telegramWebhookEndpoint, bytes.NewReader(bodyBytes))
	req.Header.Set("Content-Type", "application/json")
	resp, err := app.Test(req, -1)
	assert.NoError(t, err)
	return resp
}

func TestTelegramWebhookHandlerHandleTelegramWebhookCommandParsing(t *testing.T) {
	tests := []struct {
		name                string
		messageText         string
		messageEntities     []tgbotapi.MessageEntity
		expectedCommand     string
		expectedArgs        []string
		expectedHandlerCall bool
	}{
		{
			name:        "should parse help command",
			messageText: "/help",
			messageEntities: []tgbotapi.MessageEntity{
				{Type: "bot_command", Offset: 0, Length: 5},
			},
			expectedCommand:     "help",
			expectedArgs:        []string{},
			expectedHandlerCall: true,
		},
		{
			name:        "should parse status command with arguments",
			messageText: "/status my-project",
			messageEntities: []tgbotapi.MessageEntity{
				{Type: "bot_command", Offset: 0, Length: 7},
			},
			expectedCommand:     "status",
			expectedArgs:        []string{"my-project"},
			expectedHandlerCall: true,
		},
		{
			name:        "should parse subscribe command with project name",
			messageText: "/subscribe my-awesome-project",
			messageEntities: []tgbotapi.MessageEntity{
				{Type: "bot_command", Offset: 0, Length: 10},
			},
			expectedCommand:     "subscribe",
			expectedArgs:        []string{"my-awesome-project"},
			expectedHandlerCall: true,
		},
		{
			name:        "should parse command with multiple arguments",
			messageText: "/command arg1 arg2 arg3",
			messageEntities: []tgbotapi.MessageEntity{
				{Type: "bot_command", Offset: 0, Length: 8},
			},
			expectedCommand:     "command",
			expectedArgs:        []string{"arg1", "arg2", "arg3"},
			expectedHandlerCall: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockBotService := setupMockBotServiceForCommand(t, tt.expectedCommand, tt.expectedArgs, tt.expectedHandlerCall)
			mockValidator := &MockCommandValidator{}
			update := createTelegramUpdate(tt.messageText, tt.messageEntities)
			resp := performTelegramWebhookTest(t, mockBotService, mockValidator, update)
			assert.Equal(t, 200, resp.StatusCode)
			mockBotService.AssertExpectations(t)
		})
	}
}
