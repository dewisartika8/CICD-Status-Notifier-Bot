package subscription

import (
	"context"

	"github.com/dewisartika8/cicd-status-notifier-bot/internal/core/shared/domain/value_objects"
)

// ValidationService defines methods for subscription validation
type ValidationService interface {
	// ValidateUserPermissions validates if user has permission to create subscription
	ValidateUserPermissions(ctx context.Context, userID int64, projectID value_objects.ID, chatID int64) error

	// ValidateProjectExistence validates if project exists
	ValidateProjectExistence(ctx context.Context, projectID value_objects.ID) error

	// ValidateDuplicateSubscription validates if subscription already exists
	ValidateDuplicateSubscription(ctx context.Context, projectID value_objects.ID, chatID int64) error

	// ValidateChatID validates Telegram chat ID format
	ValidateChatID(ctx context.Context, chatID int64) error

	// ValidateSubscriptionParameters validates all subscription parameters
	ValidateSubscriptionParameters(ctx context.Context, projectID value_objects.ID, chatID int64, userID *int64) error
}
