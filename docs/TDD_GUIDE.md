# Test-Driven Development (TDD) Implementation Guide
## CI/CD Status Notifier Bot

### 1. TDD Overview for This Project

Test-Driven Development adalah metodologi pengembangan yang sangat sesuai untuk project ini karena:
- **Reliability:** Bot harus dapat diandalkan untuk mengirim notifikasi
- **Integration Complexity:** Banyak integrasi eksternal (GitHub, Telegram, Database)
- **Team Collaboration:** 2 developer dapat bekerja parallel dengan confidence
- **Maintainability:** Code yang mudah di-maintain dan di-refactor

### 2. TDD Cycle Implementation

#### Red-Green-Refactor Cycle
```
🔴 RED: Write a failing test
    ↓
🟢 GREEN: Write minimal code to pass
    ↓  
🔵 REFACTOR: Improve code quality
    ↓
🔁 REPEAT
```

### 3. Testing Strategy by Layer

#### 3.1 Unit Tests (70% of total tests)

**Webhook Service Testing:**
```go
// tests/unit/services/webhook_service_test.go
package services_test

import (
    "context"
    "testing"
    "github.com/stretchr/testify/assert"
    "github.com/golang/mock/gomock"
    "your-project/internal/services"
    "your-project/tests/mocks"
)

func TestWebhookService_ProcessGitHubWebhook(t *testing.T) {
    tests := []struct {
        name           string
        payload        services.GitHubWebhookPayload
        mockSetup      func(*mocks.MockBuildService, *mocks.MockTelegramService)
        expectedError  error
        expectedCalls  int
    }{
        {
            name: "should process successful build event",
            payload: services.GitHubWebhookPayload{
                Action: "completed",
                WorkflowRun: services.WorkflowRun{
                    Status:     "completed",
                    Conclusion: "success",
                    HeadBranch: "main",
                },
            },
            mockSetup: func(buildSvc *mocks.MockBuildService, telegramSvc *mocks.MockTelegramService) {
                buildSvc.EXPECT().
                    CreateBuildEvent(gomock.Any(), gomock.Any()).
                    Return(nil).
                    Times(1)
                
                telegramSvc.EXPECT().
                    SendNotification(gomock.Any(), gomock.Any(), gomock.Any()).
                    Return(nil).
                    Times(1)
            },
            expectedError: nil,
            expectedCalls: 2,
        },
        {
            name: "should handle build service error gracefully",
            payload: services.GitHubWebhookPayload{
                Action: "completed",
                WorkflowRun: services.WorkflowRun{
                    Status: "completed",
                    Conclusion: "failure",
                },
            },
            mockSetup: func(buildSvc *mocks.MockBuildService, telegramSvc *mocks.MockTelegramService) {
                buildSvc.EXPECT().
                    CreateBuildEvent(gomock.Any(), gomock.Any()).
                    Return(errors.New("database error")).
                    Times(1)
                
                // Should not call telegram service if build service fails
                telegramSvc.EXPECT().
                    SendNotification(gomock.Any(), gomock.Any(), gomock.Any()).
                    Times(0)
            },
            expectedError: errors.New("failed to create build event: database error"),
            expectedCalls: 1,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            // Arrange
            ctrl := gomock.NewController(t)
            defer ctrl.Finish()

            mockBuildService := mocks.NewMockBuildService(ctrl)
            mockTelegramService := mocks.NewMockTelegramService(ctrl)
            
            tt.mockSetup(mockBuildService, mockTelegramService)
            
            webhookService := services.NewWebhookService(mockBuildService, mockTelegramService)

            // Act
            err := webhookService.ProcessGitHubWebhook(context.Background(), "test-project-id", tt.payload)

            // Assert
            if tt.expectedError != nil {
                assert.Error(t, err)
                assert.Contains(t, err.Error(), tt.expectedError.Error())
            } else {
                assert.NoError(t, err)
            }
        })
    }
}
```

**Telegram Service Testing:**
```go
// tests/unit/services/telegram_service_test.go
func TestTelegramService_SendNotification(t *testing.T) {
    tests := []struct {
        name           string
        buildEvent     models.BuildEvent
        mockSetup      func(*mocks.MockTelegramClient, *mocks.MockSubscriptionRepo)
        expectedError  error
    }{
        {
            name: "should send notification to subscribed users",
            buildEvent: models.BuildEvent{
                ProjectID:  "proj-123",
                EventType:  "build_success",
                Status:     "success",
                Branch:     "main",
                AuthorName: "John Doe",
            },
            mockSetup: func(client *mocks.MockTelegramClient, subRepo *mocks.MockSubscriptionRepo) {
                subscriptions := []models.TelegramSubscription{
                    {ChatID: 123456789, ProjectID: "proj-123"},
                    {ChatID: 987654321, ProjectID: "proj-123"},
                }
                
                subRepo.EXPECT().
                    GetActiveSubscriptionsByProject(gomock.Any(), "proj-123").
                    Return(subscriptions, nil).
                    Times(1)
                
                client.EXPECT().
                    SendMessage(gomock.Any(), int64(123456789), gomock.Any()).
                    Return(nil).
                    Times(1)
                
                client.EXPECT().
                    SendMessage(gomock.Any(), int64(987654321), gomock.Any()).
                    Return(nil).
                    Times(1)
            },
            expectedError: nil,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            // Test implementation
        })
    }
}
```

