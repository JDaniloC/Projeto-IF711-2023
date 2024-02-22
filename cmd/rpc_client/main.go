package main

import (
	"fmt"
	"log"
	"net/rpc"

	utils "github.com/JDaniloC/Projeto-IF711-2023/pkg/server"
)

const (
	serverAddress = "localhost:1234"
)

func main() {
	var response utils.Response
	request := &utils.Request{
		Link:  "https://hackerspaces.org/",
		Depth: 2,
	}

	client, err := rpc.DialHTTP("tcp", serverAddress)
	if err != nil {
		log.Fatalf("could not connect to server: " + err.Error())
	}
	defer client.Close()

	if err := client.Call("CrawlerRPC.Crawl", request, &response); err != nil {
		log.Fatalf("Erro ao chamar método remoto: %s", err)
	}

	// Exibe os resultados recebidos do servidor RPC
	fmt.Println("Links válidos do servidor:")
	for _, link := range response.ValidLinks {
		fmt.Println(link)
	}

	fmt.Println("\nLinks inválidos do servidor:")
	for _, link := range response.InvalidLinks {
		fmt.Println(link)
	}
}
