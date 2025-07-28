# Database Migration Scripts

Untuk pengguna Windows yang tidak memiliki `make` tersedia, kami menyediakan PowerShell scripts untuk menjalankan database migration.

## Prerequisites

1. PostgreSQL installed dan running
2. Go installed 
3. Database `cicd_notifier` sudah dibuat

## Scripts

### migrate-up.ps1
Menjalankan database migrations up.

```powershell
.\migrate-up.ps1
```

### db-fresh.ps1  
Reset database dan jalankan fresh migrations (memerlukan `psql` di PATH).

```powershell
.\db-fresh.ps1
```

## Manual Commands

Jika scripts tidak bekerja, Anda bisa menjalankan command manual:

```powershell
# Install migrate tool
go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest

# Add Go bin to PATH
$env:PATH += ";$(go env GOPATH)\bin"

# Check migration version
migrate -path scripts/migrations -database "postgres://postgres:jQlwjVKoQw@localhost:5432/cicd_notifier?sslmode=disable" version

# Run migrations up
migrate -path scripts/migrations -database "postgres://postgres:jQlwjVKoQw@localhost:5432/cicd_notifier?sslmode=disable" up

# Run migrations down
migrate -path scripts/migrations -database "postgres://postgres:jQlwjVKoQw@localhost:5432/cicd_notifier?sslmode=disable" down
```

## Troubleshooting

### Error: "migrate: command not found"
```powershell
# Install migrate tool
go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest

# Pastikan Go bin di PATH
$env:PATH += ";$(go env GOPATH)\bin"
```

### Error: "Dirty database version"
```powershell
# Force set version (contoh ke version 1)
migrate -path scripts/migrations -database "postgres://postgres:jQlwjVKoQw@localhost:5432/cicd_notifier?sslmode=disable" force 1
```

### Error: "psql: command not found" (untuk db-fresh.ps1)
Install PostgreSQL command line tools atau gunakan manual migration dengan migrate command saja.
