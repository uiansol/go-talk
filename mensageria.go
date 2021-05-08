package main

import (
	"log"

	"github.com/streadway/amqp"
)

func publishMail(userName string) {
	// Open connection with RabbitMQ server
	conn, err := amqp.Dial("amqp://guest:guest@rabbitmq/")
	if err != nil {
		println()
		log.Fatalf("fail to dial to rabbitmq: %v", err)
	}
	defer conn.Close()

	// Open a channel to RabbitMQ instance
	ch, err := conn.Channel()
	if err != nil {
		log.Fatalf("fail to create channel: %v", err)
	}
	defer ch.Close()

	// Declare a queue to be used
	q, err := ch.QueueDeclare(
		"MailQueue",
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		log.Fatalf("fail to declare queue: %v", err)
		log.Fatal(q)
	}
	// Publish a message to the queue
	err = ch.Publish(
		"",
		"MailQueue",
		false,
		false,
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(userName),
		},
	)
	if err != nil {
		log.Fatalf("fail to publish message: %v", err)
	}
}
