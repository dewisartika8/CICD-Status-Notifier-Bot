#!/bin/bash

# ================================
# System Monitoring and Alerting Script
# ================================

set -e

# Colors
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
BLUE='\033[0;34m'
NC='\033[0m'

# Configuration
ALERT_EMAIL="${ALERT_EMAIL:-admin@example.com}"
ALERT_WEBHOOK="${ALERT_WEBHOOK:-}"
CHECK_INTERVAL="${CHECK_INTERVAL:-300}"  # 5 minutes
LOG_FILE="logs/monitoring.log"

# Thresholds
CPU_THRESHOLD=80
MEMORY_THRESHOLD=80
DISK_THRESHOLD=90
RESPONSE_TIME_THRESHOLD=5000  # 5 seconds

log_with_timestamp() {
    echo "$(date '+%Y-%m-%d %H:%M:%S') - $1" | tee -a $LOG_FILE
}

send_alert() {
    local level=$1
    local message=$2
    local timestamp=$(date '+%Y-%m-%d %H:%M:%S')
    
    # Log the alert
    log_with_timestamp "[$level] ALERT: $message"
    
    # Send email alert if configured
    if command -v mail &> /dev/null && [ -n "$ALERT_EMAIL" ]; then
        echo "Alert: $message at $timestamp" | mail -s "CI/CD Notifier Alert - $level" $ALERT_EMAIL
    fi
    
    # Send webhook alert if configured
    if command -v curl &> /dev/null && [ -n "$ALERT_WEBHOOK" ]; then
        curl -X POST "$ALERT_WEBHOOK" \
            -H "Content-Type: application/json" \
            -d "{\"level\":\"$level\",\"message\":\"$message\",\"timestamp\":\"$timestamp\",\"service\":\"cicd-notifier\"}" \
            > /dev/null 2>&1 || true
    fi
}

check_system_resources() {
    local alerts=0
    
    # Check CPU usage (Linux)
    if command -v top &> /dev/null; then
        local cpu_usage=$(top -bn1 | grep "Cpu(s)" | sed "s/.*, *\([0-9.]*\)%* id.*/\1/" | awk '{print 100 - $1}')
        if (( $(echo "$cpu_usage > $CPU_THRESHOLD" | bc -l) )); then
            send_alert "WARNING" "High CPU usage: ${cpu_usage}%"
            alerts=$((alerts + 1))
        fi
    fi
    
    # Check memory usage
    if command -v free &> /dev/null; then
        local mem_usage=$(free | grep Mem | awk '{printf "%.1f", $3/$2 * 100.0}')
        if (( $(echo "$mem_usage > $MEMORY_THRESHOLD" | bc -l) )); then
            send_alert "WARNING" "High memory usage: ${mem_usage}%"
            alerts=$((alerts + 1))
        fi
    fi
    
    # Check disk usage
    local disk_usage=$(df / | awk 'NR==2 {print $5}' | sed 's/%//')
    if [ "$disk_usage" -gt "$DISK_THRESHOLD" ]; then
        send_alert "CRITICAL" "High disk usage: ${disk_usage}%"
        alerts=$((alerts + 1))
    fi
    
    return $alerts
}

