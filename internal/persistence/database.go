package persistence

import (
	"errors"
	"gocache/internal/core/resp"
	"sync"
)

type setStorage struct {
	store map[string]string
	mutex sync.RWMutex
}
type hsetStorage struct {
	store map[string]map[string]string
	mutex sync.RWMutex
}

// TODO transactions for disk persistence rollback
type DatabaseImpl struct {
	setStorage  setStorage
	hsetStorage hsetStorage

	diskPersistence diskPersistence
}

func NewDatabase(path string) (*DatabaseImpl, error) {
	aof, err := newAof(path)
	if err != nil {
		return nil, err
	}

	return &DatabaseImpl{
		setStorage: setStorage{
			store: map[string]string{},
		},
		hsetStorage: hsetStorage{
			store: map[string]map[string]string{},
		},

		diskPersistence: aof,
	}, nil
}

func (db *DatabaseImpl) SaveSet(requestValue resp.Value, key string, value string) error {
	db.setStorage.mutex.Lock()
	db.setStorage.store[key] = value
	db.setStorage.mutex.Unlock()

	return db.diskPersistence.Save(requestValue)
}

func (db *DatabaseImpl) GetSet(key string) (string, error) {
	db.setStorage.mutex.RLock()
	value, ok := db.setStorage.store[key]
	db.setStorage.mutex.RUnlock()

	if !ok {
		return "", errors.New("No value with key: " + key)
	}

	return value, nil
}

func (db *DatabaseImpl) SaveHSet(requestValue resp.Value, hash string, key string, value string) error {
	db.hsetStorage.mutex.Lock()
	if _, ok := db.hsetStorage.store[hash]; !ok {
		db.hsetStorage.store[hash] = map[string]string{}
	}
	db.hsetStorage.store[hash][key] = value
	db.hsetStorage.mutex.Unlock()

	return db.diskPersistence.Save(requestValue)
}

func (db *DatabaseImpl) GetHSet(hash string) (map[string]string, error) {
	db.hsetStorage.mutex.RLock()
	value, ok := db.hsetStorage.store[hash]
	db.hsetStorage.mutex.RUnlock()

	if !ok {
		return nil, errors.New("Did not find any value with hash " + hash)
	}

	return value, nil
}

func (db *DatabaseImpl) GetInit() ([]resp.Value, error) {
	return db.diskPersistence.GetInit()
}

func (db *DatabaseImpl) Close() error {
	return db.diskPersistence.Close()
}
