package delivery

import (
	"context"
	"testing"
	"time"

	"github.com/dewisartika8/cicd-status-notifier-bot/internal/adapter/repository/memory"
	"github.com/dewisartika8/cicd-status-notifier-bot/internal/core/notification/domain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type RateLimiterTestSuite struct {
	suite.Suite
	rateLimiter domain.RateLimiter
	ctx         context.Context
}

func (suite *RateLimiterTestSuite) SetupTest() {
	suite.rateLimiter = memory.NewInMemoryRateLimiter()
	suite.ctx = context.Background()
}

func (suite *RateLimiterTestSuite) TestAllowFirstRequest() {
	// Act
	allowed, err := suite.rateLimiter.Allow(suite.ctx, "user1", domain.NotificationChannelTelegram)

	// Assert
	assert.NoError(suite.T(), err)
	assert.True(suite.T(), allowed)
}

func (suite *RateLimiterTestSuite) TestAllowWithinLimits() {
	// Arrange
	channel := domain.NotificationChannelTelegram
	key := "user1"

	// Act - Make multiple requests within limit
	for i := 0; i < 5; i++ {
		allowed, err := suite.rateLimiter.Allow(suite.ctx, key, channel)
		assert.NoError(suite.T(), err)
		assert.True(suite.T(), allowed)
	}
}

func (suite *RateLimiterTestSuite) TestDenyExceedsLimits() {
	// Arrange
	channel := domain.NotificationChannelEmail // Has lower limit (10 per minute)
	key := "user1"

	// Act - Make requests up to limit
	for i := 0; i < 10; i++ {
		allowed, err := suite.rateLimiter.Allow(suite.ctx, key, channel)
		assert.NoError(suite.T(), err)
		assert.True(suite.T(), allowed)
	}

	// This should be denied
	allowed, err := suite.rateLimiter.Allow(suite.ctx, key, channel)
	assert.NoError(suite.T(), err)
	assert.False(suite.T(), allowed)
}

func (suite *RateLimiterTestSuite) TestGetRemainingRequests() {
	// Arrange
	channel := domain.NotificationChannelEmail // 10 requests per minute
	key := "user1"

	// Make 3 requests
	for i := 0; i < 3; i++ {
		suite.rateLimiter.Allow(suite.ctx, key, channel)
	}

	// Act
	remaining, err := suite.rateLimiter.GetRemainingRequests(suite.ctx, key, channel)

	// Assert
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), 7, remaining) // 10 - 3 = 7
}

func (suite *RateLimiterTestSuite) TestSetAndGetRule() {
	// Arrange
	channel := domain.NotificationChannelWebhook
	expectedRule := &domain.RateLimitRule{
		Channel:     channel,
		MaxRequests: 50,
		WindowSize:  time.Minute * 5,
		BurstLimit:  10,
	}

	// Act
	err := suite.rateLimiter.SetRule(suite.ctx, channel, 50, time.Minute*5, 10)
	assert.NoError(suite.T(), err)

	rule, err := suite.rateLimiter.GetRule(suite.ctx, channel)

	// Assert
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), expectedRule.MaxRequests, rule.MaxRequests)
	assert.Equal(suite.T(), expectedRule.WindowSize, rule.WindowSize)
	assert.Equal(suite.T(), expectedRule.BurstLimit, rule.BurstLimit)
}

func (suite *RateLimiterTestSuite) TestResetRateLimit() {
	// Arrange
	channel := domain.NotificationChannelEmail
	key := "user1"

	// Make maximum requests
	for i := 0; i < 10; i++ {
		suite.rateLimiter.Allow(suite.ctx, key, channel)
	}

	// Verify we're at limit
	allowed, _ := suite.rateLimiter.Allow(suite.ctx, key, channel)
	assert.False(suite.T(), allowed)

	// Act - Reset
	err := suite.rateLimiter.Reset(suite.ctx, key, channel)
	assert.NoError(suite.T(), err)

	// Assert - Should be allowed again
	allowed, err = suite.rateLimiter.Allow(suite.ctx, key, channel)
	assert.NoError(suite.T(), err)
	assert.True(suite.T(), allowed)
}

func (suite *RateLimiterTestSuite) TestGetStats() {
	// Arrange
	channel := domain.NotificationChannelTelegram
	suite.rateLimiter.Allow(suite.ctx, "user1", channel)
	suite.rateLimiter.Allow(suite.ctx, "user2", channel)

	// Act
	stats, err := suite.rateLimiter.GetStats(suite.ctx, channel)

	// Assert
	assert.NoError(suite.T(), err)
	assert.Contains(suite.T(), stats, "max_requests")
	assert.Contains(suite.T(), stats, "window_size")
	assert.Contains(suite.T(), stats, "burst_limit")
	assert.Contains(suite.T(), stats, "active_entries")
	assert.Equal(suite.T(), 2, stats["active_entries"])
}

func TestRateLimiterTestSuite(t *testing.T) {
	suite.Run(t, new(RateLimiterTestSuite))
}

// Test default rate limit rules
func TestDefaultRateLimitRulesLoaded(t *testing.T) {
	rateLimiter := memory.NewInMemoryRateLimiter()
	ctx := context.Background()

	// Test Telegram rule
	rule, err := rateLimiter.GetRule(ctx, domain.NotificationChannelTelegram)
	assert.NoError(t, err)
	assert.NotNil(t, rule)
	assert.Equal(t, 30, rule.MaxRequests)
	assert.Equal(t, time.Minute, rule.WindowSize)

	// Test Email rule
	rule, err = rateLimiter.GetRule(ctx, domain.NotificationChannelEmail)
	assert.NoError(t, err)
	assert.NotNil(t, rule)
	assert.Equal(t, 10, rule.MaxRequests)
}
