# CI/CD Status Notifier Bot - Deployment Guide

Complete deployment and development setup for the CI/CD Status Notifier Bot application with PostgreSQL, Go backend, and React frontend.

## 🚀 Quick Start

### Production Deployment
```bash
# Make scripts executable
chmod +x *.sh

# Deploy everything
./deploy.sh
```

### Development Environment
```bash
# Setup development environment
./setup-dev.sh

# Start development servers
make dev-start
```

## 📋 Prerequisites

### Required Software
- **Docker** (20.10+) and **Docker Compose** (2.0+)
- **Go** (1.19+) 
- **Node.js** (18+) and **npm** (9+)
- **Git** (for updates)
- **curl** (for health checks)

### System Requirements
- **RAM**: 4GB minimum, 8GB recommended
- **Disk**: 10GB free space minimum
- **OS**: Linux, macOS, or Windows with WSL2

## 🏗️ Architecture

```
┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐
│   Frontend      │    │    Backend      │    │   PostgreSQL    │
│   React + Vite  │◄──►│   Go + Fiber    │◄──►│   Docker        │
│   Port: 3000    │    │   Port: 8080    │    │   Port: 5432    │
└─────────────────┘    └─────────────────┘    └─────────────────┘
```

## 📂 Directory Structure

```
CICD-Status-Notifier-Bot/
├── backend/                 # Go backend application
│   ├── cmd/main.go         # Application entry point
│   ├── config/             # Configuration files
│   ├── internal/           # Internal packages
│   └── scripts/            # Database migrations
├── frontend/               # React frontend application
│   ├── src/                # Source code
│   ├── public/             # Static assets
│   └── dist/               # Built application
├── docs/                   # Documentation
├── logs/                   # Application logs
├── backups/                # Database backups
├── deploy.sh               # Production deployment
├── setup-dev.sh            # Development setup
├── maintenance.sh          # Backup and monitoring
└── Makefile               # Quick commands
```

## 🛠️ Scripts Overview

### Production Scripts

#### `deploy.sh` - Full Production Deployment
Comprehensive deployment script that:
- ✅ Checks system requirements
- ✅ Sets up PostgreSQL with Docker
- ✅ Configures and builds backend
- ✅ Builds and optimizes frontend
- ✅ Creates systemd services
- ✅ Starts all applications
- ✅ Performs health checks

```bash
./deploy.sh
```

#### Management Scripts (Auto-generated)
- `start.sh` - Start all production services
- `stop.sh` - Stop all production services  
- `status.sh` - Check service status

### Development Scripts

#### `setup-dev.sh` - Development Environment Setup
Sets up development environment with:
- ✅ Separate development database (port 5433)
- ✅ Hot reload for backend (Air)
- ✅ Hot module replacement for frontend (Vite HMR)
- ✅ Debug logging enabled
- ✅ Development configuration

```bash
./setup-dev.sh
make dev-start
```

#### Development Management
- `dev-start.sh` - Start development environment
- `dev-stop.sh` - Stop development environment
- `dev-status.sh` - Check development status

### Maintenance Scripts

#### `maintenance.sh` - Backup and Monitoring
Comprehensive maintenance operations:

```bash
# Create backup
./maintenance.sh backup

# Restore from backup
./maintenance.sh restore backups/db_backup_20250801_120000.sql.gz

# Clean old backups (older than 7 days)
./maintenance.sh cleanup

# System monitoring
./maintenance.sh monitor

# Health checks
./maintenance.sh health

# Update applications
./maintenance.sh update

# View logs
./maintenance.sh logs backend 100
./maintenance.sh logs frontend 50
./maintenance.sh logs database
```

## 🎯 Makefile Commands

```bash
# Development
make dev-start      # Start development environment
make dev-stop       # Stop development environment
make dev-status     # Check development status

# Production
make deploy         # Deploy production environment
make start          # Start production services
make stop           # Stop production services
make status         # Check production status

# Maintenance
make test           # Run all tests
make build          # Build applications
make clean          # Clean up containers and volumes

# Logs
make logs-backend   # View backend logs
make logs-frontend  # View frontend logs
make logs-db        # View database logs

# Help
make help           # Show all available commands
```

## 🔧 Configuration

### Backend Configuration
Located in `backend/config/config.yaml`:

```yaml
database:
  host: localhost
  port: 5432
  user: postgres
  password: postgres123
  dbname: cicd_notifier
  sslmode: disable

server:
  port: 8080
  host: 0.0.0.0

telegram:
  bot_token: "your_bot_token"
  
webhook:
  secret: "your_webhook_secret"
```

### Frontend Configuration
Environment variables in `frontend/.env`:

```env
VITE_API_BASE_URL=http://localhost:8080/api/v1
VITE_WS_URL=ws://localhost:8080/ws
VITE_APP_NAME=CI/CD Status Notifier
VITE_APP_VERSION=1.0.0
```

### Database Configuration
PostgreSQL runs in Docker with:
- **Database**: `cicd_notifier`
- **User**: `postgres`
- **Password**: `postgres123`
- **Port**: `5432` (production), `5433` (development)

## 📊 Service Endpoints

### Production URLs
- **Frontend**: http://localhost:3000
- **Backend API**: http://localhost:8080
- **Database**: postgresql://localhost:5432/cicd_notifier

### Development URLs  
- **Frontend**: http://localhost:5174 (with HMR)
- **Backend API**: http://localhost:8080 (with hot reload)
- **Database**: postgresql://localhost:5433/cicd_notifier_dev

