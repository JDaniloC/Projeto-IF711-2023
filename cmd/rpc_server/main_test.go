package main

import (
	"net/rpc"
	"testing"

	utils "github.com/JDaniloC/Projeto-IF711-2023/pkg/server"
)

const (
	rpcServerAddr = ":1124"
)

func init() {
	rpcServer := utils.NewRPCServer(rpcServerAddr)
	go rpcServer.Start()
}

func BenchmarkRPCServer(b *testing.B) {
	var response utils.Response
	request := &utils.Request{
		Link:  "https://hackerspaces.org/",
		Depth: 2,
	}

	client, err := rpc.DialHTTP("tcp", "localhost"+rpcServerAddr)
	if err != nil {
		b.Error("could not connect to server: " + err.Error())
	}
	defer client.Close()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if err := client.Call("CrawlerRPC.Crawl", request, &response); err != nil {
			b.Error(err.Error())
		}
	}
}
