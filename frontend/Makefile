# Frontend Makefile

.PHONY: build test clean run docker-build docker-run

SERVICE_NAME := frontend
DOCKER_IMAGE := microservices-$(SERVICE_NAME)
DOCKER_TAG := latest
PORT := 8080

# Build the service
build:
	@echo "Building $(SERVICE_NAME)..."
	go build -o bin/$(SERVICE_NAME) .

# Run tests
test:
	@echo "Running tests for $(SERVICE_NAME)..."
	go test -v .

# Clean build artifacts
clean:
	rm -f bin/$(SERVICE_NAME)

# Run the service locally
run:
	@echo "Starting $(SERVICE_NAME) on port $(PORT)..."
	go run main.go templates.go

# Run with environment variables
run-env:
	@echo "Starting $(SERVICE_NAME) with custom environment..."
	API_SERVER_URL=${API_SERVER_URL} GIN_MODE=debug go run main.go templates.go

# Docker build
docker-build:
	@echo "Building Docker image for $(SERVICE_NAME)..."
	cd .. && docker build -f $(SERVICE_NAME)/Dockerfile -t $(DOCKER_IMAGE):$(DOCKER_TAG) .

# Docker run
docker-run: docker-build
	@echo "Running $(SERVICE_NAME) in Docker..."
	docker run --rm -p $(PORT):$(PORT) \
		-e API_SERVER_URL=http://host.docker.internal:8081 \
		$(DOCKER_IMAGE):$(DOCKER_TAG)

# Docker run with network
docker-run-network: docker-build
	@echo "Running $(SERVICE_NAME) in Docker with network..."
	docker run --rm -p $(PORT):$(PORT) \
		--network microservices-demo_microservices-net \
		-e API_SERVER_URL=http://api-server:8081 \
		$(DOCKER_IMAGE):$(DOCKER_TAG)

# Check if service is healthy
health-check:
	@echo "Checking $(SERVICE_NAME) health..."
	curl -f http://localhost:$(PORT)/ || echo "Service not available"

# Open frontend in browser
open:
	@echo "Opening $(SERVICE_NAME) in browser..."
	@if command -v xdg-open > /dev/null; then \
		xdg-open http://localhost:$(PORT); \
	elif command -v open > /dev/null; then \
		open http://localhost:$(PORT); \
	else \
		echo "Please open http://localhost:$(PORT) in your browser"; \
	fi

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
	@echo "  make run            - Run service locally"
	@echo "  make run-env        - Run with custom environment"
	@echo "  make docker-build   - Build Docker image"
	@echo "  make docker-run     - Run in Docker (standalone)"
	@echo "  make docker-run-network - Run in Docker (with network)"
	@echo "  make health-check   - Check service health"
	@echo "  make open           - Open frontend in browser"
	@echo "  make dev            - Run with hot reload"
	@echo "  make clean          - Clean build artifacts"
	@echo "  make fmt            - Format code"
	@echo "  make lint           - Run linter"
	@echo "  make help           - Show this help"
