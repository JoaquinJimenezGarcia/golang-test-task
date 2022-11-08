package main

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/streadway/amqp"
)

type Message struct {
	Sender   string `json:"sender"`
	Receiver string `json:"receiver"`
	Message  string `json:"message"`
}

func main() {
	r := gin.Default()

	r.GET("/test", func(c *gin.Context) {
		c.JSON(200, "worked")
	})

	r.POST("/message", SendMessage)

	r.Run("localhost:8080")
}

func SendMessage(c *gin.Context) {
	var message Message

	// Define RabbitMQ server URL.
	amqpServerURL := "amqp://user:password@localhost:7001"

	// Create a new RabbitMQ connection.
	connectRabbitMQ, err := amqp.Dial(amqpServerURL)
	if err != nil {
		panic(err)
	}
	defer connectRabbitMQ.Close()

	// Let's start by opening a channel to our RabbitMQ
	// instance over the connection we have already
	// established.
	channelRabbitMQ, err := connectRabbitMQ.Channel()
	if err != nil {
		panic(err)
	}
	defer channelRabbitMQ.Close()

	// With the instance and declare Queues that we can
	// publish and subscribe to.
	_, err = channelRabbitMQ.QueueDeclare(
		"MessagesQueue", // queue name
		true,            // durable
		false,           // auto delete
		false,           // exclusive
		false,           // no wait
		nil,             // arguments
	)
	if err != nil {
		panic(err)
	}

	if err := c.BindJSON(&message); err != nil {
		c.IndentedJSON(http.StatusBadRequest, err)
		return
	}

	RabbitMessage := amqp.Publishing{
		ContentType: "application/json",
		Body:        []byte(fmt.Sprintf("%v", message)),
	}

	// Attempt to publish a message to the queue.
	if err := channelRabbitMQ.Publish(
		"",              // exchange
		"MessagesQueue", // queue name
		false,           // mandatory
		false,           // immediate
		RabbitMessage,   // message to publish
	); err != nil {
		return
	}

	c.IndentedJSON(http.StatusCreated, message)
}
