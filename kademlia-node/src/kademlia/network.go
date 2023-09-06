package kademlia

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"net"
	"os"
	"time"
)

type Network struct {
	ip             string
	port           int
	messageHandler *MessageHandler
}

func Listen(ip string, port int) error {
	// listen to incoming udp packets
	conn, err := net.ListenUDP("udp", &net.UDPAddr{
		IP:   net.ParseIP(ip),
		Port: port,
	})
	if err != nil {
		return err
	}

	defer conn.Close()
	hostname, err := os.Hostname()
	if err != nil {
		return err
	}
	ips, err := net.LookupIP(hostname)
	if err != nil {
		return err
	}

	fmt.Printf("server listening %s\n", ips[0])

	network := &Network{
		ips[0].String(),
		port,
		&MessageHandler{},
	}

	for {
		data := make([]byte, 1024)
		len, remote, err := conn.ReadFromUDP(data[:])
		if err != nil {
			return err
		}

		go func(myConn *net.UDPConn) {
			response := network.messageHandler.HandleMessage(data[:len])
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
	ping := NewPingMessage(network.ip)
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
	if errUnmarshal != nil || message.MessageType != ACK_PING {
		fmt.Println("Ping failed")
		return false
	}

	var ackPing AckPing

	errUnmarshalAckPing := json.Unmarshal(response, &ackPing)
	if errUnmarshalAckPing != nil {
		fmt.Println("Ping failed: " + errUnmarshalAckPing.Error())
		return false
	}

	fmt.Println(ackPing.FromAddress + " acknowledged your ping")
	return true

}

func (network *Network) SendFindContactMessage(contact *Contact) {
	// TODO
}

func (network *Network) SendFindDataMessage(hash string) {
	// TODO
}

func (network *Network) SendStoreMessage(data []byte) {
	// TODO
}
