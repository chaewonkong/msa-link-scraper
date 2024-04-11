package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"
	"time"

	config "github.com/chaewonkong/msa-link-api/config"
	"github.com/chaewonkong/msa-link-api/link"
	"github.com/chaewonkong/msa-link-scraper/fetch"
	amqp "github.com/rabbitmq/amqp091-go"
)

func NewHTTPClient() *http.Client {
	return &http.Client{
		Timeout: 5 * time.Second,
	}
}

func main() {
	cfg := config.NewAppConfig()
	client := NewHTTPClient()

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
			// todo: scrape link then save to db

			q := link.QueuePayload{}

			err := json.Unmarshal(d.Body, &q)
			if err != nil {
				logger.Error("Failed to unmarshal message", err)
			}

			fetcher := fetch.NewFetcher(logger)
			ogData, err := fetcher.GetOpenGraphTags(q.URL)
			if err != nil {
				logger.Error("Failed to fetch Open Graph tags", err)
			}

			// save img: call Link API, make Link API to save to db
			updatePayload := link.UpdatePayload{
				ID: q.ID,
			}

			if img, exists := ogData["og:image"]; exists {
				updatePayload.ThumbnailImage = img
			}

			if title, exists := ogData["og:title"]; exists {
				updatePayload.Title = title
			}

			if description, exists := ogData["og:description"]; exists {
				updatePayload.Description = description
			}

			jsonPayload, err := json.Marshal(updatePayload)
			if err != nil {
				logger.Error("Failed to marshal payload", err)
			}

			req, err := http.NewRequest("PATCH", "http://localhost:8080/link", bytes.NewBuffer(jsonPayload))
			if err != nil {
				logger.Error("Failed to create request", err)
			}

			req.Header.Set("Content-Type", "application/json")

			res, err := client.Do(req)
			if err != nil {
				logger.Error("Failed to send request", err)
			}

			logger.Info("Response", fmt.Sprintf("%v", res.StatusCode), res.Body)
			res.Body.Close()
		}
	}()

	<-forever
}
