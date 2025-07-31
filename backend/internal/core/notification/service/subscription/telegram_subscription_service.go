package subscription

import (
	"context"
	"fmt"

	"github.com/dewisartika8/cicd-status-notifier-bot/internal/core/notification/domain"
	"github.com/dewisartika8/cicd-status-notifier-bot/internal/core/notification/port"
	"github.com/dewisartika8/cicd-status-notifier-bot/internal/core/shared/domain/value_objects"
	"github.com/sirupsen/logrus"
)

// Resource type constants
const (
	resourceSubscription = "telegram subscription"
)

type Dep struct {
	// put your dependencies here
	TelegramRepo port.TelegramSubscriptionRepository
	Logger       *logrus.Logger
}

// telegramSubscriptionService implements telegram subscription business logic
type telegramSubscriptionService struct {
	Dep
}

// NewTelegramSubscriptionService creates a new telegram subscription service
func NewTelegramSubscriptionService(d Dep) port.TelegramSubscriptionService {
	return &telegramSubscriptionService{
		Dep: d,
	}
}

// CreateTelegramSubscription creates a new telegram subscription
func (s *telegramSubscriptionService) CreateTelegramSubscription(
	ctx context.Context,
	projectID value_objects.ID,
	chatID int64,
) (*domain.TelegramSubscription, error) {
	s.Logger.WithFields(logrus.Fields{
		"project_id": projectID.String(),
		"chat_id":    chatID,
	}).Info("Creating telegram subscription")

	// Check if subscription already exists
	exists, err := s.TelegramRepo.ExistsByProjectAndChatID(ctx, projectID, chatID)
	if err != nil {
		s.Logger.WithError(err).Error("Failed to check if subscription exists")
		return nil, fmt.Errorf("failed to check if subscription exists: %w", err)
	}

	if exists {
		s.Logger.Warn("Subscription already exists")
		return nil, domain.ErrSubscriptionAlreadyExists
	}

	// Create new telegram subscription entity
	subscription, err := domain.NewTelegramSubscription(projectID, chatID)
	if err != nil {
		s.Logger.WithError(err).Error("Failed to create telegram subscription entity")
		return nil, fmt.Errorf(domain.ErrMsgCreate, resourceSubscription, err)
	}

	// Persist the subscription
	if err := s.TelegramRepo.Create(ctx, subscription); err != nil {
		s.Logger.WithError(err).Error("Failed to persist telegram subscription")
		return nil, fmt.Errorf("failed to persist telegram subscription: %w", err)
	}

	s.Logger.WithField("subscription_id", subscription.ID().String()).Info("Telegram subscription created successfully")
	return subscription, nil
}

// GetTelegramSubscription retrieves a telegram subscription by its ID
func (s *telegramSubscriptionService) GetTelegramSubscription(ctx context.Context, id value_objects.ID) (*domain.TelegramSubscription, error) {
	s.Logger.WithField("id", id.String()).Info("Getting telegram subscription")

	subscription, err := s.TelegramRepo.GetByID(ctx, id)
	if err != nil {
		s.Logger.WithError(err).Error(domain.LogMsgGetSubscription)
		return nil, fmt.Errorf(domain.ErrMsgGet, resourceSubscription, err)
	}

	return subscription, nil
}

// GetTelegramSubscriptionsByProject retrieves telegram subscriptions for a project
func (s *telegramSubscriptionService) GetTelegramSubscriptionsByProject(ctx context.Context, projectID value_objects.ID) ([]*domain.TelegramSubscription, error) {
	s.Logger.WithField("project_id", projectID.String()).Info("Getting telegram subscriptions by project")

	subscriptions, err := s.TelegramRepo.GetByProjectID(ctx, projectID)
	if err != nil {
		s.Logger.WithError(err).Error("Failed to get telegram subscriptions by project")
		return nil, fmt.Errorf("failed to get telegram subscriptions by project: %w", err)
	}

	return subscriptions, nil
}

// GetTelegramSubscriptionByChatID retrieves a telegram subscription by chat ID
func (s *telegramSubscriptionService) GetTelegramSubscriptionByChatID(ctx context.Context, chatID int64) (*domain.TelegramSubscription, error) {
	s.Logger.WithField("chat_id", chatID).Info("Getting telegram subscription by chat ID")

	subscription, err := s.TelegramRepo.GetByChatID(ctx, chatID)
	if err != nil {
		s.Logger.WithError(err).Error("Failed to get telegram subscription by chat ID")
		return nil, fmt.Errorf("failed to get telegram subscription by chat ID: %w", err)
	}

	return subscription, nil
}

