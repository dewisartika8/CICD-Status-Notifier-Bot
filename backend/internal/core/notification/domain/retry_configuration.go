package domain

import (
	"time"

	"github.com/dewisartika8/cicd-status-notifier-bot/internal/core/shared/domain/value_objects"
)

// RetryConfiguration represents retry configuration for failed notifications
type RetryConfiguration struct {
	id                       value_objects.ID
	maxRetryAttempts         int
	initialRetryDelay        time.Duration
	maxRetryDelay            time.Duration
	retryDelayMultiplier     float64
	retryTimeoutDuration     time.Duration
	enableExponentialBackoff bool
	enableDeadLetterQueue    bool
	isActive                 bool
	createdAt                value_objects.Timestamp
	updatedAt                value_objects.Timestamp
}

// RetryPolicy represents a retry policy for different notification types
type RetryPolicy struct {
	channel       NotificationChannel
	maxAttempts   int
	baseDelay     time.Duration
	maxDelay      time.Duration
	backoffFactor float64
}

// Getter methods for RetryPolicy
func (rp *RetryPolicy) Channel() NotificationChannel {
	return rp.channel
}

func (rp *RetryPolicy) MaxAttempts() int {
	return rp.maxAttempts
}

func (rp *RetryPolicy) BaseDelay() time.Duration {
	return rp.baseDelay
}

func (rp *RetryPolicy) MaxDelay() time.Duration {
	return rp.maxDelay
}

func (rp *RetryPolicy) BackoffFactor() float64 {
	return rp.backoffFactor
}

// NewRetryConfiguration creates a new retry configuration entity
func NewRetryConfiguration(
	maxRetryAttempts int,
	initialRetryDelay, maxRetryDelay, retryTimeoutDuration time.Duration,
	retryDelayMultiplier float64,
	enableExponentialBackoff, enableDeadLetterQueue bool,
) (*RetryConfiguration, error) {
	// Validate configuration
	if maxRetryAttempts < 0 || maxRetryAttempts > 10 {
		return nil, NewInvalidRetryConfigurationError("max retry attempts must be between 0 and 10")
	}

	if initialRetryDelay < 0 {
		return nil, NewInvalidRetryConfigurationError("initial retry delay cannot be negative")
	}

	if maxRetryDelay < initialRetryDelay {
		return nil, NewInvalidRetryConfigurationError("max retry delay must be greater than or equal to initial delay")
	}

	if retryDelayMultiplier < 1.0 {
		return nil, NewInvalidRetryConfigurationError("retry delay multiplier must be at least 1.0")
	}

	if retryTimeoutDuration < 0 {
		return nil, NewInvalidRetryConfigurationError("retry timeout duration cannot be negative")
	}

	now := value_objects.NewTimestamp()
	id := value_objects.NewID()

	return &RetryConfiguration{
		id:                       id,
		maxRetryAttempts:         maxRetryAttempts,
		initialRetryDelay:        initialRetryDelay,
		maxRetryDelay:            maxRetryDelay,
		retryDelayMultiplier:     retryDelayMultiplier,
		retryTimeoutDuration:     retryTimeoutDuration,
		enableExponentialBackoff: enableExponentialBackoff,
		enableDeadLetterQueue:    enableDeadLetterQueue,
		isActive:                 true,
		createdAt:                now,
		updatedAt:                now,
	}, nil
}

// RestoreRetryConfiguration restores a retry configuration from persistence
func RestoreRetryConfiguration(params RestoreRetryConfigurationParams) *RetryConfiguration {
	return &RetryConfiguration{
		id:                       params.ID,
		maxRetryAttempts:         params.MaxRetryAttempts,
		initialRetryDelay:        params.InitialRetryDelay,
		maxRetryDelay:            params.MaxRetryDelay,
		retryDelayMultiplier:     params.RetryDelayMultiplier,
		retryTimeoutDuration:     params.RetryTimeoutDuration,
		enableExponentialBackoff: params.EnableExponentialBackoff,
		enableDeadLetterQueue:    params.EnableDeadLetterQueue,
		isActive:                 params.IsActive,
		createdAt:                params.CreatedAt,
		updatedAt:                params.UpdatedAt,
	}
}

