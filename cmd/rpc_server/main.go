package server

import (
	rpc "github.com/JDaniloC/Projeto-IF711-2023/pkg/server"
)

const (
	addr = ":1234" // Endereço para o servidor RPC
)

func main() {
	server := rpc.NewRPCServer(addr)
	server.Start()
}
