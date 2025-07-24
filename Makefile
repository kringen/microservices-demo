.PHONY: all build test clean run run-api run-job run-frontend stop-all deps rabbitmq-up rabbitmq-down docker-build docker-up docker-down docker-build-push docker-push docker-build-push-tag docker-clean docker-clean-all

# Default target
all: deps test build

# Install dependencies
deps:
	go mod tidy
	go mod download

# Build all services
build:
	@echo "Building all services..."
	cd api-server && make build
	cd job-runner && make build
	cd frontend && make build
	@echo "Build complete! Binaries are in each service's bin/ directory"

# Run tests for all services
test:
	@echo "Running tests..."
	go test ./...
	@echo "All tests passed!"

# Run tests with coverage
test-coverage:
	@echo "Running tests with coverage..."
	go test -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

# Run tests for CI (includes individual service tests)
test-ci:
	@echo "Running CI tests..."
	@echo "Testing root module..."
	go test -v -race -coverprofile=coverage.out ./...
	@echo "Testing API server..."
	cd api-server && go test -v -race -coverprofile=coverage-api.out ./...
	@echo "Testing frontend..."
	cd frontend && go test -v -race -coverprofile=coverage-frontend.out ./...
	@echo "Testing job runner..."
	cd job-runner && go test -v -race -coverprofile=coverage-job-runner.out ./...
	@echo "Testing shared package..."
	cd shared && go test -v -race -coverprofile=coverage-shared.out ./...
	@echo "All CI tests passed!"

# Run integration tests
test-integration:
	@echo "Running integration tests..."
	go test -v -tags=integration ./...
	@echo "Integration tests passed!"

# Clean build artifacts
clean:
	@echo "Cleaning build artifacts..."
	cd api-server && make clean
	cd job-runner && make clean
	cd frontend && make clean
	rm -f coverage.out coverage.html
	@echo "Clean complete!"

# Create bin directory
bin:
	mkdir -p bin

# Start RabbitMQ
rabbitmq-up:
	@echo "Starting RabbitMQ container..."
	@if [ "$$(docker ps -q -f name=microservices-rabbitmq)" ]; then \
		echo "RabbitMQ container is already running"; \
	elif [ "$$(docker ps -aq -f name=microservices-rabbitmq)" ]; then \
		echo "Starting existing RabbitMQ container..."; \
		docker start microservices-rabbitmq; \
	else \
		echo "Creating new RabbitMQ container..."; \
		docker run -d --name microservices-rabbitmq \
			-p 5672:5672 -p 15672:15672 \
			-e RABBITMQ_DEFAULT_USER=guest \
			-e RABBITMQ_DEFAULT_PASS=guest \
			rabbitmq:3-management; \
	fi
	@echo "Waiting for RabbitMQ to be ready..."
	@sleep 10
	@echo "RabbitMQ is ready!"
	@echo "Management UI: http://localhost:15672 (guest/guest)"

# Stop RabbitMQ
rabbitmq-down:
	@echo "Stopping RabbitMQ container..."
	@docker stop microservices-rabbitmq || true
	@docker rm microservices-rabbitmq || true

# Run API Server
run-api:
	@echo "Starting API Server on :8081..."
	cd api-server && make run

# Run Job Runner
run-job:
	@echo "Starting Job Runner..."
	cd job-runner && make run

# Run Frontend
run-frontend:
	@echo "Starting Frontend on :8080..."
	cd frontend && make run

# Docker: Build all service images
docker-build:
	@echo "Building all Docker images..."
	cd api-server && make docker-build
	cd job-runner && make docker-build
	cd frontend && make docker-build
	@echo "All Docker images built!"

# Docker: Build and push to custom repository
# Usage: make docker-build-push REPO=localhost:5000
# Usage: make docker-build-push REPO=registry.company.com
docker-build-push:
	@if [ -z "$(REPO)" ]; then \
		echo "❌ Error: REPO parameter is required"; \
		echo "Usage: make docker-build-push REPO=your-registry.com"; \
		echo "Examples:"; \
		echo "  make docker-build-push REPO=localhost:5000"; \
		echo "  make docker-build-push REPO=docker.io/username"; \
		echo "  make docker-build-push REPO=gcr.io/project-id"; \
		exit 1; \
	fi
	@echo "Building and pushing all Docker images to $(REPO)..."
	@echo "Building API Server..."
	docker build -f api-server/Dockerfile -t $(REPO)/microservices-api-server:latest .
	docker push $(REPO)/microservices-api-server:latest
	@echo "Building Job Runner..."
	docker build -f job-runner/Dockerfile -t $(REPO)/microservices-job-runner:latest .
	docker push $(REPO)/microservices-job-runner:latest
	@echo "Building Frontend..."
	docker build -f frontend/Dockerfile -t $(REPO)/microservices-frontend:latest .
	docker push $(REPO)/microservices-frontend:latest
	@echo "✅ All images built and pushed to $(REPO)!"

