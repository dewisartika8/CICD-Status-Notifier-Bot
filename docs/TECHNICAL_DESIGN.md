# Technical Design Document
## CI/CD Status Notifier Bot

### 1. System Architecture Overview

```
┌─────────────────┐    ┌──────────────────┐    ┌─────────────────┐
│   GitHub        │    │   Webhook        │    │   Telegram      │
│   Actions       │───▶│   Handler        │───▶│   Bot           │
│                 │    │   (Fiber API)    │    │   Notifications │
└─────────────────┘    └──────────────────┘    └─────────────────┘
                                │
                                ▼
                       ┌──────────────────┐    ┌─────────────────┐
                       │   PostgreSQL     │    │   Web Dashboard │
                       │   Database       │◀───│   (React)       │
                       └──────────────────┘    └─────────────────┘
                                │
                                ▼
                       ┌──────────────────┐
                       │   Background     │
                       │   Jobs Queue     │
                       └──────────────────┘
```

### 2. Technology Stack

#### Backend
- **Language:** Go 1.21+
- **Framework:** Fiber v2.x (Fast, Express-inspired web framework)
- **Database:** PostgreSQL 15+
- **ORM:** GORM v2
- **Message Queue:** Built-in Go channels (Phase 1)
- **Testing:** Testify, GoMock
- **Documentation:** Swagger/OpenAPI

#### Frontend (Dashboard)
- **Framework:** React 18+ with TypeScript
- **State Management:** React Query + Zustand
- **UI Library:** Tailwind CSS + Headless UI
- **Charts:** Recharts
- **Build Tool:** Vite

#### Infrastructure
- **Containerization:** Docker + Docker Compose
- **Database Migrations:** Golang-migrate
- **Environment Management:** Viper
- **Logging:** Logrus
- **Monitoring:** Prometheus + Grafana (Future)

### 3. Database Schema

```sql
-- Projects table
CREATE TABLE projects (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(255) NOT NULL UNIQUE,
    repository_url VARCHAR(500) NOT NULL,
    webhook_secret VARCHAR(255),
    telegram_chat_id BIGINT,
    is_active BOOLEAN DEFAULT true,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Build status events
CREATE TABLE build_events (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    project_id UUID NOT NULL REFERENCES projects(id) ON DELETE CASCADE,
    event_type VARCHAR(50) NOT NULL, -- build_started, build_success, etc.
    status VARCHAR(20) NOT NULL, -- success, failed, pending
    branch VARCHAR(255) NOT NULL,
    commit_sha VARCHAR(40),
    commit_message TEXT,
    author_name VARCHAR(255),
    author_email VARCHAR(255),
    build_url VARCHAR(500),
    duration_seconds INTEGER,
    webhook_payload JSONB,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Telegram subscriptions
CREATE TABLE telegram_subscriptions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    project_id UUID NOT NULL REFERENCES projects(id) ON DELETE CASCADE,
    chat_id BIGINT NOT NULL,
    user_id BIGINT,
    username VARCHAR(255),
    event_types TEXT[], -- Array of event types to subscribe to
    is_active BOOLEAN DEFAULT true,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    UNIQUE(project_id, chat_id)
);

-- Notification logs
CREATE TABLE notification_logs (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    build_event_id UUID NOT NULL REFERENCES build_events(id) ON DELETE CASCADE,
    chat_id BIGINT NOT NULL,
    message_id INTEGER,
    status VARCHAR(20) NOT NULL, -- sent, failed, pending
    error_message TEXT,
    sent_at TIMESTAMP WITH TIME ZONE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Indexes for performance
CREATE INDEX idx_build_events_project_id ON build_events(project_id);
CREATE INDEX idx_build_events_created_at ON build_events(created_at DESC);
CREATE INDEX idx_build_events_status ON build_events(status);
CREATE INDEX idx_telegram_subscriptions_chat_id ON telegram_subscriptions(chat_id);
CREATE INDEX idx_notification_logs_build_event_id ON notification_logs(build_event_id);
```

### 4. API Design

#### 4.1 Webhook Endpoints

```go
// POST /api/v1/webhooks/github/:projectId
type GitHubWebhookPayload struct {
    Action     string `json:"action"`
    Repository struct {
        Name     string `json:"name"`
        FullName string `json:"full_name"`
        HTMLURL  string `json:"html_url"`
    } `json:"repository"`
    WorkflowRun struct {
        ID         int64  `json:"id"`
        Status     string `json:"status"`
        Conclusion string `json:"conclusion"`
        HTMLURL    string `json:"html_url"`
        HeadBranch string `json:"head_branch"`
        HeadSHA    string `json:"head_sha"`
        HeadCommit struct {
            Message string `json:"message"`
            Author  struct {
                Name  string `json:"name"`
                Email string `json:"email"`
            } `json:"author"`
        } `json:"head_commit"`
    } `json:"workflow_run"`
}
```

