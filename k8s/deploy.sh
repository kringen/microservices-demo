#!/bin/bash

# Kubernetes deployment script for microservices-demo
# Usage: ./deploy.sh [environment] [action] [registry] [tag]
# Example: ./deploy.sh development apply localhost:5000
# Example: ./deploy.sh production apply registry.company.com v1.2.3

set -e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
K8S_DIR="$SCRIPT_DIR"

# Show usage if requested
if [[ "$1" == "-h" || "$1" == "--help" ]]; then
    echo "🚀 Kubernetes Deployment Script"
    echo "================================"
    echo ""
    echo "Usage: $0 [environment] [action] [registry] [tag] [hostname]"
    echo ""
    echo "Arguments:"
    echo "  environment    Target environment (development|production)"
    echo "  action         Action to perform (apply|delete|diff|build)"  
    echo "  registry       Optional: Container registry URL"
    echo "  tag            Optional: Image tag (default: latest)"
    echo "  hostname       Optional: Hostname for ingress (default: environment-specific)"
    echo ""
    echo "Examples:"
    echo "  $0 development apply"
    echo "  $0 development apply localhost:5000"
    echo "  $0 development apply localhost:5000 v1.2.3"
    echo "  $0 development apply localhost:5000 v1.2.3 microservices-demo.local"
    echo "  $0 production apply registry.company.com"
    echo "  $0 production apply registry.company.com v2.0.0 microservices-demo.kringen.io"
    echo "  $0 production apply registry.company.com v2.0.0"
    echo "  $0 development diff"
    echo "  $0 development build"
    exit 0
fi

# Default values
ENVIRONMENT="${1:-development}"
ACTION="${2:-apply}"
REGISTRY="${3:-}"
TAG="${4:-latest}"
HOSTNAME="${5:-}"

# Set default hostname based on environment if not provided
if [[ -z "$HOSTNAME" ]]; then
    case "$ENVIRONMENT" in
        development)
            HOSTNAME="microservices-demo.local"
            ;;
        production)
            HOSTNAME="microservices-demo.kringen.io"
            ;;
        *)
            HOSTNAME="microservices-demo.local"
            ;;
    esac
fi

# Validate environment
if [[ ! -d "$K8S_DIR/overlays/$ENVIRONMENT" ]]; then
    echo "❌ Environment '$ENVIRONMENT' not found!"
    echo "Available environments:"
    ls -1 "$K8S_DIR/overlays/"
    exit 1
fi

# Validate action
case "$ACTION" in
    apply|delete|diff|build)
        ;;
    *)
        echo "❌ Invalid action '$ACTION'!"
        echo "Valid actions: apply, delete, diff, build"
        exit 1
        ;;
esac

echo "🚀 Kubernetes Deployment Script" >&2
echo "================================" >&2
echo "Environment: $ENVIRONMENT" >&2
echo "Action: $ACTION" >&2
if [[ -n "$REGISTRY" ]]; then
    echo "Registry: $REGISTRY" >&2
fi
echo "Tag: $TAG" >&2
echo "Hostname: $HOSTNAME" >&2
echo "Directory: $K8S_DIR/overlays/$ENVIRONMENT" >&2
echo "" >&2

# Check if kubectl is available
if ! command -v kubectl &> /dev/null; then
    echo "❌ kubectl is not installed or not in PATH" >&2
    exit 1
fi

# Check if kustomize is available
if ! command -v kustomize &> /dev/null; then
    echo "⚠️  kustomize not found, using kubectl kustomize instead" >&2
    KUSTOMIZE_CMD="kubectl kustomize"
else
    KUSTOMIZE_CMD="kustomize"
fi

# Check cluster connectivity
echo "🔍 Checking Kubernetes cluster connectivity..." >&2
if ! kubectl cluster-info &> /dev/null; then
    echo "❌ Cannot connect to Kubernetes cluster" >&2
    echo "Please check your kubeconfig and cluster status" >&2
    exit 1
fi

echo "✅ Connected to cluster: $(kubectl config current-context)" >&2
echo "" >&2

