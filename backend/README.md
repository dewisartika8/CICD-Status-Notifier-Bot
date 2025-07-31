# CI/CD Status Notifier Bot - Backend

> **Lihat juga:** [README utama di root project](../README.md) untuk penjelasan arsitektur, fitur, dan integrasi lintas komponen.

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

#### Konfigurasi dengan Prioritas Hierarki ⚡

Aplikasi menggunakan **sistem konfigurasi hierarki** dengan prioritas sebagai berikut:

1. **Environment Variables** (Prioritas Tertinggi)
2. **Configuration File** (config/config.yaml)
3. **Default Values** (Prioritas Terendah)

> 📖 **Detail lengkap**: Lihat [Configuration Priority Guide](../docs/CONFIGURATION_PRIORITY_GUIDE.md)

#### Option 1: Menggunakan Environment Variables (Recommended untuk Production)

```bash
# Contoh untuk Windows PowerShell
$env:TELEGRAM_BOT_TOKEN="your_bot_token_here"
$env:SERVER_PORT=8080
$env:DB_HOST="localhost"
$env:DB_PASSWORD="your_secure_password"

# Contoh untuk Linux/macOS
export TELEGRAM_BOT_TOKEN="your_bot_token_here"
export SERVER_PORT=8080
export DB_HOST="localhost"
export DB_PASSWORD="your_secure_password"
```

#### Option 2: Menggunakan Configuration File (Recommended untuk Development)

```bash
# Copy file example dan edit
cp config/config-example.yaml config/config.yaml
```

Edit file `config/config.yaml` sesuai dengan environment Anda:

```yaml
server:
  port: 8081
  host: "localhost"

database:
  host: "127.0.0.1"
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
  webhook_url: "https://yourdomain.com/webhooks/telegram"

github:
  webhook_secret: "your_github_secret"

gitlab:
  webhook_secret: "your_gitlab_secret"
  
logging:
  level: "info"
  format: "json"
  output: "stdout"
```

#### Option 3: Hybrid Approach (Recommended untuk Docker)

Kombinasi environment variables untuk sensitive data dan config file untuk konfigurasi umum:

```bash
# Sensitive data via environment
export TELEGRAM_BOT_TOKEN="your_bot_token"
export DB_PASSWORD="secure_password"
export GITHUB_WEBHOOK_SECRET="secret123"

# Non-sensitive config tetap di config.yaml
```

### 4. Setup Database

#### Quick Setup (Recommended)
```bash
# Setup everything: dependencies + tools + database + migrations
make setup-fresh
```

#### Manual Setup
```bash
# Setup PostgreSQL dengan Docker (recommended)
make db-setup

# Install migrate tool
make install-migrate

# Jalankan migrations
make migrate-up
```

#### Alternative: Manual PostgreSQL Setup
```bash
# Jika Anda sudah memiliki PostgreSQL running
createdb cicd_notifier
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

### Setup Database

Ada beberapa cara untuk setup database:

#### 1. Setup Development Environment (Recommended)
```bash
# Setup environment lengkap dengan dependencies dan database
make setup-dev

# Atau setup fresh dengan reset database
make setup-fresh
```

#### 2. Manual Setup
```bash
# Install tools yang diperlukan
make install-migrate

# Setup database (buat database jika belum ada)
make db-setup

# Jalankan migrations
make migrate-up
```

### Database Migrations

#### Basic Migration Commands

```bash
# Jalankan semua migrations ke versi terbaru
make migrate-up

# Rollback migration (satu step)
make migrate-down

# Cek status migration saat ini
make migrate-status

# Reset database dan jalankan fresh migration
make db-fresh

# Force rollback semua migrations
make migrate-force-down
```

#### Advanced Migration Commands

```bash
# Install migrate tool jika belum ada
make install-migrate

# Reset database (drop dan buat ulang)
make db-reset

# Setup database tanpa migration
make db-setup
```

#### Untuk Windows (tanpa Make)

Jika menggunakan Windows PowerShell dan tidak memiliki `make`, gunakan script PowerShell:

```powershell
# Jalankan migration
.\migrate-up.ps1

# Reset database dan migration
.\db-fresh.ps1
```

> **📝 Note:** Lihat [MIGRATION_SCRIPTS.md](./scripts/migrations/MIGRATION_SCRIPTS.md) untuk dokumentasi lengkap PowerShell scripts.

#### Troubleshooting Migration

**Error: "Dirty database version"**
```bash
# Force set ke versi tertentu (misalnya versi 1)
migrate -path scripts/migrations -database "postgres://postgres:password@localhost:5432/cicd_notifier?sslmode=disable" force 1

# Lalu jalankan migration normal
make migrate-up
```

**Error: "migrate: command not found"**
```bash
# Install migrate tool
make install-migrate

# Atau manual install
go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest
```

### Database Schema

Database terdiri dari 4 tabel utama:

#### 1. `projects` Table
```sql
CREATE TABLE projects (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name VARCHAR(255) NOT NULL UNIQUE,
    repository_url VARCHAR(500) NOT NULL,
    webhook_secret VARCHAR(255),
    telegram_chat_id BIGINT,
    is_active BOOLEAN DEFAULT true,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);
