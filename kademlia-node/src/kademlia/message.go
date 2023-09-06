package kademlia

type MessageType string

const (
	ERROR    MessageType = "ERROR"
	PING     MessageType = "PING"
	ACK_PING MessageType = "ACK_PING"
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

type AckPing struct {
	Message
	FromAddress string `json:"fromAddress"`
}

func NewAckPingMessage(fromAddress string) AckPing {
	message := Message{
		MessageType: ACK_PING,
	}
	return AckPing{
		message,
		fromAddress,
	}
}
