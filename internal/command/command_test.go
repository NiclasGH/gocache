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

	ping, ok := Commands[PING]
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

	ping, ok := Commands[PING]
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

	set, ok := Commands[SET]
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

	set, ok := Commands[SET]
	if !ok {
		t.Error("Command does not exist")
		return
	}

	// when
	result := set(args)

	// then
	assert.Equal(t, "error", result.Typ)
}

func Test_del(t *testing.T) {
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
		Typ: "integer",
		Num: 1,
	}

	del, ok := Commands[DEL]
	if !ok {
		t.Error("Command does not exist")
		return
	}

	// when
	result := del(args)

	// then
	assert.EqualValues(t, expected, result)

	_, ok = setStorage["Tira"]
	if ok {
		t.Error("Set Storage Key 'Tira' was not deleted")
		return
	}
}

func Test_del_multipleKeys(t *testing.T) {
	// given
	for k := range setStorage {
		delete(setStorage, k)
	}
	setStorage["Tira"] = "Misu"
	setStorage["Misu"] = "Tira"

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
		Typ: "integer",
		Num: 2,
	}

	del, ok := Commands[DEL]
	if !ok {
		t.Error("Command does not exist")
		return
	}

	// when
	result := del(args)

	// then
	assert.EqualValues(t, expected, result)

	_, ok = setStorage["Tira"]
	if ok {
		t.Error("Set Storage Key 'Tira' was not deleted")
		return
	}

	_, ok = setStorage["Misu"]
	if ok {
		t.Error("Set Storage Key 'Misu' was not deleted")
		return
	}
}

func Test_del_needsAtLeastOneKey(t *testing.T) {
	// given
	args := []resp.Value{}

	del, ok := Commands[DEL]
	if !ok {
		t.Error("Command does not exist")
		return
	}

	// when
	result := del(args)

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

	get, ok := Commands[GET]
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

	get, ok := Commands[GET]
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

	get, ok := Commands[GET]
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

	hset, ok := Commands[HSET]
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

	hset, ok := Commands[HSET]
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

	hget, ok := Commands[HGET]
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

	hget, ok := Commands[HGET]
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

	hget, ok := Commands[HGET]
	if !ok {
		t.Error("Command does not exist")
		return
	}

	// when
	result := hget(args)

	// then
	assert.EqualValues(t, expected, result)
}

func Test_hdel(t *testing.T) {
	// given
	for k := range hsetStorage {
		delete(hsetStorage, k)
	}
	hsetStorage["tira"] = map[string]string{}
	hsetStorage["tira"]["misu"] = "cute"
	hsetStorage["tira"]["void"] = "scary"

	args := []resp.Value{
		// hash
		{
			Typ:  resp.BULK.Typ,
			Bulk: "tira",
		},
		// field
		{
			Typ:  resp.BULK.Typ,
			Bulk: "misu",
		},
	}

	expected := resp.Value{
		Typ: "integer",
		Num: 1,
	}

	hdel, ok := Commands[HDEL]
	if !ok {
		t.Error("Command does not exist")
		return
	}

	// when
	result := hdel(args)

	// then
	assert.EqualValues(t, expected, result)

	_, ok = hsetStorage["tira"]["misu"]
	if ok {
		t.Error("HSet Storage did not get key 'misu' deleted")
		return
	}

	value, ok := hsetStorage["tira"]["void"]
	if !ok {
		t.Error("HSet Storage did get key 'void' deleted but wasn't supposed to")
		return
	}
	assert.Equal(t, value, "scary")
}

func Test_hdel_lastKeyDeletesHash(t *testing.T) {
	// given
	for k := range hsetStorage {
		delete(hsetStorage, k)
	}
	hsetStorage["tira"] = map[string]string{}
	hsetStorage["tira"]["misu"] = "cute"

	args := []resp.Value{
		// hash
		{
			Typ:  resp.BULK.Typ,
			Bulk: "tira",
		},
		// field
		{
			Typ:  resp.BULK.Typ,
			Bulk: "misu",
		},
	}

	expected := resp.Value{
		Typ: "integer",
		Num: 1,
	}

	hdel, ok := Commands[HDEL]
	if !ok {
		t.Error("Command does not exist")
		return
	}

	// when
	result := hdel(args)

	// then
	assert.EqualValues(t, expected, result)

	_, ok = hsetStorage["tira"]
	if ok {
		t.Error("HSet Storage did not get hash 'tira' deleted")
		return
	}
}

