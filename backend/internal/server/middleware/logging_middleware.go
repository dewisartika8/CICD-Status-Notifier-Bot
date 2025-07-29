package middleware

import (
	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
)

func trafficLogMiddleware(logger *logrus.Logger) fiber.Handler {
	return func(c *fiber.Ctx) error {
		logger.Info("Incoming request", map[string]interface{}{
			"method": c.Method(),
			"url":    c.OriginalURL(),
		})
		return c.Next()
	}
}
