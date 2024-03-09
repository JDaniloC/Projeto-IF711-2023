package server

import (
    "encoding/json"
    "fmt"
    "log"
    "net/http"
    "net/rpc"

    "github.com/JDaniloC/Projeto-IF711-2023/pkg/runner"
    "github.com/streadway/amqp"
)

type CrawlerRabbitMQ struct{}

type RequestMQ struct {
    Link  string
    Depth int
}

type ResponseMQ struct {
    ValidLinks   []string
    InvalidLinks []string
}

func (c *CrawlerRabbitMQ) Crawl(request *RequestMQ, response *ResponseMQ) error {
    controller := runner.TimeoutCrawl(request.Link, request.Depth)
    response.ValidLinks = controller.ValidLinks.ToArray()
    response.InvalidLinks = controller.InvalidLinks.ToArray()
    return nil
}

type RabbitMQServer struct {
    addr     string
    rabbitMQ *amqp.Connection
}

func (r *RabbitMQServer) Start() (err error) {
    crawler := new(CrawlerRabbitMQ)
    if err := rpc.Register(crawler); err != nil {
        return fmt.Errorf("error registering service: %s", err)
    }

    http.HandleFunc("/rpc", r.handleRPCRequest)

    r.rabbitMQ, err = amqp.Dial("amqp://guest:guest@localhost:5672/")
    if err != nil {
        return fmt.Errorf("failed to connect to RabbitMQ: %v", err)
    }
    defer r.rabbitMQ.Close()

    fmt.Println("Listening RPC server at", r.addr)
    log.Fatal(http.ListenAndServe(r.addr, nil))
    return nil
}

func (r *RabbitMQServer) handleRPCRequest(w http.ResponseWriter, req *http.Request) {
    var request RequestMQ
    var response ResponseMQ

    err := json.NewDecoder(req.Body).Decode(&request)
    if err != nil {
        http.Error(w, "failed to decode request", http.StatusBadRequest)
        return
    }

    crawler := new(CrawlerRabbitMQ)
    err = crawler.Crawl(&request, &response)
    if err != nil {
        http.Error(w, "failed to process request", http.StatusInternalServerError)
        return
    }

    err = json.NewEncoder(w).Encode(&response)
    if err != nil {
        http.Error(w, "failed to encode response", http.StatusInternalServerError)
        return
    }
}

func NewRabbitMQServer(addr string) *RabbitMQServer {
    server := RabbitMQServer{
        addr: addr,
    }
    return &server
}
