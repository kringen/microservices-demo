#!/bin/bash

# Pre-commit script to run the same checks as CI
# Usage: ./scripts/pre-commit.sh

set -e

echo "🔍 Running pre-commit checks..."
echo "================================"

# Change to project directory
cd "$(dirname "$0")/.."

echo ""
echo "📝 Checking Go code formatting..."
if ! make fmt-check; then
    echo "❌ Code formatting check failed!"
    echo "💡 Run 'make fmt' to fix formatting issues."
    exit 1
fi
echo "✅ Code formatting check passed!"

echo ""
echo "🔧 Running go vet..."
if ! make vet; then
    echo "❌ go vet failed!"
    exit 1
fi
echo "✅ go vet passed!"

echo ""
echo "🧪 Running tests..."
if ! make test; then
    echo "❌ Tests failed!"
    exit 1
fi
echo "✅ All tests passed!"

echo ""
echo "🔍 Running linter..."
if ! make lint; then
    echo "❌ Linter failed!"
    exit 1
fi
echo "✅ Linter passed!"

echo ""
echo "🎉 All pre-commit checks passed!"
echo "Your code is ready for commit and CI will likely pass."