#### 4.2 Dashboard API Endpoints

```go
// GET /api/v1/projects
type ProjectResponse struct {
    ID            string    `json:"id"`
    Name          string    `json:"name"`
    RepositoryURL string    `json:"repository_url"`
    LastBuild     *BuildEvent `json:"last_build"`
    IsActive      bool      `json:"is_active"`
    CreatedAt     time.Time `json:"created_at"`
}

// GET /api/v1/projects/:id/builds
type BuildEventsResponse struct {
    Events     []BuildEvent `json:"events"`
    Pagination Pagination   `json:"pagination"`
}

// GET /api/v1/projects/:id/metrics
type ProjectMetrics struct {
    TotalBuilds    int     `json:"total_builds"`
    SuccessRate    float64 `json:"success_rate"`
    AvgDuration    int     `json:"avg_duration_seconds"`
    RecentBuilds   []BuildEvent `json:"recent_builds"`
    BuildsByStatus map[string]int `json:"builds_by_status"`
}
```

#### 4.3 Telegram Bot Commands

```go
type BotCommand struct {
    Command     string
    Description string
    Handler     func(ctx context.Context, message *tgbotapi.Message) error
}

var botCommands = []BotCommand{
    {"/start", "Initialize bot and show welcome message", handleStart},
    {"/help", "Show available commands", handleHelp},
    {"/status", "Show current status of all projects", handleStatus},
    {"/status <project>", "Show status of specific project", handleProjectStatus},
    {"/projects", "List all monitored projects", handleProjects},
    {"/subscribe <project>", "Subscribe to project notifications", handleSubscribe},
    {"/unsubscribe <project>", "Unsubscribe from project notifications", handleUnsubscribe},
    {"/history <project> [limit]", "Show recent build history", handleHistory},
}
```

### 5. Service Architecture

#### 5.1 Core Services

```go
// Webhook service - handles incoming GitHub webhooks
type WebhookService interface {
    ProcessGitHubWebhook(ctx context.Context, projectID string, payload GitHubWebhookPayload) error
    ValidateWebhookSignature(payload []byte, signature string, secret string) error
}

// Telegram service - manages bot interactions and notifications
type TelegramService interface {
    SendNotification(ctx context.Context, chatID int64, event BuildEvent) error
    HandleCommand(ctx context.Context, message *tgbotapi.Message) error
    RegisterWebhook() error
}

// Project service - manages project configurations
type ProjectService interface {
    CreateProject(ctx context.Context, req CreateProjectRequest) (*Project, error)
    GetProject(ctx context.Context, id string) (*Project, error)
    ListProjects(ctx context.Context) ([]Project, error)
    UpdateProject(ctx context.Context, id string, req UpdateProjectRequest) (*Project, error)
    DeleteProject(ctx context.Context, id string) error
}

// Build service - manages build events and history
type BuildService interface {
    CreateBuildEvent(ctx context.Context, event BuildEvent) error
    GetBuildEvents(ctx context.Context, projectID string, pagination Pagination) ([]BuildEvent, error)
    GetProjectMetrics(ctx context.Context, projectID string) (*ProjectMetrics, error)
    GetLatestBuildStatus(ctx context.Context, projectID string) (*BuildEvent, error)
}
```

#### 5.2 Repository Pattern

```go
// Project repository
type ProjectRepository interface {
    Create(ctx context.Context, project *Project) error
    GetByID(ctx context.Context, id string) (*Project, error)
    GetByName(ctx context.Context, name string) (*Project, error)
    List(ctx context.Context) ([]Project, error)
    Update(ctx context.Context, project *Project) error
    Delete(ctx context.Context, id string) error
}

// Build event repository
type BuildEventRepository interface {
    Create(ctx context.Context, event *BuildEvent) error
    GetByProjectID(ctx context.Context, projectID string, limit, offset int) ([]BuildEvent, error)
    GetLatestByProjectID(ctx context.Context, projectID string) (*BuildEvent, error)
    GetMetrics(ctx context.Context, projectID string) (*ProjectMetrics, error)
}
```

### 6. Message Format Design

#### 6.1 Telegram Notification Templates

