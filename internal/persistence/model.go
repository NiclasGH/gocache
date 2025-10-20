package persistence

import (
	"time"
)

type Expirationable struct {
	ExpiresAt time.Time
}

func newExpiration(now time.Time, expireDuration time.Duration) *Expirationable {
	return &Expirationable{
		ExpiresAt: now.Add(expireDuration),
	}
}

func (e *Expirationable) isExpired(now time.Time) bool {
	return now.After(e.ExpiresAt)
}

type StringEntity struct {
	Value      string
	Expiration *Expirationable
}

func NewString(value string, expireDuration time.Duration) StringEntity {
	entity := StringEntity{
		Value:      value,
		Expiration: nil,
	}

	if expireDuration > 0 {
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

func (s *StringEntity) IsExpired() bool {
	if s.Expiration == nil {
		return false
	}

	return s.Expiration.isExpired(time.Now().UTC())
}
