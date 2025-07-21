# Development Guide

This guide provides comprehensive information for developers working on the Microservices Demo application.

## Table of Contents

- [Getting Started](#getting-started)
- [Development Environment](#development-environment)
- [Code Organization](#code-organization)
- [Development Workflow](#development-workflow)
- [Testing Strategy](#testing-strategy)
- [Debugging](#debugging)
- [Best Practices](#best-practices)
- [Contributing](#contributing)

## Getting Started

### Prerequisites

#### Required Tools
```bash
# Core development tools
go version          # Go 1.21 or later
docker --version    # Docker for containerization
git --version       # Git for version control
make --version      # Make for build automation

# Optional but recommended
golangci-lint --version  # Code linting
air --version           # Hot reloading (go install github.com/cosmtrek/air@latest)
```

#### IDE Setup

**VS Code (Recommended)**
```json
// .vscode/settings.json
{
  "go.useLanguageServer": true,
  "go.formatTool": "goimports",
  "go.lintTool": "golangci-lint",
  "go.testFlags": ["-v", "-race"],
  "go.buildTags": "integration",
  "files.exclude": {
    "**/bin": true,
    "**/.git": true
  }
}
```

**Extensions:**
- Go (official Google extension)
- Docker
- YAML
- GitLens
- REST Client (for API testing)

### Quick Setup

```bash
# 1. Clone the repository
git clone https://github.com/kringen/homelab.git
cd microservices-demo

# 2. Install dependencies
make deps

# 3. Start development environment
make docker-up

# 4. Verify setup
make health-check
```

## Development Environment

### Local Development Stack

#### Option 1: Full Docker Environment (Recommended for beginners)
```bash
# Start everything with Docker
make docker-up

# View logs
make docker-logs

# Stop everything
make docker-down
```

#### Option 2: Hybrid Development (Recommended for active development)
```bash
# Start dependencies in Docker
make rabbitmq-up

# Run services locally for development
# Terminal 1
cd api-server && make dev

# Terminal 2
cd job-runner && make dev

# Terminal 3
cd frontend && make dev
```

#### Option 3: Full Local Development
```bash
# Install and start RabbitMQ locally
# macOS: brew install rabbitmq && brew services start rabbitmq
# Linux: sudo apt-get install rabbitmq-server

# Run services
make run-all-background
```

### Environment Variables

#### Development Configuration
```bash
# .env (create in root directory)
RABBITMQ_URL=amqp://guest:guest@localhost:5672/
API_SERVER_URL=http://localhost:8081
FRONTEND_URL=http://localhost:8080
GO_ENV=development
LOG_LEVEL=debug
```

#### Service-Specific Variables
```bash
# API Server
API_SERVER_PORT=8081
API_SERVER_HOST=localhost

# Frontend
FRONTEND_PORT=8080
FRONTEND_HOST=localhost

# Job Runner
JOB_RUNNER_WORKERS=3
JOB_TIMEOUT=60s
```

### Hot Reloading Setup

#### Install Air (Go Hot Reloader)
```bash
go install github.com/cosmtrek/air@latest
```

#### Configuration Files
```toml
# api-server/.air.toml
root = "."
testdata_dir = "testdata"
tmp_dir = "tmp"

[build]
args_bin = []
bin = "./tmp/main"
cmd = "go build -o ./tmp/main ."
delay = 1000
exclude_dir = ["assets", "tmp", "vendor", "testdata"]
exclude_file = []
exclude_regex = ["_test.go"]
exclude_unchanged = false
follow_symlink = false
full_bin = ""
include_dir = []
include_ext = ["go", "tpl", "tmpl", "html"]
include_file = []
kill_delay = "0s"
log = "build-errors.log"
send_interrupt = false
stop_on_root = false
```

## Code Organization

### Project Structure
```
microservices-demo/
├── api-server/              # API Server service
│   ├── handlers/            # HTTP handlers
│   ├── models/              # Data models
│   ├── services/            # Business logic
│   ├── main.go              # Entry point
│   └── main_test.go         # Integration tests
├── frontend/                # Frontend service
│   ├── handlers/            # HTTP handlers
│   ├── templates/           # HTML templates
│   ├── static/              # Static assets
│   └── main.go              # Entry point
├── job-runner/              # Job Runner service
│   ├── processor/           # Job processing logic
│   ├── consumer/            # Message queue consumer
│   └── main.go              # Entry point
├── shared/                  # Shared libraries
│   ├── types.go             # Common data types
│   ├── rabbitmq.go          # RabbitMQ utilities
│   └── config.go            # Configuration helpers
├── scripts/                 # Development scripts
├── docs/                    # Documentation
└── k8s/                     # Kubernetes manifests
```

### Package Guidelines

#### Import Organization
```go
package main

import (
    // Standard library imports
    "context"
    "fmt"
    "net/http"
    
    // Third-party imports
    "github.com/gin-gonic/gin"
    "github.com/streadway/amqp"
    
    // Local imports
    "github.com/kringen/microservices-demo/shared"
)
```

#### Package Naming
- **Lowercase**: Package names should be lowercase
- **Descriptive**: Clear purpose (handlers, models, services)
- **Singular**: Use singular nouns (handler, not handlers in package name)

### Code Style

#### Go Formatting
```bash
# Format all code
make fmt

# Check formatting
gofmt -d .

# Imports formatting
goimports -w .
```

#### Naming Conventions
```go
// Good naming examples
type JobProcessor interface {
    ProcessJob(ctx context.Context, job Job) error
}

type APIServer struct {
    jobService *JobService
    config     Config
}

func (s *APIServer) CreateJob(c *gin.Context) {
    // Implementation
}

// Constants
const (
    JobStatusPending    = "pending"
    JobStatusProcessing = "processing"
    JobStatusCompleted  = "completed"
    JobStatusFailed     = "failed"
)
```

## Development Workflow

### Feature Development

#### 1. Branch Strategy
```bash
# Create feature branch
git checkout -b feature/add-job-priority

# Make changes and commit frequently
git add .
git commit -m "Add job priority field to model"

# Push and create PR
git push origin feature/add-job-priority
```

#### 2. Development Cycle
```bash
# 1. Write/modify code
vim api-server/models/job.go

# 2. Run tests
make test

# 3. Check code quality
make lint

# 4. Test locally
make run-all-background
./scripts/demo.sh

# 5. Commit changes
git add . && git commit -m "Descriptive commit message"
```

#### 3. Pre-commit Checklist
```bash
# Automated checks (matches CI)
./scripts/pre-commit.sh

# Manual verification
make docker-build    # Ensure Docker builds work
make test-coverage   # Check test coverage
```

### Testing During Development

#### Unit Testing
```bash
# Test specific service
cd api-server && go test -v ./...

# Test with coverage
cd api-server && go test -v -coverprofile=coverage.out ./...

# View coverage
go tool cover -html=coverage.out
```

#### Integration Testing
```bash
# Start dependencies
make rabbitmq-up

# Run integration tests
go test -v -tags=integration ./...

# Or use make target
make test-integration
```

#### Manual Testing
```bash
# Start services
make run-all-background

# Test API endpoints
curl -X POST http://localhost:8081/api/jobs \
  -H "Content-Type: application/json" \
  -d '{"description": "Test job"}'

# Test frontend
open http://localhost:8080

# Stop services
make stop-all
```

### Code Generation

#### Mock Generation (if using mockery)
```bash
# Install mockery
go install github.com/vektra/mockery/v2@latest

# Generate mocks
mockery --name=JobProcessor --dir=./shared --output=./mocks
```

#### API Documentation Generation
```bash
# If using swag for Swagger docs
swag init -g api-server/main.go -o docs/swagger
```

## Testing Strategy

### Test Categories

#### 1. Unit Tests
**Purpose:** Test individual components in isolation

```go
// Example unit test
func TestJobService_CreateJob(t *testing.T) {
    // Arrange
    mockStore := &MockJobStore{}
    mockPublisher := &MockMessagePublisher{}
    service := NewJobService(mockStore, mockPublisher)
    
    jobReq := CreateJobRequest{
        Description: "Test job",
    }
    
    // Act
    job, err := service.CreateJob(context.Background(), jobReq)
    
    // Assert
    assert.NoError(t, err)
    assert.NotEmpty(t, job.ID)
    assert.Equal(t, "Test job", job.Description)
    assert.Equal(t, JobStatusPending, job.Status)
}
```

#### 2. Integration Tests
**Purpose:** Test service interactions with real dependencies

```go
// +build integration

func TestJobWorkflow_Integration(t *testing.T) {
    // Setup real RabbitMQ connection
    conn, err := amqp.Dial(os.Getenv("RABBITMQ_URL"))
    require.NoError(t, err)
    defer conn.Close()
    
    // Test complete job workflow
    // 1. Create job via API
    // 2. Verify job is queued
    // 3. Process job
    // 4. Verify completion
}
```

#### 3. End-to-End Tests
**Purpose:** Test complete user workflows

```go
func TestCompleteJobWorkflow_E2E(t *testing.T) {
    // Start all services
    // Submit job via frontend
    // Wait for completion
    // Verify results
}
```

### Test Data Management

#### Test Fixtures
```go
// testdata/fixtures.go
package testdata

var SampleJobs = []shared.Job{
    {
        ID:          "test-job-1",
        Description: "Sample job 1",
        Status:      shared.JobStatusPending,
        CreatedAt:   time.Now(),
    },
    // More fixtures...
}
```

#### Test Helpers
```go
// testutils/helpers.go
package testutils

func SetupTestRabbitMQ(t *testing.T) *amqp.Connection {
    conn, err := amqp.Dial(os.Getenv("RABBITMQ_URL"))
    require.NoError(t, err)
    
    t.Cleanup(func() {
        conn.Close()
    })
    
    return conn
}
```

### Test Coverage

#### Coverage Goals
- **Unit Tests**: > 80% coverage
- **Integration Tests**: > 60% coverage
- **Critical Paths**: 100% coverage

#### Measuring Coverage
```bash
# Generate coverage report
make test-coverage

# View in browser
go tool cover -html=coverage.out

# Coverage by package
go tool cover -func=coverage.out
```

## Debugging

### Local Debugging

#### Delve Debugger
```bash
# Install delve
go install github.com/go-delve/delve/cmd/dlv@latest

# Debug API server
cd api-server
dlv debug . -- --port=8081

# Debug with breakpoints
(dlv) break main.main
(dlv) continue
```

#### IDE Debugging
**VS Code launch.json:**
```json
{
    "version": "0.2.0",
    "configurations": [
        {
            "name": "Debug API Server",
            "type": "go",
            "request": "launch",
            "mode": "debug",
            "program": "${workspaceFolder}/api-server",
            "env": {
                "RABBITMQ_URL": "amqp://guest:guest@localhost:5672/"
            },
            "args": []
        }
    ]
}
```

### Application Logging

#### Structured Logging
```go
// Use structured logging
log.WithFields(log.Fields{
    "job_id": job.ID,
    "status": job.Status,
    "duration": time.Since(startTime),
}).Info("Job processing completed")
```

#### Log Levels
```go
// Development: DEBUG level
log.SetLevel(log.DebugLevel)

// Production: INFO level
log.SetLevel(log.InfoLevel)
```

### Troubleshooting Common Issues

#### RabbitMQ Connection Issues
```bash
# Check RabbitMQ status
docker ps | grep rabbitmq

# Check RabbitMQ logs
docker logs microservices-demo_rabbitmq_1

# Access management UI
open http://localhost:15672
```

#### Port Conflicts
```bash
# Check what's using ports
lsof -i :8080
lsof -i :8081
lsof -i :5672

# Kill processes on ports
kill -9 $(lsof -t -i:8080)
```

#### Build Issues
```bash
# Clean build cache
go clean -cache
go clean -modcache

# Rebuild everything
make clean
make build
```

## Best Practices

### Code Quality

#### Error Handling
```go
// Good error handling
func (s *JobService) CreateJob(ctx context.Context, req CreateJobRequest) (*Job, error) {
    if req.Description == "" {
        return nil, fmt.Errorf("job description is required")
    }
    
    job := &Job{
        ID:          uuid.New().String(),
        Description: req.Description,
        Status:      JobStatusPending,
        CreatedAt:   time.Now(),
    }
    
    if err := s.store.Save(ctx, job); err != nil {
        return nil, fmt.Errorf("failed to save job: %w", err)
    }
    
    if err := s.publisher.Publish(ctx, JobCreatedEvent{Job: job}); err != nil {
        // Log error but don't fail the request
        log.WithError(err).Error("Failed to publish job created event")
    }
    
    return job, nil
}
```

#### Interface Design
```go
// Good interface design - small and focused
type JobStore interface {
    Save(ctx context.Context, job *Job) error
    Get(ctx context.Context, id string) (*Job, error)
    List(ctx context.Context, filters JobFilters) ([]Job, error)
}

type MessagePublisher interface {
    Publish(ctx context.Context, event Event) error
}
```

#### Context Usage
```go
// Always use context for cancellation and timeouts
func (s *JobService) ProcessJob(ctx context.Context, job *Job) error {
    // Set timeout for job processing
    ctx, cancel := context.WithTimeout(ctx, 60*time.Second)
    defer cancel()
    
    // Check context in long-running operations
    select {
    case <-ctx.Done():
        return ctx.Err()
    default:
        // Continue processing
    }
    
    return nil
}
```

### Performance

#### Memory Management
```go
// Use sync.Pool for frequently allocated objects
var jobPool = sync.Pool{
    New: func() interface{} {
        return &Job{}
    },
}

func (s *JobService) processJob() {
    job := jobPool.Get().(*Job)
    defer jobPool.Put(job)
    
    // Use job...
}
```

#### Goroutine Management
```go
// Use worker pools for bounded concurrency
func (r *JobRunner) Start(ctx context.Context) error {
    sem := make(chan struct{}, r.maxWorkers)
    
    for {
        select {
        case <-ctx.Done():
            return ctx.Err()
        case msg := <-r.messages:
            sem <- struct{}{}
            go func(msg amqp.Delivery) {
                defer func() { <-sem }()
                r.processMessage(ctx, msg)
            }(msg)
        }
    }
}
```

### Security

#### Input Validation
```go
func validateJobRequest(req CreateJobRequest) error {
    if len(req.Description) == 0 {
        return errors.New("description is required")
    }
    
    if len(req.Description) > 1000 {
        return errors.New("description too long")
    }
    
    // Sanitize input
    req.Description = html.EscapeString(req.Description)
    
    return nil
}
```

#### Configuration Management
```go
// Use environment variables for configuration
type Config struct {
    Port        int    `env:"PORT" envDefault:"8080"`
    RabbitMQURL string `env:"RABBITMQ_URL" envDefault:"amqp://localhost:5672/"`
    LogLevel    string `env:"LOG_LEVEL" envDefault:"info"`
}
```

## Contributing

### Contribution Guidelines

#### Code Review Checklist
- [ ] Tests added for new functionality
- [ ] Documentation updated
- [ ] Code follows style guidelines
- [ ] No security vulnerabilities introduced
- [ ] Performance impact considered
- [ ] Error handling implemented
- [ ] Logging added for debugging

#### Pull Request Template
```markdown
## Description
Brief description of changes

## Type of Change
- [ ] Bug fix
- [ ] New feature
- [ ] Documentation update
- [ ] Performance improvement

## Testing
- [ ] Unit tests pass
- [ ] Integration tests pass
- [ ] Manual testing completed

## Checklist
- [ ] Code follows style guidelines
- [ ] Self-review completed
- [ ] Documentation updated
```

### Development Standards

#### Commit Messages
```bash
# Good commit messages
feat: add job priority support
fix: resolve memory leak in job runner
docs: update API documentation
test: add integration tests for job workflow
refactor: extract common validation logic
```

#### Branch Naming
```bash
# Feature branches
feature/add-job-priority
feature/improve-error-handling

# Bug fix branches
fix/memory-leak-job-runner
fix/frontend-polling-issue

# Documentation branches
docs/api-documentation
docs/deployment-guide
```

---

## Related Documentation

- [API Documentation](API.md) - REST API reference
- [Architecture Documentation](ARCHITECTURE.md) - System design and patterns
- [Deployment Guide](DEPLOYMENT.md) - Production deployment
- [CI/CD Documentation](CICD.md) - Pipeline and automation
