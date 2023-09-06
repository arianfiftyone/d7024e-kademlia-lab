package kademlia

import (
	"fmt"
	"net"
	"os"
	"strconv"
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

	hostname := os.Getenv("CONTAINER_NODE")
	IS_BOOTSTRAP_STR := os.Getenv("IS_BOOTSTRAP")
	IS_BOOTSTRAP := strings.ToLower(IS_BOOTSTRAP_STR) == "true"
	NODE_PORT_STR := os.Getenv("NODE_PORT")
	NODE_PORT, err := strconv.Atoi(NODE_PORT_STR)
	if err != nil {
		fmt.Println("Error", err)
	}

	localIps, err := net.LookupIP(hostname)
	if err != nil {
		fmt.Println("Error", err)
	}
	localIp := localIps[0].String()

	if IS_BOOTSTRAP {
		kademliaID = *NewKademliaID(BootstrapKademliaID)
	} else {
		kademliaID = *NewRandomKademliaID()
	}

	// Create a new Contact instance based on the Kademlia ID and local IP
	contact := NewContact(&kademliaID, localIp, NODE_PORT)

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
