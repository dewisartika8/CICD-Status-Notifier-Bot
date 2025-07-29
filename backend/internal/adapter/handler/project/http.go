package project

import (
	"context"
	"strconv"

	"github.com/dewisartika8/cicd-status-notifier-bot/internal/core/project/domain"
	"github.com/dewisartika8/cicd-status-notifier-bot/internal/core/project/dto"
	"github.com/dewisartika8/cicd-status-notifier-bot/internal/core/shared/domain/value_objects"
	"github.com/dewisartika8/cicd-status-notifier-bot/pkg/exception"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

// Constants for error messages and responses
const (
	// Error messages
	ErrorFailedToParseRequestBody = "Failed to parse request body"
	ErrorInvalidRequestBody       = "Invalid request body"
	ErrorRequestValidationFailed  = "Request validation failed"
	ErrorValidationFailed         = "Validation failed"
	ErrorInvalidProjectID         = "Invalid project ID"
	ErrorInternalServer           = "Internal server error"

	// Success messages
	MessageProjectCreatedSuccessfully       = "Project created successfully"
	MessageProjectsRetrievedSuccessfully    = "Projects retrieved successfully"
	MessageProjectRetrievedSuccessfully     = "Project retrieved successfully"
	MessageProjectUpdatedSuccessfully       = "Project updated successfully"
	MessageProjectDeletedSuccessfully       = "Project deleted successfully"
	MessageProjectStatusUpdatedSuccessfully = "Project status updated successfully"

	// Log messages
	LogCreatingNewProject          = "Creating new project"
	LogListingProjects             = "Listing projects"
	LogGettingProject              = "Getting project"
	LogUpdatingProject             = "Updating project"
	LogDeletingProject             = "Deleting project"
	LogUpdatingProjectStatus       = "Updating project status"
	LogFilterValidationFailed      = "Filter validation failed"
	LogFailedToListProjects        = "Failed to list projects"
	LogFailedToCountProjects       = "Failed to count projects"
	LogProjectsListedSuccessfully  = "Projects listed successfully"
	LogFailedToCreateProject       = "Failed to create project"
	LogFailedToGetProject          = "Failed to get project"
	LogFailedToUpdateProject       = "Failed to update project"
	LogFailedToDeleteProject       = "Failed to delete project"
	LogFailedToUpdateProjectStatus = "Failed to update project status"
)

// HTTP Routing registerer
func (h *Handler) RegisterRoutes(r fiber.Router) {
	projects := r.Group("/projects")

	projects.Post("/", h.CreateProject)
	projects.Get("/", h.ListProjects)
	projects.Get("/:id", h.GetProject)
	projects.Put("/:id", h.UpdateProject)
	projects.Delete("/:id", h.DeleteProject)
	projects.Patch("/:id/status", h.UpdateProjectStatus)
}

// CreateProject creates a new project
func (h *Handler) CreateProject(c *fiber.Ctx) error {
	ctx := context.Background()
	h.Logger.Info(LogCreatingNewProject)

	var req dto.CreateProjectRequest
	if err := c.BodyParser(&req); err != nil {
		h.Logger.WithError(err).Error(ErrorFailedToParseRequestBody)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": ErrorInvalidRequestBody,
		})
	}

	validator := validator.New()
	if err := validator.Struct(&req); err != nil {
		h.Logger.WithError(err).Error(ErrorRequestValidationFailed)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   ErrorValidationFailed,
			"details": err.Error(),
		})
	}

	project, err := h.ProjectService.CreateProject(ctx, req)
	if err != nil {
		h.Logger.WithError(err).Error(LogFailedToCreateProject)
		return h.handleError(c, err)
	}

	response := dto.ToProjectResponse(project)
	h.Logger.WithField("project_id", project.ID().String()).Info(MessageProjectCreatedSuccessfully)

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": MessageProjectCreatedSuccessfully,
		"data":    response,
	})
}

