package kademlia

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"net"
	"time"
)

type Network struct {
	Ip             string
	Port           int
	MessageHandler MessageHandler
}

func (network *Network) Listen() error {
	// listen to incoming udp packets
	conn, err := net.ListenUDP("udp", &net.UDPAddr{
		IP:   net.ParseIP(network.Ip),
		Port: network.Port,
	})
	if err != nil {
		return err
	}

	defer conn.Close()

	fmt.Printf("server listening %s:%d\n", network.Ip, network.Port)

	contact := NewContact(NewRandomKademliaID(), network.Ip, network.Port)
	go func() {
		network.SendFindDataMessage(&contact, HashToKey("value"))
	}()

	for {
		data := make([]byte, 1024)
		len, remote, err := conn.ReadFromUDP(data[:])
		if err != nil {
			return err
		}

		go func(myConn *net.UDPConn) {
			response := network.MessageHandler.HandleMessage(data[:len])
			if err != nil {
				fmt.Println(err)
				return
			}
			myConn.WriteToUDP([]byte(string(response)+"\n"), remote)

		}(conn)

	}

}

func (network *Network) Send(ip string, port int, message []byte, timeOut time.Duration) ([]byte, error) {
	conn, err := net.DialUDP("udp", nil, &net.UDPAddr{
		IP:   net.ParseIP(ip),
		Port: port,
	})

	if err != nil {
		return nil, err
	}

	// Send a message to the server
	_, err = conn.Write(message)
	if err != nil {
		fmt.Println(err)
	}

	responseChannel := make(chan []byte)
	go func() {
		// Read from the connection untill a new line is send
		data, _ := bufio.NewReader(conn).ReadString('\n')
		responseChannel <- []byte(data)

	}()

	select {
	case response := <-responseChannel:
		return response, nil
	case <-time.After(timeOut):
		return nil, errors.New("Time Out Error")

	}

}

func (network *Network) SendPingMessage(contact *Contact) bool {
	ping := NewPingMessage(network.Ip)
	bytes, err := json.Marshal(ping)
	if err != nil {
		return false
	}

	response, err := network.Send(contact.Ip, contact.Port, bytes, time.Second*3)
	if err != nil {
		fmt.Println("Ping failed: " + err.Error())
		return false
	}
	var message Message
	errUnmarshal := json.Unmarshal(response, &message)
	if errUnmarshal != nil || message.MessageType != PONG {
		fmt.Println("Ping failed")
		return false
	}

	var pong Pong

	errUnmarshalAckPing := json.Unmarshal(response, &pong)
	if errUnmarshalAckPing != nil {
		fmt.Println("Ping failed: " + errUnmarshalAckPing.Error())
		return false
	}

	fmt.Println(pong.FromAddress + " acknowledged your ping")
	return true

}

func (network *Network) SendFindContactMessage(contact *Contact) ([]Contact, error) {
	findN := NewFindNodeMessage(network.Ip, contact.ID)
	bytes, err := json.Marshal(findN)
	if err != nil {
		return nil, err
	}

	response, err := network.Send(contact.Ip, contact.Port, bytes, time.Second*3)
	if err != nil {
		fmt.Println("Find node failed: " + err.Error())
		return nil, err
	}

	var arrayOfContacts []Contact
	json.Unmarshal(response, &arrayOfContacts)

	return arrayOfContacts, nil
}

func (network *Network) SendFindDataMessage(contact *Contact, key *Key) ([]Contact, string, error) {
	findData := NewFindDataMessage(network.Ip, contact.ID, key)
	bytes, err := json.Marshal(findData)
	if err != nil {
		return nil, "", err
	}

	response, err := network.Send(contact.Ip, contact.Port, bytes, time.Second*3)
	if err != nil {
		fmt.Println("Find data failed: " + err.Error())
		return nil, "", err
	}

	var data string
	json.Unmarshal(response, &data)
	fmt.Println("hello3:" + data)
	if data == "" {
		var arrayOfContacts []Contact
		json.Unmarshal(response, &arrayOfContacts)

		return arrayOfContacts, "", nil
	} else {
		return nil, data, nil
	}

}

func (network *Network) SendStoreMessage(data []byte) {
	// TODO
}
