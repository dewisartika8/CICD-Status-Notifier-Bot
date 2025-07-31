package mocks

import (
	"context"

	"github.com/stretchr/testify/mock"

	"github.com/dewisartika8/cicd-status-notifier-bot/internal/core/project/domain"
	"github.com/dewisartika8/cicd-status-notifier-bot/internal/core/project/dto"
	"github.com/dewisartika8/cicd-status-notifier-bot/internal/core/shared/domain/value_objects"
)

// MockProjectService is a comprehensive mock implementation of port.ProjectService
// This mock can be used across all test packages to ensure consistency
type MockProjectService struct {
	mock.Mock
}

// Core project operations
func (m *MockProjectService) CreateProject(ctx context.Context, req dto.CreateProjectRequest) (*domain.Project, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Project), args.Error(1)
}

func (m *MockProjectService) GetProject(ctx context.Context, id value_objects.ID) (*domain.Project, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Project), args.Error(1)
}

func (m *MockProjectService) GetProjectByName(ctx context.Context, name string) (*domain.Project, error) {
	args := m.Called(ctx, name)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Project), args.Error(1)
}

func (m *MockProjectService) GetProjectByRepositoryURL(ctx context.Context, repositoryURL string) (*domain.Project, error) {
	args := m.Called(ctx, repositoryURL)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Project), args.Error(1)
}

func (m *MockProjectService) UpdateProject(ctx context.Context, id value_objects.ID, req dto.UpdateProjectRequest) (*domain.Project, error) {
	args := m.Called(ctx, id, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Project), args.Error(1)
}

func (m *MockProjectService) UpdateProjectStatus(ctx context.Context, id value_objects.ID, status domain.ProjectStatus) (*domain.Project, error) {
	args := m.Called(ctx, id, status)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Project), args.Error(1)
}

func (m *MockProjectService) DeleteProject(ctx context.Context, id value_objects.ID) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

// Project listing and filtering
func (m *MockProjectService) ListProjects(ctx context.Context, filters dto.ListProjectFilters) ([]*domain.Project, error) {
	args := m.Called(ctx, filters)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*domain.Project), args.Error(1)
}

func (m *MockProjectService) GetActiveProjects(ctx context.Context) ([]*domain.Project, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*domain.Project), args.Error(1)
}

func (m *MockProjectService) GetProjectsWithTelegramChat(ctx context.Context) ([]*domain.Project, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*domain.Project), args.Error(1)
}

// Project lifecycle operations
func (m *MockProjectService) ActivateProject(ctx context.Context, id value_objects.ID) (*domain.Project, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Project), args.Error(1)
}

func (m *MockProjectService) DeactivateProject(ctx context.Context, id value_objects.ID) (*domain.Project, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Project), args.Error(1)
}

func (m *MockProjectService) ArchiveProject(ctx context.Context, id value_objects.ID) (*domain.Project, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Project), args.Error(1)
}

// Additional utility methods
func (m *MockProjectService) ValidateWebhookSecret(ctx context.Context, id value_objects.ID, secret string) (bool, error) {
	args := m.Called(ctx, id, secret)
	return args.Bool(0), args.Error(1)
}

func (m *MockProjectService) CountProjects(ctx context.Context, filters dto.ListProjectFilters) (int64, error) {
	args := m.Called(ctx, filters)
	return args.Get(0).(int64), args.Error(1)
}
