package telegram

import (
	"context"
	"testing"

	"github.com/dewisartika8/cicd-status-notifier-bot/internal/core/bot/service"
	projectDomain "github.com/dewisartika8/cicd-status-notifier-bot/internal/core/project/domain"
	"github.com/dewisartika8/cicd-status-notifier-bot/internal/core/project/dto"
	"github.com/dewisartika8/cicd-status-notifier-bot/internal/core/project/port"
	"github.com/dewisartika8/cicd-status-notifier-bot/internal/core/shared/domain/value_objects"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockProjectService implements port.ProjectService for testing
type MockProjectService struct {
	mock.Mock
}

// GetActiveProjects mocks the GetActiveProjects method
func (m *MockProjectService) GetActiveProjects(ctx context.Context) ([]*projectDomain.Project, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*projectDomain.Project), args.Error(1)
}

// GetProjectByName mocks the GetProjectByName method
func (m *MockProjectService) GetProjectByName(ctx context.Context, name string) (*projectDomain.Project, error) {
	args := m.Called(ctx, name)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*projectDomain.Project), args.Error(1)
}

// CreateProject mocks the CreateProject method
func (m *MockProjectService) CreateProject(ctx context.Context, req dto.CreateProjectRequest) (*projectDomain.Project, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*projectDomain.Project), args.Error(1)
}

// GetProject mocks the GetProject method
func (m *MockProjectService) GetProject(ctx context.Context, id value_objects.ID) (*projectDomain.Project, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*projectDomain.Project), args.Error(1)
}

// GetProjectByRepositoryURL mocks the GetProjectByRepositoryURL method
func (m *MockProjectService) GetProjectByRepositoryURL(ctx context.Context, url string) (*projectDomain.Project, error) {
	args := m.Called(ctx, url)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*projectDomain.Project), args.Error(1)
}

// UpdateProject mocks the UpdateProject method
func (m *MockProjectService) UpdateProject(ctx context.Context, id value_objects.ID, req dto.UpdateProjectRequest) (*projectDomain.Project, error) {
	args := m.Called(ctx, id, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*projectDomain.Project), args.Error(1)
}

// UpdateProjectStatus mocks the UpdateProjectStatus method
func (m *MockProjectService) UpdateProjectStatus(ctx context.Context, id value_objects.ID, status projectDomain.ProjectStatus) (*projectDomain.Project, error) {
	args := m.Called(ctx, id, status)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*projectDomain.Project), args.Error(1)
}

// DeleteProject mocks the DeleteProject method
func (m *MockProjectService) DeleteProject(ctx context.Context, id value_objects.ID) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

// ListProjects mocks the ListProjects method
func (m *MockProjectService) ListProjects(ctx context.Context, filters dto.ListProjectFilters) ([]*projectDomain.Project, error) {
	args := m.Called(ctx, filters)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*projectDomain.Project), args.Error(1)
}

// GetProjectsWithTelegramChat mocks the GetProjectsWithTelegramChat method
func (m *MockProjectService) GetProjectsWithTelegramChat(ctx context.Context) ([]*projectDomain.Project, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*projectDomain.Project), args.Error(1)
}

// ActivateProject mocks the ActivateProject method
func (m *MockProjectService) ActivateProject(ctx context.Context, id value_objects.ID) (*projectDomain.Project, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*projectDomain.Project), args.Error(1)
}

// DeactivateProject mocks the DeactivateProject method
func (m *MockProjectService) DeactivateProject(ctx context.Context, id value_objects.ID) (*projectDomain.Project, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*projectDomain.Project), args.Error(1)
}

// ArchiveProject mocks the ArchiveProject method
func (m *MockProjectService) ArchiveProject(ctx context.Context, id value_objects.ID) (*projectDomain.Project, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*projectDomain.Project), args.Error(1)
}

// ValidateWebhookSecret mocks the ValidateWebhookSecret method
func (m *MockProjectService) ValidateWebhookSecret(ctx context.Context, id value_objects.ID, secret string) (bool, error) {
	args := m.Called(ctx, id, secret)
	return args.Bool(0), args.Error(1)
}

// CountProjects mocks the CountProjects method
func (m *MockProjectService) CountProjects(ctx context.Context, filters dto.ListProjectFilters) (int64, error) {
	args := m.Called(ctx, filters)
	return args.Get(0).(int64), args.Error(1)
}

// Ensure MockProjectService implements port.ProjectService
var _ port.ProjectService = (*MockProjectService)(nil)

func TestStatusCommandHandler_HandleStatusAllProjects(t *testing.T) {
	// Create mock service
	mockService := new(MockProjectService)
	handler := service.NewStatusCommandService(mockService)

	// Create test project
	chatID := int64(123456789)
	project, err := projectDomain.NewProject("test-project", "https://github.com/test/repo", "secret", &chatID)
	assert.NoError(t, err)

	projects := []*projectDomain.Project{project}

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
		mockService.On("GetActiveProjects", mock.Anything).Return([]*projectDomain.Project{}, nil).Once()

		response, err := handler.HandleStatusAllProjects()

		assert.NoError(t, err)
		assert.Contains(t, response, "ℹ️ No projects are currently being monitored")
		mockService.AssertExpectations(t)
	})

	// Test case 3: Service error
	t.Run("Service error", func(t *testing.T) {
		mockService.On("GetActiveProjects", mock.Anything).Return(nil, assert.AnError).Once()

		response, err := handler.HandleStatusAllProjects()

		assert.NoError(t, err)
		assert.Contains(t, response, "❌ **Error fetching project status**")
		mockService.AssertExpectations(t)
	})
}

func TestStatusCommandHandler_HandleStatusSpecificProject(t *testing.T) {
	// Create mock service
	mockService := new(MockProjectService)
	handler := service.NewStatusCommandService(mockService)

	// Create test project
	chatID := int64(123456789)
	project, err := projectDomain.NewProject("test-project", "https://github.com/test/repo", "secret", &chatID)
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
	handler := service.NewStatusCommandService(mockService)

	// Create test projects with different statuses
	chatID := int64(123456789)
	activeProject, err := projectDomain.NewProject("active-project", "https://github.com/test/active", "secret", &chatID)
	assert.NoError(t, err)

	projects := []*projectDomain.Project{activeProject}

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
		mockService.On("GetActiveProjects", mock.Anything).Return([]*projectDomain.Project{}, nil).Once()

		response, err := handler.HandleProjectsList()

		assert.NoError(t, err)
		assert.Contains(t, response, "ℹ️ No projects are currently being monitored")
		mockService.AssertExpectations(t)
	})

	// Test case 3: Service error
	t.Run("Service error", func(t *testing.T) {
		mockService.On("GetActiveProjects", mock.Anything).Return(nil, assert.AnError).Once()

		response, err := handler.HandleProjectsList()

		assert.NoError(t, err)
		assert.Contains(t, response, "❌ **Error fetching projects**")
		mockService.AssertExpectations(t)
	})
}
