package notification

// Test constants for consistent testing across all notification test suites
const (
	// Type strings for mock assertions
	TelegramSubscriptionType = "*domain.TelegramSubscription"

	// Common error messages for testing
	DatabaseError        = "database error"
	SubscriptionNotFound = "subscription not found"
	InvalidUserID        = "invalid user ID"
	InvalidProjectID     = "invalid project ID"
	InvalidChatID        = "invalid chat ID"
	SubscriptionExists   = "subscription already exists"
	NetworkError         = "network error"
	ValidationError      = "validation error"

	// Test data constants
	TestChatID1     = int64(123456789)
	TestChatID2     = int64(987654321)
	TestUserID1     = int64(111111)
	TestUserID2     = int64(222222)
	TestUsername1   = "testuser1"
	TestUsername2   = "testuser2"
	TestMessage     = "Test notification message"
	TestSubject     = "Test Subject"
	TestRecipient   = "test@example.com"
	TestProjectName = "Test Project"
)
