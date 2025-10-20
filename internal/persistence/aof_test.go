package persistence

import (
	"io"
	"os"
	"testing"

	"gocache/internal/core/resp"

	"github.com/stretchr/testify/assert"
)

func Test_savePersistsCommand(t *testing.T) {
	// given
	request := resp.Value{
		Typ: resp.ARRAY.Typ,
		Array: []resp.Value{
			{
				Typ:  resp.BULK.Typ,
				Bulk: "SET",
			},
			{
				Typ:  resp.BULK.Typ,
				Bulk: "Tira",
			},
			{
				Typ:  resp.BULK.Typ,
				Bulk: "Misu",
			},
		},
	}

	expected := "*3\r\n$3\r\nSET\r\n$4\r\nTira\r\n$4\r\nMisu\r\n"

	file, err := os.CreateTemp("", "database.test.aof")
	if err != nil {
		t.Error(err)
		return
	}
	defer os.Remove(file.Name())

	aof, err := NewAof(file.Name())
	if err != nil {
		t.Error(err)
		return
	}

	// when
	err = aof.Save(request)
	if err != nil {
		t.Error(err)
		return
	}

	// then
	result, err := io.ReadAll(file)
	if err != nil {
		t.Error(err)
		return
	}

	assert.Equal(t, expected, string(result))
}

func Test_initializeReturnsCommands(t *testing.T) {
	// given
	request := resp.Value{
		Typ: resp.ARRAY.Typ,
		Array: []resp.Value{
			{
				Typ:  resp.BULK.Typ,
				Bulk: "SET",
			},
			{
				Typ:  resp.BULK.Typ,
				Bulk: "Tira",
			},
			{
				Typ:  resp.BULK.Typ,
				Bulk: "Misu",
			},
		},
	}

	file, err := os.CreateTemp("", "database.test.aof")
	if err != nil {
		t.Error(err)
		return
	}
	defer os.Remove(file.Name())
	file.Write(request.Marshal())

	aof, err := NewAof(file.Name())
	if err != nil {
		t.Error(err)
		return
	}

	// when
	result, err := aof.ReadPersistedCommands()
	if err != nil {
		t.Error(err)
		return
	}

	// then
	assert.ElementsMatch(t, []resp.Value{request}, result)
}