#### 3.2 Integration Tests (20% of total tests)

**Database Repository Testing:**
```go
// tests/integration/repositories/project_repository_test.go
package repositories_test

import (
    "context"
    "testing"
    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/require"
    "your-project/internal/repositories"
    "your-project/internal/models"
    "your-project/tests/testutils"
)

func TestProjectRepository_Create(t *testing.T) {
    // Arrange
    db := testutils.SetupTestDB(t)
    defer testutils.TeardownTestDB(t, db)
    
    repo := repositories.NewProjectRepository(db)
    
    project := &models.Project{
        Name:          "test-project",
        RepositoryURL: "https://github.com/user/repo",
        WebhookSecret: "secret123",
        IsActive:      true,
    }

    // Act
    err := repo.Create(context.Background(), project)

    // Assert
    require.NoError(t, err)
    assert.NotEmpty(t, project.ID)
    assert.NotZero(t, project.CreatedAt)
    
    // Verify project was actually saved
    savedProject, err := repo.GetByID(context.Background(), project.ID)
    require.NoError(t, err)
    assert.Equal(t, project.Name, savedProject.Name)
    assert.Equal(t, project.RepositoryURL, savedProject.RepositoryURL)
}

func TestProjectRepository_GetByName(t *testing.T) {
    // Test implementation for finding projects by name
    db := testutils.SetupTestDB(t)
    defer testutils.TeardownTestDB(t, db)
    
    repo := repositories.NewProjectRepository(db)
    
    // Create test project
    project := &models.Project{
        Name:          "unique-project-name",
        RepositoryURL: "https://github.com/user/repo",
    }
    
    err := repo.Create(context.Background(), project)
    require.NoError(t, err)
    
    // Test finding by name
    foundProject, err := repo.GetByName(context.Background(), "unique-project-name")
    require.NoError(t, err)
    assert.Equal(t, project.ID, foundProject.ID)
    
    // Test not found case
    _, err = repo.GetByName(context.Background(), "non-existent-project")
    assert.Error(t, err)
    assert.True(t, errors.Is(err, repositories.ErrProjectNotFound))
}
```

**API Integration Testing:**
```go
// tests/integration/handlers/webhook_handler_test.go
func TestWebhookHandler_HandleGitHubWebhook(t *testing.T) {
    // Setup test server
    app := testutils.SetupTestApp(t)
    
    tests := []struct {
        name           string
        projectID      string
        payload        string
        signature      string
        expectedStatus int
        expectedBody   string
    }{
        {
            name:      "should process valid webhook",
            projectID: "valid-project-id",
            payload:   `{"action":"completed","workflow_run":{"status":"completed","conclusion":"success"}}`,
            signature: "sha256=valid_signature",
            expectedStatus: 200,
            expectedBody:   `{"message":"webhook processed successfully"}`,
        },
        {
            name:           "should reject invalid signature",
            projectID:      "valid-project-id", 
            payload:        `{"action":"completed"}`,
            signature:      "sha256=invalid_signature",
            expectedStatus: 401,
            expectedBody:   `{"error":"invalid webhook signature"}`,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            // Create request
            req := httptest.NewRequest("POST", "/api/v1/webhooks/github/"+tt.projectID, 
                strings.NewReader(tt.payload))
            req.Header.Set("Content-Type", "application/json")
            req.Header.Set("X-Hub-Signature-256", tt.signature)

            // Perform request
            resp, err := app.Test(req)
            require.NoError(t, err)
            
            // Assert response
            assert.Equal(t, tt.expectedStatus, resp.StatusCode)
            
            body, err := io.ReadAll(resp.Body)
            require.NoError(t, err)
            assert.JSONEq(t, tt.expectedBody, string(body))
        })
    }
}
```

#### 3.3 End-to-End Tests (10% of total tests)

