package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	"microservices-demo/shared"

	"github.com/gin-gonic/gin"
)

var apiServerURL string

func init() {
	apiServerURL = os.Getenv("API_SERVER_URL")
	if apiServerURL == "" {
		apiServerURL = "http://localhost:8081"
	}
}

type Frontend struct {
	templates *template.Template
}

func NewFrontend() *Frontend {
	return &Frontend{}
}

func (f *Frontend) loadTemplates() error {
	var err error
	f.templates, err = template.New("").Funcs(template.FuncMap{
		"formatTime": func(t time.Time) string {
			return t.Format("2006-01-02 15:04:05")
		},
		"formatDuration": func(start, end *time.Time) string {
			if start == nil || end == nil {
				return "N/A"
			}
			duration := end.Sub(*start)
			return duration.String()
		},
		"statusColor": func(status shared.JobStatus) string {
			switch status {
			case shared.JobStatusPending:
				return "warning"
			case shared.JobStatusProcessing:
				return "info"
			case shared.JobStatusCompleted:
				return "success"
			case shared.JobStatusFailed:
				return "danger"
			default:
				return "secondary"
			}
		},
	}).ParseGlob("templates/*.html")

	if err != nil {
		// If templates directory doesn't exist, create inline templates
		f.createInlineTemplates()
	}

	return nil
}

func (f *Frontend) createInlineTemplates() {
	f.templates = template.Must(template.New("").Funcs(template.FuncMap{
		"formatTime": func(t time.Time) string {
			return t.Format("2006-01-02 15:04:05")
		},
		"formatDuration": func(start, end *time.Time) string {
			if start == nil || end == nil {
				return "N/A"
			}
			duration := end.Sub(*start)
			return duration.String()
		},
		"statusColor": func(status shared.JobStatus) string {
			switch status {
			case shared.JobStatusPending:
				return "warning"
			case shared.JobStatusProcessing:
				return "info"
			case shared.JobStatusCompleted:
				return "success"
			case shared.JobStatusFailed:
				return "danger"
			default:
				return "secondary"
			}
		},
		"multiply": func(a, b float64) float64 {
			return a * b
		},
	}).Parse(indexTemplate + jobStatusTemplate))
}

func (f *Frontend) homePage(c *gin.Context) {
	// Get recent jobs
	jobs, err := f.fetchJobs()
	if err != nil {
		log.Printf("Failed to fetch jobs: %v", err)
		jobs = []shared.Job{} // Empty slice on error
	}

	data := gin.H{
		"Title": "Microservices Demo",
		"Jobs":  jobs,
	}

	c.Header("Content-Type", "text/html")
	if err := f.templates.ExecuteTemplate(c.Writer, "index", data); err != nil {
		log.Printf("Template execution error: %v", err)
		c.String(http.StatusInternalServerError, "Template error")
	}
}

func (f *Frontend) submitJob(c *gin.Context) {
	title := c.PostForm("title")
	description := c.PostForm("description")

	if title == "" {
		c.Redirect(http.StatusSeeOther, "/?error=Title is required")
		return
	}

	jobRequest := shared.JobRequest{
		Title:       title,
		Description: description,
	}

	job, err := f.createJob(jobRequest)
	if err != nil {
		log.Printf("Failed to create job: %v", err)
		c.Redirect(http.StatusSeeOther, "/?error=Failed to create job")
		return
	}

	// For AJAX requests, return JSON
	if c.GetHeader("Accept") == "application/json" || c.GetHeader("Content-Type") == "application/json" {
		c.JSON(http.StatusCreated, gin.H{
			"success": true,
			"job":     job,
			"message": "Job created successfully",
		})
		return
	}

	// For form submissions, redirect to job status page
	c.Redirect(http.StatusSeeOther, fmt.Sprintf("/status/%s?created=true", job.ID))
}

