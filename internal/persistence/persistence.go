package persistence

import (
	"gocache/internal/core/resp"
)

type Database interface {
	SaveSet(request resp.Value, key string, value string) error
	DeleteAllSet(request resp.Value, keys []string) int
	GetSet(key string) (string, error)

	SaveHSet(request resp.Value, hash string, key string, value string) error
	DeleteAllHSet(request resp.Value, hash string, keys []string) int
	GetHSet(hash string) (map[string]string, error)

	GetInit() ([]resp.Value, error)
	Close() error
}

type diskPersistence interface {
	Save(resp.Value) error
	GetInit() ([]resp.Value, error)
	Close() error
}
