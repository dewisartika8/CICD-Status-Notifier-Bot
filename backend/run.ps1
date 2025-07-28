# PowerShell script untuk Windows development
# run.ps1 - Helper script untuk development di Windows

param(
    [Parameter(Position=0)]
    [string]$Command = "help"
)

$ErrorActionPreference = "Stop"

# Configuration
$BINARY_NAME = "cicd-notifier-bot.exe"
$DB_HOST = "localhost"
$DB_PORT = "5432"
$DB_USER = "postgres"
$DB_PASSWORD = "password"
$DB_NAME = "cicd_notifier"

function Show-Help {
    Write-Host "CI/CD Status Notifier Bot - Development Helper" -ForegroundColor Green
    Write-Host ""
    Write-Host "Usage: .\run.ps1 <command>" -ForegroundColor Yellow
    Write-Host ""
    Write-Host "Commands:" -ForegroundColor Cyan
    Write-Host "  dev           Run application in development mode"
    Write-Host "  build         Build the application"
    Write-Host "  test          Run tests"
    Write-Host "  coverage      Run tests with coverage"
    Write-Host "  clean         Clean build files"
    Write-Host "  deps          Download dependencies"
    Write-Host "  db-setup      Setup PostgreSQL with Docker"
    Write-Host "  db-stop       Stop PostgreSQL container"
    Write-Host "  migrate-up    Run database migrations up"
    Write-Host "  migrate-down  Run database migrations down"
    Write-Host "  help          Show this help"
    Write-Host ""
}

function Invoke-Dev {
    Write-Host "Starting application in development mode..." -ForegroundColor Green
    go run cmd/main.go
}

function Invoke-Build {
    Write-Host "Building application..." -ForegroundColor Green
    go build -o $BINARY_NAME -v ./cmd/main.go
    Write-Host "Build completed: $BINARY_NAME" -ForegroundColor Green
}

function Invoke-Test {
    Write-Host "Running tests..." -ForegroundColor Green
    go test -v -race -timeout 30s ./...
}

function Invoke-Coverage {
    Write-Host "Running tests with coverage..." -ForegroundColor Green
    go test -v -race -coverprofile=coverage.out ./...
    go tool cover -html=coverage.out -o coverage.html
    Write-Host "Coverage report generated: coverage.html" -ForegroundColor Green
}

function Invoke-Clean {
    Write-Host "Cleaning build files..." -ForegroundColor Green
    if (Test-Path $BINARY_NAME) {
        Remove-Item $BINARY_NAME
    }
    if (Test-Path "coverage.out") {
        Remove-Item "coverage.out"
    }
    if (Test-Path "coverage.html") {
        Remove-Item "coverage.html"
    }
    go clean
    Write-Host "Clean completed" -ForegroundColor Green
}

function Invoke-Deps {
    Write-Host "Downloading dependencies..." -ForegroundColor Green
    go mod download
    go mod tidy
    Write-Host "Dependencies updated" -ForegroundColor Green
}

function Invoke-DbSetup {
    Write-Host "Setting up PostgreSQL with Docker..." -ForegroundColor Green
    
    # Stop existing container if running
    docker stop postgres-dev 2>$null
    docker rm postgres-dev 2>$null
    
    # Start new container
    docker run --name postgres-dev `
        -e POSTGRES_USER=$DB_USER `
        -e POSTGRES_PASSWORD=$DB_PASSWORD `
        -e POSTGRES_DB=$DB_NAME `
        -p "${DB_PORT}:5432" `
        -d postgres:15-alpine
    
    Write-Host "Waiting for database to be ready..." -ForegroundColor Yellow
    Start-Sleep -Seconds 5
    
    Write-Host "Database setup completed" -ForegroundColor Green
    Write-Host "Connection: postgresql://${DB_USER}:${DB_PASSWORD}@${DB_HOST}:${DB_PORT}/${DB_NAME}" -ForegroundColor Cyan
}

function Invoke-DbStop {
    Write-Host "Stopping PostgreSQL container..." -ForegroundColor Green
    docker stop postgres-dev 2>$null
    docker rm postgres-dev 2>$null
    Write-Host "Database stopped" -ForegroundColor Green
}

function Invoke-MigrateUp {
    Write-Host "Running database migrations up..." -ForegroundColor Green
    $DB_URL = "postgres://${DB_USER}:${DB_PASSWORD}@${DB_HOST}:${DB_PORT}/${DB_NAME}?sslmode=disable"
    
    # Check if migrate tool exists
    $migratePath = Get-Command migrate -ErrorAction SilentlyContinue
    if (-not $migratePath) {
        Write-Host "Error: migrate tool not found" -ForegroundColor Red
        Write-Host "Install it from: https://github.com/golang-migrate/migrate" -ForegroundColor Yellow
        return
    }
    
    migrate -path scripts/migrations -database $DB_URL up
    Write-Host "Migrations completed" -ForegroundColor Green
}

function Invoke-MigrateDown {
    Write-Host "Running database migrations down..." -ForegroundColor Green
    $DB_URL = "postgres://${DB_USER}:${DB_PASSWORD}@${DB_HOST}:${DB_PORT}/${DB_NAME}?sslmode=disable"
    
    # Check if migrate tool exists
    $migratePath = Get-Command migrate -ErrorAction SilentlyContinue
    if (-not $migratePath) {
        Write-Host "Error: migrate tool not found" -ForegroundColor Red
        Write-Host "Install it from: https://github.com/golang-migrate/migrate" -ForegroundColor Yellow
        return
    }
    
    migrate -path scripts/migrations -database $DB_URL down 1
    Write-Host "Migration rollback completed" -ForegroundColor Green
}

# Main command dispatcher
switch ($Command.ToLower()) {
    "dev" { Invoke-Dev }
    "build" { Invoke-Build }
    "test" { Invoke-Test }
    "coverage" { Invoke-Coverage }
    "clean" { Invoke-Clean }
    "deps" { Invoke-Deps }
    "db-setup" { Invoke-DbSetup }
    "db-stop" { Invoke-DbStop }
    "migrate-up" { Invoke-MigrateUp }
    "migrate-down" { Invoke-MigrateDown }
    "help" { Show-Help }
    default { 
        Write-Host "Unknown command: $Command" -ForegroundColor Red
        Write-Host ""
        Show-Help
    }
}
