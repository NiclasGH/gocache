package persistence

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func Test_expirationCalculation(t *testing.T) {
	// given
	now := time.Now()
	in60Seconds := now.Add(time.Minute)

	// when
	expiration := newExpiration(now, time.Minute)

	// then
	assert.Equal(t, in60Seconds, expiration.ExpiresAt)
}
