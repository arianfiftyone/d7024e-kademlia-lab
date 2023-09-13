package kademlia

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net"
	"time"
)

type Network interface {
	Listen() error
	Send(ip string, port int, message []byte, timeOut time.Duration) ([]byte, error)
	SendPingMessage(from *Contact, contact *Contact) error
	SendFindContactMessage(from *Contact, contact *Contact, id *KademliaID) ([]Contact, error)
	SendFindDataMessage(from *Contact, contact *Contact, key *Key) ([]Contact, string, error)
	SendStoreMessage(from *Contact, contact *Contact, value string) bool
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
		log.Printf("Failed to listen for UDP packets: %v\n", err)
		return err
	}

	defer conn.Close()

	fmt.Printf("server listening %s:%d\n", network.Ip, network.Port)

	// //TEST
	// kademlia := &Kademlia{
	// 	network,
	// 	NewKademliaNode(network.Ip, network.Port),
	// }
	// contact := Test(kademlia)
	// list, _ := kademlia.LookupContact(contact)
	// fmt.Println(list)
	// //

	for {
		data := make([]byte, 1024)
		len, remote, err := conn.ReadFromUDP(data[:])
		if err != nil {
			log.Printf("Failed to read from UDP: %v\n", err)
			return err
		}

		go func(myConn *net.UDPConn) {
			response, err := network.MessageHandler.HandleMessage(data[:len])
			if err != nil {
				log.Printf("Failed to handle response message: %v\n", err)
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
		log.Printf("Failed to connect via UDP: %v\n", err)
		return nil, err
	}

	// Send a message to the server
	_, err = conn.Write(message)
	if err != nil {
		log.Printf("Failed to send a message to the server: %v\n", err)
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
		log.Printf("Failed to send ping message to the server: %v\n", err)
		return err
	}

	response, err := network.Send(contact.Ip, contact.Port, bytes, time.Second*3)
	if err != nil {
		log.Printf("Ping failed: %v\n", err)
		return err
	}
	var message Message
	errUnmarshal := json.Unmarshal(response, &message)
	if errUnmarshal != nil || message.MessageType != PONG {
		log.Printf("Ping failed: %v\n", errUnmarshal)
		return errUnmarshal
	}

	var pong Pong

	errUnmarshalAckPing := json.Unmarshal(response, &pong)
	if errUnmarshalAckPing != nil {
		log.Printf("Ping failed: %v\n", errUnmarshalAckPing)
		return errUnmarshalAckPing
	}

	fmt.Println(pong.From.Ip + " acknowledged your ping")
	return nil

}

func (network *NetworkImplementation) SendFindContactMessage(from *Contact, contact *Contact, id *KademliaID) ([]Contact, error) {
	findN := NewFindNodeMessage(*from, id)
	bytes, err := json.Marshal(findN)
	if err != nil {
		log.Printf("Error when marshaling `findN`: %v\n", err)
		return nil, err
	}

	response, err := network.Send(contact.Ip, contact.Port, bytes, time.Second*3)
	if err != nil {
		log.Printf("Find node failed: %v\n", err)
		return nil, err
	}

	var message Message
	errUnmarshal := json.Unmarshal(response, &message)
	if errUnmarshal != nil || message.MessageType != FOUND_CONTACTS {
		log.Printf("Failed to find contacts: %v\n", errUnmarshal)
		return nil, errUnmarshal
	}

	var arrayOfContacts FoundContacts
	errUnmarshalFoundContacts := json.Unmarshal(response, &arrayOfContacts)
	if errUnmarshalFoundContacts != nil {
		log.Printf("Failed to find contacts: %v\n", errUnmarshalFoundContacts)
		return nil, errUnmarshalFoundContacts
	}

	return arrayOfContacts.Contacts, nil
}

func (network *NetworkImplementation) SendFindDataMessage(from *Contact, contact *Contact, key *Key) ([]Contact, string, error) {
	findData := NewFindDataMessage(*from, network.Ip, contact.ID, key)
	bytes, err := json.Marshal(findData)
	if err != nil {
		return nil, "", err
	}

	response, err := network.Send(contact.Ip, contact.Port, bytes, time.Second*3)
	if err != nil {
		log.Printf("Find data failed: %v\n", err)
		return nil, "", err
	}

	var data string
	json.Unmarshal(response, &data)
	if data == "" {
		var arrayOfContacts []Contact
		json.Unmarshal(response, &arrayOfContacts)
		return arrayOfContacts, "", nil
	} else {
		return nil, data, nil
	}

}

func (network *NetworkImplementation) SendStoreMessage(from *Contact, contact *Contact, value string) bool {
	key := HashToKey(value)
	store := NewStoreMessage(*from, network.Ip, key, contact.ID, value)
	bytes, err := json.Marshal(store)
	if err != nil {
		log.Printf("Error when marshaling `store` message: %v\n", err)
		return false
	}

	response, err := network.Send(contact.Ip, contact.Port, bytes, time.Second*3)
	if err != nil {
		log.Printf("Store failed: %v\n", err)
		return false
	}

	var storeResponse StoreResponse
	err = json.Unmarshal(response, &storeResponse)
	if err != nil {
		log.Printf("Error when unmarshaling `storeResponse` message: %v\n", err)
		return false
	}

	return storeResponse.StoreSuccess

}
