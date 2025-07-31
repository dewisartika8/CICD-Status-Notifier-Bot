package telegram

import (
	"context"
	"testing"

	"github.com/dewisartika8/cicd-status-notifier-bot/internal/core/project/domain"
	"github.com/dewisartika8/cicd-status-notifier-bot/internal/core/project/dto"
	"github.com/dewisartika8/cicd-status-notifier-bot/internal/core/shared/domain/value_objects"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockProjectService for testing
type MockProjectService struct {
	mock.Mock
}

func (m *MockProjectService) GetActiveProjects(ctx context.Context) ([]*domain.Project, error) {
	args := m.Called(ctx)
	return args.Get(0).([]*domain.Project), args.Error(1)
}

func (m *MockProjectService) GetProjectByName(ctx context.Context, name string) (*domain.Project, error) {
	args := m.Called(ctx, name)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Project), args.Error(1)
}

// Implement other required methods as no-ops for testing
func (m *MockProjectService) CreateProject(ctx context.Context, req dto.CreateProjectRequest) (*domain.Project, error) {
	args := m.Called(ctx, req)
	return args.Get(0).(*domain.Project), args.Error(1)
}
func (m *MockProjectService) GetProject(ctx context.Context, id value_objects.ID) (*domain.Project, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(*domain.Project), args.Error(1)
}
func (m *MockProjectService) GetProjectByRepositoryURL(ctx context.Context, url string) (*domain.Project, error) {
	args := m.Called(ctx, url)
	return args.Get(0).(*domain.Project), args.Error(1)
}
func (m *MockProjectService) UpdateProject(ctx context.Context, id value_objects.ID, req dto.UpdateProjectRequest) (*domain.Project, error) {
	args := m.Called(ctx, id, req)
	return args.Get(0).(*domain.Project), args.Error(1)
}
func (m *MockProjectService) UpdateProjectStatus(ctx context.Context, id value_objects.ID, status domain.ProjectStatus) (*domain.Project, error) {
	args := m.Called(ctx, id, status)
	return args.Get(0).(*domain.Project), args.Error(1)
}
func (m *MockProjectService) DeleteProject(ctx context.Context, id value_objects.ID) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}
func (m *MockProjectService) ListProjects(ctx context.Context, filters dto.ListProjectFilters) ([]*domain.Project, error) {
	args := m.Called(ctx, filters)
	return args.Get(0).([]*domain.Project), args.Error(1)
}
func (m *MockProjectService) GetProjectsWithTelegramChat(ctx context.Context) ([]*domain.Project, error) {
	args := m.Called(ctx)
	return args.Get(0).([]*domain.Project), args.Error(1)
}
func (m *MockProjectService) ActivateProject(ctx context.Context, id value_objects.ID) (*domain.Project, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(*domain.Project), args.Error(1)
}
func (m *MockProjectService) DeactivateProject(ctx context.Context, id value_objects.ID) (*domain.Project, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(*domain.Project), args.Error(1)
}
func (m *MockProjectService) ArchiveProject(ctx context.Context, id value_objects.ID) (*domain.Project, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(*domain.Project), args.Error(1)
}

func (m *MockProjectService) ValidateWebhookSecret(ctx context.Context, id value_objects.ID, secret string) (bool, error) {
	args := m.Called(ctx, id, secret)
	return args.Bool(0), args.Error(1)
}

func (m *MockProjectService) CountProjects(ctx context.Context, filters dto.ListProjectFilters) (int64, error) {
	args := m.Called(ctx, filters)
	return args.Get(0).(int64), args.Error(1)
}