### API Endpoints
- `GET /api/v1/health` - Health check
- `GET /api/v1/projects` - List projects
- `GET /api/v1/dashboard/analytics` - Dashboard analytics
- `GET /api/v1/webhooks` - Webhook configurations
- `GET /api/v1/telegram/subscriptions` - Telegram subscriptions

## 🚦 Health Monitoring

### Automated Health Checks
The deployment includes comprehensive health monitoring:

```bash
# Check all services
./maintenance.sh health

# Monitor system resources
./maintenance.sh monitor
```

### Manual Health Checks
```bash
# Database
docker-compose exec postgres pg_isready -U postgres -d cicd_notifier

# Backend
curl http://localhost:8080/health

# Frontend
curl http://localhost:3000
```

## 🗃️ Backup and Recovery

### Automated Backups
```bash
# Create backup (database + configuration)
./maintenance.sh backup

# Files created:
# - backups/db_backup_YYYYMMDD_HHMMSS.sql.gz
# - backups/config_backup_YYYYMMDD_HHMMSS.tar.gz
```

### Restore from Backup
```bash
# Restore database
./maintenance.sh restore backups/db_backup_20250801_120000.sql.gz

# Restore configuration (manual)
tar -xzf backups/config_backup_20250801_120000.tar.gz
```

### Cleanup Old Backups
```bash
# Clean backups older than 7 days (default)
./maintenance.sh cleanup

# Clean backups older than 14 days
./maintenance.sh cleanup 14
```

## 🐳 Docker Operations

### View Container Status
```bash
docker-compose ps
```

### View Logs
```bash
# Database logs
docker-compose logs -f postgres

# All services
docker-compose logs -f
```

### Database Access
```bash
# Connect to database
docker-compose exec postgres psql -U postgres -d cicd_notifier

# Run SQL commands
docker-compose exec postgres psql -U postgres -d cicd_notifier -c "SELECT * FROM projects;"
```

## 🔄 Updates and Maintenance

### Update Applications
```bash
# Pull latest code and update dependencies
./maintenance.sh update

# Restart services
make stop && make start
```

### Database Migrations
```bash
cd backend
chmod +x scripts/migrate.sh
./scripts/migrate.sh up
```

## 🛡️ Security Considerations

### Production Security
- Change default passwords in configuration
- Use environment variables for secrets
- Enable SSL/TLS in production
- Configure firewall rules
- Regular security updates

### Development Security
- Development database uses separate port
- Debug mode enabled (disable in production)
- CORS configured for development

## 🚨 Troubleshooting

### Common Issues

#### Database Connection Failed
```bash
# Check if database is running
docker-compose ps postgres

# Check database logs
docker-compose logs postgres

# Restart database
docker-compose restart postgres
```

#### Backend Won't Start
```bash
# Check backend logs
tail -f logs/backend.log

# Check configuration
cat backend/config/config.yaml

# Rebuild backend
cd backend && go build -o bin/cicd-notifier cmd/main.go
```

#### Frontend Build Errors
```bash
# Check frontend logs
tail -f logs/frontend.log

# Clear cache and reinstall
cd frontend
rm -rf node_modules package-lock.json
npm install
npm run build
```

#### Port Already in Use
```bash
# Find process using port
lsof -i :8080
lsof -i :3000
lsof -i :5432

# Kill process
kill -9 <PID>
```

### Log Locations
- Backend: `logs/backend.log`
- Frontend: `logs/frontend.log`  
- Database: `docker-compose logs postgres`
- Development Backend: `logs/dev-backend.log`
- Development Frontend: `logs/dev-frontend.log`

## 🎮 Development Workflow

### Starting Development
```bash
# First time setup
./setup-dev.sh

# Daily development
make dev-start
# Edit code (auto-reload enabled)
make dev-stop
```

### Testing
```bash
# Run all tests
make test

# Backend tests only
cd backend && go test -v ./...

# Frontend tests only
cd frontend && npm test
```

### Building for Production
```bash
# Build both applications
make build

# Backend only
cd backend && go build -o bin/cicd-notifier cmd/main.go

# Frontend only
cd frontend && npm run build
```

## 📈 Performance Optimization

### Backend Optimization
- Go binary compilation with optimizations
- Connection pooling for database
- Request logging and monitoring
- Memory management

### Frontend Optimization
- Vite build optimization
- Code splitting (vendor, mui, charts)
- Gzip compression
- Asset optimization

### Database Optimization
- Connection pooling
- Index optimization
- Regular maintenance
- Backup scheduling

## 🤝 Contributing

### Development Setup
1. Fork the repository
2. Run `./setup-dev.sh`
3. Start development with `make dev-start`
4. Make changes with hot reload
5. Run tests with `make test`
6. Submit pull request

### Code Standards
- Go: Follow `go fmt` and `golint`
- TypeScript: ESLint and Prettier configured
- Git: Conventional commit messages

## 📞 Support

### Getting Help
- Check logs: `./maintenance.sh logs <service>`
- Run health checks: `./maintenance.sh health`
- Monitor system: `./maintenance.sh monitor`
- Check documentation in `docs/`

### Reporting Issues
Include in your report:
- Operating system and version
- Docker and Docker Compose versions
- Go and Node.js versions
- Complete error logs
- Steps to reproduce

---

## 📄 License

This project is licensed under the MIT License - see the LICENSE file for details.

---

**🎉 Happy Coding! The CI/CD Status Notifier Bot is ready for deployment and development!**
