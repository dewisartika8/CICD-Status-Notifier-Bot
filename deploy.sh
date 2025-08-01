#!/bin/bash

# ================================
# CI/CD Status Notifier Bot Deployment Script
# ================================

set -e  # Exit on any error

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Configuration
POSTGRES_DB="cicd_notifier"
POSTGRES_USER="postgres"
POSTGRES_PASSWORD="postgres123"
POSTGRES_PORT="5432"
BACKEND_PORT="8080"
FRONTEND_PORT="3000"

# Functions
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

check_requirements() {
    log_info "Checking system requirements..."
    
    # Check if Docker is installed
    if ! command -v docker &> /dev/null; then
        log_error "Docker is not installed. Please install Docker first."
        exit 1
    fi
    
    # Check if Docker Compose is installed
    if ! command -v docker-compose &> /dev/null; then
        log_error "Docker Compose is not installed. Please install Docker Compose first."
        exit 1
    fi
    
    # Check if Go is installed
    if ! command -v go &> /dev/null; then
        log_error "Go is not installed. Please install Go 1.19+ first."
        exit 1
    fi
    
    # Check if Node.js is installed
    if ! command -v node &> /dev/null; then
        log_error "Node.js is not installed. Please install Node.js 18+ first."
        exit 1
    fi
    
    # Check if npm is installed
    if ! command -v npm &> /dev/null; then
        log_error "npm is not installed. Please install npm first."
        exit 1
    fi
    
    log_success "All requirements are satisfied"
}

setup_database() {
    log_info "Setting up PostgreSQL database with Docker..."
    
    # Stop existing containers if running
    docker-compose down 2>/dev/null || true
    
    # Create docker-compose.yml if it doesn't exist
    if [ ! -f "docker-compose.yml" ]; then
        log_info "Creating docker-compose.yml..."
        cat > docker-compose.yml << EOF
version: '3.8'

services:
  postgres:
    image: postgres:15-alpine
    container_name: cicd-notifier-postgres
    environment:
      POSTGRES_DB: ${POSTGRES_DB}
      POSTGRES_USER: ${POSTGRES_USER}
      POSTGRES_PASSWORD: ${POSTGRES_PASSWORD}
    ports:
      - "${POSTGRES_PORT}:5432"
    volumes:
      - postgres_data:/var/lib/postgresql/data
      - ./backend/scripts/migrations:/docker-entrypoint-initdb.d
    healthcheck:
      test: ["CMD-SHELL", "pg_isready -U ${POSTGRES_USER} -d ${POSTGRES_DB}"]
      interval: 10s
      timeout: 5s
      retries: 5
    restart: unless-stopped

volumes:
  postgres_data:
EOF
        log_success "docker-compose.yml created"
    fi
    
    # Start PostgreSQL
    log_info "Starting PostgreSQL container..."
    docker-compose up -d postgres
    
    # Wait for PostgreSQL to be ready
    log_info "Waiting for PostgreSQL to be ready..."
    timeout=60
    while ! docker-compose exec -T postgres pg_isready -U ${POSTGRES_USER} -d ${POSTGRES_DB} > /dev/null 2>&1; do
        sleep 2
        timeout=$((timeout - 2))
        if [ $timeout -le 0 ]; then
            log_error "PostgreSQL failed to start within 60 seconds"
            exit 1
        fi
    done
    
    log_success "PostgreSQL is ready"
}

