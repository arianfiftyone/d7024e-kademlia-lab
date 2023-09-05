package kademlia

import (
	"fmt"
	"net"
	"os"
	"strings"
)

type Network struct {
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
	localIp, err := net.LookupIP(hostname)
	if err != nil {
		return err
	}

	fmt.Printf("server listening %s\n", localIp[0])

	for {
		data := make([]byte, 20)
		rlen, remote, err := conn.ReadFromUDP(data[:])
		if err != nil {
			return err
		}

		message := strings.TrimSpace(string(data[:rlen]))
		fmt.Printf("received: %s from %s\n", message, remote)
	}

}

func (network *Network) SendPingMessage(contact *Contact) {
	// TODO
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
