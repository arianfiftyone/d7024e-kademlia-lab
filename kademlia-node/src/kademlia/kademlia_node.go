package kademlia

const (
	NumberOfClosestNodesToRetrieved = 3
)

type KademliaNode struct {
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
