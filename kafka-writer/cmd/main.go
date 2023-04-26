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
	// Create a new Kafka writer.
	writer := kafka.NewWriter(kafka.WriterConfig{
		Brokers: []string{"localhost:9092"},
		Topic:   "user-writes",
	})

	// Set up a signal handler to handle graceful shutdown.
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt)

	// Start an infinite loop that reads from stdin and sends messages to Kafka.
	fmt.Println("Enter user data to write:")
	for {
		// Read input from stdin.
		var id int
		var name string
		fmt.Print("ID: ")
		fmt.Scanln(&id)
		fmt.Print("Name: ")
		fmt.Scanln(&name)

		// Create a new user object and marshal it to JSON.
		user := User{id, name}
		jsonBytes, err := json.Marshal(user)
		if err != nil {
			log.Fatalf("Error marshaling JSON: %v", err)
		}

		// Send the message to Kafka.
		err = writer.WriteMessages(context.Background(), kafka.Message{
			Key:   []byte(fmt.Sprintf("%d", id)),
			Value: jsonBytes,
		})
		if err != nil {
			log.Fatalf("Error sending message to Kafka: %v", err)
		}

		// Wait for a signal or for the next input.
		select {
		case <-signals:
			fmt.Println("Shutting down...")
			writer.Close()
			return
		default:
			fmt.Println("Enter user data to write:")
		}
	}
}
