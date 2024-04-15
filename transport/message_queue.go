package transport

import (
	"fmt"
	"log"

	"github.com/chaewonkong/msa-link-scraper/config"
	amqp "github.com/rabbitmq/amqp091-go"
)

// RabbitMQ represents the message queue
type RabbitMQ struct {
	// Conn is the connection
	Conn *amqp.Connection

	// Ch is the channel
	Ch *amqp.Channel
}

// NewRabbitMQ creates a new RabbitMQ
func NewRabbitMQ(cfg *config.App) *RabbitMQ {
	conStr := fmt.Sprintf("amqp://%s:%s@%s:%s/", cfg.Queue.User, cfg.Queue.Password, cfg.Queue.Host, cfg.Queue.Port)
	conn, err := amqp.Dial(conStr)
	if err != nil {
		log.Fatal(err)
	}

	ch, err := conn.Channel()
	if err != nil {
		log.Fatal(err)
	}

	return &RabbitMQ{conn, ch}
}

// Close closes the connection and channel
func (mq *RabbitMQ) Close() {
	mq.Conn.Close()
	mq.Ch.Close()
}
