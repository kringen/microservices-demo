# CI/CD Pipeline Documentation

This document provides comprehensive information about the CI/CD pipeline for the Microservices Demo application.

## Table of Contents

- [Overview](#overview)
- [Pipeline Architecture](#pipeline-architecture)
- [Workflow Details](#workflow-details)
- [Configuration](#configuration)
- [Local Development](#local-development)
- [Troubleshooting](#troubleshooting)
- [Best Practices](#best-practices)

## Overview

The project uses **GitHub Actions** for CI/CD with a comprehensive pipeline that ensures:
- 🧪 **Quality**: Automated testing, linting, and formatting checks
- 🔨 **Building**: Parallel compilation and Docker image creation
- 🚀 **Deployment**: Automated Kubernetes deployment with environment support
- 📊 **Monitoring**: Coverage reporting and artifact management

## Pipeline Architecture

```mermaid
graph TD
    A[Code Push/PR] --> B[CI Workflow]
    B --> C[Test Job]
    B --> D[Build Job]
    B --> E[Docker Build]
    B --> F[Integration Tests]
    
    C --> G[Quality Checks]
    C --> H[Unit Tests]
    C --> I[Coverage Report]
    
    E --> J[Docker Hub]
    
    K[Main Branch Push] --> L[Deploy Workflow]
    L --> M[Kubernetes Deployment]
    
    style A fill:#e1f5fe
    style J fill:#fff3e0
    style M fill:#e8f5e8
```

## Workflow Details

### CI Workflow (`.github/workflows/ci.yml`)

#### **Triggers**
```yaml
on:
  push:
    branches: [ main, develop ]
  pull_request:
    branches: [ main ]
```

#### **Jobs Structure**

##### 1. Test Job 🧪
- **Platform**: `ubuntu-latest`
- **Go Version**: 1.21
- **Services**: RabbitMQ 3.12 with management interface
- **Timeout**: Health checks with retry logic

**Test Matrix:**
```bash
# Root module tests
go test -v -race -coverprofile=coverage.out ./...

# Service-specific tests
cd api-server && go test -v -race -coverprofile=coverage-api.out ./...
cd frontend && go test -v -race -coverprofile=coverage-frontend.out ./...
cd job-runner && go test -v -race -coverprofile=coverage-job-runner.out ./...
cd shared && go test -v -race -coverprofile=coverage-shared.out ./...
```

**Quality Checks:**
- `go vet ./...` - Static analysis
- `gofmt -s -l .` - Code formatting verification
- `golangci-lint run ./...` - Advanced linting

##### 2. Build Job 🔨
- **Strategy**: Matrix build for `[api-server, frontend, job-runner]`
- **Artifacts**: Compiled binaries uploaded for each service
- **Dependencies**: Requires test job to pass

##### 3. Docker Build Job 🐳
- **Condition**: Only on `main` branch pushes
- **Strategy**: Matrix build for all services
- **Registry**: Docker Hub (`kringen/microservices-{service}`)

**Tagging Strategy:**
```yaml
tags: |
  type=ref,event=branch          # branch-name
  type=ref,event=pr              # pr-123
  type=sha,prefix={{branch}}-    # main-abc1234
  type=raw,value=latest,enable={{is_default_branch}}  # latest
```

##### 4. Integration Tests 🔗
- **Environment**: Full RabbitMQ integration
- **Scope**: Cross-service communication validation
- **Tags**: Uses `integration` build tag

### Deploy Workflow (`.github/workflows/deploy.yml`)

#### **Triggers**
```yaml
# Automatic deployment
on:
  push:
    branches: [ main ]

# Manual deployment with options
on:
  workflow_dispatch:
    inputs:
      environment: [development, production]
      tag: string
      hostname: string
```

#### **Deployment Process**
1. **Variable Setup**: Environment-specific configuration
2. **Kubernetes Apply**: Uses enhanced deployment script
3. **Health Verification**: Post-deployment validation

## Configuration

### Required GitHub Secrets

#### **Docker Hub Authentication**
```bash
DOCKER_USERNAME=your-dockerhub-username
DOCKER_PASSWORD=your-dockerhub-access-token
```

**Setup Instructions:**
1. Go to [Docker Hub Access Tokens](https://hub.docker.com/settings/security)
2. Create a new access token with "Read, Write, Delete" permissions
3. Add to GitHub repository secrets

#### **Kubernetes Access (for GitHub-hosted runners)**
```bash
KUBECONFIG=base64-encoded-kubeconfig-file
```

**Setup Instructions:**
```bash
# Encode your kubeconfig
cat ~/.kube/config | base64 -w 0

# Add the output to GitHub secrets as KUBECONFIG
```

### Environment Variables

#### **CI Workflow**
```yaml
env:
  GO_VERSION: '1.21'
  RABBITMQ_URL: amqp://guest:guest@localhost:5672/
  API_SERVER_URL: http://localhost:8081
  FRONTEND_URL: http://localhost:8080
```

#### **Deploy Workflow**
```yaml
env:
  REGISTRY: kringen
```

## Local Development

### Pre-commit Validation

The `scripts/pre-commit.sh` script mirrors CI checks:

```bash
#!/bin/bash
# Comprehensive pre-commit checks

# 1. Dependency verification
echo "🔍 Verifying Go dependencies..."
go mod verify || exit 1

# 2. Code formatting
echo "📝 Checking code formatting..."
if [ "$(gofmt -s -l . | wc -l)" -gt 0 ]; then
    echo "❌ Code formatting issues found"
    gofmt -s -l .
    exit 1
fi

# 3. Static analysis
echo "🔬 Running static analysis..."
go vet ./... || exit 1

# 4. Linting (optional)
if command -v golangci-lint &> /dev/null; then
    echo "🧹 Running linter..."
    golangci-lint run ./... || exit 1
fi

# 5. Testing
echo "🧪 Running tests..."
go test -v -race ./... || exit 1

echo "✅ All checks passed!"
```

### Local Docker Testing

```bash
# Build all images locally
make docker-build

# Test with local registry
docker tag kringen/microservices-api-server:latest localhost:5000/microservices-api-server:test
docker push localhost:5000/microservices-api-server:test

# Deploy to local Kubernetes
./k8s/deploy.sh development apply localhost:5000 test
```

### Simulating CI Environment

```bash
# Start RabbitMQ (matches CI setup)
docker run -d --name rabbitmq-test \
  -p 5672:5672 -p 15672:15672 \
  -e RABBITMQ_DEFAULT_USER=guest \
  -e RABBITMQ_DEFAULT_PASS=guest \
  rabbitmq:3.12-management

# Wait for readiness (matches CI)
timeout 60s bash -c 'until curl -f http://localhost:15672; do sleep 2; done'

# Run tests with same environment
export RABBITMQ_URL=amqp://guest:guest@localhost:5672/
make test

# Cleanup
docker rm -f rabbitmq-test
```

## Troubleshooting

### Common Issues

#### **Docker Authentication Failures**
```bash
# Symptoms:
Error: buildx failed with: ERROR: failed to solve: failed to push

# Solutions:
1. Verify DOCKER_USERNAME and DOCKER_PASSWORD secrets
2. Ensure Docker Hub token has push permissions
3. Check token expiration date
```

#### **RabbitMQ Connection Issues**
```bash
# Symptoms:
dial tcp [::1]:5672: connect: connection refused

# Solutions:
1. CI: Health checks automatically handle this
2. Local: Ensure RabbitMQ is running
   make rabbitmq-up
3. Check port availability (5672, 15672)
```

#### **Kubernetes Deployment Failures**
```bash
# Symptoms:
error: You must be logged in to the server (the server has asked for the client to provide credentials)

# Solutions:
1. Verify KUBECONFIG secret (if using GitHub-hosted runners)
2. Check self-hosted runner kubectl access
3. Validate cluster connectivity:
   kubectl cluster-info
```

#### **Go Module Issues**
```bash
# Symptoms:
go: github.com/some/module@version: reading module: unknown revision

# Solutions:
1. Verify go.mod and go.sum consistency
2. Run go mod download locally
3. Check for private module access issues
```

### Debugging Pipeline Failures

#### **Viewing Detailed Logs**
1. **GitHub Actions**: Go to Actions tab → Select workflow run → Click job
2. **Expand failing step**: Review full command output
3. **Download artifacts**: Access compiled binaries for debugging

#### **Testing Pipeline Changes**
```bash
# 1. Create feature branch
git checkout -b fix/ci-pipeline

# 2. Modify workflow files
vim .github/workflows/ci.yml

# 3. Test locally first
./scripts/pre-commit.sh

# 4. Push and monitor
git push origin fix/ci-pipeline
# Check GitHub Actions for results
```

### Performance Optimization

#### **Cache Strategy**
```yaml
# Go modules cache
- uses: actions/cache@v4
  with:
    path: |
      ~/.cache/go-build
      ~/go/pkg/mod
    key: ${{ runner.os }}-go-${{ hashFiles('**/go.sum') }}
```

#### **Parallel Execution**
- **Matrix builds**: Services build in parallel
- **Job dependencies**: `needs:` ensures proper order
- **Resource allocation**: Optimized for GitHub Actions limits

## Best Practices

### Code Quality Gates

1. **Pre-commit hooks**: Run `scripts/pre-commit.sh` before pushing
2. **Pull request reviews**: Require approvals for main branch
3. **Status checks**: Enforce CI passage before merge

### Security Considerations

1. **Secret management**: Use GitHub secrets, never hardcode
2. **Token rotation**: Regularly update Docker Hub tokens
3. **Kubernetes RBAC**: Limit deployment permissions
4. **Image scanning**: Consider adding vulnerability scanning

### Deployment Strategy

1. **Environment promotion**: 
   - Develop → Main (auto-deploy to dev)
   - Main → Production (manual approval)
2. **Rollback capability**: Keep previous image tags
3. **Health checks**: Verify deployment success

### Monitoring and Alerting

1. **Pipeline notifications**: Configure GitHub notifications
2. **Deployment monitoring**: Set up Kubernetes health checks
3. **Coverage tracking**: Monitor Codecov trends
4. **Performance baselines**: Track build times and test duration

### Documentation Maintenance

1. **Keep workflows updated**: Document any configuration changes
2. **Version compatibility**: Update Go version, action versions
3. **Runbook updates**: Maintain troubleshooting guides
4. **Team knowledge**: Regular pipeline review sessions

## Additional Resources

- [GitHub Actions Documentation](https://docs.github.com/en/actions)
- [Docker Hub Access Tokens](https://docs.docker.com/hub/access-tokens/)
- [Kubernetes kubectl Reference](https://kubernetes.io/docs/reference/kubectl/)
- [Go Testing Package](https://pkg.go.dev/testing)
- [golangci-lint Configuration](https://golangci-lint.run/usage/configuration/)
- **[CI/CD Best Practices](CICD-BEST-PRACTICES.md)** - Optimization strategies and build efficiency
