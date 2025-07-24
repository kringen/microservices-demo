# MCP Server Deployment Configuration Guide

This guide explains how to configure MCP (Model Context Protocol) server addresses and ports for different deployment environments.

## Configuration Overview

The Research Agent supports both **Test Mode** (simulated data) and **Production Mode** (real MCP servers). Configuration is managed through Kubernetes ConfigMaps and environment variables.

## Environment Variables

### Core MCP Configuration
- `MCP_TEST_MODE`: Set to `"true"` for test mode (simulated data), `"false"` for production mode
- `MCP_WEB_SERVER_URL`: URL for web search MCP server
- `MCP_GITHUB_SERVER_URL`: URL for GitHub MCP server  
- `MCP_FILES_SERVER_URL`: URL for file system MCP server

### Additional Configuration
- `MCP_TIMEOUT`: Timeout for MCP server requests (e.g., "120s")
- `OLLAMA_URL`: Ollama AI server endpoint
- `OLLAMA_MODEL`: AI model to use (e.g., "llama3.2")

## Deployment Scenarios

### 1. Development Environment (Test Mode)
**Configuration**: Uses simulated data, no real MCP servers needed.
```yaml
# k8s/overlays/development/configmap-patch.yaml
MCP_TEST_MODE: "true"
```

**Deploy**:
```bash
kubectl apply -k k8s/overlays/development
```

### 2. Production with In-Cluster MCP Servers
**Configuration**: MCP servers deployed within the same Kubernetes cluster.
```yaml
# k8s/overlays/production/configmap-patch.yaml
MCP_TEST_MODE: "false"
MCP_WEB_SERVER_URL: "http://mcp-web-service:3001"
MCP_GITHUB_SERVER_URL: "http://mcp-github-service:3002"
MCP_FILES_SERVER_URL: "http://mcp-files-service:3003"
```

**Deploy**:
```bash
# Deploy MCP servers first
kubectl apply -f k8s/base/mcp-secrets.yaml
kubectl apply -f k8s/base/mcp-web-deployment.yaml
kubectl apply -f k8s/base/mcp-github-deployment.yaml
kubectl apply -f k8s/base/mcp-files-deployment.yaml

# Deploy main application
kubectl apply -k k8s/overlays/production
```

### 3. Production with External MCP Servers
**Configuration**: MCP servers hosted outside the cluster (e.g., managed services).
```yaml
# k8s/overlays/production-external/configmap-patch.yaml
MCP_TEST_MODE: "false"
MCP_WEB_SERVER_URL: "https://mcp-web.internal.company.com:443"
MCP_GITHUB_SERVER_URL: "https://mcp-github.internal.company.com:443"
MCP_FILES_SERVER_URL: "https://mcp-files.internal.company.com:443"
```

**Deploy**:
```bash
kubectl apply -k k8s/overlays/production-external
```

## MCP Server Requirements

Each MCP server must expose a REST API endpoint at `/api/mcp` that accepts:

```json
{
  "method": "search|search_repositories|search_files",
  "params": {
    "query": "search query",
    "limit": 10
  }
}
```

And returns:

```json
{
  "data": "formatted research data",
  "sources": ["list", "of", "source", "urls"],
  "error": "optional error message"
}
```

## Secrets Management

### Required Secrets

Create the `mcp-secrets` secret with necessary API keys:

```bash
kubectl create secret generic mcp-secrets \\
  --from-literal=search-api-key="your-search-api-key" \\
  --from-literal=github-token="your-github-token" \\
  --namespace=microservices-demo
```

### Secrets for Different Providers

**Google Custom Search**:
```yaml
search-api-key: "your-google-custom-search-api-key"
```

**Bing Search API**:
```yaml
search-api-key: "your-bing-search-api-key"
```

**GitHub Personal Access Token**:
```yaml
github-token: "ghp_your-github-personal-access-token"
```

## Port Configuration

### Default Ports
- **Web Search MCP**: Port 3001
- **GitHub MCP**: Port 3002
- **Files MCP**: Port 3003
- **Ollama AI**: Port 11434

### Custom Ports
To use custom ports, update the service definitions:

```yaml
# In mcp-web-deployment.yaml
spec:
  ports:
    - port: 8001  # Custom port
      targetPort: 8001
```

And update the ConfigMap:
```yaml
MCP_WEB_SERVER_URL: "http://mcp-web-service:8001"
```

## Service Discovery

### Kubernetes DNS
Services automatically get DNS names:
- `mcp-web-service.microservices-demo.svc.cluster.local`
- `mcp-github-service.microservices-demo.svc.cluster.local`
- `mcp-files-service.microservices-demo.svc.cluster.local`

### External Services
For external MCP servers, use ExternalName services:

```yaml
apiVersion: v1
kind: Service
metadata:
  name: mcp-web-service
  namespace: microservices-demo
spec:
  type: ExternalName
  externalName: mcp-web.external.com
  ports:
    - port: 443
      targetPort: 443
```

## Monitoring and Health Checks

### Health Check Endpoints
MCP servers should expose:
- `/health` - Liveness probe
- `/ready` - Readiness probe

### Monitoring
Add Prometheus monitoring annotations:

```yaml
metadata:
  annotations:
    prometheus.io/scrape: "true"
    prometheus.io/port: "3001"
    prometheus.io/path: "/metrics"
```

## Scaling Configuration

### Horizontal Pod Autoscaling
```yaml
apiVersion: autoscaling/v2
kind: HorizontalPodAutoscaler
metadata:
  name: mcp-web-server-hpa
spec:
  scaleTargetRef:
    apiVersion: apps/v1
    kind: Deployment
    name: mcp-web-server
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

## Troubleshooting

### Common Issues

1. **MCP Server Connection Failed**
   - Check service DNS resolution
   - Verify port configuration
   - Check firewall rules

2. **Authentication Errors**
   - Verify API keys in secrets
   - Check token permissions
   - Ensure secrets are mounted correctly

3. **Timeout Issues**
   - Increase `MCP_TIMEOUT` value
   - Check network latency
   - Scale up MCP server replicas

### Debug Commands

```bash
# Check MCP server status
kubectl get pods -l app.kubernetes.io/component=mcp-web-server

# View MCP server logs
kubectl logs -l app.kubernetes.io/component=mcp-web-server

# Test MCP server connectivity
kubectl exec -it research-agent-xxx -- curl http://mcp-web-service:3001/health

# Check configuration
kubectl get configmap app-config -o yaml
```

## Security Considerations

### Network Policies
Implement network policies to restrict MCP server access:

```yaml
apiVersion: networking.k8s.io/v1
kind: NetworkPolicy
metadata:
  name: mcp-servers-policy
spec:
  podSelector:
    matchLabels:
      app.kubernetes.io/component: mcp-web-server
  policyTypes:
    - Ingress
  ingress:
    - from:
        - podSelector:
            matchLabels:
              app.kubernetes.io/component: research-agent
      ports:
        - protocol: TCP
          port: 3001
```

### TLS Configuration
For external MCP servers, always use HTTPS:

```yaml
MCP_WEB_SERVER_URL: "https://mcp-web.company.com:443"
```

### API Key Rotation
Regularly rotate API keys and update secrets:

```bash
kubectl patch secret mcp-secrets -p '{"stringData":{"github-token":"new-token"}}'
kubectl rollout restart deployment/mcp-github-server
```
