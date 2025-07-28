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
│   │   └── main.go               # Application entry point
│   ├── internal/
│   │   ├── config/               # Configuration management
│   │   ├── core/                 # Domain layer (Clean Architecture)
│   │   │   ├── build/            # Build domain module
│   │   │   │   ├── domain/       # Business entities and domain logic
│   │   │   │   ├── port/         # Repository and service interfaces
│   │   │   │   ├── dto/          # Data Transfer Objects
│   │   │   │   └── service/      # Business logic implementation
│   │   │   ├── project/          # Project domain module
│   │   │   │   ├── domain/       # Business entities and domain logic
│   │   │   │   ├── port/         # Repository and service interfaces
│   │   │   │   ├── dto/          # Data Transfer Objects
│   │   │   │   └── service/      # Business logic implementation
│   │   │   ├── notification/     # Notification domain module (similar structure)
│   │   │   ├── webhook/          # Webhook domain module (similar structure)
│   │   │   └── shared/           # Shared domain components
│   │   │       ├── domain/
│   │   │       │   ├── value_objects/  # Common value objects (ID, Timestamp, etc)
│   │   │       │   └── events/         # Domain events
│   │   ├── adapter/              # External adapters
│   │   │   ├── handler/          # HTTP handlers
│   │   │   └── repository/       # Repository implementations
│   │   └── middleware/           # HTTP middleware
│   ├── pkg/                      # Shared packages
│   │   ├── logger/               # Logging utilities
│   │   ├── database/             # Database connection utilities
│   │   └── exception/            # Common error definitions and domain errors
│   ├── scripts/                  # Database migrations and build scripts
│   │   └── migrations/           # Database migrations
│   ├── tests/                    # Test files
│   │   ├── unit/                 # Unit tests
│   │   ├── integration/          # Integration tests
│   │   └── testutils/            # Test utilities and fixtures
│   ├── go.mod                    # Go module file
│   ├── go.sum                    # Go module checksums
│   ├── Makefile                  # Build automation
│   ├── README.md                 # Backend documentation
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
cp backend/internal/config/config.yaml.example backend/internal/config/config.yaml
cp frontend/.env.example frontend/.env

# Setup backend with Makefile
cd backend
make setup-dev    # This will install dependencies, tools, and setup database

# Start backend development server
make dev          # or make watch for hot reloading

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
// Unit test example for domain entity
func TestProject_UpdateName(t *testing.T) {
    project, _ := domain.NewProject("test", "https://github.com/test/repo", "secret", nil)
    
    err := project.UpdateName("new-name")
    
    assert.NoError(t, err)
    assert.Equal(t, "new-name", project.Name())
}

// Unit test example for service with mocked repository
func TestProjectService_CreateProject(t *testing.T) {
    // Arrange
    mockRepo := mocks.NewMockProjectRepository(ctrl)
    service := service.NewProjectService(mockRepo)
    
    req := dto.CreateProjectRequest{
        Name:          "test-project",
        RepositoryURL: "https://github.com/user/repo",
        WebhookSecret: "secret123456",
    }
    
    mockRepo.EXPECT().ExistsByName(gomock.Any(), req.Name).Return(false, nil)
    mockRepo.EXPECT().ExistsByRepositoryURL(gomock.Any(), req.RepositoryURL).Return(false, nil)
    mockRepo.EXPECT().Create(gomock.Any(), gomock.Any()).Return(nil)
    
    // Act
    project, err := service.CreateProject(context.Background(), req)
    
    // Assert
    assert.NoError(t, err)
    assert.NotNil(t, project)
    assert.Equal(t, req.Name, project.Name())
}

// Integration test example
func TestProjectRepository_Create(t *testing.T) {
    db := testutils.SetupTestDB(t)
    repo := repositories.NewProjectRepository(db)
    
    project, _ := domain.NewProject("test-project", "https://github.com/user/repo", "secret", nil)
    err := repo.Create(context.Background(), project)
    
    assert.NoError(t, err)
    assert.NotEmpty(t, project.ID())
}

// Running tests with Makefile
// make test           # Run all tests
// make test-coverage  # Run tests with coverage
// make coverage       # Generate HTML coverage report
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

### 6. Domain Module Architecture Pattern

#### Domain Module Structure
Each domain module (like `build`, `project`, `notification`) follows a consistent structure:

```
domain_module/
├── domain/          # Business entities, domain logic, and domain-specific errors
│   ├── entity.go    # Domain entity with business logic methods
│   └── errors.go    # Domain-specific error definitions
├── port/            # Interfaces (ports in hexagonal architecture)
│   ├── repository.go # Repository interface (data access)
│   └── service.go   # Service interface (business logic)
├── dto/             # Data Transfer Objects
│   └── requests.go  # DTOs for API requests and responses
└── service/         # Business logic implementation
    └── service.go   # Service implementation
```

#### Domain Layer Guidelines

**`/domain` - Business Entities & Domain Logic:**
- Contains domain entities with encapsulated business logic
- Domain-specific error definitions using DomainError pattern
- No external dependencies (pure business logic)
- Immutable state with controlled mutations through methods
- Validation rules as part of domain entity

