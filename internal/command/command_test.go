package command

import (
	"gocache/internal/resp"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_ping(t *testing.T) {
	// given
	for k := range setStorage {
		delete(setStorage, k)
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
	assert.EqualValues(t, expected, result)
}

func Test_pingWithArg(t *testing.T) {
	// given
	for k := range setStorage {
		delete(setStorage, k)
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
	assert.EqualValues(t, expected, result)
}

func Test_set(t *testing.T) {
	// given
	for k := range setStorage {
		delete(setStorage, k)
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
	assert.EqualValues(t, expected, result)

	value, ok := setStorage["Tira"]
	if !ok {
		t.Error("Set Storage did not contain key 'Tira'")
		return
	}
	assert.Equal(t, "Misu", value)
}

func Test_setNeedsTwoArgs(t *testing.T) {
	// given
	for k := range setStorage {
		delete(setStorage, k)
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
	for k := range setStorage {
		delete(setStorage, k)
	}
	setStorage["Tira"] = "Misu"

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
	assert.EqualValues(t, expected, result)

	value, ok := setStorage["Tira"]
	if !ok {
		t.Error("Set Storage did not contain key 'Tira'")
		return
	}
	assert.Equal(t, "Misu", value)
}

func Test_getCanOnlyReceiveOneArg(t *testing.T) {
	// given
	for k := range setStorage {
		delete(setStorage, k)
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
	for k := range setStorage {
		delete(setStorage, k)
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
	assert.EqualValues(t, expected, result)
}

func Test_hset(t *testing.T) {
	// given
	for k := range hsetStorage {
		delete(hsetStorage, k)
	}
	args := []resp.Value{
		{
			Typ:  resp.BULK.Typ,
			Bulk: "tira",
		},
		{
			Typ:  resp.BULK.Typ,
			Bulk: "misu",
		},
		{
			Typ:  resp.BULK.Typ,
			Bulk: "cute",
		},
	}

	expected := resp.Value{
		Typ: "string",
		Str: "OK",
	}

	hset, ok := Commands["HSET"]
	if !ok {
		t.Error("Command does not exist")
		return
	}

	// when
	result := hset(args)

	// then
	assert.EqualValues(t, expected, result)

	value, ok := hsetStorage["tira"]["misu"]
	if !ok {
		t.Error("HSet Storage did not contain hash 'tira' or key 'misu'")
		return
	}
	assert.Equal(t, "cute", value)
}

func Test_hsetNeedsThreeArgs(t *testing.T) {
	// given
	for k := range hsetStorage {
		delete(hsetStorage, k)
	}

	args := []resp.Value{
		{
			Typ:  resp.BULK.Typ,
			Bulk: "tira",
		},
	}

	hset, ok := Commands["HSET"]
	if !ok {
		t.Error("Command does not exist")
		return
	}

	// when
	result := hset(args)

	// then
	assert.Equal(t, "error", result.Typ)
}

func Test_hget(t *testing.T) {
	// given
	for k := range hsetStorage {
		delete(hsetStorage, k)
	}

	hsetStorage["tira"] = map[string]string{}
	hsetStorage["tira"]["misu"] = "cute"

	args := []resp.Value{
		{
			Typ:  resp.BULK.Typ,
			Bulk: "tira",
		},
		{
			Typ:  resp.BULK.Typ,
			Bulk: "misu",
		},
	}

	expected := resp.Value{
		Typ:  "bulk",
		Bulk: "cute",
	}

	hget, ok := Commands["HGET"]
	if !ok {
		t.Error("Command does not exist")
		return
	}

	// when
	result := hget(args)

	// then
	assert.EqualValues(t, expected, result)

	value, ok := hsetStorage["tira"]["misu"]
	if !ok {
		t.Error("HSet Storage did not contain hash 'tira' or key 'misu'")
		return
	}
	assert.Equal(t, "cute", value)
}

func Test_hgetCanOnlyReceiveTwoArg(t *testing.T) {
	// given
	for k := range hsetStorage {
		delete(hsetStorage, k)
	}

	args := []resp.Value{
		{
			Typ:  resp.BULK.Typ,
			Bulk: "tira",
		},
	}

	hget, ok := Commands["HGET"]
	if !ok {
		t.Error("Command does not exist")
		return
	}

	// when
	result := hget(args)

	// then
	assert.Equal(t, "error", result.Typ)
}

func Test_hgetNoValueAvailable(t *testing.T) {
	// given
	for k := range hsetStorage {
		delete(hsetStorage, k)
	}

	args := []resp.Value{
		{
			Typ:  resp.BULK.Typ,
			Bulk: "tira",
		},
		{
			Typ:  resp.BULK.Typ,
			Bulk: "misu",
		},
	}

	expected := resp.Value{
		Typ: "null",
	}

	hget, ok := Commands["HGET"]
	if !ok {
		t.Error("Command does not exist")
		return
	}

	// when
	result := hget(args)

	// then
	assert.EqualValues(t, expected, result)
}

func Test_command_returnsSpecs(t *testing.T) {
	// given
	expected := []resp.Value{
		{
			Typ: resp.ARRAY.Typ,
			Array: []resp.Value{
				// 1. command
				{Typ: resp.BULK.Typ, Bulk: "PING"},
				// 2. arg count
				{Typ: resp.INTEGER.Typ, Num: -1},
				// 3. flags
				{
					Typ: resp.ARRAY.Typ,
					Array: []resp.Value{
						{Typ: resp.BULK.Typ, Bulk: "readonly"},
						{Typ: resp.BULK.Typ, Bulk: "fast"},
					},
				},
				// 4. first key
				{Typ: resp.INTEGER.Typ, Num: 1},
				// 5. last key
				{Typ: resp.INTEGER.Typ, Num: 1},
				// 6. steps between keys
				{Typ: resp.INTEGER.Typ, Num: 1},
				// 7. ACL flags
				{
					Typ: resp.ARRAY.Typ,
					Array: []resp.Value{
						{Typ: resp.BULK.Typ, Bulk: "@connection"},
						{Typ: resp.BULK.Typ, Bulk: "@fast"},
					},
				},
			},
		},
		{
			Typ: resp.ARRAY.Typ,
			Array: []resp.Value{
				// 1. command
				{Typ: resp.BULK.Typ, Bulk: "GET"},
				// 2. arg count
				{Typ: resp.INTEGER.Typ, Num: 2},
				// 3. flags
				{
					Typ: resp.ARRAY.Typ,
					Array: []resp.Value{
						{Typ: resp.BULK.Typ, Bulk: "readonly"},
						{Typ: resp.BULK.Typ, Bulk: "fast"},
					},
				},
				// 4. first key
				{Typ: resp.INTEGER.Typ, Num: 1},
				// 5. last key
				{Typ: resp.INTEGER.Typ, Num: 1},
				// 6. steps between keys
				{Typ: resp.INTEGER.Typ, Num: 1},
				// 7. ACL flags
				{
					Typ: resp.ARRAY.Typ,
					Array: []resp.Value{
						{Typ: resp.BULK.Typ, Bulk: "@read"},
						{Typ: resp.BULK.Typ, Bulk: "@fast"},
						{Typ: resp.BULK.Typ, Bulk: "@string"},
					},
				},
			},
		},
		{
			Typ: resp.ARRAY.Typ,
			Array: []resp.Value{
				// 1. command
				{Typ: resp.BULK.Typ, Bulk: "SET"},
				// 2. arg count
				{Typ: resp.INTEGER.Typ, Num: 3},
				// 3. flags
				{
					Typ: resp.ARRAY.Typ,
					Array: []resp.Value{
						{Typ: resp.BULK.Typ, Bulk: "write"},
						{Typ: resp.BULK.Typ, Bulk: "fast"},
					},
				},
				// 4. first key
				{Typ: resp.INTEGER.Typ, Num: 1},
				// 5. last key
				{Typ: resp.INTEGER.Typ, Num: 2},
				// 6. steps between keys
				{Typ: resp.INTEGER.Typ, Num: 1},
				// 7. ACL flags
				{
					Typ: resp.ARRAY.Typ,
					Array: []resp.Value{
						{Typ: resp.BULK.Typ, Bulk: "@write"},
						{Typ: resp.BULK.Typ, Bulk: "@slow"},
						{Typ: resp.BULK.Typ, Bulk: "@string"},
					},
				},
			},
		},
		{
			Typ: resp.ARRAY.Typ,
			Array: []resp.Value{
				// 1. command
				{Typ: resp.BULK.Typ, Bulk: "HGET"},
				// 2. arg count
				{Typ: resp.INTEGER.Typ, Num: 3},
				// 3. flags
				{
					Typ: resp.ARRAY.Typ,
					Array: []resp.Value{
						{Typ: resp.BULK.Typ, Bulk: "readonly"},
						{Typ: resp.BULK.Typ, Bulk: "fast"},
					},
				},
				// 4. first key
				{Typ: resp.INTEGER.Typ, Num: 1},
				// 5. last key
				{Typ: resp.INTEGER.Typ, Num: 2},
				// 6. steps between keys
				{Typ: resp.INTEGER.Typ, Num: 1},
				// 7. ACL flags
				{
					Typ: resp.ARRAY.Typ,
					Array: []resp.Value{
						{Typ: resp.BULK.Typ, Bulk: "@read"},
						{Typ: resp.BULK.Typ, Bulk: "@hash"},
						{Typ: resp.BULK.Typ, Bulk: "@fast"},
					},
				},
			},
		},
		{
			Typ: resp.ARRAY.Typ,
			Array: []resp.Value{
				// 1. command
				{Typ: resp.BULK.Typ, Bulk: "HSET"},
				// 2. arg count
				{Typ: resp.INTEGER.Typ, Num: 4},
				// 3. flags
				{
					Typ: resp.ARRAY.Typ,
					Array: []resp.Value{
						{Typ: resp.BULK.Typ, Bulk: "write"},
						{Typ: resp.BULK.Typ, Bulk: "fast"},
					},
				},
				// 4. first key
				{Typ: resp.INTEGER.Typ, Num: 1},
				// 5. last key
				{Typ: resp.INTEGER.Typ, Num: 3},
				// 6. steps between keys
				{Typ: resp.INTEGER.Typ, Num: 1},
				// 7. ACL flags
				{
					Typ: resp.ARRAY.Typ,
					Array: []resp.Value{
						{Typ: resp.BULK.Typ, Bulk: "@write"},
						{Typ: resp.BULK.Typ, Bulk: "@hash"},
						{Typ: resp.BULK.Typ, Bulk: "@fast"},
					},
				},
			},
		},
		{
			Typ: resp.ARRAY.Typ,
			Array: []resp.Value{
				// 1. command
				{Typ: resp.BULK.Typ, Bulk: "HGETALL"},
				// 2. arg count
				{Typ: resp.INTEGER.Typ, Num: 2},
				// 3. flags
				{
					Typ: resp.ARRAY.Typ,
					Array: []resp.Value{
						{Typ: resp.BULK.Typ, Bulk: "readonly"},
					},
				},
				// 4. first key
				{Typ: resp.INTEGER.Typ, Num: 1},
				// 5. last key
				{Typ: resp.INTEGER.Typ, Num: 1},
				// 6. steps between keys
				{Typ: resp.INTEGER.Typ, Num: 1},
				// 7. ACL flags
				{
					Typ: resp.ARRAY.Typ,
					Array: []resp.Value{
						{Typ: resp.BULK.Typ, Bulk: "@read"},
						{Typ: resp.BULK.Typ, Bulk: "@hash"},
						{Typ: resp.BULK.Typ, Bulk: "@slow"},
					},
				},
			},
		},
	}

	commandSpecs, ok := Commands["COMMAND"]
	if !ok {
		t.Error("Command does not exist")
		return
	}

	// when
	result := commandSpecs([]resp.Value{})

	// then
	assert.ElementsMatch(t, result.Array, expected)
}

func Test_command_withFilter_caseInsensitive_returnsSpecOfFilter(t *testing.T) {
	// given
	args := []resp.Value{
		{
			Typ:  resp.BULK.Typ,
			Bulk: "PiNg",
		},
	}

	expected := []resp.Value{
		{
			Typ: resp.ARRAY.Typ,
			Array: []resp.Value{
				// 1. command
				{Typ: resp.BULK.Typ, Bulk: "PING"},
				// 2. arg count
				{Typ: resp.INTEGER.Typ, Num: -1},
				// 3. flags
				{
					Typ: resp.ARRAY.Typ,
					Array: []resp.Value{
						{Typ: resp.BULK.Typ, Bulk: "readonly"},
						{Typ: resp.BULK.Typ, Bulk: "fast"},
					},
				},
				// 4. first key
				{Typ: resp.INTEGER.Typ, Num: 1},
				// 5. last key
				{Typ: resp.INTEGER.Typ, Num: 1},
				// 6. steps between keys
				{Typ: resp.INTEGER.Typ, Num: 1},
				// 7. ACL flags
				{
					Typ: resp.ARRAY.Typ,
					Array: []resp.Value{
						{Typ: resp.BULK.Typ, Bulk: "@connection"},
						{Typ: resp.BULK.Typ, Bulk: "@fast"},
					},
				},
			},
		},
	}

	commandSpecs, ok := Commands["COMMAND"]
	if !ok {
		t.Error("Command does not exist")
		return
	}

	// when
	result := commandSpecs(args)

	// then
	assert.ElementsMatch(t, result.Array, expected)
}

// func Test_commandDocs_returnsDocs(t *testing.T) {
// 	// given
// 	args := []resp.Value{
// 		{
// 			Typ:  resp.BULK.Typ,
// 			Bulk: "tira",
// 		},
// 		{
// 			Typ:  resp.BULK.Typ,
// 			Bulk: "misu",
// 		},
// 	}
//
// 	expected := resp.Value{
// 		Typ: "null",
// 	}
//
// 	hget, ok := Commands["HGET"]
// 	if !ok {
// 		t.Error("Command does not exist")
// 		return
// 	}
//
// 	// when
// 	result := hget(args)
//
// 	// then
// 	assert.DeepEqual(t, expected, result)
// }
