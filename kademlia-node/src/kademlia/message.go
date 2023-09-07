package kademlia

type MessageType string

const (
	ERROR     MessageType = "ERROR"
	PING      MessageType = "PING"
	PONG      MessageType = "PONG"
	FIND_NODE MessageType = "FIND_NODE"
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
