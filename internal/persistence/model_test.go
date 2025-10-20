package persistence

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func Test_string_expirationCreation(t *testing.T) {
	// given
	now := time.Now()
	in60Seconds := now.Add(time.Minute)

	// when
	expiration := newExpiration(now, time.Minute)

	// then
	assert.Equal(t, in60Seconds, expiration.ExpiresAt)
}

func Test_string_expirationCalculation_false(t *testing.T) {
	// given
	str := NewString("value", time.Minute)

	// when
	expired := str.IsExpired()

	// then
	assert.Equal(t, false, expired)
}

func Test_string_expirationCalculation_true(t *testing.T) {
	// given
	str := NewString("value", time.Nanosecond)
	// wait 2 nanoseconds to expire the string
	time.Sleep(time.Nanosecond * 2)

	// when
	expired := str.IsExpired()

	// then
	assert.Equal(t, true, expired)
}
