package dashboard_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	dashboardDomain "github.com/dewisartika8/cicd-status-notifier-bot/internal/core/dashboard/domain"
	dashboardService "github.com/dewisartika8/cicd-status-notifier-bot/internal/core/dashboard/service"
	"github.com/dewisartika8/cicd-status-notifier-bot/internal/core/shared/domain/value_objects"
)

// MockBuildEventRepository is a mock implementation of BuildEventRepositoryInterface
type MockBuildEventRepository struct {
	mock.Mock
}

func (m *MockBuildEventRepository) GetOverviewMetrics() (*dashboardDomain.OverviewMetrics, error) {
	args := m.Called()
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*dashboardDomain.OverviewMetrics), args.Error(1)
}

func (m *MockBuildEventRepository) GetProjectStatistics(projectID value_objects.ID) (*dashboardDomain.ProjectStatistics, error) {
	args := m.Called(projectID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*dashboardDomain.ProjectStatistics), args.Error(1)
}

func (m *MockBuildEventRepository) GetBuildAnalytics(timeRange string) (*dashboardDomain.BuildAnalytics, error) {
	args := m.Called(timeRange)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*dashboardDomain.BuildAnalytics), args.Error(1)
}

// MockProjectRepository is a mock implementation of ProjectRepositoryInterface
type MockProjectRepository struct {
	mock.Mock
}

func (m *MockProjectRepository) GetProjectNameByID(projectID value_objects.ID) (string, error) {
	args := m.Called(projectID)
	return args.String(0), args.Error(1)
}

func (m *MockProjectRepository) GetActiveProjectsCount() (int, error) {
	args := m.Called("GetActiveProjectsCount")
	return args.Int(0), args.Error(1)
}

func (m *MockProjectRepository) GetTotalProjectsCount() (int, error) {
	args := m.Called("GetTotalProjectsCount")
	return args.Int(0), args.Error(1)
}

// MockCacheService is a mock implementation of CacheServiceInterface
type MockCacheService struct {
	mock.Mock
}

func (m *MockCacheService) Get(key string) (interface{}, bool) {
	args := m.Called(key)
	return args.Get(0), args.Bool(1)
}

func (m *MockCacheService) Set(key string, value interface{}, ttl time.Duration) {
	m.Called(key, value, ttl)
}

func (m *MockCacheService) Delete(key string) {
	m.Called(key)
}

func TestDashboardServiceGetOverview(t *testing.T) {
	t.Run("successful overview retrieval", func(t *testing.T) {
		mockBuildRepo := &MockBuildEventRepository{}
		mockProjectRepo := &MockProjectRepository{}
		mockCache := &MockCacheService{}

		service := dashboardService.NewService(mockBuildRepo, mockProjectRepo, mockCache)

		// Setup mocks
		mockCache.On("Get", "dashboard_overview").Return(nil, false)
		mockBuildRepo.On("GetOverviewMetrics").Return(&dashboardDomain.OverviewMetrics{
			TotalBuilds:      100,
			SuccessfulBuilds: 85,
			FailedBuilds:     15,
			AverageDuration:  300,
		}, nil)
		mockProjectRepo.On("GetTotalProjectsCount", "GetTotalProjectsCount").Return(5, nil)
		mockProjectRepo.On("GetActiveProjectsCount", "GetActiveProjectsCount").Return(4, nil)
		mockCache.On("Set", "dashboard_overview", mock.Anything, 5*time.Minute).Return()

		result, err := service.GetOverview()

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, 5, result.TotalProjects)
		assert.Equal(t, 4, result.ActiveProjects)
		assert.Equal(t, 100, result.TotalBuilds)
		assert.Equal(t, 85, result.SuccessfulBuilds)
		assert.Equal(t, 15, result.FailedBuilds)
		assert.InDelta(t, 85.0, result.SuccessRate, 0.01)
		assert.Equal(t, 300, result.AverageDuration)

		mockBuildRepo.AssertExpectations(t)
		mockProjectRepo.AssertExpectations(t)
		mockCache.AssertExpectations(t)
	})

	t.Run("repository error", func(t *testing.T) {
		mockBuildRepo := &MockBuildEventRepository{}
		mockProjectRepo := &MockProjectRepository{}
		mockCache := &MockCacheService{}

		service := dashboardService.NewService(mockBuildRepo, mockProjectRepo, mockCache)

		mockCache.On("Get", "dashboard_overview").Return(nil, false)
		mockBuildRepo.On("GetOverviewMetrics").Return(nil, fmt.Errorf("database error"))

		result, err := service.GetOverview()

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "failed to get overview metrics")

		mockBuildRepo.AssertExpectations(t)
		mockProjectRepo.AssertExpectations(t)
		mockCache.AssertExpectations(t)
	})
}

