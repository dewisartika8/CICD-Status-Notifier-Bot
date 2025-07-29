package domain

import (
	"time"

	"github.com/dewisartika8/cicd-status-notifier-bot/internal/core/shared/domain/value_objects"
)

// RetryConfigurationModel represents the database model for retry configurations
type RetryConfigurationModel struct {
	ID               string    `gorm:"column:id;primaryKey" json:"id"`
	ProjectID        string    `gorm:"column:project_id;not null" json:"project_id"`
	Channel          string    `gorm:"column:channel;not null" json:"channel"`
	MaxRetries       int       `gorm:"column:max_retries;not null;default:3" json:"max_retries"`
	BaseDelaySeconds int       `gorm:"column:base_delay_seconds;not null;default:5" json:"base_delay_seconds"`
	MaxDelaySeconds  int       `gorm:"column:max_delay_seconds;not null;default:300" json:"max_delay_seconds"`
	BackoffFactor    float64   `gorm:"column:backoff_factor;not null;default:2.0" json:"backoff_factor"`
	RetryableErrors  []string  `gorm:"column:retryable_errors;type:text[]" json:"retryable_errors"`
	IsActive         bool      `gorm:"column:is_active;default:true" json:"is_active"`
	CreatedAt        time.Time `gorm:"column:created_at;autoCreateTime" json:"created_at"`
	UpdatedAt        time.Time `gorm:"column:updated_at;autoUpdateTime" json:"updated_at"`
}

// TableName returns the table name for retry configurations
func (RetryConfigurationModel) TableName() string {
	return "retry_configurations"
}

// ToEntity converts the model to domain entity
func (m *RetryConfigurationModel) ToEntity() *RetryConfiguration {
	id, _ := value_objects.NewIDFromString(m.ID)
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
	m.ID = config.ID().String()
	m.ProjectID = "" // This field doesn't exist in the current domain
	m.Channel = ""   // This field doesn't exist in the current domain
	m.MaxRetries = config.MaxRetryAttempts()
	m.BaseDelaySeconds = int(config.InitialRetryDelay().Seconds())
	m.MaxDelaySeconds = int(config.MaxRetryDelay().Seconds())
	m.BackoffFactor = config.RetryDelayMultiplier()
	m.RetryableErrors = []string{} // This field doesn't exist in the current domain
	m.IsActive = config.IsActive()
	m.CreatedAt = config.CreatedAt().ToTime()
	m.UpdatedAt = config.UpdatedAt().ToTime()
}
