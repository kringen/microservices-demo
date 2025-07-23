# AI Research Agent

The AI Research Agent is an intelligent worker service that consumes research requests from RabbitMQ, uses AI and external data sources to conduct comprehensive research, and reports detailed findings back through the message queue.

## 🤖 Purpose

- **AI-Powered Research**: Conduct intelligent research using Ollama LLM models
- **Multi-Source Data Gathering**: Integrate web search, GitHub, and file system data via MCP services
- **Comprehensive Analysis**: Generate detailed research reports with confidence scoring
- **Real-time Updates**: Send status updates throughout the research process
- **Scalable Processing**: Support multiple research agents for high throughput

## 🏗️ Architecture

```
RabbitMQ Research Queue → AI Research Agent → Ollama AI Analysis → RabbitMQ Results Queue
                             ↓                    ↑
                       MCP Data Services    AI Processing
                       (Web, GitHub, Files)  (llama3.2)
                             ↓                    ↓
                       Information Gathering → Confidence Scoring
                                              → Source Citations
                                              → Token Usage Tracking
```

## ⚡ Features

### AI-Powered Analysis
- **Ollama Integration**: Uses llama3.2 model for intelligent text analysis
- **Multi-Step Research**: Combines data gathering with AI reasoning
- **Contextual Understanding**: Processes queries with research type awareness
- **Quality Assessment**: Automatic confidence scoring for all research results

### Data Source Integration (MCP Services)
- **Web Search**: Simulated web search and content extraction
- **GitHub Analysis**: Repository search and code pattern analysis  
- **File System Access**: Local document and configuration analysis
- **Extensible Framework**: Ready for additional MCP service integrations

### Research Types Supported
- **General Research**: Broad information gathering and analysis
- **Technical Analysis**: Deep-dive technical investigations
- **Market Research**: Business and market intelligence
- **Competitive Analysis**: Competitor and industry analysis
- **Code & Development**: Software development insights
- **Data Analysis**: Statistical and data-driven research

### Advanced Processing
- **Token Usage Tracking**: Monitor AI model resource consumption
- **Source Citation**: Track and reference all information sources
- **Confidence Scoring**: Rate reliability of research findings (0.0-1.0)
- **Error Handling**: Graceful degradation when services are unavailable

## 🚀 Quick Start

### Local Development
```bash
# Run AI research agent
make run

# Run with custom Ollama server
OLLAMA_URL=http://custom-ollama:11434 make run

# Run with different AI model
OLLAMA_MODEL=llama3.1 make run

# Run with hot reload
make dev
```

### Docker (Recommended)
```bash
# Run full AI stack with docker-compose
make docker-up

# Check research agent logs
docker logs microservices-research-agent

# Test research endpoint
curl -X POST http://localhost:8081/api/jobs \
  -H "Content-Type: application/json" \
  -d '{"title":"AI Research Test","query":"What are the latest developments in AI?","research_type":"technical","mcp_services":["web"]}'
```

### Testing
```bash
# Run tests
make test

# Test AI integration
go test -v . -run TestOllamaIntegration

# Test MCP services
go test -v . -run TestMCPServices
```

## ⚙️ Configuration

### Environment Variables

| Variable | Default | Description |
|----------|---------|-------------|
| `RABBITMQ_URL` | `amqp://guest:guest@localhost:5672/` | RabbitMQ connection URL |
| `OLLAMA_URL` | `http://localhost:11434` | Ollama AI server endpoint |
| `OLLAMA_MODEL` | `llama3.2` | AI model to use for research |
| `DAPR_HTTP_ENDPOINT` | `http://localhost:3500` | Dapr service mesh endpoint |

### Example Configuration
```bash
export RABBITMQ_URL="amqp://user:password@rabbitmq-host:5672/"
export OLLAMA_URL="http://ollama-server:11434"
export OLLAMA_MODEL="llama3.2"
export DAPR_HTTP_ENDPOINT="http://dapr-sidecar:3500"
```

