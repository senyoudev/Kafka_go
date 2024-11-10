package main

import (
	"encoding/binary"
	"fmt"
	"net"
	"os"

	"github.com/codecrafters-io/kafka-starter-go/errors"
)


func main() {
	l, err := net.Listen("tcp", "0.0.0.0:9092")
	if err != nil {
		fmt.Println("Failed to bind to port 9092")
		os.Exit(1)
	}
	defer l.Close()
	// Close the listener when the application closes
	for {
		conn, err := l.Accept()
		if err != nil {
			fmt.Println("Error accepting connection: ", err.Error())
			os.Exit(1)
		}
		go handleRequest(conn)
	}
}

func handleRequest(conn net.Conn) {
	defer conn.Close()
	const maxPacketSize = 1024
	buff := make([]byte, maxPacketSize)
	_, err:= conn.Read(buff)
	if err != nil {
		fmt.Println("Error writing response:", err)
		return
	}

	// Read the api version
    api_version := binary.BigEndian.Uint16(buff[6:8])
	response := make([]byte, 10)
	
	if api_version <= 4 {
    	fmt.Println("Api version is within the range 0 to 4:", api_version)
	} else {
		fmt.Println(errors.UNSUPPORTED_VERSION.Message())
		errorBytes := errors.ErrorCodeToBytes(errors.UNSUPPORTED_VERSION)
		copy(response[8:], errorBytes)
	}
	

	binary.BigEndian.PutUint32(response[:4], 0)
	// Extract correlation_id (4 bytes starting from the 9th byte)
	correlationID := binary.BigEndian.Uint32(buff[8:12])
	binary.BigEndian.PutUint32(response[4:], correlationID)
	_, error := conn.Write(response)
	if error != nil {
		fmt.Println("Error writing response:", err)
		return
	}
}
