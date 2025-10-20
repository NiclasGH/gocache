package main

import (
	"net"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_pingIntegrationTest(t *testing.T) {
	// PREPARATIONS
	port := "8101"
	file, err := prepare(port)
	if err != nil {
		t.Error(err)
		return
	}
	defer os.Remove(file.Name())

	// STEP 1: Starting the server
	client, err := setupServer(port)
	if err != nil {
		t.Error(err)
		return
	}
	defer client.Close()

	// STEP 2: Call Ping
	err = func() error {
		// given
		request := "*2\r\n$4\r\nPING\r\n$8\r\nTiramisu\r\n"
		expected := "+Tiramisu\r\n"

		// when
		client.Write([]byte(request))

		// then
		buf := make([]byte, 1024)
		length, err := client.Read(buf)
		if err != nil {
			t.Error("The client was not able to read the handled connection response")
		}
		res := string(buf[:length])

		assert.Equal(t, expected, res)

		return nil
	}()
	if err != nil {
		t.Error(err)
		return
	}
}

func Test_setAndGet(t *testing.T) {
	// PREPARATIONS
	port := "8102"
	file, err := prepare(port)
	if err != nil {
		t.Error(err)
		return
	}
	defer os.Remove(file.Name())

	// STEP 1: Starting the server
	client, err := setupServer(port)
	if err != nil {
		t.Error(err)
		return
	}
	defer client.Close()

	// STEP 2: Call Set
	err = func() error {
		// given
		request := "*3\r\n$3\r\nSET\r\n$4\r\nTira\r\n$4\r\nMisu\r\n"
		expected := "+OK\r\n"

		// when
		client.Write([]byte(request))

		// then
		buf := make([]byte, 1024)
		length, err := client.Read(buf)
		if err != nil {
			t.Error("The client was not able to read the handled connection response")
		}
		res := string(buf[:length])

		assert.Equal(t, expected, res)

		return nil
	}()
	if err != nil {
		t.Error(err)
		return
	}

	// STEP 3: Call Get
	err = func() error {
		// given
		request := "*2\r\n$3\r\nGET\r\n$4\r\nTira\r\n"
		expected := "$4\r\nMisu\r\n"

		// when
		client.Write([]byte(request))

		// then
		buf := make([]byte, 1024)
		length, err := client.Read(buf)
		if err != nil {
			t.Error("The client was not able to read the handled connection response")
		}
		res := string(buf[:length])

		assert.Equal(t, expected, res)

		return nil
	}()
	if err != nil {
		t.Error(err)
		return
	}
}

// TODO enable test when startup use case exists
// func Test_getWithInitialization(t *testing.T) {
// 	// PREPARATIONS
// 	port := "8103"
// 	file, err := prepare(port)
// 	if err != nil {
// 		t.Error(err)
// 		return
// 	}
// 	defer os.Remove(file.Name())
//
// 	request := resp.Value{
// 		Typ: resp.ARRAY.Typ,
// 		Array: []resp.Value{
// 			{
// 				Typ:  resp.BULK.Typ,
// 				Bulk: "SET",
// 			},
// 			{
// 				Typ:  resp.BULK.Typ,
// 				Bulk: "Tira",
// 			},
// 			{
// 				Typ:  resp.BULK.Typ,
// 				Bulk: "Misu",
// 			},
// 		},
// 	}
// 	file.Write(request.Marshal())
//
// 	// STEP 1: Starting the server
// 	client, err := setupServer(port)
// 	if err != nil {
// 		t.Error(err)
// 		return
// 	}
// 	defer client.Close()
//
// 	// STEP 2: Call Get
// 	err = func() error {
// 		// given
// 		request := "*2\r\n$3\r\nGET\r\n$4\r\nTira\r\n"
// 		expected := "$4\r\nMisu\r\n"
//
// 		// when
// 		client.Write([]byte(request))
//
// 		// then
// 		buf := make([]byte, 1024)
// 		length, err := client.Read(buf)
// 		if err != nil {
// 			t.Error("The client was not able to read the handled connection response")
// 		}
// 		res := string(buf[:length])
//
// 		assert.Equal(t, expected, res)
//
// 		return nil
// 	}()
// 	if err != nil {
// 		t.Error(err)
// 		return
// 	}
// }

func prepare(port string) (*os.File, error) {
	f, err := os.CreateTemp("", "database.test.aof")
	if err != nil {
		return nil, err
	}
	os.Setenv("GC_DATABASE_PATH", f.Name())
	os.Setenv("GC_PORT", port)

	return f, nil
}

func setupServer(port string) (net.Conn, error) {
	// given
	port = ":" + port

	// when
	ready = make(chan struct{})
	go main()

	// then
	<-ready
	conn, err := net.Dial("tcp", port)
	if err != nil {
		return nil, err
	}

	return conn, nil
}
