package handler

import (
	"fmt"
	"io"
	"net"
	"os"
)

func HandleConnection(connection net.Conn) error {
	for {
		buf := make([]byte, 1024)

		_, err := connection.Read(buf)
		if err != nil {
			if err == io.EOF {
				break
			}
			fmt.Println("error reading from client: ", err.Error())
			os.Exit(1)
		}

		connection.Write([]byte("+OK\r\n"))
	}

	return nil
}
