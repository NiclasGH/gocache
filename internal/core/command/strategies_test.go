package command

import (
	"gocache/internal/core/resp"
	"gocache/internal/persistence"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_ping(t *testing.T) {
	// given
	expected := resp.Value{Typ: "string", Str: "PONG"}

	ping, ok := Strategies[PING]
	if !ok {
		t.Error("Command does not exist")
		return
	}

	args := request(PING, []resp.Value{})

	// when
	result := ping(args, defaultDb())

	// then
	assert.EqualValues(t, expected, result)
}

func Test_pingWithArg(t *testing.T) {
	// given
	args := []resp.Value{
		{
			Typ:  resp.BULK.Typ,
			Bulk: "Tiramisu",
		},
	}

	ping, ok := Strategies[PING]
	if !ok {
		t.Error("Command does not exist")
		return
	}

	expected := resp.Value{Typ: "string", Str: "Tiramisu"}

	// when
	result := ping(request(PING, args), defaultDb())

	// then
	assert.EqualValues(t, expected, result)
}

func Test_set(t *testing.T) {
	// given
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

	set, ok := Strategies[SET]
	if !ok {
		t.Error("Command does not exist")
		return
	}

	db := defaultDb()

	// when
	result := set(request(SET, args), db)

	// then
	assert.EqualValues(t, expected, result)

	value, err := db.GetString("Tira")
	if err != nil {
		t.Error("Set Storage did not contain key 'Tira'")
		return
	}
	assert.Equal(t, "Misu", value.Value)
}

func Test_setNeedsTwoArgs(t *testing.T) {
	// given
	args := []resp.Value{
		{
			Typ:  resp.BULK.Typ,
			Bulk: "Tira",
		},
	}

	set, ok := Strategies[SET]
	if !ok {
		t.Error("Command does not exist")
		return
	}

	// when
	result := set(request(SET, args), defaultDb())

	// then
	assert.Equal(t, "error", result.Typ)
}

func Test_incr(t *testing.T) {
	// given
	db := defaultDb()
	db.SaveString(resp.Value{}, "Tira", persistence.NewString("5", 0))

	args := []resp.Value{
		{
			Typ:  resp.BULK.Typ,
			Bulk: "Tira",
		},
	}

	expected := resp.Value{
		Typ: "integer",
		Num: 6,
	}

	incr, ok := Strategies[INCR]
	if !ok {
		t.Error("Command does not exist")
		return
	}

	// when
	result := incr(request(INCR, args), db)

	// then
	assert.EqualValues(t, expected, result)

	value, err := db.GetString("Tira")
	if err != nil {
		t.Error("Set Storage Key 'Tira' does not exist")
		return
	}
	assert.Equal(t, "6", value.Value)
}

func Test_incr_needsOneArg(t *testing.T) {
	// given
	args := []resp.Value{}

	incr, ok := Strategies[INCR]
	if !ok {
		t.Error("Command does not exist")
		return
	}

	// when
	result := incr(request(INCR, args), defaultDb())

	// then
	assert.Equal(t, resp.ERROR.Typ, result.Typ)
}

func Test_incr_needsStringToBeNumber(t *testing.T) {
	// given
	db := defaultDb()
	db.SaveString(resp.Value{}, "Tira", persistence.NewString("number", 0))

	args := []resp.Value{
		{
			Typ:  resp.BULK.Typ,
			Bulk: "Tira",
		},
	}

	incr, ok := Strategies[INCR]
	if !ok {
		t.Error("Command does not exist")
		return
	}

	// when
	result := incr(request(INCR, args), db)

	// then
	assert.Equal(t, resp.ERROR.Typ, result.Typ)
}

func Test_incr_createsKeyIfNotExists(t *testing.T) {
	// given
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

	incr, ok := Strategies[INCR]
	if !ok {
		t.Error("Command does not exist")
		return
	}

	db := defaultDb()

	// when
	result := incr(request(INCR, args), db)

	// then
	assert.EqualValues(t, expected, result)

	value, err := db.GetString("Tira")
	if err != nil {
		t.Error("Set Storage Key 'Tira' does not exist")
		return
	}
	assert.Equal(t, "1", value.Value)
}

func Test_del(t *testing.T) {
	// given
	db := defaultDb()
	db.SaveString(resp.Value{}, "Tira", persistence.NewString("Misu", 0))

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

	del, ok := Strategies[DEL]
	if !ok {
		t.Error("Command does not exist")
		return
	}

	// when
	result := del(request(DEL, args), db)

	// then
	assert.EqualValues(t, expected, result)

	_, err := db.GetString("Tira")
	if err == nil {
		t.Error("Set Storage Key 'Tira' was not deleted")
		return
	}
}

func Test_del_multipleKeys(t *testing.T) {
	// given
	db := defaultDb()
	db.SaveString(resp.Value{}, "Tira", persistence.NewString("Misu", 0))
	db.SaveString(resp.Value{}, "Misu", persistence.NewString("Tira", 0))

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

	del, ok := Strategies[DEL]
	if !ok {
		t.Error("Command does not exist")
		return
	}

	// when
	result := del(request(DEL, args), db)

	// then
	assert.EqualValues(t, expected, result)

	_, err := db.GetString("Tira")
	if err == nil {
		t.Error("Set Storage Key 'Tira' was not deleted")
		return
	}

	_, err = db.GetString("Tira")
	if err == nil {
		t.Error("Set Storage Key 'Misu' was not deleted")
		return
	}
}

func Test_del_needsAtLeastOneKey(t *testing.T) {
	// given
	args := []resp.Value{}

	del, ok := Strategies[DEL]
	if !ok {
		t.Error("Command does not exist")
		return
	}

	// when
	result := del(request(DEL, args), defaultDb())

	// then
	assert.Equal(t, "error", result.Typ)
}

func Test_get(t *testing.T) {
	// given
	db := defaultDb()
	db.SaveString(resp.Value{}, "Tira", persistence.NewString("Misu", 0))

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

	get, ok := Strategies[GET]
	if !ok {
		t.Error("Command does not exist")
		return
	}

	// when
	result := get(request(GET, args), db)

	// then
	assert.EqualValues(t, expected, result)

	value, err := db.GetString("Tira")
	if err != nil {
		t.Error("Set Storage did not contain key 'Tira'")
		return
	}
	assert.Equal(t, "Misu", value.Value)
}

func Test_getCanOnlyReceiveOneArg(t *testing.T) {
	// given
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

	get, ok := Strategies[GET]
	if !ok {
		t.Error("Command does not exist")
		return
	}

	// when
	result := get(request(GET, args), defaultDb())

	// then
	assert.Equal(t, "error", result.Typ)
}

func Test_getNoValueAvailable(t *testing.T) {
	// given
	args := []resp.Value{
		{
			Typ:  resp.BULK.Typ,
			Bulk: "Tira",
		},
	}

	expected := resp.Value{
		Typ: "null",
	}

	get, ok := Strategies[GET]
	if !ok {
		t.Error("Command does not exist")
		return
	}

	// when
	result := get(request(GET, args), defaultDb())

	// then
	assert.EqualValues(t, expected, result)
}

func Test_hset(t *testing.T) {
	// given
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

	hset, ok := Strategies[HSET]
	if !ok {
		t.Error("Command does not exist")
		return
	}

	db := defaultDb()

	// when
	result := hset(request(HSET, args), db)

	// then
	assert.EqualValues(t, expected, result)

	valueMap, err := db.GetHash("tira")
	if err != nil {
		t.Error("HSet Storage did not contain hash 'tira'")
		return
	}
	value, ok := valueMap["misu"]
	if !ok {
		t.Error("HSet Storage did not contain key 'misu'")
		return
	}
	assert.Equal(t, "cute", value)
}

func Test_hsetNeedsThreeArgs(t *testing.T) {
	// given
	args := []resp.Value{
		{
			Typ:  resp.BULK.Typ,
			Bulk: "tira",
		},
	}

	hset, ok := Strategies[HSET]
	if !ok {
		t.Error("Command does not exist")
		return
	}

	// when
	result := hset(request(HSET, args), defaultDb())

	// then
	assert.Equal(t, "error", result.Typ)
}

func Test_hget(t *testing.T) {
	// given
	db := defaultDb()
	db.SaveHash(resp.Value{}, "tira", "misu", "cute")

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

	hget, ok := Strategies[HGET]
	if !ok {
		t.Error("Command does not exist")
		return
	}

	// when
	result := hget(request(HGET, args), db)

	// then
	assert.EqualValues(t, expected, result)

	valueMap, err := db.GetHash("tira")
	if err != nil {
		t.Error("HSet Storage did not contain hash 'tira'")
		return
	}
	value, ok := valueMap["misu"]
	if !ok {
		t.Error("HSet Storage did not contain key 'misu'")
		return
	}

	assert.Equal(t, "cute", value)
}

func Test_hgetCanOnlyReceiveTwoArg(t *testing.T) {
	// given
	args := []resp.Value{
		{
			Typ:  resp.BULK.Typ,
			Bulk: "tira",
		},
	}

	hget, ok := Strategies[HGET]
	if !ok {
		t.Error("Command does not exist")
		return
	}

	// when
	result := hget(request(HGET, args), defaultDb())

	// then
	assert.Equal(t, "error", result.Typ)
}

func Test_hgetNoValueAvailable(t *testing.T) {
	// given
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

	hget, ok := Strategies[HGET]
	if !ok {
		t.Error("Command does not exist")
		return
	}

	// when
	result := hget(request(HGET, args), defaultDb())

	// then
	assert.EqualValues(t, expected, result)
}

func Test_hdel(t *testing.T) {
	// given
	db := defaultDb()
	db.SaveHash(resp.Value{}, "tira", "misu", "cute")
	db.SaveHash(resp.Value{}, "tira", "void", "scary")

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

	hdel, ok := Strategies[HDEL]
	if !ok {
		t.Error("Command does not exist")
		return
	}

	// when
	result := hdel(request(HDEL, args), db)

	// then
	assert.EqualValues(t, expected, result)

	valueMap, err := db.GetHash("tira")
	if err != nil {
		t.Error("HSet Storage did get hash 'tira' deleted")
		return
	}

	_, ok = valueMap["misu"]
	if ok {
		t.Error("HSet Storage did not get key 'misu' deleted")
		return
	}

	value, ok := valueMap["void"]
	if !ok {
		t.Error("HSet Storage did get key 'void' deleted but wasn't supposed to")
		return
	}
	assert.Equal(t, "scary", value)
}

func Test_hdel_lastKeyDeletesHash(t *testing.T) {
	// given
	db := defaultDb()
	db.SaveHash(resp.Value{}, "tira", "misu", "cute")

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

	hdel, ok := Strategies[HDEL]
	if !ok {
		t.Error("Command does not exist")
		return
	}

	// when
	result := hdel(request(HDEL, args), db)

	// then
	assert.EqualValues(t, expected, result)

	_, err := db.GetHash("tira")
	if err == nil {
		t.Error("HSet Storage did not get hash 'tira' deleted")
		return
	}
}

func Test_hdel_deleteMultipleFields(t *testing.T) {
	// given
	db := defaultDb()
	db.SaveHash(resp.Value{}, "tira", "misu", "cute")
	db.SaveHash(resp.Value{}, "tira", "void", "scary")
	db.SaveHash(resp.Value{}, "tira", "isSpider", "true")

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

	hdel, ok := Strategies[HDEL]
	if !ok {
		t.Error("Command does not exist")
		return
	}

	// when
	result := hdel(request(HDEL, args), db)

	// then
	assert.EqualValues(t, expected, result)

	valueMap, err := db.GetHash("tira")
	if err != nil {
		t.Error("Hash 'tira' got deleted but wasnt supposed to")
		return
	}

	_, ok = valueMap["misu"]
	if ok {
		t.Error("HSet Storage did not get key 'misu' deleted")
		return
	}

	_, ok = valueMap["void"]
	if ok {
		t.Error("HSet Storage did not get key 'void' deleted")
		return
	}

	value, ok := valueMap["isSpider"]
	if !ok {
		t.Error("HSet Storage did get key 'isSpider' deleted but wasn't supposed to")
		return
	}
	assert.Equal(t, "true", value)
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

	hdel, ok := Strategies[HDEL]
	if !ok {
		t.Error("Command does not exist")
		return
	}

	// when
	result := hdel(request(HDEL, args), defaultDb())

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
				{Typ: resp.BULK.Typ, Bulk: "INCR"},
				// 2. arg count
				{Typ: resp.INTEGER.Typ, Num: 2},
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
				{Typ: resp.INTEGER.Typ, Num: 1},
				// 6. steps between keys
				{Typ: resp.INTEGER.Typ, Num: 1},
				// 7. ACL flags
				{
					Typ: resp.ARRAY.Typ,
					Array: []resp.Value{
						{Typ: resp.BULK.Typ, Bulk: "@write"},
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

	commandSpecs, ok := Strategies["COMMAND"]
	if !ok {
		t.Error("Command does not exist")
		return
	}

	// when
	result := commandSpecs(request(COMMAND, []resp.Value{}), defaultDb())

	// then
	assert.ElementsMatch(t, expected, result.Array)
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

	commandSpecs, ok := Strategies["COMMAND"]
	if !ok {
		t.Error("Command does not exist")
		return
	}

	// when
	result := commandSpecs(request(COMMAND, args), defaultDb())

	// then
	assert.ElementsMatch(t, expected, result.Array)
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
			Bulk: "INCR",
		},
		{
			Typ: resp.ARRAY.Typ,
			Array: []resp.Value{
				{Typ: resp.BULK.Typ, Bulk: "summary"},
				{Typ: resp.BULK.Typ, Bulk: "Increments the number stored at key by one."},

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

	commandSpecs, ok := Strategies["COMMAND"]
	if !ok {
		t.Error("Command does not exist")
		return
	}

	// when
	result := commandSpecs(request(COMMAND, args), defaultDb())

	// then
	assert.Equal(t, expected, result.Array)
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

	commandSpecs, ok := Strategies["COMMAND"]
	if !ok {
		t.Error("Command does not exist")
		return
	}

	// when
	result := commandSpecs(request(COMMAND, args), defaultDb())

	// then
	assert.Equal(t, expected, result.Array)
}

func request(command string, args []resp.Value) resp.Value {
	request := resp.Value{
		Typ: resp.ARRAY.Typ,
		Array: []resp.Value{
			{
				Typ:  resp.BULK.Typ,
				Bulk: command,
			},
		},
	}
	request.Array = append(request.Array, args...)

	return request
}

func defaultDb() persistence.Database {
	return persistence.NewDatabase(nil)
}
