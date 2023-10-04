package kademlia

import (
	"reflect"
	"testing"
	"time"
)

// TestNewDataStore tests the NewDataStore function.
func TestNewDataStore(t *testing.T) {
	dataStore := NewDataStore()

	if dataStore.data == nil {
		t.Error("NewDataStore: Data map should be initialized.")
	}
}

// TestInsert tests the Insert method.
func TestInsert(t *testing.T) {
	dataStore := NewDataStore()

	value := string("testValue")
	key := NewKey(value)

	dataStore.Insert(key, value)

	if !reflect.DeepEqual(dataStore.data[key.Hash], value) {
		t.Errorf("Insert: Expected %v, got %v", value, dataStore.data[key.Hash])
	}
}

func TestInsertAndGet(t *testing.T) {
	dataStore := NewDataStore()

	value := "testValue"
	key := NewKey(value)

	dataStore.Insert(key, value)

	retrievedValue, err := dataStore.Get(key)
	if err != nil {
		t.Errorf("Get: Unexpected error: %v", err)
	}

	if retrievedValue != value {
		t.Errorf("Get: Expected %v, got %v", value, retrievedValue)
	}

	// Test case for a non-existent key
	value2 := "testValue2" // refers to key that has not been previously inserted into data store
	keyNotExisting := NewKey(value2)
	_, err = dataStore.Get(keyNotExisting)
	if err == nil {
		t.Errorf("Get: Expected error for non-existent key, but got none")
	}
}

func TestInsertAndGetTime(t *testing.T) {
	dataStore := NewDataStore()

	// Insert a key-value pair
	value := "testValue"
	key := HashToKey(value)
	dataStore.Insert(key, value)

	// Retrieve the time associated with the key
	insertedTime, err := dataStore.GetTime(key)
	if err != nil {
		t.Errorf("Expected no error, but got %v", err)
	}

	// Calculate the expected expiration time
	expectedTime := time.Now().Unix() + dataStore.ttl

	// Check if the retrieved time is within a small tolerance (1 second) of the expected time
	if insertedTime < expectedTime-1 || insertedTime > expectedTime+1 {
		t.Errorf("Expected time to be roughly %v, but got %v", expectedTime, insertedTime)
	}
}
