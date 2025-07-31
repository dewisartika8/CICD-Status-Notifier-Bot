package dto

// CreateTelegramSubscriptionRequest represents the request to create a telegram subscription
type CreateTelegramSubscriptionRequest struct {
	ProjectID string `json:"project_id" validate:"required,uuid"`
	ChatID    int64  `json:"chat_id" validate:"required"`
}

// UpdateTelegramSubscriptionRequest represents the request to update a telegram subscription
type UpdateTelegramSubscriptionRequest struct {
	IsActive bool `json:"is_active"`
}

// TelegramSubscriptionResponse represents the response for telegram subscription operations
type TelegramSubscriptionResponse struct {
	ID        string `json:"id"`
	ProjectID string `json:"project_id"`
	ChatID    int64  `json:"chat_id"`
	IsActive  bool   `json:"is_active"`
	CreatedAt int64  `json:"created_at"`
	UpdatedAt int64  `json:"updated_at"`
}

// TelegramSubscriptionStatsResponse represents subscription statistics
type TelegramSubscriptionStatsResponse struct {
	Count     int64  `json:"count"`
	ProjectID string `json:"project_id,omitempty"`
	IsActive  *bool  `json:"is_active,omitempty"`
}

// TelegramSubscriptionListResponse represents a list of subscriptions with metadata
type TelegramSubscriptionListResponse struct {
	Subscriptions []TelegramSubscriptionResponse `json:"subscriptions"`
	Total         int64                          `json:"total"`
	ProjectID     string                         `json:"project_id,omitempty"`
}