**Complete Workflow Testing:**
```go
// tests/e2e/webhook_to_notification_test.go
func TestCompleteWebhookToNotificationFlow(t *testing.T) {
    // This test verifies the complete flow from webhook reception to Telegram notification
    
    // Setup
    app := testutils.SetupTestApp(t)
    mockTelegramClient := testutils.SetupMockTelegramClient(t)
    
    // Create test project
    project := testutils.CreateTestProject(t, "test-repo", "webhook-secret")
    
    // Create subscription
    subscription := testutils.CreateTestSubscription(t, project.ID, 123456789)
    
    // Prepare webhook payload
    payload := testutils.CreateGitHubWebhookPayload("completed", "success", "main")
    signature := testutils.GenerateWebhookSignature(payload, "webhook-secret")
    
    // Mock Telegram client to capture sent message
    var sentMessage string
    mockTelegramClient.EXPECT().
        SendMessage(gomock.Any(), int64(123456789), gomock.Any()).
        DoAndReturn(func(ctx context.Context, chatID int64, message string) error {
            sentMessage = message
            return nil
        }).
        Times(1)
    
    // Send webhook
    req := httptest.NewRequest("POST", "/api/v1/webhooks/github/"+project.ID,
        strings.NewReader(payload))
    req.Header.Set("X-Hub-Signature-256", signature)
    
    resp, err := app.Test(req)
    require.NoError(t, err)
    assert.Equal(t, 200, resp.StatusCode)
    
    // Wait for async processing
    time.Sleep(100 * time.Millisecond)
    
    // Verify notification was sent
    assert.Contains(t, sentMessage, "Build Successful")
    assert.Contains(t, sentMessage, "test-repo")
    assert.Contains(t, sentMessage, "main")
    
    // Verify build event was saved
    buildEvents := testutils.GetBuildEventsByProject(t, project.ID)
    assert.Len(t, buildEvents, 1)
    assert.Equal(t, "build_success", buildEvents[0].EventType)
}
```

### 4. TDD Implementation Timeline

#### Sprint 1: Foundation Testing
**Week 1:**
- Setup testing framework and mocks
- Write failing tests for basic models
- Implement models to pass tests
- Write failing tests for project repository
- Implement repository to pass tests

**Week 2:**
- Write failing tests for webhook signature verification
- Implement signature verification
- Write failing tests for webhook payload parsing
- Implement payload parsing

#### Sprint 2: Service Layer Testing
**Week 3:**
- Write failing tests for Telegram service
- Implement basic Telegram functionality
- Write failing tests for notification formatting
- Implement notification templates

**Week 4:**
- Write failing tests for bot command parsing
- Implement command router
- Write failing tests for subscription management
- Implement subscription logic

#### Sprint 3: Integration Testing
**Week 5:**
- Write failing integration tests for API endpoints
- Implement API handlers
- Write failing tests for metrics calculation
- Implement metrics services

**Week 6:**
- Write failing tests for advanced bot commands
- Implement advanced features
- Write failing tests for webhook processing pipeline
- Implement complete processing pipeline

#### Sprint 4: E2E Testing
**Week 7:**
- Write failing E2E tests for complete user journeys
- Fix any integration issues found
- Write failing tests for dashboard API
- Implement dashboard backend

**Week 8:**
- Write failing tests for frontend components
- Implement frontend components
- Write failing tests for deployment scenarios
- Fix deployment and configuration issues

### 5. Test Utilities and Helpers

```go
// tests/testutils/database.go
package testutils

import (
    "testing"
    "gorm.io/gorm"
    "gorm.io/driver/sqlite"
)

func SetupTestDB(t *testing.T) *gorm.DB {
    db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
    if err != nil {
        t.Fatalf("Failed to connect to test database: %v", err)
    }
    
    // Run migrations
    err = db.AutoMigrate(&models.Project{}, &models.BuildEvent{}, &models.TelegramSubscription{})
    if err != nil {
        t.Fatalf("Failed to migrate test database: %v", err)
    }
    
    return db
}

func TeardownTestDB(t *testing.T, db *gorm.DB) {
    sqlDB, err := db.DB()
    if err != nil {
        t.Errorf("Failed to get underlying sql.DB: %v", err)
        return
    }
    
    err = sqlDB.Close()
    if err != nil {
        t.Errorf("Failed to close test database: %v", err)
    }
}

// tests/testutils/fixtures.go
func CreateTestProject(t *testing.T, name, secret string) *models.Project {
    project := &models.Project{
        Name:          name,
        RepositoryURL: "https://github.com/test/" + name,
        WebhookSecret: secret,
        IsActive:      true,
    }
    
    db := GetTestDB(t)
    err := db.Create(project).Error
    require.NoError(t, err)
    
    return project
}

func CreateGitHubWebhookPayload(status, conclusion, branch string) string {
    return fmt.Sprintf(`{
        "action": "completed",
        "workflow_run": {
            "status": "%s",
            "conclusion": "%s",
            "head_branch": "%s",
            "head_sha": "abc123",
            "html_url": "https://github.com/test/repo/actions/runs/123"
        },
        "repository": {
            "name": "test-repo",
            "full_name": "test/repo"
        }
    }`, status, conclusion, branch)
}
```

### 6. Mock Generation

