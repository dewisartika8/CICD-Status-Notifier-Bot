#!/bin/bash

# Manual Staging Deployment Test Script
# This script helps test the staging deployment manually before relying on GitHub Actions

set -e

# Configuration
TARGET_HOST="172.16.19.11"
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
    
    if ! command_exists "ssh"; then
        missing_tools+=("ssh")
    fi
    
    if ! command_exists "curl"; then
        missing_tools+=("curl")
    fi
    
    if ! command_exists "docker"; then
        missing_tools+=("docker")
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
    
    print_success "All prerequisites are met"
}

# Test SSH connectivity
test_ssh_connectivity() {
    print_header "🔐 Testing SSH Connectivity"
    
    if ssh -o ConnectTimeout=10 -o StrictHostKeyChecking=no "$TARGET_HOST" "echo 'SSH connection successful'" 2>/dev/null; then
        print_success "SSH connection to $TARGET_HOST successful"
    else
        print_error "Failed to connect to $TARGET_HOST via SSH"
        print_error "Please ensure:"
        print_error "1. SSH key is properly configured"
        print_error "2. Server is accessible"
        print_error "3. User has proper permissions"
        exit 1
    fi
}

# Check server prerequisites
check_server_prerequisites() {
    print_header "🖥️ Checking Server Prerequisites"
    
    ssh "$TARGET_HOST" << 'EOF'
        echo "Checking server environment..."
        
        # Check OS
        echo "OS: $(uname -a)"
        
        # Check Docker
        if command -v docker >/dev/null 2>&1; then
            echo "✅ Docker: $(docker --version)"
        else
            echo "❌ Docker not installed"
            exit 1
        fi
        
        # Check Docker Compose
        if command -v docker-compose >/dev/null 2>&1; then
            echo "✅ Docker Compose: $(docker-compose --version)"
        else
            echo "❌ Docker Compose not installed"
            exit 1
        fi
        
        # Check disk space
        echo "Disk usage:"
        df -h /
        available_space=$(df / | tail -1 | awk '{print $4}' | sed 's/G//')
        if [ "${available_space%.*}" -lt 5 ]; then
            echo "⚠️ Warning: Less than 5GB available space"
        fi
        
        # Check memory
        echo "Memory usage:"
        free -h
        
        # Check if staging directory exists
        if [ -d "/opt/cicd-notifier-staging" ]; then
            echo "✅ Staging directory exists"
        else
            echo "Creating staging directory..."
            sudo mkdir -p /opt/cicd-notifier-staging
            sudo chown $USER:$USER /opt/cicd-notifier-staging
        fi
EOF
    
    if [ $? -eq 0 ]; then
        print_success "Server prerequisites check passed"
    else
        print_error "Server prerequisites check failed"
        exit 1
    fi
}

# Build and push Docker images
build_and_push_images() {
    print_header "🐳 Building and Pushing Docker Images"
    
    # Check if GitHub CLI is available for authentication
    if command_exists "gh"; then
        print_status "Using GitHub CLI for authentication..."
        echo $(gh auth token) | docker login ghcr.io -u $(gh api user --jq .login) --password-stdin
    else
        print_warning "GitHub CLI not available. Please ensure you're logged into ghcr.io manually:"
        print_warning "echo \$GITHUB_TOKEN | docker login ghcr.io -u \$GITHUB_USERNAME --password-stdin"
    fi
    
    # Build backend image
    print_status "Building backend image..."
    docker build -t ghcr.io/${{ github.repository }}-backend:staging-manual ./backend
    
    # Build frontend image
    print_status "Building frontend image..."
    docker build -t ghcr.io/${{ github.repository }}-frontend:staging-manual ./frontend
    
    # Push images
    print_status "Pushing images to registry..."
    docker push ghcr.io/${{ github.repository }}-backend:staging-manual
    docker push ghcr.io/${{ github.repository }}-frontend:staging-manual
    
    print_success "Images built and pushed successfully"
}

