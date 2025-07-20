# Shared Module

The Shared module contains common data structures, utilities, and RabbitMQ client functionality used across all microservices. This promotes code reuse and ensures consistency in data models and communication protocols.

## 🎯 Purpose

- **Data Models**: Common structs for jobs, requests, and results
- **RabbitMQ Client**: Unified message queue communication
- **Type Safety**: Consistent data types across services
- **Code Reuse**: Shared utilities and helper functions

## 🏗️ Module Structure

```
shared/
├── types.go          # Data models and enums
├── rabbitmq.go       # RabbitMQ client implementation
├── types_test.go     # Data model tests
├── rabbitmq_test.go  # RabbitMQ client tests
└── README.md         # This file
```

## 📊 Data Models

### Job Status Enum
```go
type JobStatus string

const (
    JobStatusPending    JobStatus = "pending"
    JobStatusProcessing JobStatus = "processing"
    JobStatusCompleted  JobStatus = "completed"
    JobStatusFailed     JobStatus = "failed"
)
```

### Core Data Structures

#### Job Request
```go
type JobRequest struct {
    Title       string `json:"title" binding:"required"`
    Description string `json:"description"`
}
```

#### Job Entity
```go
type Job struct {
    ID          string     `json:"id"`
    Title       string     `json:"title"`
    Description string     `json:"description"`
    Status      JobStatus  `json:"status"`
    CreatedAt   time.Time  `json:"created_at"`
    StartedAt   *time.Time `json:"started_at,omitempty"`
    CompletedAt *time.Time `json:"completed_at,omitempty"`
    Result      string     `json:"result,omitempty"`
    Error       string     `json:"error,omitempty"`
}
```

#### Job Message (Queue)
```go
type JobMessage struct {
    JobID       string `json:"job_id"`
    Title       string `json:"title"`
    Description string `json:"description"`
}
```

#### Job Result (Queue)
```go
type JobResult struct {
    JobID       string    `json:"job_id"`
    Status      JobStatus `json:"status"`
    Result      string    `json:"result,omitempty"`
    Error       string    `json:"error,omitempty"`
    CompletedAt time.Time `json:"completed_at"`
}
```

## 🐰 RabbitMQ Client

### Client Features
- **Connection Management**: Automatic reconnection and error handling
- **Queue Management**: Dynamic queue creation and configuration
- **Message Publishing**: Reliable job and result publishing
- **Message Consumption**: Efficient job and result consumption
- **Error Handling**: Robust error recovery and logging

### Client Interface
```go
type RabbitMQClient struct {
    conn    *amqp.Connection
    channel *amqp.Channel
    url     string
}

// Core methods
func NewRabbitMQClient(url string) (*RabbitMQClient, error)
func (r *RabbitMQClient) PublishJob(job JobMessage) error
func (r *RabbitMQClient) PublishResult(result JobResult) error
func (r *RabbitMQClient) ConsumeJobs() (<-chan amqp.Delivery, error)
func (r *RabbitMQClient) ConsumeResults() (<-chan amqp.Delivery, error)
func (r *RabbitMQClient) Close() error
```

### Queue Configuration
```go
// Job queue for pending work
const JobQueueName = "microservices.jobs"

// Result queue for completed work
const ResultQueueName = "microservices.results"

// Queue properties
QueueDurable    = true   // Survive server restart
QueueAutoDelete = false  // Don't delete when unused
QueueExclusive  = false  // Allow multiple consumers
```

## 🔄 Message Flow

### Job Publishing (API Server → Job Runner)
```go
// Create job message
jobMessage := shared.JobMessage{
    JobID:       job.ID,
    Title:       job.Title,
    Description: job.Description,
}

// Publish to jobs queue
err := rabbitmqClient.PublishJob(jobMessage)
```

### Result Publishing (Job Runner → API Server)
```go
// Create result message
result := shared.JobResult{
    JobID:       jobMessage.JobID,
    Status:      shared.JobStatusCompleted,
    Result:      "Job completed successfully",
    CompletedAt: time.Now(),
}

// Publish to results queue
err := rabbitmqClient.PublishResult(result)
```

### Job Consumption (Job Runner)
```go
// Start consuming jobs
jobs, err := rabbitmqClient.ConsumeJobs()
if err != nil {
    return err
}

// Process each job
for delivery := range jobs {
    var jobMessage shared.JobMessage
    json.Unmarshal(delivery.Body, &jobMessage)
    
    // Process job...
    
    delivery.Ack(false)
}
```

