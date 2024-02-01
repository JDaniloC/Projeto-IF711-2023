package main

import (
	"encoding/json"
	"fmt"
	"net"
)

func main() {
	request := map[string]interface{}{
		"link":  "https://hackerspaces.org/",
		"depth": 2,
	}
	req, err := json.Marshal(request)
	if err != nil {
		println("Error on serialize request", err)
		return
	}
	addr, err := net.ResolveUDPAddr("udp4", ":6250")
	if err != nil {
		println("could not resolve address:", err)
		return
	}

	conn, err := net.DialUDP("udp4", nil, addr)
	if err != nil {
		println("could not connect to server: ", err)
	}
	defer conn.Close()

	_, err = conn.Write(append(req, '\n'))
	if err != nil {
		println("could not write payload to server: ", err)
	}

	buffer := make([]byte, 2048)
	n, _, err := conn.ReadFromUDP(buffer)
	if err != nil {
		println("could not read payload from server: ", err)
	}

	var response map[string]interface{}
	err = json.Unmarshal([]byte(buffer[:n]), &response)
	if err != nil {
		println("Could not unmarshal response", err)
	}

	errMsg, hasError := response["error"].(string)
	if hasError {
		println(errMsg)
	}

	vl := response["validLinks"]
	ivl := response["invalidLinks"]
	fmt.Printf("Valid links from server: %s\n", vl)
	fmt.Printf("Invalid links from server: %s\n", ivl)
}
