package main

import (
	amqp "github.com/JDaniloC/Projeto-IF711-2023/pkg/server"
)

const (
	rabbitMQAddr = "amqp://guest:guest@localhost:5672/"
)

func main() {
	amqpServer := amqp.NewRabbitMQ(rabbitMQAddr)
	amqpServer.Start()
	amqpServer.Serve()
}
