package webhook

import (
	"encoding/json"
	"strconv"

	"github.com/dewisartika8/cicd-status-notifier-bot/internal/core/shared/domain/value_objects"
	"github.com/dewisartika8/cicd-status-notifier-bot/internal/core/webhook/domain"
	"github.com/dewisartika8/cicd-status-notifier-bot/internal/core/webhook/dto"
	"github.com/gofiber/fiber/v2"
)

// Error messages
const (
	ErrorProjectIDRequired        = "project_id is required"
	ErrorInvalidProjectIDFormat   = "invalid project_id format"
	ErrorMissingSignatureHeader   = "missing X-Hub-Signature-256 header"
	ErrorMissingEventTypeHeader   = "missing X-GitHub-Event header"
	ErrorUnsupportedEventType     = "unsupported event type: "
	ErrorEmptyRequestBody         = "empty request body"
	ErrorInvalidJSONPayload       = "invalid JSON payload"
	ErrorInvalidWebhookSignature  = "invalid webhook signature"
	ErrorProjectNotFound          = "project not found"
	ErrorInvalidWebhookPayload    = "invalid webhook payload"
	ErrorInternalServerError      = "internal server error"
	ErrorEventIDRequired          = "event_id is required"
	ErrorInvalidEventIDFormat     = "invalid event_id format"
	ErrorWebhookEventNotFound     = "webhook event not found"
	ErrorFailedToGetWebhookEvents = "failed to get webhook events"
)

// Success messages
const (
	MessageWebhookProcessedSuccessfully = "webhook processed successfully"
)

// Log messages
const (
	LogFailedToProcessWebhook       = "Failed to process webhook"
	LogWebhookProcessedSuccessfully = "Webhook processed successfully"
	LogFailedToGetWebhookEvents     = "Failed to get webhook events"
	LogFailedToGetWebhookEvent      = "Failed to get webhook event"
)

// RegisterRoutes registers webhook routes following the health handler pattern
func (h *WebhookHandler) RegisterRoutes(r fiber.Router) {
	// GitHub webhook endpoint
	r.Post("/github/:projectId", h.ProcessGitHubWebhook)

	// Webhook events endpoints
	r.Get("/events/:projectId", h.GetWebhookEvents)
	r.Get("/events/:projectId/:eventId", h.GetWebhookEvent)
}

// ProcessGitHubWebhook handles incoming GitHub webhook requests
func (h *WebhookHandler) ProcessGitHubWebhook(c *fiber.Ctx) error {
	// Extract project ID from URL parameters
	projectIDStr := c.Params("projectId")
	if projectIDStr == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": ErrorProjectIDRequired,
		})
	}

	projectID, err := value_objects.NewIDFromString(projectIDStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": ErrorInvalidProjectIDFormat,
		})
	}

	// Extract signature from headers
	signature := c.Get("X-Hub-Signature-256")
	if signature == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": ErrorMissingSignatureHeader,
		})
	}

	// Extract event type from headers
	eventTypeStr := c.Get("X-GitHub-Event")
	if eventTypeStr == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": ErrorMissingEventTypeHeader,
		})
	}

	// Validate and convert event type
	eventType := domain.WebhookEventType(eventTypeStr)
	if !h.isValidEventType(eventType) {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": ErrorUnsupportedEventType + eventTypeStr,
		})
	}

	// Extract delivery ID (optional but recommended for idempotency)
	deliveryID := c.Get("X-GitHub-Delivery")

	// Get request body
	body := c.Body()
	if len(body) == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": ErrorEmptyRequestBody,
		})
	}

	// Parse payload
	var payload dto.GitHubActionsPayload
	if err := json.Unmarshal(body, &payload); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": ErrorInvalidJSONPayload,
		})
	}

	// Create process webhook request
	processReq := dto.ProcessWebhookRequest{
		ProjectID:  projectID,
		EventType:  eventType,
		Signature:  signature,
		DeliveryID: deliveryID,
		Body:       body,
		Payload:    payload,
	}

	// Process webhook
	webhookEvent, err := h.webhookService.ProcessWebhook(c.Context(), processReq)
	if err != nil {
		h.logger.Error(LogFailedToProcessWebhook, map[string]interface{}{
			"project_id":  projectIDStr,
			"event_type":  eventTypeStr,
			"delivery_id": deliveryID,
			"error":       err.Error(),
		})

		// Handle different error types
		switch err.(type) {
		case interface{ Code() string }:
			domainErr := err.(interface{ Code() string })
			switch domainErr.Code() {
			case domain.WebhookErrInvalidSignature:
				return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
					"error": ErrorInvalidWebhookSignature,
				})
			case domain.WebhookErrProjectNotFound:
				return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
					"error": ErrorProjectNotFound,
				})
			case domain.WebhookErrInvalidPayload:
				return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
					"error": ErrorInvalidWebhookPayload,
				})
			default:
				return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
					"error": ErrorInternalServerError,
				})
			}
		default:
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": ErrorInternalServerError,
			})
		}
	}

	// Log successful processing
	h.logger.Info(LogWebhookProcessedSuccessfully, map[string]interface{}{
		"project_id":       projectIDStr,
		"event_type":       eventTypeStr,
		"delivery_id":      deliveryID,
		"webhook_event_id": webhookEvent.ID().String(),
	})

	// Return response
	response := dto.ToWebhookEventResponse(webhookEvent)
	return c.Status(fiber.StatusAccepted).JSON(fiber.Map{
		"message": MessageWebhookProcessedSuccessfully,
		"data":    response,
	})
}

