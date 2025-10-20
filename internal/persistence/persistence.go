package persistence

import (
	"gocache/internal/core/resp"
)

type Database interface {
	SaveSet(request resp.Value, key string, value string) error
	DeleteAllSet(request resp.Value, keys []string) (int, error)
	GetSet(key string) (string, error)

	SaveHSet(request resp.Value, hash string, key string, value string) error
	DeleteAllHSet(request resp.Value, hash string, keys []string) (int, error)
	GetHSet(hash string) (map[string]string, error)

	EnablePersistence(diskPersistence DiskPersistence)

	Close() error
}

type DiskPersistence interface {
	Save(resp.Value) error
	ReadPersistedCommands() ([]resp.Value, error)
	Close() error
}
