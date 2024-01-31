package server

import (
	"encoding/json"
	"fmt"
	"net"
	"testing"
)

const (
	address = ":6250"
)

func init() {
	server := NewUDPServer(address)

	go func() {
		server.Start()
	}()
}

func TestUDPServer_Connection(t *testing.T) {
	s, err := net.ResolveUDPAddr("udp4", address)
	if err != nil {
		t.Error("could not resolve address:", err)
	}
	conn, err := net.DialUDP("udp4", nil, s)
	if err != nil {
		t.Error("could not connect to server:", err)
	}
	conn.Close()
}

func TestUDPServer_Request(t *testing.T) {
	request := map[string]interface{}{
		"link":  "https://hackerspaces.org/",
		"depth": 3,
	}
	req, err := json.Marshal(request)
	if err != nil {
		t.Error("Error on serialize request", err)
	}
	addr, err := net.ResolveUDPAddr("udp4", address)
	if err != nil {
		t.Error("could not resolve address:", err)
	}

	t.Run("Send UDP request", func(t *testing.T) {
		conn, err := net.DialUDP("udp4", nil, addr)
		if err != nil {
			t.Error("could not connect to server: ", err)
		}
		defer conn.Close()

		_, err = conn.Write(append(req, '\n'))
		if err != nil {
			t.Error("could not write payload to server: ", err)
		}

		buffer := make([]byte, 2048)
		n, _, err := conn.ReadFromUDP(buffer)
		if err != nil {
			t.Error("could not read payload from server: ", err)
		}

		var response map[string]interface{}
		err = json.Unmarshal([]byte(buffer[:n]), &response)
		if err != nil {
			t.Error("Could not unmarshal response", err)
		}

		errMsg, hasError := response["error"].(string)
		if hasError {
			t.Error(errMsg)
		}

		fmt.Printf("Answer from server: %s\n", response)
	})
}
