package main

import (
	"encoding/json"
	"fmt"
	"log"
	"log/slog"
	"os"

	config "github.com/chaewonkong/msa-link-api/config"
	"github.com/chaewonkong/msa-link-api/link"
	"github.com/chaewonkong/msa-link-scraper/convert"
	"github.com/chaewonkong/msa-link-scraper/fetch"
	"github.com/chaewonkong/msa-link-scraper/transport"
	amqp "github.com/rabbitmq/amqp091-go"
)

func main() {
	cfg := config.NewAppConfig()
	requester := transport.NewHTTPRequester("http://localhost:8080")

	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

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

			q := link.QueuePayload{}

			err := json.Unmarshal(d.Body, &q)
			if err != nil {
				logger.Error("Failed to unmarshal message", err)
			}

			// fetch og tags
			fetcher := fetch.NewFetcher(logger)
			ogData, err := fetcher.GetOpenGraphTags(q.URL)
			if err != nil {
				logger.Error("Failed to fetch Open Graph tags", err)
			}

			// save og tags
			updatePayload := convert.MapToUpdatePayload(ogData)
			updatePayload.ID = q.ID

			resp, err := requester.UpdateLink(updatePayload)
			if err != nil {
				logger.Error("failed to send request", err)
			}

			logger.Info("Response", "body", fmt.Sprintf("%v", resp.Body), "status", resp.StatusCode, "message", resp.Message)
		}
	}()

	<-forever
}
