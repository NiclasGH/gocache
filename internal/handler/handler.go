package handler

import (
	"fmt"
	"gocache/internal/resp"
	"net"
)

func HandleConnection(connection net.Conn) error {
	// the server allows long lived connections with many commands, until the client closes the connection
	for {
		resp := resp.NewReader(connection)
		value, err := resp.Read()
		if err != nil {
			fmt.Println(err)
			return err
		}

		fmt.Println(value)

		// ignore request and send back a PONG
		connection.Write([]byte("+OK\r\n"))
	}
}
