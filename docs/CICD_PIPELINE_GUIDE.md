# GitHub Actions CI/CD Pipeline Documentation

## Overview

This repository contains a comprehensive CI/CD pipeline using GitHub Actions for the CICD Status Notifier Bot application. The pipeline includes testing, building, security scanning, SonarQube analysis, and deployment to server IP `172.16.19.11`.

## Pipeline Components

### 1. Main CI/CD Workflow (`ci-cd.yml`)

#### Triggers
- Push to `main`, `develop`, and `staging` branches
- Pull requests to `main`, `develop`, and `staging` branches

#### Jobs

**Backend Test Job**
- Runs Go tests with PostgreSQL service
- Performs static analysis (go vet, staticcheck, gosec)
- Generates test coverage reports
- Runs integration tests
- Uploads coverage to Codecov

**Frontend Test Job**
- Runs Node.js tests with coverage
- Performs ESLint and Prettier checks
- Builds production frontend
- Uploads build artifacts

**Security Scan Job**
- Trivy vulnerability scanner
- Snyk security checks
- Uploads results to GitHub Security tab

**Build Job**
- Builds Docker images for backend and frontend
- Pushes images to GitHub Container Registry
- Uses caching for optimization

**Deploy Staging Job**
- Deploys to staging environment on `staging` branch
- Runs on port 8082 (backend) and 3002 (frontend)  
- Uses staging-specific configuration
- Comprehensive health checks

**Deploy Test Job**
- Deploys to test environment on `develop` branch or PRs
- Runs on port 8081 (backend) and 3001 (frontend)
- Performs health checks

**Deploy Production Job**
- Deploys to production on `main` branch
- Zero-downtime deployment strategy
- Runs database migrations
- Comprehensive health checks

**Smoke Tests Job**
- End-to-end functionality tests
- API endpoint validation
- Frontend accessibility checks

**Notification Job**
- Slack notifications for deployment status
- Always runs regardless of success/failure

### 2. SonarQube Analysis (`sonarqube.yml`)

- Dedicated workflow for code quality analysis
- Separate coverage reports for Go and JavaScript
- Quality gate enforcement
- Detailed code metrics and security analysis

### 4. Remote Deployment Test (`staging-deployment-test.yml`)

- Comprehensive testing workflow for staging environment
- Pre-deployment server validation
- Automated deployment testing on remote server (172.16.19.11)
- Post-deployment health checks and performance tests
- External accessibility validation
- Automated notifications and PR comments

### 5. Health Check & Monitoring (`monitoring.yml`)

- Scheduled health checks every 15 minutes
- Performance monitoring
- Disk space monitoring
- Security checks
- SSL certificate expiry monitoring
- Automated alerts via Slack

## Server Setup

### Prerequisites

1. Ubuntu 20.04+ server at IP `172.16.19.11`
2. SSH access with sudo privileges
3. Domain name (optional, for SSL)

### Automated Setup

Run the server setup script:

```bash
# On the target server
wget https://raw.githubusercontent.com/your-username/CICD-Status-Notifier-Bot/main/scripts/server-setup.sh
chmod +x server-setup.sh
./server-setup.sh
```

### Manual Setup Steps

1. **Install Dependencies**
   ```bash
   sudo apt update && sudo apt upgrade -y
   sudo apt install -y docker.io docker-compose nginx ufw fail2ban
   ```

2. **Configure Firewall**
   ```bash
   sudo ufw enable
   sudo ufw allow ssh
   sudo ufw allow 80/tcp
   sudo ufw allow 443/tcp
   sudo ufw allow 8080/tcp
   ```

3. **Create Application Directories**
   ```bash
   sudo mkdir -p /opt/cicd-notifier
   sudo mkdir -p /opt/cicd-notifier-test
   sudo chown -R $USER:$USER /opt/cicd-notifier*
   ```

## GitHub Repository Configuration

### Required Secrets

Configure the following secrets in your GitHub repository settings:

#### Server Access
- `SSH_USERNAME`: Server username for deployment
- `SSH_PRIVATE_KEY`: Private SSH key for server access
- `SSH_PORT`: SSH port (default: 22)

#### Application Configuration
- `POSTGRES_PASSWORD`: PostgreSQL password for production
- `POSTGRES_PASSWORD_STAGING`: PostgreSQL password for staging
- `JWT_SECRET`: JWT secret key for authentication
- `JWT_SECRET_STAGING`: JWT secret key for staging
- `TELEGRAM_BOT_TOKEN`: Telegram bot token for notifications

