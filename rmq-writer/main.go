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

	_, err = ch.QueueDeclare("DataBase",
		false,
		false,
		false,
		false,
		nil,
	)

	if err != nil {
		fmt.Println(err)
		panic(err)
	}

	err = ch.Publish("", "DataBase", false, false, amqp.Publishing{
		ContentType: "text/plain",
		Body:        []byte("World"),
	})
	if err != nil {
		fmt.Println(err)
		panic(err)
	}
	fmt.Println("Publish message")

}
