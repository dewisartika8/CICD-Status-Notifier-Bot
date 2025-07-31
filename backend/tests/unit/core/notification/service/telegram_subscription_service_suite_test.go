package service_test

import (
	"context"
	"errors"
	"testing"

	"github.com/dewisartika8/cicd-status-notifier-bot/internal/core/notification/domain"
	"github.com/dewisartika8/cicd-status-notifier-bot/internal/core/shared/domain/value_objects"
	testmocks "github.com/dewisartika8/cicd-status-notifier-bot/tests/testutils/mocks/notification"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

// TelegramSubscriptionServiceTestSuite defines the test suite for Telegram subscription service
type TelegramSubscriptionServiceTestSuite struct {
	suite.Suite
	helpers   *testmocks.TestHelpers
	projectID value_objects.ID
	chatID1   int64
	chatID2   int64
}

// SetupSuite runs once before all tests in the suite
func (suite *TelegramSubscriptionServiceTestSuite) SetupSuite() {
	suite.helpers = testmocks.NewTestHelpers()
	suite.projectID = value_objects.NewID()
	suite.chatID1 = testmocks.TestChatID1
	suite.chatID2 = testmocks.TestChatID2
}

// TestCreateTelegramSubscription tests subscription creation scenarios
func (suite *TelegramSubscriptionServiceTestSuite) TestCreateTelegramSubscription() {
	tests := []struct {
		name          string
		setupMocks    func(*testmocks.MockTelegramRepository)
		expectedError string
		shouldSucceed bool
	}{
		{
			name: "successful_creation",
			setupMocks: func(mockRepo *testmocks.MockTelegramRepository) {
				mockRepo.On("ExistsByProjectAndChatID", mock.Anything, suite.projectID, suite.chatID1).Return(false, nil)
				mockRepo.On("Create", mock.Anything, mock.AnythingOfType(testmocks.TelegramSubscriptionType)).Return(nil)
			},
			shouldSucceed: true,
		},
		{
			name: "subscription_already_exists",
			setupMocks: func(mockRepo *testmocks.MockTelegramRepository) {
				mockRepo.On("ExistsByProjectAndChatID", mock.Anything, suite.projectID, suite.chatID1).Return(true, nil)
			},
			expectedError: "subscription already exists",
			shouldSucceed: false,
		},
		{
			name: "repository_error_on_exists_check",
			setupMocks: func(mockRepo *testmocks.MockTelegramRepository) {
				mockRepo.On("ExistsByProjectAndChatID", mock.Anything, suite.projectID, suite.chatID1).
					Return(false, errors.New(testmocks.DatabaseError))
			},
			expectedError: "failed to check subscription existence",
			shouldSucceed: false,
		},
		{
			name: "repository_error_on_create",
			setupMocks: func(mockRepo *testmocks.MockTelegramRepository) {
				mockRepo.On("ExistsByProjectAndChatID", mock.Anything, suite.projectID, suite.chatID1).Return(false, nil)
				mockRepo.On("Create", mock.Anything, mock.AnythingOfType(testmocks.TelegramSubscriptionType)).
					Return(errors.New(testmocks.DatabaseError))
			},
			expectedError: "failed to persist telegram subscription",
			shouldSucceed: false,
		},
	}

	for _, tt := range tests {
		suite.Run(tt.name, func() {
			service, mockRepo := suite.helpers.CreateTelegramSubscriptionService()
			tt.setupMocks(mockRepo)

			result, err := service.CreateTelegramSubscription(context.Background(), suite.projectID, suite.chatID1)

			if tt.shouldSucceed {
				require.NoError(suite.T(), err)
				assert.NotNil(suite.T(), result)
				assert.Equal(suite.T(), suite.projectID, result.ProjectID())
				assert.Equal(suite.T(), suite.chatID1, result.ChatID())
				assert.True(suite.T(), result.IsActive())
			} else {
				assert.Error(suite.T(), err)
				assert.Nil(suite.T(), result)
				if tt.expectedError != "" {
					assert.Contains(suite.T(), err.Error(), tt.expectedError)
				}
			}

			mockRepo.AssertExpectations(suite.T())
		})
	}
}

// TestGetTelegramSubscription tests subscription retrieval scenarios
func (suite *TelegramSubscriptionServiceTestSuite) TestGetTelegramSubscription() {
	subscriptionID := value_objects.NewID()

	tests := []struct {
		name          string
		setupMocks    func(*testmocks.MockTelegramRepository)
		expectedError string
		shouldSucceed bool
	}{
		{
			name: "successful_retrieval",
			setupMocks: func(mockRepo *testmocks.MockTelegramRepository) {
				expectedSub, _ := domain.NewTelegramSubscription(suite.projectID, suite.chatID1)
				mockRepo.On("GetByID", mock.Anything, subscriptionID).Return(expectedSub, nil)
			},
			shouldSucceed: true,
		},
		{
			name: "subscription_not_found",
			setupMocks: func(mockRepo *testmocks.MockTelegramRepository) {
				mockRepo.On("GetByID", mock.Anything, subscriptionID).
					Return(nil, errors.New(testmocks.SubscriptionNotFound))
			},
			expectedError: "failed to get telegram subscription",
			shouldSucceed: false,
		},
	}

	for _, tt := range tests {
		suite.Run(tt.name, func() {
			service, mockRepo := suite.helpers.CreateTelegramSubscriptionService()
			tt.setupMocks(mockRepo)

			result, err := service.GetTelegramSubscription(context.Background(), subscriptionID)

			if tt.shouldSucceed {
				require.NoError(suite.T(), err)
				assert.NotNil(suite.T(), result)
				assert.Equal(suite.T(), suite.projectID, result.ProjectID())
				assert.Equal(suite.T(), suite.chatID1, result.ChatID())
			} else {
				assert.Error(suite.T(), err)
				assert.Nil(suite.T(), result)
				if tt.expectedError != "" {
					assert.Contains(suite.T(), err.Error(), tt.expectedError)
				}
			}

			mockRepo.AssertExpectations(suite.T())
		})
	}
}

// TestUpdateTelegramSubscription tests subscription update scenarios
func (suite *TelegramSubscriptionServiceTestSuite) TestUpdateTelegramSubscription() {
	subscriptionID := value_objects.NewID()
	newChatID := suite.chatID2
	isActive := false

	tests := []struct {
		name          string
		setupMocks    func(*testmocks.MockTelegramRepository)
		expectedError string
		shouldSucceed bool
	}{
		{
			name: "successful_update",
			setupMocks: func(mockRepo *testmocks.MockTelegramRepository) {
				originalSub, _ := domain.NewTelegramSubscription(suite.projectID, suite.chatID1)
				mockRepo.On("GetByID", mock.Anything, subscriptionID).Return(originalSub, nil)
				mockRepo.On("Update", mock.Anything, mock.AnythingOfType(testmocks.TelegramSubscriptionType)).Return(nil)
			},
			shouldSucceed: true,
		},
		{
			name: "subscription_not_found",
			setupMocks: func(mockRepo *testmocks.MockTelegramRepository) {
				mockRepo.On("GetByID", mock.Anything, subscriptionID).
					Return(nil, errors.New(testmocks.SubscriptionNotFound))
			},
			expectedError: "failed to get telegram subscription",
			shouldSucceed: false,
		},
		{
			name: "repository_error_on_update",
			setupMocks: func(mockRepo *testmocks.MockTelegramRepository) {
				originalSub, _ := domain.NewTelegramSubscription(suite.projectID, suite.chatID1)
				mockRepo.On("GetByID", mock.Anything, subscriptionID).Return(originalSub, nil)
				mockRepo.On("Update", mock.Anything, mock.AnythingOfType(testmocks.TelegramSubscriptionType)).
					Return(errors.New(testmocks.DatabaseError))
			},
			expectedError: "failed to update telegram subscription",
			shouldSucceed: false,
		},
	}

	for _, tt := range tests {
		suite.Run(tt.name, func() {
			service, mockRepo := suite.helpers.CreateTelegramSubscriptionService()
			tt.setupMocks(mockRepo)

			result, err := service.UpdateTelegramSubscription(context.Background(), subscriptionID, &newChatID, &isActive)

			if tt.shouldSucceed {
				require.NoError(suite.T(), err)
				assert.NotNil(suite.T(), result)
				assert.Equal(suite.T(), newChatID, result.ChatID())
				assert.Equal(suite.T(), isActive, result.IsActive())
			} else {
				assert.Error(suite.T(), err)
				assert.Nil(suite.T(), result)
				if tt.expectedError != "" {
					assert.Contains(suite.T(), err.Error(), tt.expectedError)
				}
			}

			mockRepo.AssertExpectations(suite.T())
		})
	}
}

// TestDeleteTelegramSubscription tests subscription deletion scenarios
func (suite *TelegramSubscriptionServiceTestSuite) TestDeleteTelegramSubscription() {
	subscriptionID := value_objects.NewID()

	tests := []struct {
		name          string
		setupMocks    func(*testmocks.MockTelegramRepository)
		expectedError string
		shouldSucceed bool
	}{
		{
			name: "successful_deletion",
			setupMocks: func(mockRepo *testmocks.MockTelegramRepository) {
				mockRepo.On("Delete", mock.Anything, subscriptionID).Return(nil)
			},
			shouldSucceed: true,
		},
		{
			name: "repository_error",
			setupMocks: func(mockRepo *testmocks.MockTelegramRepository) {
				mockRepo.On("Delete", mock.Anything, subscriptionID).
					Return(errors.New(testmocks.DatabaseError))
			},
			expectedError: "failed to delete telegram subscription",
			shouldSucceed: false,
		},
	}

	for _, tt := range tests {
		suite.Run(tt.name, func() {
			service, mockRepo := suite.helpers.CreateTelegramSubscriptionService()
			tt.setupMocks(mockRepo)

			err := service.DeleteTelegramSubscription(context.Background(), subscriptionID)

			if tt.shouldSucceed {
				assert.NoError(suite.T(), err)
			} else {
				assert.Error(suite.T(), err)
				if tt.expectedError != "" {
					assert.Contains(suite.T(), err.Error(), tt.expectedError)
				}
			}

			mockRepo.AssertExpectations(suite.T())
		})
	}
}

// TestGetTelegramSubscriptionsByProject tests getting subscriptions by project
func (suite *TelegramSubscriptionServiceTestSuite) TestGetTelegramSubscriptionsByProject() {
	tests := []struct {
		name          string
		setupMocks    func(*testmocks.MockTelegramRepository)
		expectedError string
		shouldSucceed bool
		expectedCount int
	}{
		{
			name: "successful_retrieval",
			setupMocks: func(mockRepo *testmocks.MockTelegramRepository) {
				sub1, _ := domain.NewTelegramSubscription(suite.projectID, suite.chatID1)
				sub2, _ := domain.NewTelegramSubscription(suite.projectID, suite.chatID2)
				expectedSubs := []*domain.TelegramSubscription{sub1, sub2}
				mockRepo.On("GetByProjectID", mock.Anything, suite.projectID).Return(expectedSubs, nil)
			},
			shouldSucceed: true,
			expectedCount: 2,
		},
		{
			name: "repository_error",
			setupMocks: func(mockRepo *testmocks.MockTelegramRepository) {
				mockRepo.On("GetByProjectID", mock.Anything, suite.projectID).
					Return(nil, errors.New(testmocks.DatabaseError))
			},
			expectedError: "failed to get telegram subscriptions by project",
			shouldSucceed: false,
		},
		{
			name: "empty_result",
			setupMocks: func(mockRepo *testmocks.MockTelegramRepository) {
				mockRepo.On("GetByProjectID", mock.Anything, suite.projectID).
					Return([]*domain.TelegramSubscription{}, nil)
			},
			shouldSucceed: true,
			expectedCount: 0,
		},
	}

	for _, tt := range tests {
		suite.Run(tt.name, func() {
			service, mockRepo := suite.helpers.CreateTelegramSubscriptionService()
			tt.setupMocks(mockRepo)

			result, err := service.GetTelegramSubscriptionsByProject(context.Background(), suite.projectID)

			if tt.shouldSucceed {
				require.NoError(suite.T(), err)
				assert.Len(suite.T(), result, tt.expectedCount)
			} else {
				assert.Error(suite.T(), err)
				assert.Nil(suite.T(), result)
				if tt.expectedError != "" {
					assert.Contains(suite.T(), err.Error(), tt.expectedError)
				}
			}

			mockRepo.AssertExpectations(suite.T())
		})
	}
}

// TestValidation tests validation methods
func (suite *TelegramSubscriptionServiceTestSuite) TestValidation() {
	suite.Run("ValidateUserPermissions", func() {
		service, _ := suite.helpers.CreateValidationService()

		tests := []struct {
			name          string
			userID        int64
			expectedError string
		}{
			{
				name:   "valid_user_permissions",
				userID: testmocks.TestUserID1,
			},
			{
				name:          "invalid_user_id_zero",
				userID:        0,
				expectedError: "invalid user ID",
			},
			{
				name:          "invalid_user_id_negative",
				userID:        -1,
				expectedError: "invalid user ID",
			},
		}

		for _, tt := range tests {
			suite.Run(tt.name, func() {
				err := service.ValidateUserPermissions(context.Background(), tt.userID, suite.projectID, suite.chatID1)

				if tt.expectedError != "" {
					assert.Error(suite.T(), err)
					assert.Contains(suite.T(), err.Error(), tt.expectedError)
				} else {
					assert.NoError(suite.T(), err)
				}
			})
		}
	})

	suite.Run("ValidateChatID", func() {
		service, _ := suite.helpers.CreateValidationService()

		tests := []struct {
			name          string
			chatID        int64
			expectedError string
		}{
			{
				name:   "positive_chat_id",
				chatID: 12345,
			},
			{
				name:   "negative_chat_id",
				chatID: -12345,
			},
			{
				name:          "zero_chat_id",
				chatID:        0,
				expectedError: "invalid chat ID",
			},
		}

		for _, tt := range tests {
			suite.Run(tt.name, func() {
				err := service.ValidateChatID(context.Background(), tt.chatID)

				if tt.expectedError != "" {
					assert.Error(suite.T(), err)
					assert.Contains(suite.T(), err.Error(), tt.expectedError)
				} else {
					assert.NoError(suite.T(), err)
				}
			})
		}
	})
}

// TestEdgeCases tests edge cases and additional functionality
func (suite *TelegramSubscriptionServiceTestSuite) TestEdgeCases() {
	suite.Run("ActivateSubscription", func() {
		service, mockRepo := suite.helpers.CreateTelegramSubscriptionService()
		subscriptionID := value_objects.NewID()

		originalSub, _ := domain.NewTelegramSubscription(suite.projectID, suite.chatID1)
		originalSub.Deactivate() // Start as inactive

		mockRepo.On("GetByID", mock.Anything, subscriptionID).Return(originalSub, nil)
		mockRepo.On("Update", mock.Anything, mock.AnythingOfType(testmocks.TelegramSubscriptionType)).Return(nil)

		err := service.ActivateTelegramSubscription(context.Background(), subscriptionID)

		assert.NoError(suite.T(), err)
		mockRepo.AssertExpectations(suite.T())
	})

	suite.Run("DeactivateSubscription", func() {
		service, mockRepo := suite.helpers.CreateTelegramSubscriptionService()
		subscriptionID := value_objects.NewID()

		originalSub, _ := domain.NewTelegramSubscription(suite.projectID, suite.chatID1)

		mockRepo.On("GetByID", mock.Anything, subscriptionID).Return(originalSub, nil)
		mockRepo.On("Update", mock.Anything, mock.AnythingOfType(testmocks.TelegramSubscriptionType)).Return(nil)

		err := service.DeactivateTelegramSubscription(context.Background(), subscriptionID)

		assert.NoError(suite.T(), err)
		mockRepo.AssertExpectations(suite.T())
	})

	suite.Run("CheckSubscriptionExists", func() {
		service, mockRepo := suite.helpers.CreateTelegramSubscriptionService()

		mockRepo.On("ExistsByProjectAndChatID", mock.Anything, suite.projectID, suite.chatID1).Return(true, nil)

		exists, err := service.CheckSubscriptionExists(context.Background(), suite.projectID, suite.chatID1)

		assert.NoError(suite.T(), err)
		assert.True(suite.T(), exists)
		mockRepo.AssertExpectations(suite.T())
	})
}

// TestTelegramSubscriptionServiceSuite runs the test suite
func TestTelegramSubscriptionServiceSuite(t *testing.T) {
	suite.Run(t, new(TelegramSubscriptionServiceTestSuite))
}
