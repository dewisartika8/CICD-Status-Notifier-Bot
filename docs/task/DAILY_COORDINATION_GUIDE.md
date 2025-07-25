# Quick Reference Guide
## Daily Coordination for Arif & Dewi

> **🚀 Quick access to daily tasks, coordination points, and progress tracking**

---

## 📋 Today's Focus

### Current Sprint: **Sprint 1** | Day: **___** | Date: **___________**

#### 🧑‍💻 Arif's Today Tasks:
- [ ] **Primary:** ________________________________
- [ ] **Secondary:** ________________________________
- [ ] **Integration:** ________________________________

#### 👩‍💻 Dewi's Today Tasks:
- [ ] **Primary:** ________________________________
- [ ] **Secondary:** ________________________________
- [ ] **Integration:** ________________________________

---

## 🤝 Daily Coordination Checklist

### Morning Standup (9:00 AM) - 15 minutes:
- [ ] **Arif Update:**
  - Yesterday: ____________________________
  - Today: ______________________________
  - Blockers: ____________________________

- [ ] **Dewi Update:**
  - Yesterday: ____________________________
  - Today: ______________________________
  - Blockers: ____________________________

- [ ] **Integration Points:** ____________________________
- [ ] **Decisions Needed:** ____________________________

### Daily Sync Points:
- [ ] **API Contract Review** (if needed)
- [ ] **Code Review** (active PRs)
- [ ] **Testing Coordination** (integration tests)
- [ ] **Blocker Resolution** (immediate help needed)

---

## 🔄 Quick Status Check

### Arif's Component Status:
| Component | Status | Next Action |
|-----------|---------|-------------|
| Database Schema | ⬜ Not Started / 🟡 In Progress / ✅ Complete | _____________ |
| Repository Layer | ⬜ Not Started / 🟡 In Progress / ✅ Complete | _____________ |
| API Endpoints | ⬜ Not Started / 🟡 In Progress / ✅ Complete | _____________ |
| Tests | ⬜ Not Started / 🟡 In Progress / ✅ Complete | _____________ |

### Dewi's Component Status:
| Component | Status | Next Action |
|-----------|---------|-------------|
| Project Setup | ⬜ Not Started / 🟡 In Progress / ✅ Complete | _____________ |
| Webhook Handler | ⬜ Not Started / 🟡 In Progress / ✅ Complete | _____________ |
| Bot Foundation | ⬜ Not Started / 🟡 In Progress / ✅ Complete | _____________ |
| DevOps Setup | ⬜ Not Started / 🟡 In Progress / ✅ Complete | _____________ |

---

## 🎯 Sprint Goals Quick Reference

### Sprint 1 (Current) - Foundation:
**Arif Focus:** Database + API  
**Dewi Focus:** Setup + Webhooks  
**Integration:** API contracts defined

### Sprint 2 (Next) - Bot + Notifications:
**Arif Focus:** Notification system  
**Dewi Focus:** Telegram bot  
**Integration:** Bot commands with API

### Sprint 3 (Week 5-6) - Advanced Features:
**Arif Focus:** Dashboard API + Analytics  
**Dewi Focus:** Advanced bot commands  
**Integration:** Complete feature set

### Sprint 4 (Week 7-8) - Polish + Deploy:
**Arif Focus:** Documentation + Performance  
**Dewi Focus:** Deployment + Monitoring  
**Integration:** Production readiness

---

## 🚨 Quick Emergency Contacts

### Technical Issues:
- **Go/Fiber Issues:** [Documentation Link] / [Community Forum]
- **PostgreSQL Issues:** [Documentation Link] / [Stack Overflow]
- **Telegram Bot API:** [Documentation Link] / [Support]
- **Docker Issues:** [Documentation Link] / [Community Help]

### Team Communication:
- **Immediate Help:** WhatsApp/Slack
- **Code Review:** GitHub PR mentions
- **Design Questions:** Schedule quick call
- **Blocker Resolution:** Daily standup or immediate sync

---

## 📊 Quick Progress Tracking

### This Week's Targets:
- [ ] **Week Goal 1:** ____________________________
- [ ] **Week Goal 2:** ____________________________
- [ ] **Week Goal 3:** ____________________________

### Daily Progress:
| Day | Arif Progress | Dewi Progress | Integration Points |
|-----|---------------|---------------|-------------------|
| Mon | _____________ | _____________ | _________________ |
| Tue | _____________ | _____________ | _________________ |
| Wed | _____________ | _____________ | _________________ |
| Thu | _____________ | _____________ | _________________ |
| Fri | _____________ | _____________ | _________________ |

---

## 🔧 Development Commands Quick Reference

### Arif's Common Commands:
```bash
# Database operations
go run cmd/server/main.go migrate up
go run cmd/server/main.go migrate down

# Testing
go test ./internal/repository/... -v
go test ./internal/services/... -v

# API testing
curl -X POST http://localhost:8080/api/v1/projects
go run cmd/server/main.go # Start server
```

### Dewi's Common Commands:
```bash
# Docker operations
docker-compose up -d
docker-compose down
docker-compose logs -f app

# Bot testing
curl -X POST http://localhost:8080/webhook/github/test-project
# Test Telegram bot locally

# CI/CD
git push origin feature/webhook-integration
# Trigger GitHub Actions
```

### Shared Commands:
```bash
# Code quality
go fmt ./...
go vet ./...
golangci-lint run

# Testing
go test ./... -v
go test ./... -race -coverprofile=coverage.out

# Git workflow
git checkout -b feature/new-feature
git add .
git commit -m "feat: description"
git push origin feature/new-feature
```

---

## 📝 Daily Notes Template

### Date: **___________**

#### Morning Plan:
- **Arif:** ________________________________
- **Dewi:** ________________________________
- **Together:** ________________________________

#### Midday Check:
- **Progress:** ________________________________
- **Issues:** ________________________________
- **Adjustments:** ________________________________

#### End of Day:
- **Completed:** ________________________________
- **Blocked:** ________________________________
- **Tomorrow:** ________________________________

#### Integration Points:
- **API Changes:** ________________________________
- **Shared Code:** ________________________________
- **Testing Needs:** ________________________________

---

## 🎯 Quick Decision Framework

### When to Pair Program:
- ✅ Complex integration points
- ✅ New technology/framework
- ✅ Debugging difficult issues
- ✅ Architecture decisions

### When to Work Independently:
- ✅ Well-defined tasks
- ✅ Individual component development
- ✅ Writing tests
- ✅ Documentation

### When to Ask for Help:
- 🚨 Blocked for >2 hours
- 🚨 Architecture uncertainty
- 🚨 Integration conflicts
- 🚨 Timeline concerns

---

## 📋 End-of-Day Checklist

### Before Finishing:
- [ ] Push code changes to repository
- [ ] Update task status in progress tracker
- [ ] Review any PRs from teammate
- [ ] Plan tomorrow's priorities
- [ ] Update any blockers or concerns
- [ ] Commit to shared progress document

### Weekly Review (Fridays):
- [ ] Review sprint progress
- [ ] Update team progress tracker
- [ ] Plan weekend work (if any)
- [ ] Prepare for next week
- [ ] Document lessons learned

---

**Last Updated:** ___________  
**Next Standup:** ___________  
**Sprint End:** ___________
