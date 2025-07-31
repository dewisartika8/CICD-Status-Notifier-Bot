package dto

import (
	"time"

	"github.com/dewisartika8/cicd-status-notifier-bot/internal/core/notification/domain"
	"github.com/dewisartika8/cicd-status-notifier-bot/internal/core/shared/domain/value_objects"
)

// CreateNotificationLogRequest represents a request to create a notification log
type CreateNotificationLogRequest struct {
	BuildEventID string                     `json:"build_event_id" validate:"required,uuid"`
	ProjectID    string                     `json:"project_id" validate:"required,uuid"`
	Channel      domain.NotificationChannel `json:"channel" validate:"required,oneof=telegram email slack webhook"`
	Recipient    string                     `json:"recipient" validate:"required"`
	Message      string                     `json:"message" validate:"required"`
}

// SendNotificationRequest represents a request to send a notification
type SendNotificationRequest struct {
	NotificationLogID string `json:"notification_log_id" validate:"required,uuid"`
	MessageID         string `json:"message_id,omitempty"`
}

// UpdateNotificationStatusRequest represents a request to update notification status
type UpdateNotificationStatusRequest struct {
	Status       domain.NotificationStatus `json:"status" validate:"required,oneof=pending sent failed retrying"`
	ErrorMessage string                    `json:"error_message,omitempty"`
	MessageID    string                    `json:"message_id,omitempty"`
}

// ListNotificationLogsRequest represents a request to list notification logs
type ListNotificationLogsRequest struct {
	ProjectID    string                     `form:"project_id" validate:"omitempty,uuid"`
	BuildEventID string                     `form:"build_event_id" validate:"omitempty,uuid"`
	Channel      domain.NotificationChannel `form:"channel" validate:"omitempty,oneof=telegram email slack webhook"`
	Status       domain.NotificationStatus  `form:"status" validate:"omitempty,oneof=pending sent failed retrying"`
	Recipient    string                     `form:"recipient"`
	Page         int                        `form:"page" validate:"min=1" default:"1"`
	Limit        int                        `form:"limit" validate:"min=1,max=100" default:"20"`
	SortBy       string                     `form:"sort_by" validate:"omitempty,oneof=created_at updated_at sent_at"`
	SortOrder    string                     `form:"sort_order" validate:"omitempty,oneof=asc desc" default:"desc"`
}

// NotificationLogResponse represents a notification log response
type NotificationLogResponse struct {
	ID           string                     `json:"id"`
	BuildEventID string                     `json:"build_event_id"`
	ProjectID    string                     `json:"project_id"`
	Channel      domain.NotificationChannel `json:"channel"`
	Recipient    string                     `json:"recipient"`
	Message      string                     `json:"message"`
	Status       domain.NotificationStatus  `json:"status"`
	ErrorMessage string                     `json:"error_message,omitempty"`
	RetryCount   int                        `json:"retry_count"`
	MessageID    *string                    `json:"message_id,omitempty"`
	SentAt       *time.Time                 `json:"sent_at,omitempty"`
	CreatedAt    time.Time                  `json:"created_at"`
	UpdatedAt    time.Time                  `json:"updated_at"`
}

// NotificationLogListResponse represents a paginated list of notification logs
type NotificationLogListResponse struct {
	Data       []NotificationLogResponse `json:"data"`
	Total      int64                     `json:"total"`
	Page       int                       `json:"page"`
	Limit      int                       `json:"limit"`
	TotalPages int                       `json:"total_pages"`
}

// CreateTelegramSubscriptionRequest represents a request to create a telegram subscription
type CreateTelegramSubscriptionRequest struct {
	ProjectID string `json:"project_id" validate:"required,uuid"`
	ChatID    int64  `json:"chat_id" validate:"required"`
}

// UpdateTelegramSubscriptionRequest represents a request to update a telegram subscription
type UpdateTelegramSubscriptionRequest struct {
	ChatID   int64 `json:"chat_id,omitempty"`
	IsActive *bool `json:"is_active,omitempty"`
}

