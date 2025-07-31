package retry

import (
	"context"
	"fmt"
	"time"

	"github.com/dewisartika8/cicd-status-notifier-bot/internal/core/notification/domain"
	"github.com/dewisartika8/cicd-status-notifier-bot/internal/core/notification/dto"
	"github.com/dewisartika8/cicd-status-notifier-bot/internal/core/notification/port"
	"github.com/dewisartika8/cicd-status-notifier-bot/internal/core/shared/domain/value_objects"
	"github.com/sirupsen/logrus"
)

// Resource type constants
const (
	resourceRetryConfig = "retry configuration"
)

// Service operation messages
const (
	RetryConfigurationServiceName = "RetryConfigurationService"
)

type Dep struct {
	RetryRepo port.RetryConfigurationRepository
	Logger    *logrus.Logger
}

// retryService implements the RetryService interface
type retryService struct {
	Dep
}

// NewRetryService creates a new retry service instance
func NewRetryService(d Dep) port.RetryService {
	return &retryService{
		Dep: d,
	}
}

// CreateRetryConfiguration creates a new retry configuration
func (s *retryService) CreateRetryConfiguration(ctx context.Context, req dto.CreateRetryConfigurationRequest) (*domain.RetryConfiguration, error) {
	s.Logger.WithFields(logrus.Fields{
		"service": RetryConfigurationServiceName,
		"method":  "CreateRetryConfiguration",
		"request": req,
	}).Info("Creating retry configuration")

	// Create new retry configuration with validation
	config, err := domain.NewRetryConfiguration(
		req.MaxRetryAttempts,
		req.InitialRetryDelay,
		req.MaxRetryDelay,
		req.RetryTimeoutDuration,
		req.RetryDelayMultiplier,
		req.EnableExponentialBackoff,
		req.EnableDeadLetterQueue,
	)
	if err != nil {
		s.Logger.WithError(err).Error("Failed to create retry configuration entity")
		return nil, fmt.Errorf(domain.ErrMsgCreate, resourceRetryConfig, err)
	}

	// Persist the configuration
	if err := s.RetryRepo.Create(ctx, config); err != nil {
		s.Logger.WithError(err).Error("Failed to persist retry configuration")
		return nil, fmt.Errorf(domain.ErrMsgPersist, resourceRetryConfig, err)
	}

	s.Logger.WithField("config_id", config.ID().String()).Info(domain.RetryConfigCreated)
	return config, nil
}

// GetRetryConfiguration retrieves a retry configuration by ID
func (s *retryService) GetRetryConfiguration(ctx context.Context, id value_objects.ID) (*domain.RetryConfiguration, error) {
	s.Logger.WithField("id", id.String()).Info("Getting retry configuration")

	config, err := s.RetryRepo.GetByID(ctx, id)
	if err != nil {
		s.Logger.WithError(err).Error(domain.LogMsgGetRetryConfig)
		return nil, fmt.Errorf(domain.ErrMsgGet, resourceRetryConfig, err)
	}

	if config == nil {
		s.Logger.Error("Retry configuration not found")
		return nil, domain.NewRetryConfigurationNotFoundError(id.String())
	}

	return config, nil
}

// GetRetryConfigurationByChannel retrieves retry configuration for a channel
func (s *retryService) GetRetryConfigurationByChannel(ctx context.Context, channel domain.NotificationChannel) (*domain.RetryConfiguration, error) {
	s.Logger.WithFields(logrus.Fields{
		"channel": channel,
	}).Info("Getting retry configuration by channel")

	config, err := s.RetryRepo.GetByChannel(ctx, channel)
	if err != nil {
		s.Logger.WithError(err).Error("Failed to get retry configuration by channel")
		return nil, fmt.Errorf("failed to get retry configuration by channel: %w", err)
	}

	// Return default configuration if no specific configuration found
	if config == nil {
		s.Logger.Info("No specific configuration found, returning default configuration")
		return domain.GetDefaultRetryConfiguration(), nil
	}

	return config, nil
}

