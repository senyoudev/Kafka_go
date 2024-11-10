package main

import (
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
	header := make([]byte, 1024)
	conn.Read(header)
	conn.Write([]byte{0, 0, 0, 0, 0, 0, 0, 7})
}
