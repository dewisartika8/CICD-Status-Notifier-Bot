package dto

import (
	"time"

	"github.com/dewisartika8/cicd-status-notifier-bot/internal/core/shared/domain/value_objects"
)

// OverviewResponse represents the dashboard overview data
type OverviewResponse struct {
	TotalProjects    int       `json:"total_projects"`
	ActiveProjects   int       `json:"active_projects"`
	TotalBuilds      int       `json:"total_builds"`
	SuccessfulBuilds int       `json:"successful_builds"`
	FailedBuilds     int       `json:"failed_builds"`
	SuccessRate      float64   `json:"success_rate"`
	AverageDuration  int       `json:"average_duration"` // in seconds
	LastUpdated      time.Time `json:"last_updated"`
}

// ProjectStatisticsResponse represents project-specific statistics
type ProjectStatisticsResponse struct {
	ProjectID        value_objects.ID `json:"project_id"`
	ProjectName      string           `json:"project_name"`
	TotalBuilds      int              `json:"total_builds"`
	SuccessfulBuilds int              `json:"successful_builds"`
	FailedBuilds     int              `json:"failed_builds"`
	SuccessRate      float64          `json:"success_rate"`
	AverageDuration  int              `json:"average_duration"` // in seconds
	LastBuildTime    time.Time        `json:"last_build_time"`
	BuildTrends      []BuildTrendData `json:"build_trends"`
}

// BuildTrendData represents build trend information for a specific date
type BuildTrendData struct {
	Date         time.Time `json:"date"`
	Count        int       `json:"count"`
	SuccessCount int       `json:"success_count"`
}

// BuildAnalyticsResponse represents comprehensive build analytics
type BuildAnalyticsResponse struct {
	TimeRange        string               `json:"time_range"`
	TotalBuilds      int                  `json:"total_builds"`
	SuccessfulBuilds int                  `json:"successful_builds"`
	FailedBuilds     int                  `json:"failed_builds"`
	SuccessRate      float64              `json:"success_rate"`
	AverageDuration  int                  `json:"average_duration"` // in seconds
	BuildsByDay      []BuildAnalyticsData `json:"builds_by_day"`
	DurationTrends   []DurationTrendData  `json:"duration_trends"`
}

// BuildAnalyticsData represents daily build analytics
type BuildAnalyticsData struct {
	Date             time.Time `json:"date"`
	TotalBuilds      int       `json:"total_builds"`
	SuccessfulBuilds int       `json:"successful_builds"`
	FailedBuilds     int       `json:"failed_builds"`
}

// DurationTrendData represents build duration trends
type DurationTrendData struct {
	Date            time.Time `json:"date"`
	AverageDuration int       `json:"average_duration"` // in seconds
}

// MetricsData represents metrics for caching
type MetricsData struct {
	Key         string        `json:"key"`
	Value       interface{}   `json:"value"`
	LastUpdated time.Time     `json:"last_updated"`
	TTL         time.Duration `json:"ttl"`
}