Example entity structure:
```go
// domain/project.go
type Project struct {
    id         value_objects.ID
    name       string
    // ... other fields
}

func NewProject(name string) (*Project, error) {
    // Factory method with validation
}

func (p *Project) UpdateName(name string) error {
    // Business logic for name updates
}
```

**`/port` - Interface Definitions:**
- Repository interfaces for data persistence contracts
- Service interfaces for business logic contracts
- No implementation details, only contracts
- Used for dependency injection and testing

Example interface:
```go
// port/repository.go
type ProjectRepository interface {
    Create(ctx context.Context, project *domain.Project) error
    GetByID(ctx context.Context, id value_objects.ID) (*domain.Project, error)
    // ... other repository methods
}
```

**`/dto` - Data Transfer Objects:**
- Request/response structures for API layer
- Validation tags for input validation
- Conversion methods to/from domain entities
- No business logic, only data structure

Example DTO:
```go
// dto/project.go
type CreateProjectRequest struct {
    Name          string `json:"name" validate:"required,min=1,max=100"`
    RepositoryURL string `json:"repository_url" validate:"required,url"`
}

func ToProjectResponse(project *domain.Project) *ProjectResponse {
    // Conversion logic
}
```

**`/service` - Business Logic Implementation:**
- Implements service interfaces from `/port`
- Orchestrates business operations
- Handles domain entity lifecycle
- Enforces business rules and validations

Example service:
```go
// service/project_service.go
type ProjectService struct {
    repo port.ProjectRepository
}

func (s *ProjectService) CreateProject(ctx context.Context, req dto.CreateProjectRequest) (*domain.Project, error) {
    // Business logic implementation
}
```

#### Error Handling Pattern

**Domain-Specific Errors** (in `/domain/errors.go`):
```go
var (
    ErrProjectNotFound = exception.NewDomainError(
        "PROJECT_NOT_FOUND",
        "project not found",
    )
)
```

**Generic Errors** (in `/pkg/exception`):
```go
// Common errors that can be reused across modules
var (
    ErrValidationFailed = NewDomainError("VALIDATION_ERROR", "validation failed")
    ErrNotFound         = NewDomainError("NOT_FOUND", "resource not found")
)
```

#### Implementation Benefits

1. **Consistency**: All domain modules follow the same pattern
2. **Testability**: Clear interfaces enable easy mocking
3. **Maintainability**: Separation of concerns makes code easier to maintain
4. **Scalability**: Easy to add new domain modules
5. **Clean Architecture**: Dependencies point inward (towards domain)

### 7. Key Implementation Guidelines

#### Go Backend Best Practices

**Domain Module Organization:**
- Follow consistent domain module pattern: `/domain`, `/port`, `/dto`, `/service`
- Keep domain entities pure (no external dependencies)
- Use repository pattern with dependency injection
- Separate domain entities from database models
- Implement domain-specific errors in `/domain/errors.go`
- Use generic errors from `/pkg/exception` for common cases

**Clean Architecture Implementation:**
- Use Clean Architecture with domain-driven design
- Dependencies point inward (towards domain layer)
- Domain layer contains business logic, no infrastructure concerns
- Ports (interfaces) define contracts between layers
- Services orchestrate business operations
- DTOs handle data transformation between layers

**Error Handling:**
- Use DomainError pattern for structured error handling
- Define error codes as constants for API consistency
- Separate domain-specific errors from generic errors
- Implement error wrapping for context preservation

**Repository Implementation Pattern:**
```go
// Implementation should be in /internal/adapter/repository/
type projectRepository struct {
    db *gorm.DB
}

func NewProjectRepository(db *gorm.DB) port.ProjectRepository {
    return &projectRepository{db: db}
}

func (r *projectRepository) Create(ctx context.Context, project *domain.Project) error {
    model := toProjectModel(project) // Convert domain to database model
    if err := r.db.WithContext(ctx).Create(&model).Error; err != nil {
        return exception.NewDomainErrorWithCause("DB_CREATE_FAILED", "failed to create project", err)
    }
    return nil
}

// Database model (separate from domain entity)
type ProjectModel struct {
    ID             string `gorm:"primaryKey"`
    Name           string `gorm:"uniqueIndex;not null"`
    RepositoryURL  string `gorm:"uniqueIndex;not null"`
    WebhookSecret  string `gorm:"not null"`
    TelegramChatID *int64
    Status         string `gorm:"not null;default:'active'"`
    CreatedAt      time.Time
    UpdatedAt      time.Time
}

func toProjectModel(p *domain.Project) *ProjectModel {
    return &ProjectModel{
        ID:             p.ID().String(),
        Name:           p.Name(),
        RepositoryURL:  p.RepositoryURL(),
        // ... conversion logic
    }
}
```

**Code Organization Guidelines:**
- Use context for request tracing and cancellation
- Handle errors explicitly (no silent failures)
- Use structured logging with appropriate levels
- Validate all inputs at service layer
- Use database transactions where needed
- Keep business logic in service layer, not repositories
- Use Makefile for consistent build and test commands

**Value Objects and Entities:**
- Use value objects for IDs, timestamps, and other domain concepts
- Implement immutable entities with controlled mutation methods
- Encapsulate business rules within domain entities
- Use factory methods for entity creation with validation

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

### 8. Deployment Strategy

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

### 9. Monitoring & Observability

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

### 10. Security Considerations

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

### 11. Performance Optimization

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