func TestStatusCommandHandler_HandleStatusAllProjects(t *testing.T) {
	// Create mock service
	mockService := new(MockProjectService)
	handler := NewStatusCommandHandler(mockService)

	// Create test project
	chatID := int64(123456789)
	project, err := domain.NewProject("test-project", "https://github.com/test/repo", "secret", &chatID)
	assert.NoError(t, err)

	projects := []*domain.Project{project}

	// Test case 1: Successful response with projects
	t.Run("Success with projects", func(t *testing.T) {
		mockService.On("GetActiveProjects", mock.Anything).Return(projects, nil).Once()

		response, err := handler.HandleStatusAllProjects()

		assert.NoError(t, err)
		assert.Contains(t, response, "📊 **Overall Project Status**")
		assert.Contains(t, response, "test-project")
		assert.Contains(t, response, "✅")
		assert.Contains(t, response, "Active")
		mockService.AssertExpectations(t)
	})

	// Test case 2: No projects
	t.Run("No projects", func(t *testing.T) {
		mockService.On("GetActiveProjects", mock.Anything).Return([]*domain.Project{}, nil).Once()

		response, err := handler.HandleStatusAllProjects()

		assert.NoError(t, err)
		assert.Contains(t, response, "ℹ️ No projects are currently being monitored")
		mockService.AssertExpectations(t)
	})
}

func TestStatusCommandHandler_HandleStatusSpecificProject(t *testing.T) {
	// Create mock service
	mockService := new(MockProjectService)
	handler := NewStatusCommandHandler(mockService)

	// Create test project
	chatID := int64(123456789)
	project, err := domain.NewProject("test-project", "https://github.com/test/repo", "secret", &chatID)
	assert.NoError(t, err)

	// Test case 1: Project found
	t.Run("Project found", func(t *testing.T) {
		mockService.On("GetProjectByName", mock.Anything, "test-project").Return(project, nil).Once()

		response, err := handler.HandleStatusSpecificProject("test-project")

		assert.NoError(t, err)
		assert.Contains(t, response, "📊 **Project Status: test-project**")
		assert.Contains(t, response, "✅")
		assert.Contains(t, response, "Active")
		assert.Contains(t, response, "https://github.com/test/repo")
		mockService.AssertExpectations(t)
	})

	// Test case 2: Project not found
	t.Run("Project not found", func(t *testing.T) {
		mockService.On("GetProjectByName", mock.Anything, "unknown-project").Return(nil, assert.AnError).Once()

		response, err := handler.HandleStatusSpecificProject("unknown-project")

		assert.NoError(t, err)
		assert.Contains(t, response, "❌ **Project not found**")
		assert.Contains(t, response, "unknown-project")
		mockService.AssertExpectations(t)
	})

	// Test case 3: Empty project name
	t.Run("Empty project name", func(t *testing.T) {
		response, err := handler.HandleStatusSpecificProject("")

		assert.NoError(t, err)
		assert.Contains(t, response, "❌ **Invalid command**")
		assert.Contains(t, response, "Please specify a project name")
	})
}

func TestStatusCommandHandler_HandleProjectsList(t *testing.T) {
	// Create mock service
	mockService := new(MockProjectService)
	handler := NewStatusCommandHandler(mockService)

	// Create test projects with different statuses
	chatID := int64(123456789)
	activeProject, _ := domain.NewProject("active-project", "https://github.com/test/active", "secret", &chatID)

	projects := []*domain.Project{activeProject}

	// Test case 1: Projects list success
	t.Run("Projects list success", func(t *testing.T) {
		mockService.On("GetActiveProjects", mock.Anything).Return(projects, nil).Once()

		response, err := handler.HandleProjectsList()

		assert.NoError(t, err)
		assert.Contains(t, response, "📋 **Monitored Projects**")
		assert.Contains(t, response, "active-project")
		assert.Contains(t, response, "Quick Commands")
		mockService.AssertExpectations(t)
	})

	// Test case 2: No projects
	t.Run("No projects", func(t *testing.T) {
		mockService.On("GetActiveProjects", mock.Anything).Return([]*domain.Project{}, nil).Once()

		response, err := handler.HandleProjectsList()

		assert.NoError(t, err)
		assert.Contains(t, response, "ℹ️ No projects are currently being monitored")
		mockService.AssertExpectations(t)
	})
}
