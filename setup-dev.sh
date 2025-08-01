#!/bin/bash

# ================================
# Development Environment Setup Script
# ================================

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m'

# Configuration
POSTGRES_DB="cicd_notifier_dev"
POSTGRES_USER="postgres"
POSTGRES_PASSWORD="postgres123"
POSTGRES_PORT="5433"  # Different port for dev
BACKEND_PORT="8080"
FRONTEND_PORT="5174"  # Vite default

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

setup_dev_database() {
    log_info "Setting up development PostgreSQL database..."
    
    # Create development docker-compose
    cat > docker-compose.dev.yml << EOF
version: '3.8'

services:
  postgres-dev:
    image: postgres:15-alpine
    container_name: cicd-notifier-postgres-dev
    environment:
      POSTGRES_DB: ${POSTGRES_DB}
      POSTGRES_USER: ${POSTGRES_USER}
      POSTGRES_PASSWORD: ${POSTGRES_PASSWORD}
    ports:
      - "${POSTGRES_PORT}:5432"
    volumes:
      - postgres_dev_data:/var/lib/postgresql/data
      - ./backend/scripts/migrations:/docker-entrypoint-initdb.d
    healthcheck:
      test: ["CMD-SHELL", "pg_isready -U ${POSTGRES_USER} -d ${POSTGRES_DB}"]
      interval: 5s
      timeout: 3s
      retries: 5

volumes:
  postgres_dev_data:
EOF
    
    # Start development database
    docker-compose -f docker-compose.dev.yml up -d postgres-dev
    
    # Wait for database
    log_info "Waiting for development database to be ready..."
    timeout=30
    while ! docker-compose -f docker-compose.dev.yml exec -T postgres-dev pg_isready -U ${POSTGRES_USER} -d ${POSTGRES_DB} > /dev/null 2>&1; do
        sleep 1
        timeout=$((timeout - 1))
        if [ $timeout -le 0 ]; then
            log_error "Development database failed to start"
            exit 1
        fi
    done
    
    log_success "Development database is ready"
}

setup_dev_backend() {
    log_info "Setting up backend for development..."
    
    cd backend
    
    # Create development config
    if [ ! -f "config/config.dev.yaml" ]; then
        cp config/config-example.yaml config/config.dev.yaml
        
        # Update for development
        sed -i.bak "s/host: localhost/host: localhost/g" config/config.dev.yaml
        sed -i.bak "s/port: 5432/port: ${POSTGRES_PORT}/g" config/config.dev.yaml
        sed -i.bak "s/user: postgres/user: ${POSTGRES_USER}/g" config/config.dev.yaml
        sed -i.bak "s/password: your_password/password: ${POSTGRES_PASSWORD}/g" config/config.dev.yaml
        sed -i.bak "s/dbname: cicd_notifier/dbname: ${POSTGRES_DB}/g" config/config.dev.yaml
        sed -i.bak "s/port: 8080/port: ${BACKEND_PORT}/g" config/config.dev.yaml
        
        # Enable debug mode
        echo "debug: true" >> config/config.dev.yaml
        echo "log_level: debug" >> config/config.dev.yaml
        
        rm config/config.dev.yaml.bak 2>/dev/null || true
    fi
    
    # Install Air for hot reload (if not installed)
    if ! command -v air &> /dev/null; then
        log_info "Installing Air for hot reload..."
        go install github.com/cosmtrek/air@latest
    fi
    
    # Create Air configuration
    if [ ! -f ".air.toml" ]; then
        cat > .air.toml << 'EOF'
root = "."
testdata_dir = "testdata"
tmp_dir = "tmp"

[build]
  args_bin = []
  bin = "./tmp/main"
  cmd = "go build -o ./tmp/main ./cmd/main.go"
  delay = 1000
  exclude_dir = ["assets", "tmp", "vendor", "testdata"]
  exclude_file = []
  exclude_regex = ["_test.go"]
  exclude_unchanged = false
  follow_symlink = false
  full_bin = ""
  include_dir = []
  include_ext = ["go", "tpl", "tmpl", "html", "yaml", "yml"]
  kill_delay = "0s"
  log = "build-errors.log"
  send_interrupt = false
  stop_on_root = false

[color]
  app = ""
  build = "yellow"
  main = "magenta"
  runner = "green"
  watcher = "cyan"

[log]
  time = false

[misc]
  clean_on_exit = false

[screen]
  clear_on_rebuild = false
EOF
    fi
    
    cd ..
    log_success "Backend development setup completed"
}

