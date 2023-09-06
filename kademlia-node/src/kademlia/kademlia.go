package kademlia

import (
	"fmt"
	"net"
	"os"
	"strings"
)

type Kademlia struct {
	rt *RoutingTable
}

const (
	BootstrapKademliaID = "FFFFFFFF00000000000000000000000000000000"
)

func NewKademliaInstance() *Kademlia {

	var routingTable *RoutingTable
	var kademliaID KademliaID

	envValue := os.Getenv("IS_BOOTSTRAP")
	// converting string to bool
	IS_BOOTSTRAP := strings.ToLower(envValue) == "true"
	hostname := os.Getenv("CONTAINER_NODE")

	localIps, err := net.LookupIP(hostname)
	if err != nil {
		// TODO: Handle the error more gracefully(log or return)
		fmt.Println("Error", err)
	}
	localIp := localIps[0].String()

	if IS_BOOTSTRAP {

		kademliaID = *NewKademliaID(BootstrapKademliaID)
	} else {

		kademliaID = *NewRandomKademliaID()
	}

	// Create a new Contact instance based on the Kademlia ID and local IP
	contact := NewContact(&kademliaID, localIp)

	// Create a new RoutingTable instance and add the initial contact
	routingTable = NewRoutingTable(contact)

	return &Kademlia{
		rt: routingTable,
	}

}

func (kademlia *Kademlia) LookupContact(target *Contact) {
	// TODO
}

func (kademlia *Kademlia) LookupData(hash string) {
	// TODO
}

func (kademlia *Kademlia) Store(data []byte) {
	// TODO
}
