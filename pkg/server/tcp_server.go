package server

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net"

	utils "github.com/JDaniloC/Projeto-IF711-2023/internal/utils"
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
			errorMsg, _ := json.Marshal(map[string]interface{}{
				"error": fmt.Sprintf("Error on read the request: %s", err),
			})
			rw.Write(append(errorMsg, '\n'))
			rw.Flush()
			return
		}

		// Deserialize the request
		request := &utils.Request{}
		if err := json.Unmarshal([]byte(s), &request); err != nil {
			errorMsg, _ := json.Marshal(map[string]interface{}{
				"error": fmt.Sprintf("Erro on deserialize to JSON: %s", err),
			})
			rw.Write(append(errorMsg, '\n'))
			rw.Flush()
			return
		}

		// Extract the request params
		link := request.Link
		depth := request.Depth

		// Run the crawler and mount the result
		controller := runner.TimeoutCrawl(link, depth)
		response := map[string]interface{}{
			"validLinks":   controller.ValidLinks,
			"invalidLinks": controller.InvalidLinks,
		}

		// Send response to the client
		responseJSON, _ := json.Marshal(response)
		rw.Write(append(responseJSON, '\n'))
		rw.Flush()
	}
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
	fmt.Println("Listening TCP server at", t.addr)
	defer t.Close()

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

		go t.handleConnection(conn)
	}
	return nil
}

func NewTCPServer(addr string) *TCPServer {
	server := TCPServer{
		addr: addr,
	}
	return &server
}
