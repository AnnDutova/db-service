package main

import (
	"fmt"
	"github.com/confluentinc/confluent-kafka-go/kafka"
	"log"
	"os"
	"time"
)

type OrderPlacer struct {
	producer   *kafka.Producer
	topic      string
	deliveryCh chan kafka.Event
}

func NewOrderPlacer(k *kafka.Producer, topic string) *OrderPlacer {
	return &OrderPlacer{
		producer:   k,
		topic:      topic,
		deliveryCh: make(chan kafka.Event, 10000),
	}
}

func (op *OrderPlacer) placeOrder(orderType string, size int) error {
	var (
		format  = fmt.Sprintf("%s - %d", orderType, size)
		payload = []byte(format)
	)

	err := op.producer.Produce(
		&kafka.Message{
			TopicPartition: kafka.TopicPartition{Topic: &op.topic, Partition: kafka.PartitionAny},
			Value:          payload}, /*[]byte("FOO")}*/ //better to use protobuf
		op.deliveryCh,
	)
	if err != nil {
		log.Fatal(err)
	}
	<-op.deliveryCh
	fmt.Printf("placed order %s in the queue\n", format)
	return nil
}

func main() {
	topic := "HVSE"
	p, err := kafka.NewProducer(&kafka.ConfigMap{
		"bootstrap.servers": "localhost:9090",
		"client.id":         "foo",
		"acks":              "all"})

	if err != nil {
		fmt.Printf("Failed to create producer: %s\n", err)
		os.Exit(1)
	}

	op := NewOrderPlacer(p, topic)
	for i := 0; i < 1000; i++ {
		if err := op.placeOrder("market", i); err != nil {
			log.Fatal(err)
		}
		time.Sleep(time.Second * 3)
	}
	//fmt.Printf("%+v\n", e.String())
	//fmt.Printf("%+v\n", p)

}