### Result Consumption (API Server)
```go
// Start consuming results
results, err := rabbitmqClient.ConsumeResults()
if err != nil {
    return err
}

// Process each result
for delivery := range results {
    var result shared.JobResult
    json.Unmarshal(delivery.Body, &result)
    
    // Update job status...
}
```

## 🧪 Testing

### Unit Tests
```bash
# Run all shared module tests
go test -v .

# Test specific components
go test -run TestJobStatus
go test -run TestRabbitMQClient
```

### Integration Tests
```bash
# Test with real RabbitMQ instance
go test -v . -tags=integration

# Test message publishing/consuming
go test -run TestMessageFlow
```

### Test Coverage
```bash
# Generate coverage report
go test -coverprofile=coverage.out .
go tool cover -html=coverage.out
```

## 🔧 Development

### Prerequisites
- Go 1.21+
- RabbitMQ server (for integration tests)
- AMQP Go library

### Building
```bash
# The shared module is imported by other services
# No standalone build required
```

### Code Quality
```bash
# Format code
go fmt .

# Run linter
golangci-lint run
```

## 📈 Usage Patterns

### Importing in Services
```go
import "microservices-demo/shared"

// Use data types
job := &shared.Job{
    Title:  "Example Job",
    Status: shared.JobStatusPending,
}

// Use RabbitMQ client
client, err := shared.NewRabbitMQClient(rabbitmqURL)
defer client.Close()
```

### Error Handling
```go
// Connection errors
client, err := shared.NewRabbitMQClient(url)
if err != nil {
    log.Fatalf("Failed to connect to RabbitMQ: %v", err)
}

// Publishing errors
if err := client.PublishJob(jobMessage); err != nil {
    log.Printf("Failed to publish job: %v", err)
    // Handle error appropriately
}
```

## 🔍 Monitoring

### Connection Health
```go
// Check if connection is closed
if client.IsConnectionClosed() {
    // Reconnect or handle error
}
```

### Queue Statistics
- Monitor queue depth
- Track message rates
- Watch for failed deliveries
- Alert on connection issues

## ⚡ Performance

### Best Practices
- **Connection Reuse**: Share client instances across goroutines
- **Channel Management**: Use separate channels for publishing/consuming
- **Message Acknowledgment**: Properly ack/nack messages
- **Error Recovery**: Implement exponential backoff for reconnection

### Optimization Tips
```go
// Prefetch count for consumers
err := channel.Qos(
    1,     // prefetch count
    0,     // prefetch size
    false, // global
)

// Publish with confirmation
if err := channel.Confirm(false); err != nil {
    return err
}
```

## 🛠️ Troubleshooting

### Common Issues

**Connection Refused**
```bash
# Check RabbitMQ server status
rabbitmq-diagnostics ping

# Verify connection URL
echo $RABBITMQ_URL
```

**Message Not Delivered**
```bash
# Check queue existence
rabbitmqctl list_queues

# Verify queue bindings
rabbitmqctl list_bindings
```

**Memory Leaks**
```bash
# Monitor connection count
rabbitmqctl list_connections

# Check for unclosed channels
rabbitmqctl list_channels
```

## 🔐 Security

### Connection Security
```bash
# Use TLS for production
RABBITMQ_URL="amqps://user:pass@host:5671/"

# Secure credentials
# Store in environment variables
# Use connection URI parsing
```

### Message Security
- Validate message content
- Implement message signing
- Use secure serialization
- Monitor for malicious content

## 📚 Dependencies

### Required Packages
```go
import (
    "encoding/json"
    "log"
    "time"
    
    amqp "github.com/rabbitmq/amqp091-go"
)
```

### Version Compatibility
- Go 1.21+
- RabbitMQ 3.8+
- AMQP 0.9.1 protocol

## 🎯 Extension Points

### Adding New Job Types
```go
// Extend JobMessage with type field
type JobMessage struct {
    JobID       string    `json:"job_id"`
    Type        string    `json:"type"`
    Title       string    `json:"title"`
    Description string    `json:"description"`
    Payload     map[string]interface{} `json:"payload,omitempty"`
}
```

### Custom Result Data
```go
// Extend JobResult with metadata
type JobResult struct {
    JobID       string                 `json:"job_id"`
    Status      JobStatus              `json:"status"`
    Result      string                 `json:"result,omitempty"`
    Metadata    map[string]interface{} `json:"metadata,omitempty"`
    CompletedAt time.Time              `json:"completed_at"`
}
```

## 📈 Future Enhancements

- **Dead Letter Queues**: Handle failed message processing
- **Message TTL**: Implement message expiration
- **Priority Queues**: Support job prioritization
- **Routing Keys**: Enable message routing patterns
- **Metrics Collection**: Built-in performance metrics
