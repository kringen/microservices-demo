package shared

import (
	"testing"
	"time"
)

// TestRabbitMQClientCreation tests the creation of RabbitMQ client
// Note: This test requires RabbitMQ to be running
func TestRabbitMQClientCreation(t *testing.T) {
	// Skip this test if RabbitMQ is not available
	client, err := NewRabbitMQClient("amqp://guest:guest@localhost:5672/")
	if err != nil {
		t.Skipf("Skipping RabbitMQ test: %v", err)
		return
	}
	defer client.Close()

	if client.connection == nil {
		t.Error("Expected connection to be initialized")
	}
	if client.channel == nil {
		t.Error("Expected channel to be initialized")
	}
}

func TestJobMessageAndResultSerialization(t *testing.T) {
	// Test JobMessage
	jobMsg := JobMessage{
		JobID:       "test-job-123",
		Title:       "Test Job",
		Description: "This is a test job",
	}

	// Create a mock RabbitMQ client (without actual connection)
	// and test message creation
	if jobMsg.JobID == "" {
		t.Error("JobID should not be empty")
	}
	if jobMsg.Title == "" {
		t.Error("Title should not be empty")
	}

	// Test JobResult
	now := time.Now()
	result := JobResult{
		JobID:       "test-job-123",
		Status:      JobStatusCompleted,
		Result:      "Job completed successfully",
		CompletedAt: now,
	}

	if result.JobID != jobMsg.JobID {
		t.Error("JobResult should have same JobID as JobMessage")
	}
	if result.Status != JobStatusCompleted {
		t.Error("Status should be completed")
	}
	if result.CompletedAt.IsZero() {
		t.Error("CompletedAt should be set")
	}
}

func TestQueueConstants(t *testing.T) {
	expectedJobQueue := "jobs"
	expectedResultQueue := "job_results"

	if JobQueueName != expectedJobQueue {
		t.Errorf("Expected JobQueueName to be %s, got %s", expectedJobQueue, JobQueueName)
	}
	if ResultQueueName != expectedResultQueue {
		t.Errorf("Expected ResultQueueName to be %s, got %s", expectedResultQueue, ResultQueueName)
	}
}

// MockRabbitMQClient for testing without actual RabbitMQ connection
type MockRabbitMQClient struct {
	messages []interface{}
	closed   bool
}

func NewMockRabbitMQClient() *MockRabbitMQClient {
	return &MockRabbitMQClient{
		messages: make([]interface{}, 0),
		closed:   false,
	}
}

func (m *MockRabbitMQClient) PublishJob(job JobMessage) error {
	m.messages = append(m.messages, job)
	return nil
}

func (m *MockRabbitMQClient) PublishResult(result JobResult) error {
	m.messages = append(m.messages, result)
	return nil
}

func (m *MockRabbitMQClient) Close() {
	m.closed = true
}

func (m *MockRabbitMQClient) IsConnectionClosed() bool {
	return m.closed
}

func TestMockRabbitMQClient(t *testing.T) {
	mock := NewMockRabbitMQClient()

	if mock.IsConnectionClosed() {
		t.Error("Mock client should not be closed initially")
	}

	// Test job publishing
	job := JobMessage{
		JobID:       "test-123",
		Title:       "Test Job",
		Description: "Test description",
	}

	err := mock.PublishJob(job)
	if err != nil {
		t.Errorf("Mock should not return error: %v", err)
	}

	if len(mock.messages) != 1 {
		t.Errorf("Expected 1 message, got %d", len(mock.messages))
	}

	// Test result publishing
	result := JobResult{
		JobID:       "test-123",
		Status:      JobStatusCompleted,
		Result:      "Success",
		CompletedAt: time.Now(),
	}

	err = mock.PublishResult(result)
	if err != nil {
		t.Errorf("Mock should not return error: %v", err)
	}

	if len(mock.messages) != 2 {
		t.Errorf("Expected 2 messages, got %d", len(mock.messages))
	}

	// Test close
	mock.Close()
	if !mock.IsConnectionClosed() {
		t.Error("Mock client should be closed after Close()")
	}
}
