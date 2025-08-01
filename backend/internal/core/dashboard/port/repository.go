package port

import (
	"time"

	"github.com/dewisartika8/cicd-status-notifier-bot/internal/core/dashboard/domain"
	"github.com/dewisartika8/cicd-status-notifier-bot/internal/core/dashboard/dto"
	"github.com/dewisartika8/cicd-status-notifier-bot/internal/core/shared/domain/value_objects"
)

// BuildEventRepositoryInterface defines methods for build event data access
type BuildEventRepositoryInterface interface {
	// GetOverviewMetrics returns overall build metrics
	GetOverviewMetrics() (*domain.OverviewMetrics, error)

	// GetProjectStatistics returns statistics for a specific project
	GetProjectStatistics(projectID value_objects.ID) (*domain.ProjectStatistics, error)

	// GetBuildAnalytics returns build analytics for a given time range
	GetBuildAnalytics(timeRange string) (*domain.BuildAnalytics, error)
}

// ProjectRepositoryInterface defines methods for project data access
type ProjectRepositoryInterface interface {
	// GetProjectNameByID returns project name by ID
	GetProjectNameByID(projectID value_objects.ID) (string, error)

	// GetActiveProjectsCount returns count of active projects
	GetActiveProjectsCount() (int, error)

	// GetTotalProjectsCount returns total count of projects
	GetTotalProjectsCount() (int, error)
}

// CacheServiceInterface defines methods for caching
type CacheServiceInterface interface {
	// Get retrieves a value from cache
	Get(key string) (interface{}, bool)

	// Set stores a value in cache with TTL
	Set(key string, value interface{}, ttl time.Duration)

	// Delete removes a value from cache
	Delete(key string)
}

// DashboardService defines the contract for dashboard service
type DashboardService interface {
	// GetOverview returns dashboard overview metrics
	GetOverview() (*dto.OverviewResponse, error)

	// GetProjectStatistics returns statistics for a specific project
	GetProjectStatistics(projectID value_objects.ID) (*dto.ProjectStatisticsResponse, error)

	// GetBuildAnalytics returns build analytics for a given time range
	GetBuildAnalytics(timeRange string) (*dto.BuildAnalyticsResponse, error)
}
