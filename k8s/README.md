# Kubernetes Deployment

This directory contains Kubernetes manifests for deploying the microservices demo application to Kubernetes clusters.

## 🏗️ Architecture

The Kubernetes deployment includes:

- **Namespace**: Isolated environment for the application
- **Secrets**: Secure storage for RabbitMQ credentials
- **ConfigMaps**: Application configuration
- **Deployments**: Container orchestration for each service
- **Services**: Network access and service discovery
- **Ingress**: External access routing
- **PersistentVolumeClaims**: Storage for RabbitMQ data

## 📁 Structure

```
k8s/
├── base/                    # Base Kustomize configuration
│   ├── namespace.yaml       # Application namespace
│   ├── secrets.yaml         # RabbitMQ credentials
│   ├── configmap.yaml       # Application configuration
│   ├── rabbitmq-*           # RabbitMQ deployment and service
│   ├── api-server-*         # API server deployment and service
│   ├── job-runner-*         # Job runner deployment
│   ├── frontend-*           # Frontend deployment and service
│   ├── ingress.yaml         # External access routing
│   └── kustomization.yaml   # Base Kustomize config
├── overlays/                # Environment-specific configurations
│   ├── development/         # Development environment
│   │   ├── kustomization.yaml
│   │   ├── deployment-patches.yaml
│   │   └── configmap-patch.yaml
│   └── production/          # Production environment
│       ├── kustomization.yaml
│       ├── deployment-patches.yaml
│       ├── configmap-patch.yaml
│       └── secrets-patch.yaml
├── deploy.sh               # Deployment script
└── README.md               # This file
```

## 🚀 Quick Start

### Prerequisites

- Kubernetes cluster (v1.20+)
- kubectl configured
- kustomize (optional, kubectl has built-in support)
- Ingress controller (nginx recommended)

#### Deploy Script

The `deploy.sh` script simplifies deployment across environments with flexible Ollama configurations:

```bash
# Show help and available options
./deploy.sh --help

# Deploy development with local Ollama
./deploy.sh --environment development --action apply --ollama local

# Deploy production with external Ollama  
./deploy.sh --environment production --action apply --ollama 192.168.1.100:11434

# Build manifests without applying (dry-run)
./deploy.sh --environment development --action build --ollama local

# Clean up deployment
./deploy.sh --environment development --action delete --ollama local

# Using short options
./deploy.sh -e dev -a apply -o local

# Deploy with custom registry and tag
./deploy.sh -e prod -a apply -o local -r my-registry.com -t v1.2.3

# Production with external hostname
./deploy.sh -e prod -a apply -o local -h my-ai-system.example.com
```

### Parameters

**Required:**
- `-e, --environment`: `development` (or `dev`) / `production` (or `prod`)
- `-a, --action`: `apply`, `delete`, `diff`, or `build`

**Optional:**
- `-o, --ollama`: `local` (default) or `<host:port>` for external Ollama
- `-r, --registry`: Container registry (default: docker.io)
- `-t, --tag`: Image tag (default: latest)  
- `-h, --hostname`: External hostname for production ingress
- `--help`: Show usage information

## 🛠️ Deployment Script Usage

The `deploy.sh` script provides a convenient way to manage deployments:

```bash
# Syntax
./deploy.sh [environment] [action] [registry]

# Examples
./deploy.sh development apply                    # Deploy with default registry
./deploy.sh development apply localhost:5000     # Deploy with local registry
./deploy.sh production apply registry.company.com # Deploy with company registry
./deploy.sh development delete                   # Delete development environment
./deploy.sh development diff                     # Show what would change
./deploy.sh development build                    # Build and show manifests
./deploy.sh --help                              # Show help information
```

### Available Environments
- `development` - Single replicas, debug mode, lower resources, uses `localhost:5000` registry
- `production` - Multiple replicas, release mode, production resources, uses `registry.company.com` registry

### Available Actions
- `apply` - Deploy or update the environment
- `delete` - Remove the environment
- `diff` - Show differences that would be applied
- `build` - Generate and display the final manifests

### Registry Configuration
The deployment script supports custom container registries:

- **Default**: Uses registry specified in kustomization.yaml files
- **Override**: Pass registry as third parameter to override default
- **Development**: Defaults to `localhost:5000` (for local development)
- **Production**: Defaults to `registry.company.com` (configurable)

## 🎯 Manual Deployment with kubectl

If you prefer manual control:

```bash
# Development
kubectl kustomize overlays/development | kubectl apply -f -

# Production
kubectl kustomize overlays/production | kubectl apply -f -

# Delete
kubectl kustomize overlays/development | kubectl delete -f -
```

## 🔧 Service Configuration

### Resource Allocation

