package server

import (
	"fmt"
	"net"

	"github.com/codecrafters-io/kafka-starter-go/handler"
)

func StartServer() error {
	l, err := net.Listen("tcp", "0.0.0.0:9092")
	if err != nil {
		return fmt.Errorf("failed to bind to port 9092: %v", err)
	}
	defer l.Close()

	for {
		conn, err := l.Accept()
		if err != nil {
			return fmt.Errorf("error accepting connection: %v", err)
		}
		go handler.HandleRequest(conn)
	}
}