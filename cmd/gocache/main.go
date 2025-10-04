package main

import (
	"fmt"
	handler "gocache/internal"
	"net"
)

func main() {
	listener, err := net.Listen("tcp", ":6379")
	if err != nil {
		fmt.Println(err)
		return
	}

	// Waits until a message is received. Then returns connection
	connection, err := listener.Accept()
	if err != nil {
		fmt.Println(err)
		return
	}

	defer connection.Close()

	handler.HandleConnection(connection)
}
