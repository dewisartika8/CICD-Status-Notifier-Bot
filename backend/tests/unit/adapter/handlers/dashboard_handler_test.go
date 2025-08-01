package dashboard

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/dewisartika8/cicd-status-notifier-bot/internal/adapter/handler/dashboard"
	"github.com/dewisartika8/cicd-status-notifier-bot/internal/core/dashboard/dto"
	"github.com/dewisartika8/cicd-status-notifier-bot/internal/core/dashboard/port"
	"github.com/dewisartika8/cicd-status-notifier-bot/internal/core/shared/domain/value_objects"
)

const (
	// Test constants for error messages
	ServiceErrorMsg         = "service error"
	InvalidUUID             = "invalid-uuid"
	ProjectNotFoundMsg      = "project not found"
	AnalyticsCalculationMsg = "analytics calculation error"

	// Test constants for success messages
	OverviewSuccessMsg  = "Dashboard overview retrieved successfully"
	AnalyticsSuccessMsg = "Build analytics retrieved successfully"

	// Test constants for test names
	SuccessfulOverviewTest         = "successful overview retrieval"
	ServiceErrorTest               = "service error"
	InvalidProjectIDTest           = "invalid project ID"
	SuccessfulProjectStatsTest     = "successful project statistics retrieval"
	SuccessfulAnalyticsDefaultTest = "successful build analytics retrieval with default range"
	SuccessfulAnalytics30dTest     = "successful build analytics retrieval with 30d range"

	// Test constants for time ranges
	DefaultTimeRange  = "7d"
	ExtendedTimeRange = "30d"

	// Test constants for project names
	TestProjectName = "Test Project"

	// HTTP status codes
	StatusOK                  = 200
	StatusBadRequest          = 400
	StatusNotFound            = 404
	StatusInternalServerError = 500

	// Test data constants
	TotalProjectsCount     = 5
	ActiveProjectsCount    = 4
	TotalBuildsCount       = 100
	SuccessfulBuildsCount  = 85
	FailedBuildsCount      = 15
	SuccessRatePercentage  = 85.0
	AverageDurationSeconds = 300

	ProjectTotalBuilds      = 50
	ProjectSuccessfulBuilds = 42
	ProjectFailedBuilds     = 8
	ProjectSuccessRate      = 84.0
	ProjectAverageDuration  = 250

	AnalyticsTotalBuilds      = 200
	AnalyticsSuccessfulBuilds = 170
	AnalyticsFailedBuilds     = 30
	AnalyticsSuccessRate      = 85.0
	AnalyticsAverageDuration  = 300

	ExtendedTotalBuilds      = 500
	ExtendedSuccessfulBuilds = 425
	ExtendedFailedBuilds     = 75
	ExtendedAverageDuration  = 310

	// Build trend test data
	BuildTrendCount1        = 5
	BuildTrendSuccessCount1 = 4
	BuildTrendCount2        = 3
	BuildTrendSuccessCount2 = 3

	// Analytics daily data
	DailyBuilds1        = 30
	DailySuccessBuilds1 = 25
	DailyFailedBuilds1  = 5
	DailyBuilds2        = 28
	DailySuccessBuilds2 = 24
	DailyFailedBuilds2  = 4
	DailyDuration1      = 280
	DailyDuration2      = 320
)

// MockDashboardService is a mock implementation of port.DashboardService
type MockDashboardService struct {
	mock.Mock
}

// Ensure MockDashboardService implements port.DashboardService
var _ port.DashboardService = (*MockDashboardService)(nil)

func (m *MockDashboardService) GetOverview() (*dto.OverviewResponse, error) {
	args := m.Called()
	return args.Get(0).(*dto.OverviewResponse), args.Error(1)
}

func (m *MockDashboardService) GetProjectStatistics(projectID value_objects.ID) (*dto.ProjectStatisticsResponse, error) {
	args := m.Called(projectID)
	return args.Get(0).(*dto.ProjectStatisticsResponse), args.Error(1)
}

func (m *MockDashboardService) GetBuildAnalytics(timeRange string) (*dto.BuildAnalyticsResponse, error) {
	args := m.Called(timeRange)
	return args.Get(0).(*dto.BuildAnalyticsResponse), args.Error(1)
}

