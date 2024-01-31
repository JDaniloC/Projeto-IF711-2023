package server

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net"
	"testing"
)

const (
	addr = ":1123"
)

func init() {
	tcp := NewTCPServer(addr)

	go func() {
		tcp.Start()
	}()
}

func TestTCPServer_Connection(t *testing.T) {
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		t.Error("could not connect to server: ", err)
	}
	defer conn.Close()
}

func TestTCPServer_Request(t *testing.T) {
	request := map[string]interface{}{
		"link":  "https://hackerspaces.org/",
		"depth": 1,
	}
	json, err := json.Marshal(request)
	if err != nil {
		t.Error("Error on serialize request", err)
	}

	conn, err := net.Dial("tcp", addr)
	if err != nil {
		t.Error("could not connect to server: ", err)
	}
	defer conn.Close()

	t.Run("Send a simple request", func(t *testing.T) {
		_, err = conn.Write(append(json, '\n'))
		if err != nil {
			t.Error("could not write payload to server: ", err)
		}

		response, err := bufio.NewReader(conn).ReadString('\n')
		if err != nil {
			t.Error("Erro ao ler resposta do servidor:", err)
			return
		}

		fmt.Println("Resposta recebida do servidor: ", response)
	})
}
