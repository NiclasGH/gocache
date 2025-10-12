package main

import (
	"fmt"
	"gocache/internal/handler"
	"net"
)

func main() {
	fmt.Println("Listening on port :6379")

	listener, err := net.Listen("tcp", ":6379")
	if err != nil {
		fmt.Println(err)
		return
	}

	for {
		// Waits until a message is received. Then returns connection
		connection, err := listener.Accept()
		if err != nil {
			fmt.Println(err)
			return
		}

		defer connection.Close()

		handler.HandleConnection(connection)
	}
}
