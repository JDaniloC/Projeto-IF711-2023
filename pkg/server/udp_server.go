package server

import (
	"encoding/json"
	"fmt"
	"net"

	utils "github.com/JDaniloC/Projeto-IF711-2023/internal/utils"
	runner "github.com/JDaniloC/Projeto-IF711-2023/pkg/runner"
)

type UDPServer struct {
	addr   string
	server *net.UDPConn
}

func (s *UDPServer) handleConnection(data []byte, addr *net.UDPAddr) {
	request := &utils.Request{}
	if err := json.Unmarshal(data, &request); err != nil {
		errorMsg, _ := json.Marshal(map[string]interface{}{
			"error": fmt.Sprintf("Erro ao desserializar para JSON: %s", err),
		})
		s.server.WriteToUDP(append(errorMsg, '\n'), addr)
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
	s.server.WriteToUDP(append(responseJSON, '\n'), addr)
}

func (s *UDPServer) Start() (err error) {
	addr, err := net.ResolveUDPAddr("udp", s.addr)
	if err != nil {
		println("could not resolve UDP addr")
		return err
	}

	s.server, err = net.ListenUDP("udp", addr)
	if err != nil {
		println("could not listen on UDP")
		return err
	}
	fmt.Println("Listening UDP server at", s.addr)

	buf := make([]byte, 2048)
	for {
		n, uAddr, err := s.server.ReadFromUDP(buf)
		if err != nil {
			println(err)
			return err
		}
		if uAddr == nil {
			continue
		}

		go s.handleConnection(buf[:n], uAddr)
	}
}

func NewUDPServer(addr string) *UDPServer {
	server := UDPServer{
		addr: addr,
	}
	return &server
}
