package repositories_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"github.com/dewisartika8/cicd-status-notifier-bot/internal/adapter/repository/postgres"
	"github.com/dewisartika8/cicd-status-notifier-bot/internal/core/dashboard/port"
	"github.com/dewisartika8/cicd-status-notifier-bot/internal/core/shared/domain/value_objects"
)

type DashboardBuildEventRepositoryTestSuite struct {
	suite.Suite
	db   *gorm.DB
	repo port.BuildEventRepositoryInterface
}

func (suite *DashboardBuildEventRepositoryTestSuite) SetupSuite() {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	suite.Require().NoError(err)

	// Create tables
	err = db.Exec(`
		CREATE TABLE build_events (
			id TEXT PRIMARY KEY,
			project_id TEXT NOT NULL,
			event_type TEXT NOT NULL,
			status TEXT NOT NULL,
			branch TEXT NOT NULL,
			commit_sha TEXT,
			commit_message TEXT,
			author_name TEXT,
			author_email TEXT,
			build_url TEXT,
			duration_seconds INTEGER,
			webhook_payload TEXT,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP
		)
	`).Error
	suite.Require().NoError(err)

	err = db.Exec(`
		CREATE TABLE projects (
			id TEXT PRIMARY KEY,
			name TEXT NOT NULL,
			repository_url TEXT NOT NULL,
			webhook_secret TEXT,
			telegram_chat_id INTEGER,
			is_active BOOLEAN DEFAULT true,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
		)
	`).Error
	suite.Require().NoError(err)

	suite.db = db
	suite.repo = postgres.NewDashboardBuildEventRepository(db)
}

func (suite *DashboardBuildEventRepositoryTestSuite) TearDownSuite() {
	db, _ := suite.db.DB()
	db.Close()
}

func (suite *DashboardBuildEventRepositoryTestSuite) SetupTest() {
	// Clean up tables
	suite.db.Exec("DELETE FROM build_events")
	suite.db.Exec("DELETE FROM projects")
}

func (suite *DashboardBuildEventRepositoryTestSuite) TestGetOverviewMetrics() {
	// Insert test data
	projectID := value_objects.NewID()
	now := time.Now()

	suite.db.Exec(`
		INSERT INTO build_events (id, project_id, event_type, status, branch, duration_seconds, created_at)
		VALUES 
		(?, ?, 'build_completed', 'success', 'main', 300, ?),
		(?, ?, 'build_completed', 'success', 'main', 250, ?),
		(?, ?, 'build_completed', 'failed', 'main', 180, ?),
		(?, ?, 'build_completed', 'success', 'develop', 320, ?)
	`,
		value_objects.NewID().String(), projectID.String(), now.Add(-1*time.Hour),
		value_objects.NewID().String(), projectID.String(), now.Add(-2*time.Hour),
		value_objects.NewID().String(), projectID.String(), now.Add(-3*time.Hour),
		value_objects.NewID().String(), projectID.String(), now.Add(-4*time.Hour),
	)

	metrics, err := suite.repo.GetOverviewMetrics()

	assert.NoError(suite.T(), err)
	assert.NotNil(suite.T(), metrics)
	assert.Equal(suite.T(), 4, metrics.TotalBuilds)
	assert.Equal(suite.T(), 3, metrics.SuccessfulBuilds)
	assert.Equal(suite.T(), 1, metrics.FailedBuilds)
	assert.Greater(suite.T(), metrics.AverageDuration, 0)
}

func (suite *DashboardBuildEventRepositoryTestSuite) TestGetProjectStatistics() {
	projectID := value_objects.NewID()
	now := time.Now()

	// Insert test data
	suite.db.Exec(`
		INSERT INTO build_events (id, project_id, event_type, status, branch, duration_seconds, created_at)
		VALUES 
		(?, ?, 'build_completed', 'success', 'main', 300, ?),
		(?, ?, 'build_completed', 'failed', 'main', 180, ?)
	`,
		value_objects.NewID().String(), projectID.String(), now.Add(-1*time.Hour),
		value_objects.NewID().String(), projectID.String(), now.Add(-2*time.Hour),
	)

	stats, err := suite.repo.GetProjectStatistics(projectID)

	assert.NoError(suite.T(), err)
	assert.NotNil(suite.T(), stats)
	assert.Equal(suite.T(), 2, stats.TotalBuilds)
	assert.Equal(suite.T(), 1, stats.SuccessfulBuilds)
	assert.Equal(suite.T(), 1, stats.FailedBuilds)
	assert.Greater(suite.T(), stats.AverageDuration, 0)
	assert.NotNil(suite.T(), stats.BuildTrends)
}

func (suite *DashboardBuildEventRepositoryTestSuite) TestGetProjectStatisticsNotFound() {
	nonExistentProjectID := value_objects.NewID()

	stats, err := suite.repo.GetProjectStatistics(nonExistentProjectID)

	assert.Error(suite.T(), err)
	assert.Nil(suite.T(), stats)
	assert.Contains(suite.T(), err.Error(), "project not found")
}

func (suite *DashboardBuildEventRepositoryTestSuite) TestGetBuildAnalytics() {
	projectID := value_objects.NewID()
	now := time.Now()

	// Insert test data for the last 7 days
	suite.db.Exec(`
		INSERT INTO build_events (id, project_id, event_type, status, branch, duration_seconds, created_at)
		VALUES 
		(?, ?, 'build_completed', 'success', 'main', 300, ?),
		(?, ?, 'build_completed', 'success', 'main', 250, ?),
		(?, ?, 'build_completed', 'failed', 'main', 180, ?)
	`,
		value_objects.NewID().String(), projectID.String(), now.Add(-1*time.Hour),
		value_objects.NewID().String(), projectID.String(), now.Add(-25*time.Hour),
		value_objects.NewID().String(), projectID.String(), now.Add(-49*time.Hour),
	)

	analytics, err := suite.repo.GetBuildAnalytics("7d")

	assert.NoError(suite.T(), err)
	assert.NotNil(suite.T(), analytics)
	assert.Equal(suite.T(), 3, analytics.TotalBuilds)
	assert.Equal(suite.T(), 2, analytics.SuccessfulBuilds)
	assert.Equal(suite.T(), 1, analytics.FailedBuilds)
	assert.Greater(suite.T(), analytics.AverageDuration, 0)
	assert.NotNil(suite.T(), analytics.BuildsByDay)
	assert.NotNil(suite.T(), analytics.DurationTrends)
}

func TestDashboardBuildEventRepositoryTestSuite(t *testing.T) {
	suite.Run(t, new(DashboardBuildEventRepositoryTestSuite))
}
