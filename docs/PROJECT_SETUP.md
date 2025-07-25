# Project Structure & Setup Guide
## CI/CD Status Notifier Bot

### 1. Recommended Project Structure

```
CICD-Status-Notifier-Bot/
├── docs/                           # Documentation
│   ├── PRD.md                     # Product Requirements Document
│   ├── TECHNICAL_DESIGN.md        # Technical Design Document
│   ├── SPRINT_PLANNING.md         # Sprint Planning Document
│   ├── API_DOCS.md               # API Documentation
│   └── USER_GUIDE.md             # User Setup Guide
├── backend/                       # Go backend application
│   ├── cmd/
│   │   └── server/
│   │       └── main.go           # Application entry point
│   ├── internal/
│   │   ├── config/               # Configuration management
│   │   ├── models/               # Database models
│   │   ├── repositories/         # Data access layer
│   │   ├── services/             # Business logic layer
│   │   ├── handlers/             # HTTP handlers
│   │   ├── middleware/           # HTTP middleware
│   │   └── utils/                # Utility functions
│   ├── pkg/                      # Shared packages
│   │   ├── telegram/             # Telegram bot client
│   │   ├── database/             # Database connection
│   │   └── logger/               # Logging utilities
│   ├── migrations/               # Database migrations
│   ├── tests/                    # Test files
│   │   ├── unit/                 # Unit tests
│   │   ├── integration/          # Integration tests
│   │   └── fixtures/             # Test data
│   ├── scripts/                  # Build and deployment scripts
│   ├── go.mod                    # Go module file
│   ├── go.sum                    # Go module checksums
│   ├── Dockerfile               # Docker configuration
│   └── .env.example             # Environment variables template
├── frontend/                     # React dashboard
│   ├── public/                   # Static assets
│   ├── src/
│   │   ├── components/           # React components
│   │   ├── pages/                # Page components
│   │   ├── hooks/                # Custom React hooks
│   │   ├── services/             # API service layer
│   │   ├── types/                # TypeScript type definitions
│   │   ├── utils/                # Utility functions
│   │   └── App.tsx               # Main application component
│   ├── package.json              # Node.js dependencies
│   ├── tsconfig.json             # TypeScript configuration
│   ├── tailwind.config.js        # Tailwind CSS configuration
│   ├── vite.config.ts            # Vite build configuration
│   └── Dockerfile               # Docker configuration
├── docker-compose.yml            # Docker Compose configuration
├── docker-compose.dev.yml        # Development Docker Compose
├── .github/
│   └── workflows/
│       ├── backend-ci.yml        # Backend CI pipeline
│       ├── frontend-ci.yml       # Frontend CI pipeline
│       └── deploy.yml            # Deployment pipeline
├── scripts/                      # Project scripts
│   ├── setup.sh                  # Initial setup script
│   ├── test.sh                   # Run all tests
│   └── deploy.sh                 # Deployment script
├── README.md                     # Project overview
├── .gitignore                    # Git ignore rules
└── LICENSE                       # License file
```

### 2. Technology Stack Summary

#### Backend Stack
- **Language:** Go 1.21+
- **Framework:** Fiber v2 (fast HTTP framework)
- **Database:** PostgreSQL 15+ with GORM v2
- **Testing:** Testify + GoMock
- **Configuration:** Viper
- **Logging:** Logrus
- **Documentation:** Swagger/OpenAPI

#### Frontend Stack
- **Framework:** React 18 + TypeScript
- **Build Tool:** Vite
- **Styling:** Tailwind CSS + Headless UI
- **State Management:** React Query + Zustand
- **Charts:** Recharts
- **Testing:** Vitest + Testing Library

#### DevOps Stack
- **Containerization:** Docker + Docker Compose
- **CI/CD:** GitHub Actions
- **Database Migrations:** golang-migrate
- **Monitoring:** Structured logging (future: Prometheus)

### 3. Development Environment Setup

#### Prerequisites
- Go 1.21+ installed
- Node.js 18+ and npm/yarn installed
- Docker and Docker Compose installed
- PostgreSQL 15+ (or use Docker)
- Git for version control

#### Initial Setup Commands
```bash
# Clone the repository
git clone https://github.com/your-org/CICD-Status-Notifier-Bot.git
cd CICD-Status-Notifier-Bot

# Copy environment files
cp backend/.env.example backend/.env
cp frontend/.env.example frontend/.env

# Start development environment
docker-compose -f docker-compose.dev.yml up -d

# Setup backend
cd backend
go mod download
go run cmd/server/main.go migrate

# Setup frontend
cd ../frontend
npm install
npm run dev
```

### 4. Development Workflow

#### TDD Workflow
1. **Red Phase:** Write failing test first
2. **Green Phase:** Write minimal code to pass test
3. **Refactor Phase:** Improve code while keeping tests green

#### Git Workflow
```bash
# Feature development
git checkout -b feature/webhook-integration
# Make changes
git add .
git commit -m "feat: add GitHub webhook signature verification"
git push origin feature/webhook-integration
# Create Pull Request

# Testing before commit
./scripts/test.sh

# Code review process
# - At least one reviewer approval required
# - All CI checks must pass
# - Coverage threshold met
```

#### Daily Development Process
1. **Stand-up (15 min):** What did you do yesterday? What will you do today? Any blockers?
2. **Pair Programming:** Rotate pairs daily for knowledge sharing
3. **TDD Development:** Red-Green-Refactor cycle
4. **Code Review:** All code reviewed before merge
5. **Integration Testing:** Test features together regularly

