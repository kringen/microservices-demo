package shared

import (
	"encoding/json"
	"log"

	amqp "github.com/rabbitmq/amqp091-go"
)

const (
	JobQueueName    = "jobs"
	ResultQueueName = "job_results"
)

// RabbitMQClient wraps the RabbitMQ connection and channel
type RabbitMQClient struct {
	connection *amqp.Connection
	channel    *amqp.Channel
}

// NewRabbitMQClient creates a new RabbitMQ client
func NewRabbitMQClient(url string) (*RabbitMQClient, error) {
	conn, err := amqp.Dial(url)
	if err != nil {
		return nil, err
	}

	ch, err := conn.Channel()
	if err != nil {
		conn.Close()
		return nil, err
	}

	client := &RabbitMQClient{
		connection: conn,
		channel:    ch,
	}

	// Declare queues
	if err := client.declareQueues(); err != nil {
		client.Close()
		return nil, err
	}

	return client, nil
}

// declareQueues declares the required queues
func (c *RabbitMQClient) declareQueues() error {
	queues := []string{JobQueueName, ResultQueueName}

	for _, queueName := range queues {
		_, err := c.channel.QueueDeclare(
			queueName, // name
			true,      // durable
			false,     // delete when unused
			false,     // exclusive
			false,     // no-wait
			nil,       // arguments
		)
		if err != nil {
			return err
		}
	}

	return nil
}

// PublishJob publishes a job message to the job queue
func (c *RabbitMQClient) PublishJob(job JobMessage) error {
	body, err := json.Marshal(job)
	if err != nil {
		return err
	}

	return c.channel.Publish(
		"",           // exchange
		JobQueueName, // routing key
		false,        // mandatory
		false,        // immediate
		amqp.Publishing{
			ContentType: "application/json",
			Body:        body,
		})
}

// PublishResult publishes a job result to the result queue
func (c *RabbitMQClient) PublishResult(result JobResult) error {
	body, err := json.Marshal(result)
	if err != nil {
		return err
	}

	return c.channel.Publish(
		"",              // exchange
		ResultQueueName, // routing key
		false,           // mandatory
		false,           // immediate
		amqp.Publishing{
			ContentType: "application/json",
			Body:        body,
		})
}

// ConsumeJobs consumes job messages from the job queue
func (c *RabbitMQClient) ConsumeJobs() (<-chan amqp.Delivery, error) {
	return c.channel.Consume(
		JobQueueName, // queue
		"",           // consumer
		false,        // auto-ack
		false,        // exclusive
		false,        // no-local
		false,        // no-wait
		nil,          // args
	)
}

// ConsumeResults consumes job result messages from the result queue
func (c *RabbitMQClient) ConsumeResults() (<-chan amqp.Delivery, error) {
	return c.channel.Consume(
		ResultQueueName, // queue
		"",              // consumer
		true,            // auto-ack
		false,           // exclusive
		false,           // no-local
		false,           // no-wait
		nil,             // args
	)
}

// Close closes the RabbitMQ connection and channel
func (c *RabbitMQClient) Close() {
	if c.channel != nil {
		c.channel.Close()
	}
	if c.connection != nil {
		c.connection.Close()
	}
}

// IsConnectionClosed checks if the connection is closed
func (c *RabbitMQClient) IsConnectionClosed() bool {
	return c.connection.IsClosed()
}

// LogError is a helper function to log errors
func LogError(err error, message string) {
	if err != nil {
		log.Printf("ERROR: %s: %v", message, err)
	}
}
