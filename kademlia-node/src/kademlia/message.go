package kademlia

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

type Message struct {
	MessageType MessageType `json:"messageType"`
}

type Error struct {
	Message
}

func NewErrorMessage() Error {
	message := Message{
		MessageType: ERROR,
	}
	return Error{
		message,
	}
}

type Ping struct {
	Message
	FromAddress string `json:"fromAddress"`
}

func NewPingMessage(fromAddress string) Ping {
	message := Message{
		MessageType: PING,
	}
	return Ping{
		message,
		fromAddress,
	}
}

type Pong struct {
	Message
	FromAddress string `json:"fromAddress"`
}

func NewAckPingMessage(fromAddress string) Pong {
	message := Message{
		MessageType: PONG,
	}
	return Pong{
		message,
		fromAddress,
	}
}

type FindNode struct {
	Message
	FromAddress string `json:"fromAddress"`
	ID          *KademliaID
}

func NewFindNodeMessage(fromAddress string, id *KademliaID) FindNode {
	message := Message{
		MessageType: FIND_NODE,
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

func NewFindDataMessage(fromAddress string, id *KademliaID, key *Key) FindData {
	message := Message{
		MessageType: FIND_DATA,
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

func NewStoreMessage(fromAddress string, key *Key, id *KademliaID, value string) Store {
	message := Message{
		MessageType: STORE,
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

func NewStoreResponseMessage() StoreResponse {
	message := Message{
		MessageType: STORE_RESPONSE,
	}

	// change this value in message_handler
	storeSuccess := true

	return StoreResponse{
		message,
		storeSuccess,
	}

}
