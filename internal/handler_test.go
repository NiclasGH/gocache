package handler

import (
	"net"
	"testing"
)

func HandlesConnection_respondsWithOk(t *testing.T) {
	// given
	client, server := net.Pipe()
	go HandleConnection(server)

	expectedResponse := "+OK\r\n"

	// when
	client.Write([]byte("PING\r\n"))

	// then
	buf := make([]byte, 1024)
	length, err := client.Read(buf)
	if err != nil {
		t.Error("The client was not able to read the handled connection response")
	}
	
	res := string(buf[:length])
	if res != expectedResponse {
		t.Errorf("Response [%s] didn't match expected [%s]", res, expectedResponse)
	}

	client.Close()
	server.Close()
}
