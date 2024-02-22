package main

import (
	rpc "github.com/JDaniloC/Projeto-IF711-2023/pkg/server"
)

const (
	addr = ":1234"
)

func main() {
	server := rpc.NewRPCServer(addr)
	server.Start()
}