# Function to build manifests
build_manifests() {
    echo "🔨 Building manifests for $ENVIRONMENT..." >&2
    echo "🌐 Using hostname: $HOSTNAME" >&2
    
    if [[ -n "$REGISTRY" ]]; then
        echo "🏷️  Using custom registry: $REGISTRY" >&2
        echo "🏷️  Using image tag: $TAG" >&2
        
        # Build and then replace image references and hostname with sed
        # Redirect warnings to stderr, keep only YAML on stdout
        $KUSTOMIZE_CMD "$K8S_DIR/overlays/$ENVIRONMENT" 2>&1 | \
        sed -e '/^#.*Warning:/d' \
            -e "s|localhost:5000/microservices-\([^:]*\):latest|$REGISTRY/microservices-\1:$TAG|g" \
            -e "s|localhost:5000/microservices-\([^:]*\):\([^[:space:]]*\)|$REGISTRY/microservices-\1:$TAG|g" \
            -e "s|registry\.company\.com/microservices-\([^:]*\):v[0-9.]*|$REGISTRY/microservices-\1:$TAG|g" \
            -e "s|registry\.company\.com/microservices-\([^:]*\):\([^[:space:]]*\)|$REGISTRY/microservices-\1:$TAG|g" \
            -e "s|HOSTNAME_PLACEHOLDER|$HOSTNAME|g"
    else
        echo "🏷️  Using image tag: $TAG" >&2
        
        # Replace only the tag and hostname, keep default registry
        # Redirect warnings to stderr, keep only YAML on stdout
        $KUSTOMIZE_CMD "$K8S_DIR/overlays/$ENVIRONMENT" 2>&1 | \
        sed -e '/^#.*Warning:/d' \
            -e "s|localhost:5000/microservices-\([^:]*\):latest|localhost:5000/microservices-\1:$TAG|g" \
            -e "s|localhost:5000/microservices-\([^:]*\):\([^[:space:]]*\)|localhost:5000/microservices-\1:$TAG|g" \
            -e "s|registry\.company\.com/microservices-\([^:]*\):v[0-9.]*|registry.company.com/microservices-\1:$TAG|g" \
            -e "s|registry\.company\.com/microservices-\([^:]*\):\([^[:space:]]*\)|registry.company.com/microservices-\1:$TAG|g" \
            -e "s|HOSTNAME_PLACEHOLDER|$HOSTNAME|g"
    fi
}

# Function to apply manifests
apply_manifests() {
    echo "📦 Applying manifests for $ENVIRONMENT..." >&2
    
    # First, build and validate
    if ! build_manifests > /tmp/k8s-manifests.yaml; then
        echo "❌ Failed to build manifests" >&2
        exit 1
    fi
    
    # Apply the manifests
    kubectl apply -f /tmp/k8s-manifests.yaml
    
    # Clean up temp file
    rm -f /tmp/k8s-manifests.yaml
    
    echo "" >&2
    echo "✅ Deployment completed!" >&2
    echo "" >&2
    echo "📋 Checking deployment status..." >&2
    kubectl get pods -n microservices-demo -l environment=$ENVIRONMENT
    echo "" >&2
    echo "🔗 Service URLs:" >&2
    
    if [[ "$ENVIRONMENT" == "development" ]]; then
        echo "  Frontend: http://localhost:31080"
        echo "  API Server: http://localhost:31081"
        echo "  RabbitMQ Management: http://localhost:31567"
        echo ""
        echo "💡 For port-forward access:"
        echo "  kubectl port-forward -n microservices-demo svc/dev-frontend-service 8080:8080"
        echo "  kubectl port-forward -n microservices-demo svc/dev-api-server-service 8081:8081"
    else
        echo "  Frontend: http://microservices-demo.local"
        echo "  API Server: http://api.microservices-demo.local"
        echo "  RabbitMQ Management: http://rabbitmq.microservices-demo.local"
        echo ""
        echo "💡 Add to /etc/hosts or configure DNS:"
        echo "  <INGRESS_IP> microservices-demo.local api.microservices-demo.local rabbitmq.microservices-demo.local"
    fi
}

# Function to delete manifests
delete_manifests() {
    echo "🗑️  Deleting manifests for $ENVIRONMENT..."
    
    if build_manifests | kubectl delete -f - 2>/dev/null; then
        echo "✅ Resources deleted successfully"
    else
        echo "⚠️  Some resources may not have been deleted (they might not exist)"
    fi
}

# Function to show diff
show_diff() {
    echo "📊 Showing diff for $ENVIRONMENT..."
    
    if ! build_manifests > /tmp/k8s-manifests.yaml; then
        echo "❌ Failed to build manifests"
        exit 1
    fi
    
    kubectl diff -f /tmp/k8s-manifests.yaml || true
    rm -f /tmp/k8s-manifests.yaml
}

# Execute action
case "$ACTION" in
    apply)
        apply_manifests
        ;;
    delete)
        delete_manifests
        ;;
    diff)
        show_diff
        ;;
    build)
        build_manifests
        ;;
esac

echo "" >&2
echo "🎉 Operation completed successfully!" >&2
