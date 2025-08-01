package service

import (
	"fmt"
	"time"

	"github.com/dewisartika8/cicd-status-notifier-bot/internal/core/dashboard/dto"
	"github.com/dewisartika8/cicd-status-notifier-bot/internal/core/dashboard/port"
	"github.com/dewisartika8/cicd-status-notifier-bot/internal/core/shared/domain/value_objects"
)

// Service implements dashboard business logic
type Service struct {
	buildEventRepo port.BuildEventRepositoryInterface
	projectRepo    port.ProjectRepositoryInterface
	cacheService   port.CacheServiceInterface
}

// Ensure Service implements port.DashboardService
var _ port.DashboardService = (*Service)(nil)

// NewService creates a new dashboard service
func NewService(
	buildEventRepo port.BuildEventRepositoryInterface,
	projectRepo port.ProjectRepositoryInterface,
	cacheService port.CacheServiceInterface,
) *Service {
	return &Service{
		buildEventRepo: buildEventRepo,
		projectRepo:    projectRepo,
		cacheService:   cacheService,
	}
}

// GetOverview returns dashboard overview metrics
func (s *Service) GetOverview() (*dto.OverviewResponse, error) {
	const cacheKey = "dashboard_overview"
	const cacheTTL = 5 * time.Minute

	// Try to get from cache first
	if cached, found := s.cacheService.Get(cacheKey); found {
		if overview, ok := cached.(*dto.OverviewResponse); ok {
			return overview, nil
		}
	}

	// Get metrics from repository
	metrics, err := s.buildEventRepo.GetOverviewMetrics()
	if err != nil {
		return nil, fmt.Errorf("failed to get overview metrics: %w", err)
	}

	// Get project counts
	totalProjects, err := s.projectRepo.GetTotalProjectsCount()
	if err != nil {
		return nil, fmt.Errorf("failed to get total projects count: %w", err)
	}

	activeProjects, err := s.projectRepo.GetActiveProjectsCount()
	if err != nil {
		return nil, fmt.Errorf("failed to get active projects count: %w", err)
	}

	// Calculate success rate
	successRate := 0.0
	if metrics.TotalBuilds > 0 {
		successRate = float64(metrics.SuccessfulBuilds) / float64(metrics.TotalBuilds) * 100
	}

	overview := &dto.OverviewResponse{
		TotalProjects:    totalProjects,
		ActiveProjects:   activeProjects,
		TotalBuilds:      metrics.TotalBuilds,
		SuccessfulBuilds: metrics.SuccessfulBuilds,
		FailedBuilds:     metrics.FailedBuilds,
		SuccessRate:      successRate,
		AverageDuration:  metrics.AverageDuration,
		LastUpdated:      time.Now(),
	}

	// Cache the result
	s.cacheService.Set(cacheKey, overview, cacheTTL)

	return overview, nil
}

// GetProjectStatistics returns statistics for a specific project
func (s *Service) GetProjectStatistics(projectID value_objects.ID) (*dto.ProjectStatisticsResponse, error) {
	cacheKey := fmt.Sprintf("project_statistics_%s", projectID.String())
	const cacheTTL = 2 * time.Minute

	// Try to get from cache first
	if cached, found := s.cacheService.Get(cacheKey); found {
		if statistics, ok := cached.(*dto.ProjectStatisticsResponse); ok {
			return statistics, nil
		}
	}

	// Get statistics from repository
	stats, err := s.buildEventRepo.GetProjectStatistics(projectID)
	if err != nil {
		return nil, fmt.Errorf("failed to get project statistics: %w", err)
	}

	// Get project name
	projectName, err := s.projectRepo.GetProjectNameByID(projectID)
	if err != nil {
		return nil, fmt.Errorf("failed to get project name: %w", err)
	}

	// Calculate success rate
	successRate := 0.0
	if stats.TotalBuilds > 0 {
		successRate = float64(stats.SuccessfulBuilds) / float64(stats.TotalBuilds) * 100
	}

	// Convert build trends
	buildTrends := make([]dto.BuildTrendData, len(stats.BuildTrends))
	for i, trend := range stats.BuildTrends {
		buildTrends[i] = dto.BuildTrendData{
			Date:         trend.Date,
			Count:        trend.Count,
			SuccessCount: trend.SuccessCount,
		}
	}

	statistics := &dto.ProjectStatisticsResponse{
		ProjectID:        projectID,
		ProjectName:      projectName,
		TotalBuilds:      stats.TotalBuilds,
		SuccessfulBuilds: stats.SuccessfulBuilds,
		FailedBuilds:     stats.FailedBuilds,
		SuccessRate:      successRate,
		AverageDuration:  stats.AverageDuration,
		LastBuildTime:    stats.LastBuildTime,
		BuildTrends:      buildTrends,
	}

	// Cache the result
	s.cacheService.Set(cacheKey, statistics, cacheTTL)

	return statistics, nil
}

// GetBuildAnalytics returns build analytics for a given time range
func (s *Service) GetBuildAnalytics(timeRange string) (*dto.BuildAnalyticsResponse, error) {
	cacheKey := fmt.Sprintf("build_analytics_%s", timeRange)
	const cacheTTL = 10 * time.Minute

	// Try to get from cache first
	if cached, found := s.cacheService.Get(cacheKey); found {
		if analytics, ok := cached.(*dto.BuildAnalyticsResponse); ok {
			return analytics, nil
		}
	}

	// Get analytics from repository
	analytics, err := s.buildEventRepo.GetBuildAnalytics(timeRange)
	if err != nil {
		return nil, fmt.Errorf("failed to get build analytics: %w", err)
	}

	// Calculate success rate
	successRate := 0.0
	if analytics.TotalBuilds > 0 {
		successRate = float64(analytics.SuccessfulBuilds) / float64(analytics.TotalBuilds) * 100
	}

	// Convert builds by day
	buildsByDay := make([]dto.BuildAnalyticsData, len(analytics.BuildsByDay))
	for i, data := range analytics.BuildsByDay {
		buildsByDay[i] = dto.BuildAnalyticsData{
			Date:             data.Date,
			TotalBuilds:      data.TotalBuilds,
			SuccessfulBuilds: data.SuccessfulBuilds,
			FailedBuilds:     data.FailedBuilds,
		}
	}

	// Convert duration trends
	durationTrends := make([]dto.DurationTrendData, len(analytics.DurationTrends))
	for i, trend := range analytics.DurationTrends {
		durationTrends[i] = dto.DurationTrendData{
			Date:            trend.Date,
			AverageDuration: trend.AverageDuration,
		}
	}

	result := &dto.BuildAnalyticsResponse{
		TimeRange:        timeRange,
		TotalBuilds:      analytics.TotalBuilds,
		SuccessfulBuilds: analytics.SuccessfulBuilds,
		FailedBuilds:     analytics.FailedBuilds,
		SuccessRate:      successRate,
		AverageDuration:  analytics.AverageDuration,
		BuildsByDay:      buildsByDay,
		DurationTrends:   durationTrends,
	}

	// Cache the result
	s.cacheService.Set(cacheKey, result, cacheTTL)

	return result, nil
}
