package main

import "github.com/JDaniloC/Projeto-IF711-2023/pkg/server"

func main() {
	server := server.NewUDPServer(":6250")
	server.Start()
}
