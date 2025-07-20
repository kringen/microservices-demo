# API Server

The API Server is a REST API service that handles job submission requests and provides job status endpoints. It acts as the central coordinator between the frontend and job processing system.

## 🎯 Purpose

- **Job Management**: Create, track, and manage job lifecycle
- **API Gateway**: Provide REST endpoints for frontend and external clients
- **Message Coordination**: Publish jobs to RabbitMQ and consume results
- **Status Tracking**: Maintain real-time job status and completion data

## 🏗️ Architecture

```
Frontend → API Server → RabbitMQ → Job Runner
     ↑         ↓           ↓           ↓
     └─────────┴───────────┴───────────┘
              Status Updates
```

## 📡 API Endpoints

### Job Management
- `POST /api/jobs` - Submit a new job
- `GET /api/jobs/{id}` - Get specific job status
- `GET /api/jobs` - List all jobs

### Health & Monitoring
- `GET /api/health` - Service health check

## 🚀 Quick Start

### Local Development
```bash
# Run with default settings
make run

# Run with custom environment
RABBITMQ_URL=amqp://user:pass@localhost:5672/ make run-env

# Run with hot reload
make dev
```

### Docker
```bash
# Build and run standalone
make docker-run

# Run with Docker network (for full stack)
make docker-run-network
```

### Testing
```bash
# Run tests
make test

# Check service health
make health-check
```

## ⚙️ Configuration

### Environment Variables

| Variable | Default | Description |
|----------|---------|-------------|
| `RABBITMQ_URL` | `amqp://guest:guest@localhost:5672/` | RabbitMQ connection URL |
| `GIN_MODE` | `debug` | Gin framework mode (debug/release) |
| `PORT` | `8081` | Server port |

### Example Configuration
```bash
export RABBITMQ_URL="amqp://user:password@rabbitmq-host:5672/"
export GIN_MODE="release"
```

## 🔄 Message Flow

### Job Submission
1. Receive job request via REST API
2. Create job record with "pending" status
3. Publish job message to RabbitMQ queue
4. Return job details to client

### Status Updates
1. Consume job results from RabbitMQ
2. Update job status in memory
3. Track processing times and completion

## 📊 Data Models

### Job Request
```json
{
  "title": "Data Analysis Task",
  "description": "Analyze customer data for monthly report"
}
```

### Job Response
```json
{
  "id": "550e8400-e29b-41d4-a716-446655440000",
  "title": "Data Analysis Task",
  "description": "Analyze customer data for monthly report",
  "status": "pending",
  "created_at": "2025-07-17T21:30:00Z",
  "started_at": null,
  "completed_at": null,
  "result": null,
  "error": null
}
```

## 🔧 Development

### Prerequisites
- Go 1.21+
- RabbitMQ server
- Access to shared module

### Building
```bash
# Build binary
make build

# Clean artifacts
make clean
```

### Code Quality
```bash
# Format code
make fmt

# Run linter
make lint
```

## 🐳 Docker

### Dockerfile Features
- **Multi-stage build** for smaller images
- **Health checks** for container orchestration
- **Non-root user** for security
- **Alpine base** for minimal size

### Building
```bash
# Build image
make docker-build

# Run container
make docker-run
```

## 🧪 Testing

### API Testing
```bash
# Create a job
curl -X POST http://localhost:8081/api/jobs \
  -H "Content-Type: application/json" \
  -d '{"title": "Test Job", "description": "Test description"}'

# Get job status
curl http://localhost:8081/api/jobs/{job-id}

# List all jobs
curl http://localhost:8081/api/jobs

# Health check
curl http://localhost:8081/api/health
```

### Load Testing
```bash
# Multiple concurrent requests
for i in {1..10}; do
  curl -X POST http://localhost:8081/api/jobs \
    -H "Content-Type: application/json" \
    -d "{\"title\": \"Load Test $i\", \"description\": \"Test job $i\"}" &
done
```

## 🔍 Monitoring

### Health Checks
- HTTP endpoint: `/api/health`
- RabbitMQ connection status
- Service uptime and version

### Logging
- Structured logging with timestamps
- Job lifecycle events
- Error tracking and debugging

## 📈 Scaling

### Horizontal Scaling
- Stateless design allows multiple instances
- Shared job state via RabbitMQ
- Load balancer friendly

### Performance Tips
- Use connection pooling for RabbitMQ
- Implement caching for frequently accessed jobs
- Monitor memory usage for job storage

## 🔐 Security

### Best Practices
- Input validation on all endpoints
- CORS configuration for frontend
- Health check endpoint security
- Environment variable configuration

## 🛠️ Troubleshooting

### Common Issues

**Connection Refused**
```bash
# Check if RabbitMQ is running
make health-check
```

**Port Already in Use**
```bash
# Find process using port 8081
lsof -i :8081
# Kill process if needed
```

**Memory Issues**
- Monitor job storage size
- Implement job cleanup for old jobs
- Use pagination for job listing

## 📚 Related Services

- **[Frontend](../frontend/README.md)** - Web interface
- **[Job Runner](../job-runner/README.md)** - Job processing
- **[Shared](../shared/README.md)** - Common utilities
