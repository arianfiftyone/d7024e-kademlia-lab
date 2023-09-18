package kademlia

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestJoinWithBootstrapOnly(t *testing.T) {
	kademliaBootsrap := NewKademlia("127.0.0.1", 2000, true, "", 0)

	go kademliaBootsrap.Start()

	time.Sleep(time.Second)

	kademlia := NewKademlia("127.0.0.1", 2001, false, "127.0.0.1", 2000)

	kademlia.Join()

	contact := kademlia.KademliaNode.RoutingTable.me
	assert.Equal(t, contact.ID, kademliaBootsrap.KademliaNode.RoutingTable.FindClosestContacts(contact.ID, 1)[0].ID, "The new node must be in the bootstraps routing table.")

	bootstrapContact := kademliaBootsrap.KademliaNode.RoutingTable.me
	assert.Equal(t, bootstrapContact.ID, kademlia.KademliaNode.RoutingTable.FindClosestContacts(bootstrapContact.ID, 1)[0].ID, "The bootsrapt must be in the new nodes routing table.")

}

func CreateMockedJoinKademlia(kademliaID *KademliaID, ip string, port int, bootstrapContact *Contact) Kademlia {
	routingTable := NewRoutingTable(NewContact(kademliaID, ip, port))
	dataStore := NewDataStore()
	kademliaNode := &KademliaNode{
		RoutingTable: routingTable,
		DataStore:    &dataStore,
	}

	network := &NetworkImplementation{
		ip,
		port,
		&MessageHandlerImplementation{
			kademliaNode,
		},
	}
	kademliaNode.setNetwork(network)
	kademlia := Kademlia{
		Network:          network,
		KademliaNode:     kademliaNode,
		isBootstrap:      false,
		bootstrapContact: bootstrapContact,
	}

	return kademlia
}

func TestJoinWithMultipleNodes(t *testing.T) {
	kademliaBootsrap := NewKademlia("127.0.0.1", 2002, true, "", 0)

	kademlia1 := CreateMockedKademlia(NewKademliaID("0000000000000000000000000000000000000001"), "127.0.0.1", 7001)
	kademlia2 := CreateMockedKademlia(NewKademliaID("0000000000000000000000000000000000000002"), "127.0.0.1", 7002)
	kademlia3 := CreateMockedKademlia(NewKademliaID("FFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFF"), "127.0.0.1", 7003)

	kademliaBootsrap.KademliaNode.RoutingTable.AddContact(kademlia1.KademliaNode.RoutingTable.me)
	kademliaBootsrap.KademliaNode.RoutingTable.AddContact(kademlia2.KademliaNode.RoutingTable.me)
	kademliaBootsrap.KademliaNode.RoutingTable.AddContact(kademlia3.KademliaNode.RoutingTable.me)

	go kademliaBootsrap.Start()

	time.Sleep(time.Second)

	kademlia := CreateMockedJoinKademlia(NewKademliaID("0000000000000000000000000000000000000000"), "127.0.0.1", 2001, &kademliaBootsrap.KademliaNode.RoutingTable.me)

	kademlia.Join()

	contacts := kademlia.KademliaNode.RoutingTable.FindClosestContacts(kademlia.KademliaNode.RoutingTable.me.ID, 20)
	fmt.Println(contacts)
	doesContainAll := kademlia.firstSetContainsAllContactsOfSecondSet(contacts, []Contact{kademlia1.KademliaNode.RoutingTable.me, kademlia2.KademliaNode.RoutingTable.me})
	containsKademlia3 := kademlia.firstSetContainsAllContactsOfSecondSet(contacts, []Contact{kademlia3.KademliaNode.RoutingTable.me})

	assert.True(t, doesContainAll && containsKademlia3)

}

type networkMock struct{}

