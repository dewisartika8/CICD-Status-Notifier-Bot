#!/bin/bash

# Migration helper script
# Usage: ./migrate.sh up|down|version|force

set -e

# Database connection parameters
DB_HOST=${DB_HOST:-localhost}
DB_PORT=${DB_PORT:-5432}
DB_USER=${DB_USER:-postgres}
DB_PASSWORD=${DB_PASSWORD:-password}
DB_NAME=${DB_NAME:-cicd_notifier}
DB_URL="postgres://${DB_USER}:${DB_PASSWORD}@${DB_HOST}:${DB_PORT}/${DB_NAME}?sslmode=disable"

MIGRATION_PATH="./migrations"

# Check if migrate tool is installed
if ! command -v migrate &> /dev/null; then
    echo "Error: migrate tool is not installed"
    echo "Install it with: make install-migrate"
    exit 1
fi

# Check if migration directory exists
if [[ ! -d "$MIGRATION_PATH" ]]; then
    echo "Error: Migration directory $MIGRATION_PATH not found"
    exit 1
fi

case "$1" in
    up)
        echo "Running migrations up..."
        migrate -path $MIGRATION_PATH -database "$DB_URL" up
        ;;
    down)
        echo "Running migrations down..."
        if [[ -n "$2" ]]; then
            migrate -path $MIGRATION_PATH -database "$DB_URL" down $2
        else
            migrate -path $MIGRATION_PATH -database "$DB_URL" down 1
        fi
        ;;
    version)
        echo "Current migration version:"
        migrate -path $MIGRATION_PATH -database "$DB_URL" version
        ;;
    force)
        if [[ -z "$2" ]]; then
            echo "Error: Please provide version number to force"
            echo "Usage: ./migrate.sh force <version>"
            exit 1
        fi
        echo "Forcing migration to version $2..."
        migrate -path $MIGRATION_PATH -database "$DB_URL" force $2
        ;;
    *)
        echo "Usage: $0 {up|down|version|force}"
        echo ""
        echo "Commands:"
        echo "  up              Run all pending migrations"
        echo "  down [N]        Rollback N migrations (default: 1)"
        echo "  version         Show current migration version"
        echo "  force <version> Force database to specific version"
        echo ""
        echo "Environment variables:"
        echo "  DB_HOST     Database host (default: localhost)"
        echo "  DB_PORT     Database port (default: 5432)"
        echo "  DB_USER     Database user (default: postgres)"
        echo "  DB_PASSWORD Database password (default: password)"
        echo "  DB_NAME     Database name (default: cicd_notifier)"
        exit 1
        ;;
esac

echo "Migration operation completed successfully!"