#### External Services
- `SONAR_TOKEN`: SonarQube authentication token
- `SONAR_HOST_URL`: SonarQube server URL
- `SLACK_WEBHOOK_URL`: Slack webhook for notifications
- `SNYK_TOKEN`: Snyk authentication token

### Branch Protection Rules

Configure branch protection for `main` branch:
- Require status checks to pass
- Require up-to-date branches
- Include administrators
- Restrict pushes to specific people/teams

## Deployment Environments

### Staging Environment
- **URL**: `http://172.16.19.11:3002`
- **API**: `http://172.16.19.11:8082`
- **Database**: PostgreSQL on port 5434
- **Triggers**: Push to `staging` branch

### Test Environment
- **URL**: `http://172.16.19.11:3001`
- **API**: `http://172.16.19.11:8081`
- **Database**: PostgreSQL on port 5433
- **Triggers**: Push to `develop` or pull requests

### Production Environment
- **URL**: `http://172.16.19.11`
- **API**: `http://172.16.19.11:8080`
- **Database**: PostgreSQL on port 5432
- **Triggers**: Push to `main` branch

## Monitoring and Observability

### Health Checks
- **Backend**: `/health` endpoint
- **Database**: Connection test via API
- **Frontend**: Root endpoint availability

### Metrics Collection
- Node Exporter on port 9100
- Docker metrics
- Application metrics via API

### Log Management
- Application logs in `/var/log/cicd-notifier/`
- Log rotation configured
- Centralized logging with structured format

### Backup Strategy
- Daily database backups at 2:00 AM
- Application data backups
- 7-day retention policy
- Automated backup verification

## Security Measures

### Infrastructure Security
- UFW firewall configuration
- Fail2ban for SSH protection
- Regular security updates
- SSL/TLS encryption (configurable)

### Application Security
- Dependency vulnerability scanning
- Static code analysis
- Container image scanning
- Security-focused linting rules

### Access Control
- SSH key-based authentication
- Limited sudo access
- Container isolation
- Network segmentation

## Troubleshooting

### Common Issues

**Deployment Failures**
```bash
# Check container status
docker ps -a

# View container logs
docker logs cicd_backend
docker logs cicd_frontend
docker logs cicd_postgres

# Restart services
docker-compose down && docker-compose up -d
```

**Database Issues**
```bash
# Check database connectivity
docker exec cicd_postgres pg_isready -U postgres

# Run migrations manually
docker exec cicd_backend /app/migrate up

# Backup database
docker exec cicd_postgres pg_dump -U postgres cicd_notifier > backup.sql
```

**Performance Issues**
```bash
# Monitor resource usage
htop
docker stats

# Check disk space
df -h
docker system df

# Clean up Docker resources
docker system prune -f
```

### Log Locations

- **Application Logs**: `/var/log/cicd-notifier/`
- **Nginx Logs**: `/var/log/nginx/`
- **Docker Logs**: `docker logs <container_name>`
- **System Logs**: `/var/log/syslog`

## Development Workflow

### Development Workflow

### Feature Development
1. Create feature branch from `develop`
2. Implement changes with tests
3. Push to trigger CI pipeline
4. Create pull request to `develop`
5. Review and merge after CI passes

### Staging Testing
1. Create staging branch from `develop`
2. Push to `staging` branch
3. Automated deployment to staging environment
4. Manual testing on http://172.16.19.11:3002
5. Validate all functionality

### Release Process
1. Merge `staging` to `main` after validation
2. Tag release version
3. Automatic production deployment
4. Monitor deployment health
5. Rollback if necessary

### Hotfix Process
1. Create hotfix branch from `main`
2. Implement fix with tests
3. Test on staging first
4. Create PR to `main`
5. Emergency deployment after approval

## Performance Optimization

### Build Optimization
- Multi-stage Docker builds
- Layer caching strategies
- Parallel job execution
- Artifact caching

### Runtime Optimization
- Container resource limits
- Database connection pooling
- Static asset optimization
- CDN integration (future)

## Future Enhancements

### Planned Features
- Kubernetes deployment option
- Multi-environment support
- Blue-green deployment
- Canary releases
- Advanced monitoring dashboard
- Automated scaling

### Integration Opportunities
- Prometheus + Grafana monitoring
- ELK stack for logging
- ArgoCD for GitOps
- Terraform for infrastructure
- Istio service mesh

## Support and Maintenance

### Regular Maintenance Tasks
- Weekly security updates
- Monthly dependency updates
- Quarterly architecture reviews
- Bi-annual disaster recovery tests

### Contact Information
- **Team Lead**: [Your Name]
- **DevOps**: [DevOps Team]
- **On-call**: [Emergency Contact]

For questions or issues, please create a GitHub issue or contact the development team.
