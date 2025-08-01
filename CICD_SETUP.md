# CICD GitHub Actions Quick Start

## ✨ Overview

Komprehensif GitHub Actions CI/CD pipeline untuk CICD Status Notifier Bot dengan deployment otomatis ke server IP `172.16.19.11`.

## 🚀 Quick Start

### 1. Setup Server

Jalankan script setup di server target:

```bash
# SSH ke server 172.16.19.11
ssh user@172.16.19.11

# Download dan jalankan setup script
wget https://raw.githubusercontent.com/dewisartika8/CICD-Status-Notifier-Bot/main/scripts/server-setup.sh
chmod +x server-setup.sh
sudo ./server-setup.sh
```

### 2. Configure GitHub Secrets

Tambahkan secrets berikut di GitHub repository settings:

#### 🔐 Server Access
```
SSH_USERNAME=your_server_username
SSH_PRIVATE_KEY=your_private_ssh_key
SSH_PORT=22
```

#### 🗄️ Database & App
```
POSTGRES_PASSWORD=your_secure_password
POSTGRES_PASSWORD_STAGING=your_staging_password
JWT_SECRET=your_jwt_secret_key
JWT_SECRET_STAGING=your_staging_jwt_secret
TELEGRAM_BOT_TOKEN=your_telegram_bot_token
```

#### 🔍 Code Quality & Monitoring
```
SONAR_TOKEN=your_sonarqube_token
SONAR_HOST_URL=https://your-sonarqube-server.com
SLACK_WEBHOOK_URL=https://hooks.slack.com/services/your/webhook/url
SNYK_TOKEN=your_snyk_token
```

### 3. Test Locally (Optional)

```bash
# Test staging deployment manually
./scripts/manual-staging-test.sh

# Test general pipeline
./scripts/test-pipeline.sh
```

### 4. Deploy

```bash
# Push ke staging untuk staging environment
git push origin staging

# Push ke develop untuk test environment
git push origin develop

# Push ke main untuk production
git push origin main
```

## 📊 Pipeline Features

### ✅ Automated Testing
- **Backend**: Go unit tests, integration tests, static analysis
- **Frontend**: React tests, ESLint, build validation
- **Coverage**: Codecov integration dengan detailed reports

### 🔒 Security Scanning
- **Trivy**: Container vulnerability scanning
- **Snyk**: Dependency vulnerability checks
- **gosec**: Go security analyzer
- **SonarQube**: Code quality dan security analysis

### 🐳 Containerization
- **Multi-stage builds** untuk optimized images
- **GitHub Container Registry** untuk image storage
- **Layer caching** untuk faster builds

### 🚀 Deployment
- **Staging Environment**: Deployment dari `staging` branch ke port 8082/3002
- **Test Environment**: Automatic deployment dari `develop` branch
- **Production**: Zero-downtime deployment dari `main` branch
- **Health checks** sebelum dan sesudah deployment
- **Rollback capability** jika deployment gagal

### 📈 Monitoring
- **Health checks** setiap 15 menit
- **Performance monitoring** 
- **Disk space monitoring**
- **SSL certificate expiry alerts**
- **Slack notifications** untuk semua events

## 🌐 Access Points

### Staging Environment
- **Frontend**: http://localhost:3002
- **Backend API**: http://localhost:8082
- **Health Check**: http://localhost:8082/health

### Test Environment
- **Frontend**: http://localhost:3001
- **Backend API**: http://localhost:8081
- **Health Check**: http://localhost:8081/health

### Production Environment  
- **Frontend**: http://localhost
- **Backend API**: http://localhost:8080
- **Health Check**: http://localhost:8080/health

## 📋 Workflow Triggers

### Main CI/CD Pipeline
- ✅ Push ke `main` branch → Production deployment
- ✅ Push ke `staging` branch → Staging deployment
- ✅ Push ke `develop` branch → Test deployment  
- ✅ Pull requests → Testing only

### Remote Deployment Test
- ✅ Push ke `staging` branch → Comprehensive staging tests
- ✅ Manual trigger available
- ✅ PR comments with test results

