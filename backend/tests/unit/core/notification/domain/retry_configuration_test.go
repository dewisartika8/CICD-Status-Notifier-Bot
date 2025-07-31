package domain_test

import (
	"errors"
	"testing"
	"time"

	"github.com/dewisartika8/cicd-status-notifier-bot/internal/core/notification/domain"
	"github.com/dewisartika8/cicd-status-notifier-bot/internal/core/shared/domain/value_objects"
	"github.com/stretchr/testify/assert"
)

func TestNewRetryConfiguration(t *testing.T) {
	tests := []struct {
		name                     string
		maxRetryAttempts         int
		initialRetryDelay        time.Duration
		maxRetryDelay            time.Duration
		retryTimeoutDuration     time.Duration
		retryDelayMultiplier     float64
		enableExponentialBackoff bool
		enableDeadLetterQueue    bool
		wantErr                  bool
		expectedErr              string
	}{
		{
			name:                     "valid_configuration",
			maxRetryAttempts:         3,
			initialRetryDelay:        30 * time.Second,
			maxRetryDelay:            5 * time.Minute,
			retryTimeoutDuration:     10 * time.Minute,
			retryDelayMultiplier:     2.0,
			enableExponentialBackoff: true,
			enableDeadLetterQueue:    true,
			wantErr:                  false,
		},
		{
			name:                     "invalid_max_retry_attempts_negative",
			maxRetryAttempts:         -1,
			initialRetryDelay:        30 * time.Second,
			maxRetryDelay:            5 * time.Minute,
			retryTimeoutDuration:     10 * time.Minute,
			retryDelayMultiplier:     2.0,
			enableExponentialBackoff: true,
			enableDeadLetterQueue:    true,
			wantErr:                  true,
			expectedErr:              "max retry attempts must be between 0 and 10",
		},
		{
			name:                     "invalid_max_retry_attempts_too_high",
			maxRetryAttempts:         15,
			initialRetryDelay:        30 * time.Second,
			maxRetryDelay:            5 * time.Minute,
			retryTimeoutDuration:     10 * time.Minute,
			retryDelayMultiplier:     2.0,
			enableExponentialBackoff: true,
			enableDeadLetterQueue:    true,
			wantErr:                  true,
			expectedErr:              "max retry attempts must be between 0 and 10",
		},
		{
			name:                     "invalid_initial_delay_negative",
			maxRetryAttempts:         3,
			initialRetryDelay:        -1 * time.Second,
			maxRetryDelay:            5 * time.Minute,
			retryTimeoutDuration:     10 * time.Minute,
			retryDelayMultiplier:     2.0,
			enableExponentialBackoff: true,
			enableDeadLetterQueue:    true,
			wantErr:                  true,
			expectedErr:              "initial retry delay cannot be negative",
		},
		{
			name:                     "invalid_max_delay_less_than_initial",
			maxRetryAttempts:         3,
			initialRetryDelay:        5 * time.Minute,
			maxRetryDelay:            1 * time.Minute,
			retryTimeoutDuration:     10 * time.Minute,
			retryDelayMultiplier:     2.0,
			enableExponentialBackoff: true,
			enableDeadLetterQueue:    true,
			wantErr:                  true,
			expectedErr:              "max retry delay must be greater than or equal to initial delay",
		},
		{
			name:                     "invalid_multiplier_less_than_one",
			maxRetryAttempts:         3,
			initialRetryDelay:        30 * time.Second,
			maxRetryDelay:            5 * time.Minute,
			retryTimeoutDuration:     10 * time.Minute,
			retryDelayMultiplier:     0.5,
			enableExponentialBackoff: true,
			enableDeadLetterQueue:    true,
			wantErr:                  true,
			expectedErr:              "retry delay multiplier must be at least 1.0",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config, err := domain.NewRetryConfiguration(
				tt.maxRetryAttempts,
				tt.initialRetryDelay,
				tt.maxRetryDelay,
				tt.retryTimeoutDuration,
				tt.retryDelayMultiplier,
				tt.enableExponentialBackoff,
				tt.enableDeadLetterQueue,
			)

			if tt.wantErr {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedErr)
				assert.Nil(t, config)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, config)
				assert.Equal(t, tt.maxRetryAttempts, config.MaxRetryAttempts())
				assert.Equal(t, tt.initialRetryDelay, config.InitialRetryDelay())
				assert.Equal(t, tt.maxRetryDelay, config.MaxRetryDelay())
				assert.Equal(t, tt.retryDelayMultiplier, config.RetryDelayMultiplier())
				assert.Equal(t, tt.enableExponentialBackoff, config.EnableExponentialBackoff())
				assert.Equal(t, tt.enableDeadLetterQueue, config.EnableDeadLetterQueue())
				assert.True(t, config.IsActive())
				assert.False(t, config.ID().IsNil())
			}
		})
	}
}

