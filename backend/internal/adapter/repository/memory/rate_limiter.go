package memory

import (
	"context"
	"sync"
	"time"

	"github.com/dewisartika8/cicd-status-notifier-bot/internal/core/notification/domain"
)

// inMemoryRateLimiter implements RateLimiter using in-memory storage
type inMemoryRateLimiter struct {
	entries map[string]*domain.RateLimitEntry
	rules   map[domain.NotificationChannel]*domain.RateLimitRule
	mutex   sync.RWMutex
}

// NewInMemoryRateLimiter creates a new in-memory rate limiter
func NewInMemoryRateLimiter() domain.RateLimiter {
	limiter := &inMemoryRateLimiter{
		entries: make(map[string]*domain.RateLimitEntry),
		rules:   make(map[domain.NotificationChannel]*domain.RateLimitRule),
	}

	// Initialize with default rules
	defaultRules := domain.DefaultRateLimitRules()
	for channel, rule := range defaultRules {
		limiter.rules[channel] = rule
	}

	return limiter
}

// Allow checks if a request is allowed for the given key and channel
func (r *inMemoryRateLimiter) Allow(ctx context.Context, key string, channel domain.NotificationChannel) (bool, error) {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	// Get rule for channel
	rule, exists := r.rules[channel]
	if !exists {
		// If no rule exists, allow the request
		return true, nil
	}

	entryKey := r.makeEntryKey(key, channel)
	entry, exists := r.entries[entryKey]

	if !exists {
		// First request for this key-channel combination
		entry = &domain.RateLimitEntry{
			Key:         entryKey,
			Channel:     channel,
			Count:       1,
			WindowStart: time.Now(),
			LastRequest: time.Now(),
		}
		r.entries[entryKey] = entry
		return true, nil
	}

	// Check if window has expired
	if entry.ShouldReset(rule.WindowSize) {
		entry.Reset()
		entry.IncrementCount()
		return true, nil
	}

	// Check if within limits
	if entry.Count >= rule.MaxRequests {
		return false, nil
	}

	// Allow and increment
	entry.IncrementCount()
	return true, nil
}

// GetRemainingRequests returns the number of remaining requests for the key and channel
func (r *inMemoryRateLimiter) GetRemainingRequests(ctx context.Context, key string, channel domain.NotificationChannel) (int, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	rule, exists := r.rules[channel]
	if !exists {
		return 0, nil
	}

	entryKey := r.makeEntryKey(key, channel)
	entry, exists := r.entries[entryKey]

	if !exists {
		return rule.MaxRequests, nil
	}

	if entry.ShouldReset(rule.WindowSize) {
		return rule.MaxRequests, nil
	}

	remaining := rule.MaxRequests - entry.Count
	if remaining < 0 {
		remaining = 0
	}

	return remaining, nil
}

// GetResetTime returns when the rate limit will reset for the key and channel
func (r *inMemoryRateLimiter) GetResetTime(ctx context.Context, key string, channel domain.NotificationChannel) (time.Time, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	rule, exists := r.rules[channel]
	if !exists {
		return time.Now(), nil
	}

	entryKey := r.makeEntryKey(key, channel)
	entry, exists := r.entries[entryKey]

	if !exists {
		return time.Now().Add(rule.WindowSize), nil
	}

	return entry.WindowStart.Add(rule.WindowSize), nil
}

// SetRule sets a rate limiting rule for a channel
func (r *inMemoryRateLimiter) SetRule(ctx context.Context, channel domain.NotificationChannel, maxRequests int, windowSize time.Duration, burstLimit int) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	r.rules[channel] = &domain.RateLimitRule{
		Channel:     channel,
		MaxRequests: maxRequests,
		WindowSize:  windowSize,
		BurstLimit:  burstLimit,
	}

	return nil
}

// GetRule gets the rate limiting rule for a channel
func (r *inMemoryRateLimiter) GetRule(ctx context.Context, channel domain.NotificationChannel) (*domain.RateLimitRule, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	rule, exists := r.rules[channel]
	if !exists {
		return nil, nil
	}

	return rule, nil
}

// RemoveRule removes the rate limiting rule for a channel
func (r *inMemoryRateLimiter) RemoveRule(ctx context.Context, channel domain.NotificationChannel) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	delete(r.rules, channel)
	return nil
}

// Reset resets the rate limit for a specific key and channel
func (r *inMemoryRateLimiter) Reset(ctx context.Context, key string, channel domain.NotificationChannel) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	entryKey := r.makeEntryKey(key, channel)
	entry, exists := r.entries[entryKey]

	if exists {
		entry.Reset()
	}

	return nil
}

// GetStats returns rate limiting statistics
func (r *inMemoryRateLimiter) GetStats(ctx context.Context, channel domain.NotificationChannel) (map[string]interface{}, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	stats := make(map[string]interface{})

	rule, exists := r.rules[channel]
	if exists {
		stats["max_requests"] = rule.MaxRequests
		stats["window_size"] = rule.WindowSize.String()
		stats["burst_limit"] = rule.BurstLimit
	}

	// Count active entries for this channel
	activeEntries := 0
	for _, entry := range r.entries {
		if entry.Channel == channel {
			activeEntries++
		}
	}

	stats["active_entries"] = activeEntries
	return stats, nil
}

// makeEntryKey creates a unique key for rate limit entries
func (r *inMemoryRateLimiter) makeEntryKey(key string, channel domain.NotificationChannel) string {
	return string(channel) + ":" + key
}
