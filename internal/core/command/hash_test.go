package command

import (
	"gocache/internal/core/resp"
	"testing"

	"github.com/stretchr/testify/assert"
)

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
