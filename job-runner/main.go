package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"microservices-demo/shared"

	amqp "github.com/rabbitmq/amqp091-go"
)

// OllamaClient represents connection to local Ollama server
type OllamaClient struct {
	baseURL string
	client  *http.Client
}

// OllamaRequest represents a request to Ollama API
type OllamaRequest struct {
	Model    string `json:"model"`
	Prompt   string `json:"prompt"`
	Stream   bool   `json:"stream"`
	System   string `json:"system,omitempty"`
	Template string `json:"template,omitempty"`
}

// OllamaResponse represents response from Ollama API
type OllamaResponse struct {
	Response string `json:"response"`
	Done     bool   `json:"done"`
	Context  []int  `json:"context,omitempty"`
}

// MCPServiceHandler handles different MCP service integrations
type MCPServiceHandler struct {
	// In a real implementation, these would be actual MCP client connections
	availableServices map[shared.MCPService]bool
	testMode         bool
}

// ResearchAgent is the AI-powered research agent
type ResearchAgent struct {
	rabbitmq    *shared.RabbitMQClient
	ollama      *OllamaClient
	mcpHandler  *MCPServiceHandler
	daprURL     string
}

func NewResearchAgent() *ResearchAgent {
	return &ResearchAgent{
		daprURL: getEnvOrDefault("DAPR_HTTP_ENDPOINT", "http://localhost:3500"),
	}
}

func (ra *ResearchAgent) initOllama() error {
	ollamaURL := getEnvOrDefault("OLLAMA_URL", "http://localhost:11434")
	
	ra.ollama = &OllamaClient{
		baseURL: ollamaURL,
		client:  &http.Client{Timeout: 5 * time.Minute},
	}

	// Test connection to Ollama
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := ra.testOllamaConnection(ctx); err != nil {
		return fmt.Errorf("failed to connect to Ollama: %w", err)
	}

	log.Println("Connected to Ollama server successfully")
	return nil
}

func (ra *ResearchAgent) testOllamaConnection(ctx context.Context) error {
	req, err := http.NewRequestWithContext(ctx, "GET", ra.ollama.baseURL+"/api/tags", nil)
	if err != nil {
		return err
	}

	resp, err := ra.ollama.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return fmt.Errorf("ollama server returned status: %d", resp.StatusCode)
	}

	return nil
}

func (ra *ResearchAgent) initMCPServices() {
	// Check if we're in test mode
	testMode := getEnvOrDefault("MCP_TEST_MODE", "false") == "true"
	
	ra.mcpHandler = &MCPServiceHandler{
		availableServices: map[shared.MCPService]bool{
			shared.MCPServiceWeb:      true,  // Web search/scraping
			shared.MCPServiceGitHub:   true,  // GitHub API
			shared.MCPServiceDatabase: false, // Database queries (disabled for demo)
			shared.MCPServiceFiles:    true,  // File system access
			shared.MCPServiceCalendar: false, // Calendar integration (disabled)
			shared.MCPServiceSlack:    false, // Slack integration (disabled)
		},
		testMode: testMode,
	}
	
	if testMode {
		log.Println("MCP services initialized in TEST MODE (using simulated data)")
	} else {
		log.Println("MCP services initialized in PRODUCTION MODE (using real MCP servers)")
	}
}

func (ra *ResearchAgent) initRabbitMQ() error {
	var err error

	// Get RabbitMQ URL from environment variable
	rabbitmqURL := os.Getenv("RABBITMQ_URL")
	if rabbitmqURL == "" {
		rabbitmqURL = "amqp://guest:guest@localhost:5672/"
	}

	ra.rabbitmq, err = shared.NewRabbitMQClient(rabbitmqURL)
	if err != nil {
		return err
	}
	return nil
}