func setupTestApp() (*fiber.App, *MockDashboardService) {
	app := fiber.New()
	mockService := &MockDashboardService{}
	handler := dashboard.NewHandler(mockService)

	// Setup API routes like in main application
	api := app.Group("/api/v1")
	handler.RegisterRoutes(api)

	return app, mockService
}

func TestDashboardHandlerGetOverview(t *testing.T) {
	tests := []struct {
		name         string
		mockResponse *dto.OverviewResponse
		mockError    error
		expectedCode int
		expectError  bool
	}{
		{
			name: SuccessfulOverviewTest,
			mockResponse: &dto.OverviewResponse{
				TotalProjects:    TotalProjectsCount,
				ActiveProjects:   ActiveProjectsCount,
				TotalBuilds:      TotalBuildsCount,
				SuccessfulBuilds: SuccessfulBuildsCount,
				FailedBuilds:     FailedBuildsCount,
				SuccessRate:      SuccessRatePercentage,
				AverageDuration:  AverageDurationSeconds,
				LastUpdated:      time.Now(),
			},
			mockError:    nil,
			expectedCode: StatusOK,
			expectError:  false,
		},
		{
			name:         ServiceErrorTest,
			mockResponse: nil,
			mockError:    fmt.Errorf(ServiceErrorMsg),
			expectedCode: StatusInternalServerError,
			expectError:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			app, mockService := setupTestApp()

			if tt.mockError != nil {
				mockService.On("GetOverview").Return((*dto.OverviewResponse)(nil), tt.mockError)
			} else {
				mockService.On("GetOverview").Return(tt.mockResponse, nil)
			}

			req := httptest.NewRequest("GET", "/api/v1/dashboard/overview", nil)
			resp, err := app.Test(req)

			assert.NoError(t, err)
			assert.Equal(t, tt.expectedCode, resp.StatusCode)

			if !tt.expectError {
				body, err := io.ReadAll(resp.Body)
				assert.NoError(t, err)

				var response map[string]interface{}
				err = json.Unmarshal(body, &response)
				assert.NoError(t, err)

				assert.Equal(t, OverviewSuccessMsg, response["message"])
				assert.NotNil(t, response["data"])
			}

			mockService.AssertExpectations(t)
		})
	}
}

func TestDashboardHandlerGetProjectStatistics(t *testing.T) {
	projectID := value_objects.NewID()

	tests := []struct {
		name         string
		projectID    string
		mockResponse *dto.ProjectStatisticsResponse
		mockError    error
		expectedCode int
		expectError  bool
	}{
		{
			name:      SuccessfulProjectStatsTest,
			projectID: projectID.String(),
			mockResponse: &dto.ProjectStatisticsResponse{
				ProjectID:        projectID,
				ProjectName:      TestProjectName,
				TotalBuilds:      ProjectTotalBuilds,
				SuccessfulBuilds: ProjectSuccessfulBuilds,
				FailedBuilds:     ProjectFailedBuilds,
				SuccessRate:      ProjectSuccessRate,
				AverageDuration:  ProjectAverageDuration,
				LastBuildTime:    time.Now(),
				BuildTrends: []dto.BuildTrendData{
					{Date: time.Now().AddDate(0, 0, -1), Count: BuildTrendCount1, SuccessCount: BuildTrendSuccessCount1},
					{Date: time.Now(), Count: BuildTrendCount2, SuccessCount: BuildTrendSuccessCount2},
				},
			},
			mockError:    nil,
			expectedCode: StatusOK,
			expectError:  false,
		},
		{
			name:         InvalidProjectIDTest,
			projectID:    InvalidUUID,
			mockResponse: nil,
			mockError:    nil,
			expectedCode: StatusBadRequest,
			expectError:  true,
		},
		{
			name:         ServiceErrorTest,
			projectID:    projectID.String(),
			mockResponse: nil,
			mockError:    fmt.Errorf(ProjectNotFoundMsg),
			expectedCode: StatusNotFound,
			expectError:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			app, mockService := setupTestApp()

			if tt.mockError != nil && tt.projectID != InvalidUUID {
				parsedID, _ := value_objects.NewIDFromString(tt.projectID)
				mockService.On("GetProjectStatistics", parsedID).Return((*dto.ProjectStatisticsResponse)(nil), tt.mockError)
			} else if tt.mockError == nil && tt.projectID != InvalidUUID {
				parsedID, _ := value_objects.NewIDFromString(tt.projectID)
				mockService.On("GetProjectStatistics", parsedID).Return(tt.mockResponse, nil)
			}

			url := fmt.Sprintf("/api/v1/projects/%s/statistics", tt.projectID)
			req := httptest.NewRequest("GET", url, nil)
			resp, err := app.Test(req)

			assert.NoError(t, err)
			assert.Equal(t, tt.expectedCode, resp.StatusCode)

			if tt.projectID != InvalidUUID {
				mockService.AssertExpectations(t)
			}
		})
	}
}

