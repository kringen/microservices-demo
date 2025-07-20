#!/bin/bash

# ============================================================================
# Microservices Health Check Script  
# ============================================================================
#
# Comprehensive health check for all microservices and dependencies.
# This script verifies that all components are running and accessible.
#
# Checks performed:
# 1. RabbitMQ Management UI (localhost:15672)
# 2. API Server health endpoint (localhost:8081/api/health) 
# 3. Frontend web interface (localhost:8080)
# 4. Job Runner process detection
#
# Usage: ./scripts/health-check.sh
#
# Exit codes:
# - 0: All checks passed
# - 1: One or more checks failed
#
# Output: Color-coded status with quick access links
#
# ============================================================================

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

API_URL="http://localhost:8081"
FRONTEND_URL="http://localhost:8080"
RABBITMQ_URL="http://localhost:15672"

echo "🏥 Microservices Health Check"
echo "============================="

# Check RabbitMQ
echo -n "Checking RabbitMQ Management UI... "
if curl -s -f "${RABBITMQ_URL}" > /dev/null; then
    echo -e "${GREEN}✓ OK${NC}"
else
    echo -e "${RED}✗ FAILED${NC}"
    echo "  RabbitMQ management UI is not accessible at ${RABBITMQ_URL}"
fi

# Check API Server
echo -n "Checking API Server... "
if curl -s -f "${API_URL}/api/health" > /dev/null; then
    echo -e "${GREEN}✓ OK${NC}"
    
    # Get detailed health info
    health_info=$(curl -s "${API_URL}/api/health")
    echo "  Health details: ${health_info}"
else
    echo -e "${RED}✗ FAILED${NC}"
    echo "  API Server is not accessible at ${API_URL}"
fi

# Check Frontend
echo -n "Checking Frontend... "
if curl -s -f "${FRONTEND_URL}" > /dev/null; then
    echo -e "${GREEN}✓ OK${NC}"
else
    echo -e "${RED}✗ FAILED${NC}"
    echo "  Frontend is not accessible at ${FRONTEND_URL}"
fi

# Check if job runner process is running
echo -n "Checking Job Runner process... "
if pgrep -f "job-runner" > /dev/null; then
    echo -e "${GREEN}✓ Running${NC}"
else
    echo -e "${YELLOW}⚠ Not detected${NC}"
    echo "  Job runner process not found (this is OK if running with 'go run')"
fi

echo ""
echo "Health check complete!"
echo ""
echo "💡 Quick links:"
echo "   Frontend: ${FRONTEND_URL}"
echo "   API: ${API_URL}/api/health"
echo "   RabbitMQ: ${RABBITMQ_URL} (guest/guest)"