// UpdateRetryConfiguration updates an existing retry configuration
func (s *retryService) UpdateRetryConfiguration(ctx context.Context, id value_objects.ID, req dto.UpdateRetryConfigurationRequest) (*domain.RetryConfiguration, error) {
	s.Logger.WithField("id", id.String()).Info("Updating retry configuration")

	// Get the configuration
	config, err := s.RetryRepo.GetByID(ctx, id)
	if err != nil {
		s.Logger.WithError(err).Error(domain.LogMsgGetRetryConfig)
		return nil, fmt.Errorf(domain.ErrMsgGet, resourceRetryConfig, err)
	}

	// Create updated configuration with new values
	updatedParams := domain.RestoreRetryConfigurationParams{
		ID:                       config.ID(),
		MaxRetryAttempts:         config.MaxRetryAttempts(),
		InitialRetryDelay:        config.InitialRetryDelay(),
		MaxRetryDelay:            config.MaxRetryDelay(),
		RetryDelayMultiplier:     config.RetryDelayMultiplier(),
		RetryTimeoutDuration:     config.RetryTimeoutDuration(),
		EnableExponentialBackoff: config.EnableExponentialBackoff(),
		EnableDeadLetterQueue:    config.EnableDeadLetterQueue(),
		IsActive:                 config.IsActive(),
		CreatedAt:                config.CreatedAt(),
		UpdatedAt:                value_objects.NewTimestamp(),
	}

	// Apply updates
	if req.MaxRetryAttempts != nil {
		updatedParams.MaxRetryAttempts = *req.MaxRetryAttempts
	}
	if req.InitialRetryDelay != nil {
		updatedParams.InitialRetryDelay = *req.InitialRetryDelay
	}
	if req.MaxRetryDelay != nil {
		updatedParams.MaxRetryDelay = *req.MaxRetryDelay
	}
	if req.RetryDelayMultiplier != nil {
		updatedParams.RetryDelayMultiplier = *req.RetryDelayMultiplier
	}
	if req.RetryTimeoutDuration != nil {
		updatedParams.RetryTimeoutDuration = *req.RetryTimeoutDuration
	}
	if req.EnableExponentialBackoff != nil {
		updatedParams.EnableExponentialBackoff = *req.EnableExponentialBackoff
	}
	if req.EnableDeadLetterQueue != nil {
		updatedParams.EnableDeadLetterQueue = *req.EnableDeadLetterQueue
	}

	// Create updated configuration
	updatedConfig := domain.RestoreRetryConfiguration(updatedParams)

	// Update the configuration in repository
	if err := s.RetryRepo.Update(ctx, updatedConfig); err != nil {
		s.Logger.WithError(err).Error(domain.LogMsgUpdateRetryConfig)
		return nil, fmt.Errorf(domain.ErrMsgUpdate, resourceRetryConfig, err)
	}

	s.Logger.Info(domain.RetryConfigUpdated)
	return updatedConfig, nil
}

// ActivateRetryConfiguration activates a retry configuration
func (s *retryService) ActivateRetryConfiguration(ctx context.Context, id value_objects.ID) error {
	s.Logger.WithField("id", id.String()).Info("Activating retry configuration")

	config, err := s.RetryRepo.GetByID(ctx, id)
	if err != nil {
		s.Logger.WithError(err).Error(domain.LogMsgGetRetryConfig)
		return fmt.Errorf(domain.ErrMsgGet, resourceRetryConfig, err)
	}

	config.Activate()

	if err := s.RetryRepo.Update(ctx, config); err != nil {
		s.Logger.WithError(err).Error(domain.LogMsgUpdateRetryConfig)
		return fmt.Errorf(domain.ErrMsgUpdate, resourceRetryConfig, err)
	}

	s.Logger.Info(domain.RetryConfigActivated)
	return nil
}

// DeactivateRetryConfiguration deactivates a retry configuration
func (s *retryService) DeactivateRetryConfiguration(ctx context.Context, id value_objects.ID) error {
	s.Logger.WithField("id", id.String()).Info("Deactivating retry configuration")

	config, err := s.RetryRepo.GetByID(ctx, id)
	if err != nil {
		s.Logger.WithError(err).Error(domain.LogMsgGetRetryConfig)
		return fmt.Errorf(domain.ErrMsgGet, resourceRetryConfig, err)
	}

	config.Deactivate()

	if err := s.RetryRepo.Update(ctx, config); err != nil {
		s.Logger.WithError(err).Error(domain.LogMsgUpdateRetryConfig)
		return fmt.Errorf(domain.ErrMsgUpdate, resourceRetryConfig, err)
	}

	s.Logger.Info(domain.RetryConfigDeactivated)
	return nil
}

