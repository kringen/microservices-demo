#!/bin/bash

# MCP Server Deployment Helper Script

set -e

echo "🔧 MCP Server Deployment Helper"
echo "==============================="

# Function to display usage
usage() {
    echo "Usage: $0 [COMMAND] [OPTIONS]"
    echo ""
    echo "Commands:"
    echo "  dev           Deploy in development mode (test mode, no real MCP servers)"
    echo "  prod          Deploy in production mode with in-cluster MCP servers"
    echo "  external      Deploy in production mode with external MCP servers"
    echo "  status        Check deployment status"
    echo "  secrets       Create MCP secrets (interactive)"
    echo "  clean         Clean up MCP deployments"
    echo ""
    echo "Options:"
    echo "  -n, --namespace   Kubernetes namespace (default: microservices-demo)"
    echo "  -h, --help        Show this help message"
}

# Default values
NAMESPACE="microservices-demo"
COMMAND=""

# Parse arguments
while [[ $# -gt 0 ]]; do
    case $1 in
        dev|prod|external|status|secrets|clean)
            COMMAND="$1"
            shift
            ;;
        -n|--namespace)
            NAMESPACE="$2"
            shift 2
            ;;
        -h|--help)
            usage
            exit 0
            ;;
        *)
            echo "Unknown option: $1"
            usage
            exit 1
            ;;
    esac
done

# Check if kubectl is available
if ! command -v kubectl &> /dev/null; then
    echo "❌ kubectl is not installed or not in PATH"
    exit 1
fi

# Function to create namespace if it doesn't exist
ensure_namespace() {
    if ! kubectl get namespace "$NAMESPACE" &> /dev/null; then
        echo "📁 Creating namespace: $NAMESPACE"
        kubectl create namespace "$NAMESPACE"
    fi
}

# Function to create MCP secrets interactively
create_secrets() {
    echo "🔐 Creating MCP Secrets"
    echo "----------------------"
    
    ensure_namespace
    
    echo "Enter your API credentials (press Enter to skip):"
    
    read -p "GitHub Personal Access Token: " -s GITHUB_TOKEN
    echo ""
    
    read -p "Search API Key (Google/Bing): " -s SEARCH_API_KEY
    echo ""
    
    # Create secret
    kubectl create secret generic mcp-secrets \
        --from-literal=github-token="${GITHUB_TOKEN:-placeholder}" \
        --from-literal=search-api-key="${SEARCH_API_KEY:-placeholder}" \
        --namespace="$NAMESPACE" \
        --dry-run=client -o yaml | kubectl apply -f -
    
    echo "✅ MCP secrets created/updated"
}

# Function to deploy development environment
deploy_dev() {
    echo "🚀 Deploying Development Environment (Test Mode)"
    echo "------------------------------------------------"
    
    ensure_namespace
    
    echo "📦 Applying development configuration..."
    kubectl apply -k k8s/overlays/development
    
    echo "✅ Development deployment complete!"
    echo "💡 In test mode - using simulated MCP data"
}

# Function to deploy production with in-cluster MCP servers
deploy_prod() {
    echo "🚀 Deploying Production Environment (In-Cluster MCP)"
    echo "---------------------------------------------------"
    
    ensure_namespace
    
    # Check if secrets exist
    if ! kubectl get secret mcp-secrets -n "$NAMESPACE" &> /dev/null; then
        echo "⚠️  MCP secrets not found. Creating with placeholders..."
        kubectl create secret generic mcp-secrets \
            --from-literal=github-token="your-github-token-here" \
            --from-literal=search-api-key="your-search-api-key-here" \
            --namespace="$NAMESPACE"
        echo "🔐 Update secrets with: $0 secrets"
    fi
    
    echo "📦 Deploying MCP servers..."
    kubectl apply -f k8s/base/mcp-web-deployment.yaml
    kubectl apply -f k8s/base/mcp-github-deployment.yaml
    kubectl apply -f k8s/base/mcp-files-deployment.yaml
    
    echo "📦 Deploying main application..."
    kubectl apply -k k8s/overlays/production
    
    echo "✅ Production deployment complete!"
    echo "🔗 MCP servers running in-cluster"
}

# Function to deploy production with external MCP servers
deploy_external() {
    echo "🚀 Deploying Production Environment (External MCP)"
    echo "--------------------------------------------------"
    
    ensure_namespace
    
    echo "📦 Applying external production configuration..."
    kubectl apply -k k8s/overlays/production-external
    
    echo "✅ External production deployment complete!"
    echo "🌐 Using external MCP servers - ensure they're accessible"
}

# Function to check deployment status
check_status() {
    echo "📊 Deployment Status"
    echo "-------------------"
    
    echo ""
    echo "🏃 Running Pods:"
    kubectl get pods -n "$NAMESPACE" -l app.kubernetes.io/name=microservices-demo
    
    echo ""
    echo "🔗 Services:"
    kubectl get services -n "$NAMESPACE" -l app.kubernetes.io/name=microservices-demo
    
    echo ""
    echo "⚙️  ConfigMaps:"
    kubectl get configmap app-config -n "$NAMESPACE" -o jsonpath='{.data.MCP_TEST_MODE}' 2>/dev/null && echo " (MCP_TEST_MODE)" || echo ""
    
    echo ""
    echo "🔐 Secrets:"
    kubectl get secrets mcp-secrets -n "$NAMESPACE" &> /dev/null && echo "✅ mcp-secrets exists" || echo "❌ mcp-secrets missing"
}

# Function to clean up deployments
clean_up() {
    echo "🧹 Cleaning Up MCP Deployments"
    echo "------------------------------"
    
    read -p "Are you sure you want to delete all MCP deployments? (y/N): " -n 1 -r
    echo
    
    if [[ $REPLY =~ ^[Yy]$ ]]; then
        echo "🗑️  Deleting MCP servers..."
        kubectl delete -f k8s/base/mcp-web-deployment.yaml --ignore-not-found
        kubectl delete -f k8s/base/mcp-github-deployment.yaml --ignore-not-found  
        kubectl delete -f k8s/base/mcp-files-deployment.yaml --ignore-not-found
        
        echo "🗑️  Deleting main application..."
        kubectl delete -k k8s/overlays/production --ignore-not-found
        kubectl delete -k k8s/overlays/development --ignore-not-found
        kubectl delete -k k8s/overlays/production-external --ignore-not-found
        
        echo "✅ Cleanup complete!"
    else
        echo "❌ Cleanup cancelled"
    fi
}

# Main command execution
case $COMMAND in
    dev)
        deploy_dev
        ;;
    prod)
        deploy_prod
        ;;
    external)
        deploy_external
        ;;
    status)
        check_status
        ;;
    secrets)
        create_secrets
        ;;
    clean)
        clean_up
        ;;
    "")
        echo "❌ No command specified"
        usage
        exit 1
        ;;
    *)
        echo "❌ Unknown command: $COMMAND"
        usage
        exit 1
        ;;
esac
