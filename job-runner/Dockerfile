# Build stage
FROM golang:1.21-alpine AS builder

WORKDIR /app

# Copy go mod files
COPY go.mod go.sum ./
COPY shared/ ./shared/

# Download dependencies
RUN go mod download

# Copy source code
COPY job-runner/ ./job-runner/

# Build the application
RUN cd job-runner && CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main .

# Final stage
FROM alpine:latest

# Install ca-certificates for HTTPS calls, curl for health checks
RUN apk --no-cache add ca-certificates tzdata curl

# Create a non-root user
RUN addgroup -g 1000 appgroup && \
    adduser -D -s /bin/sh -u 1000 -G appgroup appuser

# Copy the binary from builder stage to a standard location
COPY --from=builder /app/job-runner/main /usr/local/bin/main

# Make the binary executable (this should work in /usr/local/bin)
RUN chmod +x /usr/local/bin/main

# Switch to non-root user
USER appuser

# Set working directory
WORKDIR /home/appuser

# Set environment variables for AI Research Agent
ENV RABBITMQ_URL=amqp://guest:guest@rabbitmq:5672/
ENV OLLAMA_URL=http://ollama:11434
ENV OLLAMA_MODEL=llama3.2
ENV DAPR_HTTP_ENDPOINT=http://localhost:3500

# Health check to verify Ollama connectivity
HEALTHCHECK --interval=30s --timeout=10s --start-period=60s --retries=3 \
  CMD curl -f $OLLAMA_URL/api/tags || exit 1

# Run the AI research agent
CMD ["/usr/local/bin/main"]
