package kademlia

import "errors"

type MessageType string

const (
	ERROR          MessageType = "ERROR"
	PING           MessageType = "PING"
	PONG           MessageType = "PONG"
	FIND_NODE      MessageType = "FIND_NODE"
	FIND_DATA      MessageType = "FIND_DATA"
	STORE          MessageType = "STORE"
	STORE_RESPONSE MessageType = "STORE_RESPONSE"
)

func (messageType MessageType) IsValid() error {
	switch messageType {
	case ERROR, PING, PONG, FIND_NODE, FIND_DATA, STORE, STORE_RESPONSE: // Add new messageTypes to the case, so it is seen as a valid type
		return nil
	}
	return errors.New("Invalid message type")
}

type Message struct {
	MessageType MessageType `json:"messageType"`
	Contact     Contact     `json:"contact"`
}

type Error struct {
	Message
}

func NewErrorMessage(contact Contact) Error {
	message := Message{
		MessageType: ERROR,
		Contact:     contact,
	}
	return Error{
		message,
	}
}

type Ping struct {
	Message
}

func NewPingMessage(contact Contact) Ping {
	message := Message{
		MessageType: PING,
		Contact:     contact,
	}
	return Ping{
		message,
	}
}

type Pong struct {
	Message
}

func NewPongMessage(contact Contact) Pong {
	message := Message{
		MessageType: PONG,
		Contact:     contact,
	}
	return Pong{
		message,
	}
}

type FindNode struct {
	Message
	FromAddress string `json:"fromAddress"`
	ID          *KademliaID
}

func NewFindNodeMessage(contact Contact, fromAddress string, id *KademliaID) FindNode {
	message := Message{
		MessageType: FIND_NODE,
		Contact:     contact,
	}
	return FindNode{
		message,
		fromAddress,
		id,
	}
}

type FindData struct {
	Message
	FromAddress string `json:"fromAddress"`
	ID          *KademliaID
	Key         *Key
}

func NewFindDataMessage(contact Contact, fromAddress string, id *KademliaID, key *Key) FindData {
	message := Message{
		MessageType: FIND_DATA,
		Contact:     contact,
	}
	return FindData{
		message,
		fromAddress,
		id,
		key,
	}
}

type Store struct {
	Message
	FromAddress string `json:"fromAddress"`
	Key         *Key
	ID          *KademliaID
	Value       string
}

func NewStoreMessage(contact Contact, fromAddress string, key *Key, id *KademliaID, value string) Store {
	message := Message{
		MessageType: STORE,
		Contact:     contact,
	}

	return Store{
		message,
		fromAddress,
		key,
		id,
		value,
	}
}

type StoreResponse struct {
	Message
	StoreSuccess bool `json:"storeSuccess"`
}

func NewStoreResponseMessage(contact Contact) StoreResponse {
	message := Message{
		MessageType: STORE_RESPONSE,
		Contact:     contact,
	}

	// change this value in message_handler
	storeSuccess := true

	return StoreResponse{
		message,
		storeSuccess,
	}

}
