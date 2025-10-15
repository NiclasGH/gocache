package handler

import (
	"gocache/internal/resp"
	"net"
	"strings"
	"testing"

	"gotest.tools/v3/assert"
)

func Test_handlesConnection_noArray_err(t *testing.T) {
	// given
	client, server := net.Pipe()
	defer client.Close()
	defer server.Close()

	go HandleConnection(server)

	// when
	client.Write([]byte("$4\r\nTira\r\n"))

	// then
	buf := make([]byte, 1024)
	length, err := client.Read(buf)
	if err != nil {
		t.Error("The client was not able to read the handled connection response")
	}

	res := string(buf[:length])
	if !strings.Contains(res, string(resp.ERROR.RespCode)) {
		t.Error("Response didn't contain an error")
		return
	}
	if !strings.Contains(res, "array") {
		t.Error("Response didn't contain the correct error")
		return
	}
}

func Test_handlesConnection_noBulkInArray_err(t *testing.T) {
	// given
	client, server := net.Pipe()
	defer client.Close()
	defer server.Close()

	go HandleConnection(server)

	// when
	client.Write([]byte("*1\r\n*1\r\n$4\r\nTira\r\n"))

	// then
	buf := make([]byte, 1024)
	length, err := client.Read(buf)
	if err != nil {
		t.Error("The client was not able to read the handled connection response")
	}

	res := string(buf[:length])
	if !strings.Contains(res, string(resp.ERROR.RespCode)) {
		t.Error("Response didn't contain an error")
		return
	}
	if !strings.Contains(res, "Unable to read command") {
		t.Error("Response didn't contain the correct error")
		return
	}
}

func Test_handlesConnection_unknownCommand_err(t *testing.T) {
	// given
	client, server := net.Pipe()
	defer client.Close()
	defer server.Close()

	go HandleConnection(server)

	// when
	client.Write([]byte("*1\r\n$7\r\nUNKNOWN\r\n"))

	// then
	buf := make([]byte, 1024)
	length, err := client.Read(buf)
	if err != nil {
		t.Error("The client was not able to read the handled connection response")
	}

	res := string(buf[:length])
	if !strings.Contains(res, string(resp.ERROR.RespCode)) {
		t.Error("Response didn't contain an error")
		return
	}
	if !strings.Contains(res, "Command is unknown") {
		t.Error("Response didn't contain the correct error")
		return
	}
}

func Test_handlesConnection_ping(t *testing.T) {
	// given
	client, server := net.Pipe()
	defer client.Close()
	defer server.Close()

	go HandleConnection(server)

	expectedResponse := "+PONG\r\n"

	// when
	client.Write([]byte("*1\r\n$4\r\nPING\r\n"))

	// then
	buf := make([]byte, 1024)
	length, err := client.Read(buf)
	if err != nil {
		t.Error("The client was not able to read the handled connection response")
	}
	res := string(buf[:length])

	assert.Equal(t, expectedResponse, res)
}
