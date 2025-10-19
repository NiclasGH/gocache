package infrastructure

import (
	"errors"
	"gocache/internal/core/resp"
	"net"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_handlesConnection_noArray_err(t *testing.T) {
	// given
	client, server := net.Pipe()
	defer client.Close()
	defer server.Close()

	testDb := TestDatabase{[]resp.Value{}}

	go HandleConnection(server, testDb)

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

	testDb := TestDatabase{[]resp.Value{}}

	go HandleConnection(server, testDb)

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

	testDb := TestDatabase{[]resp.Value{}}

	go HandleConnection(server, testDb)

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

	testDb := TestDatabase{[]resp.Value{}}

	go HandleConnection(server, testDb)

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
	assert.Equal(t, len(testDb.executedCommands), 0)
}

// Could be replaced with actual mocks
// TODO apparently I have no test for this? lol
type TestDatabase struct {
	executedCommands []resp.Value
}

func (db TestDatabase) GetInit() ([]resp.Value, error) {
	return nil, errors.New("Should never run this unmocked method Initialize()")
}
func (db TestDatabase) Close() error {
	return errors.New("Should never run this unmocked method Close()")
}
func (db TestDatabase) SaveSet(value resp.Value, _ string, _ string) error {
	db.executedCommands = append(db.executedCommands, value)
	return nil
}

func (db TestDatabase) GetSet(string) (string, error) {
	return "", errors.New("Should never run this unmocked method GetSet()")
}

func (db TestDatabase) SaveHSet(value resp.Value, _ string, _ string, _ string) error {
	db.executedCommands = append(db.executedCommands, value)
	return nil
}

func (db TestDatabase) GetHSet(string) (map[string]string, error) {
	return nil, errors.New("Should never run this unmocked method GetHSet()")
}
