#!/bin/bash

# ================================
# Backup and Maintenance Script
# ================================

set -e

# Colors
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m'

# Configuration
BACKUP_DIR="./backups"
POSTGRES_DB="cicd_notifier"
POSTGRES_USER="postgres"
DATE=$(date +%Y%m%d_%H%M%S)

log_info() {
    echo -e "${BLUE}[INFO]${NC} $1"
}

log_success() {
    echo -e "${GREEN}[SUCCESS]${NC} $1"
}

log_warning() {
    echo -e "${YELLOW}[WARNING]${NC} $1"
}

log_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

create_backup() {
    log_info "Creating database backup..."
    
    # Create backup directory
    mkdir -p ${BACKUP_DIR}
    
    # Create database backup
    docker-compose exec -T postgres pg_dump -U ${POSTGRES_USER} -d ${POSTGRES_DB} > ${BACKUP_DIR}/db_backup_${DATE}.sql
    
    # Compress backup
    gzip ${BACKUP_DIR}/db_backup_${DATE}.sql
    
    log_success "Database backup created: ${BACKUP_DIR}/db_backup_${DATE}.sql.gz"
    
    # Backup configuration files
    tar -czf ${BACKUP_DIR}/config_backup_${DATE}.tar.gz \
        backend/config/ \
        frontend/.env* \
        docker-compose.yml \
        *.sh \
        Makefile 2>/dev/null || true
    
    log_success "Configuration backup created: ${BACKUP_DIR}/config_backup_${DATE}.tar.gz"
}

restore_backup() {
    local backup_file=$1
    
    if [ -z "$backup_file" ]; then
        log_error "Please specify backup file: $0 restore <backup_file>"
        exit 1
    fi
    
    if [ ! -f "$backup_file" ]; then
        log_error "Backup file not found: $backup_file"
        exit 1
    fi
    
    log_info "Restoring database from: $backup_file"
    
    # Extract if compressed
    if [[ $backup_file == *.gz ]]; then
        log_info "Extracting compressed backup..."
        gunzip -c $backup_file | docker-compose exec -T postgres psql -U ${POSTGRES_USER} -d ${POSTGRES_DB}
    else
        docker-compose exec -T postgres psql -U ${POSTGRES_USER} -d ${POSTGRES_DB} < $backup_file
    fi
    
    log_success "Database restored successfully"
}

cleanup_old_backups() {
    local days=${1:-7}  # Default 7 days
    
    log_info "Cleaning up backups older than $days days..."
    
    find ${BACKUP_DIR} -name "*.sql.gz" -mtime +$days -delete 2>/dev/null || true
    find ${BACKUP_DIR} -name "*.tar.gz" -mtime +$days -delete 2>/dev/null || true
    
    log_success "Old backups cleaned up"
}

monitor_system() {
    log_info "System monitoring report..."
    echo ""
    
    # Docker containers status
    echo "📦 Container Status:"
    docker-compose ps
    echo ""
    
    # Database status
    echo "🗄️ Database Status:"
    if docker-compose exec -T postgres pg_isready -U ${POSTGRES_USER} -d ${POSTGRES_DB} > /dev/null 2>&1; then
        echo "✅ Database: Healthy"
        
        # Database statistics
        echo ""
        echo "📊 Database Statistics:"
        docker-compose exec -T postgres psql -U ${POSTGRES_USER} -d ${POSTGRES_DB} -c "
        SELECT 
            schemaname,
            tablename,
            n_tup_ins as inserts,
            n_tup_upd as updates,
            n_tup_del as deletes,
            n_live_tup as live_rows
        FROM pg_stat_user_tables;" 2>/dev/null || true
    else
        echo "❌ Database: Unhealthy"
    fi
    
    echo ""
    
    # Disk usage
    echo "💾 Disk Usage:"
    df -h | grep -E "(Filesystem|/dev/)" || df -h
    echo ""
    
    # Memory usage
    echo "🧠 Memory Usage:"
    free -h 2>/dev/null || vm_stat | head -5
    echo ""
    
    # Process status
    echo "⚙️ Application Processes:"
    if [ -f logs/backend.pid ]; then
        if kill -0 $(cat logs/backend.pid) 2>/dev/null; then
            echo "✅ Backend: Running (PID: $(cat logs/backend.pid))"
        else
            echo "❌ Backend: Not running"
        fi
    else
        echo "❌ Backend: No PID file found"
    fi
    
    if [ -f logs/frontend.pid ]; then
        if kill -0 $(cat logs/frontend.pid) 2>/dev/null; then
            echo "✅ Frontend: Running (PID: $(cat logs/frontend.pid))"
        else
            echo "❌ Frontend: Not running"
        fi
    else
        echo "❌ Frontend: No PID file found"
    fi
    
    echo ""
    
    # Log file sizes
    echo "📝 Log Files:"
    if [ -d logs ]; then
        ls -lh logs/ 2>/dev/null || echo "No log files found"
    else
        echo "No logs directory found"
    fi
}

