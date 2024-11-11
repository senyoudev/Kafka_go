package main

import (
	"fmt"
	"os"
	"github.com/codecrafters-io/kafka-starter-go/server"
)

const (
	maxPacketSize     = 1024
	apiVersionsKey    = 18
	maxSupportedVersion = 4
)




func main() {
	if err:= server.StartServer(); err != nil {
		fmt.Println("Server Error: ", err)
		os.Exit(1)
	}
}


