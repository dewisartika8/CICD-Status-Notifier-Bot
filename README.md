# CI/CD Status Notifier Bot

![Go](https://img.shields.io/badge/go-%2300ADD8.svg?style=for-the-badge&logo=go&logoColor=white)
![React](https://img.shields.io/badge/react-%2320232a.svg?style=for-the-badge&logo=react&logoColor=%2361DAFB)
![PostgreSQL](https://img.shields.io/badge/postgres-%23316192.svg?style=for-the-badge&logo=postgresql&logoColor=white)
![Telegram](https://img.shields.io/badge/Telegram-2CA5E0?style=for-the-badge&logo=telegram&logoColor=white)
![Docker](https://img.shields.io/badge/docker-%230db7ed.svg?style=for-the-badge&logo=docker&logoColor=white)

> **Tugas Akhir - Workshop AI Enhancement**  
> **Team:** 2 developers | **Duration:** 8 weeks | **Methodology:** Test-Driven Development (TDD)

A comprehensive CI/CD status notification system that integrates with GitHub Actions to deliver real-time build and deployment notifications via Telegram, complete with a modern React dashboard for monitoring and analytics.

## 🚀 Project Overview

### Problem Statement
Development teams need immediate visibility into their CI/CD pipeline status. Current solutions often require manual checking of multiple platforms or lack consolidated reporting. This project creates a centralized notification system that provides real-time updates through Telegram and a comprehensive React dashboard with modern UI/UX.

### Solution
CI/CD Status Notifier Bot provides:
- **Real-time Telegram notifications** for build/deployment events
- **Interactive bot commands** for status queries
- **Web dashboard** for project monitoring and analytics
- **Multi-project support** with team-based subscriptions
- **Historical data tracking** for performance insights

## ✨ Key Features

### 🔔 Notification System
- ✅ Receive webhooks from GitHub Actions
- ✅ Send formatted notifications to Telegram
- ✅ Support for multiple event types (build success/failure, test results, deployment status)
- ✅ Smart notification filtering and subscription management

### 🤖 Telegram Bot Commands
- `/status` - View current status of all projects
- `/status <project>` - Get detailed status for specific project  
- `/projects` - List all monitored projects
- `/subscribe <project>` - Subscribe to project notifications
- `/unsubscribe <project>` - Unsubscribe from notifications
- `/history <project>` - View recent build history
- `/help` - Show available commands

### 📊 React Dashboard Features
- **Real-time Overview**: Live project status dashboard with auto-refresh
- **Analytics & Metrics**: Build success rates, deployment frequency, and performance trends
- **Interactive Charts**: Modern data visualization with filtering and drill-down capabilities
- **Project Management**: Configuration interface for managing monitored projects
- **User Management**: Team member administration with role-based permissions
- **Notification Settings**: Configure alert preferences and escalation rules
- **Historical Reports**: Detailed trend analysis and export capabilities
- **Responsive Design**: Mobile-first UI that works across all devices

### 🛠 Technical Features
- RESTful API backend with Go/Fiber
- PostgreSQL database with GORM
- React TypeScript frontend
- Docker containerization
- Comprehensive test coverage (>85%)
- CI/CD pipeline with GitHub Actions

## 🏗 Architecture

```
┌─────────────────┐    ┌──────────────────┐    ┌─────────────────┐
│   GitHub        │    │   Webhook        │    │   Telegram      │
│   Actions       │───▶│   Handler        │───▶│   Bot           │
│                 │    │   (Go/Fiber)     │    │   Notifications │
└─────────────────┘    └──────────────────┘    └─────────────────┘
                                │
                                ▼
                       ┌──────────────────┐    ┌─────────────────┐
                       │   PostgreSQL     │    │   React         │
                       │   Database       │◀───│   Dashboard     │
                       └──────────────────┘    │   (TypeScript)  │
                                               └─────────────────┘
```

## 🚀 Quick Start

### Prerequisites
- Go 1.21+
- Node.js 18+ with npm/yarn
- Docker & Docker Compose
- PostgreSQL 15+ (or use Docker)

### Installation

1. **Clone the repository**
```bash
git clone https://github.com/dewisartika8/CICD-Status-Notifier-Bot.git
cd CICD-Status-Notifier-Bot
```

2. **Set up environment**
```bash
cp backend/.env.example backend/.env
cp frontend/.env.example frontend/.env
# Edit .env files with your configuration
```

3. **Start with Docker Compose**
```bash
docker-compose up -d
```

4. **Or run locally**
```bash
# Backend
cd backend
go mod download
go run cmd/server/main.go

# Frontend (new terminal)
cd frontend
npm install
npm run dev
```

### Configuration

Create your `.env` file with:
```env
# Database
DB_HOST=localhost
DB_PORT=5432
DB_NAME=cicd_notifier
DB_USER=postgres
DB_PASSWORD=your_password

# Telegram Bot
TELEGRAM_BOT_TOKEN=your_bot_token

# Security
WEBHOOK_SECRET=your_webhook_secret
JWT_SECRET=your_jwt_secret
```

## 📋 Supported Status Types

| Status | Description | Icon | Trigger |
|--------|-------------|------|---------|
| build_started | Build process initiated | 🔄 | Workflow started |
| build_success | Build completed successfully | ✅ | Workflow completed (success) |
| build_failed | Build failed | ❌ | Workflow completed (failure) |
| test_passed | All tests passed | ✅ | Test job completed (success) |
| test_failed | Tests failed | ❌ | Test job completed (failure) |
| deployment_started | Deployment initiated | 🚀 | Deploy job started |
| deployment_success | Deployment successful | 🎉 | Deploy job completed (success) |
| deployment_failed | Deployment failed | 💥 | Deploy job completed (failure) |

## 🧪 Test-Driven Development

This project follows TDD methodology with comprehensive test coverage:

- **Unit Tests (70%):** Service logic, business rules, utilities
- **Integration Tests (20%):** Database operations, API endpoints
- **End-to-End Tests (10%):** Complete user workflows

### Running Tests

```bash
# Backend tests
cd backend
go test -v ./...
go test -race -coverprofile=coverage.out ./...

# Frontend tests  
cd frontend
npm test
npm run test:coverage

# Integration tests
docker-compose -f docker-compose.test.yml up --abort-on-container-exit
```

## 📈 Sprint Planning

**8-week development timeline with 4 sprints:**

### Sprint 1 (Weeks 1-2): Foundation
- Project setup and infrastructure
- Database schema and basic models
- Webhook endpoint with signature verification
- Basic project management API

### Sprint 2 (Weeks 3-4): Telegram Integration
- Telegram bot foundation
- Notification system implementation
- Basic bot commands (/status, /help)
- Subscription management

### Sprint 3 (Weeks 5-6): Dashboard Backend
- Dashboard API development
- Build metrics and analytics
- Advanced bot commands
- Enhanced webhook processing

### Sprint 4 (Weeks 7-8): Frontend & Polish
- React dashboard implementation
- Real-time updates
- Deployment setup
- Final testing and documentation

## 🛠 Tech Stack

### Backend
- **Language:** Go 1.21+
- **Framework:** Fiber v2
- **Database:** PostgreSQL 15+ with GORM v2
- **Testing:** Testify, GoMock
- **Documentation:** Swagger/OpenAPI

### Frontend
- **Framework:** React 18 + TypeScript
- **Build Tool:** Vite
- **Styling:** Tailwind CSS + Headless UI
- **State Management:** React Query + Zustand
- **Charts:** Recharts

### DevOps
- **Containerization:** Docker + Docker Compose
- **CI/CD:** GitHub Actions
- **Monitoring:** Structured logging
- **Security:** Environment-based secrets

## 📁 Project Structure

```
CICD-Status-Notifier-Bot/
├── docs/                     # 📚 Documentation
│   ├── PRD.md               # Product Requirements
│   ├── TECHNICAL_DESIGN.md  # Technical Architecture
│   ├── SPRINT_PLANNING.md   # Development Timeline
│   ├── TDD_GUIDE.md        # Testing Strategy
│   └── PROJECT_SETUP.md    # Setup Instructions
├── backend/                 # 🚀 Go Backend
│   ├── cmd/server/         # Application entry point
│   ├── internal/           # Internal packages
│   ├── pkg/                # Shared packages
│   ├── tests/              # Test files
│   └── migrations/         # Database migrations
├── frontend/               # 🎨 React Dashboard
│   ├── src/                # Source code
│   ├── public/             # Static assets
│   └── tests/              # Frontend tests
└── docker-compose.yml     # 🐳 Container orchestration
```

## 🤝 Contributing

This is an educational project for Binar Academy's final assignment. The development follows Agile Scrum methodology with TDD approach.

### Development Workflow
1. **Sprint Planning:** Define user stories and acceptance criteria
2. **TDD Cycle:** Red → Green → Refactor
3. **Code Review:** All changes require peer review
4. **Integration:** Continuous integration with automated testing

### Team Responsibilities
- **Developer 1:** Backend API, database, webhook processing
- **Developer 2:** Telegram bot, frontend dashboard, DevOps

## 📊 Success Metrics

### Technical KPIs
- **Test Coverage:** >85%
- **Notification Delivery:** >99% success rate
- **API Response Time:** <2 seconds
- **System Uptime:** >99.5%

### Product KPIs
- **Feature Completion:** 100% of MVP features
- **User Adoption:** >80% team member usage
- **Bug Rate:** <5 bugs per sprint in final phases

## 📚 Documentation

Comprehensive documentation is available in the `/docs` folder:

### 📋 Project Documentation
- **[PRD.md](docs/PRD.md)** - Product Requirements Document
- **[TECHNICAL_DESIGN.md](docs/TECHNICAL_DESIGN.md)** - Technical Architecture
- **[SPRINT_PLANNING.md](docs/SPRINT_PLANNING.md)** - 8-week Development Plan
- **[TDD_GUIDE.md](docs/TDD_GUIDE.md)** - Test-Driven Development Guide
- **[PROJECT_SETUP.md](docs/PROJECT_SETUP.md)** - Detailed Setup Instructions

### 👥 Team Task Management
- **[TASK_ASSIGNMENT_ARIF.md](docs/task/TASK_ASSIGNMENT_ARIF.md)** - Arif's Individual Task List (Backend Core Lead)
- **[TASK_ASSIGNMENT_DEWI.md](docs/task/TASK_ASSIGNMENT_DEWI.md)** - Dewi's Individual Task List (Integration & Frontend Lead)
- **[TEAM_PROGRESS_TRACKER.md](docs/task/TEAM_PROGRESS_TRACKER.md)** - Shared Progress Tracking Dashboard
- **[DAILY_COORDINATION_GUIDE.md](docs/task/DAILY_COORDINATION_GUIDE.md)** - Quick Reference for Daily Coordination

## 📄 License

This project is created for educational purposes as part of Binar Academy's AI Enhancement Workshop final assignment.

## 👥 Team

**Project Team:** 2 Developers  
**Duration:** 8 weeks (January - March 2025)  
**Methodology:** Agile Scrum with Test-Driven Development

---

> **Note:** This is a learning project focused on implementing best practices in software development, including TDD, clean architecture, and DevOps practices. The bot is designed to handle real-world CI/CD notification scenarios while maintaining high code quality and comprehensive test coverage.