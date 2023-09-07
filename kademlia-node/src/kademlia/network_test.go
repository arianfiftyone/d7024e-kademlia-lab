package kademlia

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

type MockMessageHandler struct{}

const OK_MESSAGE MessageType = "OK"

type OkMessage struct {
	Message
	DebugMessage string
}

func NewOkMessage(debugMessage string) OkMessage {
	message := Message{
		MessageType: OK_MESSAGE,
	}
	return OkMessage{
		message,
		debugMessage,
	}
}

func (messageHandler *MockMessageHandler) HandleMessage(rawMessage []byte) []byte {
	var message Message
	json.Unmarshal(rawMessage, &message)
	fmt.Println(message.MessageType)
	if message.MessageType != "" {
		ok := NewOkMessage("")
		bytes, _ := json.Marshal(ok)
		return bytes

	} else {
		return make([]byte, 0)

	}
}

func mockSend(t *testing.T, ip string, port int, message []byte, timeOut time.Duration) {
	conn, _ := net.DialUDP("udp", nil, &net.UDPAddr{
		IP:   net.ParseIP(ip),
		Port: port,
	})

	// Send a message to the server
	_, _ = conn.Write(message)

	responseChannel := make(chan []byte)
	go func() {
		// Read from the connection untill a new line is send
		data, _ := bufio.NewReader(conn).ReadString('\n')
		responseChannel <- []byte(data)

	}()

	select {
	case response := <-responseChannel:
		var message OkMessage
		json.Unmarshal(response, &message)
		fmt.Println(message)
		assert.True(t, message.MessageType == OK_MESSAGE, "Communication message must be an OK message!")

	case <-time.After(timeOut):
		assert.Fail(t, "Communication failed!")

	}
}

func TestServer(t *testing.T) {
	network := Network{
		"localhost",
		3000,
		&MockMessageHandler{},
	}

	go network.Listen()
	time.Sleep(time.Second)

	ping := NewPingMessage(network.Ip)
	bytes, _ := json.Marshal(ping)
	mockSend(t, network.Ip, network.Port, bytes, time.Second*3)
}

func TestClient(t *testing.T) {
	network := Network{
		"localhost",
		3000,
		&MockMessageHandler{},
	}

	go network.Listen()
	time.Sleep(time.Second)

	ping := NewPingMessage(network.Ip)
	bytes, _ := json.Marshal(ping)
	response, err := network.Send(network.Ip, network.Port, bytes, time.Second*3)

	if err != nil {
		assert.Fail(t, "Error sending message!: "+err.Error())
	}

	var message OkMessage
	json.Unmarshal(response, &message)
	fmt.Println(message)
	assert.True(t, message.MessageType == OK_MESSAGE, "Communication message must be an OK message!")
}

type MockMessageHandler2 struct {
}

func (mockMessageHandler *MockMessageHandler2) HandleMessage(rawMessage []byte) []byte {
	var findN FindNode

	json.Unmarshal(rawMessage, &findN)
	if findN.MessageType == FIND_NODE {
		var arrayC [1]Contact
		arrayC[0] = NewContact(NewKademliaID("FFFFFFFF00000000000000000000000000000000"), "localhost", 8001)
		bytes, _ := json.Marshal(arrayC)
		return bytes

	} else {
		return make([]byte, 0)

	}
}

func TestSendNodeContactMessage(t *testing.T) {
	// Create a mock Contact for testing
	mockContact := Contact{
		ID:       NewRandomKademliaID(),
		Ip:       "127.0.0.1",
		Port:     5000,
		distance: nil,
	}

	// Create a mock Network instance
	mockNetwork := &Network{
		Ip:             "127.0.0.1",
		Port:           5000,
		MessageHandler: &MockMessageHandler2{},
	}

	go mockNetwork.Listen()
	time.Sleep(time.Second * 5)

	response, _ := mockNetwork.SendFindContactMessage(&mockContact)
	println(response)
	assert.Equal(t, response[0], NewContact(NewKademliaID("FFFFFFFF00000000000000000000000000000000"), "localhost", 8001))
}