#### Development Environment
- **API Server**: 1 replica, 64Mi-128Mi memory, 50m-200m CPU
- **Job Runner**: 1 replica, 64Mi-128Mi memory, 50m-200m CPU  
- **Frontend**: 1 replica, 64Mi-128Mi memory, 50m-200m CPU
- **RabbitMQ**: 1 replica, 256Mi-512Mi memory, 100m-500m CPU

#### Production Environment
- **API Server**: 3 replicas, 256Mi-512Mi memory, 200m-1000m CPU
- **Job Runner**: 5 replicas, 256Mi-512Mi memory, 200m-1000m CPU
- **Frontend**: 3 replicas, 256Mi-512Mi memory, 200m-1000m CPU
- **RabbitMQ**: 1 replica, 512Mi memory, 500m CPU

### Network Access

#### Development (NodePort)
- Frontend: `localhost:31080`
- API Server: `localhost:31081`
- RabbitMQ Management: `localhost:31567`

#### Production (Ingress)
- Frontend: `microservices-demo.local`
- API Server: `api.microservices-demo.local`
- RabbitMQ Management: `rabbitmq.microservices-demo.local`

## 🔐 Security Features

- **Non-root containers**: All services run as non-root users
- **Security contexts**: Restricted capabilities and privileges
- **Network policies**: (Can be added for additional security)
- **Secret management**: Credentials stored in Kubernetes secrets
- **Resource limits**: CPU and memory limits prevent resource exhaustion

## 📊 Monitoring & Observability

### Health Checks
- **Liveness probes**: Restart containers if unhealthy
- **Readiness probes**: Remove from service if not ready
- **Startup probes**: Allow for longer startup times

### Metrics Collection
- Prometheus annotations for scraping
- Health check endpoints exposed
- Service metrics available

### Logging
- Structured logging with configurable levels
- Centralized log collection ready
- Debug mode available for development

## 🎛️ Configuration Management

### Container Images
The application uses custom-built container images that can be configured per environment:

```yaml
# Base configuration (base/kustomization.yaml)
images:
  - name: microservices-api-server
    newName: localhost:5000/microservices-api-server
    newTag: latest
  - name: microservices-job-runner
    newName: localhost:5000/microservices-job-runner
    newTag: latest
  - name: microservices-frontend
    newName: localhost:5000/microservices-frontend
    newTag: latest
```

#### Registry Options
- **Local Development**: `localhost:5000` (requires local registry)
- **Docker Hub**: `docker.io/username` or just `username`
- **Google Container Registry**: `gcr.io/project-id`
- **Amazon ECR**: `123456789012.dkr.ecr.region.amazonaws.com`
- **Azure Container Registry**: `myregistry.azurecr.io`
- **Private Registry**: `registry.company.com`

#### Building and Pushing Images
Before deploying, ensure your images are built and available:

```bash
# Build images (from project root)
docker build -f api-server/Dockerfile -t localhost:5000/microservices-api-server:latest .
docker build -f job-runner/Dockerfile -t localhost:5000/microservices-job-runner:latest .
docker build -f frontend/Dockerfile -t localhost:5000/microservices-frontend:latest .

# Push to registry
docker push localhost:5000/microservices-api-server:latest
docker push localhost:5000/microservices-job-runner:latest
docker push localhost:5000/microservices-frontend:latest

# Or use deployment script with custom registry
./deploy.sh development apply your-registry.com
```

### Environment Variables
All configuration is managed through ConfigMaps and Secrets:

```yaml
# ConfigMap - Application settings
GIN_MODE: "release"
LOG_LEVEL: "info"
API_SERVER_URL: "http://api-server-service:8081"
MAX_CONCURRENT_JOBS: "20"

# Secret - Credentials
RABBITMQ_USERNAME: "guest"
RABBITMQ_PASSWORD: "guest"
```

### Customization
To customize for your environment:

1. **Update base configuration**: Modify `base/configmap.yaml`
2. **Configure container registry**: Update `images` section in kustomization files
3. **Create environment overlay**: Add new directory under `overlays/`
4. **Update image tags**: Modify image tags in kustomization files
5. **Adjust resources**: Update deployment patches

#### Example: Adding a New Environment

```bash
# Create new environment overlay
mkdir -p overlays/staging
cp -r overlays/development/* overlays/staging/

# Update the kustomization.yaml
cat > overlays/staging/kustomization.yaml << EOF
apiVersion: kustomize.config.k8s.io/v1beta1
kind: Kustomization

resources:
  - ../../base

namePrefix: staging-

labels:
  - includeSelectors: true
    pairs:
      environment: staging

patchesStrategicMerge:
  - deployment-patches.yaml
  - configmap-patch.yaml

images:
  - name: microservices-api-server
    newName: registry.company.com/microservices-api-server
    newTag: staging-v1.2.0
  - name: microservices-job-runner
    newName: registry.company.com/microservices-job-runner
    newTag: staging-v1.2.0
  - name: microservices-frontend
    newName: registry.company.com/microservices-frontend
    newTag: staging-v1.2.0
EOF

# Deploy staging environment
./deploy.sh staging apply
```

