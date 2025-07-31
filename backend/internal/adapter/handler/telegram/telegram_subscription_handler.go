package telegram

import (
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"

	"github.com/dewisartika8/cicd-status-notifier-bot/internal/core/bot/dto"
	"github.com/dewisartika8/cicd-status-notifier-bot/internal/core/notification/port"
	"github.com/dewisartika8/cicd-status-notifier-bot/internal/core/shared/domain/value_objects"
)

// Error message constants
const (
	ErrInvalidRequest           = "Invalid request"
	ErrValidationFailed         = "Validation failed"
	ErrInvalidProjectID         = "Invalid project ID"
	ErrInvalidSubscriptionID    = "Invalid subscription ID"
	ErrProjectIDRequired        = "project_id is required"
	ErrChatIDCannotBeZero       = "chat_id cannot be zero"
	ErrProjectIDMustBeUUID      = "Project ID must be a valid UUID"
	ErrSubscriptionIDRequired   = "Subscription ID is required"
	ErrSubscriptionIDMustBeUUID = "Subscription ID must be a valid UUID"
	ErrInvalidRequestBody       = "Invalid request body"
	ErrSubscriptionNotFound     = "Subscription not found"
	ErrFailedToCreateSub        = "Failed to create subscription"
	ErrFailedToGetSub           = "Failed to get telegram subscription"
	ErrFailedToGetSubsByProject = "Failed to get subscriptions by project"
	ErrFailedToGetSubs          = "Failed to get subscriptions"
	ErrFailedToUpdateSub        = "Failed to update subscription"
	ErrFailedToDeleteSub        = "Failed to delete subscription"
	ErrFailedToGetActiveSubs    = "Failed to get active subscriptions"
	ErrFailedToGetSubCount      = "Failed to get subscription count"
	ErrFailedToGetSubStats      = "Failed to get subscription stats"
	ErrInvalidIsActiveParam     = "Invalid is_active parameter"
	ErrIsActiveMustBeBoolean    = "is_active must be a boolean (true/false)"

	// Log message constants
	LogInvalidProjectIDFormat      = "Invalid project ID format"
	LogInvalidSubscriptionIDFormat = "Invalid subscription ID format"
	LogInvalidIsActiveFormat       = "Invalid is_active format"
	LogFailedToParseCreateReq      = "Failed to parse create subscription request"
	LogFailedToParseUpdateReq      = "Failed to parse update subscription request"
	LogFailedToCreateSub           = "Failed to create telegram subscription"
	LogFailedToGetSub              = "Failed to get telegram subscription"
	LogFailedToGetSubsByProject    = "Failed to get subscriptions by project"
	LogFailedToUpdateSub           = "Failed to update telegram subscription"
	LogFailedToDeleteSub           = "Failed to delete telegram subscription"
	LogFailedToGetActiveSubs       = "Failed to get active subscriptions"
	LogFailedToGetSubCount         = "Failed to get subscription count"

	// Success message constants
	MsgSubCreatedSuccessfully          = "Subscription created successfully"
	MsgSubRetrievedSuccessfully        = "Subscription retrieved successfully"
	MsgSubsRetrievedSuccessfully       = "Subscriptions retrieved successfully"
	MsgSubUpdatedSuccessfully          = "Subscription updated successfully"
	MsgSubDeletedSuccessfully          = "Subscription deleted successfully"
	MsgActiveSubsRetrievedSuccessfully = "Active subscriptions retrieved successfully"
	MsgSubStatsRetrievedSuccessfully   = "Subscription stats retrieved successfully"
)

type TelegramSubscriptionHandler struct {
	subscriptionService port.TelegramSubscriptionService
	logger              *logrus.Logger
}

func NewTelegramSubscriptionHandler(subscriptionService port.TelegramSubscriptionService, logger *logrus.Logger) *TelegramSubscriptionHandler {
	return &TelegramSubscriptionHandler{
		subscriptionService: subscriptionService,
		logger:              logger,
	}
}

