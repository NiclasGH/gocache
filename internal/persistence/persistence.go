package persistence

import "gocache/internal/resp"

type Database interface {
	Initialize(func(resp.Value)) error
	Save(resp.Value) error
	Close() error
}
