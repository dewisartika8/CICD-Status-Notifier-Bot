#!/bin/bash

# Server Setup Script for CICD Status Notifier Bot
# Target: 172.16.19.11

set -e

echo "🚀 Setting up server for CICD Status Notifier Bot..."

# Update system
echo "📦 Updating system packages..."
sudo apt update && sudo apt upgrade -y

# Install required packages
echo "📥 Installing required packages..."
sudo apt install -y \
    apt-transport-https \
    ca-certificates \
    curl \
    gnupg \
    lsb-release \
    software-properties-common \
    unzip \
    wget \
    htop \
    git \
    ufw \
    fail2ban

# Install Docker
echo "🐳 Installing Docker..."
if ! command -v docker &> /dev/null; then
    curl -fsSL https://download.docker.com/linux/ubuntu/gpg | sudo gpg --dearmor -o /usr/share/keyrings/docker-archive-keyring.gpg
    echo "deb [arch=$(dpkg --print-architecture) signed-by=/usr/share/keyrings/docker-archive-keyring.gpg] https://download.docker.com/linux/ubuntu $(lsb_release -cs) stable" | sudo tee /etc/apt/sources.list.d/docker.list > /dev/null
    sudo apt update
    sudo apt install -y docker-ce docker-ce-cli containerd.io docker-buildx-plugin docker-compose-plugin
    
    # Add current user to docker group
    sudo usermod -aG docker $USER
    
    # Enable Docker service
    sudo systemctl enable docker
    sudo systemctl start docker
fi

# Install Docker Compose (standalone)
echo "📦 Installing Docker Compose..."
if ! command -v docker-compose &> /dev/null; then
    sudo curl -L "https://github.com/docker/compose/releases/latest/download/docker-compose-$(uname -s)-$(uname -m)" -o /usr/local/bin/docker-compose
    sudo chmod +x /usr/local/bin/docker-compose
fi

# Configure firewall
echo "🔥 Configuring firewall..."
sudo ufw --force enable
sudo ufw default deny incoming
sudo ufw default allow outgoing

# Allow SSH
sudo ufw allow ssh

# Allow HTTP and HTTPS
sudo ufw allow 80/tcp
sudo ufw allow 443/tcp

# Allow application ports
sudo ufw allow 8080/tcp  # Backend
sudo ufw allow 3000/tcp  # Frontend dev
sudo ufw allow 5432/tcp  # PostgreSQL (if needed externally)

# Configure fail2ban
echo "🛡️ Configuring fail2ban..."
sudo systemctl enable fail2ban
sudo systemctl start fail2ban

# Create jail for SSH
sudo tee /etc/fail2ban/jail.local << EOF
[DEFAULT]
bantime = 3600
findtime = 600
maxretry = 3

[sshd]
enabled = true
port = ssh
filter = sshd
logpath = /var/log/auth.log
maxretry = 3
bantime = 3600
EOF

sudo systemctl restart fail2ban

# Create application directories
echo "📁 Creating application directories..."
sudo mkdir -p /opt/cicd-notifier
sudo mkdir -p /opt/cicd-notifier-test
sudo mkdir -p /opt/cicd-notifier/backups
sudo mkdir -p /opt/cicd-notifier/ssl
sudo mkdir -p /var/log/cicd-notifier

# Set proper permissions
sudo chown -R $USER:$USER /opt/cicd-notifier
sudo chown -R $USER:$USER /opt/cicd-notifier-test
sudo chown -R $USER:$USER /var/log/cicd-notifier

# Install Nginx (for SSL termination and reverse proxy)
echo "🌐 Installing and configuring Nginx..."
sudo apt install -y nginx
sudo systemctl enable nginx

