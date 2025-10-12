package handler

import (
	"net"
	"testing"

	"gotest.tools/v3/assert"
)

func Test_handlesConnection_respondsWithOk(t *testing.T) {
	// given
	client, server := net.Pipe()
	go HandleConnection(server)

	expectedResponse := "+OK\r\n"

	// when
	client.Write([]byte("$4\r\nTira\r\n"))

	// then
	buf := make([]byte, 1024)
	length, err := client.Read(buf)
	if err != nil {
		t.Error("The client was not able to read the handled connection response")
	}

	res := string(buf[:length])
	assert.Equal(t, expectedResponse, res)

	client.Close()
	server.Close()
}