### 5. Testing Strategy

#### Backend Testing
```go
// Unit test example
func TestWebhookService_ProcessGitHubWebhook(t *testing.T) {
    // Arrange
    mockBuildService := mocks.NewMockBuildService(ctrl)
    mockTelegramService := mocks.NewMockTelegramService(ctrl)
    service := NewWebhookService(mockBuildService, mockTelegramService)
    
    // Act & Assert
    err := service.ProcessGitHubWebhook(ctx, "project-id", payload)
    assert.NoError(t, err)
}

// Integration test example
func TestProjectRepository_Create(t *testing.T) {
    db := setupTestDB(t)
    repo := NewProjectRepository(db)
    
    project := &Project{Name: "test-project"}
    err := repo.Create(context.Background(), project)
    
    assert.NoError(t, err)
    assert.NotEmpty(t, project.ID)
}
```

#### Frontend Testing
```typescript
// Component test example
import { render, screen } from '@testing-library/react'
import { ProjectCard } from './ProjectCard'

test('renders project card with status', () => {
  const project = {
    id: '1',
    name: 'Test Project',
    lastBuild: { status: 'success' }
  }
  
  render(<ProjectCard project={project} />)
  
  expect(screen.getByText('Test Project')).toBeInTheDocument()
  expect(screen.getByText('success')).toBeInTheDocument()
})
```

### 6. Key Implementation Guidelines

#### Go Backend Best Practices
- Use dependency injection for testability
- Implement repository pattern for data access
- Use context for request tracing
- Handle errors explicitly (no silent failures)
- Use structured logging
- Validate all inputs
- Use database transactions where needed

#### React Frontend Best Practices
- Use TypeScript for type safety
- Implement proper error boundaries
- Use React Query for server state
- Follow component composition patterns
- Implement proper loading states
- Use proper accessibility attributes
- Optimize for performance (memoization, lazy loading)

#### Database Best Practices
- Use UUIDs for primary keys
- Add proper indexes for performance
- Use database constraints for data integrity
- Implement soft deletes where appropriate
- Use migrations for schema changes
- Regular backup strategy

### 7. Deployment Strategy

#### Development Deployment
```yaml
# docker-compose.dev.yml
version: '3.8'
services:
  app-dev:
    build:
      context: ./backend
      dockerfile: Dockerfile.dev
    volumes:
      - ./backend:/app
    environment:
      - GO_ENV=development
      - DB_HOST=postgres-dev
    depends_on:
      - postgres-dev
    
  frontend-dev:
    build:
      context: ./frontend
      dockerfile: Dockerfile.dev
    volumes:
      - ./frontend:/app
    ports:
      - "3000:3000"
```

#### Production Deployment
```yaml
# docker-compose.yml
version: '3.8'
services:
  app:
    image: cicd-notifier-bot:latest
    environment:
      - GO_ENV=production
      - DB_HOST=postgres
    restart: unless-stopped
    
  postgres:
    image: postgres:15-alpine
    environment:
      - POSTGRES_DB=cicd_notifier
    volumes:
      - postgres_data:/var/lib/postgresql/data
    restart: unless-stopped
```

### 8. Monitoring & Observability

#### Logging Strategy
```go
// Structured logging example
log.WithFields(logrus.Fields{
    "webhook_id": webhookID,
    "project_id": projectID,
    "event_type": eventType,
    "duration_ms": duration.Milliseconds(),
}).Info("Webhook processed successfully")
```

#### Health Checks
```go
// Health check endpoint
func (h *HealthHandler) Check(c *fiber.Ctx) error {
    checks := map[string]string{
        "database": h.checkDatabase(),
        "telegram": h.checkTelegramAPI(),
    }
    
    allHealthy := true
    for _, status := range checks {
        if status != "healthy" {
            allHealthy = false
            break
        }
    }
    
    if allHealthy {
        return c.JSON(fiber.Map{
            "status": "healthy",
            "checks": checks,
        })
    }
    
    return c.Status(503).JSON(fiber.Map{
        "status": "unhealthy",
        "checks": checks,
    })
}
```

### 9. Security Considerations

#### Environment Variables
```bash
# backend/.env
DB_PASSWORD=your_secure_password
TELEGRAM_BOT_TOKEN=your_bot_token
WEBHOOK_SECRET=your_webhook_secret
JWT_SECRET=your_jwt_secret
```

#### Security Checklist
- [ ] All secrets in environment variables
- [ ] Webhook signature verification implemented
- [ ] Input validation on all endpoints
- [ ] Rate limiting on public endpoints
- [ ] HTTPS in production
- [ ] Database connection encryption
- [ ] Regular dependency updates
- [ ] Security headers configured

### 10. Performance Optimization

#### Database Optimization
- Add indexes on frequently queried columns
- Use connection pooling
- Implement query result caching
- Use database-level constraints
- Monitor slow queries

#### API Optimization
- Implement response caching
- Use pagination for large result sets
- Compress responses
- Optimize JSON serialization
- Monitor response times

#### Frontend Optimization
- Code splitting and lazy loading
- Image optimization
- Bundle size monitoring
- Performance budgets
- Core Web Vitals tracking

This comprehensive setup provides a solid foundation for your 2-person team to build the CI/CD Status Notifier Bot efficiently while maintaining high code quality and following best practices.
