# PowerShell script untuk reset database dan migration
# DB_URL for migration
$DB_HOST = "localhost"
$DB_PORT = "5432"
$DB_USER = "postgres"
$DB_PASSWORD = "jQlwjVKoQw"
$DB_NAME = "cicd_notifier"
$DB_URL = "postgres://$DB_USER`:$DB_PASSWORD@$DB_HOST`:$DB_PORT/$DB_NAME`?sslmode=disable"

Write-Host "Resetting database..." -ForegroundColor Yellow

# Drop and recreate database
try {
    Write-Host "Dropping database if exists..." -ForegroundColor Gray
    psql -U $DB_USER -h $DB_HOST -p $DB_PORT -c "DROP DATABASE IF EXISTS $DB_NAME;" 2>$null
    
    Write-Host "Creating database..." -ForegroundColor Gray
    psql -U $DB_USER -h $DB_HOST -p $DB_PORT -c "CREATE DATABASE $DB_NAME;"
    
    Write-Host "Database $DB_NAME has been reset" -ForegroundColor Green
    
    # Install migrate tool if not exists
    $MIGRATE_BIN = "$(go env GOPATH)\bin\migrate.exe"
    if (-not (Test-Path $MIGRATE_BIN)) {
        Write-Host "Installing migrate tool..." -ForegroundColor Yellow
        go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest
    }
    
    # Run migrations
    Write-Host "Running database migrations..." -ForegroundColor Yellow
    & $MIGRATE_BIN -path scripts/migrations -database $DB_URL up
    
    Write-Host "Database fresh migration completed!" -ForegroundColor Green
    
} catch {
    Write-Host "Error: $($_.Exception.Message)" -ForegroundColor Red
    exit 1
}
