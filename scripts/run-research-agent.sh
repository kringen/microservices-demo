#!/bin/bash

# Script to run the Research Agent in different modes

echo "Research Agent MCP Mode Selector"
echo "================================"

case "${1:-test}" in
    "test")
        echo "Starting Research Agent in TEST MODE"
        echo "Using simulated MCP data sources..."
        export MCP_TEST_MODE=true
        ;;
    "prod"|"production")
        echo "Starting Research Agent in PRODUCTION MODE"
        echo "Using real MCP server connections..."
        export MCP_TEST_MODE=false
        # Set default MCP server URLs if not already set
        export MCP_WEB_SERVER_URL=${MCP_WEB_SERVER_URL:-http://localhost:3001}
        export MCP_GITHUB_SERVER_URL=${MCP_GITHUB_SERVER_URL:-http://localhost:3002}
        export MCP_FILES_SERVER_URL=${MCP_FILES_SERVER_URL:-http://localhost:3003}
        echo "MCP Server URLs:"
        echo "  Web: $MCP_WEB_SERVER_URL"
        echo "  GitHub: $MCP_GITHUB_SERVER_URL"
        echo "  Files: $MCP_FILES_SERVER_URL"
        ;;
    *)
        echo "Usage: $0 [test|prod]"
        echo "  test - Use simulated data (default)"
        echo "  prod - Use real MCP servers"
        exit 1
        ;;
esac

echo ""
echo "Starting job-runner..."
echo "Press Ctrl+C to stop"

# Run the job-runner
exec ./bin/job-runner
