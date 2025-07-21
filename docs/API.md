# API Documentation

This document provides comprehensive API documentation for the Microservices Demo application.

## Table of Contents

- [Overview](#overview)
- [API Server Endpoints](#api-server-endpoints)
- [Frontend Endpoints](#frontend-endpoints)
- [Message Queue Integration](#message-queue-integration)
- [Data Models](#data-models)
- [Error Handling](#error-handling)
- [Examples](#examples)

## Overview

The Microservices Demo application exposes two main HTTP interfaces:
- **API Server** (Port 8081): RESTful JSON API for programmatic access
- **Frontend** (Port 8080): Web interface for human interaction

## API Server Endpoints

### Base URL
```
http://localhost:8081
```

### Authentication
Currently, no authentication is required. This is a demo application.

### Content Type
All API endpoints expect and return `application/json` unless otherwise specified.

---

### Jobs API

#### Create Job
Creates a new job and queues it for processing.

**Endpoint:** `POST /api/jobs`

**Request Body:**
```json
{
  "description": "string"
}
```

**Response:** `201 Created`
```json
{
  "id": "550e8400-e29b-41d4-a716-446655440000",
  "description": "Process data file",
  "status": "pending",
  "created_at": "2025-07-20T10:30:00Z",
  "updated_at": "2025-07-20T10:30:00Z"
}
```

**Example:**
```bash
curl -X POST http://localhost:8081/api/jobs \
  -H "Content-Type: application/json" \
  -d '{"description": "Process user data"}'
```

---

#### Get Job by ID
Retrieves the current status and details of a specific job.

**Endpoint:** `GET /api/jobs/{id}`

**Path Parameters:**
- `id` (string): UUID of the job

**Response:** `200 OK`
```json
{
  "id": "550e8400-e29b-41d4-a716-446655440000",
  "description": "Process data file",
  "status": "completed",
  "result": "Successfully processed 1,234 records",
  "created_at": "2025-07-20T10:30:00Z",
  "updated_at": "2025-07-20T10:35:30Z",
  "started_at": "2025-07-20T10:30:15Z",
  "completed_at": "2025-07-20T10:35:30Z"
}
```

**Example:**
```bash
curl http://localhost:8081/api/jobs/550e8400-e29b-41d4-a716-446655440000
```

---

#### List All Jobs
Retrieves a list of all jobs with pagination support.

**Endpoint:** `GET /api/jobs`

**Query Parameters:**
- `limit` (integer, optional): Maximum number of jobs to return (default: 50, max: 100)
- `offset` (integer, optional): Number of jobs to skip (default: 0)
- `status` (string, optional): Filter by job status (`pending`, `processing`, `completed`, `failed`)

**Response:** `200 OK`
```json
{
  "jobs": [
    {
      "id": "550e8400-e29b-41d4-a716-446655440000",
      "description": "Process data file",
      "status": "completed",
      "created_at": "2025-07-20T10:30:00Z",
      "updated_at": "2025-07-20T10:35:30Z"
    }
  ],
  "total": 1,
  "limit": 50,
  "offset": 0
}
```

**Examples:**
```bash
# Get all jobs
curl http://localhost:8081/api/jobs

# Get first 10 jobs
curl "http://localhost:8081/api/jobs?limit=10"

# Get only pending jobs
curl "http://localhost:8081/api/jobs?status=pending"

# Pagination example
curl "http://localhost:8081/api/jobs?limit=10&offset=20"
```

---

### Health Check

#### API Health
Returns the health status of the API server.

**Endpoint:** `GET /health`

**Response:** `200 OK`
```json
{
  "status": "healthy",
  "timestamp": "2025-07-20T10:30:00Z",
  "version": "1.0.0",
  "dependencies": {
    "rabbitmq": "connected",
    "database": "not_configured"
  }
}
```

**Example:**
```bash
curl http://localhost:8081/health
```

## Frontend Endpoints

### Base URL
```
http://localhost:8080
```

### Web Interface

#### Main Page
**Endpoint:** `GET /`

Returns the main HTML interface for job submission and monitoring.

**Response:** `200 OK` (HTML content)

---

#### Submit Job (Form)
**Endpoint:** `POST /submit`

**Content-Type:** `application/x-www-form-urlencoded`

**Form Data:**
- `description` (string): Job description

**Response:** `302 Found` (Redirect to status page)

**Example:**
```bash
curl -X POST http://localhost:8080/submit \
  -H "Content-Type: application/x-www-form-urlencoded" \
  -d "description=Process+user+data"
```

---

#### Job Status Page
**Endpoint:** `GET /status/{id}`

**Path Parameters:**
- `id` (string): UUID of the job

**Response:** `200 OK` (HTML content)

Returns an HTML page displaying job status with auto-refresh functionality.

## Message Queue Integration

### RabbitMQ Topology

```
Exchange: job_exchange (direct)
├── Queue: job_queue
│   ├── Routing Key: job.created
│   └── Consumer: Job Runner
└── Queue: status_queue
    ├── Routing Key: job.status
    └── Consumer: API Server (for status updates)
```

### Message Formats

#### Job Creation Message
**Exchange:** `job_exchange`  
**Routing Key:** `job.created`

```json
{
  "id": "550e8400-e29b-41d4-a716-446655440000",
  "description": "Process data file",
  "created_at": "2025-07-20T10:30:00Z"
}
```

#### Job Status Update Message
**Exchange:** `job_exchange`  
**Routing Key:** `job.status`

```json
{
  "id": "550e8400-e29b-41d4-a716-446655440000",
  "status": "processing",
  "updated_at": "2025-07-20T10:30:15Z",
  "started_at": "2025-07-20T10:30:15Z"
}
```

#### Job Completion Message
**Exchange:** `job_exchange`  
**Routing Key:** `job.status`

```json
{
  "id": "550e8400-e29b-41d4-a716-446655440000",
  "status": "completed",
  "result": "Successfully processed 1,234 records",
  "updated_at": "2025-07-20T10:35:30Z",
  "completed_at": "2025-07-20T10:35:30Z"
}
```

## Data Models

### Job Object
```go
type Job struct {
    ID          string    `json:"id"`
    Description string    `json:"description"`
    Status      string    `json:"status"`
    Result      string    `json:"result,omitempty"`
    CreatedAt   time.Time `json:"created_at"`
    UpdatedAt   time.Time `json:"updated_at"`
    StartedAt   *time.Time `json:"started_at,omitempty"`
    CompletedAt *time.Time `json:"completed_at,omitempty"`
}
```

### Job Status Values
- `pending`: Job created and queued for processing
- `processing`: Job is currently being processed by a worker
- `completed`: Job finished successfully
- `failed`: Job encountered an error and could not complete

### Error Response
```go
type ErrorResponse struct {
    Error   string `json:"error"`
    Code    int    `json:"code"`
    Message string `json:"message"`
}
```

## Error Handling

### HTTP Status Codes

- `200 OK`: Request successful
- `201 Created`: Resource created successfully
- `400 Bad Request`: Invalid request format or parameters
- `404 Not Found`: Resource not found
- `500 Internal Server Error`: Server error

### Error Response Format
All errors return a JSON object with error details:

```json
{
  "error": "invalid_request",
  "code": 400,
  "message": "Job description is required"
}
```

### Common Error Scenarios

#### Invalid Job Creation
**Request:**
```bash
curl -X POST http://localhost:8081/api/jobs \
  -H "Content-Type: application/json" \
  -d '{}'
```

**Response:** `400 Bad Request`
```json
{
  "error": "validation_failed",
  "code": 400,
  "message": "Job description is required"
}
```

#### Job Not Found
**Request:**
```bash
curl http://localhost:8081/api/jobs/invalid-id
```

**Response:** `404 Not Found`
```json
{
  "error": "job_not_found",
  "code": 404,
  "message": "Job with ID 'invalid-id' not found"
}
```

#### Service Unavailable
**Response:** `500 Internal Server Error`
```json
{
  "error": "service_unavailable",
  "code": 500,
  "message": "Unable to connect to message queue"
}
```

## Examples

### Complete Job Lifecycle

#### 1. Create a Job
```bash
# Create job
RESPONSE=$(curl -s -X POST http://localhost:8081/api/jobs \
  -H "Content-Type: application/json" \
  -d '{"description": "Process customer data"}')

# Extract job ID
JOB_ID=$(echo $RESPONSE | jq -r '.id')
echo "Created job: $JOB_ID"
```

#### 2. Monitor Job Progress
```bash
# Poll for status updates
while true; do
  STATUS=$(curl -s http://localhost:8081/api/jobs/$JOB_ID | jq -r '.status')
  echo "Job status: $STATUS"
  
  if [[ "$STATUS" == "completed" || "$STATUS" == "failed" ]]; then
    break
  fi
  
  sleep 2
done
```

#### 3. Get Final Result
```bash
# Get final job details
curl -s http://localhost:8081/api/jobs/$JOB_ID | jq .
```

### Batch Job Management

#### Create Multiple Jobs
```bash
#!/bin/bash
JOB_DESCRIPTIONS=(
  "Process user registrations"
  "Generate monthly reports"
  "Cleanup temporary files"
  "Backup user data"
)

for desc in "${JOB_DESCRIPTIONS[@]}"; do
  curl -s -X POST http://localhost:8081/api/jobs \
    -H "Content-Type: application/json" \
    -d "{\"description\": \"$desc\"}" | jq '.id'
done
```

#### Monitor All Jobs
```bash
# Get current job queue status
curl -s "http://localhost:8081/api/jobs?limit=100" | \
  jq '.jobs[] | {id: .id, status: .status, description: .description}'
```

### Integration with Frontend

#### Submit Job via Web Form
```bash
# Simulate form submission
curl -X POST http://localhost:8080/submit \
  -H "Content-Type: application/x-www-form-urlencoded" \
  -d "description=Web+submitted+job" \
  -L  # Follow redirect
```

## Rate Limiting and Performance

### Current Limitations
- No rate limiting implemented
- In-memory job storage (data lost on restart)
- Single-threaded job processing per worker

### Recommended Usage
- **API calls**: No specific limits, but avoid excessive polling
- **Job creation**: Reasonable rate for demo purposes
- **Concurrent jobs**: Limited by available Job Runner instances

### Performance Characteristics
- **Job creation**: < 10ms typical response time
- **Job processing**: 5-60 seconds (simulated work)
- **Status queries**: < 5ms typical response time

## Testing the API

### Using curl
```bash
# Health check
curl http://localhost:8081/health

# Create and monitor job
JOB_ID=$(curl -s -X POST http://localhost:8081/api/jobs \
  -H "Content-Type: application/json" \
  -d '{"description": "Test job"}' | jq -r '.id')

# Monitor until completion
watch curl -s http://localhost:8081/api/jobs/$JOB_ID
```

### Using httpie
```bash
# Install httpie: pip install httpie

# Create job
http POST localhost:8081/api/jobs description="Test job"

# Get job status
http GET localhost:8081/api/jobs/550e8400-e29b-41d4-a716-446655440000

# List jobs
http GET localhost:8081/api/jobs limit==10 status==pending
```

### Using Postman
Import the following collection for easy API testing:

```json
{
  "info": {
    "name": "Microservices Demo API",
    "description": "Collection for testing the microservices demo API"
  },
  "item": [
    {
      "name": "Create Job",
      "request": {
        "method": "POST",
        "header": [{"key": "Content-Type", "value": "application/json"}],
        "body": {
          "mode": "raw",
          "raw": "{\"description\": \"Test job from Postman\"}"
        },
        "url": {
          "raw": "http://localhost:8081/api/jobs"
        }
      }
    }
  ]
}
```

---

For more information about the underlying implementation, see:
- [Architecture Documentation](ARCHITECTURE.md)
- [Deployment Guide](DEPLOYMENT.md)
- [Development Setup](DEVELOPMENT.md)
