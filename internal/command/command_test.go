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

	// when
	result := Ping([]resp.Value{})

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

	expected := resp.Value{Typ: "string", Str: "Tiramisu"}

	// when
	result := Ping(args)

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

	// when
	result := Set(args)

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

	// when
	result := Set(args)

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

	// when
	result := Get(args)

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

	// when
	result := Get(args)

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

	// when
	result := Get(args)

	// then
	assert.DeepEqual(t, expected, result)
}
