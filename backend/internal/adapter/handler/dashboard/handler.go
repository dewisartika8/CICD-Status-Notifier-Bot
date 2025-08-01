package dashboard

import (
	"net/http"

	"github.com/gofiber/fiber/v2"

	"github.com/dewisartika8/cicd-status-notifier-bot/internal/core/dashboard/port"
	"github.com/dewisartika8/cicd-status-notifier-bot/internal/core/shared/domain/value_objects"
)

// Handler handles dashboard HTTP requests
type Handler struct {
	dashboardService port.DashboardService
}

// NewHandler creates a new dashboard handler
func NewHandler(dashboardService port.DashboardService) *Handler {
	return &Handler{
		dashboardService: dashboardService,
	}
}

// RegisterRoutes registers dashboard routes
func (h *Handler) RegisterRoutes(router fiber.Router) {
	dashboard := router.Group("/dashboard")
	dashboard.Get("/overview", h.GetOverview)

	// Project-specific routes
	projects := router.Group("/projects")
	projects.Get("/:id/statistics", h.GetProjectStatistics)

	// Build analytics routes
	builds := router.Group("/builds")
	builds.Get("/analytics", h.GetBuildAnalytics)
}

// GetOverview handles GET /api/v1/dashboard/overview
func (h *Handler) GetOverview(c *fiber.Ctx) error {
	overview, err := h.dashboardService.GetOverview()
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error":   "Failed to retrieve dashboard overview",
			"details": err.Error(),
		})
	}

	return c.Status(http.StatusOK).JSON(fiber.Map{
		"message": "Dashboard overview retrieved successfully",
		"data":    overview,
	})
}

// GetProjectStatistics handles GET /api/v1/projects/:id/statistics
func (h *Handler) GetProjectStatistics(c *fiber.Ctx) error {
	projectIDStr := c.Params("id")

	projectID, err := value_objects.NewIDFromString(projectIDStr)
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error":   "Invalid project ID format",
			"details": err.Error(),
		})
	}

	statistics, err := h.dashboardService.GetProjectStatistics(projectID)
	if err != nil {
		// Check if it's a not found error
		statusCode := http.StatusInternalServerError
		if err.Error() == "project not found" {
			statusCode = http.StatusNotFound
		}

		return c.Status(statusCode).JSON(fiber.Map{
			"error":   "Failed to retrieve project statistics",
			"details": err.Error(),
		})
	}

	return c.Status(http.StatusOK).JSON(fiber.Map{
		"message": "Project statistics retrieved successfully",
		"data":    statistics,
	})
}

// GetBuildAnalytics handles GET /api/v1/builds/analytics
func (h *Handler) GetBuildAnalytics(c *fiber.Ctx) error {
	timeRange := c.Query("range", "7d") // default to 7 days

	// Validate time range
	validRanges := map[string]bool{
		"7d":  true,
		"30d": true,
		"90d": true,
	}

	if !validRanges[timeRange] {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error":   "Invalid time range",
			"details": "Supported ranges are: 7d, 30d, 90d",
		})
	}

	analytics, err := h.dashboardService.GetBuildAnalytics(timeRange)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error":   "Failed to retrieve build analytics",
			"details": err.Error(),
		})
	}

	return c.Status(http.StatusOK).JSON(fiber.Map{
		"message": "Build analytics retrieved successfully",
		"data":    analytics,
	})
}