func Test_hdel_deleteMultipleFields(t *testing.T) {
	// given
	for k := range hsetStorage {
		delete(hsetStorage, k)
	}
	hsetStorage["tira"] = map[string]string{}
	hsetStorage["tira"]["misu"] = "cute"
	hsetStorage["tira"]["void"] = "scary"
	hsetStorage["tira"]["isSpider"] = "yes"

	args := []resp.Value{
		// hash
		{
			Typ:  resp.BULK.Typ,
			Bulk: "tira",
		},
		// field
		{
			Typ:  resp.BULK.Typ,
			Bulk: "misu",
		},
		// field
		{
			Typ:  resp.BULK.Typ,
			Bulk: "void",
		},
	}

	expected := resp.Value{
		Typ: "integer",
		Num: 2,
	}

	hdel, ok := Commands[HDEL]
	if !ok {
		t.Error("Command does not exist")
		return
	}

	// when
	result := hdel(args)

	// then
	assert.EqualValues(t, expected, result)

	_, ok = hsetStorage["tira"]["misu"]
	if ok {
		t.Error("HSet Storage did not get key 'misu' deleted")
		return
	}

	_, ok = hsetStorage["tira"]["void"]
	if ok {
		t.Error("HSet Storage did not get key 'void' deleted")
		return
	}

	value, ok := hsetStorage["tira"]["isSpider"]
	if !ok {
		t.Error("HSet Storage did get key 'isSpider' deleted but wasn't supposed to")
		return
	}
	assert.Equal(t, value, "yes")
}

