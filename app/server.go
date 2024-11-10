package main

import (
	"encoding/binary"
	"fmt"
	"net"
	"os"
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

	response := make([]byte, 8)
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
