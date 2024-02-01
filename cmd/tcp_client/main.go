package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net"
)

const (
	addr = ":1123"
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

	conn, err := net.Dial("tcp", ":1123")
	if err != nil {
		println("could not connect to server: ", err)
		return
	}
	defer conn.Close()

	_, err = conn.Write(append(req, '\n'))
	if err != nil {
		println("could not write payload to server: ", err)
		return
	}

	reader := bufio.NewReader(conn)
	s, err := reader.ReadString('\n')
	if err != nil {
		println("Error on read server answer:", err)
		return
	}

	var response map[string]interface{}
	err = json.Unmarshal([]byte(s), &response)
	if err != nil {
		println("Could not unmarshal response", err)
		return
	}

	errMsg, hasError := response["error"].(string)
	if hasError {
		println(errMsg)
		return
	}

	vl := response["validLinks"]
	ivl := response["invalidLinks"]
	fmt.Printf("Valid links from server: %s\n", vl)
	fmt.Printf("Invalid links from server: %s\n", ivl)
}