// DeleteRetryConfiguration deletes a retry configuration
func (s *retryService) DeleteRetryConfiguration(ctx context.Context, id value_objects.ID) error {
	s.Logger.WithField("id", id.String()).Info("Deleting retry configuration")

	if err := s.RetryRepo.Delete(ctx, id); err != nil {
		s.Logger.WithError(err).Error(domain.LogMsgDeleteRetryConfig)
		return fmt.Errorf(domain.ErrMsgDelete, resourceRetryConfig, err)
	}

	s.Logger.Info(domain.RetryConfigDeleted)
	return nil
}

// ListActiveRetryConfigurations lists all active retry configurations
func (s *retryService) ListActiveRetryConfigurations(ctx context.Context) ([]*domain.RetryConfiguration, error) {
	s.Logger.Info("Getting active retry configurations")

	configs, err := s.RetryRepo.GetActiveConfigurations(ctx)
	if err != nil {
		s.Logger.WithError(err).Error("Failed to get active retry configurations")
		return nil, fmt.Errorf("failed to get active retry configurations: %w", err)
	}

	return configs, nil
}

// InitializeDefaultRetryConfigurations sets up default retry configurations
func (s *retryService) InitializeDefaultRetryConfigurations(ctx context.Context) error {
	s.Logger.Info("Initializing default retry configurations")

	// Define default configurations for each channel
	defaultConfigs := []domain.RestoreRetryConfigurationParams{
		{
			ID:                       value_objects.NewID(),
			MaxRetryAttempts:         3,
			InitialRetryDelay:        time.Second * 5,
			MaxRetryDelay:            time.Minute * 5,
			RetryDelayMultiplier:     2.0,
			RetryTimeoutDuration:     time.Minute * 10,
			EnableExponentialBackoff: true,
			EnableDeadLetterQueue:    false,
			IsActive:                 true,
			CreatedAt:                value_objects.NewTimestamp(),
			UpdatedAt:                value_objects.NewTimestamp(),
		},
		{
			ID:                       value_objects.NewID(),
			MaxRetryAttempts:         5,
			InitialRetryDelay:        time.Second * 10,
			MaxRetryDelay:            time.Minute * 10,
			RetryDelayMultiplier:     1.5,
			RetryTimeoutDuration:     time.Minute * 15,
			EnableExponentialBackoff: true,
			EnableDeadLetterQueue:    false,
			IsActive:                 true,
			CreatedAt:                value_objects.NewTimestamp(),
			UpdatedAt:                value_objects.NewTimestamp(),
		},
		{
			ID:                       value_objects.NewID(),
			MaxRetryAttempts:         3,
			InitialRetryDelay:        time.Second * 3,
			MaxRetryDelay:            time.Minute * 3,
			RetryDelayMultiplier:     2.0,
			RetryTimeoutDuration:     time.Minute * 5,
			EnableExponentialBackoff: true,
			EnableDeadLetterQueue:    false,
			IsActive:                 true,
			CreatedAt:                value_objects.NewTimestamp(),
			UpdatedAt:                value_objects.NewTimestamp(),
		},
		{
			ID:                       value_objects.NewID(),
			MaxRetryAttempts:         2,
			InitialRetryDelay:        time.Second * 2,
			MaxRetryDelay:            time.Minute * 2,
			RetryDelayMultiplier:     2.0,
			RetryTimeoutDuration:     time.Minute * 3,
			EnableExponentialBackoff: true,
			EnableDeadLetterQueue:    false,
			IsActive:                 true,
			CreatedAt:                value_objects.NewTimestamp(),
			UpdatedAt:                value_objects.NewTimestamp(),
		},
	}

	// Create configurations
	var configs []*domain.RetryConfiguration
	for _, params := range defaultConfigs {
		config := domain.RestoreRetryConfiguration(params)
		configs = append(configs, config)
	}

	// Bulk create configurations
	if err := s.RetryRepo.BulkCreate(ctx, configs); err != nil {
		s.Logger.WithError(err).Error("Failed to persist default retry configurations")
		return fmt.Errorf("failed to create default retry configurations: %w", err)
	}

	s.Logger.Info(domain.DefaultRetryConfigsInitialized)
	return nil
}

