package main

import (
	tcp "github.com/JDaniloC/Projeto-IF711-2023/pkg/server"
)

const (
	addr = ":1123"
)

func main() {
	server := tcp.NewTCPServer(addr)
	server.Start()
}
