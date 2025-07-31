package domain

import (
	"strings"
	"time"

	"github.com/dewisartika8/cicd-status-notifier-bot/internal/core/shared/domain/value_objects"
)

// NotificationStatus represents the status of a notification
type NotificationStatus value_objects.Status

const (
	NotificationStatusPending   NotificationStatus = "pending"
	NotificationStatusSent      NotificationStatus = "sent"
	NotificationStatusDelivered NotificationStatus = "delivered"
	NotificationStatusFailed    NotificationStatus = "failed"
	NotificationStatusRetrying  NotificationStatus = "retrying"
	NotificationStatusCancelled NotificationStatus = "cancelled"
	NotificationStatusExpired   NotificationStatus = "expired"
)

// IsValid checks if the notification status is valid
func (s NotificationStatus) IsValid() bool {
	switch s {
	case NotificationStatusPending, NotificationStatusSent, NotificationStatusDelivered,
		NotificationStatusFailed, NotificationStatusRetrying, NotificationStatusCancelled,
		NotificationStatusExpired:
		return true
	default:
		return false
	}
}

// NotificationChannel represents the notification channel type
type NotificationChannel value_objects.Status

const (
	NotificationChannelTelegram NotificationChannel = "telegram"
	NotificationChannelEmail    NotificationChannel = "email"
	NotificationChannelSlack    NotificationChannel = "slack"
	NotificationChannelWebhook  NotificationChannel = "webhook"
)

// IsValid checks if the notification channel is valid
func (c NotificationChannel) IsValid() bool {
	switch c {
	case NotificationChannelTelegram, NotificationChannelEmail, NotificationChannelSlack, NotificationChannelWebhook:
		return true
	default:
		return false
	}
}

// IsValidNotificationChannel is a helper function for external validation
func IsValidNotificationChannel(channel NotificationChannel) bool {
	return channel.IsValid()
}

// NotificationMetrics represents metrics tracking for notifications
type NotificationMetrics struct {
	deliveryAttempts    int
	totalRetries        int
	averageDeliveryTime time.Duration
	lastAttemptAt       *value_objects.Timestamp
	firstAttemptAt      *value_objects.Timestamp
	deliveredAt         *value_objects.Timestamp
	failedAt            *value_objects.Timestamp
}

// NewNotificationMetrics creates a new notification metrics instance
func NewNotificationMetrics() *NotificationMetrics {
	return &NotificationMetrics{
		deliveryAttempts:    0,
		totalRetries:        0,
		averageDeliveryTime: 0,
	}
}

// Getters for NotificationMetrics
func (nm *NotificationMetrics) DeliveryAttempts() int {
	return nm.deliveryAttempts
}

func (nm *NotificationMetrics) TotalRetries() int {
	return nm.totalRetries
}

func (nm *NotificationMetrics) AverageDeliveryTime() time.Duration {
	return nm.averageDeliveryTime
}

func (nm *NotificationMetrics) LastAttemptAt() *value_objects.Timestamp {
	return nm.lastAttemptAt
}

func (nm *NotificationMetrics) FirstAttemptAt() *value_objects.Timestamp {
	return nm.firstAttemptAt
}

func (nm *NotificationMetrics) DeliveredAt() *value_objects.Timestamp {
	return nm.deliveredAt
}

func (nm *NotificationMetrics) FailedAt() *value_objects.Timestamp {
	return nm.failedAt
}

// Business methods for NotificationMetrics
func (nm *NotificationMetrics) RecordAttempt() {
	nm.deliveryAttempts++
	now := value_objects.NewTimestamp()
	nm.lastAttemptAt = &now

	if nm.firstAttemptAt == nil {
		nm.firstAttemptAt = &now
	}
}

func (nm *NotificationMetrics) RecordRetry() {
	nm.totalRetries++
	nm.RecordAttempt()
}

func (nm *NotificationMetrics) RecordDelivery() {
	now := value_objects.NewTimestamp()
	nm.deliveredAt = &now

	if nm.firstAttemptAt != nil {
		deliveryTime := now.ToTime().Sub(nm.firstAttemptAt.ToTime())
		nm.averageDeliveryTime = deliveryTime
	}
}

func (nm *NotificationMetrics) RecordFailure() {
	now := value_objects.NewTimestamp()
	nm.failedAt = &now
}