check_services() {
    local alerts=0
    
    # Check database
    if ! docker-compose exec -T postgres pg_isready -U postgres -d cicd_notifier > /dev/null 2>&1; then
        send_alert "CRITICAL" "Database is not responding"
        alerts=$((alerts + 1))
    fi
    
    # Check backend
    local backend_response_time=$(curl -o /dev/null -s -w '%{time_total}' http://localhost:8080/health || echo "timeout")
    if [ "$backend_response_time" = "timeout" ]; then
        send_alert "CRITICAL" "Backend service is not responding"
        alerts=$((alerts + 1))
    else
        local response_ms=$(echo "$backend_response_time * 1000" | bc -l | cut -d. -f1)
        if [ "$response_ms" -gt "$RESPONSE_TIME_THRESHOLD" ]; then
            send_alert "WARNING" "Backend slow response time: ${response_ms}ms"
            alerts=$((alerts + 1))
        fi
    fi
    
    # Check frontend
    if ! curl -f http://localhost:3000 > /dev/null 2>&1; then
        send_alert "WARNING" "Frontend service is not responding"
        alerts=$((alerts + 1))
    fi
    
    # Check processes
    if [ -f logs/backend.pid ]; then
        if ! kill -0 $(cat logs/backend.pid) 2>/dev/null; then
            send_alert "CRITICAL" "Backend process is not running"
            alerts=$((alerts + 1))
        fi
    fi
    
    if [ -f logs/frontend.pid ]; then
        if ! kill -0 $(cat logs/frontend.pid) 2>/dev/null; then
            send_alert "WARNING" "Frontend process is not running"
            alerts=$((alerts + 1))
        fi
    fi
    
    return $alerts
}

check_log_errors() {
    local alerts=0
    
    # Check backend logs for errors
    if [ -f logs/backend.log ]; then
        local error_count=$(tail -100 logs/backend.log | grep -i error | wc -l)
        if [ "$error_count" -gt 5 ]; then
            send_alert "WARNING" "High error count in backend logs: $error_count errors in last 100 lines"
            alerts=$((alerts + 1))
        fi
    fi
    
    # Check for panic/fatal errors
    if [ -f logs/backend.log ]; then
        local fatal_count=$(tail -100 logs/backend.log | grep -i "panic\|fatal" | wc -l)
        if [ "$fatal_count" -gt 0 ]; then
            send_alert "CRITICAL" "Fatal errors detected in backend logs: $fatal_count"
            alerts=$((alerts + 1))
        fi
    fi
    
    return $alerts
}

generate_status_report() {
    local timestamp=$(date '+%Y-%m-%d %H:%M:%S')
    
    echo "======================================"
    echo "CI/CD Notifier Status Report"
    echo "Generated: $timestamp"
    echo "======================================"
    echo ""
    
    # System resources
    echo "📊 System Resources:"
    if command -v free &> /dev/null; then
        echo "Memory: $(free -h | awk 'NR==2{printf "%.1f%% used\n", $3*100/$2}')"
    fi
    echo "Disk: $(df -h / | awk 'NR==2{print $5 " used"}')"
    echo ""
    
    # Services status
    echo "🔧 Services Status:"
    
    # Database
    if docker-compose exec -T postgres pg_isready -U postgres -d cicd_notifier > /dev/null 2>&1; then
        echo "✅ Database: Healthy"
    else
        echo "❌ Database: Unhealthy"
    fi
    
    # Backend
    local backend_status=$(curl -o /dev/null -s -w '%{http_code}' http://localhost:8080/health 2>/dev/null || echo "000")
    if [ "$backend_status" = "200" ]; then
        echo "✅ Backend: Healthy"
    else
        echo "❌ Backend: Unhealthy (HTTP $backend_status)"
    fi
    
    # Frontend
    local frontend_status=$(curl -o /dev/null -s -w '%{http_code}' http://localhost:3000 2>/dev/null || echo "000")
    if [ "$frontend_status" = "200" ]; then
        echo "✅ Frontend: Healthy"
    else
        echo "❌ Frontend: Unhealthy (HTTP $frontend_status)"
    fi
    
    echo ""
    
    # Recent errors
    echo "🚨 Recent Errors (last hour):"
    if [ -f logs/backend.log ]; then
        local recent_errors=$(grep "$(date '+%Y-%m-%d %H')" logs/backend.log | grep -i error | wc -l)
        echo "Backend errors: $recent_errors"
    fi
    
    echo ""
    
    # Performance metrics
    echo "⚡ Performance:"
    local backend_response=$(curl -o /dev/null -s -w '%{time_total}' http://localhost:8080/health 2>/dev/null || echo "timeout")
    if [ "$backend_response" != "timeout" ]; then
        echo "Backend response time: ${backend_response}s"
    else
        echo "Backend response time: timeout"
    fi
    
    echo ""
}

monitor_continuous() {
    log_with_timestamp "Starting continuous monitoring (interval: ${CHECK_INTERVAL}s)"
    
    while true; do
        local total_alerts=0
        
        # Run all checks
        check_system_resources
        total_alerts=$((total_alerts + $?))
        
        check_services  
        total_alerts=$((total_alerts + $?))
        
        check_log_errors
        total_alerts=$((total_alerts + $?))
        
        if [ $total_alerts -eq 0 ]; then
            log_with_timestamp "All checks passed - system healthy"
        else
            log_with_timestamp "$total_alerts alerts detected"
        fi
        
        sleep $CHECK_INTERVAL
    done
}

install_monitoring_service() {
    echo "Installing monitoring as systemd service..."
    
    cat > cicd-notifier-monitor.service << EOF
[Unit]
Description=CI/CD Status Notifier Monitoring
After=network.target
Requires=docker.service

[Service]
Type=simple
User=$USER
WorkingDirectory=$(pwd)
ExecStart=$(pwd)/monitor.sh continuous
Restart=always
RestartSec=30
Environment=PATH=/usr/local/bin:/usr/bin:/bin

[Install]
WantedBy=multi-user.target
EOF
    
    echo "Service file created: cicd-notifier-monitor.service"
    echo ""
    echo "To install:"
    echo "  sudo cp cicd-notifier-monitor.service /etc/systemd/system/"
    echo "  sudo systemctl daemon-reload"
    echo "  sudo systemctl enable cicd-notifier-monitor"
    echo "  sudo systemctl start cicd-notifier-monitor"
}

usage() {
    echo "Usage: $0 <command>"
    echo ""
    echo "Commands:"
    echo "  check        - Run all health checks once"
    echo "  continuous   - Run continuous monitoring"
    echo "  report       - Generate status report"
    echo "  install      - Install monitoring service"
    echo ""
    echo "Environment Variables:"
    echo "  ALERT_EMAIL=admin@example.com     # Email for alerts"
    echo "  ALERT_WEBHOOK=http://webhook.url  # Webhook for alerts"
    echo "  CHECK_INTERVAL=300                # Check interval in seconds"
}

# Create logs directory
mkdir -p logs

case "${1:-}" in
    "check")
        echo "Running health checks..."
        total_alerts=0
        
        check_system_resources
        total_alerts=$((total_alerts + $?))
        
        check_services
        total_alerts=$((total_alerts + $?))
        
        check_log_errors
        total_alerts=$((total_alerts + $?))
        
        if [ $total_alerts -eq 0 ]; then
            echo -e "${GREEN}✅ All checks passed - system healthy${NC}"
            exit 0
        else
            echo -e "${RED}❌ $total_alerts alerts detected${NC}"
            exit 1
        fi
        ;;
    "continuous")
        monitor_continuous
        ;;
    "report")
        generate_status_report
        ;;
    "install")
        install_monitoring_service
        ;;
    *)
        usage
        exit 1
        ;;
esac