setup_dev_frontend() {
    log_info "Setting up frontend for development..."
    
    cd frontend
    
    # Create development environment
    cat > .env.development << EOF
VITE_API_BASE_URL=http://localhost:${BACKEND_PORT}/api/v1
VITE_WS_URL=ws://localhost:${BACKEND_PORT}/ws
VITE_APP_NAME=CI/CD Status Notifier (Dev)
VITE_APP_VERSION=1.0.0-dev
VITE_AUTO_REFRESH_INTERVAL=10000
VITE_WEBSOCKET_RECONNECT_DELAY=3000
VITE_DEBUG=true
EOF
    
    # Install dependencies if not done
    if [ ! -d "node_modules" ]; then
        log_info "Installing frontend dependencies..."
        npm install
    fi
    
    cd ..
    log_success "Frontend development setup completed"
}

create_dev_scripts() {
    log_info "Creating development scripts..."
    
    # Create development start script
    cat > dev-start.sh << 'EOF'
#!/bin/bash

echo "🚀 Starting CI/CD Status Notifier in Development Mode..."

# Start development database
echo "📦 Starting development database..."
docker-compose -f docker-compose.dev.yml up -d postgres-dev

# Wait for database
sleep 5

# Start backend with hot reload
echo "🔧 Starting backend with hot reload..."
cd backend
air -c .air.toml &
BACKEND_PID=$!
cd ..

# Start frontend with hot reload  
echo "🌐 Starting frontend with hot reload..."
cd frontend
npm run dev &
FRONTEND_PID=$!
cd ..

# Save PIDs
mkdir -p logs
echo $BACKEND_PID > logs/dev-backend.pid
echo $FRONTEND_PID > logs/dev-frontend.pid

echo ""
echo "✅ Development environment started!"
echo "📊 Services:"
echo "  Frontend:  http://localhost:5174 (with HMR)"
echo "  Backend:   http://localhost:8080 (with hot reload)"
echo "  Database:  postgresql://localhost:5433/cicd_notifier_dev"
echo ""
echo "🛑 To stop: ./dev-stop.sh"
echo "📊 To check status: ./dev-status.sh"
EOF

    # Create development stop script
    cat > dev-stop.sh << 'EOF'
#!/bin/bash

echo "🛑 Stopping development environment..."

# Stop backend
if [ -f logs/dev-backend.pid ]; then
    kill $(cat logs/dev-backend.pid) 2>/dev/null || true
    rm logs/dev-backend.pid
fi

# Stop frontend
if [ -f logs/dev-frontend.pid ]; then
    kill $(cat logs/dev-frontend.pid) 2>/dev/null || true
    rm logs/dev-frontend.pid
fi

# Stop database
docker-compose -f docker-compose.dev.yml down

echo "✅ Development environment stopped!"
EOF

    # Create development status script
    cat > dev-status.sh << 'EOF'
#!/bin/bash

echo "📊 Development Environment Status:"
echo "=================================="

# Check database
if docker-compose -f docker-compose.dev.yml ps postgres-dev | grep -q "Up"; then
    echo "✅ Database: Running (Port 5433)"
else
    echo "❌ Database: Stopped"
fi

# Check backend
if [ -f logs/dev-backend.pid ] && kill -0 $(cat logs/dev-backend.pid) 2>/dev/null; then
    echo "✅ Backend: Running with hot reload (PID: $(cat logs/dev-backend.pid))"
else
    echo "❌ Backend: Stopped"
fi

# Check frontend
if [ -f logs/dev-frontend.pid ] && kill -0 $(cat logs/dev-frontend.pid) 2>/dev/null; then
    echo "✅ Frontend: Running with HMR (PID: $(cat logs/dev-frontend.pid))"
else
    echo "❌ Frontend: Stopped"
fi

echo ""
echo "🌐 Services:"
echo "  Frontend:  http://localhost:5174"
echo "  Backend:   http://localhost:8080"
echo "  Database:  postgresql://localhost:5433/cicd_notifier_dev"
EOF

    # Create test script
    cat > test.sh << 'EOF'
#!/bin/bash

echo "🧪 Running tests..."

# Backend tests
echo "🔧 Running backend tests..."
cd backend
go test -v ./...
BACKEND_TEST_EXIT=$?
cd ..

# Frontend tests
echo "🌐 Running frontend tests..."
cd frontend
npm test
FRONTEND_TEST_EXIT=$?
cd ..

if [ $BACKEND_TEST_EXIT -eq 0 ] && [ $FRONTEND_TEST_EXIT -eq 0 ]; then
    echo "✅ All tests passed!"
    exit 0
else
    echo "❌ Some tests failed!"
    exit 1
fi
EOF

    # Make all scripts executable
    chmod +x dev-start.sh dev-stop.sh dev-status.sh test.sh
    
    log_success "Development scripts created"
}