func (network *networkMock) Listen() error {
	return nil
}
func (network *networkMock) Send(ip string, port int, message []byte, timeOut time.Duration) ([]byte, error) {
	return nil, nil
}
func (network *networkMock) SendPingMessage(from *Contact, contact *Contact) error {
	return nil
}
func (network *networkMock) SendFindContactMessage(from *Contact, contact *Contact, id *KademliaID) ([]Contact, error) {
	return nil, nil
}
func (network *networkMock) SendFindDataMessage(from *Contact, contact *Contact, key *Key) ([]Contact, string, error) {
	return nil, "", nil
}
func (network *networkMock) SendStoreMessage(from *Contact, contact *Contact, value string) bool {
	return false
}

func TestLookupContact(t *testing.T) {

	bootstrap := NewKademlia("127.0.0.1", 6000, true, "", 0)

	contact1 := NewContact(NewKademliaID("0000000000000000000000000000000000000001"), "198.168.1.1", 3000)
	contact2 := NewContact(NewKademliaID("0000000000000000000000000000000000000002"), "198.168.1.2", 3000)
	contact3 := NewContact(NewKademliaID("000000000000000000000000000000000000000F"), "198.168.1.3", 3000)
	contact4 := NewContact(NewKademliaID("FFFFFFFF00000000000000000000000000000001"), "198.168.1.4", 3000)
	contact5 := NewContact(NewKademliaID("FFFFFFFF00000000000000000000000000000002"), "198.168.1.5", 3000)
	contact6 := NewContact(NewKademliaID("FFFFFFFF00000000000000000000000000000003"), "198.168.1.6", 3000)
	contact7 := NewContact(NewKademliaID("FFFFFFFF00000000000000000000000000000004"), "198.168.1.7", 3000)
	contact8 := NewContact(NewKademliaID("FFFFFFFF00000000000000000000000000000005"), "198.168.1.8", 3000)
	contact9 := NewContact(NewKademliaID("FFFFFFFF00000000000000000000000000000006"), "198.168.1.9", 3000)
	bootstrap.KademliaNode.RoutingTable.AddContact(contact1)
	bootstrap.KademliaNode.RoutingTable.AddContact(contact3)
	bootstrap.KademliaNode.RoutingTable.AddContact(contact4)
	bootstrap.KademliaNode.RoutingTable.AddContact(contact5)
	bootstrap.KademliaNode.RoutingTable.AddContact(contact6)
	bootstrap.KademliaNode.RoutingTable.AddContact(contact7)
	bootstrap.KademliaNode.RoutingTable.AddContact(contact8)
	bootstrap.KademliaNode.RoutingTable.AddContact(contact9)
	bootstrap.KademliaNode.RoutingTable.AddContact(contact2)

	go bootstrap.Start()

	time.Sleep(time.Second)

	kademlia := NewKademlia("127.0.0.1", 3001, false, "", 0)

	kademlia.KademliaNode.RoutingTable.AddContact(bootstrap.KademliaNode.RoutingTable.me)

	list, _ := kademlia.LookupContact(contact1.ID)

	fmt.Println("closest To Target List")
	fmt.Println(list)

	doesContainAll := bootstrap.firstSetContainsAllContactsOfSecondSet(list, []Contact{contact1, contact2, contact3})
	assert.True(t, doesContainAll)

}

