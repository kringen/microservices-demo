package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"sort"
	"sync"
	"time"

	"microservices-demo/shared"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type APIServer struct {
	jobs      map[string]*shared.Job
	jobsMutex sync.RWMutex
	rabbitmq  *shared.RabbitMQClient
}

func NewAPIServer() *APIServer {
	return &APIServer{
		jobs: make(map[string]*shared.Job),
	}
}

func (s *APIServer) initRabbitMQ() error {
	var err error

	// Get RabbitMQ URL from environment variable
	rabbitmqURL := os.Getenv("RABBITMQ_URL")
	if rabbitmqURL == "" {
		rabbitmqURL = "amqp://guest:guest@localhost:5672/"
	}

	s.rabbitmq, err = shared.NewRabbitMQClient(rabbitmqURL)
	if err != nil {
		return err
	}

	// Start consuming job results
	go s.consumeJobResults()

	return nil
}

func (s *APIServer) consumeJobResults() {
	results, err := s.rabbitmq.ConsumeResults()
	if err != nil {
		log.Printf("Failed to consume results: %v", err)
		return
	}

	for delivery := range results {
		var result shared.JobResult
		if err := json.Unmarshal(delivery.Body, &result); err != nil {
			log.Printf("Failed to unmarshal job result: %v", err)
			continue
		}

		s.updateJobStatus(result)
	}
}

func (s *APIServer) updateJobStatus(result shared.JobResult) {
	s.jobsMutex.Lock()
	defer s.jobsMutex.Unlock()

	if job, exists := s.jobs[result.JobID]; exists {
		// Update job status and timing
		previousStatus := job.Status
		job.Status = result.Status
		job.Result = result.Result
		job.Error = result.Error
		job.Sources = result.Sources
		job.Confidence = result.Confidence
		job.TokensUsed = result.TokensUsed

		// Handle different status updates
		switch result.Status {
		case shared.JobStatusProcessing:
			// Set started time when research begins processing
			if job.StartedAt == nil {
				startTime := result.CompletedAt // Using CompletedAt as the timestamp for when processing started
				job.StartedAt = &startTime
				log.Printf("Research %s started processing at %v", result.JobID, startTime)
			}
		case shared.JobStatusCompleted, shared.JobStatusFailed:
			// Set completion time for final states
			job.CompletedAt = &result.CompletedAt

			// Calculate and log processing duration
			if job.StartedAt != nil {
				duration := result.CompletedAt.Sub(*job.StartedAt)
				log.Printf("Research %s completed in %v with confidence %.2f", result.JobID, duration, result.Confidence)
			}
		}

		// Log status change for monitoring
		if previousStatus != result.Status {
			log.Printf("Research %s status changed: %s -> %s", result.JobID, previousStatus, result.Status)
		}
	} else {
		log.Printf("Received result for unknown research: %s", result.JobID)
	}
}

func (s *APIServer) createJob(c *gin.Context) {
	var req shared.ResearchRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	job := &shared.Job{
		ID:           uuid.New().String(),
		Title:        req.Title,
		Query:        req.Query,
		ResearchType: req.ResearchType,
		MCPServices:  req.MCPServices,
		Status:       shared.JobStatusPending,
		CreatedAt:    time.Now(),
	}

	s.jobsMutex.Lock()
	s.jobs[job.ID] = job
	s.jobsMutex.Unlock()

	// Send research request to queue (only if RabbitMQ is initialized)
	if s.rabbitmq != nil {
		jobMessage := shared.JobMessage{
			JobID:        job.ID,
			Title:        job.Title,
			Query:        job.Query,
			ResearchType: job.ResearchType,
			MCPServices:  job.MCPServices,
		}

		if err := s.rabbitmq.PublishJob(jobMessage); err != nil {
			log.Printf("Failed to publish research request: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to queue research request"})
			return
		}
	} else {
		log.Println("RabbitMQ not initialized - research request not queued (test mode?)")
	}

	c.JSON(http.StatusCreated, job)
}

func (s *APIServer) getJob(c *gin.Context) {
	jobID := c.Param("id")

	s.jobsMutex.RLock()
	job, exists := s.jobs[jobID]
	s.jobsMutex.RUnlock()

	if !exists {
		c.JSON(http.StatusNotFound, gin.H{"error": "Job not found"})
		return
	}

	c.JSON(http.StatusOK, job)
}

func (s *APIServer) listJobs(c *gin.Context) {
	s.jobsMutex.RLock()
	jobs := make([]*shared.Job, 0, len(s.jobs))
	for _, job := range s.jobs {
		jobs = append(jobs, job)
	}
	s.jobsMutex.RUnlock()

	// Sort jobs by CreatedAt in descending order (newest first)
	sort.Slice(jobs, func(i, j int) bool {
		return jobs[i].CreatedAt.After(jobs[j].CreatedAt)
	})

	c.JSON(http.StatusOK, gin.H{"jobs": jobs})
}

func (s *APIServer) healthCheck(c *gin.Context) {
	status := gin.H{
		"status":    "healthy",
		"timestamp": time.Now(),
		"service":   "api-server",
	}

	if s.rabbitmq.IsConnectionClosed() {
		status["status"] = "unhealthy"
		status["rabbitmq"] = "disconnected"
		c.JSON(http.StatusServiceUnavailable, status)
		return
	}

	status["rabbitmq"] = "connected"

	// Add Ollama server endpoint information
	ollamaURL := getEnvOrDefault("OLLAMA_URL", "http://localhost:11434")
	status["ollama"] = gin.H{
		"endpoint": ollamaURL,
		"model":    getEnvOrDefault("OLLAMA_MODEL", "llama3.2"),
	}

	// Add MCP configuration information
	mcpTestMode := getEnvOrDefault("MCP_TEST_MODE", "false") == "true"
	mcpInfo := gin.H{
		"test_mode": mcpTestMode,
	}

	if !mcpTestMode {
		mcpInfo["endpoints"] = gin.H{
			"web_search": getEnvOrDefault("MCP_WEB_SERVER_URL", "http://localhost:3001"),
			"github":     getEnvOrDefault("MCP_GITHUB_SERVER_URL", "http://localhost:3002"),
			"files":      getEnvOrDefault("MCP_FILES_SERVER_URL", "http://localhost:3003"),
		}
	}

	status["mcp"] = mcpInfo

	c.JSON(http.StatusOK, status)
}

func (s *APIServer) setupRoutes() *gin.Engine {
	r := gin.Default()

	// Enable CORS for frontend
	r.Use(func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	})

	api := r.Group("/api")
	{
		api.POST("/jobs", s.createJob)
		api.GET("/jobs/:id", s.getJob)
		api.GET("/jobs", s.listJobs)
		api.GET("/health", s.healthCheck)
	}

	return r
}

// Utility function to get environment variable with default
func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func main() {
	server := NewAPIServer()

	if err := server.initRabbitMQ(); err != nil {
		log.Fatalf("Failed to initialize RabbitMQ: %v", err)
	}
	defer server.rabbitmq.Close()

	r := server.setupRoutes()

	log.Println("API Server starting on :8081")
	if err := r.Run(":8081"); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
