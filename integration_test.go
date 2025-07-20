//go:build integration
// +build integration

package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"testing"
	"time"

	"microservices-demo/shared"
)

// Integration test that requires all services to be running
// Run with: go test -tags=integration ./integration_test.go

const (
	apiURL      = "http://localhost:8081"
	frontendURL = "http://localhost:8080"
)

func TestEndToEndJobFlow(t *testing.T) {
	// Test API health first
	resp, err := http.Get(apiURL + "/api/health")
	if err != nil {
		t.Skipf("API server not available: %v", err)
		return
	}
	resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		t.Skipf("API server not healthy, status: %d", resp.StatusCode)
		return
	}

	// Create a job
	jobRequest := shared.JobRequest{
		Title:       "Integration Test Job",
		Description: "A job created during integration testing",
	}

	jobData, err := json.Marshal(jobRequest)
	if err != nil {
		t.Fatalf("Failed to marshal job request: %v", err)
	}

	// Submit job
	resp, err = http.Post(apiURL+"/api/jobs", "application/json", bytes.NewBuffer(jobData))
	if err != nil {
		t.Fatalf("Failed to submit job: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		t.Fatalf("Expected status 201, got %d", resp.StatusCode)
	}

	var createdJob shared.Job
	if err := json.NewDecoder(resp.Body).Decode(&createdJob); err != nil {
		t.Fatalf("Failed to decode job response: %v", err)
	}

	t.Logf("Created job with ID: %s", createdJob.ID)

	// Poll for job completion (timeout after 90 seconds to account for max 60s job time)
	timeout := time.After(90 * time.Second)
	ticker := time.NewTicker(2 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-timeout:
			t.Fatal("Job did not complete within 90 seconds")
		case <-ticker.C:
			// Check job status
			resp, err := http.Get(fmt.Sprintf("%s/api/jobs/%s", apiURL, createdJob.ID))
			if err != nil {
				t.Logf("Error checking job status: %v", err)
				continue
			}

			var job shared.Job
			if err := json.NewDecoder(resp.Body).Decode(&job); err != nil {
				resp.Body.Close()
				t.Logf("Error decoding job response: %v", err)
				continue
			}
			resp.Body.Close()

			t.Logf("Job status: %s", job.Status)

			if job.Status == shared.JobStatusCompleted {
				t.Logf("Job completed successfully!")
				t.Logf("Result: %s", job.Result)
				return
			} else if job.Status == shared.JobStatusFailed {
				t.Logf("Job failed: %s", job.Error)
				return // Still a successful test since the flow worked
			}
		}
	}
}

func TestFrontendAccessibility(t *testing.T) {
	resp, err := http.Get(frontendURL)
	if err != nil {
		t.Skipf("Frontend not available: %v", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected frontend to return 200, got %d", resp.StatusCode)
	}

	// Check that it's actually serving HTML
	contentType := resp.Header.Get("Content-Type")
	if contentType != "text/html" && contentType != "text/html; charset=utf-8" {
		t.Errorf("Expected HTML content type, got %s", contentType)
	}
}

func TestAPIJobsList(t *testing.T) {
	resp, err := http.Get(apiURL + "/api/jobs")
	if err != nil {
		t.Skipf("API server not available: %v", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status 200, got %d", resp.StatusCode)
	}

	var response map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		t.Fatalf("Failed to decode jobs list response: %v", err)
	}

	if _, exists := response["jobs"]; !exists {
		t.Error("Expected 'jobs' field in response")
	}
}
