# Kubernetes AI Research Agent Deployment Guide

This document explains how to deploy the AI Research Agent system with different Ollama configurations.

## Deployment Options

### 1. Local Ollama (In-Cluster)
Deploys Ollama as a pod within the Kubernetes cluster with local model storage.

**Pros:**
- No external dependencies
- Consistent performance
- Full control over AI models

**Cons:**
- Requires cluster resources (4-8GB memory)
- Longer initial deployment (model download)
- Storage requirements (10GB for models)

### 2. External Ollama (Network Server)
Uses an external Ollama server running on the same network as your Kubernetes nodes.

**Pros:**
- No cluster resource overhead for AI
- Shared AI server across multiple workloads
- Faster cluster deployments

**Cons:**
- Network dependency
- External server maintenance
- Potential network latency

## New Resources Added

### 1. Ollama AI Server (`ollama-deployment.yaml`)
- **Purpose**: Provides local LLM capabilities using the llama3.2 model
- **Key Features**:
  - Init container for automatic model download
  - Persistent volume claim for 10GB model storage
  - Health checks using `ollama list` command
  - Resource limits: 4-8GB memory for AI workloads
  - Non-root security context

### 2. Ollama Service (`ollama-service.yaml`)
- **Purpose**: Exposes Ollama on port 11434 for internal communication
- **Configuration**: ClusterIP service for pod-to-pod communication

### 3. Ollama Persistent Volume Claim (`ollama-pvc.yaml`)
- **Purpose**: Provides persistent storage for LLM models
- **Configuration**: 10GB storage with ReadWriteOnce access

## Updated Resources

### 1. Research Agent Deployment (formerly `job-runner-deployment.yaml`)
- **Name Change**: `job-runner` → `research-agent`
- **Container Updates**:
  - Updated labels and selectors
  - Added AI-specific environment variables:
    - `OLLAMA_URL`
    - `OLLAMA_MODEL`
    - `MCP_TIMEOUT`
    - `ENABLE_AI_FEATURES`
- **Resource Updates**:
  - Increased memory: 512Mi request, 1Gi limit
  - Increased CPU: 200m request, 1000m limit

### 2. ConfigMap (`configmap.yaml`)
- **Added AI Configuration**:
  - `OLLAMA_URL`: "http://ollama-service:11434"
  - `OLLAMA_MODEL`: "llama3.2"
  - `MCP_TIMEOUT`: "120s"
  - `ENABLE_AI_FEATURES`: "true"
- **Updated Job Settings**:
  - Increased `JOB_TIMEOUT` to 300s for AI processing
  - Reduced `MAX_CONCURRENT_JOBS` to 5 for AI workloads

### 3. Kustomization (`kustomization.yaml`)
- **Added Resources**:
  - `ollama-pvc.yaml`
  - `ollama-deployment.yaml`
  - `ollama-service.yaml`
- **Added Image**: `ollama/ollama:latest`

## Environment-Specific Updates

### Development Overlay
- **Resource Reductions**: Smaller memory/CPU limits for development
- **Service Names**: Prefixed with `dev-` for proper service discovery
- **AI Settings**: Reduced timeouts and concurrent jobs for development
- **Ollama Resources**: 2-4GB memory for development workloads

### Production Overlay
- **Full Resources**: Production-ready memory and CPU limits
- **Scaling**: 
  - Research agent: 2 replicas (reduced from 5 for AI workloads)
  - Ollama: 1 replica (single instance)
- **Timeouts**: Increased for production AI processing
- **Ollama Resources**: 4-8GB memory for production workloads

## Deployment Instructions

### Local Ollama Deployment

#### Development Environment
```bash
# Deploy development with local Ollama
kubectl apply -k k8s/overlays/development/

# Wait for Ollama model download (may take 5-10 minutes)
kubectl logs -n microservices-demo deployment/dev-ollama -c ollama-init -f

# Verify AI system is ready
kubectl get pods -n microservices-demo
kubectl logs -n microservices-demo deployment/dev-research-agent
```

#### Production Environment
```bash
# Deploy production with local Ollama
kubectl apply -k k8s/overlays/production/

# Monitor deployment
kubectl rollout status deployment/ollama -n microservices-demo
kubectl rollout status deployment/research-agent -n microservices-demo

# Verify AI features
kubectl exec -n microservices-demo deployment/ollama -- ollama list
```

### External Ollama Deployment

#### Development Environment
```bash
# 1. First, ensure your external Ollama server is running and accessible
# Test connectivity from your cluster:
kubectl run test-pod --image=curlimages/curl --rm -it -- curl http://YOUR_OLLAMA_SERVER:11434/api/tags

# 2. Update the Ollama URL in the configmap patch
# Edit k8s/overlays/development-external/configmap-patch.yaml
# Update OLLAMA_URL to your server's address

# 3. Deploy development with external Ollama
kubectl apply -k k8s/overlays/development-external/

# 4. Verify connectivity
kubectl logs -n microservices-demo deployment/dev-research-agent
```

#### Production Environment
```bash
# 1. Update production Ollama URL
# Edit k8s/overlays/production-external/configmap-patch.yaml
# Update OLLAMA_URL to your production server

# 2. Deploy production with external Ollama
kubectl apply -k k8s/overlays/production-external/

# 3. Monitor deployment
kubectl rollout status deployment/research-agent -n microservices-demo

# 4. Test AI functionality
kubectl exec -n microservices-demo deployment/research-agent -- curl http://YOUR_OLLAMA_SERVER:11434/api/tags
```

## Available Overlays

| Overlay | Description | Ollama Location | Resource Usage |
|---------|-------------|-----------------|----------------|
| `development` | Dev environment with local Ollama | In-cluster | High (includes AI) |
| `development-external` | Dev environment with external Ollama | External server | Low (no AI overhead) |
| `production` | Prod environment with local Ollama | In-cluster | High (includes AI) |
| `production-external` | Prod environment with external Ollama | External server | Medium (no AI overhead) |
| `local-ollama` | Base local Ollama setup | In-cluster | High |
| `external-ollama` | Base external Ollama setup | External server | Low |

## Key Considerations

### Local Ollama Deployment
1. **Storage**: Ensure your cluster has adequate storage for the 10GB model PVC
2. **Resources**: AI workloads require significant memory - ensure nodes have sufficient capacity (4-8GB)
3. **Init Time**: First deployment will take longer due to model download (2-10 minutes)
4. **Scaling**: Ollama should remain at 1 replica to avoid model download conflicts

### External Ollama Deployment
1. **Network**: Ensure external Ollama server is accessible from cluster nodes
2. **Firewall**: Open port 11434 on external server for cluster access
3. **Performance**: Consider network latency for AI API calls
4. **Reliability**: External server becomes a single point of failure

### General
1. **Monitoring**: Monitor resource usage and network connectivity
2. **Security**: Consider network policies and access controls for AI services

## Troubleshooting

### Common Issues
1. **Model Download Failure**: Check init container logs and network connectivity
2. **Out of Memory**: Increase node capacity or reduce resource requests
3. **Storage Issues**: Verify storage class exists and has sufficient capacity
4. **Service Discovery**: Ensure DNS is working for service-to-service communication

### Verification Commands
```bash
# Check Ollama health
kubectl exec deployment/ollama -n microservices-demo -- ollama list

# Check research agent connectivity
kubectl logs deployment/research-agent -n microservices-demo

# Monitor resource usage
kubectl top pods -n microservices-demo
```
