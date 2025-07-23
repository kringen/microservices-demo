#!/bin/bash

# Example script to test the enhanced status endpoint with different configurations

echo "Testing AI Research Agent Status Display"
echo "======================================="

echo ""
echo "1. Testing with default configuration (test mode):"
echo "   MCP_TEST_MODE=true (default)"
echo "   Expected: Shows test mode, no MCP endpoints"

echo ""
echo "2. Testing with production mode:"
echo "   MCP_TEST_MODE=false"
echo "   Expected: Shows production mode with MCP endpoints"

echo ""
echo "3. Testing with custom Ollama configuration:"
echo "   OLLAMA_URL=http://custom-ollama:11434"
echo "   OLLAMA_MODEL=llama3.1"
echo "   Expected: Shows custom Ollama endpoint and model"

echo ""
echo "4. Testing with custom MCP endpoints:"
echo "   MCP_WEB_SERVER_URL=http://search-service:8001"
echo "   MCP_GITHUB_SERVER_URL=http://github-service:8002"
echo "   MCP_FILES_SERVER_URL=http://files-service:8003"
echo "   Expected: Shows custom MCP server endpoints"

echo ""
echo "To test these configurations, set the environment variables and start the services:"
echo ""
echo "# Test Mode (default)"
echo "export MCP_TEST_MODE=true"
echo "make run"
echo ""
echo "# Production Mode with custom endpoints"
echo "export MCP_TEST_MODE=false"
echo "export OLLAMA_URL=http://my-ollama-server:11434"
echo "export OLLAMA_MODEL=llama3.1"
echo "export MCP_WEB_SERVER_URL=http://search-service:8001"
echo "export MCP_GITHUB_SERVER_URL=http://github-service:8002"
echo "export MCP_FILES_SERVER_URL=http://files-service:8003"
echo "make run"
echo ""
echo "Then visit http://localhost:8080 to see the enhanced status display!"
