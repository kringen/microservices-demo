#!/bin/bash

# Kubernetes deployment script for microservices-demo
# Usage: ./deploy.sh --environment development --action apply --ollama local
# Usage: ./deploy.sh -e development -a apply -o 192.168.1.100:11434 --registry kringen --tag v1.2.3

set -e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
K8S_DIR="$SCRIPT_DIR"

# Default values
ENVIRONMENT="development"
ACTION="apply"
OLLAMA_TYPE="local"
REGISTRY=""
TAG="latest"
HOSTNAME=""
FORCE=false

# Show usage function
show_usage() {
    echo "🚀 Kubernetes Deployment Script"
    echo "================================"
    echo ""
    echo "Usage: $0 [OPTIONS]"
    echo ""
    echo "Required Options:"
    echo "  -e, --environment ENV    Target environment (development|production)"
    echo "  -a, --action ACTION      Action to perform (apply|delete|diff|build)"
    echo ""
    echo "Optional Parameters:"
    echo "  -o, --ollama TYPE        Ollama deployment type:"
    echo "                             'local' - Deploy Ollama as a pod in the cluster"
    echo "                             '<host:port>' - Use external Ollama server"
    echo "  -r, --registry URL       Container registry URL"
    echo "  -t, --tag TAG            Image tag (default: latest)"
    echo "  -h, --hostname HOST      Hostname for ingress (default: environment-specific)"
    echo "      --force              Force deployment by deleting conflicting resources"
    echo "      --help               Show this help message"
    echo ""
    echo "Ollama Deployment Types:"
    echo "  local                    Deploy Ollama as a pod in the cluster (requires 4-8GB memory)"
    echo "  <host:port>             Use external Ollama server at specified address"
    echo "  192.168.1.100:11434     Example: External server on local network"
    echo "  ollama.company.com:11434 Example: External server with hostname"
    echo ""
    echo "Examples:"
    echo "  $0 --environment development --action apply --ollama local"
    echo "  $0 -e development -a apply -o 192.168.1.100:11434"
    echo "  $0 -e development -a apply -o local -r localhost:5000 -t v1.2.3"
    echo "  $0 -e development -a apply -o 10.0.1.50:11434 -r localhost:5000 -t v1.2.3 -h microservices-demo.local"
    echo "  $0 -e production -a apply -o ollama.company.com:11434 -r registry.company.com -t v2.0.0"
    echo "  $0 -e production -a diff -o local"
    echo "  $0 -e development -a build -o 192.168.1.100:11434"
    echo ""
    echo "Short form examples:"
    echo "  $0 -e dev -a apply -o local"
    echo "  $0 -e prod -a apply -o 192.168.1.100:11434 -r kringen -t latest"
    echo ""
    echo "Force deployment (handles label conflicts):"
    echo "  $0 -e dev -a apply -o 192.168.1.100:11434 --force"
    echo ""
    echo "Note: External Ollama servers must be accessible from the Kubernetes cluster."
}

# Parse command line arguments
PARSED_ARGS=$(getopt -o e:a:o:r:t:h: --long environment:,action:,ollama:,registry:,tag:,hostname:,force,help -n "$0" -- "$@")
if [[ $? -ne 0 ]]; then
    echo "❌ Invalid arguments. Use --help for usage information." >&2
    exit 1
fi

eval set -- "$PARSED_ARGS"

while true; do
    case "$1" in
        -e|--environment)
            ENVIRONMENT="$2"
            shift 2
            ;;
        -a|--action)
            ACTION="$2"
            shift 2
            ;;
        -o|--ollama)
            OLLAMA_TYPE="$2"
            shift 2
            ;;
        -r|--registry)
            REGISTRY="$2"
            shift 2
            ;;
        -t|--tag)
            TAG="$2"
            shift 2
            ;;
        -h|--hostname)
            HOSTNAME="$2"
            shift 2
            ;;
        --force)
            FORCE=true
            shift
            ;;
        --help)
            show_usage
            exit 0
            ;;
        --)
            shift
            break
            ;;
        *)
            echo "❌ Unknown option: $1" >&2
            echo "Use --help for usage information." >&2
            exit 1
            ;;
    esac
done

# Handle short environment names
case "$ENVIRONMENT" in
    dev|development)
        ENVIRONMENT="development"
        ;;
    prod|production)
        ENVIRONMENT="production"
        ;;
esac

# Validate required parameters
if [[ -z "$ENVIRONMENT" ]]; then
    echo "❌ Environment is required. Use -e or --environment." >&2
    echo "Use --help for usage information." >&2
    exit 1
fi

if [[ -z "$ACTION" ]]; then
    echo "❌ Action is required. Use -a or --action." >&2
    echo "Use --help for usage information." >&2
    exit 1
fi

# Parse Ollama type to determine if it's local or external
if [[ "$OLLAMA_TYPE" == "local" ]]; then
    DEPLOYMENT_TYPE="local"
    OVERLAY_PATH="$ENVIRONMENT"
    OLLAMA_URL=""