// UpdateTelegramSubscription updates a telegram subscription
func (s *telegramSubscriptionService) UpdateTelegramSubscription(
	ctx context.Context,
	id value_objects.ID,
	chatID *int64,
	isActive *bool,
) (*domain.TelegramSubscription, error) {
	s.Logger.WithField("id", id.String()).Info("Updating telegram subscription")

	// Get the subscription
	subscription, err := s.TelegramRepo.GetByID(ctx, id)
	if err != nil {
		s.Logger.WithError(err).Error(domain.LogMsgGetSubscription)
		return nil, fmt.Errorf(domain.ErrMsgGet, resourceSubscription, err)
	}

	// Update chat ID if provided
	if chatID != nil {
		if err := subscription.UpdateChatID(*chatID); err != nil {
			s.Logger.WithError(err).Error("Failed to update chat ID")
			return nil, fmt.Errorf("failed to update chat ID: %w", err)
		}
	}

	// Update active status if provided
	if isActive != nil {
		if *isActive {
			subscription.Activate()
		} else {
			subscription.Deactivate()
		}
	}

	// Update the subscription in repository
	if err := s.TelegramRepo.Update(ctx, subscription); err != nil {
		s.Logger.WithError(err).Error(domain.LogMsgUpdateSubscription)
		return nil, fmt.Errorf(domain.ErrMsgUpdate, resourceSubscription, err)
	}

	s.Logger.Info("Telegram subscription updated successfully")
	return subscription, nil
}

// DeleteTelegramSubscription deletes a telegram subscription
func (s *telegramSubscriptionService) DeleteTelegramSubscription(ctx context.Context, id value_objects.ID) error {
	s.Logger.WithField("id", id.String()).Info("Deleting telegram subscription")

	if err := s.TelegramRepo.Delete(ctx, id); err != nil {
		s.Logger.WithError(err).Error("Failed to delete telegram subscription")
		return fmt.Errorf("failed to delete telegram subscription: %w", err)
	}

	s.Logger.Info("Telegram subscription deleted successfully")
	return nil
}

// ActivateTelegramSubscription activates a telegram subscription
func (s *telegramSubscriptionService) ActivateTelegramSubscription(ctx context.Context, id value_objects.ID) error {
	s.Logger.WithField("id", id.String()).Info("Activating telegram subscription")

	subscription, err := s.TelegramRepo.GetByID(ctx, id)
	if err != nil {
		s.Logger.WithError(err).Error(domain.LogMsgGetSubscription)
		return fmt.Errorf(domain.ErrMsgGet, resourceSubscription, err)
	}

	subscription.Activate()

	if err := s.TelegramRepo.Update(ctx, subscription); err != nil {
		s.Logger.WithError(err).Error(domain.LogMsgUpdateSubscription)
		return fmt.Errorf(domain.ErrMsgUpdate, resourceSubscription, err)
	}

	s.Logger.Info("Telegram subscription activated successfully")
	return nil
}

// DeactivateTelegramSubscription deactivates a telegram subscription
func (s *telegramSubscriptionService) DeactivateTelegramSubscription(ctx context.Context, id value_objects.ID) error {
	s.Logger.WithField("id", id.String()).Info("Deactivating telegram subscription")

	subscription, err := s.TelegramRepo.GetByID(ctx, id)
	if err != nil {
		s.Logger.WithError(err).Error(domain.LogMsgGetSubscription)
		return fmt.Errorf(domain.ErrMsgGet, resourceSubscription, err)
	}

	subscription.Deactivate()

	if err := s.TelegramRepo.Update(ctx, subscription); err != nil {
		s.Logger.WithError(err).Error(domain.LogMsgUpdateSubscription)
		return fmt.Errorf(domain.ErrMsgUpdate, resourceSubscription, err)
	}

	s.Logger.Info("Telegram subscription deactivated successfully")
	return nil
}

// GetActiveSubscriptionsForProject retrieves active telegram subscriptions for a project
func (s *telegramSubscriptionService) GetActiveSubscriptionsForProject(ctx context.Context, projectID value_objects.ID) ([]*domain.TelegramSubscription, error) {
	s.Logger.WithField("project_id", projectID.String()).Info("Getting active subscriptions for project")

	subscriptions, err := s.TelegramRepo.GetActiveSubscriptionsByProject(ctx, projectID)
	if err != nil {
		s.Logger.WithError(err).Error("Failed to get active subscriptions for project")
		return nil, fmt.Errorf("failed to get active subscriptions for project: %w", err)
	}

	return subscriptions, nil
}

// GetAllActiveSubscriptions retrieves all active telegram subscriptions
func (s *telegramSubscriptionService) GetAllActiveSubscriptions(ctx context.Context) ([]*domain.TelegramSubscription, error) {
	s.Logger.Info("Getting all active subscriptions")

	subscriptions, err := s.TelegramRepo.GetActiveSubscriptions(ctx)
	if err != nil {
		s.Logger.WithError(err).Error("Failed to get all active subscriptions")
		return nil, fmt.Errorf("failed to get all active subscriptions: %w", err)
	}

	return subscriptions, nil
}

// CheckSubscriptionExists checks if a subscription exists for a project and chat ID
func (s *telegramSubscriptionService) CheckSubscriptionExists(ctx context.Context, projectID value_objects.ID, chatID int64) (bool, error) {
	s.Logger.WithFields(logrus.Fields{
		"project_id": projectID.String(),
		"chat_id":    chatID,
	}).Info("Checking if subscription exists")

	exists, err := s.TelegramRepo.ExistsByProjectAndChatID(ctx, projectID, chatID)
	if err != nil {
		s.Logger.WithError(err).Error("Failed to check if subscription exists")
		return false, fmt.Errorf("failed to check if subscription exists: %w", err)
	}

	return exists, nil
}
