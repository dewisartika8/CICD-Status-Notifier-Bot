package main

import (
	"log"

	"cicd-status-notifier-bot/internal/config"
	"cicd-status-notifier-bot/internal/logger"

	"github.com/gofiber/fiber/v2"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Config error: %v", err)
	}

	log := logger.NewLogger()
	log.Info("Starting CI/CD Status Notifier Bot...")

	app := fiber.New()
	app.Get("/", func(c *fiber.Ctx) error {
		log.Info("Received request on /")
		return c.SendString("CI/CD Status Notifier Bot is running 🚀")
	})

	log.Fatal(app.Listen(":" + cfg.ServerPort))
}