else
    DEPLOYMENT_TYPE="external"
    OVERLAY_PATH="$ENVIRONMENT-external"
    # Validate external server format (host:port)
    if [[ ! "$OLLAMA_TYPE" =~ ^[a-zA-Z0-9.-]+:[0-9]+$ ]]; then
        echo "❌ Invalid external Ollama server format: '$OLLAMA_TYPE'" >&2
        echo "Expected format: <host:port> (e.g., '192.168.1.100:11434')" >&2
        exit 1
    fi
    OLLAMA_URL="http://$OLLAMA_TYPE"
fi

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

# Validate environment and Ollama type combination
if [[ ! -d "$K8S_DIR/overlays/$OVERLAY_PATH" ]]; then
    echo "❌ Overlay '$OVERLAY_PATH' not found!" >&2
    echo "Available overlays:" >&2
    ls -1 "$K8S_DIR/overlays/" >&2
    echo "" >&2
    echo "Valid combinations:" >&2
    echo "  --environment development --ollama local           → overlays/development/" >&2
    echo "  --environment development --ollama <host:port>     → overlays/development-external/" >&2
    echo "  --environment production --ollama local            → overlays/production/" >&2
    echo "  --environment production --ollama <host:port>      → overlays/production-external/" >&2
    exit 1
fi

# Validate Ollama type
if [[ "$DEPLOYMENT_TYPE" != "local" && "$DEPLOYMENT_TYPE" != "external" ]]; then
    echo "❌ Invalid Ollama type '$OLLAMA_TYPE'!" >&2
    echo "Valid types: 'local' or '<host:port>' (e.g., '192.168.1.100:11434')" >&2
    exit 1
fi

# Validate action
case "$ACTION" in
    apply|delete|diff|build)
        ;;
    *)
        echo "❌ Invalid action '$ACTION'!" >&2
        echo "Valid actions: apply, delete, diff, build" >&2
        exit 1
        ;;
esac

echo "🚀 Kubernetes Deployment Script" >&2
echo "================================" >&2
echo "Environment: $ENVIRONMENT" >&2
echo "Deployment Type: $DEPLOYMENT_TYPE" >&2
if [[ "$DEPLOYMENT_TYPE" == "external" ]]; then
    echo "Ollama Server: $OLLAMA_URL" >&2
fi
echo "Overlay Path: $OVERLAY_PATH" >&2
echo "Action: $ACTION" >&2
if [[ -n "$REGISTRY" ]]; then
    echo "Registry: $REGISTRY" >&2
fi
echo "Tag: $TAG" >&2
echo "Hostname: $HOSTNAME" >&2
echo "Directory: $K8S_DIR/overlays/$OVERLAY_PATH" >&2
echo "" >&2

# Show Ollama deployment info
if [[ "$DEPLOYMENT_TYPE" == "local" ]]; then
    echo "🤖 Deploying with local Ollama (in-cluster AI)" >&2
    echo "   • Requires 4-8GB memory for AI workloads" >&2
    echo "   • Initial deployment takes 5-10 minutes (model download)" >&2
    echo "   • Persistent storage required for models" >&2
else
    echo "🌐 Deploying with external Ollama server: $OLLAMA_URL" >&2
    echo "   • Server must be accessible from Kubernetes cluster" >&2
    echo "   • Reduced cluster resource requirements" >&2
    echo "   • Will automatically configure OLLAMA_URL in deployment" >&2
fi
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

# Function to handle deployment conflicts
handle_deployment_conflicts() {
    echo "🔍 Checking for deployment label conflicts..." >&2
    
    # Get deployments that might have conflicting selectors
    local conflicting_deployments=$(kubectl get deployments -n microservices-demo -o jsonpath='{.items[?(@.metadata.labels.environment=="'$ENVIRONMENT'")].metadata.name}' 2>/dev/null || true)
    
    if [[ -n "$conflicting_deployments" && "$FORCE" == "true" ]]; then
        echo "⚠️  Force mode enabled: deleting conflicting deployments..." >&2
        echo "   Deployments to delete: $conflicting_deployments" >&2
        
        for deployment in $conflicting_deployments; do
            echo "   Deleting deployment: $deployment" >&2
            kubectl delete deployment "$deployment" -n microservices-demo --ignore-not-found=true
        done
        
        echo "✅ Conflicting deployments removed" >&2
        echo "" >&2
    elif [[ -n "$conflicting_deployments" ]]; then
        echo "⚠️  Detected existing deployments that may conflict:" >&2
        echo "   $conflicting_deployments" >&2
        echo "" >&2
        echo "💡 If you encounter selector immutable errors, try:" >&2
        echo "   $0 $* --force" >&2
        echo "" >&2
    fi
}

