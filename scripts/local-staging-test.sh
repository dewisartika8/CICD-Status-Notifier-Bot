#!/bin/bash

# Local Staging Deployment Test Script (Without SSH)
# This script tests the staging deployment locally using Docker directly

set -e

# Configuration
TARGET_HOST="localhost"
STAGING_BACKEND_PORT="8082"
STAGING_FRONTEND_PORT="3002"
STAGING_DB_PORT="5434"

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
PURPLE='\033[0;35m'
NC='\033[0m' # No Color

# Function to print colored output
print_status() {
    echo -e "${BLUE}[INFO]${NC} $1"
}

print_success() {
    echo -e "${GREEN}[SUCCESS]${NC} $1"
}

print_warning() {
    echo -e "${YELLOW}[WARNING]${NC} $1"
}

print_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

print_header() {
    echo -e "${PURPLE}$1${NC}"
}

# Function to check if command exists
command_exists() {
    command -v "$1" >/dev/null 2>&1
}

# Check prerequisites
check_prerequisites() {
    print_header "🔍 Checking Prerequisites"
    
    local missing_tools=()
    
    if ! command_exists "curl"; then
        missing_tools+=("curl")
    fi
    
    if ! command_exists "docker"; then
        missing_tools+=("docker")
    fi
    
    if ! command_exists "docker-compose"; then
        missing_tools+=("docker-compose")
    fi
    
    if ! command_exists "git"; then
        missing_tools+=("git")
    fi
    
    if [ ${#missing_tools[@]} -ne 0 ]; then
        print_error "Missing required tools: ${missing_tools[*]}"
        exit 1
    fi
    
    # Check if we're on staging branch
    current_branch=$(git branch --show-current 2>/dev/null || echo "unknown")
    if [ "$current_branch" != "staging" ]; then
        print_warning "Current branch is '$current_branch', not 'staging'"
        print_status "Switching to staging branch..."
        git checkout staging || {
            print_error "Failed to switch to staging branch"
            exit 1
        }
    fi
    
    # Check Docker daemon
    if ! docker info >/dev/null 2>&1; then
        print_error "Docker daemon is not running"
        print_error "Please start Docker and try again"
        exit 1
    fi
    
    print_success "All prerequisites are met"
}

# Check local environment
check_local_environment() {
    print_header "🖥️ Checking Local Environment"
    
    print_status "OS: $(uname -a)"
    print_status "Docker: $(docker --version)"
    print_status "Docker Compose: $(docker-compose --version)"
    
    # Check disk space
    print_status "Disk usage:"
    df -h /
    
    # Check memory
    print_status "Memory usage:"
    if command_exists "free"; then
        free -h
    else
        # macOS alternative
        echo "Memory: $(vm_stat | grep "Pages free" | awk '{print $3}' | sed 's/\.//')K free"
    fi
    
    # Create staging directory
    print_status "Creating staging directory..."
    mkdir -p ./staging-deployment
    cd ./staging-deployment
    
    print_success "Local environment check passed"
}

# Build Docker images locally
build_images_locally() {
    print_header "🐳 Building Docker Images Locally"
    
    # Build backend image
    print_status "Building backend image..."
    docker build -t cicd-backend:staging-local ../backend
    
    # Build frontend image  
    print_status "Building frontend image..."
    docker build -t cicd-frontend:staging-local ../frontend
    
    print_success "Images built successfully"
}

# Deploy to local staging
deploy_local_staging() {
    print_header "🚀 Deploying to Local Staging Environment"
    
    # Stop existing containers
    print_status "Stopping existing containers..."
    docker-compose -f docker-compose.staging.yml down 2>/dev/null || true
    
    # Create environment file
    print_status "Creating environment configuration..."
    cat > .env.staging << 'EOF'
POSTGRES_DB=cicd_notifier_staging
POSTGRES_USER=postgres
POSTGRES_PASSWORD=staging_password_123
BACKEND_IMAGE=cicd-backend:staging-local
FRONTEND_IMAGE=cicd-frontend:staging-local
JWT_SECRET=staging_jwt_secret_key_123
TELEGRAM_BOT_TOKEN=your_telegram_bot_token_here
EOF

    # Create docker-compose file
    print_status "Creating docker-compose configuration..."
    cat > docker-compose.staging.yml << 'EOF'
version: '3.8'

services:
  postgres-staging:
    image: postgres:15-alpine
    container_name: cicd_postgres_staging_local
    environment:
      POSTGRES_DB: ${POSTGRES_DB}
      POSTGRES_USER: ${POSTGRES_USER}
      POSTGRES_PASSWORD: ${POSTGRES_PASSWORD}
    ports:
      - "5434:5432"
    volumes:
      - postgres_staging_data:/var/lib/postgresql/data
    healthcheck:
      test: ["CMD-SHELL", "pg_isready -U postgres"]
      interval: 10s
      timeout: 5s
      retries: 10
    restart: unless-stopped

  backend-staging:
    image: ${BACKEND_IMAGE}
    container_name: cicd_backend_staging_local
    depends_on:
      postgres-staging:
        condition: service_healthy
    ports:
      - "8082:8080"
    environment:
      DATABASE_URL: postgres://${POSTGRES_USER}:${POSTGRES_PASSWORD}@postgres-staging:5432/${POSTGRES_DB}?sslmode=disable
      PORT: 8080
      ENV: staging
      JWT_SECRET: ${JWT_SECRET}
      TELEGRAM_BOT_TOKEN: ${TELEGRAM_BOT_TOKEN}
    restart: unless-stopped
    healthcheck:
      test: ["CMD", "wget", "--no-verbose", "--tries=1", "--spider", "http://localhost:8080/health"]
      interval: 30s
      timeout: 10s
      retries: 5
      start_period: 60s

  frontend-staging:
    image: ${FRONTEND_IMAGE}
    container_name: cicd_frontend_staging_local
    depends_on:
      backend-staging:
        condition: service_healthy
    ports:
      - "3002:80"
    environment:
      REACT_APP_API_URL: http://localhost:8082
    restart: unless-stopped

volumes:
  postgres_staging_data:
EOF

    # Start services
    print_status "Starting staging services..."
    docker-compose -f docker-compose.staging.yml --env-file .env.staging up -d
    
    print_status "Waiting for services to start..."
    sleep 60
    
    # Check status
    print_status "Container status:"
    docker-compose -f docker-compose.staging.yml ps
    
    print_success "Local staging deployment completed"
}

# Run health checks
run_health_checks() {
    print_header "🔍 Running Health Checks"
    
    print_status "Waiting for services to stabilize..."
    sleep 30
    
    # Check database
    print_status "Checking database connectivity..."
    if docker exec cicd_postgres_staging_local pg_isready -U postgres; then
        print_success "Database is healthy"
    else
        print_error "Database health check failed"
        docker logs cicd_postgres_staging_local --tail=20
        return 1
    fi
    
    # Check backend
    print_status "Checking backend health..."
    max_attempts=10
    attempt=1
    while [ $attempt -le $max_attempts ]; do
        if curl -f http://localhost:$STAGING_BACKEND_PORT/health >/dev/null 2>&1; then
            print_success "Backend is healthy"
            break
        else
            print_warning "Attempt $attempt/$max_attempts: Backend not ready yet..."
            sleep 15
            attempt=$((attempt + 1))
        fi
    done
    
    if [ $attempt -gt $max_attempts ]; then
        print_error "Backend health check failed after $max_attempts attempts"
        docker logs cicd_backend_staging_local --tail=30
        return 1
    fi
    
    # Check frontend
    print_status "Checking frontend availability..."
    if curl -f http://localhost:$STAGING_FRONTEND_PORT >/dev/null 2>&1; then
        print_success "Frontend is accessible"
    else
        print_error "Frontend accessibility check failed"
        docker logs cicd_frontend_staging_local --tail=20
        return 1
    fi
    
    print_success "All health checks passed"
}

# Test external accessibility
test_external_access() {
    print_header "🌍 Testing External Accessibility"
    
    # Test backend from external
    print_status "Testing backend external access..."
    if curl -f "http://$TARGET_HOST:$STAGING_BACKEND_PORT/health" >/dev/null 2>&1; then
        print_success "Backend is externally accessible"
        
        # Get response time
        response_time=$(curl -o /dev/null -s -w '%{time_total}\n' "http://$TARGET_HOST:$STAGING_BACKEND_PORT/health")
        print_status "Backend response time: ${response_time}s"
    else
        print_error "Backend is not externally accessible"
        return 1
    fi
    
    # Test frontend from external
    print_status "Testing frontend external access..."
    if curl -f "http://$TARGET_HOST:$STAGING_FRONTEND_PORT" >/dev/null 2>&1; then
        print_success "Frontend is externally accessible"
        
        # Get response time
        response_time=$(curl -o /dev/null -s -w '%{time_total}\n' "http://$TARGET_HOST:$STAGING_FRONTEND_PORT")
        print_status "Frontend response time: ${response_time}s"
    else
        print_error "Frontend is not externally accessible"
        return 1
    fi
    
    print_success "External accessibility tests passed"
}

# Generate test report
generate_report() {
    print_header "📊 Generating Test Report"
    
    # Get container status
    container_status=$(docker-compose -f docker-compose.staging.yml ps)
    
    # Get system resources
    system_info="Disk Usage:\n$(df -h /)\n\nMemory Usage:\n$(if command_exists 'free'; then free -h; else echo 'Memory info not available on macOS'; fi)"
    
    cat > ../staging-test-report.md << EOF
# Local Staging Deployment Test Report

**Date**: $(date)
**Branch**: $(git branch --show-current 2>/dev/null || echo "unknown")
**Commit**: $(git rev-parse --short HEAD 2>/dev/null || echo "unknown")
**Target**: Local Docker Environment

## Test Results ✅

- ✅ Prerequisites Check
- ✅ Local Environment Check
- ✅ Docker Images Build
- ✅ Local Staging Deployment
- ✅ Health Checks
- ✅ External Accessibility

## Environment Access

- **Frontend**: http://$TARGET_HOST:$STAGING_FRONTEND_PORT
- **Backend API**: http://$TARGET_HOST:$STAGING_BACKEND_PORT
- **Health Check**: http://$TARGET_HOST:$STAGING_BACKEND_PORT/health
- **Database**: $TARGET_HOST:$STAGING_DB_PORT

## Container Status

\`\`\`
$container_status
\`\`\`

## System Resources

\`\`\`
$system_info
\`\`\`

## Next Steps

1. ✅ Local testing completed successfully
2. 🚀 GitHub Actions will handle remote deployment
3. 🔍 Monitor the GitHub Actions workflow
4. 📝 Review logs and metrics
5. 🎯 Ready for production deployment

## Commands for Further Testing

\`\`\`bash
# Test API endpoints
curl http://$TARGET_HOST:$STAGING_BACKEND_PORT/health
curl http://$TARGET_HOST:$STAGING_BACKEND_PORT/api/v1/projects

# View logs
docker logs cicd_backend_staging_local
docker logs cicd_frontend_staging_local
docker logs cicd_postgres_staging_local

# Check container status
cd staging-deployment && docker-compose -f docker-compose.staging.yml ps
\`\`\`

---
Generated by local staging deployment test script
EOF
    
    print_success "Test report generated: staging-test-report.md"
}

# Cleanup function
cleanup() {
    print_status "Cleaning up..."
    cd ..
}

# Main execution
main() {
    clear
    print_header "🧪 Local Staging Deployment Test"
    print_header "================================="
    echo ""
    
    # Set trap for cleanup on exit
    trap cleanup EXIT
    
    check_prerequisites
    echo ""
    
    check_local_environment
    echo ""
    
    build_images_locally
    echo ""
    
    deploy_local_staging
    echo ""
    
    run_health_checks
    echo ""
    
    test_external_access
    echo ""
    
    generate_report
    echo ""
    
    print_success "🎉 Local staging deployment test completed successfully!"
    echo ""
    print_status "Environment is ready at:"
    print_status "  - Frontend: http://$TARGET_HOST:$STAGING_FRONTEND_PORT"
    print_status "  - Backend:  http://$TARGET_HOST:$STAGING_BACKEND_PORT"
    print_status "  - Health:   http://$TARGET_HOST:$STAGING_BACKEND_PORT/health"
    echo ""
    print_status "Next steps:"
    print_status "  1. Review the staging-test-report.md file"
    print_status "  2. Test the application manually"
    print_status "  3. GitHub Actions is already running from the push"
    print_status "  4. Monitor GitHub Actions workflow at:"
    print_status "     https://github.com/dewisartika8/CICD-Status-Notifier-Bot/actions"
    echo ""
    print_status "To stop the staging environment:"
    print_status "  cd staging-deployment && docker-compose -f docker-compose.staging.yml down"
    echo ""
}

# Check if script is being sourced or executed
if [[ "${BASH_SOURCE[0]}" == "${0}" ]]; then
    main "$@"
fi