setup_backend() {
    log_info "Setting up backend application..."
    
    cd backend
    
    # Create config.yaml if it doesn't exist
    if [ ! -f "config/config.yaml" ]; then
        log_info "Creating backend configuration..."
        cp config/config-example.yaml config/config.yaml
        
        # Update database configuration
        sed -i.bak "s/host: localhost/host: localhost/g" config/config.yaml
        sed -i.bak "s/port: 5432/port: ${POSTGRES_PORT}/g" config/config.yaml
        sed -i.bak "s/user: postgres/user: ${POSTGRES_USER}/g" config/config.yaml
        sed -i.bak "s/password: your_password/password: ${POSTGRES_PASSWORD}/g" config/config.yaml
        sed -i.bak "s/dbname: cicd_notifier/dbname: ${POSTGRES_DB}/g" config/config.yaml
        sed -i.bak "s/port: 8080/port: ${BACKEND_PORT}/g" config/config.yaml
        
        rm config/config.yaml.bak 2>/dev/null || true
        log_success "Backend configuration created"
    fi
    
    # Download dependencies
    log_info "Downloading Go dependencies..."
    go mod download
    go mod tidy
    
    # Run database migrations
    log_info "Running database migrations..."
    chmod +x scripts/migrate.sh
    ./scripts/migrate.sh up || {
        log_warning "Migration script failed, trying manual migration..."
        # Manual migration fallback
        for migration in scripts/migrations/*.up.sql; do
            if [ -f "$migration" ]; then
                log_info "Running migration: $(basename $migration)"
                docker-compose exec -T postgres psql -U ${POSTGRES_USER} -d ${POSTGRES_DB} -f - < "$migration" || true
            fi
        done
    }
    
    # Build backend
    log_info "Building backend application..."
    go build -o bin/cicd-notifier cmd/main.go
    
    cd ..
    log_success "Backend setup completed"
}

setup_frontend() {
    log_info "Setting up frontend application..."
    
    cd frontend
    
    # Create .env file if it doesn't exist
    if [ ! -f ".env" ]; then
        log_info "Creating frontend environment configuration..."
        cat > .env << EOF
VITE_API_BASE_URL=http://localhost:${BACKEND_PORT}/api/v1
VITE_WS_URL=ws://localhost:${BACKEND_PORT}/ws
VITE_APP_NAME=CI/CD Status Notifier
VITE_APP_VERSION=1.0.0
VITE_AUTO_REFRESH_INTERVAL=30000
VITE_WEBSOCKET_RECONNECT_DELAY=5000
EOF
        log_success "Frontend environment configuration created"
    fi
    
    # Install dependencies
    log_info "Installing npm dependencies..."
    npm install
    
    # Build frontend
    log_info "Building frontend application..."
    npm run build
    
    cd ..
    log_success "Frontend setup completed"
}

start_applications() {
    log_info "Starting all applications..."
    
    # Start backend in background
    log_info "Starting backend server on port ${BACKEND_PORT}..."
    cd backend
    nohup ./bin/cicd-notifier > ../logs/backend.log 2>&1 &
    BACKEND_PID=$!
    echo $BACKEND_PID > ../logs/backend.pid
    cd ..
    
    # Wait a moment for backend to start
    sleep 3
    
    # Check if backend is running
    if kill -0 $BACKEND_PID 2>/dev/null; then
        log_success "Backend server started (PID: $BACKEND_PID)"
    else
        log_error "Backend server failed to start"
        exit 1
    fi
    
    # Start frontend in background
    log_info "Starting frontend server on port ${FRONTEND_PORT}..."
    cd frontend
    nohup npm run preview -- --port ${FRONTEND_PORT} --host > ../logs/frontend.log 2>&1 &
    FRONTEND_PID=$!
    echo $FRONTEND_PID > ../logs/frontend.pid
    cd ..
    
    # Wait a moment for frontend to start
    sleep 3
    
    # Check if frontend is running
    if kill -0 $FRONTEND_PID 2>/dev/null; then
        log_success "Frontend server started (PID: $FRONTEND_PID)"
    else
        log_error "Frontend server failed to start"
        exit 1
    fi
}

create_systemd_services() {
    log_info "Creating systemd service files (optional)..."
    
    # Create systemd service for backend
    cat > cicd-notifier-backend.service << EOF
[Unit]
Description=CI/CD Status Notifier Backend
After=network.target postgresql.service

[Service]
Type=simple
User=$USER
WorkingDirectory=$(pwd)/backend
ExecStart=$(pwd)/backend/bin/cicd-notifier
Restart=always
RestartSec=5
Environment=PATH=/usr/local/bin:/usr/bin:/bin
Environment=GO_ENV=production

[Install]
WantedBy=multi-user.target
EOF
    
    # Create systemd service for frontend
    cat > cicd-notifier-frontend.service << EOF
[Unit]
Description=CI/CD Status Notifier Frontend
After=network.target

[Service]
Type=simple
User=$USER
WorkingDirectory=$(pwd)/frontend
ExecStart=/usr/bin/npm run preview -- --port ${FRONTEND_PORT} --host
Restart=always
RestartSec=5
Environment=NODE_ENV=production

[Install]
WantedBy=multi-user.target
EOF
    
    log_info "Systemd service files created. To install them:"
    log_info "  sudo cp cicd-notifier-*.service /etc/systemd/system/"
    log_info "  sudo systemctl daemon-reload"
    log_info "  sudo systemctl enable cicd-notifier-backend cicd-notifier-frontend"
    log_info "  sudo systemctl start cicd-notifier-backend cicd-notifier-frontend"
}

create_management_scripts() {
    log_info "Creating management scripts..."
    
    # Create logs directory
    mkdir -p logs
    
    # Create start script
    cat > start.sh << 'EOF'
#!/bin/bash
echo "Starting CI/CD Status Notifier..."
docker-compose up -d postgres
sleep 5

cd backend && nohup ./bin/cicd-notifier > ../logs/backend.log 2>&1 &
echo $! > ../logs/backend.pid
cd ..

cd frontend && nohup npm run preview -- --port 3000 --host > ../logs/frontend.log 2>&1 &
echo $! > ../logs/frontend.pid
cd ..

echo "All services started!"
echo "Frontend: http://localhost:3000"
echo "Backend: http://localhost:8080"
echo "Database: localhost:5432"
EOF
    
    # Create stop script
    cat > stop.sh << 'EOF'
#!/bin/bash
echo "Stopping CI/CD Status Notifier..."

# Stop backend
if [ -f logs/backend.pid ]; then
    kill $(cat logs/backend.pid) 2>/dev/null || true
    rm logs/backend.pid
fi

# Stop frontend
if [ -f logs/frontend.pid ]; then
    kill $(cat logs/frontend.pid) 2>/dev/null || true
    rm logs/frontend.pid
fi

# Stop database
docker-compose down

echo "All services stopped!"
EOF
    
    # Create status script
    cat > status.sh << 'EOF'
#!/bin/bash
echo "CI/CD Status Notifier Status:"
echo "============================="

# Check database
if docker-compose ps postgres | grep -q "Up"; then
    echo "✅ Database: Running"
else
    echo "❌ Database: Stopped"
fi

# Check backend
if [ -f logs/backend.pid ] && kill -0 $(cat logs/backend.pid) 2>/dev/null; then
    echo "✅ Backend: Running (PID: $(cat logs/backend.pid))"
else
    echo "❌ Backend: Stopped"
fi

# Check frontend
if [ -f logs/frontend.pid ] && kill -0 $(cat logs/frontend.pid) 2>/dev/null; then
    echo "✅ Frontend: Running (PID: $(cat logs/frontend.pid))"
else
    echo "❌ Frontend: Stopped"
fi

echo ""
echo "Services:"
echo "- Frontend: http://localhost:3000"
echo "- Backend API: http://localhost:8080"
echo "- Database: postgresql://localhost:5432/cicd_notifier"
EOF
    
    # Make scripts executable
    chmod +x start.sh stop.sh status.sh
    
    log_success "Management scripts created (start.sh, stop.sh, status.sh)"
}

health_check() {
    log_info "Performing health checks..."
    
    # Check database
    if docker-compose exec -T postgres pg_isready -U ${POSTGRES_USER} -d ${POSTGRES_DB} > /dev/null 2>&1; then
        log_success "✅ Database: Healthy"
    else
        log_error "❌ Database: Unhealthy"
        return 1
    fi
    
    # Check backend
    sleep 5  # Give backend time to start
    if curl -f http://localhost:${BACKEND_PORT}/health > /dev/null 2>&1; then
        log_success "✅ Backend: Healthy"
    else
        log_warning "⚠️  Backend: Health check failed (might still be starting)"
    fi
    
    # Check frontend
    if curl -f http://localhost:${FRONTEND_PORT} > /dev/null 2>&1; then
        log_success "✅ Frontend: Healthy"
    else
        log_warning "⚠️  Frontend: Health check failed (might still be starting)"
    fi
}

print_summary() {
    echo ""
    echo "======================================"
    echo "🎉 DEPLOYMENT COMPLETED SUCCESSFULLY! 🎉"
    echo "======================================"
    echo ""
    echo "📋 Service Information:"
    echo "  🌐 Frontend:  http://localhost:${FRONTEND_PORT}"
    echo "  🔧 Backend:   http://localhost:${BACKEND_PORT}"
    echo "  🗄️  Database: postgresql://localhost:${POSTGRES_PORT}/${POSTGRES_DB}"
    echo ""
    echo "📁 Management Commands:"
    echo "  ./start.sh   - Start all services"
    echo "  ./stop.sh    - Stop all services" 
    echo "  ./status.sh  - Check service status"
    echo ""
    echo "📊 Logs Location:"
    echo "  Backend:  logs/backend.log"
    echo "  Frontend: logs/frontend.log"
    echo "  Database: docker-compose logs postgres"
    echo ""
    echo "🔧 Useful Commands:"
    echo "  docker-compose logs -f postgres  # View database logs"
    echo "  tail -f logs/backend.log         # View backend logs"
    echo "  tail -f logs/frontend.log        # View frontend logs"
    echo ""
    log_success "All services are now running!"
}

# Main deployment process
main() {
    log_info "Starting CI/CD Status Notifier Bot deployment..."
    echo ""
    
    check_requirements
    setup_database
    setup_backend
    setup_frontend
    create_management_scripts
    start_applications
    create_systemd_services
    
    # Wait a bit before health check
    log_info "Waiting for services to fully start..."
    sleep 10
    
    health_check
    print_summary
}

# Handle script termination
trap 'log_error "Deployment interrupted"; exit 1' SIGINT SIGTERM

# Run main function
main "$@"
