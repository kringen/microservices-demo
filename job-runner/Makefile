# Job Runner Makefile

.PHONY: build test clean run docker-build docker-run

SERVICE_NAME := job-runner
DOCKER_IMAGE := microservices-$(SERVICE_NAME)
DOCKER_TAG := latest

# Build the service
build:
	@echo "Building $(SERVICE_NAME)..."
	go build -o bin/$(SERVICE_NAME) .

# Run tests
test:
	@echo "Running tests for $(SERVICE_NAME)..."
	go test -v .

# Run tests with coverage
test-coverage:
	@echo "Running tests with coverage for $(SERVICE_NAME)..."
	go test -v -race -coverprofile=coverage.out .
	go tool cover -html=coverage.out -o coverage.html

# Run tests for CI
test-ci:
	@echo "Running CI tests for $(SERVICE_NAME)..."
	go test -v -race -coverprofile=coverage-$(SERVICE_NAME).out .

# Clean build artifacts
clean:
	rm -f bin/$(SERVICE_NAME)

# Run the service locally
run:
	@echo "Starting $(SERVICE_NAME)..."
	go run main.go

# Run with environment variables
run-env:
	@echo "Starting $(SERVICE_NAME) with custom environment..."
	RABBITMQ_URL=${RABBITMQ_URL} go run main.go

# Docker build
docker-build:
	@echo "Building Docker image for $(SERVICE_NAME)..."
	cd .. && docker build -f $(SERVICE_NAME)/Dockerfile -t $(DOCKER_IMAGE):$(DOCKER_TAG) .

# Docker run
docker-run: docker-build
	@echo "Running $(SERVICE_NAME) in Docker..."
	docker run --rm \
		-e RABBITMQ_URL=amqp://guest:guest@host.docker.internal:5672/ \
		$(DOCKER_IMAGE):$(DOCKER_TAG)

# Docker run with network
docker-run-network: docker-build
	@echo "Running $(SERVICE_NAME) in Docker with network..."
	docker run --rm \
		--network microservices-demo_microservices-net \
		-e RABBITMQ_URL=amqp://guest:guest@rabbitmq:5672/ \
		$(DOCKER_IMAGE):$(DOCKER_TAG)

# Run multiple instances (for load testing)
run-multiple:
	@echo "Starting 3 instances of $(SERVICE_NAME)..."
	RABBITMQ_URL=${RABBITMQ_URL} go run main.go &
	RABBITMQ_URL=${RABBITMQ_URL} go run main.go &
	RABBITMQ_URL=${RABBITMQ_URL} go run main.go &
	@echo "Started 3 job runner instances"

# Development mode with hot reload (requires air)
dev:
	@if ! command -v air > /dev/null; then \
		echo "Installing air for hot reload..."; \
		go install github.com/cosmtrek/air@latest; \
	fi
	air

# Format code
fmt:
	go fmt .

# Lint code
lint:
	@if command -v golangci-lint > /dev/null; then \
		golangci-lint run; \
	else \
		echo "golangci-lint not found. Install it with:"; \
		echo "go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest"; \
	fi

# Show help
help:
	@echo "Available commands for $(SERVICE_NAME):"
	@echo "  make build          - Build the service binary"
	@echo "  make test           - Run tests"
	@echo "  make test-coverage  - Run tests with coverage"
	@echo "  make test-ci        - Run tests for CI"
	@echo "  make run            - Run service locally"
	@echo "  make run-env        - Run with custom environment"
	@echo "  make run-multiple   - Run multiple instances"
	@echo "  make docker-build   - Build Docker image"
	@echo "  make docker-run     - Run in Docker (standalone)"
	@echo "  make docker-run-network - Run in Docker (with network)"
	@echo "  make dev            - Run with hot reload"
	@echo "  make clean          - Clean build artifacts"
	@echo "  make fmt            - Format code"
	@echo "  make lint           - Run linter"
	@echo "  make help           - Show this help"