func TestDashboardHandlerGetBuildAnalytics(t *testing.T) {
	tests := []struct {
		name         string
		timeRange    string
		mockResponse *dto.BuildAnalyticsResponse
		mockError    error
		expectedCode int
		expectError  bool
	}{
		{
			name:      SuccessfulAnalyticsDefaultTest,
			timeRange: "",
			mockResponse: &dto.BuildAnalyticsResponse{
				TimeRange:        DefaultTimeRange,
				TotalBuilds:      AnalyticsTotalBuilds,
				SuccessfulBuilds: AnalyticsSuccessfulBuilds,
				FailedBuilds:     AnalyticsFailedBuilds,
				SuccessRate:      AnalyticsSuccessRate,
				AverageDuration:  AnalyticsAverageDuration,
				BuildsByDay: []dto.BuildAnalyticsData{
					{Date: time.Now().AddDate(0, 0, -6), TotalBuilds: DailyBuilds1, SuccessfulBuilds: DailySuccessBuilds1, FailedBuilds: DailyFailedBuilds1},
					{Date: time.Now().AddDate(0, 0, -5), TotalBuilds: DailyBuilds2, SuccessfulBuilds: DailySuccessBuilds2, FailedBuilds: DailyFailedBuilds2},
				},
				DurationTrends: []dto.DurationTrendData{
					{Date: time.Now().AddDate(0, 0, -6), AverageDuration: DailyDuration1},
					{Date: time.Now().AddDate(0, 0, -5), AverageDuration: DailyDuration2},
				},
			},
			mockError:    nil,
			expectedCode: StatusOK,
			expectError:  false,
		},
		{
			name:      SuccessfulAnalytics30dTest,
			timeRange: ExtendedTimeRange,
			mockResponse: &dto.BuildAnalyticsResponse{
				TimeRange:        ExtendedTimeRange,
				TotalBuilds:      ExtendedTotalBuilds,
				SuccessfulBuilds: ExtendedSuccessfulBuilds,
				FailedBuilds:     ExtendedFailedBuilds,
				SuccessRate:      AnalyticsSuccessRate,
				AverageDuration:  ExtendedAverageDuration,
				BuildsByDay:      []dto.BuildAnalyticsData{},
				DurationTrends:   []dto.DurationTrendData{},
			},
			mockError:    nil,
			expectedCode: StatusOK,
			expectError:  false,
		},
		{
			name:         ServiceErrorTest,
			timeRange:    DefaultTimeRange,
			mockResponse: nil,
			mockError:    fmt.Errorf(AnalyticsCalculationMsg),
			expectedCode: StatusInternalServerError,
			expectError:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			app, mockService := setupTestApp()

			expectedRange := tt.timeRange
			if expectedRange == "" {
				expectedRange = DefaultTimeRange
			}

			if tt.mockError != nil {
				mockService.On("GetBuildAnalytics", expectedRange).Return((*dto.BuildAnalyticsResponse)(nil), tt.mockError)
			} else {
				mockService.On("GetBuildAnalytics", expectedRange).Return(tt.mockResponse, nil)
			}

			url := "/api/v1/builds/analytics"
			if tt.timeRange != "" {
				url += "?range=" + tt.timeRange
			}

			req := httptest.NewRequest("GET", url, nil)
			resp, err := app.Test(req)

			assert.NoError(t, err)
			assert.Equal(t, tt.expectedCode, resp.StatusCode)

			if !tt.expectError {
				body, err := io.ReadAll(resp.Body)
				assert.NoError(t, err)

				var response map[string]interface{}
				err = json.Unmarshal(body, &response)
				assert.NoError(t, err)

				assert.Equal(t, AnalyticsSuccessMsg, response["message"])
				assert.NotNil(t, response["data"])
			}

			mockService.AssertExpectations(t)
		})
	}
}
