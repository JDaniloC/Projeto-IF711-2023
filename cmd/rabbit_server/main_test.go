package main

import (
	"encoding/json"
	"testing"

	utils "github.com/JDaniloC/Projeto-IF711-2023/pkg/server"
	"github.com/streadway/amqp"
)

func init() {
	amqpServer := utils.NewRabbitMQ(rabbitMQAddr)

	go func() {
		amqpServer.Start()
		amqpServer.Serve()
	}()
}

func BenchmarkRPCServer(b *testing.B) {
	request := &utils.Request{
		Link: "https://hackerspaces.org/", Depth: 1,
	}
	bytes, err := json.Marshal(request)
	if err != nil {
		b.Errorf("filed to mashall request: %v", err)
	}

	amqpClient := utils.NewRabbitMQ(rabbitMQAddr)
	err = amqpClient.Start()
	if err != nil {
		b.Errorf("failed to connect to RabbitMQ: %v", err)
	}
	defer amqpClient.Close()

	msgs, err := amqpClient.Channel.Consume(
		amqpClient.Response.Name, // queue
		"",                       // consumer
		false,                    // auto-ack
		false,                    // exclusive
		false,                    // no-local
		false,                    // no-wait
		nil,                      // args
	)
	if err != nil {
		b.Errorf("failed to register a consumer: %v", err)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if err := amqpClient.Channel.Publish(
			"",                      // exchange
			amqpClient.Request.Name, // routing key
			false,                   // mandatory
			false,                   // immediate
			amqp.Publishing{
				ContentType: "application/json",
				Body:        bytes,
			}); err != nil {
			b.Errorf("failed to publish the request: %v", err)
		}

		for b := range msgs {
			b.Ack(false)
			break
		}
	}
}