### AI Model Requirements
- **Model Size**: llama3.2 requires ~2GB storage
- **Memory**: Minimum 4GB RAM recommended for processing
- **Network**: Internet access for initial model download

## 🔄 AI Research Processing Flow

### 1. Research Request Consumption
```go
// Consume research request from queue
researchMessage := consumeFromQueue()
log.Printf("Received research request: %s - %s", researchMessage.JobID, researchMessage.Title)
```

### 2. Status Update - Processing
```go
// Send "processing" status
publishResult(JobResult{
    JobID: researchMessage.JobID,
    Status: "processing",
    CompletedAt: time.Now()
})
```

### 3. Data Gathering Phase
```go
// Gather information using MCP services
mcpData, sources, err := gatherInformationWithMCP(ctx, researchMessage)
// Web search, GitHub analysis, file system access
```

### 4. AI Analysis Phase
```go
// Process with Ollama AI
response, confidence, tokens, err := analyzeWithOllama(ctx, researchMessage, mcpData)
// Generate comprehensive research report using llama3.2
```

### 5. Status Update - Completion
```go
// Send final research results
publishResult(JobResult{
    JobID: researchMessage.JobID,
    Status: "completed",
    Result: response,
    Sources: sources,
    Confidence: confidence,
    TokensUsed: tokens,
    CompletedAt: time.Now()
})
```

## 🤖 AI & MCP Integration Details

### Ollama AI Integration
```go
// AI analysis with comprehensive prompting
type OllamaRequest struct {
    Model    string `json:"model"`    // llama3.2
    Prompt   string `json:"prompt"`   // Research query + context
    System   string `json:"system"`   // Research agent instructions
    Stream   bool   `json:"stream"`   // false for complete responses
}

// Generate research report
response, tokens, err := callOllama(ctx, systemPrompt, researchPrompt)
```

### MCP Service Simulations
```go
// Web search simulation
func simulateWebSearch(query string) (string, []string, error) {
    data := fmt.Sprintf(`Web Search Results for "%s":
    1. Comprehensive overview found on multiple authoritative sources
    2. Recent developments and trends identified
    3. Technical specifications and best practices documented`, query)
    
    sources := []string{
        "https://example.com/research-1",
        "https://example.com/research-2",
    }
    return data, sources, nil
}

// GitHub analysis simulation
func simulateGitHubSearch(query string) (string, []string, error) {
    data := fmt.Sprintf(`GitHub Repository Analysis for "%s":
    - 5 repositories with 1000+ stars
    - Modern architecture patterns prevalent
    - Active community contributions`, query)
    
    sources := []string{
        "https://github.com/example/repo1",
        "https://github.com/example/repo2",
    }
    return data, sources, nil
}
```

### Confidence Scoring Algorithm
```go
func calculateConfidence(response, mcpData string, mcpServiceCount int) float64 {
    baseConfidence := 0.6
    
    // Increase confidence based on response quality
    responseWords := len(strings.Fields(response))
    if responseWords > 100 { baseConfidence += 0.1 }
    if responseWords > 300 { baseConfidence += 0.1 }
    
    // Increase confidence based on data sources
    dataWords := len(strings.Fields(mcpData))
    if dataWords > 200 { baseConfidence += 0.1 }
    
    // Factor in number of MCP services used
    baseConfidence += float64(mcpServiceCount) * 0.05
    
    // Cap at 0.95 to account for inherent uncertainty
    if baseConfidence > 0.95 { baseConfidence = 0.95 }
    
    return baseConfidence
}
```

## 📊 Performance Characteristics

### AI Processing Times
- **Data Gathering**: 1-5 seconds per MCP service
- **AI Analysis**: 10-60 seconds depending on query complexity
- **Total Processing**: Typically 15-120 seconds per research request
- **Timeout Protection**: 5-minute hard limit with context cancellation

### Quality Metrics
- **Confidence Scoring**: 0.6-0.95 range based on data quality and sources
- **Token Usage**: Tracked for cost management and optimization
- **Source Attribution**: All information sources properly cited
- **Error Handling**: Graceful degradation when services unavailable

