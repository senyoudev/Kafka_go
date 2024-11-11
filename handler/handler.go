package handler

import (
	"fmt"
	"net"
	"time"

	"github.com/codecrafters-io/kafka-starter-go/errors"
	"github.com/codecrafters-io/kafka-starter-go/protocol"
	"github.com/codecrafters-io/kafka-starter-go/response"
)

func HandleRequest(conn net.Conn) {
	defer conn.Close()

	// Send small packets to detect if the other part is still reachable
	if tcpConn, ok := conn.(*net.TCPConn); ok {
        tcpConn.SetKeepAlive(true)
        tcpConn.SetKeepAlivePeriod(30 * time.Second)
    }

	// Since the client is not disconnected, we can process serial requests
	for {
		buff := make([]byte, protocol.MaxPacketSize)
		if _, err := conn.Read(buff); err != nil {
			fmt.Printf("Error reading request: %v\n", err)
			return
		}
		// Set Read Deadline to detect inactive users
		conn.SetReadDeadline(time.Now().Add(5 * time.Second))

		req := protocol.ParseRequest(buff)
		
		if !protocol.IsVersionSupported(req.APIVersion) {
			fmt.Println("Unsupported version detected:", req.APIVersion)
			response.SendErrorResponse(conn, req.CorrelationID, errors.UNSUPPORTED_VERSION)
			return
		}

		switch req.APIKey {
		case protocol.APIVersionsKey:
			handleAPIVersions(conn, req)
		default:
			response.SendErrorResponse(conn, req.CorrelationID, nil)
		}
	}


}

func handleAPIVersions(conn net.Conn, req protocol.Request) {
	resp := response.CreateAPIVersionsResponse(req.CorrelationID)
	response.SendResponse(conn, resp)
}