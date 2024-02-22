package server

import (
	"fmt"
	"net"
	"net/http"
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

func (c *CrawlerRPC) Crawl(request *Request, response *Response) error {
	controller := runner.TimeoutCrawl(request.Link, request.Depth)
	response.ValidLinks = controller.ValidLinks.ToArray()
	response.InvalidLinks = controller.InvalidLinks.ToArray()
	return nil
}

type RPCServer struct {
	addr   string
	server net.Listener
}

func (r *RPCServer) Start() (err error) {
	rpcCrawler := new(CrawlerRPC)
	if err := rpc.Register(rpcCrawler); err != nil {
		return fmt.Errorf("error registering service: %s", err)
	}

	rpc.HandleHTTP()

	r.server, err = net.Listen("tcp", r.addr)
	if err != nil {
		return fmt.Errorf("error opening listener: %s", err)
	}
	fmt.Println("Listening RPC server at", r.addr)
	defer r.server.Close()

	http.Serve(r.server, nil)
	return nil
}

func NewRPCServer(addr string) *RPCServer {
	server := RPCServer{
		addr: addr,
	}
	return &server
}
