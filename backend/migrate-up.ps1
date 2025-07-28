# PowerShell script untuk migration saja
$DB_HOST = "localhost"
$DB_PORT = "5432"
$DB_USER = "postgres"
$DB_PASSWORD = "jQlwjVKoQw"
$DB_NAME = "cicd_notifier"
$DB_URL = "postgres://$DB_USER`:$DB_PASSWORD@$DB_HOST`:$DB_PORT/$DB_NAME`?sslmode=disable"

# Add Go bin to PATH
$env:PATH += ";$(go env GOPATH)\bin"

Write-Host "Running database migrations..." -ForegroundColor Yellow

try {
    # Run migrations
    migrate -path scripts/migrations -database $DB_URL up
    Write-Host "Database migration completed!" -ForegroundColor Green
} catch {
    Write-Host "Error: $($_.Exception.Message)" -ForegroundColor Red
    exit 1
}
