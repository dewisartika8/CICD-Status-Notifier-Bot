package domain

import (
	"time"

	"github.com/dewisartika8/cicd-status-notifier-bot/internal/core/shared/domain/value_objects"
	"github.com/google/uuid"
)

// RetryConfigurationModel represents the database model for retry configurations
type RetryConfigurationModel struct {
	ID               uuid.UUID `gorm:"column:id;primaryKey;type:uuid;default:uuid_generate_v4()"`
	ProjectID        uuid.UUID `gorm:"column:project_id;not null;type:uuid;index:idx_retry_configurations_project;uniqueIndex:unique_project_channel;constraint:OnDelete:CASCADE"`
	Channel          string    `gorm:"column:channel;not null;index:idx_retry_configurations_channel;uniqueIndex:unique_project_channel"`
	MaxRetries       int       `gorm:"column:max_retries;not null;default:3"`
	BaseDelaySeconds int       `gorm:"column:base_delay_seconds;not null;default:5"`
	MaxDelaySeconds  int       `gorm:"column:max_delay_seconds;not null;default:300"`
	BackoffFactor    float64   `gorm:"column:backoff_factor;not null;default:2.0;type:decimal(3,2)"`
	RetryableErrors  []string  `gorm:"column:retryable_errors;type:text[]"`
	IsActive         bool      `gorm:"column:is_active;default:true;index:idx_retry_configurations_active"`
	CreatedAt        time.Time `gorm:"column:created_at;type:timestamp with time zone;default:current_timestamp"`
	UpdatedAt        time.Time `gorm:"column:updated_at;type:timestamp with time zone;default:current_timestamp"`
}

// TableName returns the table name for retry configurations
func (RetryConfigurationModel) TableName() string {
	return "retry_configurations"
}

// ToEntity converts the model to domain entity
func (m *RetryConfigurationModel) ToEntity() *RetryConfiguration {
	id, _ := value_objects.NewIDFromString(m.ID.String())
	baseDelay := time.Duration(m.BaseDelaySeconds) * time.Second
	maxDelay := time.Duration(m.MaxDelaySeconds) * time.Second
	createdAt := value_objects.NewTimestampFromTime(m.CreatedAt)
	updatedAt := value_objects.NewTimestampFromTime(m.UpdatedAt)

	config := RestoreRetryConfiguration(RestoreRetryConfigurationParams{
		ID:                       id,
		MaxRetryAttempts:         m.MaxRetries,
		InitialRetryDelay:        baseDelay,
		MaxRetryDelay:            maxDelay,
		RetryDelayMultiplier:     m.BackoffFactor,
		RetryTimeoutDuration:     time.Hour, // Default timeout
		EnableExponentialBackoff: true,      // Default value
		EnableDeadLetterQueue:    false,     // Default value
		IsActive:                 m.IsActive,
		CreatedAt:                createdAt,
		UpdatedAt:                updatedAt,
	})

	return config
}

// FromEntity converts domain entity to model
func (m *RetryConfigurationModel) FromEntity(config *RetryConfiguration) {
	// Parse UUID from string
	if id, err := uuid.Parse(config.ID().String()); err == nil {
		m.ID = id
	}

	// Note: ProjectID and Channel don't exist in current domain entity
	// These would need to be passed separately or added to the domain
	m.ProjectID = uuid.New() // Placeholder
	m.Channel = "telegram"   // Default value
	m.MaxRetries = config.MaxRetryAttempts()
	m.BaseDelaySeconds = int(config.InitialRetryDelay().Seconds())
	m.MaxDelaySeconds = int(config.MaxRetryDelay().Seconds())
	m.BackoffFactor = config.RetryDelayMultiplier()
	m.RetryableErrors = []string{} // Default empty array
	m.IsActive = config.IsActive()
	m.CreatedAt = config.CreatedAt().ToTime()
	m.UpdatedAt = config.UpdatedAt().ToTime()
}
