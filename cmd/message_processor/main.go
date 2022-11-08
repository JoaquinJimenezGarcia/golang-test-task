package main

import (
	"fmt"
	"log"
	"strconv"

	"github.com/go-redis/redis"
	"github.com/streadway/amqp"
)

type Message struct {
	Sender   string
	Receiver string
	Message  string
}

func main() {
	// Define RabbitMQ server URL.
	amqpServerURL := "amqp://user:password@localhost:7001"

	// Create a new RabbitMQ connection.
	connectRabbitMQ, err := amqp.Dial(amqpServerURL)
	if err != nil {
		panic(err)
	}
	defer connectRabbitMQ.Close()

	// Opening a channel to our RabbitMQ instance over
	// the connection we have already established.
	channelRabbitMQ, err := connectRabbitMQ.Channel()
	if err != nil {
		panic(err)
	}
	defer channelRabbitMQ.Close()

	// Subscribing to QueueService1 for getting messages.
	messages, err := channelRabbitMQ.Consume(
		"MessagesQueue", // queue name
		"",              // consumer
		true,            // auto-ack
		false,           // exclusive
		false,           // no local
		false,           // no wait
		nil,             // arguments
	)
	if err != nil {
		log.Println(err)
	}

	// Build a welcome message.
	log.Println("Successfully connected to RabbitMQ")
	log.Println("Waiting for messages")

	// Make a channel to receive messages into infinite loop.
	forever := make(chan bool)

	message_id := 0

	go func() {
		for message := range messages {
			// For example, show received message in a console.
			log.Printf(" > Received message: %s\n", message.Body)
			client := redis.NewClient(&redis.Options{
				Addr:     "localhost:6379",
				Password: "",
				DB:       0,
			})

			err = client.Set("message_"+strconv.Itoa(message_id), message.Body, 0).Err()
			// if there has been an error setting the value
			// handle the error
			if err != nil {
				fmt.Println(err)
			}

			val, err := client.Get("message_" + strconv.Itoa(message_id)).Result()
			if err != nil {
				fmt.Println(err)
			}

			fmt.Println(val)
			message_id = message_id + 1

			err = client.Set("total_messages", message_id, 0).Err()
			// if there has been an error setting the value
			// handle the error
			if err != nil {
				fmt.Println(err)
			}
		}
	}()

	<-forever
}
