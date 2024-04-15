package main

import (
	"encoding/json"
	"fmt"
	"log"
	"log/slog"
	"os"

	"github.com/chaewonkong/msa-link-api/link"
	"github.com/chaewonkong/msa-link-scraper/config"
	"github.com/chaewonkong/msa-link-scraper/convert"
	"github.com/chaewonkong/msa-link-scraper/scraper"
	"github.com/chaewonkong/msa-link-scraper/transport"
)

func main() {
	cfg := config.NewAppConfig()
	requester := transport.NewHTTPRequester(cfg.APIHost)

	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	// rabbitMQ
	mq := transport.NewRabbitMQ(cfg)
	defer mq.Close()

	q, err := mq.Ch.QueueDeclare(
		cfg.Queue.Name, // name
		true,           // durable
		false,          // delete when unused
		false,          // exclusive
		false,          // no-wait
		nil,            // arguments
	)
	if err != nil {
		log.Fatal(err)
	}

	msgs, err := mq.Ch.Consume(
		q.Name, // queue
		"",     // consumer
		false,  // auto-ack
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
			ogData, err := scraper.GetOpenGraphTags(q.URL)
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

			logger.Info("Response", "body", fmt.Sprintf("%v", resp.String()), "status", resp.GetStatusCode(), "message", resp.GetMessage())
			err = d.Ack(false)
			if err != nil {
				logger.Error("Failed to ack message", err)
			}
		}
	}()

	<-forever
}

// kafka message handler
// message
// scraping
// convert object
// send request
// return
