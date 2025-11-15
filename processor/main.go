package main

import (
	"context"
	"encoding/json"
	"log"
	"time"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/segmentio/kafka-go"

	"go-analytics/pkg/config"
	db "go-analytics/db/sqlc"
	"go-analytics/pkg/types"
)

func main() {
	log.Println("Starting processor service...")
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("could not load config: %v", err)
	}

	// Connect to PostgreSQL
	dbpool, err := pgxpool.New(context.Background(), cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("Unable to connect to database: %v\n", err)
	}
	defer dbpool.Close()

	// `sqlc` queries
	queries := db.New(dbpool)

	// Connect to Kafka Consumer
	kafkaReader := kafka.NewReader(kafka.ReaderConfig{
		Brokers:        cfg.KafkaBrokerURLs,
		Topic:          cfg.KafkaTopic,
		GroupID:        cfg.KafkaGroupID,
		MinBytes:       10e3, // 10KB
		MaxBytes:       10e6, // 10MB
		CommitInterval: 1 * time.Second, // Commit offsets every second
	})
	defer kafkaReader.Close()

	log.Println("Processor is running and waiting for messages...")

	// Infinite loop to process messages
	for {
		// This blocks until a message is available
		msg, err := kafkaReader.FetchMessage(context.Background())
		if err != nil {
			log.Printf("could not fetch message: %v\n", err)
			continue
		}

		var event types.Event
		if err := json.Unmarshal(msg.Value, &event); err != nil {
			log.Printf("could not unmarshal message: %v\n", err)
			// Don't commit, but log the error
			continue
		}

		// Create the event in the database
		params := db.CreateEventParams{
			SiteID:    event.SiteID,
			EventType: event.EventType,
			Path:      event.Path,
			UserID:    event.UserID,
			Timestamp: pgtype.Timestamptz{Time: event.Timestamp, Valid: true},
		}

		_, err = queries.CreateEvent(context.Background(), params)
		if err != nil {
			log.Printf("could not create event in db: %v\n", err)
			// Do not commit, so we can retry
			continue
		}

		// Log and commit the message offset to Kafka
		log.Printf("Processed event for site: %s, path: %s", event.SiteID, event.Path)
		if err := kafkaReader.CommitMessages(context.Background(), msg); err != nil {
			log.Printf("failed to commit message: %v\n", err)
		}
	}
}