func TestRetryConfiguration_CalculateRetryDelay(t *testing.T) {
	tests := []struct {
		name                     string
		enableExponentialBackoff bool
		initialDelay             time.Duration
		maxDelay                 time.Duration
		multiplier               float64
		attemptNumber            int
		expectedDelay            time.Duration
	}{
		{
			name:                     "exponential_backoff_disabled",
			enableExponentialBackoff: false,
			initialDelay:             30 * time.Second,
			maxDelay:                 5 * time.Minute,
			multiplier:               2.0,
			attemptNumber:            3,
			expectedDelay:            30 * time.Second,
		},
		{
			name:                     "exponential_backoff_first_attempt",
			enableExponentialBackoff: true,
			initialDelay:             30 * time.Second,
			maxDelay:                 5 * time.Minute,
			multiplier:               2.0,
			attemptNumber:            1,
			expectedDelay:            30 * time.Second,
		},
		{
			name:                     "exponential_backoff_second_attempt",
			enableExponentialBackoff: true,
			initialDelay:             30 * time.Second,
			maxDelay:                 5 * time.Minute,
			multiplier:               2.0,
			attemptNumber:            2,
			expectedDelay:            60 * time.Second,
		},
		{
			name:                     "exponential_backoff_third_attempt",
			enableExponentialBackoff: true,
			initialDelay:             30 * time.Second,
			maxDelay:                 5 * time.Minute,
			multiplier:               2.0,
			attemptNumber:            3,
			expectedDelay:            120 * time.Second,
		},
		{
			name:                     "exponential_backoff_capped_at_max",
			enableExponentialBackoff: true,
			initialDelay:             30 * time.Second,
			maxDelay:                 1 * time.Minute,
			multiplier:               2.0,
			attemptNumber:            5,
			expectedDelay:            1 * time.Minute,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config, err := domain.NewRetryConfiguration(
				5,
				tt.initialDelay,
				tt.maxDelay,
				10*time.Minute,
				tt.multiplier,
				tt.enableExponentialBackoff,
				true,
			)
			assert.NoError(t, err)

			delay := config.CalculateRetryDelay(tt.attemptNumber)
			assert.Equal(t, tt.expectedDelay, delay)
		})
	}
}

func TestRetryConfiguration_ShouldRetry(t *testing.T) {
	config, err := domain.NewRetryConfiguration(
		3,
		30*time.Second,
		5*time.Minute,
		10*time.Minute,
		2.0,
		true,
		true,
	)
	assert.NoError(t, err)

	tests := []struct {
		name         string
		attemptCount int
		lastError    error
		isActive     bool
		expected     bool
	}{
		{
			name:         "should_retry_first_attempt_with_retryable_error",
			attemptCount: 1,
			lastError:    errors.New("network timeout"),
			isActive:     true,
			expected:     true,
		},
		{
			name:         "should_not_retry_max_attempts_reached",
			attemptCount: 3,
			lastError:    errors.New("network timeout"),
			isActive:     true,
			expected:     false,
		},
		{
			name:         "should_not_retry_non_retryable_error",
			attemptCount: 1,
			lastError:    errors.New("invalid recipient"),
			isActive:     true,
			expected:     false,
		},
		{
			name:         "should_not_retry_config_inactive",
			attemptCount: 1,
			lastError:    errors.New("network timeout"),
			isActive:     false,
			expected:     false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if !tt.isActive {
				err := config.Deactivate()
				assert.NoError(t, err)
			}

			result := config.ShouldRetry(tt.attemptCount, tt.lastError)
			assert.Equal(t, tt.expected, result)

			// Reactivate for next test
			if !tt.isActive {
				err := config.Activate()
				assert.NoError(t, err)
			}
		})
	}
}

