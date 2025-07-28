package testutils

import (
	"testing"

	"github.com/dewisartika8/cicd-status-notifier-bot/internal/adapters/database"
	"github.com/dewisartika8/cicd-status-notifier-bot/internal/domain/entities"
	"github.com/google/uuid"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// SetupTestDB creates an in-memory SQLite database for testing
func SetupTestDB(t *testing.T) *gorm.DB {
	// Skip SQLite tests if CGO is disabled
	db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		t.Skipf("Skipping SQLite tests (CGO required): %v", err)
		return nil
	}

	// Run migrations
	if err := database.AutoMigrate(db); err != nil {
		t.Fatalf("Failed to migrate test database: %v", err)
	}

	return db
}

// TeardownTestDB closes the test database
func TeardownTestDB(t *testing.T, db *gorm.DB) {
	if db == nil {
		return
	}

	sqlDB, err := db.DB()
	if err != nil {
		t.Errorf("Failed to get underlying sql.DB: %v", err)
		return
	}

	if err := sqlDB.Close(); err != nil {
		t.Errorf("Failed to close test database: %v", err)
	}
}

// CreateTestProject creates a test project entity
func CreateTestProject(name, repoURL string) *entities.Project {
	return entities.NewProject(name, repoURL, "test-secret")
}

// CreateTestBuildEvent creates a test build event entity
func CreateTestBuildEvent(projectID uuid.UUID, eventType entities.EventType, status entities.BuildStatus, branch string) *entities.BuildEvent {
	return entities.NewBuildEvent(projectID, eventType, status, branch)
}

// CreateTestTelegramSubscription creates a test telegram subscription entity
func CreateTestTelegramSubscription(projectID uuid.UUID, chatID int64) *entities.TelegramSubscription {
	return entities.NewTelegramSubscription(projectID, chatID, nil, "test_user")
}

// CreateTestNotificationLog creates a test notification log entity
func CreateTestNotificationLog(buildEventID uuid.UUID, chatID int64) *entities.NotificationLog {
	return entities.NewNotificationLog(buildEventID, chatID)
}
