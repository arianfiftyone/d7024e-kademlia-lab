package kademlia

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestJoinWithBootstrapOnly(t *testing.T) {
	kademliaBootsrap := NewKademlia("127.0.0.1", 3001, true, "", 0)

	go kademliaBootsrap.Start()

	time.Sleep(time.Second)

	kademlia := NewKademlia("127.0.0.1", 3002, false, "127.0.0.1", 3001)

	kademlia.Join()

	contact := kademlia.KademliaNode.RoutingTable.me
	assert.Equal(t, contact.ID, kademliaBootsrap.KademliaNode.RoutingTable.FindClosestContacts(contact.ID, 1)[0].ID, "The new node must be in the bootstraps routing table.")

	bootstrapContact := kademliaBootsrap.KademliaNode.RoutingTable.me
	assert.Equal(t, bootstrapContact.ID, kademlia.KademliaNode.RoutingTable.FindClosestContacts(bootstrapContact.ID, 1)[0].ID, "The bootsrapt must be in the new nodes routing table.")

}

func TestJoinWithMultipleNodes(t *testing.T) {
	kademliaBootsrap := NewKademlia("127.0.0.1", 3001, true, "", 0)

	var mockContacts = []Contact{
		NewContact(NewRandomKademliaID(), "198.162.1.1", 3000),
		NewContact(NewRandomKademliaID(), "198.162.1.2", 3000),
		NewContact(NewRandomKademliaID(), "198.162.1.3", 3000),
		NewContact(NewRandomKademliaID(), "198.162.1.4", 3000),
	}

	for _, contact := range mockContacts {
		kademliaBootsrap.KademliaNode.RoutingTable.AddContact(contact)
	}

	go kademliaBootsrap.Start()

	time.Sleep(time.Second)

	kademlia := NewKademlia("127.0.0.1", 3002, false, "127.0.0.1", 3001)

	kademlia.Join()

	// TODO: add assert after lookup is completed
}
