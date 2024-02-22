package main

import (
	"fmt"
	"log"
	"net/rpc"
)

const (
	serverAddress = "localhost:1234"
)

type Request struct {
	Link  string `json:"link"`
	Depth int    `json:"depth"`
}

type Response struct {
	ValidLinks   []string `json:"validLinks"`
	InvalidLinks []string `json:"invalidLinks"`
}

func main() {
	client, err := rpc.Dial("rpc", serverAddress)
	if err != nil {
		log.Fatalf("Erro ao conectar ao servidor RPC: %s", err)
	}
	defer client.Close()

	request := Request{
		Link:  "https://hackerspaces.org/",
		Depth: 2,
	}

	// Chama o método remoto no servidor RPC
	var response Response
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
