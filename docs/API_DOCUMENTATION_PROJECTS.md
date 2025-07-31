# CI/CD Status Notifier Bot - API Documentation

## Overview
The CI/CD Status Notifier Bot API provides comprehensive endpoints for managing CI/CD projects, processing webhooks, managing Telegram bot interactions, and monitoring system health.

**Base URL:** `http://localhost:8080`
**API Version:** `v1`

---

## Table of Contents
1. [Health Check Endpoints](#health-check-endpoints)
2. [Project Management API](#project-management-api)
3. [Webhook API](#webhook-api)
4. [Telegram Bot API](#telegram-bot-api)
5. [Common Data Models](#common-data-models)
6. [Error Handling](#error-handling)
7. [Authentication & Security](#authentication--security)

---

## Health Check Endpoints

### 1. Root Health Check
**GET** `/`

Returns basic service status and version information.

#### Response (200 OK)
```json
{
    "message": "CI/CD Status Notifier Bot is running 🚀",
    "status": "healthy",
    "version": "1.0.0"
}
```

### 2. Detailed Health Check
**GET** `/health`

Returns detailed health status including database connectivity.

#### Response (200 OK)
```json
{
    "status": "healthy",
    "database": "connected",
    "timestamp": "2025-07-31T10:00:00Z"
}
```

---

# Project Management API

**Base URL:** `/api/v1/projects`

## Project Management Endpoints

### 1. Create Project
**POST** `/api/v1/projects`

Creates a new project with the provided information.

#### Request Body
```json
{
    "name": "My CI/CD Project",
    "repository_url": "https://github.com/user/repo",
    "webhook_secret": "my-secure-webhook-secret",
    "telegram_chat_id": -1001234567890
}
```

#### Response (201 Created)
```json
{
    "message": "Project created successfully",
    "data": {
        "id": "550e8400-e29b-41d4-a716-446655440000",
        "name": "My CI/CD Project",
        "repository_url": "https://github.com/user/repo",
        "status": "active",
        "telegram_chat_id": -1001234567890,
        "created_at": "2025-07-29T10:00:00Z",
        "updated_at": "2025-07-29T10:00:00Z"
    }
}
```

#### Error Response (400 Bad Request)
```json
{
    "error": "Validation failed",
    "details": "Key: 'CreateProjectRequest.Name' Error:Field validation for 'Name' failed on the 'required' tag"
}
```

---

### 2. List Projects
**GET** `/api/v1/projects`

Retrieves a list of projects with optional filtering and pagination.

#### Query Parameters
- `status` (optional): Filter by project status (active, inactive, archived)
- `name` (optional): Filter by project name (partial match)
- `repository_url` (optional): Filter by repository URL (partial match)
- `has_telegram_chat` (optional): Filter projects with/without telegram chat (true/false)
- `limit` (optional): Number of results to return (default: 10, max: 100)
- `offset` (optional): Number of results to skip (default: 0)
- `sort_by` (optional): Sort field (name, repository_url, status, created_at, updated_at)
- `sort_order` (optional): Sort order (asc, desc)

#### Example Request
```
GET /api/v1/projects?status=active&limit=5&offset=0&sort_by=name&sort_order=asc
```

#### Response (200 OK)
```json
{
    "message": "Projects retrieved successfully",
    "data": {
        "projects": [
            {
                "id": "550e8400-e29b-41d4-a716-446655440000",
                "name": "My CI/CD Project",
                "repository_url": "https://github.com/user/repo",
                "status": "active",
                "telegram_chat_id": -1001234567890,
                "created_at": "2025-07-29T10:00:00Z",
                "updated_at": "2025-07-29T10:00:00Z"
            },
            {
                "id": "550e8400-e29b-41d4-a716-446655440001",
                "name": "Another Project",
                "repository_url": "https://github.com/user/another-repo",
                "status": "active",
                "telegram_chat_id": null,
                "created_at": "2025-07-29T09:00:00Z",
                "updated_at": "2025-07-29T09:00:00Z"
            }
        ],
        "total": 25,
        "limit": 5,
        "offset": 0
    }
}
```

---

### 3. Get Project by ID
**GET** `/api/v1/projects/{id}`

Retrieves a specific project by its ID.

#### Path Parameters
- `id` (required): Project UUID

#### Example Request
```
GET /api/v1/projects/550e8400-e29b-41d4-a716-446655440000
```

#### Response (200 OK)
```json
{
    "message": "Project retrieved successfully",
    "data": {
        "id": "550e8400-e29b-41d4-a716-446655440000",
        "name": "My CI/CD Project",
        "repository_url": "https://github.com/user/repo",
        "status": "active",
        "telegram_chat_id": -1001234567890,
        "created_at": "2025-07-29T10:00:00Z",
        "updated_at": "2025-07-29T10:00:00Z"
    }
}
```

#### Error Response (404 Not Found)
```json
{
    "error": "project not found"
}
```

---

### 4. Update Project
**PUT** `/api/v1/projects/{id}`

Updates an existing project with the provided information.

#### Path Parameters
- `id` (required): Project UUID

#### Request Body
```json
{
    "name": "Updated Project Name",
    "repository_url": "https://github.com/user/updated-repo",
    "webhook_secret": "updated-webhook-secret",
    "telegram_chat_id": -1001234567891
}
```

#### Response (200 OK)
```json
{
    "message": "Project updated successfully",
    "data": {
        "id": "550e8400-e29b-41d4-a716-446655440000",
        "name": "Updated Project Name",
        "repository_url": "https://github.com/user/updated-repo",
        "status": "active",
        "telegram_chat_id": -1001234567891,
        "created_at": "2025-07-29T10:00:00Z",
        "updated_at": "2025-07-29T11:00:00Z"
    }
}
```

---

### 5. Delete Project
**DELETE** `/api/v1/projects/{id}`

Deletes a project by its ID.

#### Path Parameters
- `id` (required): Project UUID

#### Example Request
```
DELETE /api/v1/projects/550e8400-e29b-41d4-a716-446655440000
```

#### Response (200 OK)
```json
{
    "message": "Project deleted successfully"
}
```

#### Error Response (404 Not Found)
```json
{
    "error": "project not found"
}
```

---

### 6. Update Project Status
**PATCH** `/api/v1/projects/{id}/status`

Updates the status of a specific project.

#### Path Parameters
- `id` (required): Project UUID

#### Request Body
```json
{
    "status": "inactive"
}
```

#### Response (200 OK)
```json
{
    "message": "Project status updated successfully",
    "data": {
        "id": "550e8400-e29b-41d4-a716-446655440000",
        "name": "My CI/CD Project",
        "repository_url": "https://github.com/user/repo",
        "status": "inactive",
        "telegram_chat_id": -1001234567890,
        "created_at": "2025-07-29T10:00:00Z",
        "updated_at": "2025-07-29T11:30:00Z"
    }
}
```

---

# Webhook API

**Base URL:** `/api/v1/webhooks`

The Webhook API handles incoming webhook events from GitHub and provides endpoints to retrieve webhook event data.

## Webhook Endpoints

### 1. Process GitHub Webhook
**POST** `/api/v1/webhooks/github/{projectId}`

Processes incoming GitHub webhook events for a specific project.

#### Path Parameters
- `projectId` (required): Project UUID

#### Headers Required
- `X-Hub-Signature-256` (required): GitHub HMAC-SHA256 signature for payload verification
- `X-GitHub-Event` (required): Type of GitHub event (workflow_run, push, pull_request)
- `X-GitHub-Delivery` (optional): GitHub delivery ID for idempotency
- `Content-Type`: application/json

#### Supported Event Types
- `workflow_run`: GitHub Actions workflow events
- `push`: Repository push events
- `pull_request`: Pull request events

#### Request Body
GitHub webhook payload (varies by event type)

#### Example Request
```bash
POST /api/v1/webhooks/github/550e8400-e29b-41d4-a716-446655440000
X-Hub-Signature-256: sha256=1234567890abcdef...
X-GitHub-Event: workflow_run
X-GitHub-Delivery: 12345-67890-abcdef
Content-Type: application/json

{
    "action": "completed",
    "workflow_run": {
        "id": 12345,
        "name": "CI",
        "status": "completed",
        "conclusion": "success"
    },
    "repository": {
        "name": "my-repo",
        "full_name": "user/my-repo"
    }
}
```

#### Response (202 Accepted)
```json
{
    "message": "webhook processed successfully",
    "data": {
        "id": "webhook-event-uuid",
        "project_id": "550e8400-e29b-41d4-a716-446655440000",
        "event_type": "workflow_run",
        "delivery_id": "12345-67890-abcdef",
        "processed_at": "2025-07-31T10:00:00Z",
        "created_at": "2025-07-31T10:00:00Z"
    }
}
```

#### Error Responses

**400 Bad Request**
```json
{
    "error": "project_id is required"
}
```

**401 Unauthorized**
```json
{
    "error": "missing X-Hub-Signature-256 header"
}
```

**404 Not Found**
```json
{
    "error": "project not found"
}
```

### 2. Get Webhook Events by Project
**GET** `/api/v1/webhooks/events/{projectId}`

Retrieves webhook events for a specific project with pagination support.

#### Path Parameters
- `projectId` (required): Project UUID

#### Query Parameters
- `limit` (optional): Number of results to return (default: 20, max: 100)
- `offset` (optional): Number of results to skip (default: 0)

#### Example Request
```
GET /api/v1/webhooks/events/550e8400-e29b-41d4-a716-446655440000?limit=10&offset=0
```

#### Response (200 OK)
```json
{
    "data": [
        {
            "id": "webhook-event-uuid-1",
            "project_id": "550e8400-e29b-41d4-a716-446655440000",
            "event_type": "workflow_run",
            "delivery_id": "12345-67890-abcdef",
            "processed_at": "2025-07-31T10:00:00Z",
            "created_at": "2025-07-31T10:00:00Z"
        },
        {
            "id": "webhook-event-uuid-2",
            "project_id": "550e8400-e29b-41d4-a716-446655440000",
            "event_type": "push",
            "delivery_id": "67890-abcdef-12345",
            "processed_at": "2025-07-31T09:30:00Z",
            "created_at": "2025-07-31T09:30:00Z"
        }
    ],
    "pagination": {
        "limit": 10,
        "offset": 0,
        "count": 2
    }
}
```

### 3. Get Specific Webhook Event
**GET** `/api/v1/webhooks/events/{projectId}/{eventId}`

Retrieves details of a specific webhook event including the full payload.

#### Path Parameters
- `projectId` (required): Project UUID
- `eventId` (required): Webhook event UUID

#### Example Request
```
GET /api/v1/webhooks/events/550e8400-e29b-41d4-a716-446655440000/webhook-event-uuid-1
```

#### Response (200 OK)
```json
{
    "data": {
        "id": "webhook-event-uuid-1",
        "project_id": "550e8400-e29b-41d4-a716-446655440000",
        "event_type": "workflow_run",
        "delivery_id": "12345-67890-abcdef",
        "payload": {
            "action": "completed",
            "workflow_run": {
                "id": 12345,
                "name": "CI",
                "status": "completed",
                "conclusion": "success"
            }
        },
        "processed_at": "2025-07-31T10:00:00Z",
        "created_at": "2025-07-31T10:00:00Z"
    }
}
```

#### Error Response (404 Not Found)
```json
{
    "error": "webhook event not found"
}
```

---

# Telegram Bot API

**Base URL:** `/api/v1/telegram`

The Telegram Bot API handles Telegram webhook events and provides webhook management endpoints.

## Telegram Endpoints

### 1. Telegram Webhook Handler
**POST** `/api/v1/telegram/webhook`

Handles incoming webhook updates from Telegram Bot API.

#### Request Body
Telegram Update object (automatically sent by Telegram)

#### Response (200 OK)
```json
{
    "status": "ok"
}
```

### 2. Set Telegram Webhook
**POST** `/api/v1/telegram/webhook/set`

Sets the webhook URL for the Telegram bot.

#### Request Body
```json
{
    "webhook_url": "https://your-domain.com"
}
```

#### Response (200 OK)
```json
{
    "message": "Webhook set successfully",
    "webhook_url": "https://your-domain.com/api/v1/telegram/webhook"
}
```

#### Error Response (400 Bad Request)
```json
{
    "error": "webhook_url is required"
}
```

### 3. Delete Telegram Webhook
**DELETE** `/api/v1/telegram/webhook`

Removes the webhook URL for the Telegram bot.

#### Response (200 OK)
```json
{
    "message": "Webhook deleted successfully"
}
```

---

## Data Models

## Common Data Models

### Project
```json
{
    "id": "string (UUID)",
    "name": "string (1-100 characters, unique)",
    "repository_url": "string (valid URL, unique)",
    "status": "string (active|inactive|archived)",
    "telegram_chat_id": "integer (optional)",
    "created_at": "string (ISO 8601 timestamp)",
    "updated_at": "string (ISO 8601 timestamp)"
}
```

### Webhook Event
```json
{
    "id": "string (UUID)",
    "project_id": "string (UUID)",
    "event_type": "string (workflow_run|push|pull_request)",
    "delivery_id": "string (optional)",
    "payload": "object (GitHub webhook payload)",
    "processed_at": "string (ISO 8601 timestamp, nullable)",
    "created_at": "string (ISO 8601 timestamp)"
}
```

### Project Status Values
- `active`: Project is actively monitored
- `inactive`: Project monitoring is paused
- `archived`: Project is archived (read-only)

### Webhook Event Types
- `workflow_run`: GitHub Actions workflow events
- `push`: Repository push events  
- `pull_request`: Pull request events

---

## Error Handling

### HTTP Status Codes
- `200 OK`: Request successful
- `201 Created`: Resource created successfully
- `202 Accepted`: Request accepted for processing (webhooks)
- `400 Bad Request`: Invalid request data or validation failed
- `401 Unauthorized`: Authentication failed (webhook signature invalid)
- `404 Not Found`: Resource not found
- `409 Conflict`: Resource already exists (duplicate name or repository URL)
- `500 Internal Server Error`: Server error

### Common Error Responses

#### Validation Error (400)
```json
{
    "error": "Validation failed",
    "details": "Detailed validation error message"
}
```

#### Authentication Error (401)
```json
{
    "error": "missing X-Hub-Signature-256 header"
}
```

#### Not Found Error (404)
```json
{
    "error": "project not found"
}
```

#### Conflict Error (409)
```json
{
    "error": "project with this name already exists"
}
```

#### Server Error (500)
```json
{
    "error": "Internal server error"
}
```

### Webhook-Specific Errors

#### Invalid Webhook Signature (401)
```json
{
    "error": "invalid webhook signature"
}
```

#### Unsupported Event Type (400)
```json
{
    "error": "unsupported event type: issues"
}
```

#### Empty Request Body (400)
```json
{
    "error": "empty request body"
}
```

---

## Examples

### cURL Examples

#### Health Check Examples

##### Root health check
```bash
curl http://localhost:8080/
```

##### Detailed health check
```bash
curl http://localhost:8080/health
```

#### Project Management Examples

##### Create a new project
```bash
curl -X POST http://localhost:8080/api/v1/projects \
  -H "Content-Type: application/json" \
  -d '{
    "name": "My Test Project",
    "repository_url": "https://github.com/user/test-repo",
    "webhook_secret": "my-secure-secret"
  }'
```

##### List all active projects
```bash
curl "http://localhost:8080/api/v1/projects?status=active&limit=10"
```

##### Get a specific project
```bash
curl "http://localhost:8080/api/v1/projects/550e8400-e29b-41d4-a716-446655440000"
```

##### Update a project
```bash
curl -X PUT http://localhost:8080/api/v1/projects/550e8400-e29b-41d4-a716-446655440000 \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Updated Project Name"
  }'
```

##### Delete a project
```bash
curl -X DELETE "http://localhost:8080/api/v1/projects/550e8400-e29b-41d4-a716-446655440000"
```

##### Update project status
```bash
curl -X PATCH http://localhost:8080/api/v1/projects/550e8400-e29b-41d4-a716-446655440000/status \
  -H "Content-Type: application/json" \
  -d '{
    "status": "inactive"
  }'
```

#### Webhook Examples

##### Process GitHub webhook (automatically called by GitHub)
```bash
curl -X POST http://localhost:8080/api/v1/webhooks/github/550e8400-e29b-41d4-a716-446655440000 \
  -H "Content-Type: application/json" \
  -H "X-Hub-Signature-256: sha256=1234567890abcdef..." \
  -H "X-GitHub-Event: workflow_run" \
  -H "X-GitHub-Delivery: 12345-67890-abcdef" \
  -d '{
    "action": "completed",
    "workflow_run": {
      "id": 12345,
      "name": "CI",
      "status": "completed",
      "conclusion": "success"
    }
  }'
```

##### List webhook events for a project
```bash
curl "http://localhost:8080/api/v1/webhooks/events/550e8400-e29b-41d4-a716-446655440000?limit=10&offset=0"
```

##### Get specific webhook event
```bash
curl "http://localhost:8080/api/v1/webhooks/events/550e8400-e29b-41d4-a716-446655440000/webhook-event-uuid"
```

#### Telegram Bot Examples

##### Set Telegram webhook
```bash
curl -X POST http://localhost:8080/api/v1/telegram/webhook/set \
  -H "Content-Type: application/json" \
  -d '{
    "webhook_url": "https://your-domain.com"
  }'
```

##### Delete Telegram webhook
```bash
curl -X DELETE http://localhost:8080/api/v1/telegram/webhook
```

---

## Authentication & Security

### Current Implementation
- **Webhook Signature Verification**: GitHub webhooks are verified using HMAC-SHA256 signatures
- **No API Authentication**: Other endpoints currently do not require authentication
- **Input Validation**: All endpoints perform request validation
- **SQL Injection Protection**: GORM ORM with prepared statements

### Security Headers Required for Webhooks
- `X-Hub-Signature-256`: Required for GitHub webhook signature verification
- `X-GitHub-Event`: Required to identify the event type
- `X-GitHub-Delivery`: Optional but recommended for idempotency

### Production Security Recommendations
Consider implementing the following for production environments:

- **API Key Authentication**: For project management endpoints
- **JWT Tokens**: For user session management
- **Rate Limiting**: To prevent abuse
- **HTTPS Enforcement**: For secure communication
- **Input Sanitization**: Additional validation layers
- **CORS Configuration**: Proper cross-origin request handling
- **Logging & Monitoring**: Security event tracking

### Webhook Security
- **Signature Verification**: All GitHub webhooks must include valid HMAC-SHA256 signatures
- **Event Type Validation**: Only supported event types are processed
- **Project Validation**: Webhooks are validated against existing projects
- **Idempotency**: Duplicate webhooks are handled gracefully using delivery IDs

---

**API Documentation**  
**Version:** v1  
**Last Updated:** July 31, 2025  
**Backend Implementation:** Go with Fiber framework  
**Database:** PostgreSQL with GORM  
**Architecture:** Hexagonal/Clean Architecture

---

## Additional Resources

- [Webhook Implementation Guide](./WEBHOOK_IMPLEMENTATION.md)
- [Technical Design Document](./TECHNICAL_DESIGN.md)
- [Project Setup Guide](./PROJECT_SETUP.md)
- [Implementation Summary](./IMPLEMENTATION_SUMMARY_1.4.md)
