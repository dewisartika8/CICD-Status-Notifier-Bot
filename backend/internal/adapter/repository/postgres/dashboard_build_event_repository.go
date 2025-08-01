package postgres

import (
	"fmt"
	"time"

	"gorm.io/gorm"

	"github.com/dewisartika8/cicd-status-notifier-bot/internal/core/dashboard/domain"
	"github.com/dewisartika8/cicd-status-notifier-bot/internal/core/dashboard/port"
	"github.com/dewisartika8/cicd-status-notifier-bot/internal/core/shared/domain/value_objects"
)

// dashboardBuildEventRepository implements dashboard specific queries
type dashboardBuildEventRepository struct {
	db *gorm.DB
}

// NewDashboardBuildEventRepository creates a new dashboard build event repository
func NewDashboardBuildEventRepository(db *gorm.DB) port.BuildEventRepositoryInterface {
	return &dashboardBuildEventRepository{db: db}
}

// GetOverviewMetrics returns overall build metrics
func (r *dashboardBuildEventRepository) GetOverviewMetrics() (*domain.OverviewMetrics, error) {
	var result struct {
		TotalBuilds      int `json:"total_builds"`
		SuccessfulBuilds int `json:"successful_builds"`
		FailedBuilds     int `json:"failed_builds"`
		AverageDuration  int `json:"average_duration"`
	}

	query := `
		SELECT 
			COUNT(*) as total_builds,
			COUNT(CASE WHEN status = 'success' THEN 1 END) as successful_builds,
			COUNT(CASE WHEN status = 'failed' THEN 1 END) as failed_builds,
			COALESCE(AVG(CASE WHEN duration_seconds IS NOT NULL THEN duration_seconds END), 0) as average_duration
		FROM build_events 
		WHERE created_at >= NOW() - INTERVAL '30 days'
	`

	if err := r.db.Raw(query).Scan(&result).Error; err != nil {
		return nil, fmt.Errorf("failed to get overview metrics: %w", err)
	}

	return &domain.OverviewMetrics{
		TotalBuilds:      result.TotalBuilds,
		SuccessfulBuilds: result.SuccessfulBuilds,
		FailedBuilds:     result.FailedBuilds,
		AverageDuration:  result.AverageDuration,
	}, nil
}

// GetProjectStatistics returns statistics for a specific project
func (r *dashboardBuildEventRepository) GetProjectStatistics(projectID value_objects.ID) (*domain.ProjectStatistics, error) {
	var result struct {
		TotalBuilds      int       `json:"total_builds"`
		SuccessfulBuilds int       `json:"successful_builds"`
		FailedBuilds     int       `json:"failed_builds"`
		AverageDuration  int       `json:"average_duration"`
		LastBuildTime    time.Time `json:"last_build_time"`
	}

	query := `
		SELECT 
			COUNT(*) as total_builds,
			COUNT(CASE WHEN status = 'success' THEN 1 END) as successful_builds,
			COUNT(CASE WHEN status = 'failed' THEN 1 END) as failed_builds,
			COALESCE(AVG(CASE WHEN duration_seconds IS NOT NULL THEN duration_seconds END), 0) as average_duration,
			COALESCE(MAX(created_at), NOW()) as last_build_time
		FROM build_events 
		WHERE project_id = ? AND created_at >= NOW() - INTERVAL '30 days'
	`

	if err := r.db.Raw(query, projectID.Value()).Scan(&result).Error; err != nil {
		return nil, fmt.Errorf("failed to get project statistics: %w", err)
	}

	// Check if project has any builds
	if result.TotalBuilds == 0 {
		return nil, fmt.Errorf("project not found")
	}

	// Get build trends for the last 7 days
	buildTrends, err := r.getBuildTrends(projectID, 7)
	if err != nil {
		return nil, fmt.Errorf("failed to get build trends: %w", err)
	}

	return &domain.ProjectStatistics{
		TotalBuilds:      result.TotalBuilds,
		SuccessfulBuilds: result.SuccessfulBuilds,
		FailedBuilds:     result.FailedBuilds,
		AverageDuration:  result.AverageDuration,
		LastBuildTime:    result.LastBuildTime,
		BuildTrends:      buildTrends,
	}, nil
}