// CalculateRetryDelay calculates the delay for a retry attempt
func (s *retryService) CalculateRetryDelay(ctx context.Context, channel domain.NotificationChannel, attemptNumber int) (time.Duration, error) {
	s.Logger.WithFields(logrus.Fields{
		"channel": channel,
		"attempt": attemptNumber,
	}).Info(domain.LogMsgCalculatingDelay)

	// Get retry configuration
	config, err := s.RetryRepo.GetByChannel(ctx, channel)
	if err != nil {
		s.Logger.WithError(err).Error("Failed to get retry configuration for delay calculation")
		return 0, fmt.Errorf(domain.ErrMsgCalculateRetryDelay, err)
	}

	// Calculate delay using exponential backoff
	delay := config.CalculateRetryDelay(attemptNumber)

	s.Logger.WithFields(logrus.Fields{
		"channel": channel,
		"attempt": attemptNumber,
		"delay":   delay,
	}).Info(domain.RetryDelayCalculated)

	return delay, nil
}

// ShouldRetryNotification determines if a notification should be retried
func (s *retryService) ShouldRetryNotification(ctx context.Context, channel domain.NotificationChannel, attemptCount int, lastError error) (bool, error) {
	s.Logger.WithFields(logrus.Fields{
		"channel":       channel,
		"attempt_count": attemptCount,
		"last_error":    lastError.Error(),
	}).Info(domain.LogMsgRetryDecision)

	// Get retry configuration
	config, err := s.RetryRepo.GetByChannel(ctx, channel)
	if err != nil {
		s.Logger.WithError(err).Error("Failed to get retry configuration for retry decision")
		return false, fmt.Errorf(domain.ErrMsgGet, resourceRetryConfig, err)
	}

	// Check if should retry
	shouldRetry := config.ShouldRetry(attemptCount, lastError)

	s.Logger.WithFields(logrus.Fields{
		"channel":      channel,
		"attempt":      attemptCount,
		"should_retry": shouldRetry,
	}).Info(domain.RetryDecisionMade)

	return shouldRetry, nil
}

// ProcessRetryableNotification processes a notification that can be retried
func (s *retryService) ProcessRetryableNotification(ctx context.Context, req dto.ProcessRetryableNotificationRequest) (*dto.ProcessRetryableNotificationResponse, error) {
	s.Logger.WithFields(logrus.Fields{
		"notification_id": req.NotificationID.String(),
		"channel":         req.Channel,
		"attempt_count":   req.AttemptCount,
	}).Info(domain.LogMsgProcessingRetryable)

	// Validate retry attempt
	if req.AttemptCount < 1 {
		err := fmt.Errorf("invalid retry attempt: %d", req.AttemptCount)
		s.Logger.WithError(err).Error("Invalid retry attempt")
		return nil, fmt.Errorf(domain.ErrMsgInvalidRetryAttempt, err)
	}

	// Check if should retry
	shouldRetry, err := s.ShouldRetryNotification(ctx, req.Channel, req.AttemptCount, req.LastError)
	if err != nil {
		s.Logger.WithError(err).Error("Failed to determine if should retry")
		return nil, fmt.Errorf(domain.ErrMsgProcessRetryable, err)
	}

	response := &dto.ProcessRetryableNotificationResponse{
		ShouldRetry:      shouldRetry,
		SendToDeadLetter: false,
	}

	if shouldRetry {
		// Calculate retry delay
		delay, err := s.CalculateRetryDelay(ctx, req.Channel, req.AttemptCount)
		if err != nil {
			s.Logger.WithError(err).Error("Failed to calculate retry delay")
			return nil, fmt.Errorf(domain.ErrMsgProcessRetryable, err)
		}

		response.RetryDelay = delay
		nextAttemptAt := time.Now().Add(delay)
		response.NextAttemptAt = &nextAttemptAt

		// Get retry configuration for response
		config, err := s.RetryRepo.GetByChannel(ctx, req.Channel)
		if err == nil {
			response.RetryConfiguration = config
		}
	} else {
		// If cannot retry, send to dead letter queue if enabled
		config, err := s.RetryRepo.GetByChannel(ctx, req.Channel)
		if err == nil && config.EnableDeadLetterQueue() {
			response.SendToDeadLetter = true
		}
	}

	// Log the retry processing
	s.Logger.WithFields(logrus.Fields{
		"notification_id": req.NotificationID.String(),
		"channel":         req.Channel,
		"attempt_count":   req.AttemptCount,
		"should_retry":    response.ShouldRetry,
		"retry_delay":     response.RetryDelay,
	}).Info(domain.RetryableNotificationProcessed)

	return response, nil
}
