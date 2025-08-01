# GitHub Actions Workflow Update

## 📋 Overview

The CI/CD pipeline has been restructured to better handle local testing and different deployment environments. The previous SSH-based localhost deployment approach has been replaced with a more practical Docker-based testing strategy.

## 🔄 Workflow Changes

### 1. Simple CI Testing (`simple-ci.yml`)
- **Triggers**: Push to `staging`, `develop` and PRs to `main`, `staging`, `develop`
- **Purpose**: Basic build verification and testing
- **Components**:
  - Backend compilation check
  - Frontend build verification  
  - Docker image build testing
  - Simple Docker Compose validation

### 2. Local Testing (`local-ci-cd.yml`)
- **Triggers**: Push to `main`, `develop`, `staging` and PRs
- **Purpose**: Comprehensive testing with database integration
- **Features**:
  - Full backend/frontend testing
  - Docker image building and pushing to GHCR
  - Integration testing with PostgreSQL
  - Health checks and API validation

### 3. Production Pipeline (`ci-cd.yml`)
- **Triggers**: Push to `main`, release tags
- **Purpose**: Production deployment (for actual server deployment)
- **Note**: Modified to only run for production releases

## 🧪 Local Testing Strategy

### Option 1: GitHub Actions Simulation
Use the provided script to simulate GitHub Actions locally:

```bash
# Make executable
chmod +x scripts/local-test.sh

# Run local testing
./scripts/local-test.sh
```

### Option 2: Manual Docker Testing
```bash
# Build and start services
docker-compose -f docker-compose.test.yml up --build -d

# Check services
docker-compose -f docker-compose.test.yml ps

# View logs
docker-compose -f docker-compose.test.yml logs -f

# Cleanup
docker-compose -f docker-compose.test.yml down --volumes
```

## 🎯 Testing URLs

When running locally:
- **Backend Health**: http://localhost:8082/health
- **Backend API**: http://localhost:8082/api/v1/status  
- **Frontend**: http://localhost:3002
- **Database**: localhost:5434

## 🔐 Required Secrets

For GitHub Actions to work properly, add these secrets to your repository:

```
POSTGRES_PASSWORD=your_secure_password
```

Go to: Repository Settings → Secrets and variables → Actions → New repository secret

## 📊 Workflow Results

### Simple CI (`simple-ci.yml`)
- ✅ Quick validation of code compilation
- ✅ Basic build checks
- ✅ Docker image creation verification

### Local CI/CD (`local-ci-cd.yml`) 
- ✅ Comprehensive testing suite
- ✅ Database integration testing
- ✅ Full deployment simulation
- ✅ Health check validation

## 🚀 Branch Strategy

| Branch | Workflow | Purpose |
|--------|----------|---------|
| `develop` | simple-ci.yml | Development validation |
| `staging` | simple-ci.yml + local-ci-cd.yml | Pre-production testing |
| `main` | All workflows | Production ready |

## 🔧 Troubleshooting

### Common Issues

1. **Docker not running**
   ```bash
   # Start Docker Desktop or Docker daemon
   sudo systemctl start docker  # Linux
   # or open Docker Desktop      # macOS/Windows
   ```

2. **Port conflicts**
   ```bash
   # Check what's using the ports
   lsof -i :8082  # Backend
   lsof -i :3002  # Frontend  
   lsof -i :5434  # Database
   ```

3. **Container build failures**
   ```bash
   # Clean Docker cache
   docker system prune -a
   docker volume prune
   ```

### Viewing Logs

```bash
# Backend logs
docker-compose -f docker-compose.test.yml logs backend-test

# Frontend logs  
docker-compose -f docker-compose.test.yml logs frontend-test

# Database logs
docker-compose -f docker-compose.test.yml logs postgres-test
```

## 📈 Monitoring

The workflows generate:
- Test coverage reports
- Build artifacts
- Integration test results
- Health check status
- Performance metrics

## 🔄 Next Steps

1. **Test the new workflow**:
   ```bash
   git add .
   git commit -m "feat: Updated CI/CD workflows for local testing"
   git push origin staging
   ```

2. **Monitor GitHub Actions**: Check the Actions tab in your repository

3. **Run local tests**: Use `./scripts/local-test.sh` for local validation

4. **Review results**: Check generated reports and logs

## 💡 Best Practices

- Always test locally before pushing
- Use the staging branch for pre-production testing
- Monitor resource usage during testing
- Keep Docker images optimized
- Regular cleanup of test artifacts

---

**Generated**: $(date)
**Version**: 2.0
**Status**: Active