# Deploy to staging
deploy_to_staging() {
    print_header "🚀 Deploying to Staging Environment"
    
    # Copy deployment files to server
    print_status "Preparing deployment files..."
    
    # Create temporary deployment script
    cat > /tmp/staging-deploy.sh << 'EOF'
#!/bin/bash
set -e

echo "🚀 Starting staging deployment on server..."

# Navigate to staging directory
cd /opt/cicd-notifier-staging

# Stop existing containers
if [ -f "docker-compose.staging.yml" ]; then
    echo "🛑 Stopping existing containers..."
    sudo docker-compose -f docker-compose.staging.yml down || true
fi

# Pull latest images
echo "📥 Pulling latest images..."
sudo docker pull ghcr.io/dewisartika8/cicd-status-notifier-bot-backend:staging-manual
sudo docker pull ghcr.io/dewisartika8/cicd-status-notifier-bot-frontend:staging-manual

# Create environment file
echo "⚙️ Creating environment configuration..."
cat > .env.staging << 'ENVEOF'
POSTGRES_DB=cicd_notifier_staging
POSTGRES_USER=postgres
POSTGRES_PASSWORD=staging_password_123
BACKEND_IMAGE=ghcr.io/dewisartika8/cicd-status-notifier-bot-backend:staging-manual
FRONTEND_IMAGE=ghcr.io/dewisartika8/cicd-status-notifier-bot-frontend:staging-manual
JWT_SECRET=staging_jwt_secret_key_123
TELEGRAM_BOT_TOKEN=your_telegram_bot_token_here
ENVEOF

# Create docker-compose file
echo "📝 Creating docker-compose configuration..."
cat > docker-compose.staging.yml << 'COMPOSEEOF'
version: '3.8'

services:
  postgres-staging:
    image: postgres:15-alpine
    container_name: cicd_postgres_staging
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
    container_name: cicd_backend_staging
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
    container_name: cicd_frontend_staging
    depends_on:
      backend-staging:
        condition: service_healthy
    ports:
      - "3002:80"
    environment:
      REACT_APP_API_URL: http://172.16.19.11:8082
    restart: unless-stopped

volumes:
  postgres_staging_data:
COMPOSEEOF

# Start services
echo "🔄 Starting services..."
sudo docker-compose -f docker-compose.staging.yml --env-file .env.staging up -d

echo "⏳ Waiting for services to start..."
sleep 60

# Check status
echo "📊 Service status:"
sudo docker-compose -f docker-compose.staging.yml ps

echo "✅ Staging deployment completed!"
EOF
    
    # Copy and execute deployment script
    print_status "Copying deployment script to server..."
    scp /tmp/staging-deploy.sh "$TARGET_HOST:/tmp/"
    
    print_status "Executing deployment on server..."
    ssh "$TARGET_HOST" "chmod +x /tmp/staging-deploy.sh && /tmp/staging-deploy.sh"
    
    # Cleanup
    rm /tmp/staging-deploy.sh
    
    print_success "Deployment completed successfully"
}

