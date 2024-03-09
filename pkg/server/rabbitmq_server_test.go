package server

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/JDaniloC/Projeto-IF711-2023/pkg/runner"
	"github.com/streadway/amqp"
)

const (
	rabbitMQAddr = "amqp://guest:guest@localhost:5672/"
	queueName    = "rpc_requests"
)

type RequestMQTest struct {
	Link  string
	Depth int
}

type ResponseMQTest struct {
	ValidLinks   []string
	InvalidLinks []string
}

type RabbitMQServerMQTest struct {
	addr    string
	channel *amqp.Channel
}

func NewRabbitMQServerMQTest(addr string) *RabbitMQServerMQTest {
	return &RabbitMQServerMQTest{
		addr: addr,
	}
}

func (r *RabbitMQServerMQTest) StartMQ() error {
	rabbitConn, err := amqp.Dial(rabbitMQAddr)
	if err != nil {
		return fmt.Errorf("failed to connect to RabbitMQ: %v", err)
	}
	defer rabbitConn.Close()

	channel, err := rabbitConn.Channel()
	if err != nil {
		return fmt.Errorf("failed to open a channel: %v", err)
	}
	defer channel.Close()
	r.channel = channel

	queue, err := channel.QueueDeclare(
		queueName, // name
		false,     // durable
		false,     // delete when unused
		false,     // exclusive
		false,     // no-wait
		nil,       // arguments
	)
	if err != nil {
		return fmt.Errorf("failed to declare a queue: %v", err)
	}

	msgs, err := channel.Consume(
		queue.Name, // queue
		"",         // consumer
		true,       // auto-ack
		false,      // exclusive
		false,      // no-local
		false,      // no-wait
		nil,        // args
	)
	if err != nil {
		return fmt.Errorf("failed to register a consumer: %v", err)
	}

	forever := make(chan bool)

	go func() {
		for d := range msgs {
			var request RequestMQTest
			if err := json.Unmarshal(d.Body, &request); err != nil {
				log.Printf("failed to decode message: %v", err)
				continue
			}

			// Process the request
			response := r.processRequest(&request)

			// Send the response back
			if err := r.channel.Publish(
				"",          // exchange
				d.ReplyTo,   // routing key
				false,       // mandatory
				false,       // immediate
				amqp.Publishing{
					ContentType:   "application/json",
					CorrelationId: d.CorrelationId,
					Body:          responseToBytes(response),
				}); err != nil {
				log.Printf("failed to publish a response: %v", err)
				continue
			}
		}
	}()

	fmt.Printf("Listening RPC server at %s\n", r.addr)
	<-forever

	return nil
}

func responseToBytes(response *ResponseMQTest) []byte {
	bytes, _ := json.Marshal(response)
	return bytes
}

func (r *RabbitMQServerMQTest) processRequest(request *RequestMQTest) *ResponseMQTest {
	controller := runner.TimeoutCrawl(request.Link, request.Depth)
	return &ResponseMQTest{
		ValidLinks:   controller.ValidLinks.ToArray(),
		InvalidLinks: controller.InvalidLinks.ToArray(),
	}
}
