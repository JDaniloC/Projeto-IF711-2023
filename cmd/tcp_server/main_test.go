package main

import (
	"bufio"
	"encoding/json"
	"net"
	"testing"

	tcp "github.com/JDaniloC/Projeto-IF711-2023/pkg/server"
)

func init() {
	tcpServer := tcp.NewTCPServer(":1123")
	go tcpServer.Start()
}

func BenchmarkTCPServer(b *testing.B) {
	request := map[string]interface{}{
		"link":  "https://hackerspaces.org/",
		"depth": 1,
	}
	req, err := json.Marshal(request)
	if err != nil {
		b.Error("Error on serialize request", err)
	}

	conn, err := net.Dial("tcp", ":1123")
	if err != nil {
		b.Error("could not connect to server: ", err)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err = conn.Write(append(req, '\n'))
		if err != nil {
			b.Error("could not write payload to server: ", err)
		}

		reader := bufio.NewReader(conn)
		s, err := reader.ReadString('\n')
		if err != nil {
			b.Error("Error on read server answer:", err)
		}

		var response map[string]interface{}
		err = json.Unmarshal([]byte(s), &response)
		if err != nil {
			b.Error("Could not unmarshal response", err)
		}

		errMsg, hasError := response["error"].(string)
		if hasError {
			b.Error(errMsg)
		}
	}
	conn.Close()
}
