package main

import (
	"github.com/JDaniloC/Projeto-IF711-2023/pkg/server"
)

const (
	addr          = ":1234"
	rabbitMQAddr  = "amqp://guest:guest@localhost:5672/"
	rabbitMQQueue = "rpc_requests"
)

func main() {
	rabbitMQServer := server.NewRabbitMQServer(rabbitMQAddr)
	if err := rabbitMQServer.Start(); err != nil {
		panic(err)
	}
}
