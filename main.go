package main

import (
	"context"
	"log"

	"github.com/proofchronicle/content-indexer/config"
	"github.com/proofchronicle/content-indexer/internal/consumer"
	"github.com/proofchronicle/content-indexer/internal/svc"
	amqp "github.com/rabbitmq/amqp091-go"
)

func main() {
	config := config.NewConfig()

	// Connect to RabbitMQ
	conn, err := amqp.Dial(config.RabbitMQURL)
	if err != nil {
		log.Fatalf("Failed to connect to RabbitMQ: %s", err)
	}
	defer conn.Close()

	consumer, err := consumer.NewConsumer(conn, config.QueueName, svc.NewProcessor(*config))
	if err != nil {
		log.Fatalf("Failed to create consumer: %s", err)
	}
	consumer.Start(context.TODO())
}
