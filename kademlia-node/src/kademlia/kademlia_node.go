package kademlia

import "github.com/arianfiftyone/src/logger"

const (
	NumberOfClosestNodesToRetrieved = 3 // Must be atleast 3, otherwize some tests will fail
)

type KademliaNode interface {
	setNetwork(network Network)
	GetRoutingTable() *RoutingTable
	GetDataStore() *DataStore
	updateRoutingTable(contact Contact)
}

type KademliaNodeImplementation struct {
	Network      Network
	RoutingTable *RoutingTable
	DataStore    *DataStore
}

func NewKademliaNode(ip string, port int, isBootstrap bool) *KademliaNodeImplementation {
	var routingTable *RoutingTable
	var kademliaID KademliaID

	if isBootstrap {

		kademliaIDPointer, err := NewKademliaID(BootstrapKademliaID)

		if err != nil {
			logger.Log("cannot create KademliaID from given string")
		}
		kademliaID = *kademliaIDPointer
	} else {
		kademliaID = *NewRandomKademliaID()
	}

	// Create a new Contact instance based on the Kademlia ID and local IP
	contact := NewContact(&kademliaID, ip, port)

	// Create a new RoutingTable instance and add the initial contact
	routingTable = NewRoutingTable(contact)

	// Create new DataStore instance
	dataStore := NewDataStore()

	return &KademliaNodeImplementation{
		RoutingTable: routingTable,
		DataStore:    &dataStore,
	}
}

func (kademliaNode *KademliaNodeImplementation) setNetwork(network Network) {
	kademliaNode.Network = network
}

func (kademliaNode *KademliaNodeImplementation) GetRoutingTable() *RoutingTable {
	return kademliaNode.RoutingTable
}

func (kademliaNode *KademliaNodeImplementation) GetDataStore() *DataStore {
	return kademliaNode.DataStore
}

func (kademliaNode *KademliaNodeImplementation) updateRoutingTable(contact Contact) {
	if kademliaNode.Network == nil {
		return
	}

	bucket := kademliaNode.RoutingTable.buckets[kademliaNode.RoutingTable.getBucketIndex(contact.ID)]
	if bucket.Len() < bucketSize {
		kademliaNode.RoutingTable.AddContact(contact)

	} else {
		lastContact := bucket.list.Back().Value.(Contact)

		// Ping the last node in the bucket, replace if it does not respond otherwize do nothing
		err := kademliaNode.Network.SendPingMessage(&kademliaNode.RoutingTable.Me, &lastContact)
		if err != nil {
			return

		}
		bucket.RemoveLastIfFull()

		bucket.AddContact(contact)

	}
}
