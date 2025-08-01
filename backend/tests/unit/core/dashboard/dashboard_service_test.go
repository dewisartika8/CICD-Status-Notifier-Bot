package dashboard_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	dashboardDomain "github.com/dewisartika8/cicd-status-notifier-bot/internal/core/dashboard/domain"
	"github.com/dewisartika8/cicd-status-notifier-bot/internal/core/dashboard/dto"
	dashboardService "github.com/dewisartika8/cicd-status-notifier-bot/internal/core/dashboard/service"
)

// Note: Using shared mock types from service_test.go in the same package
// to avoid redeclaration errors

// Test GetOverview
func TestGetOverview(t *testing.T) {
	tests := []struct {
		name           string
		cacheHit       bool
		cachedValue    interface{}
		mockMetrics    *dashboardDomain.OverviewMetrics
		mockError      error
		projectsCount  int
		activeProjects int
		expectedResult *dto.OverviewResponse
		expectError    bool
	}{
		{
			name:     "successful overview retrieval without cache",
			cacheHit: false,
			mockMetrics: &dashboardDomain.OverviewMetrics{
				TotalBuilds:      100,
				SuccessfulBuilds: 85,
				FailedBuilds:     15,
				AverageDuration:  300,
			},
			mockError:      nil,
			projectsCount:  5,
			activeProjects: 4,
			expectedResult: &dto.OverviewResponse{
				TotalProjects:    5,
				ActiveProjects:   4,
				TotalBuilds:      100,
				SuccessfulBuilds: 85,
				FailedBuilds:     15,
				SuccessRate:      85.0,
				AverageDuration:  300,
			},
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockBuildEventRepo := new(MockBuildEventRepository)
			mockProjectRepo := new(MockProjectRepository)
			mockCacheService := new(MockCacheService)

			service := dashboardService.NewService(mockBuildEventRepo, mockProjectRepo, mockCacheService)

			if tt.cacheHit {
				mockCacheService.On("Get", "dashboard_overview").Return(tt.cachedValue, true)
			} else {
				mockCacheService.On("Get", "dashboard_overview").Return(nil, false)
				mockBuildEventRepo.On("GetOverviewMetrics").Return(tt.mockMetrics, tt.mockError)
				mockProjectRepo.On("GetTotalProjectsCount").Return(tt.projectsCount, nil)
				mockProjectRepo.On("GetActiveProjectsCount").Return(tt.activeProjects, nil)
				mockCacheService.On("Set", "dashboard_overview", mock.Anything, mock.Anything).Return()
			}

			result, err := service.GetOverview()

			if tt.expectError {
				assert.Error(t, err)
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, result)
				assert.Equal(t, tt.expectedResult.TotalProjects, result.TotalProjects)
				assert.Equal(t, tt.expectedResult.ActiveProjects, result.ActiveProjects)
				assert.Equal(t, tt.expectedResult.TotalBuilds, result.TotalBuilds)
				assert.Equal(t, tt.expectedResult.SuccessfulBuilds, result.SuccessfulBuilds)
				assert.Equal(t, tt.expectedResult.FailedBuilds, result.FailedBuilds)
				assert.Equal(t, tt.expectedResult.SuccessRate, result.SuccessRate)
				assert.Equal(t, tt.expectedResult.AverageDuration, result.AverageDuration)
			}

			mockBuildEventRepo.AssertExpectations(t)
			mockProjectRepo.AssertExpectations(t)
			mockCacheService.AssertExpectations(t)
		})
	}
}
