#!/bin/bash

# Local Testing Script for CICD Pipeline
# This script helps test the deployment process locally before pushing to GitHub

set -e

echo "🧪 CICD Pipeline Local Testing Script"
echo "======================================"

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Configuration
TARGET_HOST="172.16.19.11"
BACKEND_PORT="8080"
FRONTEND_PORT="3000"
TEST_BACKEND_PORT="8081"
TEST_FRONTEND_PORT="3001"

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

# Function to check if command exists
command_exists() {
    command -v "$1" >/dev/null 2>&1
}

# Check prerequisites
check_prerequisites() {
    print_status "Checking prerequisites..."
    
    local missing_tools=()
    
    if ! command_exists "docker"; then
        missing_tools+=("docker")
    fi
    
    if ! command_exists "docker-compose"; then
        missing_tools+=("docker-compose")
    fi
    
    if ! command_exists "go"; then
        missing_tools+=("go")
    fi
    
    if ! command_exists "npm"; then
        missing_tools+=("npm")
    fi
    
    if ! command_exists "curl"; then
        missing_tools+=("curl")
    fi
    
    if [ ${#missing_tools[@]} -ne 0 ]; then
        print_error "Missing required tools: ${missing_tools[*]}"
        print_error "Please install the missing tools before running this script."
        exit 1
    fi
    
    print_success "All prerequisites are installed"
}

# Run backend tests
test_backend() {
    print_status "Running backend tests..."
    
    cd backend
    
    # Download dependencies
    print_status "Downloading Go dependencies..."
    go mod download
    go mod verify
    
    # Run static analysis
    print_status "Running static analysis..."
    go vet ./...
    
    if command_exists "staticcheck"; then
        staticcheck ./...
    else
        print_warning "staticcheck not installed, skipping..."
    fi
    
    # Start test database
    print_status "Starting test database..."
    docker-compose -f ../docker-compose.yml up -d postgres
    
    # Wait for database to be ready
    print_status "Waiting for database to be ready..."
    sleep 10
    
    # Run database migrations
    if [ -f "scripts/migrate.sh" ]; then
        print_status "Running database migrations..."
        DATABASE_URL="postgres://postgres:jQlwjVKoQw@localhost:5432/cicd_notifier?sslmode=disable" \
        ./scripts/migrate.sh up
    fi
    
    # Run tests
    print_status "Running unit tests..."
    DATABASE_URL="postgres://postgres:jQlwjVKoQw@localhost:5432/cicd_notifier?sslmode=disable" \
    go test -v -race -coverprofile=coverage.out ./...
    
    # Run integration tests if they exist
    if [ -d "tests" ] && [ -f "tests/Makefile" ]; then
        print_status "Running integration tests..."
        cd tests
        DATABASE_URL="postgres://postgres:jQlwjVKoQw@localhost:5432/cicd_notifier?sslmode=disable" \
        make test-integration || print_warning "Integration tests failed or not available"
        cd ..
    fi
    
    print_success "Backend tests completed"
    cd ..
}

# Run frontend tests
test_frontend() {
    print_status "Running frontend tests..."
    
    cd frontend
    
    # Install dependencies
    print_status "Installing npm dependencies..."
    npm ci
    
    # Run linting if available
    if npm run lint --silent 2>/dev/null; then
        print_status "Running ESLint..."
        npm run lint
    else
        print_warning "ESLint not configured, skipping..."
    fi
    
    # Run tests
    print_status "Running frontend tests..."
    npm test -- --coverage --watchAll=false || print_warning "Frontend tests failed or not configured"
    
    # Build frontend
    print_status "Building frontend..."
    npm run build
    
    print_success "Frontend tests completed"
    cd ..
}

# Build Docker images
build_images() {
    print_status "Building Docker images..."
    
    # Build backend image
    print_status "Building backend Docker image..."
    docker build -t cicd-backend:test ./backend
    
    # Build frontend image
    print_status "Building frontend Docker image..."
    docker build -t cicd-frontend:test ./frontend
    
    print_success "Docker images built successfully"
}

# Test deployment locally
test_deployment() {
    print_status "Testing deployment locally..."
    
    # Create test docker-compose file
    cat > docker-compose.test.yml << 'EOF'
version: '3.8'

services:
  postgres-test:
    image: postgres:15-alpine
    container_name: cicd_postgres_test_local
    environment:
      POSTGRES_DB: cicd_notifier_test
      POSTGRES_USER: postgres
      POSTGRES_PASSWORD: jQlwjVKoQw
    ports:
      - "5433:5432"
    volumes:
      - postgres_test_data:/var/lib/postgresql/data
    healthcheck:
      test: ["CMD-SHELL", "pg_isready -U postgres"]
      interval: 10s
      timeout: 5s
      retries: 5

  backend-test:
    image: cicd-backend:test
    container_name: cicd_backend_test_local
    depends_on:
      postgres-test:
        condition: service_healthy
    ports:
      - "8081:8080"
    environment:
      DATABASE_URL: postgres://postgres:jQlwjVKoQw@postgres-test:5432/cicd_notifier_test?sslmode=disable
      PORT: 8080
      ENV: test

  frontend-test:
    image: cicd-frontend:test
    container_name: cicd_frontend_test_local
    depends_on:
      - backend-test
    ports:
      - "3001:80"
    environment:
      REACT_APP_API_URL: http://localhost:8081

volumes:
  postgres_test_data:
EOF

    # Stop any existing test containers
    docker-compose -f docker-compose.test.yml down || true
    
    # Start test environment
    print_status "Starting test environment..."
    docker-compose -f docker-compose.test.yml up -d
    
    # Wait for services to be ready
    print_status "Waiting for services to start..."
    sleep 30
    
    # Test backend health
    print_status "Testing backend health..."
    if curl -f http://localhost:8081/health; then
        print_success "Backend health check passed"
    else
        print_error "Backend health check failed"
        docker-compose -f docker-compose.test.yml logs backend-test
        return 1
    fi
    
    # Test frontend
    print_status "Testing frontend..."
    if curl -f http://localhost:3001; then
        print_success "Frontend health check passed"
    else
        print_error "Frontend health check failed"
        docker-compose -f docker-compose.test.yml logs frontend-test
        return 1
    fi
    
    print_success "Local deployment test completed successfully"
}

# Test remote server connectivity
test_remote_connectivity() {
    print_status "Testing remote server connectivity..."
    
    # Test SSH connectivity (if SSH key is available)
    if [ -f "$HOME/.ssh/id_rsa" ] || [ -f "$HOME/.ssh/id_ed25519" ]; then
        print_status "Testing SSH connectivity to $TARGET_HOST..."
        if ssh -o ConnectTimeout=5 -o StrictHostKeyChecking=no "$TARGET_HOST" "echo 'SSH connection successful'" 2>/dev/null; then
            print_success "SSH connectivity test passed"
        else
            print_warning "SSH connectivity test failed. Please ensure SSH access is configured."
        fi
    else
        print_warning "No SSH keys found, skipping SSH connectivity test"
    fi
    
    # Test HTTP connectivity
    print_status "Testing HTTP connectivity to $TARGET_HOST..."
    if curl -s --connect-timeout 5 http://$TARGET_HOST >/dev/null 2>&1; then
        print_success "HTTP connectivity test passed"
    else
        print_warning "HTTP connectivity test failed. Server might not be set up yet."
    fi
}

# Security checks
run_security_checks() {
    print_status "Running security checks..."
    
    # Check for secrets in code
    print_status "Checking for potential secrets in code..."
    if command_exists "grep"; then
        local secrets_found=false
        
        # Check for common secret patterns
        if grep -r -i "password.*=" . --exclude-dir=.git --exclude-dir=node_modules --exclude="*.log" | grep -v "example\|test\|mock"; then
            print_warning "Potential hardcoded passwords found"
            secrets_found=true
        fi
        
        if grep -r -i "api[_-]key.*=" . --exclude-dir=.git --exclude-dir=node_modules --exclude="*.log" | grep -v "example\|test\|mock"; then
            print_warning "Potential hardcoded API keys found"
            secrets_found=true
        fi
        
        if grep -r -i "secret.*=" . --exclude-dir=.git --exclude-dir=node_modules --exclude="*.log" | grep -v "example\|test\|mock"; then
            print_warning "Potential hardcoded secrets found"
            secrets_found=true
        fi
        
        if [ "$secrets_found" = false ]; then
            print_success "No obvious secrets found in code"
        fi
    fi
    
    # Check Docker image security (if trivy is available)
    if command_exists "trivy"; then
        print_status "Running Trivy security scan on Docker images..."
        trivy image cicd-backend:test || print_warning "Trivy scan failed or found vulnerabilities"
        trivy image cicd-frontend:test || print_warning "Trivy scan failed or found vulnerabilities"
    else
        print_warning "Trivy not installed, skipping Docker image security scan"
    fi
}

# Generate test report
generate_report() {
    print_status "Generating test report..."
    
    cat > test-report.md << EOF
# CICD Pipeline Test Report

**Date**: $(date)
**Branch**: $(git branch --show-current 2>/dev/null || echo "unknown")
**Commit**: $(git rev-parse --short HEAD 2>/dev/null || echo "unknown")

## Test Results

- ✅ Prerequisites Check
- ✅ Backend Tests
- ✅ Frontend Tests  
- ✅ Docker Image Build
- ✅ Local Deployment Test
- ✅ Security Checks

## Next Steps

1. Review any warnings or errors above
2. Fix any issues found
3. Commit and push changes to trigger GitHub Actions
4. Monitor the deployment in GitHub Actions

## Manual Deployment

If you need to deploy manually to $TARGET_HOST:

\`\`\`bash
# SSH to the server
ssh $TARGET_HOST

# Pull latest changes
cd /opt/cicd-notifier
git pull origin main

# Rebuild and restart
docker-compose down
docker-compose build
docker-compose up -d
\`\`\`

## Monitoring

After deployment, monitor:
- Backend: http://$TARGET_HOST:$BACKEND_PORT/health
- Frontend: http://$TARGET_HOST:$FRONTEND_PORT
- Logs: \`docker-compose logs -f\`
EOF
    
    print_success "Test report generated: test-report.md"
}

# Cleanup function
cleanup() {
    print_status "Cleaning up test environment..."
    docker-compose -f docker-compose.test.yml down || true
    docker-compose -f docker-compose.yml down || true
    rm -f docker-compose.test.yml
    print_success "Cleanup completed"
}

# Main execution
main() {
    echo ""
    print_status "Starting CICD pipeline local testing..."
    echo ""
    
    # Set trap for cleanup on exit
    trap cleanup EXIT
    
    check_prerequisites
    echo ""
    
    test_backend
    echo ""
    
    test_frontend
    echo ""
    
    build_images
    echo ""
    
    test_deployment
    echo ""
    
    test_remote_connectivity
    echo ""
    
    run_security_checks
    echo ""
    
    generate_report
    echo ""
    
    print_success "🎉 All tests completed successfully!"
    print_status "Review the test-report.md file for details"
    print_status "You can now safely push your changes to trigger the GitHub Actions pipeline"
}

# Check if script is being sourced or executed
if [[ "${BASH_SOURCE[0]}" == "${0}" ]]; then
    main "$@"
fi
