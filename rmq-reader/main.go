package main

import (
	"fmt"
	"github.com/streadway/amqp"
)

func main() {
	conn, err := amqp.Dial("amqp://user:password@localhost:5672/")
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	fmt.Println("Success")
	ch, err := conn.Channel()
	if err != nil {
		panic(err)
	}
	defer ch.Close()

	msg, err := ch.Consume(
		"DataBase",
		"",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		panic(err)
	}

	forever := make(chan bool)
	go func() {
		for d := range msg {
			fmt.Printf("Recived msg: %s\n", d.Body)
		}
	}()
	fmt.Println("Successfully connect to our Rabbitmq instance")
	fmt.Println("waiting for msgs")
	<-forever

}
