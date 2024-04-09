package main

import (
	"fmt"
	"log"

	config "github.com/chaewonkong/msa-link-api/config"
	amqp "github.com/rabbitmq/amqp091-go"
)

func main() {
	cfg := config.NewAppConfig()

	conStr := fmt.Sprintf("amqp://%s:%s@%s:%s/", cfg.Queue.User, cfg.Queue.Password, cfg.Queue.Host, cfg.Queue.Port)
	queueConn, err := amqp.Dial(conStr)
	if err != nil {
		log.Fatal(err)
	}
	defer queueConn.Close()

	ch, err := queueConn.Channel()
	if err != nil {
		log.Fatal(err)
	}
	defer ch.Close()

	q, err := ch.QueueDeclare(
		"link", // name
		false,  // durable
		false,  // delete when unused
		false,  // exclusive
		false,  // no-wait
		nil,    // arguments
	)
	if err != nil {
		log.Fatal(err)
	}

	msgs, err := ch.Consume(
		q.Name, // queue
		"",     // consumer
		true,   // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
	if err != nil {
		log.Fatal(err)
	}

	var forever chan struct{}

	go func() {
		for d := range msgs {
			log.Printf("Received a message: %s", d.Body)
		}
	}()

	<-forever
}
