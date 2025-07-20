# Build stage
FROM golang:1.21-alpine AS builder

WORKDIR /app

# Copy go mod files
COPY go.mod go.sum ./
COPY shared/ ./shared/

# Download dependencies
RUN go mod download

# Copy source code
COPY frontend/ ./frontend/

# Build the application
RUN cd frontend && CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main .

# Final stage
FROM alpine:latest

# Install ca-certificates for HTTPS calls
RUN apk --no-cache add ca-certificates tzdata

# Create a non-root user
RUN addgroup -g 1000 appgroup && \
    adduser -D -s /bin/sh -u 1000 -G appgroup appuser

# Copy the binary from builder stage to a standard location
COPY --from=builder /app/frontend/main /usr/local/bin/main

# Make the binary executable (this should work in /usr/local/bin)
RUN chmod +x /usr/local/bin/main

# Switch to non-root user
USER appuser

# Set working directory
WORKDIR /home/appuser

# Expose port
EXPOSE 8080

# Add health check
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
    CMD wget --no-verbose --tries=1 --spider http://localhost:8080/ || exit 1

# Set environment variables
ENV GIN_MODE=release
ENV API_SERVER_URL=http://api-server:8081

# Run the binary
CMD ["/usr/local/bin/main"]
