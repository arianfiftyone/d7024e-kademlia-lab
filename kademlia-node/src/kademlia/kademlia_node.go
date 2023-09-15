package kademlia

const (
	NumberOfClosestNodesToRetrieved = 3 // Must be atleast 3, otherwize some tests will fail
)

type KademliaNode struct {
	Network      Network
	RoutingTable *RoutingTable
	DataStore    *DataStore
}

func NewKademliaNode(ip string, port int, isBootstrap bool) *KademliaNode {
	var routingTable *RoutingTable
	var kademliaID KademliaID

	if isBootstrap {
		kademliaID = *NewKademliaID(BootstrapKademliaID)
	} else {
		kademliaID = *NewRandomKademliaID()
	}

	// Create a new Contact instance based on the Kademlia ID and local IP
	contact := NewContact(&kademliaID, ip, port)

	// Create a new RoutingTable instance and add the initial contact
	routingTable = NewRoutingTable(contact)

	// Create new DataStore instance
	dataStore := NewDataStore()

	return &KademliaNode{
		RoutingTable: routingTable,
		DataStore:    &dataStore,
	}
}

func (kademliaNode *KademliaNode) setNetwork(network Network) {
	kademliaNode.Network = network
}

func (kademliaNode *KademliaNode) updateRoutingTable(contact Contact) {
	if kademliaNode.Network == nil {
		return
	}

	bucket := kademliaNode.RoutingTable.buckets[kademliaNode.RoutingTable.getBucketIndex(contact.ID)]
	if bucket.Len() < bucketSize {
		kademliaNode.RoutingTable.AddContact(contact)

	} else {
		lastContact := bucket.list.Back().Value.(Contact)

		// Ping the last node in the bucket, replace if it does not respond otherwize do nothing
		err := kademliaNode.Network.SendPingMessage(&kademliaNode.RoutingTable.me, &lastContact)
		if err != nil {
			return

		}
		bucket.RemoveLastIfFull()

		bucket.AddContact(contact)

	}
}
