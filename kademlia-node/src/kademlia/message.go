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
	FOUND_CONTACTS MessageType = "FOUND_CONTACTS"
)

func (messageType MessageType) IsValid() error {
	switch messageType {
	case ERROR, PING, PONG, FIND_NODE, FIND_DATA, STORE, STORE_RESPONSE, FOUND_CONTACTS: // Add new messageTypes to the case, so it is seen as a valid type
		return nil
	}
	return errors.New("Invalid message type")
}

type Message struct {
	MessageType MessageType `json:"messageType"`
	From        Contact     `json:"contact"`
}

type Error struct {
	Message
}

func NewErrorMessage(from Contact) Error {
	message := Message{
		MessageType: ERROR,
		From:        from,
	}
	return Error{
		message,
	}
}

type Ping struct {
	Message
}

func NewPingMessage(from Contact) Ping {
	message := Message{
		MessageType: PING,
		From:        from,
	}
	return Ping{
		message,
	}
}

type Pong struct {
	Message
}

func NewPongMessage(from Contact) Pong {
	message := Message{
		MessageType: PONG,
		From:        from,
	}
	return Pong{
		message,
	}
}

type FindNode struct {
	Message
	ID *KademliaID
}

func NewFindNodeMessage(from Contact, id *KademliaID) FindNode {
	message := Message{
		MessageType: FIND_NODE,
		From:        from,
	}
	return FindNode{
		message,
		id,
	}
}

type FindData struct {
	Message
	FromAddress string `json:"fromAddress"`
	ID          *KademliaID
	Key         *Key
}

func NewFindDataMessage(from Contact, fromAddress string, id *KademliaID, key *Key) FindData {
	message := Message{
		MessageType: FIND_DATA,
		From:        from,
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

func NewStoreMessage(from Contact, fromAddress string, key *Key, id *KademliaID, value string) Store {
	message := Message{
		MessageType: STORE,
		From:        from,
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

func NewStoreResponseMessage(from Contact) StoreResponse {
	message := Message{
		MessageType: STORE_RESPONSE,
		From:        from,
	}

	// change this value in message_handler
	storeSuccess := true

	return StoreResponse{
		message,
		storeSuccess,
	}

}

type FoundContacts struct {
	Message
	Contacts []Contact `json:"contacts"`
}

func NewFoundContactsMessage(from Contact, contacts []Contact) FoundContacts {
	message := Message{
		MessageType: FOUND_CONTACTS,
		From:        from,
	}

	return FoundContacts{
		message,
		contacts,
	}

}
