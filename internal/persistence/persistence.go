package persistence

import (
	"gocache/internal/core/resp"
)

type Database interface {
	SaveSet(resp.Value, string, string) error
	GetSet(string) (string, error)

	SaveHSet(resp.Value, string, string, string) error
	GetHSet(string) (map[string]string, error)

	GetInit() ([]resp.Value, error)
	Close() error
}

type diskPersistence interface {
	Save(resp.Value) error
	GetInit() ([]resp.Value, error)
	Close() error
}
