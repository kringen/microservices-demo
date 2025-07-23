# AI Research Agent - Microservices Demo

[![CI](https://github.com/kringen/homelab/actions/workflows/ci.yml/badge.svg)](https://github.com/kringen/homelab/actions/workflows/ci.yml)
[![Deploy](https://github.com/kringen/homelab/actions/workflows/deploy.yml/badge.svg)](https://github.com/kringen/homelab/actions/workflows/deploy.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/kringen/homelab)](https://goreportcard.com/report/github.com/kringen/homelab)

An AI-powered research agent built with Go microservices, featuring Dapr service mesh, Ollama LLM integration, and Model Context Protocol (MCP) services for intelligent information gathering and analysis.

## 🤖 AI Features

- **Intelligent Research**: AI-powered research agent using local Ollama LLM
- **MCP Integration**: Model Context Protocol services for web search, GitHub, and file access
- **Dapr Service Mesh**: Modern microservices communication and state management
- **Real-time Results**: Live research progress and confidence scoring
- **Multi-source Analysis**: Combines data from multiple sources for comprehensive insights

## 📋 Documentation

**Comprehensive documentation is available in the [`docs/`](docs/) directory:**

| Document | Description |
|----------|-------------|
| **[API Reference](docs/API.md)** | Complete REST API documentation with examples |
| **[Architecture Guide](docs/ARCHITECTURE.md)** | System design, patterns, and technical decisions |
| **[Development Guide](docs/DEVELOPMENT.md)** | Local development setup and workflows |
| **[Deployment Guide](docs/DEPLOYMENT.md)** | Production deployment strategies |
| **[CI/CD Documentation](docs/CICD.md)** | Pipeline automation and GitHub Actions |
| **[Troubleshooting Guide](docs/TROUBLESHOOTING.md)** | Issue diagnosis and resolution |

**Quick Links:**
- 🚀 **New Developer?** Start with [Development Guide](docs/DEVELOPMENT.md)
- 🔌 **API Integration?** See [API Reference](docs/API.md)
- 🚢 **Deploying?** Check [Deployment Guide](docs/DEPLOYMENT.md)
- 🐛 **Issues?** Use [Troubleshooting Guide](docs/TROUBLESHOOTING.md)

## Architecture

This AI-powered application consists of three main microservices:

1. **Frontend** - A web interface for submitting research requests and viewing results
2. **API Server** - A REST API that handles research requests and manages queue
3. **AI Research Agent** - An intelligent agent that processes research using Ollama LLM and MCP services

## Components

### Frontend (`frontend/`)
- Modern web application built with Gin and Bootstrap
- Research request form with MCP service selection
- Real-time research progress and confidence display
- Source citation and result presentation

### API Server (`api-server/`)
- REST API built with Gin
- Handles research request submissions
- Publishes research tasks to RabbitMQ queue
- Tracks research progress and results

### AI Research Agent (`job-runner/`)
- **Ollama Integration**: Uses local LLM for intelligent analysis
- **MCP Services**: Integrates with Model Context Protocol for:
  - Web search and scraping
  - GitHub repository analysis
  - Local file system access
  - Database queries (configurable)
- **Dapr Ready**: Prepared for service mesh integration
- **Intelligent Processing**: 
  - Multi-source data gathering
  - AI-powered analysis and synthesis
  - Confidence scoring for results
  - Source citation tracking
  - Token usage monitoring

### Shared (`shared/`)
- Enhanced data structures for research requests and results
- RabbitMQ client utilities with research-specific messaging
- AI and MCP service type definitions

## AI Infrastructure

### Ollama Integration
- **Local LLM**: Uses Ollama for privacy-preserving AI analysis
- **Model Flexibility**: Configurable model selection (default: llama3.2)
- **Streaming Support**: Ready for real-time response streaming
- **Resource Management**: Token usage tracking and optimization

### MCP Services
- **Web Search**: Intelligent web information gathering
- **GitHub Integration**: Repository and code analysis
- **File System**: Local document and data access
- **Extensible**: Easy addition of new MCP service types

## Prerequisites

- Go 1.21+
- Docker and Docker Compose
- **Ollama** (for AI functionality): Install from [ollama.ai](https://ollama.ai)
- **Dapr CLI** (optional): For enhanced service mesh features
- Minimum 8GB RAM (for Ollama LLM models)
- `sudo` access for Docker commands (or user added to docker group)

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

## Production Deployment

### Kubernetes Deployment Options

The AI Research Agent supports multiple deployment configurations for Kubernetes:

#### 🏠 Local Ollama (In-Cluster AI)
Deploy with Ollama running as a pod within your Kubernetes cluster:

```bash
# Development environment
kubectl apply -k k8s/overlays/development/

# Production environment  
kubectl apply -k k8s/overlays/production/
```

**Features:**
- ✅ Complete AI stack in cluster
- ✅ No external dependencies
- ✅ Persistent model storage (10GB PVC)
- ⚠️ Requires 4-8GB memory per node
- ⚠️ Initial deployment takes 5-10 minutes (model download)

#### 🌐 External Ollama (Network AI Server)
Deploy using an external Ollama server on your network:

```bash
# Development with external AI
kubectl apply -k k8s/overlays/development-external/

# Production with external AI
kubectl apply -k k8s/overlays/production-external/
```

**Before deployment, update the Ollama URL:**
```bash
# Edit the configmap patch
vim k8s/overlays/development-external/configmap-patch.yaml
# Update OLLAMA_URL to your server: http://your-ollama-server:11434
```

**Features:**
- ✅ Reduced cluster resource usage
- ✅ Shared AI server across workloads  
- ✅ Faster deployments (no model download)
- ✅ Higher concurrency (15 vs 5 jobs)
- ⚠️ Network dependency on external server

#### 📚 Deployment Documentation
- **[Kubernetes AI Deployment Guide](k8s/AI_DEPLOYMENT_GUIDE.md)** - Comprehensive deployment instructions
- **Available overlays:** development, production, development-external, production-external
- **Resource requirements:** Memory, storage, and scaling considerations
- **Troubleshooting:** Common issues and verification steps

### Docker Compose (Development)

For local development and testing:

```bash
# Full AI stack with local Ollama
make docker-up

# View AI system logs
make docker-logs

# Clean up everything including AI models
make docker-clean-all
```

**Cleanup Options:**
- `make docker-down` - Stop services (keep data)
- `make docker-clean` - Remove containers and images
- `make docker-clean-ollama` - Remove AI models and data
- `make docker-clean-all` - Complete cleanup (⚠️ removes everything)

### Environment Configuration

Copy and customize the environment template:

```bash
cp .env.example .env
vim .env  # Configure for your environment
```

**Key configurations:**
- **Local Ollama:** `OLLAMA_URL=http://localhost:11434`
- **External Ollama:** `OLLAMA_URL=http://your-ollama-server:11434`
- **AI Model:** `OLLAMA_MODEL=llama3.2` (or your preferred model)
- **Development:** `LOG_LEVEL=debug` and `GIN_MODE=debug`
- **Production:** `LOG_LEVEL=info` and `GIN_MODE=release`

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

#### Cleanup Commands (Free Disk Space)
- `make docker-clean` - Remove containers & images, keep data volumes  
- `make docker-clean-ollama` - **Remove 2GB Ollama model** (frees significant space)
- `make docker-clean-all` - **⚠️ Full cleanup** including volumes (removes model & data)

> **💡 Tip:** The Ollama AI model is ~2GB. Use `make docker-clean-ollama` to free space if you don't need the AI features temporarily. The model will be automatically re-downloaded on next startup.

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
# Deploy to development environment with local Ollama (uses defaults)
./k8s/deploy.sh --environment development --action apply --ollama local

# Deploy to development with external Ollama server
./k8s/deploy.sh --environment development --action apply --ollama 192.168.1.100:11434

# Deploy to production environment with local Ollama
./k8s/deploy.sh --environment production --action apply --ollama local

# Deploy with custom registry and version
./k8s/deploy.sh -e dev -a apply -o local -r kringen -t v1.2.3

# Deploy with external Ollama and custom hostname
./k8s/deploy.sh -e prod -a apply -o 10.0.1.50:11434 -r kringen -t v1.2.3 -h my-app.example.com
```

### Deployment Script Usage

**Syntax**: `./k8s/deploy.sh [options]`

**Required Parameters**:
- `-e, --environment` - Target environment (`development`/`dev` or `production`/`prod`)
- `-a, --action` - Action to perform (`apply`, `delete`, `diff`, or `build`)

**Optional Parameters**:
- `-o, --ollama` - Ollama deployment type:
  - `local` - Deploy Ollama as a pod in the cluster (default)
  - `<host:port>` - Use external Ollama server (e.g., `192.168.1.100:11434`)
- `-r, --registry` - Container registry URL (defaults to `docker.io`)
- `-t, --tag` - Image tag (defaults to `latest`)
- `-h, --hostname` - Custom hostname for ingress (uses environment defaults)
- `--help` - Show usage information

**Examples**:
```bash
# Development deployment with local Ollama (default)
./k8s/deploy.sh --environment development --action apply --ollama local
# → Uses: local Ollama in cluster, docker.io registry, latest tag

# Development deployment with external Ollama server
./k8s/deploy.sh --environment development --action apply --ollama 192.168.1.100:11434
# → Uses: external Ollama at 192.168.1.100:11434, automatically configures OLLAMA_URL

# Production deployment with external Ollama and custom settings
./k8s/deploy.sh -e prod -a apply -o ollama.company.com:11434 -r kringen -t v2.1.0 -h microservices.kringen.io
# → Uses: external Ollama server, kringen registry, v2.1.0 tag, custom hostname

# View what would be deployed without applying
./k8s/deploy.sh -e dev -a diff -o 10.0.1.50:11434 -r kringen -t v1.0.0

# Build manifests only (for debugging)
./k8s/deploy.sh --environment development --action build --ollama local

# Clean up deployment
./k8s/deploy.sh --environment development --action delete --ollama 192.168.1.100:11434
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

## CI/CD Pipeline

This project includes a comprehensive **GitHub Actions CI/CD pipeline** that automatically tests, builds, and deploys the microservices. The pipeline ensures code quality and provides automated Docker image building with seamless Kubernetes deployment.

### 🔄 Pipeline Overview

The CI/CD pipeline consists of two main workflows:

1. **[CI Workflow](.github/workflows/ci.yml)** - Continuous Integration
2. **[Deploy Workflow](.github/workflows/deploy.yml)** - Continuous Deployment

### 🧪 CI Workflow (`ci.yml`)

**Triggers:**
- Push to `main` or `develop` branches
- Pull requests to `main` branch

**Pipeline Stages:**

#### 1. **Test Job** 🧪
Comprehensive testing with RabbitMQ integration:
- **Environment**: Ubuntu Latest with Go 1.21+
- **Services**: RabbitMQ 3.12 with management interface
- **Quality Checks**:
  - `go vet` - Static analysis
  - `gofmt` - Code formatting verification
  - `golangci-lint` - Advanced linting
  - **Comprehensive test coverage** across all modules:
    - Root module integration tests
    - API Server unit tests
    - Frontend unit tests  
    - Job Runner unit tests
    - Shared package tests
- **Coverage Reporting**: Uploads to Codecov with detailed reports

#### 2. **Build Job** 🔨
Parallel binary compilation:
- **Strategy**: Matrix build for all services (`api-server`, `frontend`, `job-runner`)
- **Artifacts**: Uploads compiled binaries for each service
- **Caching**: Go module and build cache optimization

#### 3. **Docker Build Job** 🐳
Container image creation and registry push:
- **Trigger**: Only on `main` branch pushes
- **Strategy**: Matrix build for all services
- **Features**:
  - Docker Buildx for multi-platform support
  - Automatic Docker Hub authentication
  - **Smart tagging strategy**:
    - `latest` for main branch
    - Branch-based tags
    - SHA-based tags for traceability
  - **Build caching**: GitHub Actions cache for faster builds
  - **Registry**: Pushes to Docker Hub (`kringen/microservices-{service}`)

#### 4. **Integration Tests** 🔗
End-to-end validation:
- **Real services**: Tests with actual RabbitMQ instance
- **Cross-service communication**: Validates message queue integration
- **Environment**: Full microservices environment simulation

### 🚀 Deploy Workflow (`deploy.yml`)

**Triggers:**
- Automatic: Push to `main` branch
- Manual: `workflow_dispatch` with custom parameters

**Deployment Features:**

#### **Flexible Deployment Options**
```yaml
# Manual deployment with options
workflow_dispatch:
  environment: development|production
  tag: custom-tag-or-latest
  hostname: custom.domain.com
```

#### **Environment Support**
- **Development**: Fast deployment with basic configuration
- **Production**: Hardened deployment with security contexts and scaling

#### **Kubernetes Integration**
- **Self-hosted runner**: Assumes kubectl access to your cluster
- **Smart tagging**: Uses specified tag or defaults to `latest`
- **Custom hostnames**: Supports custom domain configuration

### 🔧 CI/CD Configuration

#### **Required Secrets**
Configure these in your GitHub repository settings:

```bash
# Docker Hub Authentication
DOCKER_USERNAME=your-dockerhub-username
DOCKER_PASSWORD=your-dockerhub-password-or-token

# Kubernetes Deployment (if using GitHub-hosted runners)
KUBECONFIG=base64-encoded-kubeconfig
```

#### **Workflow Status Badges**
Track pipeline status with badges (already included in README header):
- [![CI](https://github.com/kringen/homelab/actions/workflows/ci.yml/badge.svg)](https://github.com/kringen/homelab/actions/workflows/ci.yml)
- [![Deploy](https://github.com/kringen/homelab/actions/workflows/deploy.yml/badge.svg)](https://github.com/kringen/homelab/actions/workflows/deploy.yml)

📋 **For detailed CI/CD documentation, troubleshooting, and advanced configuration, see [docs/CICD.md](docs/CICD.md)**

### 🛠️ Local CI/CD Simulation

#### **Pre-commit Validation**
Run the same checks as CI locally:
```bash
./scripts/pre-commit.sh
```
This script mirrors the CI quality checks:
- Code formatting (`gofmt`)
- Static analysis (`go vet`)
- Linting (`golangci-lint`)
- Comprehensive testing
- **Smart RabbitMQ handling** (starts if needed)

#### **Manual Docker Build**
Test Docker operations locally:
```bash
# Build all images
make docker-build

# Push to registry (requires Docker login)
docker push kringen/microservices-api-server:latest
docker push kringen/microservices-frontend:latest
docker push kringen/microservices-job-runner:latest
```

#### **Local Kubernetes Testing**
```bash
# Deploy to local cluster
./k8s/deploy.sh development apply localhost:5000 local-test
```

### 📊 Pipeline Benefits

#### **Quality Assurance**
- ✅ **Automated testing**: Comprehensive test coverage with real dependencies
- ✅ **Code quality**: Multiple linting and formatting checks
- ✅ **Integration validation**: End-to-end testing with message queues

#### **DevOps Efficiency**
- ✅ **Fast feedback**: Parallel job execution for quick results
- ✅ **Artifact management**: Automatic binary and image creation
- ✅ **Environment parity**: Same deployment process for dev/prod

#### **Production Readiness**
- ✅ **Container registry**: Automatic image pushing to Docker Hub
- ✅ **Kubernetes deployment**: Seamless cluster deployment
- ✅ **Version management**: Smart tagging with SHA and branch tracking

### 🐛 Troubleshooting CI/CD

#### **Common Issues**

1. **Docker authentication failures**:
   - Verify `DOCKER_USERNAME` and `DOCKER_PASSWORD` secrets
   - Ensure Docker Hub token has push permissions

2. **Test failures with RabbitMQ**:
   - CI automatically handles RabbitMQ service startup
   - Local tests: ensure RabbitMQ is running (`make rabbitmq-up`)

3. **Kubernetes deployment failures**:
   - Verify self-hosted runner has `kubectl` access
   - Check cluster connectivity and permissions

4. **Build cache issues**:
   - GitHub Actions cache is automatic
   - Local: use `make clean` to reset build state

#### **Monitoring Pipeline Health**
- **GitHub Actions tab**: View real-time pipeline execution
- **Codecov dashboard**: Monitor test coverage trends
- **Docker Hub**: Verify image push success and tags

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