check_health() {
    log_info "Performing health checks..."
    
    local errors=0
    
    # Check database
    if docker-compose exec -T postgres pg_isready -U ${POSTGRES_USER} -d ${POSTGRES_DB} > /dev/null 2>&1; then
        log_success "✅ Database: Healthy"
    else
        log_error "❌ Database: Unhealthy"
        errors=$((errors + 1))
    fi
    
    # Check backend
    if curl -f http://localhost:8080/health > /dev/null 2>&1; then
        log_success "✅ Backend: Healthy"
    else
        log_warning "⚠️ Backend: Health check failed"
        errors=$((errors + 1))
    fi
    
    # Check frontend
    if curl -f http://localhost:3000 > /dev/null 2>&1; then
        log_success "✅ Frontend: Healthy"
    else
        log_warning "⚠️ Frontend: Health check failed"
        errors=$((errors + 1))
    fi
    
    # Check disk space
    local disk_usage=$(df / | awk 'NR==2 {print $5}' | sed 's/%//')
    if [ "$disk_usage" -gt 90 ]; then
        log_error "❌ Disk usage is high: ${disk_usage}%"
        errors=$((errors + 1))
    else
        log_success "✅ Disk usage is normal: ${disk_usage}%"
    fi
    
    if [ $errors -eq 0 ]; then
        log_success "All health checks passed!"
        return 0
    else
        log_error "$errors health check(s) failed!"
        return 1
    fi
}

update_apps() {
    log_info "Updating applications..."
    
    # Pull latest changes
    git pull origin development || log_warning "Git pull failed or not in git repository"
    
    # Update backend dependencies
    log_info "Updating backend dependencies..."
    cd backend
    go mod download
    go mod tidy
    cd ..
    
    # Update frontend dependencies
    log_info "Updating frontend dependencies..."
    cd frontend
    npm update
    cd ..
    
    # Rebuild applications
    log_info "Rebuilding applications..."
    make build
    
    log_success "Applications updated successfully"
    log_info "Please restart services with: make stop && make start"
}

show_logs() {
    local service=$1
    local lines=${2:-50}
    
    case $service in
        "backend")
            if [ -f logs/backend.log ]; then
                tail -n $lines logs/backend.log
            else
                log_error "Backend log file not found"
            fi
            ;;
        "frontend")
            if [ -f logs/frontend.log ]; then
                tail -n $lines logs/frontend.log
            else
                log_error "Frontend log file not found"
            fi
            ;;
        "database"|"db")
            docker-compose logs --tail=$lines postgres
            ;;
        *)
            echo "Usage: $0 logs <backend|frontend|database> [lines]"
            echo "Example: $0 logs backend 100"
            ;;
    esac
}

usage() {
    echo "Usage: $0 <command> [options]"
    echo ""
    echo "Commands:"
    echo "  backup                    - Create database and config backup"
    echo "  restore <file>           - Restore database from backup file"
    echo "  cleanup [days]           - Clean up old backups (default: 7 days)"
    echo "  monitor                  - Show system monitoring information"
    echo "  health                   - Perform health checks"
    echo "  update                   - Update applications and dependencies"
    echo "  logs <service> [lines]   - Show logs for service (backend/frontend/database)"
    echo ""
    echo "Examples:"
    echo "  $0 backup"
    echo "  $0 restore backups/db_backup_20250801_120000.sql.gz"
    echo "  $0 cleanup 14"
    echo "  $0 logs backend 100"
}

main() {
    case "${1:-}" in
        "backup")
            create_backup
            ;;
        "restore")
            restore_backup "$2"
            ;;
        "cleanup")
            cleanup_old_backups "$2"
            ;;
        "monitor")
            monitor_system
            ;;
        "health")
            check_health
            ;;
        "update")
            update_apps
            ;;
        "logs")
            show_logs "$2" "$3"
            ;;
        *)
            usage
            exit 1
            ;;
    esac
}

main "$@"
