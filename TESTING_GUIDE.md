# Testing Guide for CICD Pipeline

## 🧪 Testing Strategy

Proyek ini menggunakan strategi testing multi-level untuk memastikan deployment yang reliable ke server 172.16.19.11.

## 🔄 Testing Environments

### 1. Staging Environment (Port 8082/3002)
- **Purpose**: Testing comprehensive sebelum production
- **Trigger**: Push ke branch `staging`
- **URL**: http://localhost:3002
- **API**: http://localhost:8082

### 2. Test Environment (Port 8081/3001)  
- **Purpose**: Development testing dan PR validation
- **Trigger**: Push ke branch `develop` atau pull requests
- **URL**: http://localhost:3001
- **API**: http://localhost:8081

### 3. Production Environment (Port 80/8080)
- **Purpose**: Live production application
- **Trigger**: Push ke branch `main`
- **URL**: http://localhost
- **API**: http://localhost:8080

## 🚀 How to Test

### Method 1: Manual Staging Test (Recommended)

```bash
# Pastikan berada di branch staging
git checkout staging

# Jalankan manual testing script
./scripts/manual-staging-test.sh
```

Script ini akan:
1. ✅ Check prerequisites (SSH, Docker, dll)
2. ✅ Test SSH connectivity ke server
3. ✅ Validate server environment
4. ✅ Build dan push Docker images
5. ✅ Deploy ke staging environment
6. ✅ Run comprehensive health checks
7. ✅ Test external accessibility
8. ✅ Generate detailed report

### Method 2: GitHub Actions Automated Test

```bash
# Push ke staging branch untuk trigger automated test
git add .
git commit -m "Test staging deployment"
git push origin staging
```

Automated workflow akan:
1. ✅ Pre-deployment validation
2. ✅ Build dan deploy ke remote server
3. ✅ Post-deployment health checks
4. ✅ Performance baseline tests
5. ✅ Notification ke Slack dan PR comments

### Method 3: Local Pipeline Test

```bash
# Test pipeline components locally
./scripts/test-pipeline.sh
```

## 📊 Test Results Interpretation

### ✅ Success Indicators
- All containers running (`docker ps`)
- Health endpoints responding (200 status)
- External accessibility confirmed
- Response times < 2s for backend, < 3s for frontend
- Database connectivity verified

### ❌ Failure Indicators
- Container status `Exited` atau `Restarting`
- Health check timeouts
- External access failures
- High response times
- Database connection errors

## 🔍 Debugging Failed Tests

### 1. Check Container Status
```bash
ssh localhost "cd /opt/cicd-notifier-staging && docker-compose -f docker-compose.staging.yml ps"
```

### 2. View Container Logs
```bash
# Backend logs
ssh localhost "docker logs cicd_backend_staging --tail=50"

# Frontend logs  
ssh localhost "docker logs cicd_frontend_staging --tail=50"

# Database logs
ssh localhost "docker logs cicd_postgres_staging --tail=50"
```

### 3. Check Server Resources
```bash
ssh localhost "df -h && free -h && docker system df"
```

### 4. Test Network Connectivity
```bash
# From outside server
curl -v http://localhost:8082/health
curl -v http://localhost:3002

# From inside server
ssh localhost "curl -v http://localhost:8082/health"
```

## 🛠️ Common Issues & Solutions

### Issue 1: SSH Connection Failed
**Symptoms**: Cannot connect to server
**Solutions**:
```bash
# Check SSH key permissions
chmod 600 ~/.ssh/id_rsa

# Test SSH connection
ssh -v localhost

# Verify server is accessible
ping localhost
```

### Issue 2: Docker Image Pull Failed
**Symptoms**: Cannot pull images from registry
**Solutions**:
```bash
# Login to GitHub Container Registry
echo $GITHUB_TOKEN | docker login ghcr.io -u $GITHUB_USERNAME --password-stdin

# Check image exists
docker pull ghcr.io/dewisartika8/cicd-status-notifier-bot-backend:staging-test
```

### Issue 3: Port Already in Use
**Symptoms**: Deployment fails with port binding error
**Solutions**:
```bash
# Stop existing containers
ssh localhost "cd /opt/cicd-notifier-staging && docker-compose -f docker-compose.staging.yml down"

# Check what's using the port
ssh localhost "netstat -tulpn | grep :8082"
```

### Issue 4: Health Check Timeout
**Symptoms**: Services not responding to health checks
**Solutions**:
```bash
# Increase wait time in deployment script
# Check application logs for startup errors
# Verify environment variables are set correctly
```

## 📈 Performance Benchmarks

### Expected Response Times
- **Backend Health Check**: < 500ms
- **Frontend Load**: < 2s
- **API Endpoints**: < 1s
- **Database Queries**: < 100ms

### Resource Usage Limits
- **CPU**: < 80% sustained usage
- **Memory**: < 85% total RAM
- **Disk**: < 80% total space
- **Network**: < 100MB/s sustained

## 🔄 Test Automation Schedule

### Continuous Integration
- **Every Push**: Unit tests, linting, security scans
- **Every PR**: Full test suite + staging deployment
- **Every Merge**: Production deployment + smoke tests

### Scheduled Testing
- **Every 15 min**: Health checks via monitoring workflow
- **Daily**: Backup validation
- **Weekly**: Security updates check
- **Monthly**: Performance benchmarking

## 📝 Test Reporting

### Automated Reports
- **GitHub Actions**: Workflow summaries
- **Slack**: Real-time notifications
- **PR Comments**: Deployment status updates

### Manual Reports
- **staging-test-report.md**: Generated by manual test script
- **Container logs**: Available via SSH
- **System metrics**: Accessible through monitoring

## 🎯 Best Practices

### Before Testing
1. ✅ Ensure all required secrets are configured
2. ✅ Verify server accessibility and resources
3. ✅ Check for conflicting processes on target ports
4. ✅ Review recent changes and dependencies

### During Testing
1. ✅ Monitor logs in real-time
2. ✅ Verify each stage completes successfully
3. ✅ Test all critical user journeys
4. ✅ Validate external integrations

### After Testing
1. ✅ Review all logs and metrics
2. ✅ Document any issues or improvements
3. ✅ Clean up test artifacts if needed
4. ✅ Update documentation if processes change

## 🚨 Emergency Procedures

### Rollback Staging Deployment
```bash
ssh localhost "cd /opt/cicd-notifier-staging && docker-compose -f docker-compose.staging.yml down"
```

### Emergency Access to Server
```bash
# Direct SSH access
ssh localhost

# Check all running services
docker ps -a

# Emergency stop all services
docker stop $(docker ps -q)
```

### Contact Information
- **Primary**: Create GitHub issue with `urgent` label
- **Secondary**: Check team communication channels
- **Emergency**: Server admin contact (if available)

---

## 🎉 Success Criteria

Your staging deployment is successful when:

1. ✅ Manual test script completes without errors
2. ✅ All health checks pass consistently
3. ✅ External access works from multiple locations
4. ✅ Performance metrics meet benchmarks
5. ✅ No critical errors in logs
6. ✅ Database connectivity confirmed
7. ✅ All containers running stably

Ready to proceed with GitHub Actions automated deployment! 🚀