// RestoreRetryConfigurationParams holds parameters for restoring retry configuration
type RestoreRetryConfigurationParams struct {
	ID                       value_objects.ID
	MaxRetryAttempts         int
	InitialRetryDelay        time.Duration
	MaxRetryDelay            time.Duration
	RetryDelayMultiplier     float64
	RetryTimeoutDuration     time.Duration
	EnableExponentialBackoff bool
	EnableDeadLetterQueue    bool
	IsActive                 bool
	CreatedAt                value_objects.Timestamp
	UpdatedAt                value_objects.Timestamp
}

// Getters
func (rc *RetryConfiguration) ID() value_objects.ID {
	return rc.id
}

func (rc *RetryConfiguration) MaxRetryAttempts() int {
	return rc.maxRetryAttempts
}

func (rc *RetryConfiguration) InitialRetryDelay() time.Duration {
	return rc.initialRetryDelay
}

func (rc *RetryConfiguration) MaxRetryDelay() time.Duration {
	return rc.maxRetryDelay
}

func (rc *RetryConfiguration) RetryDelayMultiplier() float64 {
	return rc.retryDelayMultiplier
}

func (rc *RetryConfiguration) RetryTimeoutDuration() time.Duration {
	return rc.retryTimeoutDuration
}

func (rc *RetryConfiguration) EnableExponentialBackoff() bool {
	return rc.enableExponentialBackoff
}

func (rc *RetryConfiguration) EnableDeadLetterQueue() bool {
	return rc.enableDeadLetterQueue
}

func (rc *RetryConfiguration) IsActive() bool {
	return rc.isActive
}

func (rc *RetryConfiguration) CreatedAt() value_objects.Timestamp {
	return rc.createdAt
}

func (rc *RetryConfiguration) UpdatedAt() value_objects.Timestamp {
	return rc.updatedAt
}

// Business methods

// CalculateRetryDelay calculates the retry delay for the given attempt
func (rc *RetryConfiguration) CalculateRetryDelay(attemptNumber int) time.Duration {
	if !rc.enableExponentialBackoff {
		return rc.initialRetryDelay
	}

	// For first attempt, use initial delay
	// For subsequent attempts, apply exponential backoff
	delay := float64(rc.initialRetryDelay.Nanoseconds())
	for i := 0; i < attemptNumber-1; i++ {
		delay *= rc.retryDelayMultiplier
	}

	calculatedDelay := time.Duration(delay)
	if calculatedDelay > rc.maxRetryDelay {
		return rc.maxRetryDelay
	}

	return calculatedDelay
}

// ShouldRetry determines if a notification should be retried
func (rc *RetryConfiguration) ShouldRetry(attemptCount int, lastError error) bool {
	if !rc.isActive {
		return false
	}

	if attemptCount >= rc.maxRetryAttempts {
		return false
	}

	// Check if error is retryable (implementation specific)
	return rc.isRetryableError(lastError)
}

// UpdateConfiguration updates the retry configuration
func (rc *RetryConfiguration) UpdateConfiguration(
	maxRetryAttempts int,
	initialRetryDelay, maxRetryDelay, retryTimeoutDuration time.Duration,
	retryDelayMultiplier float64,
	enableExponentialBackoff, enableDeadLetterQueue bool,
) error {
	// Validate new configuration
	if maxRetryAttempts < 0 || maxRetryAttempts > 10 {
		return NewInvalidRetryConfigurationError("max retry attempts must be between 0 and 10")
	}

	if initialRetryDelay < 0 {
		return NewInvalidRetryConfigurationError("initial retry delay cannot be negative")
	}

	if maxRetryDelay < initialRetryDelay {
		return NewInvalidRetryConfigurationError("max retry delay must be greater than or equal to initial delay")
	}

	if retryDelayMultiplier < 1.0 {
		return NewInvalidRetryConfigurationError("retry delay multiplier must be at least 1.0")
	}

	if retryTimeoutDuration < 0 {
		return NewInvalidRetryConfigurationError("retry timeout duration cannot be negative")
	}

	// Update fields
	rc.maxRetryAttempts = maxRetryAttempts
	rc.initialRetryDelay = initialRetryDelay
	rc.maxRetryDelay = maxRetryDelay
	rc.retryDelayMultiplier = retryDelayMultiplier
	rc.retryTimeoutDuration = retryTimeoutDuration
	rc.enableExponentialBackoff = enableExponentialBackoff
	rc.enableDeadLetterQueue = enableDeadLetterQueue
	rc.updatedAt = value_objects.NewTimestamp()

	return nil
}

