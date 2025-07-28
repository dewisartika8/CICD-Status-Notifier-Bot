# CI/CD Status Notifier Bot - Backend

Backend service untuk CI/CD Status Notifier Bot yang dibangun dengan Go, menggunakan Fiber web framework dan PostgreSQL sebagai database.

## 📋 Daftar Isi

- [Arsitektur](#arsitektur)
- [Teknologi](#teknologi)
- [Prasyarat](#prasyarat)
- [Instalasi](#instalasi)
- [Menjalankan Aplikasi](#menjalankan-aplikasi)
- [Testing](#testing)
- [Struktur Project](#struktur-project)
- [API Endpoints](#api-endpoints)
- [Database](#database)
- [Docker](#docker)
- [Development](#development)

## 🏗️ Arsitektur

Backend menggunakan Clean Architecture dengan struktur berikut:

```
backend/
├── cmd/                    # Application entry points
├── internal/               # Private application code
│   ├── config/            # Configuration management
│   ├── domain/            # Business entities and interfaces
│   │   ├── entities/      # Domain entities
│   │   └── ports/         # Repository interfaces
│   ├── repositories/      # Repository implementations
│   ├── services/          # Business logic services
│   └── adapters/          # External adapters (database models)
├── pkg/                   # Shared packages
│   ├── logger/           # Logging utilities
│   └── database/         # Database connection
├── scripts/              # Database migrations and scripts
└── tests/                # Test files
```

## 🚀 Teknologi

- **Go 1.21+** - Programming language
- **Fiber v2** - Web framework
- **GORM v2** - ORM
- **PostgreSQL 15+** - Database
- **Docker** - Containerization
- **Testify** - Testing framework
- **Air** - Hot reloading untuk development

## 📋 Prasyarat

Pastikan Anda memiliki tools berikut terinstall:

- [Go 1.21+](https://golang.org/doc/install)
- [PostgreSQL 15+](https://www.postgresql.org/download/) atau Docker
- [Make](https://www.gnu.org/software/make/) (optional, tapi direkomendasikan)
- [Docker](https://docs.docker.com/get-docker/) (optional)

## 🛠️ Instalasi

### 1. Clone Repository

```bash
git clone https://github.com/dewisartika8/cicd-status-notifier-bot.git
cd cicd-status-notifier-bot/backend
```

### 2. Install Dependencies

```bash
# Menggunakan Makefile
make deps

# Atau manual
go mod download
go mod tidy
```

### 3. Setup Environment

```bash
# Copy dan edit file konfigurasi
cp internal/config/config.yaml.example internal/config/config.yaml
```

Edit file `config.yaml` sesuai dengan environment Anda:

```yaml
server:
  port: "8080"
  host: "localhost"

database:
  host: "localhost"
  port: "5432"
  user: "postgres"
  password: "password"
  dbname: "cicd_notifier"
  sslmode: "disable"
  max_open_conns: 10
  max_idle_conns: 5
  max_lifetime: 300

telegram:
  bot_token: "your_bot_token"
  
logging:
  level: "info"
  format: "json"
```

### 4. Setup Database

```bash
# Setup PostgreSQL dengan Docker (recommended)
make db-setup

# Atau setup manual PostgreSQL, lalu jalankan migrations
make migrate-up
```

## 🏃‍♂️ Menjalankan Aplikasi

### Development Mode

```bash
# Menjalankan dengan hot reload
make watch

# Atau menjalankan sekali
make dev

# Atau manual
go run cmd/main.go
```

### Production Mode

```bash
# Build aplikasi
make build

# Jalankan binary
./cicd-notifier-bot
```

## 🧪 Testing

### Menjalankan Unit Tests

```bash
# Menjalankan semua tests
make test

# Test dengan coverage
make test-coverage

# Generate coverage report HTML
make coverage
```

### Structure Test Files

```
tests/
├── unit/                     # Unit tests
│   ├── entities/            # Entity tests
│   ├── repositories/        # Repository tests
│   └── services/           # Service tests
├── integration/            # Integration tests
└── testutils/             # Test utilities
```

### Contoh Unit Test

```go
// tests/unit/entities/project_test.go
func TestProject_Validate(t *testing.T) {
    tests := []struct {
        name    string
        project *entities.Project
        wantErr bool
    }{
        {
            name: "valid project",
            project: &entities.Project{
                Name:          "test-project",
                RepositoryURL: "https://github.com/user/repo",
            },
            wantErr: false,
        },
        {
            name: "invalid project - empty name",
            project: &entities.Project{
                Name:          "",
                RepositoryURL: "https://github.com/user/repo",
            },
            wantErr: true,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            err := tt.project.Validate()
            if (err != nil) != tt.wantErr {
                t.Errorf("Project.Validate() error = %v, wantErr %v", err, tt.wantErr)
            }
        })
    }
}
```

### Menjalankan Integration Tests

```bash
# Pastikan database test siap
make db-setup

# Jalankan integration tests
go test -v ./tests/integration/...
```

## 📁 Struktur Project Detail

### Domain Layer (`internal/domain/`)

**Entities**: Business objects murni tanpa dependencies
```go
// entities/project.go
type Project struct {
    ID            uuid.UUID
    Name          string
    RepositoryURL string
    WebhookSecret string
    TelegramChatID int64
    IsActive      bool
    CreatedAt     time.Time
    UpdatedAt     time.Time
}
```

**Ports**: Interfaces untuk repository dan external services
```go
// ports/project_repository.go
type ProjectRepository interface {
    Create(ctx context.Context, project *entities.Project) error
    GetByID(ctx context.Context, id uuid.UUID) (*entities.Project, error)
    // ... other methods
}
```

### Repository Layer (`internal/repositories/`)

Implementasi dari domain ports untuk data access:

```go
// repositories/project_repository.go
type projectRepository struct {
    db *gorm.DB
}

func (r *projectRepository) Create(ctx context.Context, project *entities.Project) error {
    // Implementation
}
```

### Service Layer (`internal/services/`)

Business logic dan use cases:

```go
// services/project_service.go
type ProjectService struct {
    repo ports.ProjectRepository
}

func (s *ProjectService) CreateProject(ctx context.Context, project *entities.Project) error {
    // Validation and business logic
    if err := project.Validate(); err != nil {
        return err
    }
    return s.repo.Create(ctx, project)
}
```

## 🌐 API Endpoints

### Health Check
```
GET /health
```
Response:
```json
{
    "status": "healthy",
    "database": "connected",
    "timestamp": "2024-01-20T10:30:00Z"
}
```

### Projects
```
GET    /api/v1/projects          # List all projects
POST   /api/v1/projects          # Create new project
GET    /api/v1/projects/:id      # Get project by ID
PUT    /api/v1/projects/:id      # Update project
DELETE /api/v1/projects/:id      # Delete project
```

### Build Events
```
GET    /api/v1/build-events                    # List build events
POST   /api/v1/build-events                    # Create build event
GET    /api/v1/build-events/:id                # Get build event
GET    /api/v1/projects/:id/build-events       # Get project build events
GET    /api/v1/projects/:id/metrics            # Get project metrics
```

### Webhooks
```
POST   /webhooks/github          # GitHub webhook endpoint
POST   /webhooks/gitlab          # GitLab webhook endpoint
```

## 🗄️ Database

### Migrations

```bash
# Run migrations up
make migrate-up

# Run migrations down
make migrate-down

# Reset database
make db-reset
```

### Schema

Database terdiri dari 4 tabel utama:
- `projects` - Project information
- `build_events` - CI/CD build events
- `telegram_subscriptions` - Telegram chat subscriptions
- `notification_logs` - Notification delivery logs

## 🐳 Docker

### Build Docker Image

```bash
make docker-build
```

### Run dengan Docker

```bash
# Single container
make docker-run

# Dengan Docker Compose
make compose-up
```

### Docker Compose untuk Development

```yaml
# docker-compose.dev.yml
version: '3.8'
services:
  app:
    build: 
      context: .
      dockerfile: Dockerfile.dev
    ports:
      - "8080:8080"
    volumes:
      - .:/app
    environment:
      - GO_ENV=development
    depends_on:
      - postgres

  postgres:
    image: postgres:15-alpine
    environment:
      - POSTGRES_DB=cicd_notifier
      - POSTGRES_USER=postgres
      - POSTGRES_PASSWORD=password
    ports:
      - "5432:5432"
```

## 🔧 Development

### Hot Reloading

```bash
# Install air (sekali saja)
make install-air

# Jalankan dengan hot reload
make watch
```

### Code Quality

```bash
# Format code
make fmt

# Vet code
make vet

# Lint code (requires golangci-lint)
make lint

# Security check (requires gosec)
make sec

# Run all quality checks
make quality
```

### Database Development

```bash
# Setup development database
make db-setup

# Stop database
make db-stop

# Reset database
make db-reset
```

## 📝 Makefile Commands

Gunakan `make help` untuk melihat semua commands yang tersedia:

```bash
make help
```

Common commands:
- `make dev` - Run development server
- `make test` - Run tests
- `make build` - Build application
- `make docker-build` - Build Docker image
- `make db-setup` - Setup development database
- `make migrate-up` - Run database migrations

## 🤝 Contributing

1. Fork the repository
2. Create feature branch (`git checkout -b feature/amazing-feature`)
3. Write tests for your changes
4. Ensure all tests pass (`make test`)
5. Run quality checks (`make quality`)
6. Commit your changes (`git commit -m 'Add amazing feature'`)
7. Push to branch (`git push origin feature/amazing-feature`)
8. Create Pull Request

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](../LICENSE) file for details.
