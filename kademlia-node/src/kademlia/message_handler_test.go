package kademlia

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

type KademliaNodeMock struct {
	me        *Contact
	DataStore *DataStore
}

func (kademliaNode *KademliaNodeMock) setNetwork(network Network) {

}

func (kademliaNode *KademliaNodeMock) GetRoutingTable() *RoutingTable {
	return NewRoutingTable(*kademliaNode.me)
}

func (kademliaNode *KademliaNodeMock) GetDataStore() *DataStore {
	return kademliaNode.DataStore
}

func (kademliaNode *KademliaNodeMock) updateRoutingTable(contact Contact) {

}

func TestPongMessage(t *testing.T) {
	contact := NewContact(NewRandomKademliaID(), "127.0.0.1", 80)
	messageHandler := &MessageHandlerImplementation{
		kademliaNode: &KademliaNodeMock{
			me: &contact,
		},
	}

	pong := NewPongMessage(NewContact(NewRandomKademliaID(), "127.0.0.1", 80))
	bytes, err := json.Marshal(pong)
	if err != nil {
		assert.Fail(t, err.Error())
	}
	response, err := messageHandler.HandleMessage(bytes)
	if err != nil {
		assert.Fail(t, err.Error())
	}
	var message Message
	errUnmarshal := json.Unmarshal(response, &message)
	if errUnmarshal != nil {
		assert.Fail(t, errUnmarshal.Error())
	}
	assert.Equal(t, ERROR, message.MessageType)

}

func TestPingMessage(t *testing.T) {
	contact := NewContact(NewRandomKademliaID(), "127.0.0.1", 80)
	messageHandler := &MessageHandlerImplementation{
		kademliaNode: &KademliaNodeMock{
			me: &contact,
		},
	}

	ping := NewPingMessage(NewContact(NewRandomKademliaID(), "127.0.0.1", 80))
	bytes, err := json.Marshal(ping)
	if err != nil {
		assert.Fail(t, err.Error())
	}
	response, err := messageHandler.HandleMessage(bytes)
	if err != nil {
		assert.Fail(t, err.Error())
	}
	var message Message
	errUnmarshal := json.Unmarshal(response, &message)
	if errUnmarshal != nil {
		assert.Fail(t, errUnmarshal.Error())
	}
	assert.Equal(t, PONG, message.MessageType)

}

func TestFindNodeMessage(t *testing.T) {
	contact := NewContact(NewRandomKademliaID(), "127.0.0.1", 80)
	messageHandler := &MessageHandlerImplementation{
		kademliaNode: &KademliaNodeMock{
			me: &contact,
		},
	}

	target := NewRandomKademliaID()
	findNode := NewFindNodeMessage(NewContact(NewRandomKademliaID(), "127.0.0.1", 80), target)
	bytes, err := json.Marshal(findNode)
	if err != nil {
		assert.Fail(t, err.Error())
	}
	response, err := messageHandler.HandleMessage(bytes)
	if err != nil {
		assert.Fail(t, err.Error())
	}
	var message Message
	errUnmarshal := json.Unmarshal(response, &message)
	if errUnmarshal != nil {
		assert.Fail(t, errUnmarshal.Error())
	}
	assert.Equal(t, FOUND_CONTACTS, message.MessageType)

}

func TestFindDataMessage(t *testing.T) {
	contact := NewContact(NewRandomKademliaID(), "127.0.0.1", 80)
	messageHandler := &MessageHandlerImplementation{
		kademliaNode: &KademliaNodeMock{
			me:        &contact,
			DataStore: &DataStore{},
		},
	}

	target := NewRandomKademliaID()
	findData := NewFindDataMessage(NewContact(NewRandomKademliaID(), "127.0.0.1", 80), GetKeyRepresentationOfKademliaId(target))
	bytes, err := json.Marshal(findData)
	if err != nil {
		assert.Fail(t, err.Error())
	}
	response, err := messageHandler.HandleMessage(bytes)
	if err != nil {
		assert.Fail(t, err.Error())
	}
	var message Message
	errUnmarshal := json.Unmarshal(response, &message)
	if errUnmarshal != nil {
		assert.Fail(t, errUnmarshal.Error())
	}
	assert.Equal(t, FOUND_DATA, message.MessageType)

}

func TestFindDataMessageDataExists(t *testing.T) {
	contact := NewContact(NewRandomKademliaID(), "127.0.0.1", 80)

	target := NewRandomKademliaID()
	dataStore := NewDataStore()
	value := "test"
	dataStore.Insert(GetKeyRepresentationOfKademliaId(target), value)

	messageHandler := &MessageHandlerImplementation{
		kademliaNode: &KademliaNodeMock{
			me:        &contact,
			DataStore: &dataStore,
		},
	}

	findData := NewFindDataMessage(NewContact(NewRandomKademliaID(), "127.0.0.1", 80), GetKeyRepresentationOfKademliaId(target))
	bytes, err := json.Marshal(findData)
	if err != nil {
		assert.Fail(t, err.Error())
	}
	response, err := messageHandler.HandleMessage(bytes)
	if err != nil {
		assert.Fail(t, err.Error())
	}
	var message Message
	errUnmarshal := json.Unmarshal(response, &message)
	if errUnmarshal != nil {
		assert.Fail(t, errUnmarshal.Error())
	}
	assert.Equal(t, FOUND_DATA, message.MessageType)

	var data FoundData
	errUnmarshalFoundData := json.Unmarshal(response, &data)
	if errUnmarshalFoundData != nil {
		assert.Fail(t, errUnmarshal.Error())

	}

	assert.Equal(t, value, data.Value)

}

func TestStoreMessage(t *testing.T) {
	contact := NewContact(NewRandomKademliaID(), "127.0.0.1", 80)
	dataStore := NewDataStore()
	value := "test"

	messageHandler := &MessageHandlerImplementation{
		kademliaNode: &KademliaNodeMock{
			me:        &contact,
			DataStore: &dataStore,
		},
	}

	target := NewRandomKademliaID()
	store := NewStoreMessage(NewContact(NewRandomKademliaID(), "127.0.0.1", 80), GetKeyRepresentationOfKademliaId(target), value)
	bytes, err := json.Marshal(store)
	if err != nil {
		assert.Fail(t, err.Error())
	}
	response, err := messageHandler.HandleMessage(bytes)
	if err != nil {
		assert.Fail(t, err.Error())
	}
	var message Message
	errUnmarshal := json.Unmarshal(response, &message)
	if errUnmarshal != nil {
		assert.Fail(t, errUnmarshal.Error())
	}
	assert.Equal(t, STORE_RESPONSE, message.MessageType)

}
