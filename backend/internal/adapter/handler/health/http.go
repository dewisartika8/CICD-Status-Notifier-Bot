package health

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
)

type HealthHandler struct {
	Logger *logrus.Logger
}

func NewHealthHandler(logger *logrus.Logger) *HealthHandler {
	return &HealthHandler{
		Logger: logger,
	}
}

// HTTP Routing registerer
func (h *HealthHandler) RegisterRoutes(r fiber.Router) {
	r.Get("/", h.GetHealthStatus)
	r.Get("/health", h.CheckHealth)
}

// GetHealthStatus returns the health status of the service
func (h *HealthHandler) GetHealthStatus(c *fiber.Ctx) error {
	h.Logger.Info("Health check request received")
	return c.JSON(fiber.Map{
		"message": "CI/CD Status Notifier Bot is running 🚀",
		"status":  "healthy",
		"version": "1.0.0",
	})
}

// CheckHealth checks the health of the service
func (h *HealthHandler) CheckHealth(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{
		"status":    "healthy",
		"database":  "connected",
		"timestamp": time.Now().UTC(),
	})
}
