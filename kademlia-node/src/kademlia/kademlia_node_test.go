package kademlia

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestUpdateRoutingTableEmptyTable(t *testing.T) {
	kademliaNode := NewKademliaNode("127.0.0.1", 3002, false)

	kademliaNode.setNetwork(&NetworkImplementation{})

	contact := NewContact(NewKademliaID("FFFFFFFF00000000000000000000000000000000"), "198.168.1.1", 5000)

	fmt.Println("Updating an empty routing table with: " + contact.String())

	kademliaNode.updateRoutingTable(contact)

	bucket := kademliaNode.RoutingTable.buckets[kademliaNode.RoutingTable.getBucketIndex(contact.ID)]
	front := bucket.list.Front().Value.(Contact)
	fmt.Println("New front is: " + front.String())
	assert.Equal(t, contact, bucket.list.Front().Value.(Contact))
}

type NetworkMock struct{}

func (network *NetworkMock) Listen() error {
	return nil
}
func (network *NetworkMock) Send(ip string, port int, message []byte, timeOut time.Duration) ([]byte, error) {
	return nil, nil
}
func (network *NetworkMock) SendPingMessage(from *Contact, contact *Contact) error {
	return nil
}
func (network *NetworkMock) SendFindContactMessage(from *Contact, contact *Contact, id *KademliaID) ([]Contact, error) {
	return nil, nil
}
func (network *NetworkMock) SendFindDataMessage(from *Contact, contact *Contact, key *Key) ([]Contact, string, error) {
	return nil, "", nil
}
func (network *NetworkMock) SendStoreMessage(from *Contact, contact *Contact, key *Key, value string) bool {
	return false
}

func TestUpdateRoutingTableFullTable(t *testing.T) {
	kademliaNode := NewKademliaNode("127.0.0.1", 3002, false)
	kademliaNode.setNetwork(&NetworkMock{})

	contact := NewContact(NewKademliaID("FFFFFFFF00000000000000000000000000000000"), "198.168.1.1", 5000)

	bucket := kademliaNode.RoutingTable.buckets[kademliaNode.RoutingTable.getBucketIndex(contact.ID)]

	for bucket.Len() < bucketSize {
		randomContact := NewContact(NewRandomKademliaID(), "198.168.1.1", 5000)

		bucket.AddContact(randomContact)
	}
	fmt.Println("Updating a full routing table with: " + contact.String())

	kademliaNode.updateRoutingTable(contact)

	front := bucket.list.Front().Value.(Contact)
	fmt.Println("New front is: " + front.String())

	assert.Equal(t, contact, bucket.list.Front().Value.(Contact))
}