// NotificationLog represents a notification log domain entity
type NotificationLog struct {
	id           value_objects.ID
	buildEventID value_objects.ID
	projectID    value_objects.ID
	channel      NotificationChannel
	recipient    string
	message      string
	status       NotificationStatus
	errorMessage string
	retryCount   int
	maxRetries   int
	messageID    *string // For storing external message ID (e.g., Telegram message ID)
	templateID   *value_objects.ID
	metadata     map[string]interface{}
	metrics      *NotificationMetrics
	nextRetryAt  *value_objects.Timestamp
	expiresAt    *value_objects.Timestamp
	sentAt       *value_objects.Timestamp
	createdAt    value_objects.Timestamp
	updatedAt    value_objects.Timestamp
}

// NewNotificationLog creates a new notification log entity
func NewNotificationLog(
	buildEventID, projectID value_objects.ID,
	channel NotificationChannel,
	recipient, message string,
	maxRetries int,
) (*NotificationLog, error) {
	// Initialize metadata
	metadata := make(map[string]interface{})

	notificationLog := &NotificationLog{
		id:           value_objects.NewID(),
		buildEventID: buildEventID,
		projectID:    projectID,
		channel:      channel,
		recipient:    strings.TrimSpace(recipient),
		message:      strings.TrimSpace(message),
		status:       NotificationStatusPending,
		retryCount:   0,
		maxRetries:   maxRetries,
		metadata:     metadata,
		metrics:      NewNotificationMetrics(),
		createdAt:    value_objects.NewTimestamp(),
		updatedAt:    value_objects.NewTimestamp(),
	}

	if err := notificationLog.validate(); err != nil {
		return nil, err
	}

	return notificationLog, nil
}

// RestoreNotificationLog restores a notification log from persistence
func RestoreNotificationLog(params RestoreNotificationLogParams) *NotificationLog {
	return &NotificationLog{
		id:           params.ID,
		buildEventID: params.BuildEventID,
		projectID:    params.ProjectID,
		channel:      params.Channel,
		recipient:    params.Recipient,
		message:      params.Message,
		status:       params.Status,
		errorMessage: params.ErrorMessage,
		retryCount:   params.RetryCount,
		maxRetries:   params.MaxRetries,
		messageID:    params.MessageID,
		templateID:   params.TemplateID,
		metadata:     params.Metadata,
		metrics:      params.Metrics,
		nextRetryAt:  params.NextRetryAt,
		expiresAt:    params.ExpiresAt,
		sentAt:       params.SentAt,
		createdAt:    params.CreatedAt,
		updatedAt:    params.UpdatedAt,
	}
}

// RestoreNotificationLogParams holds parameters for restoring a notification log
type RestoreNotificationLogParams struct {
	ID           value_objects.ID
	BuildEventID value_objects.ID
	ProjectID    value_objects.ID
	Channel      NotificationChannel
	Recipient    string
	Message      string
	Status       NotificationStatus
	ErrorMessage string
	RetryCount   int
	MaxRetries   int
	MessageID    *string
	TemplateID   *value_objects.ID
	Metadata     map[string]interface{}
	Metrics      *NotificationMetrics
	NextRetryAt  *value_objects.Timestamp
	ExpiresAt    *value_objects.Timestamp
	SentAt       *value_objects.Timestamp
	CreatedAt    value_objects.Timestamp
	UpdatedAt    value_objects.Timestamp
}

// ID returns the notification log ID
func (nl *NotificationLog) ID() value_objects.ID {
	return nl.id
}

// BuildEventID returns the build event ID
func (nl *NotificationLog) BuildEventID() value_objects.ID {
	return nl.buildEventID
}

// ProjectID returns the project ID
func (nl *NotificationLog) ProjectID() value_objects.ID {
	return nl.projectID
}

// Channel returns the notification channel
func (nl *NotificationLog) Channel() NotificationChannel {
	return nl.channel
}

// Recipient returns the notification recipient
func (nl *NotificationLog) Recipient() string {
	return nl.recipient
}

// Message returns the notification message
func (nl *NotificationLog) Message() string {
	return nl.message
}

// Status returns the notification status
func (nl *NotificationLog) Status() NotificationStatus {
	return nl.status
}

