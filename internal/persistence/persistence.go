package persistence

import (
	"gocache/internal/core/resp"
)

type Database interface {
	SaveString(request resp.Value, key string, value string) error
	DeleteAllStrings(request resp.Value, keys []string) (int, error)
	GetString(key string) (string, error)

	SaveHash(request resp.Value, hash string, key string, value string) error
	DeleteAllHashKeys(request resp.Value, hash string, keys []string) (int, error)
	GetHash(hash string) (map[string]string, error)

	EnablePersistence(diskPersistence DiskPersistence)

	Close() error
}

type DiskPersistence interface {
	Save(resp.Value) error
	ReadPersistedCommands() ([]resp.Value, error)
	Close() error
}
