package datateam

import (
	"fmt"
	"github.com/confluentinc/confluent-kafka-go/kafka"
	"log"
)

func main() {
	topic := "HVSE"

	consumer, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers": "localhost:9090",
		"group.id":          "foo_team",
		"auto.offset.reset": "smallest"})
	if err != nil {
		log.Fatal(err)
	}

	err = consumer.Subscribe(topic, nil)
	if err != nil {
		log.Fatal(err)
	}

	for {
		ev := consumer.Poll(100)
		switch e := ev.(type) {
		case *kafka.Message:
			fmt.Printf("datateam send message to the queue %s\n", string(e.Value))
		case *kafka.Error:
			fmt.Printf("%v\n", e)
		}
	}
}
