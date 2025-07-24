# MCP (Model Context Protocol) Server Configuration

The Research Agent can work in two modes, and the status is now displayed in the web UI.

## Status Display

The AI Research Agent Status panel now shows:
- **API Server**: Health status
- **RabbitMQ**: Connection status  
- **Ollama AI**: Model name and endpoint URL
- **MCP Services**: Mode and endpoint URLs (if in production mode)

## Test Mode (Default)
Uses simulated data sources with placeholder URLs.

To enable test mode:
```bash
export MCP_TEST_MODE=true
```

Status display will show: **MCP Services: Test Mode (Simulated)**

## Production Mode
Connects to real MCP servers for data gathering.

To enable production mode:
```bash
export MCP_TEST_MODE=false
```

Status display will show: **MCP Services: Production Mode** with individual endpoint URLs listed.

### MCP Server URLs

Configure the URLs for your MCP servers:

```bash
# Web Search MCP Server (e.g., using DuckDuckGo, Google Custom Search, etc.)
export MCP_WEB_SERVER_URL=http://localhost:3001

# GitHub MCP Server (for repository searches)
export MCP_GITHUB_SERVER_URL=http://localhost:3002

# File System MCP Server (for local file searches)
export MCP_FILES_SERVER_URL=http://localhost:3003
```

### MCP Server API Format

Each MCP server should expose a POST endpoint at `/api/mcp` that accepts:

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

### Fallback Behavior

If a real MCP server is unavailable, the Research Agent will automatically fall back to simulated data to ensure the system remains functional.

### Example MCP Servers

You can create MCP servers using:
- **Web Search**: Use search APIs like DuckDuckGo, Google Custom Search, or Bing
- **GitHub**: Use the GitHub REST API or GraphQL API
- **Files**: Create a file indexing service that searches local or networked files
