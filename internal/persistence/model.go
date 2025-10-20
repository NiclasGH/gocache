package persistence

import "time"

type Expirationable struct {
	ExpiresAt time.Time
}

func newExpiration(now time.Time, expireDuration time.Duration) Expirationable {
	return Expirationable{
		ExpiresAt: now.Add(expireDuration),
	}
}

type StringEntity struct {
	Value      string
	Expiration Expirationable
}

func NewString(value string, expireDuration time.Duration) StringEntity {
	entity := StringEntity{Value: value}

	if expireDuration <= 0 {
		entity.Expiration = newExpiration(time.Now().UTC(), expireDuration)
	}

	return entity
}

func (s *StringEntity) SetValue(value string) {
	s.Value = value
}

func (s *StringEntity) SetExpiration(expireDuration time.Duration) {
	s.Expiration = newExpiration(time.Now().UTC(), expireDuration)
}