create_makefile() {
    log_info "Creating Makefile for easy commands..."
    
    cat > Makefile << 'EOF'
.PHONY: help dev-start dev-stop dev-status deploy start stop status test clean build

help: ## Show this help message
	@echo "Available commands:"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-15s\033[0m %s\n", $$1, $$2}'

dev-start: ## Start development environment
	./dev-start.sh

dev-stop: ## Stop development environment  
	./dev-stop.sh

dev-status: ## Check development environment status
	./dev-status.sh

deploy: ## Deploy production environment
	./deploy.sh

start: ## Start production environment
	./start.sh

stop: ## Stop production environment
	./stop.sh

status: ## Check production environment status
	./status.sh

test: ## Run all tests
	./test.sh

clean: ## Clean up containers and volumes
	docker-compose down -v
	docker-compose -f docker-compose.dev.yml down -v
	docker system prune -f

build: ## Build applications
	cd backend && go build -o bin/cicd-notifier cmd/main.go
	cd frontend && npm run build

logs-backend: ## View backend logs
	tail -f logs/backend.log

logs-frontend: ## View frontend logs  
	tail -f logs/frontend.log

logs-db: ## View database logs
	docker-compose logs -f postgres

dev-logs-backend: ## View development backend logs
	tail -f logs/dev-backend.log

dev-logs-frontend: ## View development frontend logs
	tail -f logs/dev-frontend.log

dev-logs-db: ## View development database logs
	docker-compose -f docker-compose.dev.yml logs -f postgres-dev
EOF
    
    log_success "Makefile created"
}

print_dev_summary() {
    echo ""
    echo "======================================"
    echo "🛠️  DEVELOPMENT SETUP COMPLETED! 🛠️"
    echo "======================================"
    echo ""
    echo "📋 Quick Start Commands:"
    echo "  make dev-start   - Start development environment"
    echo "  make dev-stop    - Stop development environment"
    echo "  make dev-status  - Check status"
    echo "  make test        - Run tests"
    echo "  make help        - Show all commands"
    echo ""
    echo "🔧 Development Features:"
    echo "  ✅ Hot reload for backend (Air)"
    echo "  ✅ Hot module replacement for frontend (Vite HMR)"
    echo "  ✅ Separate development database"
    echo "  ✅ Debug logging enabled"
    echo "  ✅ Development environment variables"
    echo ""
    echo "🚀 To start developing:"
    echo "  1. make dev-start"
    echo "  2. Open http://localhost:5174 (Frontend)"
    echo "  3. Backend API available at http://localhost:8080"
    echo ""
    log_success "Development environment is ready!"
}

main() {
    log_info "Setting up development environment..."
    echo ""
    
    setup_dev_database
    setup_dev_backend 
    setup_dev_frontend
    create_dev_scripts
    create_makefile
    
    print_dev_summary
}

# Run main function
main "$@"