# Create Nginx configuration for the application
sudo tee /etc/nginx/sites-available/cicd-notifier << 'EOF'
server {
    listen 80;
    server_name 172.16.19.11;

    # Redirect HTTP to HTTPS (uncomment when SSL is configured)
    # return 301 https://$server_name$request_uri;

    # For now, serve directly
    location / {
        proxy_pass http://localhost:3000;
        proxy_http_version 1.1;
        proxy_set_header Upgrade $http_upgrade;
        proxy_set_header Connection 'upgrade';
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
        proxy_cache_bypass $http_upgrade;
    }

    location /api/ {
        proxy_pass http://localhost:8080;
        proxy_http_version 1.1;
        proxy_set_header Upgrade $http_upgrade;
        proxy_set_header Connection 'upgrade';
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
        proxy_cache_bypass $http_upgrade;
    }

    # Health check endpoint
    location /health {
        proxy_pass http://localhost:8080/health;
        access_log off;
    }
}

# HTTPS configuration (uncomment when SSL certificates are available)
# server {
#     listen 443 ssl http2;
#     server_name 172.16.19.11;
#
#     ssl_certificate /opt/cicd-notifier/ssl/cert.pem;
#     ssl_certificate_key /opt/cicd-notifier/ssl/key.pem;
#     ssl_protocols TLSv1.2 TLSv1.3;
#     ssl_ciphers ECDHE-RSA-AES256-GCM-SHA512:DHE-RSA-AES256-GCM-SHA512:ECDHE-RSA-AES256-GCM-SHA384:DHE-RSA-AES256-GCM-SHA384;
#     ssl_prefer_server_ciphers off;
#
#     location / {
#         proxy_pass http://localhost:3000;
#         proxy_http_version 1.1;
#         proxy_set_header Upgrade $http_upgrade;
#         proxy_set_header Connection 'upgrade';
#         proxy_set_header Host $host;
#         proxy_set_header X-Real-IP $remote_addr;
#         proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
#         proxy_set_header X-Forwarded-Proto $scheme;
#         proxy_cache_bypass $http_upgrade;
#     }
#
#     location /api/ {
#         proxy_pass http://localhost:8080;
#         proxy_http_version 1.1;
#         proxy_set_header Host $host;
#         proxy_set_header X-Real-IP $remote_addr;
#         proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
#         proxy_set_header X-Forwarded-Proto $scheme;
#     }
# }
EOF

# Enable the site
sudo ln -sf /etc/nginx/sites-available/cicd-notifier /etc/nginx/sites-enabled/
sudo rm -f /etc/nginx/sites-enabled/default

# Test Nginx configuration
sudo nginx -t
sudo systemctl restart nginx

# Install monitoring tools
echo "📊 Installing monitoring tools..."

# Install node_exporter for Prometheus monitoring
if ! command -v node_exporter &> /dev/null; then
    wget https://github.com/prometheus/node_exporter/releases/latest/download/node_exporter-1.6.1.linux-amd64.tar.gz
    tar xvfz node_exporter-1.6.1.linux-amd64.tar.gz
    sudo mv node_exporter-1.6.1.linux-amd64/node_exporter /usr/local/bin/
    rm -rf node_exporter-1.6.1.linux-amd64*

    # Create systemd service for node_exporter
    sudo tee /etc/systemd/system/node_exporter.service << EOF
[Unit]
Description=Node Exporter
After=network.target

[Service]
User=nobody
Group=nogroup
Type=simple
ExecStart=/usr/local/bin/node_exporter --web.listen-address=:9100

[Install]
WantedBy=multi-user.target
EOF

    sudo systemctl daemon-reload
    sudo systemctl enable node_exporter
    sudo systemctl start node_exporter
fi

