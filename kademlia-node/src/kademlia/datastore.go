package kademlia

import (
	"errors"
)

// DataStore represents a key-value data store.
type DataStore struct {
	data map[*Key]string // Map to store key-value pairs.
}

// NewDataStore initializes a new DataStore instance.
func NewDataStore() DataStore {
	dataStore := DataStore{}
	dataStore.data = make(map[*Key]string)
	return dataStore
}

// Insert inserts a key-value pair into the DataStore.
func (dataStore DataStore) Insert(key *Key, value string) {
	dataStore.data[key] = value
}

// Get retrieves the value associated with a key from the DataStore.
func (dataStore DataStore) Get(key *Key) (string, error) {
	value, ok := dataStore.data[key]
	if !ok {
		return "", errors.New("key not found")
	}
	return value, nil
}
