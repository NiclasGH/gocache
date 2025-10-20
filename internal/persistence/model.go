package persistence

import "time"

type Expirationable struct {
	expiresAt time.Time
}

func newExpiration(now time.Time, expireDuration time.Duration) Expirationable {
	return Expirationable{
		expiresAt: now.Add(expireDuration),
	}
}

type StringEntity struct {
	value      string
	expiration Expirationable
}

func NewString(value string, expireDuration time.Duration) StringEntity {
	entity := StringEntity{value: value}

	if expireDuration <= 0 {
		entity.expiration = newExpiration(time.Now().UTC(), expireDuration)
	}

	return entity
}
