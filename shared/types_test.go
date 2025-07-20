package shared

import (
	"encoding/json"
	"testing"
)

func TestJobMessageSerialization(t *testing.T) {
	job := JobMessage{
		JobID:       "test-123",
		Title:       "Test Job",
		Description: "A test job description",
	}

	// Test marshaling
	data, err := json.Marshal(job)
	if err != nil {
		t.Fatalf("Failed to marshal job message: %v", err)
	}

	// Test unmarshaling
	var unmarshaled JobMessage
	err = json.Unmarshal(data, &unmarshaled)
	if err != nil {
		t.Fatalf("Failed to unmarshal job message: %v", err)
	}

	if unmarshaled.JobID != job.JobID {
		t.Errorf("Expected JobID %s, got %s", job.JobID, unmarshaled.JobID)
	}
	if unmarshaled.Title != job.Title {
		t.Errorf("Expected Title %s, got %s", job.Title, unmarshaled.Title)
	}
	if unmarshaled.Description != job.Description {
		t.Errorf("Expected Description %s, got %s", job.Description, unmarshaled.Description)
	}
}

func TestJobResultSerialization(t *testing.T) {
	result := JobResult{
		JobID:  "test-123",
		Status: JobStatusCompleted,
		Result: "Job completed successfully",
	}

	// Test marshaling
	data, err := json.Marshal(result)
	if err != nil {
		t.Fatalf("Failed to marshal job result: %v", err)
	}

	// Test unmarshaling
	var unmarshaled JobResult
	err = json.Unmarshal(data, &unmarshaled)
	if err != nil {
		t.Fatalf("Failed to unmarshal job result: %v", err)
	}

	if unmarshaled.JobID != result.JobID {
		t.Errorf("Expected JobID %s, got %s", result.JobID, unmarshaled.JobID)
	}
	if unmarshaled.Status != result.Status {
		t.Errorf("Expected Status %s, got %s", result.Status, unmarshaled.Status)
	}
	if unmarshaled.Result != result.Result {
		t.Errorf("Expected Result %s, got %s", result.Result, unmarshaled.Result)
	}
}

func TestJobStatusValues(t *testing.T) {
	statuses := []JobStatus{
		JobStatusPending,
		JobStatusProcessing,
		JobStatusCompleted,
		JobStatusFailed,
	}

	expectedValues := []string{
		"pending",
		"processing",
		"completed",
		"failed",
	}

	for i, status := range statuses {
		if string(status) != expectedValues[i] {
			t.Errorf("Expected status %s, got %s", expectedValues[i], string(status))
		}
	}
}
