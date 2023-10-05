package kademlia

import (
	"fmt"
	"reflect"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
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

	if !reflect.DeepEqual(dataStore.data[key.Hash], value) {
		t.Errorf("Insert: Expected %v, got %v", value, dataStore.data[key.Hash])
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
	expectedTime := time.Now().Add(dataStore.ttl)

	fmt.Println(insertedTime)
	fmt.Println(expectedTime)

	// Check if the retrieved time is within a small tolerance (1 second) of the expected time
	if expectedTime.Sub((insertedTime)) > time.Millisecond*500 {
		t.Errorf("Expected time to be roughly %v, but got %v", expectedTime, insertedTime)
	}
}

func TestRefreshExpirationTime(t *testing.T) {
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

	fmt.Println(time.Now())
	delay := time.Second * 2
	time.Sleep(delay + time.Millisecond*100)

	dataStore.RefreshExpirationTime(key)

	newTime, err := dataStore.GetTime(key)
	if err != nil {
		t.Errorf("Expected no error, but got %v", err)
	}

	fmt.Println(insertedTime)
	fmt.Println(newTime)

	assert.Greater(t, newTime, insertedTime.Add(delay))
}
