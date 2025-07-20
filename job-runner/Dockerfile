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

# Install ca-certificates for HTTPS calls
RUN apk --no-cache add ca-certificates tzdata

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

# Set environment variables
ENV RABBITMQ_URL=amqp://guest:guest@rabbitmq:5672/

# Run the binary
CMD ["/usr/local/bin/main"]
