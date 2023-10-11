package kademlia

import (
	"encoding/hex"
	"errors"
	"time"

	"github.com/arianfiftyone/src/logger"
)

// DataStore represents a key-value data store.
type DataStore struct {
	data map[[KeySize]byte]string    // Map to store key-value pairs.
	time map[[KeySize]byte]time.Time // Map to store key-time for expiration pairs. Unix time.
	ttl  time.Duration               //seconds
}

// NewDataStore initializes a new DataStore instance.
func NewDataStore() DataStore {
	dataStore := DataStore{}
	dataStore.data = make(map[[KeySize]byte]string)
	dataStore.time = make(map[[KeySize]byte]time.Time)
	dataStore.ttl = time.Second * 10
	return dataStore
}

// Insert inserts a key-value pair into the DataStore.
func (dataStore DataStore) Insert(key *Key, value string) {
	dataStore.time[key.Hash] = dataStore.calculateExpirationTime()
	dataStore.data[key.Hash] = value

	go dataStore.deleteAfterExpirationTime(dataStore.ttl, key.Hash)
}

func (dataStore DataStore) deleteAfterExpirationTime(timer time.Duration, hash [KeySize]byte) {

	select {
	case currTime := <-time.After(timer):
		differenceInTime := dataStore.time[hash].Sub(currTime)
		if differenceInTime <= 0 {
			dataStore.DeleteExpiredData(hash)
		} else {
			dataStore.deleteAfterExpirationTime(differenceInTime, hash)
		}
	}
}

// Get retrieves the value associated with a key from the DataStore.
func (dataStore DataStore) Get(key *Key) (string, error) {
	value, ok := dataStore.data[key.Hash]
	if !ok {
		return "", errors.New("key not found")
	}
	return value, nil
}

func (dataStore DataStore) GetTime(key *Key) (time.Time, error) {
	time, ok := dataStore.time[key.Hash]
	if !ok {
		return time, errors.New("key not found")
	}
	return time, nil
}

func (dataStore DataStore) calculateExpirationTime() time.Time {
	return time.Now().Add(dataStore.ttl)
}

func (dataStore DataStore) RefreshExpirationTime(key *Key) error {
	_, ok := dataStore.data[key.Hash]
	if !ok {
		return errors.New("key not found")
	}
	ttl := dataStore.calculateExpirationTime()
	dataStore.time[key.Hash] = ttl
	return nil
}

func (dataStore DataStore) DeleteExpiredData(dataToDelete [KeySize]byte) {
	value := dataStore.data[dataToDelete]
	delete(dataStore.time, dataToDelete)
	delete(dataStore.data, dataToDelete)
	logger.Log("The data object " + hex.EncodeToString(dataToDelete[:]) + " with the value " + value + " has been deleted due to the expired TTL.")
}
