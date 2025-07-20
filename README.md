# Microservices Demo Application

[![CI](https://github.com/kringen/homelab/actions/workflows/ci.yml/badge.svg)](https://github.com/kringen/homelab/actions/workflows/ci.yml)
[![Deploy](https://github.com/kringen/homelab/actions/workflows/deploy.yml/badge.svg)](https://github.com/kringen/homelab/actions/workflows/deploy.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/kringen/homelab)](https://goreportcard.com/report/github.com/kringen/homelab)

A demonstration of Go microservices architecture using RabbitMQ for message queuing.

## Architecture

This application consists of three microservices:

1. **Frontend** - A lightweight web application that provides a user interface
2. **API Server** - A REST API that handles requests and queues jobs
3. **Job Runner** - A worker service that processes jobs and reports completion

## Components

### Frontend (`frontend/`)
- Lightweight web application built with Gin
- Provides HTML interface to submit jobs
- Displays job status and results

### API Server (`api-server/`)
- REST API built with Gin
- Handles job submission requests
- Publishes jobs to RabbitMQ queue
- Provides job status endpoints

### Job Runner (`job-runner/`)
- Worker service that consumes jobs from RabbitMQ
- Processes jobs asynchronously with 5-60 second durations (max 1 minute)
- Includes timeout protection to ensure no job exceeds 1 minute
- Reports job status changes (pending → processing → completed/failed)
- Simulates different types of work based on job description
- Reports job completion back to message queue

### Shared (`shared/`)
- Common data structures and message types
- RabbitMQ client utilities
- Configuration helpers

## Prerequisites

- Go 1.21+
- Docker and Docker Compose (for containerized deployment)
- `sudo` access for Docker commands (or user added to docker group)
- RabbitMQ server (automatically started with Docker commands)

## Project Structure

This project uses a **multi-Makefile architecture** for better organization:

```
microservices-demo/
├── Makefile              # Main project commands (docker-up, test-all, etc.)
├── docker-compose.yml    # Full stack container orchestration
├── api-server/
│   ├── Makefile         # API server specific commands
│   ├── Dockerfile       # Container definition
│   └── README.md        # Detailed API server docs
├── job-runner/
│   ├── Makefile         # Job runner specific commands  
│   ├── Dockerfile       # Container definition
│   └── README.md        # Detailed job runner docs
├── frontend/
│   ├── Makefile         # Frontend specific commands
│   ├── Dockerfile       # Container definition
│   └── README.md        # Detailed frontend docs
├── shared/
│   └── README.md        # Shared utilities documentation
└── scripts/
    ├── README.md        # Script documentation
    ├── demo.sh          # End-to-end demo
    └── health-check.sh  # System health verification
```

Each service is **fully independent** with its own:
- Dockerfile for containerization
- Makefile for build/run/test commands
- README.md with detailed documentation
- Environment-based configuration

## Quick Start

The easiest way to get started:

### Option 1: Docker (Recommended)

**Note: All Docker commands require `sudo` unless your user is in the docker group.**

```bash
# Start entire stack with one command
make docker-up

# View logs
make docker-logs

# Stop everything
make docker-down
```

### Option 2: Local Development
```bash
make quick-start
```

This will start RabbitMQ and show you the setup instructions for running services individually.

## Running the Application

### Option 1: Using Docker Compose (Recommended)

**Note: Docker commands require `sudo` unless your user is in the docker group.**

1. **Start entire stack**:
```bash
make docker-up
```

This will:
- Build Docker images for all services
- Start RabbitMQ, API Server, Job Runner, and Frontend
- Set up networking between containers
- Expose services on standard ports

2. **Access the application**:
- Frontend: http://localhost:8080
- API Server: http://localhost:8081  
- RabbitMQ Management: http://localhost:15672 (guest/guest)

3. **View logs**:
```bash
make docker-logs
```

4. **Stop everything**:
```bash
make docker-down
```

### Option 2: Using Make Commands (Local Development)

1. **Start RabbitMQ**:
```bash
make rabbitmq-up
```

2. **Install dependencies and build**:
```bash
make deps
make build
```

3. **Start the services** (in separate terminals):
```bash
# Terminal 1: Start API Server
make run-api

# Terminal 2: Start Job Runner  
make run-job

# Terminal 3: Start Frontend
make run-frontend
```

4. **Access the application**:
- Frontend: http://localhost:8080
- API Server: http://localhost:8081
- RabbitMQ Management: http://localhost:15672 (guest/guest)

### Option 3: Individual Service Development

Each service is **completely independent** with its own Makefile and can be developed/deployed separately:

```bash
# API Server (in api-server/ directory)
cd api-server
make help              # See all available commands
make run               # Run locally
make docker-build      # Build Docker image
make docker-run        # Run in container
make test              # Run service tests

# Job Runner (in job-runner/ directory)  
cd job-runner
make help              # See all available commands
make run               # Run locally
make run-multiple      # Run multiple instances
make docker-run        # Run in container

# Frontend (in frontend/ directory)
cd frontend
make help              # See all available commands
make run               # Run locally
make open              # Open in browser
make docker-run        # Run in container
```

**Service Independence Benefits:**
- ✅ Each service can be built, tested, and deployed separately
- ✅ Independent versioning and release cycles
- ✅ Technology stack flexibility per service
- ✅ Isolated development environments
- ✅ Microservices best practices

See individual service READMEs for detailed instructions:
- [API Server](api-server/README.md) - REST API, job management, RabbitMQ integration
- [Job Runner](job-runner/README.md) - Async processing, scaling, job simulation
- [Frontend](frontend/README.md) - Web UI, AJAX, real-time updates
- [Shared](shared/README.md) - Common utilities, data models, RabbitMQ client

## Available Make Commands

The project uses a **hierarchical Makefile structure**:

### 🏠 Main Project Commands (Root Makefile)
**Global operations affecting the entire stack:**

#### Quick Start & Setup
- `make quick-start` - Start RabbitMQ and show setup instructions
- `make deps` - Install Go dependencies for all services
- `make build` - Build all services into individual bin/ directories
- `make clean` - Clean build artifacts from all services

#### Docker Commands (Recommended for Full Stack)
- `make docker-build` - Build all Docker images
- `make docker-up` - **🚀 Start entire stack with docker-compose**
- `make docker-down` - Stop docker-compose stack  
- `make docker-logs` - View container logs
- `make docker-restart` - Restart entire stack

#### Running Services (Local Development)
- `make run-api` - Run API server on port 8081
- `make run-job` - Run job runner service
- `make run-frontend` - Run frontend web app on port 8080
- `make run-all-background` - Run all services in background (for testing)
- `make stop-all` - Stop all background services

#### Docker/RabbitMQ
- `make rabbitmq-up` - Start RabbitMQ container with Docker
- `make rabbitmq-down` - Stop and remove RabbitMQ container

#### Development
- `make dev-api` - Run API server with hot reload (requires air)
- `make dev-job` - Run job runner with hot reload (requires air)  
- `make dev-frontend` - Run frontend with hot reload (requires air)

#### Testing & Quality
- `make test` - Run all tests across all services
- `make test-coverage` - Run tests with coverage report
- `make test-ci` - Run comprehensive CI tests (includes individual service tests)
- `make test-integration` - Run integration tests with real services
- `make benchmark` - Run benchmark tests
- `make fmt` - Format all Go code
- `make lint` - Run linter (requires golangci-lint)

### 🔧 Individual Service Commands
**Each service has its own comprehensive Makefile:**

#### Get Service-Specific Help
```bash
cd api-server && make help     # 📡 API server commands
cd job-runner && make help     # ⚙️ Job runner commands  
cd frontend && make help       # 🎨 Frontend commands
```

#### Common Service Commands Available in Each Directory
- `make run` - Run service locally
- `make build` - Build service binary
- `make test` - Run service tests
- `make docker-build` - Build service Docker image
- `make docker-run` - Run service in container
- `make dev` - Run with hot reload
- `make clean` - Clean service artifacts

#### Service-Specific Commands
**API Server** (`api-server/`):
- `make health-check` - Check API server health
- `make run-env` - Run with custom environment

**Job Runner** (`job-runner/`):
- `make run-multiple` - Run multiple instances for load testing
- `make docker-run-network` - Run with Docker network

**Frontend** (`frontend/`):
- `make open` - Open frontend in browser
- `make health-check` - Check frontend accessibility

### 📋 Command Usage Examples

```bash
# Start entire stack (recommended)
make docker-up

# Develop individual service with hot reload
cd api-server && make dev

# Build and test everything
make build && make test

# Scale job runners
cd job-runner && make run-multiple

# Get help for any level
make help                    # Main project help
cd frontend && make help     # Frontend-specific help
```

### Help
- `make help` - Show all main project commands

## Scripts

The `scripts/` directory contains useful utilities for development and demonstration:

### 🏥 Health Check
```bash
./scripts/health-check.sh
```
Comprehensive health check for all services and dependencies. Verifies:
- RabbitMQ Management UI accessibility
- API Server health endpoint  
- Frontend web interface
- Job Runner process status

### 🚀 Demo Script
```bash
./scripts/demo.sh
```
Interactive demonstration of the complete job processing workflow:
- Creates 4 test jobs with different types
- Monitors real-time status updates
- Shows final results with detailed output
- Perfect for demos and testing end-to-end functionality

### ✅ Pre-commit Checks
```bash
./scripts/pre-commit.sh
```
Runs the same checks as CI to ensure your code will pass:
- Code formatting verification (`gofmt`)
- Static analysis (`go vet`)
- All unit tests
- Linting (`golangci-lint`)
- Perfect for running before committing code

See [scripts/README.md](scripts/README.md) for detailed documentation.

## Testing

Run tests for all services:
```bash
make test
```

Run tests with coverage:
```bash
make test-coverage
```

Run individual module tests:
```bash
go test ./shared/
go test ./api-server/
go test ./job-runner/
go test ./frontend/
```

## API Endpoints

### API Server (Port 8081)

- `POST /api/jobs` - Submit a new job
- `GET /api/jobs/{id}` - Get job status
- `GET /api/jobs` - List all jobs

### Frontend (Port 8080)

- `GET /` - Main interface
- `POST /submit` - Submit job form
- `GET /status/{id}` - Job status page

## Kubernetes Deployment

Deploy to Kubernetes using the enhanced deployment script with support for multiple environments, container registries, and custom hostnames.

### Quick Start

```bash
# Deploy to development environment (uses defaults)
./k8s/deploy.sh development apply

# Deploy to production environment
./k8s/deploy.sh production apply

# Deploy with custom registry and version
./k8s/deploy.sh development apply kringen v1.2.3

# Deploy with custom hostname
./k8s/deploy.sh development apply kringen v1.2.3 my-app.example.com
```

### Deployment Script Usage

**Syntax**: `./k8s/deploy.sh [environment] [action] [registry] [tag] [hostname]`

**Parameters**:
- `environment` - Target environment (`development` or `production`)
- `action` - Action to perform (`apply`, `delete`, `diff`, or `build`)
- `registry` - Container registry URL (optional, defaults to `localhost:5000`)
- `tag` - Image tag (optional, defaults to `latest`)
- `hostname` - Custom hostname for ingress (optional, uses environment defaults)

**Examples**:
```bash
# Development deployment with defaults
./k8s/deploy.sh development apply
# → Uses: localhost:5000 registry, latest tag, microservices-demo.local hostname

# Production deployment with custom settings
./k8s/deploy.sh production apply kringen v2.1.0 microservices.kringen.io
# → Uses: kringen registry, v2.1.0 tag, custom hostname

# View what would be deployed without applying
./k8s/deploy.sh development diff kringen v1.0.0

# Build manifests only (for debugging)
./k8s/deploy.sh development build

# Clean up deployment
./k8s/deploy.sh development delete
```

### Environment Configuration

The deployment system supports environment-specific configurations:

#### Development Environment
- **Default hostname**: `microservices-demo.local`
- **Default registry**: `localhost:5000`
- **Resource limits**: Lower CPU/memory for development
- **Replicas**: Single replica for each service
- **Debug mode**: Enabled with verbose logging

#### Production Environment  
- **Default hostname**: `microservices-demo.kringen.io`
- **Default registry**: `registry.company.com`
- **Resource limits**: Production-grade CPU/memory
- **Replicas**: Multiple replicas with autoscaling
- **Security**: Hardened security contexts

### TLS and Certificate Management

The ingress is configured with **cert-manager** for automatic TLS certificate provisioning:

- **Development**: Uses `.local` domains (requires manual certificate or DNS setup)
- **Production**: Uses Let's Encrypt with HTTP-01 challenge for automatic certificates
- **Custom domains**: Automatically provisions certificates for any valid domain

**Certificate Features**:
- Automatic certificate issuance and renewal
- TLS redirect (HTTP → HTTPS)
- Support for custom domains via hostname parameter

### Container Registry Support

The deployment script supports multiple container registries:

```bash
# Local registry (development)
./k8s/deploy.sh development apply localhost:5000 v1.0.0

# Docker Hub
./k8s/deploy.sh production apply kringen v1.0.0

# Private registry
./k8s/deploy.sh production apply registry.company.com v1.0.0

# Google Container Registry
./k8s/deploy.sh production apply gcr.io/project-id v1.0.0
```

### Monitoring Deployment Status

```bash
# Check deployment status
kubectl get pods -n microservices-demo

# View deployment details
kubectl describe deployment -n microservices-demo

# Check service endpoints
kubectl get svc -n microservices-demo

# Monitor logs
kubectl logs -f deployment/dev-frontend -n microservices-demo
kubectl logs -f deployment/dev-api-server -n microservices-demo
kubectl logs -f deployment/dev-job-runner -n microservices-demo
```

### Prerequisites
- Kubernetes cluster (v1.20+)
- kubectl configured with cluster access
- **cert-manager** installed (for TLS certificates)
- **nginx-ingress-controller** installed (for ingress)
- Container registry access (Docker Hub, private registry, etc.)
- **kustomize** or `kubectl kustomize` support

### What Gets Deployed
- **RabbitMQ server** with persistent storage (NFS-backed PVC)
- **API Server** with health checks and resource limits
- **Job Runner** with configurable replica count and scaling
- **Frontend** web interface with ingress and TLS
- **ConfigMaps** with environment-specific configuration
- **Secrets** for RabbitMQ credentials and application config
- **Services** for internal communication and external access
- **Ingress** with automatic TLS certificate provisioning
- **Network policies** for secure pod-to-pod communication (production)

For detailed Kubernetes documentation, see [k8s/README.md](k8s/README.md).

## Continuous Integration & Deployment

This project includes automated CI/CD pipelines via GitHub Actions:

### 🔄 CI Workflow (`.github/workflows/ci.yml`)
**Triggers**: Push to `main`/`develop`, Pull Requests to `main`

**Features**:
- ✅ **Comprehensive Testing**: Tests all microservices with RabbitMQ integration
- 🔍 **Code Quality**: go vet, go fmt, golangci-lint checks
- 📊 **Coverage Reports**: Uploads coverage to Codecov
- 🐳 **Docker Build**: Builds and pushes images on main branch
- 🧪 **Integration Tests**: End-to-end testing with real services

### 🚀 Deploy Workflow (`.github/workflows/deploy.yml`)  
**Triggers**: Push to `main` (auto-deploy), Manual workflow dispatch

**Features**:
- 🎯 **Environment Selection**: Deploy to development or production
- 🏷️ **Version Control**: Deploy specific image tags
- 🌐 **Custom Hostnames**: Override default hostnames
- ✅ **Deployment Verification**: Waits for pods to be ready
- 📊 **Health Checks**: Post-deployment verification

**Manual Deployment**:
1. Go to **GitHub Actions** → **Deploy** workflow → **"Run workflow"**
2. Select environment (`development` or `production`)
3. Specify image tag (default: `latest`)  
4. Add custom hostname (optional)

### 📊 Status Monitoring

Check workflow status locally:
```bash
# Install GitHub CLI (if not installed)
brew install gh  # macOS
apt install gh    # Ubuntu

# Check workflow status
./scripts/check-workflows.sh

# View specific workflow run
gh run list --workflow=ci.yml
gh run view [run-id]
```

## Troubleshooting

### Common Issues

1. **Port already in use errors**:
   ```bash
   make stop-all  # Stop any running services
   ```

2. **RabbitMQ connection failures**:
   ```bash
   make rabbitmq-down
   make rabbitmq-up
   ```

3. **Build failures**:
   ```bash
   make clean
   make deps
   make build
   ```

4. **Docker permission issues**:
   - All Docker commands in this project require `sudo` by default
   - Alternative: Add your user to the docker group and restart your session:
     ```bash
     sudo usermod -aG docker $USER
     newgrp docker  # Or logout/login
     ```

### Job Processing Behavior

- Jobs are processed **asynchronously** - the API returns immediately with status "pending"
- Job status updates happen via RabbitMQ messaging
- Frontend polls for updates every 3 seconds
- Jobs take 5-60 seconds to complete (simulated processing time)
- Jobs have a 1-minute timeout protection