# Run health checks
run_health_checks() {
    print_header "🔍 Running Health Checks"
    
    print_status "Waiting for services to stabilize..."
    sleep 30
    
    # Check database
    print_status "Checking database connectivity..."
    if ssh "$TARGET_HOST" "docker exec cicd_postgres_staging pg_isready -U postgres"; then
        print_success "Database is healthy"
    else
        print_error "Database health check failed"
        ssh "$TARGET_HOST" "docker logs cicd_postgres_staging --tail=20"
        return 1
    fi
    
    # Check backend
    print_status "Checking backend health..."
    max_attempts=10
    attempt=1
    while [ $attempt -le $max_attempts ]; do
        if ssh "$TARGET_HOST" "curl -f http://localhost:$STAGING_BACKEND_PORT/health" >/dev/null 2>&1; then
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
        ssh "$TARGET_HOST" "docker logs cicd_backend_staging --tail=30"
        return 1
    fi
    
    # Check frontend
    print_status "Checking frontend availability..."
    if ssh "$TARGET_HOST" "curl -f http://localhost:$STAGING_FRONTEND_PORT" >/dev/null 2>&1; then
        print_success "Frontend is accessible"
    else
        print_error "Frontend accessibility check failed"
        ssh "$TARGET_HOST" "docker logs cicd_frontend_staging --tail=20"
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
    container_status=$(ssh "$TARGET_HOST" "cd /opt/cicd-notifier-staging && docker-compose -f docker-compose.staging.yml ps")
    
    # Get system resources
    system_info=$(ssh "$TARGET_HOST" "echo 'Disk Usage:' && df -h / && echo '' && echo 'Memory Usage:' && free -h")
    
    cat > staging-test-report.md << EOF
# Staging Deployment Test Report

**Date**: $(date)
**Branch**: $(git branch --show-current 2>/dev/null || echo "unknown")
**Commit**: $(git rev-parse --short HEAD 2>/dev/null || echo "unknown")
**Target Server**: $TARGET_HOST

## Test Results ✅

- ✅ Prerequisites Check
- ✅ SSH Connectivity
- ✅ Server Prerequisites
- ✅ Docker Images Build & Push
- ✅ Staging Deployment
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

1. ✅ Manual testing completed successfully
2. 🚀 Ready to push to staging branch for automated GitHub Actions
3. 🔍 Monitor the GitHub Actions workflow
4. 📝 Review logs and metrics
5. 🎯 Proceed with production deployment after validation

## Commands for Further Testing

\`\`\`bash
# Test API endpoints
curl http://$TARGET_HOST:$STAGING_BACKEND_PORT/health
curl http://$TARGET_HOST:$STAGING_BACKEND_PORT/api/v1/projects

# View logs
ssh $TARGET_HOST "docker logs cicd_backend_staging"
ssh $TARGET_HOST "docker logs cicd_frontend_staging"
ssh $TARGET_HOST "docker logs cicd_postgres_staging"

# Check container status
ssh $TARGET_HOST "cd /opt/cicd-notifier-staging && docker-compose -f docker-compose.staging.yml ps"
\`\`\`

---
Generated by staging deployment test script
EOF
    
    print_success "Test report generated: staging-test-report.md"
}

# Cleanup function
cleanup() {
    print_status "Cleaning up temporary files..."
    rm -f /tmp/staging-deploy.sh
}

# Main execution
main() {
    clear
    print_header "🧪 Manual Staging Deployment Test"
    print_header "=================================="
    echo ""
    
    # Set trap for cleanup on exit
    trap cleanup EXIT
    
    check_prerequisites
    echo ""
    
    test_ssh_connectivity
    echo ""
    
    check_server_prerequisites
    echo ""
    
    build_and_push_images
    echo ""
    
    deploy_to_staging
    echo ""
    
    run_health_checks
    echo ""
    
    test_external_access
    echo ""
    
    generate_report
    echo ""
    
    print_success "🎉 Manual staging deployment test completed successfully!"
    echo ""
    print_status "Environment is ready at:"
    print_status "  - Frontend: http://$TARGET_HOST:$STAGING_FRONTEND_PORT"
    print_status "  - Backend:  http://$TARGET_HOST:$STAGING_BACKEND_PORT"
    print_status "  - Health:   http://$TARGET_HOST:$STAGING_BACKEND_PORT/health"
    echo ""
    print_status "Next steps:"
    print_status "  1. Review the staging-test-report.md file"
    print_status "  2. Test the application manually"
    print_status "  3. Push to staging branch: git push origin staging"
    print_status "  4. Monitor GitHub Actions workflow"
    echo ""
}

# Check if script is being sourced or executed
if [[ "${BASH_SOURCE[0]}" == "${0}" ]]; then
    main "$@"
fi