```go
// Build Success Message
🎉 *Build Successful*

📦 *Project:* {{.ProjectName}}
🌿 *Branch:* {{.Branch}}
👤 *Author:* {{.AuthorName}}
⏱️ *Duration:* {{.Duration}}
📝 *Commit:* {{.CommitMessage}}

🔗 [View Build]({{.BuildURL}})

// Build Failed Message
❌ *Build Failed*

📦 *Project:* {{.ProjectName}}
🌿 *Branch:* {{.Branch}}
👤 *Author:* {{.AuthorName}}
⏱️ *Duration:* {{.Duration}}
📝 *Commit:* {{.CommitMessage}}

🔗 [View Build]({{.BuildURL}})
💡 [View Logs]({{.LogsURL}})

// Status Command Response
📊 *Project Status Overview*

{{range .Projects}}
📦 *{{.Name}}*
   Status: {{.Status}} {{.StatusIcon}}
   Branch: {{.LastBuild.Branch}}
   Last Updated: {{.LastBuild.CreatedAt.Format "15:04, 02 Jan"}}
{{end}}

🔄 Use `/status <project>` for detailed info
```

### 7. Configuration Management

```yaml
# config.yaml
server:
  port: 8080
  host: "0.0.0.0"
  read_timeout: 30s
  write_timeout: 30s

database:
  host: ${DB_HOST:localhost}
  port: ${DB_PORT:5432}
  name: ${DB_NAME:cicd_notifier}
  user: ${DB_USER:postgres}
  password: ${DB_PASSWORD}
  ssl_mode: ${DB_SSL_MODE:disable}
  max_open_conns: 25
  max_idle_conns: 10
  max_lifetime: 5m

telegram:
  bot_token: ${TELEGRAM_BOT_TOKEN}
  webhook_url: ${TELEGRAM_WEBHOOK_URL}
  webhook_secret: ${TELEGRAM_WEBHOOK_SECRET}

logging:
  level: ${LOG_LEVEL:info}
  format: ${LOG_FORMAT:json}

security:
  jwt_secret: ${JWT_SECRET}
  webhook_timeout: 30s
```

### 8. Testing Strategy

#### 8.1 Unit Testing
- Service layer unit tests with mocked dependencies
- Repository layer tests with test database
- Utility function tests
- Message formatting tests

#### 8.2 Integration Testing
- Database integration tests
- Telegram API integration tests
- GitHub webhook integration tests
- End-to-end API tests

#### 8.3 Test Structure
```go
// Example test structure
func TestWebhookService_ProcessGitHubWebhook(t *testing.T) {
    tests := []struct {
        name          string
        payload       GitHubWebhookPayload
        projectID     string
        mockSetup     func(*mocks.MockBuildService, *mocks.MockTelegramService)
        expectedError error
    }{
        {
            name: "successful build event processing",
            payload: GitHubWebhookPayload{
                Action: "completed",
                WorkflowRun: WorkflowRun{
                    Status: "completed",
                    Conclusion: "success",
                },
            },
            mockSetup: func(buildSvc *mocks.MockBuildService, telegramSvc *mocks.MockTelegramService) {
                buildSvc.EXPECT().CreateBuildEvent(gomock.Any(), gomock.Any()).Return(nil)
                telegramSvc.EXPECT().SendNotification(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
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

### 9. Security Considerations

#### 9.1 Webhook Security
- GitHub webhook signature verification using HMAC-SHA256
- Request payload size limits
- Rate limiting on webhook endpoints
- Input validation and sanitization

#### 9.2 Telegram Bot Security
- Bot token stored in environment variables
- Chat ID validation
- Command rate limiting
- User authorization for admin commands

#### 9.3 Database Security
- Connection encryption (SSL)
- Prepared statements to prevent SQL injection
- Database user with minimal required permissions
- Regular security updates

### 10. Deployment Architecture

```yaml
# docker-compose.yml
version: '3.8'
services:
  app:
    build: .
    ports:
      - "8080:8080"
    environment:
      - DB_HOST=postgres
      - TELEGRAM_BOT_TOKEN=${TELEGRAM_BOT_TOKEN}
    depends_on:
      - postgres
    restart: unless-stopped

  postgres:
    image: postgres:15-alpine
    environment:
      - POSTGRES_DB=cicd_notifier
      - POSTGRES_USER=postgres
      - POSTGRES_PASSWORD=${DB_PASSWORD}
    volumes:
      - postgres_data:/var/lib/postgresql/data
    restart: unless-stopped

  dashboard:
    build:
      context: ./frontend
    ports:
      - "3000:3000"
    environment:
      - REACT_APP_API_URL=http://localhost:8080
    depends_on:
      - app

volumes:
  postgres_data:
```

### 11. Monitoring and Observability

#### 11.1 Metrics to Track
- Webhook processing time
- Notification delivery success rate
- Database query performance
- Bot command response time
- System resource usage

#### 11.2 Logging Strategy
- Structured logging with JSON format
- Log levels: DEBUG, INFO, WARN, ERROR
- Request/response logging
- Error tracking and alerting

#### 11.3 Health Checks
- Database connectivity check
- Telegram API connectivity check
- System resource monitoring
- Application startup verification
