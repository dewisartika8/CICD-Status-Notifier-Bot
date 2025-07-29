package testutils

import (
	"testing"

	builddomain "github.com/dewisartika8/cicd-status-notifier-bot/internal/core/build/domain"
	notificationdomain "github.com/dewisartika8/cicd-status-notifier-bot/internal/core/notification/domain"
	projectdomain "github.com/dewisartika8/cicd-status-notifier-bot/internal/core/project/domain"
	"github.com/dewisartika8/cicd-status-notifier-bot/internal/core/shared/domain/value_objects"
	"github.com/dewisartika8/cicd-status-notifier-bot/pkg/database"
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
func CreateTestProject(name, repoURL string) *projectdomain.Project {
	project, _ := projectdomain.NewProject(name, repoURL, "test-secret", nil)
	return project
}

// CreateTestBuildEvent creates a test build event entity
func CreateTestBuildEvent(projectID uuid.UUID, eventType builddomain.EventType, status builddomain.BuildStatus, branch string) *builddomain.BuildEvent {
	params := builddomain.BuildEventParams{
		ProjectID: value_objects.NewIDFromUUID(projectID),
		EventType: eventType,
		Status:    status,
		Branch:    branch,
	}
	event, _ := builddomain.NewBuildEvent(params)
	return event
}

// CreateTestTelegramSubscription creates a test telegram subscription entity
func CreateTestTelegramSubscription(projectID uuid.UUID, chatID int64) *notificationdomain.TelegramSubscription {
	sub, _ := notificationdomain.NewTelegramSubscription(value_objects.NewIDFromUUID(projectID), chatID)
	return sub
}

// CreateTestNotificationLog creates a test notification log entity
func CreateTestNotificationLog(buildEventID uuid.UUID, chatID int64) *notificationdomain.NotificationLog {
	// For test, use dummy projectID, channel, recipient, and message
	projectID := value_objects.NewID()
	channel := notificationdomain.NotificationChannelTelegram
	recipient := "test_user"
	message := "test message"
	log, _ := notificationdomain.NewNotificationLog(
		value_objects.NewIDFromUUID(buildEventID),
		projectID,
		channel,
		recipient,
		message,
	)
	return log
}
