package server

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net"
	"testing"
)

func init() {
	tcp := NewTCPServer(":1123")

	go func() {
		tcp.Start()
	}()
}

func TestTCPServer_Connection(t *testing.T) {
	conn, err := net.Dial("tcp", ":1123")
	if err != nil {
		t.Error("could not connect to server: ", err)
	}
	conn.Close()
}

func TestTCPServer_Request(t *testing.T) {
	request := map[string]interface{}{
		"link":  "https://hackerspaces.org/",
		"depth": 2,
	}
	req, err := json.Marshal(request)
	if err != nil {
		t.Error("Error on serialize request", err)
	}

	t.Run("Send a simple request", func(t *testing.T) {
		conn, err := net.Dial("tcp", ":1123")
		if err != nil {
			t.Error("could not connect to server: ", err)
		}
		defer conn.Close()

		_, err = conn.Write(append(req, '\n'))
		if err != nil {
			t.Error("could not write payload to server: ", err)
		}

		reader := bufio.NewReader(conn)
		s, err := reader.ReadString('\n')
		if err != nil {
			t.Error("Error on read server answer:", err)
		}

		var response map[string]interface{}
		err = json.Unmarshal([]byte(s), &response)
		if err != nil {
			t.Error("Could not unmarshal response", err)
		}

		errMsg, hasError := response["error"].(string)
		if hasError {
			t.Error(errMsg)
		}

		fmt.Printf("Resposta recebida do servidor: %s\n", response)
	})
}