### Resource Requirements
- **Memory**: 2-4GB for Ollama model loading
- **CPU**: Moderate during AI inference
- **Network**: Required for MCP service calls and model downloads
- **Storage**: 2GB+ for AI model persistence

## 🔧 Development

### Prerequisites
- Go 1.21+
- RabbitMQ server
- **Ollama server with llama3.2 model**
- Access to shared module
- Docker and Docker Compose (for full stack)

### Building
```bash
# Build binary
make build

# Clean artifacts
make clean

# Build Docker image
make docker-build
```

### Local Development Setup
```bash
# Start dependencies (RabbitMQ + Ollama)
make docker-up

# Or start individual components
docker run -d --name ollama -p 11434:11434 ollama/ollama:latest
docker exec ollama ollama pull llama3.2

# Run research agent locally
make run
```

### Code Quality
```bash
# Format code
make fmt

# Run linter
make lint

# Security scan
make security-check
```

## 🧪 Testing

### Unit Tests
```bash
# Run all tests
make test

# Test AI integration
go test -run TestOllamaConnection
go test -run TestResearchProcessing

# Test MCP services
go test -run TestMCPServices
go test -run TestConfidenceScoring
```

### Integration Testing
```bash
# Test with real Ollama server
INTEGRATION_TEST=true go test -v ./...

# Test research pipeline end-to-end
curl -X POST http://localhost:8081/api/jobs \
  -H "Content-Type: application/json" \
  -d '{"title":"Integration Test","query":"Test AI research capabilities","research_type":"technical","mcp_services":["web","github"]}'
```

### Performance Testing
```bash
# Load test with multiple research requests
make load-test

# Monitor AI response times
go test -bench=BenchmarkAIAnalysis -benchmem

# Test with different model configurations
OLLAMA_MODEL=llama3.1 go test -run TestModelPerformance
```

## 🐳 Docker & Deployment

### Dockerfile Features
- **Multi-stage build** for optimized images
- **Alpine base** with curl for health checks
- **Non-root execution** for security
- **AI model volume mounting** for persistence

### Docker Compose Integration
```yaml
# Full AI research stack
services:
  ollama:
    image: ollama/ollama:latest
    volumes:
      - ollama_data:/root/.ollama
  
  research-agent:
    build: ./job-runner
    environment:
      - OLLAMA_URL=http://ollama:11434
      - OLLAMA_MODEL=llama3.2
    depends_on:
      - ollama
      - rabbitmq
```

### Scaling with Kubernetes
```yaml
# research-agent-deployment.yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: research-agent
spec:
  replicas: 3  # Multiple AI research agents
  template:
    spec:
      containers:
      - name: research-agent
        image: microservices-research-agent:latest
        env:
        - name: OLLAMA_URL
          value: "http://ollama-service:11434"
```

## 📈 Scaling & Performance

### Horizontal Scaling
```bash
# Scale research agents with Docker Compose
docker-compose up --scale research-agent=5

# Monitor AI processing load
docker stats microservices-research-agent
```

### Performance Optimization
- **Model Caching**: Keep Ollama model loaded in memory
- **Connection Pooling**: Reuse HTTP connections to Ollama
- **Batch Processing**: Group similar research requests
- **Resource Limits**: Configure appropriate CPU/memory limits

### Monitoring Key Metrics
- **Research requests processed per minute**
- **Average AI analysis time**
- **Token usage and costs**
- **Confidence score distributions**
- **MCP service response times**
- **Queue depth and processing lag**

### Cost Optimization
- **Model Selection**: Balance accuracy vs. speed/cost
- **Token Management**: Monitor and optimize prompt engineering
- **Resource Scheduling**: Scale agents based on demand
- **Caching Strategies**: Cache similar research results

## 🔍 Monitoring & Observability

### Comprehensive Logging
```go
log.Printf("Research %s started processing at %v", jobID, startTime)
log.Printf("Connected to Ollama server successfully")
log.Printf("MCP services initialized: %v", availableServices)
log.Printf("Research %s completed in %v with confidence %.2f", jobID, duration, confidence)
log.Printf("AI analysis used %d tokens", tokenCount)
```