// CreateSubscription creates a new telegram subscription
func (h *TelegramSubscriptionHandler) CreateSubscription(c *fiber.Ctx) error {
	var req dto.CreateTelegramSubscriptionRequest
	if err := c.BodyParser(&req); err != nil {
		h.logger.WithError(err).Error(LogFailedToParseCreateReq)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   ErrInvalidRequestBody,
			"message": err.Error(),
		})
	}

	// Validate request
	if req.ProjectID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   ErrValidationFailed,
			"message": ErrProjectIDRequired,
		})
	}

	// ChatID of 0 is invalid, but negative values are valid for groups
	if req.ChatID == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   ErrValidationFailed,
			"message": ErrChatIDCannotBeZero,
		})
	}

	// Parse project ID
	projectID, err := value_objects.NewIDFromString(req.ProjectID)
	if err != nil {
		h.logger.WithError(err).Error(LogInvalidProjectIDFormat)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   ErrInvalidProjectID,
			"message": ErrProjectIDMustBeUUID,
		})
	}

	// Create subscription
	subscription, err := h.subscriptionService.CreateTelegramSubscription(c.Context(), projectID, req.ChatID)
	if err != nil {
		h.logger.WithError(err).Error(LogFailedToCreateSub)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   ErrFailedToCreateSub,
			"message": err.Error(),
		})
	}

	response := dto.TelegramSubscriptionResponse{
		ID:        subscription.ID().String(),
		ProjectID: subscription.ProjectID().String(),
		ChatID:    subscription.ChatID(),
		IsActive:  subscription.IsActive(),
		CreatedAt: subscription.CreatedAt().Unix(),
		UpdatedAt: subscription.UpdatedAt().Unix(),
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": MsgSubCreatedSuccessfully,
		"data":    response,
	})
}

// GetSubscriptionByID gets a subscription by ID
func (h *TelegramSubscriptionHandler) GetSubscriptionByID(c *fiber.Ctx) error {
	subscriptionIDStr := c.Params("id")
	if subscriptionIDStr == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   ErrInvalidRequest,
			"message": ErrSubscriptionIDRequired,
		})
	}

	subscriptionID, err := value_objects.NewIDFromString(subscriptionIDStr)
	if err != nil {
		h.logger.WithError(err).Error(LogInvalidSubscriptionIDFormat)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   ErrInvalidSubscriptionID,
			"message": ErrSubscriptionIDMustBeUUID,
		})
	}

	subscription, err := h.subscriptionService.GetTelegramSubscription(c.Context(), subscriptionID)
	if err != nil {
		h.logger.WithError(err).Error(LogFailedToGetSub)
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error":   ErrSubscriptionNotFound,
			"message": err.Error(),
		})
	}

	response := dto.TelegramSubscriptionResponse{
		ID:        subscription.ID().String(),
		ProjectID: subscription.ProjectID().String(),
		ChatID:    subscription.ChatID(),
		IsActive:  subscription.IsActive(),
		CreatedAt: subscription.CreatedAt().Unix(),
		UpdatedAt: subscription.UpdatedAt().Unix(),
	}

	return c.JSON(fiber.Map{
		"message": MsgSubRetrievedSuccessfully,
		"data":    response,
	})
}

// GetSubscriptionsByProject gets all subscriptions for a project
func (h *TelegramSubscriptionHandler) GetSubscriptionsByProject(c *fiber.Ctx) error {
	projectIDStr := c.Params("projectId")
	if projectIDStr == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   ErrInvalidRequest,
			"message": ErrProjectIDRequired,
		})
	}

	projectID, err := value_objects.NewIDFromString(projectIDStr)
	if err != nil {
		h.logger.WithError(err).Error(LogInvalidProjectIDFormat)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   ErrInvalidProjectID,
			"message": ErrProjectIDMustBeUUID,
		})
	}

	subscriptions, err := h.subscriptionService.GetTelegramSubscriptionsByProject(c.Context(), projectID)
	if err != nil {
		h.logger.WithError(err).Error(LogFailedToGetSubsByProject)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   ErrFailedToGetSubs,
			"message": err.Error(),
		})
	}

	responses := make([]dto.TelegramSubscriptionResponse, len(subscriptions))
	for i, subscription := range subscriptions {
		responses[i] = dto.TelegramSubscriptionResponse{
			ID:        subscription.ID().String(),
			ProjectID: subscription.ProjectID().String(),
			ChatID:    subscription.ChatID(),
			IsActive:  subscription.IsActive(),
			CreatedAt: subscription.CreatedAt().Unix(),
			UpdatedAt: subscription.UpdatedAt().Unix(),
		}
	}

	return c.JSON(fiber.Map{
		"message": MsgSubsRetrievedSuccessfully,
		"data":    responses,
	})
}

