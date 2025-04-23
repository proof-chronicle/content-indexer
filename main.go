package main

import (
	"log"
	"os"

	amqp "github.com/rabbitmq/amqp091-go"
)

// failOnError logs and exits on error
func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}

func main() {
	// RabbitMQ connection settings
	rabbitURL := os.Getenv("RABBITMQ_URL")
	if rabbitURL == "" {
		failOnError(nil, "RABBITMQ_URL environment variable not set")
	}
	queueName := os.Getenv("QUEUE_NAME")
	if queueName == "" {
		failOnError(nil, "QUEUE_NAME environment variable not set")
	}

	// Connect to RabbitMQ
	conn, err := amqp.Dial(rabbitURL)
	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	// Open a channel
	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()

	// Consume messages
	msgs, err := ch.Consume(
		queueName, // queue
		"",        // consumer
		true,      // auto-ack
		false,     // exclusive
		false,     // no-local
		false,     // no-wait
		nil,       // args
	)
	failOnError(err, "Failed to register a consumer")

	log.Printf("[*] Waiting for messages on queue '%s' (URL: %s)...", queueName, rabbitURL)

	// Process messages
	for d := range msgs {
		log.Printf("[x] Received new message: %s", string(d.Body))
	}
}