// Activate activates the retry configuration
func (rc *RetryConfiguration) Activate() error {
	if rc.isActive {
		return ErrRetryConfigurationAlreadyActive
	}

	rc.isActive = true
	rc.updatedAt = value_objects.NewTimestamp()
	return nil
}

// Deactivate deactivates the retry configuration
func (rc *RetryConfiguration) Deactivate() error {
	if !rc.isActive {
		return ErrRetryConfigurationAlreadyInactive
	}

	rc.isActive = false
	rc.updatedAt = value_objects.NewTimestamp()
	return nil
}

// Helper methods

// isRetryableError determines if an error is retryable
func (rc *RetryConfiguration) isRetryableError(err error) bool {
	if err == nil {
		return false
	}

	// Check for specific error types that should not be retried
	errorString := err.Error()

	// Non-retryable errors
	nonRetryableErrors := []string{
		"invalid recipient",
		"invalid message",
		"invalid chat ID",
		"forbidden",
		"unauthorized",
		"not found",
	}

	for _, nonRetryable := range nonRetryableErrors {
		if containsIgnoreCase(errorString, nonRetryable) {
			return false
		}
	}

	// All other errors are considered retryable (network issues, temporary failures, etc.)
	return true
}

// containsIgnoreCase checks if a string contains a substring (case-insensitive)
func containsIgnoreCase(str, substr string) bool {
	return len(str) >= len(substr) &&
		(str == substr ||
			(len(str) > len(substr) &&
				findIgnoreCase(str, substr) >= 0))
}

// findIgnoreCase finds the index of substr in str (case-insensitive)
func findIgnoreCase(str, substr string) int {
	for i := 0; i <= len(str)-len(substr); i++ {
		match := true
		for j := 0; j < len(substr); j++ {
			if toLower(str[i+j]) != toLower(substr[j]) {
				match = false
				break
			}
		}
		if match {
			return i
		}
	}
	return -1
}

// toLower converts a byte to lowercase
func toLower(b byte) byte {
	if b >= 'A' && b <= 'Z' {
		return b + ('a' - 'A')
	}
	return b
}

// GetDefaultRetryPolicies returns default retry policies for different channels
func GetDefaultRetryPolicies() map[NotificationChannel]RetryPolicy {
	return map[NotificationChannel]RetryPolicy{
		NotificationChannelTelegram: {
			channel:       NotificationChannelTelegram,
			maxAttempts:   3,
			baseDelay:     30 * time.Second,
			maxDelay:      5 * time.Minute,
			backoffFactor: 2.0,
		},
		NotificationChannelEmail: {
			channel:       NotificationChannelEmail,
			maxAttempts:   5,
			baseDelay:     1 * time.Minute,
			maxDelay:      10 * time.Minute,
			backoffFactor: 1.5,
		},
		NotificationChannelSlack: {
			channel:       NotificationChannelSlack,
			maxAttempts:   3,
			baseDelay:     30 * time.Second,
			maxDelay:      5 * time.Minute,
			backoffFactor: 2.0,
		},
		NotificationChannelWebhook: {
			channel:       NotificationChannelWebhook,
			maxAttempts:   3,
			baseDelay:     15 * time.Second,
			maxDelay:      2 * time.Minute,
			backoffFactor: 2.0,
		},
	}
}

// GetDefaultRetryConfiguration returns a default retry configuration
func GetDefaultRetryConfiguration() *RetryConfiguration {
	config, _ := NewRetryConfiguration(
		3,              // maxRetryAttempts
		30*time.Second, // initialRetryDelay
		5*time.Minute,  // maxRetryDelay
		30*time.Minute, // retryTimeoutDuration
		2.0,            // retryDelayMultiplier
		true,           // enableExponentialBackoff
		true,           // enableDeadLetterQueue
	)
	return config
}