# Set up log rotation
echo "📜 Setting up log rotation..."
sudo tee /etc/logrotate.d/cicd-notifier << EOF
/var/log/cicd-notifier/*.log {
    daily
    missingok
    rotate 52
    compress
    delaycompress
    notifempty
    create 644 $USER $USER
    postrotate
        systemctl reload nginx > /dev/null 2>&1 || true
    endscript
}
EOF

# Create backup script
echo "💾 Creating backup script..."
sudo tee /opt/cicd-notifier/backup.sh << 'EOF'
#!/bin/bash

# Backup script for CICD Status Notifier Bot
BACKUP_DIR="/opt/cicd-notifier/backups"
DATE=$(date +%Y%m%d_%H%M%S)

echo "Starting backup at $(date)"

# Create backup directory
mkdir -p $BACKUP_DIR

# Backup database
docker exec cicd_postgres pg_dump -U postgres cicd_notifier > $BACKUP_DIR/database_$DATE.sql

# Backup application data
tar -czf $BACKUP_DIR/app_data_$DATE.tar.gz /opt/cicd-notifier --exclude=/opt/cicd-notifier/backups

# Keep only last 7 days of backups
find $BACKUP_DIR -name "*.sql" -mtime +7 -delete
find $BACKUP_DIR -name "*.tar.gz" -mtime +7 -delete

echo "Backup completed at $(date)"
EOF

sudo chmod +x /opt/cicd-notifier/backup.sh

# Add backup to crontab
(crontab -l 2>/dev/null; echo "0 2 * * * /opt/cicd-notifier/backup.sh >> /var/log/cicd-notifier/backup.log 2>&1") | crontab -

# Create health check script
echo "🔍 Creating health check script..."
tee /opt/cicd-notifier/health-check.sh << 'EOF'
#!/bin/bash

# Health check script
LOG_FILE="/var/log/cicd-notifier/health-check.log"

echo "$(date): Starting health check" >> $LOG_FILE

# Check backend
if curl -f http://localhost:8080/health > /dev/null 2>&1; then
    echo "$(date): Backend: OK" >> $LOG_FILE
else
    echo "$(date): Backend: FAILED" >> $LOG_FILE
    # Restart backend container
    docker restart cicd_backend
fi

# Check frontend
if curl -f http://localhost:3000 > /dev/null 2>&1; then
    echo "$(date): Frontend: OK" >> $LOG_FILE
else
    echo "$(date): Frontend: FAILED" >> $LOG_FILE
    # Restart frontend container
    docker restart cicd_frontend
fi

# Check database
if docker exec cicd_postgres pg_isready -U postgres > /dev/null 2>&1; then
    echo "$(date): Database: OK" >> $LOG_FILE
else
    echo "$(date): Database: FAILED" >> $LOG_FILE
    # Restart database container
    docker restart cicd_postgres
fi

echo "$(date): Health check completed" >> $LOG_FILE
EOF

chmod +x /opt/cicd-notifier/health-check.sh

# Add health check to crontab (every 5 minutes)
(crontab -l 2>/dev/null; echo "*/5 * * * * /opt/cicd-notifier/health-check.sh") | crontab -

# Install Let's Encrypt for SSL (optional)
echo "🔒 Installing Certbot for SSL certificates..."
sudo apt install -y certbot python3-certbot-nginx

echo "✅ Server setup completed successfully!"
echo ""
echo "📋 Next steps:"
echo "1. Configure your GitHub repository secrets:"
echo "   - SSH_USERNAME: your server username"
echo "   - SSH_PRIVATE_KEY: your private SSH key"
echo "   - SSH_PORT: SSH port (default: 22)"
echo "   - POSTGRES_PASSWORD: secure database password"
echo "   - JWT_SECRET: random JWT secret key"
echo "   - TELEGRAM_BOT_TOKEN: your Telegram bot token"
echo "   - SONAR_TOKEN: SonarQube token"
echo "   - SONAR_HOST_URL: SonarQube server URL"
echo "   - SLACK_WEBHOOK_URL: Slack webhook for notifications"
echo ""
echo "2. To enable SSL, run:"
echo "   sudo certbot --nginx -d your-domain.com"
echo ""
echo "3. Test the setup by pushing to your repository"
echo ""
echo "🎉 Happy deploying!"
