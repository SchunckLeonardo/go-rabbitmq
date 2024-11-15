package main

import (
	"fmt"
	"github.com/SchunckLeonardo/handling-events/pkg/rabbitmq"
	amqp "github.com/rabbitmq/amqp091-go"
)

func main() {
	ch, err := rabbitmq.OpenChannel()
	if err != nil {
		panic(err)
	}
	defer ch.Close()

	msgs := make(chan amqp.Delivery)

	go func() {
		err = rabbitmq.Consume(ch, msgs)
		if err != nil {
			panic(err)
		}
	}()

	for msg := range msgs {
		fmt.Println(string(msg.Body))
		msg.Ack(false)
	}

}
