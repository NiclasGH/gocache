package command

import (
	"gocache/internal/resp"
	"testing"

	"gotest.tools/v3/assert"
)

func Test_ping(t *testing.T) {
	// given
	for k := range storage {
		delete(storage, k)
	}
	expected := resp.Value{Typ: "string", Str: "PONG"}

	ping, ok := Commands["PING"]
	if !ok {
		t.Error("Command does not exist")
		return
	}

	// when
	result := ping([]resp.Value{})

	// then
	assert.DeepEqual(t, expected, result)
}

func Test_pingWithArg(t *testing.T) {
	// given
	for k := range storage {
		delete(storage, k)
	}
	args := []resp.Value{
		{
			Typ:  resp.BULK.Typ,
			Bulk: "Tiramisu",
		},
	}

	ping, ok := Commands["PING"]
	if !ok {
		t.Error("Command does not exist")
		return
	}

	expected := resp.Value{Typ: "string", Str: "Tiramisu"}

	// when
	result := ping(args)

	// then
	assert.DeepEqual(t, expected, result)
}

func Test_set(t *testing.T) {
	// given
	for k := range storage {
		delete(storage, k)
	}
	args := []resp.Value{
		{
			Typ:  resp.BULK.Typ,
			Bulk: "Tira",
		},
		{
			Typ:  resp.BULK.Typ,
			Bulk: "Misu",
		},
	}

	expected := resp.Value{
		Typ: "string",
		Str: "OK",
	}

	set, ok := Commands["SET"]
	if !ok {
		t.Error("Command does not exist")
		return
	}

	// when
	result := set(args)

	// then
	assert.DeepEqual(t, expected, result)
	assert.Equal(t, "Misu", storage["Tira"])
}

func Test_setNeedsTwoArgs(t *testing.T) {
	// given
	for k := range storage {
		delete(storage, k)
	}

	args := []resp.Value{
		{
			Typ:  resp.BULK.Typ,
			Bulk: "Tira",
		},
	}

	set, ok := Commands["SET"]
	if !ok {
		t.Error("Command does not exist")
		return
	}

	// when
	result := set(args)

	// then
	assert.Equal(t, "error", result.Typ)
}

func Test_get(t *testing.T) {
	// given
	for k := range storage {
		delete(storage, k)
	}
	storage["Tira"] = "Misu"

	args := []resp.Value{
		{
			Typ:  resp.BULK.Typ,
			Bulk: "Tira",
		},
	}

	expected := resp.Value{
		Typ:  "bulk",
		Bulk: "Misu",
	}

	get, ok := Commands["GET"]
	if !ok {
		t.Error("Command does not exist")
		return
	}

	// when
	result := get(args)

	// then
	assert.DeepEqual(t, expected, result)
	assert.Equal(t, "Misu", storage["Tira"])
}

func Test_getCanOnlyReceiveOneArg(t *testing.T) {
	// given
	for k := range storage {
		delete(storage, k)
	}

	args := []resp.Value{
		{
			Typ:  resp.BULK.Typ,
			Bulk: "Tira",
		},
		{
			Typ:  resp.BULK.Typ,
			Bulk: "Tira",
		},
	}

	get, ok := Commands["GET"]
	if !ok {
		t.Error("Command does not exist")
		return
	}

	// when
	result := get(args)

	// then
	assert.Equal(t, "error", result.Typ)
}

func Test_getNoValueAvailable(t *testing.T) {
	// given
	for k := range storage {
		delete(storage, k)
	}

	args := []resp.Value{
		{
			Typ:  resp.BULK.Typ,
			Bulk: "Tira",
		},
	}

	expected := resp.Value{
		Typ: "null",
	}

	get, ok := Commands["GET"]
	if !ok {
		t.Error("Command does not exist")
		return
	}

	// when
	result := get(args)

	// then
	assert.DeepEqual(t, expected, result)
}
