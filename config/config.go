package config

import "os"

type Config struct {
	RabbitMQURL string
	QueueName   string
	DBDSN       string
	GatewayAddr string
}

func NewConfig() *Config {
	return &Config{
		RabbitMQURL: os.Getenv("RABBITMQ_URL"),
		QueueName:   os.Getenv("QUEUE_NAME"),
		DBDSN:       os.Getenv("DB_DSN"),
		GatewayAddr: "chain-gateway:50051",
	}
}
