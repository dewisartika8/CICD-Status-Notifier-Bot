package domain_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/dewisartika8/cicd-status-notifier-bot/internal/core/bot/domain"
)

func TestCommandValidator_ValidateCommand(t *testing.T) {
	tests := []struct {
		name          string
		ctx           *domain.CommandContext
		setupUser     func(*domain.CommandValidator)
		expectedError string
	}{
		{
			name: "should validate basic command successfully",
			ctx: &domain.CommandContext{
				Command:  "help",
				Args:     []string{},
				UserID:   12345,
				ChatID:   67890,
				Username: "testuser",
			},
			setupUser:     func(cv *domain.CommandValidator) {},
			expectedError: "",
		},
		{
			name: "should validate status command with project argument",
			ctx: &domain.CommandContext{
				Command:  "status",
				Args:     []string{"my-project"},
				UserID:   12345,
				ChatID:   67890,
				Username: "testuser",
			},
			setupUser:     func(cv *domain.CommandValidator) {},
			expectedError: "",
		},
		{
			name: "should reject invalid command",
			ctx: &domain.CommandContext{
				Command:  "invalid",
				Args:     []string{},
				UserID:   12345,
				ChatID:   67890,
				Username: "testuser",
			},
			setupUser:     func(cv *domain.CommandValidator) {},
			expectedError: "invalid command",
		},
		{
			name: "should reject subscribe command without project name",
			ctx: &domain.CommandContext{
				Command:  "subscribe",
				Args:     []string{},
				UserID:   12345,
				ChatID:   67890,
				Username: "testuser",
			},
			setupUser: func(cv *domain.CommandValidator) {
				cv.AddAllowedUser(12345)
			},
			expectedError: "project name is required",
		},
		{
			name: "should reject subscribe command with short project name",
			ctx: &domain.CommandContext{
				Command:  "subscribe",
				Args:     []string{"a"},
				UserID:   12345,
				ChatID:   67890,
				Username: "testuser",
			},
			setupUser: func(cv *domain.CommandValidator) {
				cv.AddAllowedUser(12345)
			},
			expectedError: "project name too short",
		},
		{
			name: "should accept subscribe command with valid project name",
			ctx: &domain.CommandContext{
				Command:  "subscribe",
				Args:     []string{"my-project"},
				UserID:   12345,
				ChatID:   67890,
				Username: "testuser",
			},
			setupUser: func(cv *domain.CommandValidator) {
				cv.AddAllowedUser(12345)
			},
			expectedError: "",
		},
		{
			name: "should reject subscribe command for unauthorized user",
			ctx: &domain.CommandContext{
				Command:  "subscribe",
				Args:     []string{"my-project"},
				UserID:   99999,
				ChatID:   67890,
				Username: "testuser",
			},
			setupUser:     func(cv *domain.CommandValidator) {},
			expectedError: "insufficient permissions",
		},
		{
			name: "should reject status command with too many arguments",
			ctx: &domain.CommandContext{
				Command:  "status",
				Args:     []string{"project1", "project2"},
				UserID:   12345,
				ChatID:   67890,
				Username: "testuser",
			},
			setupUser:     func(cv *domain.CommandValidator) {},
			expectedError: "too many arguments for status command",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			validator := domain.NewCommandValidator()
			tt.setupUser(validator)

			err := validator.ValidateCommand(tt.ctx)

			if tt.expectedError == "" {
				assert.NoError(t, err)
			} else {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedError)
			}
		})
	}
}

func TestCommandValidator_AddAllowedUser(t *testing.T) {
	validator := domain.NewCommandValidator()
	userID := int64(12345)

	// Initially user should not have permissions for restricted commands
	ctx := &domain.CommandContext{
		Command: "subscribe",
		Args:    []string{"test-project"},
		UserID:  userID,
	}

	err := validator.ValidateCommand(ctx)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "insufficient permissions")

	// After adding user, they should have permissions
	validator.AddAllowedUser(userID)

	err = validator.ValidateCommand(ctx)
	assert.NoError(t, err)
}

func TestCommandRouter_RegisterHandler(t *testing.T) {
	router := domain.NewCommandRouter()
	handler := &mockCommandHandler{handleFunc: func(ctx *domain.CommandContext) error { return nil }}

	router.RegisterHandler("test", handler)

	ctx := &domain.CommandContext{
		Command: "test",
		Args:    []string{},
		UserID:  12345,
		ChatID:  67890,
	}

	err := router.RouteCommand(ctx)
	assert.NoError(t, err)
	assert.True(t, handler.called)
}

func TestCommandRouter_RouteCommand_UnknownCommand(t *testing.T) {
	router := domain.NewCommandRouter()

	ctx := &domain.CommandContext{
		Command: "unknown",
		Args:    []string{},
		UserID:  12345,
		ChatID:  67890,
	}

	err := router.RouteCommand(ctx)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "command handler not found")
}

// Mock command handler for testing
type mockCommandHandler struct {
	called     bool
	handleFunc func(ctx *domain.CommandContext) error
}

func (m *mockCommandHandler) Handle(ctx *domain.CommandContext) error {
	m.called = true
	if m.handleFunc != nil {
		return m.handleFunc(ctx)
	}
	return nil
}
