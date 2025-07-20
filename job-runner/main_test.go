package main

import (
	"strings"
	"testing"
	"time"

	"microservices-demo/shared"
)

func TestProcessJob(t *testing.T) {
	runner := NewJobRunner()

	jobMessage := shared.JobMessage{
		JobID:       "test-123",
		Title:       "Test Job",
		Description: "A test job for unit testing",
	}

	result := runner.processJob(jobMessage)

	if result.JobID != jobMessage.JobID {
		t.Errorf("Expected JobID %s, got %s", jobMessage.JobID, result.JobID)
	}

	if result.Status != shared.JobStatusCompleted && result.Status != shared.JobStatusFailed {
		t.Errorf("Expected status to be completed or failed, got %s", result.Status)
	}

	if result.CompletedAt.IsZero() {
		t.Error("Expected CompletedAt to be set")
	}

	// Check that either result or error is set
	if result.Status == shared.JobStatusCompleted && result.Result == "" {
		t.Error("Expected result to be set for completed job")
	}
	if result.Status == shared.JobStatusFailed && result.Error == "" {
		t.Error("Expected error to be set for failed job")
	}
}

func TestSimulateCalculation(t *testing.T) {
	runner := NewJobRunner()
	result := runner.simulateCalculation()

	if result == "" {
		t.Error("Expected calculation result to be non-empty")
	}

	// Should contain "Calculation completed"
	if len(result) < 10 {
		t.Error("Expected longer calculation result")
	}
}

func TestSimulateDataProcessing(t *testing.T) {
	runner := NewJobRunner()
	result := runner.simulateDataProcessing()

	if result == "" {
		t.Error("Expected data processing result to be non-empty")
	}

	// Should contain "Processed"
	if len(result) < 10 {
		t.Error("Expected longer data processing result")
	}
}

func TestSimulateReportGeneration(t *testing.T) {
	runner := NewJobRunner()
	result := runner.simulateReportGeneration()

	if result == "" {
		t.Error("Expected report generation result to be non-empty")
	}

	// Should contain "Generated report"
	if len(result) < 10 {
		t.Error("Expected longer report generation result")
	}
}

func TestSimulateGenericWork(t *testing.T) {
	runner := NewJobRunner()
	result := runner.simulateGenericWork()

	if result == "" {
		t.Error("Expected generic work result to be non-empty")
	}

	// Should contain "Performed"
	if len(result) < 10 {
		t.Error("Expected longer generic work result")
	}
}

func TestContains(t *testing.T) {
	testCases := []struct {
		s        string
		substr   string
		expected bool
	}{
		{"calculation task", "calculation", true},
		{"data processing", "data", true},
		{"report generation", "report", true},
		{"simple task", "calculation", false},
		{"", "test", false},
		{"test", "", true},
	}

	for _, tc := range testCases {
		result := contains(tc.s, tc.substr)
		if result != tc.expected {
			t.Errorf("contains(%q, %q) = %v, expected %v", tc.s, tc.substr, result, tc.expected)
		}
	}
}

func BenchmarkProcessJob(b *testing.B) {
	runner := NewJobRunner()
	jobMessage := shared.JobMessage{
		JobID:       "benchmark-job",
		Title:       "Benchmark Job",
		Description: "A job for benchmarking",
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		runner.processJob(jobMessage)
	}
}

func TestGenerateJobResult(t *testing.T) {
	runner := NewJobRunner()

	testCases := []struct {
		title       string
		description string
		expectType  string
	}{
		{"Data Analysis", "analyze customer data", "analysis"},
		{"Email Campaign", "send marketing emails", "email"},
		{"Backup Task", "backup database files", "backup"},
		{"Calculate Report", "financial calculations", "calculation"},
		{"Process Data", "process user data", "data"},
		{"Generate Report", "monthly report generation", "report"},
		{"Generic Task", "some other work", "generic"},
	}

	for _, tc := range testCases {
		jobMessage := shared.JobMessage{
			JobID:       "test-job",
			Title:       tc.title,
			Description: tc.description,
		}

		result := runner.generateJobResult(jobMessage, 5*time.Second)

		if result == "" {
			t.Errorf("Expected non-empty result for job type %s", tc.expectType)
		}

		if !strings.Contains(result, "completed successfully") {
			t.Errorf("Expected result to contain success message, got: %s", result)
		}

		// Verify it contains job title
		if !strings.Contains(result, tc.title) {
			t.Errorf("Expected result to contain job title %s, got: %s", tc.title, result)
		}
	}
}

func TestJobTimeout(t *testing.T) {
	runner := NewJobRunner()
	
	// Create a job that would normally take longer than 1 minute
	// For testing, we'll verify the timeout structure exists
	jobMessage := shared.JobMessage{
		JobID:       "timeout-test",
		Title:       "Long Running Job",
		Description: "A job that should timeout",
	}
	
	// We can't easily test the full timeout in unit tests (takes too long)
	// but we can verify the job completes normally within our test timeframe
	start := time.Now()
	result := runner.processJob(jobMessage)
	duration := time.Since(start)
	
	// Verify the job completed
	if result.JobID != jobMessage.JobID {
		t.Errorf("Expected JobID %s, got %s", jobMessage.JobID, result.JobID)
	}
	
	// Verify it didn't take longer than the 60 second timeout
	if duration > 61*time.Second {
		t.Errorf("Job took longer than expected timeout: %v", duration)
	}
	
	// Log the actual time for reference
	t.Logf("Job completed in %v", duration)
}

func TestNewJobTypes(t *testing.T) {
	runner := NewJobRunner()

	// Test email processing
	emailResult := runner.simulateEmailProcessing()
	if !strings.Contains(emailResult, "emails") {
		t.Errorf("Expected email result to mention emails, got: %s", emailResult)
	}

	// Test backup job
	backupResult := runner.simulateBackupJob()
	if !strings.Contains(backupResult, "Backed up") {
		t.Errorf("Expected backup result to mention backup, got: %s", backupResult)
	}

	// Test analysis job
	analysisResult := runner.simulateAnalysisJob()
	if !strings.Contains(analysisResult, "Analyzed") {
		t.Errorf("Expected analysis result to mention analysis, got: %s", analysisResult)
	}
}
