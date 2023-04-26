package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/segmentio/kafka-go"
	"log"
	"os"
	"os/signal"
)

type User struct {
	ID   int
	Name string
}

func main() {
	// Create a new Kafka reader.
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers: []string{"localhost:9092"},
		Topic:   "user-writes",
	})

	// Set up a signal handler to handle graceful shutdown.
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt)

	// Start an infinite loop that reads messages from Kafka and processes them.
	fmt.Println("Waiting for user data to read...")
	for {
		// Read the next message from Kafka.
		message, err := reader.ReadMessage(context.Background())
		if err != nil {
			log.Fatalf("Error reading message from Kafka: %v", err)
		}

		// Unmarshal the message body to a user object.
		var user User
		err = json.Unmarshal(message.Value, &user)
		if err != nil {
			log.Fatalf("Error unmarshaling JSON: %v", err)
		}

		// Print the user data.
		fmt.Printf("Read user: %v\n", user)

		// Wait for a signal or for the next message.
		select {
		case <-signals:
			fmt.Println("Shutting down...")
			reader.Close()
			return
		default:
			fmt.Println("Waiting for user data to read...")
		}
	}
}
