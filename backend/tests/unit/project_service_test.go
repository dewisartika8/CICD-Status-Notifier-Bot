package repository_test

import (
	"context"
	"testing"

	"github.com/dewisartika8/cicd-status-notifier-bot/internal/core/project/domain"
	"github.com/dewisartika8/cicd-status-notifier-bot/internal/core/project/dto"
	"github.com/dewisartika8/cicd-status-notifier-bot/internal/core/project/port"
	"github.com/dewisartika8/cicd-status-notifier-bot/internal/core/project/service"
	"github.com/dewisartika8/cicd-status-notifier-bot/internal/core/shared/domain/value_objects"
	"github.com/dewisartika8/cicd-status-notifier-bot/pkg/exception"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

const (
	serviceTestProjectName   = "Test Project"
	serviceTestRepositoryURL = "https://github.com/test/repo"
	serviceTestWebhookSecret = "test-secret-123"
	domainProjectType        = "*domain.Project"
)

// MockProjectRepository implements the ProjectRepository interface for testing
type MockProjectRepository struct {
	mock.Mock
}

func (m *MockProjectRepository) Create(ctx context.Context, project *domain.Project) error {
	args := m.Called(ctx, project)
	return args.Error(0)
}

func (m *MockProjectRepository) GetByID(ctx context.Context, id value_objects.ID) (*domain.Project, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Project), args.Error(1)
}

func (m *MockProjectRepository) GetByName(ctx context.Context, name string) (*domain.Project, error) {
	args := m.Called(ctx, name)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Project), args.Error(1)
}

func (m *MockProjectRepository) GetByRepositoryURL(ctx context.Context, repositoryURL string) (*domain.Project, error) {
	args := m.Called(ctx, repositoryURL)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Project), args.Error(1)
}

func (m *MockProjectRepository) Update(ctx context.Context, project *domain.Project) error {
	args := m.Called(ctx, project)
	if err := args.Error(0); err != nil {
		return err
	}
	return nil
}

func (m *MockProjectRepository) Delete(ctx context.Context, id value_objects.ID) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockProjectRepository) List(ctx context.Context, filters dto.ListProjectFilters) ([]*domain.Project, error) {
	args := m.Called(ctx, filters)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*domain.Project), args.Error(1)
}

func (m *MockProjectRepository) Count(ctx context.Context, filters dto.ListProjectFilters) (int64, error) {
	args := m.Called(ctx, filters)
	return args.Get(0).(int64), args.Error(1)
}

func (m *MockProjectRepository) GetActiveProjects(ctx context.Context) ([]*domain.Project, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*domain.Project), args.Error(1)
}

