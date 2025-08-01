# 🚀 CICD Pipeline Testing Instructions

## ✅ Pipeline Fixed - Deployment ke Server Remote (172.16.19.11)

Masalah deployment localhost sudah diperbaiki! Sekarang semua deployment berjalan di server remote **172.16.19.11**.

## 🔧 Yang Sudah Diperbaiki:

### 1. **Remote Deployment Configuration**
- ✅ Semua deployment sekarang ke server 172.16.19.11
- ✅ SSH-based deployment menggunakan `appleboy/ssh-action`
- ✅ Health checks berjalan di server remote
- ✅ External accessibility validation

### 2. **Multi-Environment Setup**
- 🟢 **Staging**: Port 8082/3002 (branch `staging`)
- 🟡 **Test**: Port 8081/3001 (branch `develop`) 
- 🔴 **Production**: Port 80/8080 (branch `main`)

### 3. **Comprehensive Testing**
- ✅ Manual staging test script
- ✅ Automated GitHub Actions testing
- ✅ Pre-deployment validation
- ✅ Post-deployment health checks

## 🧪 Cara Test Pipeline:

### **Option 1: Manual Test (Direkomendasikan)**
```bash
# Pastikan di branch staging
git checkout staging

# Jalankan manual test
./scripts/manual-staging-test.sh
```

### **Option 2: GitHub Actions Automated**
```bash
# Push ke staging untuk trigger automated test
git push origin staging
```

### **Option 3: Test Individual Components**
```bash
# Test pipeline components
./scripts/test-pipeline.sh
```

## 📋 Prerequisites Setup:

### 1. **Server Setup** (Jalankan di server 172.16.19.11):
```bash
wget https://raw.githubusercontent.com/dewisartika8/CICD-Status-Notifier-Bot/staging/scripts/server-setup.sh
chmod +x server-setup.sh
sudo ./server-setup.sh
```

### 2. **GitHub Secrets** (Tambahkan di repository settings):
```
# Server Access
SSH_USERNAME=your_server_username
SSH_PRIVATE_KEY=your_private_ssh_key
SSH_PORT=22

# Application Config
POSTGRES_PASSWORD=your_secure_password
POSTGRES_PASSWORD_STAGING=your_staging_password
JWT_SECRET=your_jwt_secret
JWT_SECRET_STAGING=your_staging_jwt_secret
TELEGRAM_BOT_TOKEN=your_telegram_bot_token

# External Services
SONAR_TOKEN=your_sonarqube_token
SONAR_HOST_URL=https://your-sonarqube-server.com
SLACK_WEBHOOK_URL=https://hooks.slack.com/services/your/webhook/url
SNYK_TOKEN=your_snyk_token
```

## 🔍 Test Results Access:

### **Staging Environment**
- **Frontend**: http://172.16.19.11:3002
- **Backend**: http://172.16.19.11:8082
- **Health**: http://172.16.19.11:8082/health

### **Test Environment**  
- **Frontend**: http://172.16.19.11:3001
- **Backend**: http://172.16.19.11:8081
- **Health**: http://172.16.19.11:8081/health

### **Production Environment**
- **Frontend**: http://172.16.19.11
- **Backend**: http://172.16.19.11:8080
- **Health**: http://172.16.19.11:8080/health

## 🎯 Testing Workflow:

1. **Development** → Push ke `develop` → Test Environment
2. **Staging** → Push ke `staging` → Staging Environment 
3. **Production** → Push ke `main` → Production Environment

## 📊 Monitoring & Alerts:

- **Health Checks**: Setiap 15 menit
- **Slack Notifications**: Real-time deployment status
- **PR Comments**: Automated test result updates
- **GitHub Actions**: Comprehensive workflow monitoring

## 🚨 Jika Ada Masalah:

### Debug Commands:
```bash
# Check container status
ssh 172.16.19.11 "docker ps"

# View logs
ssh 172.16.19.11 "docker logs cicd_backend_staging"

# Check server resources
ssh 172.16.19.11 "df -h && free -h"

# Test connectivity
curl -v http://172.16.19.11:8082/health
```

### Rollback:
```bash
ssh 172.16.19.11 "cd /opt/cicd-notifier-staging && docker-compose down"
```

## 📚 Dokumentasi Lengkap:

- **[CICD_PIPELINE_GUIDE.md](docs/CICD_PIPELINE_GUIDE.md)** - Dokumentasi lengkap
- **[TESTING_GUIDE.md](TESTING_GUIDE.md)** - Panduan testing detail
- **[CICD_SETUP.md](CICD_SETUP.md)** - Quick start guide

---

## 🎉 Ready to Test!

Pipeline sudah siap untuk testing dengan deployment ke server remote. Tidak ada lagi deployment ke localhost!

**Next Steps:**
1. Setup server prerequisites
2. Configure GitHub secrets  
3. Run manual test: `./scripts/manual-staging-test.sh`
4. Push to staging: `git push origin staging`
5. Monitor GitHub Actions workflow
6. Access staging environment: http://172.16.19.11:3002

**Happy Testing! 🚀**
