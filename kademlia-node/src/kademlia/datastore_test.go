package kademlia

import (
	"reflect"
	"testing"
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
	key := HashToKey(value)

	dataStore.Insert(key, value)

	if !reflect.DeepEqual(dataStore.data[key], value) {
		t.Errorf("Insert: Expected %v, got %v", value, dataStore.data[key])
	}
}

func TestInsertAndGet(t *testing.T) {
	dataStore := NewDataStore()

	value := "testValue"
	key := HashToKey(value)

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
	keyNotExisting := HashToKey(value2)
	_, err = dataStore.Get(keyNotExisting)
	if err == nil {
		t.Errorf("Get: Expected error for non-existent key, but got none")
	}
}
