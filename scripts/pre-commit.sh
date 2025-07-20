#!/bin/bash

# Pre-commit checks script
# Runs the same checks as CI to ensure your code will pass

set -e  # Exit on any error

echo "🔍 Running pre-commit checks..."
echo "==============================="

# Color codes for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Function to print status
print_status() {
    if [ $1 -eq 0 ]; then
        echo -e "${GREEN}✅ $2${NC}"
    else
        echo -e "${RED}❌ $2${NC}"
        exit 1
    fi
}

# Change to project root
cd "$(dirname "$0")/.."

echo -e "${BLUE}📁 Working directory: $(pwd)${NC}"
echo

# 1. Download dependencies
echo -e "${YELLOW}📦 Downloading dependencies...${NC}"
go mod download
print_status $? "Dependencies downloaded"

# 2. Verify dependencies
echo -e "${YELLOW}🔍 Verifying dependencies...${NC}"
go mod verify
print_status $? "Dependencies verified"

# 3. Run go vet
echo -e "${YELLOW}🕵️  Running go vet...${NC}"
go vet ./...
print_status $? "go vet passed"

# 4. Check formatting
echo -e "${YELLOW}📝 Checking code formatting...${NC}"
UNFORMATTED=$(gofmt -s -l . | grep -v vendor | wc -l)
if [ "$UNFORMATTED" -gt 0 ]; then
    echo -e "${RED}❌ The following files are not formatted:${NC}"
    gofmt -s -l . | grep -v vendor
    echo -e "${YELLOW}💡 Run 'gofmt -s -w .' to fix formatting${NC}"
    exit 1
else
    print_status 0 "Code formatting check passed"
fi

# 5. Run golangci-lint (if available)
echo -e "${YELLOW}🔬 Running golangci-lint...${NC}"
if command -v golangci-lint >/dev/null 2>&1; then
    golangci-lint run ./...
    print_status $? "golangci-lint passed"
else
    echo -e "${YELLOW}⚠️  golangci-lint not found, installing...${NC}"
    go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
    if command -v golangci-lint >/dev/null 2>&1; then
        golangci-lint run ./...
        print_status $? "golangci-lint passed"
    else
        echo -e "${YELLOW}⚠️  golangci-lint installation failed, skipping...${NC}"
    fi
fi

# 6. Run tests
echo -e "${YELLOW}🧪 Running tests...${NC}"

# Check if RabbitMQ is running
echo -e "${BLUE}🐰 Checking RabbitMQ availability...${NC}"
if ! curl -f http://localhost:15672 >/dev/null 2>&1; then
    echo -e "${YELLOW}⚠️  RabbitMQ not running, starting with Docker...${NC}"
    make rabbitmq-up >/dev/null 2>&1 || {
        echo -e "${YELLOW}⚠️  Could not start RabbitMQ, running tests without integration tests${NC}"
        echo -e "${BLUE}🧪 Running unit tests only...${NC}"
        go test -v -short ./...
        print_status $? "Unit tests passed"
    }
else
    echo -e "${GREEN}✅ RabbitMQ is running${NC}"
fi

# Run full test suite if RabbitMQ is available
if curl -f http://localhost:15672 >/dev/null 2>&1; then
    echo -e "${BLUE}🧪 Running full test suite...${NC}"
    
    # Root module tests
    echo -e "${BLUE}  Testing root module...${NC}"
    go test -v -race ./...
    print_status $? "Root module tests passed"
    
    # Individual service tests
    echo -e "${BLUE}  Testing API server...${NC}"
    (cd api-server && go test -v -race ./...)
    print_status $? "API server tests passed"
    
    echo -e "${BLUE}  Testing frontend...${NC}"
    (cd frontend && go test -v -race ./...)
    print_status $? "Frontend tests passed"
    
    echo -e "${BLUE}  Testing job runner...${NC}"
    (cd job-runner && go test -v -race ./...)
    print_status $? "Job runner tests passed"
    
    echo -e "${BLUE}  Testing shared package...${NC}"
    (cd shared && go test -v -race ./...)
    print_status $? "Shared package tests passed"
fi

echo
echo -e "${GREEN}🎉 All pre-commit checks passed!${NC}"
echo -e "${GREEN}✨ Your code is ready to commit!${NC}"
echo
