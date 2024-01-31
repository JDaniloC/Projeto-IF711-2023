package server

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net"

	runner "github.com/JDaniloC/Projeto-IF711-2023/pkg/runner"
)

type TCPServer struct {
	addr   string
	server net.Listener
}

func (t *TCPServer) handleConnection(conn net.Conn) {
	defer conn.Close()

	// Create the object to receive and answer requests
	rw := bufio.NewReadWriter(bufio.NewReader(conn), bufio.NewWriter(conn))

	for {
		// Get the serialized string message
		s, err := rw.ReadString('\n')

		if err != nil {
			rw.WriteString(fmt.Sprintf("Error on read the request: %s", err))
			rw.Flush()
			return
		}

		// Deserialize the request
		var request map[string]interface{}
		err = json.Unmarshal([]byte(s), &request)
		if err != nil {
			rw.WriteString(fmt.Sprintf("Erro on deserialize to JSON: %s", err))
			rw.Flush()
			return
		}

		// Extract the request params
		link, ok1 := request["link"].(string)
		depth, ok2 := request["depth"].(int)
		if !ok1 || !ok2 {
			rw.WriteString("Missing request params link or depth")
			rw.Flush()
			return
		}

		// Run the crawler and mount the result
		controller := runner.TimeoutCrawl(link, depth)
		response := map[string]interface{}{
			"validLinks":   controller.ValidLinks,
			"invalidLinks": controller.InvalidLinks,
		}

		// Send response to the client
		responseJSON, _ := json.Marshal(response)
		rw.Write([]byte(responseJSON))
		rw.Flush()
	}
}

func (t *TCPServer) handleConnections() (err error) {
	for {
		conn, err := t.server.Accept()
		if err != nil || conn == nil {
			println("Error on accept connection")
			break
		}

		go t.handleConnection(conn)
	}
	return
}

func (t *TCPServer) Close() (err error) {
	return t.server.Close()
}

func (t *TCPServer) Start() (err error) {
	t.server, err = net.Listen("tcp", t.addr)
	if err != nil {
		fmt.Println("Erro on open the connection:", err)
		return err
	}
	defer t.Close()
	fmt.Println("Listening TCP server at", t.addr)

	for {
		conn, err := t.server.Accept()
		if err != nil {
			println("Error on accept connection")
			break
		}
		if conn == nil {
			println("Error on create connection")
			break
		}
		t.handleConnections()
	}
	return
}

func NewTCPServer(addr string) *TCPServer {
	server := TCPServer{
		addr: addr,
	}
	return &server
}
