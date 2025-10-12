package handler

import (
	"fmt"
	"gocache/internal/resp"
	"io"
	"net"
)

func HandleConnection(connection net.Conn) error {
	// the server allows long lived connections with many commands, until the client closes the connection
	for {
		reader := resp.NewReader(connection)
		value, err := reader.Read()
		if err != nil {
			if err == io.EOF {
				return nil
			}
			fmt.Println(err)
			return err
		}

		fmt.Println(value)

		writer := resp.NewWriter(connection)
		writer.Write(resp.Value{Typ: "string", Str: "OK"})
	}
}
