package infrastructure

import (
	"gocache/internal/core/expiration"
	"gocache/internal/persistence"
	"time"
)

// Could theoretically make this configurable but eh
const amountOfKeys = 10

func ExpirationJob(delay time.Duration, db persistence.Database) {
	for {
		expiration.ExpireRandomKeys(amountOfKeys, db)
		time.Sleep(delay)
	}
}
