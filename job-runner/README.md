# Job Runner

The Job Runner is a worker service that consumes jobs from RabbitMQ, processes them asynchronously, and reports completion status back through the message queue.

## 🎯 Purpose

- **Job Processing**: Execute jobs asynchronously with 5-60 second durations
- **Status Reporting**: Send real-time updates (pending → processing → completed/failed)
- **Timeout Protection**: Ensure no job exceeds 1-minute maximum
- **Load Balancing**: Support multiple instances for high throughput

## 🏗️ Architecture

```
RabbitMQ Jobs Queue → Job Runner → Job Processing → RabbitMQ Results Queue
                         ↓
                   Status Updates
                    (processing,
                     completed,
                     failed)
```

## ⚡ Features

### Job Processing
- **Realistic Simulation**: Different job types with varied processing logic
- **Random Duration**: 5-60 seconds to simulate real work
- **Timeout Protection**: Hard 1-minute limit on all jobs
- **Error Simulation**: 10% failure rate for testing error handling

### Job Types Supported
- **Data Analysis**: Customer data processing, insights generation
- **Report Generation**: Document creation with page counts
- **Email Campaigns**: Bulk email sending with bounce tracking
- **Backup Operations**: File and database backup simulation
- **Generic Processing**: Flexible operation simulation

### Scaling
- **Multiple Instances**: Run multiple job runners for load balancing
- **Queue-based Distribution**: RabbitMQ handles job distribution
- **Concurrent Processing**: Each instance processes jobs independently

## 🚀 Quick Start

### Local Development
```bash
# Run single instance
make run

# Run with custom environment
RABBITMQ_URL=amqp://user:pass@localhost:5672/ make run-env

# Run multiple instances for load testing
make run-multiple

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

# Test job processing logic
go test -v . -run TestJobProcessing
```

## ⚙️ Configuration

### Environment Variables

| Variable | Default | Description |
|----------|---------|-------------|
| `RABBITMQ_URL` | `amqp://guest:guest@localhost:5672/` | RabbitMQ connection URL |

### Example Configuration
```bash
export RABBITMQ_URL="amqp://user:password@rabbitmq-host:5672/"
```

## 🔄 Job Processing Flow

### 1. Job Consumption
```go
// Consume job from queue
jobMessage := consumeFromQueue()
```

### 2. Status Update - Processing
```go
// Send "processing" status
publishResult(JobResult{
    JobID: jobMessage.JobID,
    Status: "processing",
    CompletedAt: time.Now()
})
```

### 3. Job Execution
```go
// Process job with timeout protection
result := processJobWithTimeout(jobMessage)
```

### 4. Status Update - Completion
```go
// Send final result
publishResult(JobResult{
    JobID: jobMessage.JobID,
    Status: "completed", // or "failed"
    Result: "Job completed successfully...",
    CompletedAt: time.Now()
})
```

## 🎲 Job Simulation Details

### Data Analysis Jobs
```go
func simulateAnalysisJob() string {
    dataPoints := rand.Intn(1000000) + 10000
    insights := rand.Intn(50) + 5
    return fmt.Sprintf("Analyzed %d data points and generated %d insights", 
        dataPoints, insights)
}
```

### Email Campaign Jobs
```go
func simulateEmailProcessing() string {
    emailsSent := rand.Intn(500) + 50
    bounces := rand.Intn(emailsSent / 20)
    return fmt.Sprintf("Sent %d emails with %d bounces", emailsSent, bounces)
}
```

### Backup Jobs
```go
func simulateBackupJob() string {
    sizeMB := rand.Intn(5000) + 100
    files := rand.Intn(10000) + 1000
    return fmt.Sprintf("Backed up %d files (%.1f GB)", files, float64(sizeMB)/1024.0)
}
```

## 📊 Performance Characteristics

### Processing Times
- **Minimum**: 5 seconds
- **Maximum**: 60 seconds  
- **Timeout**: Hard limit at 60 seconds
- **Distribution**: Random uniform distribution

### Success Rates
- **Success**: 90% of jobs complete successfully
- **Failure**: 10% fail with simulated errors
- **Timeout**: Rare, only if processing logic hangs

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

## 🧪 Testing

### Unit Tests
```bash
# Run all tests
make test

# Test specific functions
go test -run TestJobProcessing
go test -run TestJobSimulation
```

### Load Testing
```bash
# Run multiple job runners
make run-multiple

# Submit multiple jobs via API server
# Each runner will process jobs concurrently
```

### Timeout Testing
```bash
# Test timeout protection
go test -run TestJobTimeout -timeout 70s
```

## 🐳 Docker

### Dockerfile Features
- **Multi-stage build** for optimized images
- **Alpine base** for minimal size
- **Non-root execution** for security
- **Resource limits** for container orchestration

### Scaling with Docker
```yaml
# docker-compose.yml example
job-runner:
  deploy:
    replicas: 3  # Run 3 instances
```

## 📈 Scaling & Performance

### Horizontal Scaling
```bash
# Run multiple instances
docker-compose up --scale job-runner=5
```

### Performance Monitoring
- Track job processing times
- Monitor queue depth
- Watch memory and CPU usage
- Alert on failed jobs

### Optimization Tips
- **Queue Prefetch**: Limit concurrent jobs per runner
- **Memory Management**: Clean up after job completion
- **Connection Pooling**: Reuse RabbitMQ connections

## 🔍 Monitoring & Observability

### Logging
```go
log.Printf("Job %s started processing", jobID)
log.Printf("Job %s completed in %v", jobID, duration)
log.Printf("Job %s failed: %v", jobID, error)
```

### Metrics to Track
- Jobs processed per second
- Average processing time
- Success/failure rates
- Queue depth and lag

## 🛠️ Troubleshooting

### Common Issues

**No Jobs Being Processed**
```bash
# Check RabbitMQ connection
# Verify queue exists
# Check API server is publishing jobs
```

**Memory Leaks**
```bash
# Monitor memory usage
docker stats microservices-job-runner
# Check for goroutine leaks
```

**Slow Processing**
```bash
# Check system resources
# Monitor job complexity
# Verify network connectivity
```

## 🔐 Security

### Best Practices
- Run as non-root user in containers
- Validate job input data
- Implement resource limits
- Secure RabbitMQ credentials

## 📚 Related Services

- **[API Server](../api-server/README.md)** - Job coordination
- **[Frontend](../frontend/README.md)** - User interface
- **[Shared](../shared/README.md)** - Common utilities

## 🎯 Use Cases

### Development
- Testing asynchronous job processing
- Demonstrating microservices patterns
- Learning message queue integration

### Production Patterns
- Background task processing
- Image/video processing
- Report generation
- Data analysis pipelines
- Batch processing systems
