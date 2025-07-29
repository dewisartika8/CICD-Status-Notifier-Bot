package main

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"log"
	"time"

	"github.com/dewisartika8/cicd-status-notifier-bot/internal/adapter/handler/health"
	"github.com/dewisartika8/cicd-status-notifier-bot/internal/adapter/repository"
	"github.com/dewisartika8/cicd-status-notifier-bot/internal/config"
	bs "github.com/dewisartika8/cicd-status-notifier-bot/internal/core/build/service"
	ps "github.com/dewisartika8/cicd-status-notifier-bot/internal/core/project/service"
	"github.com/dewisartika8/cicd-status-notifier-bot/internal/server/app"
	"github.com/dewisartika8/cicd-status-notifier-bot/pkg/database"
	"github.com/dewisartika8/cicd-status-notifier-bot/pkg/logger"
)

func verifyGitHubSignature(secret, signature string, body []byte) bool {
	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write(body)
	expected := "sha256=" + hex.EncodeToString(mac.Sum(nil))
	return hmac.Equal([]byte(expected), []byte(signature))
}

// Struct untuk payload GitHub Actions (contoh sederhana)
type GitHubActionsPayload struct {
	Action     string `json:"action"`
	Workflow   string `json:"workflow"`
	Repository struct {
		Name string `json:"name"`
	} `json:"repository"`
	Sender struct {
		Login string `json:"login"`
	} `json:"sender"`
	// Tambahkan field lain sesuai kebutuhan
}

// Handler untuk event processing
func processGitHubEvent(eventType string, payload GitHubActionsPayload) error {
	// Event validation & filtering
	if payload.Repository.Name == "" || payload.Workflow == "" {
		return fmt.Errorf("invalid payload: missing repository or workflow")
	}

	switch eventType {
	case "workflow_run":
		// Handler untuk workflow_run
		fmt.Printf("Processing workflow_run for %s/%s\n", payload.Repository.Name, payload.Workflow)
		// Tambahkan logika sesuai kebutuhan
	case "push":
		// Handler untuk push event
		fmt.Printf("Processing push for %s\n", payload.Repository.Name)
		// Tambahkan logika sesuai kebutuhan
	default:
		fmt.Printf("Unhandled event type: %s\n", eventType)
	}
	return nil
}

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
	projectRepo := repository.NewProjectRepository(db)
	buildEventRepo := repository.NewBuildEventRepository(db)

	// Initialize services
	_ = ps.NewProjectService(ps.Dep{
		ProjectRepo: projectRepo,
	})
	_ = bs.NewBuildEventService(bs.Dep{
		BuildEventRepo: buildEventRepo,
	})

	// Initialize handlers
	healthHandler := health.NewHealthHandler(logger)

	// run APP in http server
	// inject all usecases here
	appService := app.Init(app.Dep{
		AppConfig:     cfg,
		HealthHandler: healthHandler,
		Logger:        logger,
	})
	appService.Run() // start http server
}
