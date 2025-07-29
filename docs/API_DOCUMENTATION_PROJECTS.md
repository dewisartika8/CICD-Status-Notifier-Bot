# Project Management API Documentation

## Overview
The Project Management API provides CRUD operations for managing CI/CD projects in the notification system.

**Base URL:** `/api/v1/projects`

---

## Endpoints

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

## Data Models

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

### Project Status Values
- `active`: Project is actively monitored
- `inactive`: Project monitoring is paused
- `archived`: Project is archived (read-only)

---

## Error Codes

### HTTP Status Codes
- `200 OK`: Request successful
- `201 Created`: Resource created successfully
- `400 Bad Request`: Invalid request data or validation failed
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

---

## Examples

### cURL Examples

#### Create a new project
```bash
curl -X POST http://localhost:8080/api/v1/projects \
  -H "Content-Type: application/json" \
  -d '{
    "name": "My Test Project",
    "repository_url": "https://github.com/user/test-repo",
    "webhook_secret": "my-secure-secret"
  }'
```

#### List all active projects
```bash
curl "http://localhost:8080/api/v1/projects?status=active&limit=10"
```

#### Get a specific project
```bash
curl "http://localhost:8080/api/v1/projects/550e8400-e29b-41d4-a716-446655440000"
```

#### Update a project
```bash
curl -X PUT http://localhost:8080/api/v1/projects/550e8400-e29b-41d4-a716-446655440000 \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Updated Project Name"
  }'
```

#### Delete a project
```bash
curl -X DELETE "http://localhost:8080/api/v1/projects/550e8400-e29b-41d4-a716-446655440000"
```

#### Update project status
```bash
curl -X PATCH http://localhost:8080/api/v1/projects/550e8400-e29b-41d4-a716-446655440000/status \
  -H "Content-Type: application/json" \
  -d '{
    "status": "inactive"
  }'
```

---

## Authentication & Security

Currently, the API does not require authentication. In production environments, consider implementing:

- API key authentication
- JWT tokens
- Rate limiting
- Input sanitization
- HTTPS enforcement

---

**API Version:** v1  
**Last Updated:** July 29, 2025  
**Implemented in:** Sprint 1 - Story 1.4
