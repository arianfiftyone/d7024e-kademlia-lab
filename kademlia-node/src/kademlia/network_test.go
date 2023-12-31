package kademlia

import (
	"encoding/json"
	"fmt"
	"net"
	"strconv"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

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

type MockMessageHandler struct{}

func (messageHandler *MockMessageHandler) HandleMessage(rawMessage []byte) ([]byte, error) {
	var message Message
	json.Unmarshal(rawMessage, &message)
	fmt.Println(message.MessageType)
	if message.MessageType != "" {
		ok := NewOkMessage("")
		bytes, _ := json.Marshal(ok)
		return bytes, nil

	} else {
		return make([]byte, 0), nil

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
		// Read from the connection
		data := make([]byte, 1024)
		len, _, err := conn.ReadFromUDP(data[:])
		if err != nil {
			return
		}
		responseChannel <- data[:len]

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
	network := NetworkImplementation{
		"localhost",
		3000,
		&MockMessageHandler{},
	}

	go network.Listen()
	time.Sleep(time.Second)

	ping := NewPingMessage(NewContact(GenerateNewKademliaID("FFFFFFFF00000000000000000000000000000000"), network.Ip, network.Port))
	bytes, _ := json.Marshal(ping)
	mockSend(t, network.Ip, network.Port, bytes, time.Second*3)
}

func TestClient(t *testing.T) {
	network := NetworkImplementation{
		"localhost",
		4000,
		&MockMessageHandler{},
	}

	go network.Listen()
	time.Sleep(time.Second)

	ping := NewPingMessage(NewContact(GenerateNewKademliaID("FFFFFFFF00000000000000000000000000000000"), network.Ip, network.Port))
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

func (mockMessageHandler *MockMessageHandler2) HandleMessage(rawMessage []byte) ([]byte, error) {
	var findN FindNode

	json.Unmarshal(rawMessage, &findN)
	if findN.MessageType == FIND_NODE {
		var arrayC [1]Contact
		arrayC[0] = NewContact(GenerateNewKademliaID("FFFFFFFF00000000000000000000000000000000"), "localhost", 8001)
		bytes, _ := json.Marshal(NewFoundContactsMessage(findN.From, arrayC[:]))
		return bytes, nil

	} else if findN.MessageType == FIND_DATA {
		var arrayC [1]Contact
		arrayC[0] = NewContact(GenerateNewKademliaID("FFFFFFFF00000000000000000000000000000000"), "localhost", 8001)
		bytes, _ := json.Marshal(NewFoundContactsMessage(findN.From, arrayC[:]))
		return bytes, nil

	} else {
		return make([]byte, 0), nil

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
	mockNetwork := &NetworkImplementation{
		Ip:             "127.0.0.1",
		Port:           5000,
		MessageHandler: &MockMessageHandler2{},
	}

	go mockNetwork.Listen()
	time.Sleep(time.Second)

	from := NewContact(GenerateNewKademliaID("FFFFFFFF00000000000000000000000000000000"), mockNetwork.Ip, mockNetwork.Port)
	response, _ := mockNetwork.SendFindContactMessage(&from, &mockContact, mockContact.ID)
	fmt.Println("First contact: " + response[0].ID.String())
	assert.Equal(t, response[0], NewContact(GenerateNewKademliaID("FFFFFFFF00000000000000000000000000000000"), "localhost", 8001))
}

func TestSendNodeDataMessage(t *testing.T) {

	bootstrap := CreateMockedKademlia(GenerateNewKademliaID("FFFFFFFF00000000000000000000000000000000"), "127.0.0.1", 7010)

	value := "data"
	key := NewKey(value)

	bootstrap.KademliaNode.GetDataStore().Insert(key, value)
	go bootstrap.Start()
	time.Sleep(time.Second)
	_, str, _ := bootstrap.Network.SendFindDataMessage(&bootstrap.KademliaNode.GetRoutingTable().Me, &bootstrap.KademliaNode.GetRoutingTable().Me, key)

	assert.Equal(t, str, value)
}

func TestSendNodeDataMessageNoData(t *testing.T) {

	bootstrap := CreateMockedKademlia(GenerateNewKademliaID("FFFFFFFF00000000000000000000000000000000"), "127.0.0.1", 7020)

	kademlia1 := CreateMockedKademlia(GenerateNewKademliaID("0000000000000000000000000000000000000001"), "127.0.0.1", 7011)
	kademlia2 := CreateMockedKademlia(GenerateNewKademliaID("0000000000000000000000000000000000000002"), "127.0.0.1", 7012)

	bootstrap.KademliaNode.GetRoutingTable().AddContact(kademlia1.KademliaNode.GetRoutingTable().Me)
	bootstrap.KademliaNode.GetRoutingTable().AddContact(kademlia2.KademliaNode.GetRoutingTable().Me)

	go bootstrap.Start()
	time.Sleep(time.Second)
	value := "data"
	key := NewKey(value)
	contacts, _, _ := bootstrap.Network.SendFindDataMessage(&bootstrap.KademliaNode.GetRoutingTable().Me, &bootstrap.KademliaNode.GetRoutingTable().Me, key)

	kClosest := bootstrap.KademliaNode.GetRoutingTable().FindClosestContacts(bootstrap.KademliaNode.GetRoutingTable().Me.ID, NumberOfClosestNodesToRetrieved)
	doesContainAll := bootstrap.FirstSetContainsAllContactsOfSecondSet(kClosest, contacts)
	assert.True(t, doesContainAll)
}

type MockSlowMessageHandler struct{}

func (messageHandler *MockSlowMessageHandler) HandleMessage(rawMessage []byte) ([]byte, error) {
	var message Message
	json.Unmarshal(rawMessage, &message)
	fmt.Println(message.MessageType)
	time.Sleep(time.Second * 5)
	if message.MessageType != "" {
		ok := NewOkMessage("")
		bytes, _ := json.Marshal(ok)
		return bytes, nil

	} else {
		return make([]byte, 0), nil

	}
}
func TestTimeout(t *testing.T) {
	network := NetworkImplementation{
		"localhost",
		8000,
		&MockSlowMessageHandler{},
	}

	go network.Listen()
	time.Sleep(time.Second)

	ping := NewPingMessage(NewContact(GenerateNewKademliaID("FFFFFFFF00000000000000000000000000000000"), network.Ip, network.Port))
	bytes, _ := json.Marshal(ping)
	_, err := network.Send(network.Ip, network.Port, bytes, time.Second*3)
	fmt.Println(err)
	assert.EqualError(t, err, "time out error")

}

type MockMessageHandlerConcurrentSend struct{}

func (messageHandler *MockMessageHandlerConcurrentSend) HandleMessage(rawMessage []byte) ([]byte, error) {
	var message Message
	json.Unmarshal(rawMessage, &message)
	if message.MessageType == OK_MESSAGE {
		var okMessage OkMessage
		json.Unmarshal(rawMessage, &okMessage)
		fmt.Println(okMessage)
		bytes, _ := json.Marshal(okMessage)
		return bytes, nil

	} else {
		return make([]byte, 0), nil

	}
}

func (network *NetworkImplementation) sendOkMessage(t *testing.T, startNumber int) {
	debugMessage := "Start number: " + strconv.Itoa(startNumber)
	outMessage := NewOkMessage("Start number: " + strconv.Itoa(startNumber))
	bytes, _ := json.Marshal(outMessage)

	response, err := network.Send(network.Ip, network.Port, bytes, time.Second*3)

	if err != nil {
		assert.Fail(t, "Error sending message!: "+err.Error())
	}

	var inMessage OkMessage
	json.Unmarshal(response, &inMessage)
	fmt.Println("Expected message: " + debugMessage + ", In message: " + inMessage.DebugMessage)

	assert.True(t, inMessage.DebugMessage == debugMessage, "Communication message must be: "+debugMessage)
}

func TestConcurrentSends(t *testing.T) {
	network := NetworkImplementation{
		"localhost",
		7000,
		&MockMessageHandlerConcurrentSend{},
	}

	go network.Listen()

	time.Sleep(time.Second)

	i := 1
	for i < 10 {
		go network.sendOkMessage(t, i)
		i += 1
	}

	time.Sleep(time.Second)

}

type MockMessageHandlerRefresh struct {
}

func (mockMessageHandler *MockMessageHandlerRefresh) HandleMessage(rawMessage []byte) ([]byte, error) {
	var refreshExpirationTime RefreshExpirationTime

	json.Unmarshal(rawMessage, &refreshExpirationTime)
	if refreshExpirationTime.MessageType == REFRESH_EXPIRATION_TIME {
		bytes, _ := json.Marshal(NewExpirationTimeHasBeenRefreshedMessage(NewContact(GenerateNewKademliaID("FFFFFFFF00000000000000000000000000000000"), "localhost", 8001)))
		return bytes, nil

	} else {
		return make([]byte, 0), nil

	}
}

func TestSendRefreshMessage(t *testing.T) {
	// Create a mock Contact for testing
	mockContact := Contact{
		ID:       NewRandomKademliaID(),
		Ip:       "127.0.0.1",
		Port:     30000,
		distance: nil,
	}

	// Create a mock Network instance
	mockNetwork := &NetworkImplementation{
		Ip:             "127.0.0.1",
		Port:           30000,
		MessageHandler: &MockMessageHandlerRefresh{},
	}

	key := NewKey("test")

	go mockNetwork.Listen()
	time.Sleep(time.Second)

	from := NewContact(GenerateNewKademliaID("FFFFFFFF00000000000000000000000000000000"), mockNetwork.Ip, mockNetwork.Port)
	response := mockNetwork.SendRefreshExpirationTimeMessage(&from, &mockContact, key)
	assert.True(t, response)
}
