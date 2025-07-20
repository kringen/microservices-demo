package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"os"
	"strings"
	"time"

	"microservices-demo/shared"

	amqp "github.com/rabbitmq/amqp091-go"
)

type JobRunner struct {
	rabbitmq *shared.RabbitMQClient
}

func NewJobRunner() *JobRunner {
	return &JobRunner{}
}

func (jr *JobRunner) initRabbitMQ() error {
	var err error

	// Get RabbitMQ URL from environment variable
	rabbitmqURL := os.Getenv("RABBITMQ_URL")
	if rabbitmqURL == "" {
		rabbitmqURL = "amqp://guest:guest@localhost:5672/"
	}

	jr.rabbitmq, err = shared.NewRabbitMQClient(rabbitmqURL)
	if err != nil {
		return err
	}
	return nil
}

func (jr *JobRunner) start() error {
	jobs, err := jr.rabbitmq.ConsumeJobs()
	if err != nil {
		return err
	}

	log.Println("Job Runner started. Waiting for jobs...")

	for delivery := range jobs {
		var jobMessage shared.JobMessage
		if err := json.Unmarshal(delivery.Body, &jobMessage); err != nil {
			log.Printf("Failed to unmarshal job message: %v", err)
			if err := delivery.Nack(false, false); err != nil {
				log.Printf("Failed to nack message: %v", err)
			}
			continue
		}

		log.Printf("Received job: %s - %s", jobMessage.JobID, jobMessage.Title)

		// Process the job in a goroutine
		go func(msg shared.JobMessage, d amqp.Delivery) {
			// Send processing status update
			processingUpdate := shared.JobResult{
				JobID:       msg.JobID,
				Status:      shared.JobStatusProcessing,
				CompletedAt: time.Now(), // Use this as "started at" timestamp
			}

			if err := jr.rabbitmq.PublishResult(processingUpdate); err != nil {
				log.Printf("Failed to publish processing status for job %s: %v", msg.JobID, err)
			} else {
				log.Printf("Job %s marked as processing", msg.JobID)
			}

			// Process the job and get final result
			result := jr.processJob(msg)

			// Publish the final result
			if err := jr.rabbitmq.PublishResult(result); err != nil {
				log.Printf("Failed to publish result for job %s: %v", msg.JobID, err)
			} else {
				log.Printf("Published final result for job %s", msg.JobID)
			}

			// Acknowledge the message
			if err := d.Ack(false); err != nil {
				log.Printf("Failed to ack message: %v", err)
			}
		}(jobMessage, delivery)
	}

	return nil
}

func (jr *JobRunner) processJob(jobMessage shared.JobMessage) shared.JobResult {
	log.Printf("Starting to process job: %s", jobMessage.JobID)

	// Simulate random processing time (5-60 seconds, max 1 minute)
	processingTime := time.Duration(rand.Intn(56)+5) * time.Second
	log.Printf("Job %s will take %v to complete", jobMessage.JobID, processingTime)

	// Process the job with timeout protection
	startTime := time.Now()

	// Use a timeout context to ensure jobs never exceed 1 minute
	timeout := time.After(60 * time.Second)
	done := make(chan bool)

	go func() {
		time.Sleep(processingTime)
		done <- true
	}()

	var result shared.JobResult

	select {
	case <-done:
		// Job completed within time limit
		actualDuration := time.Since(startTime)

		// Simulate random success/failure (90% success rate)
		success := rand.Float32() < 0.9

		result = shared.JobResult{
			JobID:       jobMessage.JobID,
			CompletedAt: time.Now(),
		}

		if success {
			result.Status = shared.JobStatusCompleted
			result.Result = jr.generateJobResult(jobMessage, actualDuration)
			log.Printf("Job %s completed successfully in %v", jobMessage.JobID, actualDuration)
		} else {
			result.Status = shared.JobStatusFailed
			result.Error = fmt.Sprintf("Job '%s' failed during processing: simulated random failure after %v",
				jobMessage.Title, actualDuration)
			log.Printf("Job %s failed after %v", jobMessage.JobID, actualDuration)
		}

	case <-timeout:
		// Job exceeded 1 minute timeout
		result = shared.JobResult{
			JobID:       jobMessage.JobID,
			Status:      shared.JobStatusFailed,
			Error:       fmt.Sprintf("Job '%s' timed out after 1 minute", jobMessage.Title),
			CompletedAt: time.Now(),
		}
		log.Printf("Job %s timed out after 1 minute", jobMessage.JobID)
	}

	return result
}

