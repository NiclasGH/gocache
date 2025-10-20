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

type DatabaseImpl struct {
	setStorage  setStorage
	hsetStorage hsetStorage

	// could also be a list, to enable multiple forms of disk persistence (aof, snapshots etc)
	diskPersistence DiskPersistence
}

func NewDatabase(diskPersistence DiskPersistence) *DatabaseImpl {
	return &DatabaseImpl{
		setStorage: setStorage{
			store: map[string]string{},
		},
		hsetStorage: hsetStorage{
			store: map[string]map[string]string{},
		},

		diskPersistence: diskPersistence,
	}
}
func (db *DatabaseImpl) EnablePersistence(diskPersistence DiskPersistence) {
	db.diskPersistence = diskPersistence
}

func (db *DatabaseImpl) SaveSet(requestValue resp.Value, key string, value string) error {
	db.setStorage.mutex.Lock()
	defer db.setStorage.mutex.Unlock()

	if db.diskPersistence != nil {
		if err := db.diskPersistence.Save(requestValue); err != nil {
			return err
		}
	}

	db.setStorage.store[key] = value

	return nil
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

func (db *DatabaseImpl) DeleteAllSet(requestValue resp.Value, keys []string) (int, error) {
	db.setStorage.mutex.Lock()
	defer db.setStorage.mutex.Unlock()

	if db.diskPersistence != nil {
		if err := db.diskPersistence.Save(requestValue); err != nil {
			return 0, err
		}
	}

	amountDeleted := 0
	for _, key := range keys {
		if _, ok := db.setStorage.store[key]; ok {
			delete(db.setStorage.store, key)
			amountDeleted += 1
		}
	}

	return amountDeleted, nil
}

func (db *DatabaseImpl) SaveHSet(requestValue resp.Value, hash string, key string, value string) error {
	db.hsetStorage.mutex.Lock()
	defer db.hsetStorage.mutex.Unlock()

	if db.diskPersistence != nil {
		if err := db.diskPersistence.Save(requestValue); err != nil {
			return err
		}
	}

	if _, ok := db.hsetStorage.store[hash]; !ok {
		db.hsetStorage.store[hash] = map[string]string{}
	}
	db.hsetStorage.store[hash][key] = value

	return nil
}

func (db *DatabaseImpl) DeleteAllHSet(requestValue resp.Value, hash string, keys []string) (int, error) {
	db.hsetStorage.mutex.Lock()
	defer db.hsetStorage.mutex.Unlock()

	if db.diskPersistence != nil {
		if err := db.diskPersistence.Save(requestValue); err != nil {
			return 0, err
		}
	}

	hashMap, ok := db.hsetStorage.store[hash]
	if !ok {
		return 0, nil
	}

	amountDeleted := 0
	for _, key := range keys {
		if _, ok := hashMap[key]; ok {
			delete(hashMap, key)
			amountDeleted++
		}
	}

	if len(hashMap) == 0 {
		delete(db.hsetStorage.store, hash)
	}

	return amountDeleted, nil
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

func (db *DatabaseImpl) Close() error {
	return db.diskPersistence.Close()
}
