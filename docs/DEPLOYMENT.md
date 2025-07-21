# Deployment Guide

This guide provides comprehensive information for deploying the Microservices Demo application in various environments.

## Table of Contents

- [Deployment Overview](#deployment-overview)
- [Local Development](#local-development)
- [Docker Deployment](#docker-deployment)
- [Kubernetes Deployment](#kubernetes-deployment)
- [Production Considerations](#production-considerations)
- [Monitoring and Observability](#monitoring-and-observability)
- [Troubleshooting](#troubleshooting)

## Deployment Overview

The Microservices Demo application supports multiple deployment strategies:

| Environment | Method | Use Case | Complexity |
|-------------|--------|----------|------------|
| **Local Development** | Native binaries | Development, debugging | Low |
| **Local Testing** | Docker Compose | Integration testing | Medium |
| **Staging/Production** | Kubernetes | Production workloads | High |
| **Cloud Platforms** | Managed services | Scalable production | Medium |

## Local Development

### Prerequisites
```bash
# Required tools
go version    # Go 1.21+
make --version
docker --version (optional)

# Optional tools
air --version          # Hot reloading
golangci-lint --version # Code quality
```

### Native Binary Deployment

#### Quick Start
```bash
# 1. Clone and setup
git clone https://github.com/kringen/homelab.git
cd microservices-demo
make deps

# 2. Start dependencies
make rabbitmq-up

# 3. Build and run services
make build
make run-all-background

# 4. Verify deployment
make health-check

# 5. Access application
open http://localhost:8080    # Frontend
open http://localhost:8081    # API
open http://localhost:15672   # RabbitMQ Management (guest/guest)
```

#### Manual Service Startup
```bash
# Terminal 1: API Server
cd api-server
go run main.go

# Terminal 2: Job Runner
cd job-runner
go run main.go

# Terminal 3: Frontend
cd frontend
go run main.go
```

### Configuration

#### Environment Variables
```bash
# .env file (create in project root)
RABBITMQ_URL=amqp://guest:guest@localhost:5672/
API_SERVER_URL=http://localhost:8081
FRONTEND_URL=http://localhost:8080
LOG_LEVEL=debug
```

#### Service-Specific Configuration
```bash
# API Server
export API_SERVER_PORT=8081
export API_SERVER_HOST=0.0.0.0

# Frontend
export FRONTEND_PORT=8080
export FRONTEND_HOST=0.0.0.0

# Job Runner
export JOB_RUNNER_WORKERS=3
export JOB_PROCESSING_TIMEOUT=60s
```

## Docker Deployment

### Docker Compose (Recommended for Local Testing)

#### Quick Start
```bash
# Start entire stack
make docker-up

# View logs
make docker-logs

# Stop everything
make docker-down

# Restart services
make docker-restart
```

#### Manual Docker Compose
```bash
# Start with docker-compose directly
docker-compose up -d

# View specific service logs
docker-compose logs -f api-server
docker-compose logs -f job-runner
docker-compose logs -f frontend

# Scale job runners
docker-compose up -d --scale job-runner=3

# Stop and remove everything
docker-compose down -v
```

### Individual Container Deployment

#### Build Images
```bash
# Build all images
make docker-build

# Or build individually
docker build -t kringen/microservices-api-server -f api-server/Dockerfile .
docker build -t kringen/microservices-frontend -f frontend/Dockerfile .
docker build -t kringen/microservices-job-runner -f job-runner/Dockerfile .
```

#### Run Containers
```bash
# 1. Start RabbitMQ
docker run -d --name rabbitmq \
  -p 5672:5672 -p 15672:15672 \
  -e RABBITMQ_DEFAULT_USER=guest \
  -e RABBITMQ_DEFAULT_PASS=guest \
  rabbitmq:3.12-management

# 2. Start API Server
docker run -d --name api-server \
  -p 8081:8081 \
  -e RABBITMQ_URL=amqp://guest:guest@rabbitmq:5672/ \
  --link rabbitmq:rabbitmq \
  kringen/microservices-api-server

# 3. Start Job Runner
docker run -d --name job-runner \
  -e RABBITMQ_URL=amqp://guest:guest@rabbitmq:5672/ \
  --link rabbitmq:rabbitmq \
  kringen/microservices-job-runner

# 4. Start Frontend
docker run -d --name frontend \
  -p 8080:8080 \
  -e API_SERVER_URL=http://api-server:8081 \
  --link api-server:api-server \
  kringen/microservices-frontend
```

### Docker Network Configuration
```bash
# Create custom network
docker network create microservices-net

# Run with custom network
docker run -d --name rabbitmq --network microservices-net \
  -p 5672:5672 -p 15672:15672 \
  rabbitmq:3.12-management

docker run -d --name api-server --network microservices-net \
  -p 8081:8081 \
  -e RABBITMQ_URL=amqp://guest:guest@rabbitmq:5672/ \
  kringen/microservices-api-server
```

## Kubernetes Deployment

### Prerequisites

#### Kubernetes Cluster
```bash
# Local development clusters
minikube start --memory=4096 --cpus=2
# OR
kind create cluster --config=k8s/kind-config.yaml
# OR
k3d cluster create microservices-demo --servers 1 --agents 2
```

#### kubectl Configuration
```bash
# Verify cluster access
kubectl cluster-info
kubectl get nodes

# Create namespace
kubectl create namespace microservices-demo
kubectl config set-context --current --namespace=microservices-demo
```

### Quick Deployment

#### Using Deployment Script
```bash
# Development deployment
./k8s/deploy.sh development apply

# Production deployment
./k8s/deploy.sh production apply

# Custom registry and tag
./k8s/deploy.sh development apply my-registry.com v1.0.0

# With custom hostname
./k8s/deploy.sh production apply kringen v1.0.0 my-app.example.com
```

#### Manual Kubectl Deployment
```bash
# Apply base manifests
kubectl apply -k k8s/base/

# Apply environment-specific overlays
kubectl apply -k k8s/overlays/development/
# OR
kubectl apply -k k8s/overlays/production/
```

### Kubernetes Manifests Overview

#### Base Resources (`k8s/base/`)
```yaml
# Core resources
- namespace.yaml           # Namespace definition
- configmap.yaml          # Application configuration
- secrets.yaml            # Sensitive configuration
- rabbitmq-deployment.yaml # RabbitMQ StatefulSet
- rabbitmq-service.yaml   # RabbitMQ Service
- api-server-deployment.yaml # API Server Deployment
- api-server-service.yaml    # API Server Service
- frontend-deployment.yaml   # Frontend Deployment
- frontend-service.yaml     # Frontend Service
- job-runner-deployment.yaml # Job Runner Deployment
- ingress.yaml             # Ingress Controller
```

#### Environment Overlays

**Development (`k8s/overlays/development/`)**
```yaml
# Development-specific settings
Resources:
  - Lower resource requests/limits
  - Single replica for each service
  - Debug logging enabled
  - Simplified ingress configuration

ConfigMap patches:
  - LOG_LEVEL: debug
  - REPLICAS: 1
```

**Production (`k8s/overlays/production/`)**
```yaml
# Production-specific settings
Resources:
  - Higher resource requests/limits
  - Multiple replicas with HPA
  - Info/error logging only
  - TLS-enabled ingress

ConfigMap patches:
  - LOG_LEVEL: info
  - REPLICAS: 3
  - ENABLE_METRICS: true

Security:
  - Non-root security contexts
  - Read-only root filesystems
  - Network policies
```

### Scaling and High Availability

#### Horizontal Pod Autoscaling
```yaml
# HPA for API Server
apiVersion: autoscaling/v2
kind: HorizontalPodAutoscaler
metadata:
  name: api-server-hpa
spec:
  scaleTargetRef:
    apiVersion: apps/v1
    kind: Deployment
    name: api-server
  minReplicas: 2
  maxReplicas: 10
  metrics:
  - type: Resource
    resource:
      name: cpu
      target:
        type: Utilization
        averageUtilization: 70
  - type: Resource
    resource:
      name: memory
      target:
        type: Utilization
        averageUtilization: 80
```

#### Job Runner Scaling Based on Queue Depth
```yaml
# KEDA ScaledObject (if using KEDA)
apiVersion: keda.sh/v1alpha1
kind: ScaledObject
metadata:
  name: job-runner-scaler
spec:
  scaleTargetRef:
    name: job-runner
  minReplicaCount: 1
  maxReplicaCount: 20
  triggers:
  - type: rabbitmq
    metadata:
      host: amqp://guest:guest@rabbitmq:5672/
      queueName: job_queue
      queueLength: '5'
```

### Storage Configuration

#### RabbitMQ Persistent Storage
```yaml
# Production RabbitMQ with persistence
apiVersion: apps/v1
kind: StatefulSet
metadata:
  name: rabbitmq
spec:
  serviceName: rabbitmq
  replicas: 1
  selector:
    matchLabels:
      app: rabbitmq
  template:
    # ... pod template
  volumeClaimTemplates:
  - metadata:
      name: rabbitmq-data
    spec:
      accessModes: ["ReadWriteOnce"]
      resources:
        requests:
          storage: 10Gi
      storageClassName: fast-ssd
```

### Network Configuration

#### Ingress Configuration
```yaml
# Production ingress with TLS
apiVersion: networking.k8s.io/v1
kind: Ingress
metadata:
  name: microservices-ingress
  annotations:
    kubernetes.io/ingress.class: nginx
    cert-manager.io/cluster-issuer: letsencrypt-prod
    nginx.ingress.kubernetes.io/ssl-redirect: "true"
spec:
  tls:
  - hosts:
    - microservices.example.com
    secretName: microservices-tls
  rules:
  - host: microservices.example.com
    http:
      paths:
      - path: /api
        pathType: Prefix
        backend:
          service:
            name: api-server
            port:
              number: 8081
      - path: /
        pathType: Prefix
        backend:
          service:
            name: frontend
            port:
              number: 8080
```

#### Network Policies
```yaml
# Restrict network access
apiVersion: networking.k8s.io/v1
kind: NetworkPolicy
metadata:
  name: microservices-netpol
spec:
  podSelector: {}
  policyTypes:
  - Ingress
  - Egress
  ingress:
  - from:
    - namespaceSelector:
        matchLabels:
          name: ingress-nginx
  - from:
    - podSelector: {}
  egress:
  - to:
    - podSelector: {}
  - to: []
    ports:
    - protocol: UDP
      port: 53
```

## Production Considerations

### Security Hardening

#### Container Security
```yaml
# Security context for production
securityContext:
  runAsNonRoot: true
  runAsUser: 1000
  runAsGroup: 1000
  fsGroup: 1000
  seccompProfile:
    type: RuntimeDefault
  capabilities:
    drop:
    - ALL
  readOnlyRootFilesystem: true
  allowPrivilegeEscalation: false
```

#### Secrets Management
```bash
# Create secrets from files
kubectl create secret generic app-secrets \
  --from-file=rabbitmq-url=./secrets/rabbitmq-url.txt \
  --from-file=api-key=./secrets/api-key.txt

# Or from command line
kubectl create secret generic app-secrets \
  --from-literal=rabbitmq-url='amqp://user:pass@rabbitmq:5672/' \
  --from-literal=api-key='your-api-key'
```

#### RBAC Configuration
```yaml
# Service account with minimal permissions
apiVersion: v1
kind: ServiceAccount
metadata:
  name: microservices-sa
---
apiVersion: rbac.authorization.k8s.io/v1
kind: Role
metadata:
  name: microservices-role
rules:
- apiGroups: [""]
  resources: ["configmaps", "secrets"]
  verbs: ["get", "list"]
---
apiVersion: rbac.authorization.k8s.io/v1
kind: RoleBinding
metadata:
  name: microservices-binding
subjects:
- kind: ServiceAccount
  name: microservices-sa
roleRef:
  kind: Role
  name: microservices-role
  apiGroup: rbac.authorization.k8s.io
```

### Resource Management

#### Resource Requests and Limits
```yaml
# Production resource configuration
resources:
  requests:
    memory: "128Mi"
    cpu: "100m"
  limits:
    memory: "256Mi"
    cpu: "500m"
```

#### Quality of Service Classes
```yaml
# Guaranteed QoS (requests = limits)
resources:
  requests:
    memory: "256Mi"
    cpu: "200m"
  limits:
    memory: "256Mi"
    cpu: "200m"
```

### Health Checks

#### Liveness and Readiness Probes
```yaml
# Health check configuration
livenessProbe:
  httpGet:
    path: /health
    port: 8081
  initialDelaySeconds: 30
  periodSeconds: 10
  timeoutSeconds: 5
  failureThreshold: 3

readinessProbe:
  httpGet:
    path: /ready
    port: 8081
  initialDelaySeconds: 5
  periodSeconds: 5
  timeoutSeconds: 3
  failureThreshold: 3
```

### Backup and Disaster Recovery

#### RabbitMQ Backup
```bash
# Backup RabbitMQ configuration and data
kubectl exec -it rabbitmq-0 -- rabbitmqctl export_definitions /tmp/backup.json
kubectl cp rabbitmq-0:/tmp/backup.json ./rabbitmq-backup.json

# Backup persistent volume
kubectl create job rabbitmq-backup --image=busybox \
  --restart=OnFailure -- sh -c \
  'tar czf /backup/rabbitmq-$(date +%Y%m%d).tar.gz /data'
```

#### Application State Backup
```bash
# Since this demo uses in-memory storage, implement proper backup
# for production databases/persistent storage
```

## Monitoring and Observability

### Metrics Collection

#### Prometheus Configuration
```yaml
# ServiceMonitor for Prometheus
apiVersion: monitoring.coreos.com/v1
kind: ServiceMonitor
metadata:
  name: microservices-metrics
spec:
  selector:
    matchLabels:
      app: microservices
  endpoints:
  - port: metrics
    interval: 30s
    path: /metrics
```

#### Application Metrics
```go
// Add metrics to your Go applications
import "github.com/prometheus/client_golang/prometheus"

var (
    jobsCreated = prometheus.NewCounterVec(
        prometheus.CounterOpts{
            Name: "jobs_created_total",
            Help: "Total number of jobs created",
        },
        []string{"status"},
    )
    
    jobDuration = prometheus.NewHistogramVec(
        prometheus.HistogramOpts{
            Name: "job_duration_seconds",
            Help: "Job processing duration",
        },
        []string{"status"},
    )
)
```

### Logging

#### Centralized Logging with ELK Stack
```yaml
# Filebeat DaemonSet for log collection
apiVersion: apps/v1
kind: DaemonSet
metadata:
  name: filebeat
spec:
  selector:
    matchLabels:
      name: filebeat
  template:
    spec:
      containers:
      - name: filebeat
        image: elastic/filebeat:7.15.0
        args: [
          "-c", "/etc/filebeat.yml",
          "-e",
        ]
        volumeMounts:
        - name: config
          mountPath: /etc/filebeat.yml
          readOnly: true
          subPath: filebeat.yml
        - name: varlibdockercontainers
          mountPath: /var/lib/docker/containers
          readOnly: true
        - name: varlog
          mountPath: /var/log
          readOnly: true
```

### Distributed Tracing

#### Jaeger Integration
```yaml
# Jaeger all-in-one deployment
apiVersion: apps/v1
kind: Deployment
metadata:
  name: jaeger
spec:
  replicas: 1
  selector:
    matchLabels:
      app: jaeger
  template:
    spec:
      containers:
      - name: jaeger
        image: jaegertracing/all-in-one:1.29
        ports:
        - containerPort: 16686
        - containerPort: 14268
        env:
        - name: COLLECTOR_ZIPKIN_HTTP_PORT
          value: "9411"
```

## Troubleshooting

### Common Deployment Issues

#### Pod Startup Failures
```bash
# Check pod status
kubectl get pods -o wide

# Describe pod for events
kubectl describe pod <pod-name>

# Check pod logs
kubectl logs <pod-name> -f

# Check previous container logs
kubectl logs <pod-name> --previous
```

#### Service Discovery Issues
```bash
# Test service connectivity
kubectl run debug --image=busybox -it --rm --restart=Never -- sh

# Inside debug pod
nslookup api-server
wget -O- http://api-server:8081/health
telnet rabbitmq 5672
```

#### Ingress Issues
```bash
# Check ingress status
kubectl get ingress

# Check ingress controller logs
kubectl logs -n ingress-nginx deployment/nginx-ingress-controller

# Test ingress rules
curl -H "Host: microservices.example.com" http://<ingress-ip>/health
```

### Performance Troubleshooting

#### Resource Usage
```bash
# Check resource usage
kubectl top pods
kubectl top nodes

# Get detailed resource info
kubectl describe node <node-name>
```

#### Application Performance
```bash
# Check application metrics
kubectl port-forward svc/api-server 8081:8081
curl http://localhost:8081/metrics

# Load testing
kubectl run loadtest --image=busybox -it --rm --restart=Never -- sh
# Use wget or similar tools for load testing
```

### Recovery Procedures

#### Rolling Back Deployments
```bash
# Check rollout history
kubectl rollout history deployment/api-server

# Rollback to previous version
kubectl rollout undo deployment/api-server

# Rollback to specific revision
kubectl rollout undo deployment/api-server --to-revision=2
```

#### Emergency Procedures
```bash
# Scale down problematic services
kubectl scale deployment api-server --replicas=0

# Restart all pods
kubectl rollout restart deployment/api-server
kubectl rollout restart deployment/frontend
kubectl rollout restart deployment/job-runner

# Emergency access to nodes
kubectl debug node/<node-name> -it --image=busybox
```

---

## Related Documentation

- [Architecture Documentation](ARCHITECTURE.md) - System design and components
- [API Documentation](API.md) - REST API reference
- [Development Guide](DEVELOPMENT.md) - Local development setup
- [CI/CD Documentation](CICD.md) - Automated deployment pipeline