// generateJobResult creates a detailed result message based on job content
func (jr *JobRunner) generateJobResult(jobMessage shared.JobMessage, duration time.Duration) string {
	baseResult := fmt.Sprintf("Job '%s' completed successfully after %v", jobMessage.Title, duration)

	// Generate specific results based on job description/title
	var specificResult string
	description := strings.ToLower(jobMessage.Description)
	title := strings.ToLower(jobMessage.Title)

	switch {
	case contains(description, "calculation") || contains(title, "calculation"):
		specificResult = jr.simulateCalculation()
	case contains(description, "data") || contains(title, "data"):
		specificResult = jr.simulateDataProcessing()
	case contains(description, "report") || contains(title, "report"):
		specificResult = jr.simulateReportGeneration()
	case contains(description, "email") || contains(title, "email"):
		specificResult = jr.simulateEmailProcessing()
	case contains(description, "backup") || contains(title, "backup"):
		specificResult = jr.simulateBackupJob()
	case contains(description, "analysis") || contains(title, "analysis"):
		specificResult = jr.simulateAnalysisJob()
	default:
		specificResult = jr.simulateGenericWork()
	}

	return fmt.Sprintf("%s. %s", baseResult, specificResult)
}

func (jr *JobRunner) simulateCalculation() string {
	// Simulate some mathematical work
	result := 0
	for i := 0; i < 1000000; i++ {
		result += i
	}
	return fmt.Sprintf("Calculation completed. Result: %d", result)
}

func (jr *JobRunner) simulateDataProcessing() string {
	// Simulate data processing
	records := rand.Intn(1000) + 100
	return fmt.Sprintf("Processed %d data records successfully", records)
}

func (jr *JobRunner) simulateReportGeneration() string {
	// Simulate report generation
	pages := rand.Intn(50) + 10
	return fmt.Sprintf("Generated report with %d pages", pages)
}

func (jr *JobRunner) simulateGenericWork() string {
	// Generic work simulation
	operations := rand.Intn(100) + 20
	return fmt.Sprintf("Performed %d operations successfully", operations)
}

func (jr *JobRunner) simulateEmailProcessing() string {
	emailsSent := rand.Intn(500) + 50
	bounces := rand.Intn(emailsSent / 20)
	return fmt.Sprintf("Sent %d emails with %d bounces", emailsSent, bounces)
}

func (jr *JobRunner) simulateBackupJob() string {
	sizeMB := rand.Intn(5000) + 100
	files := rand.Intn(10000) + 1000
	return fmt.Sprintf("Backed up %d files (%.1f GB)", files, float64(sizeMB)/1024.0)
}

func (jr *JobRunner) simulateAnalysisJob() string {
	dataPoints := rand.Intn(1000000) + 10000
	insights := rand.Intn(50) + 5
	return fmt.Sprintf("Analyzed %d data points and generated %d insights", dataPoints, insights)
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) && s[:len(substr)] == substr
}

func main() {
	runner := NewJobRunner()

	if err := runner.initRabbitMQ(); err != nil {
		log.Fatalf("Failed to initialize RabbitMQ: %v", err)
	}
	defer runner.rabbitmq.Close()

	log.Println("Job Runner is starting...")

	if err := runner.start(); err != nil {
		log.Fatalf("Failed to start job runner: %v", err)
	}
}