func TestDashboardServiceGetProjectStatistics(t *testing.T) {
	projectID := value_objects.NewID()

	t.Run("successful project statistics retrieval", func(t *testing.T) {
		mockBuildRepo := &MockBuildEventRepository{}
		mockProjectRepo := &MockProjectRepository{}
		mockCache := &MockCacheService{}

		service := dashboardService.NewService(mockBuildRepo, mockProjectRepo, mockCache)

		cacheKey := fmt.Sprintf("project_statistics_%s", projectID.String())
		mockCache.On("Get", cacheKey).Return(nil, false)
		mockBuildRepo.On("GetProjectStatistics", projectID).Return(&dashboardDomain.ProjectStatistics{
			TotalBuilds:      50,
			SuccessfulBuilds: 42,
			FailedBuilds:     8,
			AverageDuration:  250,
			LastBuildTime:    time.Now(),
			BuildTrends: []dashboardDomain.BuildTrendInfo{
				{Date: time.Now().AddDate(0, 0, -1), Count: 5, SuccessCount: 4},
			},
		}, nil)
		mockProjectRepo.On("GetProjectNameByID", projectID).Return("Test Project", nil)
		mockCache.On("Set", cacheKey, mock.Anything, 2*time.Minute).Return()

		result, err := service.GetProjectStatistics(projectID)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, projectID, result.ProjectID)
		assert.Equal(t, "Test Project", result.ProjectName)
		assert.Equal(t, 50, result.TotalBuilds)
		assert.Equal(t, 42, result.SuccessfulBuilds)
		assert.Equal(t, 8, result.FailedBuilds)
		assert.InDelta(t, 84.0, result.SuccessRate, 0.01)
		assert.Equal(t, 250, result.AverageDuration)
		assert.Equal(t, 1, len(result.BuildTrends))

		mockBuildRepo.AssertExpectations(t)
		mockProjectRepo.AssertExpectations(t)
		mockCache.AssertExpectations(t)
	})
}

func TestDashboardServiceGetBuildAnalytics(t *testing.T) {
	t.Run("successful build analytics retrieval", func(t *testing.T) {
		mockBuildRepo := &MockBuildEventRepository{}
		mockProjectRepo := &MockProjectRepository{}
		mockCache := &MockCacheService{}

		service := dashboardService.NewService(mockBuildRepo, mockProjectRepo, mockCache)

		timeRange := "7d"
		cacheKey := fmt.Sprintf("build_analytics_%s", timeRange)
		mockCache.On("Get", cacheKey).Return(nil, false)
		mockBuildRepo.On("GetBuildAnalytics", timeRange).Return(&dashboardDomain.BuildAnalytics{
			TotalBuilds:      200,
			SuccessfulBuilds: 170,
			FailedBuilds:     30,
			AverageDuration:  300,
			BuildsByDay: []dashboardDomain.BuildAnalyticsInfo{
				{Date: time.Now().AddDate(0, 0, -6), TotalBuilds: 30, SuccessfulBuilds: 25, FailedBuilds: 5},
			},
			DurationTrends: []dashboardDomain.DurationTrendInfo{
				{Date: time.Now().AddDate(0, 0, -6), AverageDuration: 280},
			},
		}, nil)
		mockCache.On("Set", cacheKey, mock.Anything, 10*time.Minute).Return()

		result, err := service.GetBuildAnalytics(timeRange)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, timeRange, result.TimeRange)
		assert.Equal(t, 200, result.TotalBuilds)
		assert.Equal(t, 170, result.SuccessfulBuilds)
		assert.Equal(t, 30, result.FailedBuilds)
		assert.InDelta(t, 85.0, result.SuccessRate, 0.01)
		assert.Equal(t, 300, result.AverageDuration)
		assert.Equal(t, 1, len(result.BuildsByDay))
		assert.Equal(t, 1, len(result.DurationTrends))

		mockBuildRepo.AssertExpectations(t)
		mockProjectRepo.AssertExpectations(t)
		mockCache.AssertExpectations(t)
	})
}