// ErrorMessage returns the error message
func (nl *NotificationLog) ErrorMessage() string {
	return nl.errorMessage
}

// RetryCount returns the retry count
func (nl *NotificationLog) RetryCount() int {
	return nl.retryCount
}

// MaxRetries returns the maximum retry count
func (nl *NotificationLog) MaxRetries() int {
	return nl.maxRetries
}

// MessageID returns the external message ID
func (nl *NotificationLog) MessageID() *string {
	return nl.messageID
}

// TemplateID returns the template ID
func (nl *NotificationLog) TemplateID() *value_objects.ID {
	return nl.templateID
}

// Metadata returns a copy of the metadata
func (nl *NotificationLog) Metadata() map[string]interface{} {
	metadataCopy := make(map[string]interface{})
	for k, v := range nl.metadata {
		metadataCopy[k] = v
	}
	return metadataCopy
}

// Metrics returns the notification metrics
func (nl *NotificationLog) Metrics() *NotificationMetrics {
	return nl.metrics
}

// NextRetryAt returns the next retry timestamp
func (nl *NotificationLog) NextRetryAt() *value_objects.Timestamp {
	return nl.nextRetryAt
}

// ExpiresAt returns the expiration timestamp
func (nl *NotificationLog) ExpiresAt() *value_objects.Timestamp {
	return nl.expiresAt
}

// SentAt returns the sent timestamp
func (nl *NotificationLog) SentAt() *value_objects.Timestamp {
	return nl.sentAt
}

// CreatedAt returns the creation timestamp
func (nl *NotificationLog) CreatedAt() value_objects.Timestamp {
	return nl.createdAt
}

// UpdatedAt returns the last update timestamp
func (nl *NotificationLog) UpdatedAt() value_objects.Timestamp {
	return nl.updatedAt
}

// MarkAsSent marks the notification as successfully sent
func (nl *NotificationLog) MarkAsSent(messageID *string) error {
	if nl.status == NotificationStatusSent {
		return nil // Already sent, no-op
	}

	nl.status = NotificationStatusSent
	nl.messageID = messageID
	ts := value_objects.NewTimestamp()
	nl.sentAt = &ts
	nl.updatedAt = value_objects.NewTimestamp()
	nl.errorMessage = "" // Clear any previous error

	// Record metrics
	nl.metrics.RecordAttempt()

	return nil
}

// MarkAsDelivered marks the notification as delivered
func (nl *NotificationLog) MarkAsDelivered() error {
	nl.status = NotificationStatusDelivered
	nl.updatedAt = value_objects.NewTimestamp()
	nl.errorMessage = "" // Clear any previous error

	// Record delivery metrics
	nl.metrics.RecordDelivery()

	return nil
}

// MarkAsFailed marks the notification as failed with an error message
func (nl *NotificationLog) MarkAsFailed(errorMessage string) error {
	if strings.TrimSpace(errorMessage) == "" {
		return ErrInvalidMessage
	}

	nl.status = NotificationStatusFailed
	nl.errorMessage = strings.TrimSpace(errorMessage)
	nl.updatedAt = value_objects.NewTimestamp()

	// Record failure metrics
	nl.metrics.RecordFailure()

	return nil
}

// MarkAsRetrying marks the notification for retry
func (nl *NotificationLog) MarkAsRetrying() error {
	if nl.retryCount >= nl.maxRetries {
		return ErrMaxRetryAttemptsExceeded
	}

	nl.status = NotificationStatusRetrying
	nl.retryCount++
	nl.updatedAt = value_objects.NewTimestamp()

	// Record retry metrics
	nl.metrics.RecordRetry()

	return nil
}

// MarkAsExpired marks the notification as expired
func (nl *NotificationLog) MarkAsExpired() error {
	nl.status = NotificationStatusExpired
	nl.updatedAt = value_objects.NewTimestamp()
	nl.ClearRetrySchedule()

	return nil
}

// MarkAsCancelled marks the notification as cancelled
func (nl *NotificationLog) MarkAsCancelled() error {
	nl.status = NotificationStatusCancelled
	nl.updatedAt = value_objects.NewTimestamp()
	nl.ClearRetrySchedule()

	return nil
}

