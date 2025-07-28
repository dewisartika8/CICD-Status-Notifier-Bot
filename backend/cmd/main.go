package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/dewisartika8/cicd-status-notifier-bot/internal/config"
	"github.com/dewisartika8/cicd-status-notifier-bot/internal/repositories"
	"github.com/dewisartika8/cicd-status-notifier-bot/internal/services"
	"github.com/dewisartika8/cicd-status-notifier-bot/pkg/database"
	"github.com/dewisartika8/cicd-status-notifier-bot/pkg/logger"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/recover"
)

func main() {
	// Load configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Config error: %v", err)
	}

	// Initialize logger
	logger := logger.NewLogger()
	logger.Info("Starting CI/CD Status Notifier Bot...")

	// Connect to database
	dbConfig := database.Config{
		Host:         cfg.Database.Host,
		Port:         cfg.Database.Port,
		User:         cfg.Database.User,
		Password:     cfg.Database.Password,
		DBName:       cfg.Database.DBName,
		SSLMode:      cfg.Database.SSLMode,
		MaxOpenConns: cfg.Database.MaxOpenConns,
		MaxIdleConns: cfg.Database.MaxIdleConns,
		MaxLifetime:  time.Duration(cfg.Database.MaxLifetime) * time.Second,
	}

	db, err := database.Connect(dbConfig)
	if err != nil {
		logger.Fatalf("Database connection error: %v", err)
	}

	logger.Info("Database connected successfully")

	// Initialize repositories
	projectRepo := repositories.NewProjectRepository(db)
	buildEventRepo := repositories.NewBuildEventRepository(db)

	// Initialize services
	_ = services.NewProjectService(projectRepo)
	_ = services.NewBuildEventService(buildEventRepo)

	// Create Fiber app
	app := fiber.New(fiber.Config{
		AppName: "CI/CD Status Notifier Bot",
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			logger.Errorf("Request error: %v", err)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Internal Server Error",
			})
		},
	})

	// Middleware
	app.Use(recover.New())
	app.Use(cors.New())

	// Routes
	app.Get("/", func(c *fiber.Ctx) error {
		logger.Info("Health check request received")
		return c.JSON(fiber.Map{
			"message": "CI/CD Status Notifier Bot is running 🚀",
			"status":  "healthy",
			"version": "1.0.0",
		})
	})

	app.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"status":    "healthy",
			"database":  "connected",
			"timestamp": time.Now().UTC(),
		})
	})

	// Graceful shutdown
	go func() {
		if err := app.Listen(":" + cfg.ServerPort); err != nil {
			logger.Fatalf("Server failed to start: %v", err)
		}
	}()

	logger.Infof("Server started on port %s", cfg.ServerPort)

	// Wait for interrupt signal to gracefully shutdown the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Info("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := app.ShutdownWithContext(ctx); err != nil {
		logger.Fatalf("Server forced to shutdown: %v", err)
	}

	logger.Info("Server exited")
}
