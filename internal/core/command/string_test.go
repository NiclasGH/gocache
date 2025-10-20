package command

import (
	"gocache/internal/core/resp"
	"gocache/internal/persistence"
	"testing"

	"github.com/stretchr/testify/assert"
)

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
