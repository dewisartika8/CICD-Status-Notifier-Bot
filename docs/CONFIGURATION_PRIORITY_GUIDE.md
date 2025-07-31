# Configuration Priority Guide

Sistem konfigurasi aplikasi ini menggunakan **hierarki prioritas** yang memungkinkan flexibilitas dalam pengaturan konfigurasi berdasarkan environment dan kebutuhan deployment.

## Urutan Prioritas Konfigurasi

Konfigurasi akan diambil berdasarkan urutan prioritas berikut (dari tertinggi ke terendah):

### 1. **Environment Variables** (Prioritas Tertinggi) 🔥
Environment variables akan selalu menimpa nilai dari sumber lainnya.

**Contoh:**
```bash
# Windows PowerShell
$env:TELEGRAM_BOT_TOKEN="your_bot_token"
$env:SERVER_PORT=9000
$env:DB_HOST="production-db.example.com"

# Linux/macOS
export TELEGRAM_BOT_TOKEN="your_bot_token"
export SERVER_PORT=9000
export DB_HOST="production-db.example.com"
```

**Daftar Environment Variables yang Didukung:**
- `SERVER_PORT` - Port server aplikasi
- `TELEGRAM_BOT_TOKEN` - Token bot Telegram (prioritas utama)
- `DB_HOST` - Host database
- `DB_PORT` - Port database
- `DB_USER` - Username database
- `DB_PASSWORD` - Password database
- `DB_NAME` - Nama database
- `GITHUB_WEBHOOK_SECRET` - Secret webhook GitHub
- `GITLAB_WEBHOOK_SECRET` - Secret webhook GitLab
- `LOG_LEVEL` - Level logging (debug, info, warn, error)
- `ENVIRONMENT` - Environment aplikasi (development, staging, production)

### 2. **Configuration File** (Prioritas Menengah) 📄
Jika environment variable tidak ada, sistem akan mengambil nilai dari file `config/config.yaml`.

**Lokasi pencarian file konfigurasi:**
- `../../config/config.yaml`
- `./config/config.yaml`
- `./internal/config/config.yaml`
- `./config.yaml`

**Contoh konfigurasi di `config/config.yaml`:**
```yaml
server:
  port: 8081
  host: "localhost"

database:
  host: "127.0.0.1"
  port: "5432"
  user: "postgres"
  password: "password"
  dbname: "cicd_notifier"
  sslmode: "disable"

telegram:
  bot_token: "your_bot_token_here"
  webhook_url: "https://yourdomain.com/webhooks/telegram"

logging:
  level: "info"
  format: "json"
```

### 3. **Default Values** (Prioritas Terendah) ⚙️
Jika tidak ada environment variable atau konfigurasi file, sistem akan menggunakan nilai default.

**Default Values:**
```go
DefaultEnvironment        = "development"
DefaultServerPort         = 8080
DefaultServerHost         = "localhost"
DefaultDBHost             = "localhost"
DefaultDBPort             = "5432"
DefaultDBUser             = "postgres"
DefaultDBPassword         = "password"
DefaultDBName             = "cicd_notifier"
DefaultLogLevel           = "info"
```

## Contoh Penggunaan

### Development Environment
```yaml
# config/config.yaml
environment: "development"
server:
  port: 8081
database:
  host: "localhost"
  password: "dev_password"
telegram:
  bot_token: "dev_bot_token"
```

### Production Environment
```bash
# Environment variables untuk production
export ENVIRONMENT="production"
export SERVER_PORT=80
export DB_HOST="prod-db.example.com"
export DB_PASSWORD="super_secure_password"
export TELEGRAM_BOT_TOKEN="prod_bot_token"
```

### Testing Environment
Environment variables akan menimpa konfigurasi file:
```bash
# Override hanya beberapa nilai untuk testing
export DB_NAME="test_db"
export LOG_LEVEL="debug"
# Nilai lainnya akan tetap menggunakan config file atau default
```

## Validasi Konfigurasi

Sistem akan melakukan validasi terhadap konfigurasi yang dimuat:

### Required Fields
- `telegram.bot_token` - Wajib diisi untuk bot telegram
- `database.host` - Host database
- `database.port` - Port database  
- `database.user` - Username database
- `database.dbname` - Nama database

### Validation Rules
- `server.port` - Harus > 0
- `database.max_open_conns` - Harus > 0
- `database.max_idle_conns` - Harus >= 0
- `logging.level` - Harus salah satu dari: debug, info, warn, error, fatal, panic
- `logging.format` - Harus salah satu dari: json, text

## Error Handling

Jika konfigurasi tidak valid, aplikasi akan:
1. Menampilkan error detail tentang field yang bermasalah
2. Keluar dengan exit code 1
3. Mencatat error ke log

**Contoh error message:**
```
Config error: config validation failed: validation errors: config validation error for field 'telegram.bot_token': bot token is required; config validation error for field 'server.port': port must be set
```

## Best Practices

1. **Development**: Gunakan file konfigurasi untuk kemudahan
2. **Production**: Gunakan environment variables untuk keamanan
3. **Docker**: Combine environment variables dengan config file mounting
4. **Sensitive Data**: Selalu gunakan environment variables untuk password, token, dan secret
5. **Default Values**: Tidak perlu set nilai default di environment atau config jika sudah sesuai

## Troubleshooting

### Bot Token Issues
```bash
# Cek apakah environment variable sudah diset
echo $TELEGRAM_BOT_TOKEN  # Linux/macOS
echo $env:TELEGRAM_BOT_TOKEN  # Windows PowerShell

# Cek isi config file
cat config/config.yaml | grep bot_token
```

### Database Connection Issues
```bash
# Cek konfigurasi database
echo "Host: $DB_HOST, Port: $DB_PORT, User: $DB_USER"
```

### Server Port Conflicts
```bash
# Cek port yang digunakan
netstat -an | grep :8081  # Check if port is in use
```
