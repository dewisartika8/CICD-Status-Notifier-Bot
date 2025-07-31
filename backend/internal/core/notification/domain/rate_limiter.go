package domain

import (
	"context"
	"time"
)

// RateLimitRule represents a rate limiting rule
type RateLimitRule struct {
	Channel     NotificationChannel `json:"channel"`
	MaxRequests int                 `json:"max_requests"`
	WindowSize  time.Duration       `json:"window_size"`
	BurstLimit  int                 `json:"burst_limit"`
}

// RateLimitEntry represents a rate limit entry for tracking
type RateLimitEntry struct {
	Key         string              `json:"key"`
	Channel     NotificationChannel `json:"channel"`
	Count       int                 `json:"count"`
	WindowStart time.Time           `json:"window_start"`
	LastRequest time.Time           `json:"last_request"`
}

// IsExpired checks if the rate limit window has expired
func (rle *RateLimitEntry) IsExpired(windowSize time.Duration) bool {
	return time.Now().After(rle.WindowStart.Add(windowSize))
}

// ShouldReset checks if the rate limit should be reset
func (rle *RateLimitEntry) ShouldReset(windowSize time.Duration) bool {
	return rle.IsExpired(windowSize)
}

// IncrementCount increments the request count
func (rle *RateLimitEntry) IncrementCount() {
	rle.Count++
	rle.LastRequest = time.Now()
}

// Reset resets the rate limit entry
func (rle *RateLimitEntry) Reset() {
	rle.Count = 0
	rle.WindowStart = time.Now()
	rle.LastRequest = time.Now()
}

// RateLimiter defines the interface for rate limiting operations
type RateLimiter interface {
	// Allow checks if a request is allowed for the given key and channel
	Allow(ctx context.Context, key string, channel NotificationChannel) (bool, error)

	// GetRemainingRequests returns the number of remaining requests for the key and channel
	GetRemainingRequests(ctx context.Context, key string, channel NotificationChannel) (int, error)

	// GetResetTime returns when the rate limit will reset for the key and channel
	GetResetTime(ctx context.Context, key string, channel NotificationChannel) (time.Time, error)

	// SetRule sets a rate limiting rule for a channel
	SetRule(ctx context.Context, channel NotificationChannel, maxRequests int, windowSize time.Duration, burstLimit int) error

	// GetRule gets the rate limiting rule for a channel
	GetRule(ctx context.Context, channel NotificationChannel) (*RateLimitRule, error)

	// RemoveRule removes the rate limiting rule for a channel
	RemoveRule(ctx context.Context, channel NotificationChannel) error

	// Reset resets the rate limit for a specific key and channel
	Reset(ctx context.Context, key string, channel NotificationChannel) error

	// GetStats returns rate limiting statistics
	GetStats(ctx context.Context, channel NotificationChannel) (map[string]interface{}, error)
}

// DefaultRateLimitRules returns default rate limiting rules for different channels
func DefaultRateLimitRules() map[NotificationChannel]*RateLimitRule {
	return map[NotificationChannel]*RateLimitRule{
		NotificationChannelTelegram: {
			Channel:     NotificationChannelTelegram,
			MaxRequests: 30,          // 30 messages per minute
			WindowSize:  time.Minute, // 1 minute window
			BurstLimit:  5,           // Allow burst of 5 messages
		},
		NotificationChannelEmail: {
			Channel:     NotificationChannelEmail,
			MaxRequests: 10,          // 10 emails per minute
			WindowSize:  time.Minute, // 1 minute window
			BurstLimit:  3,           // Allow burst of 3 emails
		},
		NotificationChannelSlack: {
			Channel:     NotificationChannelSlack,
			MaxRequests: 50,          // 50 messages per minute
			WindowSize:  time.Minute, // 1 minute window
			BurstLimit:  10,          // Allow burst of 10 messages
		},
		NotificationChannelWebhook: {
			Channel:     NotificationChannelWebhook,
			MaxRequests: 100,         // 100 webhooks per minute
			WindowSize:  time.Minute, // 1 minute window
			BurstLimit:  20,          // Allow burst of 20 webhooks
		},
	}
}
