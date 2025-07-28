package database

import (
	"fmt"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// Config holds database configuration
type Config struct {
	Host         string
	Port         string
	User         string
	Password     string
	DBName       string
	SSLMode      string
	MaxOpenConns int
	MaxIdleConns int
	MaxLifetime  time.Duration
}

// Connect creates a new database connection
func Connect(config Config) (*gorm.DB, error) {
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		config.Host, config.Port, config.User, config.Password, config.DBName, config.SSLMode)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	// Configure connection pool
	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to get underlying sql.DB: %w", err)
	}

	sqlDB.SetMaxOpenConns(config.MaxOpenConns)
	sqlDB.SetMaxIdleConns(config.MaxIdleConns)
	sqlDB.SetConnMaxLifetime(config.MaxLifetime)

	return db, nil
}

// ConnectSQLite creates a SQLite connection for testing
func ConnectSQLite(path string) (*gorm.DB, error) {
	db, err := gorm.Open(sqlite.Open(path), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent), // Silent for tests
	})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to SQLite database: %w", err)
	}

	return db, nil
}

// AutoMigrate runs database migrations
func AutoMigrate(db *gorm.DB) error {
	return db.AutoMigrate(
		&ProjectModel{},
		&BuildEventModel{},
		&TelegramSubscriptionModel{},
		&NotificationLogModel{},
	)
}

// CreateUniqueConstraints adds unique constraints that GORM might miss
func CreateUniqueConstraints(db *gorm.DB) error {
	// Ensure unique constraint on telegram_subscriptions (project_id, chat_id)
	if err := db.Exec(`
		ALTER TABLE telegram_subscriptions 
		ADD CONSTRAINT unique_project_chat 
		UNIQUE (project_id, chat_id)
	`).Error; err != nil {
		// Ignore error if constraint already exists
		if !isConstraintAlreadyExistsError(err) {
			return err
		}
	}

	return nil
}

func isConstraintAlreadyExistsError(err error) bool {
	if err == nil {
		return false
	}
	errStr := err.Error()
	return contains(errStr, "already exists") || contains(errStr, "constraint")
}