func Test_hdel_needsAtLeastTwoArgs(t *testing.T) {
	// given
	args := []resp.Value{
		// hash
		{
			Typ:  resp.BULK.Typ,
			Bulk: "tira",
		},
	}

	hdel, ok := Commands[HDEL]
	if !ok {
		t.Error("Command does not exist")
		return
	}

	// when
	result := hdel(args)

	// then
	assert.Equal(t, "error", result.Typ)
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
				{Typ: resp.BULK.Typ, Bulk: "DEL"},
				// 2. arg count
				{Typ: resp.INTEGER.Typ, Num: -2},
				// 3. flags
				{
					Typ: resp.ARRAY.Typ,
					Array: []resp.Value{
						{Typ: resp.BULK.Typ, Bulk: "write"},
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
						{Typ: resp.BULK.Typ, Bulk: "@write"},
						{Typ: resp.BULK.Typ, Bulk: "@slow"},
						{Typ: resp.BULK.Typ, Bulk: "@keyspace"},
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
				{Typ: resp.BULK.Typ, Bulk: "HDEL"},
				// 2. arg count
				{Typ: resp.INTEGER.Typ, Num: -3},
				// 3. flags
				{
					Typ: resp.ARRAY.Typ,
					Array: []resp.Value{
						{Typ: resp.BULK.Typ, Bulk: "write"},
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
						{Typ: resp.BULK.Typ, Bulk: "@fast"},
						{Typ: resp.BULK.Typ, Bulk: "@hash"},
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
		{
			Typ: resp.ARRAY.Typ,
			Array: []resp.Value{
				// 1. command
				{Typ: resp.BULK.Typ, Bulk: "COMMAND"},
				// 2. arg count
				{Typ: resp.INTEGER.Typ, Num: -1},
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
						{Typ: resp.BULK.Typ, Bulk: "@connection"},
						{Typ: resp.BULK.Typ, Bulk: "@slow"},
					},
				},
			},
		},
		{
			Typ: resp.ARRAY.Typ,
			Array: []resp.Value{
				// 1. command
				{Typ: resp.BULK.Typ, Bulk: "COMMAND DOCS"},
				// 2. arg count
				{Typ: resp.INTEGER.Typ, Num: -2},
				// 3. flags
				{
					Typ: resp.ARRAY.Typ,
					Array: []resp.Value{
						{Typ: resp.BULK.Typ, Bulk: "readonly"},
					},
				},
				// 4. first key
				{Typ: resp.INTEGER.Typ, Num: 2},
				// 5. last key
				{Typ: resp.INTEGER.Typ, Num: 2},
				// 6. steps between keys
				{Typ: resp.INTEGER.Typ, Num: 1},
				// 7. ACL flags
				{
					Typ: resp.ARRAY.Typ,
					Array: []resp.Value{
						{Typ: resp.BULK.Typ, Bulk: "@connection"},
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

func Test_commandDocs_returnsDocs(t *testing.T) {
	// given
	args := []resp.Value{
		{
			Typ:  resp.BULK.Typ,
			Bulk: "DOCS",
		},
	}

	expected := []resp.Value{
		{
			Typ:  resp.BULK.Typ,
			Bulk: "PING",
		},
		{
			Typ: resp.ARRAY.Typ,
			Array: []resp.Value{
				{Typ: resp.BULK.Typ, Bulk: "summary"},
				{Typ: resp.BULK.Typ, Bulk: "Returns PONG if no argument is provided, otherwise return a copy of the argument as a bulk."},

				{Typ: resp.BULK.Typ, Bulk: "since"},
				{Typ: resp.BULK.Typ, Bulk: "1.0.0"},

				{Typ: resp.BULK.Typ, Bulk: "group"},
				{Typ: resp.BULK.Typ, Bulk: "connection"},

				{Typ: resp.BULK.Typ, Bulk: "complexity"},
				{Typ: resp.BULK.Typ, Bulk: "O(1)"},
			},
		},
		{
			Typ:  resp.BULK.Typ,
			Bulk: "GET",
		},
		{
			Typ: resp.ARRAY.Typ,
			Array: []resp.Value{
				{Typ: resp.BULK.Typ, Bulk: "summary"},
				{Typ: resp.BULK.Typ, Bulk: "Get the value of key."},

				{Typ: resp.BULK.Typ, Bulk: "since"},
				{Typ: resp.BULK.Typ, Bulk: "1.0.0"},

				{Typ: resp.BULK.Typ, Bulk: "group"},
				{Typ: resp.BULK.Typ, Bulk: "string"},

				{Typ: resp.BULK.Typ, Bulk: "complexity"},
				{Typ: resp.BULK.Typ, Bulk: "O(1)"},
			},
		},
		{
			Typ:  resp.BULK.Typ,
			Bulk: "SET",
		},
		{
			Typ: resp.ARRAY.Typ,
			Array: []resp.Value{
				{Typ: resp.BULK.Typ, Bulk: "summary"},
				{Typ: resp.BULK.Typ, Bulk: "Set key to hold the string value."},

				{Typ: resp.BULK.Typ, Bulk: "since"},
				{Typ: resp.BULK.Typ, Bulk: "1.0.0"},

				{Typ: resp.BULK.Typ, Bulk: "group"},
				{Typ: resp.BULK.Typ, Bulk: "string"},

				{Typ: resp.BULK.Typ, Bulk: "complexity"},
				{Typ: resp.BULK.Typ, Bulk: "O(1)"},
			},
		},
		{
			Typ:  resp.BULK.Typ,
			Bulk: "DEL",
		},
		{
			Typ: resp.ARRAY.Typ,
			Array: []resp.Value{
				{Typ: resp.BULK.Typ, Bulk: "summary"},
				{Typ: resp.BULK.Typ, Bulk: "Removes the specified keys."},

				{Typ: resp.BULK.Typ, Bulk: "since"},
				{Typ: resp.BULK.Typ, Bulk: "1.0.0"},

				{Typ: resp.BULK.Typ, Bulk: "group"},
				{Typ: resp.BULK.Typ, Bulk: "keyspace"},

				{Typ: resp.BULK.Typ, Bulk: "complexity"},
				{Typ: resp.BULK.Typ, Bulk: "O(1) - O(N)"},
			},
		},
		{
			Typ:  resp.BULK.Typ,
			Bulk: "HGET",
		},
		{
			Typ: resp.ARRAY.Typ,
			Array: []resp.Value{
				{Typ: resp.BULK.Typ, Bulk: "summary"},
				{Typ: resp.BULK.Typ, Bulk: "Returns the value associated with field in the hash stored at key."},

				{Typ: resp.BULK.Typ, Bulk: "since"},
				{Typ: resp.BULK.Typ, Bulk: "2.0.0"},

				{Typ: resp.BULK.Typ, Bulk: "group"},
				{Typ: resp.BULK.Typ, Bulk: "hash"},

				{Typ: resp.BULK.Typ, Bulk: "complexity"},
				{Typ: resp.BULK.Typ, Bulk: "O(1)"},
			},
		},
		{
			Typ:  resp.BULK.Typ,
			Bulk: "HSET",
		},
		{
			Typ: resp.ARRAY.Typ,
			Array: []resp.Value{
				{Typ: resp.BULK.Typ, Bulk: "summary"},
				{Typ: resp.BULK.Typ, Bulk: "Sets the specified fields to their respective values in the hash stored at key."},

				{Typ: resp.BULK.Typ, Bulk: "since"},
				{Typ: resp.BULK.Typ, Bulk: "2.0.0"},

				{Typ: resp.BULK.Typ, Bulk: "group"},
				{Typ: resp.BULK.Typ, Bulk: "hash"},

				{Typ: resp.BULK.Typ, Bulk: "complexity"},
				{Typ: resp.BULK.Typ, Bulk: "O(1)"},
			},
		},
		{
			Typ:  resp.BULK.Typ,
			Bulk: "HDEL",
		},
		{
			Typ: resp.ARRAY.Typ,
			Array: []resp.Value{
				{Typ: resp.BULK.Typ, Bulk: "summary"},
				{Typ: resp.BULK.Typ, Bulk: "Removes the specified fields from the hash stored at key."},

				{Typ: resp.BULK.Typ, Bulk: "since"},
				{Typ: resp.BULK.Typ, Bulk: "2.0.0"},

				{Typ: resp.BULK.Typ, Bulk: "group"},
				{Typ: resp.BULK.Typ, Bulk: "keyspace"},

				{Typ: resp.BULK.Typ, Bulk: "complexity"},
				{Typ: resp.BULK.Typ, Bulk: "O(N)"},
			},
		},
		{
			Typ:  resp.BULK.Typ,
			Bulk: "HGETALL",
		},
		{
			Typ: resp.ARRAY.Typ,
			Array: []resp.Value{
				{Typ: resp.BULK.Typ, Bulk: "summary"},
				{Typ: resp.BULK.Typ, Bulk: "Returns all fields and values of the hash stored at key."},

				{Typ: resp.BULK.Typ, Bulk: "since"},
				{Typ: resp.BULK.Typ, Bulk: "2.0.0"},

				{Typ: resp.BULK.Typ, Bulk: "group"},
				{Typ: resp.BULK.Typ, Bulk: "hash"},

				{Typ: resp.BULK.Typ, Bulk: "complexity"},
				{Typ: resp.BULK.Typ, Bulk: "O(N)"},
			},
		},
		{
			Typ:  resp.BULK.Typ,
			Bulk: "COMMAND",
		},
		{
			Typ: resp.ARRAY.Typ,
			Array: []resp.Value{
				{Typ: resp.BULK.Typ, Bulk: "summary"},
				{Typ: resp.BULK.Typ, Bulk: "Return an array with details about every Redis command."},
				{Typ: resp.BULK.Typ, Bulk: "since"},
				{Typ: resp.BULK.Typ, Bulk: "2.8.13"},
				{Typ: resp.BULK.Typ, Bulk: "group"},
				{Typ: resp.BULK.Typ, Bulk: "connection"},
				{Typ: resp.BULK.Typ, Bulk: "complexity"},
				{Typ: resp.BULK.Typ, Bulk: "O(N)"},
			},
		},
		{
			Typ:  resp.BULK.Typ,
			Bulk: "COMMAND DOCS",
		},
		{
			Typ: resp.ARRAY.Typ,
			Array: []resp.Value{
				{Typ: resp.BULK.Typ, Bulk: "summary"},
				{Typ: resp.BULK.Typ, Bulk: "Return documentary information about commands."},

				{Typ: resp.BULK.Typ, Bulk: "since"},
				{Typ: resp.BULK.Typ, Bulk: "7.0.0"},

				{Typ: resp.BULK.Typ, Bulk: "group"},
				{Typ: resp.BULK.Typ, Bulk: "connection"},

				{Typ: resp.BULK.Typ, Bulk: "complexity"},
				{Typ: resp.BULK.Typ, Bulk: "O(N)"},
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
	assert.Equal(t, result.Array, expected)
}

func Test_commandDocs_withFilter_caseInsensitive_returnsDocsOfFilter(t *testing.T) {
	// given
	args := []resp.Value{
		{
			Typ:  resp.BULK.Typ,
			Bulk: "DOCS",
		},
		{
			Typ:  resp.BULK.Typ,
			Bulk: "PiNg",
		},
	}

	expected := []resp.Value{
		{
			Typ:  resp.BULK.Typ,
			Bulk: "PING",
		},
		{
			Typ: resp.ARRAY.Typ,
			Array: []resp.Value{
				{Typ: resp.BULK.Typ, Bulk: "summary"},
				{Typ: resp.BULK.Typ, Bulk: "Returns PONG if no argument is provided, otherwise return a copy of the argument as a bulk."},

				{Typ: resp.BULK.Typ, Bulk: "since"},
				{Typ: resp.BULK.Typ, Bulk: "1.0.0"},

				{Typ: resp.BULK.Typ, Bulk: "group"},
				{Typ: resp.BULK.Typ, Bulk: "connection"},

				{Typ: resp.BULK.Typ, Bulk: "complexity"},
				{Typ: resp.BULK.Typ, Bulk: "O(1)"},
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
	assert.Equal(t, result.Array, expected)
}