func (f *Frontend) jobStatus(c *gin.Context) {
	jobID := c.Param("id")

	job, err := f.fetchJob(jobID)
	if err != nil {
		log.Printf("Failed to fetch job: %v", err)
		c.String(http.StatusNotFound, "Job not found")
		return
	}

	data := gin.H{
		"Title": fmt.Sprintf("Job Status - %s", job.Title),
		"Job":   job,
	}

	c.Header("Content-Type", "text/html")
	if err := f.templates.ExecuteTemplate(c.Writer, "job-status", data); err != nil {
		log.Printf("Template execution error: %v", err)
		c.String(http.StatusInternalServerError, "Template error")
	}
}

func (f *Frontend) apiStatus(c *gin.Context) {
	resp, err := http.Get(apiServerURL + "/api/health")
	if err != nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{
			"status": "unhealthy",
			"error":  err.Error(),
		})
		return
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status": "error",
			"error":  err.Error(),
		})
		return
	}

	var healthStatus map[string]interface{}
	if err := json.Unmarshal(body, &healthStatus); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status": "error",
			"error":  err.Error(),
		})
		return
	}

	c.JSON(resp.StatusCode, healthStatus)
}

func (f *Frontend) createJob(jobRequest shared.JobRequest) (*shared.Job, error) {
	body, err := json.Marshal(jobRequest)
	if err != nil {
		return nil, err
	}

	resp, err := http.Post(apiServerURL+"/api/jobs", "application/json", bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		return nil, fmt.Errorf("API returned status %d", resp.StatusCode)
	}

	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var job shared.Job
	if err := json.Unmarshal(responseBody, &job); err != nil {
		return nil, err
	}

	return &job, nil
}

func (f *Frontend) fetchJob(jobID string) (*shared.Job, error) {
	resp, err := http.Get(fmt.Sprintf("%s/api/jobs/%s", apiServerURL, jobID))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API returned status %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var job shared.Job
	if err := json.Unmarshal(body, &job); err != nil {
		return nil, err
	}

	return &job, nil
}

func (f *Frontend) fetchJobs() ([]shared.Job, error) {
	resp, err := http.Get(apiServerURL + "/api/jobs")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API returned status %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var response struct {
		Jobs []shared.Job `json:"jobs"`
	}
	if err := json.Unmarshal(body, &response); err != nil {
		return nil, err
	}

	return response.Jobs, nil
}

func (f *Frontend) apiJobs(c *gin.Context) {
	jobs, err := f.fetchJobs()
	if err != nil {
		log.Printf("Failed to fetch jobs: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to fetch jobs",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"jobs": jobs,
	})
}

func (f *Frontend) apiGetJob(c *gin.Context) {
	jobID := c.Param("id")

	job, err := f.fetchJob(jobID)
	if err != nil {
		log.Printf("Failed to fetch job %s: %v", jobID, err)
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Job not found",
		})
		return
	}

	c.JSON(http.StatusOK, job)
}

func (f *Frontend) submitJobAPI(c *gin.Context) {
	var jobRequest shared.JobRequest
	if err := c.ShouldBindJSON(&jobRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request format",
		})
		return
	}

	if jobRequest.Title == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Title is required",
		})
		return
	}

	job, err := f.createJob(jobRequest)
	if err != nil {
		log.Printf("Failed to create job: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to create job",
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"success": true,
		"job":     job,
		"message": "Job created successfully",
	})
}

func (f *Frontend) setupRoutes() *gin.Engine {
	r := gin.Default()

	r.GET("/", f.homePage)
	r.POST("/submit", f.submitJob)
	r.GET("/status/:id", f.jobStatus)
	r.GET("/api/status", f.apiStatus)
	r.GET("/api/jobs", f.apiJobs)
	r.GET("/api/jobs/:id", f.apiGetJob)
	r.POST("/api/jobs", f.submitJobAPI)

	// Static files (if needed)
	r.Static("/static", "./static")

	return r
}

func main() {
	frontend := NewFrontend()

	if err := frontend.loadTemplates(); err != nil {
		log.Printf("Warning: Failed to load templates: %v", err)
	}

	r := frontend.setupRoutes()

	log.Println("Frontend starting on :8080")
	if err := r.Run(":8080"); err != nil {
		log.Fatalf("Failed to start frontend: %v", err)
	}
}
