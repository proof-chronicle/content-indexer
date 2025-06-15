package consumer

import (
	"context"
	"encoding/json"
	"log"

	amqp "github.com/rabbitmq/amqp091-go"
)

type Message struct {
	Uid             string `json:"uid,omitempty"`
	CreatedAt       string `json:"created_at,omitempty"`
	Hash            string `json:"hash,omitempty"`
	Url             string `json:"url,omitempty"`
	ContentLength   uint64 `json:"content_length,omitempty"`
	ContentSelector string `json:"content_selector,omitempty"`
}

type MessageHandler interface {
	Handle(msg Message) error
}
type Consumer struct {
	conn    *amqp.Connection
	channel *amqp.Channel
	queue   string
	handler MessageHandler
}

func NewConsumer(conn *amqp.Connection, queueName string, handler MessageHandler) (*Consumer, error) {
	ch, err := conn.Channel()
	if err != nil {
		return nil, err
	}

	return &Consumer{
		conn:    conn,
		channel: ch,
		queue:   queueName,
		handler: handler,
	}, nil
}

func (c *Consumer) Start(ctx context.Context) error {
	msgs, err := c.channel.Consume(
		c.queue,
		"",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return err
	}

	go func() {
		for msg := range msgs {
			message := Message{}
			if err := json.Unmarshal(msg.Body, &message); err != nil {
				log.Printf("Error unmarshalling message: %s", err)
				continue
			}
			if err := c.handler.Handle(message); err != nil {
				log.Printf("Error handling message: %s", err)
			}
		}
	}()

	<-ctx.Done()
	return c.Close()
}

func (c *Consumer) Close() error {
	if err := c.channel.Close(); err != nil {
		return err
	}
	return c.conn.Close()
}
