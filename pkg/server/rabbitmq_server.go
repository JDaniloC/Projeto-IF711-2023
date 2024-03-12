package server

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/JDaniloC/Projeto-IF711-2023/pkg/runner"
	"github.com/streadway/amqp"
)

type CrawlerAMQP = CrawlerRPC

type RabbitMQ struct {
	addr     string
	channel  *amqp.Channel
	rabbitMQ *amqp.Connection
	request  amqp.Queue
	response amqp.Queue
}

func (r *RabbitMQ) processRequest(request *Request) *Response {
	controller := runner.TimeoutCrawl(request.Link, request.Depth)
	return &Response{
		ValidLinks:   controller.ValidLinks.ToArray(),
		InvalidLinks: controller.InvalidLinks.ToArray(),
	}
}

func (r *RabbitMQ) Start() (err error) {
	r.rabbitMQ, err = amqp.Dial(r.addr)
	if err != nil {
		log.Fatalf("%s: %s", "Failed to serve RabbitMQ", err)
		return err
	}

	r.channel, err = r.rabbitMQ.Channel()
	if err != nil {
		log.Fatalf("%s: %s", "Failed to open a channel", err)
		return err
	}

	r.request, err = r.channel.QueueDeclare(
		"request", // name
		false,     // durable
		false,     // delete when unused
		false,     // exclusive
		false,     // no-wait
		nil,       // arguments
	)
	if err != nil {
		log.Fatalf("%s: %s", "Failed to declare a queue", err)
	}

	r.response, err = r.channel.QueueDeclare(
		"reponse", // name
		false,     // durable
		false,     // delete when unused
		false,     // exclusive
		false,     // no-wait
		nil,       // arguments
	)
	if err != nil {
		log.Fatalf("%s: %s", "Failed to declare a queue", err)
	}

	fmt.Println("Listening RPC server at", r.addr)
	return nil
}

func (r *RabbitMQ) Serve() {
	msgs, err := r.channel.Consume(
		r.request.Name, // queue
		"",             // consumer
		true,           // auto-ack
		false,          // exclusive
		false,          // no-local
		false,          // no-wait
		nil,            // args
	)
	if err != nil {
		log.Printf("Failed to register a consumer: %v", err)
		return
	}

	for d := range msgs {
		var request Request
		if err := json.Unmarshal(d.Body, &request); err != nil {
			log.Printf("failed to decode message: %v", err)
			continue
		}

		// Process the request
		response := r.processRequest(&request)
		bytes, err := json.Marshal(response)
		if err != nil {
			log.Printf("failed to mashal response: %v", err)
			continue
		}

		// Send the response back
		if err := r.channel.Publish(
			"",              // exchange
			r.response.Name, // routing key
			false,           // mandatory
			false,           // immediate
			amqp.Publishing{
				ContentType:   "application/json",
				CorrelationId: d.CorrelationId,
				Body:          bytes,
			}); err != nil {
			log.Printf("failed to publish a response: %v", err)
			continue
		}
	}
}

func (r *RabbitMQ) Close() {
	r.channel.Close()
	r.rabbitMQ.Close()
}

func NewRabbitMQ(addr string) *RabbitMQ {
	server := RabbitMQ{
		addr: addr,
	}
	return &server
}