// ListProjects retrieves a list of projects with filtering and pagination
func (h *Handler) ListProjects(c *fiber.Ctx) error {
	ctx := context.Background()
	h.Logger.Info(LogListingProjects)

	filters := h.parseListFilters(c)

	validator := validator.New()
	if err := validator.Struct(&filters); err != nil {
		h.Logger.WithError(err).Error(LogFilterValidationFailed)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   "Invalid filters",
			"details": err.Error(),
		})
	}

	projects, err := h.ProjectService.ListProjects(ctx, filters)
	if err != nil {
		h.Logger.WithError(err).Error(LogFailedToListProjects)
		return h.handleError(c, err)
	}

	total, err := h.ProjectService.CountProjects(ctx, filters)
	if err != nil {
		h.Logger.WithError(err).Error(LogFailedToCountProjects)
		return h.handleError(c, err)
	}

	response := dto.ListProjectResponse{
		Projects: dto.ToProjectResponseList(projects),
		Total:    total,
		Limit:    getIntValueOrDefault(filters.Limit, 10),
		Offset:   getIntValueOrDefault(filters.Offset, 0),
	}

	h.Logger.WithField("count", len(projects)).Info(LogProjectsListedSuccessfully)

	return c.JSON(fiber.Map{
		"message": MessageProjectsRetrievedSuccessfully,
		"data":    response,
	})
}

// GetProject retrieves a specific project by ID
func (h *Handler) GetProject(c *fiber.Ctx) error {
	ctx := context.Background()

	projectID, err := h.parseProjectID(c)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": ErrorInvalidProjectID,
		})
	}

	h.Logger.WithField("project_id", projectID.String()).Info(LogGettingProject)

	project, err := h.ProjectService.GetProject(ctx, projectID)
	if err != nil {
		h.Logger.WithError(err).WithField("project_id", projectID.String()).Error(LogFailedToGetProject)
		return h.handleError(c, err)
	}

	response := dto.ToProjectResponse(project)
	h.Logger.WithField("project_id", projectID.String()).Info(MessageProjectRetrievedSuccessfully)

	return c.JSON(fiber.Map{
		"message": MessageProjectRetrievedSuccessfully,
		"data":    response,
	})
}

// UpdateProject updates an existing project
func (h *Handler) UpdateProject(c *fiber.Ctx) error {
	ctx := context.Background()

	projectID, err := h.parseProjectID(c)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": ErrorInvalidProjectID,
		})
	}

	h.Logger.WithField("project_id", projectID.String()).Info(LogUpdatingProject)

	var req dto.UpdateProjectRequest
	if err := c.BodyParser(&req); err != nil {
		h.Logger.WithError(err).Error(ErrorFailedToParseRequestBody)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": ErrorInvalidRequestBody,
		})
	}

	validator := validator.New()
	if err := validator.Struct(&req); err != nil {
		h.Logger.WithError(err).Error(ErrorRequestValidationFailed)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   ErrorValidationFailed,
			"details": err.Error(),
		})
	}

	project, err := h.ProjectService.UpdateProject(ctx, projectID, req)
	if err != nil {
		h.Logger.WithError(err).WithField("project_id", projectID.String()).Error(LogFailedToUpdateProject)
		return h.handleError(c, err)
	}

	response := dto.ToProjectResponse(project)
	h.Logger.WithField("project_id", projectID.String()).Info(MessageProjectUpdatedSuccessfully)

	return c.JSON(fiber.Map{
		"message": MessageProjectUpdatedSuccessfully,
		"data":    response,
	})
}

// DeleteProject deletes a project
func (h *Handler) DeleteProject(c *fiber.Ctx) error {
	ctx := context.Background()

	projectID, err := h.parseProjectID(c)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": ErrorInvalidProjectID,
		})
	}

	h.Logger.WithField("project_id", projectID.String()).Info(LogDeletingProject)

	err = h.ProjectService.DeleteProject(ctx, projectID)
	if err != nil {
		h.Logger.WithError(err).WithField("project_id", projectID.String()).Error(LogFailedToDeleteProject)
		return h.handleError(c, err)
	}

	h.Logger.WithField("project_id", projectID.String()).Info(MessageProjectDeletedSuccessfully)

	return c.JSON(fiber.Map{
		"message": MessageProjectDeletedSuccessfully,
	})
}