# Function to build manifests
build_manifests() {
    echo "🔨 Building manifests for $ENVIRONMENT with $DEPLOYMENT_TYPE Ollama..." >&2
    echo "🌐 Using hostname: $HOSTNAME" >&2
    
    if [[ -n "$REGISTRY" ]]; then
        echo "🏷️  Using custom registry: $REGISTRY" >&2
        echo "🏷️  Using image tag: $TAG" >&2
        
        # Build and then replace image references and hostname with sed
        # Redirect warnings to stderr, keep only YAML on stdout
        MANIFEST_OUTPUT=$($KUSTOMIZE_CMD "$K8S_DIR/overlays/$OVERLAY_PATH" 2>&1 | \
        sed -e '/^#.*Warning:/d' \
            -e "s|localhost:5000/microservices-\([^:]*\):latest|$REGISTRY/microservices-\1:$TAG|g" \
            -e "s|localhost:5000/microservices-\([^:]*\):\([^[:space:]]*\)|$REGISTRY/microservices-\1:$TAG|g" \
            -e "s|registry\.company\.com/microservices-\([^:]*\):v[0-9.]*|$REGISTRY/microservices-\1:$TAG|g" \
            -e "s|registry\.company\.com/microservices-\([^:]*\):\([^[:space:]]*\)|$REGISTRY/microservices-\1:$TAG|g" \
            -e "s|HOSTNAME_PLACEHOLDER|$HOSTNAME|g")
    else
        echo "🏷️  Using image tag: $TAG" >&2
        
        # Replace only the tag and hostname, keep default registry
        # Redirect warnings to stderr, keep only YAML on stdout
        MANIFEST_OUTPUT=$($KUSTOMIZE_CMD "$K8S_DIR/overlays/$OVERLAY_PATH" 2>&1 | \
        sed -e '/^#.*Warning:/d' \
            -e "s|localhost:5000/microservices-\([^:]*\):latest|localhost:5000/microservices-\1:$TAG|g" \
            -e "s|localhost:5000/microservices-\([^:]*\):\([^[:space:]]*\)|localhost:5000/microservices-\1:$TAG|g" \
            -e "s|registry\.company\.com/microservices-\([^:]*\):v[0-9.]*|registry.company.com/microservices-\1:$TAG|g" \
            -e "s|registry\.company\.com/microservices-\([^:]*\):\([^[:space:]]*\)|registry.company.com/microservices-\1:$TAG|g" \
            -e "s|HOSTNAME_PLACEHOLDER|$HOSTNAME|g")
    fi
    
    # If using external Ollama, replace the placeholder OLLAMA_URL with the actual server
    if [[ "$DEPLOYMENT_TYPE" == "external" ]]; then
        echo "🔗 Setting external Ollama URL: $OLLAMA_URL" >&2
        echo "$MANIFEST_OUTPUT" | sed "s|http://your-ollama-server:11434|$OLLAMA_URL|g" | \
            sed "s|http://192\.168\.1\.100:11434|$OLLAMA_URL|g" | \
            sed "s|http://ollama\.internal\.company\.com:11434|$OLLAMA_URL|g"
    else
        echo "$MANIFEST_OUTPUT"
    fi
}

# Function to apply manifests
apply_manifests() {
    echo "📦 Applying manifests for $ENVIRONMENT with $DEPLOYMENT_TYPE Ollama..." >&2
    
    # Handle potential deployment conflicts
    handle_deployment_conflicts
    
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
    
    # Show deployment-specific information
    if [[ "$DEPLOYMENT_TYPE" == "local" ]]; then
        echo "🤖 Local Ollama deployment notes:" >&2
        echo "   • Model download may take 5-10 minutes on first deployment" >&2
        echo "   • Check Ollama init container logs: kubectl logs -n microservices-demo deployment/ollama -c ollama-init" >&2
        echo "   • Verify model ready: kubectl exec -n microservices-demo deployment/ollama -- ollama list" >&2
    else
        echo "🌐 External Ollama deployment notes:" >&2
        echo "   • Using external Ollama server: $OLLAMA_URL" >&2
        echo "   • Ensure server is running and accessible from cluster" >&2
        echo "   • Test connectivity: kubectl exec -n microservices-demo deployment/research-agent -- curl $OLLAMA_URL/api/tags" >&2
    fi
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
    echo "🗑️  Deleting manifests for $ENVIRONMENT with $DEPLOYMENT_TYPE Ollama..." >&2
    
    if build_manifests | kubectl delete -f - 2>/dev/null; then
        echo "✅ Resources deleted successfully" >&2
        
        if [[ "$DEPLOYMENT_TYPE" == "local" ]]; then
            echo "💾 Note: Ollama model data persists in PVC (ollama-models-pvc)" >&2
            echo "    To also delete model data: kubectl delete pvc ollama-models-pvc -n microservices-demo" >&2
        fi
    else
        echo "⚠️  Some resources may not have been deleted (they might not exist)" >&2
    fi
}

# Function to show diff
show_diff() {
    echo "📊 Showing diff for $ENVIRONMENT with $DEPLOYMENT_TYPE Ollama..." >&2
    
    if ! build_manifests > /tmp/k8s-manifests.yaml; then
        echo "❌ Failed to build manifests" >&2
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