// CanRetry checks if the notification can be retried
func (nl *NotificationLog) CanRetry() bool {
	return nl.status == NotificationStatusFailed &&
		nl.retryCount < nl.maxRetries &&
		!nl.IsExpired()
}

// ScheduleRetry schedules the next retry attempt
func (nl *NotificationLog) ScheduleRetry(retryAt value_objects.Timestamp) error {
	if nl.retryCount >= nl.maxRetries {
		return NewMaxRetriesExceededError(nl.maxRetries)
	}

	nl.nextRetryAt = &retryAt
	nl.updatedAt = value_objects.NewTimestamp()

	return nil
}

// ClearRetrySchedule clears the retry schedule
func (nl *NotificationLog) ClearRetrySchedule() {
	nl.nextRetryAt = nil
	nl.updatedAt = value_objects.NewTimestamp()
}

// SetExpiration sets the expiration time for the notification
func (nl *NotificationLog) SetExpiration(expiresAt value_objects.Timestamp) {
	nl.expiresAt = &expiresAt
	nl.updatedAt = value_objects.NewTimestamp()
}

// IsExpired checks if the notification has expired
func (nl *NotificationLog) IsExpired() bool {
	if nl.expiresAt == nil {
		return false
	}

	return time.Now().After(nl.expiresAt.ToTime())
}

// SetTemplateID sets the template ID
func (nl *NotificationLog) SetTemplateID(templateID value_objects.ID) {
	nl.templateID = &templateID
	nl.updatedAt = value_objects.NewTimestamp()
}

// UpdateMetadata updates specific metadata fields
func (nl *NotificationLog) UpdateMetadata(key string, value interface{}) {
	if nl.metadata == nil {
		nl.metadata = make(map[string]interface{})
	}

	nl.metadata[key] = value
	nl.updatedAt = value_objects.NewTimestamp()
}

// RemoveMetadata removes a metadata field
func (nl *NotificationLog) RemoveMetadata(key string) {
	if nl.metadata != nil {
		delete(nl.metadata, key)
		nl.updatedAt = value_objects.NewTimestamp()
	}
}

// GetMetadataValue gets a specific metadata value
func (nl *NotificationLog) GetMetadataValue(key string) (interface{}, bool) {
	if nl.metadata == nil {
		return nil, false
	}

	value, exists := nl.metadata[key]
	return value, exists
}

// UpdateMessage updates the notification message
func (nl *NotificationLog) UpdateMessage(message string) error {
	message = strings.TrimSpace(message)
	if message == "" {
		return ErrInvalidMessage
	}

	nl.message = message
	nl.updatedAt = value_objects.NewTimestamp()

	return nil
}

// validate validates the notification log entity
func (nl *NotificationLog) validate() error {
	if nl.buildEventID.IsNil() {
		return ErrInvalidNotificationLog
	}

	if nl.projectID.IsNil() {
		return ErrInvalidNotificationLog
	}

	if !nl.channel.IsValid() {
		return ErrInvalidNotificationChannel
	}

	if strings.TrimSpace(nl.recipient) == "" {
		return ErrInvalidRecipient
	}

	if strings.TrimSpace(nl.message) == "" {
		return ErrInvalidMessage
	}

	if !nl.status.IsValid() {
		return ErrInvalidNotificationStatus
	}

	// Channel-specific validation
	if err := nl.validateChannelSpecific(); err != nil {
		return err
	}

	return nil
}

// validateChannelSpecific performs channel-specific validation
func (nl *NotificationLog) validateChannelSpecific() error {
	switch nl.channel {
	case NotificationChannelTelegram:
		// Telegram recipient should be a chat ID (numeric string) or username
		if nl.recipient == "" {
			return NewInvalidRecipientError("telegram chat ID cannot be empty")
		}
	case NotificationChannelEmail:
		// Basic email validation (could be enhanced)
		if !strings.Contains(nl.recipient, "@") {
			return NewInvalidRecipientError("invalid email format")
		}
	case NotificationChannelSlack:
		// Slack channel should start with # or be a user ID
		if nl.recipient == "" {
			return NewInvalidRecipientError("slack channel cannot be empty")
		}
	case NotificationChannelWebhook:
		// Webhook URL validation (basic)
		if !strings.HasPrefix(nl.recipient, "http") {
			return NewInvalidRecipientError("webhook URL must start with http")
		}
	}

	return nil
}

