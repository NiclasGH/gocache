package persistence

import (
	"errors"
	"gocache/internal/core/resp"
	"sync"
)

type stringStorage struct {
	store map[string]StringEntity
	mutex sync.RWMutex
}
type hashStorage struct {
	store map[string]map[string]string
	mutex sync.RWMutex
}

type DatabaseImpl struct {
	stringStorage stringStorage
	hashStorage   hashStorage

	// could also be a list, to enable multiple forms of disk persistence (aof, snapshots etc)
	diskPersistence DiskPersistence
}

func NewDatabase(diskPersistence DiskPersistence) *DatabaseImpl {
	return &DatabaseImpl{
		stringStorage: stringStorage{
			store: map[string]StringEntity{},
		},
		hashStorage: hashStorage{
			store: map[string]map[string]string{},
		},

		diskPersistence: diskPersistence,
	}
}
func (db *DatabaseImpl) EnablePersistence(diskPersistence DiskPersistence) {
	db.diskPersistence = diskPersistence
}

func (db *DatabaseImpl) SaveString(requestValue resp.Value, key string, value StringEntity) error {
	db.stringStorage.mutex.Lock()
	defer db.stringStorage.mutex.Unlock()

	if db.diskPersistence != nil {
		if err := db.diskPersistence.Save(requestValue); err != nil {
			return err
		}
	}

	db.stringStorage.store[key] = value

	return nil
}

func (db *DatabaseImpl) GetString(key string) (StringEntity, error) {
	db.stringStorage.mutex.RLock()
	value, ok := db.stringStorage.store[key]
	db.stringStorage.mutex.RUnlock()

	if !ok {
		return StringEntity{}, errors.New("No value with key: " + key)
	}

	return value, nil
}

func (db *DatabaseImpl) GetRandomString() (string, StringEntity, bool) {
	db.stringStorage.mutex.RLock()
	defer db.stringStorage.mutex.RUnlock()

	for k, v := range db.stringStorage.store {
		return k, v, true
	}
	return "", StringEntity{}, false
}

func (db *DatabaseImpl) DeleteAllStrings(requestValue resp.Value, keys []string) (int, error) {
	db.stringStorage.mutex.Lock()
	defer db.stringStorage.mutex.Unlock()

	if db.diskPersistence != nil {
		if err := db.diskPersistence.Save(requestValue); err != nil {
			return 0, err
		}
	}

	amountDeleted := 0
	for _, key := range keys {
		if _, ok := db.stringStorage.store[key]; ok {
			delete(db.stringStorage.store, key)
			amountDeleted += 1
		}
	}

	return amountDeleted, nil
}

func (db *DatabaseImpl) SaveHash(requestValue resp.Value, hash string, key string, value string) error {
	db.hashStorage.mutex.Lock()
	defer db.hashStorage.mutex.Unlock()

	if db.diskPersistence != nil {
		if err := db.diskPersistence.Save(requestValue); err != nil {
			return err
		}
	}

	if _, ok := db.hashStorage.store[hash]; !ok {
		db.hashStorage.store[hash] = map[string]string{}
	}
	db.hashStorage.store[hash][key] = value

	return nil
}

func (db *DatabaseImpl) DeleteAllHashKeys(requestValue resp.Value, hash string, keys []string) (int, error) {
	db.hashStorage.mutex.Lock()
	defer db.hashStorage.mutex.Unlock()

	if db.diskPersistence != nil {
		if err := db.diskPersistence.Save(requestValue); err != nil {
			return 0, err
		}
	}

	hashMap, ok := db.hashStorage.store[hash]
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
		delete(db.hashStorage.store, hash)
	}

	return amountDeleted, nil
}

func (db *DatabaseImpl) GetHash(hash string) (map[string]string, error) {
	db.hashStorage.mutex.RLock()
	value, ok := db.hashStorage.store[hash]
	db.hashStorage.mutex.RUnlock()

	if !ok {
		return nil, errors.New("Did not find any value with hash " + hash)
	}

	return value, nil
}

func (db *DatabaseImpl) Close() error {
	return db.diskPersistence.Close()
}