# Docker: Push existing images to custom repository
# Usage: make docker-push REPO=localhost:5000
docker-push:
	@if [ -z "$(REPO)" ]; then \
		echo "❌ Error: REPO parameter is required"; \
		echo "Usage: make docker-push REPO=your-registry.com"; \
		echo "Examples:"; \
		echo "  make docker-push REPO=localhost:5000"; \
		echo "  make docker-push REPO=docker.io/username"; \
		echo "  make docker-push REPO=gcr.io/project-id"; \
		exit 1; \
	fi
	@echo "Tagging and pushing existing images to $(REPO)..."
	@echo "Tagging and pushing API Server..."
	docker tag microservices-api-server:latest $(REPO)/microservices-api-server:latest
	docker push $(REPO)/microservices-api-server:latest
	@echo "Tagging and pushing Job Runner..."
	docker tag microservices-job-runner:latest $(REPO)/microservices-job-runner:latest
	docker push $(REPO)/microservices-job-runner:latest
	@echo "Tagging and pushing Frontend..."
	docker tag microservices-frontend:latest $(REPO)/microservices-frontend:latest
	docker push $(REPO)/microservices-frontend:latest
	@echo "✅ All images tagged and pushed to $(REPO)!"

# Docker: Build and push with custom tag
# Usage: make docker-build-push-tag REPO=localhost:5000 TAG=v1.0.0
docker-build-push-tag:
	@if [ -z "$(REPO)" ]; then \
		echo "❌ Error: REPO parameter is required"; \
		echo "Usage: make docker-build-push-tag REPO=your-registry.com TAG=version"; \
		exit 1; \
	fi
	@if [ -z "$(TAG)" ]; then \
		echo "❌ Error: TAG parameter is required"; \
		echo "Usage: make docker-build-push-tag REPO=your-registry.com TAG=version"; \
		exit 1; \
	fi
	@echo "Building and pushing all Docker images to $(REPO) with tag $(TAG)..."
	@echo "Building API Server..."
	docker build -f api-server/Dockerfile -t $(REPO)/microservices-api-server:$(TAG) .
	docker push $(REPO)/microservices-api-server:$(TAG)
	@echo "Building Job Runner..."
	docker build -f job-runner/Dockerfile -t $(REPO)/microservices-job-runner:$(TAG) .
	docker push $(REPO)/microservices-job-runner:$(TAG)
	@echo "Building Frontend..."
	docker build -f frontend/Dockerfile -t $(REPO)/microservices-frontend:$(TAG) .
	docker push $(REPO)/microservices-frontend:$(TAG)
	@echo "✅ All images built and pushed to $(REPO) with tag $(TAG)!"

# Docker: Start entire stack with docker-compose
docker-up: docker-build
	@echo "Starting full microservices stack with Docker Compose..."
	docker-compose up -d
	@echo "Stack started! Services available at:"
	@echo "  Frontend: http://localhost:8080"
	@echo "  API Server: http://localhost:8081"
	@echo "  RabbitMQ Management: http://localhost:15672 (guest/guest)"

# Docker: Stop entire stack
docker-down:
	@echo "Stopping Docker Compose stack..."
	docker-compose down
	@echo "Stack stopped!"

# Docker: View logs
docker-logs:
	docker-compose logs -f

# Docker: Restart stack
docker-restart: docker-down docker-up

# Docker: Clean up Docker images and containers
docker-clean:
	@echo "🧹 Cleaning up Docker images and containers..."
	@echo "Stopping and removing all containers..."
	docker-compose down --remove-orphans || true
	@echo "Removing dangling images..."
	docker image prune -f
	@echo "Removing microservices-demo images..."
	docker images | grep "microservices-" | awk '{print $$3}' | xargs -r docker rmi -f || true
	@echo "Removing unused volumes..."
	docker volume prune -f
	@echo "Removing unused networks..."
	docker network prune -f
	@echo "✅ Docker cleanup complete!"

# Docker: Clean everything (including all unused images)
docker-clean-all: docker-clean
	@echo "🧹 Performing deep Docker cleanup..."
	@echo "Removing all unused images (not just dangling)..."
	docker image prune -a -f
	@echo "Removing all stopped containers..."
	docker container prune -f
	@echo "✅ Deep Docker cleanup complete!"

# Run all services (requires multiple terminals)
run-all:
	@echo "To run all services, open 3 separate terminals and run:"
	@echo "  Terminal 1: make run-api"
	@echo "  Terminal 2: make run-job" 
	@echo "  Terminal 3: make run-frontend"
	@echo ""
	@echo "Or use Docker to run everything:"
	@echo "  make docker-up"

