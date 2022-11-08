package main

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis"
)

type Message struct {
	Sender   string `json:"sender"`
	Receiver string `json:"receiver"`
	Message  string `json:"message"`
}

type MessagesList struct {
	TotalMessages int       `json:"total_messages"`
	Messages      []Message `json:"messages_list"`
}

func main() {
	r := gin.Default()

	r.GET("/message/list", ReadRedis)

	r.Run("localhost:8081")
}

func ReadRedis(c *gin.Context) {
	sender_requested := c.Query("sender")
	receiver_requested := c.Query("receiver")

	client := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})

	val, err := client.Get("total_messages").Result()
	if err != nil {
		fmt.Println(err)
	}

	total_messages, err := strconv.Atoi(val)
	real_conversation_messages := total_messages

	messages_list := MessagesList{}

	for i := total_messages - 1; i >= 0; i-- {
		val, err = client.Get("message_" + strconv.Itoa(i)).Result()
		if err != nil {
			fmt.Println(err)
		}

		sender := strings.Split(strings.Split(val, " ")[0], "{")[1]
		receiver := strings.Split(val, " ")[1]
		message_content := strings.Split(strings.Split(val, " ")[2], "}")[0]

		message := Message{
			Sender:   sender,
			Receiver: receiver,
			Message:  message_content,
		}

		if sender_requested == message.Sender && receiver_requested == message.Receiver {
			messages_list.Messages = append(messages_list.Messages, message)
		} else {
			real_conversation_messages = real_conversation_messages - 1
		}
	}

	messages_list.TotalMessages = real_conversation_messages

	c.IndentedJSON(http.StatusFound, messages_list)
}