## 🔄 Scaling

### Manual Scaling
```bash
# Scale API server
kubectl scale deployment api-server -n microservices-demo --replicas=5

# Scale job runners for high load
kubectl scale deployment job-runner -n microservices-demo --replicas=10

# Scale frontend
kubectl scale deployment frontend -n microservices-demo --replicas=3
```

### Horizontal Pod Autoscaling
```yaml
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
```

## 🗄️ Storage

### RabbitMQ Persistence
- **PersistentVolumeClaim**: 2Gi storage for RabbitMQ data
- **StorageClass**: Uses `standard` (adjust for your cluster)
- **Backup**: Consider regular backups for production

### StatefulSet Alternative
For production RabbitMQ clustering, consider using StatefulSet:

```bash
# Example StatefulSet deployment (not included but recommended for production)
# - Multiple RabbitMQ nodes
# - Persistent storage per node
# - Automatic clustering
# - Rolling updates
```

## 🌐 Ingress Configuration

### nginx-ingress
```bash
# Install nginx-ingress controller
kubectl apply -f https://raw.githubusercontent.com/kubernetes/ingress-nginx/controller-v1.8.2/deploy/static/provider/cloud/deploy.yaml

# Or using Helm
helm upgrade --install ingress-nginx ingress-nginx \
  --repo https://kubernetes.github.io/ingress-nginx \
  --namespace ingress-nginx --create-namespace
```

### DNS Configuration
Add to `/etc/hosts` for local testing:
```
127.0.0.1 microservices-demo.local
127.0.0.1 api.microservices-demo.local
127.0.0.1 rabbitmq.microservices-demo.local
```

## 🔍 Troubleshooting

### Common Issues

#### Pods Not Starting
```bash
# Check pod status
kubectl get pods -n microservices-demo

# View pod logs
kubectl logs -n microservices-demo deployment/api-server

# Describe pod for events
kubectl describe pod -n microservices-demo <pod-name>
```

#### Image Pull Issues
```bash
# Check if images exist in registry
docker pull localhost:5000/microservices-api-server:latest

# Check image pull secrets (if using private registry)
kubectl get secrets -n microservices-demo

# Check image names in deployment
kubectl get deployment dev-api-server -n microservices-demo -o yaml | grep image

# Common solutions:
# 1. Build and push images to registry
# 2. Update imagePullPolicy to IfNotPresent for local images
# 3. Add image pull secrets for private registries
```

#### Service Connectivity Issues
```bash
# Test service connectivity
kubectl run test-pod -n microservices-demo --image=busybox --rm -it -- sh

# Inside the pod:
nslookup api-server-service
wget -qO- http://api-server-service:8081/api/health
```

#### Storage Issues
```bash
# Check PVC status
kubectl get pvc -n microservices-demo

# Check storage class
kubectl get storageclass
```

#### Ingress Not Working
```bash
# Check ingress status
kubectl get ingress -n microservices-demo

# Check ingress controller
kubectl get pods -n ingress-nginx

# Check ingress controller logs
kubectl logs -n ingress-nginx deployment/ingress-nginx-controller
```

### Debug Commands

```bash
# Port forward for local access
kubectl port-forward -n microservices-demo svc/frontend-service 8080:8080
kubectl port-forward -n microservices-demo svc/api-server-service 8081:8081
kubectl port-forward -n microservices-demo svc/rabbitmq-service 15672:15672

# View all resources
kubectl get all -n microservices-demo

# View events
kubectl get events -n microservices-demo --sort-by='.lastTimestamp'

# Execute into container
kubectl exec -it -n microservices-demo deployment/api-server -- sh
```

## 🎯 Production Considerations

### Security Hardening
- Use proper secrets management (HashiCorp Vault, AWS Secrets Manager)
- Implement network policies
- Enable Pod Security Standards
- Use service mesh for mTLS (Istio, Linkerd)

### High Availability
- Deploy across multiple availability zones
- Use anti-affinity rules for pod distribution
- Implement proper backup strategies
- Use external load balancers

### Monitoring
- Deploy Prometheus and Grafana
- Set up alerting rules
- Use distributed tracing (Jaeger, Zipkin)
- Implement log aggregation (ELK, Loki)

### Performance
- Use resource quotas and limits
- Implement caching strategies
- Consider using admission controllers
- Regular performance testing