// NotificationStats represents notification statistics for a project
type NotificationStats struct {
	ProjectID           value_objects.ID
	TotalNotifications  int64
	StatusCounts        map[NotificationStatus]int64
	ChannelCounts       map[NotificationChannel]int64
	SuccessRate         float64
	FailureRate         float64
	RetryRate           float64
	AverageDeliveryTime time.Duration
	GeneratedAt         value_objects.Timestamp
}

// NewNotificationStats creates a new notification stats instance
func NewNotificationStats(projectID value_objects.ID) *NotificationStats {
	return &NotificationStats{
		ProjectID:     projectID,
		StatusCounts:  make(map[NotificationStatus]int64),
		ChannelCounts: make(map[NotificationChannel]int64),
		GeneratedAt:   value_objects.NewTimestamp(),
	}
}

// UpdateStatusCount updates the count for a specific status
func (s *NotificationStats) UpdateStatusCount(status NotificationStatus, count int64) {
	s.StatusCounts[status] = count
	s.calculateTotalAndRates()
}

// UpdateChannelCount updates the count for a specific channel
func (s *NotificationStats) UpdateChannelCount(channel NotificationChannel, count int64) {
	s.ChannelCounts[channel] = count
}

// SetAverageDeliveryTime sets the average delivery time
func (s *NotificationStats) SetAverageDeliveryTime(duration time.Duration) {
	s.AverageDeliveryTime = duration
}

// calculateTotalAndRates calculates total notifications and rates
func (s *NotificationStats) calculateTotalAndRates() {
	s.TotalNotifications = 0
	for _, count := range s.StatusCounts {
		s.TotalNotifications += count
	}

	if s.TotalNotifications > 0 {
		successCount := s.StatusCounts[NotificationStatusSent] + s.StatusCounts[NotificationStatusDelivered]
		failureCount := s.StatusCounts[NotificationStatusFailed] + s.StatusCounts[NotificationStatusExpired]
		retryCount := s.StatusCounts[NotificationStatusRetrying]

		s.SuccessRate = float64(successCount) / float64(s.TotalNotifications) * 100
		s.FailureRate = float64(failureCount) / float64(s.TotalNotifications) * 100
		s.RetryRate = float64(retryCount) / float64(s.TotalNotifications) * 100
	}
}

// GetSuccessCount returns the total count of successful notifications
func (s *NotificationStats) GetSuccessCount() int64 {
	return s.StatusCounts[NotificationStatusSent] + s.StatusCounts[NotificationStatusDelivered]
}

// GetFailureCount returns the total count of failed notifications
func (s *NotificationStats) GetFailureCount() int64 {
	return s.StatusCounts[NotificationStatusFailed] + s.StatusCounts[NotificationStatusExpired]
}

// GetPendingCount returns the total count of pending notifications
func (s *NotificationStats) GetPendingCount() int64 {
	return s.StatusCounts[NotificationStatusPending]
}

// GetRetryingCount returns the total count of retrying notifications
func (s *NotificationStats) GetRetryingCount() int64 {
	return s.StatusCounts[NotificationStatusRetrying]
}

// RestoreNotificationMetricsParams represents parameters for restoring notification metrics
type RestoreNotificationMetricsParams struct {
	DeliveryAttempts    int
	TotalRetries        int
	AverageDeliveryTime time.Duration
	LastAttemptAt       *value_objects.Timestamp
	FirstAttemptAt      *value_objects.Timestamp
	DeliveredAt         *value_objects.Timestamp
	FailedAt            *value_objects.Timestamp
}

// RestoreNotificationMetrics restores notification metrics from persistence
func RestoreNotificationMetrics(params RestoreNotificationMetricsParams) *NotificationMetrics {
	return &NotificationMetrics{
		deliveryAttempts:    params.DeliveryAttempts,
		totalRetries:        params.TotalRetries,
		averageDeliveryTime: params.AverageDeliveryTime,
		lastAttemptAt:       params.LastAttemptAt,
		firstAttemptAt:      params.FirstAttemptAt,
		deliveredAt:         params.DeliveredAt,
		failedAt:            params.FailedAt,
	}
}