```go
// Generate mocks using gomock
//go:generate mockgen -source=internal/services/interfaces.go -destination=tests/mocks/mock_services.go

// internal/services/interfaces.go
package services

type BuildService interface {
    CreateBuildEvent(ctx context.Context, event models.BuildEvent) error
    GetBuildEvents(ctx context.Context, projectID string, pagination Pagination) ([]models.BuildEvent, error)
    GetProjectMetrics(ctx context.Context, projectID string) (*ProjectMetrics, error)
}

type TelegramService interface {
    SendNotification(ctx context.Context, chatID int64, event models.BuildEvent) error
    HandleCommand(ctx context.Context, message *tgbotapi.Message) error
}

type ProjectService interface {
    CreateProject(ctx context.Context, req CreateProjectRequest) (*models.Project, error)
    GetProject(ctx context.Context, id string) (*models.Project, error)
    ListProjects(ctx context.Context) ([]models.Project, error)
}
```

### 7. Continuous Integration Testing

```yaml
# .github/workflows/backend-ci.yml
name: Backend CI

on:
  push:
    branches: [ main, develop ]
  pull_request:
    branches: [ main ]

jobs:
  test:
    runs-on: ubuntu-latest
    
    services:
      postgres:
        image: postgres:15
        env:
          POSTGRES_PASSWORD: password
          POSTGRES_DB: cicd_notifier_test
        options: >-
          --health-cmd pg_isready
          --health-interval 10s
          --health-timeout 5s
          --health-retries 5

    steps:
    - uses: actions/checkout@v3
    
    - name: Set up Go
      uses: actions/setup-go@v3
      with:
        go-version: 1.21
        
    - name: Cache Go modules
      uses: actions/cache@v3
      with:
        path: ~/go/pkg/mod
        key: ${{ runner.os }}-go-${{ hashFiles('**/go.sum') }}
        
    - name: Install dependencies
      run: go mod download
      working-directory: ./backend
      
    - name: Generate mocks
      run: go generate ./...
      working-directory: ./backend
      
    - name: Run unit tests
      run: go test -v -race -coverprofile=coverage.out ./...
      working-directory: ./backend
      env:
        GO_ENV: test
        
    - name: Run integration tests
      run: go test -v -tags=integration ./tests/integration/...
      working-directory: ./backend
      env:
        GO_ENV: test
        DB_HOST: localhost
        DB_PORT: 5432
        DB_NAME: cicd_notifier_test
        DB_USER: postgres
        DB_PASSWORD: password
        
    - name: Check test coverage
      run: |
        go tool cover -func=coverage.out
        go tool cover -html=coverage.out -o coverage.html
      working-directory: ./backend
      
    - name: Upload coverage reports
      uses: codecov/codecov-action@v3
      with:
        file: ./backend/coverage.out
        
    - name: Check code quality
      run: |
        go vet ./...
        go fmt ./...
        golangci-lint run
      working-directory: ./backend
```

### 8. Testing Best Practices untuk Project Ini

#### DO's:
- ✅ Write tests before implementation (Red-Green-Refactor)
- ✅ Test one thing at a time
- ✅ Use descriptive test names
- ✅ Mock external dependencies (Telegram API, GitHub API)
- ✅ Test both happy path and error scenarios
- ✅ Use table-driven tests for multiple scenarios
- ✅ Test edge cases (empty payloads, invalid signatures)
- ✅ Verify mock expectations
- ✅ Test database transactions and rollbacks
- ✅ Test concurrent scenarios (webhook processing)

#### DON'Ts:
- ❌ Don't test external APIs directly in unit tests
- ❌ Don't write tests that depend on test order
- ❌ Don't ignore test failures
- ❌ Don't mock everything (test real database operations in integration tests)
- ❌ Don't write overly complex test setups
- ❌ Don't forget to test error handling
- ❌ Don't skip testing for time-sensitive operations
- ❌ Don't hardcode test data that changes

### 9. Measuring TDD Success

#### Code Coverage Targets:
- **Unit Tests:** >90% coverage
- **Integration Tests:** >80% coverage
- **Overall:** >85% coverage

#### Quality Metrics:
- **Bug Rate:** <2 bugs per 100 lines of code
- **Test Maintenance:** <10% of development time spent fixing tests
- **Refactoring Confidence:** Can refactor without breaking functionality
- **Development Speed:** Stable velocity after initial learning curve

#### Team Metrics:
- **Test-First Development:** >80% of features developed with tests first
- **Regression Rate:** <5% of bugs are regressions
- **Code Review Time:** Reduced time due to test coverage
- **Deployment Confidence:** High confidence in releases

TDD implementation ini akan memastikan bahwa CI/CD Status Notifier Bot yang Anda bangun memiliki kualitas tinggi, mudah di-maintain, dan dapat diandalkan dalam environment production.