// TelegramSubscriptionResponse represents a telegram subscription response
type TelegramSubscriptionResponse struct {
	ID        string    `json:"id"`
	ProjectID string    `json:"project_id"`
	ChatID    int64     `json:"chat_id"`
	IsActive  bool      `json:"is_active"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// TelegramSubscriptionListResponse represents a paginated list of telegram subscriptions
type TelegramSubscriptionListResponse struct {
	Data       []TelegramSubscriptionResponse `json:"data"`
	Total      int64                          `json:"total"`
	Page       int                            `json:"page"`
	Limit      int                            `json:"limit"`
	TotalPages int                            `json:"total_pages"`
}

// RetryNotificationRequest represents a request to retry a failed notification
type RetryNotificationRequest struct {
	NotificationLogID string `json:"notification_log_id" validate:"required,uuid"`
}

// NotificationStatsResponse represents notification statistics
type NotificationStatsResponse struct {
	Total    int64 `json:"total"`
	Sent     int64 `json:"sent"`
	Failed   int64 `json:"failed"`
	Pending  int64 `json:"pending"`
	Retrying int64 `json:"retrying"`
}

// ToNotificationLogResponse converts domain entity to response DTO
func ToNotificationLogResponse(entity *domain.NotificationLog) NotificationLogResponse {
	response := NotificationLogResponse{
		ID:           entity.ID().String(),
		BuildEventID: entity.BuildEventID().String(),
		ProjectID:    entity.ProjectID().String(),
		Channel:      entity.Channel(),
		Recipient:    entity.Recipient(),
		Message:      entity.Message(),
		Status:       entity.Status(),
		ErrorMessage: entity.ErrorMessage(),
		RetryCount:   entity.RetryCount(),
		MessageID:    entity.MessageID(),
		CreatedAt:    entity.CreatedAt().ToTime(),
		UpdatedAt:    entity.UpdatedAt().ToTime(),
	}

	if entity.SentAt() != nil {
		sentAt := entity.SentAt().ToTime()
		response.SentAt = &sentAt
	}

	return response
}

// ToTelegramSubscriptionResponse converts domain entity to response DTO
func ToTelegramSubscriptionResponse(entity *domain.TelegramSubscription) TelegramSubscriptionResponse {
	return TelegramSubscriptionResponse{
		ID:        entity.ID().String(),
		ProjectID: entity.ProjectID().String(),
		ChatID:    entity.ChatID(),
		IsActive:  entity.IsActive(),
		CreatedAt: entity.CreatedAt().ToTime(),
		UpdatedAt: entity.UpdatedAt().ToTime(),
	}
}

// ToCreateNotificationLogParams converts request DTO to service parameters
func (req CreateNotificationLogRequest) ToCreateNotificationLogParams() (value_objects.ID, value_objects.ID, domain.NotificationChannel, string, string, error) {
	buildEventID, err := value_objects.NewIDFromString(req.BuildEventID)
	if err != nil {
		return value_objects.ID{}, value_objects.ID{}, "", "", "", err
	}

	projectID, err := value_objects.NewIDFromString(req.ProjectID)
	if err != nil {
		return value_objects.ID{}, value_objects.ID{}, "", "", "", err
	}

	return buildEventID, projectID, req.Channel, req.Recipient, req.Message, nil
}

// CreateRetryConfigurationRequest represents the request to create a retry configuration
type CreateRetryConfigurationRequest struct {
	MaxRetryAttempts         int                        `json:"max_retry_attempts" validate:"required,min=0,max=10"`
	InitialRetryDelay        time.Duration              `json:"initial_retry_delay" validate:"required"`
	MaxRetryDelay            time.Duration              `json:"max_retry_delay" validate:"required"`
	RetryTimeoutDuration     time.Duration              `json:"retry_timeout_duration" validate:"required"`
	RetryDelayMultiplier     float64                    `json:"retry_delay_multiplier" validate:"required,min=1.0"`
	EnableExponentialBackoff bool                       `json:"enable_exponential_backoff"`
	EnableDeadLetterQueue    bool                       `json:"enable_dead_letter_queue"`
	Channel                  domain.NotificationChannel `json:"channel,omitempty"`
}

// UpdateRetryConfigurationRequest represents the request to update a retry configuration
type UpdateRetryConfigurationRequest struct {
	MaxRetryAttempts         *int           `json:"max_retry_attempts,omitempty" validate:"omitempty,min=0,max=10"`
	InitialRetryDelay        *time.Duration `json:"initial_retry_delay,omitempty"`
	MaxRetryDelay            *time.Duration `json:"max_retry_delay,omitempty"`
	RetryTimeoutDuration     *time.Duration `json:"retry_timeout_duration,omitempty"`
	RetryDelayMultiplier     *float64       `json:"retry_delay_multiplier,omitempty" validate:"omitempty,min=1.0"`
	EnableExponentialBackoff *bool          `json:"enable_exponential_backoff,omitempty"`
	EnableDeadLetterQueue    *bool          `json:"enable_dead_letter_queue,omitempty"`
}

// ProcessRetryableNotificationRequest represents a request to process a retryable notification
type ProcessRetryableNotificationRequest struct {
	NotificationID  value_objects.ID           `json:"notification_id" validate:"required"`
	Channel         domain.NotificationChannel `json:"channel" validate:"required"`
	AttemptCount    int                        `json:"attempt_count" validate:"required,min=1"`
	LastError       error                      `json:"last_error"`
	OriginalPayload map[string]interface{}     `json:"original_payload"`
	ScheduledAt     *time.Time                 `json:"scheduled_at,omitempty"`
}

// ProcessRetryableNotificationResponse represents the response from processing a retryable notification
type ProcessRetryableNotificationResponse struct {
	ShouldRetry        bool                       `json:"should_retry"`
	NextAttemptAt      *time.Time                 `json:"next_attempt_at,omitempty"`
	RetryDelay         time.Duration              `json:"retry_delay"`
	SendToDeadLetter   bool                       `json:"send_to_dead_letter"`
	RetryConfiguration *domain.RetryConfiguration `json:"retry_configuration"`
}
