package entities

import (
	"time"

	"github.com/dewisartika8/cicd-status-notifier-bot/internal/core/shared/domain/value_objects"
)

// BuildMetrics represents build metrics for a project
type BuildMetrics struct {
	projectID        value_objects.ID
	totalBuilds      int64
	successfulBuilds int64
	failedBuilds     int64
	averageDuration  time.Duration
	lastBuildStatus  BuildStatus
	lastBuildTime    value_objects.Timestamp
	successRate      float64
}

// NewBuildMetrics creates new build metrics
func NewBuildMetrics(projectID value_objects.ID) *BuildMetrics {
	return &BuildMetrics{
		projectID:   projectID,
		successRate: 0.0,
	}
}

// RestoreBuildMetrics restores build metrics from calculated values
func RestoreBuildMetrics(
	projectID value_objects.ID,
	totalBuilds, successfulBuilds, failedBuilds int64,
	averageDuration time.Duration,
	lastBuildStatus BuildStatus,
	lastBuildTime value_objects.Timestamp,
) *BuildMetrics {
	successRate := 0.0
	if totalBuilds > 0 {
		successRate = float64(successfulBuilds) / float64(totalBuilds) * 100
	}

	return &BuildMetrics{
		projectID:        projectID,
		totalBuilds:      totalBuilds,
		successfulBuilds: successfulBuilds,
		failedBuilds:     failedBuilds,
		averageDuration:  averageDuration,
		lastBuildStatus:  lastBuildStatus,
		lastBuildTime:    lastBuildTime,
		successRate:      successRate,
	}
}

// ProjectID returns the project ID
func (bm *BuildMetrics) ProjectID() value_objects.ID {
	return bm.projectID
}

// TotalBuilds returns the total number of builds
func (bm *BuildMetrics) TotalBuilds() int64 {
	return bm.totalBuilds
}

// SuccessfulBuilds returns the number of successful builds
func (bm *BuildMetrics) SuccessfulBuilds() int64 {
	return bm.successfulBuilds
}

// FailedBuilds returns the number of failed builds
func (bm *BuildMetrics) FailedBuilds() int64 {
	return bm.failedBuilds
}

// AverageDuration returns the average build duration
func (bm *BuildMetrics) AverageDuration() time.Duration {
	return bm.averageDuration
}

// LastBuildStatus returns the last build status
func (bm *BuildMetrics) LastBuildStatus() BuildStatus {
	return bm.lastBuildStatus
}

// LastBuildTime returns the last build time
func (bm *BuildMetrics) LastBuildTime() value_objects.Timestamp {
	return bm.lastBuildTime
}

// SuccessRate returns the success rate as a percentage
func (bm *BuildMetrics) SuccessRate() float64 {
	return bm.successRate
}

// UpdateWithBuildEvent updates metrics with a new build event
func (bm *BuildMetrics) UpdateWithBuildEvent(buildEvent *BuildEvent) {
	// Update counts
	bm.totalBuilds++
	if buildEvent.IsSuccessful() {
		bm.successfulBuilds++
	} else if buildEvent.IsFailed() {
		bm.failedBuilds++
	}

	// Update last build info
	bm.lastBuildStatus = buildEvent.Status()
	bm.lastBuildTime = buildEvent.CreatedAt()

	// Recalculate success rate
	if bm.totalBuilds > 0 {
		bm.successRate = float64(bm.successfulBuilds) / float64(bm.totalBuilds) * 100
	}

	// Update average duration if duration is available
	if buildEvent.DurationSeconds() != nil {
		newDuration := time.Duration(*buildEvent.DurationSeconds()) * time.Second
		if bm.totalBuilds == 1 {
			bm.averageDuration = newDuration
		} else {
			// Calculate new average: ((old_avg * (count-1)) + new_duration) / count
			totalDuration := bm.averageDuration*time.Duration(bm.totalBuilds-1) + newDuration
			bm.averageDuration = totalDuration / time.Duration(bm.totalBuilds)
		}
	}
}

// IsHealthy checks if the project build health is good
func (bm *BuildMetrics) IsHealthy() bool {
	// Consider healthy if success rate is above 80% and we have at least 5 builds
	return bm.totalBuilds >= 5 && bm.successRate >= 80.0
}

// GetHealthStatus returns a health status description
func (bm *BuildMetrics) GetHealthStatus() string {
	if bm.totalBuilds == 0 {
		return "No builds"
	}

	if bm.successRate >= 90.0 {
		return "Excellent"
	} else if bm.successRate >= 80.0 {
		return "Good"
	} else if bm.successRate >= 60.0 {
		return "Fair"
	} else {
		return "Poor"
	}
}