```

#### 2. `build_events` Table
```sql
CREATE TABLE build_events (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    project_id UUID NOT NULL REFERENCES projects(id) ON DELETE CASCADE,
    event_type VARCHAR(50) NOT NULL,
    status VARCHAR(20) NOT NULL,
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
```

#### 3. `telegram_subscriptions` Table
```sql
CREATE TABLE telegram_subscriptions (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    project_id UUID NOT NULL REFERENCES projects(id) ON DELETE CASCADE,
    chat_id BIGINT NOT NULL,
    user_id BIGINT,
    username VARCHAR(255),
    event_types TEXT[],
    is_active BOOLEAN DEFAULT true,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);
```

#### 4. `notification_logs` Table
```sql
CREATE TABLE notification_logs (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    build_event_id UUID NOT NULL REFERENCES build_events(id) ON DELETE CASCADE,
    chat_id BIGINT NOT NULL,
    message_id INTEGER,
    status VARCHAR(20) NOT NULL,
    error_message TEXT,
    sent_at TIMESTAMP WITH TIME ZONE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);
```

### Database Configuration

Edit file `config/config.yaml` untuk mengatur koneksi database:

```yaml
database:
  host: "localhost"
  port: "5432"
  user: "postgres"
  password: "your_password"
  dbname: "cicd_notifier"
  sslmode: "disable"
  max_open_conns: 10
  max_idle_conns: 5
  max_lifetime: 300
```

### Migration Files

Migration files terletak di `scripts/migrations/`:
- `000_initial.up.sql` / `000_initial.down.sql` - Placeholder migration
- `001_initial_schema.up.sql` / `001_initial_schema.down.sql` - Schema utama
- `002_notification_tables.up.sql` / `002_notification_tables.down.sql` - Enhancement (placeholder)

## 🔧 Troubleshooting

### Common Issues

#### 1. Migration Errors

**Error: "Dirty database version"**
```bash
# Cek versi migration saat ini
make migrate-status

# Force set ke versi tertentu (contoh: versi 1)
migrate -path scripts/migrations -database "postgres://postgres:password@localhost:5432/cicd_notifier?sslmode=disable" force 1

# Jalankan migration normal
make migrate-up
```

**Error: "migrate: command not found"**
```bash
# Install migrate tool
make install-migrate

# Jika masih error, tambahkan Go bin ke PATH
export PATH="$(go env GOPATH)/bin:$PATH"
```

**Error: "relation already exists"**
```bash
# Reset database dan jalankan fresh migration
make db-fresh
```

#### 2. Database Connection Issues

**Error: "connection refused"**
- Pastikan PostgreSQL running
- Cek konfigurasi di `config/config.yaml`
- Cek environment variables

**Error: "database does not exist"**
```bash
# Buat database
make db-setup
```

#### 3. Windows Specific Issues

**Error: "make: command not found"**
```powershell
# Gunakan PowerShell scripts
.\migrate-up.ps1
.\db-fresh.ps1

# Atau install make untuk Windows
scoop install make
```

**Error: "psql: command not found"**
- Install PostgreSQL command line tools
- Atau gunakan Docker untuk database setup

#### 4. Development Issues

**Error: "air: command not found"**
```bash
make install-air
```

**Hot reload tidak bekerja**
```bash
# Hapus file .air.toml dan generate ulang
rm .air.toml
make watch
```

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

### Development Commands
- `make dev` - Run development server
- `make watch` - Run with hot reload (requires air)
- `make setup-dev` - Setup development environment
- `make setup-fresh` - Setup fresh environment with database reset

### Build Commands
- `make build` - Build application for current OS
- `make build-linux` - Build for Linux
- `make clean` - Clean build files

### Test Commands
- `make test` - Run unit tests
- `make test-coverage` - Run tests with coverage in terminal
- `make coverage` - Generate HTML coverage report

### Database Commands
- `make db-setup` - Create database if not exists
- `make db-reset` - Drop and recreate database
- `make db-fresh` - Reset database and run fresh migrations
- `make migrate-up` - Run database migrations up
- `make migrate-down` - Rollback migrations (one step)
- `make migrate-force-down` - Force rollback all migrations
- `make migrate-status` - Check current migration version

### Docker Commands
- `make docker-build` - Build Docker image
- `make docker-run` - Run Docker container
- `make docker-stop` - Stop Docker container
- `make compose-up` - Start with Docker Compose
- `make compose-down` - Stop Docker Compose

### Code Quality Commands
- `make fmt` - Format code
- `make vet` - Vet code
- `make lint` - Run linter (requires golangci-lint)
- `make sec` - Run security check (requires gosec)
- `make quality` - Run all quality checks
- `make generate` - Generate mock files

### Installation Commands
- `make deps` - Download and tidy dependencies
- `make install-air` - Install air for hot reloading
- `make install-migrate` - Install migrate tool
- `make install-lint` - Install golangci-lint
- `make install-sec` - Install gosec security scanner

### Quick Start Examples

```bash
# Setup development environment dari awal
make setup-fresh

# Development dengan hot reload
make watch

# Build dan test
make build test

# Quality check lengkap
make quality

# Reset database jika ada masalah
make db-fresh
```

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
