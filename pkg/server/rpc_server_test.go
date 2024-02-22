package server

import (
	"fmt"
	"net/rpc"
	"testing"
)

const (
	rpcAddr = ":1124"
)

func init() {
	rpcServer := NewRPCServer(rpcAddr)

	go func() {
		rpcServer.Start()
	}()
}

func TestRPCServer_Connection(t *testing.T) {
	conn, err := rpc.DialHTTP("tcp", "localhost"+rpcAddr)
	if err != nil {
		t.Error("could not connect to server: ", err)
	}
	conn.Close()
}

func TestRPCServer_Request(t *testing.T) {
	var reply Response
	request := &Request{
		"https://hackerspaces.org/", 2,
	}

	client, err := rpc.DialHTTP("tcp", "localhost"+rpcAddr)
	if err != nil {
		t.Error("could not connect to server: ", err)
	}
	defer client.Close()

	t.Run("Send RPC request", func(t *testing.T) {
		err = client.Call("CrawlerRPC.Crawl", request, &reply)
		if err != nil {
			t.Error("Client invocation error: ", err)
		}
		fmt.Printf("Answer from server: %s\n", reply)
	})
}