### Health Checks
```bash
# Check research agent health
curl http://localhost:8080/health

# Verify Ollama connectivity
curl http://localhost:11434/api/tags

# Monitor RabbitMQ queue depth
curl -u guest:guest http://localhost:15672/api/queues
```

### AI-Specific Metrics
- **Model performance**: Response time and quality
- **Token usage patterns**: Cost tracking and optimization
- **Confidence score trends**: Research quality over time
- **MCP service availability**: External dependency health
- **Research type distributions**: Usage patterns analysis

### Alerting Scenarios
- Ollama server unavailable or slow
- High token usage indicating cost issues
- Low confidence scores suggesting quality problems
- MCP service failures affecting research quality

## 🛠️ Troubleshooting

### Common Issues

**Research Agent Not Processing Requests**
```bash
# Check Ollama server connectivity
curl http://localhost:11434/api/tags

# Verify RabbitMQ connection
docker logs microservices-research-agent | grep "RabbitMQ"

# Check if AI model is loaded
docker exec microservices-ollama ollama list
```

**AI Analysis Failing**
```bash
# Check Ollama server logs
docker logs microservices-ollama

# Verify model is downloaded
docker exec microservices-ollama ollama list | grep llama3.2

# Test direct Ollama API
curl -X POST http://localhost:11434/api/generate \
  -d '{"model":"llama3.2","prompt":"test","stream":false}'
```

**High Memory Usage**
```bash
# Monitor container resources
docker stats microservices-ollama microservices-research-agent

# Check for memory leaks in AI processing
# Ollama models require 2-4GB baseline memory
```

**Slow Research Processing**
```bash
# Check AI model response times
# Monitor MCP service simulation delays
# Verify system resources available for AI inference
# Consider using smaller/faster models for testing
```

**Model Download Issues**
```bash
# Manually pull model if auto-download fails
docker exec microservices-ollama ollama pull llama3.2

# Check available disk space (models are ~2GB)
df -h

# Verify internet connectivity for downloads
```

## 🔐 Security

### AI Security Best Practices
- **Input Validation**: Sanitize research queries to prevent prompt injection
- **Model Access Control**: Restrict access to AI endpoints
- **Token Limits**: Implement usage quotas to prevent abuse
- **Content Filtering**: Monitor and filter AI-generated content

### Infrastructure Security
- **Container Security**: Run as non-root user in containers
- **Network Policies**: Isolate AI services in secure networks
- **Secrets Management**: Secure RabbitMQ credentials and API keys
- **Resource Limits**: Prevent resource exhaustion attacks

### Data Privacy
- **Research Data**: Ensure sensitive queries are handled appropriately
- **AI Model Data**: Consider data residency for AI processing
- **Logging**: Avoid logging sensitive research content
- **Compliance**: Meet regulatory requirements for AI processing

## 📚 Related Services

- **[API Server](../api-server/README.md)** - Job coordination
- **[Frontend](../frontend/README.md)** - User interface
- **[Shared](../shared/README.md)** - Common utilities

## 🎯 Use Cases

### Development & Learning
- **AI Integration Patterns**: Learn to integrate LLMs into microservices
- **MCP Protocol**: Understand Model Context Protocol for AI agents
- **Research Workflows**: Build intelligent information gathering systems
- **Asynchronous AI**: Handle long-running AI tasks with message queues

### Production Applications
- **Intelligent Research Assistants**: Automated information gathering and analysis
- **Content Generation**: AI-powered report and document creation
- **Knowledge Management**: Organize and analyze large information datasets
- **Decision Support**: Provide AI-assisted insights for business decisions
- **Technical Documentation**: Auto-generate technical summaries and guides

### Advanced Scenarios
- **Multi-Agent Systems**: Coordinate multiple AI research agents
- **Specialized Research**: Configure domain-specific AI models
- **Real-time Intelligence**: Stream processing with AI analysis
- **Hybrid AI/Human Workflows**: Combine AI research with human review
