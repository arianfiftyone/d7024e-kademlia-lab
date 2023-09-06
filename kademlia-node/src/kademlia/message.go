package kademlia

type MessageType string

const (
	ERROR MessageType = "ERROR"
	PING  MessageType = "PING"
	PONG  MessageType = "PONG"
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
