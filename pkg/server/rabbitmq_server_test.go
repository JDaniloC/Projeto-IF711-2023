package server

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/streadway/amqp"
)

const (
	rabbitMQAddr = "amqp://guest:guest@localhost:5672/"
	queueName    = "crawl"
)

func init() {
	amqpServer := NewRabbitMQ(rabbitMQAddr)

	go func() {
		amqpServer.Start()
		amqpServer.Serve()
	}()
}

func TestAMQPServer_Connection(t *testing.T) {
	conn, err := amqp.Dial(rabbitMQAddr)
	if err != nil {
		t.Errorf("failed to connect to RabbitMQ: %v", err)
	}
	defer conn.Close()
}

func TestAMQPServer_Request(t *testing.T) {
	request := &Request{
		"https://hackerspaces.org/", 1,
	}
	bytes, err := json.Marshal(request)
	if err != nil {
		t.Errorf("filed to mashall request: %v", err)
	}

	amqpClient := NewRabbitMQ(rabbitMQAddr)
	err = amqpClient.Start()
	if err != nil {
		t.Errorf("failed to connect to RabbitMQ: %v", err)
	}
	defer amqpClient.Close()

	msgs, err := amqpClient.channel.Consume(
		amqpClient.response.Name, // queue
		"",                       // consumer
		true,                     // auto-ack
		false,                    // exclusive
		false,                    // no-local
		false,                    // no-wait
		nil,                      // args
	)
	if err != nil {
		t.Errorf("failed to register a consumer: %v", err)
	}

	t.Run("Send AMQP request", func(t *testing.T) {
		if err := amqpClient.channel.Publish(
			"",                      // exchange
			amqpClient.request.Name, // routing key
			false,                   // mandatory
			false,                   // immediate
			amqp.Publishing{
				ContentType: "application/json",
				Body:        bytes,
			}); err != nil {
			t.Errorf("failed to publish the request: %v", err)
		}

		for d := range msgs {
			fmt.Printf("Received a message: %s", d.Body)
			break
		}
	})
}
