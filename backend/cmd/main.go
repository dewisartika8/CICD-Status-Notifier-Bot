package main

import (
	"log"
	"time"

	"github.com/dewisartika8/cicd-status-notifier-bot/internal/adapter/handler/health"
	"github.com/dewisartika8/cicd-status-notifier-bot/internal/adapter/handler/project"
	"github.com/dewisartika8/cicd-status-notifier-bot/internal/adapter/handler/webhook"
	"github.com/dewisartika8/cicd-status-notifier-bot/internal/adapter/repository"
	"github.com/dewisartika8/cicd-status-notifier-bot/internal/config"
	bs "github.com/dewisartika8/cicd-status-notifier-bot/internal/core/build/service"
	ps "github.com/dewisartika8/cicd-status-notifier-bot/internal/core/project/service"
	ws "github.com/dewisartika8/cicd-status-notifier-bot/internal/core/webhook/service"
	"github.com/dewisartika8/cicd-status-notifier-bot/internal/server/app"
	"github.com/dewisartika8/cicd-status-notifier-bot/pkg/crypto"
	"github.com/dewisartika8/cicd-status-notifier-bot/pkg/database"
	"github.com/dewisartika8/cicd-status-notifier-bot/pkg/logger"
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
	projectRepo := repository.NewProjectRepository(db)
	buildEventRepo := repository.NewBuildEventRepository(db)
	webhookEventRepo := repository.NewWebhookEventRepository(db)

	// Initialize services
	projectService := ps.NewProjectService(ps.Dep{
		ProjectRepo: projectRepo,
	})
	buildService := bs.NewBuildEventService(bs.Dep{
		BuildEventRepo: buildEventRepo,
	})

	// Initialize crypto components
	signatureVerifier := crypto.NewGitHubSignatureVerifier()

	// Initialize webhook service
	webhookService := ws.NewWebhookService(ws.Dep{
		WebhookEventRepo:  webhookEventRepo,
		ProjectService:    projectService,
		BuildService:      buildService,
		SignatureVerifier: signatureVerifier,
	})

	// Initialize handlers
	healthHandler := health.NewHealthHandler(logger)
	projectHandler := project.NewProjectHandler(project.ProjectHandlerDep{
		ProjectService: projectService,
		Logger:         logger,
	})
	webhookHandler := webhook.NewWebhookHandler(webhookService, logger)

	// run APP in http server
	// inject all usecases here
	appService := app.Init(app.Dep{
		AppConfig:      cfg,
		HealthHandler:  healthHandler,
		ProjectHandler: projectHandler,
		WebhookHandler: webhookHandler,
		Logger:         logger,
	})
	appService.Run() // start http server
}
