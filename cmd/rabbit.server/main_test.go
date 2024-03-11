package main

import (
    "os"
	"net/rpc"
    "testing"

	utils "github.com/JDaniloC/Projeto-IF711-2023/pkg/server"
)

const (
    rpcServerAddr        = ":1124"
    rabbitMQAddrTest     = "amqp://guest:guest@localhost:5672/"
    rabbitMQQueueNameTest = "rpc_requests"
)

func TestMain(m *testing.M) {
    // Iniciar servidor RPC
    rpcServer := utils.NewRPCServer(rpcServerAddr)
    go rpcServer.Start()

    // Iniciar servidor RabbitMQ
    rabbitMQServer := utils.NewRabbitMQServer(rabbitMQAddrTest)
    if err := rabbitMQServer.Start(); err != nil {
        panic(err)
    }

    // Executar os testes
    exitCode := m.Run()

    // Encerrar servidores
    os.Exit(exitCode)
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
