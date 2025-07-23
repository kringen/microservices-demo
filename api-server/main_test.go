package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"microservices-demo/shared"

	"github.com/gin-gonic/gin"
)

func TestCreateJob(t *testing.T) {
	// Set Gin to test mode
	gin.SetMode(gin.TestMode)

	server := NewAPIServer()
	// Don't initialize RabbitMQ for unit tests
	router := server.setupRoutes()

	researchRequest := shared.ResearchRequest{
		Title:        "Test Research",
		Query:        "Research about AI and machine learning",
		ResearchType: shared.ResearchTypeGeneral,
		MCPServices:  []shared.MCPService{shared.MCPServiceWeb},
	}

	body, _ := json.Marshal(researchRequest)
	req, _ := http.NewRequest("POST", "/api/jobs", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Since RabbitMQ is not initialized, job creation should still succeed
	// but the job won't be queued (test mode)
	if w.Code != http.StatusCreated {
		t.Errorf("Expected status code %d, got %d", http.StatusCreated, w.Code)
	}

	// Verify the response contains the job
	var response shared.Job
	if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
		t.Errorf("Failed to parse response: %v", err)
	}

	if response.Title != researchRequest.Title {
		t.Errorf("Expected research title %s, got %s", researchRequest.Title, response.Title)
	}

	if response.Query != researchRequest.Query {
		t.Errorf("Expected research query %s, got %s", researchRequest.Query, response.Query)
	}

	if response.Status != "pending" {
		t.Errorf("Expected research status 'pending', got %s", response.Status)
	}
}

func TestGetJob(t *testing.T) {
	gin.SetMode(gin.TestMode)

	server := NewAPIServer()
	router := server.setupRoutes()

	// Create a test research job
	testJob := &shared.Job{
		ID:           "test-123",
		Title:        "Test Research",
		Query:        "Research about AI",
		ResearchType: shared.ResearchTypeGeneral,
		MCPServices:  []shared.MCPService{shared.MCPServiceWeb},
		Status:       shared.JobStatusPending,
		CreatedAt:    time.Now(),
	}

	server.jobsMutex.Lock()
	server.jobs[testJob.ID] = testJob
	server.jobsMutex.Unlock()

	req, _ := http.NewRequest("GET", "/api/jobs/test-123", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, w.Code)
	}

	var job shared.Job
	err := json.Unmarshal(w.Body.Bytes(), &job)
	if err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	if job.ID != testJob.ID {
		t.Errorf("Expected ID %s, got %s", testJob.ID, job.ID)
	}
}

func TestGetJobNotFound(t *testing.T) {
	gin.SetMode(gin.TestMode)

	server := NewAPIServer()
	router := server.setupRoutes()

	req, _ := http.NewRequest("GET", "/api/jobs/nonexistent", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("Expected status code %d, got %d", http.StatusNotFound, w.Code)
	}
}

func TestListJobs(t *testing.T) {
	gin.SetMode(gin.TestMode)

	server := NewAPIServer()
	router := server.setupRoutes()

	// Create test research jobs with different timestamps
	now := time.Now()
	testJobs := []*shared.Job{
		{
			ID:           "test-1",
			Title:        "Research 1",
			Query:        "First test research",
			ResearchType: shared.ResearchTypeGeneral,
			MCPServices:  []shared.MCPService{shared.MCPServiceWeb},
			Status:       shared.JobStatusPending,
			CreatedAt:    now.Add(-2 * time.Hour), // Created 2 hours ago
		},
		{
			ID:           "test-2",
			Title:        "Research 2",
			Query:        "Second test research",
			ResearchType: shared.ResearchTypeTechnical,
			MCPServices:  []shared.MCPService{shared.MCPServiceWeb, shared.MCPServiceGitHub},
			Status:       shared.JobStatusCompleted,
			CreatedAt:    now.Add(-1 * time.Hour), // Created 1 hour ago
		},
		{
			ID:           "test-3",
			Title:        "Research 3",
			Query:        "Third test research",
			ResearchType: shared.ResearchTypeGeneral,
			MCPServices:  []shared.MCPService{shared.MCPServiceWeb},
			Status:       shared.JobStatusProcessing,
			CreatedAt:    now, // Created now (most recent)
		},
	}

	server.jobsMutex.Lock()
	for _, job := range testJobs {
		server.jobs[job.ID] = job
	}
	server.jobsMutex.Unlock()

	req, _ := http.NewRequest("GET", "/api/jobs", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, w.Code)
	}

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	if err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	jobs, ok := response["jobs"].([]interface{})
	if !ok {
		t.Fatal("Expected jobs array in response")
	}

	if len(jobs) != 3 {
		t.Errorf("Expected 3 jobs, got %d", len(jobs))
	}

	// Verify jobs are sorted by CreatedAt in descending order (newest first)
	firstJob := jobs[0].(map[string]interface{})
	secondJob := jobs[1].(map[string]interface{})
	thirdJob := jobs[2].(map[string]interface{})

	if firstJob["id"].(string) != "test-3" {
		t.Errorf("Expected first job to be test-3 (newest), got %s", firstJob["id"].(string))
	}

	if secondJob["id"].(string) != "test-2" {
		t.Errorf("Expected second job to be test-2 (middle), got %s", secondJob["id"].(string))
	}

	if thirdJob["id"].(string) != "test-1" {
		t.Errorf("Expected third job to be test-1 (oldest), got %s", thirdJob["id"].(string))
	}
}

func TestUpdateJobStatus(t *testing.T) {
	server := NewAPIServer()

	// Create a test research job
	testJob := &shared.Job{
		ID:           "test-123",
		Title:        "Test Research",
		Query:        "Research about AI",
		ResearchType: shared.ResearchTypeGeneral,
		MCPServices:  []shared.MCPService{shared.MCPServiceWeb},
		Status:       shared.JobStatusPending,
		CreatedAt:    time.Now(),
	}

	server.jobsMutex.Lock()
	server.jobs[testJob.ID] = testJob
	server.jobsMutex.Unlock()

	// Update job status
	result := shared.JobResult{
		JobID:       "test-123",
		Status:      shared.JobStatusCompleted,
		Result:      "Job completed successfully",
		CompletedAt: time.Now(),
	}

	server.updateJobStatus(result)

	server.jobsMutex.RLock()
	updatedJob := server.jobs["test-123"]
	server.jobsMutex.RUnlock()

	if updatedJob.Status != shared.JobStatusCompleted {
		t.Errorf("Expected status %s, got %s", shared.JobStatusCompleted, updatedJob.Status)
	}
	if updatedJob.Result != "Job completed successfully" {
		t.Errorf("Expected result 'Job completed successfully', got %s", updatedJob.Result)
	}
	if updatedJob.CompletedAt == nil {
		t.Error("Expected CompletedAt to be set")
	}
}
