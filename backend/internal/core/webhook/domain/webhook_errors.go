package domain

import "github.com/dewisartika8/cicd-status-notifier-bot/pkg/exception"

const (
	// Error codes for webhook domain
	WebhookErrInvalidSignature = "WEBHOOK_INVALID_SIGNATURE"
	WebhookErrInvalidPayload   = "WEBHOOK_INVALID_PAYLOAD"
	WebhookErrInvalidEvent     = "WEBHOOK_INVALID_EVENT"
	WebhookErrProjectNotFound  = "WEBHOOK_PROJECT_NOT_FOUND"
	WebhookErrProcessingFailed = "WEBHOOK_PROCESSING_FAILED"
)

// Webhook-specific domain errors
var (
	ErrWebhookInvalidSignature = exception.NewDomainError(
		WebhookErrInvalidSignature,
		"webhook signature verification failed",
	)
)

// NewWebhookInvalidPayloadError creates an error for invalid webhook payload
func NewWebhookInvalidPayloadError(message string) exception.DomainError {
	return exception.NewDomainError(
		WebhookErrInvalidPayload,
		message,
	)
}

// NewWebhookInvalidEventError creates an error for invalid webhook event type
func NewWebhookInvalidEventError(eventType string) exception.DomainError {
	return exception.NewDomainError(
		WebhookErrInvalidEvent,
		"unsupported webhook event type: "+eventType,
	)
}

// NewWebhookProjectNotFoundError creates an error when project is not found
func NewWebhookProjectNotFoundError(projectID string) exception.DomainError {
	return exception.NewDomainError(
		WebhookErrProjectNotFound,
		"project not found: "+projectID,
	)
}

// NewWebhookProcessingFailedError creates an error for webhook processing failures
func NewWebhookProcessingFailedError(message string) exception.DomainError {
	return exception.NewDomainError(
		WebhookErrProcessingFailed,
		"webhook processing failed: "+message,
	)
}