// GetWebhookEvents retrieves webhook events for a project
func (h *WebhookHandler) GetWebhookEvents(c *fiber.Ctx) error {
	// Extract project ID from URL parameters
	projectIDStr := c.Params("projectId")
	if projectIDStr == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": ErrorProjectIDRequired,
		})
	}

	projectID, err := value_objects.NewIDFromString(projectIDStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": ErrorInvalidProjectIDFormat,
		})
	}

	// Extract pagination parameters
	limitStr := c.Query("limit", "20")
	offsetStr := c.Query("offset", "0")

	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit <= 0 || limit > 100 {
		limit = 20
	}

	offset, err := strconv.Atoi(offsetStr)
	if err != nil || offset < 0 {
		offset = 0
	}

	// Get webhook events
	webhookEvents, err := h.webhookService.GetWebhookEventsByProject(c.Context(), projectID, limit, offset)
	if err != nil {
		h.logger.Error(LogFailedToGetWebhookEvents, map[string]interface{}{
			"project_id": projectIDStr,
			"error":      err.Error(),
		})
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": ErrorFailedToGetWebhookEvents,
		})
	}

	// Convert to response DTOs
	responses := make([]*dto.WebhookEventResponse, len(webhookEvents))
	for i, event := range webhookEvents {
		responses[i] = dto.ToWebhookEventResponse(event)
	}

	return c.JSON(fiber.Map{
		"data": responses,
		"pagination": fiber.Map{
			"limit":  limit,
			"offset": offset,
			"count":  len(responses),
		},
	})
}

// GetWebhookEvent retrieves a specific webhook event
func (h *WebhookHandler) GetWebhookEvent(c *fiber.Ctx) error {
	// Extract webhook event ID from URL parameters
	eventIDStr := c.Params("eventId")
	if eventIDStr == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": ErrorEventIDRequired,
		})
	}

	eventID, err := value_objects.NewIDFromString(eventIDStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": ErrorInvalidEventIDFormat,
		})
	}

	// Get webhook event
	webhookEvent, err := h.webhookService.GetWebhookEvent(c.Context(), eventID)
	if err != nil {
		h.logger.Error(LogFailedToGetWebhookEvent, map[string]interface{}{
			"event_id": eventIDStr,
			"error":    err.Error(),
		})
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": ErrorWebhookEventNotFound,
		})
	}

	// Convert to response DTO
	response := dto.ToWebhookEventResponse(webhookEvent)
	return c.JSON(fiber.Map{
		"data": response,
	})
}

// isValidEventType checks if the event type is supported
func (h *WebhookHandler) isValidEventType(eventType domain.WebhookEventType) bool {
	switch eventType {
	case domain.WorkflowRunEvent, domain.PushEvent, domain.PullRequestEvent:
		return true
	default:
		return false
	}
}