func (m *MockProjectRepository) GetProjectsWithTelegramChat(ctx context.Context) ([]*domain.Project, error) {
	args := m.Called(ctx)
	if projects := args.Get(0); projects != nil {
		return projects.([]*domain.Project), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockProjectRepository) ExistsByName(ctx context.Context, name string) (bool, error) {
	args := m.Called(ctx, name)
	return args.Bool(0), args.Error(1)
}

func (m *MockProjectRepository) ExistsByRepositoryURL(ctx context.Context, repositoryURL string) (bool, error) {
	args := m.Called(ctx, repositoryURL)
	return args.Bool(0), args.Error(1)
}

type ProjectServiceTestSuite struct {
	suite.Suite
	mockRepo       *MockProjectRepository
	projectService port.ProjectService
	ctx            context.Context
}

func (suite *ProjectServiceTestSuite) SetupTest() {
	suite.mockRepo = &MockProjectRepository{}
	suite.projectService = service.NewProjectService(service.Dep{
		ProjectRepo: suite.mockRepo,
	})
	suite.ctx = context.Background()
}

func (suite *ProjectServiceTestSuite) TestCreateProjectSuccess() {
	// Test data
	req := dto.CreateProjectRequest{
		Name:          serviceTestProjectName,
		RepositoryURL: serviceTestRepositoryURL,
		WebhookSecret: serviceTestWebhookSecret,
	}

	// Setup mock expectations
	suite.mockRepo.On("ExistsByName", suite.ctx, req.Name).Return(false, nil)
	suite.mockRepo.On("ExistsByRepositoryURL", suite.ctx, req.RepositoryURL).Return(false, nil)
	suite.mockRepo.On("Create", suite.ctx, mock.AnythingOfType(domainProjectType)).Return(nil)

	// Execute
	project, err := suite.projectService.CreateProject(suite.ctx, req)

	// Assert
	suite.NoError(err)
	suite.NotNil(project)
	assert.Equal(suite.T(), req.Name, project.Name())
	assert.Equal(suite.T(), req.RepositoryURL, project.RepositoryURL())
	assert.Equal(suite.T(), domain.ProjectStatusActive, project.Status())

	// Verify mock calls
	suite.mockRepo.AssertExpectations(suite.T())
}

func (suite *ProjectServiceTestSuite) TestCreateProjectDuplicateName() {
	// Test data
	req := dto.CreateProjectRequest{
		Name:          "Existing Project",
		RepositoryURL: serviceTestRepositoryURL,
		WebhookSecret: serviceTestWebhookSecret,
	}

	// Setup mock expectations
	suite.mockRepo.On("ExistsByName", suite.ctx, req.Name).Return(true, nil)

	// Execute
	project, err := suite.projectService.CreateProject(suite.ctx, req)

	// Assert
	suite.Error(err)
	suite.Nil(project)

	// Verify mock calls
	suite.mockRepo.AssertExpectations(suite.T())
}

func (suite *ProjectServiceTestSuite) TestCreateProjectDuplicateRepositoryURL() {
	// Test data
	req := dto.CreateProjectRequest{
		Name:          serviceTestProjectName,
		RepositoryURL: "https://github.com/existing/repo",
		WebhookSecret: serviceTestWebhookSecret,
	}

	// Setup mock expectations
	suite.mockRepo.On("ExistsByName", suite.ctx, req.Name).Return(false, nil)
	suite.mockRepo.On("ExistsByRepositoryURL", suite.ctx, req.RepositoryURL).Return(true, nil)

	// Execute
	project, err := suite.projectService.CreateProject(suite.ctx, req)

	// Assert
	suite.Error(err)
	suite.Nil(project)

	// Verify mock calls
	suite.mockRepo.AssertExpectations(suite.T())
}

func (suite *ProjectServiceTestSuite) TestGetProjectSuccess() {
	// Test data
	projectID := value_objects.NewID()
	project, _ := domain.NewProject(serviceTestProjectName, serviceTestRepositoryURL, "secret", nil)

	// Setup mock expectations
	suite.mockRepo.On("GetByID", suite.ctx, projectID).Return(project, nil)

	// Execute
	result, err := suite.projectService.GetProject(suite.ctx, projectID)

	// Assert
	suite.NoError(err)
	suite.NotNil(result)
	assert.Equal(suite.T(), project.Name(), result.Name())
	assert.Equal(suite.T(), project.RepositoryURL(), result.RepositoryURL())

	// Verify mock calls
	suite.mockRepo.AssertExpectations(suite.T())
}

func (suite *ProjectServiceTestSuite) TestGetProjectNotFound() {
	// Test data
	projectID := value_objects.NewID()

	// Setup mock expectations
	suite.mockRepo.On("GetByID", suite.ctx, projectID).Return(nil, exception.NewDomainError("PROJECT_NOT_FOUND", "project not found"))

	// Execute
	result, err := suite.projectService.GetProject(suite.ctx, projectID)

	// Assert
	suite.Error(err)
	suite.Nil(result)

	// Verify mock calls
	suite.mockRepo.AssertExpectations(suite.T())
}

func (suite *ProjectServiceTestSuite) TestUpdateProjectSuccess() {
	// Test data
	projectID := value_objects.NewID()
	originalProject, _ := domain.NewProject("Original Project", serviceTestRepositoryURL, "secret", nil)
	newName := "Updated Project Name"
	newRepoURL := "https://github.com/test/updated-repo"

	req := dto.UpdateProjectRequest{
		Name:          &newName,
		RepositoryURL: &newRepoURL,
	}

	// Setup mock expectations
	suite.mockRepo.On("GetByID", suite.ctx, projectID).Return(originalProject, nil)
	suite.mockRepo.On("GetByName", suite.ctx, newName).Return(nil, nil)
	suite.mockRepo.On("GetByRepositoryURL", suite.ctx, newRepoURL).Return(nil, nil)
	suite.mockRepo.On("Update", suite.ctx, mock.AnythingOfType(domainProjectType)).Return(nil)

	// Execute
	result, err := suite.projectService.UpdateProject(suite.ctx, projectID, req)

	// Assert
	suite.NoError(err)
	suite.NotNil(result)
	assert.Equal(suite.T(), newName, result.Name())
	assert.Equal(suite.T(), newRepoURL, result.RepositoryURL())

	// Verify mock calls
	suite.mockRepo.AssertExpectations(suite.T())
}

func (suite *ProjectServiceTestSuite) TestDeleteProjectSuccess() {
	// Test data
	projectID := value_objects.NewID()
	project, _ := domain.NewProject(serviceTestProjectName, serviceTestRepositoryURL, "secret", nil)

	// Setup mock expectations
	suite.mockRepo.On("GetByID", suite.ctx, projectID).Return(project, nil)
	suite.mockRepo.On("Delete", suite.ctx, projectID).Return(nil)

	// Execute
	err := suite.projectService.DeleteProject(suite.ctx, projectID)

	// Assert
	suite.NoError(err)

	// Verify mock calls
	suite.mockRepo.AssertExpectations(suite.T())
}

func (suite *ProjectServiceTestSuite) TestListProjectsSuccess() {
	// Test data
	filters := dto.ListProjectFilters{}
	projects := []*domain.Project{}
	project1, _ := domain.NewProject("Project 1", "https://github.com/test/repo1", "secret1", nil)
	project2, _ := domain.NewProject("Project 2", "https://github.com/test/repo2", "secret2", nil)
	projects = append(projects, project1, project2)

	// Setup mock expectations
	suite.mockRepo.On("List", suite.ctx, filters).Return(projects, nil)

	// Execute
	result, err := suite.projectService.ListProjects(suite.ctx, filters)

	// Assert
	suite.NoError(err)
	suite.NotNil(result)
	assert.Len(suite.T(), result, 2)
	assert.Equal(suite.T(), "Project 1", result[0].Name())
	assert.Equal(suite.T(), "Project 2", result[1].Name())

	// Verify mock calls
	suite.mockRepo.AssertExpectations(suite.T())
}

func (suite *ProjectServiceTestSuite) TestCountProjectsSuccess() {
	// Test data
	filters := dto.ListProjectFilters{}
	expectedCount := int64(5)

	// Setup mock expectations
	suite.mockRepo.On("Count", suite.ctx, filters).Return(expectedCount, nil)

	// Execute
	result, err := suite.projectService.CountProjects(suite.ctx, filters)

	// Assert
	suite.NoError(err)
	assert.Equal(suite.T(), expectedCount, result)

	// Verify mock calls
	suite.mockRepo.AssertExpectations(suite.T())
}

func (suite *ProjectServiceTestSuite) TestUpdateProjectStatusSuccess() {
	// Test data
	projectID := value_objects.NewID()
	project, _ := domain.NewProject(serviceTestProjectName, serviceTestRepositoryURL, "secret", nil)
	newStatus := domain.ProjectStatusInactive

	// Setup mock expectations
	suite.mockRepo.On("GetByID", suite.ctx, projectID).Return(project, nil)
	suite.mockRepo.On("Update", suite.ctx, mock.AnythingOfType(domainProjectType)).Return(nil)

	// Execute
	result, err := suite.projectService.UpdateProjectStatus(suite.ctx, projectID, newStatus)

	// Assert
	suite.NoError(err)
	suite.NotNil(result)
	assert.Equal(suite.T(), newStatus, result.Status())

	// Verify mock calls
	suite.mockRepo.AssertExpectations(suite.T())
}

func TestProjectServiceTestSuite(t *testing.T) {
	suite.Run(t, new(ProjectServiceTestSuite))
}
