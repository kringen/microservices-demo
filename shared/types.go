package shared

import "time"

// JobStatus represents the current status of a job
type JobStatus string

const (
	JobStatusPending    JobStatus = "pending"
	JobStatusProcessing JobStatus = "processing"
	JobStatusCompleted  JobStatus = "completed"
	JobStatusFailed     JobStatus = "failed"
)

// Job represents a job in the system
type Job struct {
	ID          string     `json:"id"`
	Title       string     `json:"title"`
	Description string     `json:"description"`
	Status      JobStatus  `json:"status"`
	CreatedAt   time.Time  `json:"created_at"`
	StartedAt   *time.Time `json:"started_at,omitempty"`
	CompletedAt *time.Time `json:"completed_at,omitempty"`
	Result      string     `json:"result,omitempty"`
	Error       string     `json:"error,omitempty"`
}

// JobRequest represents a request to create a new job
type JobRequest struct {
	Title       string `json:"title" binding:"required"`
	Description string `json:"description"`
}

// JobMessage represents a message sent to the job queue
type JobMessage struct {
	JobID       string `json:"job_id"`
	Title       string `json:"title"`
	Description string `json:"description"`
}

// JobResult represents the result of a completed job
type JobResult struct {
	JobID       string    `json:"job_id"`
	Status      JobStatus `json:"status"`
	Result      string    `json:"result,omitempty"`
	Error       string    `json:"error,omitempty"`
	CompletedAt time.Time `json:"completed_at"`
}