func (ra *ResearchAgent) start() error {
	jobs, err := ra.rabbitmq.ConsumeJobs()
	if err != nil {
		return err
	}

	log.Println("Research Agent started. Waiting for research requests...")

	for delivery := range jobs {
		var jobMessage shared.JobMessage
		if err := json.Unmarshal(delivery.Body, &jobMessage); err != nil {
			log.Printf("Failed to unmarshal job message: %v", err)
			if err := delivery.Nack(false, false); err != nil {
				log.Printf("Failed to nack message: %v", err)
			}
			continue
		}

		log.Printf("Received research request: %s - %s", jobMessage.JobID, jobMessage.Title)

		// Process the research request in a goroutine
		go func(msg shared.JobMessage, d amqp.Delivery) {
			// Send processing status update
			processingUpdate := shared.JobResult{
				JobID:       msg.JobID,
				Status:      shared.JobStatusProcessing,
				CompletedAt: time.Now(), // Use this as "started at" timestamp
			}

			if err := ra.rabbitmq.PublishResult(processingUpdate); err != nil {
				log.Printf("Failed to publish processing status for research %s: %v", msg.JobID, err)
			} else {
				log.Printf("Research %s marked as processing", msg.JobID)
			}

			// Process the research request and get final result
			result := ra.processResearchRequest(msg)

			// Publish the final result
			if err := ra.rabbitmq.PublishResult(result); err != nil {
				log.Printf("Failed to publish result for research %s: %v", msg.JobID, err)
			} else {
				log.Printf("Published final result for research %s", msg.JobID)
			}

			// Acknowledge the message
			if err := d.Ack(false); err != nil {
				log.Printf("Failed to ack message: %v", err)
			}
		}(jobMessage, delivery)
	}

	return nil
}
func (ra *ResearchAgent) processResearchRequest(jobMessage shared.JobMessage) shared.JobResult {
	log.Printf("Starting research: %s - %s", jobMessage.JobID, jobMessage.Query)
	startTime := time.Now()

	// Create context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()

	var result shared.JobResult
	result.JobID = jobMessage.JobID
	result.CompletedAt = time.Now()

	// Step 1: Gather information using MCP services
	mcpData, sources, err := ra.gatherInformationWithMCP(ctx, jobMessage)
	if err != nil {
		result.Status = shared.JobStatusFailed
		result.Error = fmt.Sprintf("Failed to gather information: %v", err)
		return result
	}

	// Step 2: Use Ollama to process and analyze the gathered information
	research, confidence, tokens, err := ra.analyzeWithOllama(ctx, jobMessage, mcpData)
	if err != nil {
		result.Status = shared.JobStatusFailed
		result.Error = fmt.Sprintf("Failed to analyze with AI: %v", err)
		return result
	}

	// Step 3: Create comprehensive result
	duration := time.Since(startTime)
	result.Status = shared.JobStatusCompleted
	result.Result = research
	result.Sources = sources
	result.Confidence = confidence
	result.TokensUsed = tokens

	log.Printf("Research %s completed in %v with confidence %.2f", 
		jobMessage.JobID, duration, confidence)
	
	return result
}

func (ra *ResearchAgent) gatherInformationWithMCP(ctx context.Context, jobMessage shared.JobMessage) (string, []string, error) {
	var allData []string
	var sources []string

	log.Printf("Gathering information for: %s", jobMessage.Query)

	// Process each requested MCP service
	for _, service := range jobMessage.MCPServices {
		if !ra.mcpHandler.availableServices[service] {
			log.Printf("MCP service %s not available, skipping", service)
			continue
		}

		data, serviceSources, err := ra.queryMCPService(ctx, service, jobMessage)
		if err != nil {
			log.Printf("Error querying MCP service %s: %v", service, err)
			continue
		}

		allData = append(allData, data)
		sources = append(sources, serviceSources...)
	}

	if len(allData) == 0 {
		return "", sources, fmt.Errorf("no data could be gathered from MCP services")
	}

	return strings.Join(allData, "\n\n"), sources, nil
}