// UpdateProjectStatus updates the project status
func (h *Handler) UpdateProjectStatus(c *fiber.Ctx) error {
	ctx := context.Background()

	projectID, err := h.parseProjectID(c)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": ErrorInvalidProjectID,
		})
	}

	h.Logger.WithField("project_id", projectID.String()).Info(LogUpdatingProjectStatus)

	var req dto.ProjectStatusUpdateRequest
	if err := c.BodyParser(&req); err != nil {
		h.Logger.WithError(err).Error(ErrorFailedToParseRequestBody)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": ErrorInvalidRequestBody,
		})
	}

	validator := validator.New()
	if err := validator.Struct(&req); err != nil {
		h.Logger.WithError(err).Error(ErrorRequestValidationFailed)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   ErrorValidationFailed,
			"details": err.Error(),
		})
	}

	project, err := h.ProjectService.UpdateProjectStatus(ctx, projectID, req.Status)
	if err != nil {
		h.Logger.WithError(err).WithField("project_id", projectID.String()).Error(LogFailedToUpdateProjectStatus)
		return h.handleError(c, err)
	}

	response := dto.ToProjectResponse(project)
	h.Logger.WithField("project_id", projectID.String()).WithField("status", req.Status).Info(MessageProjectStatusUpdatedSuccessfully)

	return c.JSON(fiber.Map{
		"message": MessageProjectStatusUpdatedSuccessfully,
		"data":    response,
	})
}

// Helper methods

func (h *Handler) parseProjectID(c *fiber.Ctx) (value_objects.ID, error) {
	idStr := c.Params("id")
	if idStr == "" {
		return value_objects.ID{}, exception.NewDomainError("VALIDATION_ERROR", "project ID is required")
	}

	return value_objects.NewIDFromString(idStr)
}

func (h *Handler) parseListFilters(c *fiber.Ctx) dto.ListProjectFilters {
	filters := dto.ListProjectFilters{}

	h.parseStatusFilter(c, &filters)
	h.parseStringFilters(c, &filters)
	h.parseBooleanFilters(c, &filters)
	h.parsePaginationFilters(c, &filters)
	h.parseSortingFilters(c, &filters)

	return filters
}

func (h *Handler) parseStatusFilter(c *fiber.Ctx, filters *dto.ListProjectFilters) {
	if statusStr := c.Query("status"); statusStr != "" {
		if status := parseProjectStatus(statusStr); status != nil {
			filters.Status = status
		}
	}
}

func (h *Handler) parseStringFilters(c *fiber.Ctx, filters *dto.ListProjectFilters) {
	if name := c.Query("name"); name != "" {
		filters.Name = &name
	}

	if repoURL := c.Query("repository_url"); repoURL != "" {
		filters.RepositoryURL = &repoURL
	}
}

func (h *Handler) parseBooleanFilters(c *fiber.Ctx, filters *dto.ListProjectFilters) {
	if hasTelegramStr := c.Query("has_telegram_chat"); hasTelegramStr != "" {
		if hasTelegram, err := strconv.ParseBool(hasTelegramStr); err == nil {
			filters.HasTelegramChat = &hasTelegram
		}
	}
}

func (h *Handler) parsePaginationFilters(c *fiber.Ctx, filters *dto.ListProjectFilters) {
	if limitStr := c.Query("limit"); limitStr != "" {
		if limit, err := strconv.Atoi(limitStr); err == nil {
			filters.Limit = &limit
		}
	}

	if offsetStr := c.Query("offset"); offsetStr != "" {
		if offset, err := strconv.Atoi(offsetStr); err == nil {
			filters.Offset = &offset
		}
	}
}

func (h *Handler) parseSortingFilters(c *fiber.Ctx, filters *dto.ListProjectFilters) {
	if sortBy := c.Query("sort_by"); sortBy != "" {
		filters.SortBy = &sortBy
	}

	if sortOrder := c.Query("sort_order"); sortOrder != "" {
		filters.SortOrder = &sortOrder
	}
}

func (h *Handler) handleError(c *fiber.Ctx, err error) error {
	if domainErr, ok := err.(*exception.DomainError); ok {
		switch domainErr.Code {
		case "PROJECT_NOT_FOUND":
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": domainErr.Message,
			})
		case "PROJECT_ALREADY_EXISTS":
			return c.Status(fiber.StatusConflict).JSON(fiber.Map{
				"error": domainErr.Message,
			})
		case "VALIDATION_ERROR":
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": domainErr.Message,
			})
		default:
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": ErrorInternalServer,
			})
		}
	}

	return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
		"error": ErrorInternalServer,
	})
}

// Helper functions

func parseProjectStatus(statusStr string) *domain.ProjectStatus {
	switch statusStr {
	case "active":
		status := domain.ProjectStatusActive
		return &status
	case "inactive":
		status := domain.ProjectStatusInactive
		return &status
	case "archived":
		status := domain.ProjectStatusArchived
		return &status
	default:
		return nil
	}
}

func getIntValueOrDefault(value *int, defaultValue int) int {
	if value != nil {
		return *value
	}
	return defaultValue
}
