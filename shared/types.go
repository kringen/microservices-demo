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

// ResearchType represents different types of research requests
type ResearchType string

const (
	ResearchTypeGeneral     ResearchType = "general"
	ResearchTypeCode        ResearchType = "code"
	ResearchTypeData        ResearchType = "data"
	ResearchTypeMarket      ResearchType = "market"
	ResearchTypeTechnical   ResearchType = "technical"
	ResearchTypeCompetitive ResearchType = "competitive"
)

// MCPService represents available MCP services for research
type MCPService string

const (
	MCPServiceWeb      MCPService = "web"
	MCPServiceGitHub   MCPService = "github"
	MCPServiceDatabase MCPService = "database"
	MCPServiceFiles    MCPService = "files"
	MCPServiceCalendar MCPService = "calendar"
	MCPServiceSlack    MCPService = "slack"
)

// Job represents a research job in the system
type Job struct {
	ID           string       `json:"id"`
	Title        string       `json:"title"`
	Query        string       `json:"query"`
	ResearchType ResearchType `json:"research_type"`
	MCPServices  []MCPService `json:"mcp_services"`
	Status       JobStatus    `json:"status"`
	CreatedAt    time.Time    `json:"created_at"`
	StartedAt    *time.Time   `json:"started_at,omitempty"`
	CompletedAt  *time.Time   `json:"completed_at,omitempty"`
	Result       string       `json:"result,omitempty"`
	Sources      []string     `json:"sources,omitempty"`
	Error        string       `json:"error,omitempty"`
	Confidence   float64      `json:"confidence,omitempty"`
	TokensUsed   int          `json:"tokens_used,omitempty"`
}

// ResearchRequest represents a request to create a new research job
type ResearchRequest struct {
	Title        string       `json:"title" binding:"required"`
	Query        string       `json:"query" binding:"required"`
	ResearchType ResearchType `json:"research_type"`
	MCPServices  []MCPService `json:"mcp_services"`
}

// JobMessage represents a message sent to the research queue
type JobMessage struct {
	JobID        string       `json:"job_id"`
	Title        string       `json:"title"`
	Query        string       `json:"query"`
	ResearchType ResearchType `json:"research_type"`
	MCPServices  []MCPService `json:"mcp_services"`
}

// JobResult represents the result of a completed research job
type JobResult struct {
	JobID       string    `json:"job_id"`
	Status      JobStatus `json:"status"`
	Result      string    `json:"result,omitempty"`
	Sources     []string  `json:"sources,omitempty"`
	Error       string    `json:"error,omitempty"`
	CompletedAt time.Time `json:"completed_at"`
	Confidence  float64   `json:"confidence,omitempty"`
	TokensUsed  int       `json:"tokens_used,omitempty"`
}