func TestRetryConfiguration_UpdateConfiguration(t *testing.T) {
	config, err := domain.NewRetryConfiguration(
		3,
		30*time.Second,
		5*time.Minute,
		10*time.Minute,
		2.0,
		true,
		true,
	)
	assert.NoError(t, err)

	// Test valid update
	err = config.UpdateConfiguration(
		5,
		1*time.Minute,
		10*time.Minute,
		20*time.Minute,
		1.5,
		false,
		false,
	)
	assert.NoError(t, err)
	assert.Equal(t, 5, config.MaxRetryAttempts())
	assert.Equal(t, 1*time.Minute, config.InitialRetryDelay())
	assert.Equal(t, 10*time.Minute, config.MaxRetryDelay())
	assert.Equal(t, 1.5, config.RetryDelayMultiplier())
	assert.False(t, config.EnableExponentialBackoff())
	assert.False(t, config.EnableDeadLetterQueue())

	// Test invalid update
	err = config.UpdateConfiguration(
		-1,
		1*time.Minute,
		10*time.Minute,
		20*time.Minute,
		1.5,
		false,
		false,
	)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "max retry attempts must be between 0 and 10")
}

func TestRetryConfiguration_ActivateDeactivate(t *testing.T) {
	config, err := domain.NewRetryConfiguration(
		3,
		30*time.Second,
		5*time.Minute,
		10*time.Minute,
		2.0,
		true,
		true,
	)
	assert.NoError(t, err)
	assert.True(t, config.IsActive())

	// Test deactivate
	err = config.Deactivate()
	assert.NoError(t, err)
	assert.False(t, config.IsActive())

	// Test deactivate already inactive
	err = config.Deactivate()
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "retry configuration is already inactive")

	// Test activate
	err = config.Activate()
	assert.NoError(t, err)
	assert.True(t, config.IsActive())

	// Test activate already active
	err = config.Activate()
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "retry configuration is already active")
}

func TestGetDefaultRetryPolicies(t *testing.T) {
	policies := domain.GetDefaultRetryPolicies()

	expectedChannels := []domain.NotificationChannel{
		domain.NotificationChannelTelegram,
		domain.NotificationChannelEmail,
		domain.NotificationChannelSlack,
		domain.NotificationChannelWebhook,
	}

	for _, channel := range expectedChannels {
		assert.Contains(t, policies, channel)
		policy := policies[channel]
		assert.Greater(t, policy.MaxAttempts(), 0)
		assert.Greater(t, policy.BaseDelay(), time.Duration(0))
		assert.Greater(t, policy.MaxDelay(), policy.BaseDelay())
		assert.GreaterOrEqual(t, policy.BackoffFactor(), 1.0)
	}
}

func TestGetDefaultRetryConfiguration(t *testing.T) {
	config := domain.GetDefaultRetryConfiguration()

	assert.NotNil(t, config)
	assert.Equal(t, 3, config.MaxRetryAttempts())
	assert.Equal(t, 30*time.Second, config.InitialRetryDelay())
	assert.Equal(t, 5*time.Minute, config.MaxRetryDelay())
	assert.Equal(t, 2.0, config.RetryDelayMultiplier())
	assert.True(t, config.EnableExponentialBackoff())
	assert.True(t, config.EnableDeadLetterQueue())
	assert.True(t, config.IsActive())
}

func TestRestoreRetryConfiguration(t *testing.T) {
	id := value_objects.NewID()
	now := value_objects.NewTimestamp()

	params := domain.RestoreRetryConfigurationParams{
		ID:                       id,
		MaxRetryAttempts:         5,
		InitialRetryDelay:        1 * time.Minute,
		MaxRetryDelay:            10 * time.Minute,
		RetryDelayMultiplier:     1.5,
		RetryTimeoutDuration:     30 * time.Minute,
		EnableExponentialBackoff: false,
		EnableDeadLetterQueue:    false,
		IsActive:                 false,
		CreatedAt:                now,
		UpdatedAt:                now,
	}

	config := domain.RestoreRetryConfiguration(params)
	assert.NotNil(t, config)
	assert.Equal(t, id, config.ID())
	assert.Equal(t, 5, config.MaxRetryAttempts())
	assert.Equal(t, 1*time.Minute, config.InitialRetryDelay())
	assert.Equal(t, 10*time.Minute, config.MaxRetryDelay())
	assert.Equal(t, 1.5, config.RetryDelayMultiplier())
	assert.False(t, config.EnableExponentialBackoff())
	assert.False(t, config.EnableDeadLetterQueue())
	assert.False(t, config.IsActive())
}
