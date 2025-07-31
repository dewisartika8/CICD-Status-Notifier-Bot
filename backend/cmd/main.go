package main

import (
	"log"
	"time"

	"github.com/dewisartika8/cicd-status-notifier-bot/internal/adapter/handler/health"
	"github.com/dewisartika8/cicd-status-notifier-bot/internal/adapter/handler/project"
	"github.com/dewisartika8/cicd-status-notifier-bot/internal/adapter/handler/telegram"
	"github.com/dewisartika8/cicd-status-notifier-bot/internal/adapter/handler/webhook"
	"github.com/dewisartika8/cicd-status-notifier-bot/internal/adapter/repository/postgres"
	"github.com/dewisartika8/cicd-status-notifier-bot/internal/config"
	bs "github.com/dewisartika8/cicd-status-notifier-bot/internal/core/build/service"
	subscription "github.com/dewisartika8/cicd-status-notifier-bot/internal/core/notification/service/subscription"
	ps "github.com/dewisartika8/cicd-status-notifier-bot/internal/core/project/service"
	ws "github.com/dewisartika8/cicd-status-notifier-bot/internal/core/webhook/service"
	"github.com/dewisartika8/cicd-status-notifier-bot/internal/server/app"
	"github.com/dewisartika8/cicd-status-notifier-bot/pkg/crypto"
	"github.com/dewisartika8/cicd-status-notifier-bot/pkg/database"
	loggerPkg "github.com/dewisartika8/cicd-status-notifier-bot/pkg/logger"
)

func main() {
	// Load configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Config error: %v", err)
	}

	// Initialize logger
	logger := loggerPkg.NewLogger()
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
	projectRepo := postgres.NewProjectRepository(db)
	buildEventRepo := postgres.NewBuildEventRepository(db)
	webhookEventRepo := postgres.NewWebhookEventRepository(db)
	telegramSubscriptionRepo := postgres.NewTelegramSubscriptionRepository(db)

	// Initialize services
	projectService := ps.NewProjectService(ps.Dep{
		ProjectRepo: projectRepo,
	})
	buildService := bs.NewBuildEventService(bs.Dep{
		BuildEventRepo: buildEventRepo,
	})

	// Initialize telegram subscription service
	telegramSubscriptionService := subscription.NewTelegramSubscriptionService(subscription.Dep{
		TelegramRepo: telegramSubscriptionRepo,
		Logger:       logger,
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
	telegramHandler := telegram.NewTelegramHandler(cfg, telegramSubscriptionService, logger)

	// run APP in http server
	// inject all usecases here
	appService := app.Init(app.Dep{
		AppConfig:       cfg,
		HealthHandler:   healthHandler,
		ProjectHandler:  projectHandler,
		WebhookHandler:  webhookHandler,
		TelegramHandler: telegramHandler,
		Logger:          logger,
	})
	appService.Run() // start http server
}
