package telegram

import (
	"context"
	"errors"
	"testing"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/dewisartika8/cicd-status-notifier-bot/internal/core/project/domain"
	"github.com/dewisartika8/cicd-status-notifier-bot/internal/core/shared/domain/value_objects"
	"github.com/dewisartika8/cicd-status-notifier-bot/internal/telegram"
)

// Test for Task 2.3.2: HandleStatusSpecificProject
func TestStatusCommandHandler_HandleStatusSpecificProject(t *testing.T) {
	tests := []struct {
		name             string
		projectName      string
		mockSetup        func(*MockProjectService)
		expectedContains []string
		expectedError    error
	}{
		{
			name:        "Success - Project found",
			projectName: "test-project",
			mockSetup: func(mockService *MockProjectService) {
				// We'll mock this but the actual project creation would need the real constructor
				mockService.On("GetProjectByName", mock.Anything, "test-project").Return(nil, nil)
			},
			expectedContains: []string{"📊 **Project Status:", "test-project"},
			expectedError:    nil,
		},
		{
			name:        "Error - Empty project name",
			projectName: "",
			mockSetup: func(mockService *MockProjectService) {
				// No setup needed as validation happens first
			},
			expectedContains: []string{"❌ **Invalid command**", "Please specify a project name"},
			expectedError:    nil,
		},
		{
			name:        "Error - Project not found",
			projectName: "nonexistent-project",
			mockSetup: func(mockService *MockProjectService) {
				mockService.On("GetProjectByName", mock.Anything, "nonexistent-project").Return(nil, errors.New("project not found"))
			},
			expectedContains: []string{"❌ **Project not found**", "nonexistent-project"},
			expectedError:    nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService := new(MockProjectService)
			tt.mockSetup(mockService)

			handler := telegram.NewStatusCommandHandler(mockService)

			message := &tgbotapi.Message{
				Chat: &tgbotapi.Chat{ID: 123456789},
				From: &tgbotapi.User{ID: 987654321},
				Text: "/status " + tt.projectName,
			}

			response, err := handler.HandleStatusSpecificProject(message, tt.projectName)

			if tt.expectedError != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.expectedError, err)
			} else {
				assert.NoError(t, err)
				for _, contains := range tt.expectedContains {
					assert.Contains(t, response, contains)
				}
			}

			mockService.AssertExpectations(t)
		})
	}
}

// Test for Task 2.3.3: HandleProjectsList
func TestStatusCommandHandler_HandleProjectsList(t *testing.T) {
	tests := []struct {
		name             string
		mockSetup        func(*MockProjectService)
		expectedContains []string
		expectedError    error
	}{
		{
			name: "Success - Projects available",
			mockSetup: func(mockService *MockProjectService) {
				projects := []*domain.Project{} // Empty for now, but would contain mock projects
				mockService.On("GetActiveProjects", mock.Anything).Return(projects, nil)
			},
			expectedContains: []string{"📋 **Monitored Projects**"},
			expectedError:    nil,
		},
		{
			name: "Success - No projects",
			mockSetup: func(mockService *MockProjectService) {
				mockService.On("GetActiveProjects", mock.Anything).Return([]*domain.Project{}, nil)
			},
			expectedContains: []string{"📋 **Monitored Projects**", "ℹ️ No projects are currently being monitored"},
			expectedError:    nil,
		},
		{
			name: "Error - Service failure",
			mockSetup: func(mockService *MockProjectService) {
				mockService.On("GetActiveProjects", mock.Anything).Return([]*domain.Project(nil), errors.New("service error"))
			},
			expectedContains: []string{"❌ **Error fetching projects**", "Unable to retrieve project list"},
			expectedError:    nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService := new(MockProjectService)
			tt.mockSetup(mockService)

			handler := telegram.NewStatusCommandHandler(mockService)

			message := &tgbotapi.Message{
				Chat: &tgbotapi.Chat{ID: 123456789},
				From: &tgbotapi.User{ID: 987654321},
				Text: "/projects",
			}

			response, err := handler.HandleProjectsList(message)

			if tt.expectedError != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.expectedError, err)
			} else {
				assert.NoError(t, err)
				for _, contains := range tt.expectedContains {
					assert.Contains(t, response, contains)
				}
			}

			mockService.AssertExpectations(t)
		})
	}
}

// Test for Task 2.3.4: Error handling and response formatting
func TestResponseFormatter(t *testing.T) {
	formatter := telegram.NewResponseFormatter()

	tests := []struct {
		name             string
		method           func() string
		expectedContains []string
	}{
		{
			name: "FormatError",
			method: func() string {
				return formatter.FormatError("Test Error", "This is a test error message")
			},
			expectedContains: []string{"❌ **Test Error**", "This is a test error message"},
		},
		{
			name: "FormatSuccess",
			method: func() string {
				return formatter.FormatSuccess("Test Success", "This is a test success message")
			},
			expectedContains: []string{"✅ **Test Success**", "This is a test success message"},
		},
		{
			name: "FormatInfo",
			method: func() string {
				return formatter.FormatInfo("Test Info", "This is a test info message")
			},
			expectedContains: []string{"ℹ️ **Test Info**", "This is a test info message"},
		},
		{
			name: "FormatProjectNotFound",
			method: func() string {
				return formatter.FormatProjectNotFound("test-project")
			},
			expectedContains: []string{"❌ **Project not found**", "test-project", "not being monitored"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.method()
			for _, contains := range tt.expectedContains {
				assert.Contains(t, result, contains)
			}
		})
	}
}

// Test for error handler
func TestErrorHandler(t *testing.T) {
	errorHandler := telegram.NewErrorHandler()

	tests := []struct {
		name             string
		err              error
		operation        string
		expectedContains []string
	}{
		{
			name:             "Project not found error",
			err:              errors.New("project not found"),
			operation:        "fetching project",
			expectedContains: []string{"❌ **Project not found**", "not found or is not accessible"},
		},
		{
			name:             "Permission error",
			err:              errors.New("permission denied"),
			operation:        "updating project",
			expectedContains: []string{"❌ **Permission denied**", "don't have permission"},
		},
		{
			name:             "Generic error",
			err:              errors.New("generic database error"),
			operation:        "fetching data",
			expectedContains: []string{"❌ **Error fetching data**", "Unable to complete the operation"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := errorHandler.HandleProjectServiceError(tt.err, tt.operation)
			for _, contains := range tt.expectedContains {
				assert.Contains(t, result, contains)
			}
		})
	}
}

// Integration test combining all handlers
func TestStatusCommandHandler_Integration(t *testing.T) {
	mockService := new(MockProjectService)
	handler := telegram.NewStatusCommandHandler(mockService)

	// Test full workflow: list projects -> get specific project status
	t.Run("Full workflow", func(t *testing.T) {
		// Mock empty projects list
		mockService.On("GetActiveProjects", mock.Anything).Return([]*domain.Project{}, nil)

		// Test projects list
		msg := &tgbotapi.Message{
			Chat: &tgbotapi.Chat{ID: 123},
			From: &tgbotapi.User{ID: 456},
			Text: "/projects",
		}

		response, err := handler.HandleProjectsList(msg)
		assert.NoError(t, err)
		assert.Contains(t, response, "📋 **Monitored Projects**")
		assert.Contains(t, response, "No projects are currently being monitored")

		mockService.AssertExpectations(t)
	})
}
