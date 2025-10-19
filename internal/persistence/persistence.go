package persistence

import "gocache/internal/core/resp"

type Database interface {
	Initialize(func(resp.Value)) error
	Save(resp.Value) error
	Close() error
}
