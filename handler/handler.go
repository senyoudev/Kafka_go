package handler

import (
	"fmt"
	"net"
	"github.com/codecrafters-io/kafka-starter-go/errors"
	"github.com/codecrafters-io/kafka-starter-go/protocol"
	"github.com/codecrafters-io/kafka-starter-go/response"
)

func HandleRequest(conn net.Conn) {
	defer conn.Close()

	buff := make([]byte, protocol.MaxPacketSize)
	if _, err := conn.Read(buff); err != nil {
		fmt.Printf("Error reading request: %v\n", err)
		return
	}

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

func handleAPIVersions(conn net.Conn, req protocol.Request) {
	resp := response.CreateAPIVersionsResponse(req.CorrelationID)
	response.SendResponse(conn, resp)
}