package expiration

import (
	"gocache/internal/core/resp"
	"gocache/internal/persistence"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func Test_expiresKey(t *testing.T) {
	// given
	db := defaultDb()
	// expires immidiately
	db.SaveString(resp.Value{}, "tira", persistence.NewString("misu", time.Nanosecond))

	// when
	ExpireRandomKeys(1, db)

	// then
	_, err := db.GetString("tira")
	if err == nil {
		t.Error("Key tira was not expired")
	}
}

func Test_doesntExpireKeyWithRemainingExpiration(t *testing.T) {
	// given
	db := defaultDb()
	// expires far into future
	db.SaveString(resp.Value{}, "tira", persistence.NewString("misu", time.Hour))

	// when
	ExpireRandomKeys(1, db)

	// then
	value, err := db.GetString("tira")
	if err != nil {
		t.Error("Key tira was expired")
		return
	}
	assert.Equal(t, "misu", value.Value)
}

func Test_ignoredKeyWithoutExpiration(t *testing.T) {
	// given
	db := defaultDb()
	// doesnt expire
	db.SaveString(resp.Value{}, "tira", persistence.NewString("misu", time.Duration(0)))

	// when
	ExpireRandomKeys(1, db)

	// then
	value, err := db.GetString("tira")
	if err != nil {
		t.Error("Tira wasn't supposed to be deleted")
		return
	}
	assert.Equal(t, "misu", value.Value)
}

func defaultDb() persistence.Database {
	return persistence.NewDatabase(nil)
}
