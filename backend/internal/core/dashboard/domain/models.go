package domain

import (
	"time"
)

// OverviewMetrics represents raw metrics data from repository
type OverviewMetrics struct {
	TotalBuilds      int `json:"total_builds"`
	SuccessfulBuilds int `json:"successful_builds"`
	FailedBuilds     int `json:"failed_builds"`
	AverageDuration  int `json:"average_duration"` // in seconds
}

// ProjectStatistics represents project-specific build statistics
type ProjectStatistics struct {
	TotalBuilds      int              `json:"total_builds"`
	SuccessfulBuilds int              `json:"successful_builds"`
	FailedBuilds     int              `json:"failed_builds"`
	AverageDuration  int              `json:"average_duration"` // in seconds
	LastBuildTime    time.Time        `json:"last_build_time"`
	BuildTrends      []BuildTrendInfo `json:"build_trends"`
}

// BuildTrendInfo represents build trend information for a specific date
type BuildTrendInfo struct {
	Date         time.Time `json:"date"`
	Count        int       `json:"count"`
	SuccessCount int       `json:"success_count"`
}

// BuildAnalytics represents comprehensive build analytics
type BuildAnalytics struct {
	TotalBuilds      int                  `json:"total_builds"`
	SuccessfulBuilds int                  `json:"successful_builds"`
	FailedBuilds     int                  `json:"failed_builds"`
	AverageDuration  int                  `json:"average_duration"` // in seconds
	BuildsByDay      []BuildAnalyticsInfo `json:"builds_by_day"`
	DurationTrends   []DurationTrendInfo  `json:"duration_trends"`
}

// BuildAnalyticsInfo represents daily build analytics
type BuildAnalyticsInfo struct {
	Date             time.Time `json:"date"`
	TotalBuilds      int       `json:"total_builds"`
	SuccessfulBuilds int       `json:"successful_builds"`
	FailedBuilds     int       `json:"failed_builds"`
}

// DurationTrendInfo represents build duration trends
type DurationTrendInfo struct {
	Date            time.Time `json:"date"`
	AverageDuration int       `json:"average_duration"` // in seconds
}