func CreateMockedKademlia(kademliaID *KademliaID, ip string, port int) Kademlia {
	routingTable := NewRoutingTable(NewContact(kademliaID, ip, port))
	dataStore := NewDataStore()
	kademliaNode := &KademliaNode{
		RoutingTable: routingTable,
		DataStore:    &dataStore,
	}

	network := &NetworkImplementation{
		ip,
		port,
		&MessageHandlerImplementation{
			kademliaNode,
		},
	}
	kademliaNode.setNetwork(network)
	kademlia := Kademlia{
		Network:      network,
		KademliaNode: kademliaNode,
		isBootstrap:  true,
	}

	return kademlia
}
func TestLookupContact2(t *testing.T) {

	bootstrap := CreateMockedKademlia(NewKademliaID("FFFFFFFF00000000000000000000000000000000"), "127.0.0.1", 7000)

	kademlia1 := CreateMockedKademlia(NewKademliaID("0000000000000000000000000000000000000001"), "127.0.0.1", 7001)
	kademlia2 := CreateMockedKademlia(NewKademliaID("0000000000000000000000000000000000000002"), "127.0.0.1", 7002)
	kademlia3 := CreateMockedKademlia(NewKademliaID("0000000000000000000000000000000000000003"), "127.0.0.1", 7003)
	kademlia4 := CreateMockedKademlia(NewKademliaID("FFFFFFFF00000000000000000000000000000001"), "127.0.0.1", 7004)
	kademlia5 := CreateMockedKademlia(NewKademliaID("FFFFFFFF00000000000000000000000000000002"), "127.0.0.1", 7005)
	kademlia6 := CreateMockedKademlia(NewKademliaID("FFFFFFFF00000000000000000000000000000003"), "127.0.0.1", 7006)

	bootstrap.KademliaNode.RoutingTable.AddContact(kademlia6.KademliaNode.RoutingTable.me)
	bootstrap.KademliaNode.RoutingTable.AddContact(kademlia5.KademliaNode.RoutingTable.me)

	kademlia6.KademliaNode.RoutingTable.AddContact(kademlia5.KademliaNode.RoutingTable.me)
	kademlia6.KademliaNode.RoutingTable.AddContact(kademlia4.KademliaNode.RoutingTable.me)
	kademlia6.KademliaNode.RoutingTable.AddContact(kademlia2.KademliaNode.RoutingTable.me)

	kademlia5.KademliaNode.RoutingTable.AddContact(kademlia6.KademliaNode.RoutingTable.me)
	kademlia5.KademliaNode.RoutingTable.AddContact(kademlia4.KademliaNode.RoutingTable.me)
	kademlia5.KademliaNode.RoutingTable.AddContact(kademlia3.KademliaNode.RoutingTable.me)

	kademlia4.KademliaNode.RoutingTable.AddContact(kademlia1.KademliaNode.RoutingTable.me)
	kademlia4.KademliaNode.RoutingTable.AddContact(kademlia5.KademliaNode.RoutingTable.me)
	kademlia4.KademliaNode.RoutingTable.AddContact(kademlia6.KademliaNode.RoutingTable.me)

	kademlia3.KademliaNode.RoutingTable.AddContact(kademlia5.KademliaNode.RoutingTable.me)
	kademlia3.KademliaNode.RoutingTable.AddContact(kademlia2.KademliaNode.RoutingTable.me)
	kademlia3.KademliaNode.RoutingTable.AddContact(kademlia1.KademliaNode.RoutingTable.me)

	kademlia2.KademliaNode.RoutingTable.AddContact(kademlia6.KademliaNode.RoutingTable.me)
	kademlia2.KademliaNode.RoutingTable.AddContact(kademlia3.KademliaNode.RoutingTable.me)
	kademlia2.KademliaNode.RoutingTable.AddContact(kademlia1.KademliaNode.RoutingTable.me)

	kademlia1.KademliaNode.RoutingTable.AddContact(kademlia4.KademliaNode.RoutingTable.me)
	kademlia1.KademliaNode.RoutingTable.AddContact(kademlia3.KademliaNode.RoutingTable.me)
	kademlia1.KademliaNode.RoutingTable.AddContact(kademlia2.KademliaNode.RoutingTable.me)

	go bootstrap.Start()
	go kademlia1.Start()
	go kademlia2.Start()
	go kademlia3.Start()
	go kademlia4.Start()
	go kademlia5.Start()
	go kademlia6.Start()
	time.Sleep(time.Second)

	kademlia := NewKademlia("127.0.0.1", 4000, false, "", 0)

	kademlia.KademliaNode.RoutingTable.AddContact(bootstrap.KademliaNode.RoutingTable.me)

	list, _ := kademlia.LookupContact(kademlia1.KademliaNode.RoutingTable.me.ID)

	fmt.Println("closest To Target List")
	fmt.Println(list)

	doesContainAll := bootstrap.firstSetContainsAllContactsOfSecondSet(list, []Contact{kademlia1.KademliaNode.RoutingTable.me, kademlia2.KademliaNode.RoutingTable.me, kademlia3.KademliaNode.RoutingTable.me})
	assert.True(t, doesContainAll)

}