# Run all services in background (for testing)
run-all-background: rabbitmq-up
	@echo "Starting all services in background..."
	cd api-server && make run &
	sleep 2
	cd job-runner && make run &
	sleep 2
	cd frontend && make run &
	@echo "All services started. Check http://localhost:8080"
	@echo "To stop all services: make stop-all"

# Stop all background services
stop-all:
	@echo "Stopping all services..."
	pkill -f "go run main.go" || true
	make rabbitmq-down

# Quick start - start RabbitMQ and show instructions
quick-start: rabbitmq-up
	@echo ""
	@echo "🚀 Quick Start Instructions:"
	@echo "================================"
	@echo "Choose one of the following options:"
	@echo ""
	@echo "Option 1 - Local Development (3 terminals):"
	@echo "  Terminal 1: make run-api"
	@echo "  Terminal 2: make run-job"
	@echo "  Terminal 3: make run-frontend"
	@echo ""
	@echo "Option 2 - Docker (single command):"
	@echo "  make docker-up"
	@echo ""
	@echo "Then visit: http://localhost:8080"
	@echo ""
	@echo "Service URLs:"
	@echo "  Frontend: http://localhost:8080"
	@echo "  API Server: http://localhost:8081"
	@echo "  RabbitMQ Management: http://localhost:15672 (guest/guest)"
	@echo ""

# Development mode - run with hot reload (requires air)
dev-api:
	cd api-server && make dev

dev-job:
	cd job-runner && make dev

dev-frontend:
	cd frontend && make dev

# Benchmark tests
benchmark:
	@echo "Running benchmarks..."
	go test -bench=. -benchmem ./...

# Check code formatting
fmt:
	go fmt ./...

# Run linter (requires golangci-lint)
lint:
	@if command -v golangci-lint > /dev/null; then \
		golangci-lint run; \
	else \
		echo "golangci-lint not found. Install it with:"; \
		echo "go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest"; \
	fi

# Show help
help:
	@echo "Available commands:"
	@echo ""
	@echo "🚀 Getting Started:"
	@echo "  make quick-start    - Start RabbitMQ and show setup instructions"
	@echo "  make docker-up      - Start entire stack with Docker Compose"
	@echo ""
	@echo "🔧 Development:"
	@echo "  make deps           - Install Go dependencies"
	@echo "  make build          - Build all services"
	@echo "  make run-api        - Run API server (port 8081)"
	@echo "  make run-job        - Run job runner"
	@echo "  make run-frontend   - Run frontend (port 8080)"
	@echo ""
	@echo "🐳 Docker Commands:"
	@echo "  make docker-build   - Build all Docker images"
	@echo "  make docker-up      - Start stack with docker-compose"
	@echo "  make docker-down    - Stop docker-compose stack"
	@echo "  make docker-logs    - View container logs"
	@echo "  make docker-restart - Restart entire stack"
	@echo "  make docker-clean   - Clean up Docker images and containers"
	@echo "  make docker-clean-all - Deep clean (removes all unused images)"
	@echo ""
	@echo "📦 Docker Registry Commands:"
	@echo "  make docker-build-push REPO=<registry> - Build and push to registry"
	@echo "  make docker-push REPO=<registry>       - Push existing images to registry"
	@echo "  make docker-build-push-tag REPO=<registry> TAG=<tag> - Build and push with custom tag"
	@echo "  Examples:"
	@echo "    make docker-build-push REPO=localhost:5000"
	@echo "    make docker-build-push REPO=docker.io/username"
	@echo "    make docker-build-push-tag REPO=gcr.io/project-id TAG=v1.0.0"
	@echo ""
	@echo "🧪 Testing:"
	@echo "  make test           - Run all tests"
	@echo "  make test-coverage  - Run tests with coverage"
	@echo ""
	@echo "🛠️ Development Tools:"
	@echo "  make dev-api        - Run API with hot reload"
	@echo "  make dev-job        - Run job runner with hot reload"
	@echo "  make dev-frontend   - Run frontend with hot reload"
	@echo ""
	@echo "🐰 RabbitMQ:"
	@echo "  make rabbitmq-up    - Start RabbitMQ container"
	@echo "  make rabbitmq-down  - Stop RabbitMQ container"
	@echo ""
	@echo "🧹 Cleanup:"
	@echo "  make clean          - Clean build artifacts"
	@echo "  make stop-all       - Stop all background services"
	@echo "  make docker-clean   - Clean Docker images and containers"
	@echo "  make docker-clean-all - Deep Docker cleanup (removes all unused images)"
	@echo ""
	@echo "📚 Individual Service Help:"
	@echo "  cd api-server && make help"
	@echo "  cd job-runner && make help"
	@echo "  cd frontend && make help"
