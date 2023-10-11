package kademlia

import (
	"errors"
	"sync"
	"time"

	"github.com/arianfiftyone/src/logger"
)

// DataStore represents a key-value data store.
type DataStore struct {
	data       map[[KeySize]byte]string      // Map to store key-value pairs.
	time       map[[KeySize]byte]time.Time   // Map to store key-time for expiration pairs. Unix time.
	mutexLocks map[[KeySize]byte]*sync.Mutex // Map to store a mutex lock for each key.
	ttl        time.Duration                 //seconds
}

// NewDataStore initializes a new DataStore instance.
func NewDataStore() DataStore {
	dataStore := DataStore{}
	dataStore.data = make(map[[KeySize]byte]string)
	dataStore.time = make(map[[KeySize]byte]time.Time)
	dataStore.mutexLocks = make(map[[20]byte]*sync.Mutex)
	dataStore.ttl = time.Second * 10
	return dataStore
}

// Insert inserts a key-value pair into the DataStore.
func (dataStore DataStore) Insert(key *Key, value string) {
	dataStore.mutexLocks[key.Hash] = &sync.Mutex{}
	mutex := dataStore.mutexLocks[key.Hash]
	mutex.Lock()

	dataStore.time[key.Hash] = dataStore.calculateExpirationTime()
	dataStore.data[key.Hash] = value

	mutex.Unlock()
	go dataStore.deleteAfterExpirationTime(dataStore.ttl, key)
}

func (dataStore DataStore) deleteAfterExpirationTime(timer time.Duration, key *Key) {

	select {
	case currTime := <-time.After(timer):
		mutex, ok := dataStore.mutexLocks[key.Hash]
		if !ok {
			return
		}
		mutex.Lock()

		_, ok = dataStore.time[key.Hash]
		if !ok {
			mutex.Unlock()
			return
		}

		differenceInTime := dataStore.time[key.Hash].Sub(currTime)
		mutex.Unlock()
		if differenceInTime <= 0 {
			dataStore.Delete(key)
		} else {
			dataStore.deleteAfterExpirationTime(differenceInTime, key)
		}

	}
}

// Get retrieves the value associated with a key from the DataStore.
func (dataStore DataStore) Get(key *Key) (string, error) {
	mutex, ok := dataStore.mutexLocks[key.Hash]
	if !ok {
		return "", errors.New("key not found")
	}
	mutex.Lock()

	value, ok := dataStore.data[key.Hash]
	if !ok {
		mutex.Unlock()
		return "", errors.New("key not found")
	}
	mutex.Unlock()
	err := dataStore.RefreshExpirationTime(key)
	if err != nil {
		return "", errors.New("Refresh failed")
	}

	return value, nil
}

func (dataStore DataStore) GetTime(key *Key) (time.Time, error) {
	mutex, ok := dataStore.mutexLocks[key.Hash]
	if !ok {
		return time.Now(), errors.New("key not found")
	}
	mutex.Lock()

	time, ok := dataStore.time[key.Hash]
	if !ok {
		mutex.Unlock()
		return time, errors.New("key not found")
	}
	mutex.Unlock()
	return time, nil
}

func (dataStore DataStore) calculateExpirationTime() time.Time {
	return time.Now().Add(dataStore.ttl)
}

func (dataStore DataStore) RefreshExpirationTime(key *Key) error {
	mutex, ok := dataStore.mutexLocks[key.Hash]
	if !ok {
		return errors.New("key not found")
	}
	mutex.Lock()

	_, ok = dataStore.data[key.Hash]
	if !ok {
		mutex.Unlock()
		return errors.New("key not found")
	}
	ttl := dataStore.calculateExpirationTime()
	dataStore.time[key.Hash] = ttl
	mutex.Unlock()
	return nil
}

func (dataStore DataStore) Delete(key *Key) error {
	mutex, ok := dataStore.mutexLocks[key.Hash]
	if !ok {
		return errors.New("key not found")
	}
	mutex.Lock()

	value, ok := dataStore.data[key.Hash]
	if !ok {
		mutex.Unlock()
		return errors.New("key not found")
	}
	_, ok = dataStore.time[key.Hash]
	if !ok {
		mutex.Unlock()
		return errors.New("key not found")
	}

	delete(dataStore.time, key.Hash)
	delete(dataStore.data, key.Hash)
	delete(dataStore.mutexLocks, key.Hash)
	mutex.Unlock()
	logger.Log("The data object " + key.GetHashString() + " with the value " + value + " has been deleted due to the expired TTL.")
	return nil
}
