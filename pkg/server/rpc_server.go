package server

import (
	"fmt"
	"net"
	"net/rpc"

	"github.com/JDaniloC/Projeto-IF711-2023/pkg/runner"
)

type CrawlerRPC struct{}

type Request struct {
	Link  string
	Depth int
}

type Response struct {
	ValidLinks   []string
	InvalidLinks []string
}

func (c *CrawlerRPC) Crawl(request Request, response *Response) error {
	controller := runner.TimeoutCrawl(request.Link, request.Depth)
	response.ValidLinks = controller.ValidLinks.ToArray()
	response.InvalidLinks = controller.InvalidLinks.ToArray()
	return nil
}

type RPCServer struct {
	addr   string
	server *rpc.Server
}

func (r *RPCServer) Close() {
	// O servidor RPC não precisa ser fechado explicitamente
}

func (r *RPCServer) Start() error {
	r.server = rpc.NewServer()
	if err := r.server.Register(&CrawlerRPC{}); err != nil {
		return fmt.Errorf("error registering service: %s", err)
	}

	listener, err := net.Listen("tcp", r.addr)
	if err != nil {
		return fmt.Errorf("error opening listener: %s", err)
	}
	defer listener.Close()

	fmt.Println("Listening RPC server at", r.addr)

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Printf("Error accepting connection: %s", err)
			continue
		}
		go r.server.ServeConn(conn)
	}
}

func NewRPCServer(addr string) *RPCServer {
	server := RPCServer{
		addr: addr,
	}
	return &server
}