// UpdateSubscription updates an existing subscription
func (h *TelegramSubscriptionHandler) UpdateSubscription(c *fiber.Ctx) error {
	subscriptionIDStr := c.Params("id")
	if subscriptionIDStr == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   ErrInvalidRequest,
			"message": ErrSubscriptionIDRequired,
		})
	}

	subscriptionID, err := value_objects.NewIDFromString(subscriptionIDStr)
	if err != nil {
		h.logger.WithError(err).Error(LogInvalidSubscriptionIDFormat)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   ErrInvalidSubscriptionID,
			"message": ErrSubscriptionIDMustBeUUID,
		})
	}

	var req dto.UpdateTelegramSubscriptionRequest
	if err := c.BodyParser(&req); err != nil {
		h.logger.WithError(err).Error(LogFailedToParseUpdateReq)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   ErrInvalidRequestBody,
			"message": err.Error(),
		})
	}

	subscription, err := h.subscriptionService.UpdateTelegramSubscription(c.Context(), subscriptionID, nil, &req.IsActive)
	if err != nil {
		h.logger.WithError(err).Error(LogFailedToUpdateSub)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   ErrFailedToUpdateSub,
			"message": err.Error(),
		})
	}

	response := dto.TelegramSubscriptionResponse{
		ID:        subscription.ID().String(),
		ProjectID: subscription.ProjectID().String(),
		ChatID:    subscription.ChatID(),
		IsActive:  subscription.IsActive(),
		CreatedAt: subscription.CreatedAt().Unix(),
		UpdatedAt: subscription.UpdatedAt().Unix(),
	}

	return c.JSON(fiber.Map{
		"message": MsgSubUpdatedSuccessfully,
		"data":    response,
	})
}

// DeleteSubscription deletes a subscription
func (h *TelegramSubscriptionHandler) DeleteSubscription(c *fiber.Ctx) error {
	subscriptionIDStr := c.Params("id")
	if subscriptionIDStr == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   ErrInvalidRequest,
			"message": ErrSubscriptionIDRequired,
		})
	}

	subscriptionID, err := value_objects.NewIDFromString(subscriptionIDStr)
	if err != nil {
		h.logger.WithError(err).Error(LogInvalidSubscriptionIDFormat)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   ErrInvalidSubscriptionID,
			"message": ErrSubscriptionIDMustBeUUID,
		})
	}

	err = h.subscriptionService.DeleteTelegramSubscription(c.Context(), subscriptionID)
	if err != nil {
		h.logger.WithError(err).Error(LogFailedToDeleteSub)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   ErrFailedToDeleteSub,
			"message": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"message": MsgSubDeletedSuccessfully,
	})
}

// GetActiveSubscriptions gets all active subscriptions
func (h *TelegramSubscriptionHandler) GetActiveSubscriptions(c *fiber.Ctx) error {
	subscriptions, err := h.subscriptionService.GetAllActiveSubscriptions(c.Context())
	if err != nil {
		h.logger.WithError(err).Error(LogFailedToGetActiveSubs)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   ErrFailedToGetActiveSubs,
			"message": err.Error(),
		})
	}

	responses := make([]dto.TelegramSubscriptionResponse, len(subscriptions))
	for i, subscription := range subscriptions {
		responses[i] = dto.TelegramSubscriptionResponse{
			ID:        subscription.ID().String(),
			ProjectID: subscription.ProjectID().String(),
			ChatID:    subscription.ChatID(),
			IsActive:  subscription.IsActive(),
			CreatedAt: subscription.CreatedAt().Unix(),
			UpdatedAt: subscription.UpdatedAt().Unix(),
		}
	}

	return c.JSON(fiber.Map{
		"message": MsgActiveSubsRetrievedSuccessfully,
		"data":    responses,
	})
}

// GetSubscriptionStats gets subscription statistics
func (h *TelegramSubscriptionHandler) GetSubscriptionStats(c *fiber.Ctx) error {
	projectIDStr := c.Query("project_id")
	isActiveStr := c.Query("is_active")

	var projectID *value_objects.ID
	var isActive *bool

	if projectIDStr != "" {
		pid, err := value_objects.NewIDFromString(projectIDStr)
		if err != nil {
			h.logger.WithError(err).Error(LogInvalidProjectIDFormat)
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error":   ErrInvalidProjectID,
				"message": ErrProjectIDMustBeUUID,
			})
		}
		projectID = &pid
	}

	if isActiveStr != "" {
		active, err := strconv.ParseBool(isActiveStr)
		if err != nil {
			h.logger.WithError(err).Error(LogInvalidIsActiveFormat)
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error":   ErrInvalidIsActiveParam,
				"message": ErrIsActiveMustBeBoolean,
			})
		}
		isActive = &active
	}

	count, err := h.subscriptionService.GetSubscriptionCount(c.Context(), projectID, isActive)
	if err != nil {
		h.logger.WithError(err).Error(LogFailedToGetSubCount)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   ErrFailedToGetSubStats,
			"message": err.Error(),
		})
	}

	stats := dto.TelegramSubscriptionStatsResponse{
		Count:     count,
		ProjectID: projectIDStr,
		IsActive:  isActive,
	}

	return c.JSON(fiber.Map{
		"message": MsgSubStatsRetrievedSuccessfully,
		"data":    stats,
	})
}