// GetBuildAnalytics returns build analytics for a given time range
func (r *dashboardBuildEventRepository) GetBuildAnalytics(timeRange string) (*domain.BuildAnalytics, error) {
	days := r.parseDaysFromTimeRange(timeRange)

	var result struct {
		TotalBuilds      int `json:"total_builds"`
		SuccessfulBuilds int `json:"successful_builds"`
		FailedBuilds     int `json:"failed_builds"`
		AverageDuration  int `json:"average_duration"`
	}

	query := `
		SELECT 
			COUNT(*) as total_builds,
			COUNT(CASE WHEN status = 'success' THEN 1 END) as successful_builds,
			COUNT(CASE WHEN status = 'failed' THEN 1 END) as failed_builds,
			COALESCE(AVG(CASE WHEN duration_seconds IS NOT NULL THEN duration_seconds END), 0) as average_duration
		FROM build_events 
		WHERE created_at >= NOW() - INTERVAL '%d days'
	`

	if err := r.db.Raw(fmt.Sprintf(query, days), days).Scan(&result).Error; err != nil {
		return nil, fmt.Errorf("failed to get build analytics: %w", err)
	}

	// Get builds by day
	buildsByDay, err := r.getBuildsByDay(days)
	if err != nil {
		return nil, fmt.Errorf("failed to get builds by day: %w", err)
	}

	// Get duration trends
	durationTrends, err := r.getDurationTrends(days)
	if err != nil {
		return nil, fmt.Errorf("failed to get duration trends: %w", err)
	}

	return &domain.BuildAnalytics{
		TotalBuilds:      result.TotalBuilds,
		SuccessfulBuilds: result.SuccessfulBuilds,
		FailedBuilds:     result.FailedBuilds,
		AverageDuration:  result.AverageDuration,
		BuildsByDay:      buildsByDay,
		DurationTrends:   durationTrends,
	}, nil
}

// getBuildTrends returns build trends for a project
func (r *dashboardBuildEventRepository) getBuildTrends(projectID value_objects.ID, days int) ([]domain.BuildTrendInfo, error) {
	var trends []struct {
		Date         time.Time `json:"date"`
		Count        int       `json:"count"`
		SuccessCount int       `json:"success_count"`
	}

	query := `
		SELECT 
			DATE(created_at) as date,
			COUNT(*) as count,
			COUNT(CASE WHEN status = 'success' THEN 1 END) as success_count
		FROM build_events 
		WHERE project_id = ? AND created_at >= NOW() - INTERVAL '%d days'
		GROUP BY DATE(created_at)
		ORDER BY date DESC
	`

	if err := r.db.Raw(fmt.Sprintf(query, days), projectID.Value()).Scan(&trends).Error; err != nil {
		return nil, err
	}

	result := make([]domain.BuildTrendInfo, len(trends))
	for i, trend := range trends {
		result[i] = domain.BuildTrendInfo{
			Date:         trend.Date,
			Count:        trend.Count,
			SuccessCount: trend.SuccessCount,
		}
	}

	return result, nil
}

// getBuildsByDay returns daily build analytics
func (r *dashboardBuildEventRepository) getBuildsByDay(days int) ([]domain.BuildAnalyticsInfo, error) {
	var builds []struct {
		Date             time.Time `json:"date"`
		TotalBuilds      int       `json:"total_builds"`
		SuccessfulBuilds int       `json:"successful_builds"`
		FailedBuilds     int       `json:"failed_builds"`
	}

	query := `
		SELECT 
			DATE(created_at) as date,
			COUNT(*) as total_builds,
			COUNT(CASE WHEN status = 'success' THEN 1 END) as successful_builds,
			COUNT(CASE WHEN status = 'failed' THEN 1 END) as failed_builds
		FROM build_events 
		WHERE created_at >= NOW() - INTERVAL '%d days'
		GROUP BY DATE(created_at)
		ORDER BY date DESC
	`

	if err := r.db.Raw(fmt.Sprintf(query, days)).Scan(&builds).Error; err != nil {
		return nil, err
	}

	result := make([]domain.BuildAnalyticsInfo, len(builds))
	for i, build := range builds {
		result[i] = domain.BuildAnalyticsInfo{
			Date:             build.Date,
			TotalBuilds:      build.TotalBuilds,
			SuccessfulBuilds: build.SuccessfulBuilds,
			FailedBuilds:     build.FailedBuilds,
		}
	}

	return result, nil
}

// getDurationTrends returns build duration trends
func (r *dashboardBuildEventRepository) getDurationTrends(days int) ([]domain.DurationTrendInfo, error) {
	var trends []struct {
		Date            time.Time `json:"date"`
		AverageDuration int       `json:"average_duration"`
	}

	query := `
		SELECT 
			DATE(created_at) as date,
			COALESCE(AVG(CASE WHEN duration_seconds IS NOT NULL THEN duration_seconds END), 0) as average_duration
		FROM build_events 
		WHERE created_at >= NOW() - INTERVAL '%d days' AND duration_seconds IS NOT NULL
		GROUP BY DATE(created_at)
		ORDER BY date DESC
	`

	if err := r.db.Raw(fmt.Sprintf(query, days)).Scan(&trends).Error; err != nil {
		return nil, err
	}

	result := make([]domain.DurationTrendInfo, len(trends))
	for i, trend := range trends {
		result[i] = domain.DurationTrendInfo{
			Date:            trend.Date,
			AverageDuration: trend.AverageDuration,
		}
	}

	return result, nil
}

// parseDaysFromTimeRange converts time range string to days
func (r *dashboardBuildEventRepository) parseDaysFromTimeRange(timeRange string) int {
	switch timeRange {
	case "7d":
		return 7
	case "30d":
		return 30
	case "90d":
		return 90
	default:
		return 7 // default to 7 days
	}
}
