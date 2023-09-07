package kademlia

import (
	"os"
	"strings"
)

const (
	NumberOfClosestNodesToRetrieved = 3
)

type KademliaNode struct {
	RoutingTable *RoutingTable
}

func NewKademliaNode(ip string, port int) *KademliaNode {
	var routingTable *RoutingTable
	var kademliaID KademliaID

	IS_BOOTSTRAP_STR := os.Getenv("IS_BOOTSTRAP")
	IS_BOOTSTRAP := strings.ToLower(IS_BOOTSTRAP_STR) == "true"

	if IS_BOOTSTRAP {
		kademliaID = *NewKademliaID(BootstrapKademliaID)
	} else {
		kademliaID = *NewRandomKademliaID()
	}

	// Create a new Contact instance based on the Kademlia ID and local IP
	contact := NewContact(&kademliaID, ip, port)

	// Create a new RoutingTable instance and add the initial contact
	routingTable = NewRoutingTable(contact)

	return &KademliaNode{
		RoutingTable: routingTable,
	}
}
