package kademlia

import (
	"encoding/json"
	"errors"
	"net"
	"strconv"
	"time"

	"github.com/arianfiftyone/src/logger"
)

type Network interface {
	Listen() error
	Send(ip string, port int, message []byte, timeOut time.Duration) ([]byte, error)
	SendPingMessage(from *Contact, contact *Contact) error
	SendFindContactMessage(from *Contact, contact *Contact, id *KademliaID) ([]Contact, error)
	SendFindDataMessage(from *Contact, contact *Contact, key *Key) ([]Contact, string, error)
	SendStoreMessage(from *Contact, contact *Contact, key *Key, value string) bool
}

type NetworkImplementation struct {
	Ip             string
	Port           int
	MessageHandler MessageHandler
}

func (network *NetworkImplementation) Listen() error {
	// listen to incoming udp packets
	conn, err := net.ListenUDP("udp", &net.UDPAddr{
		IP:   net.ParseIP(network.Ip),
		Port: network.Port,
	})
	if err != nil {
		logger.Log("Failed to listen for UDP packets: " + err.Error())
		return err
	}

	defer conn.Close()

	logger.Log("Server listening " + network.Ip + ":" + strconv.Itoa(network.Port))

	for {
		data := make([]byte, 1024)
		len, remote, err := conn.ReadFromUDP(data[:])
		if err != nil {
			logger.Log("Failed to read from UDP: " + err.Error())
			return err
		}

		go func(myConn *net.UDPConn) {
			response, err := network.MessageHandler.HandleMessage(data[:len])
			if err != nil {
				logger.Log("Failed to handle response message: " + err.Error())
				return
			}
			myConn.WriteToUDP(response, remote)

		}(conn)

	}

}

func (network *NetworkImplementation) Send(ip string, port int, message []byte, timeOut time.Duration) ([]byte, error) {
	conn, err := net.DialUDP("udp", nil, &net.UDPAddr{
		IP:   net.ParseIP(ip),
		Port: port,
	})

	if err != nil {
		logger.Log("Failed to connect via UDP: " + err.Error())
		return nil, err
	}

	// Send a message to the server
	_, err = conn.Write(message)
	if err != nil {
		logger.Log("Failed to send a message to the server: " + err.Error())
		return nil, err
	}

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
		if _, err := network.MessageHandler.HandleMessage(response); err != nil {
			return nil, err
		} else {
			return response, nil
		}
	case <-time.After(timeOut):
		return nil, errors.New("time out error")

	}

}

func (network *NetworkImplementation) SendPingMessage(from *Contact, contact *Contact) error {
	ping := NewPingMessage(*from)
	bytes, err := json.Marshal(ping)
	if err != nil {
		logger.Log("Failed to send a ping message to the server: " + err.Error())
		return err
	}

	response, err := network.Send(contact.Ip, contact.Port, bytes, time.Second*3)
	if err != nil {
		logger.Log("Ping failed: " + err.Error())
		return err
	}
	var message Message
	errUnmarshal := json.Unmarshal(response, &message)
	if errUnmarshal != nil || message.MessageType != PONG {
		logger.Log("Ping failed: " + errUnmarshal.Error())
		return errUnmarshal
	}

	var pong Pong

	errUnmarshalAckPing := json.Unmarshal(response, &pong)
	if errUnmarshalAckPing != nil {
		logger.Log("Ping failed: " + errUnmarshalAckPing.Error())
		return errUnmarshalAckPing
	}

	logger.Log(pong.From.Ip + " acknowledged your ping")
	return nil

}

func (network *NetworkImplementation) SendFindContactMessage(from *Contact, contact *Contact, id *KademliaID) ([]Contact, error) {
	findN := NewFindNodeMessage(*from, id)
	bytes, err := json.Marshal(findN)
	if err != nil {
		logger.Log("Error when marshaling `findN`: " + err.Error())
		return nil, err
	}

	response, err := network.Send(contact.Ip, contact.Port, bytes, time.Second*3)
	if err != nil {
		logger.Log("Find node failed: " + err.Error())
		return nil, err
	}

	var message Message
	errUnmarshal := json.Unmarshal(response, &message)
	if errUnmarshal != nil || message.MessageType != FOUND_CONTACTS {
		logger.Log("Find contact failed: " + errUnmarshal.Error())
		return nil, errUnmarshal
	}

	var arrayOfContacts FoundContacts
	errUnmarshalFoundContacts := json.Unmarshal(response, &arrayOfContacts)
	if errUnmarshalFoundContacts != nil {
		logger.Log("Error when unmarshaling 'foundContacts' message: " + errUnmarshalFoundContacts.Error())
		return nil, errUnmarshalFoundContacts
	}

	return arrayOfContacts.Contacts, nil
}

func (network *NetworkImplementation) SendFindDataMessage(from *Contact, contact *Contact, key *Key) ([]Contact, string, error) {
	findData := NewFindDataMessage(*from, key)
	bytes, err := json.Marshal(findData)
	if err != nil {
		return nil, "", err
	}

	response, err := network.Send(contact.Ip, contact.Port, bytes, time.Second*3)
	if err != nil {
		logger.Log("Find data failed: " + err.Error())
		return nil, "", err
	}

	var message Message
	errUnmarshal := json.Unmarshal(response, &message)
	if errUnmarshal != nil || message.MessageType != FOUND_DATA {
		logger.Log("Find data failed: " + errUnmarshal.Error())
		return nil, "", errUnmarshal
	}

	var data FoundData
	errUnmarshalFoundData := json.Unmarshal(response, &data)
	if errUnmarshalFoundData != nil {
		logger.Log("Error when unmarshaling 'foundData' message: " + errUnmarshalFoundData.Error())
		return nil, "", errUnmarshalFoundData
	}

	json.Unmarshal(response, &data)
	if data.Value == "" {
		return data.Contacts, "", nil
	} else {
		return nil, data.Value, nil
	}

}

func (network *NetworkImplementation) SendStoreMessage(from *Contact, contact *Contact, key *Key, value string) bool {
	store := NewStoreMessage(*from, key, value)
	bytes, err := json.Marshal(store)
	if err != nil {
		logger.Log("Error when marshaling `store` message: " + err.Error())
		return false
	}

	response, err := network.Send(contact.Ip, contact.Port, bytes, time.Second*3)
	if err != nil {
		logger.Log("Store failed: " + err.Error())
		return false
	}

	var storeResponse StoreResponse
	err = json.Unmarshal(response, &storeResponse)
	if err != nil {
		logger.Log("Error when unmarshaling `storeResponse` message: " + err.Error())
		return false
	}

	return storeResponse.StoreSuccess

}
