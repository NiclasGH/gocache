package startup

import (
	"errors"
	"gocache/internal/core/resp"
	"gocache/internal/persistence"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_startup_repeatsSet(t *testing.T) {
	// given
	key := "Tira"
	expected := "Misu"

	request := []resp.Value{
		{
			Typ: resp.ARRAY.Typ,
			Array: []resp.Value{
				{
					Typ:  resp.BULK.Typ,
					Bulk: "SET",
				},
				{
					Typ:  resp.BULK.Typ,
					Bulk: key,
				},
				{
					Typ:  resp.BULK.Typ,
					Bulk: expected,
				},
			},
		},
	}

	db := defaultDb()
	disk := defaultDisk(request)

	// when
	err := ReplayCommands(disk, db)

	// then
	if err != nil {
		t.Error(err.Error())
		return
	}

	value, err := db.GetString(key)
	if err != nil {
		t.Error("Value was not set")
	}
	assert.Equal(t, expected, value)
}

func Test_startup_repeatsHSet(t *testing.T) {
	// given
	hash := "Tira"
	key := "Misu"
	expected := "Cute"

	request := []resp.Value{
		{
			Typ: resp.ARRAY.Typ,
			Array: []resp.Value{
				{
					Typ:  resp.BULK.Typ,
					Bulk: "HSET",
				},
				{
					Typ:  resp.BULK.Typ,
					Bulk: hash,
				},
				{
					Typ:  resp.BULK.Typ,
					Bulk: key,
				},
				{
					Typ:  resp.BULK.Typ,
					Bulk: expected,
				},
			},
		},
	}

	db := defaultDb()
	disk := defaultDisk(request)

	// when
	err := ReplayCommands(disk, db)

	// then
	if err != nil {
		t.Error(err.Error())
		return
	}

	value, err := db.GetHash(hash)
	if err != nil {
		t.Error("Value was not set")
	}
	assert.Equal(t, expected, value[key])
}

func Test_startup_repeatsDel(t *testing.T) {
	// given
	key := "Tira"
	expected := "Misu"

	request := []resp.Value{
		{
			Typ: resp.ARRAY.Typ,
			Array: []resp.Value{
				{
					Typ:  resp.BULK.Typ,
					Bulk: "SET",
				},
				{
					Typ:  resp.BULK.Typ,
					Bulk: key,
				},
				{
					Typ:  resp.BULK.Typ,
					Bulk: expected,
				},
			},
		},
		{
			Typ: resp.ARRAY.Typ,
			Array: []resp.Value{
				{
					Typ:  resp.BULK.Typ,
					Bulk: "DEL",
				},
				{
					Typ:  resp.BULK.Typ,
					Bulk: key,
				},
			},
		},
	}

	db := defaultDb()
	disk := defaultDisk(request)

	// when
	err := ReplayCommands(disk, db)

	// then
	if err != nil {
		t.Error(err.Error())
		return
	}

	_, err = db.GetString(key)
	if err == nil {
		t.Error("Value was not deleted")
	}
}

func defaultDb() persistence.Database {
	return persistence.NewDatabase(nil)
}

func defaultDisk(request []resp.Value) persistence.DiskPersistence {
	return simpleDisk{request}
}

type simpleDisk struct {
	request []resp.Value
}

func (_ simpleDisk) Save(resp.Value) error {
	return errors.New("Save called but shouldnt be by the startup")
}

func (d simpleDisk) ReadPersistedCommands() ([]resp.Value, error) {
	return d.request, nil
}

func (_ simpleDisk) Close() error {
	return errors.New("Save called but shouldnt be by the startup")
}