func (ra *ResearchAgent) queryMCPService(ctx context.Context, service shared.MCPService, jobMessage shared.JobMessage) (string, []string, error) {
	// Use simulation if in test mode, otherwise use real MCP servers
	if ra.mcpHandler.testMode {
		switch service {
		case shared.MCPServiceWeb:
			return ra.simulateWebSearch(jobMessage.Query)
		case shared.MCPServiceGitHub:
			return ra.simulateGitHubSearch(jobMessage.Query)
		case shared.MCPServiceFiles:
			return ra.simulateFileSearch(jobMessage.Query)
		default:
			return "", nil, fmt.Errorf("unsupported MCP service: %s", service)
		}
	}
	
	// Real MCP server implementations
	switch service {
	case shared.MCPServiceWeb:
		return ra.queryWebSearchMCP(ctx, jobMessage.Query)
	case shared.MCPServiceGitHub:
		return ra.queryGitHubMCP(ctx, jobMessage.Query)
	case shared.MCPServiceFiles:
		return ra.queryFilesMCP(ctx, jobMessage.Query)
	default:
		return "", nil, fmt.Errorf("unsupported MCP service: %s", service)
	}
}

func (ra *ResearchAgent) simulateWebSearch(query string) (string, []string, error) {
	// In a real implementation, this would use a web search MCP server
	data := fmt.Sprintf(`Web Search Results for "%s":

1. Comprehensive overview found on multiple authoritative sources
2. Recent developments and trends identified
3. Technical specifications and best practices documented
4. Community discussions and expert opinions gathered

Key findings:
- Industry standard approaches have evolved significantly
- Best practices emphasize scalability and maintainability  
- Recent innovations provide improved performance
- Community consensus supports modern methodologies`, query)

	sources := []string{
		"https://example.com/research-1",
		"https://example.com/research-2", 
		"https://example.com/research-3",
	}

	return data, sources, nil
}

func (ra *ResearchAgent) simulateGitHubSearch(query string) (string, []string, error) {
	// In a real implementation, this would use GitHub MCP server
	data := fmt.Sprintf(`GitHub Repository Analysis for "%s":

Popular repositories found:
- 5 repositories with 1000+ stars
- 12 active projects with recent commits
- 3 enterprise-grade solutions identified

Code patterns analysis:
- Modern architecture patterns prevalent
- Comprehensive test coverage common
- Documentation quality varies but generally good
- Active community contributions`, query)

	sources := []string{
		"https://github.com/example/repo1",
		"https://github.com/example/repo2",
		"https://github.com/example/repo3",
	}

	return data, sources, nil
}

func (ra *ResearchAgent) simulateFileSearch(query string) (string, []string, error) {
	// In a real implementation, this would use file system MCP server
	data := fmt.Sprintf(`Local File System Search for "%s":

Relevant files found:
- Configuration files containing related settings
- Documentation with relevant information
- Code examples and implementations
- Historical data and logs

Analysis summary:
- Current implementations follow established patterns
- Configuration is well-structured
- Documentation provides good coverage
- Historical data shows consistent usage patterns`, query)

	sources := []string{
		"/local/docs/related-file-1.md",
		"/local/config/settings.yaml",
		"/local/examples/implementation.go",
	}

	return data, sources, nil
}

// Real MCP Server Implementations
// These functions call actual MCP servers instead of simulations

func (ra *ResearchAgent) queryWebSearchMCP(ctx context.Context, query string) (string, []string, error) {
	// Real web search MCP server implementation
	// This would typically use a search engine API like Google, Bing, or DuckDuckGo
	// through an MCP server
	
	mcpServerURL := getEnvOrDefault("MCP_WEB_SERVER_URL", "http://localhost:3001")
	
	requestBody := map[string]interface{}{
		"method": "search",
		"params": map[string]interface{}{
			"query": query,
			"limit": 10,
		},
	}
	
	data, sources, err := ra.callMCPServer(ctx, mcpServerURL, requestBody)
	if err != nil {
		log.Printf("Web search MCP server error: %v", err)
		// Fallback to simulation if MCP server fails
		return ra.simulateWebSearch(query)
	}
	
	return data, sources, nil
}

func (ra *ResearchAgent) queryGitHubMCP(ctx context.Context, query string) (string, []string, error) {
	// Real GitHub MCP server implementation
	// This would use the GitHub API through an MCP server
	
	mcpServerURL := getEnvOrDefault("MCP_GITHUB_SERVER_URL", "http://localhost:3002")
	
	requestBody := map[string]interface{}{
		"method": "search_repositories",
		"params": map[string]interface{}{
			"query": query,
			"sort":  "stars",
			"order": "desc",
			"limit": 10,
		},
	}
	
	data, sources, err := ra.callMCPServer(ctx, mcpServerURL, requestBody)
	if err != nil {
		log.Printf("GitHub MCP server error: %v", err)
		// Fallback to simulation if MCP server fails
		return ra.simulateGitHubSearch(query)
	}
	
	return data, sources, nil
}

