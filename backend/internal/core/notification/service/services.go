package service

// Package level constructors for easy import and backward compatibility

import (
	"github.com/dewisartika8/cicd-status-notifier-bot/internal/core/notification/service/formatter"
	"github.com/dewisartika8/cicd-status-notifier-bot/internal/core/notification/service/log"
	"github.com/dewisartika8/cicd-status-notifier-bot/internal/core/notification/service/retry"
	"github.com/dewisartika8/cicd-status-notifier-bot/internal/core/notification/service/sender"
	"github.com/dewisartika8/cicd-status-notifier-bot/internal/core/notification/service/subscription"
	"github.com/dewisartika8/cicd-status-notifier-bot/internal/core/notification/service/template"
)

// Type aliases for backward compatibility
type (
	// Notification Log Service
	NotificationLogDep = log.Dep

	// Notification Sender Service
	NotificationSenderDep = sender.Dep

	// Notification Template Service
	NotificationTemplateDep = template.Dep

	// Notification Formatter Service
	NotificationFormatterDep = formatter.Dep

	// Retry Service
	RetryDep = retry.Dep

	// Telegram Subscription Service
	TelegramSubscriptionDep = subscription.Dep
)

// Constructor aliases for backward compatibility
var (
	NewNotificationLogService       = log.NewNotificationLogService
	NewNotificationSenderService    = sender.NewNotificationSenderService
	NewNotificationTemplateService  = template.NewNotificationTemplateService
	NewNotificationFormatterService = formatter.NewNotificationFormatterService
	NewRetryService                 = retry.NewRetryService
	NewTelegramSubscriptionService  = subscription.NewTelegramSubscriptionService
)