### SonarQube Analysis
- ✅ Push ke `main`/`staging`/`develop` branches
- ✅ Pull requests ke `main`/`staging`

### Health Monitoring
- ✅ Scheduled: Every 15 minutes
- ✅ Manual trigger available

## 🛠️ Local Development

### Prerequisites
- Docker & Docker Compose
- Go 1.23+
- Node.js 18+
- Git

### Development Workflow
```bash
# 1. Clone repository
git clone https://github.com/dewisartika8/CICD-Status-Notifier-Bot.git
cd CICD-Status-Notifier-Bot

# 2. Start local development
docker-compose up -d postgres

# 3. Run backend
cd backend
go run cmd/main.go

# 4. Run frontend (new terminal)
cd frontend
npm start

# 5. Test changes
./scripts/test-pipeline.sh
```

## 🔧 Troubleshooting

### Common Issues

**❌ Deployment Failed**
```bash
# Check logs
ssh 172.16.19.11 "docker logs cicd_backend"
ssh 172.16.19.11 "docker logs cicd_frontend"

# Restart services
ssh 172.16.19.11 "cd /opt/cicd-notifier && docker-compose restart"
```

**❌ Tests Failing**
```bash
# Run tests locally
./scripts/test-pipeline.sh

# Check specific test
cd backend && go test -v ./...
cd frontend && npm test
```

**❌ Server Issues**
```bash
# Check server status
ssh 172.16.19.11 "systemctl status docker nginx"

# Check disk space
ssh 172.16.19.11 "df -h"

# Check logs
ssh 172.16.19.11 "tail -f /var/log/cicd-notifier/health-check.log"
```

## 📚 Documentation

- **[Complete Pipeline Guide](docs/CICD_PIPELINE_GUIDE.md)** - Detailed documentation
- **[API Documentation](docs/API_DOCUMENTATION_PROJECTS.md)** - API endpoints
- **[Technical Design](docs/TECHNICAL_DESIGN.md)** - Architecture overview

## 🔄 Backup & Recovery

### Automated Backups
- **Daily backups** at 2:00 AM
- **7-day retention** policy
- **Database + application data**

### Manual Backup
```bash
ssh 172.16.19.11 "/opt/cicd-notifier/backup.sh"
```

### Recovery
```bash
# Restore from backup
ssh 172.16.19.11 "cd /opt/cicd-notifier && docker exec cicd_postgres psql -U postgres -d cicd_notifier < backups/database_YYYYMMDD_HHMMSS.sql"
```

## 🚨 Monitoring & Alerts

### Health Checks
- **Endpoint availability**: Frontend, Backend, Database
- **Response time monitoring**
- **SSL certificate expiry**
- **Disk space usage**

### Alert Channels
- **Slack**: Real-time notifications
- **GitHub**: Issues untuk persistent problems
- **Email**: Critical alerts (configure via Slack)

## 📈 Performance Metrics

### Build Performance
- **Average build time**: ~8-12 minutes
- **Test execution**: ~3-5 minutes
- **Image build**: ~2-3 minutes
- **Deployment**: ~2-3 minutes

### Runtime Performance
- **API response time**: <200ms average
- **Frontend load time**: <2s
- **Database queries**: <100ms average

## 🎯 Next Steps

1. **SSL Configuration**: Setup Let's Encrypt certificates
2. **Domain Setup**: Configure custom domain
3. **Advanced Monitoring**: Prometheus + Grafana
4. **Scaling**: Kubernetes migration plan
5. **Multi-region**: Additional deployment targets

## 🤝 Contributing

1. Fork repository
2. Create feature branch
3. Run tests locally
4. Submit pull request
5. Monitor CI pipeline

## 📞 Support

- **Issues**: GitHub Issues
- **Documentation**: `/docs` folder
- **Team Contact**: Create GitHub issue dengan label `support`

---

**🎉 Happy Deploying!**

Pipeline ini dirancang untuk memberikan development experience yang smooth dengan deployment yang reliable ke server 172.16.19.11.