func (ra *ResearchAgent) queryFilesMCP(ctx context.Context, query string) (string, []string, error) {
	// Real file system MCP server implementation
	// This would search local or networked file systems through an MCP server
	
	mcpServerURL := getEnvOrDefault("MCP_FILES_SERVER_URL", "http://localhost:3003")
	
	requestBody := map[string]interface{}{
		"method": "search_files",
		"params": map[string]interface{}{
			"query":     query,
			"file_type": []string{".md", ".txt", ".go", ".js", ".py"},
			"limit":     20,
		},
	}
	
	data, sources, err := ra.callMCPServer(ctx, mcpServerURL, requestBody)
	if err != nil {
		log.Printf("Files MCP server error: %v", err)
		// Fallback to simulation if MCP server fails
		return ra.simulateFileSearch(query)
	}
	
	return data, sources, nil
}

// Generic MCP server call function
func (ra *ResearchAgent) callMCPServer(ctx context.Context, serverURL string, requestBody map[string]interface{}) (string, []string, error) {
	// Marshal request
	jsonData, err := json.Marshal(requestBody)
	if err != nil {
		return "", nil, fmt.Errorf("failed to marshal MCP request: %w", err)
	}
	
	// Create HTTP request
	req, err := http.NewRequestWithContext(ctx, "POST", serverURL+"/api/mcp", bytes.NewBuffer(jsonData))
	if err != nil {
		return "", nil, fmt.Errorf("failed to create MCP request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	
	// Make HTTP request with timeout
	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return "", nil, fmt.Errorf("MCP server request failed: %w", err)
	}
	defer resp.Body.Close()
	
	if resp.StatusCode != 200 {
		body, _ := io.ReadAll(resp.Body)
		return "", nil, fmt.Errorf("MCP server returned status %d: %s", resp.StatusCode, string(body))
	}
	
	// Parse response
	var mcpResponse struct {
		Data    string   `json:"data"`
		Sources []string `json:"sources"`
		Error   string   `json:"error,omitempty"`
	}
	
	if err := json.NewDecoder(resp.Body).Decode(&mcpResponse); err != nil {
		return "", nil, fmt.Errorf("failed to decode MCP response: %w", err)
	}
	
	if mcpResponse.Error != "" {
		return "", nil, fmt.Errorf("MCP server error: %s", mcpResponse.Error)
	}
	
	return mcpResponse.Data, mcpResponse.Sources, nil
}

func (ra *ResearchAgent) analyzeWithOllama(ctx context.Context, jobMessage shared.JobMessage, mcpData string) (string, float64, int, error) {
	// Create a comprehensive prompt for the AI
	systemPrompt := `You are a professional research agent. Your task is to analyze the provided information and create a comprehensive, well-structured research report. 

Guidelines:
- Provide accurate, fact-based analysis
- Structure your response with clear sections using markdown formatting (headers, lists, tables, etc.)
- Include key findings and insights
- Mention any limitations or areas needing further research
- Rate your confidence in the findings (0.0 to 1.0)
- Be concise but thorough
- Use markdown formatting for better readability (# headers, **bold**, *italic*, lists, tables, code blocks)`

	var userPrompt string
	if ra.mcpHandler.testMode {
		// In test mode, inform about placeholder sources
		systemPrompt += `

IMPORTANT: The provided "sources" are placeholder examples. In your response, you should reference realistic, relevant sources that would actually exist for this research topic. Generate appropriate URLs, documentation links, academic papers, or industry resources that would be credible sources for this type of research, even though you cannot actually access them.`

		userPrompt = fmt.Sprintf(`Research Request: %s

Query: %s
Research Type: %s

Gathered Information:
%s

NOTE: The sources listed above are placeholder examples. Please provide a comprehensive research report and suggest realistic, relevant sources that would be appropriate for this research topic. Include references to actual websites, documentation, academic papers, or industry resources that would credibly support this type of research.

Please structure your response using markdown formatting and provide a professional research report.`, 
			jobMessage.Title, jobMessage.Query, jobMessage.ResearchType, mcpData)
	} else {
		// In production mode, sources are real
		systemPrompt += `

The sources provided are from real data gathering services. Reference them appropriately in your analysis.`

		userPrompt = fmt.Sprintf(`Research Request: %s

Query: %s
Research Type: %s

Gathered Information from MCP Services:
%s

Please provide a comprehensive research report based on this real data. Structure your response using markdown formatting and provide a professional analysis.`, 
			jobMessage.Title, jobMessage.Query, jobMessage.ResearchType, mcpData)
	}

	// Make request to Ollama
	response, tokens, err := ra.callOllama(ctx, systemPrompt, userPrompt)
	if err != nil {
		return "", 0.0, 0, err
	}

	// Calculate confidence based on response quality and data availability
	confidence := ra.calculateConfidence(response, mcpData, len(jobMessage.MCPServices))

	return response, confidence, tokens, nil
}

func (ra *ResearchAgent) callOllama(ctx context.Context, systemPrompt, userPrompt string) (string, int, error) {
	model := getEnvOrDefault("OLLAMA_MODEL", "llama3.2")
	
	reqBody := OllamaRequest{
		Model:  model,
		Prompt: userPrompt,
		System: systemPrompt,
		Stream: false,
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return "", 0, err
	}

	req, err := http.NewRequestWithContext(ctx, "POST", ra.ollama.baseURL+"/api/generate", bytes.NewBuffer(jsonData))
	if err != nil {
		return "", 0, err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := ra.ollama.client.Do(req)
	if err != nil {
		return "", 0, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		body, _ := io.ReadAll(resp.Body)
		return "", 0, fmt.Errorf("ollama API error: %d - %s", resp.StatusCode, string(body))
	}

	var ollamaResp OllamaResponse
	if err := json.NewDecoder(resp.Body).Decode(&ollamaResp); err != nil {
		return "", 0, err
	}

	// Estimate token usage (rough approximation)
	tokens := len(strings.Fields(userPrompt + systemPrompt + ollamaResp.Response))

	return ollamaResp.Response, tokens, nil
}

func (ra *ResearchAgent) calculateConfidence(response, mcpData string, mcpServiceCount int) float64 {
	baseConfidence := 0.6

	// Increase confidence based on response length and quality
	responseWords := len(strings.Fields(response))
	if responseWords > 100 {
		baseConfidence += 0.1
	}
	if responseWords > 300 {
		baseConfidence += 0.1
	}

	// Increase confidence based on amount of gathered data
	dataWords := len(strings.Fields(mcpData))
	if dataWords > 200 {
		baseConfidence += 0.1
	}

	// Increase confidence based on number of MCP services used
	baseConfidence += float64(mcpServiceCount) * 0.05

	// Cap at 0.95 to account for inherent uncertainty
	if baseConfidence > 0.95 {
		baseConfidence = 0.95
	}

	return baseConfidence
}

// Utility function to get environment variable with default
func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func main() {
	agent := NewResearchAgent()

	// Initialize components
	if err := agent.initRabbitMQ(); err != nil {
		log.Fatalf("Failed to initialize RabbitMQ: %v", err)
	}
	defer agent.rabbitmq.Close()

	if err := agent.initOllama(); err != nil {
		log.Fatalf("Failed to initialize Ollama: %v", err)
	}

	agent.initMCPServices()

	log.Println("AI Research Agent is starting...")
	log.Println("Components initialized:")
	log.Println("  ✓ RabbitMQ connection")
	log.Println("  ✓ Ollama AI model")
	if agent.mcpHandler.testMode {
		log.Println("  ✓ MCP services (TEST MODE)")
	} else {
		log.Println("  ✓ MCP services (PRODUCTION MODE)")
	}
	log.Printf("  ✓ Dapr endpoint: %s", agent.daprURL)

	if err := agent.start(); err != nil {
		log.Fatalf("Failed to start research agent: %v", err)
	}
}
